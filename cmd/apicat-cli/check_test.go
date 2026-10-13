package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/apicat/apicat/v3/internal/project"
)

const cleanCheckSpec = `openapi: 3.1.2
info: {title: Check, version: '1', description: Test API.}
paths:
  /pets:
    get:
      tags: [pets, public]
      summary: List pets
      operationId: listPets
      responses:
        '200':
          description: Pets found.
          content:
            application/json:
              schema:
                $ref: './pet.yaml#/Pet'
        '400': {description: Invalid request.}
`

func writeCheckSpec(t *testing.T, dir, spec string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, content := range map[string]string{
		"openapi.yaml": spec,
		"pet.yaml": `Pet:
  type: object
  properties:
    name: {type: [string, "null"], description: Pet name.}
`,
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestCheckJSON(t *testing.T) {
	cases := []struct {
		name, spec, status string
		code, lintTotal    int
	}{
		{"clean 3.1 with external ref", cleanCheckSpec, "ok", 0, 0},
		{"clean 3.0", strings.ReplaceAll(strings.ReplaceAll(cleanCheckSpec, "3.1.2", "3.0.3"), "./pet.yaml#/Pet", "#/components/schemas/Pet") + "components:\n  schemas:\n    Pet: {type: string, nullable: true}\n", "ok", 0, 0},
		{"hints counted once across tags", strings.ReplaceAll(cleanCheckSpec, "      operationId: listPets\n", ""), "hints", 1, 1},
		{"API description", strings.ReplaceAll(cleanCheckSpec, ", description: Test API.", ""), "hints", 1, 1},
		{"invalid YAML", "openapi: [", "error", 2, 0},
		{"missing info", "openapi: 3.1.2\npaths: {}\n", "error", 2, 0},
		{"missing ref file", strings.ReplaceAll(cleanCheckSpec, "./pet.yaml#/Pet", "./missing.yaml#/Pet"), "error", 2, 0},
		{"missing ref target", strings.ReplaceAll(cleanCheckSpec, "./pet.yaml#/Pet", "./pet.yaml#/Missing"), "error", 2, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			writeCheckSpec(t, dir, tc.spec)
			var stdout, stderr bytes.Buffer
			if code := runCheck([]string{"-json", dir}, &stdout, &stderr); code != tc.code {
				t.Fatalf("exit = %d, want %d; stdout=%s stderr=%s", code, tc.code, &stdout, &stderr)
			}
			var report checkReport
			if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
				t.Fatalf("invalid JSON output: %v: %s", err, &stdout)
			}
			if report.Status != tc.status || report.SpecDir != dir || report.LintTotal != tc.lintTotal {
				t.Errorf("unexpected report: %+v", report)
			}
			if (len(report.Errors) > 0) != (tc.code == 2) {
				t.Errorf("unexpected errors: %v", report.Errors)
			}
			if report.Errors == nil || report.Hints == nil {
				t.Error("report arrays must not be null")
			}
			if tc.name == "hints counted once across tags" {
				if len(report.Hints) != 1 || report.Hints[0].Method != "GET" || report.Hints[0].Path != "/pets" || report.Hints[0].Message != "no operationId" {
					t.Errorf("missing operation context: %+v", report.Hints)
				}
			}
		})
	}
}

func TestCheckProjectMode(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Chdir(t.TempDir())
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir, err := project.DirFor(cwd)
	if err != nil {
		t.Fatal(err)
	}
	writeCheckSpec(t, dir, cleanCheckSpec)
	var stdout, stderr bytes.Buffer
	if code := runCheck(nil, &stdout, &stderr); code != 0 {
		t.Fatalf("exit = %d; stdout=%s stderr=%s", code, &stdout, &stderr)
	}
	if !strings.Contains(stdout.String(), dir) || !strings.Contains(stdout.String(), "Check passed") {
		t.Errorf("unexpected report: %s", &stdout)
	}
}

func TestCheckMissingSpec(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := runCheck([]string{"-json", t.TempDir()}, &stdout, &stderr); code != 2 {
		t.Fatalf("exit = %d, want 2", code)
	}
	var report checkReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil || len(report.Errors) == 0 {
		t.Fatalf("expected JSON loading error: %s", &stdout)
	}
}

func TestCheckUsage(t *testing.T) {
	for _, tc := range []struct {
		args []string
		code int
	}{
		{[]string{"-h"}, 0},
		{[]string{"-watch"}, 2},
		{[]string{"one", "two"}, 2},
		{[]string{"some-dir", "-json"}, 2},
	} {
		var stdout, stderr bytes.Buffer
		if code := runCheck(tc.args, &stdout, &stderr); code != tc.code {
			t.Errorf("%v exit = %d, want %d", tc.args, code, tc.code)
		}
		if stdout.Len() != 0 || stderr.Len() == 0 {
			t.Errorf("%v: stdout=%s stderr=%s", tc.args, &stdout, &stderr)
		}
	}
}

func TestCheckTextHints(t *testing.T) {
	dir := t.TempDir()
	writeCheckSpec(t, dir, strings.ReplaceAll(cleanCheckSpec, "      operationId: listPets\n", ""))
	var stdout, stderr bytes.Buffer
	if code := runCheck([]string{dir}, &stdout, &stderr); code != 1 {
		t.Fatalf("exit = %d; stdout=%s stderr=%s", code, &stdout, &stderr)
	}
	if !strings.Contains(stdout.String(), "GET /pets: no operationId") {
		t.Errorf("missing actionable hint: %s", &stdout)
	}
}
