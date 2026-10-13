// Package project maps a user's project directory to its apicat working
// directory under ~/.apicat/projects/. Specs live centrally there, not in
// the user's repository; the mapping key is the project's absolute path
// encoded Claude Code-style: every character outside [A-Za-z0-9_-] becomes
// a dash, so /Users/alice/dev/shop becomes -Users-alice-dev-shop.
package project

import (
	"fmt"
	"os"
	"path/filepath"
)

// Encode turns an absolute project path into its directory name under the
// projects root.
func Encode(absPath string) string {
	out := make([]rune, 0, len(absPath))
	for _, r := range absPath {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_', r == '-':
			out = append(out, r)
		default:
			out = append(out, '-')
		}
	}
	return string(out)
}

// Root returns ~/.apicat/projects.
func Root() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot resolve home directory: %w", err)
	}
	return filepath.Join(home, ".apicat", "projects"), nil
}

// DirFor returns the apicat working directory for the given project path
// (typically the current working directory). The directory is not created.
func DirFor(projectPath string) (string, error) {
	abs, err := filepath.Abs(projectPath)
	if err != nil {
		return "", fmt.Errorf("cannot resolve project path: %w", err)
	}
	root, err := Root()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, Encode(abs)), nil
}
