# ChugWare2 – Operator Instruction Manual

Version 1.0.0 · Last updated 2026-02-22

---

## Table of Contents

1. [Overview](#1-overview)
2. [First-Time Setup](#2-first-time-setup)
3. [Contest Wizard – Creating a New Contest](#3-contest-wizard--creating-a-new-contest)
4. [Add Participants – Managing the Competitor List](#4-add-participants--managing-the-competitor-list)
5. [Chug Manager – Running the Contest](#5-chug-manager--running-the-contest)
   - 5.1 [Opening and Loading Data](#51-opening-and-loading-data)
   - 5.2 [Selecting a Discipline](#52-selecting-a-discipline)
   - 5.3 [Loading a Participant](#53-loading-a-participant)
   - 5.4 [The Timer](#54-the-timer)
   - 5.5 [Recording a Bottle Result](#55-recording-a-bottle-result)
   - 5.6 [Recording Half Tankard and Full Tankard Results](#56-recording-half-tankard-and-full-tankard-results)
   - 5.7 [Recording Team / Relay Discipline Results](#57-recording-team--relay-discipline-results)
   - 5.8 [Entering a Result Manually](#58-entering-a-result-manually)
   - 5.9 [Disqualifying a Participant Outright](#59-disqualifying-a-participant-outright)
   - 5.10 [Skipping a Participant](#510-skipping-a-participant)
   - 5.11 [External Clock Mode](#511-external-clock-mode)
6. [Finish Contest – Viewing and Exporting Results](#6-finish-contest--viewing-and-exporting-results)
7. [Configuration](#7-configuration)
8. [Special Situations](#8-special-situations)
   - 8.1 [Bottle Passed but Disqualified for Overflow](#81-bottle-passed-but-disqualified-for-overflow)
   - 8.2 [Same Participant Competing in Both Bottle and Half Tankard](#82-same-participant-competing-in-both-bottle-and-half-tankard)
   - 8.3 [Bottle DQ (Overflow) Followed by Normal Half-Tankard Attempt](#83-bottle-dq-overflow-followed-by-normal-half-tankard-attempt)
   - 8.4 [Participant Has Multiple Tries Remaining](#84-participant-has-multiple-tries-remaining)
   - 8.5 [Correcting a Wrongly Saved Result](#85-correcting-a-wrongly-saved-result)
9. [Participant Discipline Codes (the "322" format)](#9-participant-discipline-codes-the-322-format)
10. [Time Format Reference](#10-time-format-reference)
11. [Keyboard / Workflow Quick-Reference](#11-keyboard--workflow-quick-reference)
12. [HTML Contest Browser (`htmlgen`)](#12-html-contest-browser-htmlgen)
13. [Installation & Desktop Shortcut](#13-installation--desktop-shortcut)
    - 13.1 [Windows](#131-windows)
    - 13.2 [Linux](#132-linux)
    - 13.3 [macOS](#133-macos)

---

## 1. Overview

ChugWare2 is a contest management system for competitive chugging events. It handles:

- Contest folder and file creation
- Participant registration with per-discipline try counts
- Live in-app timer (or external hardware clock)
- Result recording for six disciplines (Bottle, Half Tankard, Full Tankard, Bier Staphette, Mega Medley, Team Clash)
- Overtime / overflow penalty time addition for Bottle events
- Leaderboard generation, report export, and diploma generation

### Disciplines at a Glance

| Discipline | Default Tries | Notes |
|---|---|---|
| Bottle | 3 | Has overflow penalty time dialog |
| Half Tankard | 2 | Standard pass / DQ flow |
| Full Tankard | 1 | Standard pass / DQ flow |
| Bier Staphette | — | No try count |
| Mega Medley | — | No try count |
| Team Clash | — | No try count |

---

## 2. First-Time Setup

1. Launch `ChugWare2.exe`.
2. The app opens the **Main Menu** with six buttons.
3. Click **Configuration** to verify (or set) the top-level contest folder path. The default is a `ChugWare/` folder next to the executable. See [Section 7](#7-configuration) for details.
4. Once the folder path is correct you are ready to create a contest.

---

## 3. Contest Wizard – Creating a New Contest

**Purpose:** Creates the folder structure and empty data files for one contest.

**Steps:**

1. From the Main Menu click **Contest Wizard**.
2. Fill in:
   - **Contest Name** – e.g. `RegionalChampionship` (spaces are allowed; they are replaced with underscores in the folder name).
   - **Contest Date** – in `YYYY-MM-DD` format. Today's date is pre-filled.
3. Click **Create Contest**.
4. A confirmation dialog shows the path where files were created, e.g.:
   ```
   ChugWare/RegionalChampionship_2026-02-22_Official/
     contest/
       participants.json
       results.json
     results/
     diplomas/
     images/
     template/
   ```
5. Click **Exit** to close the wizard.

> **Important:** You must create a contest before using Add Participants or Chug Manager. If no contest exists the other windows will show a configuration error.

---

## 4. Add Participants – Managing the Competitor List

**Purpose:** Registers competitors, assigns them to disciplines, and sets their try counts.

### 4.1 Adding a New Participant

1. From the Main Menu click **Add Participants**.
2. Fill in the form fields on the left:
   - **Name** – full display name. Must be unique.
   - **Program / Course** – e.g. `Computer Science`.
   - **Team** – team or faction name.
   - **Discipline tries** – a three-digit code (see [Section 9](#9-participant-discipline-codes-the-322-format)). Default `322` means 3 bottle tries, 2 half-tankard tries, 2 full-tankard tries.
3. Click **Add Participant**.

### 4.2 Updating a Participant

1. Click the participant's row in the list to load their data into the form.
2. Edit any fields.
3. Click **Update**.

### 4.3 Deleting a Participant

1. Click the participant's row to select it.
2. Click **Delete**.

> Deleting a participant does **not** remove their already-saved results.

### 4.4 Loading an Existing List

Click **Load from File** to import a JSON file that was prepared outside ChugWare (e.g. a pre-registered list).

### 4.5 Saving

Click **Save All** to persist all changes. The participant list auto-saves when using **Add / Update / Delete**, but manual saves are recommended before closing the window.

---

## 5. Chug Manager – Running the Contest

**Purpose:** The primary screen used on contest day to time each attempt and record results.

### 5.1 Opening and Loading Data

1. From the Main Menu click **Chug Manager**.
2. The window loads the participant list and any existing results automatically from the files configured in the active contest. If required files are missing an error dialog is shown; resolve it in **Configuration** or by running **Contest Wizard** first.

### 5.2 Selecting a Discipline

Use the **Discipline** dropdown (top-left) to choose which discipline is currently being run. Changing the discipline immediately refreshes the **Participants in Discipline** tab to show only participants who still have remaining tries for that discipline.

Disciplines:
- Bottle
- Half Tankard
- Full Tankard
- Bier Staphette
- Mega Medley
- Team Clash

### 5.3 Loading a Participant

Two methods are available:

| Method | How |
|---|---|
| **Load Next Chugger** | Loads the first participant in the Participants in Discipline list (sorted: most tries remaining first, then alphabetically). |
| **Load From List** | Click a row in the **Participants in Discipline** tab to highlight it, then click **Load From List** to load that specific person. |

After loading, the **Current Chugger** card shows the participant's name, program, team, and remaining tries for the selected discipline.

### 5.4 The Timer – Three Timing Methods

ChugWare supports three independent ways to capture a participant's time. Choose the method that fits your setup before the first attempt; you can switch between runs but not mid-attempt.

---

#### Method A – Internal GUI Stopwatch (default)

The built-in software stopwatch. No hardware required.

**Setup:** No setup needed. The **Use External Clock** button must show no checkmark (default state).

**Per-attempt workflow:**

| Step | Button | What happens |
|---|---|---|
| 1 | **Ready Check** | Confirm the participant is ready. Enables the **Start** button. |
| 2 | **Start** | Stopwatch begins counting from `00:00:00.0000`. |
| 3 | **Stop** | Stopwatch freezes. Elapsed time is written into **Base Time** automatically. Try count decremented by 1 and saved. Pass / DQ action buttons are enabled. |
| 4 | *(record result)* | See Section 5.5 / 5.6 for pass and DQ flows. |
| 5 | **Reset** | Clears the timer and result form. Re-enables Load Next Chugger. Use this if you need to restart an attempt before recording a result. |

> **Note:** Reset does **not** decrement a try and does **not** save any result. Only Stop triggers a try decrement.

---

#### Method B – External Hardware Clock

A serial/USB timing device (e.g. a dedicated sports timer, an Arduino emitting time strings, or a PC running minicom piped over a COM port) sends time strings to ChugWare via a serial port. ChugWare parses any `HH:MM:SS` or `HH:MM:SS.mmmm` pattern it receives.

**One-time setup (do this before the contest):**

1. Plug the timing device into the PC.
2. Open **Configuration** from the Main Menu.
3. In the **External Equipment** section:
   - Set **Serial Port** to the port name (e.g. `COM3` on Windows, `/dev/ttyUSB0` on Linux).
   - Set **Baud Rate** to match the device (default `9600`; common values: `4800`, `9600`, `19200`, `115200`).
4. Click **Connect**. The status label changes to `✓ Connected – COM3`.
5. Optionally click **Show Log Window** (accessible from Configuration) to see the raw data stream from the device and verify time strings are being received correctly.
6. Click **Save** in Configuration to persist the port settings.

**Activating external clock in Chug Manager:**

1. Open **Chug Manager**.
2. Click **Use External Clock**. The button label changes to `✓ Using External Clock`.
   - If the device is not connected yet, an error dialog will appear when you press Start. Connect the device in Configuration first.
3. Click it again at any time to switch back to the internal stopwatch (`Use External Clock` label restored).

**Per-attempt workflow:**

| Step | Button | What happens |
|---|---|---|
| 1 | **Ready Check** | Confirms participant ready. Enables Start. |
| 2 | **Start** | ChugWare subscribes to the external clock channel. The timer display updates in real time from the device's transmitted values. |
| 3 | **Stop** | Unsubscribes from the clock channel. The **last received time string** is written into **Base Time** automatically. Try count decremented and saved. Pass / DQ buttons enabled. |
| 4 | *(record result)* | Identical to Method A from this point on. |

> **Tip:** Keep the **External Clock – Live Logs** window open on a second monitor during the contest to verify the device is transmitting correctly. Each received line is shown there in real time.

> **Fallback:** If the device disconnects mid-attempt, the last value received before disconnection is used as the base time when Stop is pressed. If nothing was received, Base Time will be empty and you will need to enter the time manually.

---

#### Method C – Fully Manual Time Entry

Use this when no timer is run at all — for example when reading times from a paper sheet, an external results board, or when correcting a prior contest session.

**Per-attempt workflow:**

1. Load the participant (do **not** start the timer).
2. Optionally click **Ready Check** (required to enable the **Enter Result Manually** button).
3. Fill in the time fields directly:
   - **Base Time** – the raw measured time (e.g. `23.540` for 23.54 seconds, or `1:03.200` for 1 min 3.2 s). See [Section 10](#10-time-format-reference) for accepted formats.
   - **Additional Time** – any penalty time, or leave blank for zero.
   - **Time** – if you already know the final combined time you can enter it here directly and leave Base/Additional blank. The Time field takes priority over the sum of Base + Additional.
4. Set the **Status** dropdown to `Pass` or `Disqualified`.
5. Add a **Comment** if needed.
6. Click **Enter Result Manually**.

> The **Calculate Final Time** button lets you preview `Base + Additional` in the Time field before committing. The Time field locks after Calculate is used — to unlock it, click **Reset**.

---

#### Timing Method Comparison

| | Internal Timer | External Clock | Manual Entry |
|---|---|---|---|
| Hardware needed | None | Serial/USB device | None |
| Setup | None | Configure port + connect in Configuration | None |
| Base Time auto-filled | Yes (at Stop) | Yes (last received value at Stop) | No — you type it |
| Try count decremented | At Stop | At Stop | At Enter Result Manually (if timer never ran) |
| Best for | Most contests | High-accuracy timing hardware | Corrections / paper results |

### 5.5 Recording a Bottle Result

Bottle is the only discipline with a dedicated result dialog because it may carry an **overflow penalty time**.

#### 5.5.1 Clean Bottle (no overflow, no spill)

1. Load participant → Ready Check → Start → Stop.
2. Click **Mark as Pass**.
3. The **Result Entry** dialog opens showing the base time.
4. Click **Clean Bottle (no penalty)**.
5. Result saved as `Pass` with Additional Time = `0`. The participant moves to the next slot.

#### 5.5.2 Bottle Passed but with Overflow Penalty Time

The participant emptied the bottle but overflowed / spilled some liquid. A penalty time must be added.

1. Load participant → Ready Check → Start → Stop.
2. Click **Mark as Pass**.
3. The **Result Entry** dialog opens.
4. Enter the measured penalty time in the **Additional Time** field (e.g. `0:03.500` for 3.5 seconds).
5. Click **Save (Pass)**.
6. Result saved as `Pass`. Final time = Base Time + Additional Time.

#### 5.5.3 Bottle Disqualified (Overflow / Did Not Finish)

The participant overflowed to such an extent that they are disqualified, or they failed to empty the bottle.

1. Load participant → Ready Check → Start → Stop (or Stop whenever the attempt ends).
2. Click **Disqualify + Measure Time**.
3. The **Result Entry** dialog opens with title `Result Entry (Disqualified)`.
4. Optionally enter the measured overflow time in **Additional Time** (for record-keeping).
5. Click **Save (Disqualified)** or **Disqualify (Overflow)** — both save the result as `Disqualified` with comment `Overflow` and set the time to `NaN`.

> A Disqualified result **cannot be overwritten** once saved. If you made a mistake, see [Section 8.5](#85-correcting-a-wrongly-saved-result).

### 5.6 Recording Half Tankard and Full Tankard Results

These disciplines do not have the overflow penalty dialog. The flow is:

1. Select **Half Tankard** (or **Full Tankard**) in the Discipline dropdown.
2. Load participant → Ready Check → Start → Stop.
   - **Base Time** is filled automatically.
3. If there is an additional time (penalty, measurement delay, etc.) enter it in **Additional Time**.
4. Optionally click **Calculate Final Time** to preview `Base + Additional` in the **Time** field.
5. Choose outcome:
   - **Mark as Pass** → saves result as `Pass`, moves to next participant.
   - **Disqualify + Measure Time** → opens the result dialog (same as bottle); saves as `Disqualified`.

### 5.7 Recording Team / Relay Discipline Results

Bier Staphette, Mega Medley, and Team Clash have no individual try counts. The workflow is the same as Half Tankard (Section 5.6) but the Discipline dropdown simply has no try-count filter applied.

### 5.8 Entering a Result Manually

Use this when no timer was run (e.g. result from a separate timing device, or a correction from a paper sheet).

1. Load participant.
2. Optionally click **Ready Check** to enable the **Enter Result Manually** button.
3. Fill in **Base Time** (and optionally **Additional Time** or the final **Time** field directly).
4. Set the **Status** dropdown to `Pass` or `Disqualified`.
5. Add a comment if required.
6. Click **Enter Result Manually**.

Time resolution priority:
1. If **Time** field is filled → used as-is.
2. Else if **Base Time** is filled → Final Time = Base + Additional.
3. Otherwise → error.

### 5.9 Disqualifying a Participant Outright

If the participant did not attempt (no-show, refused to start, etc.):

1. Load participant.
2. Enter `NaN` in the **Time** field (or leave empty and use only the **Disqualify** button).
3. Click **Disqualify** (via the **Enter Result Manually** path with Status = Disqualified), or use **Disqualify + Measure Time** from the status actions.

### 5.10 Skipping a Participant

Use **Skip Participant** when you need to temporarily defer a loaded participant (e.g. they stepped away for a moment). This removes them from the current slot but they remain in the discipline list. After clearing them:

- Click **Load Next Chugger** or **Load From List** to continue with someone else.
- The skipped participant will appear in the list again for the next load.
- **Clear Skipped** resets the internal skip log (does not remove participants from the event).

### 5.11 External Clock – Troubleshooting

| Problem | Likely cause | Fix |
|---|---|---|
| "External clock is not connected" error on Start | Device not connected or port wrong | Open Configuration → set correct port → click Connect |
| Timer display doesn't move after Start | Device connected but not sending data, or baud rate mismatch | Check the Live Logs window; verify baud rate matches device |
| Base Time is empty after Stop | No time string was received before Stop was pressed | Enter the time manually in the Base Time field |
| Time values look garbled in logs | Baud rate mismatch | Adjust baud rate in Configuration to match the device spec |
| Port not listed / can't open COM3 | Driver not installed or port in use by another app | Install device driver; close any other serial terminal (minicom, PuTTY, etc.) |

The external clock parser accepts any line containing a pattern matching `H:MM:SS` or `H:MM:SS.mmmm`. If your device emits lines like `TIME=0:23.540` the value `0:23.540` will be extracted automatically.

---

## 6. Finish Contest – Viewing and Exporting Results

**Purpose:** Review the leaderboard, sort results by different criteria, and generate reports/diplomas.

1. From the Main Menu click **Finish Contest**.
2. Results are loaded automatically and displayed in per-discipline tabs:
   - Bottle · Half Tankard · Full Tankard · Bier Staphette · Mega Medley · Team Clash
3. Use the **Sort / Filter** radio group to change the view:
   | Option | Description |
   |---|---|
   | 🏆 Fastest First | Ascending by final time (DQ at bottom) |
   | 🐢 Slowest First | Descending by final time |
   | ⏳ Most Penalty Time | Most additional time at top |
   | ✨ No Penalty (Clean) | Only results with zero additional time |
   | 🕐 Longest Warm-Up | Highest base time first |
   | ⚡ Quickest Warm-Up | Lowest base time first |
   | 💀 Hall of Shame | Disqualified results only |
4. Click **Refresh** to reload from file (useful if Chug Manager is still recording results in another window).
5. Click **Generate Report** to export a text/JSON summary to the `results/` folder.
6. Click **Export Results** to write a CSV-compatible file.
7. Click **Generate Diplomas** to produce diploma files in the `diplomas/` folder.
8. Click **Save** to persist any pending changes.

---

## 7. Configuration

**Purpose:** Set file paths and hardware device settings that persist across sessions.

Open via Main Menu → **Configuration**.

### 7.1 Path Settings

| Setting | Description |
|---|---|
| **Folder Path** | Top-level ChugWare directory (where contest sub-folders are created). Use the **Browse** button to select. |
| **Current Contest Folder** | Read-only. Shows the active contest path set by Contest Wizard. |

### 7.2 File Settings

| Setting | Description |
|---|---|
| **Participant File** | Path to `participants.json`. Set automatically by Contest Wizard; override here if needed. |
| **Result File** | Path to `results.json`. Set automatically by Contest Wizard. |
| **Template File** | Path to the diploma/report template. |

### 7.3 External Equipment (Serial Clock)

This section controls the hardware timing device used with [Method B – External Hardware Clock](#method-b--external-hardware-clock).

| Setting | Description |
|---|---|
| **Serial Port** | Port name of the timing device. Windows: `COM1`–`COM99`. Linux/Mac: `/dev/ttyUSB0`, `/dev/ttyACM0`, etc. |
| **Baud Rate** | Data rate of the serial connection. Must match the device's setting. Common values: `4800`, `9600`, `19200`, `38400`, `115200`. Default: `9600`. |
| **Connect / Disconnect** | Opens or closes the serial port. Status shows `✓ Connected – COMx` when active. |
| **Show Log Window** | Opens a floating window displaying every raw line received from the device. Use this to verify the device is sending readable time strings before starting the contest. |

**Steps to connect an external clock:**

1. Plug in the device.
2. Enter the port name and baud rate.
3. Click **Connect**. The status label turns green (`✓ Connected`).
4. Click **Show Log Window** and confirm time strings appear (e.g. `0:00:03.540`).
5. Click **Save** to persist the settings.
6. In Chug Manager, click **Use External Clock** to activate it for the session.

> The connection persists across window switches (Configuration → Chug Manager) as long as the app is running. If you close and reopen ChugWare you must reconnect.

### 7.4 Saving Configuration

Click **Save** at any time. Settings are written to `~/.chugware/chugware_config.json` and take effect immediately without restarting the app.

---

## 8. Special Situations

### 8.1 Bottle Passed but Disqualified for Overflow

**Scenario:** The participant emptied the bottle (the liquid was consumed) but overflowed severely and the judges rule it a disqualification.

This is **not** a clean pass with penalty — it is a full disqualification. Procedure:

1. Load participant → Ready Check → Start → Stop.
2. Click **Disqualify + Measure Time** (do NOT click Mark as Pass).
3. In the Result Entry dialog, enter the overflow amount in **Additional Time** if required for your records.
4. Click **Disqualify (Overflow)**.
5. Result is saved: Status = `Disqualified`, Comment = `Overflow`, Time = `NaN`.

The participant's bottle try count is decremented at Stop. If tries remain, they may attempt again in the next bottle round.

### 8.2 Same Participant Competing in Both Bottle and Half Tankard

Each discipline is tracked independently. A Bottle result (pass or DQ) has **no effect** on the participant's Half Tankard try count or results. The two disciplines are run as separate events, usually at different times during the contest.

### 8.3 Bottle DQ (Overflow) Followed by Normal Half-Tankard Attempt

This is the most common dual-event scenario. Step-by-step:

#### Step A – Record the Bottle DQ

1. In the **Discipline** dropdown select **Bottle**.
2. Load the participant → time the attempt → click **Stop**.
3. Click **Disqualify + Measure Time** or **Disqualify** as appropriate.
4. Save as Disqualified (see Section 8.1). The result is written to `results.json` with discipline `Bottle` and status `Disqualified`.

#### Step B – Run the Half Tankard Attempt Normally

1. Change the **Discipline** dropdown to **Half Tankard**.
2. The Participants in Discipline list refreshes. The participant appears if their `HalfTankard` try count is still > 0 (it is unaffected by the Bottle result).
3. Load the participant using **Load Next Chugger** or **Load From List**.
4. Perform the attempt normally:
   - Ready Check → Start → Stop.
   - Enter additional time if any.
   - Click **Mark as Pass** or **Disqualify + Measure Time** as appropriate.
5. The Half Tankard result is written independently to `results.json` with discipline `Half Tankard`.

#### Summary Table

| Event Stage | Action | Result saved |
|---|---|---|
| Bottle attempt | DQ + Overflow | `Discipline: Bottle`, `Status: Disqualified`, `Time: NaN` |
| Half Tankard attempt | Pass (clean, 00:23.540) | `Discipline: Half Tankard`, `Status: Pass`, `Time: 00:00:23.5400` |

Both results appear independently in the Finish Contest leaderboards under their respective discipline tabs.

### 8.4 Participant Has Multiple Tries Remaining

**Bottle (3 tries by default):** Each time **Stop** is clicked, one try is deducted and saved. The participant remains in the Participants in Discipline list until their try count reaches 0. Each attempt is recorded as a separate result row. In Finish Contest the best (fastest passing) result is what counts for the leaderboard.

> If you reset the timer with **Reset** instead of stopping it, no try is deducted and no result is written.

### 8.5 Correcting a Wrongly Saved Result

ChugWare **does not allow overwriting a Disqualified result** through the UI as a safety measure. If a result was saved incorrectly:

1. Navigate to the contest folder:
   ```
   ChugWare/<ContestName>_<Date>_Official/contest/results.json
   ```
2. Open `results.json` in a text editor.
3. Locate the incorrect entry by participant name and discipline.
4. Edit or remove the entry.
5. Save the file.
6. In ChugWare click **Refresh** in Finish Contest (or re-open Chug Manager) to reload the corrected data.

---

## 9. Participant Discipline Codes (the "322" format)

When registering a participant in **Add Participants**, the **Discipline Tries** field accepts a three-character string:

| Position | Digit | Meaning |
|---|---|---|
| 1st | 0–9 | Bottle tries |
| 2nd | 0–9 | Half Tankard tries |
| 3rd | 0–9 | Full Tankard tries |

**Examples:**

| Code | Bottle | Half Tankard | Full Tankard |
|---|---|---|---|
| `322` | 3 | 2 | 2 |
| `100` | 1 | 0 | 0 |
| `310` | 3 | 1 | 0 |
| `000` | Not eligible for any of the three | — | — |

Participants with 0 tries in a discipline are excluded from that discipline's list in Chug Manager.

Team / relay disciplines (Bier Staphette, Mega Medley, Team Clash) are **not** counted here — all participants are eligible for those events regardless of the code.

---

## 10. Time Format Reference

ChugWare accepts several time input formats and normalises them all to `HH:MM:SS.mmmm` internally (where `mmmm` is tenths-of-milliseconds, giving 0.1 ms precision).

| Input example | Interpreted as |
|---|---|
| `3` | 3 seconds → `00:00:03.0000` |
| `3.5` | 3.5 seconds → `00:00:03.5000` |
| `3.500` | 3.5 seconds → `00:00:03.5000` |
| `00:03` | 3 seconds (MM:SS) → `00:00:03.0000` |
| `1:03` | 1 minute 3 seconds → `00:01:03.0000` |
| `0:01:03.500` | 1 minute 3.5 seconds → `00:01:03.5000` |
| `NaN` | Disqualified / no valid time |

The **Calculate Final Time** button in Chug Manager can be used to verify the sum of Base Time + Additional Time before committing a result.

---

## 11. Keyboard / Workflow Quick-Reference

### Standard Per-Participant Flow

```
Select Discipline
      ↓
Load Next Chugger (or Load From List)
      ↓
Ready Check  →  Start  →  [participant competes]  →  Stop
      ↓
Enter Additional Time (if any)
      ↓
Mark as Pass  ──or──  Disqualify + Measure Time
      ↓
(Automatic) Save result, decrement tries, move to next participant
```

### Bottle-Specific Flow

```
Stop timer
    ↓
Mark as Pass
    ├── Clean Bottle (no penalty)      → Pass, Additional = 0
    ├── Save (Pass) with extra time    → Pass, Final = Base + Additional
    └── Disqualify (Overflow)          → DQ, Time = NaN, Comment = Overflow
         ──or──
Disqualify + Measure Time
    └── Save (Disqualified)            → DQ, Time = NaN, Comment = Overflow
```

### Contest Day Checklist

- [ ] Launch ChugWare
- [ ] Verify Configuration (folder path, result file)
- [ ] Run Contest Wizard (if new contest)
- [ ] Add all participants (Add Participants)
- [ ] Open Chug Manager
- [ ] Select first discipline
- [ ] Work through all participants, recording results
- [ ] Switch disciplines and repeat
- [ ] Open Finish Contest to review leaderboard
- [ ] Generate Report and Diplomas
- [ ] Run `htmlgen` to publish the web results browser

---

## 12. HTML Contest Browser (`htmlgen`)

`htmlgen` is a standalone command-line tool (built alongside ChugWare2) that scans every contest folder in your ChugWare root directory and produces a single **self-contained HTML file** you can open in any browser — no internet connection required.

### 12.1 What it generates

- **Overview page** – cards for every contest showing date, official/unofficial status, total athletes, passes, and DQs.
- **Per-contest page** – discipline tabs (Bottle, Half Tankard, Full Tankard, Bier Staphette, Mega Medley, Team Clash) each showing a ranked results table with medal icons (🥇🥈🥉) for top 3, colour-coded Pass/DQ pills, base time, and penalty time columns.
- **Athletes panel** – every registered participant with their try counts.
- **Sidebar navigation** – jump instantly between contests.

### 12.2 Building ChugWare2 and `htmlgen`

Use the provided `build.ps1` script (PowerShell) to build both executables with version metadata stamped in:

```powershell
.\build.ps1                    # builds 1.0.0 (default)
.\build.ps1 -Version "1.1.0"  # override version
```

The script produces two executables and prints the stamped version, build date, and git commit:

```
=== ChugWare2 Build ===
  Version   : 1.0.0
  BuildDate : 2026-02-22
  GitCommit : abc1234

Building ChugWare2.exe ...
  -> ChugWare2.exe
Building htmlgen.exe ...
  -> htmlgen.exe
```

The version is injected at link time via Go's `-ldflags` mechanism:
```
-X chugware/internal/version.Version=1.0.0
-X chugware/internal/version.BuildDate=2026-02-22
-X chugware/internal/version.GitCommit=abc1234
```

To build manually without the script:
```powershell
go build -ldflags "-X chugware/internal/version.Version=1.0.0" -o ChugWare2.exe ./cmd/
go build -ldflags "-X chugware/internal/version.Version=1.0.0" -o htmlgen.exe ./cmd/htmlgen/
```

The stamped version appears in **Help → About** inside the application.

### 12.3 Running `htmlgen`

```
htmlgen.exe [--root <ChugWare folder>] [--out <output file>]
```

| Flag | Default | Description |
|---|---|---|
| `--root` | `ChugWare` | Path to the ChugWare contests folder (relative or absolute). |
| `--out` | `chugware_results.html` | Output HTML file path. |

**Examples:**

```powershell
# Run from the project folder – uses the default ChugWare/ subfolder
.\htmlgen.exe

# Explicit paths
.\htmlgen.exe --root C:\Contests\ChugWare --out C:\Shared\results.html

# Linux / macOS
./htmlgen --root ~/ChugWare --out ~/Desktop/results.html
```

Output:
```
Scanning contests in: C:\...\ChugWare
Found 3 contest(s)
HTML written to: C:\...\chugware_results.html
```

### 12.4 Viewing the results

Double-click `chugware_results.html` (or the path you chose) to open it in your default browser. The file is fully self-contained — all CSS and JavaScript are embedded — so you can:

- Copy it to a USB stick and open it on any PC.
- Email it or upload it to a shared drive.
- Drop it on a web server and share the URL.

### 12.5 Keeping results up to date

Re-run `htmlgen.exe` at any point during or after the contest to regenerate the file with the latest data. The tool always reads directly from the `results.json` files on disk, so it reflects whatever ChugWare has saved most recently.

### 12.6 Contest folder detection

`htmlgen` recognises sub-folders that follow the naming convention created by Contest Wizard:

```
<ContestName>_<YYYY-MM-DD>_Official
<ContestName>_<YYYY-MM-DD>_Unofficial
```

Folders with any other naming pattern are skipped with a warning printed to the terminal. Contests are sorted newest-first in the browser.

### 12.7 Results ordering

Within each discipline tab:
- **Pass** results are ranked by final time, fastest first (rank 1 = winner).
- **DQ** results appear at the bottom of the table, unranked.
- Top 3 places receive 🥇 🥈 🥉 medal icons.

---

## 13. Installation & Desktop Shortcut

After building with `build.ps1` you have two executables:

| File | Purpose |
|---|---|
| `ChugWare2.exe` | Main GUI application |
| `htmlgen.exe` | HTML contest browser generator |

Copy **both** files together with the `ChugWare/` data folder to the target machine. The app stores its config in `~/.chugware/chugware_config.json` and its contest data wherever the **Folder Path** setting points (default: a `ChugWare/` sub-folder next to the executable).

---

### 13.1 Windows

#### Install

1. Create an installation folder, e.g. `C:\Program Files\ChugWare2\`.
2. Copy `ChugWare2.exe`, `htmlgen.exe`, and the `ChugWare\` data folder into it.
3. (Optional) Copy any contest images / diploma templates into the matching sub-folders.

#### Desktop shortcut (GUI)

1. Open **File Explorer** and navigate to `C:\Program Files\ChugWare2\`.
2. Right-click `ChugWare2.exe` → **Send to** → **Desktop (create shortcut)**.
3. The shortcut appears on the Desktop. Right-click it → **Properties** to:
   - Change the **icon** (point it to `ChugWare2.exe` and pick the embedded icon).
   - Set **Start in** to `C:\Program Files\ChugWare2\` (important — the app looks for the `ChugWare\` folder relative to the working directory).

#### Desktop shortcut (PowerShell)

```powershell
$WshShell = New-Object -ComObject WScript.Shell
$Shortcut = $WshShell.CreateShortcut("$env:USERPROFILE\Desktop\ChugWare2.lnk")
$Shortcut.TargetPath   = "C:\Program Files\ChugWare2\ChugWare2.exe"
$Shortcut.WorkingDirectory = "C:\Program Files\ChugWare2"
$Shortcut.Description  = "ChugWare2 Contest Management"
$Shortcut.Save()
```

Run the snippet once in PowerShell — the shortcut is created immediately.

#### Windows Defender / SmartScreen

Because the binary is unsigned, Windows may show a **SmartScreen** warning on first launch. Click **More info → Run anyway**. To suppress this for all users on the machine, right-click `ChugWare2.exe` → **Properties** → tick **Unblock** → **OK**.

---

### 13.2 Linux

#### Install

```bash
# Create install directory
sudo mkdir -p /opt/chugware2

# Copy binaries (built on Linux: no .exe extension)
sudo cp ChugWare2 htmlgen /opt/chugware2/
sudo chmod +x /opt/chugware2/ChugWare2 /opt/chugware2/htmlgen

# Copy existing contest data if migrating
sudo cp -r ChugWare /opt/chugware2/
```

> **Note:** Fyne requires a display server. On a headless server the app will not run. On a desktop Linux install (X11 or Wayland) no extra dependencies are needed beyond standard system libraries.

#### Desktop shortcut (.desktop file)

Create a `.desktop` launcher so ChugWare2 appears in your application menu and can be pinned to the desktop:

```bash
cat > ~/.local/share/applications/chugware2.desktop << 'EOF'
[Desktop Entry]
Version=1.0
Type=Application
Name=ChugWare2
Comment=Contest Management System
Exec=/opt/chugware2/ChugWare2
Path=/opt/chugware2
Icon=/opt/chugware2/ChugWare2
Terminal=false
Categories=Utility;
EOF

# Make it executable
chmod +x ~/.local/share/applications/chugware2.desktop

# Refresh the application database
update-desktop-database ~/.local/share/applications/ 2>/dev/null || true
```

To also place a copy on the **Desktop**:

```bash
cp ~/.local/share/applications/chugware2.desktop ~/Desktop/
chmod +x ~/Desktop/chugware2.desktop
```

On GNOME you may need to right-click the desktop file → **Allow Launching** the first time.

---

### 13.3 macOS

#### Install

macOS requires a window manager and uses `.app` bundles, but since Fyne produces a plain Unix binary you can run it directly from a terminal or wrap it in a minimal app bundle.

**Simple install (terminal only):**

```bash
mkdir -p ~/Applications/ChugWare2
cp ChugWare2 htmlgen ~/Applications/ChugWare2/
chmod +x ~/Applications/ChugWare2/ChugWare2
cp -r ChugWare ~/Applications/ChugWare2/
```

**Minimal `.app` bundle** (enables Dock icon and Spotlight):

```bash
APP=~/Applications/ChugWare2.app
mkdir -p "$APP/Contents/MacOS"
mkdir -p "$APP/Contents/Resources"

# Binary
cp ChugWare2 "$APP/Contents/MacOS/ChugWare2"
chmod +x "$APP/Contents/MacOS/ChugWare2"

# Minimal Info.plist
cat > "$APP/Contents/Info.plist" << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
  "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
  <key>CFBundleName</key>        <string>ChugWare2</string>
  <key>CFBundleExecutable</key>  <string>ChugWare2</string>
  <key>CFBundleIdentifier</key>  <string>com.chugware2.contest</string>
  <key>CFBundleVersion</key>     <string>1.0.0</string>
  <key>CFBundlePackageType</key> <string>APPL</string>
  <key>NSHighResolutionCapable</key><true/>
</dict></plist>
EOF
```

#### Desktop shortcut / Dock pin

1. Open **Finder**, navigate to **Applications**.
2. Double-click `ChugWare2.app` to launch it.
3. The icon appears in the Dock while running — right-click it → **Options** → **Keep in Dock**.

> **Gatekeeper warning:** On first launch macOS may block an unsigned binary. Go to **System Settings → Privacy & Security → Security** and click **Open Anyway** next to the ChugWare2 entry. Alternatively run once with: `xattr -cr ~/Applications/ChugWare2.app && open ~/Applications/ChugWare2.app`

#### Working directory

When launched from the Dock or Finder the working directory defaults to `~`. Set it explicitly in the `.app` wrapper or ensure your **Folder Path** in Configuration uses an absolute path so the `ChugWare/` data folder is always found correctly.
