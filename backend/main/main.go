package main

import (
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"optimizer/globals"
	"optimizer/internal/db"
	"optimizer/internal/server"
	"optimizer/logger"
	"optimizer/material_utils"
	"optimizer/optimizer"
	"optimizer/part_utils"
)

func main() {

	rootPath := filepath.Join(".", "prod.db")

	// Check if prod.db exists at the root of the project
	if !fileExists(rootPath) {
		db.SetupDB()
	} else {
		db.InitDB("./prod.db")
	}

	db.InsertPartsIntoPartTable(globals.Parts)
	sortedGroupedPartSlice := part_utils.SortPartsByCode(globals.Parts)
	// fmt.Println(sortedGroupedPartSlice)
	for _, partSlice := range sortedGroupedPartSlice {
		// fmt.Println(partSlice)
		materialCode := partSlice[0].MaterialCode
		results, err := material_utils.SortMaterialByCode(
			globals.Materials,
			materialCode)
		if err != nil {
			logger.LogError(err.Error())
		} else {
			results, errSlice := optimizer.CreateLayout(partSlice, results)
			if len(errSlice) > 0 {
				for _, err := range errSlice {
					logger.LogError(err)
				}
			} else {
				fmt.Println(results)
			}
		}

	}
	server.StartServer()
}

func fileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
