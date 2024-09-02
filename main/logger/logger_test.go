package logger

import (
	"bufio"
	"os"
	"testing"
)

func TestLogError(t *testing.T) {
	// Path for the log file used by LogError
	logFilePath := "testErrors.log"

	// Remove any existing log file
	if err := os.Remove(logFilePath); err != nil && !os.IsNotExist(err) {
		t.Fatalf("failed to remove existing log file: %v", err)
	}

	// Call LogError function with the updated signature
	LogError(logFilePath, "This is a test error message")

	// Check if the log file was created
	fileInfo, err := os.Stat(logFilePath)
	if err != nil {
		t.Fatalf("failed to stat log file: %v", err)
	}

	// Check if the file size is greater than 0
	if fileInfo.Size() == 0 {
		t.Errorf("log file is empty")
	}

	// Read the log file and check the content
	file, err := os.Open(logFilePath)
	if err != nil {
		t.Fatalf("failed to open log file: %v", err)
	}
	defer file.Close() // Ensure file is closed before removal

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		t.Fatalf("failed to read log file")
	}

	logLine := scanner.Text()
	expectedPrefix := "Error: "

	if !startsWith(logLine, expectedPrefix) {
		t.Errorf("log line does not start with expected prefix: got %v, want prefix %v", logLine, expectedPrefix)
	}

	expectedMsg := "This is a test error message"
	if !contains(logLine, expectedMsg) {
		t.Errorf("log line does not contain expected message: got %v, want %v", logLine, expectedMsg)
	}

	// Clean up
	file.Close() // Ensure file is closed before attempting to remove
	if err := os.Remove(logFilePath); err != nil {
		t.Fatalf("failed to remove test log file: %v", err)
	}
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr
}
