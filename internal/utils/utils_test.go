package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─────────────────────────────────────────────────────────────────────────────
// DoesFileExist
// ─────────────────────────────────────────────────────────────────────────────

func TestDoesFileExist_ExistingFile(t *testing.T) {
	f, err := os.CreateTemp("", "cwtest-*.json")
	require.NoError(t, err)
	f.Close()
	defer os.Remove(f.Name())

	assert.True(t, DoesFileExist(f.Name()))
}

func TestDoesFileExist_MissingFile(t *testing.T) {
	assert.False(t, DoesFileExist("/totally/nonexistent/path/file.json"))
}

// ─────────────────────────────────────────────────────────────────────────────
// IsStringValid
// ─────────────────────────────────────────────────────────────────────────────

func TestIsStringValid(t *testing.T) {
	assert.True(t, IsStringValid("hello"))
	assert.True(t, IsStringValid("a"))
	assert.False(t, IsStringValid(""), "empty string should be invalid")
	assert.False(t, IsStringValid(strings.Repeat("x", 256)), "string > MaxStringLength should be invalid")
	assert.True(t, IsStringValid(strings.Repeat("x", 255)), "string == MaxStringLength should be valid")
}

// ─────────────────────────────────────────────────────────────────────────────
// IsNullString
// ─────────────────────────────────────────────────────────────────────────────

func TestIsNullString(t *testing.T) {
	assert.True(t, IsNullString(""))
	assert.True(t, IsNullString("   "))
	assert.True(t, IsNullString("\t\n"))
	assert.False(t, IsNullString("hello"))
	assert.False(t, IsNullString(" a "))
}

// ─────────────────────────────────────────────────────────────────────────────
// ParseAndPadTimeString
// ─────────────────────────────────────────────────────────────────────────────

func TestParseAndPadTimeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Already fully formatted
		{"00:00:05.0000", "00:00:05.0000"},
		// HH:MM:SS → padded with .0000
		{"0:00:03", "00:00:03.0000"},
		{"1:23:45", "01:23:45.0000"},
		// MM:SS → prefixed with 00:
		{"00:05", "00:00:05.0000"},
		// SS.mmmm
		{"3.5", "00:00:03.5000"},
		{"3.1234", "00:00:03.1234"},
		// Plain seconds
		{"3", "00:00:03.0000"},
		{"12", "00:00:12.0000"},
		// NaN passthrough
		{"NaN", "NaN"},
		// Whitespace trimmed
		{"  5  ", "00:00:05.0000"},
		// HH:MM:SS.mmm (3 fractional digits padded to 4)
		{"0:00:03.500", "00:00:03.5000"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := ParseAndPadTimeString(tc.input)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestParseAndPadTimeString_InvalidInput(t *testing.T) {
	// Known config.NoKey value is "No"; invalid strings return it
	result := ParseAndPadTimeString("not-a-time")
	assert.Equal(t, "No", result)
}

// ─────────────────────────────────────────────────────────────────────────────
// ParseTimeForComparison
// ─────────────────────────────────────────────────────────────────────────────

func TestParseTimeForComparison(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		// NaN → -1
		{"NaN", -1},
		// 0:00:00.0000 → 0
		{"00:00:00.0000", 0},
		// 1 second → 10000 tms
		{"00:00:01.0000", 10000},
		// 1 minute → 600000 tms
		{"00:01:00.0000", 600000},
		// 1 hour → 36000000 tms
		{"01:00:00.0000", 36000000},
		// Mixed: 1:23:45.1000
		{"01:23:45.1000", 36000000 + 23*600000 + 45*10000 + 1000},
		// Subsecond
		{"00:00:05.5000", 55000},
		// Fractional padding (3 digits)
		{"00:00:03.500", 35000},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := ParseTimeForComparison(tc.input)
			assert.Equal(t, tc.expected, got)
		})
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// RemoveDuplicateDictionaries
// ─────────────────────────────────────────────────────────────────────────────

func TestRemoveDuplicateDictionaries(t *testing.T) {
	input := []map[string]string{
		{"name": "Alice", "team": "A"},
		{"name": "Bob", "team": "B"},
		{"name": "Alice", "team": "A2"}, // duplicate name
	}

	result := RemoveDuplicateDictionaries(input, "name")
	require.Len(t, result, 2)
	assert.Equal(t, "Alice", result[0]["name"])
	assert.Equal(t, "Bob", result[1]["name"])
}

func TestRemoveDuplicateDictionaries_MissingKey(t *testing.T) {
	input := []map[string]string{
		{"name": "Alice"},
		{"other": "no-name"}, // missing key → skipped
	}

	result := RemoveDuplicateDictionaries(input, "name")
	require.Len(t, result, 1)
}

// ─────────────────────────────────────────────────────────────────────────────
// AddChugToResultList
// ─────────────────────────────────────────────────────────────────────────────

func TestAddChugToResultList_ValidEntry(t *testing.T) {
	list := []map[string]string{}
	entry := map[string]string{
		"name":       "Alice",
		"discipline": "Bottle",
		"time":       "00:00:03.0000",
		"status":     "Pass",
		"base_time":  "00:00:03.0000",
	}

	result := AddChugToResultList(list, entry)
	require.Len(t, result, 1)
	assert.Equal(t, "Alice", result[0]["name"])
}

func TestAddChugToResultList_MissingRequiredField(t *testing.T) {
	list := []map[string]string{}
	entry := map[string]string{
		"name":       "Alice",
		"discipline": "Bottle",
		// missing "time" and "status"
	}

	result := AddChugToResultList(list, entry)
	assert.Len(t, result, 0, "entry with missing required fields should not be added")
}

func TestAddChugToResultList_NilList(t *testing.T) {
	entry := map[string]string{
		"name":       "Alice",
		"discipline": "Bottle",
		"time":       "00:00:03.0000",
		"status":     "Pass",
	}

	result := AddChugToResultList(nil, entry)
	require.Len(t, result, 1)
}

// ─────────────────────────────────────────────────────────────────────────────
// SortListByKey
// ─────────────────────────────────────────────────────────────────────────────

func TestSortListByKey_Numeric(t *testing.T) {
	list := []map[string]string{
		{"rank": "3"},
		{"rank": "1"},
		{"rank": "2"},
	}

	sorted := SortListByKey(list, "rank")
	assert.Equal(t, "1", sorted[0]["rank"])
	assert.Equal(t, "2", sorted[1]["rank"])
	assert.Equal(t, "3", sorted[2]["rank"])
}

func TestSortListByKey_TimeString(t *testing.T) {
	list := []map[string]string{
		{"time": "00:00:05.0000"},
		{"time": "00:00:03.0000"},
		{"time": "00:00:08.0000"},
	}

	sorted := SortListByKey(list, "time")
	assert.Equal(t, "00:00:03.0000", sorted[0]["time"])
	assert.Equal(t, "00:00:05.0000", sorted[1]["time"])
	assert.Equal(t, "00:00:08.0000", sorted[2]["time"])
}

func TestSortListByKey_Alphabetic(t *testing.T) {
	list := []map[string]string{
		{"name": "Charlie"},
		{"name": "Alice"},
		{"name": "Bob"},
	}

	sorted := SortListByKey(list, "name")
	assert.Equal(t, "Alice", sorted[0]["name"])
	assert.Equal(t, "Bob", sorted[1]["name"])
	assert.Equal(t, "Charlie", sorted[2]["name"])
}

// ─────────────────────────────────────────────────────────────────────────────
// FillListFromJSONFile / SaveListToJSONFile round-trip
// ─────────────────────────────────────────────────────────────────────────────

func TestSaveAndFillJSONFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.json")

	original := []map[string]string{
		{"name": "Alice", "discipline": "Bottle"},
		{"name": "Bob", "discipline": "Half Tankard"},
	}

	err := SaveListToJSONFile(path, original)
	require.NoError(t, err)
	assert.True(t, DoesFileExist(path))

	loaded, err := FillListFromJSONFile(path)
	require.NoError(t, err)
	require.Len(t, loaded, 2)
	assert.Equal(t, "Alice", loaded[0]["name"])
	assert.Equal(t, "Bob", loaded[1]["name"])
}

func TestFillListFromJSONFile_MissingFile(t *testing.T) {
	result, err := FillListFromJSONFile("/no/such/file.json")
	assert.NoError(t, err, "missing file should return empty list, not error")
	assert.Empty(t, result)
}

func TestFillListFromJSONFile_InvalidJSON(t *testing.T) {
	f, err := os.CreateTemp("", "cwtest-*.json")
	require.NoError(t, err)
	f.WriteString("this is not json")
	f.Close()
	defer os.Remove(f.Name())

	_, err = FillListFromJSONFile(f.Name())
	assert.Error(t, err)
}

func TestSaveListToJSONFile_CreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "sub", "file.json")

	err := SaveListToJSONFile(path, []map[string]string{{"key": "val"}})
	require.NoError(t, err)

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	var loaded []map[string]string
	require.NoError(t, json.Unmarshal(data, &loaded))
	assert.Equal(t, "val", loaded[0]["key"])
}

// ─────────────────────────────────────────────────────────────────────────────
// CopyFile
// ─────────────────────────────────────────────────────────────────────────────

func TestCopyFile(t *testing.T) {
	src, err := os.CreateTemp("", "cwsrc-*.txt")
	require.NoError(t, err)
	src.WriteString("hello chugware")
	src.Close()
	defer os.Remove(src.Name())

	dst := filepath.Join(t.TempDir(), "copy.txt")

	err = CopyFile(src.Name(), dst)
	require.NoError(t, err)

	content, err := os.ReadFile(dst)
	require.NoError(t, err)
	assert.Equal(t, "hello chugware", string(content))
}

func TestCopyFile_MissingSource(t *testing.T) {
	err := CopyFile("/no/such/source.txt", filepath.Join(t.TempDir(), "dst.txt"))
	assert.Error(t, err)
}

// ─────────────────────────────────────────────────────────────────────────────
// GetCurrentDateString / GetCurrentTimeString
// ─────────────────────────────────────────────────────────────────────────────

func TestGetCurrentDateString_Format(t *testing.T) {
	s := GetCurrentDateString()
	_, err := time.Parse("2006-01-02", s)
	assert.NoError(t, err, "date string should match YYYY-MM-DD format")
}

func TestGetCurrentTimeString_Format(t *testing.T) {
	s := GetCurrentTimeString()
	_, err := time.Parse("15:04:05", s)
	assert.NoError(t, err, "time string should match HH:MM:SS format")
}

// ─────────────────────────────────────────────────────────────────────────────
// CreateContestFiles
// ─────────────────────────────────────────────────────────────────────────────

func TestCreateContestFiles(t *testing.T) {
	dir := t.TempDir()
	disciplines := []string{"Bottle", "Half Tankard"}

	created, err := CreateContestFiles(disciplines, dir, "TestContest", "2026-02-22")
	require.NoError(t, err)
	assert.Len(t, created, 2)

	for _, filename := range created {
		path := filepath.Join(dir, filename)
		assert.True(t, DoesFileExist(path), "expected created file: "+path)

		// File should contain a valid empty JSON array
		data := []map[string]string{}
		raw, err := os.ReadFile(path)
		require.NoError(t, err)
		require.NoError(t, json.Unmarshal(raw, &data))
		assert.Empty(t, data)
	}
}

func TestCreateContestFiles_SkipsExisting(t *testing.T) {
	dir := t.TempDir()

	// Pre-create one of the files
	existing := filepath.Join(dir, "TestContest_2026-02-22_Bottle.json")
	require.NoError(t, os.WriteFile(existing, []byte(`[{"pre":"existing"}]`), 0644))

	created, err := CreateContestFiles([]string{"Bottle"}, dir, "TestContest", "2026-02-22")
	require.NoError(t, err)
	assert.Len(t, created, 0, "existing file should not be overwritten or re-listed")

	// Verify original content was preserved
	raw, err := os.ReadFile(existing)
	require.NoError(t, err)
	assert.Contains(t, string(raw), "pre")
}

// ─────────────────────────────────────────────────────────────────────────────
// FilterFiles
// ─────────────────────────────────────────────────────────────────────────────

func TestFilterFiles_Found(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "participants_2026.json")
	require.NoError(t, os.WriteFile(target, []byte("[]"), 0644))

	found := FilterFiles(dir, "participants")
	assert.Equal(t, target, found)
}

func TestFilterFiles_NotFound(t *testing.T) {
	dir := t.TempDir()
	result := FilterFiles(dir, "nonexistent")
	assert.Equal(t, "No", result, "should return config.NoKey when not found")
}

func TestFilterFiles_MissingDirectory(t *testing.T) {
	result := FilterFiles("/no/such/dir", "pattern")
	assert.Equal(t, "No", result)
}
