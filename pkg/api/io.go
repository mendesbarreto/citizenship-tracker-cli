package api

import (
	"citizenship-tracker-cli/pkg/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// getFilePath returns the full path to the lastupdate.json file
func getFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".citizenship", "lastupdate.json"), nil
}

// LoadStatusResponse loads the StatusResponse from the JSON file
func LoadStatusResponse() (*model.StatusResponse, error) {
	filePath, err := getFilePath()
	if err != nil {
		return nil, err
	}

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &model.StatusResponse{}, nil // Return empty status if file doesn't exist
	}

	// Read the file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse the JSON
	var status model.StatusResponse
	if err := json.Unmarshal(data, &status); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &status, nil
}

// SaveStatusResponse saves the StatusResponse to the JSON file
func SaveStatusResponse(status *model.StatusResponse) error {
	filePath, err := getFilePath()
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal the status to JSON
	data, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to file
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
