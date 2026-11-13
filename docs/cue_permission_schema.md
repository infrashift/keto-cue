# CUE permission schema

keto-cue reads permission schemas in [CUE](https://cuelang.org). A schema declares the namespaces, their relations and their permissions. This document replaces the [Ory Permission Language specification](ory_permission_language_spec.md) for keto-cue.

## Load a schema

Set `namespaces.location` in the Keto config to the URI of the schema:

```yaml
namespaces:
  location: file:///etc/keto/schema.cue
```

The location can be one of these:

- A local file. Keto reads the file as CUE, whatever its extension.
- A local directory. All `*.cue` files of the directory form one CUE package. Keto ignores other files and logs a warning for each one.
- An `http` or `https` URL. Keto caches the compiled result for 30 minutes.
- A `base64` URL.

Keto watches local files. A change compiles all files again. A broken change keeps the last good namespaces and logs the errors.

## Structure

A schema declares `namespaces`. Each field of `namespaces` is one namespace. A namespace has optional `relations` and optional `permissions`.

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

### Relations

A relation lists the subject types that a relation tuple can hold:

| Form | Meaning |
|---|---|
| `"User"` | A subject in the namespace `User`. |
| `"Group#members"` | The subject set `members` of an object in `Group`. |
| `{namespace: "Group", relation: "members"}` | The same subject set, in the long form. |

Write the list of types as the value of the relation: `owners: ["User"]`. The long form is `owners: types: ["User"]`. A relation can also have a `rewrite`, with the same form as a permission. A relation with a rewrite uses the long form: `{types: ["User"], rewrite: "owners"}`.

### Permissions

A permission is a subject-set rewrite. Use the short forms:

| Rewrite | Check result |
|---|---|
| `"viewers"` | Allowed when the relation or permission `viewers` of the same object allows. |
| `"parents->view"` | Follows each object in `parents` and checks `view` on it. |
| `{or: [a, b, ...]}` | Allowed when one rewrite allows. |
| `{and: [a, b, ...]}` | Allowed when all rewrites allow. |
| `{not: a}` | Inverts the result of `a`. |

A rewrite node has one of the fields `or`, `and`, `not` and `op`. CUE lets you drop the braces of a struct with one field, so `view: {or: [...]}` and `view: or: [...]` are the same. Nest the forms to make a tree:

```cue
permissions: delete: and: ["owners", {not: "parents->locked"}]
```

The long form has an `op` field. Use it when a relation name contains `->` or starts or ends with a space:

| `op` | Fields | Short form |
|---|---|---|
| `or` | `children` (one or more rewrites) | `{or: [...]}` |
| `and` | `children` (one or more rewrites) | `{and: [...]}` |
| `direct` | `relation` | `"relation"` |
| `tupleToSubjectSet` | `tupleSetRelation`, `computedSubjectSet` | `"tupleSetRelation->computedSubjectSet"` |
| `invert` | `child` | `{not: ...}` |

A schema can mix the short and the long forms.

## OPL to CUE

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

## Reuse

The top level of a schema is open. Declare helper values next to `namespaces`, and refer to them. Use a comprehension to add the values of a list to another list:

```cue
#editors: ["owners", "parents->edit"]

namespaces: Folder: permissions: {
	edit: or: #editors
	view: or: ["viewers", "parents->view", for e in #editors {e}]
}
```

A comprehension can also make a permission for each entry of a list:

```cue
#actions: ["read", "write", "delete"]

namespaces: Doc: permissions: {
	for a in #actions {(a): or: ["owners", "parents->\(a)"]}
}
```

A schema can import the CUE standard library, for example `strings` or `list`. A schema cannot import other packages or modules.

## Validate a schema

Run `keto namespace validate` on a file or a directory:

```shell
keto namespace validate ./schema
```

Each error names the file, the line and the column, and shows the source:

```text
schema/doc.cue: error from 3:18 to 3:24: namespace "A" did not declare relation "nope"

   2 | 	relations: r: ["A"]
   3 | 	permissions: p: "nope"
     | 	                ^~~~~~
```

The syntax-check API (`POST /opl/syntax/check` and the gRPC `SyntaxService.Check`) checks one CUE file. The endpoint keeps its OPL name for API compatibility. Lines and columns start at 1.

## Rules

The compiler applies these rules:

- The base schema is closed. A misspelled field is an error.
- Each rewrite node has one of the fields `op`, `or`, `and` and `not`.
- Each rewrite with `op` has the fields of its `op` and no other rewrite fields.
- A short rewrite is `"relation"` or `"tupleSetRelation->computedSubjectSet"`, with no empty name.
- Each referenced namespace and relation exists.
- For `tupleToSubjectSet`, each type of `tupleSetRelation` declares `computedSubjectSet`.
- A name cannot be a relation and a permission in the same namespace.
- Rewrites nest to a depth of 10 or less.

## Limits

The syntax-check API accepts schemas from clients. These limits protect the server:

- The request body is 1 MiB or less. A larger body gets HTTP 413 or gRPC `ResourceExhausted`.
- The number of concurrent checks is the number of CPUs. An extra request waits up to 10 seconds, then gets HTTP 429 or gRPC `ResourceExhausted`.
- A check that takes more than 10 seconds gets HTTP 408 or gRPC `DeadlineExceeded`.

Go cannot stop a running CUE evaluation. A timed-out check keeps its slot until the evaluation ends. CUE comprehensions and functions such as `list.Range` can use much memory within these limits. Put the syntax-check listener behind an access control when untrusted clients can reach it.
