package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"

	"chugware/internal/config"
	"chugware/internal/resources"
	"chugware/internal/ui"
	"chugware/internal/version"
)

// customTheme extends the default theme to make placeholder text more visible
type customTheme struct {
	fyne.Theme
}

func (t *customTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// Make placeholder text more visible (darker grey)
	if name == theme.ColorNamePlaceHolder {
		return color.RGBA{R: 128, G: 128, B: 128, A: 255} // Medium grey
	}
	return t.Theme.Color(name, variant)
}

func main() {
	// Propagate build-stamped version into the config package so all
	// UI components (e.g. the About dialog) show the correct version.
	config.AppVersion = version.Version

	// Initialize application
	myApp := app.NewWithID("com.chugware2.contest")

	// Set the application icon (used in taskbar and – on supported platforms – the tray)
	myApp.SetIcon(resources.CoolSigge)

	// Apply custom theme with more visible placeholder text
	myApp.Settings().SetTheme(&customTheme{Theme: theme.DefaultTheme()})

	// Load configuration
	config.LoadConfig()

	// Initialize external clock manager (used by Configuration and ChugManager)
	ui.GlobalExternalClock = ui.NewExternalClockManager(myApp)

	// Create main window
	mainWindow := ui.NewMainWindow(myApp)

	// Intercept the window close button so it hides to the tray instead of quitting
	mainWindow.Window().SetCloseIntercept(func() {
		mainWindow.Window().Hide()
	})

	// Set up system tray (works on Windows, macOS and Linux desktop environments)
	if desk, ok := myApp.(desktop.App); ok {
		desk.SetSystemTrayIcon(resources.CoolSigge)
		desk.SetSystemTrayMenu(fyne.NewMenu("ChugWare2",
			fyne.NewMenuItem("Open ChugWare2", func() {
				mainWindow.Window().Show()
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Exit", func() {
				myApp.Quit()
			}),
		))
	}

	mainWindow.ShowAndRun()
}
