<script>
	import { getDomain, timeAgo } from '$lib/utils.js';
	import { bookmarks } from '$lib/stores/bookmarks.js';
	import { onMount } from 'svelte';
	import { imageCache } from '$lib/stores/imageCache.js';

	export let bookmark;

	$: imageUrl = bookmark.ogImage || '';
	$: domain = getDomain(bookmark.url);
	let imageLoaded = false;

	onMount(() => {
		if (!imageUrl) return;
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

	function remove(e) {
		e?.stopPropagation?.();
		bookmarks.remove(bookmark.id);
	}

	function fmt(iso) {
		if (!iso) return '—';
		const d = new Date(iso);
		const ms = Date.now() - d.getTime();
		const fiveDays = 5 * 24 * 60 * 60 * 1000;
		if (ms >= fiveDays) return d.toLocaleDateString();
		return timeAgo(Math.floor(d.getTime() / 1000));
	}
</script>

<article class="bookmark-card">
	<a
		href={bookmark.url}
		target="_blank"
		rel="noopener noreferrer"
		class="bookmark-card__media"
		aria-label={bookmark.title}
	>
		{#if imageUrl && imageLoaded}<img src={imageUrl} alt="" loading="lazy" />
		{:else if imageUrl}<div class="skeleton-shimmer"></div>
		{:else}<div class="image-placeholder"><span>HN30</span></div>{/if}
	</a>
	<div class="bookmark-card__body">
		<a href={bookmark.url} target="_blank" rel="noopener noreferrer" class="story-title"
			>{bookmark.title}</a
		>
		{#if bookmark.ogDescription}<p class="bookmark-card__description">{bookmark.ogDescription}</p
			>{/if}
		<div class="story-meta">
			{#if domain}<span>{domain}</span>{/if}
			<span>Saved {fmt(bookmark.savedAt)}</span>
			{#if bookmark.postedAt}<span>Posted {fmt(bookmark.postedAt)}</span>{/if}
		</div>
		<div class="bookmark-card__actions">
			<a
				href={`https://news.ycombinator.com/item?id=${bookmark.id}`}
				target="_blank"
				rel="noopener noreferrer"
				class="tldr-button">Discuss</a
			>
			<button type="button" class="bookmark-remove" on:click={remove} aria-label="Remove bookmark"
				>Remove</button
			>
		</div>
	</div>
</article>
