import type { LayoutServerLoad } from './$types.js';
import { loadDocIndex } from '$lib/content/pipeline.server.js';
import { buildNavigation } from '$lib/config/navigation.js';
import { siteConfig } from '$lib/config/site.js';
import { base } from '$app/paths';

export const prerender = true;

export const load: LayoutServerLoad = () => {
    const docs = loadDocIndex();
    const navigation = buildNavigation(docs, base);
    const topLinks = [
        { title: 'Changelog', href: `${base}/docs/changelog`, order: 0 }
    ];
    return { navigation, topLinks, siteConfig };
};
