import fs from 'node:fs';
import path from 'node:path';
import { compile } from 'mdsvex';
import Prism from 'prismjs';
import 'prismjs/components/prism-go.js';
import 'prismjs/components/prism-bash.js';
import 'prismjs/components/prism-json.js';
import type { ChangelogRelease } from '$lib/types/index.js';

export interface ChangelogGroup {
    minor: string;   // e.g. "0.4"
    label: string;   // e.g. "v0.4.x"
    releases: ChangelogRelease[];
}

function changelogPath(): string {
    return path.resolve(process.cwd(), '..', 'CHANGELOG.md');
}

function escapeHtml(text: string): string {
    return text
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;');
}

function highlightCode(code: string, lang: string | undefined): string {
    if (lang && Prism.languages[lang]) {
        const html = Prism.highlight(code, Prism.languages[lang], lang);
        return `<pre class="language-${lang}"><code class="language-${lang}">${html}</code></pre>`;
    }
    return `<pre><code>${escapeHtml(code)}</code></pre>`;
}

function resolveHtmlDirectives(code: string): string {
    return code.replace(/\{@html `([\s\S]*?)`\}/g, (_match, content) => {
        return content.replace(/\\`/g, '`');
    });
}

/**
 * Parse ../CHANGELOG.md into ChangelogGroup[] sorted by minor version
 * descending (latest first). Within each group, releases are sorted by patch
 * version descending.
 */
export async function loadChangelog(): Promise<ChangelogGroup[]> {
    const raw = fs.readFileSync(changelogPath(), 'utf-8');

    // Split on lines starting "## [x.y.z]"
    const parts = raw.split(/^(?=## \[\d)/m).filter((s) =>
        /^## \[\d+\.\d+\.\d+\]/.test(s)
    );

    const releases: ChangelogRelease[] = [];

    for (const part of parts) {
        const match = part.match(/^## \[(\d+\.\d+\.\d+)\] - (\d{4}-\d{2}-\d{2})/);
        if (!match) continue;

        const version = match[1];
        const date = match[2];
        const body = part.replace(/^## \[[\d.]+\] - \d{4}-\d{2}-\d{2}\r?\n/, '').trim();

        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        const compiled = await compile(body, {
            highlight: { highlighter: highlightCode }
        } as any);

        let renderedBody = '';
        if (compiled) {
            let code = compiled.code;
            code = code
                .replace(/<script[\s\S]*?<\/script>/gi, '')
                .replace(/<style[\s\S]*?<\/style>/gi, '');
            renderedBody = resolveHtmlDirectives(code).trim();
        }

        releases.push({ version, date, body, renderedBody });
    }

    // Group by "major.minor"
    const groups = new Map<string, ChangelogRelease[]>();
    for (const release of releases) {
        const parts = release.version.split('.');
        const key = `${parts[0]}.${parts[1]}`;
        const list = groups.get(key) ?? [];
        list.push(release);
        groups.set(key, list);
    }

    // Sort groups descending; sort releases within group by patch descending
    return Array.from(groups.entries())
        .sort(([a], [b]) => {
            const [ma, na] = a.split('.').map(Number);
            const [mb, nb] = b.split('.').map(Number);
            return mb - ma || nb - na;
        })
        .map(([key, rels]) => ({
            minor: key,
            label: `v${key}.x`,
            releases: rels.sort((a, b) => {
                const ap = Number(a.version.split('.')[2]);
                const bp = Number(b.version.split('.')[2]);
                return bp - ap;
            })
        }));
}
