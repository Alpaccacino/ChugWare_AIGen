package ui

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"go.bug.st/serial"
)

// GlobalExternalClock is the shared external clock instance used by both
// Configuration and ChugManager.
var GlobalExternalClock *ExternalClockManager

// timeLineRegex finds HH:MM:SS or HH:MM:SS.mmmm patterns anywhere in a log line.
var timeLineRegex = regexp.MustCompile(`\b(\d{1,2}:\d{2}:\d{2}(?:\.\d+)?)\b`)

// ExternalClockManager reads time data from a serial/USB device (e.g. via minicom
// or a direct COM port connection) and distributes parsed time strings to
// registered subscribers.
type ExternalClockManager struct {
	app fyne.App
	mu  sync.Mutex

	// Config (saved in ContestSettings by the Configuration window)
	portName string
	baudRate int

	// Runtime state
	connected bool
	stopChan  chan struct{}
	port      serial.Port

	// TimeChan delivers the most-recently parsed time string to the active
	// ChugManager.  It is a buffered channel; old values are discarded when
	// the consumer is busy so the UI never blocks the reader goroutine.
	TimeChan chan string

	// Logging
	logLines []string

	// UI – may be nil if the log window has not been opened yet.
	logEntry    *widget.Entry
	statusLabel *widget.Label
	connectBtn  *widget.Button
	logWindow   fyne.Window
}

// NewExternalClockManager creates the singleton manager.
func NewExternalClockManager(app fyne.App) *ExternalClockManager {
	return &ExternalClockManager{
		app:      app,
		baudRate: 9600,
		TimeChan: make(chan string, 256),
	}
}

// Connect opens the serial port and starts reading.
func (ecm *ExternalClockManager) Connect(portName string, baud int) error {
	ecm.mu.Lock()
	defer ecm.mu.Unlock()

	if ecm.connected {
		return fmt.Errorf("already connected – disconnect first")
	}
	if portName == "" {
		return fmt.Errorf("port name is required (e.g. COM3)")
	}
	if baud <= 0 {
		baud = 9600
	}

	mode := &serial.Mode{BaudRate: baud}
	p, err := serial.Open(portName, mode)
	if err != nil {
		return fmt.Errorf("cannot open %s: %w", portName, err)
	}

	ecm.portName = portName
	ecm.baudRate = baud
	ecm.port = p
	ecm.stopChan = make(chan struct{})
	ecm.connected = true
	ecm.updateStatusUI()

	go ecm.readLoop(p)
	return nil
}

// Disconnect closes the serial port.
func (ecm *ExternalClockManager) Disconnect() {
	ecm.mu.Lock()
	defer ecm.mu.Unlock()

	if !ecm.connected {
		return
	}
	close(ecm.stopChan)
	ecm.port.Close()
	ecm.port = nil
	ecm.connected = false
	ecm.updateStatusUI()
}

// SetBaud updates the preferred baud rate (takes effect on next Connect).
func (ecm *ExternalClockManager) SetBaud(baud int) {
	ecm.mu.Lock()
	defer ecm.mu.Unlock()
	if baud > 0 {
		ecm.baudRate = baud
	}
}

// IsConnected returns whether the port is open.
func (ecm *ExternalClockManager) IsConnected() bool {
	ecm.mu.Lock()
	defer ecm.mu.Unlock()
	return ecm.connected
}

// readLoop runs in a goroutine, reading lines from the serial port.
func (ecm *ExternalClockManager) readLoop(p serial.Port) {
	reader := bufio.NewReader(p)
	for {
		// Honour stop signal without blocking
		select {
		case <-ecm.stopChan:
			return
		default:
		}

		// Non-blocking 100 ms read window so we can check stopChan regularly.
		p.SetReadTimeout(100 * time.Millisecond)

		line, err := reader.ReadString('\n')
		if err != nil {
			select {
			case <-ecm.stopChan:
				return
			default:
				if line == "" {
					continue
				}
				// Partial line – process what we got
			}
		}

		line = strings.TrimRight(line, "\r\n ")
		if line == "" {
			continue
		}

		ecm.appendLog(line)

		if t := ecm.parseTime(line); t != "" {
			// Drain old values so the consumer always gets the latest
			select {
			case <-ecm.TimeChan:
			default:
			}
			ecm.TimeChan <- t
		}
	}
}

// parseTime extracts the first time value from a log line.
func (ecm *ExternalClockManager) parseTime(line string) string {
	m := timeLineRegex.FindStringSubmatch(line)
	if len(m) >= 2 {
		return m[1]
	}
	return ""
}

// appendLog adds a line to the log buffer and refreshes the UI log entry.
func (ecm *ExternalClockManager) appendLog(line string) {
	ecm.mu.Lock()
	ecm.logLines = append(ecm.logLines, line)
	if len(ecm.logLines) > 2000 {
		ecm.logLines = ecm.logLines[len(ecm.logLines)-2000:]
	}
	combined := strings.Join(ecm.logLines, "\n")
	ecm.mu.Unlock()

	if ecm.logEntry != nil {
		ecm.logEntry.SetText(combined)
	}
}

// updateStatusUI refreshes the status label and connect button text.
// Must be called with mu held.
func (ecm *ExternalClockManager) updateStatusUI() {
	if ecm.statusLabel != nil {
		if ecm.connected {
			ecm.statusLabel.SetText("✓ Connected – " + ecm.portName)
		} else {
			ecm.statusLabel.SetText("✗ Not connected")
		}
	}
	if ecm.connectBtn != nil {
		if ecm.connected {
			ecm.connectBtn.SetText("Disconnect")
		} else {
			ecm.connectBtn.SetText("Connect")
		}
	}
}

// ShowLogWindow opens (or focuses) the floating log window.
func (ecm *ExternalClockManager) ShowLogWindow() {
	if ecm.logWindow != nil {
		ecm.logWindow.RequestFocus()
		ecm.logWindow.Show()
		return
	}

	ecm.logWindow = ecm.app.NewWindow("External Clock – Live Logs")
	ecm.logWindow.Resize(fyne.NewSize(700, 450))

	if ecm.logEntry == nil {
		ecm.logEntry = widget.NewMultiLineEntry()
		ecm.logEntry.Wrapping = fyne.TextWrapWord
		// Populate with any lines already received
		ecm.mu.Lock()
		ecm.logEntry.SetText(strings.Join(ecm.logLines, "\n"))
		ecm.mu.Unlock()
	}
	// Make it read-only via OnChanged guard
	savedText := ""
	ecm.logEntry.OnChanged = func(s string) {
		if s != savedText {
			ecm.logEntry.SetText(savedText)
		}
	}

	clearBtn := widget.NewButton("Clear Logs", func() {
		ecm.mu.Lock()
		ecm.logLines = nil
		ecm.mu.Unlock()
		savedText = ""
		ecm.logEntry.SetText("")
	})

	ecm.logWindow.SetContent(container.NewBorder(nil, clearBtn, nil, nil,
		container.NewScroll(ecm.logEntry),
	))
	ecm.logWindow.SetOnClosed(func() {
		ecm.logWindow = nil
		ecm.logEntry = nil
	})
	ecm.logWindow.Show()
}

// BuildStatusWidget returns a small inline status label suitable for embedding
// in the Configuration layout.  The label is kept in sync via updateStatusUI.
func (ecm *ExternalClockManager) BuildStatusLabel() *widget.Label {
	ecm.mu.Lock()
	defer ecm.mu.Unlock()

	if ecm.statusLabel == nil {
		ecm.statusLabel = widget.NewLabel("✗ Not connected")
	}
	return ecm.statusLabel
}

// BuildConnectButton returns the connect/disconnect button, wired up to the
// provided port and baud getters so Configuration can pass live field values.
func (ecm *ExternalClockManager) BuildConnectButton(getPort func() string, getBaud func() int) *widget.Button {
	ecm.mu.Lock()
	defer ecm.mu.Unlock()

	ecm.connectBtn = widget.NewButton("Connect", func() {
		if ecm.IsConnected() {
			ecm.Disconnect()
		} else {
			port := getPort()
			baud := getBaud()
			if err := ecm.Connect(port, baud); err != nil {
				// Show error in status label directly (window reference not available here)
				ecm.mu.Lock()
				if ecm.statusLabel != nil {
					ecm.statusLabel.SetText("✗ " + err.Error())
				}
				ecm.mu.Unlock()
			}
		}
	})
	return ecm.connectBtn
}
