package server

import (
	"net/http"
)

func RegisterRoutes() {
	http.HandleFunc("/api/hello", HelloHandler)
	http.HandleFunc("/api/v1/job", HandleGetJob)
	http.HandleFunc("/api/v1/localjobs", HandleGetLocalJobs)
	http.HandleFunc("/api/v1/settings", SettingsHandler)
}
