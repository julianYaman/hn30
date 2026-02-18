<script>
  import { onMount } from 'svelte';
  import { slide } from 'svelte/transition';
  import { notifications } from '$lib/stores/notifications.js';
  import {
    isInstalledStandalonePWA,
    hasCookiesAccepted,
    isNotificationsPromptDismissed,
    dismissNotificationsPrompt,
  } from '$lib/utils.js';

  let visible = false;
  let state = 'idle'; // 'idle' | 'loading' | 'success' | 'error'
  let errorMessage = '';

  // Check visibility conditions
  function checkVisibility() {
    if (typeof window === 'undefined') return false;

    // Must be installed as PWA
    if (!isInstalledStandalonePWA()) return false;

    // Must have accepted cookies
    if (!hasCookiesAccepted()) return false;

    // Must not be dismissed
    if (isNotificationsPromptDismissed()) return false;

    // Must not already want/have notifications enabled
    if ($notifications.wantsNotifications) return false;

    // Must not have notifications blocked
    if ($notifications.permission === 'denied') return false;

    return true;
  }

  onMount(() => {
    // Initial check
    visible = checkVisibility();

    // Re-check when notifications store updates
    const unsubscribe = notifications.subscribe(() => {
      if (state === 'idle') {
        visible = checkVisibility();
      }
    });

    return unsubscribe;
  });

  // Also react to store changes
  $: if ($notifications.isInitialized && state === 'idle') {
    visible = checkVisibility();
  }

  function handleDismiss() {
    dismissNotificationsPrompt();
    visible = false;
  }

  async function handleEnableNotifications() {
    state = 'loading';
    errorMessage = '';

    try {
      await notifications.setNotificationsEnabled(true);

      // Check if it succeeded
      if ($notifications.wantsNotifications && !$notifications.error) {
        state = 'success';
        // Hide after showing success message
        setTimeout(() => {
          visible = false;
        }, 2500);
      } else if ($notifications.error) {
        state = 'error';
        errorMessage = $notifications.error;
      } else {
        // Fallback - might have been denied
        state = 'error';
        errorMessage = 'Could not enable notifications. Please try again.';
      }
    } catch (err) {
      state = 'error';
      errorMessage = err?.message || 'Something went wrong. Please try again.';
    }
  }
</script>

{#if visible}
  <div 
    class="bg-[var(--color-secondary-accent)]/10 border border-[var(--color-secondary-accent)]/30 rounded-lg p-4 mb-4"
    transition:slide={{ duration: 300 }}
  >
    {#if state === 'success'}
      <!-- Success state -->
      <div class="flex items-center gap-3">
        <div class="flex-shrink-0">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <p class="text-sm font-medium text-[var(--color-primary-text)]">
          Notifications enabled! You'll be notified when top stories hit the front page.
        </p>
      </div>
    {:else if state === 'error'}
      <!-- Error state -->
      <div class="flex items-start gap-3">
        <div class="flex-shrink-0 mt-0.5">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-medium text-[var(--color-primary-text)] mb-1">
            Could not enable notifications
          </p>
          <p class="text-sm text-[var(--color-secondary-text)] mb-3">
            {errorMessage}
          </p>
          <div class="flex gap-2">
            <button
              on:click={() => { state = 'idle'; }}
              class="bg-[var(--color-secondary-accent)] text-white text-sm font-medium px-3 py-2 rounded-md hover:bg-[var(--color-secondary-accent)]/90 transition-colors whitespace-nowrap"
            >
              Try Again
            </button>
            <button
              on:click={handleDismiss}
              class="text-sm font-medium px-3 py-2 rounded-md text-[var(--color-secondary-text)] hover:bg-[var(--color-border)]/50 transition-colors whitespace-nowrap"
            >
              Maybe Later
            </button>
          </div>
        </div>
      </div>
    {:else}
      <!-- Idle / Loading state -->
      <div class="flex items-start gap-3">
        <!-- Icon -->
        <div class="flex-shrink-0 mt-0.5">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-[var(--color-secondary-accent)]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
          </svg>
        </div>
        
        <!-- Content -->
        <div class="flex-1 min-w-0">
          <h3 class="text-base font-semibold text-[var(--color-primary-text)] mb-1">
            Never miss a top story
          </h3>
          <p class="text-sm text-[var(--color-secondary-text)] mb-3">
            Get notified when stories hit the front page.
          </p>
          
          <!-- Buttons -->
          <div class="flex gap-2">
            <button
              on:click={handleEnableNotifications}
              disabled={state === 'loading'}
              class="bg-[var(--color-secondary-accent)] text-white text-sm font-medium px-3 py-2 rounded-md hover:bg-[var(--color-secondary-accent)]/90 transition-colors whitespace-nowrap disabled:opacity-70 disabled:cursor-not-allowed flex items-center gap-2"
            >
              {#if state === 'loading'}
                <svg class="animate-spin h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Enabling...
              {:else}
                Enable Notifications
              {/if}
            </button>
            <button
              on:click={handleDismiss}
              disabled={state === 'loading'}
              class="text-sm font-medium px-3 py-2 rounded-md text-[var(--color-secondary-text)] hover:bg-[var(--color-border)]/50 transition-colors whitespace-nowrap disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Maybe Later
            </button>
          </div>
        </div>
        
        <!-- Close button -->
        <button
          on:click={handleDismiss}
          disabled={state === 'loading'}
          class="flex-shrink-0 text-[var(--color-tertiary-text)] hover:text-[var(--color-secondary-text)] transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          aria-label="Dismiss"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
          </svg>
        </button>
      </div>
    {/if}
  </div>
{/if}
