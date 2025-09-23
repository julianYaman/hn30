<script>
  import { generatePlaceholder, getDomain, timeAgo, proxyImageUrl } from '$lib/utils.js';
  import { bookmarks, toggleBookmark } from '$lib/stores/bookmarks.js';
  import { getSummary } from '$lib/api.js';
  import { onMount } from 'svelte';
  import { fade } from 'svelte/transition';
  import { imageCache } from '$lib/stores/imageCache.js';

  export let story;

  const placeholderUrl = generatePlaceholder(story.title, 400, 225);
  $: imageUrl = proxyImageUrl(story.ogImage) || placeholderUrl;
  $: domain = getDomain(story.url);

  let summary = null;
  let isSummaryVisible = false;
  let isLoadingSummary = false;
  let error = null;

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
      imageCache.add(imageUrl); // Cache errors too to prevent re-fetches
      imageLoaded = true;
    };
    img.src = imageUrl;

    if (img.complete) {
      imageCache.add(imageUrl);
      imageLoaded = true;
    }
  });

  function handleImageError(e) {
    if (e.target.src !== placeholderUrl) {
      e.target.src = placeholderUrl;
    }
  }

  async function handleSummaryToggle() {
    if (summary || error) {
      isSummaryVisible = !isSummaryVisible;
      return;
    }

    isLoadingSummary = true;
    error = null;
    try {
      const result = await getSummary(story.id);
      summary = result.summary;
      isSummaryVisible = true;
    } catch (e) {
      error = e.message;
      isSummaryVisible = true;
    } finally {
      isLoadingSummary = false;
    }
  }
</script>

<div class="block bg-[var(--color-background-card)] rounded-lg shadow-md hover:shadow-xl transition-shadow duration-300 overflow-hidden h-full flex flex-col">
  <a href={story.url} target="_blank" rel="noopener noreferrer" class="block group flex-grow">
    <div class="aspect-video overflow-hidden relative">
      {#if imageLoaded}
        <img
          src={imageUrl}
          alt={story.title}
          on:error={handleImageError}
          loading="lazy"
          class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300 ease-in-out"
          transition:fade={{ duration: wasInitiallyLoaded ? 0 : 300 }}
        />
      {:else}
        <div class="skeleton-shimmer w-full h-full"></div>
      {/if}
    </div>
    <div class="p-5">
  <h3 class="text-xl font-bold text-[var(--color-primary-text)] mb-1 group-hover:text-[var(--color-primary-accent)] transition-colors">{story.title}</h3>
      {#if domain}
          <p class="text-sm text-[var(--color-secondary-text)] mb-2">{domain}</p>
      {/if}
      {#if story.ogDescription}
        <p class="text-sm text-[var(--color-secondary-text)] mb-4 line-clamp-3 leading-relaxed">{story.ogDescription}</p>
      {/if}
    </div>
  </a>

  <!-- Summary Section -->
  <div class="px-5 pb-5">
    {#if isSummaryVisible && summary}
      <div class="text-sm p-3 bg-[var(--color-background-dark-sections)] rounded-md border-l-4 border-[var(--color-secondary-accent)]">
        <p class="font-semibold text-md mb-1">TL;DR</p>
        <p class="text-[var(--color-secondary-text)] whitespace-pre-line">{summary}</p>
      </div>
    {:else if isLoadingSummary}
      <div class="text-sm p-3 bg-[var(--color-background-dark-sections)] rounded-md animate-pulse">Loading summary...</div>
    {:else if isSummaryVisible && error}
      <div class="text-sm p-3 bg-red-100 text-red-700 rounded-md">{error}</div>
    {/if}
  </div>

  <div class="p-5 pt-0 mt-auto">
    <div class="text-xs text-[var(--color-secondary-text)] pt-4 border-t border-[var(--color-border)] flex items-center space-x-2">
      <span class="font-semibold">{story.score} points</span>
      <span class="text-gray-400 mx-1">•</span>
  <a href={`https://news.ycombinator.com/item?id=${story.id}`} target="_blank" rel="noopener noreferrer" class="hover:text-[var(--color-primary-accent)] transition-colors flex items-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
        </svg>
        <span>{story.descendants}</span>
      </a>
      {#if story.time}
        <span class="text-gray-400 mx-1">•</span>
        <span>{timeAgo(story.time)}</span>
      {/if}
      
      <!-- Button Group -->
      <div class="ml-auto flex items-center space-x-2">
        <!-- Summarize/Hide Button -->
        {#if !isLoadingSummary}
        <button on:click|stopPropagation={handleSummaryToggle} title={isSummaryVisible ? 'Hide Summary' : 'Generate AI Summary'} class="flex items-center space-x-2 px-3 py-1 rounded-full text-sm font-semibold text-white transition-transform hover:scale-105" style="background-color: var({isSummaryVisible ? '--color-secondary-text' : '--color-secondary-accent'});">
          {#if isSummaryVisible}
            <span>🙈</span>
            <span>Hide</span>
          {:else}
            <span>💡</span>
            <span>TL;DR</span>
          {/if}
        </button>
        {/if}

        <!-- Bookmark Button -->
        <button on:click|stopPropagation={() => toggleBookmark({ ...story, ogImage: imageUrl })} class="text-[var(--color-secondary-text)] hover:text-[var(--color-primary-accent)] transition-colors">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="{$bookmarks.some(b => b.id === story.id) ? 'currentColor' : 'none'}" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</div>

<style>
  .skeleton-shimmer {
    position: relative;
    background-color: #e2e8f0; /* slate-200 */
    overflow: hidden;
  }

  .skeleton-shimmer::after {
    content: '';
    position: absolute;
    top: 0;
    right: 0;
    bottom: 0;
    left: 0;
    transform: translateX(-100%);
    background-image: linear-gradient(
      90deg,
      rgba(255, 255, 255, 0) 0,
      rgba(255, 255, 255, 0.2) 20%,
      rgba(255, 255, 255, 0.5) 60%,
      rgba(255, 255, 255, 0)
    );
    animation: shimmer 2s infinite;
  }

  @keyframes shimmer {
    100% {
      transform: translateX(100%);
    }
  }
</style>
