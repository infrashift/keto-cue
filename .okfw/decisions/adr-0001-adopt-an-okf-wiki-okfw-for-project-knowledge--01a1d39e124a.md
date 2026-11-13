---
id: "01a1d39e-124a-7bb5-998e-30689010c652"
title: "ADR-0001: Adopt an OKF wiki (okfw) for project knowledge"
type: "ADR"
tags:
  - "adr"
  - "process"
  - "okf"
created: "2026-11-13T00:35:14-06:00"
timestamp: "2026-11-13T00:35:14-06:00"
description: "An okfw-managed OKF wiki in .okfw/ holds keto-cue project knowledge: architecture, specs, ADRs, issues and daily notes."
---

# ADR-0001: Adopt an OKF wiki (okfw) for project knowledge

**Status:** Accepted
**Date:** 2026-11-13
**Source:** User request

## Context
keto-cue changes the schema language of Ory Keto. Chat history and commit messages hide the reasons for each design choice. Sibling infrashift repos (`zmake`, `cuix`, `whac-a-leak`) use okfw for this knowledge.

## Decision
- An OKF wiki lives at `.okfw/` in the repo root. The `okfw` CLI manages it. The repository includes it.
- `okfw install --platform claude` adds the Claude Code skill and the managed `CLAUDE.md` block. The `.mcp.json` file registers `okfw mcp`.
- The wiki uses the same bundles as zmake: `architecture/`, `specs/`, `decisions/`, `issues/` and `daily/`.
- `PLAN.md` stays the phase checklist. Each phase has a mirror issue in `issues/`.
- Conventions: [Wiki conventions](/wiki-conventions--01a1d39e123f.md). Loop: [Development workflow](/development-workflow--01a1d39e1245.md).

## Consequences
- `okfw search` finds decisions and open work. Links connect them into a graph.
- `PLAN.md` and the issues can drift apart. The workflow guide requires an update to both.
- `make wiki` and `make docs-style` become commit gates.

## Alternatives considered
- **Plain `docs/` markdown.** Rejected: no frontmatter, no typed search and no generated indexes.
