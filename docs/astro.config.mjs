// @ts-check
import { fileURLToPath } from 'node:url';
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// Deployed to GitHub Pages by .github/workflows/docs-release.yml (ADR-0008).
// Each PLAN.md phase adds or updates the pages for the features it delivers
// (see the Development workflow guide in .okfw/).
// https://astro.build/config

// The schema reference imports schema/cue_base.cue from the repository root.
const repoRoot = fileURLToPath(new URL('..', import.meta.url));

export default defineConfig({
	site: 'https://infrashift.github.io',
	base: '/keto-cue',
	vite: {
		server: { fs: { allow: [repoRoot] } },
	},
	integrations: [
		starlight({
			title: 'keto-cue',
			description:
				'A fork of Ory Keto that reads Zanzibar permission schemas in CUE: short rewrites, comprehensions and errors with file, line and column.',
			favicon: '/favicon.svg',
			customCss: ['./src/styles/custom.css'],
			social: [
				{
					icon: 'github',
					label: 'GitHub',
					href: 'https://github.com/infrashift/keto-cue',
				},
			],
			editLink: {
				baseUrl: 'https://github.com/infrashift/keto-cue/edit/master/docs/',
			},
			sidebar: [
				{
					label: 'Start Here',
					items: [
						{ label: 'Overview', slug: 'docs' },
						{ label: 'Getting Started', slug: 'docs/getting-started' },
					],
				},
				{
					label: 'Guides',
					items: [
						{ label: 'Permission Schemas', slug: 'docs/guides/permission-schemas' },
						{ label: 'Rewrites', slug: 'docs/guides/rewrites' },
						{ label: 'Reuse & Comprehensions', slug: 'docs/guides/reuse' },
						{ label: 'Loading Schemas', slug: 'docs/guides/loading-schemas' },
						{ label: 'Migrating from OPL', slug: 'docs/guides/migrating-from-opl' },
						{ label: 'Syntax-Check API', slug: 'docs/guides/syntax-check-api' },
					],
				},
				{
					label: 'Reference',
					items: [
						{ label: 'Base Schema', slug: 'docs/reference/base-schema' },
						{ label: 'Rules & Errors', slug: 'docs/reference/rules-and-errors' },
						{ label: 'CLI', slug: 'docs/reference/cli' },
						{ label: 'Make Targets', slug: 'docs/reference/make-targets' },
					],
				},
				{
					label: 'Project',
					collapsed: true,
					items: [
						{ label: 'Roadmap', slug: 'docs/project/roadmap' },
						{ label: 'Design Decisions', slug: 'docs/project/decisions' },
						{ label: 'Contributing', slug: 'docs/project/contributing' },
					],
				},
			],
		}),
	],
});
