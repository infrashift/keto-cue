---
id: "01a1d39e-1264-763f-bae5-03cf3ec5dd7d"
title: "ISSUE-002: CUE compiler next to the OPL parser"
type: "Issue"
kind: "issue"
tags:
  - "issue"
  - "open"
  - "phase-1"
  - "schema"
created: "2026-11-13T00:35:14-06:00"
timestamp: "2026-11-13T00:35:14-06:00"
description: "Add the base schema, the CUE-to-AST compiler, the semantic checks and the new ParseError shape."
---

# ISSUE-002: CUE compiler next to the OPL parser

**Status:** Open. **Phase:** 1.

## Description
Add the base schema, the CUE-to-AST compiler, the semantic checks and the new ParseError shape.

## Acceptance Criteria
- [ ] `cuelang.org/go v0.17.1` is in `go.mod`
- [ ] `schema/cue_base.cue` is embedded
- [ ] `ParseCUE` maps `or`, `and`, `direct`, `tupleToSubjectSet` and `invert` to the AST
- [ ] The OPL type checks have CUE ports with source positions
- [ ] `ParseError` holds a file name, the source and start and end positions
- [ ] Snapshot, negative and multi-file tests pass

## Notes
Follow the [Development workflow](/development-workflow--01a1d39e1245.md).

## Related
- [ADR-0002](/decisions/adr-0002-replace-opl-with-cue-as-the-only-schema-language--01a1d39e124f.md)
