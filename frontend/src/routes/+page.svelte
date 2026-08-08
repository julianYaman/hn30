<script>
	import { onMount } from 'svelte';
	import HeroStory from '$lib/components/HeroStory.svelte';
	import SecondaryStory from '$lib/components/SecondaryStory.svelte';
	import StoryListItem from '$lib/components/StoryListItem.svelte';
	import AddToHomeScreen from '$lib/components/AddToHomeScreen.svelte';
	import EnableNotifications from '$lib/components/EnableNotifications.svelte';
	import SidebarNotifications from '$lib/components/SidebarNotifications.svelte';
	import { recordVisit } from '$lib/utils.js';
	export let data;
	$: heroStory = data.stories?.[0];
	$: secondaryStories = data.stories?.slice(1, 10) || [];
	$: remainingStories = data.stories?.slice(10) || [];
	$: firstStoryColumn = remainingStories.slice(0, 10);
	$: secondStoryColumn = remainingStories.slice(10);

	onMount(recordVisit);
</script>

<main class="page">
	{#if data.stories?.length}
		<AddToHomeScreen />
		<EnableNotifications />
		<section class="home-grid" aria-label="Today’s top Hacker News stories">
			<div class="feature-column"><HeroStory story={heroStory} /></div>
			<section class="story-rail" aria-labelledby="story-rail-title">
				<div class="rail-heading">
					<h2 id="story-rail-title">The next nine</h2>
					<span>Ranked by Hacker News</span>
				</div>
				{#each secondaryStories as story, index (story.id)}<div class="rail-item">
						<SecondaryStory {story} rank={index + 2} />
					</div>{/each}
			</section>
			<aside class="sidebar" aria-label="About hn30">
				<section class="sidebar-card">
					<h2>About hn30</h2>
					<p>
						A daily edition of the thirty Hacker News stories worth a closer look — curated for
						builders who want signal over noise.
					</p>
				</section>
				<section class="sidebar-card">
					<h2>Daily dispatch</h2>
					<p>Get the Daily Dispatch — the top 30 stories — in your inbox every morning.</p>
					<a class="sidebar-card__subscribe" href="/subscribe">Subscribe →</a>
				</section>
				<SidebarNotifications />
			</aside>
		</section>
		<section class="stories-section" aria-labelledby="all-stories-title">
			<div class="stories-section__heading">
				<h2 id="all-stories-title">More from today’s edition</h2>
				<span>Stories 11–30</span>
			</div>
			<div class="compact-story-columns">
				<div class="compact-story-column">
					{#each firstStoryColumn as story, index (story.id)}
						<div class="compact-story-item">
							<StoryListItem {story} rank={index + 11} />
						</div>
					{/each}
				</div>
				<div class="compact-story-column">
					{#each secondStoryColumn as story, index (story.id)}
						<div class="compact-story-item">
							<StoryListItem {story} rank={index + 21} />
						</div>
					{/each}
				</div>
			</div>
		</section>
	{:else}
		<div style="min-height: 55vh; display: grid; place-items: center;">
			<p class="eyebrow">Loading today’s edition…</p>
		</div>
	{/if}
</main>
