---
id: "01a1d39e-1273-7d3e-923c-390525be1f0f"
title: "ISSUE-005: Migrate the fixtures and remove OPL"
type: "Issue"
kind: "issue"
tags:
  - "issue"
  - "open"
  - "phase-4"
  - "schema"
created: "2026-11-13T00:35:14-06:00"
timestamp: "2026-11-13T00:35:14-06:00"
description: "Convert the OPL fixtures and examples to CUE, delete the OPL code and add the CUE validate command."
---

# ISSUE-005: Migrate the fixtures and remove OPL

**Status:** Open. **Phase:** 4.

## Description
Convert the OPL fixtures and examples to CUE, delete the OPL code and add the CUE validate command.

## Acceptance Criteria
- [ ] All engine and config tests use CUE schemas
- [ ] The `contrib/` examples use `.cue`
- [ ] The lexer, parser, fuzz corpus and type library are gone
- [ ] `keto namespace validate` checks CUE files and directories
- [ ] `GenerateCUEConfig` replaces `GenerateOPLConfig`

## Notes
Follow the [Development workflow](/development-workflow--01a1d39e1245.md).

## Related
- [ADR-0002](/decisions/adr-0002-replace-opl-with-cue-as-the-only-schema-language--01a1d39e124f.md)
