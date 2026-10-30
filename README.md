# apicat

English | [简体中文](README.zh-CN.md)

**API documentation infrastructure for the agent era.**

AI agents are becoming the primary authors — and consumers — of code. When an agent integrates with an API, it needs structured, machine-readable documentation. Existing API doc tools are built for humans first (fancy UIs, "Try it" buttons) and treat the spec as an afterthought.

apicat inverts that:

> **The OpenAPI spec is the first-class citizen. HTML rendering is just one view of it.**

Agents read the raw spec. Humans browse and debug the same source in a local React workbench. Both stay in sync because there is only one source.

## Features

- **Languages** — Simplified Chinese and English, with browser-language detection and a persistent manual choice in the sidebar. UI messages are translated; API terminology, schema types, document content and request/response data stay unchanged.
- **One command, zero setup** — point it at a directory containing an OpenAPI spec and get browsable docs
- **Multi-file specs** — `$ref` across files resolved out of the box (OpenAPI 3.0 / 3.1)
- **Single binary** — Go-compiled, the production React bundle is embedded; no Node or Python required at runtime
- **React workbench** — React + TypeScript + Vite + shadcn/ui + Tailwind CSS, served by Echo v5; overview, endpoint documentation, reusable definitions, environments, and HTTP debugging
- **Shared definitions** — browse `components` from the sidebar or `/components`, including schemas, parameters, request bodies, responses, and other reusable definitions
- **Direct links** — endpoint pages use `/endpoints/{operationId}`. Missing, duplicate, or non-URL-safe IDs receive a stable generated ID; old query/hash links are not supported.
- **Search** — filter endpoints by method, path, summary, operationId, or description; find shared definitions by category, name, or description
- **Collapsible navigation** — API tags and component categories start collapsed; search reveals matching groups and clearing it restores their state
- **Light / dark theme** — follows your system by default, with a manual override
- **Live reload** — `-watch` reloads the spec when files change; a broken edit keeps the last good version
- **Lint hints** — missing summaries, descriptions, operationIds, and error responses are counted on the overview page and marked inline where they occur

## Install

```bash
go install github.com/apicat/apicat/v3/cmd/apicat-cli@latest
```

Or build from source (Go 1.25+ and Node 24+ / npm):

```bash
git clone https://github.com/apicat/apicat.git
cd apicat
npm --prefix web ci
npm --prefix web run build
go build -o apicat-cli ./cmd/apicat-cli
```

## Usage

apicat is project-based: by default, specs live centrally under `~/.apicat/projects/`, one directory per project, keyed by the project's absolute path with every character outside `[A-Za-z0-9_-]` replaced by `-` (e.g. `/Users/alice/dev/shop` → `-Users-alice-dev-shop`). Run from your project root and apicat finds its spec directory automatically:

```bash
cd path/to/your-project
apicat-cli
```

To configure a project's spec directory, create `.apicat/setting.json` in the project root:

```json
{
  "spec_dir": "docs/api"
}
```

Relative paths resolve against the project root; absolute paths are also supported. Run `apicat-cli`, `apicat-cli -watch`, or `apicat-cli check` from that root to use the setting. Read `<project>/.apicat/setting.json` first. Only when that file is absent, read `~/.apicat/projects/<encoded-path>/setting.json`. Both configs resolve relative paths against the project root; parent directories are not searched. Paths are used literally, without `~` or environment-variable expansion.

If neither config exists, or the selected config has an unset, empty, or whitespace-only `spec_dir`, the existing `~/.apicat/projects/<encoded-path>/` default is used. An existing local config replaces the fallback config entirely. Invalid JSON, unknown fields, invalid field types, and unreadable configs report an error without trying another config. Configuration selects the directory; it does not create it or move existing specs. Restart the CLI after changing the config; `-watch` only reloads specs.

An explicit spec directory overrides the config (and skips reading it):

```bash
apicat-cli path/to/spec-dir
```

Either way, the spec directory must contain an `openapi.yaml`, `openapi.yml`, or `openapi.json` entry file; referenced files are resolved relative to it. Then open http://127.0.0.1:20317.

Flags:

| Flag | Default | Description |
|------|---------|-------------|
| `-port` | `20317` | Port to listen on |
| `-host` | `127.0.0.1` | Host to bind to |
| `-watch` | off | Reload when spec files in the directory change |

## Local workbench

- **Environments** (`/environments`): create development, test, and production environments with a base URL and variables. Global variables apply to every environment; enabled environment variables override matching global names. Use `{{name}}` in request URLs, parameter values, headers, and bodies. Values are substituted literally, not evaluated as scripts.
- **HTTP debugging**: open an endpoint, select **接口调试**, edit path/query parameters, headers, authentication and raw body, then send. Requests execute in Go, so target APIs do not need browser CORS configuration. Supports HTTP(S), Bearer headers, JSON/text/URL-encoded bodies, cancellation, configurable timeouts, response headers/body, and binary responses shown as Base64. No automatic redirects, application-level retries, or shared cookie jar. Responses are capped at 2 MiB and marked when truncated.
- **Parameter handling**: inherits path-level parameters and operation/path-level servers. Common query arrays/objects support form, spaceDelimited, pipeDelimited and shallow deepObject. Complex path serialization and multipart uploads are not automated; unsupported cases are reported. You can edit already-serialized values directly in the URL. This is interactive debugging, not an assertion-based automated test runner.
- **Storage**: environments are saved to `<spec-directory>/.apicat/workspace.json` with owner-only file permissions, atomic replacement, and revision conflict detection within one running workbench. UI “hide” masks values; it does not encrypt the file. Keep this file out of version control if it contains credentials. Request drafts and responses stay in page memory and are not persisted as history.
- **Local access**: the API validates Host/Origin and requires a session token for mutations and request execution. The default listener remains `127.0.0.1`; this local workbench is not a multi-user authenticated server.

### Frontend development

Start Go with `go run ./cmd/apicat-cli -watch /path/to/spec-dir`, then run `npm --prefix web run dev`. Vite serves the development UI and proxies `/api` to `127.0.0.1:20317`. Change the proxy in `web/vite.config.ts` if using another backend port.

`npm --prefix web run build` performs TypeScript checking and builds `web/dist`; rebuild the Go binary to embed those assets. The generated bundle is committed so `go install` and Go-only builds continue to work. Frontend source changes must include the regenerated bundle. `make build` installs locked frontend dependencies, builds assets, and builds the CLI. Run `npm --prefix web test` and `go test ./...` for regression checks. Spec watch reloads the document in the browser without discarding request drafts; failed spec loads keep the last good document.

## Check documentation

Check a spec without starting the web server:

```bash
apicat-cli check                         # current project's spec directory
apicat-cli check /path/to/spec-dir       # explicit directory
apicat-cli check -json                   # structured output for agents and CI
apicat-cli check -json /path/to/spec-dir
```

Put flags before the directory. A directory literally named `check` can still be served as `apicat-cli ./check`.

The command loads the spec, resolves references, and runs the same six documentation-quality rules as the UI. Each operation is checked once even if it has multiple tags. Reports include API-level and endpoint-specific hints; field-description hints aggregate missing fields per operation.

Exit codes: **0** = loaded with no quality hints; **1** = quality hints found; **2** = loading, argument, or output error. `-json` writes one JSON object to stdout with `status` (`ok`, `hints`, or `error`), `spec_dir`, `errors`, `lint_total`, and `hints` (`method`, `path`, `message`; API-level hints omit method/path). Usage errors and parser diagnostics go to stderr. On loading failure, zero lint counts mean lint was not run; inspect `status` and `errors` first.

This is a loading and documentation-quality check, not complete OpenAPI specification validation. It does not verify unique operationIds, agreement with implementation, unused schemas, or schema fields beyond the rendering depth limit. AI authors should run it after writing or updating a spec, fix issues within scope, and rerun it before reporting completion.

## Agent Skill

[`skills/apicat-spec/`](skills/apicat-spec/SKILL.md) is an [Agent Skill](https://code.claude.com/docs/en/skills) that teaches a coding agent to write OpenAPI specs the way apicat expects them — stored in apicat's project directory and rendering with zero lint hints. To use it in your project:

```bash
mkdir -p .claude/skills
cp -r skills/apicat-spec .claude/skills/
```

Then ask your agent to create or update your API spec; it picks up the conventions automatically: where the spec lives (local or fallback `setting.json` → `spec_dir`, defaulting to `~/.apicat/projects/<encoded-path>/`), entry-file naming, multi-file `$ref` layout, required descriptions, error responses, and examples. The skill file is plain Markdown, so it also works with other tools that read skill/instruction files.

## Roadmap

- **Cloud** — `apicat-cli publish` to a self-hostable server: shared docs, teams, versioned specs (`latest` + pinned tags)
- **Agent-native access** — raw spec endpoints, API discovery/search, MCP server for Claude Code / Cursor and friends

## v3 rewrite

This is a ground-up rewrite. The previous 2.x codebase is preserved on the [`old`](https://github.com/apicat/apicat/tree/old) branch; 2.x releases remain available under [tags](https://github.com/apicat/apicat/tags).

## License

[AGPL-3.0](LICENSE)
