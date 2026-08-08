<script>
	import { getDomain } from '$lib/utils.js';
	import { onMount } from 'svelte';
	import { imageCache } from '$lib/stores/imageCache.js';
	import SummaryDisplay from './SummaryDisplay.svelte';
	import StoryFooter from './StoryFooter.svelte';
	import StoryImagePlaceholder from './StoryImagePlaceholder.svelte';
	export let story;
	$: imageUrl = story?.ogImage || '';
	$: domain = story ? getDomain(story.url) : '';
	let imageLoaded = false;
	let isSummaryVisible = false;
	let summary = null;
	let summaryModel = null;
	let error = null;
	let isLoadingSummary = false;
	onMount(() => {
		if (!story?.ogImage) return;
		if ($imageCache.has(imageUrl)) {
			imageLoaded = true;
			return;
		}
		const image = new Image();
		image.onload = image.onerror = () => {
			imageCache.add(imageUrl);
			imageLoaded = true;
		};
		image.src = imageUrl;
	});
</script>

{#if story}
	<article>
		<a href={story.url} target="_blank" rel="noopener noreferrer" class="block">
			<div class="hero-image">
				{#if story.ogImage && imageLoaded}<img
						src={imageUrl}
						alt={`Image for ${story.title}`}
						loading="eager"
					/>
				{:else if story.ogImage}<div class="skeleton-shimmer"></div>
				{:else}<StoryImagePlaceholder title={story.title} size="hero" />{/if}
			</div>
			<div class="hero-copy">
				<p class="eyebrow">Featured story</p>
				<h1 class="hero-title">{story.title}</h1>
				{#if domain}<p class="domain-badge">{domain}</p>{/if}
				{#if story.ogDescription}<p class="hero-description">{story.ogDescription}</p>{/if}
			</div>
		</a>
		<StoryFooter
			{story}
			{imageUrl}
			{isLoadingSummary}
			on:summaryToggle={() => (isSummaryVisible = !isSummaryVisible)}
			on:summaryLoaded={(event) => {
				summary = event.detail.summary;
				summaryModel = event.detail.summaryModel;
			}}
			on:summaryError={(event) => (error = event.detail.error)}
			on:summaryLoading={(event) => (isLoadingSummary = event.detail.loading)}
		/>
		{#if isSummaryVisible}<SummaryDisplay
				{story}
				{summary}
				{summaryModel}
				{error}
				{isLoadingSummary}
				on:close={() => (isSummaryVisible = false)}
			/>{/if}
	</article>
{/if}
