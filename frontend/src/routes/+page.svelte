<script>
  import { onMount } from 'svelte';
  import HeroStory from '../lib/components/HeroStory.svelte';
  import SecondaryStory from '../lib/components/SecondaryStory.svelte';
  import StoryListItem from '../lib/components/StoryListItem.svelte';
  import AddToHomeScreen from '../lib/components/AddToHomeScreen.svelte';
  import EnableNotifications from '../lib/components/EnableNotifications.svelte';
  import { recordVisit } from '$lib/utils.js';
  export let data;

  $: heroStory = data.stories?.[0];
  $: secondaryStories = data.stories?.slice(1, 5);
  $: remainingStories = data.stories?.slice(5);

  onMount(() => {
    recordVisit();
  });
</script>

  <main class="max-w-7xl mx-auto p-4 sm:p-6 lg:p-8 flex-grow">
    {#if data.stories && data.stories.length > 0}
      <div class="grid grid-cols-1 lg:grid-cols-5 gap-10 mb-12 items-start">
        <section class="lg:col-span-3" aria-labelledby="hero-story-title">
          <AddToHomeScreen />
          <EnableNotifications />
          <HeroStory story={heroStory} />
        </section>
        <aside class="lg:col-span-2 bg-[var(--color-background-dark-sections)] p-6 rounded-lg">
          <h2 class="text-2xl font-bold mb-5 text-[var(--color-primary-text)] border-b-4 border-[var(--color-primary-accent)] pb-3">More Top Stories</h2>
          <div class="space-y-5">
            {#each secondaryStories as story (story.id)}
              <SecondaryStory {story} />
            {/each}
          </div>
        </aside>
      </div>
      
      <section aria-labelledby="all-stories-title">
        <h2 id="all-stories-title" class="text-3xl font-bold mb-6 text-[var(--color-primary-text)] border-b-4 border-[var(--color-secondary-accent)] pb-3">All Stories</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-8 mt-8">
          {#each remainingStories as story (story.id)}
            <StoryListItem {story} />
          {/each}
        </div>
      </section>
    {:else}
      <div class="flex justify-center items-center h-96">
        <p class="text-2xl text-[var(--color-secondary-text)]">Loading stories...</p>
      </div>
    {/if}
  </main>