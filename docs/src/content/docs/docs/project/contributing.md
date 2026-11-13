---
title: Contributing
description: The keto-cue development workflow, the gates and the rules for a change.
---

## Workflow

1. Read `PLAN.md`. Work on the first phase that is not done.
2. Search the project wiki before you write code: `okfw search "<query>"`.
3. Record each decision as an ADR in `.okfw/decisions/`.
4. Implement, lint and test. Fix each failure and run all gates again.
5. Update the site pages for the change. Run `make docs-build`.
6. Close the issue of the phase and commit the phase.

## Gates

```shell
make check
make docs-build
```

`make check` runs `lint`, `test-cue`, `wiki` and `docs-style`. See
[Make Targets](../../reference/make-targets/).

## Rules

- Do not change the production code in `internal/check`, `internal/expand`,
  `internal/persistence` or `proto/`. Test files in those packages can
  change.
- Do not weaken or delete a test to make a gate pass.
- Write the documentation in ASD-STE100 Simplified Technical English: active
  voice, simple present tense and short sentences.
