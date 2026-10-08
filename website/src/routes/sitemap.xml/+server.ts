import type { RequestHandler } from './$types.js';
import { loadDocIndex } from '$lib/content/pipeline.server.js';
import { siteConfig } from '$lib/config/site.js';
import fs from 'node:fs';
import path from 'node:path';

export const prerender = true;

function toW3CDate(date: Date): string {
    return date.toISOString().split('T')[0];
}

function docLastmod(slug: string): string {
    const filePath = path.resolve(process.cwd(), '..', 'docs', `${slug}.md`);
    try {
        const stat = fs.statSync(filePath);
        return toW3CDate(stat.mtime);
    } catch {
        return toW3CDate(new Date());
    }
}

export const GET: RequestHandler = () => {
    const baseUrl = `${siteConfig.url}${siteConfig.basePath}`;
    const buildDate = toW3CDate(new Date());

    // Static pages
    const staticPages = [{ path: '/', priority: '1.0', lastmod: buildDate }];

    // Doc pages at /docs/[slug]
    const docs = loadDocIndex();
    const docPages = docs.map((doc) => ({
        path: `/docs/${doc.slug}`,
        priority: '0.8',
        lastmod: docLastmod(doc.slug)
    }));

    const allPages = [...staticPages, ...docPages];
    const urlEntries = allPages
        .map(
            (e) =>
                `  <url>\n    <loc>${baseUrl}${e.path}</loc>\n    <lastmod>${e.lastmod}</lastmod>\n    <priority>${e.priority}</priority>\n  </url>`
        )
        .join('\n');

    const xml = `<?xml version="1.0" encoding="UTF-8"?>\n<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n${urlEntries}\n</urlset>`;

    return new Response(xml, {
        headers: { 'Content-Type': 'application/xml' }
    });
};
