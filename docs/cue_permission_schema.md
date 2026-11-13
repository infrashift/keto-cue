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
	Group: relations: members: types: ["User", "Group#members"]
	Folder: {
		relations: {
			parents: types: ["Folder"]
			viewers: types: ["User", "Group#members"]
		}
		permissions: view: {
			op: "or"
			children: [
				{op: "direct", relation: "viewers"},
				{op: "tupleToSubjectSet", tupleSetRelation: "parents", computedSubjectSet: "view"},
			]
		}
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

A relation can also have a `rewrite`, with the same form as a permission.

### Permissions

A permission is a subject-set rewrite. Each rewrite node has an `op` field:

| `op` | Fields | Check result |
|---|---|---|
| `or` | `children` (one or more rewrites) | Allowed when one child allows. |
| `and` | `children` (one or more rewrites) | Allowed when all children allow. |
| `direct` | `relation` | Allowed when the relation or permission of the same object allows. |
| `tupleToSubjectSet` | `tupleSetRelation`, `computedSubjectSet` | Follows each object in `tupleSetRelation` and checks `computedSubjectSet` on it. |
| `invert` | `child` | Inverts the result of the child. |

## OPL to CUE

| OPL | CUE |
|---|---|
| `class Doc implements Namespace {}` | `namespaces: Doc: {}` |
| `related: { owners: User[] }` | `relations: owners: types: ["User"]` |
| `owners: (User \| SubjectSet<Group, "members">)[]` | `owners: types: ["User", "Group#members"]` |
| `this.related.owners.includes(ctx.subject)` | `{op: "direct", relation: "owners"}` |
| `this.permits.edit(ctx)` | `{op: "direct", relation: "edit"}` |
| `this.related.parents.traverse((p) => p.permits.view(ctx))` | `{op: "tupleToSubjectSet", tupleSetRelation: "parents", computedSubjectSet: "view"}` |
| `a \|\| b` | `{op: "or", children: [a, b]}` |
| `a && b` | `{op: "and", children: [a, b]}` |
| `!a` | `{op: "invert", child: a}` |

## Reuse

The top level of a schema is open. Declare helper definitions next to `namespaces`. Use hidden fields (`_name`) as parameters:

```cue
#viaParent: {
	_permission:        string
	op:                 "tupleToSubjectSet"
	tupleSetRelation:   "parents"
	computedSubjectSet: _permission
}

namespaces: Folder: permissions: view: #viaParent & {_permission: "view"}
```

A schema can import the CUE standard library, for example `strings` or `list`. A schema cannot import other packages or modules.

## Validate a schema

Run `keto namespace validate` on a file or a directory:

```shell
keto namespace validate ./schema
```

Each error names the file, the line and the column, and shows the source:

```text
schema/doc.cue: error from 3:43 to 3:49: namespace "A" did not declare relation "nope"

   2 | 	relations: r: types: ["A"]
   3 | 	permissions: p: {op: "direct", relation: "nope"}
     | 	                                         ^~~~~~
```

The syntax-check API (`POST /opl/syntax/check` and the gRPC `SyntaxService.Check`) checks one CUE file. The endpoint keeps its OPL name for API compatibility. Lines and columns start at 1.

## Rules

The compiler applies these rules:

- The base schema is closed. A misspelled field is an error.
- Each rewrite has the fields of its `op` and no other rewrite fields.
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
