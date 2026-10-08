<script lang="ts">
	import '../app.css';
	import type { Snippet } from 'svelte';
	let { children }: { children: Snippet } = $props();
</script>

<a
	href="#main-content"
	class="skip-nav"
	style="position:absolute;left:-9999px;top:auto;width:1px;height:1px;overflow:hidden;"
	onfocus={(e) => { (e.currentTarget as HTMLElement).style.cssText = 'position:fixed;top:1rem;left:1rem;z-index:9999;padding:.5rem 1rem;background:var(--color-bg-secondary);border:1px solid var(--color-border);color:var(--color-text-primary);border-radius:4px;'; }}
	onblur={(e) => { (e.currentTarget as HTMLElement).style.cssText = 'position:absolute;left:-9999px;top:auto;width:1px;height:1px;overflow:hidden;'; }}
>Skip to content</a>

<div class="min-h-screen flex flex-col">
	{@render children()}
</div>
