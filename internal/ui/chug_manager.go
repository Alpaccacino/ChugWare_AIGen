package ui

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"chugware/internal/config"
	"chugware/internal/data"
	"chugware/internal/models"
	"chugware/internal/utils"
)

// ChugManager handles the contest execution and timing
type ChugManager struct {
	app    fyne.App
	window fyne.Window

	// Data managers
	participantMgr *data.ParticipantManager
	resultMgr      *data.ResultManager

	// Timer state
	timerState    models.TimerState
	timerWidget   *canvas.Text
	ticker        *time.Ticker
	stopTimerChan chan bool

	// UI Components - Contest Selection
	disciplineSelect  *widget.Select
	timePerEventEntry *widget.Entry

	// UI Components - Participant Lists
	availableList       *widget.List
	allParticipantsList *widget.List

	// UI Components - Current Contest
	currentChuggerCard  *widget.Card
	currentNameLabel    *widget.Label
	currentProgramLabel *widget.Label
	currentTeamLabel    *widget.Label
	currentTriesLabel   *widget.Label

	// UI Components - Timer Controls
	readyCheckBtn *widget.Button
	startBtn      *widget.Button
	stopBtn       *widget.Button
	resetBtn      *widget.Button

	// UI Components - Results Entry
	realTimeEntry       *widget.Entry
	baseTimeEntry       *widget.Entry
	additionalTimeEntry *widget.Entry
	commentEntry        *widget.Entry

	// Status action buttons
	passBtn        *widget.Button
	dqMeasureBtn   *widget.Button
	calcTimeBtn    *widget.Button

	// UI Components - Actions
	approveBtn       *widget.Button
	disqualifyBtn    *widget.Button
	skipBtn          *widget.Button
	clearSkippedBtn  *widget.Button
	loadChuggerBtn   *widget.Button
	loadFromListBtn  *widget.Button

	// Result entry status
	statusSelect     *widget.Select

	// Locks the Time field after Calculate Final Time is clicked
	lockedTimeValue  string

	// External clock mode
	useExternalClock    bool
	useExternalClockBtn *widget.Button
	extClockStopChan    chan struct{}
	lastExternalTime    string

	availableParticipants       []models.Participant
	allParticipants             []models.Participant
	skippedParticipants         []models.Participant
	contestQueue                []models.Participant
	currentChugger              *models.Participant
	currentResult               *models.Result
	
	// Track selected items in lists
	availableListSelected       widget.ListItemID
	// Track selected items in lists
	allParticipantsListSelected widget.ListItemID
}

// NewChugManager creates a new chug manager window
func NewChugManager(app fyne.App) *ChugManager {
	cm := &ChugManager{
		app:           app,
		window:        app.NewWindow("Chug Manager - Contest Execution"),
		stopTimerChan: make(chan bool),
	}

	cm.initializeManagers()
	cm.setupUI()
	cm.loadData()
	return cm
}

// initializeManagers sets up data managers
func (cm *ChugManager) initializeManagers() {
	cm.participantMgr = data.NewParticipantManager()
	cm.resultMgr = data.NewResultManager()
}

// setupUI initializes the chug manager UI
func (cm *ChugManager) setupUI() {
	cm.window.Resize(fyne.NewSize(1200, 800))
	cm.window.SetFixedSize(false)

	// Initialize components
	cm.createContestComponents()
	cm.createTimerComponents()
	cm.createListComponents()
	cm.createResultComponents()
	cm.createActionComponents()

	// Create layout
	content := cm.createLayout()
	cm.window.SetContent(content)

	// Initialize timer display
	cm.updateTimerDisplay()
}

// createContestComponents creates contest selection components
func (cm *ChugManager) createContestComponents() {
	disciplines := []string{
		models.DisciplineBottle,
		models.DisciplineHalfTankard,
		models.DisciplineFullTankard,
		models.DisciplineBierStaphette,
		models.DisciplineMegaMedley,
		models.DisciplineTeamClash,
	}

	cm.disciplineSelect = widget.NewSelect(disciplines, cm.onDisciplineSelected)
	// Note: Don't set selected here to avoid crash - will be set after UI setup

	cm.timePerEventEntry = widget.NewEntry()
	cm.timePerEventEntry.SetText("5")
	cm.timePerEventEntry.SetPlaceHolder("Time per event (minutes)")
}

// createTimerComponents creates timer-related components
func (cm *ChugManager) createTimerComponents() {
	cm.timerWidget = canvas.NewText("00:00:00.000", theme.ForegroundColor())
	cm.timerWidget.TextSize = theme.TextSize() * 4
	cm.timerWidget.TextStyle = fyne.TextStyle{Bold: true}
	cm.timerWidget.Alignment = fyne.TextAlignCenter

	cm.readyCheckBtn = widget.NewButton("Ready Check", cm.readyCheck)
	cm.startBtn = widget.NewButton("Start", cm.startTimer)
	cm.stopBtn = widget.NewButton("Stop", cm.stopTimerFunc)
	cm.resetBtn = widget.NewButton("Reset", cm.onResetButtonClicked)

	cm.useExternalClockBtn = widget.NewButton("Use External Clock", cm.toggleExternalClock)

	// Initially disable some buttons
	cm.startBtn.Disable()
	cm.stopBtn.Disable()
}
// createListComponents creates the participant list components
func (cm *ChugManager) createListComponents() {
	// Available participants list (filtered by discipline)
	cm.availableList = widget.NewList(
		func() int { return len(cm.availableParticipants) },
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
			if id >= 0 && id < len(cm.availableParticipants) {
				p := cm.availableParticipants[id]
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

	// Add OnSelected callback for availableList to track selections
	cm.availableList.OnSelected = func(id widget.ListItemID) {
		cm.availableListSelected = id
		cm.loadFromListBtn.Enable()
	}

	// All participants list
	cm.allParticipantsList = widget.NewList(
		func() int { return len(cm.allParticipants) },
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
			if id >= 0 && id < len(cm.allParticipants) {
				p := cm.allParticipants[id]
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

	// Add OnSelected callback for allParticipantsList to track selections
	cm.allParticipantsList.OnSelected = func(id widget.ListItemID) {
		cm.allParticipantsListSelected = id
		// Selecting in the All Participants list does NOT enable Load From List;
		// that button only responds to selections in Participants in Discipline.
	}

	// Current chugger display
	cm.currentNameLabel = widget.NewLabel("No participant loaded")
	cm.currentProgramLabel = widget.NewLabel("")
	cm.currentTeamLabel = widget.NewLabel("")
	cm.currentTriesLabel = widget.NewLabel("")

	cm.currentChuggerCard = widget.NewCard("Current Chugger", "",
		container.NewVBox(
			cm.currentNameLabel,
			cm.currentTriesLabel,
		),
	)
}

// createResultComponents creates result entry components
func (cm *ChugManager) createResultComponents() {
	cm.realTimeEntry = widget.NewEntry()
	cm.realTimeEntry.SetPlaceHolder("MM:SS.mmmm or HH:MM:SS.mmmm")
	cm.realTimeEntry.OnChanged = func(s string) {
		if cm.lockedTimeValue != "" && s != cm.lockedTimeValue {
			cm.realTimeEntry.SetText(cm.lockedTimeValue)
		}
	}

	cm.baseTimeEntry = widget.NewEntry()
	cm.baseTimeEntry.SetPlaceHolder("Base time")

	cm.additionalTimeEntry = widget.NewEntry()
	cm.additionalTimeEntry.SetPlaceHolder("Additional time")

	cm.commentEntry = widget.NewEntry()
	cm.commentEntry.SetPlaceHolder("Comments")
	cm.commentEntry.MultiLine = true

	// Status action buttons
	cm.passBtn = widget.NewButton("Mark as Pass", func() {
		if cm.disciplineSelect.Selected == models.DisciplineBottle {
			// Bottle: show the additional-time dialog; it handles save + moveToNext
			cm.showBottleResultDialog(models.StatusPass)
		} else {
			if err := cm.validateAndSaveResult(models.StatusPass); err != nil {
				dialog.ShowError(err, cm.window)
				return
			}
			dialog.ShowInformation("Result Saved", "Participant marked as passed", cm.window)
			cm.moveToNextParticipant()
		}
	})

	cm.dqMeasureBtn = widget.NewButton("Disqualify + Measure Time", func() {
		// Show the additional-time dialog; result will be saved as Disqualified
		cm.showBottleResultDialog(models.StatusDisqualified)
	})

	// Initially disable status buttons
	cm.passBtn.Disable()
	cm.dqMeasureBtn.Disable()

	// Calculate Final Time button
	cm.calcTimeBtn = widget.NewButton("Calculate Final Time", func() {
		baseStr := strings.TrimSpace(cm.baseTimeEntry.Text)
		additionalStr := strings.TrimSpace(cm.additionalTimeEntry.Text)

		var missing []string
		var nanFields []string

		if baseStr == "" {
			missing = append(missing, "Base Time")
		} else if strings.EqualFold(baseStr, "nan") {
			nanFields = append(nanFields, "Base Time")
		}
		if additionalStr == "" {
			missing = append(missing, "Additional Time")
		} else if strings.EqualFold(additionalStr, "nan") {
			nanFields = append(nanFields, "Additional Time")
		}

		if len(missing) > 0 {
			dialog.ShowError(fmt.Errorf(
				"the following fields are required:\n  • %s",
				strings.Join(missing, "\n  • "),
			), cm.window)
			return
		}
		if len(nanFields) > 0 {
			dialog.ShowError(fmt.Errorf(
				"NaN is not accepted in:\n  • %s\n\nEnter a valid time value (e.g. 00:00:05.000).",
				strings.Join(nanFields, "\n  • "),
			), cm.window)
			return
		}

		// Normalize both inputs so formats like "5" (seconds) or "5.000"
		// are correctly interpreted before parsing milliseconds.
		normBase := utils.ParseAndPadTimeString(baseStr)
		if normBase == config.NoKey {
			dialog.ShowError(fmt.Errorf("Base Time has an invalid format: %s", baseStr), cm.window)
			return
		}
		normAdditional := utils.ParseAndPadTimeString(additionalStr)
		if normAdditional == config.NoKey {
			dialog.ShowError(fmt.Errorf("Additional Time has an invalid format: %s", additionalStr), cm.window)
			return
		}

		baseMs := utils.ParseTimeForComparison(normBase)
		additionalMs := utils.ParseTimeForComparison(normAdditional)

		totalTms := baseMs + additionalMs
		h := totalTms / (60 * 60 * 10000)
		m := (totalTms / (60 * 10000)) % 60
		s := (totalTms / 10000) % 60
		ms := totalTms % 10000
		cm.realTimeEntry.SetText(fmt.Sprintf("%02d:%02d:%02d.%04d", h, m, s, ms))
		cm.lockedTimeValue = cm.realTimeEntry.Text
	})

	// Status select for manual result entry
	cm.statusSelect = widget.NewSelect([]string{"Pass", "Disqualified"}, nil)
	cm.statusSelect.SetSelected("Pass")
}

// createActionComponents creates action buttons
func (cm *ChugManager) createActionComponents() {
	cm.approveBtn = widget.NewButton("Enter Result Manually", cm.enterResultManually)
	cm.disqualifyBtn = widget.NewButton("Disqualify", cm.disqualifyParticipant)
	cm.skipBtn = widget.NewButton("Skip Participant", cm.skipParticipant)
	cm.clearSkippedBtn = widget.NewButton("Clear Skipped", cm.clearSkippedParticipants)
	cm.loadChuggerBtn = widget.NewButton("Load Next Chugger", cm.loadNextChugger)
	cm.loadFromListBtn = widget.NewButton("Load From List", cm.loadFromList)

	// Initially disable some buttons
	cm.approveBtn.Disable()
	cm.disqualifyBtn.Disable()
	cm.loadFromListBtn.Disable() // Will be enabled based on list selection
}

// createLayout creates the main layout
func (cm *ChugManager) createLayout() fyne.CanvasObject {
	// Contest setup section
	contestSetup := widget.NewCard("Contest Setup", "",
		container.NewVBox(
			widget.NewFormItem("Discipline", cm.disciplineSelect).Widget,
			widget.NewFormItem("Time per Event", cm.timePerEventEntry).Widget,
		),
	)

	// Timer section
	timerSection := widget.NewCard("Contest Timer", "",
		container.NewVBox(
			cm.timerWidget,
			container.NewHBox(
				cm.readyCheckBtn,
				cm.startBtn,
				cm.stopBtn,
				cm.resetBtn,
			),
			cm.useExternalClockBtn,
		),
	)

	// Participant management section with action buttons
	participantMgmt := widget.NewCard("Participant Management", "",
		container.NewVBox(
			cm.loadChuggerBtn,
			cm.loadFromListBtn,
			cm.skipBtn,
			cm.clearSkippedBtn,
		),
	)

	// Current chugger and result entry combined
	currentChuggerAndResultCard := widget.NewCard("Current Chugger & Result Entry", "",
		container.NewVBox(
			cm.currentChuggerCard,

			// Result entry section
			widget.NewLabel("Result Entry"),
			widget.NewFormItem("Base Time", cm.baseTimeEntry).Widget,
			widget.NewFormItem("Additional Time", cm.additionalTimeEntry).Widget,
			widget.NewFormItem("Time", cm.realTimeEntry).Widget,
			widget.NewFormItem("Comment", cm.commentEntry).Widget,
			widget.NewFormItem("Status", cm.statusSelect).Widget,

			widget.NewSeparator(),
			widget.NewLabel("Status Actions:"),
			container.NewHBox(
				cm.passBtn,
				cm.dqMeasureBtn,
				cm.calcTimeBtn,
			),

			widget.NewSeparator(),
			container.NewHBox(
				cm.approveBtn,
				cm.disqualifyBtn,
			),
		),
	)

	// Middle panel with participant lists as tabs
	participantTabs := container.NewAppTabs(
		container.NewTabItem("Participants in Discipline", cm.availableList),
		container.NewTabItem("All Participants", cm.allParticipantsList),
	)
	participantTabs.SetTabLocation(container.TabLocationTop)

	// Left panel with combined current chugger and result entry
	leftPanel := container.NewVBox(
		contestSetup,
		timerSection,
		participantMgmt,
		currentChuggerAndResultCard,
	)

	// Main layout
	mainContainer := container.NewHSplit(leftPanel, participantTabs)
	mainContainer.SetOffset(0.4)
	return mainContainer
}

// Timer methods
func (cm *ChugManager) readyCheck() {
	if cm.currentChugger == nil {
		dialog.ShowError(fmt.Errorf("no participant loaded"), cm.window)
		return
	}

	dialog.ShowInformation("Ready Check",
		fmt.Sprintf("Ready to start timer for %s?", cm.currentChugger.Name),
		cm.window)

	cm.startBtn.Enable()
	cm.approveBtn.Enable()
}

func (cm *ChugManager) startTimer() {
	cm.timerState.Running = true
	cm.timerState.StartTime = time.Now()
	cm.timerState.Paused = false
	cm.timerState.PausedTime = 0
	cm.lastExternalTime = ""

	cm.startBtn.Disable()
	cm.stopBtn.Enable()

	// External clock mode: subscribe to GlobalExternalClock.TimeChan instead of
	// starting an internal ticker.
	if cm.useExternalClock {
		if GlobalExternalClock == nil || !GlobalExternalClock.IsConnected() {
			dialog.ShowError(fmt.Errorf("external clock is not connected – connect it in Configuration first"), cm.window)
			cm.startBtn.Enable()
			cm.stopBtn.Disable()
			cm.timerState.Running = false
			return
		}
		cm.extClockStopChan = make(chan struct{})
		go cm.runExternalClockReceiver()
		return
	}

	// Internal timer
	cm.ticker = time.NewTicker(10 * time.Millisecond)
	go cm.runTimer()
}

// runExternalClockReceiver receives time values from the external clock and
// updates the timer display. It runs until stopTimerFunc closes extClockStopChan.
func (cm *ChugManager) runExternalClockReceiver() {
	for {
		select {
		case <-cm.extClockStopChan:
			return
		case t, ok := <-GlobalExternalClock.TimeChan:
			if !ok {
				return
			}
			cm.lastExternalTime = t
			cm.timerWidget.Text = t
			cm.timerWidget.Refresh()
		}
	}
}

// toggleExternalClock switches between internal and external clock modes.
func (cm *ChugManager) toggleExternalClock() {
	cm.useExternalClock = !cm.useExternalClock
	if cm.useExternalClock {
		cm.useExternalClockBtn.SetText("✓ Using External Clock")
	} else {
		cm.useExternalClockBtn.SetText("Use External Clock")
	}
}

func (cm *ChugManager) stopTimerFunc() {
	cm.timerState.Running = false

	// Stop external clock receiver if active
	if cm.extClockStopChan != nil {
		close(cm.extClockStopChan)
		cm.extClockStopChan = nil
	}

	if cm.ticker != nil {
		cm.ticker.Stop()
	}

	cm.startBtn.Enable()
	cm.stopBtn.Disable()

	// Auto-fill time fields
	if cm.useExternalClock && cm.lastExternalTime != "" {
		cm.baseTimeEntry.SetText(cm.lastExternalTime)
		// Use a non-zero sentinel so tries decrement logic triggers
		cm.timerState.Duration = 1
	} else if cm.timerState.Duration > 0 {
		timeStr := cm.formatDuration(cm.timerState.Duration)
		// Base time is the timer value
		cm.baseTimeEntry.SetText(timeStr)
		// Additional time is manual input only — never updated automatically
		// Time field is left empty so Enter Result Manually calculates base+additional
	}

	// Decrement tries for this discipline every time the timer stops (one stop = one try used).
	if cm.currentChugger != nil && cm.timerState.Duration > 0 {
		discipline := cm.disciplineSelect.Selected
		_ = cm.participantMgr.DecrementTries(cm.currentChugger.Name, discipline)
		if err := cm.participantMgr.SaveParticipants(); err != nil {
			dialog.ShowError(fmt.Errorf("error saving participant tries: %w", err), cm.window)
		}
		// Refresh cm.currentChugger with the post-decrement data from the manager.
		// Do NOT reassign cm.allParticipants — that would mistakenly re-add
		// the current chugger (who was removed when loaded).
		for _, p := range cm.participantMgr.GetParticipants() {
			if p.Name == cm.currentChugger.Name {
				*cm.currentChugger = p
				break
			}
		}
		cm.updateCurrentChuggerDisplay()
	}

	// Enable buttons — for all disciplines we wait for the user to pick an action.
	// For Bottle: passBtn opens the additional-time dialog; dqMeasureBtn does the same but DQ.
	// approveBtn/disqualifyBtn also remain available for quick actions.
	cm.approveBtn.Enable()
	cm.disqualifyBtn.Enable()
	cm.passBtn.Enable()
	cm.dqMeasureBtn.Enable()
}

func (cm *ChugManager) resetTimer() {
	cm.timerState = models.TimerState{}
	if cm.ticker != nil {
		cm.ticker.Stop()
	}

	cm.startBtn.Disable()
	cm.stopBtn.Disable()
	cm.approveBtn.Disable()
	cm.disqualifyBtn.Disable()
	cm.passBtn.Disable()
	cm.dqMeasureBtn.Disable()

	cm.updateTimerDisplay()
	cm.clearResultForm()
}

func (cm *ChugManager) onResetButtonClicked() {
	cm.resetTimer()
	cm.loadChuggerBtn.Enable()
	cm.loadFromListBtn.Enable()
}

func (cm *ChugManager) runTimer() {
	defer cm.ticker.Stop()

	for {
		select {
		case <-cm.ticker.C:
			if cm.timerState.Running && !cm.timerState.Paused {
				cm.timerState.Duration = time.Since(cm.timerState.StartTime) + cm.timerState.PausedTime
				cm.updateTimerDisplay()
			}
		case <-cm.stopTimerChan:
			return
		}
	}
}

func (cm *ChugManager) updateTimerDisplay() {
	timeStr := cm.formatDuration(cm.timerState.Duration)
	cm.timerWidget.Text = timeStr
	cm.timerWidget.Refresh()
}

func (cm *ChugManager) formatDuration(d time.Duration) string {
	// Use tenths-of-millisecond (100µs) resolution for 4 decimal places
	totalTms := int64(d / (time.Millisecond / 10))

	hours := totalTms / (60 * 60 * 10000)
	minutes := (totalTms / (60 * 10000)) % 60
	seconds := (totalTms / 10000) % 60
	subms := totalTms % 10000

	return fmt.Sprintf("%02d:%02d:%02d.%04d", hours, minutes, seconds, subms)
}

// Participant management methods
func (cm *ChugManager) loadNextChugger() {
	var nextParticipant *models.Participant

	// Check if there are participants in queue
	if len(cm.contestQueue) > 0 {
		nextParticipant = &cm.contestQueue[0]
		cm.contestQueue = cm.contestQueue[1:]
	} else if len(cm.availableParticipants) > 0 {
		nextParticipant = &cm.availableParticipants[0]
	} else {
		dialog.ShowError(fmt.Errorf("no participants available"), cm.window)
		return
	}

	cm.currentChugger = nextParticipant
	cm.updateCurrentChuggerDisplay()

	// Clear result form and reset timer for new participant
	cm.clearResultForm()
	cm.resetTimer()

	// Disable the load button until user skips or resets
	cm.loadChuggerBtn.Disable()
}

func (cm *ChugManager) enterResultManually() {
	if cm.currentChugger == nil {
		dialog.ShowError(fmt.Errorf("no participant loaded"), cm.window)
		return
	}

	// Map display value to status constant
	selectedStatus := models.StatusPass
	if cm.statusSelect.Selected == "Disqualified" {
		selectedStatus = models.StatusDisqualified
	}

	if err := cm.validateAndSaveResult(selectedStatus); err != nil {
		dialog.ShowError(err, cm.window)
		return
	}

	// Decrement tries only if the timer never ran for this participant.
	// If the timer was stopped, tries were already decremented at that point.
	if cm.timerState.Duration == 0 {
		discipline := cm.disciplineSelect.Selected
		_ = cm.participantMgr.DecrementTries(cm.currentChugger.Name, discipline)
		if err := cm.participantMgr.SaveParticipants(); err != nil {
			dialog.ShowError(fmt.Errorf("error saving participant tries: %w", err), cm.window)
		}
	}

	label := "Result Saved"
	msg := fmt.Sprintf("Result for %s saved as %s", cm.currentChugger.Name, selectedStatus)
	dialog.ShowInformation(label, msg, cm.window)
	cm.moveToNextParticipant()
}

func (cm *ChugManager) disqualifyParticipant() {
	if cm.currentChugger == nil {
		dialog.ShowError(fmt.Errorf("no participant loaded"), cm.window)
		return
	}
	
	if err := cm.validateAndSaveResult(models.StatusDisqualified); err != nil {
		dialog.ShowError(err, cm.window)
		return
	}
	
	dialog.ShowInformation("Disqualified", fmt.Sprintf("%s disqualified", cm.currentChugger.Name), cm.window)
	cm.moveToNextParticipant()
}

// loadFromList loads a participant from the "Participants in Discipline" list only.
func (cm *ChugManager) loadFromList() {
	var nextParticipant *models.Participant

	// Only source from the discipline-filtered list
	if cm.availableListSelected >= 0 && cm.availableListSelected < len(cm.availableParticipants) {
		nextParticipant = &cm.availableParticipants[cm.availableListSelected]
	}

	if nextParticipant == nil {
		dialog.ShowError(fmt.Errorf("please select a participant from the Participants in Discipline list"), cm.window)
		return
	}

	cm.currentChugger = nextParticipant
	cm.updateCurrentChuggerDisplay()

	// Clear result form and reset timer for new participant
	cm.clearResultForm()
	cm.resetTimer()

	// Disable the load buttons until user skips or resets
	cm.loadChuggerBtn.Disable()
	cm.loadFromListBtn.Disable()

	// Clear selections
	cm.availableList.Unselect(-1)
	cm.allParticipantsList.Unselect(-1)
	cm.availableListSelected = -1
	cm.allParticipantsListSelected = -1
}

func (cm *ChugManager) skipParticipant() {
	if cm.currentChugger == nil {
		dialog.ShowError(fmt.Errorf("no participant loaded"), cm.window)
		return
	}

	// Track who was skipped (informational only; they remain in the lists)
	cm.skippedParticipants = append(cm.skippedParticipants, *cm.currentChugger)

	dialog.ShowInformation("Participant Skipped", fmt.Sprintf("Skipped: %s", cm.currentChugger.Name), cm.window)

	cm.currentChugger = nil
	cm.updateCurrentChuggerDisplay()
	cm.clearResultForm()
	cm.resetTimer()

	// Re-enable the load button after skipping
	cm.loadChuggerBtn.Enable()
}

func (cm *ChugManager) clearSkippedParticipants() {
	if len(cm.skippedParticipants) == 0 {
		dialog.ShowInformation("No Skipped Participants", "There are no skipped participants to clear", cm.window)
		return
	}

	cm.skippedParticipants = nil
	cm.refreshLists()
	dialog.ShowInformation("Cleared", "Skipped list has been cleared", cm.window)
}

// showBottleResultDialog opens the additional-time entry window.
// requestedStatus is the intended final status (StatusPass or StatusDisqualified).
// When StatusDisqualified the "Save" button always stores the result as DQ.
func (cm *ChugManager) showBottleResultDialog(requestedStatus string) {
	participant := cm.currentChugger
	if participant == nil {
		return
	}

	discipline := cm.disciplineSelect.Selected

	title := "Result Entry"
	if requestedStatus == models.StatusDisqualified {
		title = "Result Entry (Disqualified)"
	}

	resultWindow := fyne.CurrentApp().NewWindow(title)
	resultWindow.Resize(fyne.NewSize(420, 280))

	additionalTimeEntry := widget.NewEntry()
	additionalTimeEntry.SetPlaceHolder("Additional Time (e.g. 0:03.500)")

	saveResultFunc := func(additionalTime string, forceStatus string) {
		baseTime := strings.TrimSpace(cm.baseTimeEntry.Text)
		additionalTime = strings.TrimSpace(additionalTime)

		baseMs := utils.ParseTimeForComparison(baseTime)
		additionalMs := utils.ParseTimeForComparison(additionalTime)

		var totalTimeStr string
		finalStatus := forceStatus
		if baseMs < 0 || additionalMs < 0 {
			totalTimeStr = "NaN"
			finalStatus = models.StatusDisqualified
		} else {
			totalTms := baseMs + additionalMs
			hours := totalTms / (60 * 60 * 10000)
			minutes := (totalTms / (60 * 10000)) % 60
			seconds := (totalTms / 10000) % 60
			subms := totalTms % 10000
			totalTimeStr = fmt.Sprintf("%02d:%02d:%02d.%04d", hours, minutes, seconds, subms)
		}

		comment := ""
		if finalStatus == models.StatusDisqualified {
			comment = "Overflow"
		}

		result := models.Result{
			Name:           participant.Name,
			Discipline:     discipline,
			Time:           utils.ParseAndPadTimeString(totalTimeStr),
			BaseTime:       utils.ParseAndPadTimeString(baseTime),
			AdditionalTime: additionalTime,
			Status:         finalStatus,
			Comment:        comment,
		}

		if err := cm.resultMgr.AddResult(result); err != nil {
			dialog.ShowError(err, cm.window)
			return
		}
		if err := cm.resultMgr.SaveResults(); err != nil {
			dialog.ShowError(err, cm.window)
			return
		}

		resultWindow.Close()
		cm.moveToNextParticipant()
	}

	saveLabel := "Save (Pass)"
	if requestedStatus == models.StatusDisqualified {
		saveLabel = "Save (Disqualified)"
	}

	saveBtn := widget.NewButton(saveLabel, func() {
		saveResultFunc(additionalTimeEntry.Text, requestedStatus)
	})

	// Clean Bottle button only makes sense for a pass attempt
	cleanBottleBtn := widget.NewButton("Clean Bottle (no penalty)", func() {
		saveResultFunc("0", models.StatusPass)
	})

	disqualifyOverflowBtn := widget.NewButton("Disqualify (Overflow)", func() {
		saveResultFunc(additionalTimeEntry.Text, models.StatusDisqualified)
	})

	content := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("Participant: %s  |  Discipline: %s", participant.Name, discipline)),
		widget.NewLabel(fmt.Sprintf("Base Time: %s", cm.baseTimeEntry.Text)),
		widget.NewForm(
			widget.NewFormItem("Additional Time", additionalTimeEntry),
		),
		container.NewHBox(
			saveBtn,
			cleanBottleBtn,
			disqualifyOverflowBtn,
		),
	)

	resultWindow.SetContent(container.NewPadded(content))
	resultWindow.Show()
}

func (cm *ChugManager) validateAndSaveResult(status string) error {
	if cm.currentChugger == nil {
		return fmt.Errorf("no participant loaded")
	}

	// Ensure result manager has a file path even if LoadResults was skipped
	if config.Settings.ResultFile != "" {
		cm.resultMgr.SetFilePath(config.Settings.ResultFile)
	}

	baseTimeStr := strings.TrimSpace(cm.baseTimeEntry.Text)
	additionalTimeStr := strings.TrimSpace(cm.additionalTimeEntry.Text)
	realTime := strings.TrimSpace(cm.realTimeEntry.Text)

	// Determine final time:
	//   1. If Time field is filled → use it directly (highest priority).
	//   2. Else if Base Time is filled → calculate Base + Additional.
	//   3. Otherwise → error (no time available).
	var finalTime string
	discipline := cm.disciplineSelect.Selected

	if realTime != "" {
		// Validate format
		if realTime != "NaN" && utils.ParseAndPadTimeString(realTime) == config.NoKey {
			return fmt.Errorf("invalid time format")
		}
		finalTime = realTime
	} else if baseTimeStr != "" && baseTimeStr != config.NoKey {
		// Normalize both inputs so formats like "5" (seconds) are correctly parsed.
		normBase := utils.ParseAndPadTimeString(baseTimeStr)
		if normBase == config.NoKey {
			return fmt.Errorf("invalid base time format: %s", baseTimeStr)
		}
		normAdditional := ""
		if additionalTimeStr != "" {
			normAdditional = utils.ParseAndPadTimeString(additionalTimeStr)
			if normAdditional == config.NoKey {
				return fmt.Errorf("invalid additional time format: %s", additionalTimeStr)
			}
		} else {
			normAdditional = "00:00:00.000"
		}

		baseMs := utils.ParseTimeForComparison(normBase)
		additionalMs := utils.ParseTimeForComparison(normAdditional)

		// If either operand is NaN the total is NaN → disqualified
		if baseMs < 0 || additionalMs < 0 {
			finalTime = "NaN"
			if status == models.StatusPass {
				status = models.StatusDisqualified
			}
		} else {
			totalTms := baseMs + additionalMs
			hours := totalTms / (60 * 60 * 10000)
			minutes := (totalTms / (60 * 10000)) % 60
			seconds := (totalTms / 10000) % 60
			subms := totalTms % 10000
			finalTime = fmt.Sprintf("%02d:%02d:%02d.%04d", hours, minutes, seconds, subms)
		}
	} else if status == models.StatusPass {
		return fmt.Errorf("time is required: fill in Time or Base Time + Additional Time")
	}

	// Create result
	result := models.Result{
		Name:           cm.currentChugger.Name,
		Discipline:     discipline,
		Time:           utils.ParseAndPadTimeString(finalTime),
		BaseTime:       utils.ParseAndPadTimeString(baseTimeStr),
		AdditionalTime: additionalTimeStr,
		Status:         status,
		Comment:        strings.TrimSpace(cm.commentEntry.Text),
	}

	// Add to result manager
	if err := cm.resultMgr.AddResult(result); err != nil {
		return fmt.Errorf("error adding result: %w", err)
	}

	// Save results
	if err := cm.resultMgr.SaveResults(); err != nil {
		return fmt.Errorf("error saving results: %w", err)
	}

	return nil
}

// calculateAndSetTotalTime adds base time and additional time to get total time
func (cm *ChugManager) calculateAndSetTotalTime() {
	baseTimeStr := strings.TrimSpace(cm.baseTimeEntry.Text)
	additionalTimeStr := strings.TrimSpace(cm.additionalTimeEntry.Text)

	// Parse base time
	baseMs := utils.ParseTimeForComparison(baseTimeStr)
	// Parse additional time
	additionalMs := utils.ParseTimeForComparison(additionalTimeStr)

	// Sum the times
	totalTms := baseMs + additionalMs

	// Convert back to time string format
	hours := totalTms / (60 * 60 * 10000)
	minutes := (totalTms / (60 * 10000)) % 60
	seconds := (totalTms / 10000) % 60
	subms := totalTms % 10000

	totalTimeStr := fmt.Sprintf("%02d:%02d:%02d.%04d", hours, minutes, seconds, subms)
	cm.realTimeEntry.SetText(totalTimeStr)
}

func (cm *ChugManager) moveToNextParticipant() {
	cm.currentChugger = nil
	cm.updateCurrentChuggerDisplay()
	cm.clearResultForm()
	cm.resetTimer()

	// Always rebuild lists from the manager so try-count changes are reflected.
	// Participants stay in the list until their tries reach 0.
	cm.allParticipants = cm.participantMgr.GetParticipants()
	sort.Slice(cm.allParticipants, func(i, j int) bool {
		return cm.allParticipants[i].Name < cm.allParticipants[j].Name
	})
	cm.loadAvailableParticipants()

	// Re-enable the load button so user can load next participant
	cm.loadChuggerBtn.Enable()
}

// Event handlers
func (cm *ChugManager) onDisciplineSelected(discipline string) {
	cm.loadAvailableParticipants()
	cm.updateCurrentChuggerDisplay()
}

// Helper methods
// updateCurrentChuggerDisplay helper is declared later in file, removing duplicate

func (cm *ChugManager) clearResultForm() {
	cm.lockedTimeValue = ""
	cm.realTimeEntry.SetText("")
	cm.baseTimeEntry.SetText("")
	cm.additionalTimeEntry.SetText("")
	cm.commentEntry.SetText("")
	cm.statusSelect.SetSelected("Pass")
}

func (cm *ChugManager) refreshLists() {
	// Safety check - ensure lists are initialized
	if cm.availableList == nil || cm.allParticipantsList == nil {
		return
	}

	cm.availableList.Refresh()
	cm.allParticipantsList.Refresh()
}
func (cm *ChugManager) updateCurrentChuggerDisplay() {
	if cm.currentChugger == nil {
		cm.currentNameLabel.SetText("No participant loaded")
		cm.currentProgramLabel.SetText("")
		cm.currentTeamLabel.SetText("")
		cm.currentTriesLabel.SetText("")
		cm.currentChuggerCard.SetTitle("Current Chugger")
	} else {
		// Combine name, program, and team on one line
		combinedText := fmt.Sprintf("%s | %s | %s", cm.currentChugger.Name, cm.currentChugger.Program, cm.currentChugger.Team)
		cm.currentNameLabel.SetText(combinedText)
		cm.currentChuggerCard.SetTitle("Current Chugger: " + cm.currentChugger.Name)

		// Show remaining tries for the selected discipline.
		// Use the up-to-date in-memory participant record (post-decrement).
		discipline := cm.disciplineSelect.Selected
		triesStr := ""
		switch discipline {
		case models.DisciplineBottle:
			triesStr = cm.currentChugger.Bottle
		case models.DisciplineHalfTankard:
			triesStr = cm.currentChugger.HalfTankard
		case models.DisciplineFullTankard:
			triesStr = cm.currentChugger.FullTankard
		}
		if triesStr != "" {
			cm.currentTriesLabel.SetText(fmt.Sprintf("Tries remaining (%s): %s", discipline, triesStr))
		} else {
			cm.currentTriesLabel.SetText("")
		}
	}
}
func (cm *ChugManager) loadData() {
	// Check if participant file is configured
	if config.Settings.ParticipantFile == "" {
		dialog.ShowError(fmt.Errorf("no participant file configured - please use Contest Wizard to create a contest or set the file path in Configuration"), cm.window)
		return
	}

	// Load participants
	if utils.DoesFileExist(config.Settings.ParticipantFile) {
		if err := cm.participantMgr.LoadParticipants(config.Settings.ParticipantFile); err != nil {
			dialog.ShowError(fmt.Errorf("error loading participants: %v", err), cm.window)
		} else {
			// Copy participant file into the contest folder if it lives elsewhere,
			// then update config so saves always target the local copy.
			if config.Settings.ResultFile != "" {
				contestDir := filepath.Dir(config.Settings.ResultFile)
				destPath := filepath.Join(contestDir, "participants.json")
				srcAbs, _ := filepath.Abs(config.Settings.ParticipantFile)
				dstAbs, _ := filepath.Abs(destPath)
				if srcAbs != dstAbs {
					if err := utils.CopyFile(config.Settings.ParticipantFile, destPath); err != nil {
						dialog.ShowError(fmt.Errorf("error copying participant file to contest folder: %w", err), cm.window)
					} else {
						config.Settings.ParticipantFile = destPath
						config.SaveConfig()
						cm.participantMgr.SetFilePath(destPath)
					}
				}
			}

			cm.allParticipants = cm.participantMgr.GetParticipants()
			
			sort.Slice(cm.allParticipants, func(i, j int) bool {
				return cm.allParticipants[i].Name < cm.allParticipants[j].Name
			})
			
			cm.loadAvailableParticipants()
		}
	} else {
		dialog.ShowError(fmt.Errorf("participant file not found: %s\n\nPlease create a contest using Contest Wizard or check Configuration", config.Settings.ParticipantFile), cm.window)
		return
	}

	// Load results
	if config.Settings.ResultFile != "" && utils.DoesFileExist(config.Settings.ResultFile) {
		if err := cm.resultMgr.LoadResults(config.Settings.ResultFile); err != nil {
			dialog.ShowError(fmt.Errorf("error loading results: %w", err), cm.window)
		}
	} else if config.Settings.ResultFile != "" {
		// Create empty results file if it doesn't exist
		if err := utils.SaveListToJSONFile(config.Settings.ResultFile, []map[string]string{}); err != nil {
			dialog.ShowError(fmt.Errorf("error creating result file %s: %w", config.Settings.ResultFile, err), cm.window)
		} else {
			// Load the newly-created file so filePath is set on resultMgr
			if err := cm.resultMgr.LoadResults(config.Settings.ResultFile); err != nil {
				dialog.ShowError(fmt.Errorf("error initialising result manager: %w", err), cm.window)
			}
		}
	}

	// Set initial discipline selection after everything is loaded
	cm.disciplineSelect.SetSelected(models.DisciplineBottle)
}

func (cm *ChugManager) loadAvailableParticipants() {
	discipline := cm.disciplineSelect.Selected
	
	// Create available participants list based on discipline
	var participantsForDiscipline []models.Participant
	
	// Start with all participants
	allParticipants := cm.allParticipants
	
	for _, p := range allParticipants {
		// Simple filter logic - if it has a valid value for the discipline, include it
		include := false
		
		switch discipline {
		case "Bottle":
			if p.Bottle != "" && p.Bottle != "0" {
				include = true
			}
		case "Half Tankard":
			if p.HalfTankard != "" && p.HalfTankard != "0" {
				include = true
			}
		case "Full Tankard":
			if p.FullTankard != "" && p.FullTankard != "0" {
				include = true
			}
		}
		
		if include {
			// Check if already has a result for this discipline
			hasResult := false
			// This would require checking results, for now we just show all eligible
			
			if !hasResult {
				participantsForDiscipline = append(participantsForDiscipline, p)
			}
		}
	}
	
	cm.availableParticipants = participantsForDiscipline

	// Sort by most remaining tries descending so the next chugger loaded is always
	// the one with the most tries, with name as a tiebreaker.
	sort.Slice(cm.availableParticipants, func(i, j int) bool {
		var ti, tj int
		switch discipline {
		case models.DisciplineBottle:
			ti, _ = strconv.Atoi(cm.availableParticipants[i].Bottle)
			tj, _ = strconv.Atoi(cm.availableParticipants[j].Bottle)
		case models.DisciplineHalfTankard:
			ti, _ = strconv.Atoi(cm.availableParticipants[i].HalfTankard)
			tj, _ = strconv.Atoi(cm.availableParticipants[j].HalfTankard)
		case models.DisciplineFullTankard:
			ti, _ = strconv.Atoi(cm.availableParticipants[i].FullTankard)
			tj, _ = strconv.Atoi(cm.availableParticipants[j].FullTankard)
		}
		if ti != tj {
			return ti > tj // most tries first
		}
		return cm.availableParticipants[i].Name < cm.availableParticipants[j].Name
	})

	cm.refreshLists()
}

// Show displays the chug manager window
func (cm *ChugManager) Show() {
	cm.window.Show()
}


