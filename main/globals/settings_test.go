package globals

import (
	"os"
	"testing"
)

const testFilePath = "./test_globals/test_settings.json"

// Helper function to create necessary directories
func createTestDir() error {
	dir := "./globals_test"
	return os.MkdirAll(dir, os.ModePerm)
}

// Helper function to remove test files and directories
func cleanupTestDir() error {
	return os.RemoveAll("./globals_test")
}

func TestSettingsFunctions(t *testing.T) {
	// Define the path for the test settings file
	testFilePath := "./globals_test/test_settings.json"

	// Create the necessary directory
	err := createTestDir()
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer func() {
		// Clean up the directory after the test
		if err := cleanupTestDir(); err != nil {
			t.Errorf("Failed to clean up test directory: %v", err)
		}
	}()

	// Define test settings configuration
	testSettings := SettingsConfig{
		Kerf:  0.1,
		Units: "metric",
		Excel: ExcelSettings{
			FilesPath:          "./test",
			PartSheet:          "Parts",
			PartHeaderRows:     1,
			MaterialSheet:      "Materials",
			MaterialHeaderRows: 1,
		},
	}

	// Save the settings
	err = SaveSettings(testSettings, testFilePath)
	if err != nil {
		t.Fatalf("Failed to save settings: %v", err)
	}

	// Load the settings
	loadedSettings, err := LoadSettings(testFilePath)
	if err != nil {
		t.Fatalf("Failed to load settings: %v", err)
	}

	// Verify the settings file contents
	if loadedSettings != testSettings {
		t.Errorf("Loaded settings do not match saved settings: got %+v, want %+v", loadedSettings, testSettings)
	}
	if Settings != testSettings {
		t.Errorf("Loaded settings did not update globals: got %+v, want %+v", loadedSettings, testSettings)
	}
}

func TestSaveSettingsInvalidPath(t *testing.T) {
	invalidPath := "/invalid/path/to/file.json"

	// Call SaveSettings with an invalid path
	err := SaveSettings(SettingsConfig{}, invalidPath)
	if err == nil {
		t.Fatalf("Expected error while saving settings to path %s, got nil", invalidPath)
	}

	// Check the specific type of error if needed
	if !os.IsNotExist(err) {
		t.Errorf("Expected file not exist error, got %v", err)
	}
}
