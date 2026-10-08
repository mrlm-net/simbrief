/**
 * Navigation item within a sidebar section.
 */
export interface NavItem {
    title: string;
    href: string;
    order: number;
}

/**
 * Top-level navigation section grouping one or more NavItems.
 */
export interface NavSection {
    title: string;
    id: string;
    items: NavItem[];
    defaultOpen?: boolean;
}

/**
 * Single entry in a page's table of contents, derived from heading elements.
 */
export interface TocEntry {
    depth: number;
    text: string;
    id: string;
}

/**
 * Global site configuration used across layout and SEO components.
 *
 * Adapted from mrlm-net/simconnect: ogImage is a structured object rather
 * than a flat string with separate ogImageWidth/ogImageHeight fields.
 */
export interface SiteConfig {
    title: string;
    description: string;
    repoUrl: string;
    basePath: string;
    url: string;
    ogImage: {
        width: number;
        height: number;
    };
    locale: string;
    license: string;
    /** Short label for the license pill, e.g. "BSL 1.1 · non-commercial". */
    licenseLabel: string;
    /** First year of the copyright line. */
    since: number;
    /** One or two characters in the brand mark. */
    glyph: string;
    /** Header links; match is the path prefix that marks a link active (default href). */
    nav: { title: string; href: string; match?: string }[];
}

/**
 * Metadata for a documentation page, used in navigation and page listings.
 */
export interface DocMeta {
    slug: string;
    title: string;
    description: string;
    order: number;
    section: string;
}

/**
 * Full documentation page, extending DocMeta with rendered content and
 * extracted headings for the TableOfContents component.
 */
export interface DocPage extends DocMeta {
    /** Compiled HTML output from the mdsvex pipeline. */
    renderedContent: string;
    /** Extracted h2 headings used to populate the TableOfContents. */
    headings: TocEntry[];
}

/**
 * A single changelog entry.
 *
 * Simplified from mrlm-net/simconnect: our changelog is hand-authored
 * markdown rather than fetched from the GitHub Releases API, so GitHub-
 * specific fields (tag, name, url, prerelease) are omitted.
 */
export interface ChangelogRelease {
    version: string;
    date: string;
    body: string;
    renderedBody: string;
}
