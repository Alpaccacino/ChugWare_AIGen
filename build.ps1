# build.ps1 – Builds ChugWare2 with version metadata stamped into the binary.
#
# Usage:
#   .\build.ps1                      # uses APP_VERSION default (1.0.0)
#   .\build.ps1 -Version "1.1.0"     # override version
#
# The script injects three variables at link time:
#   chugware/internal/version.Version   – semantic version  (e.g. "1.0.0")
#   chugware/internal/version.BuildDate – UTC build date    (e.g. "2026-02-22")
#   chugware/internal/version.GitCommit – short git hash    (e.g. "abc1234")
#
# Both executables are built:
#   ChugWare2.exe   – the main contest-management GUI application
#   htmlgen.exe     – the standalone HTML contest browser generator

param(
    [string]$Version = "1.0.0"
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

# ── Collect build metadata ────────────────────────────────────────────────────

$BuildDate = (Get-Date -Format "yyyy-MM-dd")

# Try to grab the current short git hash; fall back gracefully if git is absent.
try {
    $GitCommit = (git rev-parse --short HEAD 2>$null).Trim()
    if (-not $GitCommit) { $GitCommit = "unknown" }
} catch {
    $GitCommit = "unknown"
}

Write-Host ""
Write-Host "=== ChugWare2 Build ===" -ForegroundColor Cyan
Write-Host "  Version   : $Version"
Write-Host "  BuildDate : $BuildDate"
Write-Host "  GitCommit : $GitCommit"
Write-Host ""

# ── Common ldflags ────────────────────────────────────────────────────────────

$LdFlags = "-X chugware/internal/version.Version=$Version " +
           "-X chugware/internal/version.BuildDate=$BuildDate " +
           "-X chugware/internal/version.GitCommit=$GitCommit"

# ── Build main application ────────────────────────────────────────────────────

Write-Host "Building ChugWare2.exe ..." -ForegroundColor Yellow
go build -ldflags $LdFlags -o ChugWare2.exe ./cmd/
if ($LASTEXITCODE -ne 0) {
    Write-Error "Build failed for ChugWare2.exe"
    exit 1
}
Write-Host "  -> ChugWare2.exe" -ForegroundColor Green

# ── Build htmlgen ─────────────────────────────────────────────────────────────

Write-Host "Building htmlgen.exe ..." -ForegroundColor Yellow
go build -ldflags $LdFlags -o htmlgen.exe ./cmd/htmlgen/
if ($LASTEXITCODE -ne 0) {
    Write-Error "Build failed for htmlgen.exe"
    exit 1
}
Write-Host "  -> htmlgen.exe" -ForegroundColor Green

# ── Done ──────────────────────────────────────────────────────────────────────

Write-Host ""
Write-Host "Build complete." -ForegroundColor Cyan
Write-Host ""
Write-Host "To generate the HTML contest browser after running a contest:"
Write-Host '  .\htmlgen.exe --root .\ChugWare --out chugware_results.html'
Write-Host ""
