// @ts-nocheck
import { visit } from 'unist-util-visit';

function slugify(text) {
    return text
        .toLowerCase()
        .replace(/[^\w\s-]/g, '')
        .replace(/\s+/g, '-')
        .replace(/-+/g, '-')
        .trim();
}

function decodeEntities(str) {
    return str
        .replace(/&amp;/g, '&')
        .replace(/&lt;/g, '<')
        .replace(/&gt;/g, '>')
        .replace(/&quot;/g, '"')
        .replace(/&#39;/g, "'");
}

function extractText(node) {
    if (node.type === 'text') {
        return decodeEntities(node.value ?? '');
    }
    if (node.children) {
        return node.children.map(extractText).join('');
    }
    return '';
}

export default function rehypeSlug() {
    return function transformer(tree) {
        visit(tree, 'element', (node) => {
            if (/^h[1-6]$/.test(node.tagName)) {
                if (!node.properties) {
                    node.properties = {};
                }
                if (!node.properties.id) {
                    const text = extractText(node);
                    node.properties.id = slugify(text);
                }
            }
        });
    };
}
