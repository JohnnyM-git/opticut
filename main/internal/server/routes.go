package server

import (
	"net/http"
)

func RegisterRoutes() {
	http.HandleFunc("/api/v1/job", HandleGetJob)
	http.HandleFunc("/api/v1/local-jobs", HandleGetLocalJobs)
	http.HandleFunc("/api/v1/settings", SettingsHandler)
	http.HandleFunc("/api/v1/toggle-star", ToggleStar)
	http.HandleFunc("/api/v1/run-project", RunProject)
	http.HandleFunc("/api/v1/health", CheckHealth)
	http.HandleFunc("/api/v1/file-upload", FileUpload)
	http.HandleFunc("/api/v1/batch-run-files", BatchProcessFiles)
	http.HandleFunc("/api/v1/batch-add-files", BatchAddFiles)
}
