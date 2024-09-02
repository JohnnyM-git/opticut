package db

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"main/globals"
)

// MockDB is a mock implementation of the Database interface for testing
type MockDB struct{}

func (m *MockDB) GetJobInfoFromDB(jobNumber string) (globals.JobType, *int, error) {
	if jobNumber == "invalid" {
		return globals.JobType{}, nil, errors.New("job not found")
	}
	// Simulate a valid job
	job := globals.JobType{JobId: 1, Job: "Test Job", Customer: "Test Customer", Star: 5}
	return job, &job.JobId, nil
}

func TestGetJobInfoFromDB(t *testing.T) {
	// Create a mock database instance
	mockDB := new(MockDBInterface)

	// Define the expected values
	expectedJob := globals.JobType{}
	var expectedJobId *int
	expectedErr := fmt.Errorf("job not found")

	// Set up the mock to return these values
	mockDB.On("GetJobInfoFromDB", "invalid").Return(expectedJob, expectedJobId, expectedErr)

	// Call the method
	actualJob, actualJobId, actualErr := mockDB.GetJobInfoFromDB("invalid")

	// Check if the returned values match the expected values
	if actualJob != expectedJob || actualJobId != expectedJobId || actualErr != expectedErr {
		t.Errorf(
			"GetJobInfoFromDB() = %v, %v, %v; want %v, %v, %v",
			actualJob,
			actualJobId,
			actualErr,
			expectedJob,
			expectedJobId,
			expectedErr)
	}
}

// Helper function to create an *int from an int
func intPointer(i int) *int {
	return &i
}

func TestGetLocalJobs(t *testing.T) {
	mockDB := new(MockDBInterface)
	expectedJobs := []globals.LocalJobsList{
		{JobNumber: "Job1", Customer: "Customer1", Star: 0},
		{JobNumber: "Job2", Customer: "Customer2", Star: 1},
	}
	expectedErr := errors.New("database error")

	// Test case: successful query
	mockDB.On("GetLocalJobs").Return(expectedJobs, nil)

	jobs, err := mockDB.GetLocalJobs()
	assert.NoError(t, err)
	assert.Equal(t, expectedJobs, jobs)

	// Reset the mock for the next test case
	mockDB = new(MockDBInterface)

	// Test case: query error
	mockDB.On("GetLocalJobs").Return(([]globals.LocalJobsList)(nil), expectedErr)

	jobs, err = mockDB.GetLocalJobs()
	assert.Error(t, err)
	assert.Nil(t, jobs)
	assert.Equal(t, expectedErr, err)

	// Ensure all expectations were met
	mockDB.AssertExpectations(t)
}
