# decisions

<!-- okfw:generated — maintained by okfw; edits here are overwritten -->

## Concepts

* [ADR-0001: Adopt an OKF wiki (okfw) for project knowledge](adr-0001-adopt-an-okf-wiki-okfw-for-project-knowledge--01a1d39e124a.md) - An okfw-managed OKF wiki in .okfw/ holds keto-cue project knowledge: architecture, specs, ADRs, issues and daily notes.
* [ADR-0002: Replace OPL with CUE as the only schema language](adr-0002-replace-opl-with-cue-as-the-only-schema-language--01a1d39e124f.md) - keto-cue compiles CUE schemas into the existing Keto AST and removes the OPL lexer and parser.
* [ADR-0003: Build CUE instances in memory without cue/load or a module registry](adr-0003-build-cue-instances-in-memory-without-cue-load-or-a-module-registry--01a1d39e1254.md) - The schema compiler builds one CUE instance from source bytes and allows only standard-library imports.
* [ADR-0004: Add cuelang.org/go v0.17.1 as a dependency](adr-0004-add-cuelang-org-go-v0-17-1-as-a-dependency--01a1d39e125a.md) - keto-cue pins cuelang.org/go v0.17.1, the latest stable release, to compile CUE schemas.
* [ADR-0005: Parse errors carry the file name and 1-based positions](adr-0005-parse-errors-carry-the-file-name-and-1-based-positions--01a1d3adf8cf.md) - ParseError holds a file name, the source and 1-based start and end positions, and keeps the ToAPI and ToProto shapes.
* [ADR-0006: keto namespace validate checks CUE schemas](adr-0006-keto-namespace-validate-checks-cue-schemas--01a1d3cd0fe7.md) - keto namespace validate checks CUE files and directories, and the legacy validator keeps only the validate-legacy name.

