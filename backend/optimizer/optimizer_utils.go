package optimizer

import (
	"errors"
	"fmt"
	"sort"

	"optimizer/globals"
	"optimizer/internal/db"
	"optimizer/logger"
)

func CreateLayout(
	parts []globals.Part,
	materials []globals.Material,
	JobInfo globals.JobType) (
	results []globals.CutMaterial,
	errSlice []string) {
	results = []globals.CutMaterial{}
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
				if _, exists := cutMaterial.Parts[p.PartNumber]; exists {
					cutMaterial.Parts[p.PartNumber] += 1
				} else {
					cutMaterial.Parts[p.PartNumber] = 1
				}
				cutMaterial.Length -= p.Length + globals.Settings.Kerf
				remainingQty--
			}
		}
	}
	mergeDuplicateCutMaterialsInPlace(&results)
	db.SavePartsToDB(&results, JobInfo)
	// ops.SaveResultsJSONFile(&results, results[0].Job)
	return results, errSlice
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
					Parts:        map[string]uint16{},
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
	mergedResults := []globals.CutMaterial{}

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
		}
	}

	*cutMaterials = mergedResults
}

func arePartsEqual(parts1, parts2 map[string]uint16) bool {
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
