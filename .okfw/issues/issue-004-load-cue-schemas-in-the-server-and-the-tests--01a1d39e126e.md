---
id: "01a1d39e-126e-7a94-9d61-e28b476c0008"
title: "ISSUE-004: Load CUE schemas in the server and the tests"
type: "Issue"
kind: "issue"
tags:
  - "issue"
  - "phase-3"
  - "config"
  - "in-progress"
created: "2026-11-13T00:35:14-06:00"
timestamp: "2026-11-13T00:54:38-06:00"
description: "Route schema.Parse, the watcher, the syntax-check handler and the test helpers through the CUE compiler."
---

# ISSUE-004: Load CUE schemas in the server and the tests

**Status:** In progress. **Phase:** 3.

## Description
Route schema.Parse, the watcher, the syntax-check handler and the test helpers through the CUE compiler.

## Acceptance Criteria
- [ ] `schema.Parse` takes CUE source files
- [ ] The watcher keeps file bytes and compiles all `.cue` files of a directory as one package
- [ ] The syntax-check handler caps the request body
- [ ] `WithSchema` and `Scenario.Schema` replace the OPL names
- [ ] The config schema examples use `.cue`

## Notes
Follow the [Development workflow](/development-workflow--01a1d39e1245.md).

## Related
- [ADR-0002](/decisions/adr-0002-replace-opl-with-cue-as-the-only-schema-language--01a1d39e124f.md)
