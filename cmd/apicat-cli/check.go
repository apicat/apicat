package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/apicat/apicat/v3/internal/converter"
	"github.com/apicat/apicat/v3/internal/lint"
	"github.com/apicat/apicat/v3/internal/loader"
	"github.com/apicat/apicat/v3/internal/model"
	"github.com/apicat/apicat/v3/internal/project"
)

type checkHint struct {
	Method  string `json:"method,omitempty"`
	Path    string `json:"path,omitempty"`
	Message string `json:"message"`
}

type checkReport struct {
	Status    string      `json:"status"`
	SpecDir   string      `json:"spec_dir"`
	Errors    []string    `json:"errors"`
	LintTotal int         `json:"lint_total"`
	Hints     []checkHint `json:"hints"`
}

// runCheck never starts a server. Exit codes: 0 clean, 1 quality hints,
// 2 loading, usage, or output errors. JSON reports go only to stdout.
func runCheck(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("check", flag.ContinueOnError)
	flags.SetOutput(stderr)
	jsonOutput := flags.Bool("json", false, "output a JSON report")
	flags.Usage = func() {
		fmt.Fprintln(stderr, "Usage: apicat-cli check [-json] [spec-directory]")
		fmt.Fprintln(stderr, "Check loading, references, and documentation quality without starting a server.")
		flags.PrintDefaults()
		fmt.Fprintln(stderr, "Exit codes: 0 clean, 1 quality hints, 2 loading or usage error.")
	}
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if flags.NArg() > 1 {
		fmt.Fprintln(stderr, "Error: expected at most one spec directory; put flags before the directory")
		return 2
	}

	dir := flags.Arg(0)
	var err error
	if dir == "" {
		var cwd string
		cwd, err = os.Getwd()
		if err == nil {
			dir, err = project.DirFor(cwd)
		}
	} else {
		dir, err = filepath.Abs(dir)
	}
	report := checkReport{Status: "error", SpecDir: dir, Errors: []string{}, Hints: []checkHint{}}
	code := 2
	if err == nil {
		var doc *model.APIDoc
		doc, err = loadCheckDocument(dir)
		if err == nil {
			report.LintTotal = doc.LintTotal
			if doc.Description == "" {
				report.Hints = append(report.Hints, checkHint{Message: "API has no description"})
			}
			for _, group := range doc.TagGroups {
				for _, ep := range group.Endpoints {
					for _, hint := range ep.LintHints {
						report.Hints = append(report.Hints, checkHint{Method: ep.Method, Path: ep.Path, Message: hint})
					}
				}
			}
			report.Status, code = "ok", 0
			if report.LintTotal > 0 {
				report.Status, code = "hints", 1
			}
		}
	}
	if err != nil {
		report.Errors = append(report.Errors, err.Error())
	}

	if *jsonOutput {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		err = enc.Encode(report)
	} else {
		err = writeCheckText(stdout, report)
	}
	if err != nil {
		fmt.Fprintf(stderr, "Error writing check result: %v\n", err)
		return 2
	}
	return code
}

func loadCheckDocument(dir string) (*model.APIDoc, error) {
	spec, err := loader.Load(dir)
	if err != nil {
		return nil, err
	}
	doc, err := converter.Convert(spec)
	if err != nil {
		return nil, err
	}
	// The renderer repeats operations under each tag. Check each operation
	// once, in stable order, so reports do not depend on tag map iteration.
	var endpoints []model.Endpoint
	seen := make(map[string]bool)
	for _, group := range doc.TagGroups {
		for _, ep := range group.Endpoints {
			key := ep.Method + " " + ep.Path
			if !seen[key] {
				seen[key] = true
				endpoints = append(endpoints, ep)
			}
		}
	}
	sort.Slice(endpoints, func(i, j int) bool {
		if endpoints[i].Path != endpoints[j].Path {
			return endpoints[i].Path < endpoints[j].Path
		}
		return endpoints[i].Method < endpoints[j].Method
	})
	doc.TagGroups = []model.TagGroup{{Endpoints: endpoints}}
	lint.Annotate(doc)
	return doc, nil
}

func writeCheckText(w io.Writer, report checkReport) error {
	// Buffer first so write errors are handled consistently with JSON output.
	var buf []byte
	buf = fmt.Appendf(buf, "Spec dir: %s\n", report.SpecDir)
	for _, err := range report.Errors {
		buf = fmt.Appendf(buf, "Error: %s\n", err)
	}
	for _, hint := range report.Hints {
		location := "API"
		if hint.Method != "" {
			location = hint.Method + " " + hint.Path
		}
		buf = fmt.Appendf(buf, "%s: %s\n", location, hint.Message)
	}
	if report.Status == "ok" {
		buf = fmt.Appendf(buf, "Check passed: document loaded; no quality hints.\n")
	} else if report.Status == "hints" {
		buf = fmt.Appendf(buf, "Check found %d quality hints.\n", report.LintTotal)
	}
	_, err := w.Write(buf)
	return err
}
