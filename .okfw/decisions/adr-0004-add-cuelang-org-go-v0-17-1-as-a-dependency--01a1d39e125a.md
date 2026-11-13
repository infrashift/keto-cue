---
id: "01a1d39e-125a-71ec-90f4-2c5840b312f7"
title: "ADR-0004: Add cuelang.org/go v0.17.1 as a dependency"
type: "ADR"
tags:
  - "adr"
  - "dependency"
  - "cue"
created: "2026-11-13T00:35:14-06:00"
timestamp: "2026-11-13T00:35:14-06:00"
description: "keto-cue pins cuelang.org/go v0.17.1, the latest stable release, to compile CUE schemas."
---

# ADR-0004: Add cuelang.org/go v0.17.1 as a dependency

**Status:** Accepted
**Date:** 2026-11-13
**Source:** [ADR-0002](/decisions/adr-0002-replace-opl-with-cue-as-the-only-schema-language--01a1d39e124f.md)

## Context
[ADR-0002](/decisions/adr-0002-replace-opl-with-cue-as-the-only-schema-language--01a1d39e124f.md) replaces OPL with CUE. The compiler needs the CUE evaluator. The module cache holds v0.15.x, v0.16.0, v0.17.1 and a v0.18.0 release candidate.

## Decision
- Pin `cuelang.org/go v0.17.1` in `go.mod`. Do not use `@latest` or a release candidate.
- Use only the public packages: `cue`, `cue/cuecontext`, `cue/build`, `cue/errors`, `cue/token`, `cue/ast` and `cue/format`.

## Consequences
- The keto binary grows by the size of the CUE evaluator.
- A CUE upgrade is a deliberate change. Run `make test-cue` after each upgrade.
- `go get` also raised three indirect modules: `go-toml/v2` 2.2.3 → 2.3.1, `go-internal` 1.14.1 → 1.15.0 and `wazero` 1.11.0 → 1.12.0. CUE requires these versions.

## Alternatives considered
- **Write a CUE subset parser.** Rejected: it loses unification, constraints and the standard library.
