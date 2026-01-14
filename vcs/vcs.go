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

// Version returns build info as "time-revision-int" or "time-revision-int-dirty".
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

// revisionToInt converts the first 8 hex chars of a git hash to int32.
// Suitable for SQLite user_version. Returns error if revision < 8 chars or invalid hex.
func revisionToInt(revision string) (int32, error) {
	if len(revision) < 8 {
		return 0, fmt.Errorf("revision too short: need at least 8 characters, got %d", len(revision))
	}

	// Take first 8 hex characters (4 bytes)
	hexStr := revision[:8]

	val, err := strconv.ParseUint(hexStr, 16, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse revision hex: %w", err)
	}

	return int32(val), nil
}

// VersionInt32 returns the git revision as int32 for use as SQLite user_version.
// The integer is derived from the first 8 hex chars of the vcs.revision hash.
func VersionInt32() (int32, error) {
	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range bi.Settings {
			if s.Key == "vcs.revision" {
				return revisionToInt(s.Value)
			}
		}
	}

	return 0, nil
}
