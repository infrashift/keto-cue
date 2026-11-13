---
id: "01a1d39e-124f-7d4a-a732-a53fee135847"
title: "ADR-0002: Replace OPL with CUE as the only schema language"
type: "ADR"
tags:
  - "adr"
  - "schema"
  - "cue"
  - "opl"
created: "2026-11-13T00:35:14-06:00"
timestamp: "2026-11-13T00:35:14-06:00"
description: "keto-cue compiles CUE schemas into the existing Keto AST and removes the OPL lexer and parser."
---

# ADR-0002: Replace OPL with CUE as the only schema language

**Status:** Accepted
**Date:** 2026-11-13
**Source:** User request, `PLAN.md`

## Context
Upstream Keto reads the Ory Permission Language (OPL), a TypeScript-like syntax. A custom lexer and parser in `schema/` turn OPL into `[]namespace.Namespace`. The check and expand engines walk the `ast.Relation` trees in each namespace. CUE gives types, constraints, packages and reuse that OPL does not have.

## Decision
- A CUE compiler in `schema/` writes the same `[]namespace.Namespace` and `ast` trees that the OPL parser wrote.
- CUE is the only schema language. The OPL lexer, parser, fuzz corpus and TypeScript type library go away.
- The check engine, the expand engine, the persistence layer and the protos do not change.
- The `/opl/syntax/check` REST and gRPC endpoints keep their names and contracts. They check CUE now.
- Before the OPL parser goes, an equivalence test proves that each OPL fixture and its CUE port give the same AST.

## Consequences
- Existing OPL schemas do not load. Users convert them with the CUE formatter (`FormatCUE`) or by hand.
- About 126 OPL fixtures in the engine tests change to CUE. Only test files change in `internal/check` and `internal/expand`.
- The endpoint name `/opl/syntax/check` no longer matches the language. The name stays for API stability.

## Alternatives considered
- **Dispatch on the `.cue` extension and keep OPL.** Rejected by the user. The fork has one schema language.
