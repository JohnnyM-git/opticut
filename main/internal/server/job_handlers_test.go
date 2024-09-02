package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"main/globals"
	"main/internal/db"
)

func MakeHandleGetJob(dbInterface db.DBInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobNumber := r.URL.Query().Get("job_id")

		// Use the query parameters in your logic
		if jobNumber == "" {
			http.Error(w, "Missing job_id parameter", http.StatusBadRequest)
			return
		}

		job, jobId, err := dbInterface.GetJobInfoFromDB(jobNumber)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jobDataMaterials, err := dbInterface.GetJobData(jobId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		materialTotals, err := dbInterface.GetMaterialTotals(jobId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jobDataParts, err := dbInterface.GetPartData(jobId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := JobResponse{
			Message:          job.Job,
			Job:              job,
			JobDataMaterials: jobDataMaterials,
			MaterialData:     materialTotals,
			JobDataParts:     jobDataParts,
		}

		// Set the Content-Type header and encode the response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func TestHandleGetJob_MissingJobID(t *testing.T) {
	// Create a new HTTP request with no job_id query parameter
	req, err := http.NewRequest("GET", "/getJob", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Record the response using httptest
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleGetJob)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect (400)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Check the response body is what we expect
	expected := "Missing job_id parameter\n"
	if rr.Body.String() != expected {
		t.Errorf(
			"handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHandleGetJob_InvalidJobID(t *testing.T) {
	// Create a new HTTP request with an invalid job_id query parameter
	req, err := http.NewRequest("GET", "/getJob?job_id=invalid", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Record the response using httptest
	rr := httptest.NewRecorder()

	// Create a mock DBInterface
	mockDB := new(db.MockDBInterface)
	mockDB.On("GetJobInfoFromDB", "invalid").Return(globals.JobType{}, nil, fmt.Errorf("job not found"))
	mockDB.On("GetJobData", mock.Anything).Return(nil, fmt.Errorf("job not found"))
	mockDB.On("GetMaterialTotals", mock.Anything).Return(nil, fmt.Errorf("job not found"))
	mockDB.On("GetPartData", mock.Anything).Return(nil, fmt.Errorf("job not found"))

	// Define a handler function that accepts a dbInterface parameter
	handler := MakeHandleGetJob(mockDB)

	// Call the handler
	handler(rr, req)

	// Check the status code is what we expect (500 Internal Server Error)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	// Check the response body is what we expect
	expected := "job not found\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// func TestHandleGetJob_DatabaseError(t *testing.T) {
// 	// Create a new HTTP request with a valid job_id query parameter
// 	req, err := http.NewRequest("GET", "/getJob?job_id=123", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	// Record the response using httptest
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(HandleGetJob)
//
// 	// Mock database functions to return errors
// 	originalGetJobInfoFromDB := db.GetJobInfoFromDB
// 	originalGetJobData := db.GetJobData
// 	originalGetMaterialTotals := db.GetMaterialTotals
// 	originalGetPartData := db.GetPartData
//
// 	defer func() {
// 		db.GetJobInfoFromDB = originalGetJobInfoFromDB
// 		db.GetJobData = originalGetJobData
// 		db.GetMaterialTotals = originalGetMaterialTotals
// 		db.GetPartData = originalGetPartData
// 	}()
//
// 	db.GetJobInfoFromDB = func(jobNumber string) (Job, string, error) {
// 		return Job{}, "", fmt.Errorf("db error")
// 	}
// 	db.GetJobData = func(jobId string) ([]JobDataMaterial, error) {
// 		return nil, fmt.Errorf("db error")
// 	}
// 	db.GetMaterialTotals = func(jobId string) ([]MaterialTotal, error) {
// 		return nil, fmt.Errorf("db error")
// 	}
// 	db.GetPartData = func(jobId string) ([]JobDataPart, error) {
// 		return nil, fmt.Errorf("db error")
// 	}
//
// 	// Call the handler
// 	handler.ServeHTTP(rr, req)
//
// 	// Check the status code is what we expect (500 Internal Server Error)
// 	if status := rr.Code; status != http.StatusInternalServerError {
// 		t.Errorf(
// 			"handler returned wrong status code: got %v want %v",
// 			status, http.StatusInternalServerError)
// 	}
//
// 	// Check the response body is what we expect
// 	expected := "db error\n"
// 	if rr.Body.String() != expected {
// 		t.Errorf(
// 			"handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), expected)
// 	}
// }
//
// func TestHandleGetJob_Success(t *testing.T) {
// 	// Create a new HTTP request with a valid job_id query parameter
// 	req, err := http.NewRequest("GET", "/getJob?job_id=123", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	// Record the response using httptest
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(HandleGetJob)
//
// 	// Mock the database functions to return valid data
// 	originalGetJobInfoFromDB := db.GetJobInfoFromDB
// 	originalGetJobData := db.GetJobData
// 	originalGetMaterialTotals := db.GetMaterialTotals
// 	originalGetPartData := db.GetPartData
//
// 	defer func() {
// 		db.GetJobInfoFromDB = originalGetJobInfoFromDB
// 		db.GetJobData = originalGetJobData
// 		db.GetMaterialTotals = originalGetMaterialTotals
// 		db.GetPartData = originalGetPartData
// 	}()
//
// 	db.GetJobInfoFromDB = func(jobNumber string) (Job, string, error) {
// 		return Job{Job: "Test Job"}, "123", nil
// 	}
// 	db.GetJobData = func(jobId string) ([]JobDataMaterial, error) {
// 		return []JobDataMaterial{{}}, nil
// 	}
// 	db.GetMaterialTotals = func(jobId string) ([]MaterialTotal, error) {
// 		return []MaterialTotal{{}}, nil
// 	}
// 	db.GetPartData = func(jobId string) ([]JobDataPart, error) {
// 		return []JobDataPart{{}}, nil
// 	}
//
// 	// Call the handler
// 	handler.ServeHTTP(rr, req)
//
// 	// Check the status code is what we expect (200 OK)
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf(
// 			"handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}
//
// 	// Check the response body for expected content
// 	// Customize this based on your actual JobResponse structure
// 	expected := `{"Message":"Test Job","Job":{"Job":"Test Job"},"JobDataMaterials":[{}],"MaterialData":[{}],"JobDataParts":[{}]}`
// 	if rr.Body.String() != expected {
// 		t.Errorf(
// 			"handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), expected)
// 	}
// }
