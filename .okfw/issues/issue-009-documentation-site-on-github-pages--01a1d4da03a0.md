---
id: "01a1d4da-03a0-7ba9-94a8-08b2a3f8f116"
title: "ISSUE-009: Documentation site on GitHub Pages"
type: "Issue"
tags:
  - "issue"
  - "phase-7"
  - "docs"
  - "site"
  - "closed"
created: "2026-11-13T06:20:20-06:00"
timestamp: "2026-11-13T06:29:52-06:00"
description: "Publish the keto-cue documentation as an Astro Starlight site on GitHub Pages, with the zmake setup."
---

# ISSUE-009: Documentation site on GitHub Pages

**Status:** Closed. **Phase:** 7.

## Description
keto-cue has one Markdown document for users, `docs/cue_permission_schema.md`. Users need a public site with a guide, the schema reference and the CLI reference. The site uses the zmake setup: Astro Starlight in `docs/`, bun, Make targets and a GitHub Actions deploy to GitHub Pages.

## Acceptance Criteria
- [x] `docs/` is an Astro Starlight project with the same versions as zmake. `docs/bun.lock` pins the packages.
- [x] A landing page shows a real schema and the real output of `keto namespace validate`. `make docs-demo` writes that output.
- [x] The site has Start Here, Guides, Reference and Project sections. The pages replace `docs/cue_permission_schema.md`.
- [x] The schema reference page shows the embedded base schema from `schema/cue_base.cue`, so the page and the code cannot differ.
- [x] `make docs-dev`, `make docs-build` and `make docs-demo` exist.
- [x] `.github/workflows/docs-release.yml` builds each pull request and deploys `master` to GitHub Pages.
- [x] `make docs-style` checks the site pages. `make format` does not change the site.
- [x] The commands of the Getting Started page work against the built binary.

## Resolution
The site lives in `docs/` with Astro 7, Starlight 0.41 and bun 1.3.9 ([ADR-0008](/decisions/adr-0008-document-keto-cue-with-astro-starlight-on-github-pages--01a1d4da03c6.md)). It has 15 pages and a landing page with the real output of `keto namespace validate`. The base schema page imports `schema/cue_base.cue`. `TestSiteSnippets` compiles each CUE snippet of the site. The Getting Started commands gave the documented results against the built binary. `.github/workflows/docs-release.yml` deploys `master`. The dead OPL typelib publish workflow is gone.

## Notes
- The upstream OPL specification is not part of the site. The migration page links to it.

## Related
- [ADR-0008](/decisions/adr-0008-document-keto-cue-with-astro-starlight-on-github-pages--01a1d4da03c6.md)
