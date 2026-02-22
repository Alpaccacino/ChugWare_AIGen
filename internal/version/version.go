// Package version holds build-time metadata that is injected via -ldflags.
//
// To stamp a release build use:
//
//	go build -ldflags "-X chugware/internal/version.Version=1.0.0 \
//	                   -X chugware/internal/version.BuildDate=2026-02-22 \
//	                   -X chugware/internal/version.GitCommit=abc1234" \
//	         -o ChugWare2.exe ./cmd/
//
// All variables fall back to safe defaults when built without ldflags
// (e.g. during development with `go run`).
package version

// Version is the semantic version string (e.g. "1.0.0").
// Injected at link time via: -ldflags "-X chugware/internal/version.Version=<ver>"
var Version = "1.0.0"

// BuildDate is the UTC date the binary was compiled (e.g. "2026-02-22").
// Injected at link time via: -ldflags "-X chugware/internal/version.BuildDate=<date>"
var BuildDate = "unknown"

// GitCommit is the short git commit hash the binary was built from.
// Injected at link time via: -ldflags "-X chugware/internal/version.GitCommit=<hash>"
var GitCommit = "unknown"
