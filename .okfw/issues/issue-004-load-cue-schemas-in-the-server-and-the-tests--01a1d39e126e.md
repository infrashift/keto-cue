---
id: "01a1d39e-126e-7a94-9d61-e28b476c0008"
title: "ISSUE-004: Load CUE schemas in the server and the tests"
type: "Issue"
kind: "issue"
tags:
  - "issue"
  - "phase-3"
  - "config"
  - "closed"
created: "2026-11-13T00:35:14-06:00"
timestamp: "2026-11-13T01:21:52-06:00"
description: "Route schema.Parse, the watcher, the syntax-check handler and the test helpers through the CUE compiler."
---

# ISSUE-004: Load CUE schemas in the server and the tests

**Status:** Closed. **Phase:** 3.

## Description
Route schema.Parse, the watcher, the syntax-check handler and the test helpers through the CUE compiler.

## Acceptance Criteria
- [x] `schema.Parse` takes CUE source files
- [x] The watcher keeps file bytes and compiles all `.cue` files of a directory as one package
- [x] The syntax-check handler caps the request body
- [x] `WithSchema` and `Scenario.Schema` replace the OPL names
- [x] The config schema examples use `.cue`
- [x] The engine and config test fixtures, the benchmark fixture and the `contrib/` examples use CUE
- [x] A test compiles every `.cue` schema in the repository

## Resolution
Phase 3 is done. `schema.Parse` compiles CUE source files. The watcher keeps bytes and compiles a directory as one package; a test covers the multi-file change that the old reader-based watcher broke. The syntax-check handler has the limits of [ISSUE-007](/issues/issue-007-bound-cue-evaluation-on-the-syntax-check-listener--01a1d3adf8f1.md). The 51 test literals, the benchmark fixture and the four `contrib/` schemas are CUE. The check, expand, trace and step suites pass on the CUE schemas.

## Notes
Follow the [Development workflow](/development-workflow--01a1d39e1245.md).

## Related
- [ADR-0002](/decisions/adr-0002-replace-opl-with-cue-as-the-only-schema-language--01a1d39e124f.md)
