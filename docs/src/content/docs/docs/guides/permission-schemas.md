---
title: Permission Schemas
description: The structure of a keto-cue schema - namespaces, relations, subject types and permissions.
---

A keto-cue schema is a CUE value with one field, `namespaces`. Each field of
`namespaces` is one namespace. A namespace has optional `relations` and
optional `permissions`.

```cue
package keto

namespaces: {
	User: {}
	Group: relations: members: ["User", "Group#members"]
	Folder: {
		relations: {
			parents: ["Folder"]
			viewers: ["User", "Group#members"]
		}
		permissions: view: or: ["viewers", "parents->view"]
	}
}
```

The `package` clause is optional for a single file. All files of a schema
directory must have the same `package` clause.

## Namespaces

A namespace is a type of object, for example `Folder` or `Doc`. A namespace
with no relations, such as `User`, is a type of subject. Write it as `{}`.

## Relations

A relation lists the subject types that a relation tuple can hold. Write the
list as the value of the relation:

```cue
relations: {
	owners: ["User", "Group#members"]
	parents: ["Folder"]
}
```

Each subject type has one of these forms:

| Form | Meaning |
|---|---|
| `"User"` | A subject in the namespace `User`. |
| `"Group#members"` | The subject set `members` of an object in `Group`. |
| `{namespace: "Group", relation: "members"}` | The same subject set, in the long form. |

The long form of a relation is a struct with `types` and an optional
`rewrite`:

```cue
relations: editors: {
	types: ["User"]
	rewrite: or: ["editors", "owners"]
}
```

## Permissions

A permission is a rule that the check engine evaluates. It refers to the
relations and the permissions of the same namespace, and it can follow a
relation to other objects:

```cue
permissions: {
	edit: "owners"
	view: or: ["viewers", "edit", "parents->view"]
}
```

[Rewrites](../rewrites/) lists each form of a rule.

## Names

A relation and a permission cannot have the same name in one namespace.
Quote a name that is not a CUE identifier:

```cue
namespaces: {
	User: {}
	"my-app": relations: "can-read": ["User"]
}
```

## Checks

The compiler checks the schema before the server uses it. Each referenced
namespace and relation must exist. A misspelled field is an error.
[Rules & Errors](../../reference/rules-and-errors/) lists each rule.
