package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"chugware/internal/config"
)

// FillListFromJSONFile reads and parses a JSON file into a slice of maps
func FillListFromJSONFile(filename string) ([]map[string]string, error) {
	var list []map[string]string

	if !DoesFileExist(filename) {
		return list, nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return list, fmt.Errorf("error reading file %s: %w", filename, err)
	}

	err = json.Unmarshal(data, &list)
	if err != nil {
		return list, fmt.Errorf("error parsing JSON file %s: %w", filename, err)
	}

	return list, nil
}

// SaveListToJSONFile saves a slice of maps to a JSON file
func SaveListToJSONFile(filename string, list []map[string]string) error {
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", dir, err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing file %s: %w", filename, err)
	}

	return nil
}

// RemoveDuplicateDictionaries removes duplicate entries based on a specific key
func RemoveDuplicateDictionaries(list []map[string]string, key string) []map[string]string {
	seen := make(map[string]bool)
	var result []map[string]string

	for _, dict := range list {
		if value, exists := dict[key]; exists {
			if !seen[value] {
				seen[value] = true
				result = append(result, dict)
			}
		}
	}

	return result
}

// AddChugToResultList adds a new result entry to the results list
func AddChugToResultList(list []map[string]string, entry map[string]string) []map[string]string {
	if list == nil {
		list = make([]map[string]string, 0)
	}

	// Validate required fields
	requiredFields := []string{"name", "discipline", "time", "status"}
	for _, field := range requiredFields {
		if _, exists := entry[field]; !exists {
			return list // Invalid entry, don't add
		}
	}

	return append(list, entry)
}

// DoesFileExist checks if a file exists
func DoesFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// IsStringValid checks if string length is within limits
func IsStringValid(str string) bool {
	return len(str) <= config.MaxStringLength && len(str) > 0
}

// IsNullString checks if string is empty or whitespace
func IsNullString(str string) bool {
	return strings.TrimSpace(str) == ""
}

// ParseAndPadTimeString validates and formats time strings
func ParseAndPadTimeString(timeStr string) string {
	// Remove whitespace
	timeStr = strings.TrimSpace(timeStr)

	// Pass NaN through as-is (used for disqualified results)
	if timeStr == "NaN" {
		return "NaN"
	}

	// Regex patterns for time validation
	patterns := []string{
		`^\d{1,2}:\d{2}:\d{2}$`,      // HH:MM:SS or H:MM:SS
		`^\d{1,2}:\d{2}:\d{2}\.\d+$`, // HH:MM:SS.mmm
		`^\d{1,2}:\d{2}$`,            // HH:MM or H:MM
		`^\d{1,2}\.\d+$`,             // SS.mmm
		`^\d+$`,                      // SS
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, timeStr); matched {
			// Format to standardized format
			return formatTimeString(timeStr)
		}
	}

	return config.NoKey
}

// formatTimeString formats various time formats to HH:MM:SS.mmmm
func formatTimeString(timeStr string) string {
	// If already in full format, return as is
	if matched, _ := regexp.MatchString(`^\d{2}:\d{2}:\d{2}\.\d{4}$`, timeStr); matched {
		return timeStr
	}

	// Parse different formats
	if strings.Contains(timeStr, ":") {
		parts := strings.Split(timeStr, ":")
		switch len(parts) {
		case 2: // MM:SS or MM:SS.mmmm
			return fmt.Sprintf("00:%s.0000", timeStr)
		case 3: // HH:MM:SS or HH:MM:SS.mmmm
			h, _ := strconv.Atoi(parts[0])
			m, _ := strconv.Atoi(parts[1])
			if strings.Contains(parts[2], ".") {
				secParts := strings.SplitN(parts[2], ".", 2)
				s, _ := strconv.Atoi(secParts[0])
				ms := secParts[1]
				for len(ms) < 4 {
					ms += "0"
				}
				if len(ms) > 4 {
					ms = ms[:4]
				}
				msInt, _ := strconv.Atoi(ms)
				return fmt.Sprintf("%02d:%02d:%02d.%04d", h, m, s, msInt)
			}
			s, _ := strconv.Atoi(parts[2])
			return fmt.Sprintf("%02d:%02d:%02d.0000", h, m, s)
		}
	} else if strings.Contains(timeStr, ".") {
		// SS.mmmm format
		secParts := strings.SplitN(timeStr, ".", 2)
		s, _ := strconv.Atoi(secParts[0])
		ms := secParts[1]
		for len(ms) < 4 {
			ms += "0"
		}
		if len(ms) > 4 {
			ms = ms[:4]
		}
		msInt, _ := strconv.Atoi(ms)
		return fmt.Sprintf("00:00:%02d.%04d", s, msInt)
	} else {
		// Just seconds (plain integer like "3")
		s, _ := strconv.Atoi(timeStr)
		return fmt.Sprintf("00:00:%02d.0000", s)
	}

	return timeStr
}

// SortListByKey sorts a list of maps by a specific key
func SortListByKey(list []map[string]string, key string) []map[string]string {
	sort.Slice(list, func(i, j int) bool {
		valueI, existsI := list[i][key]
		valueJ, existsJ := list[j][key]

		if !existsI && !existsJ {
			return false
		}
		if !existsI {
			return false
		}
		if !existsJ {
			return true
		}

		// Try to compare as numbers first, then as strings
		numI, errI := strconv.ParseFloat(valueI, 64)
		numJ, errJ := strconv.ParseFloat(valueJ, 64)

		if errI == nil && errJ == nil {
			return numI < numJ
		}

		// If time format, parse as time
		if strings.Contains(valueI, ":") && strings.Contains(valueJ, ":") {
			timeI := ParseTimeForComparison(valueI)
			timeJ := ParseTimeForComparison(valueJ)
			return timeI < timeJ
		}

		return strings.Compare(valueI, valueJ) < 0
	})

	return list
}

// ParseTimeForComparison converts time string to tenths-of-milliseconds for comparison
func ParseTimeForComparison(timeStr string) int64 {
	timeStr = strings.TrimSpace(timeStr)

	// NaN indicates a disqualified / invalid time
	if timeStr == "NaN" {
		return -1
	}

	// Handle HH:MM:SS.mmmm format
	re := regexp.MustCompile(`^(\d{1,2}):(\d{2}):(\d{2})(?:\.(\d+))?$`)
	matches := re.FindStringSubmatch(timeStr)

	if len(matches) >= 4 {
		hours, _ := strconv.Atoi(matches[1])
		minutes, _ := strconv.Atoi(matches[2])
		seconds, _ := strconv.Atoi(matches[3])

		subms := 0
		if len(matches) > 4 && matches[4] != "" {
			// Pad or truncate fractional part to 4 digits (tenths of milliseconds)
			ms := matches[4]
			if len(ms) > 4 {
				ms = ms[:4]
			} else {
				for len(ms) < 4 {
					ms += "0"
				}
			}
			subms, _ = strconv.Atoi(ms)
		}

		return int64(hours*36000000 + minutes*600000 + seconds*10000 + subms)
	}

	return 0
}

// CreateContestFiles creates the necessary contest files and directories
func CreateContestFiles(disciplines []string, contestPath string, contestName string, contestDate string) ([]string, error) {
	var createdFiles []string

	// Ensure contest directory exists
	if err := os.MkdirAll(contestPath, 0755); err != nil {
		return createdFiles, fmt.Errorf("error creating contest directory: %w", err)
	}

	for _, discipline := range disciplines {
		filename := fmt.Sprintf("%s_%s_%s.json", contestName, contestDate, strings.Replace(discipline, " ", "_", -1))
		filepath := filepath.Join(contestPath, filename)

		// Create empty JSON array if file doesn't exist
		if !DoesFileExist(filepath) {
			emptyArray := make([]map[string]string, 0)
			if err := SaveListToJSONFile(filepath, emptyArray); err != nil {
				return createdFiles, fmt.Errorf("error creating file %s: %w", filepath, err)
			}
			createdFiles = append(createdFiles, filename)
		}
	}

	return createdFiles, nil
}

// FilterFiles finds files containing a specific pattern in a directory
func FilterFiles(directory string, pattern string) string {
	if !DoesFileExist(directory) {
		return config.NoKey
	}

	files, err := os.ReadDir(directory)
	if err != nil {
		return config.NoKey
	}

	for _, file := range files {
		if !file.IsDir() && strings.Contains(file.Name(), pattern) {
			return filepath.Join(directory, file.Name())
		}
	}

	return config.NoKey
}

// GetCurrentTimeString returns current time in HH:MM:SS format
func GetCurrentTimeString() string {
	return time.Now().Format("15:04:05")
}

// GetCurrentDateString returns current date in YYYY-MM-DD format
func GetCurrentDateString() string {
	return time.Now().Format("2006-01-02")
}

// CopyFile copies the file at src to dst, creating any missing directories.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source file: %w", err)
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("create destination directory: %w", err)
	}

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create destination file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("copy file: %w", err)
	}
	return nil
}
