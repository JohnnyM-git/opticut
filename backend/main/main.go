package main

import (
	"os"
	"path/filepath"

	"backend/globals"
	"backend/internal/db"
	"backend/internal/server"
)

func main() {

	rootPath := filepath.Join(".", "prod.db")

	// Check if prod.db exists at the root of the project
	if !fileExists(rootPath) {
		db.SetupDB()
	} else {
		db.InitDB("./prod.db")
	}

	globals.LoadSettings()

	// db.InsertPartsIntoPartTable(globals.Parts)
	// err := db.SaveJobInfoToDB(globals.JobInfo)
	// if err != nil {
	// 	errMsg := fmt.Sprintf("Failed to save job info to database: %v", err)
	// 	logger.LogError(errMsg)
	// }
	// sortedGroupedPartSlice := part_utils.SortPartsByCode(globals.Parts)
	// fmt.Println(sortedGroupedPartSlice)
	// for _, partsByCodeSlice := range sortedGroupedPartSlice {
	// 	// fmt.Println(partSlice)
	// 	materialCode := partsByCodeSlice[0].MaterialCode
	// 	results, err := material_utils.SortMaterialByCode(
	// 		globals.Materials,
	// 		materialCode)
	// 	if err != nil {
	// 		logger.LogError(err.Error())
	// 	} else {
	// 		results, errSlice := optimizer.CreateLayout(
	// 			partsByCodeSlice,
	// 			results,
	// 			globals.JobInfo)
	// 		if len(errSlice) > 0 {
	// 			for _, err := range errSlice {
	// 				logger.LogError(err)
	// 			}
	// 		} else {
	// 			fmt.Println(results)
	// 		}
	// 	}
	//
	// }
	server.StartServer()
}

func fileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
