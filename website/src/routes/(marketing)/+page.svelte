<script lang="ts">
	import { base } from '$app/paths';
	import { siteConfig } from '$lib/config/site.js';

	let copied = $state('');

	const installCommand = 'go get github.com/mrlm-net/simbrief';

	async function copy(text: string) {
		try {
			await navigator.clipboard.writeText(text);
			copied = text;
			setTimeout(() => {
				if (copied === text) copied = '';
			}, 2000);
		} catch {
			// clipboard not available
		}
	}

	const facts = [
		{ label: 'OFP formats', value: 'JSON v2 · XML' },
		{ label: 'API key', value: 'none' },
		{ label: 'Packages', value: 'client · types' },
		{ label: 'Dependencies', value: 'stdlib only' },
		{ label: 'Go', value: '1.21+' }
	];

	// Icon paths: 24×24, stroke 1.75
	const features = [
		{
			title: 'Fetch any OFP',
			icon: '<path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><path d="M12 15V3"/>',
			body: "The latest plan of a user by user ID or username, or exactly one plan by its static ID, from SimBrief's public fetcher."
		},
		{
			title: 'JSON v2 and XML',
			icon: '<polyline points="16 18 22 12 16 6"/><polyline points="8 6 2 12 8 18"/>',
			body: 'Both formats decode into the same FlightPlanResponse. Golden tests decode a real OFP in each.'
		},
		{
			title: 'The whole OFP',
			icon: '<rect x="4" y="3" width="16" height="18" rx="2"/><path d="M8 7h8M8 11h8M8 15h5"/>',
			body: 'ATC flight plan, aircraft, fuel, weights, times, weather, alternates with their own navlogs, the navlog fixes, files and links.'
		},
		{
			title: 'Take-off and landing speeds',
			icon: '<path d="M2 22h20"/><path d="M6.36 17.4 4 17l-2-4 1.1-.55a2 2 0 0 1 1.8 0l.17.1a2 2 0 0 0 1.8 0L8 12 5 6l.9-.45a2 2 0 0 1 2.09.2l4.02 3a2 2 0 0 0 2.1.2l4.19-2.06a2.41 2.41 0 0 1 1.73-.17L21 7a1.4 1.4 0 0 1 .87 1.99l-.38.76c-.23.46-.6.84-1.07 1.08L7.58 17.2a2 2 0 0 1-1.22.18Z"/>',
			body: 'V1, VR, V2, flaps, thrust and flex temperature for the planned runway, and the dry or wet landing Vref, from the runway analysis.'
		},
		{
			title: 'Numbers as SimBrief sends them',
			icon: '<path d="M4 9h16M4 15h16M10 3 8 21M16 3l-2 18"/>',
			body: 'types.Number keeps the raw text ("0365", ".57", "") with Int, Float, Bool and IsSet, so "missing" and "0" stay apart.'
		},
		{
			title: 'Dispatch URLs',
			icon: '<path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/>',
			body: 'A fluent FlightPlanBuilder and validation for new plans: open the URL in a browser logged in to SimBrief, then fetch the plan by its static ID.'
		}
	];

	const docs = [
		{ title: 'Getting Started', href: '/docs/getting-started', body: 'Install the library, create a client and fetch your first OFP.' },
		{ title: 'The OFP', href: '/docs/ofp', body: 'Every block of the decoded plan, numbers, the navlog and times.' },
		{ title: 'Runway Analysis', href: '/docs/runway-analysis', body: 'V-speeds, flaps, thrust and Vref from the take-off and landing report.' },
		{ title: 'Generating a Plan', href: '/docs/generating-plans', body: 'Dispatch URLs with the builder, validation, aircraft types and layouts.' }
	];
</script>

<svelte:head>
	<title>simbrief — Go library for SimBrief OFPs</title>
	<meta name="description" content={siteConfig.description} />
</svelte:head>

{#snippet arrow()}
	<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M5 12h14M12 5l7 7-7 7" /></svg>
{/snippet}

{#snippet command(text: string, label: string)}
	<div class="cmd">
		<code><span style="color: var(--text-3);">$&nbsp;</span>{text}</code>
		<button onclick={() => copy(text)} aria-label={label} title="Copy to clipboard" class:ok={copied === text}>
			<svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
				{#if copied === text}
					<polyline points="20 6 9 17 4 12" />
				{:else}
					<rect x="9" y="9" width="13" height="13" rx="2" /><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
				{/if}
			</svg>
		</button>
	</div>
{/snippet}

<!-- Hero -->
<section class="hero px-6 pt-16 pb-16 sm:pt-24 sm:pb-20">
	<div class="mx-auto grid max-w-6xl items-center gap-12 lg:grid-cols-[1.25fr_1fr] [&>*]:min-w-0">
		<div>
			<div class="mb-6 flex flex-wrap gap-2">
				<span class="pill"><span class="dot"></span>Go 1.21+ · Cross-platform · SimBrief API</span>
			</div>
			<h1 class="mb-5 text-4xl font-semibold tracking-tight sm:text-6xl" style="color: var(--text);">
				sim<span class="brand-word">brief</span>
			</h1>
			<p class="mb-8 max-w-xl text-base leading-relaxed sm:text-lg" style="color: var(--text-2);">
				Go library for SimBrief &mdash; fetch and decode OFPs in JSON v2 or XML, read the take-off and
				landing speeds of the runway analysis, and build the dispatch URLs that create new plans.
			</p>
			<div class="mb-8 flex flex-wrap gap-3">
				<a href="{base}/docs/getting-started" class="btn btn-primary">Get started {@render arrow()}</a>
				<a href={siteConfig.repoUrl} target="_blank" rel="noopener noreferrer" class="btn">View on GitHub</a>
			</div>
			<div class="max-w-xl">{@render command(installCommand, 'Copy install command')}</div>
			<div class="mt-4 flex flex-wrap gap-2">
				<a href="https://pkg.go.dev/github.com/mrlm-net/simbrief" target="_blank" rel="noopener noreferrer" class="pill">pkg.go.dev</a>
				<a href="{base}/docs/changelog" class="pill">changelog</a>
				<a href="{siteConfig.repoUrl}/blob/main/LICENSE" target="_blank" rel="noopener noreferrer" class="pill">{siteConfig.licenseLabel}</a>
			</div>
		</div>

		<div class="card overflow-hidden">
			<div class="px-5 py-3" style="border-bottom: 1px solid var(--border);">
				<span class="eyebrow">At a glance</span>
			</div>
			{#each facts as f (f.label)}
				<div class="flex items-center justify-between gap-4 px-5 py-3.5" style="border-bottom: 1px solid var(--border);">
					<span class="text-sm" style="color: var(--text-2);">{f.label}</span>
					<span class="mono text-sm font-medium" style="color: var(--text);">{f.value}</span>
				</div>
			{/each}
			<div class="px-5 py-3 text-xs" style="color: var(--text-3); background: var(--surface-2);">
				fetches from SimBrief's public /api/xml.fetcher.php
			</div>
		</div>
	</div>
</section>

<!-- Features -->
<section class="px-6 py-20" style="border-top: 1px solid var(--border);">
	<div class="mx-auto max-w-6xl">
		<p class="eyebrow mb-3">What it does</p>
		<h2 class="mb-10 max-w-2xl text-2xl font-semibold tracking-tight sm:text-3xl" style="color: var(--text);">
			Your SimBrief plan, typed, in your Go program
		</h2>
		<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
			{#each features as f (f.title)}
				<div class="card p-5">
					<div class="mb-3 flex items-center gap-3">
						<span class="icon">
							<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">{@html f.icon}</svg>
						</span>
						<h3 class="text-[0.95rem] font-semibold" style="color: var(--text);">{f.title}</h3>
					</div>
					<p class="text-sm leading-relaxed" style="color: var(--text-2);">{f.body}</p>
				</div>
			{/each}
		</div>
	</div>
</section>

<!-- Docs -->
<section class="px-6 py-20" style="border-top: 1px solid var(--border);">
	<div class="mx-auto max-w-6xl">
		<p class="eyebrow mb-3">Explore the docs</p>
		<h2 class="mb-10 text-2xl font-semibold tracking-tight sm:text-3xl" style="color: var(--text);">From go get to V-speeds</h2>
		<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
			{#each docs as d (d.href)}
				<a href="{base}{d.href}" class="card group block p-5">
					<div class="mb-2 flex items-center justify-between text-sm font-semibold" style="color: var(--text);">
						{d.title}
						<span class="go" style="color: var(--text-3);">{@render arrow()}</span>
					</div>
					<p class="text-sm leading-relaxed" style="color: var(--text-2);">{d.body}</p>
				</a>
			{/each}
		</div>
	</div>
</section>

<!-- Sponsor -->
<section class="px-6 pb-20">
	<div class="card mx-auto max-w-6xl p-8">
		<p class="eyebrow mb-3">Back open-source MSFS tooling</p>
		<h2 class="mb-3 flex items-center gap-2 text-xl font-semibold tracking-tight" style="color: var(--text);">
			<svg width="18" height="18" viewBox="0 0 24 24" fill="var(--danger)" aria-hidden="true"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z" /></svg>
			Support the project
		</h2>
		<p class="mb-6 max-w-2xl text-sm leading-relaxed" style="color: var(--text-2);">
			Sponsoring covers infrastructure costs, development time, and the simulator licences needed to test
			the mrlm-net flight simulation libraries against the real thing.
		</p>
		<a href="https://revolut.me/mrlm?currency=EUR" target="_blank" rel="noopener noreferrer" class="btn">Sponsor via Revolut</a>
	</div>
</section>

<style>
	.hero {
		background:
			radial-gradient(ellipse 70% 60% at 85% 0%, var(--brand-soft) 0%, transparent 70%),
			var(--bg);
	}
	.icon {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 34px;
		height: 34px;
		flex-shrink: 0;
		border-radius: 8px;
		background: var(--brand-soft);
		color: var(--brand);
	}
	.cmd {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.6rem 0.6rem 0.6rem 1rem;
		border: 1px solid var(--border);
		border-radius: 8px;
		background: var(--surface);
		box-shadow: var(--shadow-1);
	}
	.cmd code {
		flex: 1;
		min-width: 0;
		overflow-x: auto;
		white-space: nowrap;
		font-family: var(--font-mono);
		font-size: 0.8125rem;
		color: var(--text);
		scrollbar-width: none;
	}
	.cmd code::-webkit-scrollbar {
		display: none;
	}
	.cmd button {
		display: inline-flex;
		padding: 0.4rem;
		border-radius: 6px;
		color: var(--text-3);
		cursor: pointer;
	}
	.cmd button:hover {
		color: var(--text);
		background: var(--surface-2);
	}
	.cmd button.ok {
		color: var(--ok);
	}
	.group:hover .go {
		color: var(--text) !important;
	}
</style>
