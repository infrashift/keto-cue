---
id: "01a1d39e-127e-7929-9b97-ee66fb238f9e"
title: "keto-cue architecture overview"
type: "Architecture"
tags:
  - "architecture"
  - "schema"
  - "cue"
created: "2026-11-13T00:35:14-06:00"
timestamp: "2026-11-13T00:35:14-06:00"
description: "How a CUE schema flows from a config location to the AST that the check and expand engines walk."
---

# keto-cue architecture overview

## Data flow
1. The config key `namespaces.location` names a file, a directory, an `http(s)` URL or a `base64` URL.
2. `internal/driver/config/opl_config_namespace_watcher.go` reads the source bytes. It watches local files and caches remote results. A directory target gives one CUE package. The watcher compiles the `*.cue` files of the directory together. It ignores other files. Each change compiles all files again. A broken change keeps the last good namespaces.
3. `schema.Parse` builds one CUE instance in memory ([ADR-0003](/decisions/adr-0003-build-cue-instances-in-memory-without-cue-load-or-a-module-registry--01a1d39e1254.md)). It unifies the instance with the embedded base schema.
4. The compiler walks `namespaces` and writes `[]namespace.Namespace`. Each namespace holds `[]ast.Relation`.
5. The semantic checks verify each reference. Each error carries a file name, a line and a column.
6. The memory namespace manager serves the namespaces to the check and expand engines.

## Package boundaries
| Package | Role | Changes in keto-cue |
|---|---|---|
| `schema/` | Compiler, checks, errors, syntax-check handler | Yes |
| `internal/namespace/ast/` | AST types | No |
| `internal/driver/config/` | Load and watch the schema | Yes |
| `cmd/namespace/` | `keto namespace` commands | Yes |
| `internal/check/`, `internal/expand/` | Engines | Test files only |
| `internal/persistence/`, `proto/` | Storage and API | No |

## Related
- [ADR-0002](/decisions/adr-0002-replace-opl-with-cue-as-the-only-schema-language--01a1d39e124f.md)
