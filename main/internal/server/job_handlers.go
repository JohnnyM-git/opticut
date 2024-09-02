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

// Define a new type for the handler function that accepts the dbInterface

func HandleGetJob(w http.ResponseWriter, r *http.Request) {

	jobNumber := r.URL.Query().Get("job_id")

	// Use the query parameters in your logic
	if jobNumber == "" {
		http.Error(w, "Missing job_id parameter", http.StatusBadRequest)
		return
	}

	database := db.GetDatabaseInstance()

	job, jobId, err := database.GetJobInfoFromDB(jobNumber)
	if err != nil {
		fmt.Println("JOB ERR", err.Error())
		logger.LogError("JobErrors.log", err.Error())
	}
	fmt.Println("job", job)

	jobDataMaterials, err := db.GetJobData(jobId)
	if err != nil {
		logger.LogError("JobErrors.log", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	materialTotals, err := db.GetMaterialTotals(jobId)
	if err != nil {
		fmt.Println("Material Err", err.Error())
		logger.LogError("JobErrors.log", err.Error())
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

	// responseBytes, _ := json.MarshalIndent(response, "", "  ")
	// fmt.Println(string(responseBytes))
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

func RunProject(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HTTP Method:", r.Method)
	fmt.Println("Endpoint Hit: Run Project")

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
		logger.LogError("ProjectRun.log", errMsg)
	}

	sortedGroupedPartSlice := part_utils.SortPartsByCode(params.Parts)

	for _, partsByCodeSlice := range sortedGroupedPartSlice {
		materialCode := partsByCodeSlice[0].MaterialCode
		matresults, materr := material_utils.SortMaterialByCode(params.Materials, materialCode)

		if materr != nil {
			logger.LogError("ProjectRun.log", materr.Error())
			continue // Skip this iteration if there's an error
		}

		// Call CreateLayoutV2 with the current partsByCodeSlice and sorted materials
		errSlice := optimizer.CreateLayoutV2(
			partsByCodeSlice,
			matresults,
			params.JobInfo,
		)

		if len(errSlice) > 0 {
			for _, err := range errSlice {
				logger.LogError("ProjectRun.log", err)
			}
		} else {
			// Assuming results is a global or accumulated variable
			fmt.Println("completed slice")
		}
	}

	response := map[string]interface{}{
		"message": "Project run successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
