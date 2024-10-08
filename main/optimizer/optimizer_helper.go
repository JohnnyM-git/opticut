package optimizer

import (
	"errors"
	"fmt"
	"sort"

	"main/globals"
	"main/material_utils"
)

func PartsComplete(parts []globals.Part) bool {
	for _, part := range parts {
		if part.Quantity > part.CutQuantity {
			return false
		}
	}
	return true
}

func NoPartsFit(parts []globals.Part, currentMaterial *globals.CutMaterial) bool {
	if currentMaterial.MaterialCode == "" {
		return true
	}
	for _, part := range parts {
		if part.Quantity > part.CutQuantity {
			if part.Length <= currentMaterial.Length && part.CuttingOperation == currentMaterial.CuttingOperation {
				// fmt.Println("Parts are fitting")
				return false
			}
		}
	}

	// fmt.Println("No parts fit on the material")
	return true
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
				return len(*results) - 1, nil
			}
		}
	}
	errMsg := fmt.Sprint(p.MaterialCode, ": no material at the correct length of ", p.Length, " found needed")
	return 0, errors.New(errMsg)
}

func checkForMaterialV2(
	p *globals.Part, currentMaterial *globals.CutMaterial, Materials *[]globals.Material,
	JobInfo globals.JobType) error {

	if len(*Materials) == 0 {
		return errors.New("Materials list empty")
	}

	material_utils.SortMaterialByLength(*Materials)
	for i, material := range *Materials {
		if material.Length >= p.Length {
			m := globals.CutMaterial{
				Job:              JobInfo.Job,
				MaterialCode:     material.MaterialCode,
				Parts:            map[string]globals.PartQTY{},
				Quantity:         1,
				StockLength:      material.Length,
				Length:           material.Length,
				CuttingOperation: p.CuttingOperation,
			}
			*currentMaterial = m
			if (*Materials)[i].Quantity != 9999 {
				(*Materials)[i].Quantity-- // Decrement material quantity
			}
			return nil
		}
	}

	return errors.New("material with the correct length not found")
}

// func shouldAppendCurrentMaterial(sortedParts []globals.Part, currentMaterial *globals.CutMaterial) bool {
// 	return NoPartsFit(sortedParts, currentMaterial) && currentMaterial.MaterialCode != ""
// }
//
// func appendCurrentMaterial(results *[]globals.CutMaterial, currentMaterial *globals.CutMaterial) {
// 	*results = append(*results, *currentMaterial)
// }
//
// func processPart(sortedParts []globals.Part, i int, currentMaterial *globals.CutMaterial) {
// 	if sortedParts[i].CutQuantity < sortedParts[i].Quantity {
// 		if canCutPart(&sortedParts[i], currentMaterial) {
// 			if shouldSubtractKerf(&sortedParts[i], currentMaterial) {
// 				subtractKerf(&sortedParts[i], currentMaterial)
// 			} else {
// 				subtractLength(&sortedParts[i], currentMaterial)
// 			}
// 			updatePartQuantity(&sortedParts[i], currentMaterial)
// 		}
// 	}
// }
//
// func canCutPart(sortedPart *Part, currentMaterial *Material) bool {
// 	return sortedPart.Length <= currentMaterial.Length && sortedPart.CuttingOperation == currentMaterial.CuttingOperation
// }
//
// func shouldSubtractKerf(sortedPart *Part, currentMaterial *Material) bool {
// 	return sortedPart.Length+globals.Settings.Kerf <= currentMaterial.Length
// }
//
// func subtractKerf(sortedPart *Part, currentMaterial *Material) {
// 	fmt.Print("Subtracting Kerf", sortedPart.Length+globals.Settings.Kerf)
// 	fmt.Println("Kerf is: ", globals.Settings.Kerf)
// 	currentMaterial.Length -= (sortedPart.Length + globals.Settings.Kerf)
// }
//
// func subtractLength(sortedPart *Part, currentMaterial *Material) {
// 	fmt.Print("Not Subtracting Kerf", sortedPart.Length+globals.Settings.Kerf)
// 	fmt.Println("NOT Kerf is: ", globals.Settings.Kerf)
// 	currentMaterial.Length -= sortedPart.Length
// }
//
// func updatePartQuantity(sortedPart *Part, currentMaterial *Material) {
// 	if partQty, exists := currentMaterial.Parts[sortedPart.PartNumber]; exists {
// 		partQty.CurrentQty++
// 		currentMaterial.Parts[sortedPart.PartNumber] = partQty
// 	} else {
// 		currentMaterial.Parts[sortedPart.PartNumber] = globals.PartQTY{
// 			CurrentQty: 1,
// 			TotalQty:   sortedPart.Quantity,
// 		}
// 	}
// 	sortedPart.CutQuantity++
// }
