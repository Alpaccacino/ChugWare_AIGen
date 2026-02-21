package data

import (
	"chugware/internal/config"
	"chugware/internal/models"
	"chugware/internal/utils"
	"fmt"
	"strconv"
)

// calcResultTime recomputes result.Time from base_time + additional_time.
// If either value is "NaN" or unparseable the result is "NaN".
func calcResultTime(result *models.Result) {
	// Normalize both inputs through ParseAndPadTimeString so plain values
	// like "3" (seconds) are correctly handled before millisecond parsing.
	normBase := utils.ParseAndPadTimeString(result.BaseTime)
	normAdditional := result.AdditionalTime
	if normAdditional == "" || normAdditional == "0" {
		normAdditional = "00:00:00.0000"
	} else {
		normAdditional = utils.ParseAndPadTimeString(normAdditional)
	}

	baseMs := utils.ParseTimeForComparison(normBase)
	additionalMs := utils.ParseTimeForComparison(normAdditional)

	if baseMs < 0 || additionalMs < 0 {
		result.Time = "NaN"
		return
	}

	totalTms := baseMs + additionalMs
	hours := totalTms / (60 * 60 * 10000)
	minutes := (totalTms / (60 * 10000)) % 60
	seconds := (totalTms / 10000) % 60
	subms := totalTms % 10000
	result.Time = fmt.Sprintf("%02d:%02d:%02d.%04d", hours, minutes, seconds, subms)
}

// ParticipantManager handles participant data operations
type ParticipantManager struct {
	participants []models.Participant
	filePath     string
}

// NewParticipantManager creates a new participant manager
func NewParticipantManager() *ParticipantManager {
	return &ParticipantManager{
		participants: make([]models.Participant, 0),
	}
}

// LoadParticipants loads participants from JSON file
func (pm *ParticipantManager) LoadParticipants(filePath string) error {
	pm.filePath = filePath

	data, err := utils.FillListFromJSONFile(filePath)
	if err != nil {
		return fmt.Errorf("error loading participants: %w", err)
	}

	pm.participants = make([]models.Participant, 0, len(data))
	for _, entry := range data {
		participant := models.Participant{
			Name:        entry["name"],
			Program:     entry["program"],
			Team:        entry["team"],
			Bottle:      entry["bottle"],
			HalfTankard: entry["half_tankard"],
			FullTankard: entry["full_tankard"],
		}
		pm.participants = append(pm.participants, participant)
	}

	return nil
}

// SaveParticipants saves participants to JSON file
func (pm *ParticipantManager) SaveParticipants() error {
	if pm.filePath == "" {
		return fmt.Errorf("no file path set")
	}

	data := make([]map[string]string, 0, len(pm.participants))
	for _, p := range pm.participants {
		entry := map[string]string{
			"name":         p.Name,
			"program":      p.Program,
			"team":         p.Team,
			"bottle":       p.Bottle,
			"half_tankard": p.HalfTankard,
			"full_tankard": p.FullTankard,
		}
		data = append(data, entry)
	}

	return utils.SaveListToJSONFile(pm.filePath, data)
}

// SetFilePath sets the file path without loading from disk
func (pm *ParticipantManager) SetFilePath(filePath string) {
	pm.filePath = filePath
}

// DecrementTries reduces the remaining try count for the given discipline by 1
// (minimum 0). The participant data is updated in memory; call SaveParticipants to persist.
func (pm *ParticipantManager) DecrementTries(name, discipline string) error {
	for i, p := range pm.participants {
		if p.Name == name {
			switch discipline {
			case models.DisciplineBottle:
				tries, _ := strconv.Atoi(p.Bottle)
				if tries > 0 {
					pm.participants[i].Bottle = strconv.Itoa(tries - 1)
				}
			case models.DisciplineHalfTankard:
				tries, _ := strconv.Atoi(p.HalfTankard)
				if tries > 0 {
					pm.participants[i].HalfTankard = strconv.Itoa(tries - 1)
				}
			case models.DisciplineFullTankard:
				tries, _ := strconv.Atoi(p.FullTankard)
				if tries > 0 {
					pm.participants[i].FullTankard = strconv.Itoa(tries - 1)
				}
			default:
				// Other disciplines don't have try counts — silently ignore
				return nil
			}
			return nil
		}
	}
	return fmt.Errorf("participant '%s' not found", name)
}

// AddParticipant adds a new participant
func (pm *ParticipantManager) AddParticipant(participant models.Participant) error {
	// Validate participant data
	if utils.IsNullString(participant.Name) {
		return fmt.Errorf("participant name cannot be empty")
	}

	if !utils.IsStringValid(participant.Name) || !utils.IsStringValid(participant.Program) || !utils.IsStringValid(participant.Team) {
		return fmt.Errorf("participant data exceeds maximum length")
	}

	// Check for duplicates
	for _, existing := range pm.participants {
		if existing.Name == participant.Name {
			return fmt.Errorf("participant with name '%s' already exists", participant.Name)
		}
	}

	pm.participants = append(pm.participants, participant)
	return nil
}

// GetParticipants returns a copy of all participants.
// A copy is returned so callers cannot accidentally mutate the manager's
// internal slice (e.g. via append) and corrupt data like try-counts.
func (pm *ParticipantManager) GetParticipants() []models.Participant {
	cp := make([]models.Participant, len(pm.participants))
	copy(cp, pm.participants)
	return cp
}

// RemoveParticipant removes a participant by name
func (pm *ParticipantManager) RemoveParticipant(name string) error {
	for i, participant := range pm.participants {
		if participant.Name == name {
			pm.participants = append(pm.participants[:i], pm.participants[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("participant '%s' not found", name)
}

// ResultManager handles contest results
type ResultManager struct {
	results  []models.Result
	filePath string
}

// NewResultManager creates a new result manager
func NewResultManager() *ResultManager {
	return &ResultManager{
		results: make([]models.Result, 0),
	}
}

// SetFilePath sets the file path without loading from disk
func (rm *ResultManager) SetFilePath(filePath string) {
	rm.filePath = filePath
}

// LoadResults loads results from JSON file
func (rm *ResultManager) LoadResults(filePath string) error {
	rm.filePath = filePath

	data, err := utils.FillListFromJSONFile(filePath)
	if err != nil {
		return fmt.Errorf("error loading results: %w", err)
	}

	rm.results = make([]models.Result, 0, len(data))
	for _, entry := range data {
		result := models.Result{
			Name:           entry["name"],
			Discipline:     entry["discipline"],
			Time:           entry["time"],
			BaseTime:       entry["base_time"],
			AdditionalTime: entry["additional_time"],
			Status:         entry["status"],
			Comment:        entry["comment"],
		}
		rm.results = append(rm.results, result)
	}

	return nil
}

// SaveResults saves results to JSON file
func (rm *ResultManager) SaveResults() error {
	if rm.filePath == "" {
		return fmt.Errorf("no file path set")
	}

	data := make([]map[string]string, 0, len(rm.results))
	for _, r := range rm.results {
		entry := map[string]string{
			"name":            r.Name,
			"discipline":      r.Discipline,
			"time":            r.Time,
			"base_time":       r.BaseTime,
			"additional_time": r.AdditionalTime,
			"status":          r.Status,
			"comment":         r.Comment,
		}
		data = append(data, entry)
	}

	return utils.SaveListToJSONFile(rm.filePath, data)
}

// AddResult adds a new contest result
func (rm *ResultManager) AddResult(result models.Result) error {
	// Validate result data
	if utils.IsNullString(result.Name) || utils.IsNullString(result.Discipline) {
		return fmt.Errorf("name and discipline are required")
	}

	// Always recalculate time from base + additional (unless status is Disqualified)
	if result.Status != models.StatusDisqualified {
		calcResultTime(&result)
	} else {
		result.Time = "NaN"
	}

	_ = config.NoKey // keep import used
	rm.results = append(rm.results, result)
	return nil
}

// UpdateLastResult updates the last added result for a participant and discipline.
// A Disqualified result is never overwritten.
func (rm *ResultManager) UpdateLastResult(result models.Result) error {
	// Always recalculate time from base + additional
	if result.Status != models.StatusDisqualified {
		calcResultTime(&result)
	} else {
		result.Time = "NaN"
	}

	for i := len(rm.results) - 1; i >= 0; i-- {
		if rm.results[i].Name == result.Name && rm.results[i].Discipline == result.Discipline {
			// Never overwrite an existing Disqualified result
			if rm.results[i].Status == models.StatusDisqualified {
				return fmt.Errorf("%s is already disqualified in %s and cannot be overwritten", result.Name, result.Discipline)
			}
			rm.results[i] = result
			return nil
		}
	}
	return fmt.Errorf("no result found to update for %s in %s", result.Name, result.Discipline)
}

// GetResults returns all results
func (rm *ResultManager) GetResults() []models.Result {
	return rm.results
}

// GetResultsByDiscipline returns results filtered by discipline
func (rm *ResultManager) GetResultsByDiscipline(discipline string) []models.Result {
	var filtered []models.Result
	for _, result := range rm.results {
		if result.Discipline == discipline {
			filtered = append(filtered, result)
		}
	}
	return filtered
}

// GetResultsByParticipant returns results for a specific participant
func (rm *ResultManager) GetResultsByParticipant(name string) []models.Result {
	var filtered []models.Result
	for _, result := range rm.results {
		if result.Name == name {
			filtered = append(filtered, result)
		}
	}
	return filtered
}
