<script>
  import { timeAgo } from '$lib/utils.js';
  import { bookmarks, toggleBookmark } from '$lib/stores/bookmarks.js';
  import { getSummary } from '$lib/api.js';
  import { createEventDispatcher } from 'svelte';

  export let story;
  export let imageUrl;
  export let isSummaryVisible = false;
  export let isLoadingSummary = false;
  
  // Styling and layout props
  export let variant = 'hero'; // 'hero', 'secondary', 'list'
  export let containerClass = '';
  
  const dispatch = createEventDispatcher();
  
  // Internal state for summary functionality
  let summary = null;
  let summaryModel = null;
  let error = null;
  
  // Computed classes based on variant
  $: textSizeClass = variant === 'secondary' ? 'text-sm' : variant === 'list' ? 'text-xs' : 'text-sm';
  $: metadataClass = variant === 'secondary' ? 'pl-20 mt-2' : variant === 'list' ? 'p-5 pt-0 mt-auto' : 'p-6 pt-0';
  $: buttonGroupClass = variant === 'secondary' ? 'ml-auto flex items-center space-x-2' : 'ml-auto flex items-center space-x-4';
  $: summaryButtonClass = variant === 'secondary' 
    ? 'flex items-center justify-center h-8 w-8 rounded-full text-white transition-transform hover:scale-105'
    : 'flex items-center space-x-2 px-3 py-1 rounded-full text-sm font-semibold text-white transition-transform hover:scale-105';
  $: borderClass = variant === 'list' ? 'border-t border-[var(--color-border)] pt-4' : 'border-t border-[var(--color-border)] pt-4';
  $: footerContainerClass = `flex items-center ${variant === 'list' ? 'justify-between space-x-1' : 'space-x-3'} ${textSizeClass} text-[var(--color-secondary-text)] font-medium ${borderClass} ${metadataClass} ${containerClass}`;

  async function handleSummaryToggle() {
    // If summary is already loaded or there was an error, just toggle visibility
    if (summary || error) {
      dispatch('summaryToggle');
      return;
    }
    
    // First, show the summary area (even while loading)
    dispatch('summaryToggle');
    
    // Then start loading
    dispatch('summaryLoading', { loading: true });
    
    try {
      const result = await getSummary(story.id);
      summary = result.summary;
      summaryModel = result.model;
      dispatch('summaryLoaded', { summary, summaryModel });
    } catch (e) {
      error = e.message;
      dispatch('summaryError', { error });
    } finally {
      dispatch('summaryLoading', { loading: false });
    }
  }

  function handleBookmark() {
    toggleBookmark({ ...story, ogImage: imageUrl });
  }

  // Determine if summary is ready (loaded or error)
  $: isSummaryReady = summary || error;
</script>

<div class={footerContainerClass}>
  <div class="flex items-center space-x-1 lg:space-x-2">
    <span class="font-semibold">{story.score} points</span>
    <span class="text-gray-400">•</span>
    <a href={`https://news.ycombinator.com/item?id=${story.id}`} target="_blank" rel="noopener noreferrer" class="hover:text-[var(--color-primary-accent)] transition-colors flex items-center">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
      </svg>
      <span>{story.descendants}</span>
    </a>
    {#if story.time}
      <span class="text-gray-400">•</span>
      <span>{timeAgo(story.time)}</span>
    {/if}
  </div>

  <!-- Button Group -->
  <div class={buttonGroupClass}>
    {#if !isLoadingSummary}
      <button 
        on:click|stopPropagation={handleSummaryToggle} 
        title={isSummaryVisible ? 'Hide Summary' : 'Show AI Summary'} 
        class={summaryButtonClass}
        style="background-color: var({isSummaryVisible ? '--color-secondary-text' : '--color-secondary-accent'});"
      >
        {#if variant === 'secondary'}
          <span>{isSummaryVisible ? '🙈' : '💡'}</span>
        {:else}
          {#if isSummaryReady && isSummaryVisible}
            <span>🙈</span>
            <span class="text-[var(--color-background-card)]">Hide</span>
          {:else}
            <span>💡</span>
            <span>TL;DR</span>
          {/if}
        {/if}
      </button>
    {/if}

    <button on:click|stopPropagation={handleBookmark} class="text-[var(--color-secondary-text)] hover:text-[var(--color-primary-accent)] transition-colors">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="{$bookmarks.some(b => b.id === story.id) ? 'currentColor' : 'none'}" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
      </svg>
    </button>
  </div>
</div>