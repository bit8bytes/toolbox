// Package vcs provides utilities for working with version control systems.
//
// This package extracts version control information from Go binaries built
// with module support and version control integration.
package vcs

import (
	"fmt"
	"runtime/debug"
)

// Version returns the build time and revision information.
// The returned string format is "time-revision" or "time-revision-dirty"
// if the working directory had uncommitted changes at build time.
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

	if modified {
		return fmt.Sprintf("%s-%s-dirty", time, revision)
	}

	return fmt.Sprintf("%s-%s", time, revision)
}
