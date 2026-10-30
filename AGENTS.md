# Repository Guidelines

## Project Structure & Module Organization

`cmd/apicat-cli/main.go` wires the OpenAPI documentation CLI together. `internal/` contains project resolution, loading, conversion, models, linting, workspace storage, HTTP execution, Echo serving, and watching. The React API depends on `internal/model`, not libopenapi types.

`web/src/` holds the React + TypeScript workbench, shadcn/ui components and Tailwind CSS. Vite builds `web/dist/`, embedded by `web/embed.go`. Commit regenerated dist assets with frontend changes so Go-only builds work. Tests sit beside their packages. `skills/apicat-spec/` documents spec-authoring conventions. Local `testdata/` and `docs/` are gitignored.

## Build, Test, and Development Commands

Use Go 1.25.0 or newer. Frontend development/builds require Node 24+ and npm; end users only need the CLI binary.

- `make build`: install locked frontend dependencies, build the UI, and build the standalone CLI.
- `npm --prefix web run dev`: Vite development server, proxying API calls to port 8080.
- `npm --prefix web test`: frontend regression tests.
- `go build -o apicat-cli ./cmd/apicat-cli`: build using the committed frontend bundle.
- `go build ./...`: compile all packages.
- `go test ./...`: run all unit tests.
- `go run ./cmd/apicat-cli check -json /path/to/spec-dir`: check loading and documentation quality without serving HTML.
- `go run ./cmd/apicat-cli -watch /path/to/spec-dir`: serve docs at `http://127.0.0.1:8080` and reload spec changes.

Provide `openapi.yaml`, `openapi.yml`, or `openapi.json`. Without a directory argument, read `spec_dir` from `<cwd>/.apicat/setting.json`, or only if absent from `~/.apicat/projects/<encoded-cwd>/setting.json` (both relative to cwd); if unset, specs resolve under `~/.apicat/projects/<encoded-cwd>/`; encoding replaces characters outside `[A-Za-z0-9_-]` with `-`. Restart/rebuild after editing `web/`; `-watch` only reloads specs.

## Coding Style & Naming Conventions

Use `gofmt`, tabs in Go, lowercase package names, exported `PascalCase`, and unexported `camelCase`. Match two-space JavaScript/CSS indentation. No dedicated linter configuration exists.

## Architecture & Change Contracts

- Apply the single responsibility principle pragmatically: separate responsibilities that change for different reasons only when doing so does not increase the effort to understand the code or the risk of errors. Prefer boundaries that make bugs easy to locate; avoid unnecessary layers and fragmentation.
- Give independent functions and modules explicit inputs, outputs, and failure contracts so they can fulfill their responsibilities without relying on hidden upstream state or downstream behavior. Keep workflow coordination at the caller and prevent unrelated failures from blocking independent work where safe. When a required dependency fails, report the failure with enough context to identify its source; do not hide it or claim success. Preserve cancellation, ordering, and atomicity guarantees.
- Add explicit locks only where required to meet actual consistency requirements. Minimize lock scope and transaction duration; keep unrelated work and external calls outside transactions. Prefer atomic updates, unique constraints, or conditional updates when they correctly enforce the required guarantees without additional explicit locks.
- Match consistency safeguards to the impact and likelihood of inconsistency. For low-impact, recoverable cases, including infrequent operations, it is acceptable to omit explicit locks when product workflows can prevent or contain conflicts, or backend/manual repair is a practical remedy. Make the tolerated inconsistency and prevention or recovery path clear. Low frequency alone does not justify weakening required consistency or existing atomicity guarantees; avoid adding concurrency complexity when its cost exceeds the practical benefit.

- Localize application-owned UI messages through `web/src/lib/i18n.ts` and `locales/en.ts`. Preserve API/root/Schema types/HTTP methods, identifiers, spec text and request/response data. Backend validation uses `internal/problem` message keys and separate values; do not translate raw API payloads or persist translated user data on locale changes.

- Use the official light/dark PNGs in `web/src/assets/brand/` without redrawing, recoloring, or cropping them. Theme tokens live in `web/src/index.css`; brand orange uses dark button text and must stay separate from status colors.

- Load → convert → lint → render/serve. Watch reloads must atomically replace the document only on success, retaining the last good version on failure.
- Endpoint navigation uses `/endpoints/:operationId` with the converter-assigned `RouteID`: unique URL-safe operationIds are preserved; missing, duplicate, or unsafe IDs use stable generated IDs. React Router owns UI navigation; do not restore legacy query/hash routes.
- New pages require a React route and an explicit Echo SPA route. Keep `/api/*` and missing assets from falling back to HTML. Business logic must not depend on Echo contexts.
- Workspace data lives in `<spec-dir>/.apicat/workspace.json`, independently from path settings. Preserve atomic writes, revision checks, same-origin/session protection, cancellation, and response limits.
- New lint rules must update `ruleOrder` and the spec skill's Quality rules. Field-level missing-description markers are handled directly by React components.
- Encoding changes require synchronized updates to `internal/project`, the spec skill, and READMEs.
- Before cloud-server work, consult the server design decisions in local `docs/PRODUCT_DESIGN.md`, when available.

## Testing Guidelines

Use standard `testing`, `*_test.go`, `TestXxx`, table-driven cases, and `t.TempDir()` fixtures. Cover regressions; no numeric coverage threshold exists. Manually verify rendering at `/` and `/endpoints`, including search, navigation, and light/dark themes.

## Commit & Pull Request Guidelines

Follow short, imperative commit summaries, e.g. `Add spec quality lint hints`. Keep PRs focused; describe behavior changes, link relevant issues, report validation, and include screenshots for UI changes. Update both READMEs when documented usage changes.
