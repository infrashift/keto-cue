---
id: "01a1d39e-125f-72e3-9b37-daa1fd7fd3c9"
title: "ISSUE-001: Okfw wiki, MCP server and documentation gates"
type: "Issue"
kind: "issue"
tags:
  - "issue"
  - "phase-0"
  - "process"
  - "closed"
created: "2026-11-13T00:35:14-06:00"
timestamp: "2026-11-13T00:35:49-06:00"
description: "Set up the okfw wiki, the okfw MCP server, the STE checker and the keto-cue Makefile gates."
---

# ISSUE-001: Okfw wiki, MCP server and documentation gates

**Status:** Closed. **Phase:** 0.

## Description
Set up the okfw wiki, the okfw MCP server, the STE checker and the keto-cue Makefile gates.

## Acceptance Criteria
- [x] `okfw init` and `okfw install --platform claude` ran
- [x] `.mcp.json` registers `okfw mcp`
- [x] Guides, ADR-0001 to ADR-0004 and one issue per phase exist
- [x] `scripts/ste-check.py` and the `lint`, `test-cue`, `wiki`, `docs-style` and `check` targets exist
- [x] `PLAN.md` lists each phase with a `Tracking:` line

## Resolution
Phase 0 is done. `make wiki` and `make docs-style` pass. The okfw MCP server needs a Claude Code restart or an approval in `/mcp`.

## Notes
Follow the [Development workflow](/development-workflow--01a1d39e1245.md).

## Related
- [ADR-0002](/decisions/adr-0002-replace-opl-with-cue-as-the-only-schema-language--01a1d39e124f.md)
