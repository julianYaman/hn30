<script>
  import { notifications } from '$lib/stores/notifications.js';
  import { fly } from 'svelte/transition';

  export let isOpen = false;

  const handleToggleNotifications = async (e) => {
    const desiredState = e.target.checked;

    if (!$notifications.cookiesAccepted) {
      e.target.checked = false;
      return;
    }

    await notifications.setNotificationsEnabled(desiredState);
  };

  const closeModal = () => {
    isOpen = false;
    notifications.clearMessage();
  };
</script>

{#if isOpen}
  <!-- Backdrop -->
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div transition:fly={{ duration: 200, opacity: 0 }} on:click={closeModal} class="fixed inset-0 bg-black/40 backdrop-blur-sm z-40" role="button" tabindex="0"></div>

  <!-- Modal -->
  <div transition:fly={{ duration: 300, y: 50, opacity: 0 }} class="fixed top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 bg-[var(--color-background-card)] rounded-lg shadow-xl z-50 w-96 max-w-full p-6 mx-4">
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-bold text-[var(--color-primary-text)]">Settings</h2>
      <button on:click={closeModal} aria-label="Close settings" class="text-[var(--color-secondary-text)] hover:text-[var(--color-primary-text)] transition-colors">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
        </svg>
      </button>
    </div>

    <div class="space-y-4">
      <!-- Notifications Preference -->
      <div class="flex items-center justify-between p-4 rounded-lg bg-[var(--color-background-dark-sections)]">
        <div class="flex-1">
          <label for="notifications-toggle" class="text-[var(--color-primary-text)] font-medium cursor-pointer {$notifications.permission === 'denied' || (!$notifications.cookiesAccepted && !$notifications.wantsNotifications) ? 'opacity-50' : ''}">
            Receive Notifications
          </label>
          <p class="text-sm text-[var(--color-secondary-text)] mt-1">Get notified about top stories</p>
        </div>
        <div class="ml-4">
          <input
            id="notifications-toggle"
            type="checkbox"
            checked={$notifications.wantsNotifications}
            on:change={handleToggleNotifications}
            disabled={$notifications.loading || $notifications.permission === 'denied' || !$notifications.cookiesAccepted}
            class="w-5 h-5 rounded cursor-pointer accent-[var(--color-primary-accent)] disabled:opacity-50 disabled:cursor-not-allowed"
          />
        </div>
      </div>

      <!-- Warnings -->
      {#if !$notifications.cookiesAccepted}
        <div class="p-3 rounded-lg bg-yellow-100 dark:bg-yellow-900 text-yellow-800 dark:text-yellow-200 text-sm flex items-start">
          <svg class="h-4 w-4 mr-2 mt-0.5 flex-shrink-0" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
          <div>
            <p class="font-medium">Cookies required</p>
            <p class="text-xs mt-1">Please accept cookies to enable notifications.</p>
          </div>
        </div>
      {:else if $notifications.permission === 'denied'}
        <div class="p-3 rounded-lg bg-yellow-100 dark:bg-yellow-900 text-yellow-800 dark:text-yellow-200 text-sm flex items-start">
          <svg class="h-4 w-4 mr-2 mt-0.5 flex-shrink-0" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
          <div>
            <p class="font-medium">Notifications are blocked</p>
            <p class="text-xs mt-1">Please enable notifications in your browser settings to receive alerts.</p>
          </div>
        </div>
      {:else if $notifications.error}
        <div class="p-3 rounded-lg bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200 text-sm">
          {$notifications.error}
        </div>
      {/if}

      {#if $notifications.loading}
        <div class="p-3 rounded-lg bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 text-sm flex items-center">
          <svg class="animate-spin h-4 w-4 mr-2" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          {$notifications.loadingMessage || 'Updating preferences...'}
        </div>
      {/if}

      {#if $notifications.message}
        <div class={`p-3 rounded-lg text-sm flex items-center ${$notifications.messageType === 'success' ? 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200' : 'bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200'}`}>
          <svg class="h-4 w-4 mr-2" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={$notifications.messageType === 'success' ? 'M5 13l4 4L19 7' : 'M6 18L18 6M6 6l12 12'} />
          </svg>
          {$notifications.message}
        </div>
      {/if}
    </div>

    <div class="mt-6 flex justify-end">
      <button on:click={closeModal} class="px-4 py-2 bg-[var(--color-primary-accent)] text-white rounded-lg hover:opacity-90 transition-opacity font-medium">
        Done
      </button>
    </div>
  </div>
{/if}

<style>
  /* Optional: Add custom styling if needed */
</style>
