<script>
  import "../app.css";
  import { onMount } from 'svelte';
  import { bookmarks } from '$lib/stores/bookmarks';
  import { theme } from '$lib/stores/theme.js';
  import Toast from '$lib/components/Toast.svelte';
  import CookieNotice from '$lib/components/CookieNotice.svelte';
  import { pwaInfo } from 'virtual:pwa-info';
  import Header from '../lib/components/Header.svelte';
  import Footer from '../lib/components/Footer.svelte';

  import { online } from '$lib/stores/networkCheck';
  import Offline from '$lib/components/Offline.svelte';

  $: webManifestLink = pwaInfo ? pwaInfo.webManifest.linkTag : ''
  
  onMount(() => {
    bookmarks.init();
    theme.init();
  });
</script>

<svelte:head>
  {@html webManifestLink}
</svelte:head>

<CookieNotice />
<Toast />

<div class="min-h-screen flex flex-col">
  <Header />
  {#if !$online}
    <Offline />
  {:else}
  <slot />
  {/if}
  <Footer />
</div>