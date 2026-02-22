// htmlgen – ChugWare contest browser generator.
//
// Usage:
//
//	htmlgen [--root <ChugWare folder>] [--out <output.html>]
//
// Scans every sub-folder inside the ChugWare root directory, reads
// contest/participants.json and contest/results.json, then writes a
// single self-contained HTML file you can open in any browser.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// ─── data types ──────────────────────────────────────────────────────────────

type Participant struct {
	Name        string `json:"name"`
	Program     string `json:"program"`
	Team        string `json:"team"`
	Bottle      string `json:"bottle"`
	HalfTankard string `json:"half_tankard"`
	FullTankard string `json:"full_tankard"`
}

type Result struct {
	Name           string `json:"name"`
	Discipline     string `json:"discipline"`
	Time           string `json:"time"`
	BaseTime       string `json:"base_time"`
	AdditionalTime string `json:"additional_time"`
	Status         string `json:"status"`
	Comment        string `json:"comment"`
}

type RankedResult struct {
	Rank           int
	Name           string
	Program        string
	Team           string
	Discipline     string
	Time           string
	BaseTime       string
	AdditionalTime string
	Status         string
	Comment        string
}

type DisciplineTab struct {
	Name    string
	Results []RankedResult
}

type Contest struct {
	FolderName   string
	DisplayName  string
	Date         string
	Official     bool
	Participants []Participant
	Disciplines  []DisciplineTab
	// summary
	TotalResults      int
	TotalParticipants int
	TotalPass         int
	TotalDQ           int
}

// ─── folder name parser ───────────────────────────────────────────────────────

// folderNameRe matches  <Name>_<YYYY-MM-DD>_<Official|Unofficial>
var folderNameRe = regexp.MustCompile(`^(.+)_(\d{4}-\d{2}-\d{2})_(Official|Unofficial)$`)

func parseFolder(name string) (displayName, date string, official bool, ok bool) {
	m := folderNameRe.FindStringSubmatch(name)
	if m == nil {
		return
	}
	displayName = strings.ReplaceAll(m[1], "_", " ")
	date = m[2]
	official = m[3] == "Official"
	ok = true
	return
}

// ─── time helpers ─────────────────────────────────────────────────────────────

// timeRe matches HH:MM:SS.mmmm
var timeRe = regexp.MustCompile(`^(\d{1,2}):(\d{2}):(\d{2})(?:\.(\d+))?$`)

func parseTimeMs(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" || strings.EqualFold(s, "nan") {
		return math.MaxInt64 // sort DQ to the bottom
	}
	m := timeRe.FindStringSubmatch(s)
	if m == nil {
		return math.MaxInt64
	}
	h := mustAtoi(m[1])
	mi := mustAtoi(m[2])
	sec := mustAtoi(m[3])
	frac := 0
	if m[4] != "" {
		fs := m[4]
		for len(fs) < 4 {
			fs += "0"
		}
		if len(fs) > 4 {
			fs = fs[:4]
		}
		frac = mustAtoi(fs)
	}
	return int64(h*36000000 + mi*600000 + sec*10000 + frac)
}

func mustAtoi(s string) int {
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}

func formatTime(s string) string {
	if strings.EqualFold(s, "nan") || s == "" {
		return "—"
	}
	return s
}

func formatAdditional(s string) string {
	if s == "" || s == "0" || s == "00:00:00.0000" {
		return "—"
	}
	return s
}

// ─── JSON loaders ─────────────────────────────────────────────────────────────

func loadParticipants(path string) ([]Participant, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var list []map[string]string
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	out := make([]Participant, 0, len(list))
	for _, m := range list {
		out = append(out, Participant{
			Name:        m["name"],
			Program:     m["program"],
			Team:        m["team"],
			Bottle:      m["bottle"],
			HalfTankard: m["half_tankard"],
			FullTankard: m["full_tankard"],
		})
	}
	return out, nil
}

func loadResults(path string) ([]Result, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var list []map[string]string
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	out := make([]Result, 0, len(list))
	for _, m := range list {
		out = append(out, Result{
			Name:           m["name"],
			Discipline:     m["discipline"],
			Time:           m["time"],
			BaseTime:       m["base_time"],
			AdditionalTime: m["additional_time"],
			Status:         m["status"],
			Comment:        m["comment"],
		})
	}
	return out, nil
}

// ─── scanner ──────────────────────────────────────────────────────────────────

var disciplineOrder = []string{
	"Bottle", "Half Tankard", "Full Tankard",
	"Bier Staphette", "Mega Medley", "Team Clash",
}

func scanContests(root string) ([]Contest, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("cannot read root folder %q: %w", root, err)
	}

	// participant lookup keyed by name – built from the participant file
	type participantKey struct{ contestIdx int }

	var contests []Contest

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		displayName, date, official, ok := parseFolder(e.Name())
		if !ok {
			fmt.Fprintf(os.Stderr, "  skip (unrecognised folder name): %s\n", e.Name())
			continue
		}

		base := filepath.Join(root, e.Name(), "contest")
		pFile := filepath.Join(base, "participants.json")
		rFile := filepath.Join(base, "results.json")

		participants, err := loadParticipants(pFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  warning: cannot load participants for %s: %v\n", e.Name(), err)
		}

		results, err := loadResults(rFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  warning: cannot load results for %s: %v\n", e.Name(), err)
		}

		// participant lookup by name → program/team
		pLookup := make(map[string]Participant, len(participants))
		for _, p := range participants {
			pLookup[p.Name] = p
		}

		// group results by discipline
		byDisc := make(map[string][]Result, 6)
		for _, r := range results {
			byDisc[r.Discipline] = append(byDisc[r.Discipline], r)
		}

		var totalPass, totalDQ int
		var tabs []DisciplineTab

		for _, disc := range disciplineOrder {
			rs, ok := byDisc[disc]
			if !ok {
				continue
			}
			// sort: pass first by time asc, DQ at the bottom
			sort.SliceStable(rs, func(i, j int) bool {
				if rs[i].Status != rs[j].Status {
					if rs[i].Status == "Disqualified" {
						return false
					}
					if rs[j].Status == "Disqualified" {
						return true
					}
				}
				return parseTimeMs(rs[i].Time) < parseTimeMs(rs[j].Time)
			})

			ranked := make([]RankedResult, 0, len(rs))
			rank := 0
			for _, r := range rs {
				p := pLookup[r.Name]
				if r.Status == "Pass" {
					rank++
				}
				displayRank := rank
				if r.Status == "Disqualified" {
					displayRank = 0
				}
				if r.Status == "Pass" {
					totalPass++
				} else {
					totalDQ++
				}
				ranked = append(ranked, RankedResult{
					Rank:           displayRank,
					Name:           r.Name,
					Program:        p.Program,
					Team:           p.Team,
					Discipline:     r.Discipline,
					Time:           formatTime(r.Time),
					BaseTime:       formatTime(r.BaseTime),
					AdditionalTime: formatAdditional(r.AdditionalTime),
					Status:         r.Status,
					Comment:        r.Comment,
				})
			}
			tabs = append(tabs, DisciplineTab{Name: disc, Results: ranked})
		}
		// handle disciplines not in predefined order
		for disc, rs := range byDisc {
			found := false
			for _, d := range disciplineOrder {
				if d == disc {
					found = true
					break
				}
			}
			if !found {
				sort.SliceStable(rs, func(i, j int) bool {
					return parseTimeMs(rs[i].Time) < parseTimeMs(rs[j].Time)
				})
				ranked := make([]RankedResult, len(rs))
				for i, r := range rs {
					p := pLookup[r.Name]
					ranked[i] = RankedResult{
						Rank:           i + 1,
						Name:           r.Name,
						Program:        p.Program,
						Team:           p.Team,
						Discipline:     r.Discipline,
						Time:           formatTime(r.Time),
						BaseTime:       formatTime(r.BaseTime),
						AdditionalTime: formatAdditional(r.AdditionalTime),
						Status:         r.Status,
						Comment:        r.Comment,
					}
				}
				tabs = append(tabs, DisciplineTab{Name: disc, Results: ranked})
			}
		}

		contests = append(contests, Contest{
			FolderName:        e.Name(),
			DisplayName:       displayName,
			Date:              date,
			Official:          official,
			Participants:      participants,
			Disciplines:       tabs,
			TotalResults:      len(results),
			TotalParticipants: len(participants),
			TotalPass:         totalPass,
			TotalDQ:           totalDQ,
		})
	}

	// newest first
	sort.SliceStable(contests, func(i, j int) bool {
		return contests[i].Date > contests[j].Date
	})

	return contests, nil
}

// ─── template data helpers ────────────────────────────────────────────────────

type TemplateData struct {
	GeneratedAt string
	Contests    []Contest
}

// ─── HTML template ────────────────────────────────────────────────────────────

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8"/>
<meta name="viewport" content="width=device-width,initial-scale=1"/>
<title>ChugWare2 – Contest Browser</title>
<style>
:root{
  --bg:#0f1117;--surface:#1a1d27;--card:#22263a;--border:#2e3350;
  --accent:#f5a623;--accent2:#e8734a;--pass:#2ecc71;--dq:#e74c3c;
  --text:#e8eaf6;--muted:#7986a8;--gold:#ffd700;--silver:#c0c0c0;--bronze:#cd7f32;
  --radius:12px;--font:'Segoe UI',system-ui,sans-serif;
}
*{box-sizing:border-box;margin:0;padding:0;}
body{background:var(--bg);color:var(--text);font-family:var(--font);min-height:100vh;}

/* ── top bar ── */
header{
  background:linear-gradient(135deg,#1a1d27 0%,#0d1020 100%);
  border-bottom:2px solid var(--accent);
  padding:18px 32px;display:flex;align-items:center;gap:16px;
  position:sticky;top:0;z-index:100;
}
header h1{font-size:1.6rem;font-weight:800;letter-spacing:1px;
  background:linear-gradient(90deg,var(--accent),var(--accent2));
  -webkit-background-clip:text;-webkit-text-fill-color:transparent;}
.header-sub{color:var(--muted);font-size:.85rem;margin-left:auto;}

/* ── layout ── */
.sidebar{
  position:fixed;top:65px;left:0;width:260px;height:calc(100vh - 65px);
  overflow-y:auto;background:var(--surface);border-right:1px solid var(--border);
  padding:16px 0;
}
.main{margin-left:260px;padding:32px;}

/* ── sidebar nav ── */
.nav-section{padding:8px 20px 4px;font-size:.7rem;text-transform:uppercase;
  letter-spacing:1.5px;color:var(--muted);font-weight:700;}
.nav-item{
  display:block;padding:10px 20px;color:var(--text);text-decoration:none;
  border-left:3px solid transparent;transition:all .15s;font-size:.9rem;
  cursor:pointer;
}
.nav-item:hover,.nav-item.active{
  background:var(--card);border-left-color:var(--accent);color:var(--accent);
}
.nav-item .badge{
  float:right;font-size:.7rem;background:var(--border);
  padding:2px 7px;border-radius:99px;color:var(--muted);
}

/* ── contest header ── */
.contest-section{display:none;animation:fadeIn .2s;}
.contest-section.visible{display:block;}
@keyframes fadeIn{from{opacity:0;transform:translateY(8px)}to{opacity:1;transform:none}}

.contest-hero{
  background:linear-gradient(135deg,var(--card),var(--surface));
  border:1px solid var(--border);border-radius:var(--radius);
  padding:28px 32px;margin-bottom:24px;display:flex;gap:24px;align-items:flex-start;
}
.contest-hero-info h2{font-size:1.8rem;font-weight:800;margin-bottom:6px;}
.contest-hero-info .date{color:var(--accent);font-size:1rem;margin-bottom:4px;}
.contest-hero-info .badge-official{
  display:inline-block;padding:3px 12px;border-radius:99px;font-size:.75rem;
  font-weight:700;letter-spacing:.5px;
}
.badge-official.official{background:rgba(245,166,35,.15);color:var(--accent);border:1px solid var(--accent);}
.badge-official.unofficial{background:rgba(116,185,255,.1);color:#74b9ff;border:1px solid #74b9ff;}

.stat-grid{display:flex;gap:12px;margin-left:auto;flex-wrap:wrap;}
.stat-card{
  background:var(--surface);border:1px solid var(--border);border-radius:10px;
  padding:14px 20px;text-align:center;min-width:90px;
}
.stat-card .val{font-size:1.6rem;font-weight:800;}
.stat-card .lbl{font-size:.72rem;color:var(--muted);text-transform:uppercase;letter-spacing:.5px;}
.val.pass{color:var(--pass);}
.val.dq{color:var(--dq);}

/* ── discipline tabs ── */
.tab-bar{display:flex;gap:6px;margin-bottom:20px;flex-wrap:wrap;}
.tab-btn{
  padding:8px 18px;border:1px solid var(--border);border-radius:8px;
  background:var(--surface);color:var(--muted);cursor:pointer;
  font-size:.85rem;font-weight:600;transition:all .15s;
}
.tab-btn:hover{border-color:var(--accent);color:var(--text);}
.tab-btn.active{background:var(--accent);color:#111;border-color:var(--accent);}
.tab-content{display:none;}
.tab-content.active{display:block;}

/* ── participants table ── */
.card{
  background:var(--card);border:1px solid var(--border);border-radius:var(--radius);
  overflow:hidden;margin-bottom:24px;
}
.card-title{
  padding:14px 20px;font-weight:700;font-size:.95rem;
  border-bottom:1px solid var(--border);background:var(--surface);
  display:flex;align-items:center;gap:8px;
}
.card-title .count{color:var(--muted);font-weight:400;font-size:.85rem;}

table{width:100%;border-collapse:collapse;}
thead tr{background:rgba(255,255,255,.03);}
th{padding:10px 14px;text-align:left;font-size:.75rem;text-transform:uppercase;
  letter-spacing:.8px;color:var(--muted);font-weight:700;white-space:nowrap;}
td{padding:11px 14px;border-bottom:1px solid rgba(255,255,255,.04);font-size:.88rem;}
tbody tr:last-child td{border-bottom:none;}
tbody tr:hover{background:rgba(255,255,255,.03);}

/* rank medals */
.rank{font-weight:800;font-size:1rem;text-align:center;width:40px;}
.rank-1{color:var(--gold);}
.rank-2{color:var(--silver);}
.rank-3{color:var(--bronze);}
.rank-dq{color:var(--dq);font-style:italic;font-size:.8rem;}

/* status pills */
.pill{padding:3px 10px;border-radius:99px;font-size:.75rem;font-weight:700;}
.pill-pass{background:rgba(46,204,113,.15);color:var(--pass);border:1px solid rgba(46,204,113,.3);}
.pill-dq{background:rgba(231,76,60,.15);color:var(--dq);border:1px solid rgba(231,76,60,.3);}

/* time chips */
.time{font-family:'Courier New',monospace;font-size:.9rem;letter-spacing:.5px;}
.time-penalty{color:var(--accent2);font-size:.8rem;}

/* empty state */
.empty{padding:40px;text-align:center;color:var(--muted);font-size:.9rem;}

/* ── participants panel ── */
.part-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(210px,1fr));gap:12px;padding:16px;}
.part-card{
  background:var(--surface);border:1px solid var(--border);border-radius:10px;
  padding:14px 16px;
}
.part-card .pname{font-weight:700;font-size:.95rem;margin-bottom:4px;}
.part-card .pinfo{font-size:.78rem;color:var(--muted);margin-bottom:8px;}
.try-dots{display:flex;gap:6px;flex-wrap:wrap;}
.try-dot{
  font-size:.7rem;padding:2px 8px;border-radius:4px;font-weight:700;
  background:rgba(245,166,35,.12);color:var(--accent);border:1px solid rgba(245,166,35,.25);
}

/* ── overview ── */
#overview-section{display:block;}
.overview-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(280px,1fr));gap:20px;margin-top:8px;}
.ov-card{
  background:var(--card);border:1px solid var(--border);border-radius:var(--radius);
  padding:22px 24px;cursor:pointer;transition:all .18s;
  text-decoration:none;color:var(--text);display:block;
}
.ov-card:hover{border-color:var(--accent);transform:translateY(-2px);box-shadow:0 8px 24px rgba(0,0,0,.4);}
.ov-card h3{font-size:1.1rem;font-weight:800;margin-bottom:6px;}
.ov-card .ov-date{color:var(--accent);font-size:.85rem;margin-bottom:12px;}
.ov-statsrow{display:flex;gap:12px;flex-wrap:wrap;}
.ov-stat{font-size:.82rem;color:var(--muted);}
.ov-stat span{color:var(--text);font-weight:700;}

/* ── scrollbar ── */
::-webkit-scrollbar{width:6px;height:6px;}
::-webkit-scrollbar-track{background:var(--bg);}
::-webkit-scrollbar-thumb{background:var(--border);border-radius:3px;}

/* ── responsive ── */
@media(max-width:768px){
  .sidebar{display:none;}
  .main{margin-left:0;padding:16px;}
  .stat-grid{margin-left:0;}
  .contest-hero{flex-direction:column;}
}
</style>
</head>
<body>

<header>
  <h1>🍺 ChugWare2</h1>
  <span style="color:var(--muted);font-size:.9rem;">Contest Browser</span>
  <span class="header-sub">Generated {{.GeneratedAt}}</span>
</header>

<nav class="sidebar" id="sidebar">
  <div class="nav-section">Contests</div>
  <a class="nav-item active" onclick="showSection('overview-section',this)">
    🏠 Overview
    <span class="badge">{{len .Contests}}</span>
  </a>
  {{range .Contests}}
  <a class="nav-item" onclick="showSection('contest-{{.FolderName}}',this)">
    {{if .Official}}🏆{{else}}🥂{{end}} {{.DisplayName}}
    <span class="badge">{{.Date}}</span>
  </a>
  {{end}}
</nav>

<main class="main">

<!-- ░░ OVERVIEW ░░ -->
<section id="overview-section">
  <div style="margin-bottom:24px;">
    <h2 style="font-size:1.5rem;margin-bottom:8px;">All Contests</h2>
    <p style="color:var(--muted);font-size:.9rem;">{{len .Contests}} contest{{if ne (len .Contests) 1}}s{{end}} found. Click a card to view results.</p>
  </div>
  <div class="overview-grid">
  {{range .Contests}}
    <a class="ov-card" onclick="showSection('contest-{{.FolderName}}',document.querySelector('[onclick*=\'{{.FolderName}}\']'))">
      <h3>{{.DisplayName}}</h3>
      <div class="ov-date">{{.Date}} · {{if .Official}}<span style="color:var(--accent)">Official</span>{{else}}<span style="color:#74b9ff">Unofficial</span>{{end}}</div>
      <div class="ov-statsrow">
        <div class="ov-stat">Participants: <span>{{.TotalParticipants}}</span></div>
        <div class="ov-stat">Results: <span>{{.TotalResults}}</span></div>
        <div class="ov-stat" style="color:var(--pass)">Pass: <span>{{.TotalPass}}</span></div>
        <div class="ov-stat" style="color:var(--dq)">DQ: <span>{{.TotalDQ}}</span></div>
      </div>
    </a>
  {{end}}
  {{if eq (len .Contests) 0}}
    <div class="empty">No contests found. Run the Contest Wizard in ChugWare2 first.</div>
  {{end}}
  </div>
</section>

<!-- ░░ PER-CONTEST ░░ -->
{{range $ci, $c := .Contests}}
<section id="contest-{{$c.FolderName}}" class="contest-section">

  <!-- hero -->
  <div class="contest-hero">
    <div class="contest-hero-info">
      <h2>{{$c.DisplayName}}</h2>
      <div class="date">📅 {{$c.Date}}</div>
      {{if $c.Official}}
        <span class="badge-official official">OFFICIAL</span>
      {{else}}
        <span class="badge-official unofficial">UNOFFICIAL</span>
      {{end}}
    </div>
    <div class="stat-grid">
      <div class="stat-card"><div class="val">{{$c.TotalParticipants}}</div><div class="lbl">Athletes</div></div>
      <div class="stat-card"><div class="val">{{$c.TotalResults}}</div><div class="lbl">Results</div></div>
      <div class="stat-card"><div class="val pass">{{$c.TotalPass}}</div><div class="lbl">Pass</div></div>
      <div class="stat-card"><div class="val dq">{{$c.TotalDQ}}</div><div class="lbl">DQ</div></div>
      <div class="stat-card"><div class="val">{{len $c.Disciplines}}</div><div class="lbl">Disciplines</div></div>
    </div>
  </div>

  <!-- discipline tabs -->
  {{if $c.Disciplines}}
  <div class="card" style="margin-bottom:24px;">
    <div class="card-title">Results by Discipline</div>
    <div style="padding:16px 16px 0;">
      <div class="tab-bar" id="tabs-{{$ci}}">
        {{range $di, $disc := $c.Disciplines}}
        <button class="tab-btn{{if eq $di 0}} active{{end}}"
          onclick="switchTab('tabs-{{$ci}}','tabcontent-{{$ci}}',{{$di}},this)">
          {{$disc.Name}}
          <span style="opacity:.6;font-weight:400;margin-left:4px;">({{len $disc.Results}})</span>
        </button>
        {{end}}
      </div>
    </div>
    {{range $di, $disc := $c.Disciplines}}
    <div class="tab-content{{if eq $di 0}} active{{end}}" id="tabcontent-{{$ci}}-{{$di}}">
      {{if $disc.Results}}
      <table>
        <thead>
          <tr>
            <th style="width:50px;">Rank</th>
            <th>Athlete</th>
            <th>Program</th>
            <th>Team</th>
            <th>Time</th>
            <th>Base</th>
            <th>Penalty</th>
            <th>Status</th>
            <th>Comment</th>
          </tr>
        </thead>
        <tbody>
        {{range $disc.Results}}
          <tr>
            <td class="rank{{if eq .Rank 1}} rank-1{{else if eq .Rank 2}} rank-2{{else if eq .Rank 3}} rank-3{{else if eq .Rank 0}} rank-dq{{end}}">
              {{if eq .Rank 0}}DQ{{else if eq .Rank 1}}🥇{{else if eq .Rank 2}}🥈{{else if eq .Rank 3}}🥉{{else}}{{.Rank}}{{end}}
            </td>
            <td style="font-weight:600;">{{.Name}}</td>
            <td style="color:var(--muted);">{{.Program}}</td>
            <td style="color:var(--muted);">{{.Team}}</td>
            <td class="time">{{.Time}}</td>
            <td class="time" style="color:var(--muted);font-size:.82rem;">{{.BaseTime}}</td>
            <td class="time-penalty">{{.AdditionalTime}}</td>
            <td>
              {{if eq .Status "Pass"}}
                <span class="pill pill-pass">Pass</span>
              {{else}}
                <span class="pill pill-dq">DQ</span>
              {{end}}
            </td>
            <td style="color:var(--muted);font-size:.82rem;">{{.Comment}}</td>
          </tr>
        {{end}}
        </tbody>
      </table>
      {{else}}
        <div class="empty">No results recorded for this discipline.</div>
      {{end}}
    </div>
    {{end}}
  </div>
  {{else}}
    <div class="card"><div class="empty">No results recorded yet for this contest.</div></div>
  {{end}}

  <!-- participants panel -->
  <div class="card">
    <div class="card-title">
      Registered Athletes
      <span class="count">{{len $c.Participants}} total</span>
    </div>
    {{if $c.Participants}}
    <div class="part-grid">
      {{range $c.Participants}}
      <div class="part-card">
        <div class="pname">{{.Name}}</div>
        <div class="pinfo">{{.Program}}{{if and .Program .Team}} · {{end}}{{.Team}}</div>
        <div class="try-dots">
          {{if .Bottle}}<span class="try-dot">🍶 {{.Bottle}}</span>{{end}}
          {{if .HalfTankard}}<span class="try-dot">🍺 {{.HalfTankard}}</span>{{end}}
          {{if .FullTankard}}<span class="try-dot">🍻 {{.FullTankard}}</span>{{end}}
        </div>
      </div>
      {{end}}
    </div>
    {{else}}
    <div class="empty">No participants registered.</div>
    {{end}}
  </div>

</section>
{{end}}

</main>

<script>
function showSection(id, navEl) {
  // hide all sections
  document.querySelectorAll('.contest-section, #overview-section').forEach(s => {
    s.classList.remove('visible');
    if(s.id === 'overview-section') s.style.display = 'none';
    else s.style.display = 'none';
  });
  // show target
  const el = document.getElementById(id);
  if(el) {
    el.style.display = 'block';
    el.classList.add('visible');
  }
  // update nav
  document.querySelectorAll('.nav-item').forEach(n => n.classList.remove('active'));
  if(navEl) navEl.classList.add('active');
  window.scrollTo(0,0);
}

function switchTab(tabBarId, contentBaseId, idx, btn) {
  // deactivate all buttons in this tab bar
  document.querySelectorAll('#' + tabBarId + ' .tab-btn').forEach(b => b.classList.remove('active'));
  btn.classList.add('active');
  // hide all content panels for this contest
  let i = 0;
  while(true) {
    const el = document.getElementById(contentBaseId + '-' + i);
    if(!el) break;
    el.classList.remove('active');
    i++;
  }
  // show selected
  const target = document.getElementById(contentBaseId + '-' + idx);
  if(target) target.classList.add('active');
}

// init: show overview
document.getElementById('overview-section').style.display = 'block';
</script>
</body>
</html>`

// ─── entry point ──────────────────────────────────────────────────────────────

func main() {
	root := flag.String("root", "ChugWare", "Path to ChugWare contests folder")
	out := flag.String("out", "chugware_results.html", "Output HTML file path")
	flag.Parse()

	absRoot, err := filepath.Abs(*root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error resolving root path: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Scanning contests in: %s\n", absRoot)
	contests, err := scanContests(absRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error scanning contests: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d contest(s)\n", len(contests))

	absOut, err := filepath.Abs(*out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error resolving output path: %v\n", err)
		os.Exit(1)
	}

	f, err := os.Create(absOut)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	tmpl, err := template.New("html").Parse(htmlTemplate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing template: %v\n", err)
		os.Exit(1)
	}

	data := TemplateData{
		GeneratedAt: time.Now().Format("2006-01-02 15:04:05"),
		Contests:    contests,
	}

	if err := tmpl.Execute(f, data); err != nil {
		fmt.Fprintf(os.Stderr, "error rendering template: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("HTML written to: %s\n", absOut)
}
