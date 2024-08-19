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
	script, err := os.ReadFile("./internal/db/setup_db.sql")
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

func SaveJobInfoToDB(jobInfo globals.JobType) {
	if db == nil {
		log.Println("Database is not initialized")
	}
	_, err := db.Exec(
		`INSERT INTO jobs (job_number, customer) VALUES(?, ?)`,
		jobInfo.Job,
		jobInfo.Customer)
	if err != nil {
		logger.LogError(err.Error())
	}

}

func SavePartsToDB(results *[]globals.CutMaterial, jobInfo globals.JobType) {
	if db == nil {
		log.Println("Database is not initialized")
		return
	}

	var jobId int64
	jobId = 1

	for _, result := range *results {
		id, err := db.Exec(
			`INSERT INTO cut_materials (job, job_id, material_code, quantity, stock_length, length) VALUES(?, ?, ?, ?, ?, ?)`,
			result.Job,
			jobId,
			result.MaterialCode,
			result.Quantity,
			result.StockLength,
			result.Length)
		if err != nil {
			fmt.Println("Error inserting cut_materials: ", err)
			logger.LogError(err.Error())
		} else {
			lastID, err := id.LastInsertId()
			if err != nil {
				fmt.Println("Error getting last inserted id: ", err)
				logger.LogError(err.Error())
			}
			for i, part := range result.Parts {
				var partID int64
				err := db.QueryRow(
					`SELECT id FROM parts WHERE part_number = ?`,
					i).Scan(&partID)
				if err != nil {
					fmt.Println("Error getting part id: ", i, err)
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

func GetJobData(job string) ([]globals.CutMaterialPart, error) {
	if db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	query := `
    SELECT
        cm.id AS cut_material_id,
        cm.job,
        cm.job_id,
        cm.material_code AS cut_material_material_code,
        cm.quantity AS cut_material_quantity,
        cm.stock_length,
        cm.length AS cut_material_length,
        p.id AS part_id,
        p.part_number,
        p.material_code AS part_material_code,
        p.length AS part_length,
        cmp.part_qty
    FROM
        cut_materials cm
    JOIN
        cut_material_parts cmp ON cm.id = cmp.cut_material_id
    JOIN
        parts p ON cmp.part_id = p.id
    WHERE
        cm.job = ?
    `

	rows, err := db.Query(query, job)
	if err != nil {
		return nil, fmt.Errorf("query execution error: %v", err)
	}
	defer rows.Close()

	var results []globals.CutMaterialPart
	for rows.Next() {
		var cmp globals.CutMaterialPart
		err := rows.Scan(
			&cmp.CutMaterialID,
			&cmp.Job,
			&cmp.CutMaterialMaterialCode,
			&cmp.CutMaterialQuantity,
			&cmp.StockLength,
			&cmp.CutMaterialLength,
			&cmp.PartID,
			&cmp.PartNumber,
			&cmp.PartMaterialCode,
			&cmp.PartLength,
			&cmp.PartQty,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan error: %v", err)
		}
		results = append(results, cmp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return results, nil
}

func GetMaterialTotals(job string) ([]globals.CutMaterialTotals, error) {
	if db == nil {
		logger.LogError("Database is not initialized")
	}

	query := `SELECT
    material_code,
    stock_length,
    SUM(quantity) AS total_quantity,
    stock_length * SUM(quantity) AS total_length
FROM
    cut_materials
WHERE
    job = ?
GROUP BY
    material_code, stock_length;
`

	rows, err := db.Query(query, job)
	if err != nil {
		return nil, fmt.Errorf("query execution error: %v", err)
	}
	defer rows.Close()
	var results []globals.CutMaterialTotals
	for rows.Next() {
		var total globals.CutMaterialTotals
		err := rows.Scan(
			&total.MaterialCode,
			&total.StockLength,
			&total.TotalQuantity,
			&total.TotalLength,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan error: %v", err)
		}
		results = append(results, total)

	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return results, nil

}

func GetLocalJobs() ([]globals.LocalJobsList, error) {
	if db == nil {
		logger.LogError("Database is not initialized")
	}

	query := `SELECT job_number, customer FROM jobs`

	rows, err := db.Query(query)
	if err != nil {
		logger.LogError(err.Error())
	}
	defer rows.Close()
	var jobs []globals.LocalJobsList
	for rows.Next() {
		var job globals.LocalJobsList
		err := rows.Scan(
			&job.JobNumber,
			&job.Customer,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan error: %v", err)
		}
		jobs = append(jobs, job)

	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return jobs, nil
}
