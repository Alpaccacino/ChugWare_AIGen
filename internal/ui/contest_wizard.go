package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"chugware/internal/config"
	"chugware/internal/utils"
)

type ContestWizardWindow struct {
	app    fyne.App
	window fyne.Window
}

func NewContestWizard(app fyne.App) *ContestWizardWindow {
	w := &ContestWizardWindow{
		app:    app,
		window: app.NewWindow("Contest Wizard"),
	}

	w.setupUI()
	return w
}

func (w *ContestWizardWindow) setupUI() {
	title := widget.NewLabel("Contest Wizard")
	title.Alignment = fyne.TextAlignCenter

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Enter Contest Name")
	dateEntry := widget.NewEntry()
	dateEntry.SetPlaceHolder("Enter Contest Date (YYYY-MM-DD)")

	// Set today's date as default
	nameEntry.SetText("Contest")
	dateEntry.SetText(time.Now().Format("2006-01-02"))

	createButton := widget.NewButton("Create Contest", func() {
		contestName := strings.TrimSpace(nameEntry.Text)
		contestDate := strings.TrimSpace(dateEntry.Text)

		if contestName == "" {
			dialog.ShowError(fmt.Errorf("contest name cannot be empty"), w.window)
			return
		}

		if contestDate == "" {
			dialog.ShowError(fmt.Errorf("contest date cannot be empty"), w.window)
			return
		}

		if err := w.createContest(contestName, contestDate); err != nil {
			dialog.ShowError(fmt.Errorf("error creating contest: %w", err), w.window)
			return
		}

		dialog.ShowInformation("Contest Created",
			fmt.Sprintf("Contest '%s' scheduled for %s has been created.\n\nFiles created in: %s",
				contestName, contestDate, config.Settings.FolderPathContestNameAndDate),
			w.window)
	})

	exitButton := widget.NewButton("Exit", func() {
		w.window.Close()
	})

	form := container.NewVBox(
		title,
		widget.NewSeparator(),
		widget.NewLabel("Contest Name:"),
		nameEntry,
		widget.NewLabel("Contest Date:"),
		dateEntry,
		widget.NewSeparator(),
		container.NewHBox(createButton, exitButton),
	)

	content := container.NewPadded(form)
	w.window.SetContent(content)
	w.window.Resize(fyne.NewSize(400, 300))
	w.window.CenterOnScreen()
}

func (w *ContestWizardWindow) createContest(contestName, contestDate string) error {
	// Create contest folder name: Contest_Name_YYYY-MM-DD_Official
	folderName := fmt.Sprintf("%s_%s_%s",
		strings.ReplaceAll(contestName, " ", "_"),
		contestDate,
		config.OfficialKey)

	// Set up contest directory path
	contestPath := filepath.Join(config.Settings.FolderPath, folderName)
	config.Settings.FolderPathContestNameAndDate = contestPath

	// Create main contest directory
	if err := os.MkdirAll(contestPath, 0755); err != nil {
		return fmt.Errorf("failed to create contest directory: %w", err)
	}

	// Create subdirectories
	subdirs := []string{
		config.ContestDirectory,
		config.ResultsDirectory,
		config.DiplomasDirectory,
		config.ImagesDirectory,
		config.TemplateDirectory,
	}

	for _, subdir := range subdirs {
		dirPath := filepath.Join(contestPath, subdir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create %s directory: %w", subdir, err)
		}
	}

	// Create default data files
	participantFile := filepath.Join(contestPath, config.ContestDirectory, "participants.json")
	resultFile := filepath.Join(contestPath, config.ContestDirectory, "results.json")

	// Create empty participant file
	if err := utils.SaveListToJSONFile(participantFile, []map[string]string{}); err != nil {
		return fmt.Errorf("failed to create participants file: %w", err)
	}

	// Create empty results file
	if err := utils.SaveListToJSONFile(resultFile, []map[string]string{}); err != nil {
		return fmt.Errorf("failed to create results file: %w", err)
	}

	// Update configuration with new file paths
	config.Settings.ParticipantFile = participantFile
	config.Settings.ResultFile = resultFile

	// Save updated configuration
	return config.SaveConfig()
}

func (w *ContestWizardWindow) Show() {
	w.window.Show()
}
