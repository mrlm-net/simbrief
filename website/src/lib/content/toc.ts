import type { TocEntry } from '$lib/types/index.js';

export function slugify(text: string): string {
    return text
        .toLowerCase()
        .replace(/[^\w\s-]/g, '')
        .replace(/\s+/g, '-')
        .replace(/-+/g, '-')
        .trim();
}

export function extractToc(markdown: string): TocEntry[] {
    const entries: TocEntry[] = [];
    const seen = new Map<string, number>();
    const headingRegex = /^(#{2,3})\s+(.+)$/gm;
    let match: RegExpExecArray | null;

    while ((match = headingRegex.exec(markdown)) !== null) {
        const depth = match[1].length;
        const text = match[2].trim();
        const base = slugify(text);
        const count = seen.get(base) ?? 0;
        seen.set(base, count + 1);
        const id = count === 0 ? base : `${base}-${count}`;
        entries.push({ depth, text, id });
    }

    return entries;
}
