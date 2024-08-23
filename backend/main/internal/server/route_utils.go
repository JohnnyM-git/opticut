package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"main/globals"
	"main/internal/db"
	"main/logger"
	"main/material_utils"
	"main/optimizer"
	"main/part_utils"
)

type JobResponse struct {
	Message          string                      `json:"message"`
	JobDataMaterials []globals.CutMaterials      `json:"job_data_materials"`
	MaterialData     []globals.CutMaterialTotals `json:"material_data"`
	Job              globals.JobType             `json:"job_info"`
	JobDataParts     []globals.CutMaterialPart   `json:"job_data_parts"`
}

type LocalJobsResponse struct {
	Message  string                  `json:"Message"`
	JobsList []globals.LocalJobsList `json:"JobsList"`
}

// func HelloHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
//
// 	response := Response{Message: "Hello, World!"}
//
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

func HandleGetJob(w http.ResponseWriter, r *http.Request) {
	// Parse the query parameters
	query := r.URL.Query()
	jobNumber := query.Get("job_id")

	// Use the query parameters in your logic
	if jobNumber == "" {
		http.Error(w, "Missing job_id parameter", http.StatusBadRequest)
		return
	}

	job, jobId, err := db.GetJobInfoFromDB(jobNumber)
	if err != nil {
		fmt.Println("JOB ERR", err.Error())
		logger.LogError(err.Error())
	}
	fmt.Println("job", job)

	jobDataMaterials, err := db.GetJobData(jobId)
	if err != nil {
		logger.LogError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	materialTotals, err := db.GetMaterialTotals(jobId)
	if err != nil {
		fmt.Println("Material Err", err.Error())
		logger.LogError(err.Error())
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	jobDataParts, err := db.GetPartData(jobId)
	fmt.Println("jobDataParts", jobDataParts)

	response := JobResponse{
		Message:          job.Job,
		Job:              job,
		JobDataMaterials: jobDataMaterials,
		MaterialData:     materialTotals,
		JobDataParts:     jobDataParts,
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(
			w,
			"Failed to encode response",
			http.StatusInternalServerError)
	}
}

func HandleGetLocalJobs(w http.ResponseWriter, r *http.Request) {
	localJobs, err := db.GetLocalJobs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	response := LocalJobsResponse{
		Message:  "Local DB Jobs Found",
		JobsList: localJobs,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

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

	// Save the updated settings to the JSON file
	err = globals.SaveSettings(settings)
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

func RunProject(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HTTP Method:", r.Method)
	fmt.Println("Endpoint Hit: Run Project")

	type Part struct {
		PartNumber       string
		MaterialCode     string
		Length           float64
		Quantity         uint16
		CuttingOperation string
	}

	type Material struct {
		MaterialCode string
		Length       float64
		Quantity     uint16
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer r.Body.Close()
	type RunProjectParams struct {
		JobInfo   globals.JobType    `json:"jobInfo"`
		Parts     []globals.Part     `json:"parts"`
		Materials []globals.Material `json:"materials"`
	}

	var params RunProjectParams
	err = json.Unmarshal(body, &params)
	if err != nil {
		fmt.Println("Unmarshal error:", err)
	}
	fmt.Println("Request Body:", params)

	db.InsertPartsIntoPartTable(params.Parts)
	saveJobErr := db.SaveJobInfoToDB(params.JobInfo)
	if saveJobErr != nil {
		errMsg := fmt.Sprintf("Failed to save job info to database: %v", err)
		logger.LogError(errMsg)
	}

	sortedGroupedPartSlice := part_utils.SortPartsByCode(params.Parts)

	for _, partsByCodeSlice := range sortedGroupedPartSlice {
		// fmt.Println(partSlice)
		materialCode := partsByCodeSlice[0].MaterialCode
		results, err := material_utils.SortMaterialByCode(
			params.Materials,
			materialCode)
		if err != nil {
			logger.LogError(err.Error())
		} else {
			errSlice := optimizer.CreateLayout(
				partsByCodeSlice,
				results,
				globals.JobInfo)
			if len(errSlice) > 0 {
				for _, err := range errSlice {
					logger.LogError(err)
				}
			} else {
				fmt.Println(results)
			}
		}

	}

	errSlice := optimizer.CreateLayout(params.Parts, params.Materials, params.JobInfo)
	if errSlice != nil {
		fmt.Println("CreateLayout error:", errSlice)
	}
	response := map[string]interface{}{
		"message": "Project run successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}