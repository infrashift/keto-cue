---
id: "01a1d39e-1269-78d8-900c-9e4d1ed96342"
title: "ISSUE-003: CUE formatter and OPL equivalence gate"
type: "Issue"
kind: "issue"
tags:
  - "issue"
  - "phase-2"
  - "schema"
  - "closed"
created: "2026-11-13T00:35:14-06:00"
timestamp: "2026-11-13T00:54:38-06:00"
description: "Add FormatCUE and a test that proves each OPL fixture and its CUE port give the same AST."
---

# ISSUE-003: CUE formatter and OPL equivalence gate

**Status:** Closed. **Phase:** 2.

## Description
Add FormatCUE and a test that proves each OPL fixture and its CUE port give the same AST.

## Acceptance Criteria
- [x] `FormatCUE` renders namespaces as formatted CUE
- [x] The equivalence test covers every OPL fixture in the repo

## Resolution
Phase 2 is done. `FormatCUE` renders namespaces as formatted CUE in package `keto`. `TestFormatCUEEquivalence` finds 62 OPL fixtures in the repo: the parser cases, the `.ts` files, the CLI snapshots and the Go test literals. Each fixture gives the same AST through OPL and through `FormatCUE` → `ParseCUE`.

## Notes
Follow the [Development workflow](/development-workflow--01a1d39e1245.md).

## Related
- [ADR-0002](/decisions/adr-0002-replace-opl-with-cue-as-the-only-schema-language--01a1d39e124f.md)
