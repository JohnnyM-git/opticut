package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite3 driver
	"optimizer/globals"
	"optimizer/logger"
)

var db *sql.DB

func SetupDB() {
	dbFile := "prod.db"

	// Open the database (creates the file if it doesn't exist)
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Read the SQL script from file
	script, err := os.ReadFile("./db/setup_db.sql")
	if err != nil {
		log.Fatal(err)
	}

	// Execute the SQL script
	_, err = db.Exec(string(script))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database and tables created successfully!")
	InitDB("./prod.db")
}

// InitDB initializes the database connection
func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Optional: Ping the database to ensure a successful connection
	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
}

// CloseDB closes the database connection
func CloseDB() {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}
}

func InsertPartsIntoPartTable(parts []globals.Part) {
	if db == nil {
		log.Println("Database is not initialized")
		return
	}

	for _, part := range parts {
		_, err := db.Exec(
			`INSERT OR IGNORE INTO parts (part_number, material_code, length) VALUES(?, ?, ?)`,
			part.PartNumber,
			part.MaterialCode,
			part.Length)
		if err != nil {
			logger.LogError(err.Error())
		}
	}
}

func SavePartsToDB(results *[]globals.CutMaterial) {
	if db == nil {
		log.Println("Database is not initialized")
		return
	}

	for _, result := range *results {
		id, err := db.Exec(
			`INSERT INTO cut_materials (job, material_code, quantity, stock_length, length) VALUES(?, ?, ?, ?, ?)`,
			result.Job,
			result.MaterialCode,
			result.Quantity,
			result.StockLength,
			result.Length)
		if err != nil {
			logger.LogError(err.Error())
		} else {
			lastID, err := id.LastInsertId()
			if err != nil {
				logger.LogError(err.Error())
			}
			for _, part := range result.Parts {
				var partID int64
				err := db.QueryRow(
					`SELECT id FROM parts WHERE part_number = ?`,
					part).Scan(&partID)
				if err != nil {
					logger.LogError(err.Error())
					continue
				}
				_, err = db.Exec(
					`INSERT INTO cut_material_parts (cut_material_id, part_id, part_qty) VALUES(?, ?, ?)`,
					lastID,
					partID,
					part,
				)
			}
			fmt.Println("Inserted ", lastID, " into database")
		}
	}
}
