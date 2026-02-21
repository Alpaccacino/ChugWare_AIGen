package models

import "time"

// Participant represents a contest participant
type Participant struct {
	Name        string `json:"name"`
	Program     string `json:"program"`
	Team        string `json:"team"`
	Bottle      string `json:"bottle"`
	HalfTankard string `json:"half_tankard"`
	FullTankard string `json:"full_tankard"`
}

// Result represents a contest result for a participant
type Result struct {
	Name           string `json:"name"`
	Discipline     string `json:"discipline"`
	Time           string `json:"time"`
	BaseTime       string `json:"base_time"`
	AdditionalTime string `json:"additional_time"`
	Status         string `json:"status"`
	Comment        string `json:"comment"`
}

// Contest represents a contest configuration
type Contest struct {
	Name        string    `json:"name"`
	Date        time.Time `json:"date"`
	IsOfficial  bool      `json:"is_official"`
	Disciplines []string  `json:"disciplines"`
	FolderPath  string    `json:"folder_path"`
}

// ContestSettings holds contest configuration
type ContestSettings struct {
	FolderPath                   string `json:"folder_path"`
	FolderPathContestNameAndDate string `json:"folder_path_contest_name_and_date"`
	ParticipantFile              string `json:"participant_file"`
	ResultFile                   string `json:"result_file"`
	TemplateFile                 string `json:"template_file"`
	BottleFile                   string `json:"bottle_file"`
	HalfTankardFile              string `json:"half_tankard_file"`
	FullTankardFile              string `json:"full_tankard_file"`
	BierStaphetteFile            string `json:"bier_staphette_file"`
	MegaMedleyFile               string `json:"mega_medley_file"`
	TeamClashFile                string `json:"team_clash_file"`
	FactionFile                  string `json:"faction_file"`
	LeaderBoardFile              string `json:"leader_board_file"`
	DiplomasFile                 string `json:"diplomas_file"`

	// External clock / serial device
	ExternalClockPort string `json:"external_clock_port"`
	ExternalClockBaud int    `json:"external_clock_baud"`
}

// Discipline types
const (
	DisciplineBottle        = "Bottle"
	DisciplineHalfTankard   = "Half Tankard"
	DisciplineFullTankard   = "Full Tankard"
	DisciplineBierStaphette = "Bier Staphette"
	DisciplineMegaMedley    = "Mega Medley"
	DisciplineTeamClash     = "Team Clash"
)

// Status types
const (
	StatusPass         = "Pass"
	StatusDisqualified = "Disqualified"
	StatusFail         = "Fail"
)

// Timer state
type TimerState struct {
	Running    bool
	StartTime  time.Time
	Duration   time.Duration
	Paused     bool
	PausedTime time.Duration
}
