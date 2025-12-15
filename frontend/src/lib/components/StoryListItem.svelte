<script>
  import { generatePlaceholder, getDomain, proxyImageUrl } from '$lib/utils.js';
  import { onMount } from 'svelte';
  import { fade } from 'svelte/transition';
  import { imageCache } from '$lib/stores/imageCache.js';
  import SummaryDisplay from './SummaryDisplay.svelte';
  import StoryFooter from './StoryFooter.svelte';

  export let story;

  const placeholderUrl = generatePlaceholder(story.title, 400, 225);
  $: imageUrl = proxyImageUrl(story.ogImage) || placeholderUrl;
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

  <!-- Summary Display -->
  {#if isSummaryVisible}
    <div class="px-5 pb-5">
      <SummaryDisplay 
        {story}
        {summary}
        {summaryModel}
        {error}
        {isLoadingSummary}
        textSize="sm"
        padding="4"
        borderWidth="4"
        variant="list"
      />
    </div>
  {/if}

  <!-- Story Footer -->
  <div class="mt-auto">
    <StoryFooter 
      {story}
      {imageUrl}
      {isSummaryVisible}
      {isLoadingSummary}
      variant="list"
      on:summaryToggle={handleSummaryToggle}
      on:summaryLoaded={handleSummaryLoaded}
      on:summaryError={handleSummaryError}
      on:summaryLoading={handleSummaryLoading}
    />
  </div>
</div>

<style>
  .skeleton-shimmer {
    position: relative;
    background-color: #e2e8f0;
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
