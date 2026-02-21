package config

import (
	"chugware/internal/models"
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	// Application Info
	AppName    = "ChugWare"
	AppVersion = "2.0.0"
	Authors    = "Sigge McKvack, EKAK-2012 (Go Port by AI)"

	// Directories
	ImagesDirectory   = "images"
	DiplomasDirectory = "diplomas"
	TemplateDirectory = "template"
	ContestDirectory  = "contest"
	ResultsDirectory  = "results"

	// Constraints
	MaxStringLength = 255
	MaxParticipants = 200
	MaxEntries      = 999
	MaxValue        = 9999

	// Chugging Attempts
	BottleTries      = 3
	HalfTankardTries = 2
	FullTankardTries = 1

	// Keys
	OfficialKey     = "Official"
	UnOfficialKey   = "Unofficial"
	NoKey           = "No"
	ResultsKey      = "Results"
	ParticipantsKey = "Participants"
	Pass            = "Pass"
	Disqualified    = "Disqualified"

	// UI Constants
	PixelsPerChar = 8

	// Config file
	ConfigFileName = "chugware_config.json"
)

var (
	// Global settings
	Settings   models.ContestSettings
	ConfigFile string
)

// LoadConfig loads configuration from file or creates default
func LoadConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	ConfigFile = filepath.Join(homeDir, ".chugware", ConfigFileName)

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(ConfigFile)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Get current working directory (where binary is executed)
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Load existing config or create default
	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		// Create default config using current directory
		Settings = models.ContestSettings{
			FolderPath: filepath.Join(currentDir, "ChugWare"),
		}
		return SaveConfig()
	}

	// Read existing config
	data, err := os.ReadFile(ConfigFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &Settings)
	if err != nil {
		return err
	}

	// Update FolderPath to current directory if it's empty or still using old Documents path
	documentsPath := filepath.Join(homeDir, "Documents", "ChugWare")
	if Settings.FolderPath == "" || Settings.FolderPath == documentsPath {
		Settings.FolderPath = filepath.Join(currentDir, "ChugWare")
		// Save the updated config
		SaveConfig()
	}

	return nil
}

// SaveConfig saves current configuration to file
func SaveConfig() error {
	data, err := json.MarshalIndent(Settings, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(ConfigFile, data, 0644)
}

// GetImagePath returns the full path for images directory
func GetImagePath() string {
	if Settings.FolderPathContestNameAndDate == "" {
		return ""
	}
	return filepath.Join(Settings.FolderPathContestNameAndDate, ImagesDirectory)
}

// GetDiplomaPath returns the full path for diplomas directory
func GetDiplomaPath() string {
	if Settings.FolderPathContestNameAndDate == "" {
		return ""
	}
	return filepath.Join(Settings.FolderPathContestNameAndDate, DiplomasDirectory)
}

// GetTemplatePath returns the full path for template directory
func GetTemplatePath() string {
	if Settings.FolderPathContestNameAndDate == "" {
		return ""
	}
	return filepath.Join(Settings.FolderPathContestNameAndDate, TemplateDirectory)
}

// GetContestPath returns the full path for contest directory
func GetContestPath() string {
	if Settings.FolderPathContestNameAndDate == "" {
		return ""
	}
	return filepath.Join(Settings.FolderPathContestNameAndDate, ContestDirectory)
}

// GetResultsPath returns the full path for results directory
func GetResultsPath() string {
	if Settings.FolderPathContestNameAndDate == "" {
		return ""
	}
	return filepath.Join(Settings.FolderPathContestNameAndDate, ResultsDirectory)
}
