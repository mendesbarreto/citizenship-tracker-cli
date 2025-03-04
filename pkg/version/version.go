package version

import "fmt"

// Version information
var (
	Version = "0.0.1"
	Commit  = "unknown"
	Date    = "unknown"
)

// VersionInfo returns a formatted version string
func VersionInfo() string {
	return fmt.Sprintf("citizenship-tracker-cli version %s (commit: %s, built at: %s)",
		Version, Commit, Date)
}
