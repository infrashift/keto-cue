---
title: Make Targets
description: The Makefile targets that lint, test and document keto-cue.
---

The keto-cue targets are at the end of the `Makefile`. The other targets come
from Ory Keto.

| Target | Command |
|---|---|
| `make lint` | `gofmt -l` and `go vet` on the keto-cue packages |
| `make test-cue` | `go test` on the schema, configuration, CLI and engine packages |
| `make wiki` | `okfw validate --all && okfw validate --okf` on the project wiki |
| `make docs-style` | `scripts/ste-check.py`: checks the documentation against ASD-STE100 rules |
| `make check` | All of the gates above |
| `make fuzz-test` | Fuzzes the CUE compiler for 10 seconds |
| `make docs-demo` | Builds keto and writes the landing-page output of `keto namespace validate` |
| `make docs-dev` | Serves this site with live reload (bun) |
| `make docs-build` | Builds this site into `docs/dist`, the same build as the Pages deploy |

Before a commit, run `make check`. When you change the site, also run
`make docs-build`. Fix each failure and run the gates again until all pass.
