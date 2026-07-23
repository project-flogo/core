# AGENTS.md

Guidance for AI coding agents working in this repository.

## What this repo is

`github.com/project-flogo/core` — the core Go library for [Project Flogo](https://flogo.io),
a framework for building serverless functions and edge microservices.

This module ships the runtime, model, and extension points used by Flogo apps:
`action`, `activity`, `trigger`, `engine`, `app`, `data`, `api`, plus shared
`support` utilities. It is imported by the CLI (`project-flogo/cli`), contribution
repos (`project-flogo/contrib`, `project-flogo/flow`), and applications.

- Module path: `github.com/project-flogo/core`
- Go version: `1.18` (see `go.mod`)
- Current version: see `VERSION`

## Repo layout

```
action/       Action interface + registry (units executed by a trigger handler)
activity/     Activity interface, context, metadata, registry, error types
trigger/      Trigger interface, handler, config, descriptor, registry
engine/       Engine impl, runner, config, event, secret, service, channels
app/          App-level config, controller, properties, resources, resolvers
api/          High-level Go API for building/running an app programmatically
data/         Type system: coerce, expression, mapper, path, resolve, schema
support/      Shared utilities: log, trace, connection, managed, ssl, service, test
internal/     Internal-only helpers (currently JSON schema)
docs/         Model, data types, mapping, engine, properties, api docs
examples/     Minimal action, activity, trigger, and engine examples
schema.json   JSON schema for the app descriptor (flogo.json)
```

Read `docs/model.md` for the app JSON model and `docs/api.md` for the Go API
before making non-trivial changes to `app/`, `engine/`, or `api/`.

## Build, test, lint

```bash
go build ./...
go test ./...
go vet ./...
```

CI (`.travis.yml`) runs `go test ./...`. There is no separate lint step configured
in-repo. Keep new files `gofmt`-clean.

Run a single package's tests:

```bash
go test ./activity/...
```

## Contribution model (important context)

Flogo is extended via **contributions**: activities, triggers, actions, and
connections implemented as separate Go packages that register themselves via
`init()` into the registries here (`activity/registry.go`, `trigger/registry.go`,
`action/registry.go`). When adding or changing interfaces in this repo,
consider the downstream blast radius on `project-flogo/contrib` and
`project-flogo/flow`.

The `examples/` directory has the minimal shape of each contribution type — use
those as reference when reviewing or authoring registration code.

## Conventions

- Standard Go project layout; no generated code checked in.
- Logging: use `support/log` (wraps `go.uber.org/zap`), not `fmt.Println` or
  the standard `log` package, in library code.
- Errors: prefer wrapping with context; `activity/error.go` defines
  activity-specific error helpers.
- Data coercion: use `data/coerce` — do not hand-roll type conversions.
- Expressions and mappings go through `data/expression` and `data/mapper`.
- Tests live next to the code they cover (`*_test.go`), using
  `github.com/stretchr/testify`.

## Branching

- Default branch: `master`.
- Long-lived: `flogo3x` (3.x-line work; recent commits merge fixes back and
  forth). `flogo-3x` also appears in commit messages.
- Feature branches use `FLOGO-<jira-id>-<slug>` naming, e.g.
  `FLOGO-17230-schema-validation-error`.

Current working branch: `agent-resource-support`.

## When making changes

- Preserve backward compatibility of exported APIs in `action/`, `activity/`,
  `trigger/`, `app/`, `engine/`, and `api/` unless the change is coordinated
  with the contrib / flow / cli repos.
- The app JSON model is versioned via `appModel` in the descriptor; changes to
  `schema.json` or `app/config.go` may need a model-version bump.
- Update `docs/*.md` when changing the model, API surface, or data-type
  behavior — those files are the reference for contribution authors.
- Do not add third-party dependencies casually; `go.mod` is intentionally lean
  (zap, testify, gojsonschema, dateparse).

## Useful entry points when exploring

- `api/api.go` — programmatic app construction (`NewApp`, `NewTrigger`, …)
- `engine/engineimpl.go` — engine lifecycle
- `app/app.go` + `app/config.go` — app assembly from JSON descriptor
- `activity/activity.go`, `trigger/trigger.go`, `action/action.go` — the three
  extension interfaces
- `data/resolve/` and `data/expression/` — how `$property[...]`, `$flow.X`,
  `res://...`, `conn://...` etc. are resolved
