package server

import (
	"encoding/json"
	"log"
	"net/http"

	"optimizer/globals"
	"optimizer/internal/db"
)

type Response struct {
	Message      string                      `json:"message"`
	JobData      []globals.CutMaterialPart   `json:"JobData"`
	MaterialData []globals.CutMaterialTotals `json:"MaterialData"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
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

	// Respond to the client
	// response := fmt.Sprintf("Job ID: %s", jobID)
	// w.Write([]byte(response))
}

func StartServer() {
	http.HandleFunc("/api/hello", helloHandler)
	http.HandleFunc("/api/v1/job", HandleGetJob)

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
