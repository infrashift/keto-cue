---
id: "01a1d39e-1284-7285-bd69-bd08266d2727"
title: "CUE permission schema"
type: "Spec"
tags:
  - "spec"
  - "schema"
  - "cue"
created: "2026-11-13T00:35:14-06:00"
timestamp: "2026-11-13T00:35:14-06:00"
description: "The contract for keto-cue schemas: the base schema, the rewrite operations and their AST mapping."
---

# CUE permission schema

**Status:** Implemented in `schema/cue_base.cue`, `schema/cue_compiler.go` and `schema/cue_checks.go` ([ISSUE-002](/issues/issue-002-cue-compiler-next-to-the-opl-parser--01a1d39e1264.md)). The short forms came in [ISSUE-008](/issues/issue-008-short-forms-for-cue-rewrites-and-relations--01a1d4cf27e7.md).

## Shape
A schema declares `namespaces`. Each key is a namespace name. Each namespace has optional `relations` and optional `permissions`.

```cue
namespaces: {
	User: {}
	Group: relations: members: ["User"]
	Document: {
		relations: {
			owners: ["User", "Group#members"]
			parents: ["Folder"]
		}
		permissions: view: or: ["owners", "parents->view"]
	}
}
```

## Types
A relation type is `"Namespace"`, `"Namespace#relation"` or `{namespace: "Namespace", relation: "relation"}`. A relation is its list of types (`owners: ["User"]`) or a struct `{types?: [...], rewrite?: ...}`.

## Rewrite operations
A rewrite is a short form or a node with one of the fields `op`, `or`, `and` and `not` ([ADR-0007](/decisions/adr-0007-short-rewrite-forms-in-the-base-schema-compiled-in-go--01a1d4cf57cd.md)).

| Short form | `op` | Fields | AST node |
|---|---|---|---|
| `{or: [...]}` | `or` | `children` (one or more) | `*ast.SubjectSetRewrite{Operation: ast.OperatorOr}` |
| `{and: [...]}` | `and` | `children` (one or more) | `*ast.SubjectSetRewrite{Operation: ast.OperatorAnd}` |
| `"relation"` | `direct` | `relation` | `*ast.ComputedSubjectSet` |
| `"tupleSet->computed"` | `tupleToSubjectSet` | `tupleSetRelation`, `computedSubjectSet` | `*ast.TupleToSubjectSet` |
| `{not: ...}` | `invert` | `child` | `*ast.InvertResult` |

`FormatCUE` writes the short forms. A relation name that contains `->`, or that starts or ends with a space, keeps the long form.

## Helpers
The top level of a schema stays open. A schema can declare values next to `namespaces` and use comprehensions:

```cue
#actions: ["read", "write"]

namespaces: Doc: permissions: {
	for a in #actions {(a): or: ["owners", "parents->\(a)"]}
}
```

## Files and imports
- All files of one schema belong to one CUE package and declare the same `package` clause.
- A schema can import the CUE standard library. A schema cannot import other packages ([ADR-0003](/decisions/adr-0003-build-cue-instances-in-memory-without-cue-load-or-a-module-registry--01a1d39e1254.md)).

## Errors
Each error names the file, the line and the column ([ADR-0005](/decisions/adr-0005-parse-errors-carry-the-file-name-and-1-based-positions--01a1d3adf8cf.md)). The compiler reports CUE errors first. It runs the reference checks only when the schema has no CUE errors.

## User documentation
[docs/cue_permission_schema.md](https://github.com/infrashift/keto-cue/blob/master/docs/cue_permission_schema.md) is the user guide. It has the OPL to CUE table and the limits of the syntax check.

## Rules
- The compiler wraps a top-level leaf in a rewrite and flattens nested operations of the same kind.
- Relations come first in the AST, then permissions. Each group keeps the declaration order.
- A name in both `relations` and `permissions` is an error.
- Each referenced namespace and relation exists. The compiler reports each missing reference.
- Rewrites nest to a depth of 10 or less.
- A value that matches no alternative of a base-schema disjunction gives one error. When an error inside the value explains the conflict, the compiler reports that error. Otherwise the error names the allowed values.
- The compiler reports an error with the same message at the same position once.
