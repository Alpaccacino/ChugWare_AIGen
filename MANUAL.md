# ChugWare – Operator Instruction Manual

Version 2.0 · Last updated 2026-02-22

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

---

## 1. Overview

ChugWare is a contest management system for competitive chugging events. It handles:

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

1. Launch `chugware.exe`.
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

### 5.4 The Timer

| Button | Action |
|---|---|
| **Ready Check** | Confirms the participant is ready. Enables the Start button. |
| **Start** | Starts the stopwatch from zero. |
| **Stop** | Stops the stopwatch and locks the elapsed time into the **Base Time** field. Also decrements the participant's try count by 1 and saves it. Enables Pass / DQ buttons. |
| **Reset** | Resets the timer and clears the result form. Re-enables Load Next Chugger. |

> After clicking **Stop** the Base Time field is filled automatically with the stopped time. Do not edit it unless you have a correction to make.

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

### 5.11 External Clock Mode

If a hardware timing device is connected via serial port:

1. Configure the serial port and baud rate in **Configuration** (see Section 7).
2. In Chug Manager click **Use External Clock**. The button label toggles to **Use Internal Clock** to indicate external mode is active.
3. The timer display will update from the hardware clock signal instead of the internal stopwatch.
4. Stop/reset behaviour is identical to internal clock mode.

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

| Setting | Description |
|---|---|
| **Folder Path** | Top-level ChugWare directory (where contest sub-folders are created). |
| **Participant File** | Path to the active `participants.json`. Set automatically by Contest Wizard. |
| **Result File** | Path to the active `results.json`. Set automatically by Contest Wizard. |
| **External Clock Port** | Serial port name (e.g. `COM3` on Windows, `/dev/ttyUSB0` on Linux). |
| **External Clock Baud** | Baud rate of the serial clock (e.g. `9600`). |

Changes are saved to `~/.chugware/chugware_config.json` and take effect immediately.

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
