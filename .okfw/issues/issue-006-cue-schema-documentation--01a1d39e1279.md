---
id: "01a1d39e-1279-7455-ba9c-4631d7a1b1a4"
title: "ISSUE-006: CUE schema documentation"
type: "Issue"
kind: "issue"
tags:
  - "issue"
  - "phase-5"
  - "docs"
  - "closed"
created: "2026-11-13T00:35:14-06:00"
timestamp: "2026-11-13T01:28:43-06:00"
description: "Document the CUE permission schema and mark the OPL documents as superseded."
---

# ISSUE-006: CUE schema documentation

**Status:** Closed. **Phase:** 5.

## Description
Document the CUE permission schema and mark the OPL documents as superseded.

## Acceptance Criteria
- [x] `docs/cue_permission_schema.md` exists and passes `make docs-style`
- [x] `docs/ory_permission_language_spec.md` and `README.md` point to the CUE document

## Resolution
Phase 5 is done. `docs/cue_permission_schema.md` documents the schema, the OPL to CUE table, validation and the syntax-check limits. `README.md` has a fork notice with a CUE example. `docs/ory_permission_language_spec.md` has a superseded notice.

## Notes
Follow the [Development workflow](/development-workflow--01a1d39e1245.md).

## Related
- [ADR-0002](/decisions/adr-0002-replace-opl-with-cue-as-the-only-schema-language--01a1d39e124f.md)
