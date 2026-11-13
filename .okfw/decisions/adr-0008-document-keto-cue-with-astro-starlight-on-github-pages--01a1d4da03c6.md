---
id: "01a1d4da-03c6-7ab0-a2a9-e1ada141bcaa"
title: "ADR-0008: Document keto-cue with Astro Starlight on GitHub Pages"
type: "ADR"
tags:
  - "adr"
  - "docs"
  - "site"
created: "2026-11-13T06:20:20-06:00"
timestamp: "2026-11-13T06:20:20-06:00"
description: "The keto-cue documentation is an Astro Starlight site in docs/, built with bun and deployed to GitHub Pages by GitHub Actions."
---

# ADR-0008: Document keto-cue with Astro Starlight on GitHub Pages

**Status:** Accepted
**Date:** 2026-11-13
**Source:** User request, Phase 7, [ISSUE-009](/issues/issue-009-documentation-site-on-github-pages--01a1d4da03a0.md)

## Context
Users of keto-cue need a guide and a reference for CUE schemas. The upstream Ory documentation describes OPL and the Ory Network, so it does not apply. The sibling repository zmake publishes an Astro Starlight site to GitHub Pages with bun and a GitHub Actions workflow.

## Decision
- The site lives in `docs/`. It uses Astro 7, Starlight 0.41 and bun 1.3.9, the same versions as zmake. `docs/bun.lock` pins the packages.
- The Go build does not read `docs/`. `make check` does not build the site, so the Go gates do not need bun.
- `src/pages/index.astro` is a landing page outside Starlight. The Starlight pages live under `/docs/`.
- The site pages replace `docs/cue_permission_schema.md`. The upstream OPL specification leaves `docs/`, because keto-cue cannot read OPL.
- The schema reference page imports `schema/cue_base.cue` at build time. The Vite config allows reads from the repository root.
- `make docs-demo` builds the keto binary and runs `keto namespace validate` on `docs/src/demo/`. The landing page shows that output, so the demo is real.
- `.github/workflows/docs-release.yml` builds the site with `withastro/action` and deploys it with `actions/deploy-pages`. A pull request only builds the site. A push to `master` deploys it.
- The Pages source of the repository is "GitHub Actions". The site is at https://infrashift.github.io/keto-cue/ with the Astro `base` `/keto-cue`.
- `.prettierignore` and the copyright step of `make format` exclude `docs/`. Prettier wraps the prose of the pages, and then the upstream format check fails.
- `scripts/ste-check.py` checks the site pages against ASD-STE100.

## Consequences
- Each phase updates the site pages. The Development workflow guide has a step for it.
- The site does not use the Ory logo or the Ory colors. keto-cue is a fork and not an Ory product.

## Alternatives considered
- **A site in `website/` next to the upstream `docs/`.** Rejected: zmake keeps its site in `docs/`, and the upstream `docs/` holds only a pointer to the Ory documentation.
- **GitHub's Jekyll build of `docs/`.** Rejected: no search, no sidebar and a conflict with the Actions deploy.
- **A copy of the base schema in the reference page.** Rejected: the copy can differ from the code.
