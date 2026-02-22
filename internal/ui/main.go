package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"chugware/internal/config"
	"chugware/internal/resources"
)

// MainWindow represents the main application window
type MainWindow struct {
	app    fyne.App
	window fyne.Window

	// UI components
	contestWizardBtn   *widget.Button
	addParticipantsBtn *widget.Button
	chugManagerBtn     *widget.Button
	configurationBtn   *widget.Button
	finishContestBtn   *widget.Button
	exitBtn            *widget.Button
}

// NewMainWindow creates a new main window
func NewMainWindow(app fyne.App) *MainWindow {
	mw := &MainWindow{
		app:    app,
		window: app.NewWindow("ChugWare2 - Contest Management System"),
	}

	mw.setupUI()
	return mw
}

// setupUI initializes the main window UI
func (mw *MainWindow) setupUI() {
	mw.window.Resize(fyne.NewSize(900, 600))
	mw.window.CenterOnScreen()
	mw.window.SetFixedSize(false) // Allow resizing

	// Create menu buttons
	mw.createMenuButtons()

	// Create the main menu (background + buttons)
	menuContainer := mw.createMainMenu()

	// Set up the main content – no extra padding so the background fills the window
	content := menuContainer

	// Create application menu
	mainMenu := mw.createApplicationMenu()
	mw.window.SetMainMenu(mainMenu)

	mw.window.SetContent(content)
}

// createMenuButtons initializes all menu buttons
func (mw *MainWindow) createMenuButtons() {
	mw.contestWizardBtn = widget.NewButton("Contest Wizard", mw.openContestWizard)
	mw.addParticipantsBtn = widget.NewButton("Add Participants", mw.openAddParticipants)
	mw.chugManagerBtn = widget.NewButton("Chug Manager", mw.openChugManager)
	mw.configurationBtn = widget.NewButton("Configuration", mw.openConfiguration)
	mw.finishContestBtn = widget.NewButton("Finish Contest", mw.openFinishContest)
	mw.exitBtn = widget.NewButton("Exit", mw.exitApplication)
}

// createMainMenu creates the main menu layout
func (mw *MainWindow) createMainMenu() *fyne.Container {
	// Background image
	bg := canvas.NewImageFromResource(resources.CoolSigge)
	bg.FillMode = canvas.ImageFillContain

	// Title
	title := widget.NewLabel("ChugWare2")
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle.Bold = true

	// Menu buttons in a vertical layout with proper spacing
	menuButtons := container.NewVBox(
		mw.contestWizardBtn,
		mw.addParticipantsBtn,
		mw.chugManagerBtn,
		mw.configurationBtn,
		mw.finishContestBtn,
		widget.NewSeparator(),
		mw.exitBtn,
	)

	// Place buttons in a centred overlay on top of the background image
	menuContainer := container.NewStack(
		bg,
		container.NewCenter(
			container.NewVBox(
				title,
				widget.NewSeparator(),
				menuButtons,
			),
		),
	)

	return menuContainer
}

// createApplicationMenu creates the application menu bar
func (mw *MainWindow) createApplicationMenu() *fyne.MainMenu {
	// File menu
	contestWizardItem := fyne.NewMenuItem("Contest Wizard", mw.openContestWizard)
	addParticipantsItem := fyne.NewMenuItem("Add Participants", mw.openAddParticipants)
	chugManagerItem := fyne.NewMenuItem("Chug Manager", mw.openChugManager)
	configurationItem := fyne.NewMenuItem("Configuration", mw.openConfiguration)
	finishContestItem := fyne.NewMenuItem("Finish Contest", mw.openFinishContest)
	exitItem := fyne.NewMenuItem("Exit", mw.exitApplication)

	fileMenu := fyne.NewMenu("Menu",
		contestWizardItem,
		addParticipantsItem,
		chugManagerItem,
		configurationItem,
		finishContestItem,
		fyne.NewMenuItemSeparator(),
		exitItem,
	)

	// Help menu
	aboutItem := fyne.NewMenuItem("About", mw.showAbout)
	helpItem := fyne.NewMenuItem("Help", mw.showHelp)
	contactItem := fyne.NewMenuItem("Contact", mw.showContact)

	helpMenu := fyne.NewMenu("Help",
		aboutItem,
		helpItem,
		contactItem,
	)

	return fyne.NewMainMenu(fileMenu, helpMenu)
}

// Menu action handlers
func (mw *MainWindow) openContestWizard() {
	wizard := NewContestWizard(mw.app)
	wizard.Show()
}

func (mw *MainWindow) openAddParticipants() {
	participantMgr := NewParticipantManager(mw.app)
	participantMgr.Show()
}

func (mw *MainWindow) openChugManager() {
	chugMgr := NewChugManager(mw.app)
	chugMgr.Show()
}

func (mw *MainWindow) openConfiguration() {
	configWindow := NewConfigurationWindow(mw.app)
	configWindow.Show()
}

func (mw *MainWindow) openFinishContest() {
	finishWindow := NewFinishContest(mw.app)
	finishWindow.Show()
}

func (mw *MainWindow) exitApplication() {
	mw.app.Quit()
}

func (mw *MainWindow) showAbout() {
	content := widget.NewRichTextFromMarkdown(`
# About ChugWare2

**Version:** ` + config.AppVersion + `  
**Authors:** ` + config.Authors + `  

ChugWare2 is a comprehensive contest management system for drinking competitions. 

## Features
- Contest setup and configuration
- Participant management  
- Real-time contest timing
- Results tracking and analysis
- Diploma generation
- Multi-discipline support

## Disciplines Supported
- Bottle contests
- Half Tankard contests  
- Full Tankard contests
- Bier Staphette
- Mega Medley
- Team Clash competitions

Built with Go and Fyne for cross-platform compatibility.
`)

	var dialogPopup *widget.PopUp
	dialogPopup = widget.NewModalPopUp(
		container.NewVBox(
			content,
			widget.NewButton("Close", func() {
				dialogPopup.Hide()
			}),
		),
		mw.window.Canvas(),
	)

	dialogPopup.Resize(fyne.NewSize(500, 400))
	dialogPopup.Show()
}

func (mw *MainWindow) showHelp() {
	content := widget.NewRichTextFromMarkdown(`
# ChugWare2 Help

## Getting Started
1. **Contest Wizard** - Set up a new contest
2. **Add Participants** - Register contestants  
3. **Chug Manager** - Run live contests
4. **Configuration** - Adjust settings
5. **Finish Contest** - Generate final results

## Contest Workflow
1. Use Contest Wizard to create contest files and structure
2. Add participants with their details and discipline preferences
3. Use Chug Manager for real-time contest execution
4. Review and finalize results with Finish Contest

## File Management
- All contest data is stored in JSON format
- Automatic directory structure creation
- Configurable file paths and organization

For more information, contact the development team.
`)

	var dialogPopup *widget.PopUp
	dialogPopup = widget.NewModalPopUp(
		container.NewVBox(
			content,
			widget.NewButton("Close", func() {
				dialogPopup.Hide()
			}),
		),
		mw.window.Canvas(),
	)

	dialogPopup.Resize(fyne.NewSize(500, 400))
	dialogPopup.Show()
}

func (mw *MainWindow) showContact() {
	content := widget.NewLabel("Contact: Sigge McKvack, EKAK-2012\nGo Port: AI Assistant")

	var dialogPopup *widget.PopUp
	dialogPopup = widget.NewModalPopUp(
		container.NewVBox(
			content,
			widget.NewButton("Close", func() {
				dialogPopup.Hide()
			}),
		),
		mw.window.Canvas(),
	)

	dialogPopup.Show()
}

// Window returns the underlying fyne.Window so callers can set intercepts, icons, etc.
func (mw *MainWindow) Window() fyne.Window {
	return mw.window
}

// ShowAndRun displays the main window and starts the application
func (mw *MainWindow) ShowAndRun() {
	mw.window.ShowAndRun()
}
