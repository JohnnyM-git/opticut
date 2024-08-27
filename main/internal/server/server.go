package server

import (
	"log"
	"net/http"
	"time"

	"main/internal/middleware"
)

var StartTime time.Time

func StartServer() {
	RegisterRoutes()
	var port = ":2828"

	handlerWithCORS := middleware.CORS(http.DefaultServeMux)
	StartTime = time.Now()
	log.Println("Server is running on http://localhost", port)
	log.Fatal(http.ListenAndServe(port, handlerWithCORS))
}
