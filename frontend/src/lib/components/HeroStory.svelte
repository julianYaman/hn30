<script>
  import { getDomain } from '$lib/utils.js';
  import { onMount } from 'svelte';
  import { fade } from 'svelte/transition';
  import { imageCache } from '$lib/stores/imageCache.js';
  import SummaryDisplay from './SummaryDisplay.svelte';
  import StoryFooter from './StoryFooter.svelte';

  export let story;

  // Initialize all reactive values at module scope for SSR compatibility
  const placeholderUrl = story ? `https://hn30-og-image.vercel.app/api/og?id=${story.id}` : '';
  $: imageUrl = story ? (story.ogImage || placeholderUrl) : '';
  $: domain = story ? getDomain(story.url) : '';

  // Summary state - initialize with defaults for SSR
  let isSummaryVisible = false;
  let summary = null;
  let summaryModel = null;
  let error = null;
  let isLoadingSummary = false;

  // Image loading state
  const wasInitiallyLoaded = story && imageUrl ? $imageCache.has(imageUrl) : false;
  let imageLoaded = wasInitiallyLoaded;

  // Only run onMount on client side
  let mounted = false;
  
  onMount(() => {
    mounted = true;
    
    if (!story || imageLoaded) return;

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

  function handleImageError(e) {
    if (e.target.src !== placeholderUrl) {
      e.target.src = placeholderUrl;
    }
  }

  function handleSummaryToggle() {
    isSummaryVisible = !isSummaryVisible;
  }

  function handleSummaryLoaded(event) {
    summary = event.detail.summary;
    summaryModel = event.detail.summaryModel;
  }

  function handleSummaryError(event) {
    error = event.detail.error;
  }

  function handleSummaryLoading(event) {
    isLoadingSummary = event.detail.loading;
  }
</script>

{#if story}
<div class="bg-[var(--color-background-card)] rounded-lg shadow-lg hover:shadow-2xl transition-shadow duration-300 overflow-hidden">
  <div class="bg-[var(--color-primary-accent)] text-white text-sm font-bold px-4 py-2 flex items-center">
    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
      <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
    </svg>
    <span>Today's Top Story</span>
  </div>
  <a href={story.url} target="_blank" rel="noopener noreferrer" class="block group">
    <div class="aspect-video overflow-hidden relative">
      {#if imageLoaded}
        <img
          src={imageUrl}
          alt={`Image for ${story.title}`}
          on:error={handleImageError}
          loading="lazy"
          class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300 ease-in-out"
          transition:fade={{ duration: wasInitiallyLoaded ? 0 : 300 }}
        />
      {:else}
        <div class="skeleton-shimmer w-full h-full"></div>
      {/if}
    </div>
    <div class="p-6">
      <h2 class="text-3xl font-extrabold mb-2 leading-tight text-[var(--color-primary-text)] group-hover:text-[var(--color-primary-accent)] transition-colors">{story.title}</h2>
      {#if domain}
        <p class="text-md font-medium text-[var(--color-secondary-text)] mb-4">{domain}</p>
      {/if}
      {#if story.ogDescription}
        <p class="text-base mb-5 text-[var(--color-secondary-text)] max-w-3xl line-clamp-3 leading-relaxed">{story.ogDescription}</p>
      {/if}
    </div>
  </a>

  <!-- Summary Display -->
  {#if isSummaryVisible}
    <div class="px-6 pb-6">
      <SummaryDisplay 
        {story}
        {summary}
        {summaryModel}
        {error}
        {isLoadingSummary}
        textSize="md"
        padding="4"
        borderWidth="4"
        variant="hero"
      />
    </div>
  {/if}

  <!-- Story Footer -->
  <StoryFooter 
    {story}
    {imageUrl}
    {isSummaryVisible}
    {isLoadingSummary}
    variant="hero"
    on:summaryToggle={handleSummaryToggle}
    on:summaryLoaded={handleSummaryLoaded}
    on:summaryError={handleSummaryError}
    on:summaryLoading={handleSummaryLoading}
  />
</div>
{/if}