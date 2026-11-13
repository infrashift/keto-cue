# PLAN: Replace the Ory Keto OPL parser with a CUE compiler

keto-cue reads permission schemas in CUE. The CUE compiler writes the same `[]namespace.Namespace` and `ast` trees that the OPL parser wrote. The check engine, the expand engine, the persistence layer and the protos do not change. See [ADR-0002](.okfw/decisions/adr-0002-replace-opl-with-cue-as-the-only-schema-language--01a1d39e124f.md) and the [CUE permission schema](.okfw/specs/cue-permission-schema--01a1d39e1284.md) spec.

Work one phase at a time. Follow the Development workflow guide in `.okfw/`.

## Phase 0: Wiki, MCP server and gates
Tracking: [ISSUE-001](.okfw/issues/issue-001-okfw-wiki-mcp-server-and-documentation-gates--01a1d39e125f.md)

- [x] Run `okfw init` and `okfw install --platform claude`.
- [x] Register `okfw mcp` in `.mcp.json`.
- [x] Add the guides, ADR-0001 to ADR-0004, one issue per phase, the architecture note and the spec.
- [x] Add `scripts/ste-check.py` and the `lint`, `test-cue`, `wiki`, `docs-style` and `check` Makefile targets.

## Phase 1: CUE compiler next to the OPL parser
Tracking: [ISSUE-002](.okfw/issues/issue-002-cue-compiler-next-to-the-opl-parser--01a1d39e1264.md)

- [x] Add `cuelang.org/go v0.17.1` ([ADR-0004](.okfw/decisions/adr-0004-add-cuelang-org-go-v0-17-1-as-a-dependency--01a1d39e125a.md)).
- [x] Embed `schema/cue_base.cue`. Keep the definitions closed and the top level open.
- [x] Add `ParseCUE(files []SourceFile)` in `schema/cue_compiler.go`. Build the instance in memory ([ADR-0003](.okfw/decisions/adr-0003-build-cue-instances-in-memory-without-cue-load-or-a-module-registry--01a1d39e1254.md)).
- [x] Map `or` and `and` to `*ast.SubjectSetRewrite`, `direct` to `*ast.ComputedSubjectSet`, `tupleToSubjectSet` to `*ast.TupleToSubjectSet` and `invert` to `*ast.InvertResult`.
- [x] Wrap a top-level leaf with `AsRewrite()`. Run `simplifyExpression` on each rewrite.
- [x] Port the OPL type checks to `schema/cue_checks.go`. Keep the error text. Add source positions.
- [x] Change `ParseError` to hold a file name, the source and the start and end positions. Keep `Error()`, `ToAPI()` and `ToProto()`.
- [x] Add snapshot, negative and multi-file tests in `schema/cue_compiler_test.go`.

## Phase 2: CUE formatter and OPL equivalence gate
Tracking: [ISSUE-003](.okfw/issues/issue-003-cue-formatter-and-opl-equivalence-gate--01a1d39e1269.md)

- [ ] Add `FormatCUE([]namespace) []byte` in `schema/cue_format.go`.
- [ ] Add a test: OPL fixture → `Parse` → `FormatCUE` → `ParseCUE` → the same AST. Cover every OPL fixture.

## Phase 3: Load CUE schemas in the server and the tests
Tracking: [ISSUE-004](.okfw/issues/issue-004-load-cue-schemas-in-the-server-and-the-tests--01a1d39e126e.md)

- [ ] Change `schema.Parse` to take CUE source files. Update the watcher, the handler and the test helpers.
- [ ] Keep file bytes in the watcher, not readers. Compile all `.cue` files of a directory as one package.
- [ ] Cap the request body of the syntax-check handler.
- [ ] Rename `driver.WithOPL` to `WithSchema` and `Scenario.Opl` to `Scenario.Schema`.
- [ ] Change the `namespaces.location` examples in the config schema to `.cue`.

## Phase 4: Migrate the fixtures and remove OPL
Tracking: [ISSUE-005](.okfw/issues/issue-005-migrate-the-fixtures-and-remove-opl--01a1d39e1273.md)

- [ ] Convert the OPL fixtures in the engine and config tests to CUE with `FormatCUE`. Port the error cases by hand.
- [ ] Convert `internal/check/testfixtures/project_opl.ts` and the `contrib/` examples to `.cue`.
- [ ] Delete the lexer, the parser, the fuzz corpus, the lexer and parser snapshots and `contrib/namespace-type-lib/`.
- [ ] Add `FuzzParseCUE`. Keep `Test_simplify`.
- [ ] Replace `GenerateOPLConfig` with `GenerateCUEConfig`.
- [ ] Add `keto namespace validate <file.cue|dir>`. Keep the legacy command as `validate-legacy`.

## Phase 5: CUE schema documentation
Tracking: [ISSUE-006](.okfw/issues/issue-006-cue-schema-documentation--01a1d39e1279.md)

- [ ] Write `docs/cue_permission_schema.md` with an OPL to CUE table.
- [ ] Mark `docs/ory_permission_language_spec.md` as superseded. Update the OPL examples in `README.md`.

## Verification
Run these gates after each phase. Do not start the next phase while a gate fails.

- `make lint test-cue wiki docs-style`
- Phase 2 gate: the equivalence test passes for every OPL fixture.
- End to end: serve `contrib/cat-videos-example`, then compare `keto check` and `keto expand` with the results before the migration.
