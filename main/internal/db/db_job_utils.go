package db

import (
	"errors"
	"fmt"

	"main/globals"
	"main/logger"
)

func (db *Database) GetJobInfoFromDB(jobNumber string) (globals.JobType, *int, error) {

	// Create a variable to hold the result
	var job globals.JobType
	if db == nil {
		logger.LogError("DBErrors", "Database is not initialized")
		return job, nil, errors.New("database is not initialized")
	}

	query := `SELECT id, job_number, customer, 
star FROM jobs WHERE job_number = ?`

	rows, err := db.DB.Query(query, jobNumber)
	if err != nil {
		logger.LogError("DBErrors", err.Error())
		return job, nil, fmt.Errorf("query execution error: %v", err)
	}
	defer rows.Close()

	// Iterate over the result set
	if rows.Next() {
		err := rows.Scan(&job.JobId, &job.Job, &job.Customer, &job.Star)
		if err != nil {
			logger.LogError("DBErrors", err.Error())
			return job, nil, fmt.Errorf("row scan error: %v", err)
		}
		// Return the job information and no error
		return job, &job.JobId, nil
	}

	// If no rows are found, return an appropriate error
	if err = rows.Err(); err != nil {
		logger.LogError("DBErrors", err.Error())
		return job, nil, fmt.Errorf("rows error: %v", err)
	}

	return job, nil, errors.New("job not found")
}

func GetLocalJobs() ([]globals.LocalJobsList, error) {
	if db == nil {
		logger.LogError("DBErrors", "Database is not initialized")
	}

	query := `SELECT job_number, customer, star FROM jobs`

	rows, err := db.DB.Query(query)
	if err != nil {
		logger.LogError("DBErrors", err.Error())
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
