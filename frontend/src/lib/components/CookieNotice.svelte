<script>
  import { onMount } from 'svelte';
  import { notifications } from '$lib/stores/notifications.js';

  let showNotice = false;

  onMount(() => {
    if (localStorage.getItem('privacy_accepted') !== 'true') {
      showNotice = true;
    } else {
      // If already accepted, load OneSignal and initialize notifications
      loadOneSignal();
    }
  });

  function accept() {
    localStorage.setItem('privacy_accepted', 'true');
    showNotice = false;
    // Immediately update the store to reflect cookies are now accepted
    notifications.refresh();
    loadOneSignal();
  }

  function loadOneSignal() {
    if (typeof window === 'undefined') return;
    
    // Check if OneSignal is already loaded and initialized
    if (window.OneSignal && window.OneSignal.init) {
      // OneSignal already initialized, just reinitialize the notifications store
      notifications.initialize();
      return;
    }
    
    // Load OneSignal SDK script
    const script = document.createElement('script');
    script.src = 'https://cdn.onesignal.com/sdks/web/v16/OneSignalSDK.page.js';
    script.async = true;
    script.onload = () => {
      console.log('OneSignal SDK loaded');
      // Wait for OneSignal to be ready
      window.OneSignalDeferred = window.OneSignalDeferred || [];
      window.OneSignalDeferred.push(async function(OneSignal) {
        console.log('OneSignalDeferred callback executed');
        // Initialize OneSignal if not already done
        if (!OneSignal.initialized) {
          await OneSignal.init({
            appId: 'f86971bd-cd8f-418a-8c18-906803c91b99',
            allowLocalhostAsSecureOrigin: true,
          });
        }
        // Reinitialize the notifications store after OneSignal is ready
        notifications.initialize();
      });
    };
    document.head.appendChild(script);
  }
</script>

{#if showNotice}
  <div class="fixed bottom-0 left-0 right-0 bg-[var(--color-background-card)] border-t border-gray-200 dark:border-gray-700 p-4 z-50">
    <div class="max-w-7xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-4">
      <p class="text-sm text-[var(--color-secondary-text)]">
        We use analytics to understand our traffic. By using our site, you acknowledge you have read and understood our 
        <a href="/privacy" class="underline text-[var(--color-primary-accent)] hover:text-[var(--color-highlight-cta)]">Privacy Policy</a>.
      </p>
      <button 
        on:click={accept}
        class="bg-[var(--color-primary-accent)] text-white font-bold py-2 px-4 rounded-md hover:bg-opacity-90 transition-colors flex-shrink-0"
      >
        Accept
      </button>
    </div>
  </div>
{/if}
