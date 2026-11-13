# CLAUDE.md — keto-cue

## Purpose (WHY)
keto-cue is a fork of Ory Keto. It reads permission schemas in CUE, not in the Ory Permission Language (OPL). The CUE compiler writes the same AST that the OPL parser wrote. The check engine, the expand engine, the persistence layer and the protos do not change.

## Tech Stack (WHAT)
- **Language:** Go (see `go.mod`).
- **Schema language:** CUE through `cuelang.org/go`, pinned to one version in `go.mod`.
- **Wiki:** okfw (`.okfw/`). See the okfw block below.

## Layout
- `schema/` holds the schema compiler, the semantic checks, the parse errors and the syntax-check handler.
- `internal/namespace/ast/` holds the AST that the check and expand engines walk.
- `internal/driver/config/opl_config_namespace_watcher.go` loads and watches the schema files.
- `cmd/namespace/` holds the `keto namespace` commands.
- `scripts/ste-check.py` checks the keto-cue documents against ASD-STE100.

## Verification & Workflow Commands (HOW)
Run these targets after every change:
- `make lint` → `gofmt -l` and `go vet` on the keto-cue packages
- `make test-cue` → `go test` on the schema, config, CLI and engine packages
- `make wiki` → `okfw validate --all && okfw validate --okf`
- `make docs-style` → `scripts/ste-check.py`
- `make check` → all of the gates above

## Mandatory Execution Rules
1. **One phase at a time.** Read `PLAN.md`. Implement only the first unchecked phase. Follow the *Development workflow* guide in `.okfw/`.
2. **Repeat until green.** Implement, lint, verify. Fix each failure and run all gates again. Do not commit while a gate fails.
3. **Test integrity.** Do not weaken or delete tests to make a gate pass. Port an OPL test to CUE. Do not drop it.
4. **Small blast radius.** Do not change production code in `internal/check`, `internal/expand`, `internal/persistence` or `proto/`. Test files in those packages can change.
5. **Wiki stays in sync.** Each `PLAN.md` phase has a mirror issue in `.okfw/issues/`. Record each decision as an ADR when you make it. Close the issue when the phase is done.
6. **Documentation follows ASD-STE100 (Simplified Technical English).** Use active voice, simple present tense and short sentences. Run `make docs-style`.

<!-- BEGIN okfw -->
## Project knowledge base (okfw)

keto-cue keeps a local-first OKF wiki in `.okfw/`: markdown concepts with
YAML frontmatter, linked to each other. The files are the whole database —
search and the link graph are computed from them on demand, so there is no
index to build and nothing to sync. Each directory's `index.md` is generated
and gives a cheap overview before you open anything (OKF section 8).

- **Search before you derive.** `okfw search "<query>"` / `okfw_search`. The
  wiki holds what the code does not say — why things are built the way they
  are, what was considered and rejected, dependency and security findings,
  operational context, open work. Check it before explaining a design,
  assessing risk, or running an analysis someone may already have recorded and
  adjudicated.
- **Browse top-down.** Each directory's `index.md` lists its concepts with
  their descriptions (`okfw_read {"ref":"<dir>/"}` over MCP). Open only what fits.
- **Search before you write**, too, so a new concept links to what already
  exists instead of duplicating it.
- **Link actively** — `[Title](/dir/file.md)` from the wiki root, or
  `./file.md`. Links are what makes the wiki more than a folder.
- **One idea per concept**, with a one-sentence `description`.
- **Use `okfw move` / `okfw delete`, never `mv` / `rm`.** `move` repairs every
  link to the concept and its own relative links, and both keep the `index.md`
  listings current; the shell commands break the links and leave the listings
  stale.
- **`status` is the OKF lifecycle field** (`draft`/`stable`/`deprecated`), not
  task state. Mark superseded knowledge rather than deleting it.
- **Custom metadata goes under `x:`** — the schema is closed and rejects any
  other unrecognized key.
- Don't edit `index.md` or `log.md`; okfw generates them.

Results are capped and say so: `truncated`, or `… N more`, means you did not
see everything — raise the limit or page with `--offset`/`"offset"`.

`okfw mcp` serves all of this as MCP stdio tools named `okfw_*`.
<!-- END okfw -->
