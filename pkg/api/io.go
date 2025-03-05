package api

import (
	"citizenship-tracker-cli/pkg/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func getFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".citizenship", "lastupdate.json"), nil
}

func LoadStatusResponse() (*model.StatusResponse, error) {
	filePath, err := getFilePath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &model.StatusResponse{}, nil // Return empty status if file doesn't exist
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var status model.StatusResponse
	if err := json.Unmarshal(data, &status); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &status, nil
}

func SaveStatusResponse(status *model.StatusResponse) error {
	filePath, err := getFilePath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
