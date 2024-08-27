package server

import (
	"net/http"
)

func RegisterRoutes() {
	http.HandleFunc("/api/v1/hello", HelloHandler)
	http.HandleFunc("/api/v1/job", HandleGetJob)
	http.HandleFunc("/api/v1/localjobs", HandleGetLocalJobs)
	http.HandleFunc("/api/v1/settings", SettingsHandler)
	http.HandleFunc("/api/v1/togglestar", ToggleStar)
	http.HandleFunc("/api/v1/runproject", RunProject)
	http.HandleFunc("/api/v1/health", CheckHealth)
}
