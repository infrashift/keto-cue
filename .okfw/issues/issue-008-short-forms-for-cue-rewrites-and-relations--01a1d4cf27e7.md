---
id: "01a1d4cf-27e7-7293-9f50-17d34eadefae"
title: "ISSUE-008: Short forms for CUE rewrites and relations"
type: "Issue"
tags:
  - "issue"
  - "phase-6"
  - "cue"
  - "schema"
  - "closed"
created: "2026-11-13T06:08:28-06:00"
timestamp: "2026-11-13T06:15:22-06:00"
description: "Add short forms for rewrites and relation types so that a CUE schema is as short as an OPL schema."
---

# ISSUE-008: Short forms for CUE rewrites and relations

**Status:** Closed. **Phase:** 6.

## Description
The long rewrite form is verbose. One OPL traversal, `this.related.parents.traverse((p) => p.permits.view(ctx))`, becomes `{op: "tupleToSubjectSet", tupleSetRelation: "parents", computedSubjectSet: "view"}`. A schema must be as short as an OPL schema. The long form stays valid.

## Acceptance Criteria
- [x] A rewrite can be `"relation"` (direct) or `"tupleSetRelation->computedSubjectSet"` (tupleToSubjectSet).
- [x] A rewrite can be `{or: [...]}`, `{and: [...]}` or `{not: rewrite}`.
- [x] A relation can be its list of types: `owners: ["User", "Group#members"]`.
- [x] Each short form compiles to the same AST as its long form. A test compares the two forms of the full example.
- [x] Each error in a short form has the file, the line, the column and a clear message.
- [x] `FormatCUE` writes the short forms. A name that a short form cannot hold keeps the long form.
- [x] The `contrib` examples and the benchmark fixture use the short forms.
- [x] The user document, the README example and the spec concept show the short forms.

## Resolution
The base schema accepts the short forms, and the compiler converts them to the long-form AST ([ADR-0007](/decisions/adr-0007-short-rewrite-forms-in-the-base-schema-compiled-in-go--01a1d4cf57cd.md)). `FormatCUE` writes the short forms. The `contrib` examples and the benchmark fixture use them, and a test proves that each converted file compiles to the same AST as before. The compiler reports a duplicate error from a comprehension once.

## Notes
- The engine test fixtures keep the long form, so the tests cover both forms.

## Related
- [ADR-0007](/decisions/adr-0007-short-rewrite-forms-in-the-base-schema-compiled-in-go--01a1d4cf57cd.md)
- [CUE permission schema](/specs/cue-permission-schema--01a1d39e1284.md)
