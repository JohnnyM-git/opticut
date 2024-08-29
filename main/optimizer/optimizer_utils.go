package optimizer

import (
	"errors"
	"fmt"
	"sort"

	"main/globals"
	"main/internal/db"
	"main/logger"
	"main/material_utils"
	"main/ops"
	"main/part_utils"
)

func CreateLayout(
	parts []globals.Part,
	materials []globals.Material,
	JobInfo globals.JobType) (
	// results []globals.CutMaterial,
	errSlice []string) {
	var results []globals.CutMaterial
	errSlice = []string{}
	for _, part := range parts {
		p := part
		remainingQty := p.Quantity
		for {
			if remainingQty == 0 {
				break
			}
			materialIndex, err := checkForMaterial(
				&p,
				&results,
				&materials,
				JobInfo.Job)
			if err != nil {
				logger.LogError(err.Error())
				errSlice = append(errSlice, err.Error())
				break
			} else {
				cutMaterial := &results[materialIndex]
				if partQty, exists := cutMaterial.Parts[p.PartNumber]; exists {
					// Update the existing PartQTY struct
					partQty.CurrentQty += 1
					cutMaterial.Parts[p.PartNumber] = partQty // Store it back in the map
				} else {
					// Initialize and add a new PartQTY struct to the map
					cutMaterial.Parts[p.PartNumber] = globals.PartQTY{
						CurrentQty: 1,
						TotalQty:   p.Quantity,
					}
				}
				fmt.Println("SUB", globals.Settings.Kerf)
				cutMaterial.Length -= (p.Length + globals.Settings.Kerf)
				remainingQty--
			}
		}
	}
	mergeDuplicateCutMaterialsInPlace(&results)
	db.SavePartsToDB(&results, JobInfo)
	// ops.SaveResultsJSONFile(&results, results[0].Job)
	return errSlice
}

func CreateLayoutV2(Parts []globals.Part, Materials []globals.Material, JobInfo globals.JobType) (errSlice []string) {
	fmt.Println("Hitting CreateLayoutV2")

	var results []globals.CutMaterial
	errSlice = []string{}

	if len(Parts) == 0 {
		return []string{"Error: No parts to process."}
	}

	sortedParts := part_utils.SortPartByLength(Parts)
	var currentMaterial globals.CutMaterial

	complete := false
	i := 0

	for !complete {
		fmt.Println("INDEX", i)
		fmt.Println("NoPartsFit:", NoPartsFit(sortedParts, &currentMaterial), "Code:", currentMaterial.MaterialCode)

		if NoPartsFit(sortedParts, &currentMaterial) && currentMaterial.MaterialCode != "" {
			fmt.Println("Appending current because parts don't fit or all are complete", currentMaterial)
			results = append(results, currentMaterial)
			// Reset current material
			currentMaterial = globals.CutMaterial{
				Job:          "",
				MaterialCode: "",
				Parts:        map[string]globals.PartQTY{},
				Quantity:     0,
				StockLength:  0,
				Length:       0,
			}
		}

		// Update complete status
		complete = PartsComplete(sortedParts)
		fmt.Println("Complete:", complete, "Code:", currentMaterial.MaterialCode)
		if complete {
			fmt.Println("Appending current because parts are done", currentMaterial)
			if currentMaterial.MaterialCode != "" {
				results = append(results, currentMaterial)

			}
			break
		}

		fmt.Println("Current Material", currentMaterial, "Code", currentMaterial.MaterialCode)
		if currentMaterial.MaterialCode == "" {
			fmt.Println("No material available")
			err := checkForMaterialV2(&sortedParts[i], &currentMaterial, &Materials, JobInfo)
			if err != nil {
				logger.LogError(err.Error())
				// Optionally continue or break depending on error handling
			}
		}

		if sortedParts[i].CutQuantity < sortedParts[i].Quantity {
			if sortedParts[i].Length <= currentMaterial.Length {
				if sortedParts[i].Length+globals.Settings.Kerf <= currentMaterial.Length {
					currentMaterial.Length -= (sortedParts[i].Length + globals.Settings.Kerf)
				} else {
					currentMaterial.Length -= sortedParts[i].Length
				}
				sortedParts[i].CutQuantity++

				// Update or add PartQTY
				if partQty, exists := currentMaterial.Parts[sortedParts[i].PartNumber]; exists {
					partQty.CurrentQty++
					currentMaterial.Parts[sortedParts[i].PartNumber] = partQty
				} else {
					currentMaterial.Parts[sortedParts[i].PartNumber] = globals.PartQTY{
						CurrentQty: 1,
						TotalQty:   sortedParts[i].Quantity,
					}
				}
			}
		}

		// Move to the next part
		i = (i + 1) % len(sortedParts)
	}
	ops.SaveResultsJSONFile(&results, JobInfo.Job)
	mergeDuplicateCutMaterialsInPlace(&results)
	db.SavePartsToDB(&results, JobInfo)
	return errSlice
}

// Assuming PartsComplete is updated to return a bool:
func PartsComplete(parts []globals.Part) bool {
	for _, part := range parts {
		if part.Quantity > part.CutQuantity {
			return false
		}
	}
	return true
}

func NoPartsFit(parts []globals.Part, currentMaterial *globals.CutMaterial) bool {
	fmt.Println("Checking if parts fit")
	// if currentMaterial.MaterialCode == "" {
	// 	fmt.Println("No material code")
	// 	return true
	// }
	for _, part := range parts {
		if part.Quantity > part.CutQuantity {
			if part.Length <= currentMaterial.Length {
				fmt.Println("Parts are fitting")
				return false
			}
		}
	}
	fmt.Println("No parts fit on the material")
	return true
}

// func PartsComplete(parts []globals.Part) bool {
// 	fmt.Println("Checking if parts complete")
// 	for _, part := range parts {
// 		if part.Quantity > part.CutQuantity {
// 			fmt.Println("Part#", part.PartNumber, "PartQ", part.Quantity, "CutQ", part.CutQuantity)
// 			return false // Set to false if any part is not complete
// 		}
// 	}
// 	return true
// }

func checkForMaterialV2(
	p *globals.Part, currentMaterial *globals.CutMaterial, Materials *[]globals.Material,
	JobInfo globals.JobType) error {
	fmt.Println("Hitting CheckforMaterialV2")

	if len(*Materials) == 0 {
		return errors.New("Materials list empty")
	}

	material_utils.SortMaterialByLength(*Materials)
	for i, material := range *Materials {
		if material.Length >= p.Length {
			m := globals.CutMaterial{
				Job:          JobInfo.Job,
				MaterialCode: material.MaterialCode,
				Parts:        map[string]globals.PartQTY{},
				Quantity:     1,
				StockLength:  material.Length,
				Length:       material.Length,
			}
			*currentMaterial = m
			if (*Materials)[i].Quantity != 9999 {
				(*Materials)[i].Quantity-- // Decrement material quantity
			}
			return nil
		}
	}

	return errors.New("Material with the correct length not found.")
}

func checkForMaterial(
	p *globals.Part,
	results *[]globals.CutMaterial,
	materials *[]globals.Material,
	job string) (
	int, error) {
	if len(*results) != 0 {
		SortCutMaterialsByLength(*results)
		for i, m := range *results {
			if p.Length <= m.Length {
				return i, nil
			}
		}
	}

	if len(*materials) != 0 {
		SortMaterialsByLength(*materials)
		// fmt.Println(*materials)
		for i := range *materials {
			material := &(*materials)[i]
			if p.Length <= material.Length && material.Quantity != 0 {
				m := globals.CutMaterial{
					Job:          job,
					MaterialCode: material.MaterialCode,
					Parts:        map[string]globals.PartQTY{},
					Quantity:     1,
					StockLength:  material.Length,
					Length:       material.Length,
				}
				*results = append(*results, m)
				(*materials)[i].Quantity = (*materials)[i].Quantity - 1
				fmt.Println("M QTY:", material.Quantity)
				return len(*results) - 1, nil
			}
		}
	}
	errMsg := fmt.Sprint(
		p.MaterialCode,
		": no material at the correct length found ",
		p.Length,
		`"`,
		" needed")
	return 0, errors.New(errMsg)
}

func SortCutMaterialsByLength(materials []globals.CutMaterial) {
	sort.Slice(
		materials, func(i, j int) bool {
			return materials[i].Length < materials[j].Length
		})
}

func SortMaterialsByLength(materials []globals.Material) {
	sort.Slice(
		materials, func(i, j int) bool {
			return materials[i].Length < materials[j].Length
		})
}

func mergeDuplicateCutMaterialsInPlace(cutMaterials *[]globals.CutMaterial) {
	var mergedResults []globals.CutMaterial
	fmt.Println("Merging duplicates", *cutMaterials)
	for _, cm := range *cutMaterials {
		found := false
		for i, merged := range mergedResults {
			if cm.MaterialCode == merged.MaterialCode &&
				cm.StockLength == merged.StockLength &&
				cm.Length == merged.Length &&
				arePartsEqual(cm.Parts, merged.Parts) {

				mergedResults[i].Quantity += cm.Quantity
				found = true
				break
			}
		}
		if !found {
			mergedResults = append(mergedResults, cm)
			// fmt.Println("Merged duplicate:", cm)
		}
	}
	for _, merged := range mergedResults {
		fmt.Println("Merged duplicate:", merged)
	}

	*cutMaterials = mergedResults
}

func arePartsEqual(parts1, parts2 map[string]globals.PartQTY) bool {
	if len(parts1) != len(parts2) {
		return false
	}
	for part, qty := range parts1 {
		if parts2[part] != qty {
			return false
		}
	}
	return true
}
