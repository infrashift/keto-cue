# issues

<!-- okfw:generated — maintained by okfw; edits here are overwritten -->

## Concepts

* [ISSUE-001: Okfw wiki, MCP server and documentation gates](issue-001-okfw-wiki-mcp-server-and-documentation-gates--01a1d39e125f.md) - Set up the okfw wiki, the okfw MCP server, the STE checker and the keto-cue Makefile gates.
* [ISSUE-002: CUE compiler next to the OPL parser](issue-002-cue-compiler-next-to-the-opl-parser--01a1d39e1264.md) - Add the base schema, the CUE-to-AST compiler, the semantic checks and the new ParseError shape.
* [ISSUE-003: CUE formatter and OPL equivalence gate](issue-003-cue-formatter-and-opl-equivalence-gate--01a1d39e1269.md) - Add FormatCUE and a test that proves each OPL fixture and its CUE port give the same AST.
* [ISSUE-004: Load CUE schemas in the server and the tests](issue-004-load-cue-schemas-in-the-server-and-the-tests--01a1d39e126e.md) - Route schema.Parse, the watcher, the syntax-check handler and the test helpers through the CUE compiler.
* [ISSUE-005: Migrate the fixtures and remove OPL](issue-005-migrate-the-fixtures-and-remove-opl--01a1d39e1273.md) - Convert the OPL fixtures and examples to CUE, delete the OPL code and add the CUE validate command.
* [ISSUE-006: CUE schema documentation](issue-006-cue-schema-documentation--01a1d39e1279.md) - Document the CUE permission schema and mark the OPL documents as superseded.
* [ISSUE-007: Bound CUE evaluation on the syntax-check listener](issue-007-bound-cue-evaluation-on-the-syntax-check-listener--01a1d3adf8f1.md) - CUE evaluation has no time or memory bound, and the syntax-check listener accepts schemas from clients.
* [ISSUE-008: Short forms for CUE rewrites and relations](issue-008-short-forms-for-cue-rewrites-and-relations--01a1d4cf27e7.md) - Add short forms for rewrites and relation types so that a CUE schema is as short as an OPL schema.
* [ISSUE-009: Documentation site on GitHub Pages](issue-009-documentation-site-on-github-pages--01a1d4da03a0.md) - Publish the keto-cue documentation as an Astro Starlight site on GitHub Pages, with the zmake setup.

