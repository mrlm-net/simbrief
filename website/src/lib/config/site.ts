import type { SiteConfig } from '$lib/types/index.js';

export const siteConfig: SiteConfig = {
    title: 'simbrief',
    description: 'Go library for SimBrief: fetch and decode OFPs (JSON v2 and XML), take-off and landing speeds from the runway analysis, and dispatch URLs for new plans.',
    repoUrl: 'https://github.com/mrlm-net/simbrief',
    basePath: '',
    url: 'https://simbrief.mrlm.net',
    ogImage: {
        width: 1200,
        height: 630
    },
    locale: 'en_US',
    license: 'Apache-2.0',
    licenseLabel: 'Apache 2.0',
    since: 2025,
    glyph: 'sb',
    nav: [
        { title: 'Docs', href: '/docs/getting-started', match: '/docs' },
        { title: 'Examples', href: '/docs/examples' },
        { title: 'Changelog', href: '/docs/changelog' }
    ]
};
