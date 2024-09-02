package db

import (
	"database/sql"

	"github.com/stretchr/testify/mock"
	"main/globals"
)

type MockDBInterface struct {
	mock.Mock
}

// Query is a mock method for simulating database queries.
func (m *MockDBInterface) Query(query string, args ...interface{}) (*sql.Rows, error) {
	argsList := m.Called(query, args)
	return argsList.Get(0).(*sql.Rows), argsList.Error(1)
}

// MockRowScanner simulates row scanning.
type MockRowScanner struct {
	mock.Mock
}

// Scan is a mock method to simulate row scanning.
func (m *MockRowScanner) Scan(dest ...interface{}) error {
	args := m.Called(dest)
	return args.Error(0)
}

// Implement GetJobInfoFromDB
func (m *MockDBInterface) GetJobInfoFromDB(jobNumber string) (globals.JobType, *int, error) {
	args := m.Called(jobNumber)
	var jobId *int
	if args.Get(1) != nil {
		id := args.Get(1).(*int)
		jobId = id
	}
	return args.Get(0).(globals.JobType), jobId, args.Error(2)
}

// Implement InsertPartsIntoPartTable
func (m *MockDBInterface) InsertPartsIntoPartTable(parts []globals.Part) {
	m.Called(parts)
}

// Implement SaveJobInfoToDB
func (m *MockDBInterface) SaveJobInfoToDB(jobInfo globals.JobType) error {
	args := m.Called(jobInfo)
	return args.Error(0)
}

// Implement SavePartsToDB
func (m *MockDBInterface) SavePartsToDB(results *[]globals.CutMaterial, jobInfo globals.JobType) {
	m.Called(results, jobInfo)
}

// Implement GetJobData
func (m *MockDBInterface) GetJobData(jobId *int) ([]globals.CutMaterials, error) {
	args := m.Called(jobId)
	return args.Get(0).([]globals.CutMaterials), args.Error(1)
}

// Implement GetMaterialTotals
func (m *MockDBInterface) GetMaterialTotals(jobId *int) ([]globals.CutMaterialTotals, error) {
	args := m.Called(jobId)
	return args.Get(0).([]globals.CutMaterialTotals), args.Error(1)
}

// Implement GetPartData
func (m *MockDBInterface) GetPartData(jobId *int) ([]globals.CutMaterialPart, error) {
	args := m.Called(jobId)
	return args.Get(0).([]globals.CutMaterialPart), args.Error(1)
}

// Implement GetLocalJobs
func (m *MockDBInterface) GetLocalJobs() ([]globals.LocalJobsList, error) {
	args := m.Called()
	return args.Get(0).([]globals.LocalJobsList), args.Error(1)
}

// Implement ToggleStar
func (m *MockDBInterface) ToggleStar(jobNumber string, value int) error {
	args := m.Called(jobNumber, value)
	return args.Error(0)
}
