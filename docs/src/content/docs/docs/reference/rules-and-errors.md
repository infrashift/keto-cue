---
title: Rules & Errors
description: The rules that the keto-cue compiler applies to a schema, and the error for each rule.
---

The compiler reports CUE errors first. It checks the references only when the
schema has no CUE errors.

## Error format

Each error names the file, the start and the end of the value, and shows the
source with a marker:

```text
typo.cue: error from 25:30 to 25:45: relation "veiw" was not declared in namespace "Folder"

  24 | 			edit: or: ["owners", "parents->edit"]
  25 | 			view: and: [{or: ["edit", "parents->veiw"]}, {not: "blocked"}]
     | 			                          ^~~~~~~~~~~~~~~
  26 | 		}
```

Lines and columns start at 1. A schema from a `base64` location or from the
syntax-check API has no file name.

## Rules

| Rule | Example error |
|---|---|
| The base schema is closed. A misspelled field is an error. | `namespaces.A.permissions.p.orr: field not allowed` |
| A value has the type that the base schema allows. | `namespaces.A.relations.r: "A" is not allowed; use a list of types or {types?: [...], rewrite?: ...}` |
| An `op` value is one of the five operations. | `namespaces.A.permissions.p.op: "xor" is not allowed; use "or", "and", "direct", "tupleToSubjectSet" or "invert"` |
| A subject type is `"Namespace"` or `"Namespace#relation"`. | `namespaces.A.relations.r.types.0: "A#b#c" is not allowed; ...` |
| Each value is concrete. | `namespaces.A.relations.r.types.0: incomplete value ...` |
| A rewrite node has one of `op`, `or`, `and` and `not`. | `a rewrite requires one of the fields "op", "or", "and" or "not"` |
| A rewrite node has only one of them. | `a rewrite can have only one of the fields "op" and "or"` |
| A rewrite with `op` has the fields of its `op`. | `a rewrite with op "direct" requires the field "relation"` |
| A rewrite with `op` has no fields of other ops. | `field "child" is not allowed in a rewrite with op "direct"` |
| A short operator node has no other field. | `field "relation" is not allowed in a rewrite with the field "not"` |
| A short rewrite has one arrow and no empty name. | `the rewrite "r->" is not valid; use "relation" or "tupleSetRelation->computedSubjectSet"` |
| An `or` or `and` list has one or more rewrites. | `namespaces.A.permissions.p.or: incompatible list lengths (0 and 1)` |
| Each referenced relation exists in the namespace. | `namespace "A" did not declare relation "nope"` |
| Each namespace in a subject type exists. | `namespace "Nope" was not declared` |
| Each type of a traversed relation declares the computed relation. | `relation "view" was not declared in namespace "Folder"` |
| A name is not a relation and a permission in one namespace. | `namespace "A" declares "r" as a relation and as a permission` |
| Rewrites nest to a depth of 10 or less. | `expression nested too deeply; maximal nesting depth is 10` |
| All files of a directory have the same `package` clause. | `package name "b" conflicts with previous package name "a"` |
| A schema imports only the CUE standard library. | `import "example.com/foo" not found` |

## Duplicate errors

The compiler reports an error with the same message at the same position
once. A comprehension can make the same error for each entry. The compiler
reports it once.
