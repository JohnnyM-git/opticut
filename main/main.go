package main

import (
	"os"
	"path/filepath"

	"main/globals"
	"main/internal/db"
	"main/internal/server"
)

func main() {

	rootPath := filepath.Join(".", "prod.db")

	// Check if prod.db exists at the root of the project
	if !fileExists(rootPath) {
		db.SetupDB()
	} else {
		db.InitDB("./prod.db")
	}

	globals.LoadSettings()

	server.StartServer()
}

func fileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
