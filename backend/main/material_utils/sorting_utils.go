package material_utils

import (
	"errors"
	"fmt"
	"sort"

	"main/globals"
)

type ByLength []globals.Material

func (a ByLength) Len() int           { return len(a) }
func (a ByLength) Less(i, j int) bool { return a[i].Length < a[j].Length } // Descending order
func (a ByLength) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func SortMaterialByLength(materials []globals.Material) []globals.Material {
	sort.Sort(ByLength(materials))
	return materials
}

func SortMaterialByCode(materials []globals.Material, target string) ([]globals.Material, error) {
	var result []globals.Material
	for _, material := range materials {
		// fmt.Println("mat:", material, i, "target:", target)
		if material.MaterialCode == target {
			result = append(result, material)
		}
	}
	if len(result) == 0 {
		errMsg := fmt.Sprintf("%s: no material found", target)
		return result, errors.New(errMsg)
	}
	// else {
	// 	result = SortMaterialByLength(result)
	// }
	// fmt.Println(target, ":", result)
	return result, nil
}
