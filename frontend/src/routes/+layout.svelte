<script>
  import "../app.css";
  import { onMount } from 'svelte';
  import { bookmarks } from '$lib/stores/bookmarks';
  import { theme } from '$lib/stores/theme.js';
  import { invalidateAll } from '$app/navigation';
  import { browser } from '$app/environment';
  import Toast from '$lib/components/Toast.svelte';
  import CookieNotice from '$lib/components/CookieNotice.svelte';
  import PullToRefresh from '$lib/components/PullToRefresh.svelte';
  import { pwaInfo } from 'virtual:pwa-info';
  import Header from '../lib/components/Header.svelte';
  import Footer from '../lib/components/Footer.svelte';

  import { online } from '$lib/stores/networkCheck';
  import Offline from '$lib/components/Offline.svelte';

  $: webManifestLink = pwaInfo ? pwaInfo.webManifest.linkTag : ''
  
  // Auto-refresh when app becomes visible after being inactive
  const STALE_THRESHOLD_MS = 5 * 60 * 1000; // 5 minutes
  let lastVisibleTime = Date.now();

  async function handleVisibilityChange() {
    if (!browser) return;
    
    if (document.visibilityState === 'visible') {
      const now = Date.now();
      const timeSinceLastVisible = now - lastVisibleTime;
      
      // If the app was hidden for more than the threshold, refresh content
      if (timeSinceLastVisible > STALE_THRESHOLD_MS) {
        console.log(`App was inactive for ${Math.round(timeSinceLastVisible / 1000)}s. Refreshing content...`);
        try {
          await invalidateAll();
        } catch (err) {
          console.error('Auto-refresh failed:', err);
        }
      }
      
      lastVisibleTime = now;
    }
  }
  
  onMount(() => {
    bookmarks.init();
    theme.init();
    
    // Listen for visibility changes to auto-refresh stale content
    document.addEventListener('visibilitychange', handleVisibilityChange);
    window.addEventListener('focus', handleVisibilityChange);
    
    return () => {
      document.removeEventListener('visibilitychange', handleVisibilityChange);
      window.removeEventListener('focus', handleVisibilityChange);
    };
  });
</script>

<svelte:head>
  {@html webManifestLink}
</svelte:head>

<PullToRefresh />
<CookieNotice />
<Toast />

<div id="ptr-content" class="min-h-screen flex flex-col">
  <Header />
  {#if !$online}
    <Offline />
  {:else}
  <slot />
  {/if}
  <Footer />
</div>

<style>
  /* Ensure the content container can be transformed without layout issues */
  #ptr-content {
    will-change: transform;
    transition: transform 0.1s cubic-bezier(0.23, 1, 0.32, 1);
  }
</style>