---
name: apicat-spec
description: >-
  Create and update API documentation as OpenAPI 3.x specs for use with apicat.
  Use when users ask to create, generate, edit, or synchronize API/interface
  documentation, including documenting existing endpoints or updating docs
  after API changes, even if they do not mention OpenAPI or apicat.
  Do not use for merely explaining or consuming an API, or when the user
  explicitly requests a different documentation format or tool.
---

# Authoring OpenAPI specs for apicat

Create accurate OpenAPI documentation for the requested API scope. apicat displays the spec in a React workbench and checks selected documentation quality rules. Aim for zero lint hints, but do not treat that as proof of completeness or correctness.

## Establish the source and scope

- For an implemented API, inspect relevant routes, handlers, request validation, data types, authentication, and tests before writing. Use existing documentation as context and check it against the implementation.
- For a proposed API, follow the user's agreed design. If sources conflict or a missing detail affects the contract, ask a focused question rather than inventing behavior.
- Document only supported methods, paths, parameters, request/response shapes, status codes, and authentication requirements. Do not invent error responses, constraints, or business descriptions to clear lint hints.
- For updates, change only the requested endpoints and affected shared schemas. Preserve unrelated content, file organization, and existing `operationId` values unless the requested change requires otherwise. Check other consumers before editing a shared schema; do not modify application behavior to match the documentation.
- Before generating a new spec, identify the business areas, approximate endpoint count, shared models, and large examples; choose the file layout below before writing the full document.
- When updating an existing spec, run `apicat-cli check -json` from the project root before editing when available, so pre-existing hints can be distinguished from new ones.

## Where the spec lives

apicat resolves the spec directory in this order:

1. An explicit spec directory supplied by the user or to the CLI.
2. `spec_dir` in `.apicat/setting.json` in the project root (the directory where `apicat-cli` will run). Relative paths resolve against that root; absolute paths are also supported. Do not search parent directories or expand `~` or environment variables.
3. Only if the local config file is absent, read `~/.apicat/projects/<encoded-path>/setting.json`. Relative paths in this fallback config also resolve against the project root. An existing local config replaces the fallback entirely; do not merge fields.
4. If neither config exists, or the selected config has an unset, empty, or whitespace-only `spec_dir`, use `~/.apicat/projects/<encoded-path>/`. Encode the project's absolute path by replacing every character outside `[A-Za-z0-9_-]` with `-`:

```
/Users/alice/dev/shop  →  ~/.apicat/projects/-Users-alice-dev-shop/
```

Example project-root `.apicat/setting.json`:

```json
{
  "spec_dir": "docs/api"
}
```

Workflow:

1. Determine the project root and resolve the spec directory using the precedence above. Invalid JSON, unknown fields, invalid field types, or unreadable configuration must be reported rather than trying a lower-priority config or silently using the default. An explicit CLI directory skips reading the config.
2. If that directory already contains an entry file (`openapi.yaml`, `openapi.yml`, `openapi.json`), **update it in place** — never create a second copy elsewhere.
3. Otherwise create the directory and write the spec there. Changing configuration does not move existing specs; only migrate existing documentation when requested.
4. Run checks and previews from the project root. If using an explicit directory, pass it to both commands, with flags before the directory.

Running `apicat-cli` with no directory argument from the project root resolves the configured or default directory automatically. Restart the CLI after changing either `setting.json`; `-watch` only reloads specs.

## OpenAPI version and specification references

- For a new spec, default to `openapi: 3.1.2`, unless the project's toolchain requires another supported version.
- For an existing spec, preserve its exact `openapi` version, including 3.0.x. Upgrade only when the user requests a migration; do not silently change versions during documentation edits.
- This workflow targets OpenAPI 3.0 and 3.1. Do not select a newer release just because it is the latest.
- Use version-appropriate syntax: in 3.1, a nullable string can use `type: [string, "null"]`; in 3.0, use `type: string` with `nullable: true`. Quote `"null"` in YAML type declarations.

Consult the official specification for the document's declared version when uncertain about `$ref`, schemas, nullability, or security. Read relevant sections as needed, not the entire specification on every task:

- [OpenAPI 3.1.2](https://spec.openapis.org/oas/v3.1.2.html) — the default for new specs.
- [OpenAPI 3.0.3](https://spec.openapis.org/oas/v3.0.3.html) — reference for existing 3.0.3 specs.
- [Official specification index](https://spec.openapis.org/oas/) — locate other declared versions without changing the document's version.

The specification defines valid syntax and semantics; apicat's lint checks and rendering support do not replace it.

## File layout

- The spec directory must contain an entry file named `openapi.yaml`, `openapi.yml`, or `openapi.json` — checked in that order, first match wins.
- Keep small, cohesive APIs in one file. For a new, substantial API spanning multiple business areas, default to multiple files grouped by business area unless the user or consuming tool requires a single file. Do not first generate thousands of lines in one file and leave organization for later. Line count is a review signal, not a fixed splitting threshold; consider endpoint complexity and example size as well.
- Keep global metadata, servers, authentication, tags, and the path index in the entry file. Reference Path Item Objects in business files such as `paths/users.yaml` and `paths/orders.yaml`; keep related operations together. Use stable, descriptive fragment keys so the entry file makes each endpoint's definition easy to locate.
- Extract genuinely shared domain models into named schemas, exposed through `components/schemas` in the entry file. Keep simple, endpoint-specific request/response shapes near their operations. Matching fields alone do not make two business concepts the same model.
- Avoid one file per field or trivial parameter, arbitrary numbered chunks, and chains of references that only forward to another reference. Split a business file further only when it becomes difficult to navigate. Prefer a small set of related files that provides enough context to understand an endpoint.
- Keep examples representative and concise; do not copy large datasets merely to illustrate an array. Preserve user-required examples and meaningful variants. If examples dominate file size, address their organization before fragmenting otherwise cohesive operations. Never omit contract details or minify the document to meet a size target.
- For scoped updates, use the existing layout, including an existing large single file. If organization obstructs the task, report a focused restructuring recommendation instead of silently rewriting unrelated content. When restructuring is requested, preserve API semantics and operationIds and update all affected references.

Example layout for a larger API (YAML shown; JSON can use the same organization). Entry-file `$ref` paths start from the spec directory:

```
~/.apicat/projects/-Users-alice-dev-shop/
  openapi.yaml          # global settings, path index, named shared components
  paths/users.yaml      # e.g. paths/users.yaml#/collection and #/byId
  paths/orders.yaml     # related order operations
  schemas/user.yaml     # e.g. schemas/user.yaml#/User
  schemas/order.yaml    # shared order models
```

Use `$ref` only where the declared OpenAPI version permits it: each entry under `paths` may reference a Path Item Object; the entire `paths` map is not a generic include. Check fragment targets and reference bases, especially after moving definitions. A fragment-only reference in an external file refers to that file, not automatically to the entry document.

For existing docs, inspect the entry file and search for relevant paths or operationIds first. Read the affected operations, their referenced schemas, and applicable global settings; check consumers before changing shared definitions. Avoid repeatedly loading unrelated files or the entire large document. Splitting saves context only when reads are selective; full-document verification still needs to cover all relevant references.

## Quality rules

apicat aggregates these six checks in the "Spec Quality" panel and adds endpoint hints where applicable. These checks run on the converted view model, not every definition in the source spec:

1. **API description** — `info.description` must be non-empty. Rendered as CommonMark markdown on the overview page; a few sentences on what the API is for, auth model, and conventions.
2. **Operation summary** — every operation needs a `summary` (one line) or `description` (longer, markdown). Prefer both: summary for the sidebar, description for the detail page.
3. **operationId** — every operation needs a non-empty, document-wide unique `operationId`. Follow the existing style (e.g. `listUsers`, `createUser`); for new IDs, prefer descriptive camelCase using ASCII letters and digits. This ID is searchable and forms the endpoint page URL, e.g. `/endpoints/getArticleById`, so preserve existing IDs when updating an operation. Direct use in URLs requires matching `^[A-Za-z0-9_-][A-Za-z0-9_.-]*$`; missing, duplicate, or other IDs receive generated fallback route IDs without changing the spec. This fallback does not replace the authoring requirement. The lint check detects missing IDs, not duplicates or URL suitability.
4. **Error responses** — every operation must document at least one error response: a `4xx`, a `5xx`, or `default`. Document the errors a caller must actually handle (e.g. `400`, `404`, `422`), not a token `default` added to silence the rule.
5. **Parameter descriptions** — every parameter (path, query, header, cookie) needs a `description`.
6. **Schema field descriptions** — every named property in request-body and response schemas needs a `description`. apicat checks the converted tree recursively, including `allOf` / `oneOf` / `anyOf` branches and array `items`, within its conversion depth limit.

For new specs, target zero hints across the document. For scoped updates, resolve hints in the changed content without expanding into unrelated cleanup; report remaining pre-existing hints. If clearing a hint requires an unsupported assumption, preserve accuracy and report what needs clarification.

## Useful metadata and rendering limits

Include metadata when supported by the API contract, even when the current HTML view does not display it:

- **`example`, `default`, and `enum`** — provide representative examples, actual defaults, and closed value sets. Examples, defaults, and enums are displayed in the workbench. Media-type examples and schema examples also seed the request editor; use realistic values.
- **Constraints** — encode actual formats, patterns, length/range/item limits, and nullability using syntax appropriate to the spec's OpenAPI version. The schema view displays the supported annotations; rendering is not a complete validation of those constraints.
- **`required`** — the object-level `required` list, and `required: true` on parameters.
- **`deprecated: true`** — record actual deprecations. Operations and schema fields display a badge; parameter deprecation is read but currently has no badge in the parameter table.
- **Component names** — put reusable schemas in `components/schemas`; the component name (e.g. `User`) is displayed where the schema is referenced. Name them meaningfully.
- **Tags** — give every operation a tag, and declare each tag at the root `tags:` list with a `description`. Untagged operations fall into a catch-all "default" group.
- **`servers`** (with descriptions) and **`components/securitySchemes`** render on the overview page.

Structural limit: conversion stops beyond depth 10, so deeper fields are neither rendered nor included in these lint counts. `$ref` is expanded during conversion; splitting files or extracting components does not reduce nesting depth. Preserve the real API structure rather than flattening it to fit the renderer, and report relevant display limitations.

## Verify

Compare the resulting spec against the requested scope and its sources: confirm endpoint coverage, methods/paths, parameters, required fields, request/response shapes, status codes, authentication, and unique `operationId` values. Check changed references, including consumers of shared schemas. Successful loading and zero hints do not replace this review.

For new or explicitly restructured specs, review the layout before finishing: large unrelated business areas should not remain in one monolithic file without a concrete reason, and ordinary endpoint edits should not require traversing many tiny files. Confirm every intended operation is reachable from the entry file and shared definitions remain consistent. Run checks against the entry document's directory, not just individual fragments; apicat's lint rules do not assess file organization.

After creating or updating a spec, run the check command from the project root:

```bash
apicat-cli check -json
```

It resolves the configured or default project directory automatically, reports loading/reference errors and documentation-quality hints, and exits without starting a server. Read stdout as JSON and also capture the exit code and stderr:

- **0**: loaded successfully with no quality hints.
- **1**: quality hints remain. Use `hints` entries (`method`, `path`, `message`) to locate affected operations; API-level hints omit method/path. Field-description hints give a count per operation, so inspect that operation's request/response schemas to locate the fields.
- **2**: loading, argument, or output error. Read `errors` in the JSON report and stderr diagnostics. Usage errors may produce only stderr. Fix the cause before treating lint results as meaningful.

Reports include `status`, `spec_dir`, `errors`, `lint_total`, and `hints`. Check `status` and `errors` first: a failed load has not run lint, even though its count is zero. Fix issues within the requested scope and rerun after each correction until the quality target is met. If blocked by missing contract information or unrelated existing hints, report them rather than inventing content or expanding the task. A successful check does not validate the entire OpenAPI specification, unused schemas, operationId uniqueness, or agreement with implementation.

When visual verification is useful and tools are available, run a preview from the project root, reusing an existing preview for the same project when possible:

```bash
apicat-cli -watch # resolves local/fallback setting.json spec_dir or the default from the cwd
```

- Overview `http://127.0.0.1:20317` — inspect the Spec Quality panel against the quality target above.
- Endpoints `http://127.0.0.1:20317/endpoints` — inspect changed operations, field descriptions, and hint badges.
- With `-watch`, a broken edit keeps the last good version. Confirm the latest edit loaded successfully in the terminal; an unchanged page is not evidence of success.
- If the port is occupied, use an available one with `apicat-cli -watch -port 20318` and adjust the preview URLs; do not stop an unrelated server.

If the CLI is unavailable or does not support `check`, perform the source and reference checks possible in the environment and report the command check as unverified. Missing browser tools do not prevent running `check`; report visual verification separately. Do not claim a clean preview or zero hints without checking them.

Finish by reporting the spec location, changes made, checks performed, and any remaining hints or unresolved contract details. For new or restructured specs, briefly explain the chosen file layout and any reason for keeping a large single file. Completion requires documentation that covers the requested scope and matches the implementation or agreed design, with verification results and limitations stated accurately.
