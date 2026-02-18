<script>
  import { onMount } from 'svelte';
  import BookmarkStoryListItem from '../../lib/components/BookmarkStoryListItem.svelte';
  import BookmarkListItem from '../../lib/components/BookmarkListItem.svelte';
  import { bookmarks } from '$lib/stores/bookmarks.js';

  const VIEW_STORAGE_KEY = 'hn30_bookmarks_view';

  let bookmarkedStories = [];
  let isImporting = false;
  let importFileInput;
  let viewMode = 'card'; // 'card' | 'list'

  $: bookmarks.subscribe(bm => {
    bookmarkedStories = bm;
  });

  onMount(() => {
    // Load saved view preference
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

    // Reset file input to allow selecting the same file again
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


  <main class="max-w-7xl mx-auto p-4 flex-grow">
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6">
      <h1 class="text-3xl font-bold text-[var(--color-primary-text)] border-b-4 border-[var(--color-secondary-accent)] pb-3">Your Bookmarks</h1>
      
      <!-- View Toggle -->
      <div class="flex items-center gap-1 bg-[var(--color-background-dark-sections)] rounded-lg p-1">
        <button
          on:click={() => setViewMode('card')}
          class="flex items-center gap-2 px-3 py-2 rounded-md text-sm font-medium transition-colors {viewMode === 'card' ? 'bg-[var(--color-background-card)] text-[var(--color-primary-text)] shadow-sm' : 'text-[var(--color-secondary-text)] hover:text-[var(--color-primary-text)]'}"
          aria-label="Card view"
          aria-pressed={viewMode === 'card'}
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
          </svg>
          <span class="hidden sm:inline">Cards</span>
        </button>
        <button
          on:click={() => setViewMode('list')}
          class="flex items-center gap-2 px-3 py-2 rounded-md text-sm font-medium transition-colors {viewMode === 'list' ? 'bg-[var(--color-background-card)] text-[var(--color-primary-text)] shadow-sm' : 'text-[var(--color-secondary-text)] hover:text-[var(--color-primary-text)]'}"
          aria-label="List view"
          aria-pressed={viewMode === 'list'}
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
          <span class="hidden sm:inline">List</span>
        </button>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 xl:gap-8 mt-8">
      <div class="md:col-span-2 bg-[var(--color-background-dark-sections)] border-l-4 border-[var(--color-secondary-accent)] text-[var(--color-secondary-text)] p-4 rounded-md my-2" role="alert">
        <p class="font-bold">Disclaimer</p>
        <p>Your bookmarks are saved locally in your browser. Clearing your browser's cache or using incognito mode may cause them to be lost.</p>
      </div>

      <!-- Import/Export Controls -->
      <div class="md:col-span-2 xl:col-span-1 bg-[var(--color-background-dark-sections)] border border-[var(--color-secondary-accent)] rounded-lg p-4 my-2">
        <h2 class="text-xl font-bold mb-4 text-[var(--color-primary-text)] border-b border-[var(--color-border)] pb-2">Save &amp; Load</h2>
        <div class="flex gap-3">
          <!-- Export Button -->
          <button
            on:click={handleExport}
            class="flex-1 bg-[var(--color-secondary-accent)] hover:bg-[var(--color-secondary-accent)]/90 text-white font-medium py-2 px-4 rounded-lg transition-colors duration-200 flex items-center justify-center gap-2 disabled:opacity-50 cursor-pointer"
            disabled={bookmarkedStories.length === 0}
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
            </svg>
            Export ({bookmarkedStories.length})
          </button>

          <!-- Import Button -->
          <button
            on:click={handleImportClick}
            class="flex-1 bg-[var(--color-primary-accent)] hover:bg-[var(--color-primary-accent)]/90 text-white font-medium py-2 px-4 rounded-lg transition-colors duration-200 flex items-center justify-center gap-2 disabled:opacity-50 cursor-pointer"
            disabled={isImporting}
          >
            {#if isImporting}
              <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"/>
              </svg>
              Importing...
            {:else}
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10"/>
              </svg>
              Import
            {/if}
          </button>
        </div>

        <!-- Hidden file input -->
        <input
          bind:this={importFileInput}
          type="file"
          accept=".json,application/json"
          on:change={handleFileChange}
          class="hidden"
        />
      </div>
    </div>

    {#if bookmarkedStories.length > 0}
      {#if viewMode === 'card'}
        <!-- Card View -->
        <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-8 mt-8">
          {#each bookmarkedStories as bookmark (bookmark.id)}
            <BookmarkStoryListItem bookmark={bookmark} />
          {/each}
        </div>
      {:else}
        <!-- List View -->
        <div class="flex flex-col gap-3 mt-8">
          {#each bookmarkedStories as bookmark (bookmark.id)}
            <BookmarkListItem bookmark={bookmark} />
          {/each}
        </div>
      {/if}
    {:else}
      <div class="flex justify-center items-center h-96">
        <p class="text-xl text-[var(--color-secondary-text)]">You haven't bookmarked any stories yet.</p>
      </div>
    {/if}
  </main>
