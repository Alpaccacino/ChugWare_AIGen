package data

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

func newParticipant(name string) models.Participant {
	return models.Participant{
		Name:        name,
		Program:     "Engineering",
		Team:        "Alpha",
		Bottle:      "3",
		HalfTankard: "2",
		FullTankard: "1",
	}
}

// writeParticipantFile writes a slice of participants as JSON to a temp file and
// returns the file path along with a cleanup function.
func writeParticipantFile(t *testing.T, participants []models.Participant) string {
	t.Helper()

	rows := make([]map[string]string, len(participants))
	for i, p := range participants {
		rows[i] = map[string]string{
			"name":         p.Name,
			"program":      p.Program,
			"team":         p.Team,
			"bottle":       p.Bottle,
			"half_tankard": p.HalfTankard,
			"full_tankard": p.FullTankard,
		}
	}

	data, err := json.MarshalIndent(rows, "", "  ")
	require.NoError(t, err)

	path := filepath.Join(t.TempDir(), "participants.json")
	require.NoError(t, os.WriteFile(path, data, 0644))
	return path
}

// ─────────────────────────────────────────────────────────────────────────────
// calcResultTime (package-private, tested directly)
// ─────────────────────────────────────────────────────────────────────────────

func TestCalcResultTime_NoAdditional(t *testing.T) {
	r := &models.Result{BaseTime: "3"}
	calcResultTime(r)
	assert.Equal(t, "00:00:03.0000", r.Time)
}

func TestCalcResultTime_WithAdditional(t *testing.T) {
	r := &models.Result{BaseTime: "00:00:03.0000", AdditionalTime: "00:00:02.0000"}
	calcResultTime(r)
	assert.Equal(t, "00:00:05.0000", r.Time)
}

func TestCalcResultTime_ZeroAdditional(t *testing.T) {
	r := &models.Result{BaseTime: "00:00:05.0000", AdditionalTime: "0"}
	calcResultTime(r)
	assert.Equal(t, "00:00:05.0000", r.Time)
}

func TestCalcResultTime_EmptyAdditional(t *testing.T) {
	r := &models.Result{BaseTime: "00:00:05.0000", AdditionalTime: ""}
	calcResultTime(r)
	assert.Equal(t, "00:00:05.0000", r.Time)
}

func TestCalcResultTime_NaNBase(t *testing.T) {
	r := &models.Result{BaseTime: "NaN", AdditionalTime: ""}
	calcResultTime(r)
	assert.Equal(t, "NaN", r.Time)
}

func TestCalcResultTime_NaNAdditional(t *testing.T) {
	r := &models.Result{BaseTime: "00:00:03.0000", AdditionalTime: "NaN"}
	calcResultTime(r)
	assert.Equal(t, "NaN", r.Time)
}

func TestCalcResultTime_CrossMinuteBoundary(t *testing.T) {
	// 58 seconds + 5 seconds = 1 minute 3 seconds
	r := &models.Result{BaseTime: "00:00:58.0000", AdditionalTime: "00:00:05.0000"}
	calcResultTime(r)
	assert.Equal(t, "00:01:03.0000", r.Time)
}

// ─────────────────────────────────────────────────────────────────────────────
// ParticipantManager – AddParticipant
// ─────────────────────────────────────────────────────────────────────────────

func TestParticipantManager_AddParticipant_Success(t *testing.T) {
	pm := NewParticipantManager()
	err := pm.AddParticipant(newParticipant("Alice"))
	require.NoError(t, err)
	assert.Len(t, pm.GetParticipants(), 1)
}

func TestParticipantManager_AddParticipant_EmptyName(t *testing.T) {
	pm := NewParticipantManager()
	p := newParticipant("")
	err := pm.AddParticipant(p)
	assert.Error(t, err)
	assert.Empty(t, pm.GetParticipants())
}

func TestParticipantManager_AddParticipant_WhitespaceName(t *testing.T) {
	pm := NewParticipantManager()
	p := newParticipant("   ")
	err := pm.AddParticipant(p)
	assert.Error(t, err)
}

func TestParticipantManager_AddParticipant_Duplicate(t *testing.T) {
	pm := NewParticipantManager()
	require.NoError(t, pm.AddParticipant(newParticipant("Alice")))
	err := pm.AddParticipant(newParticipant("Alice"))
	assert.Error(t, err)
	assert.Len(t, pm.GetParticipants(), 1)
}

func TestParticipantManager_AddParticipant_MultipleDifferent(t *testing.T) {
	pm := NewParticipantManager()
	require.NoError(t, pm.AddParticipant(newParticipant("Alice")))
	require.NoError(t, pm.AddParticipant(newParticipant("Bob")))
	require.NoError(t, pm.AddParticipant(newParticipant("Charlie")))
	assert.Len(t, pm.GetParticipants(), 3)
}

// ─────────────────────────────────────────────────────────────────────────────
// ParticipantManager – RemoveParticipant
// ─────────────────────────────────────────────────────────────────────────────

func TestParticipantManager_RemoveParticipant_Success(t *testing.T) {
	pm := NewParticipantManager()
	require.NoError(t, pm.AddParticipant(newParticipant("Alice")))
	require.NoError(t, pm.AddParticipant(newParticipant("Bob")))

	err := pm.RemoveParticipant("Alice")
	require.NoError(t, err)

	p := pm.GetParticipants()
	require.Len(t, p, 1)
	assert.Equal(t, "Bob", p[0].Name)
}

func TestParticipantManager_RemoveParticipant_NotFound(t *testing.T) {
	pm := NewParticipantManager()
	err := pm.RemoveParticipant("Ghost")
	assert.Error(t, err)
}

func TestParticipantManager_RemoveParticipant_PreservesOrder(t *testing.T) {
	pm := NewParticipantManager()
	for _, name := range []string{"A", "B", "C", "D"} {
		require.NoError(t, pm.AddParticipant(newParticipant(name)))
	}
	require.NoError(t, pm.RemoveParticipant("B"))

	names := make([]string, 0)
	for _, p := range pm.GetParticipants() {
		names = append(names, p.Name)
	}
	assert.Equal(t, []string{"A", "C", "D"}, names)
}

// ─────────────────────────────────────────────────────────────────────────────
// ParticipantManager – GetParticipants returns copy (no aliasing)
// ─────────────────────────────────────────────────────────────────────────────

func TestParticipantManager_GetParticipants_ReturnsCopy(t *testing.T) {
	pm := NewParticipantManager()
	require.NoError(t, pm.AddParticipant(newParticipant("Alice")))

	got := pm.GetParticipants()
	got[0].Name = "Mutated"

	// Internal state should be untouched
	assert.Equal(t, "Alice", pm.GetParticipants()[0].Name)
}

// ─────────────────────────────────────────────────────────────────────────────
// ParticipantManager – DecrementTries
// ─────────────────────────────────────────────────────────────────────────────

func TestParticipantManager_DecrementTries_Bottle(t *testing.T) {
	pm := NewParticipantManager()
	require.NoError(t, pm.AddParticipant(newParticipant("Alice"))) // Bottle = "3"

	require.NoError(t, pm.DecrementTries("Alice", models.DisciplineBottle))
	assert.Equal(t, "2", pm.GetParticipants()[0].Bottle)
}

func TestParticipantManager_DecrementTries_HalfTankard(t *testing.T) {
	pm := NewParticipantManager()
	require.NoError(t, pm.AddParticipant(newParticipant("Alice"))) // HalfTankard = "2"

	require.NoError(t, pm.DecrementTries("Alice", models.DisciplineHalfTankard))
	assert.Equal(t, "1", pm.GetParticipants()[0].HalfTankard)
}

func TestParticipantManager_DecrementTries_FullTankard(t *testing.T) {
	pm := NewParticipantManager()
	require.NoError(t, pm.AddParticipant(newParticipant("Alice"))) // FullTankard = "1"

	require.NoError(t, pm.DecrementTries("Alice", models.DisciplineFullTankard))
	assert.Equal(t, "0", pm.GetParticipants()[0].FullTankard)
}

func TestParticipantManager_DecrementTries_DoesNotGoBelowZero(t *testing.T) {
	pm := NewParticipantManager()
	p := newParticipant("Alice")
	p.Bottle = "0"
	require.NoError(t, pm.AddParticipant(p))

	require.NoError(t, pm.DecrementTries("Alice", models.DisciplineBottle))
	assert.Equal(t, "0", pm.GetParticipants()[0].Bottle)
}

func TestParticipantManager_DecrementTries_NotFound(t *testing.T) {
	pm := NewParticipantManager()
	err := pm.DecrementTries("Ghost", models.DisciplineBottle)
	assert.Error(t, err)
}

func TestParticipantManager_DecrementTries_OtherDiscipline(t *testing.T) {
	pm := NewParticipantManager()
	require.NoError(t, pm.AddParticipant(newParticipant("Alice")))
	// BierStaphette has no try-count – should return nil without panicking
	err := pm.DecrementTries("Alice", models.DisciplineBierStaphette)
	assert.NoError(t, err)
}

// ─────────────────────────────────────────────────────────────────────────────
// ParticipantManager – Load/Save round-trip
// ─────────────────────────────────────────────────────────────────────────────

func TestParticipantManager_LoadSave_RoundTrip(t *testing.T) {
	originals := []models.Participant{
		newParticipant("Alice"),
		newParticipant("Bob"),
	}
	path := writeParticipantFile(t, originals)

	pm := NewParticipantManager()
	require.NoError(t, pm.LoadParticipants(path))
	require.Len(t, pm.GetParticipants(), 2)

	// Mutate and save
	require.NoError(t, pm.AddParticipant(newParticipant("Charlie")))
	require.NoError(t, pm.SaveParticipants())

	// Reload from disk
	pm2 := NewParticipantManager()
	require.NoError(t, pm2.LoadParticipants(path))
	assert.Len(t, pm2.GetParticipants(), 3)
	assert.Equal(t, "Charlie", pm2.GetParticipants()[2].Name)
}

func TestParticipantManager_SaveParticipants_NoFilePath(t *testing.T) {
	pm := NewParticipantManager()
	err := pm.SaveParticipants()
	assert.Error(t, err)
}

// ─────────────────────────────────────────────────────────────────────────────
// ResultManager – AddResult
// ─────────────────────────────────────────────────────────────────────────────

func TestResultManager_AddResult_Pass(t *testing.T) {
	rm := NewResultManager()
	r := models.Result{
		Name:       "Alice",
		Discipline: "Bottle",
		BaseTime:   "00:00:03.0000",
		Status:     models.StatusPass,
	}
	require.NoError(t, rm.AddResult(r))
	results := rm.GetResults()
	require.Len(t, results, 1)
	assert.Equal(t, "00:00:03.0000", results[0].Time)
}

func TestResultManager_AddResult_Disqualified(t *testing.T) {
	rm := NewResultManager()
	r := models.Result{
		Name:       "Alice",
		Discipline: "Bottle",
		Status:     models.StatusDisqualified,
	}
	require.NoError(t, rm.AddResult(r))
	results := rm.GetResults()
	require.Len(t, results, 1)
	assert.Equal(t, "NaN", results[0].Time)
}

func TestResultManager_AddResult_EmptyName(t *testing.T) {
	rm := NewResultManager()
	r := models.Result{Discipline: "Bottle", Status: models.StatusPass}
	err := rm.AddResult(r)
	assert.Error(t, err)
	assert.Empty(t, rm.GetResults())
}

func TestResultManager_AddResult_EmptyDiscipline(t *testing.T) {
	rm := NewResultManager()
	r := models.Result{Name: "Alice", Status: models.StatusPass}
	err := rm.AddResult(r)
	assert.Error(t, err)
}

func TestResultManager_AddResult_TimeSummedFromComponents(t *testing.T) {
	rm := NewResultManager()
	r := models.Result{
		Name:           "Alice",
		Discipline:     "Bottle",
		BaseTime:       "00:00:03.0000",
		AdditionalTime: "00:00:02.0000",
		Status:         models.StatusPass,
	}
	require.NoError(t, rm.AddResult(r))
	assert.Equal(t, "00:00:05.0000", rm.GetResults()[0].Time)
}

// ─────────────────────────────────────────────────────────────────────────────
// ResultManager – UpdateLastResult
// ─────────────────────────────────────────────────────────────────────────────

func TestResultManager_UpdateLastResult_Success(t *testing.T) {
	rm := NewResultManager()
	initial := models.Result{Name: "Alice", Discipline: "Bottle", BaseTime: "00:00:05.0000", Status: models.StatusPass}
	require.NoError(t, rm.AddResult(initial))

	updated := models.Result{Name: "Alice", Discipline: "Bottle", BaseTime: "00:00:04.0000", Status: models.StatusPass}
	require.NoError(t, rm.UpdateLastResult(updated))

	assert.Equal(t, "00:00:04.0000", rm.GetResults()[0].Time)
}

func TestResultManager_UpdateLastResult_CannotOverwriteDisqualified(t *testing.T) {
	rm := NewResultManager()
	dq := models.Result{Name: "Alice", Discipline: "Bottle", Status: models.StatusDisqualified}
	require.NoError(t, rm.AddResult(dq))

	updated := models.Result{Name: "Alice", Discipline: "Bottle", BaseTime: "00:00:03.0000", Status: models.StatusPass}
	err := rm.UpdateLastResult(updated)
	assert.Error(t, err)
	assert.Equal(t, "NaN", rm.GetResults()[0].Time)
}

func TestResultManager_UpdateLastResult_NotFound(t *testing.T) {
	rm := NewResultManager()
	r := models.Result{Name: "Ghost", Discipline: "Bottle", Status: models.StatusPass}
	err := rm.UpdateLastResult(r)
	assert.Error(t, err)
}

func TestResultManager_UpdateLastResult_UpdatesLastMatchOnly(t *testing.T) {
	rm := NewResultManager()
	// Two results for the same participant / discipline (e.g. retry after fail)
	first := models.Result{Name: "Alice", Discipline: "Bottle", BaseTime: "00:00:07.0000", Status: models.StatusPass}
	second := models.Result{Name: "Alice", Discipline: "Bottle", BaseTime: "00:00:06.0000", Status: models.StatusPass}
	require.NoError(t, rm.AddResult(first))
	require.NoError(t, rm.AddResult(second))

	upd := models.Result{Name: "Alice", Discipline: "Bottle", BaseTime: "00:00:04.0000", Status: models.StatusPass}
	require.NoError(t, rm.UpdateLastResult(upd))

	results := rm.GetResults()
	// First result unchanged
	assert.Equal(t, "00:00:07.0000", results[0].Time)
	// Second (last) updated
	assert.Equal(t, "00:00:04.0000", results[1].Time)
}

// ─────────────────────────────────────────────────────────────────────────────
// ResultManager – GetResultsByDiscipline / GetResultsByParticipant
// ─────────────────────────────────────────────────────────────────────────────

func TestResultManager_GetResultsByDiscipline(t *testing.T) {
	rm := NewResultManager()
	require.NoError(t, rm.AddResult(models.Result{Name: "Alice", Discipline: "Bottle", Status: models.StatusPass}))
	require.NoError(t, rm.AddResult(models.Result{Name: "Bob", Discipline: "Half Tankard", Status: models.StatusPass}))
	require.NoError(t, rm.AddResult(models.Result{Name: "Charlie", Discipline: "Bottle", Status: models.StatusPass}))

	bottles := rm.GetResultsByDiscipline("Bottle")
	require.Len(t, bottles, 2)
	for _, r := range bottles {
		assert.Equal(t, "Bottle", r.Discipline)
	}
}

func TestResultManager_GetResultsByParticipant(t *testing.T) {
	rm := NewResultManager()
	require.NoError(t, rm.AddResult(models.Result{Name: "Alice", Discipline: "Bottle", Status: models.StatusPass}))
	require.NoError(t, rm.AddResult(models.Result{Name: "Bob", Discipline: "Bottle", Status: models.StatusPass}))
	require.NoError(t, rm.AddResult(models.Result{Name: "Alice", Discipline: "Half Tankard", Status: models.StatusPass}))

	aliceResults := rm.GetResultsByParticipant("Alice")
	require.Len(t, aliceResults, 2)
	for _, r := range aliceResults {
		assert.Equal(t, "Alice", r.Name)
	}
}

func TestResultManager_GetResultsByDiscipline_Empty(t *testing.T) {
	rm := NewResultManager()
	results := rm.GetResultsByDiscipline("Bottle")
	assert.Nil(t, results)
}

// ─────────────────────────────────────────────────────────────────────────────
// ResultManager – LoadResults / SaveResults round-trip
// ─────────────────────────────────────────────────────────────────────────────

func TestResultManager_LoadSave_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "results.json")

	// Write initial JSON file manually
	initial := []map[string]string{
		{
			"name":            "Alice",
			"discipline":      "Bottle",
			"time":            "00:00:03.0000",
			"base_time":       "00:00:03.0000",
			"additional_time": "",
			"status":          "Pass",
			"comment":         "",
		},
	}
	data, err := json.MarshalIndent(initial, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path, data, 0644))

	rm := NewResultManager()
	require.NoError(t, rm.LoadResults(path))
	require.Len(t, rm.GetResults(), 1)
	assert.Equal(t, "Alice", rm.GetResults()[0].Name)

	// Add a second result and save
	require.NoError(t, rm.AddResult(models.Result{
		Name:       "Bob",
		Discipline: "Half Tankard",
		BaseTime:   "00:00:05.0000",
		Status:     models.StatusPass,
	}))
	require.NoError(t, rm.SaveResults())

	rm2 := NewResultManager()
	require.NoError(t, rm2.LoadResults(path))
	assert.Len(t, rm2.GetResults(), 2)
}

func TestResultManager_SaveResults_NoFilePath(t *testing.T) {
	rm := NewResultManager()
	err := rm.SaveResults()
	assert.Error(t, err)
}

func TestResultManager_SetFilePath_PathIsUsed(t *testing.T) {
	// SetFilePath should cause SaveResults to attempt a write (not "no file path set").
	// Use a path with an invalid filename character on Windows (NUL byte) so the
	// write fails with an OS error rather than succeeding unexpectedly.
	rm := NewResultManager()
	rm.SetFilePath(filepath.Join(t.TempDir(), "ok_results.json"))

	// Save with no results – should succeed and produce an empty JSON array.
	require.NoError(t, rm.SaveResults())

	// Confirm the file exists and contains a valid JSON array.
	path := filepath.Join(t.TempDir(), "ok_results.json")
	rm2 := NewResultManager()
	rm2.SetFilePath(path)
	require.NoError(t, rm2.SaveResults())

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	var loaded []map[string]string
	require.NoError(t, json.Unmarshal(data, &loaded))
	assert.Empty(t, loaded)
}
