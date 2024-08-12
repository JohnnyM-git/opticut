package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"optimizer/globals"
	"optimizer/logger"
	"optimizer/material_utils"
	"optimizer/optimizer"
	"optimizer/part_utils"
)

func main() {

	rootPath := filepath.Join(".", "prod.db")

	// Check if prod.db exists at the root of the project
	if !fileExists(rootPath) {
		initdb()
	}
	part_utils.InsertPartsIntoPartTable(globals.Parts)
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

func fileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func initdb() {
	dbFile := "prod.db"

	// Open the database (creates the file if it doesn't exist)
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Read the SQL script from file
	script, err := ioutil.ReadFile("setup_db.sql")
	if err != nil {
		log.Fatal(err)
	}

	// Execute the SQL script
	_, err = db.Exec(string(script))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database and tables created successfully!")
}
