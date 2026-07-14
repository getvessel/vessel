// @ts-check

import starlight from '@astrojs/starlight';
import { defineConfig } from 'astro/config';

export default defineConfig({
  site: 'https://docs.vessl.dev',
  integrations: [
    starlight({
      title: 'Vessl Docs',
      customCss: ['./src/styles/theme.css'],
      sidebar: [
        {
          label: 'Start here',
          items: [
            { label: 'Vessl Docs', slug: 'index' },
            { label: 'Getting Started', slug: 'getting-started' },
            { label: 'Deploy Your First App', slug: 'tutorial' },
          ],
        },
        {
          label: 'Core concepts',
          items: [
            { label: 'Deployment', slug: 'deployment' },
            { label: 'Projects & Environments', slug: 'projects' },
            { label: 'Serverless Functions', slug: 'serverless' },
          ],
        },
        {
          label: 'Resources',
          items: [
            { label: 'Databases', slug: 'databases' },
            { label: 'Storage', slug: 'storage' },
          ],
        },
        {
          label: 'Operations',
          items: [
            { label: 'Integrations', slug: 'integrations' },
            { label: 'Configuration', slug: 'configuration' },
            { label: 'Administration', slug: 'admin' },
          ],
        },
        {
          label: 'Reference',
          items: [
            { label: 'CLI Reference', slug: 'cli' },
            { label: 'No Lock-In', slug: 'adopt' },
            { label: 'API Reference', slug: 'api' },
          ],
        },
      ],
      components: {
        SiteTitle: "./src/components/docs-site-title.astro",
        ThemeSelect: "./src/components/docs-theme-select.astro",
      },
    }),
  ],
});
