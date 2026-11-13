---
id: "01a1d39e-1254-7da9-b51f-398ee2cab6b5"
title: "ADR-0003: Build CUE instances in memory without cue/load or a module registry"
type: "ADR"
tags:
  - "adr"
  - "schema"
  - "cue"
  - "security"
created: "2026-11-13T00:35:14-06:00"
timestamp: "2026-11-13T00:35:14-06:00"
description: "The schema compiler builds one CUE instance from source bytes and allows only standard-library imports."
---

# ADR-0003: Build CUE instances in memory without cue/load or a module registry

**Status:** Accepted
**Date:** 2026-11-13
**Source:** Plan review

## Context
Keto loads a schema from a local file, a local directory, an `http(s)` URL or a `base64` URL. The watcher in `internal/driver/config/opl_config_namespace_watcher.go` holds each source as bytes, not as a path. `cue/load.Instances` reads from the filesystem. It also resolves `cue.mod` dependencies, and it can fetch modules from a CUE registry.

## Decision
- The compiler builds one `build.Instance` in memory: `build.NewContext().NewInstance`, then `AddFile` for each source, then `cuecontext.BuildInstance`.
- All files of one schema share one `package` clause and compile as one package.
- A schema can import the CUE standard library. A schema cannot import other packages or modules.

## Consequences
- A remote schema cannot read local files or reach a CUE registry.
- The same code path compiles local, remote and test schemas.
- Users cannot share helper packages through `cue.mod`. They put helper definitions in a file of the same package.

## Alternatives considered
- **`cue/load.Instances` with an overlay.** Rejected: it adds filesystem and registry access that a permission server does not need.
