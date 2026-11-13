# keto-cue documentation site

This directory holds the keto-cue documentation site. The site uses Astro
Starlight and bun. GitHub Actions deploys it to
https://infrashift.github.io/keto-cue/.

- `make docs-dev` serves the site with live reload.
- `make docs-build` builds the site into `docs/dist`, the same build as the
  deploy.
- `make docs-demo` builds keto and writes the landing-page demo output.

The pages are in `src/content/docs/docs/`. The landing page is
`src/pages/index.astro`.
