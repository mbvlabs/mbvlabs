---
name: andurel
description: Use this skill for Andurel framework projects when deciding where code belongs, adding or changing resources, controllers, models, services, routes, templ views, Inertia/Vue screens, background jobs, migrations, config, clients, or framework-adjacent internals. Focuses on project structure, layer placement, command discovery, generator workflows, and agent-safe CLI usage.
---

# Andurel

Use this skill when working in an Andurel project or generating Andurel code. It helps place code in the right layer and use the `andurel` CLI safely.

## Agent Invariants

- Prefer `andurel --agent --help` and `andurel commands --json` for discovery.
- Run `andurel project info --json` before generation.
- Use `--json` or `--jq` when extracting data.
- Use `--dry-run --json` before mutating commands when intent is uncertain.
- Inspect returned artifact arrays before assuming which files changed.
- Treat `andurel generate ...` as the default path for models, controllers, Inertia pages, route constants, jobs, email templates, factories, and route helpers.
- Do not hand-create files that an Andurel generator can create unless the generator is unavailable, failed for a concrete reason, or produced an unsuitable baseline; state the exact command tried and the concrete reason before manual creation.
- For scheduled, periodic, reminder, async workflow, or queue-triggered behavior, use `andurel generate job NAME --dry-run --json --diff` and then `andurel generate job NAME --json` for the job args, worker, and registration baseline. Patch the generated worker for business-specific orchestration.
- After adding or changing Inertia routes, run `andurel generate routes --json` so frontend pages can import `resources/js/routes.ts`.
- When building URLs with query parameters, use route helper options such as `routing.QueryParam(...)` and route `URL`/`FullURL` methods instead of manual string concatenation.
- Follow the repository rules for verification.
- Prefer the local project pattern over a generic Rails, Echo, Bun, Templ, or Vue convention.
- Keep controllers as HTTP adapters: parse input, call models or services, map errors, and render a response.
- Create a service only when there is real application orchestration, not just because code exists.

## Read When Placing Code

Read [references/layer-placement.md](references/layer-placement.md) before adding or moving behavior across models, services, controllers, routes, views, queue jobs, config, clients, or internal packages.

## First Pass

1. Inspect the existing resource closest to the requested change.
2. Identify the delivery surface: public hypermedia page, admin Inertia page, API endpoint, background job, email, or CLI/command.
3. Identify the domain object or workflow being changed.
4. Keep changes in the smallest layer that can own the behavior honestly.
5. Use the CLI discovery commands before generating or mutating project files.

## Generator Gate

Before manually creating or registering framework-owned files, you MUST run CLI discovery and a generator dry-run when a generator plausibly exists.

This gate applies to:

- Models, controllers, and API controllers.
- Inertia/Vue/React pages.
- Routes, route helper files, and route registrations.
- Queue jobs/workers.
- Email templates.
- Factories.

Required sequence:

1. Run `andurel project info --json`.
2. Run `andurel commands --json` or the relevant `andurel generate ... --help`.
3. Run the closest generator with `--dry-run --json --diff` when available.
4. Inspect the returned artifacts and diff before deciding.
5. Only hand-create files if the generator is unavailable, fails, or cannot express the requested shape. Before manual creation, state the exact command tried and the concrete reason manual work is required.

## Layer Placement

- Put invariant business rules, domain validation, entity construction, persistence methods, and finder/query methods in `models/`.
- Put test factory definitions and factory helpers in `models/factories/`.
- Put transactions, cross-model coordination, external side effects, and multi-step application workflows in `services/`.
- Put HTTP-specific concerns in `controllers/`, `controllers/admin/`, or `controllers/api/`.
- Put route names, route paths, and URL builders in `router/routes/`.
- Put templ rendering helpers and presentation-specific adapters in `views/`.
- Put admin Inertia pages and reusable Vue components in `resources/js/`.
- Put River job argument types in `queue/jobs/` and worker implementations or registration in `queue/`.
- Put provider adapters in `clients/`, email templates/helpers in `email/`, and config/environment loading in `config/`.
- Put reusable framework-like support that is independent of one resource in `internal/`.
- Register new constructors in the existing `fx` modules for the package that owns them.

## Output Modes

Use structured output by default when automating:

| Flag | Use |
|------|-----|
| `--json` | Full `{ok,data,summary,breadcrumbs}` envelope |
| `--agent` | Structured output with non-essential human progress suppressed |
| `--jq '.field.path'` | Built-in simple field-path extraction |
| `--quiet` | Suppress human-only output |
| `--md` | Markdown output where supported |

Structured failures include `ok:false`, a stable `code`, `error`, optional `hint`, and `exit_code`. Prefer the `hint` and `breadcrumbs` fields over guessing the next command.

## Common Workflows

Inspect a project:

```bash
andurel project info --json
andurel routes --json
andurel models --json
andurel migrations --json
andurel commands --json
```

Preview scaffold generation:

```bash
andurel generate scaffold Product --dry-run --json
```

Generate and review artifacts:

```bash
andurel generate scaffold Product --json
```

Resource scaffold workflow:

1. Run `andurel project info --json` and `andurel commands --json` before generation.
2. For new resources, preview the appropriate generator with `--dry-run --json --diff` before manual file creation.
3. If scaffold/model generation fails because the table is missing, create only the migration manually, then rerun the dry-run.
4. If the dry-run succeeds, run the generator for the baseline and patch only deltas the generator cannot express, such as custom routes, extra methods, specialized UI, or workflow-specific behavior.
5. If the generator is not suitable, record the specific reason and then follow local patterns manually.

API endpoint workflow:

1. Do not start by hand-creating `controllers/api/*.go`, `router/routes/api_*.go`, or controller registrations.
2. Run `andurel project info --json`, `andurel commands --json`, `andurel routes --json`, and `andurel controllers --json`.
3. Attempt the closest controller or scaffold generator with `--dry-run --json --diff` when available.
4. If the generator cannot create the desired API shape, use the generated or local controller pattern manually, but first state the exact dry-run command tried and why it was insufficient.

Generate Inertia route helpers:

```bash
andurel routes --json
andurel generate routes --json
```

`andurel generate routes` reads `router/routes/*.go` as the source of truth and writes `resources/js/routes.ts`. It only runs when `andurel.lock` has `scaffoldConfig.inertia` set to `vue` or `react`. Import helpers from that file in Vue or React pages instead of hard-coding URLs.

Build URLs from route helpers:

```go
routes.AdminDiaryEntryToday.URL(routing.QueryParam("focus", "morning"))
routes.AdminDiaryEntryToday.FullURL(config.BaseURL, routing.QueryParam("focus", "morning"))
```

Do not append `?key=value` manually when a route helper can accept `routing.QueryParam`.

Generate a queue job:

```bash
andurel generate job SendWelcomeEmail --dry-run --json --diff
andurel generate job SendWelcomeEmail --json
```

Use this for new queued workflows even when the worker delegates to an existing lower-level job. For example, a periodic reminder should enqueue a reminder-specific job, and that reminder worker may enqueue `SendTransactionalEmailArgs` for delivery.

Check or sync factories:

```bash
andurel generate factory Product --check --json
andurel generate factory Product --sync --json
andurel generate factories --check --json
andurel generate factories --sync --json
```

Factory guidance:

1. Treat model `Entity` structs as the source of truth for generated factory fields.
2. Keep reusable test data builders in `models/factories/`.
3. Prefer `andurel generate factory NAME --check --json` before editing factory files by hand.
4. Use `--sync` to update Andurel generated regions and preserve custom helpers outside those regions.
5. Pass `--skip-factory` only when a generated model or scaffold should intentionally omit a factory.

Generate a named database seed:

1. Inspect the relevant models and existing factories in `models/factories`.
2. Add a seed function to `database/seeds`, using only exported model/factory/storage primitives.
3. Register it in `seeds.Registry` with a stable lowercase name.
4. Keep the seed idempotence expectations explicit in code comments when it may be re-run.
5. Verify the seed is discoverable:

```bash
andurel database seed --list
andurel database seed development
andurel database seed test
```

Check project health:

```bash
andurel doctor --json
```

In Inertia projects, `doctor` checks whether `resources/js/routes.ts` matches the current `router/routes/*.go` manifest. If the `routes.ts` check fails, run `andurel generate routes --json`.

## Validation

Use the repository's allowed validation commands and project guidance. In this repo, do not run `go test`, `go build`, or `npm run`; use `go vet`, `go fix`, and `gofmt`.
