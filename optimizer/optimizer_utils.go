package optimizer

import (
	"fmt"

	"optimizer/globals"
)

func CreateLayout(parts []globals.Part, materials []globals.Material) []globals.CutMaterials {
	var results []globals.CutMaterials
	for i, material := range materials {
		fmt.Println(i, material)
		results = append(
			results, globals.CutMaterials{
				MaterialCode:    material.MaterialCode,
				Parts:           nil,
				RemainingLength: material.StockLength,
				Quantity:        material.Quantity,
			})
	}
	// for _, part := range parts {
	// 	p := part
	// 	var remainingQty = p.Quantity
	//
	//
	// 	// for i := 0; i < int(p.Quantity); i++ {
	// 	// 	if remainingQty == 0 {
	// 	// 		break
	// 	// 	}
	// 	// 	if len(results) != 0 {
	// 	// 		for _, material := range results {
	// 	// 			if p.Length <= material.RemainingLength {
	// 	// 				m := material
	// 	// 				m.RemainingLength -= p.Length + globals.Kerf
	// 	// 				m.Parts = append(m.Parts, p)
	// 	// 				fmt.Println("TEST", m)
	// 	// 				results = append(results, m)
	// 	// 				remainingQty -= 1
	// 	// 			} else {
	// 	// 				fmt.Println("TEST", p)
	// 	// 			}
	// 	// 		}
	// 	// 	}
	// 	// }
	//
	// }
	fmt.Println("RESULTS", results)
	return results
}
