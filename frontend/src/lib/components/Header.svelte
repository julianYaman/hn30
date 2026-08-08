<script>
	import { page } from '$app/stores';
	import { theme } from '$lib/stores/theme.js';
	import { notifications } from '$lib/stores/notifications.js';
	import SettingsModal from './SettingsModal.svelte';
	import { onMount } from 'svelte';

	let settingsOpen = false;

	function toggleTheme() {
		theme.setTheme($theme === 'light' ? 'dark' : 'light');
	}

	onMount(() => notifications.initialize());
</script>

<header class="site-header">
	<div class="announcement-banner" role="status">
		<p>
			hn30.yamanlabs.com is now <strong>hn30.eu</strong>
			<a href="https://yaman.pro/blog/hn30-eu" target="_blank" rel="noopener noreferrer"
				>Read more →</a
			>
		</p>
	</div>
	<div class="masthead">
		<a href="/" class="brand" aria-label="HN30 front page">
			<span class="brand__mark">HN30</span>
			<span class="brand__tag">Top 30 stories<br />from Hacker News</span>
		</a>
		<div class="edition">
			<span>Updated regularly</span><span>{new Date().toLocaleDateString('en-CA')}</span>
		</div>
		<nav class="utility-nav" aria-label="Utility navigation">
			<button
				class="utility-button theme-toggle"
				on:click={toggleTheme}
				aria-label={$theme === 'light' ? 'Switch to dark mode' : 'Switch to light mode'}
				title={$theme === 'light' ? 'Switch to dark mode' : 'Switch to light mode'}
			>
				{#if $theme === 'light'}
					<svg class="theme-toggle__moon" viewBox="0 0 24 24" aria-hidden="true">
						<path
							d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"
						/>
					</svg>
					<span>Dark</span>
				{:else}
					<svg class="theme-toggle__sun" viewBox="0 0 24 24" aria-hidden="true">
						<circle cx="12" cy="12" r="4" />
						<path
							d="M12 2v2M12 20v2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M2 12h2M20 12h2M4.93 19.07l1.41-1.41M17.66 6.34l1.41-1.41"
						/>
					</svg>
					<span>Light</span>
				{/if}
			</button>
			<button class="utility-button" on:click={() => (settingsOpen = true)}>Settings</button>
		</nav>
	</div>
	<nav class="main-nav" aria-label="Primary navigation">
		<a href="/" aria-current={$page.url.pathname === '/' ? 'page' : undefined}>Front page</a>
		<a href="https://yaman.pro/blog" target="_blank" rel="noopener noreferrer">Blog</a>
		<a href="/bookmarks" aria-current={$page.url.pathname === '/bookmarks' ? 'page' : undefined}
			>Bookmarks</a
		>
	</nav>
</header>

<SettingsModal bind:isOpen={settingsOpen} />
