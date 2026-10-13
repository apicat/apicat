# apicat

English | [简体中文](README.zh-CN.md)

**API documentation infrastructure for the agent era.**

AI agents are becoming the primary authors — and consumers — of code. When an agent integrates with an API, it needs structured, machine-readable documentation. Existing API doc tools are built for humans first (fancy UIs, "Try it" buttons) and treat the spec as an afterthought.

apicat inverts that:

> **The OpenAPI spec is the first-class citizen. HTML rendering is just one view of it.**

Agents read the raw spec. Humans read a clean, server-rendered HTML view of the same source. Both stay in sync because there is only one source.

## Features

- **One command, zero setup** — point it at a directory containing an OpenAPI spec and get browsable docs
- **Multi-file specs** — `$ref` across files resolved out of the box (OpenAPI 3.0 / 3.1)
- **Single binary** — Go-compiled, all templates and assets embedded; no Node, no Python, no runtime dependencies
- **Server-side rendered** — fast, clean HTML designed for low reading fatigue; the rendered page doubles as a quality check for your spec (missing descriptions, incomplete examples stand out)
- **Search** — filter endpoints by method, path, summary, operationId, or description
- **Light / dark theme** — follows your system by default, with a manual override
- **Live reload** — `-watch` reloads the spec when files change; a broken edit keeps the last good version
- **Lint hints** — missing summaries, descriptions, operationIds, and error responses are counted on the overview page and marked inline where they occur

## Install

```bash
go install github.com/apicat/apicat/v3/cmd/apicat-cli@latest
```

Or build from source:

```bash
git clone https://github.com/apicat/apicat.git
cd apicat
go build -o apicat-cli ./cmd/apicat-cli
```

## Usage

apicat is project-based: specs live centrally under `~/.apicat/projects/`, one directory per project, keyed by the project's absolute path with every character outside `[A-Za-z0-9_-]` replaced by `-` (e.g. `/Users/alice/dev/shop` → `-Users-alice-dev-shop`). Run from your project root and apicat finds its spec directory automatically:

```bash
cd path/to/your-project
apicat-cli
```

Or point at any spec directory explicitly:

```bash
apicat-cli path/to/spec-dir
```

Either way, the spec directory must contain an `openapi.yaml`, `openapi.yml`, or `openapi.json` entry file; referenced files are resolved relative to it. Then open http://127.0.0.1:8080.

Flags:

| Flag | Default | Description |
|------|---------|-------------|
| `-port` | `8080` | Port to listen on |
| `-host` | `127.0.0.1` | Host to bind to |
| `-watch` | off | Reload when spec files in the directory change |

## Agent Skill

[`skills/apicat-spec/`](skills/apicat-spec/SKILL.md) is an [Agent Skill](https://code.claude.com/docs/en/skills) that teaches a coding agent to write OpenAPI specs the way apicat expects them — stored in apicat's project directory and rendering with zero lint hints. To use it in your project:

```bash
mkdir -p .claude/skills
cp -r skills/apicat-spec .claude/skills/
```

Then ask your agent to create or update your API spec; it picks up the conventions automatically: where the spec lives (`~/.apicat/projects/<encoded-path>/`), entry-file naming, multi-file `$ref` layout, required descriptions, error responses, and examples. The skill file is plain Markdown, so it also works with other tools that read skill/instruction files.

## Roadmap

- **Cloud** — `apicat-cli publish` to a self-hostable server: shared docs, teams, versioned specs (`latest` + pinned tags)
- **Agent-native access** — raw spec endpoints, API discovery/search, MCP server for Claude Code / Cursor and friends

## v3 rewrite

This is a ground-up rewrite. The previous 2.x codebase is preserved on the [`old`](https://github.com/apicat/apicat/tree/old) branch; 2.x releases remain available under [tags](https://github.com/apicat/apicat/tags).

## License

[AGPL-3.0](LICENSE)
