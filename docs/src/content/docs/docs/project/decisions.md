---
title: Design Decisions
description: A summary of the keto-cue architecture decision records (ADRs) in the project wiki.
---

The project wiki in `.okfw/decisions/` holds the full ADRs. Each ADR records
the context, the decision, the consequences and the rejected options. This
page gives a summary of each.

### ADR-0001: Adopt an OKF wiki for project knowledge

An OKF wiki in `.okfw/` holds the architecture, the specs, the decisions and
the issues. The `okfw` CLI manages it. `PLAN.md` stays the phase checklist.

### ADR-0002: Replace OPL with CUE as the only schema language

keto-cue reads CUE and does not read OPL. The CUE compiler writes the same
namespace tree as the OPL parser. The check engine, the expand engine, the
storage and the protos do not change.

### ADR-0003: Build CUE instances in memory

The compiler builds each schema from the source bytes. It does not use
`cue/load` or a module registry. A schema can import the CUE standard library
only. A remote schema cannot read local files or fetch modules.

### ADR-0004: Add cuelang.org/go as a dependency

`go.mod` pins one version of `cuelang.org/go`.

### ADR-0005: Parse errors carry the file name and 1-based positions

Each error names the file, the line and the column. Lines and columns start
at 1, as CUE reports them. The API messages start with the file name, because the API types have no file field.

### ADR-0006: keto namespace validate checks CUE schemas

`keto namespace validate` compiles a file or a directory as the server does.
The legacy validation keeps the name `validate-legacy`.

### ADR-0007: Short rewrite forms in the base schema

The base schema accepts `"relation"`, `"tupleSet->computed"`, `{or: [...]}`,
`{and: [...]}`, `{not: ...}` and a list of types for a relation. The Go
compiler converts each short form to the same tree as the long form. The ADR
rejects an importable CUE library: CUE has no functions, so a constructor call
is as long as the long form. See [Rewrites](../../guides/rewrites/).

### ADR-0008: Document keto-cue with Astro Starlight on GitHub Pages

This site lives in `docs/`. It uses Astro, Starlight and bun, the same setup
as the zmake site. GitHub Actions builds each pull request and deploys
`master` to GitHub Pages.
