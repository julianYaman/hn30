<script>
	import { onMount } from 'svelte';
	import BookmarkStoryListItem from '$lib/components/BookmarkStoryListItem.svelte';
	import BookmarkListItem from '$lib/components/BookmarkListItem.svelte';
	import { bookmarks } from '$lib/stores/bookmarks.js';

	const VIEW_STORAGE_KEY = 'hn30_bookmarks_view';

	let isImporting = false;
	let importFileInput;
	let viewMode = 'card';

	$: bookmarkedStories = $bookmarks;

	onMount(() => {
		const savedView = localStorage.getItem(VIEW_STORAGE_KEY);
		if (savedView === 'card' || savedView === 'list') {
			viewMode = savedView;
		}
	});

	function setViewMode(mode) {
		viewMode = mode;
		localStorage.setItem(VIEW_STORAGE_KEY, mode);
	}

	async function handleExport() {
		await bookmarks.exportBookmarks();
	}

	function handleImportClick() {
		importFileInput?.click();
	}

	async function handleFileChange(event) {
		const file = event.target.files[0];
		if (!file) return;
		event.target.value = '';
		try {
			isImporting = true;
			await bookmarks.importBookmarks(file, { skipDuplicates: true });
		} catch (error) {
			console.error('Import failed:', error);
		} finally {
			isImporting = false;
		}
	}
</script>

<main class="page bookmarks-page">
	<div class="stories-section__heading bookmarks-page__heading">
		<h2>Your bookmarks</h2>
		<span>{bookmarkedStories.length} saved</span>
	</div>

	<section class="bookmarks-toolbar" aria-label="Bookmark controls">
		<div class="bookmarks-toolbar__note">
			<p>
				Bookmarks are stored locally in this browser. Clearing site data or using private browsing
				can remove them — export a backup if you want to keep them.
			</p>
		</div>
		<div class="bookmarks-toolbar__actions">
			<div class="view-toggle" role="group" aria-label="View mode">
				<button
					type="button"
					class:is-active={viewMode === 'card'}
					aria-pressed={viewMode === 'card'}
					on:click={() => setViewMode('card')}>Cards</button
				>
				<button
					type="button"
					class:is-active={viewMode === 'list'}
					aria-pressed={viewMode === 'list'}
					on:click={() => setViewMode('list')}>List</button
				>
			</div>
			<div class="bookmarks-io">
				<button
					type="button"
					class="bookmarks-io-button bookmarks-io-button--export"
					on:click={handleExport}
					disabled={bookmarkedStories.length === 0}
				>
					<svg viewBox="0 0 24 24" aria-hidden="true"
						><path
							d="M12 3v12m0 0 4-4m-4 4-4-4M5 21h14"
							fill="none"
							stroke="currentColor"
							stroke-width="1.7"
							stroke-linecap="round"
							stroke-linejoin="round"
						/></svg
					>
					Export ({bookmarkedStories.length})
				</button>
				<button
					type="button"
					class="bookmarks-io-button bookmarks-io-button--import"
					on:click={handleImportClick}
					disabled={isImporting}
				>
					<svg viewBox="0 0 24 24" aria-hidden="true"
						><path
							d="M12 21V9m0 0 4 4m-4-4-4 4M5 3h14"
							fill="none"
							stroke="currentColor"
							stroke-width="1.7"
							stroke-linecap="round"
							stroke-linejoin="round"
						/></svg
					>
					{isImporting ? 'Importing…' : 'Import'}
				</button>
				<input
					bind:this={importFileInput}
					type="file"
					accept=".json,application/json"
					on:change={handleFileChange}
					hidden
				/>
			</div>
		</div>
	</section>

	{#if bookmarkedStories.length > 0}
		{#if viewMode === 'card'}
			<div class="bookmarks-grid">
				{#each bookmarkedStories as bookmark (bookmark.id)}
					<BookmarkStoryListItem {bookmark} />
				{/each}
			</div>
		{:else}
			<div class="bookmarks-list">
				{#each bookmarkedStories as bookmark (bookmark.id)}
					<BookmarkListItem {bookmark} />
				{/each}
			</div>
		{/if}
	{:else}
		<div class="bookmarks-empty">
			<p class="eyebrow">No bookmarks yet</p>
			<p>Save stories from the front page to read them later.</p>
			<a href="/">Back to front page →</a>
		</div>
	{/if}
</main>
