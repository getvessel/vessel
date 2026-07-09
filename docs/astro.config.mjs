// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

export default defineConfig({
  site: 'https://vessel.dev',
  base: '/docs',
  integrations: [
    starlight({
      title: 'Vessel Docs',
      sidebar: [
        { label: 'Getting Started', slug: 'getting-started' },
        { label: 'Deployment', slug: 'deployment' },
        { label: 'Databases', slug: 'databases' },
        { label: 'Configuration', slug: 'configuration' },
      ],
    }),
  ],
});
