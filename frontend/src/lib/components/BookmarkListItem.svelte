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
    if (!iso) return 'Unknown';
    const d = new Date(iso);
    const ms = Date.now() - d.getTime();
    const fiveDays = 5 * 24 * 60 * 60 * 1000;
    if (ms >= fiveDays) {
      return d.toLocaleDateString();
    }
    return timeAgo(Math.floor(d.getTime() / 1000));
  }
</script>

<div class="bg-[var(--color-background-card)] rounded-lg shadow-sm hover:shadow-md transition-shadow duration-200 p-4 flex items-start gap-4">
  <!-- Content -->
  <div class="flex-1 min-w-0">
    <a 
      href={bookmark.url} 
      target="_blank" 
      rel="noopener noreferrer" 
      class="group"
    >
      <h3 class="text-base font-semibold text-[var(--color-primary-text)] group-hover:text-[var(--color-primary-accent)] transition-colors line-clamp-2">
        {bookmark.title}
      </h3>
    </a>
    
    <div class="flex flex-wrap items-center gap-x-3 gap-y-1 mt-1 text-sm text-[var(--color-secondary-text)]">
      {#if domain}
        <span class="text-[var(--color-tertiary-text)]">{domain}</span>
      {/if}
      <span class="hidden sm:inline">•</span>
      <span class="hidden sm:inline">Saved {fmt(bookmark.savedAt)}</span>
      <span class="hidden sm:inline">•</span>
      <a 
        href={`https://news.ycombinator.com/item?id=${bookmark.id}`} 
        target="_blank" 
        rel="noopener noreferrer" 
        class="text-[var(--color-primary-accent)] hover:underline"
      >
        HN
      </a>
    </div>
  </div>
  
  <!-- Remove button -->
  <button
    on:click|stopPropagation={remove}
    aria-label="Remove bookmark"
    title="Remove bookmark"
    class="flex-shrink-0 text-[var(--color-tertiary-text)] hover:text-[var(--color-primary-accent)] transition-colors p-1"
  >
    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
      <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
    </svg>
  </button>
</div>
