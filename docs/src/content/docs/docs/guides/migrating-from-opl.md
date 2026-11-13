---
title: Migrating from OPL
description: Translate an Ory Permission Language (OPL) schema into a keto-cue CUE schema.
---

keto-cue does not read OPL. Translate each OPL schema to CUE. The CUE compiler
writes the same namespace tree as the OPL parser, so the relation tuples and
the check results do not change.

## Translation table

| OPL | CUE |
|---|---|
| `class Doc implements Namespace {}` | `namespaces: Doc: {}` |
| `related: { owners: User[] }` | `relations: owners: ["User"]` |
| `owners: (User \| SubjectSet<Group, "members">)[]` | `owners: ["User", "Group#members"]` |
| `this.related.owners.includes(ctx.subject)` | `"owners"` |
| `this.permits.edit(ctx)` | `"edit"` |
| `this.related.parents.traverse((p) => p.permits.view(ctx))` | `"parents->view"` |
| `a \|\| b` | `{or: [a, b]}` |
| `a && b` | `{and: [a, b]}` |
| `!a` | `{not: a}` |

## Example

This OPL schema:

```typescript
import { Namespace, Context } from "@ory/keto-namespace-types"

class User implements Namespace {}

class Folder implements Namespace {
  related: {
    parents: Folder[]
    viewers: User[]
  }

  permits = {
    view: (ctx: Context): boolean =>
      this.related.viewers.includes(ctx.subject) ||
      this.related.parents.traverse((p) => p.permits.view(ctx)),
  }
}
```

becomes this CUE schema:

```cue
namespaces: {
	User: {}
	Folder: {
		relations: {
			parents: ["Folder"]
			viewers: ["User"]
		}
		permissions: view: or: ["viewers", "parents->view"]
	}
}
```

## Steps

1. Translate each class to a field of `namespaces`.
2. Move the `related` fields to `relations` and the `permits` functions to
   `permissions`.
3. Run `keto namespace validate` on the new schema.
4. Change `namespaces.location` in the Keto configuration to the `.cue` file.
5. Compare the results of `keto check` for important tuples before and after
   the change.

## Differences

- The error positions start at line 1 and column 1. The OPL parser counted
  the columns from 0.
- The syntax-check API keeps its OPL name, `/opl/syntax/check`, and checks
  CUE.
- `keto namespace validate` checks CUE. The legacy YAML validation is
  `keto namespace validate-legacy`.

The upstream
[OPL specification](https://github.com/ory/keto/blob/master/docs/ory_permission_language_spec.md)
describes OPL in full.
