// @ts-nocheck
import { visit } from 'unist-util-visit';

/**
 * Rehype plugin that rewrites relative .md links to /docs/slug format
 * and converts ../examples/* links to the GitHub repo URL.
 */
export default function rehypeRewriteLinks() {
    const repoUrl = 'https://github.com/mrlm-net/simbrief/tree/main';

    return function transformer(tree) {
        visit(tree, 'element', (node) => {
            if (node.tagName !== 'a') return;
            const href = node.properties?.href;
            if (!href || typeof href !== 'string') return;

            if (href.startsWith('http://') || href.startsWith('https://') || href.startsWith('#')) {
                return;
            }

            if (href.startsWith('../examples')) {
                const examplePath = href.replace(/^\.\.\//, '');
                node.properties.href = `${repoUrl}/${examplePath}`;
                node.properties.target = '_blank';
                node.properties.rel = 'noopener noreferrer';
                return;
            }

            const mdMatch = href.match(/^([a-zA-Z0-9_-]+)\.md(#.*)?$/);
            if (mdMatch) {
                const slug = mdMatch[1];
                const anchor = mdMatch[2] ?? '';
                node.properties.href = `/docs/${slug}${anchor}`;
            }
        });
    };
}
