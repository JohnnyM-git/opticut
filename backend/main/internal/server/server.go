package server

import (
	"log"
	"net/http"

	"main/internal/middleware"
)

func StartServer() {
	RegisterRoutes()
	var port = ":2828"

	handlerWithCORS := middleware.CORS(http.DefaultServeMux)

	log.Println("Server is running on http://localhost", port)
	log.Fatal(http.ListenAndServe(port, handlerWithCORS))
}
