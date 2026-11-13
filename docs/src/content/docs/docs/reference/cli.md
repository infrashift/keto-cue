---
title: CLI
description: The keto-cue commands for schemas, and the changes from the Ory Keto CLI.
---

keto-cue keeps the Ory Keto CLI. This page lists the commands that work with
CUE schemas. Run `keto help` for all commands.

## keto namespace validate

```text
keto namespace validate <schema.cue | schema-dir> [...]
```

Compiles each argument as the server does. A directory is one CUE package of
its `*.cue` files.

- A valid schema prints `<path>: valid, <n> namespace(s)`.
- An invalid schema prints each error to standard error, with the file, the
  line, the column and a source snippet. The command then exits with an
  error.
- A directory with no `*.cue` files is an error.

```shell
keto namespace validate schema.cue ./more-schemas
```

## keto namespace validate-legacy

```text
keto namespace validate-legacy <namespace.yml> [...]
keto namespace validate-legacy -c <config.yaml>
```

Validates the deprecated YAML, JSON or TOML namespace files. In Ory Keto,
`validate` was an alias of this command. In keto-cue, `validate` checks CUE.

## keto serve

```text
keto serve -c keto.yml
```

Starts the server. `namespaces.location` in the configuration points at the
schema. See [Loading Schemas](../../guides/loading-schemas/).

## Client commands

The client commands are the Ory Keto commands:

| Command | Purpose |
|---|---|
| `keto relation-tuple create <subject> <relation> <object>` | Adds a relation tuple. `-f` reads JSON files. |
| `keto check <subject> <relation> <object>` | Prints `Allowed` or `Denied`. |
| `keto expand <relation> <object>` | Prints the tree of subjects. |

`KETO_READ_REMOTE` and `KETO_WRITE_REMOTE`, or `--read-remote` and
`--write-remote`, set the server addresses. Use
`--insecure-disable-transport-security` for a local server without TLS.
