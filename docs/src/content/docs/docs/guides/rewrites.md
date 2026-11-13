---
title: Rewrites
description: The short forms and the long form of a keto-cue permission rule, and how they nest.
---

A rewrite is the rule of a permission. The check engine evaluates it for one
subject and one object.

## Short forms

| Rewrite | Result |
|---|---|
| `"viewers"` | Allowed when the relation or permission `viewers` of the same object allows. |
| `"parents->view"` | Follows each object in `parents` and checks `view` on it. |
| `{or: [a, b, ...]}` | Allowed when one rewrite allows. |
| `{and: [a, b, ...]}` | Allowed when all rewrites allow. |
| `{not: a}` | Inverts the result of `a`. |

CUE lets you drop the braces of a struct with one field. `view: {or: [...]}`
and `view: or: [...]` are the same.

## Nesting

Put a rewrite in a list or under `not` to make a tree:

```cue
permissions: {
	edit: or: ["owners", "parents->edit"]
	view: and: [
		{or: ["viewers", "edit", "parents->view"]},
		{not: "blocked"},
	]
}
```

The check engine flattens nested operations of the same kind, so
`{or: ["a", {or: ["b", "c"]}]}` is the same as `{or: ["a", "b", "c"]}`.
Rewrites nest to a depth of 10 or less.

## Traversal

`"parents->view"` is a tuple-to-subject-set rewrite. The engine reads the
relation tuples of `parents` on the object. For each object in `parents`, it
checks `view`. Each type of `parents` must declare `view`.

```cue
Doc: {
	relations: parents: ["Folder"]
	permissions: view: "parents->view"
}
```

The compiler accepts spaces around the arrow: `"parents -> view"`.

## Long form

The long form has an `op` field. It is the form that keto-cue used before the
short forms. A schema can mix the two forms.

| `op` | Fields | Short form |
|---|---|---|
| `or` | `children` (one or more rewrites) | `{or: [...]}` |
| `and` | `children` (one or more rewrites) | `{and: [...]}` |
| `direct` | `relation` | `"relation"` |
| `tupleToSubjectSet` | `tupleSetRelation`, `computedSubjectSet` | `"tupleSetRelation->computedSubjectSet"` |
| `invert` | `child` | `{not: ...}` |

```cue
permissions: view: {
	op: "or"
	children: [
		{op: "direct", relation: "viewers"},
		{op: "tupleToSubjectSet", tupleSetRelation: "parents", computedSubjectSet: "view"},
	]
}
```

Use the long form for a relation name that contains `->`, or that starts or
ends with a space.

## One form per node

A rewrite node has one of the fields `or`, `and`, `not` and `op`. A node with
two of them is an error. A node with `or`, `and` or `not` has no other field.
