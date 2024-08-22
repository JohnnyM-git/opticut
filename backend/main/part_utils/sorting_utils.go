package part_utils

import (
	"sort"

	"main/globals"
)

type ByLength []globals.Part

func (a ByLength) Len() int           { return len(a) }
func (a ByLength) Less(i, j int) bool { return a[i].Length > a[j].Length } // Descending order
func (a ByLength) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func SortPartByLength(parts []globals.Part) []globals.Part {
	sort.Sort(ByLength(parts))
	return parts
}

func SortPartsByCode(parts []globals.Part) [][]globals.Part {
	partMap := make(map[string][]globals.Part)
	for _, part := range parts {
		partMap[part.MaterialCode] = append(partMap[part.MaterialCode], part)
	}

	var result [][]globals.Part
	for _, partsGroup := range partMap {
		sortedGroup := SortPartByLength(partsGroup)
		result = append(result, sortedGroup)
	}
	// fmt.Println(result)
	return result
}
