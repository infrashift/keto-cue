---
id: "01a1d39e-1273-7d3e-923c-390525be1f0f"
title: "ISSUE-005: Migrate the fixtures and remove OPL"
type: "Issue"
kind: "issue"
tags:
  - "issue"
  - "phase-4"
  - "schema"
  - "in-progress"
created: "2026-11-13T00:35:14-06:00"
timestamp: "2026-11-13T01:21:52-06:00"
description: "Convert the OPL fixtures and examples to CUE, delete the OPL code and add the CUE validate command."
---

# ISSUE-005: Migrate the fixtures and remove OPL

**Status:** In progress. **Phase:** 4.

## Description
Convert the OPL fixtures and examples to CUE, delete the OPL code and add the CUE validate command.

## Acceptance Criteria
- [ ] The lexer, parser, fuzz corpus and type library are gone
- [ ] `keto namespace validate` checks CUE files and directories
- [ ] `GenerateCUEConfig` replaces `GenerateOPLConfig`

## Notes
The fixture conversion moved to [ISSUE-004](/issues/issue-004-load-cue-schemas-in-the-server-and-the-tests--01a1d39e126e.md). The engine tests need CUE as soon as `schema.Parse` reads CUE.

Follow the [Development workflow](/development-workflow--01a1d39e1245.md).

## Related
- [ADR-0002](/decisions/adr-0002-replace-opl-with-cue-as-the-only-schema-language--01a1d39e124f.md)
