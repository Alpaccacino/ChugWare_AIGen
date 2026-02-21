package ui

import (
	"fmt"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"chugware/internal/config"
	"chugware/internal/models"
)

// ConfigurationWindow handles application configuration
type ConfigurationWindow struct {
	app    fyne.App
	window fyne.Window

	// Path settings
	folderPathEntry        *widget.Entry
	folderPathContestEntry *widget.Entry
	browseFolderBtn        *widget.Button

	// File settings
	participantFileEntry *widget.Entry
	resultFileEntry      *widget.Entry
	templateFileEntry    *widget.Entry
	browseParticipantBtn *widget.Button
	browseResultBtn      *widget.Button
	browseTemplateBtn    *widget.Button

	// Contest-specific file settings
	bottleFileEntry        *widget.Entry
	halfTankardFileEntry   *widget.Entry
	fullTankardFileEntry   *widget.Entry
	bierStaphetteFileEntry *widget.Entry
	megaMedleyFileEntry    *widget.Entry
	teamClashFileEntry     *widget.Entry
	factionFileEntry       *widget.Entry
	leaderBoardFileEntry   *widget.Entry
	diplomasFileEntry      *widget.Entry

	// Trial settings
	bottleTriesEntry      *widget.Entry
	halfTankardTriesEntry *widget.Entry
	fullTankardTriesEntry *widget.Entry

	// Application settings
	maxStringLengthEntry *widget.Entry
	maxParticipantsEntry *widget.Entry
	maxEntriesEntry      *widget.Entry
	maxValueEntry        *widget.Entry
	pixelsPerCharEntry   *widget.Entry

	// External equipment
	externalPortEntry  *widget.Entry
	externalBaudSelect *widget.Select

	// Action buttons
	saveBtn  *widget.Button
	loadBtn  *widget.Button
	resetBtn *widget.Button
	clearBtn *widget.Button

	// Status labels
	statusLabel *widget.Label
}

// NewConfigurationWindow creates a new configuration window
func NewConfigurationWindow(app fyne.App) *ConfigurationWindow {
	cw := &ConfigurationWindow{
		app:    app,
		window: app.NewWindow("ChugWare Configuration"),
	}

	cw.setupUI()
	cw.loadCurrentSettings()
	return cw
}

// setupUI initializes the configuration UI
func (cw *ConfigurationWindow) setupUI() {
	cw.window.Resize(fyne.NewSize(1600, 1200))
	cw.window.CenterOnScreen()
	cw.window.SetFixedSize(false)

	// Initialize components
	cw.createPathComponents()
	cw.createFileComponents()
	cw.createTrialComponents()
	cw.createAppComponents()
	cw.createExternalEquipmentComponents()
	cw.createActionComponents()

	// Create layout
	content := cw.createLayout()
	paddedContent := container.NewPadded(container.NewScroll(content))
	cw.window.SetContent(paddedContent)
}

// createPathComponents creates path setting components
func (cw *ConfigurationWindow) createPathComponents() {
	cw.folderPathEntry = widget.NewEntry()
	cw.folderPathEntry.SetPlaceHolder("Main contest folder path")

	cw.folderPathContestEntry = widget.NewEntry()
	cw.folderPathContestEntry.SetPlaceHolder("Current contest folder path")
	cw.folderPathContestEntry.Disable() // Read-only

	cw.browseFolderBtn = widget.NewButton("Browse", cw.browseFolderPath)
}

// createFileComponents creates file setting components
func (cw *ConfigurationWindow) createFileComponents() {
	cw.participantFileEntry = widget.NewEntry()
	cw.participantFileEntry.SetPlaceHolder("Participant file path")
	cw.browseParticipantBtn = widget.NewButton("Browse", cw.browseParticipantFile)

	cw.resultFileEntry = widget.NewEntry()
	cw.resultFileEntry.SetPlaceHolder("Result file path")
	cw.browseResultBtn = widget.NewButton("Browse", cw.browseResultFile)

	cw.templateFileEntry = widget.NewEntry()
	cw.templateFileEntry.SetPlaceHolder("Template file path")
	cw.browseTemplateBtn = widget.NewButton("Browse", cw.browseTemplateFile)

	// Contest-specific files
	cw.bottleFileEntry = widget.NewEntry()
	cw.bottleFileEntry.SetPlaceHolder("Bottle contest file")

	cw.halfTankardFileEntry = widget.NewEntry()
	cw.halfTankardFileEntry.SetPlaceHolder("Half Tankard contest file")

	cw.fullTankardFileEntry = widget.NewEntry()
	cw.fullTankardFileEntry.SetPlaceHolder("Full Tankard contest file")

	cw.bierStaphetteFileEntry = widget.NewEntry()
	cw.bierStaphetteFileEntry.SetPlaceHolder("Bier Staphette contest file")

	cw.megaMedleyFileEntry = widget.NewEntry()
	cw.megaMedleyFileEntry.SetPlaceHolder("Mega Medley contest file")

	cw.teamClashFileEntry = widget.NewEntry()
	cw.teamClashFileEntry.SetPlaceHolder("Team Clash contest file")

	cw.factionFileEntry = widget.NewEntry()
	cw.factionFileEntry.SetPlaceHolder("Faction file")

	cw.leaderBoardFileEntry = widget.NewEntry()
	cw.leaderBoardFileEntry.SetPlaceHolder("Leaderboard file")

	cw.diplomasFileEntry = widget.NewEntry()
	cw.diplomasFileEntry.SetPlaceHolder("Diplomas file")
}

// createTrialComponents creates trial setting components
func (cw *ConfigurationWindow) createTrialComponents() {
	cw.bottleTriesEntry = widget.NewEntry()
	cw.bottleTriesEntry.SetText(strconv.Itoa(config.BottleTries))
	cw.bottleTriesEntry.SetPlaceHolder("Number of bottle tries")

	cw.halfTankardTriesEntry = widget.NewEntry()
	cw.halfTankardTriesEntry.SetText(strconv.Itoa(config.HalfTankardTries))
	cw.halfTankardTriesEntry.SetPlaceHolder("Number of half tankard tries")

	cw.fullTankardTriesEntry = widget.NewEntry()
	cw.fullTankardTriesEntry.SetText(strconv.Itoa(config.FullTankardTries))
	cw.fullTankardTriesEntry.SetPlaceHolder("Number of full tankard tries")
}

// createAppComponents creates application setting components
func (cw *ConfigurationWindow) createAppComponents() {
	cw.maxStringLengthEntry = widget.NewEntry()
	cw.maxStringLengthEntry.SetText(strconv.Itoa(config.MaxStringLength))
	cw.maxStringLengthEntry.SetPlaceHolder("Maximum string length")

	cw.maxParticipantsEntry = widget.NewEntry()
	cw.maxParticipantsEntry.SetText(strconv.Itoa(config.MaxParticipants))
	cw.maxParticipantsEntry.SetPlaceHolder("Maximum participants")

	cw.maxEntriesEntry = widget.NewEntry()
	cw.maxEntriesEntry.SetText(strconv.Itoa(config.MaxEntries))
	cw.maxEntriesEntry.SetPlaceHolder("Maximum entries")

	cw.maxValueEntry = widget.NewEntry()
	cw.maxValueEntry.SetText(strconv.Itoa(config.MaxValue))
	cw.maxValueEntry.SetPlaceHolder("Maximum value")

	cw.pixelsPerCharEntry = widget.NewEntry()
	cw.pixelsPerCharEntry.SetText(strconv.Itoa(config.PixelsPerChar))
	cw.pixelsPerCharEntry.SetPlaceHolder("Pixels per character")
}

// createExternalEquipmentComponents initialises the External Equipment widgets.
func (cw *ConfigurationWindow) createExternalEquipmentComponents() {
	cw.externalPortEntry = widget.NewEntry()
	cw.externalPortEntry.SetPlaceHolder("COM3 / /dev/ttyUSB0")

	cw.externalBaudSelect = widget.NewSelect(
		[]string{"9600", "19200", "38400", "57600", "115200"},
		func(v string) {
			baud, _ := strconv.Atoi(v)
			if GlobalExternalClock != nil {
				GlobalExternalClock.SetBaud(baud)
			}
		},
	)
	cw.externalBaudSelect.SetSelected("9600")
}

// createActionComponents creates action buttons
func (cw *ConfigurationWindow) createActionComponents() {
	cw.saveBtn = widget.NewButton("Save Configuration", cw.saveConfiguration)
	cw.loadBtn = widget.NewButton("Reload Configuration", cw.loadCurrentSettings)
	cw.resetBtn = widget.NewButton("Reset to Defaults", cw.resetToDefaults)
	cw.clearBtn = widget.NewButton("Clear Persistent Data", cw.clearAllSettings)

	cw.statusLabel = widget.NewLabel("Ready")
}

// Helper to create a more descriptive settings field
func createSettingsField(label, explanation string, entry *widget.Entry, btn *widget.Button) fyne.CanvasObject {
	title := widget.NewLabelWithStyle(label, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	desc := widget.NewLabelWithStyle(explanation, fyne.TextAlignLeading, fyne.TextStyle{})
	// We want to ensure description is readable. Standard label color is fine.

	// Clear placeholder as we use description
	entry.SetPlaceHolder("")

	var content fyne.CanvasObject
	if btn != nil {
		content = container.NewBorder(nil, nil, nil, btn, entry)
	} else {
		content = entry
	}

	return container.NewVBox(
		title,
		desc,
		content,
	)
}

// createLayout creates the main configuration layout
func (cw *ConfigurationWindow) createLayout() *fyne.Container {
	// Path settings section
	pathCard := widget.NewCard("Path Settings", "",
		container.NewVBox(
			createSettingsField("Main Folder", "The root directory for ChugWare data.", cw.folderPathEntry, cw.browseFolderBtn),
			createSettingsField("Contest Folder", "Directory where specific contest data will be stored.", cw.folderPathContestEntry, nil),
		),
	)

	// Essential file settings
	essentialFilesCard := widget.NewCard("Essential Files", "",
		container.NewVBox(
			createSettingsField("Participant File", "JSON file containing list of participants.", cw.participantFileEntry, cw.browseParticipantBtn),
			createSettingsField("Result File", "JSON file where contest results are saved.", cw.resultFileEntry, cw.browseResultBtn),
			createSettingsField("Template File", "Template for generating reports (optional).", cw.templateFileEntry, cw.browseTemplateBtn),
		),
	)

	// Contest-specific files
	contestFilesCard := widget.NewCard("Contest-Specific Files", "",
		container.NewVBox(
			createSettingsField("Bottle File", "Data file for Bottle discipline.", cw.bottleFileEntry, nil),
			createSettingsField("Half Tankard File", "Data file for Half-Tankard discipline.", cw.halfTankardFileEntry, nil),
			createSettingsField("Full Tankard File", "Data file for Full-Tankard discipline.", cw.fullTankardFileEntry, nil),
			createSettingsField("Bier Staphette File", "Data file for Bier Staphette discipline.", cw.bierStaphetteFileEntry, nil),
			createSettingsField("Mega Medley File", "Data file for Mega Medley discipline.", cw.megaMedleyFileEntry, nil),
			createSettingsField("Team Clash File", "Data file for Team Clash discipline.", cw.teamClashFileEntry, nil),
			createSettingsField("Faction File", "Data file for Faction info.", cw.factionFileEntry, nil),
			createSettingsField("Leaderboard File", "Data file for storing leaderboard stats.", cw.leaderBoardFileEntry, nil),
			createSettingsField("Diplomas File", "Data file for diploma generation.", cw.diplomasFileEntry, nil),
		),
	)

	// Trial settings
	trialCard := widget.NewCard("Trial Settings", "",
		container.NewVBox(
			createSettingsField("Bottle Tries", "Number of attempts allowed for Bottle discipline.", cw.bottleTriesEntry, nil),
			createSettingsField("Half Tankard Tries", "Number of attempts allowed for Half Tankard.", cw.halfTankardTriesEntry, nil),
			createSettingsField("Full Tankard Tries", "Number of attempts allowed for Full Tankard.", cw.fullTankardTriesEntry, nil),
		),
	)

	// Application settings
	appCard := widget.NewCard("Application Settings", "",
		container.NewVBox(
			createSettingsField("Max String Length", "Maximum length for text inputs.", cw.maxStringLengthEntry, nil),
			createSettingsField("Max Participants", "Maximum number of participants allowed.", cw.maxParticipantsEntry, nil),
			createSettingsField("Max Entries", "Maximum number of entries per list.", cw.maxEntriesEntry, nil),
			createSettingsField("Max Value", "Maximum numerical value for scores.", cw.maxValueEntry, nil),
			createSettingsField("Pixels Per Char", "Display scaling factor.", cw.pixelsPerCharEntry, nil),
		),
	)

	// External Equipment card
	var extEquipContent fyne.CanvasObject
	if GlobalExternalClock != nil {
		connectBtn := GlobalExternalClock.BuildConnectButton(
			func() string { return cw.externalPortEntry.Text },
			func() int {
				baud, _ := strconv.Atoi(cw.externalBaudSelect.Selected)
				if baud == 0 {
					baud = 9600
				}
				return baud
			},
		)
		viewLogsBtn := widget.NewButton("View Logs", func() {
			GlobalExternalClock.ShowLogWindow()
		})
		statusLbl := GlobalExternalClock.BuildStatusLabel()
		extEquipContent = container.NewVBox(
			createSettingsField("Serial Port", "COM port or device path of the external clock interface.", cw.externalPortEntry, nil),
			widget.NewLabel("Baud Rate"),
			cw.externalBaudSelect,
			container.NewHBox(connectBtn, viewLogsBtn),
			statusLbl,
		)
	} else {
		extEquipContent = widget.NewLabel("External clock not initialised.")
	}
	extEquipCard := widget.NewCard("External Equipment", "", extEquipContent)

	// Action buttons
	actionContainer := container.NewHBox(
		cw.saveBtn,
		cw.loadBtn,
		cw.resetBtn,
		cw.clearBtn,
	)

	// Status
	statusContainer := container.NewCenter(cw.statusLabel)

	// Main layout
	return container.NewVBox(
		pathCard,
		essentialFilesCard,
		contestFilesCard,
		trialCard,
		appCard,
		extEquipCard,
		widget.NewSeparator(),
		actionContainer,
		statusContainer,
	)
}

// Event handlers
func (cw *ConfigurationWindow) browseFolderPath() {
	dialog.ShowFolderOpen(func(folder fyne.ListableURI, err error) {
		if err == nil && folder != nil {
			cw.folderPathEntry.SetText(folder.Path())
		}
	}, cw.window)
}

func (cw *ConfigurationWindow) browseParticipantFile() {
	cw.browseFile("Select Participant File", cw.participantFileEntry)
}

func (cw *ConfigurationWindow) browseResultFile() {
	cw.browseFile("Select Result File", cw.resultFileEntry)
}

func (cw *ConfigurationWindow) browseTemplateFile() {
	cw.browseFile("Select Template File", cw.templateFileEntry)
}

func (cw *ConfigurationWindow) browseFile(title string, entry *widget.Entry) {
	dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {
		if err == nil && file != nil {
			defer file.Close()
			entry.SetText(file.URI().Path())
		}
	}, cw.window)
}

// Configuration operations
func (cw *ConfigurationWindow) saveConfiguration() {
	// Update configuration with form values
	config.Settings.FolderPath = cw.folderPathEntry.Text
	config.Settings.ExternalClockPort = cw.externalPortEntry.Text
	if baud, err := strconv.Atoi(cw.externalBaudSelect.Selected); err == nil {
		config.Settings.ExternalClockBaud = baud
	}
	config.Settings.ParticipantFile = cw.participantFileEntry.Text
	config.Settings.ResultFile = cw.resultFileEntry.Text
	config.Settings.TemplateFile = cw.templateFileEntry.Text
	config.Settings.BottleFile = cw.bottleFileEntry.Text
	config.Settings.HalfTankardFile = cw.halfTankardFileEntry.Text
	config.Settings.FullTankardFile = cw.fullTankardFileEntry.Text
	config.Settings.BierStaphetteFile = cw.bierStaphetteFileEntry.Text
	config.Settings.MegaMedleyFile = cw.megaMedleyFileEntry.Text
	config.Settings.TeamClashFile = cw.teamClashFileEntry.Text
	config.Settings.FactionFile = cw.factionFileEntry.Text
	config.Settings.LeaderBoardFile = cw.leaderBoardFileEntry.Text
	config.Settings.DiplomasFile = cw.diplomasFileEntry.Text

	// Validate and update numeric settings
	if err := cw.validateAndUpdateNumericSettings(); err != nil {
		dialog.ShowError(err, cw.window)
		return
	}

	// Save configuration
	if err := config.SaveConfig(); err != nil {
		dialog.ShowError(fmt.Errorf("error saving configuration: %w", err), cw.window)
		return
	}

	cw.statusLabel.SetText("Configuration saved successfully")
	dialog.ShowInformation("Success", "Configuration has been saved", cw.window)
}

func (cw *ConfigurationWindow) validateAndUpdateNumericSettings() error {
	// Validate trial settings
	if _, err := strconv.Atoi(cw.bottleTriesEntry.Text); err != nil {
		return fmt.Errorf("invalid bottle tries value: %w", err)
	}

	if _, err := strconv.Atoi(cw.halfTankardTriesEntry.Text); err != nil {
		return fmt.Errorf("invalid half tankard tries value: %w", err)
	}

	if _, err := strconv.Atoi(cw.fullTankardTriesEntry.Text); err != nil {
		return fmt.Errorf("invalid full tankard tries value: %w", err)
	}

	// Validate application settings (these are read-only for display, actual constants are in config)
	if _, err := strconv.Atoi(cw.maxStringLengthEntry.Text); err != nil {
		return fmt.Errorf("invalid max string length value: %w", err)
	}

	if _, err := strconv.Atoi(cw.maxParticipantsEntry.Text); err != nil {
		return fmt.Errorf("invalid max participants value: %w", err)
	}

	if _, err := strconv.Atoi(cw.maxEntriesEntry.Text); err != nil {
		return fmt.Errorf("invalid max entries value: %w", err)
	}

	if _, err := strconv.Atoi(cw.maxValueEntry.Text); err != nil {
		return fmt.Errorf("invalid max value: %w", err)
	}

	if _, err := strconv.Atoi(cw.pixelsPerCharEntry.Text); err != nil {
		return fmt.Errorf("invalid pixels per char value: %w", err)
	}

	return nil
}

func (cw *ConfigurationWindow) loadCurrentSettings() {
	// Load current configuration
	if err := config.LoadConfig(); err != nil {
		dialog.ShowError(fmt.Errorf("error loading configuration: %w", err), cw.window)
		return
	}

	// Update form with current settings
	cw.folderPathEntry.SetText(config.Settings.FolderPath)
	if config.Settings.ExternalClockPort != "" {
		cw.externalPortEntry.SetText(config.Settings.ExternalClockPort)
	}
	if config.Settings.ExternalClockBaud > 0 {
		cw.externalBaudSelect.SetSelected(strconv.Itoa(config.Settings.ExternalClockBaud))
	}
	cw.folderPathContestEntry.SetText(config.Settings.FolderPathContestNameAndDate)
	cw.participantFileEntry.SetText(config.Settings.ParticipantFile)
	cw.resultFileEntry.SetText(config.Settings.ResultFile)
	cw.templateFileEntry.SetText(config.Settings.TemplateFile)
	cw.bottleFileEntry.SetText(config.Settings.BottleFile)
	cw.halfTankardFileEntry.SetText(config.Settings.HalfTankardFile)
	cw.fullTankardFileEntry.SetText(config.Settings.FullTankardFile)
	cw.bierStaphetteFileEntry.SetText(config.Settings.BierStaphetteFile)
	cw.megaMedleyFileEntry.SetText(config.Settings.MegaMedleyFile)
	cw.teamClashFileEntry.SetText(config.Settings.TeamClashFile)
	cw.factionFileEntry.SetText(config.Settings.FactionFile)
	cw.leaderBoardFileEntry.SetText(config.Settings.LeaderBoardFile)
	cw.diplomasFileEntry.SetText(config.Settings.DiplomasFile)

	cw.statusLabel.SetText("Configuration loaded")
}

func (cw *ConfigurationWindow) resetToDefaults() {
	dialog.ShowConfirm("Reset Configuration",
		"Are you sure you want to reset all settings to default values?",
		func(confirmed bool) {
			if confirmed {
				cw.performReset()
			}
		}, cw.window)
}

func (cw *ConfigurationWindow) performReset() {
	// Reset to default values
	cw.folderPathEntry.SetText("")
	cw.folderPathContestEntry.SetText("")
	cw.participantFileEntry.SetText("")
	cw.resultFileEntry.SetText("")
	cw.templateFileEntry.SetText("")
	cw.bottleFileEntry.SetText("")
	cw.halfTankardFileEntry.SetText("")
	cw.fullTankardFileEntry.SetText("")
	cw.bierStaphetteFileEntry.SetText("")
	cw.megaMedleyFileEntry.SetText("")
	cw.teamClashFileEntry.SetText("")
	cw.factionFileEntry.SetText("")
	cw.leaderBoardFileEntry.SetText("")
	cw.diplomasFileEntry.SetText("")

	// Reset trial settings to defaults
	cw.bottleTriesEntry.SetText(strconv.Itoa(config.BottleTries))
	cw.halfTankardTriesEntry.SetText(strconv.Itoa(config.HalfTankardTries))
	cw.fullTankardTriesEntry.SetText(strconv.Itoa(config.FullTankardTries))

	cw.statusLabel.SetText("Reset to default values")
}

func (cw *ConfigurationWindow) clearAllSettings() {
	dialog.ShowConfirm("Clear All Settings",
		"Are you sure you want to clear all persistent settings? This cannot be undone.",
		func(confirmed bool) {
			if confirmed {
				cw.performClearAll()
			}
		}, cw.window)
}

func (cw *ConfigurationWindow) performClearAll() {
	// Remove the config file from disk
	if config.ConfigFile != "" {
		if err := os.Remove(config.ConfigFile); err != nil && !os.IsNotExist(err) {
			dialog.ShowError(fmt.Errorf("error removing config file: %w", err), cw.window)
			return
		}
	}

	// Reset in-memory settings to zero value
	config.Settings = models.ContestSettings{}

	// Update form to reflect empty state
	cw.performReset()
	cw.externalPortEntry.SetText("")
	cw.externalBaudSelect.SetSelected("9600")
	cw.statusLabel.SetText("All persistent data cleared")

	dialog.ShowInformation("Data Cleared", "All persistent settings have been removed from disk", cw.window)
}

// Show displays the configuration window
func (cw *ConfigurationWindow) Show() {
	cw.window.Show()
}
