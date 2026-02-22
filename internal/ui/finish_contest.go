package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"chugware/internal/config"
	"chugware/internal/data"
	"chugware/internal/models"
	"chugware/internal/utils"
)

// FinishContest handles contest finalization and results generation
type FinishContest struct {
	app    fyne.App
	window fyne.Window

	// Data managers
	participantMgr *data.ParticipantManager
	resultMgr      *data.ResultManager

	// UI Components - Filtering
	sortFilter *widget.RadioGroup

	// UI Components - Results Display
	bottleList        *widget.List
	halfTankardList   *widget.List
	fullTankardList   *widget.List
	bierStaphetteList *widget.List
	megaMedleyList    *widget.List
	teamClashList     *widget.List
	summaryCard       *widget.Card

	// UI Components - Actions
	generateReportBtn   *widget.Button
	exportResultsBtn    *widget.Button
	generateDiplomasBtn *widget.Button
	refreshBtn          *widget.Button
	saveBtn             *widget.Button

	// UI Components - Summary
	totalParticipantsLabel *widget.Label
	totalResultsLabel      *widget.Label
	passedLabel            *widget.Label
	disqualifiedLabel      *widget.Label
	contestDateLabel       *widget.Label

	// Data
	allResults           []models.Result
	filteredResults      []models.Result
	bottleResults        []models.Result
	halfTankardResults   []models.Result
	fullTankardResults   []models.Result
	bierStaphetteResults []models.Result
	megaMedleyResults    []models.Result
	teamClashResults     []models.Result
	participants         []models.Participant
}

// LeaderboardEntry represents a leaderboard entry
type LeaderboardEntry struct {
	Rank       int
	Name       string
	Team       string
	Discipline string
	Time       string
	Status     string
}

// NewFinishContest creates a new contest finalization window
func NewFinishContest(app fyne.App) *FinishContest {
	fc := &FinishContest{
		app:    app,
		window: app.NewWindow("Finish Contest - Results & Reports"),
	}

	fc.initializeManagers()
	fc.setupUI()
	fc.loadData()
	return fc
}

// initializeManagers sets up data managers
func (fc *FinishContest) initializeManagers() {
	fc.participantMgr = data.NewParticipantManager()
	fc.resultMgr = data.NewResultManager()
}

// setupUI initializes the contest finalization UI
func (fc *FinishContest) setupUI() {
	fc.window.Resize(fyne.NewSize(1600, 1000))
	fc.window.SetFixedSize(false)

	// Initialize components
	fc.createFilterComponents()
	fc.createDisplayComponents()
	fc.createSummaryComponents()
	fc.createActionComponents()

	// Create layout
	content := fc.createLayout()
	fc.window.SetContent(content)
}

// createFilterComponents creates filtering components
func (fc *FinishContest) createFilterComponents() {
	sortOptions := []string{
		"🏆 Fastest First",
		"🐢 Slowest First",
		"⏳ Most Penalty Time",
		"✨ No Penalty (Clean)",
		"🕐 Longest Warm-Up",
		"⚡ Quickest Warm-Up",
		"💀 Hall of Shame",
	}
	fc.sortFilter = widget.NewRadioGroup(sortOptions, fc.onStatusFilterChanged)
}

// createDisplayComponents creates result display components
func (fc *FinishContest) createDisplayComponents() {
	// Bottle results list
	fc.bottleList = widget.NewList(
		func() int { return len(fc.bottleResults) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabelWithStyle("Rank", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle("Name", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle("Team", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle("Time", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle("Status", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id >= 0 && id < len(fc.bottleResults) {
				r := fc.bottleResults[id]
				containers := item.(*fyne.Container)

				containers.Objects[0].(*widget.Label).SetText(fmt.Sprintf("%d", id+1))
				containers.Objects[1].(*widget.Label).SetText(r.Name)
				// Find team mapping from participant list if possible, or leave empty if not in Result
				team := ""
				for _, p := range fc.participants {
					if p.Name == r.Name {
						team = p.Team
						break
					}
				}
				containers.Objects[2].(*widget.Label).SetText(team)
				containers.Objects[3].(*widget.Label).SetText(r.Time)
				containers.Objects[4].(*widget.Label).SetText(r.Status)
			}
		},
	)

	// Half-Tankard results list
	fc.halfTankardList = widget.NewList(
		func() int { return len(fc.halfTankardResults) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabelWithStyle("Rank", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle("Name", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle("Team", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle("Time", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle("Status", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id >= 0 && id < len(fc.halfTankardResults) {
				r := fc.halfTankardResults[id]
				containers := item.(*fyne.Container)

				containers.Objects[0].(*widget.Label).SetText(fmt.Sprintf("%d", id+1))
				containers.Objects[1].(*widget.Label).SetText(r.Name)
				team := ""
				for _, p := range fc.participants {
					if p.Name == r.Name {
						team = p.Team
						break
					}
				}
				containers.Objects[2].(*widget.Label).SetText(team)
				containers.Objects[3].(*widget.Label).SetText(r.Time)
				containers.Objects[4].(*widget.Label).SetText(r.Status)
			}
		},
	)

	// Full-Tankard results list
	fc.fullTankardList = widget.NewList(
		func() int { return len(fc.fullTankardResults) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabelWithStyle("Rank", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle("Name", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle("Team", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle("Time", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle("Status", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id >= 0 && id < len(fc.fullTankardResults) {
				r := fc.fullTankardResults[id]
				containers := item.(*fyne.Container)

				containers.Objects[0].(*widget.Label).SetText(fmt.Sprintf("%d", id+1))
				containers.Objects[1].(*widget.Label).SetText(r.Name)
				team := ""
				for _, p := range fc.participants {
					if p.Name == r.Name {
						team = p.Team
						break
					}
				}
				containers.Objects[2].(*widget.Label).SetText(team)
				containers.Objects[3].(*widget.Label).SetText(r.Time)
				containers.Objects[4].(*widget.Label).SetText(r.Status)
			}
		},
	)

	// Bier Staphette results list
	fc.bierStaphetteList = widget.NewList(
		func() int { return len(fc.bierStaphetteResults) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabelWithStyle("Rank", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle("Name", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle("Team", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle("Time", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle("Status", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id >= 0 && id < len(fc.bierStaphetteResults) {
				r := fc.bierStaphetteResults[id]
				containers := item.(*fyne.Container)

				containers.Objects[0].(*widget.Label).SetText(fmt.Sprintf("%d", id+1))
				containers.Objects[1].(*widget.Label).SetText(r.Name)
				team := ""
				for _, p := range fc.participants {
					if p.Name == r.Name {
						team = p.Team
						break
					}
				}
				containers.Objects[2].(*widget.Label).SetText(team)
				containers.Objects[3].(*widget.Label).SetText(r.Time)
				containers.Objects[4].(*widget.Label).SetText(r.Status)
			}
		},
	)

	// Mega Medley results list
	fc.megaMedleyList = widget.NewList(
		func() int { return len(fc.megaMedleyResults) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabelWithStyle("Rank", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle("Name", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle("Team", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle("Time", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle("Status", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id >= 0 && id < len(fc.megaMedleyResults) {
				r := fc.megaMedleyResults[id]
				containers := item.(*fyne.Container)

				containers.Objects[0].(*widget.Label).SetText(fmt.Sprintf("%d", id+1))
				containers.Objects[1].(*widget.Label).SetText(r.Name)
				team := ""
				for _, p := range fc.participants {
					if p.Name == r.Name {
						team = p.Team
						break
					}
				}
				containers.Objects[2].(*widget.Label).SetText(team)
				containers.Objects[3].(*widget.Label).SetText(r.Time)
				containers.Objects[4].(*widget.Label).SetText(r.Status)
			}
		},
	)

	// Team Clash results list
	fc.teamClashList = widget.NewList(
		func() int { return len(fc.teamClashResults) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabelWithStyle("Rank", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle("Name", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle("Team", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle("Time", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle("Status", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id >= 0 && id < len(fc.teamClashResults) {
				r := fc.teamClashResults[id]
				containers := item.(*fyne.Container)

				containers.Objects[0].(*widget.Label).SetText(fmt.Sprintf("%d", id+1))
				containers.Objects[1].(*widget.Label).SetText(r.Name)
				team := ""
				for _, p := range fc.participants {
					if p.Name == r.Name {
						team = p.Team
						break
					}
				}
				containers.Objects[2].(*widget.Label).SetText(team)
				containers.Objects[3].(*widget.Label).SetText(r.Time)
				containers.Objects[4].(*widget.Label).SetText(r.Status)
			}
		},
	)
}

// createSummaryComponents creates summary display components
func (fc *FinishContest) createSummaryComponents() {
	fc.totalParticipantsLabel = widget.NewLabel("0")
	fc.totalResultsLabel = widget.NewLabel("0")
	fc.passedLabel = widget.NewLabel("0")
	fc.disqualifiedLabel = widget.NewLabel("0")
	fc.contestDateLabel = widget.NewLabel(utils.GetCurrentDateString())

	fc.summaryCard = widget.NewCard("Contest Summary", "Overview of contest results and statistics",
		container.NewVBox(
			widget.NewFormItem("Contest Date", fc.contestDateLabel).Widget,
			widget.NewLabel("(The date when the contest was held)"),

			widget.NewSeparator(),
			widget.NewFormItem("Total Participants", fc.totalParticipantsLabel).Widget,
			widget.NewLabel("(Total number of participants registered for this contest)"),

			widget.NewSeparator(),
			widget.NewFormItem("Total Results", fc.totalResultsLabel).Widget,
			widget.NewLabel("(Number of completed attempts/results recorded)"),

			widget.NewSeparator(),
			widget.NewFormItem("Passed", fc.passedLabel).Widget,
			widget.NewLabel("(Number of participants who successfully completed their event)"),

			widget.NewSeparator(),
			widget.NewFormItem("Disqualified", fc.disqualifiedLabel).Widget,
			widget.NewLabel("(Number of participants disqualified for rule violations)"),
		),
	)
}

// createActionComponents creates action buttons
func (fc *FinishContest) createActionComponents() {
	fc.generateReportBtn = widget.NewButton("Generate Report", fc.generateReport)
	fc.exportResultsBtn = widget.NewButton("Export Results", fc.exportResults)
	fc.generateDiplomasBtn = widget.NewButton("Generate Diplomas", fc.generateDiplomas)
	fc.refreshBtn = widget.NewButton("Refresh Data", fc.refreshData)
	fc.saveBtn = widget.NewButton("Save Final Results", fc.saveFinalResults)
}

// createLayout creates the main layout
func (fc *FinishContest) createLayout() fyne.CanvasObject {
	// Filter section
	filterCard := widget.NewCard("Scoreboard Mode", "",
		container.NewVBox(
			fc.sortFilter,
		),
	)

	// Action buttons
	actionContainer := container.NewVBox(
		fc.generateReportBtn,
		fc.exportResultsBtn,
		fc.generateDiplomasBtn,
		widget.NewSeparator(),
		fc.refreshBtn,
		fc.saveBtn,
	)

	// Left panel with filters, summary, and actions
	leftPanel := container.NewVBox(
		filterCard,
		fc.summaryCard,
		widget.NewCard("Actions", "", actionContainer),
	)

	// Discipline tabs — each list gets the full panel height so it renders correctly
	disciplineTabs := container.NewAppTabs(
		container.NewTabItem("Bottle", fc.bottleList),
		container.NewTabItem("Half-Tankard", fc.halfTankardList),
		container.NewTabItem("Full-Tankard", fc.fullTankardList),
		container.NewTabItem("Bier Staphette", fc.bierStaphetteList),
		container.NewTabItem("Mega Medley", fc.megaMedleyList),
		container.NewTabItem("Team Clash", fc.teamClashList),
	)
	disciplineTabs.SetTabLocation(container.TabLocationTop)

	// Main layout
	mainLayout := container.NewHSplit(leftPanel, disciplineTabs)
	mainLayout.SetOffset(0.20)
	return mainLayout
}

// Event handlers
func (fc *FinishContest) onStatusFilterChanged(selected string) {
	fc.applyFilters()
}

// Data operations
func (fc *FinishContest) loadData() {
	// Load participants
	if config.Settings.ParticipantFile != "" && utils.DoesFileExist(config.Settings.ParticipantFile) {
		if err := fc.participantMgr.LoadParticipants(config.Settings.ParticipantFile); err != nil {
			dialog.ShowError(fmt.Errorf("error loading participants: %w", err), fc.window)
		} else {
			fc.participants = fc.participantMgr.GetParticipants()
		}
	}

	// Load results
	if config.Settings.ResultFile != "" && utils.DoesFileExist(config.Settings.ResultFile) {
		if err := fc.resultMgr.LoadResults(config.Settings.ResultFile); err != nil {
			dialog.ShowError(fmt.Errorf("error loading results: %w", err), fc.window)
		} else {
			fc.allResults = fc.resultMgr.GetResults()
		}
	}

	fc.updateSummary()
	fc.applyFilters()
	fc.populateDisciplineResults()

	// Set initial sort selection after data is loaded
	fc.sortFilter.SetSelected("🏆 Fastest First")
}

func (fc *FinishContest) refreshData() {
	fc.loadData()
	dialog.ShowInformation("Data Refreshed", "Contest data has been reloaded from files", fc.window)
}

func (fc *FinishContest) applyFilters() {
	mode := fc.sortFilter.Selected

	// Filter by status based on mode
	fc.filteredResults = nil
	for _, r := range fc.allResults {
		switch mode {
		case "💀 Hall of Shame":
			if r.Status == models.StatusDisqualified {
				fc.filteredResults = append(fc.filteredResults, r)
			}
		case "✨ No Penalty (Clean)":
			if r.Status == models.StatusPass {
				additionalMs := utils.ParseTimeForComparison(r.AdditionalTime)
				if additionalMs == 0 {
					fc.filteredResults = append(fc.filteredResults, r)
				}
			}
		default:
			fc.filteredResults = append(fc.filteredResults, r)
		}
	}

	fc.populateDisciplineResults()
}

func (fc *FinishContest) populateDisciplineResults() {
	// Clear all discipline-specific results
	fc.bottleResults = nil
	fc.halfTankardResults = nil
	fc.fullTankardResults = nil
	fc.bierStaphetteResults = nil
	fc.megaMedleyResults = nil
	fc.teamClashResults = nil

	// Distribute filtered results to discipline lists
	for _, result := range fc.filteredResults {
		switch result.Discipline {
		case models.DisciplineBottle:
			fc.bottleResults = append(fc.bottleResults, result)
		case models.DisciplineHalfTankard:
			fc.halfTankardResults = append(fc.halfTankardResults, result)
		case models.DisciplineFullTankard:
			fc.fullTankardResults = append(fc.fullTankardResults, result)
		case models.DisciplineBierStaphette:
			fc.bierStaphetteResults = append(fc.bierStaphetteResults, result)
		case models.DisciplineMegaMedley:
			fc.megaMedleyResults = append(fc.megaMedleyResults, result)
		case models.DisciplineTeamClash:
			fc.teamClashResults = append(fc.teamClashResults, result)
		}
	}

	// Sort each discipline's results by current mode
	mode := fc.sortFilter.Selected
	fc.bottleResults = sortByMode(fc.bottleResults, mode)
	fc.halfTankardResults = sortByMode(fc.halfTankardResults, mode)
	fc.fullTankardResults = sortByMode(fc.fullTankardResults, mode)
	fc.bierStaphetteResults = sortByMode(fc.bierStaphetteResults, mode)
	fc.megaMedleyResults = sortByMode(fc.megaMedleyResults, mode)
	fc.teamClashResults = sortByMode(fc.teamClashResults, mode)

	// Refresh all lists
	fc.bottleList.Refresh()
	fc.halfTankardList.Refresh()
	fc.fullTankardList.Refresh()
	fc.bierStaphetteList.Refresh()
	fc.megaMedleyList.Refresh()
	fc.teamClashList.Refresh()
}

func sortByMode(results []models.Result, mode string) []models.Result {
	sorted := make([]models.Result, len(results))
	copy(sorted, results)

	sort.SliceStable(sorted, func(i, j int) bool {
		ri, rj := sorted[i], sorted[j]

		switch mode {
		case "🐢 Slowest First":
			// Pass first, then slowest time first
			if ri.Status == models.StatusPass && rj.Status != models.StatusPass {
				return true
			}
			if ri.Status != models.StatusPass && rj.Status == models.StatusPass {
				return false
			}
			return utils.ParseTimeForComparison(ri.Time) > utils.ParseTimeForComparison(rj.Time)

		case "⏳ Most Penalty Time":
			// Pass first, then most additional time first
			if ri.Status == models.StatusPass && rj.Status != models.StatusPass {
				return true
			}
			if ri.Status != models.StatusPass && rj.Status == models.StatusPass {
				return false
			}
			return utils.ParseTimeForComparison(ri.AdditionalTime) > utils.ParseTimeForComparison(rj.AdditionalTime)

		case "✨ No Penalty (Clean)":
			// Sorted by base time ascending
			return utils.ParseTimeForComparison(ri.BaseTime) < utils.ParseTimeForComparison(rj.BaseTime)

		case "🕐 Longest Warm-Up":
			// Pass first, then longest base time first
			if ri.Status == models.StatusPass && rj.Status != models.StatusPass {
				return true
			}
			if ri.Status != models.StatusPass && rj.Status == models.StatusPass {
				return false
			}
			return utils.ParseTimeForComparison(ri.BaseTime) > utils.ParseTimeForComparison(rj.BaseTime)

		case "⚡ Quickest Warm-Up":
			// Pass first, then shortest base time first
			if ri.Status == models.StatusPass && rj.Status != models.StatusPass {
				return true
			}
			if ri.Status != models.StatusPass && rj.Status == models.StatusPass {
				return false
			}
			return utils.ParseTimeForComparison(ri.BaseTime) < utils.ParseTimeForComparison(rj.BaseTime)

		case "💀 Hall of Shame":
			// All disqualified — sort by name
			return ri.Name < rj.Name

		default: // "🏆 Fastest First"
			// Pass first, then fastest time first
			if ri.Status == models.StatusPass && rj.Status != models.StatusPass {
				return true
			}
			if ri.Status != models.StatusPass && rj.Status == models.StatusPass {
				return false
			}
			return utils.ParseTimeForComparison(ri.Time) < utils.ParseTimeForComparison(rj.Time)
		}
	})
	return sorted
}

func (fc *FinishContest) updateSummary() {
	totalParticipants := len(fc.participants)
	totalResults := len(fc.allResults)

	passed := 0
	disqualified := 0

	for _, result := range fc.allResults {
		switch result.Status {
		case models.StatusPass:
			passed++
		case models.StatusDisqualified:
			disqualified++
		}
	}

	fc.totalParticipantsLabel.SetText(fmt.Sprintf("%d", totalParticipants))
	fc.totalResultsLabel.SetText(fmt.Sprintf("%d", totalResults))
	fc.passedLabel.SetText(fmt.Sprintf("%d", passed))
	fc.disqualifiedLabel.SetText(fmt.Sprintf("%d", disqualified))
}

// Action implementations
func (fc *FinishContest) generateReport() {
	report := fc.createTextReport()

	// Show report in a dialog
	content := widget.NewMultiLineEntry()
	content.SetText(report)
	content.Wrapping = fyne.TextWrapWord

	var reportDialog *widget.PopUp
	reportDialog = widget.NewModalPopUp(
		container.NewBorder(
			// Top: Title and separator
			container.NewVBox(
				widget.NewLabel("Contest Report"),
				widget.NewSeparator(),
			),
			// Bottom: Buttons
			container.NewHBox(
				widget.NewButton("Copy to Clipboard", func() {
					fc.window.Clipboard().SetContent(report)
					dialog.ShowInformation("Copied", "Report copied to clipboard", fc.window)
				}),
				widget.NewButton("Close", func() {
					reportDialog.Hide()
				}),
			),
			nil, nil,
			container.NewScroll(content),
		),
		fc.window.Canvas(),
	)

	reportDialog.Resize(fyne.NewSize(1400, 900))
	reportDialog.Show()
}

func (fc *FinishContest) createTextReport() string {
	var report strings.Builder

	report.WriteString("CHUGWARE2 CONTEST REPORT\n")
	report.WriteString("=" + strings.Repeat("=", 50) + "\n\n")

	report.WriteString(fmt.Sprintf("Generated: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	report.WriteString(fmt.Sprintf("Contest Date: %s\n\n", fc.contestDateLabel.Text))

	report.WriteString("SUMMARY:\n")
	report.WriteString(fmt.Sprintf("Total Participants: %s\n", fc.totalParticipantsLabel.Text))
	report.WriteString(fmt.Sprintf("Total Results: %s\n", fc.totalResultsLabel.Text))
	report.WriteString(fmt.Sprintf("Passed: %s\n", fc.passedLabel.Text))
	report.WriteString(fmt.Sprintf("Disqualified: %s\n\n", fc.disqualifiedLabel.Text))

	// Leaderboard section
	report.WriteString("LEADERBOARD:\n" + strings.Repeat("-", 60) + "\n")

	disciplines := []struct {
		Label   string
		Results []models.Result
	}{
		{"Bottle", fc.bottleResults},
		{"Half-Tankard", fc.halfTankardResults},
		{"Full-Tankard", fc.fullTankardResults},
		{"Bier Staphette", fc.bierStaphetteResults},
		{"Mega Medley", fc.megaMedleyResults},
		{"Team Clash", fc.teamClashResults},
	}

	for _, d := range disciplines {
		if len(d.Results) == 0 {
			continue
		}

		report.WriteString(fmt.Sprintf("\n%s:\n", d.Label))
		rank := 1
		for _, result := range d.Results {
			if result.Status != models.StatusPass {
				continue
			}
			report.WriteString(fmt.Sprintf("%d. %s - %s\n", rank, result.Name, result.Time))
			rank++
		}
	}

	// Detailed results
	report.WriteString("\n\nDETAILED RESULTS:\n" + strings.Repeat("-", 60) + "\n")

	disciplineGroups := make(map[string][]models.Result)
	for _, result := range fc.allResults {
		disciplineGroups[result.Discipline] = append(disciplineGroups[result.Discipline], result)
	}

	for discipline, results := range disciplineGroups {
		report.WriteString(fmt.Sprintf("\n%s:\n", discipline))

		// Sort by status (passed first) then by time
		sort.Slice(results, func(i, j int) bool {
			if results[i].Status == models.StatusPass && results[j].Status != models.StatusPass {
				return true
			}
			if results[i].Status != models.StatusPass && results[j].Status == models.StatusPass {
				return false
			}
			if results[i].Status == models.StatusPass && results[j].Status == models.StatusPass {
				timeI := utils.ParseTimeForComparison(results[i].Time)
				timeJ := utils.ParseTimeForComparison(results[j].Time)
				return timeI < timeJ
			}
			return results[i].Name < results[j].Name
		})

		for _, result := range results {
			status := result.Status
			if result.Comment != "" {
				status += " (" + result.Comment + ")"
			}
			report.WriteString(fmt.Sprintf("  %s - %s - %s\n",
				result.Name, result.Time, status))
		}
	}

	return report.String()
}

func (fc *FinishContest) exportResults() {
	// Save to results directory
	if config.Settings.FolderPathContestNameAndDate != "" {
		resultsPath := config.GetResultsPath()
		timestamp := time.Now().Format("20060102_150405")
		filename := fmt.Sprintf("contest_export_%s.json", timestamp)

		// Convert to JSON format expected by utils
		exportList := make([]map[string]string, 0)

		// Add summary info
		exportList = append(exportList, map[string]string{
			"type":               "summary",
			"date":               fc.contestDateLabel.Text,
			"total_participants": fc.totalParticipantsLabel.Text,
			"total_results":      fc.totalResultsLabel.Text,
			"passed":             fc.passedLabel.Text,
			"disqualified":       fc.disqualifiedLabel.Text,
		})

		// Add results
		for _, result := range fc.allResults {
			exportList = append(exportList, map[string]string{
				"type":            "result",
				"name":            result.Name,
				"discipline":      result.Discipline,
				"time":            result.Time,
				"base_time":       result.BaseTime,
				"additional_time": result.AdditionalTime,
				"status":          result.Status,
				"comment":         result.Comment,
			})
		}

		fullPath := resultsPath + "/" + filename
		if err := utils.SaveListToJSONFile(fullPath, exportList); err != nil {
			dialog.ShowError(fmt.Errorf("error exporting results: %w", err), fc.window)
		} else {
			dialog.ShowInformation("Export Complete",
				fmt.Sprintf("Results exported to: %s", fullPath), fc.window)
		}
	} else {
		dialog.ShowError(fmt.Errorf("no contest directory configured"), fc.window)
	}
}

func (fc *FinishContest) generateDiplomas() {
	// Generate diploma data for winners
	diplomaData := make([]map[string]string, 0)

	disciplines := []struct {
		Label   string
		Results []models.Result
	}{
		{"Bottle", fc.bottleResults},
		{"Half-Tankard", fc.halfTankardResults},
		{"Full-Tankard", fc.fullTankardResults},
		{"Bier Staphette", fc.bierStaphetteResults},
		{"Mega Medley", fc.megaMedleyResults},
		{"Team Clash", fc.teamClashResults},
	}

	for _, d := range disciplines {
		rank := 1
		for _, result := range d.Results {
			if result.Status != models.StatusPass {
				continue
			}
			if rank > 3 {
				break
			}

			place := ""
			switch rank {
			case 1:
				place = "1st Place"
			case 2:
				place = "2nd Place"
			case 3:
				place = "3rd Place"
			}

			// We need Team info but Result struct doesn't have it directly?
			// Checking models.Result struct in types.go, it only had Name, Discipline, Time, BaseTime, Status, Comment.
			// It does NOT have Team.
			// We might need to look up participant info or ensure Result has Team.
			// Re-checking types.go:
			/*
				type Result struct {
					Name           string `json:"name"`
					Discipline     string `json:"discipline"`
					Time           string `json:"time"`
					BaseTime       string `json:"base_time"`
					AdditionalTime string `json:"additional_time"`
					Status         string `json:"status"`
					Comment        string `json:"comment"`
				}
			*/
			// Ah, missing Team. Without Team, we can't fully fill the diploma data properly if it requires Team.
			// However Result struct didn't have Team before either?
			// Wait, the OLD code for generateDiplomas used LeaderboardEntry which HAD Team.
			// LeaderboardEntry struct:
			/*
				type LeaderboardEntry struct {
					Rank       int
					Name       string
					Team       string
					Discipline string
					Time       string
				}
			*/
			// So I need to fetch Team for the result.
			// I have fc.participants map or slice?
			// fc.participants is []models.Participant.
			// I can lookup team by name.

			team := ""
			for _, p := range fc.participants {
				if p.Name == result.Name {
					team = p.Team
					break
				}
			}

			diplomaData = append(diplomaData, map[string]string{
				"name":       result.Name,
				"team":       team,
				"discipline": d.Label,
				"place":      place,
				"time":       result.Time,
				"date":       fc.contestDateLabel.Text,
			})
			rank++
		}
	}

	// Save diploma data
	if config.Settings.FolderPathContestNameAndDate != "" {
		diplomaPath := config.GetDiplomaPath()
		filename := "diploma_data.json"
		fullPath := diplomaPath + "/" + filename

		if err := utils.SaveListToJSONFile(fullPath, diplomaData); err != nil {
			dialog.ShowError(fmt.Errorf("error generating diplomas: %w", err), fc.window)
		} else {
			dialog.ShowInformation("Diplomas Generated",
				fmt.Sprintf("Diploma data saved to: %s\n\nTotal diplomas: %d", fullPath, len(diplomaData)),
				fc.window)
		}
	} else {
		dialog.ShowError(fmt.Errorf("no contest directory configured"), fc.window)
	}
}

// saveFinalResults saves the final results to file
func (fc *FinishContest) saveFinalResults() {
	if err := fc.resultMgr.SaveResults(); err != nil {
		dialog.ShowError(fmt.Errorf("error saving final results: %w", err), fc.window)
		return
	}

	dialog.ShowInformation("Results Saved", "Final contest results have been saved", fc.window)
}

// Show displays the finish contest window
func (fc *FinishContest) Show() {
	fc.window.Show()
}
