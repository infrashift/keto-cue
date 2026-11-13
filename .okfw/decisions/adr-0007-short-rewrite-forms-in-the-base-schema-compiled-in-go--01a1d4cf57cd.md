---
id: "01a1d4cf-57cd-77e3-8231-6d23ca04b2a8"
title: "ADR-0007: Short rewrite forms in the base schema, compiled in Go"
type: "ADR"
tags:
  - "adr"
  - "cue"
  - "schema"
created: "2026-11-13T06:08:40-06:00"
timestamp: "2026-11-13T06:08:40-06:00"
description: "The base schema accepts short rewrite and relation forms, and the Go compiler converts them to the same AST as the long forms."
---

# ADR-0007: Short rewrite forms in the base schema, compiled in Go

**Status:** Accepted
**Date:** 2026-11-13
**Source:** [ISSUE-008](/issues/issue-008-short-forms-for-cue-rewrites-and-relations--01a1d4cf27e7.md)

## Context
The long rewrite form needs an `op` field and named operand fields. One traversal takes 82 characters, and the short form `"parents->view"` takes 15. Users asked for a helper library that makes rewrites as short as OPL.

## Decision
The base schema accepts short forms next to the long form. The Go compiler converts each short form to the same AST node as the long form.

| Short form | Long form |
|---|---|
| `"owners"` | `{op: "direct", relation: "owners"}` |
| `"parents->view"` | `{op: "tupleToSubjectSet", tupleSetRelation: "parents", computedSubjectSet: "view"}` |
| `{or: [a, b]}` | `{op: "or", children: [a, b]}` |
| `{and: [a, b]}` | `{op: "and", children: [a, b]}` |
| `{not: a}` | `{op: "invert", child: a}` |
| `owners: ["User"]` | `owners: types: ["User"]` |

- A rewrite node has exactly one of `op`, `or`, `and` and `not`.
- The arrow `->` follows the SpiceDB notation for a tuple-to-subject-set traversal. The compiler accepts spaces around the arrow.
- A relation name that contains `->`, or that starts or ends with a space, needs the long form.
- `FormatCUE` writes the short forms.

## Consequences
- The rewrites example drops from 38 to 25 lines that are not blank. Without the license header and the white space, it has 391 characters. The long form has 681 and the OPL original has 830.
- A schema can mix the short and the long forms.
- The base schema has two kind-disjoint disjunctions: string or struct for a rewrite, and list or struct for a relation. When a struct does not match, the compiler reports the error inside the struct, for example `field not allowed`, and not the disjunction summary.

## Alternatives considered
- **An importable CUE package with constructor definitions**, such as `rw.#TTU & {#rel: "parents", #then: "view"}`. Rejected: CUE has no functions. A constructor call is as long as the long form. Hidden fields are private to their package, so they cannot carry parameters across an import. An import also needs a module resolver, against [ADR-0003](/decisions/adr-0003-build-cue-instances-in-memory-without-cue-load-or-a-module-registry--01a1d39e1254.md).
- **Comprehensions that expand a compact field into `permissions`.** Rejected: the output field must be separate from the input field, so the schema gets a second, generated structure. Errors then point at generated values and not at what the user wrote.
- **A string expression language, such as `"viewers | parents->view"`.** Rejected: it needs a second parser and gives weaker error positions. CUE lists and structs already express the tree.
