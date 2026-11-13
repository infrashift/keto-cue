---
id: "01a1d3ad-f8f1-7347-a479-c6b41ac6e533"
title: "ISSUE-007: Bound CUE evaluation on the syntax-check listener"
type: "Issue"
kind: "issue"
tags:
  - "issue"
  - "open"
  - "phase-3"
  - "security"
created: "2026-11-13T00:52:36-06:00"
timestamp: "2026-11-13T00:52:36-06:00"
description: "CUE evaluation has no time or memory bound, and the syntax-check listener accepts schemas from clients."
---

# ISSUE-007: Bound CUE evaluation on the syntax-check listener

**Status:** Open. **Phase:** 3. The handler limits are done. The document item waits for Phase 5.

## Description
CUE evaluation has no time or memory bound. A schema with nested comprehensions or a large `list.Range` can use much CPU and memory. The syntax-check listener accepts schemas from clients, so a client can send such a schema. Schemas from the config file come from the operator.

## Acceptance Criteria
- [x] The syntax-check handler caps the request body at 1 MiB (HTTP 413, gRPC `ResourceExhausted`).
- [x] The syntax-check handler limits the concurrent evaluations to `GOMAXPROCS` (HTTP 429, gRPC `ResourceExhausted`).
- [x] The syntax-check handler returns an error when an evaluation exceeds 10 seconds (HTTP 408, gRPC `DeadlineExceeded`).
- [ ] The CUE schema document states the remaining risk for memory use.

## Notes
- Go cannot stop a running CUE evaluation. A timed-out evaluation keeps its slot until it ends, so the concurrency limit also bounds the CPU use.
- A fuzz run of 4 minutes found no slow input. Each corpus input took about 1 ms or less.

## Related
- [ADR-0003](/decisions/adr-0003-build-cue-instances-in-memory-without-cue-load-or-a-module-registry--01a1d39e1254.md)
- [ISSUE-004](/issues/issue-004-load-cue-schemas-in-the-server-and-the-tests--01a1d39e126e.md)
