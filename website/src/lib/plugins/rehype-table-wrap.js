// @ts-nocheck
import { visit } from 'unist-util-visit';

/**
 * Rehype plugin that wraps <table> elements in a scrollable <div>.
 */
export default function rehypeTableWrap() {
    return function transformer(tree) {
        visit(tree, 'element', (node, index, parent) => {
            if (node.tagName === 'table' && parent && index !== null) {
                const wrapper = {
                    type: 'element',
                    tagName: 'div',
                    properties: { className: ['table-wrapper'] },
                    children: [node]
                };
                parent.children[index] = wrapper;
            }
        });
    };
}
