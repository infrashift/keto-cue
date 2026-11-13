---
title: Getting Started
description: Build keto-cue from source, write a CUE schema, start the server and check a permission.
---

This guide builds keto-cue, loads a schema and checks four permissions.

## Build keto-cue

keto-cue has no binary release. Build it from source with the Go version in
`go.mod`:

```shell
git clone https://github.com/infrashift/keto-cue
cd keto-cue
go build -o keto .
```

The result is one binary, `keto`. The commands below expect it on your `PATH`.

## Write a schema

Save this schema as `schema.cue`:

```cue
package keto

namespaces: {
	User: {}
	Group: relations: members: ["User", "Group#members"]
	Folder: {
		relations: {
			parents: ["Folder"]
			owners: ["User", "Group#members"]
			viewers: ["User", "Group#members"]
		}
		permissions: {
			edit: or: ["owners", "parents->edit"]
			view: or: ["viewers", "edit", "parents->view"]
		}
	}
	Doc: {
		relations: {
			parents: ["Folder"]
			owners: ["User"]
			blocked: ["User"]
		}
		permissions: {
			edit: or: ["owners", "parents->edit"]
			view: and: [{or: ["edit", "parents->view"]}, {not: "blocked"}]
		}
	}
}
```

- A **relation** lists the subject types that a relation tuple can hold.
  `"Group#members"` is the subject set `members` of a group.
- A **permission** is a rule. `"owners"` checks the relation `owners` of the
  same object. `"parents->edit"` checks `edit` on each parent.
- `{not: "blocked"}` inverts a result, so a blocked user cannot view a
  document.

## Validate the schema

```shell
keto namespace validate schema.cue
```

```text
schema.cue: valid, 4 namespace(s)
```

An error names the file, the line and the column. See
[Rules & Errors](../reference/rules-and-errors/).

## Start the server

Save this configuration as `keto.yml`:

```yaml
dsn: memory
namespaces:
  location: file://./schema.cue
```

Start the server:

```shell
keto serve -c keto.yml
```

The read API listens on port 4466 and the write API on port 4467. The `memory`
DSN keeps the relation tuples in memory, so they go away when the server
stops.

## Add relation tuples

Open a second terminal. Tell the CLI where the server is:

```shell
export KETO_READ_REMOTE=127.0.0.1:4466
export KETO_WRITE_REMOTE=127.0.0.1:4467
```

Add five relation tuples. The local server has no TLS, so each command needs
`--insecure-disable-transport-security`:

```shell
keto relation-tuple create Group:eng#members owners Folder:docs --insecure-disable-transport-security
keto relation-tuple create User:alice members Group:eng --insecure-disable-transport-security
keto relation-tuple create Folder:docs parents Doc:readme --insecure-disable-transport-security
keto relation-tuple create User:bob viewers Folder:docs --insecure-disable-transport-security
keto relation-tuple create User:bob blocked Doc:readme --insecure-disable-transport-security
```

The arguments are the subject, the relation and the object. The first tuple
makes the members of the group `eng` owners of the folder `docs`.

## Check permissions

```shell
keto check User:alice edit Doc:readme --insecure-disable-transport-security
keto check User:alice view Doc:readme --insecure-disable-transport-security
keto check User:bob view Doc:readme --insecure-disable-transport-security
keto check User:bob view Folder:docs --insecure-disable-transport-security
```

```text
Allowed
Allowed
Denied
Allowed
```

- Alice is a member of `eng`, `eng` owns the folder, and the folder is the
  parent of the document. So Alice can edit and view the document.
- Bob can view the folder. The `blocked` tuple denies Bob on the document.

## Expand a relation

`keto expand` shows the tree of subjects that have a relation:

```shell
keto expand owners Folder:docs --insecure-disable-transport-security
```

```text
or :#@Folder:docs#owners
└──or :#@Group:eng#members
   └──∋ :#@User:alice
```

## Next steps

- [Permission Schemas](../guides/permission-schemas/) explains the structure of a schema.
- [Rewrites](../guides/rewrites/) lists each form of a permission rule.
- [Loading Schemas](../guides/loading-schemas/) explains files, directories and remote schemas.
