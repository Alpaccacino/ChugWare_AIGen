package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"chugware/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─────────────────────────────────────────────────────────────────────────────
// helpers
// ─────────────────────────────────────────────────────────────────────────────

// saveGlobals / restoreGlobals back up and restore the package-level globals
// so individual tests remain isolated.
func saveGlobals() (models.ContestSettings, string) {
	return Settings, ConfigFile
}

func restoreGlobals(s models.ContestSettings, cf string) {
	Settings = s
	ConfigFile = cf
}

// writeConfigFile writes a ContestSettings struct to a temp file and returns its path.
func writeConfigFile(t *testing.T, s models.ContestSettings) string {
	t.Helper()
	data, err := json.MarshalIndent(s, "", "  ")
	require.NoError(t, err)

	path := filepath.Join(t.TempDir(), "chugware_config.json")
	require.NoError(t, os.WriteFile(path, data, 0644))
	return path
}

// ─────────────────────────────────────────────────────────────────────────────
// Constants smoke-test
// ─────────────────────────────────────────────────────────────────────────────

func TestConstants_ExpectedValues(t *testing.T) {
	assert.Equal(t, "ChugWare2", AppName)
	assert.Equal(t, "No", NoKey)
	assert.Equal(t, "Pass", Pass)
	assert.Equal(t, "Disqualified", Disqualified)
	assert.Equal(t, 255, MaxStringLength)
	assert.Equal(t, 3, BottleTries)
	assert.Equal(t, 2, HalfTankardTries)
	assert.Equal(t, 1, FullTankardTries)
}

// ─────────────────────────────────────────────────────────────────────────────
// SaveConfig / LoadConfig
// ─────────────────────────────────────────────────────────────────────────────

func TestSaveConfig_WritesValidJSON(t *testing.T) {
	origSettings, origFile := saveGlobals()
	defer restoreGlobals(origSettings, origFile)

	dir := t.TempDir()
	ConfigFile = filepath.Join(dir, "chugware_config.json")
	Settings = models.ContestSettings{
		FolderPath:        dir,
		ExternalClockPort: "COM3",
		ExternalClockBaud: 9600,
	}

	require.NoError(t, SaveConfig())

	data, err := os.ReadFile(ConfigFile)
	require.NoError(t, err)

	var loaded models.ContestSettings
	require.NoError(t, json.Unmarshal(data, &loaded))
	assert.Equal(t, dir, loaded.FolderPath)
	assert.Equal(t, "COM3", loaded.ExternalClockPort)
	assert.Equal(t, 9600, loaded.ExternalClockBaud)
}

func TestSaveConfig_RoundTrip(t *testing.T) {
	origSettings, origFile := saveGlobals()
	defer restoreGlobals(origSettings, origFile)

	dir := t.TempDir()
	ConfigFile = filepath.Join(dir, "chugware_config.json")
	Settings = models.ContestSettings{
		FolderPath:      "/my/contest/root",
		ParticipantFile: "participants.json",
		ResultFile:      "results.json",
	}

	require.NoError(t, SaveConfig())

	// Reload by reading file and unmarshalling (avoids touching HOME dir)
	data, err := os.ReadFile(ConfigFile)
	require.NoError(t, err)

	var loaded models.ContestSettings
	require.NoError(t, json.Unmarshal(data, &loaded))
	assert.Equal(t, "/my/contest/root", loaded.FolderPath)
	assert.Equal(t, "participants.json", loaded.ParticipantFile)
	assert.Equal(t, "results.json", loaded.ResultFile)
}

// ─────────────────────────────────────────────────────────────────────────────
// Path helpers
// ─────────────────────────────────────────────────────────────────────────────

func TestGetImagePath_WithBase(t *testing.T) {
	origSettings, origFile := saveGlobals()
	defer restoreGlobals(origSettings, origFile)

	Settings.FolderPathContestNameAndDate = "/base/contest"
	expected := filepath.Join("/base/contest", ImagesDirectory)
	assert.Equal(t, expected, GetImagePath())
}

func TestGetImagePath_NoBase(t *testing.T) {
	origSettings, origFile := saveGlobals()
	defer restoreGlobals(origSettings, origFile)

	Settings.FolderPathContestNameAndDate = ""
	assert.Equal(t, "", GetImagePath())
}

func TestGetDiplomaPath_WithBase(t *testing.T) {
	origSettings, origFile := saveGlobals()
	defer restoreGlobals(origSettings, origFile)

	Settings.FolderPathContestNameAndDate = "/base/contest"
	expected := filepath.Join("/base/contest", DiplomasDirectory)
	assert.Equal(t, expected, GetDiplomaPath())
}

func TestGetDiplomaPath_NoBase(t *testing.T) {
	origSettings, origFile := saveGlobals()
	defer restoreGlobals(origSettings, origFile)

	Settings.FolderPathContestNameAndDate = ""
	assert.Equal(t, "", GetDiplomaPath())
}

func TestGetTemplatePath_WithBase(t *testing.T) {
	origSettings, origFile := saveGlobals()
	defer restoreGlobals(origSettings, origFile)

	Settings.FolderPathContestNameAndDate = "/base/contest"
	expected := filepath.Join("/base/contest", TemplateDirectory)
	assert.Equal(t, expected, GetTemplatePath())
}

func TestGetContestPath_WithBase(t *testing.T) {
	origSettings, origFile := saveGlobals()
	defer restoreGlobals(origSettings, origFile)

	Settings.FolderPathContestNameAndDate = "/base/contest"
	expected := filepath.Join("/base/contest", ContestDirectory)
	assert.Equal(t, expected, GetContestPath())
}

func TestGetResultsPath_WithBase(t *testing.T) {
	origSettings, origFile := saveGlobals()
	defer restoreGlobals(origSettings, origFile)

	Settings.FolderPathContestNameAndDate = "/base/contest"
	expected := filepath.Join("/base/contest", ResultsDirectory)
	assert.Equal(t, expected, GetResultsPath())
}

func TestGetResultsPath_NoBase(t *testing.T) {
	origSettings, origFile := saveGlobals()
	defer restoreGlobals(origSettings, origFile)

	Settings.FolderPathContestNameAndDate = ""
	assert.Equal(t, "", GetResultsPath())
}
