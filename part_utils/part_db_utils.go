package part_utils

import (
	"database/sql"
	"fmt"
	"math"

	"optimizer/globals"
	"optimizer/logger"
)

func InsertPartsIntoPartTable(parts []globals.Part) {
	db, err := sql.Open("sqlite3", "./prod.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
	}
	defer db.Close()

	for _, part := range parts {
		formattedLength := formatFloat32(part.Length, 3)
		fmt.Println(formattedLength)
		_, err := db.Exec(
			`INSERT OR IGNORE INTO parts (part_number, material_code, length) VALUES(?, ?, ?)`,
			part.PartNumber,
			part.MaterialCode,
			formattedLength)
		if err != nil {
			logger.LogError(err.Error())
		}
	}
}

func formatFloat32(value float32, precision int) float32 {
	scale := float32(math.Pow(10, float64(precision)))
	return float32(math.Round(float64(value)*float64(scale)) / float64(scale))
}

func SavePartsToDB(results *[]globals.CutMaterial) {
	db, err := sql.Open("sqlite3", "./prod.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
	}
	defer db.Close()

	for _, result := range *results {
		id, err := db.Exec(
			`INSERT INTO cut_materials (job, material_code, quantity, stock_length, length) VALUES(?, ?, ?, ?, ?)`,
			result.Job,
			result.MaterialCode,
			result.Quantity,
			result.StockLength,
			result.Length)
		if err != nil {
			logger.LogError(err.Error())
		} else {
			lastID, err := id.LastInsertId()
			if err != nil {
				logger.LogError(err.Error())
			}
			for i, part := range result.Parts {
				var partID int64
				fmt.Println(part)
				err := db.QueryRow(`SELECT id FROM parts WHERE part_number = ?`, i).Scan(&partID)
				if err != nil {
					logger.LogError(err.Error())
					continue
				}
				_, err = db.Exec(
					`INSERT INTO cut_material_parts (cut_material_id, part_id, part_qty) VALUES(?, ?, ?)`,
					lastID,
					partID,
					part,
				)
			}
			fmt.Println("Inserted ", lastID, " into database")
		}
		// for i, part := range result.Parts {
		// 	fmt.Println(index, result.Job, i, part)
		// }
	}
}
