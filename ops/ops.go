package ops

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"optimizer/globals"
)

func SaveResultsJSONFile(results *[]globals.CutMaterial, job string) {
	// fmt.Println(results)
	resultsJSON, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
	}
	dir, err := os.Getwd()
	dir = filepath.Join(dir, "results")
	// Save JSON to a file
	filePath := filepath.Join(dir, job+".json")

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
	}
	defer file.Close()

	_, err = file.Write(resultsJSON)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}

	fmt.Println("OPTI RESULTS JSON saved to results.json")
}
