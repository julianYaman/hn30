<script>
  import { onMount } from 'svelte';

  let showNotice = false;

  onMount(() => {
    if (localStorage.getItem('privacy_accepted') !== 'true') {
      showNotice = true;
    }
  });

  function accept() {
    localStorage.setItem('privacy_accepted', 'true');
    showNotice = false;
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
