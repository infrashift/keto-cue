---
title: Roadmap
description: The keto-cue phases that are done, and the work that is open.
---

`PLAN.md` in the repository holds the phase checklist. Each phase has a
mirror issue in the project wiki.

## Done

| Phase | Result |
|---|---|
| 0 | The okfw project wiki, the okfw MCP server and the documentation gates. |
| 1 | A CUE compiler that writes the Keto namespace tree. |
| 2 | `FormatCUE` and a test that compares each OPL fixture with its CUE version. |
| 3 | The server, the file watcher and the tests load CUE schemas. The syntax-check API has limits. |
| 4 | The OPL parser is gone. `keto namespace validate` checks CUE. |
| 5 | The CUE schema documentation. |
| 6 | Short forms for rewrites and relations. |
| 7 | This site, with a deploy to GitHub Pages. |

## Open

- A binary release.
- The generated SDK clients and the OpenAPI specification still describe the
  syntax-check endpoint with its OPL name.
- A memory limit for CUE evaluation in the syntax-check API.
