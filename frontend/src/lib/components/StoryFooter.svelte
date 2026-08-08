<script>
	import { timeAgo } from '$lib/utils.js';
	import { bookmarks, toggleBookmark } from '$lib/stores/bookmarks.js';
	import { getSummary } from '$lib/api.js';
	import { createEventDispatcher } from 'svelte';
	export let story;
	export let imageUrl;
	export let isLoadingSummary = false;
	export let hideMeta = false;
	const dispatch = createEventDispatcher();
	let summary = null;
	let summaryModel = null;
	let error = null;
	async function handleSummary() {
		dispatch('summaryToggle');
		if (summary || error) return;
		dispatch('summaryLoading', { loading: true });
		try {
			const result = await getSummary(story.id);
			summary = result.summary;
			summaryModel = result.model;
			dispatch('summaryLoaded', { summary, summaryModel });
		} catch (exception) {
			error = exception.message;
			dispatch('summaryError', { error });
		} finally {
			dispatch('summaryLoading', { loading: false });
		}
	}
</script>

<div class="story-footer">
	{#if !hideMeta}
		<div class="story-footer__meta">
			<span>{story.score} pts</span><span>·</span>
			<a
				href={`https://news.ycombinator.com/item?id=${story.id}`}
				target="_blank"
				rel="noopener noreferrer">{story.descendants} comments</a
			>
			{#if story.time}<span>·</span><span>{timeAgo(story.time)}</span>{/if}
		</div>
	{/if}
	<div class="story-footer__actions">
		{#if !isLoadingSummary}<button class="tldr-button" on:click|stopPropagation={handleSummary}
				>TL;DR</button
			>{/if}
		<button
			class="bookmark-button"
			on:click|stopPropagation={() => toggleBookmark({ ...story, ogImage: imageUrl })}
			aria-label="Bookmark story"
			title="Bookmark story"
		>
			<svg
				xmlns="http://www.w3.org/2000/svg"
				width="15"
				height="15"
				fill={$bookmarks.some((bookmark) => bookmark.id === story.id) ? 'currentColor' : 'none'}
				viewBox="0 0 24 24"
				stroke="currentColor"
				><path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="1.7"
					d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"
				/></svg
			>
		</button>
	</div>
</div>
