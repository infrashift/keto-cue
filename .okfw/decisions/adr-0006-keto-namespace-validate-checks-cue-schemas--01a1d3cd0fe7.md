---
id: "01a1d3cd-0fe7-7d7d-b38e-fb3f096ffc74"
title: "ADR-0006: keto namespace validate checks CUE schemas"
type: "ADR"
tags:
  - "adr"
  - "cli"
  - "cue"
created: "2026-11-13T01:26:33-06:00"
timestamp: "2026-11-13T01:26:33-06:00"
description: "keto namespace validate checks CUE files and directories, and the legacy validator keeps only the validate-legacy name."
---

# ADR-0006: keto namespace validate checks CUE schemas

**Status:** Accepted
**Date:** 2026-11-13
**Source:** [ISSUE-005](/issues/issue-005-migrate-the-fixtures-and-remove-opl--01a1d39e1273.md)

## Context
Upstream Keto has no command that validates an OPL file. `keto namespace validate` was an alias of `validate-legacy`, which validates the deprecated YAML, JSON and TOML namespace files. keto-cue needs a command that validates CUE schemas with source positions.

## Decision
- `keto namespace validate <file.cue | dir> [...]` validates CUE schemas. A directory is one CUE package, as the server loads it.
- The command prints each error with the file, the line, the column and a source snippet. It exits with an error when a schema is not valid.
- `keto namespace validate-legacy` keeps the legacy validation. It loses the `validate` alias.
- `GenerateCUEConfig` replaces `GenerateOPLConfig`. It writes a CUE schema that declares the legacy namespace names.

## Consequences
- Scripts that call `keto namespace validate -c <config>` for legacy namespaces must call `validate-legacy`.
- Operators can validate a schema directory before they deploy it.

## Alternatives considered
- **A new command name such as `validate-cue`.** Rejected: CUE is the only schema language, so `validate` is the natural name.
