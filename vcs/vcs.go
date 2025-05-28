// Package vcs implements utility for working with version control systems.
//
// This package has only one function and can be called using vcs.Version()
// The package vcs requires a version control system such as git and will
// only work on a build binary.
package vcs

import (
	"fmt"
	"runtime/debug"
)

// Func version return the time and revision.
// If the revision is dirty, it will be appended with -dirty.
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
