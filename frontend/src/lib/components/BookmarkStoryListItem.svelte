<script>
  import { randomColorrMeBackground, getDomain, timeAgo } from '$lib/utils.js';
  import { bookmarks } from '$lib/stores/bookmarks.js';
  import { onMount } from 'svelte';
  import { fade } from 'svelte/transition';
  import { imageCache } from '$lib/stores/imageCache.js';

  export let bookmark;

  const placeholderUrl = randomColorrMeBackground();
  $: imageUrl = bookmark.ogImage || placeholderUrl;
  $: domain = getDomain(bookmark.url);

  const wasInitiallyLoaded = $imageCache.has(imageUrl);
  let imageLoaded = wasInitiallyLoaded;

  onMount(() => {
    if (imageLoaded) return;

    const img = new Image();
    img.onload = () => {
      imageCache.add(imageUrl);
      imageLoaded = true;
    };
    img.onerror = () => {
      imageCache.add(imageUrl);
      imageLoaded = true;
    };
    img.src = imageUrl;

    if (img.complete) {
      imageCache.add(imageUrl);
      imageLoaded = true;
    }
  });

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

<div class="block bg-[var(--color-background-card)] rounded-lg shadow-md hover:shadow-xl transition-shadow duration-300 overflow-hidden h-full flex flex-col">
  <a href={bookmark.url} target="_blank" rel="noopener noreferrer" class="block group flex-grow">
    <div class="aspect-video overflow-hidden relative">
      {#if imageLoaded}
        <img
          src={imageUrl}
          alt={bookmark.title}
          loading="lazy"
          class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300 ease-in-out"
          transition:fade={{ duration: wasInitiallyLoaded ? 0 : 300 }}
        />
      {:else}
        <div class="skeleton-shimmer w-full h-full"></div>
      {/if}
    </div>
    <div class="p-5">
      <h3 class="text-xl font-bold text-[var(--color-primary-text)] mb-1 group-hover:text-[var(--color-primary-accent)] transition-colors">{bookmark.title}</h3>
      {#if domain}
        <p class="text-sm text-[var(--color-secondary-text)] mb-2">{domain}</p>
      {/if}
      {#if bookmark.ogDescription}
        <p class="text-sm text-[var(--color-secondary-text)] mb-4 line-clamp-3">{bookmark.ogDescription}</p>
      {/if}
    </div>
  </a>

  <div class="p-5 pt-0 mt-auto">
    <div class="text-xs text-[var(--color-secondary-text)] pt-4 border-t border-[var(--color-border)] flex items-center space-x-2">
      <div class="flex flex-col">
        <span class="font-semibold">Saved</span>
        <span class="text-[var(--color-secondary-text)]">{fmt(bookmark.savedAt)}</span>
      </div>

      <div class="mx-2">•</div>

      <div class="flex flex-col">
        <span class="font-semibold">Posted</span>
        <span class="text-[var(--color-secondary-text)]">{fmt(bookmark.postedAt)}</span>
      </div>

      <div class="mx-2">•</div>

      <a href={`https://news.ycombinator.com/item?id=${bookmark.id}`} target="_blank" rel="noopener noreferrer" class="text-[var(--color-primary-accent)] hover:underline">
        View on Hacker News
      </a>

      <button
        on:click|stopPropagation={remove}
        aria-label="Remove bookmark"
        title="Remove bookmark"
        class="ml-auto bg-[var(--color-primary-accent)] hover:bg-[var(--color-secondary-accent)] text-white rounded-md px-3 py-1 transition-colors"
      >
        ×
      </button>
    </div>
  </div>
</div>
