package server

import (
	"log"
	"net/http"

	"optimizer/internal/middleware"
)

func StartServer() {
	RegisterRoutes()

	handlerWithCORS := middleware.CORS(http.DefaultServeMux)

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handlerWithCORS))
}
