package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"optimizer/globals"
	"optimizer/internal/db"
)

type Response struct {
	Message      string                      `json:"message"`
	JobData      []globals.CutMaterialPart   `json:"JobData"`
	MaterialData []globals.CutMaterialTotals `json:"MaterialData"`
}

type LocalJobsResponse struct {
	Message  string                  `json:"Message"`
	JobsList []globals.LocalJobsList `json:"JobsList"`
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := Response{Message: "Hello, World!"}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HandleGetJob(w http.ResponseWriter, r *http.Request) {
	// Parse the query parameters
	query := r.URL.Query()
	jobID := query.Get("job_id")

	// Use the query parameters in your logic
	if jobID == "" {
		http.Error(w, "Missing job_id parameter", http.StatusBadRequest)
		return
	}

	jobData, err := db.GetJobData(jobID)
	materialTotals, err := db.GetMaterialTotals(jobID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	response := Response{
		Message:      jobID,
		JobData:      jobData,
		MaterialData: materialTotals,
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

// func LoadSettings() (globals.SettingsConfig, error) {
// 	var settings globals.SettingsConfig
// 	var filename = "./globals/settings.json"
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return settings, err
// 	}
// 	defer file.Close()
//
// 	byteValue, _ := ioutil.ReadAll(file)
// 	json.Unmarshal(byteValue, &settings)
//
// 	return settings, nil
// }

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

	// Parse the JSON body into a Settings struct
	var newSettings globals.SettingsConfig
	err = json.Unmarshal(body, &newSettings)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert Kerf from string to float64
	kerfValue, err := strconv.ParseFloat(newSettings.Kerf, 32)
	if err != nil {
		http.Error(w, "Invalid kerf value", http.StatusBadRequest)
		return
	}

	// Load current settings
	settings, err := globals.LoadSettings()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Current Settings:", settings)

	// Update settings with new values
	settings.Kerf = kerfValue

	// Save the updated settings to the JSON file
	err = globals.SaveSettings(settings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Settings updated successfully")

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Settings updated successfully"))
}
