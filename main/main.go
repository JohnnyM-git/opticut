package main

import (
	"optimizer/globals"
	"optimizer/logger"
	"optimizer/material_utils"
	"optimizer/optimizer"
	"optimizer/part_utils"
)

func main() {
	sortedGroupedPartSlice := part_utils.SortPartsByCode(globals.Parts)
	// fmt.Println(sortedGroupedPartSlice)
	for _, partSlice := range sortedGroupedPartSlice {
		// fmt.Println(partSlice)
		materialCode := partSlice[0].MaterialCode
		results, err := material_utils.SortMaterialByCode(globals.Materials, materialCode)
		if err != nil {
			logger.LogError(err.Error())
		} else {
			optimizer.CreateLayout(partSlice, results)
			// fmt.Println(results)
		}

	}
}
