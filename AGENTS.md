# Repository Guidelines

## Project Structure & Module Organization

`cmd/apicat-cli/main.go` wires the OpenAPI documentation CLI together. `internal/` contains project resolution, loading, conversion, models, linting, rendering, serving, and watching. Templates depend only on `internal/model`, not libopenapi types.

`web/templates/` holds Go HTML templates; `web/static/` holds Bulma, custom CSS, and vanilla JavaScript, embedded by `web/embed.go`. Tests sit beside their packages. `skills/apicat-spec/` documents spec-authoring conventions. Local `testdata/` and `docs/` are gitignored.

## Build, Test, and Development Commands

Use Go 1.25.0 or newer; no Node/npm build is required.

- `go build -o apicat-cli ./cmd/apicat-cli`: build the standalone CLI.
- `go build ./...`: compile all packages.
- `go test ./...`: run all unit tests.
- `go run ./cmd/apicat-cli -watch /path/to/spec-dir`: serve docs at `http://127.0.0.1:8080` and reload spec changes.

Provide `openapi.yaml`, `openapi.yml`, or `openapi.json`. Without a directory argument, specs resolve under `~/.apicat/projects/<encoded-cwd>/`; encoding replaces characters outside `[A-Za-z0-9_-]` with `-`. Restart/rebuild after editing `web/`; `-watch` only reloads specs.

## Coding Style & Naming Conventions

Use `gofmt`, tabs in Go, lowercase package names, exported `PascalCase`, and unexported `camelCase`. Match two-space JavaScript/CSS indentation. No dedicated linter configuration exists.

## Architecture & Change Contracts

- Load → convert → lint → render/serve. Watch reloads must atomically replace the document only on success, retaining the last good version on failure.
- Keep `converter.makeAnchor`, navigation `data-anchor`, endpoint `detail-<anchor>` IDs, and JavaScript `showEndpoint()` synchronized. Intercept navigation clicks only on `/endpoints`.
- Register template helpers in `internal/render/render.go`. New pages require a server route, `pageData.Page`, and a content block dispatched by `layout.html`.
- New lint rules must update `ruleOrder` and the spec skill's Quality rules. Field-level missing-description markers are handled directly by templates.
- Encoding changes require synchronized updates to `internal/project`, the spec skill, and READMEs.
- Before cloud-server work, consult the server design decisions in local `docs/PRODUCT_DESIGN.md`, when available.

## Testing Guidelines

Use standard `testing`, `*_test.go`, `TestXxx`, table-driven cases, and `t.TempDir()` fixtures. Cover regressions; no numeric coverage threshold exists. Manually verify rendering at `/` and `/endpoints`, including search, navigation, and light/dark themes.

## Commit & Pull Request Guidelines

Follow short, imperative commit summaries, e.g. `Add spec quality lint hints`. Keep PRs focused; describe behavior changes, link relevant issues, report validation, and include screenshots for UI changes. Update both READMEs when documented usage changes.
