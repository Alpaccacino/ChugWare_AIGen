# ChugWare - Contest Management System

A comprehensive contest management system for drinking competitions, ported from VB.NET to Go using the Fyne framework.

## Features

- **Contest Wizard**: Set up new contests with customizable disciplines and configurations
- **Participant Management**: Register contestants with their details and discipline preferences  
- **Real-time Contest Execution**: Run live contests with precise timing and scoring
- **Results Tracking**: Record and analyze contest results with multiple status types
- **Configuration Management**: Customize application settings and file paths
- **Contest Finalization**: Generate reports, leaderboards, and diplomas

## Supported Disciplines

- Bottle contests
- Half Tankard contests  
- Full Tankard contests
- Bier Staphette
- Mega Medley
- Team Clash competitions

## Installation

### Prerequisites

- Go 1.21 or later
- C compiler (for Fyne dependencies on some platforms)

### Building from Source

1. Clone or download the source code
2. Navigate to the project directory
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Build the application:
   ```bash
   go build -o chugware.exe ./cmd/main.go
   ```
5. Run the application:
   ```bash
   ./chugware.exe
   ```

## Usage

### Getting Started

1. **Contest Setup**: Use the Contest Wizard to create a new contest structure
2. **Add Participants**: Register contestants with their preferred disciplines
3. **Run Contests**: Use the Chug Manager for real-time contest execution
4. **Finalize Results**: Generate reports and leaderboards

### Contest Workflow

1. **Contest Wizard**
   - Set contest name, date, and type (Official/Unofficial)
   - Select disciplines to include
   - Configure trial settings
   - Generate directory structure and files

2. **Participant Management**
   - Add participant details (name, program, team)
   - Set discipline attempt counts
   - View and manage participant lists
   - Monitor results in real-time

3. **Chug Manager**
   - Load participants for each discipline
   - Run contests with precision timing
   - Record results with status (Pass/Disqualified/Fail)
   - Handle skipped participants and queue management

4. **Configuration**
   - Set file paths and directories
   - Configure trial limits and application settings
   - Manage persistent settings

5. **Finish Contest**
   - Review all contest results
   - Generate leaderboards by discipline
   - Export results and reports
   - Create diploma data for winners

## Data Management

### File Structure

The application creates the following directory structure for each contest:

```
Contest_Name_YYYY-MM-DD_Official/
├── contest/          # Contest data files (JSON)
├── results/          # Final results and exports
├── diplomas/         # Diploma generation data
├── images/           # Contest images
└── template/         # Template files
```

### Data Format

All contest data is stored in JSON format for easy manipulation and backup:

- **Participants**: Name, program, team, discipline attempts
- **Results**: Name, discipline, timing, status, comments
- **Configuration**: File paths, settings, preferences

## Technical Details

### Architecture

- **Frontend**: Fyne v2 for cross-platform GUI
- **Data Layer**: JSON file-based persistence
- **Timing**: High-precision contest timing with millisecond accuracy
- **Configuration**: JSON-based settings management

### Key Components

- `cmd/main.go`: Application entry point
- `internal/config/`: Configuration management
- `internal/models/`: Data structures
- `internal/utils/`: Utility functions (file ops, time parsing, validation)
- `internal/data/`: Data management layers
- `internal/ui/`: User interface components

### Dependencies

- **Fyne v2**: Cross-platform GUI framework
- **Go Standard Library**: Core functionality

## Original Application

This is a Go port of the original ChugWare VB.NET application by Sigge McKvack, EKAK-2012. The Go version maintains full compatibility with the original's data format and functionality while adding cross-platform support and modern UI components.

## Features Comparison

| Feature | Original VB.NET | Go Port |
|---------|----------------|---------|
| Contest Management | ✅ | ✅ |
| Participant Tracking | ✅ | ✅ |
| Real-time Timing | ✅ | ✅ |
| Results Export | ✅ | ✅ |
| Auto-refresh | ✅ | ✅ |
| Cross-platform | ❌ | ✅ |
| Modern UI | ❌ | ✅ |
| JSON Data Format | ✅ | ✅ |

## License

This port maintains compatibility with the original ChugWare application while being implemented in Go for cross-platform support.

## Contributing

This is a faithful port of the original VB.NET application. When reporting issues or suggesting features, please consider compatibility with the original application's data format and workflow.