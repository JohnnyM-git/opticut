package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
	"main/globals"
	"main/internal/db"
	"main/logger"
	"main/material_utils"
	"main/optimizer"
	"main/part_utils"
)

func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Settings Handler")
	fmt.Println(r.Method)
	switch r.Method {
	case "GET":
		GetSettingsHandler(w, r)
	case "POST":
		UpdateSettingsHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func GetSettingsHandler(w http.ResponseWriter, r *http.Request) {
	settings, err := globals.LoadSettings()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(settings)
}

func UpdateSettingsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HTTP Method:", r.Method)
	fmt.Println("Endpoint Hit: Update Settings")

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Print the request body as a string
	fmt.Println("Request Body:", string(body))

	// Parse the JSON body into a SettingsConfig struct
	var newSettings globals.SettingsConfig
	err = json.Unmarshal(body, &newSettings)
	if err != nil {
		fmt.Println("Unmarshal error:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Println("newSettings:", newSettings)

	// Load current settings
	settings, err := globals.LoadSettings()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Current Settings:", settings)

	// Update settings with new values
	settings.Kerf = newSettings.Kerf
	settings.Units = newSettings.Units
	settings.Excel = newSettings.Excel

	// Save the updated settings to the JSON file
	err = globals.SaveSettings(settings, "./settings.json")
	if err != nil {
		fmt.Println("Save error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Saved settings:", settings)
	globals.LoadSettings()
	// Create a response object
	response := map[string]interface{}{
		"message":  "Settings updated successfully",
		"settings": newSettings,
	}

	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode the response as JSON and write it to the response body
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ToggleStar(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HTTP Method:", r.Method)
	fmt.Println("Endpoint Hit: Toggle Star")

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Define a struct with exported fields and correct JSON tags
	type ToggleStarParams struct {
		JobNumber string `json:"jobNumber"`
		Value     int    `json:"value"`
	}

	// Unmarshal the JSON body into the struct
	var params ToggleStarParams
	err = json.Unmarshal(body, &params)
	if err != nil {
		fmt.Println("Unmarshal error:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Println("Request Body:", params)

	// Call the database function
	toggleErr := db.ToggleStar(params.JobNumber, params.Value)
	if toggleErr != nil {
		fmt.Println("Toggle Error:", toggleErr)
		http.Error(w, toggleErr.Error(), http.StatusInternalServerError)
		return
	}

	// Create a response object
	response := map[string]interface{}{
		"message": "Toggle Star successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CheckHealth(w http.ResponseWriter, r *http.Request) {
	type HealthResponse struct {
		Status      string `json:"status"`
		Database    string `json:"database"`
		Version     string `json:"version"`
		Uptime      string `json:"uptime"`
		ServiceName string `json:"service_name"`
	}

	uptime := time.Since(StartTime).String()

	health := HealthResponse{
		Status:      "Healthy",
		Database:    db.DbHealthCheck(),
		Version:     "0.1.0",
		Uptime:      uptime,
		ServiceName: "Cutwise",
	}

	// If any component is unhealthy, change the overall status
	if health.Database != "Healthy" {
		health.Status = "Unhealthy"
	}

	w.Header().Set("Content-Type", "application/json")
	if health.Status == "Unhealthy" {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(health)
}

func FileUpload(w http.ResponseWriter, r *http.Request) {

	filePath := r.URL.Query().Get("filePath")
	// Define the file to process
	// cwd, _ := os.Getwd()
	// var filesDir = globals.Settings.Excel.FilesPath
	// fileName := "Book1.xlsx" // Replace with your file name
	// filePath := filepath.Join(filesDir, fileName)

	fmt.Println("File Path:", filePath)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Process the Excel file
	data, err := processExcel(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the processed data as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func BatchProcessFiles(w http.ResponseWriter, r *http.Request) {

	// Define the file to process
	// cwd, _ := os.Getwd()
	var filesDir = globals.Settings.Excel.FilesPath
	errors := make([]string, 0)

	files, err := os.ReadDir(filesDir)
	if err != nil {
		http.Error(w, "Failed to read directory", http.StatusInternalServerError)
		return
	}

	// fileName := "Book1.xlsx" // Replace with your file name
	// filePath := filepath.Join(filesDir, fileName)
	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}

	for _, fileName := range fileNames {
		filePath := filepath.Join(filesDir, fileName)
		fmt.Println("File Path:", filePath)

		// Check if the file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			errMsg := fmt.Sprintln("File not found:", fileName, "Error Message", err)
			errors = append(errors, errMsg)
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		// Process the Excel file
		data, err := processExcel(filePath)
		if err != nil {
			errMsg := fmt.Sprintln("Failed to process file:", fileName, "Error Message", err)
			errors = append(errors, errMsg)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		db.InsertPartsIntoPartTable(data.Parts)

		saveJobErr := db.SaveJobInfoToDB(data.Job)
		if saveJobErr != nil {
			errMsg := fmt.Sprintf("Failed to save job info to database: %v", err)
			logger.LogError("BatchErrors.log", errMsg)
			errors = append(errors, errMsg)
		}

		sortedGroupedPartSlice := part_utils.SortPartsByCode(data.Parts)

		for _, partsByCodeSlice := range sortedGroupedPartSlice {
			materialCode := partsByCodeSlice[0].MaterialCode
			matresults, materr := material_utils.SortMaterialByCode(data.Materials, materialCode)

			if materr != nil {
				logger.LogError("BatchErrors.log", materr.Error())
				continue // Skip this iteration if there's an error
			}

			// Call CreateLayoutV2 with the current partsByCodeSlice and sorted materials
			errSlice := optimizer.CreateLayoutV2(
				partsByCodeSlice,
				matresults,
				data.Job,
			)

			if len(errSlice) > 0 {
				for _, err := range errSlice {
					logger.LogError("BatchErrors.log", err)
					errors = append(errors, err)
				}
			} else {
				// Assuming results is a global or accumulated variable
				fmt.Println("completed slice")
			}
		}

		// errSlice := optimizer.CreateLayoutV2(data.Parts, data.Materials, data.Job)
		// if len(errSlice) > 0 {
		// 	for _, err := range errSlice {
		// 		logger.LogError(err)
		// 		errors = append(errors, err)
		// 	}
		// }
	}

	type response struct {
		Message string   `json:"message"`
		Errors  []string `json:"errors"`
	}

	resp := response{
		Message: "Batch Process Files",
		Errors:  errors,
	}

	// Send the processed data as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func BatchAddFiles(w http.ResponseWriter, r *http.Request) {
	var filesDir = globals.Settings.Excel.FilesPath
	errors := make([]string, 0)

	files, err := os.ReadDir(filesDir)
	if err != nil {
		http.Error(w, "Failed to read directory", http.StatusInternalServerError)
		return
	}

	// fileName := "Book1.xlsx" // Replace with your file name
	// filePath := filepath.Join(filesDir, fileName)
	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}
	// Initialize batchData with empty slices
	var batchData globals.ExcelFileData

	for _, fileName := range fileNames {
		filePath := filepath.Join(filesDir, fileName)
		fmt.Println("File Path:", filePath)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			errMsg := fmt.Sprintln("File not found:", fileName, "Error Message", err)
			errors = append(errors, errMsg)
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		data, err := processMultiExcel(filePath)
		if err != nil {
			errMsg := fmt.Sprintln("Failed to process file:", fileName, "Error Message", err)
			errors = append(errors, errMsg)
			continue // Skip to the next file
		}

		// Merge Parts
		for _, newPart := range data.Parts {
			partExists := false
			for i, existingPart := range batchData.Parts {
				if existingPart.PartNumber == newPart.PartNumber {
					// Part already exists, so add the quantities
					batchData.Parts[i].Quantity += newPart.Quantity
					partExists = true
					break
				}
			}
			if !partExists {
				// Part does not exist, so add it to the batch
				batchData.Parts = append(batchData.Parts, newPart)
			}
		}

		// Merge Materials
		for _, newMaterial := range data.Materials {
			materialExists := false
			for i, existingMaterial := range batchData.Materials {
				if existingMaterial.MaterialCode == newMaterial.MaterialCode && existingMaterial.Length == newMaterial.Length {
					// Material already exists, so add the quantities
					batchData.Materials[i].Quantity += newMaterial.Quantity
					materialExists = true
					break
				}
			}
			if !materialExists {
				// Material does not exist, so add it to the batch
				batchData.Materials = append(batchData.Materials, newMaterial)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(batchData)

}

func processMultiExcel(filePath string) (globals.ExcelFileData, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return globals.ExcelFileData{}, err
	}
	defer f.Close()

	// jobValue, err := f.GetCellValue(globals.Settings.Excel.PartSheet, "B1")
	// if err != nil {
	// 	return globals.ExcelFileData{}, err
	// }

	// customerValue, err := f.GetCellValue(globals.Settings.Excel.PartSheet, "B2")
	// if err != nil {
	// 	return globals.ExcelFileData{}, err
	// }
	partsHeaderRows := int(globals.Settings.Excel.PartHeaderRows)
	partrows, err := f.GetRows(globals.Settings.Excel.PartSheet)
	if err != nil {
		return globals.ExcelFileData{}, err
	}

	if partsHeaderRows < len(partrows) {
		partrows = partrows[partsHeaderRows:]
	}

	materialHeaderRows := int(globals.Settings.Excel.MaterialHeaderRows)

	materialRows, err := f.GetRows(globals.Settings.Excel.MaterialSheet)
	if err != nil {
		return globals.ExcelFileData{}, err
	}

	if materialHeaderRows < len(materialRows) {
		materialRows = materialRows[materialHeaderRows:]
	}
	var data globals.ExcelFileData

	data.Job.Job = ""
	data.Job.Customer = ""

	for _, row := range partrows {

		Length, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return globals.ExcelFileData{}, fmt.Errorf("failed to parse length '%s': %w", row[2], err)
		}
		//
		Quantity, err := strconv.ParseUint(row[3], 10, 16)
		if err != nil {
			return globals.ExcelFileData{}, fmt.Errorf("failed to parse length '%s': %w", row[2], err)
		}

		data.Parts = append(
			data.Parts, globals.Part{
				PartNumber:       row[0],
				MaterialCode:     row[1],
				Length:           Length,
				Quantity:         uint16(Quantity),
				CuttingOperation: row[4],
				// CutQuantity:      0,
			})
	}

	for _, row := range materialRows {
		Length, err := strconv.ParseFloat(row[1], 64)
		if err != nil {

		}
		Quantity, err := strconv.ParseUint(row[2], 10, 16)
		if err != nil {

		}
		data.Materials = append(
			data.Materials, globals.Material{
				MaterialCode: row[0],
				Length:       Length,
				Quantity:     uint16(Quantity),
			})
	}

	return data, nil
}

func processExcel(filePath string) (globals.ExcelFileData, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return globals.ExcelFileData{}, err
	}
	defer f.Close()

	jobValue, err := f.GetCellValue(globals.Settings.Excel.PartSheet, "B1")
	if err != nil {
		return globals.ExcelFileData{}, err
	}

	customerValue, err := f.GetCellValue(globals.Settings.Excel.PartSheet, "B2")
	if err != nil {
		return globals.ExcelFileData{}, err
	}
	partsHeaderRows := int(globals.Settings.Excel.PartHeaderRows)
	partrows, err := f.GetRows(globals.Settings.Excel.PartSheet)
	if err != nil {
		return globals.ExcelFileData{}, err
	}

	if partsHeaderRows < len(partrows) {
		partrows = partrows[partsHeaderRows:]
	}

	materialHeaderRows := int(globals.Settings.Excel.MaterialHeaderRows)

	materialRows, err := f.GetRows(globals.Settings.Excel.MaterialSheet)
	if err != nil {
		return globals.ExcelFileData{}, err
	}

	if materialHeaderRows < len(materialRows) {
		materialRows = materialRows[materialHeaderRows:]
	}

	// Skip the header rows

	// if materialHeaderRows < len(materialRows) {
	// 	materialRows = materialRows[materialHeaderRows:] // Skip the specified number of header rows
	// } else {
	// 	// Handle the case where headerRows is greater than or equal to the total number of rows
	// 	return globals.ExcelFileData{}, fmt.Errorf("header rows exceed the total number of rows")
	// }
	//
	// if partsHeaderRows < len(materialRows) {
	// 	partrows = partrows[partsHeaderRows:] // Skip the specified number of header rows
	// } else {
	// 	// Handle the case where headerRows is greater than or equal to the total number of rows
	// 	return globals.ExcelFileData{}, fmt.Errorf("header rows exceed the total number of rows")
	// }
	var data globals.ExcelFileData

	data.Job.Job = jobValue
	data.Job.Customer = customerValue

	for _, row := range partrows {

		Length, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return globals.ExcelFileData{}, fmt.Errorf("failed to parse length '%s': %w", row[2], err)
		}
		//
		Quantity, err := strconv.ParseUint(row[3], 10, 16)
		if err != nil {
			return globals.ExcelFileData{}, fmt.Errorf("failed to parse length '%s': %w", row[2], err)
		}

		data.Parts = append(
			data.Parts, globals.Part{
				PartNumber:       row[0],
				MaterialCode:     row[1],
				Length:           Length,
				Quantity:         uint16(Quantity),
				CuttingOperation: row[4],
				// CutQuantity:      0,
			})
	}

	for _, row := range materialRows {
		Length, err := strconv.ParseFloat(row[1], 64)
		if err != nil {

		}
		Quantity, err := strconv.ParseUint(row[2], 10, 16)
		if err != nil {

		}
		data.Materials = append(
			data.Materials, globals.Material{
				MaterialCode: row[0],
				Length:       Length,
				Quantity:     uint16(Quantity),
			})
	}

	return data, nil
}
