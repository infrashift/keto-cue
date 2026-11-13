---
id: "01a1d3ad-f8cf-72b0-9334-2407ef14c6aa"
title: "ADR-0005: Parse errors carry the file name and 1-based positions"
type: "ADR"
tags:
  - "adr"
  - "schema"
  - "errors"
created: "2026-11-13T00:52:36-06:00"
timestamp: "2026-11-13T00:52:36-06:00"
description: "ParseError holds a file name, the source and 1-based start and end positions, and keeps the ToAPI and ToProto shapes."
---

# ADR-0005: Parse errors carry the file name and 1-based positions

**Status:** Accepted
**Date:** 2026-11-13
**Source:** [ISSUE-002](/issues/issue-002-cue-compiler-next-to-the-opl-parser--01a1d39e1264.md)

## Context
The OPL `ParseError` held a lexer item and a pointer to the parser. It computed positions from byte offsets on demand. Its column was off by one for the caret. A CUE schema can have many files, and CUE reports positions as `token.Pos` values.

## Decision
- `ParseError` holds the message, the file name, the file source and the start and end positions.
- Lines and columns are 1-based, as CUE reports them. A zero line means that the error has no position.
- `Error()` prints `<file>: error from L:C to L:C: <message>`, then the source lines with a caret under the range.
- `ToAPI()` and `ToProto()` keep their shapes. The message starts with `<file>: ` when the error has a file, because the API types have no file field.
- Semantic errors point at the value, not at the field label.
- An error against the base schema that has no user position takes the position of the user value at the same path.

## Consequences
- Clients of `/opl/syntax/check` get 1-based columns. OPL clients that corrected for the old offset must stop the correction.
- Multi-file schemas get errors that name the file.

## Alternatives considered
- **Add a file field to the API types.** Rejected: the `proto/` package does not change in keto-cue.
