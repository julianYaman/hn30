<script>
  import Header from '../../lib/components/Header.svelte';
  import Footer from '../../lib/components/Footer.svelte';
  import BookmarkStoryListItem from '../../lib/components/BookmarkStoryListItem.svelte';
  import { bookmarks } from '$lib/stores/bookmarks.js';

  let bookmarkedStories = [];
  let isImporting = false;
  let importFileInput;
  $: bookmarks.subscribe(bm => {
    bookmarkedStories = bm;
  });

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

<div class="min-h-screen flex flex-col">
  <Header />

  <main class="max-w-7xl mx-auto p-4 flex-grow">
    <h1 class="text-3xl font-bold mb-6 text-[var(--color-primary-text)] border-b-4 border-[var(--color-secondary-accent)] pb-3">Your Bookmarks</h1>

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
      <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-8 mt-8">
        {#each bookmarkedStories as bookmark (bookmark.id)}
          <BookmarkStoryListItem bookmark={bookmark} />
        {/each}
      </div>
    {:else}
      <div class="flex justify-center items-center h-96">
        <p class="text-xl text-[var(--color-secondary-text)]">You haven't bookmarked any stories yet.</p>
      </div>
    {/if}
  </main>

  <Footer />
</div>
