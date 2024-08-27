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

func DbHealthCheck() string {
	err := db.Ping()
	if err != nil {
		return "unhealthy"
	}
	return "healthy"
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

func SaveJobInfoToDB(jobInfo globals.JobType) error {
	if db == nil {
		log.Println("Database is not initialized")
		return errors.New("Database is not initialized")
	}
	_, err := db.Exec(
		`INSERT INTO jobs (job_number, customer) VALUES(?, ?)`,
		jobInfo.Job,
		jobInfo.Customer)
	if err != nil {
		logger.LogError(err.Error())
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
	err := db.QueryRow(
		`SELECT id FROM jobs WHERE job_number = ?`,
		jobInfo.Job).Scan(&jobId)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.LogError("No job found with the provided job number")
		} else {
			logger.LogError(err.Error())
		}
	}
	// jobId = 9

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
			for i, partQty := range result.Parts {
				var partID int64
				var partLength float64
				var materialCode string

				err := db.QueryRow(
					`SELECT id, length, material_code FROM parts WHERE part_number = ?`, i).Scan(
					&partID,
					&partLength,
					&materialCode)
				if err != nil {
					fmt.Println("Error getting part id: ", i, err)
					logger.LogError(err.Error())
					continue
				}

				_, err = db.Exec(
					`INSERT INTO cut_material_parts (
            cut_material_id, 
            part_id, 
            part_qty,
            total_part_qty,
            job_id, 
            length, 
            part_cut_length, 
            material_code) 
        VALUES(?, ?, ?, ?, ?, ?, ?, ?)`,
					lastID,
					partID,
					partQty.CurrentQty, // Use CurrentQty for part_qty
					partQty.TotalQty,   // Use TotalQty for total_part_qty
					jobId,
					partLength,
					(partLength + globals.Settings.Kerf),
					materialCode,
				)
				if err != nil {
					fmt.Println("Error inserting cut material part: ", i, err)
					logger.LogError(err.Error())
					continue
				}
			}

			fmt.Println("Inserted ", lastID, " into database")
		}
	}
}

func GetJobInfoFromDB(jobNumber string) (globals.JobType, *int, error) {
	// Create a variable to hold the result
	var job globals.JobType
	if db == nil {
		logger.LogError("Database is not initialized")
		return job, nil, errors.New("database is not initialized")
	}

	query := `SELECT id, job_number, customer, 
star FROM jobs WHERE job_number = ?`

	rows, err := db.Query(query, jobNumber)
	if err != nil {
		logger.LogError(err.Error())
		return job, nil, fmt.Errorf("query execution error: %v", err)
	}
	defer rows.Close()

	// Iterate over the result set
	if rows.Next() {
		err := rows.Scan(&job.JobId, &job.Job, &job.Customer, &job.Star)
		if err != nil {
			logger.LogError(err.Error())
			return job, nil, fmt.Errorf("row scan error: %v", err)
		}
		// Return the job information and no error
		return job, &job.JobId, nil
	}

	// If no rows are found, return an appropriate error
	if err = rows.Err(); err != nil {
		logger.LogError(err.Error())
		return job, nil, fmt.Errorf("rows error: %v", err)
	}

	return job, nil, errors.New("job not found")
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

	rows, err := db.Query(query, jobId)
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
			&cmp.TotalPartsCutOnMaterial,
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

func GetMaterialTotals(jobId *int) ([]globals.CutMaterialTotals, error) {
	if db == nil {
		logger.LogError("Database is not initialized")
		return nil, fmt.Errorf("database is not initialized")
	}

	query := `SELECT
    id,
    material_code,
    stock_length,
    length AS remaining_length,
    SUM(quantity) AS total_quantity,
    stock_length * SUM(quantity) AS total_stock_length,
    (stock_length - length) * SUM(quantity) AS total_used_length
FROM
    cut_materials
WHERE
    job_id = ?
GROUP BY
    material_code, stock_length;`

	rows, err := db.Query(query, jobId)
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
		logger.LogError("Database is not initialized")
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
    cmp.total_part_qty
FROM
    cut_material_parts cmp
        JOIN
    parts p ON cmp.part_id = p.id
WHERE
    cmp.job_id = ?`

	rows, err := db.Query(query, *jobId)
	if err != nil {
		fmt.Println("func GerPartData rows, err", err)
		logger.LogError(err.Error())
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
			&cmp.TotalPartQty)
		if err != nil {
			return nil, fmt.Errorf("row scan error: %v", err)
		}
		results = append(results, cmp)
	}
	if err = rows.Err(); err != nil {
		logger.LogError(err.Error())
		return nil, fmt.Errorf("rows error: %v", err)
	}
	fmt.Println("func GetPartData result", results)
	return results, nil

}

func GetLocalJobs() ([]globals.LocalJobsList, error) {
	if db == nil {
		logger.LogError("Database is not initialized")
	}

	query := `SELECT job_number, customer, star FROM jobs`

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
			&job.Star,
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

func ToggleStar(jobNumber string, value int) error {
	if db == nil {
		logger.LogError("Database is not initialized")
		return fmt.Errorf("database is not initialized")
	}

	// Prepare the query
	query := `UPDATE jobs SET star = ? WHERE job_number = ?`

	// Execute the query
	result, err := db.Exec(query, value, jobNumber)
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
