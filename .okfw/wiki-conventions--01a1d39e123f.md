---
id: "01a1d39e-123f-7e53-b7e9-28a44690221f"
title: "Wiki conventions"
type: "Guide"
tags:
  - "conventions"
  - "meta"
created: "2026-11-13T00:35:14-06:00"
timestamp: "2026-11-13T00:35:14-06:00"
description: "Layout, writing style and authoring rules for architecture notes, specs, ADRs, issues and daily notes in this wiki."
---

# Wiki conventions

This guide defines how the keto-cue OKF wiki (`.okfw/`) stores project knowledge.

- Create concepts with `okfw create` or the okfw MCP `create` tool. Do not write frontmatter by hand.
- Search before you write. Link to related concepts.
- Keep one idea per concept. Give each concept a one-sentence `description`.
- Do not start a `--body` with `# Title`. okfw writes the H1 from the title.

## Layout
| Path | Type | Contents |
|---|---|---|
| `architecture/` | `Architecture` | System overview, package boundaries, data flow. |
| `specs/` | `Spec` | Contracts: the CUE permission schema, the AST mapping, the error format. The spec states intent. The code implements it. |
| `decisions/` | `ADR` | Architecture decision records. |
| `issues/` | `Issue` (kind `issue`) | Tracked work: one issue per `PLAN.md` phase, plus bugs and design debt. |
| `daily/` | daily note | `okfw daily`: what changed, what you learned, what is next. |

## PLAN.md and issues
`PLAN.md` is the phase checklist. Each phase has a mirror issue. Each phase heading in `PLAN.md` has a `Tracking:` line that links to its issue. When a phase is done, tick the `PLAN.md` checkboxes and the acceptance criteria of the issue.

## ADRs (`decisions/`)
- Title: `ADR-NNNN: <decision>`. Number ADRs in sequence. Do not reuse a number.
- Command: `okfw create "ADR-NNNN: <title>" --path decisions --type ADR --tags adr,<topic> -d "<one sentence>" --body "<markdown>"`
- Body: **Status** (Proposed / Accepted / Superseded by ADR-XXXX), **Date**, **Source**, then **Context**, **Decision**, **Consequences**, **Alternatives considered**.
- To supersede an ADR, write the new ADR with a link to the old ADR. Update the Status line of the old ADR. Run `okfw set-lifecycle <id> deprecated`.
- Write an ADR before you add a Go module dependency.

## Issues (`issues/`)
- Title: `ISSUE-NNN: <summary>`. Number issues in sequence.
- Command: `okfw create "ISSUE-NNN: <title>" --kind issue --path issues --type Issue --tags issue,open,phase-<n>,<area> -d "<summary>" --body "<markdown>"`
- Body: a **Status/Phase** line, **Description**, **Acceptance Criteria** (checkboxes), **Notes**, **Related**.
- okfw has no status field for issues. Tags carry the state: `open` → `in-progress` → `closed`. Use `okfw remove-tags` and `okfw add-tags`. Keep the Status line the same as the tags. Add a **Resolution** section when you close an issue.
- Find open work with `okfw search --path issues --tag open`.

## Writing style
All keto-cue documentation follows ASD-STE100 Simplified Technical English. This includes the wiki, `PLAN.md`, `CLAUDE.md` and `docs/cue_permission_schema.md`. The upstream Ory documents keep their own style.

- Use active voice and simple present tense: "The compiler rejects the field."
- Start each procedural step with an imperative verb: Run, Read, Inspect, Verify, Add, Set.
- Keep procedural steps to 20 words or fewer. Keep descriptive sentences to 25 words or fewer.
- Use a verb, not a noun form: "The compiler checks the reference", not "A check of the reference occurs".
- Use one term for one concept. This wiki uses these terms:
  - *schema*: the CUE source that declares the namespaces. A *schema package* is a directory of `.cue` files.
  - *base schema*: the embedded `schema/cue_base.cue` that each schema unifies with.
  - *namespace*, *relation* (a field in `relations`) and *permission* (a field in `permissions`).
  - *rewrite*: a subject-set rewrite tree. Its nodes are `or`, `and`, `direct`, `tupleToSubjectSet` and `invert`.
  - *AST*: the Go tree in `internal/namespace/ast` that the check and expand engines walk.
  - *relation tuple*, *subject* and *subject set*, as Keto defines them.
  - *OPL*: the Ory Permission Language that keto-cue replaces.
  - *phase*: one unit of work in `PLAN.md`.
- Do not use the words `should`, `will`, `would`, `please`, `simply`, `just`, `very`, `really` or `easily`. Do not use marketing language.
- Run `make docs-style`. The check (`scripts/ste-check.py`) reports banned words and long sentences as errors. It reports possible passive voice as warnings.

## Housekeeping
- Use `okfw move` and `okfw delete`. Do not use `mv` or `rm`. The okfw commands repair links and update the `index.md` listings.
- `status` is the OKF lifecycle field (`draft` / `stable` / `deprecated`). It is not the task state. Mark superseded knowledge `deprecated`. Do not delete it.
- Do not edit `index.md` or `log.md`. okfw generates them.
- Run `make wiki` before you commit wiki changes.
