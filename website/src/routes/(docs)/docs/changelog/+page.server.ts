import type { PageServerLoad } from './$types.js';
import { loadChangelog } from '$lib/content/changelog.server.js';

export const prerender = true;

export const load: PageServerLoad = async () => {
    const groups = await loadChangelog();
    return { groups };
};
