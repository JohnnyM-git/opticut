package globals

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var Settings SettingsConfig

func SaveSettings(settings SettingsConfig, filePath ...string) error {
	var path string
	if len(filePath) > 0 {
		path = filePath[0]
	} else {
		// Set a default file path
		path = "./globals/settings.json"
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)
	fmt.Println("Saving settings to", path, "...")
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Optional: For pretty-printing JSON
	err = encoder.Encode(settings)
	if err != nil {
		return err
	}
	// Flush the file to ensure all data is written
	err = file.Sync()
	if err != nil {
		return err
	}
	fmt.Println("Settings successfully saved to", path)
	return nil
}

func LoadSettings(filePath ...string) (SettingsConfig, error) {
	var path string
	if len(filePath) > 0 {
		path = filePath[0]
	} else {
		// Set a default file path
		path = "./globals/settings.json"
	}

	fmt.Println("Loading settings from", path, "...")
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return SettingsConfig{}, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	var loadedSettings SettingsConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&loadedSettings)
	if err != nil {
		return SettingsConfig{}, err
	}
	fmt.Println("Settings:", loadedSettings)
	Settings = loadedSettings
	return loadedSettings, nil
}
