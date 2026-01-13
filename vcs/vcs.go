// Package vcs provides utilities for working with version control systems.
//
// This package extracts version control information from Go binaries built
// with module support and version control integration.
package vcs

import (
	"fmt"
	"runtime/debug"
	"strconv"
)

// Version returns the build time, revision hash, and revision integer.
// The returned string format is "time-revision-int" or "time-revision-int-dirty"
// if the working directory had uncommitted changes at build time.
// The integer is derived from the first 4 bytes of the git revision hash.
func Version() string {
	var (
		time     string
		revision string
		modified bool
	)

	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.time":
				time = s.Value
			case "vcs.revision":
				revision = s.Value
			case "vcs.modified":
				if s.Value == "true" {
					modified = true
				}
			}
		}
	}

	revInt, err := revisionToInt(revision)
	if err != nil {
		revInt = 0
	}

	if modified {
		return fmt.Sprintf("%s-%s-%d-dirty", time, revision, revInt)
	}

	return fmt.Sprintf("%s-%s-%d", time, revision, revInt)
}

// revisionToInt converts the first 4 bytes (8 hex characters) of a git revision hash
// to a signed 32-bit integer suitable for use as SQLite user_version.
// The conversion uses the first 8 hex characters of the hash, which are
// cryptographically random and uniformly distributed.
//
// Example: "abc123def456..." -> 0xabc123de -> -1412567330
//
// Returns an error if the revision is too short (< 8 characters) or contains
// invalid hexadecimal characters.
func revisionToInt(revision string) (int32, error) {
	if len(revision) < 8 {
		return 0, fmt.Errorf("revision too short: need at least 8 characters, got %d", len(revision))
	}

	// Take first 8 hex characters (4 bytes)
	hexStr := revision[:8]

	// Parse as uint32 first to handle the full range
	val, err := strconv.ParseUint(hexStr, 16, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse revision hex: %w", err)
	}

	// Convert to int32 (may be negative due to sign bit)
	return int32(val), nil
}

// VersionInt32 returns the git revision hash and its integer representation.
// The integer is derived from the first 4 bytes (8 hex characters) of the
// revision hash and is suitable for use as SQLite user_version (signed 32-bit integer).
//
// This function extracts the vcs.revision value from the build information
// embedded in the binary during compilation with -buildvcs flag.
//
// Returns:
//   - revision: The full git commit hash (e.g., "abc123def456...")
//   - revisionInt: A signed 32-bit integer representation (e.g., -1412567330)
//   - error: Non-nil if build info is unavailable, revision not found, or parsing fails
func VersionInt32() (string, int32, error) {
	var revision string

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return "", 0, fmt.Errorf("build info not available")
	}

	for _, s := range bi.Settings {
		if s.Key == "vcs.revision" {
			revision = s.Value
			break
		}
	}

	if revision == "" {
		return "", 0, fmt.Errorf("vcs.revision not found in build info")
	}

	revInt, err := revisionToInt(revision)
	if err != nil {
		return revision, 0, err
	}

	return revision, revInt, nil
}
