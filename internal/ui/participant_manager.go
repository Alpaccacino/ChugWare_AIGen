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

// ParticipantManagerUI handles participant management interface
type ParticipantManagerUI struct {
	app    fyne.App
	window fyne.Window

	// Data managers
	participantMgr *data.ParticipantManager
	resultMgr      *data.ResultManager

	// UI Components - Participant Form
	nameEntry        *widget.Entry
	programEntry     *widget.Entry
	teamEntry        *widget.Entry
	disciplinesEntry *widget.Entry // Single entry for "322" format

	// UI Components - Lists
	participantList    *widget.List
	participantResults *widget.List
	disciplineFilter   *widget.CheckGroup

	// UI Components - Buttons
	addBtn     *widget.Button
	updateBtn  *widget.Button
	deleteBtn  *widget.Button
	refreshBtn *widget.Button
	saveBtn    *widget.Button
	loadBtn    *widget.Button

	// Data
	participants []models.Participant
	results      []models.Result

	// Selection tracking
	selectedParticipantID int

	// Auto-refresh
	refreshTimer *time.Timer
}

// NewParticipantManager creates a new participant manager window
func NewParticipantManager(app fyne.App) *ParticipantManagerUI {
	pm := &ParticipantManagerUI{
		app:                   app,
		window:                app.NewWindow("Participant Management"),
		selectedParticipantID: -1, // No selection initially
	}

	pm.initializeManagers()
	pm.setupUI()
	pm.loadData()
	return pm
}

// initializeManagers sets up data managers
func (pm *ParticipantManagerUI) initializeManagers() {
	pm.participantMgr = data.NewParticipantManager()
	pm.resultMgr = data.NewResultManager()
}

// setupUI initializes the participant manager UI
func (pm *ParticipantManagerUI) setupUI() {
	pm.window.Resize(fyne.NewSize(960, 680))
	pm.window.CenterOnScreen()
	pm.window.SetFixedSize(false)

	// Initialize components
	pm.createFormComponents()
	pm.createListComponents()
	pm.createButtonComponents()

	// Create layout
	content := pm.createLayout()
	pm.window.SetContent(content)

	// Start auto-refresh
	pm.startAutoRefresh()
}

// createFormComponents creates the participant form components
func (pm *ParticipantManagerUI) createFormComponents() {
	pm.nameEntry = widget.NewEntry()
	pm.nameEntry.SetPlaceHolder("Participant name")

	pm.programEntry = widget.NewEntry()
	pm.programEntry.SetPlaceHolder("Program/Course")

	pm.teamEntry = widget.NewEntry()
	pm.teamEntry.SetPlaceHolder("Team name")

	// Single discipline entry like reference project
	pm.disciplinesEntry = widget.NewEntry()
	pm.disciplinesEntry.SetPlaceHolder("322")
	pm.disciplinesEntry.SetText("322") // Default: 3 bottle, 2 half, 2 full
}

// createListComponents creates the list components
func (pm *ParticipantManagerUI) createListComponents() {
	// Participant list
	pm.participantList = widget.NewList(
		func() int {
			return len(pm.participants)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Name"),
				widget.NewLabel("Program"),
				widget.NewLabel("Team"),
				widget.NewLabel("B"),
				widget.NewLabel("H"),
				widget.NewLabel("F"),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id >= 0 && id < len(pm.participants) {
				p := pm.participants[id]
				containers := item.(*fyne.Container)

				containers.Objects[0].(*widget.Label).SetText(p.Name)
				containers.Objects[1].(*widget.Label).SetText(p.Program)
				containers.Objects[2].(*widget.Label).SetText(p.Team)
				containers.Objects[3].(*widget.Label).SetText(p.Bottle)
				containers.Objects[4].(*widget.Label).SetText(p.HalfTankard)
				containers.Objects[5].(*widget.Label).SetText(p.FullTankard)
			}
		},
	)

	pm.participantList.OnSelected = pm.onParticipantSelected

	// Results list
	pm.participantResults = widget.NewList(
		func() int {
			return len(pm.results)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Name"),
				widget.NewLabel("Discipline"),
				widget.NewLabel("Time"),
				widget.NewLabel("Status"),
				widget.NewLabel("Comment"),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id >= 0 && id < len(pm.results) {
				r := pm.results[id]
				containers := item.(*fyne.Container)

				containers.Objects[0].(*widget.Label).SetText(r.Name)
				containers.Objects[1].(*widget.Label).SetText(r.Discipline)
				containers.Objects[2].(*widget.Label).SetText(r.Time)
				containers.Objects[3].(*widget.Label).SetText(r.Status)
				containers.Objects[4].(*widget.Label).SetText(r.Comment)
			}
		},
	)

	// Discipline filter
	disciplines := []string{
		models.DisciplineBottle,
		models.DisciplineHalfTankard,
		models.DisciplineFullTankard,
		models.DisciplineBierStaphette,
		models.DisciplineMegaMedley,
		models.DisciplineTeamClash,
	}

	pm.disciplineFilter = widget.NewCheckGroup(disciplines, pm.onDisciplineFilterChanged)
}

// createButtonComponents creates the action buttons
func (pm *ParticipantManagerUI) createButtonComponents() {
	pm.addBtn = widget.NewButton("Add Participant", pm.addParticipant)
	pm.updateBtn = widget.NewButton("Update", pm.updateParticipant)
	pm.deleteBtn = widget.NewButton("Delete", pm.deleteParticipant)
	pm.refreshBtn = widget.NewButton("Refresh", pm.refreshData)
	pm.saveBtn = widget.NewButton("Save All", pm.saveData)
	pm.loadBtn = widget.NewButton("Load from File", pm.loadFromFile)

	// Initially disable update/delete buttons
	pm.updateBtn.Disable()
	pm.deleteBtn.Disable()
}

// createLayout creates the main layout
func (pm *ParticipantManagerUI) createLayout() fyne.CanvasObject {
	// Participant form
	form := widget.NewCard("Add/Edit Participant", "",
		container.NewVBox(
			widget.NewFormItem("Name", pm.nameEntry).Widget,
			widget.NewFormItem("Program", pm.programEntry).Widget,
			widget.NewFormItem("Team", pm.teamEntry).Widget,

			widget.NewSeparator(),
			widget.NewLabel("Disciplines (3 digits):"),
			pm.disciplinesEntry,

			widget.NewSeparator(),
			container.NewHBox(
				pm.addBtn,
				pm.updateBtn,
				pm.deleteBtn,
			),
		),
	)

	// File status — each file on two lines: label+status, then path
	pFileStatus, pFilePath := getFileStatusParts(config.Settings.ParticipantFile)
	rFileStatus, rFilePath := getFileStatusParts(config.Settings.ResultFile)

	statusCard := widget.NewCard("File Status", "",
		container.NewVBox(
			widget.NewLabel("Participant file: "+pFileStatus),
			widget.NewLabel(pFilePath),
			widget.NewSeparator(),
			widget.NewLabel("Result file: "+rFileStatus),
			widget.NewLabel(rFilePath),
		),
	)

	// Participants tab
	participantsTab := container.NewBorder(
		nil,
		container.NewHBox(pm.refreshBtn, pm.saveBtn, pm.loadBtn),
		nil, nil,
		pm.participantList,
	)

	// Results tab (discipline filter at top)
	resultsTab := container.NewBorder(
		container.NewVBox(
			widget.NewLabel("Filter by discipline:"),
			pm.disciplineFilter,
		),
		nil, nil, nil,
		pm.participantResults,
	)

	// Right panel as AppTabs — mirrors Contest Execution / Finish Contest style
	rightTabs := container.NewAppTabs(
		container.NewTabItem("Participants", participantsTab),
		container.NewTabItem("Participant Results", resultsTab),
	)
	rightTabs.SetTabLocation(container.TabLocationTop)

	// Main layout – left panel wrapped in scroll so it is reachable on small screens
	leftPanel := container.NewVBox(form, statusCard)
	mainLayout := container.NewHSplit(container.NewScroll(leftPanel), rightTabs)
	return mainLayout
}

// Event handlers
func (pm *ParticipantManagerUI) onParticipantSelected(id widget.ListItemID) {
	pm.selectedParticipantID = int(id) // Store the selected ID

	if id >= 0 && id < len(pm.participants) {
		participant := pm.participants[id]

		// Fill form with selected participant data
		pm.nameEntry.SetText(participant.Name)
		pm.programEntry.SetText(participant.Program)
		pm.teamEntry.SetText(participant.Team)
		// Combine discipline attempts into single string
		disciplineString := participant.Bottle + participant.HalfTankard + participant.FullTankard
		pm.disciplinesEntry.SetText(disciplineString)

		// Enable update/delete buttons
		pm.updateBtn.Enable()
		pm.deleteBtn.Enable()

		// Load participant results
		pm.loadParticipantResults(participant.Name)
	}
}

func (pm *ParticipantManagerUI) onDisciplineFilterChanged(selected []string) {
	pm.filterResults(selected)
}

// Data operations
func (pm *ParticipantManagerUI) addParticipant() {
	participant := pm.getParticipantFromForm()
	if participant == nil {
		return
	}

	if err := pm.participantMgr.AddParticipant(*participant); err != nil {
		dialog.ShowError(err, pm.window)
		return
	}

	// Auto-save after adding
	if err := pm.participantMgr.SaveParticipants(); err != nil {
		dialog.ShowError(fmt.Errorf("participant added but failed to save: %w", err), pm.window)
		return
	}

	pm.clearForm()
	pm.refreshParticipantList()
	dialog.ShowInformation("Success", "Participant added and saved successfully", pm.window)
}

func (pm *ParticipantManagerUI) updateParticipant() {
	// Use the stored selected ID from the list selection callback
	if pm.selectedParticipantID < 0 || pm.selectedParticipantID >= len(pm.participants) {
		dialog.ShowError(fmt.Errorf("no participant selected"), pm.window)
		return
	}

	participant := pm.getParticipantFromForm()
	if participant == nil {
		return
	}

	// Remove old and add updated
	oldParticipant := pm.participants[pm.selectedParticipantID]
	if err := pm.participantMgr.RemoveParticipant(oldParticipant.Name); err != nil {
		dialog.ShowError(err, pm.window)
		return
	}

	if err := pm.participantMgr.AddParticipant(*participant); err != nil {
		dialog.ShowError(err, pm.window)
		return
	}

	// Auto-save after updating
	if err := pm.participantMgr.SaveParticipants(); err != nil {
		dialog.ShowError(fmt.Errorf("participant updated but failed to save: %w", err), pm.window)
		return
	}

	pm.clearForm()
	pm.refreshParticipantList()
	dialog.ShowInformation("Success", "Participant updated and saved successfully", pm.window)
}

func (pm *ParticipantManagerUI) deleteParticipant() {
	// Use the stored selected ID from the list selection callback
	if pm.selectedParticipantID < 0 || pm.selectedParticipantID >= len(pm.participants) {
		dialog.ShowError(fmt.Errorf("no participant selected"), pm.window)
		return
	}

	participant := pm.participants[pm.selectedParticipantID]

	dialog.ShowConfirm("Confirm Delete",
		fmt.Sprintf("Are you sure you want to delete participant '%s'?", participant.Name),
		func(confirmed bool) {
			if confirmed {
				if err := pm.participantMgr.RemoveParticipant(participant.Name); err != nil {
					dialog.ShowError(err, pm.window)
					return
				}

				// Auto-save after deleting
				if err := pm.participantMgr.SaveParticipants(); err != nil {
					dialog.ShowError(fmt.Errorf("participant deleted but failed to save: %w", err), pm.window)
					return
				}

				pm.clearForm()
				pm.refreshParticipantList()
				dialog.ShowInformation("Success", "Participant deleted and saved successfully", pm.window)
			}
		}, pm.window)
}

func (pm *ParticipantManagerUI) refreshData() {
	pm.loadData()
}

// loadFromFile shows a file picker to select and load a participant file
func (pm *ParticipantManagerUI) loadFromFile() {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(fmt.Errorf("error selecting file: %w", err), pm.window)
			return
		}
		if reader == nil {
			return // User cancelled
		}
		defer reader.Close()

		filePath := reader.URI().Path()

		// Load participants from selected file
		if err := pm.participantMgr.LoadParticipants(filePath); err != nil {
			dialog.ShowError(fmt.Errorf("error loading participants from %s: %w", filePath, err), pm.window)
			return
		}

		// Update configuration to remember this file
		config.Settings.ParticipantFile = filePath
		if err := config.SaveConfig(); err != nil {
			dialog.ShowError(fmt.Errorf("loaded file but failed to save configuration: %w", err), pm.window)
		}

		pm.refreshParticipantList()
		dialog.ShowInformation("Success",
			fmt.Sprintf("Loaded %d participants from: %s", len(pm.participants), filePath),
			pm.window)

		// Also try to load results if result file is configured
		if config.Settings.ResultFile != "" && utils.DoesFileExist(config.Settings.ResultFile) {
			if err := pm.resultMgr.LoadResults(config.Settings.ResultFile); err == nil {
				pm.refreshResultsList()
			}
		}
	}, pm.window)
}

func (pm *ParticipantManagerUI) saveData() {
	if err := pm.participantMgr.SaveParticipants(); err != nil {
		dialog.ShowError(fmt.Errorf("error saving participants: %w", err), pm.window)
		return
	}

	dialog.ShowInformation("Success", "Data saved successfully", pm.window)
}

func (pm *ParticipantManagerUI) loadData() {
	// Check if participant file is configured
	if config.Settings.ParticipantFile == "" {
		dialog.ShowInformation("No Participant File",
			"No participant file is configured. Please use Contest Wizard to create a contest or set the file path in Configuration.",
			pm.window)
		// Clear the current list and refresh to show empty state
		pm.participants = nil
		pm.refreshParticipantList()
		return
	}

	// Load participants
	if utils.DoesFileExist(config.Settings.ParticipantFile) {
		if err := pm.participantMgr.LoadParticipants(config.Settings.ParticipantFile); err != nil {
			dialog.ShowError(fmt.Errorf("error loading participants from %s: %w", config.Settings.ParticipantFile, err), pm.window)
		} else {
			pm.refreshParticipantList()
		}
	} else {
		// File doesn't exist - create empty file and inform user
		if err := utils.SaveListToJSONFile(config.Settings.ParticipantFile, []map[string]string{}); err != nil {
			dialog.ShowError(fmt.Errorf("error creating participant file %s: %w", config.Settings.ParticipantFile, err), pm.window)
		} else {
			dialog.ShowInformation("Created New File",
				fmt.Sprintf("Participant file didn't exist. Created new empty file: %s", config.Settings.ParticipantFile),
				pm.window)
		}
		// Clear the current list and refresh to show empty state
		pm.participants = nil
		pm.refreshParticipantList()
	}

	// Load results
	if config.Settings.ResultFile != "" && utils.DoesFileExist(config.Settings.ResultFile) {
		if err := pm.resultMgr.LoadResults(config.Settings.ResultFile); err != nil {
			dialog.ShowError(fmt.Errorf("error loading results: %w", err), pm.window)
		} else {
			pm.refreshResultsList()
		}
	} else if config.Settings.ResultFile != "" {
		// Create empty results file if it doesn't exist
		if err := utils.SaveListToJSONFile(config.Settings.ResultFile, []map[string]string{}); err != nil {
			dialog.ShowError(fmt.Errorf("error creating result file %s: %w", config.Settings.ResultFile, err), pm.window)
		} else {
			// Load the newly-created file so filePath is set on resultMgr
			if err := pm.resultMgr.LoadResults(config.Settings.ResultFile); err != nil {
				dialog.ShowError(fmt.Errorf("error initialising result manager: %w", err), pm.window)
			}
		}
		pm.results = nil
		pm.refreshResultsList()
	}
}

// Helper methods
func (pm *ParticipantManagerUI) getParticipantFromForm() *models.Participant {
	name := strings.TrimSpace(pm.nameEntry.Text)
	if name == "" {
		dialog.ShowError(fmt.Errorf("participant name is required"), pm.window)
		return nil
	}

	// Parse discipline string (e.g., "322" -> 3,2,2)
	disciplineStr := strings.TrimSpace(pm.disciplinesEntry.Text)
	if len(disciplineStr) != 3 {
		dialog.ShowError(fmt.Errorf("disciplines must be 3 digits (e.g., 322)"), pm.window)
		return nil
	}

	program := strings.TrimSpace(pm.programEntry.Text)
	if program == "" {
		program = "N/A"
	}
	team := strings.TrimSpace(pm.teamEntry.Text)
	if team == "" {
		team = "N/A"
	}

	participant := &models.Participant{
		Name:        name,
		Program:     program,
		Team:        team,
		Bottle:      string(disciplineStr[0]),
		HalfTankard: string(disciplineStr[1]),
		FullTankard: string(disciplineStr[2]),
	}

	return participant
}

func (pm *ParticipantManagerUI) clearForm() {
	pm.nameEntry.SetText("")
	pm.programEntry.SetText("")
	pm.teamEntry.SetText("")
	pm.disciplinesEntry.SetText("322") // Reset to default
}

func (pm *ParticipantManagerUI) refreshParticipantList() {
	pm.participants = pm.participantMgr.GetParticipants()

	// Sort by name
	sort.Slice(pm.participants, func(i, j int) bool {
		return pm.participants[i].Name < pm.participants[j].Name
	})

	pm.participantList.Refresh()
}

func (pm *ParticipantManagerUI) refreshResultsList() {
	pm.results = pm.resultMgr.GetResults()
	pm.participantResults.Refresh()
}

func (pm *ParticipantManagerUI) loadParticipantResults(participantName string) {
	filteredResults := pm.resultMgr.GetResultsByParticipant(participantName)

	// Apply discipline filter if any
	if len(pm.disciplineFilter.Selected) > 0 {
		var finalResults []models.Result
		for _, result := range filteredResults {
			for _, discipline := range pm.disciplineFilter.Selected {
				if result.Discipline == discipline {
					finalResults = append(finalResults, result)
					break
				}
			}
		}
		pm.results = finalResults
	} else {
		pm.results = filteredResults
	}

	pm.participantResults.Refresh()
}

func (pm *ParticipantManagerUI) filterResults(disciplines []string) {
	if len(disciplines) == 0 {
		pm.results = pm.resultMgr.GetResults()
	} else {
		var filteredResults []models.Result
		allResults := pm.resultMgr.GetResults()

		for _, result := range allResults {
			for _, discipline := range disciplines {
				if result.Discipline == discipline && result.Status == models.StatusPass {
					filteredResults = append(filteredResults, result)
					break
				}
			}
		}

		pm.results = filteredResults
	}

	pm.participantResults.Refresh()
}

func (pm *ParticipantManagerUI) startAutoRefresh() {
	// Refresh every 10 seconds
	pm.refreshTimer = time.AfterFunc(10*time.Second, func() {
		pm.refreshData()
		pm.startAutoRefresh() // Restart timer
	})
}

// getFileStatusParts returns a short status word and the file path as separate strings.
func getFileStatusParts(filePath string) (status, path string) {
	if filePath == "" {
		return "Not set", ""
	}
	if utils.DoesFileExist(filePath) {
		return "OK", filePath
	}
	return "Missing", filePath
}

// getFileStatus is kept for any other callers.
func getFileStatus(filePath string) string {
	if filePath == "" {
		return "Not set"
	}
	if utils.DoesFileExist(filePath) {
		return fmt.Sprintf("OK (%s)", filePath)
	}
	return fmt.Sprintf("Missing (%s)", filePath)
}

// Show displays the participant manager window
func (pm *ParticipantManagerUI) Show() {
	pm.window.Show()
}
