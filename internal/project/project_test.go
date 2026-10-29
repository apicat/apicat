package project

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"/Users/alice/dev/shop", "-Users-alice-dev-shop"},
		{"/home/dev/my.project", "-home-dev-my-project"},
		{"/srv/api_v2", "-srv-api_v2"},
		{"/tmp/空格 dir", "-tmp----dir"},
	}
	for _, c := range cases {
		if got := Encode(c.in); got != c.want {
			t.Errorf("Encode(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestDirFor(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	projectDir := t.TempDir()
	dir, err := DirFor(projectDir)
	if err != nil {
		t.Fatalf("DirFor: %v", err)
	}
	want := filepath.Join(homeDir, ".apicat", "projects", Encode(projectDir))
	if dir != want {
		t.Errorf("DirFor = %q, want %q", dir, want)
	}
}

func TestDirForConfig(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	projectDir := t.TempDir()
	defaultDir := filepath.Join(homeDir, ".apicat", "projects", Encode(projectDir))
	if err := os.MkdirAll(filepath.Join(projectDir, ".apicat"), 0o755); err != nil {
		t.Fatal(err)
	}
	absoluteDir := filepath.Join(t.TempDir(), "API docs")
	absoluteJSON, err := json.Marshal(map[string]string{"spec_dir": absoluteDir})
	if err != nil {
		t.Fatal(err)
	}
	// Resolve relative paths against the supplied project, not the process cwd.
	t.Chdir(t.TempDir())
	for _, tc := range []struct {
		name, config, want string
		wantErr            bool
	}{
		{"relative", `{"spec_dir":"docs/api"}`, filepath.Join(projectDir, "docs/api"), false},
		{"absolute", string(absoluteJSON), absoluteDir, false},
		{"parent", `{"spec_dir":"../shared"}`, filepath.Join(projectDir, "../shared"), false},
		{"project root", `{"spec_dir":"."}`, projectDir, false},
		{"unset", `{}`, defaultDir, false},
		{"empty", `{"spec_dir":""}`, defaultDir, false},
		{"blank", `{"spec_dir":"  "}`, defaultDir, false},
		{"malformed", `{`, "", true},
		{"empty file", ``, "", true},
		{"null", `null`, "", true},
		{"wrong type", `{"spec_dir":123}`, "", true},
		{"unknown field", `{"spec_dri":"docs"}`, "", true},
		{"multiple objects", `{} {}`, "", true},
		{"trailing garbage", `{} garbage`, "", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			configPath := filepath.Join(projectDir, ".apicat", "setting.json")
			if err := os.WriteFile(configPath, []byte(tc.config), 0o644); err != nil {
				t.Fatal(err)
			}
			got, err := DirFor(projectDir)
			if tc.wantErr {
				if err == nil || !strings.Contains(err.Error(), configPath) {
					t.Fatalf("expected error identifying config, got %v", err)
				}
				return
			}
			if err != nil || got != tc.want {
				t.Fatalf("DirFor = %q, %v; want %q", got, err, tc.want)
			}
		})
	}
}

func TestDirForConfigReadError(t *testing.T) {
	projectDir := t.TempDir()
	configPath := filepath.Join(projectDir, ".apicat", "setting.json")
	if err := os.MkdirAll(configPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := DirFor(projectDir); err == nil || !strings.Contains(err.Error(), configPath) {
		t.Fatalf("expected config read error, got %v", err)
	}
}

func TestDirForConfigPrecedence(t *testing.T) {
	for _, tc := range []struct {
		name, local, fallback, want string
		wantErr                     bool
	}{
		{"neither config", "", "", "", false},
		{"fallback config", "", `{"spec_dir":"shared/docs"}`, "shared/docs", false},
		{"local overrides fallback", `{"spec_dir":"docs/api"}`, `{"spec_dir":"shared/docs"}`, "docs/api", false},
		{"empty local shadows fallback", `{}`, `{"spec_dir":"shared/docs"}`, "", false},
		{"empty fallback", "", `{}`, "", false},
		{"invalid local stops resolution", `{`, `{"spec_dir":"shared/docs"}`, "", true},
		{"invalid fallback", "", `{`, "", true},
		{"local skips invalid fallback", `{"spec_dir":"docs/api"}`, `{`, "docs/api", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			homeDir := t.TempDir()
			t.Setenv("HOME", homeDir)
			projectDir := t.TempDir()
			defaultDir := filepath.Join(homeDir, ".apicat", "projects", Encode(projectDir))
			for path, data := range map[string]string{
				filepath.Join(projectDir, ".apicat", "setting.json"): tc.local,
				filepath.Join(defaultDir, "setting.json"):            tc.fallback,
			} {
				if data == "" {
					continue
				}
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			got, err := DirFor(projectDir)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected config error")
				}
				return
			}
			want := defaultDir
			if tc.want != "" {
				want = filepath.Join(projectDir, tc.want)
			}
			if err != nil || got != want {
				t.Fatalf("DirFor = %q, %v; want %q", got, err, want)
			}
		})
	}
}
