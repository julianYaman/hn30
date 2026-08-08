<script>
	import { getDomain, timeAgo } from '$lib/utils.js';
	import { bookmarks } from '$lib/stores/bookmarks.js';

	export let bookmark;

	$: domain = getDomain(bookmark.url);

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

<article class="bookmark-row">
	<div class="bookmark-row__main">
		<a href={bookmark.url} target="_blank" rel="noopener noreferrer" class="story-title"
			>{bookmark.title}</a
		>
		<div class="story-meta">
			{#if domain}<span>{domain}</span>{/if}
			<span>Saved {fmt(bookmark.savedAt)}</span>
			<a
				href={`https://news.ycombinator.com/item?id=${bookmark.id}`}
				target="_blank"
				rel="noopener noreferrer">HN</a
			>
		</div>
	</div>
	<button type="button" class="bookmark-remove" on:click={remove} aria-label="Remove bookmark"
		>Remove</button
	>
</article>
