<script>
  import { getDomain } from '$lib/utils.js';
  import { onMount } from 'svelte';
  import { fade } from 'svelte/transition';
  import { imageCache } from '$lib/stores/imageCache.js';
  import SummaryDisplay from './SummaryDisplay.svelte';
  import StoryFooter from './StoryFooter.svelte';

  export let story;

  const placeholderUrl = story ? `/api/og?id=${story.id}` : '';
  $: imageUrl = story.ogImage || placeholderUrl;
  $: domain = getDomain(story.url);

  // Summary state managed by StoryFooter
  let isSummaryVisible = false;
  let summary = null;
  let summaryModel = null;
  let error = null;
  let isLoadingSummary = false;

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

<div class="p-4 rounded-lg hover:bg-[var(--color-background-dark-sections)] transition-colors">
  <a href={story.url} target="_blank" rel="noopener noreferrer" class="group flex items-center space-x-4">
    <div class="flex-shrink-0 w-16 h-16 rounded-md overflow-hidden relative">
      {#if imageLoaded}
        <img
          src={imageUrl}
          alt={story.title}
          on:error={handleImageError}
          loading="lazy"
          class="w-full h-full object-cover"
          transition:fade={{ duration: wasInitiallyLoaded ? 0 : 300 }}
        />
      {:else}
        <div class="skeleton-shimmer w-full h-full"></div>
      {/if}
    </div>
    <div class="flex-grow">
      <h3 class="text-lg font-bold text-[var(--color-primary-text)] group-hover:text-[var(--color-primary-accent)] transition-colors mb-1">{story.title}</h3>
      {#if domain}
        <p class="text-sm text-[var(--color-secondary-text)]">{domain}</p>
      {/if}
    </div>
  </a>

  <!-- Summary Display -->
  {#if isSummaryVisible}
    <div class="pl-20 mt-2">
      <SummaryDisplay 
        {story}
        {summary}
        {summaryModel}
        {error}
        {isLoadingSummary}
        textSize="sm"
        padding="2"
        borderWidth="4"
        variant="secondary"
      />
    </div>
  {/if}

  <!-- Story Footer -->
  <StoryFooter 
    {story}
    {imageUrl}
    {isSummaryVisible}
    {isLoadingSummary}
    variant="secondary"
    on:summaryToggle={handleSummaryToggle}
    on:summaryLoaded={handleSummaryLoaded}
    on:summaryError={handleSummaryError}
    on:summaryLoading={handleSummaryLoading}
  />
</div>
