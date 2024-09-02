package optimizer

import (
	"fmt"

	"main/globals"
	"main/internal/db"
	"main/logger"
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
				logger.LogError("MaterialErrors", err.Error())
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
				Job:              "",
				MaterialCode:     "",
				Parts:            map[string]globals.PartQTY{},
				Quantity:         0,
				StockLength:      0,
				Length:           0,
				CuttingOperation: "",
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
				logger.LogError("MaterialErrors", err.Error())
				// Optionally continue or break depending on error handling
			}
		}

		if sortedParts[i].CutQuantity < sortedParts[i].Quantity {
			if sortedParts[i].Length <= currentMaterial.Length && sortedParts[i].CuttingOperation == currentMaterial.CuttingOperation {
				if sortedParts[i].Length+globals.Settings.Kerf <= currentMaterial.Length {
					fmt.Print("Subtracting Kerf", sortedParts[i].Length+globals.Settings.Kerf)
					fmt.Println("Kerf is: ", globals.Settings.Kerf)
					currentMaterial.Length -= (sortedParts[i].Length + globals.Settings.Kerf)
				} else {
					fmt.Print("Not Subtracting Kerf", sortedParts[i].Length+globals.Settings.Kerf)
					fmt.Println("NOT Kerf is: ", globals.Settings.Kerf)
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

// func CreateLayoutV2(Parts []globals.Part, Materials []globals.Material, JobInfo globals.JobType) (errSlice []string) {
// 	var results []globals.CutMaterial
// 	errSlice = []string{}
//
// 	if len(Parts) == 0 {
// 		return []string{"Error: No parts to process."}
// 	}
//
// 	sortedParts := part_utils.SortPartByLength(Parts)
// 	var currentMaterial globals.CutMaterial
// 	complete := false
// 	i := 0
//
// 	for !complete {
// 		if shouldAppendCurrentMaterial(sortedParts, &currentMaterial) {
// 			appendCurrentMaterial(&results, &currentMaterial)
// 		}
//
// 		complete = PartsComplete(sortedParts)
// 		if complete {
// 			finalizeCurrentMaterial(&results, &currentMaterial)
// 			break
// 		}
//
// 		if currentMaterial.MaterialCode == "" {
// 			err := checkForMaterialV2(&sortedParts[i], &currentMaterial, &Materials, JobInfo)
// 			if err != nil {
// 				logger.LogError("MaterialErrors", err.Error())
// 				// Optionally handle error (e.g., continue, break, etc.)
// 			}
// 		}
//
// 		processPart(sortedParts, i, &currentMaterial)
//
// 		// Move to the next part
// 		i = (i + 1) % len(sortedParts)
// 	}
//
// 	ops.SaveResultsJSONFile(&results, JobInfo.Job)
// 	mergeDuplicateCutMaterialsInPlace(&results)
// 	db.SavePartsToDB(&results, JobInfo)
// 	return errSlice
// }
