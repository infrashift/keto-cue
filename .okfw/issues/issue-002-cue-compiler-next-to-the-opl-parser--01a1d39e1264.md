---
id: "01a1d39e-1264-763f-bae5-03cf3ec5dd7d"
title: "ISSUE-002: CUE compiler next to the OPL parser"
type: "Issue"
kind: "issue"
tags:
  - "issue"
  - "phase-1"
  - "schema"
  - "closed"
created: "2026-11-13T00:35:14-06:00"
timestamp: "2026-11-13T00:52:50-06:00"
description: "Add the base schema, the CUE-to-AST compiler, the semantic checks and the new ParseError shape."
---

# ISSUE-002: CUE compiler next to the OPL parser

**Status:** Closed. **Phase:** 1.

## Description
Add the base schema, the CUE-to-AST compiler, the semantic checks and the new ParseError shape.

## Acceptance Criteria
- [x] `cuelang.org/go v0.17.1` is in `go.mod`
- [x] `schema/cue_base.cue` is embedded
- [x] `ParseCUE` maps `or`, `and`, `direct`, `tupleToSubjectSet` and `invert` to the AST
- [x] The OPL type checks have CUE ports with source positions
- [x] `ParseError` holds a file name, the source and start and end positions
- [x] Snapshot, negative and multi-file tests pass

## Resolution
Phase 1 is done. `ParseCUE` compiles one or more CUE files into the Keto AST. The four OPL parser cases have CUE ports that give the same AST. The tests cover 17 error cases, a multi-file package and a standard-library import. `FuzzParseCUE` ran for 4 minutes without a failure. Follow-up: [ISSUE-007](/issues/issue-007-bound-cue-evaluation-on-the-syntax-check-listener--01a1d3adf8f1.md).

## Notes
Follow the [Development workflow](/development-workflow--01a1d39e1245.md).

## Related
- [ADR-0002](/decisions/adr-0002-replace-opl-with-cue-as-the-only-schema-language--01a1d39e124f.md)
