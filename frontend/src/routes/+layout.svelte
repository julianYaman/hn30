<script>
	import '../app.css';
	import { onMount } from 'svelte';
	import { bookmarks } from '$lib/stores/bookmarks';
	import { theme } from '$lib/stores/theme.js';
	import { invalidateAll } from '$app/navigation';
	import { browser } from '$app/environment';
	import { pwaInfo } from 'virtual:pwa-info';
	import Header from '$lib/components/Header.svelte';
	import Footer from '$lib/components/Footer.svelte';
	import Toast from '$lib/components/Toast.svelte';
	import CookieNotice from '$lib/components/CookieNotice.svelte';
	import PullToRefresh from '$lib/components/PullToRefresh.svelte';
	import Offline from '$lib/components/Offline.svelte';
	import { online } from '$lib/stores/networkCheck';

	$: webManifestLink = pwaInfo ? pwaInfo.webManifest.linkTag : '';
	let lastVisibleTime = Date.now();

	async function handleVisibilityChange() {
		if (!browser || document.visibilityState !== 'visible') return;
		if (Date.now() - lastVisibleTime > 5 * 60 * 1000) await invalidateAll();
		lastVisibleTime = Date.now();
	}

	onMount(() => {
		bookmarks.init();
		theme.init();
		document.addEventListener('visibilitychange', handleVisibilityChange);
		window.addEventListener('focus', handleVisibilityChange);
		return () => {
			document.removeEventListener('visibilitychange', handleVisibilityChange);
			window.removeEventListener('focus', handleVisibilityChange);
		};
	});
</script>

<svelte:head>{@html webManifestLink}</svelte:head>

<PullToRefresh />
<CookieNotice />
<Toast />

<div id="ptr-content" class="site-frame">
	<Header />
	{#if !$online}<Offline />{:else}<slot />{/if}
	<Footer />
</div>
