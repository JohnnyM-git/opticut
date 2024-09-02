package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite3 driver
	"main/globals"
	"main/logger"
)

type DBInterface interface {
	GetJobInfoFromDB(jobNumber string) (globals.JobType, *int, error)
	InsertPartsIntoPartTable(parts []globals.Part)
	SaveJobInfoToDB(jobInfo globals.JobType) error
	SavePartsToDB(results *[]globals.CutMaterial, jobInfo globals.JobType)
	GetJobData(jobId *int) ([]globals.CutMaterials, error)
	GetMaterialTotals(jobId *int) ([]globals.CutMaterialTotals, error)
	GetPartData(jobId *int) ([]globals.CutMaterialPart, error)
	GetLocalJobs() ([]globals.LocalJobsList, error)
	ToggleStar(jobNumber string, value int) error
}

var db *Database

type Database struct {
	DB *sql.DB
}

func openDB(dataSourceName string) *sql.DB {
	sqlDB, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	return sqlDB
}

func SetupDB() {
	dbFile := "prod.db"

	// Open the database (creates the file if it doesn't exist)
	sqlDB, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	// Read the SQL script from file
	script, err := os.ReadFile("./internal/db/setup_db.sql")
	if err != nil {
		log.Fatal(err)
	}

	// Execute the SQL script
	_, err = sqlDB.Exec(string(script))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database and tables created successfully!")
	InitDB("./prod.db")
}

// InitDB initializes the database connection
func InitDB(dataSourceName string) {
	var err error
	// Initialize the Database instance
	db = &Database{
		DB: openDB(dataSourceName),
	}

	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Optional: Ping the database to ensure a successful connection
	if err = db.DB.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
}

func GetDatabaseInstance() *Database {
	if db == nil {
		log.Fatal("Database is not initialized")
	}
	return db
}

// CloseDB closes the database connection
func CloseDB() {
	if db != nil {
		if err := db.DB.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}
}

func DbHealthCheck() string {
	err := db.DB.Ping()
	if err != nil {
		return "Unhealthy"
	}
	return "Healthy"
}

func InsertPartsIntoPartTable(parts []globals.Part) {
	if db == nil {
		log.Println("Database is not initialized")
		return
	}

	for _, part := range parts {
		_, err := db.DB.Exec(
			`INSERT OR IGNORE INTO parts (part_number, material_code, length, cutting_operation) VALUES(?, ?, ?, ?)`,
			part.PartNumber,
			part.MaterialCode,
			part.Length,
			part.CuttingOperation)
		if err != nil {
			logger.LogError("DBPartErrors", err.Error())
		}
	}
}

func SaveJobInfoToDB(jobInfo globals.JobType) error {
	if db == nil {
		log.Println("Database is not initialized")
		return errors.New("Database is not initialized")
	}
	_, err := db.DB.Exec(
		`INSERT INTO jobs (job_number, customer) VALUES(?, ?)`,
		jobInfo.Job,
		jobInfo.Customer)
	if err != nil {
		logger.LogError("DBSaveJobErrors", err.Error())
		return errors.New("Error inserting job into database")
	}
	return nil
}

func SavePartsToDB(results *[]globals.CutMaterial, jobInfo globals.JobType) {
	fmt.Println("Hitting SavePartsToDB")
	if db == nil {
		log.Println("Database is not initialized")
		return
	}

	var jobId int64
	err := db.DB.QueryRow(
		`SELECT id FROM jobs WHERE job_number = ?`,
		jobInfo.Job).Scan(&jobId)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.LogError("DBErrors", "No job found with the provided job number")
		} else {
			logger.LogError("DBErrors", err.Error())
		}
	}
	// jobId = 9

	for _, result := range *results {
		id, err := db.DB.Exec(
			`INSERT INTO cut_materials (job, job_id, material_code, quantity, stock_length, length, 
cutting_operation) VALUES(?, ?, ?,
?, ?, ?, ?)`,
			result.Job,
			jobId,
			result.MaterialCode,
			result.Quantity,
			result.StockLength,
			result.Length,
			result.CuttingOperation)
		if err != nil {
			fmt.Println("Error inserting cut_materials: ", err)
			logger.LogError("DBErrors", err.Error())
		} else {
			lastID, err := id.LastInsertId()
			if err != nil {
				fmt.Println("Error getting last inserted id: ", err)
				logger.LogError("DBErrors", err.Error())
			}
			for i, partQty := range result.Parts {
				var partID int64
				var partLength float64
				var materialCode string
				var cuttingOperation string

				err := db.DB.QueryRow(
					`SELECT id, length, material_code, cutting_operation FROM parts WHERE part_number = ?`, i).Scan(
					&partID,
					&partLength,
					&materialCode,
					&cuttingOperation)
				if err != nil {
					fmt.Println("Error getting part id: ", i, err)
					logger.LogError("DBErrors", err.Error())
					continue
				}

				_, err = db.DB.Exec(
					`INSERT INTO cut_material_parts (
            cut_material_id, 
            part_id, 
            part_qty,
            total_part_qty,
            job_id, 
            length, 
            part_cut_length, 
            material_code,
            cutting_operation) 
        VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					lastID,
					partID,
					partQty.CurrentQty, // Use CurrentQty for part_qty
					partQty.TotalQty,   // Use TotalQty for total_part_qty
					jobId,
					partLength,
					partLength+globals.Settings.Kerf,
					materialCode,
					cuttingOperation)
				if err != nil {
					fmt.Println("Error inserting cut material part: ", i, err)
					logger.LogError("DBErrors", err.Error())
					continue
				}
			}

			fmt.Println("Inserted ", lastID, " into database")
		}
	}
}

func GetJobData(jobId *int) ([]globals.CutMaterials, error) {
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
    	cm.stock_length - cm.length AS total_used_material,
    	cm.cutting_operation,
    	 (
        SELECT COUNT(*)
        FROM cut_material_parts cmp
        WHERE cmp.cut_material_id = cm.id
    ) AS unique_parts_qty
    FROM
        cut_materials cm
    WHERE
        cm.job_id = ?
    `

	rows, err := db.DB.Query(query, jobId)
	if err != nil {
		return nil, fmt.Errorf("query execution error: %v", err)
	}
	defer rows.Close()

	var results []globals.CutMaterials
	for rows.Next() {
		var cmp globals.CutMaterials
		err := rows.Scan(
			&cmp.CutMaterialID,
			&cmp.Job,
			&cmp.JobId,
			&cmp.CutMaterialMaterialCode,
			&cmp.CutMaterialQuantity,
			&cmp.StockLength,
			&cmp.CutMaterialLength,
			&cmp.TotalUsedLength,
			&cmp.CuttingOperation,
			&cmp.TotalPartsCutOnMaterial)
		if err != nil {
			fmt.Println("Error scanning rows: ", err)
			return nil, fmt.Errorf("row scan error: %v", err)
		}
		results = append(results, cmp)
	}
	fmt.Println("DB_UTILS RESULTS:", results)

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return results, nil
}

func GetMaterialTotals(jobId *int) ([]globals.CutMaterialTotals, error) {
	if db == nil {
		logger.LogError("DBErrors", "Database is not initialized")
		return nil, fmt.Errorf("database is not initialized")
	}

	query := `SELECT
    id,
    material_code,
    stock_length,
    length AS remaining_length,
    SUM(quantity) AS total_quantity,
    stock_length * SUM(quantity) AS total_stock_length,
    (stock_length * SUM(quantity)) - SUM(length) AS total_used_length,
    cutting_operation
FROM
    cut_materials
WHERE
    job_id = ?
GROUP BY
    material_code, stock_length;`

	rows, err := db.DB.Query(query, jobId)
	if err != nil {
		return nil, fmt.Errorf("query execution error: %v", err)
	}
	defer rows.Close()

	var results []globals.CutMaterialTotals
	for rows.Next() {
		var total globals.CutMaterialTotals
		err := rows.Scan(
			&total.Id,
			&total.MaterialCode,
			&total.StockLength,
			&total.Length, // length AS remaining_length
			&total.TotalQuantity,
			&total.TotalStockLength,
			&total.TotalUsedLength, // stock_length - length * SUM(quantity)
			&total.CuttingOperation,
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

func GetPartData(jobId *int) ([]globals.CutMaterialPart, error) {
	if db == nil {
		logger.LogError("DBErrors", "Database is not initialized")
		return nil, fmt.Errorf("database is not initialized")
	}

	query := `SELECT
	    cmp.cut_material_id,
	    cmp.part_id,
	    p.part_number,
	    p.material_code,
	    p.length,
	    cmp.part_cut_length,
	    cmp.part_qty,
	    cmp.total_part_qty,
	    cmp.cutting_operation
	FROM
	    cut_material_parts cmp
	    JOIN parts p ON cmp.part_id = p.id
	WHERE
	    cmp.job_id = ?`

	rows, err := db.DB.Query(query, *jobId)
	if err != nil {
		logger.LogError("DBErrors", fmt.Sprintf("Error executing query: %v", err))
		return nil, fmt.Errorf("query execution error: %v", err)
	}
	defer rows.Close()

	var results []globals.CutMaterialPart
	for rows.Next() {
		var cmp globals.CutMaterialPart
		err := rows.Scan(
			&cmp.CutMaterialID,
			&cmp.PartID,
			&cmp.PartNumber,
			&cmp.PartMaterialCode,
			&cmp.PartLength,
			&cmp.PartCutLength,
			&cmp.PartQty,
			&cmp.TotalPartQty,
			&cmp.CuttingOperation)
		if err != nil {
			logger.LogError("DBErrors", fmt.Sprintf("Row scan error: %v", err))
			return nil, fmt.Errorf("row scan error: %v", err)
		}
		results = append(results, cmp)
	}

	if err = rows.Err(); err != nil {
		logger.LogError("DBErrors", fmt.Sprintf("Rows error: %v", err))
		return nil, fmt.Errorf("rows error: %v", err)
	}

	fmt.Println("func GetPartData result", results)
	return results, nil
}

func ToggleStar(jobNumber string, value int) error {
	if db == nil {
		logger.LogError("DBErrors", "Database is not initialized")
		return fmt.Errorf("database is not initialized")
	}

	// Prepare the query
	query := `UPDATE jobs SET star = ? WHERE job_number = ?`

	// Execute the query
	result, err := db.DB.Exec(query, value, jobNumber)
	if err != nil {
		return fmt.Errorf("query execution error: %v", err)
	}

	// Check the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error retrieving rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf(
			"no rows updated, job number %s may not exist",
			jobNumber)
	}

	fmt.Printf("Successfully updated %d rows\n", rowsAffected)
	return nil
}
