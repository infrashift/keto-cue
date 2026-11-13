---
title: Loading Schemas
description: Load a keto-cue schema from a file, a directory, a URL or a base64 value, and how keto watches it.
---

Set `namespaces.location` in the Keto configuration to the URI of the schema:

```yaml
namespaces:
  location: file:///etc/keto/schema.cue
```

## Locations

| Location | Behavior |
|---|---|
| A local file | Keto reads the file as CUE, whatever its extension. |
| A local directory | All `*.cue` files of the directory form one CUE package. |
| An `http` or `https` URL | Keto fetches the file and keeps the compiled result for 30 minutes. |
| A `base64` URL | Keto decodes the value and reads it as one CUE file. |

## Directories

A directory schema is one CUE package. Each file must have the same `package`
clause. A file can refer to the definitions of the other files:

```cue
// schema/defs.cue
package acl

#owned: ["owners", "parents->edit"]
```

```cue
// schema/namespaces.cue
package acl

namespaces: {
	User: {}
	Folder: {
		relations: {owners: ["User"], parents: ["Folder"]}
		permissions: edit: or: #owned
	}
}
```

Keto ignores the files that do not end in `.cue` and logs a warning for each
one.

## Watching

Keto watches a local file or directory. A change compiles all files of the
schema again. When the new schema has an error, Keto keeps the last good
namespaces and logs each error with the file, the line and the column.

## Validate before you deploy

Run `keto namespace validate` on the same file or directory. It compiles the
schema as the server does:

```shell
keto namespace validate ./schema
```

See [CLI](../../reference/cli/).

## Imports

A schema can import the CUE standard library. A schema cannot import other
packages or modules. A remote schema cannot read local files.
