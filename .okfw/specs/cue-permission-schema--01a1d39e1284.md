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

**Status:** Draft. [ISSUE-002](/issues/issue-002-cue-compiler-next-to-the-opl-parser--01a1d39e1264.md) builds it in `schema/`.

## Shape
A schema declares `namespaces`. Each key is a namespace name. Each namespace has optional `relations` and optional `permissions`.

```cue
namespaces: {
	User: {}
	Group: relations: members: types: ["User"]
	Document: {
		relations: {
			owners: types: ["User", "Group#members"]
			parents: types: ["Folder"]
		}
		permissions: view: {
			op: "or"
			children: [
				{op: "direct", relation: "owners"},
				{op: "tupleToSubjectSet", tupleSetRelation: "parents", computedSubjectSet: "view"},
			]
		}
	}
}
```

## Types
A relation type is `"Namespace"`, `"Namespace#relation"` or `{namespace: "Namespace", relation: "relation"}`.

## Rewrite operations
| `op` | Fields | AST node |
|---|---|---|
| `or` | `children` (one or more) | `*ast.SubjectSetRewrite{Operation: ast.OperatorOr}` |
| `and` | `children` (one or more) | `*ast.SubjectSetRewrite{Operation: ast.OperatorAnd}` |
| `direct` | `relation` | `*ast.ComputedSubjectSet` |
| `tupleToSubjectSet` | `tupleSetRelation`, `computedSubjectSet` | `*ast.TupleToSubjectSet` |
| `invert` | `child` | `*ast.InvertResult` |

## Rules
- The compiler wraps a top-level leaf in a rewrite and flattens nested operations of the same kind.
- Relations come first in the AST, then permissions. Each group keeps the declaration order.
- A name in both `relations` and `permissions` is an error.
- Each referenced namespace and relation exists. The compiler reports each missing reference.
- Rewrites nest to a depth of 10 or less.
