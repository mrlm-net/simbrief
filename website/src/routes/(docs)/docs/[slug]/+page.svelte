<script lang="ts">
	import { base } from '$app/paths';
	import type { PageData } from './$types.js';
	import type { DocPage } from '$lib/types/index.js';
	import SeoHead from '$lib/components/layout/SeoHead.svelte';
	import JsonLd from '$lib/components/layout/JsonLd.svelte';
	import TableOfContents from '$lib/components/layout/TableOfContents.svelte';

	let { data }: { data: PageData } = $props();

	const doc: DocPage = $derived(data.doc);

	const jsonLdSchema = $derived({
		'@context': 'https://schema.org',
		'@type': 'TechArticle',
		headline: doc.title,
		description: doc.description,
		url: `${data.siteConfig.url}${base}/docs/${doc.slug}`
	});
</script>

<SeoHead
	siteConfig={data.siteConfig}
	title="{doc.title} - simbrief"
	description={doc.description}
	path="/docs/{doc.slug}"
/>

<JsonLd schema={jsonLdSchema} />

<div class="flex justify-center">
	<article class="doc prose min-w-0 flex-1 px-6 py-10 sm:px-10 lg:py-14">
		<p class="eyebrow not-prose mb-3">{data.siteConfig.title} docs</p>
		<h1>{doc.title}</h1>

		{@html doc.renderedContent}

		<nav class="not-prose mt-14 grid gap-3 sm:grid-cols-2" aria-label="Previous and next page">
			<div>
				{#if data.prev}
					<a href="{base}/docs/{data.prev.slug}" class="card block px-4 py-3">
						<span class="eyebrow block">← Previous</span>
						<span class="text-sm font-medium" style="color: var(--text);">{data.prev.title}</span>
					</a>
				{/if}
			</div>
			<div>
				{#if data.next}
					<a href="{base}/docs/{data.next.slug}" class="card block px-4 py-3 text-right">
						<span class="eyebrow block">Next →</span>
						<span class="text-sm font-medium" style="color: var(--text);">{data.next.title}</span>
					</a>
				{/if}
			</div>
		</nav>
	</article>

	{#if doc.headings.length > 0}
		<aside class="sticky top-16 hidden h-[calc(100vh-4rem)] shrink-0 overflow-y-auto xl:block">
			<TableOfContents headings={doc.headings} />
		</aside>
	{/if}
</div>

<style>
	.doc {
		max-width: 52rem;
	}
</style>
