---
id: "01a1d39e-1269-78d8-900c-9e4d1ed96342"
title: "ISSUE-003: CUE formatter and OPL equivalence gate"
type: "Issue"
kind: "issue"
tags:
  - "issue"
  - "phase-2"
  - "schema"
  - "in-progress"
created: "2026-11-13T00:35:14-06:00"
timestamp: "2026-11-13T00:52:50-06:00"
description: "Add FormatCUE and a test that proves each OPL fixture and its CUE port give the same AST."
---

# ISSUE-003: CUE formatter and OPL equivalence gate

**Status:** In progress. **Phase:** 2.

## Description
Add FormatCUE and a test that proves each OPL fixture and its CUE port give the same AST.

## Acceptance Criteria
- [ ] `FormatCUE` renders namespaces as formatted CUE
- [ ] The equivalence test covers every OPL fixture in the repo

## Notes
Follow the [Development workflow](/development-workflow--01a1d39e1245.md).

## Related
- [ADR-0002](/decisions/adr-0002-replace-opl-with-cue-as-the-only-schema-language--01a1d39e124f.md)
