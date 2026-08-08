<script>
	import { createEventDispatcher, onMount } from 'svelte';
	export let story;
	export let summary = null;
	export let summaryModel = null;
	export let error = null;
	export let isLoadingSummary = false;
	const dispatch = createEventDispatcher();
	const close = () => dispatch('close');
	function portal(node) {
		if (typeof document === 'undefined') return;
		document.body.appendChild(node);
		return { destroy: () => node.parentNode === document.body && document.body.removeChild(node) };
	}
	onMount(() => {
		const onKeydown = (event) => event.key === 'Escape' && close();
		window.addEventListener('keydown', onKeydown);
		return () => window.removeEventListener('keydown', onKeydown);
	});
</script>

<div use:portal>
	<!-- svelte-ignore a11y-click-events-have-key-events -->
	<div class="summary-backdrop" on:click={close} role="button" tabindex="0"></div>
	<div class="summary-modal" role="dialog" aria-modal="true" aria-labelledby="summary-title">
		<header class="summary-modal__header">
			<div>
				<p class="eyebrow">AI reading note</p>
				<h2 id="summary-title">TL;DR</h2>
			</div>
			<button class="icon-button" on:click={close} aria-label="Close summary">×</button>
		</header>
		<p class="summary-modal__story">{story?.title}</p>
		<div class="summary-modal__body">
			{#if summary}<p class="whitespace-pre-line">{summary}</p>
			{:else if isLoadingSummary}<p class="summary-loading">Preparing summary…</p>
			{:else if error}<p class="summary-error">{error}</p>{/if}
		</div>
		{#if summaryModel}<p class="summary-modal__model">Generated with {summaryModel}</p>{/if}
	</div>
</div>
