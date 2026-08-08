<script>
	import { getDomain, timeAgo } from '$lib/utils.js';
	import { onMount } from 'svelte';
	import { imageCache } from '$lib/stores/imageCache.js';
	import SummaryDisplay from './SummaryDisplay.svelte';
	import StoryFooter from './StoryFooter.svelte';
	import StoryImagePlaceholder from './StoryImagePlaceholder.svelte';
	export let story;
	export let rank;
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

<article class="secondary-story">
	<span class="story-rank">{rank}.</span>
	<a
		href={story.url}
		target="_blank"
		rel="noopener noreferrer"
		class="story-thumb"
		aria-label={story.title}
	>
		{#if story.ogImage && imageLoaded}<img src={imageUrl} alt="" loading="lazy" />
		{:else if story.ogImage}<div class="skeleton-shimmer"></div>
		{:else}<StoryImagePlaceholder title={story.title} />{/if}
	</a>
	<div class="secondary-story__content">
		<a href={story.url} target="_blank" rel="noopener noreferrer" class="story-title"
			>{story.title}</a
		>
		<div class="story-meta">
			<span>{domain}</span><span>{story.score} pts</span><a
				href={`https://news.ycombinator.com/item?id=${story.id}`}
				target="_blank"
				rel="noopener noreferrer">{story.descendants} cmt</a
			>{#if story.time}<span>{timeAgo(story.time)}</span>{/if}
		</div>
		<StoryFooter
			{story}
			{imageUrl}
			{isLoadingSummary}
			hideMeta={true}
			on:summaryToggle={() => (isSummaryVisible = !isSummaryVisible)}
			on:summaryLoaded={(event) => {
				summary = event.detail.summary;
				summaryModel = event.detail.summaryModel;
			}}
			on:summaryError={(event) => (error = event.detail.error)}
			on:summaryLoading={(event) => (isLoadingSummary = event.detail.loading)}
		/>
	</div>
	{#if isSummaryVisible}<SummaryDisplay
			{story}
			{summary}
			{summaryModel}
			{error}
			{isLoadingSummary}
			on:close={() => (isSummaryVisible = false)}
		/>{/if}
</article>
