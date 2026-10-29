// Package project resolves a project's spec directory from setting.json,
// falling back to ~/.apicat/projects/<encoded-project-path>.
package project

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

// DirFor reads the project's setting.json, preferring the local .apicat
// directory over the default project directory. Relative spec paths resolve
// against the project root. An existing config replaces the fallback config;
// an unset spec_dir uses the default directory. Parent configs are not searched.
func DirFor(projectPath string) (string, error) {
	abs, err := filepath.Abs(projectPath)
	if err != nil {
		return "", fmt.Errorf("cannot resolve project path: %w", err)
	}
	specDir, found, err := readSpecDir(filepath.Join(abs, ".apicat", "setting.json"))
	if err != nil {
		return "", err
	}
	if found && strings.TrimSpace(specDir) != "" {
		return resolveSpecDir(abs, specDir), nil
	}
	root, err := Root()
	if err != nil {
		return "", err
	}
	defaultDir := filepath.Join(root, Encode(abs))
	if !found {
		specDir, _, err = readSpecDir(filepath.Join(defaultDir, "setting.json"))
		if err != nil {
			return "", err
		}
		if strings.TrimSpace(specDir) != "" {
			return resolveSpecDir(abs, specDir), nil
		}
	}
	return defaultDir, nil
}

func resolveSpecDir(projectDir, specDir string) string {
	if filepath.IsAbs(specDir) {
		return filepath.Clean(specDir)
	}
	return filepath.Join(projectDir, specDir)
}

// readSpecDir distinguishes an absent config from an existing empty setting.
func readSpecDir(configPath string) (specDir string, found bool, err error) {
	data, err := os.ReadFile(configPath)
	if errors.Is(err, os.ErrNotExist) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("cannot read config %s: %w", configPath, err)
	}
	var config *struct {
		SpecDir string `json:"spec_dir"`
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&config); err != nil {
		return "", true, fmt.Errorf("invalid config %s: %w", configPath, err)
	}
	if config == nil {
		return "", true, fmt.Errorf("invalid config %s: expected a JSON object", configPath)
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return "", true, fmt.Errorf("invalid config %s: expected a single JSON object", configPath)
	}
	return config.SpecDir, true, nil
}
