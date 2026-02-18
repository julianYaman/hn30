<script>
  import { onMount } from 'svelte';
  import { slide } from 'svelte/transition';
  import {
    shouldShowA2HSPrompt,
    dismissA2HSPrompt,
    isIOS,
    isAndroid,
  } from '$lib/utils.js';

  let visible = false;
  let showIOSInstructions = false;
  let platform = 'unknown'; // 'ios' | 'android' | 'unknown'
  
  // Android's beforeinstallprompt event
  let deferredPrompt = null;

  // Render modal at <body> level so `position: fixed` stays relative to the viewport
  // even when ancestors (e.g. pull-to-refresh container) are transformed.
  function portal(node) {
    const target = typeof document !== 'undefined' ? document.body : null;
    if (!target) return;
    target.appendChild(node);

    return {
      destroy() {
        if (node.parentNode === target) {
          target.removeChild(node);
        }
      }
    };
  }

  onMount(() => {
    // Determine platform
    if (isIOS()) {
      platform = 'ios';
    } else if (isAndroid()) {
      platform = 'android';
    }

    // Listen for beforeinstallprompt (Android/Chrome)
    const handleBeforeInstallPrompt = (e) => {
      e.preventDefault();
      deferredPrompt = e;
    };

    window.addEventListener('beforeinstallprompt', handleBeforeInstallPrompt);

    // Check if we should show the prompt
    visible = shouldShowA2HSPrompt();

    return () => {
      window.removeEventListener('beforeinstallprompt', handleBeforeInstallPrompt);
    };
  });

  function handleDismiss() {
    dismissA2HSPrompt();
    visible = false;
    showIOSInstructions = false;
  }

  async function handleAddToHomeScreen() {
    if (platform === 'ios') {
      // Show iOS instructions modal
      showIOSInstructions = true;
    } else if (deferredPrompt) {
      // Trigger native Android prompt
      deferredPrompt.prompt();
      const { outcome } = await deferredPrompt.userChoice;
      if (outcome === 'accepted') {
        visible = false;
      }
      deferredPrompt = null;
    } else {
      // Fallback: show generic instructions or just dismiss
      showIOSInstructions = true;
    }
  }

  function closeIOSInstructions() {
    showIOSInstructions = false;
  }
</script>

{#if visible}
  <div 
    class="bg-[var(--color-secondary-accent)]/10 border border-[var(--color-secondary-accent)]/30 rounded-lg p-4 mb-4"
    transition:slide={{ duration: 300 }}
  >
    <div class="flex items-start gap-3">
      <!-- Icon -->
      <div class="flex-shrink-0 mt-0.5">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-[var(--color-secondary-accent)]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
        </svg>
      </div>
      
      <!-- Content -->
      <div class="flex-1 min-w-0">
        <h3 class="text-base font-semibold text-[var(--color-primary-text)] mb-1">
          Add hn30 to your home screen
        </h3>
        <p class="text-sm text-[var(--color-secondary-text)] mb-3">
          Get instant access and push notifications for top Hacker News stories.
        </p>
        
        <!-- Buttons -->
        <div class="flex gap-2">
          <button
            on:click={handleAddToHomeScreen}
            class="bg-[var(--color-secondary-accent)] text-white text-sm font-medium px-3 py-2 rounded-md hover:bg-[var(--color-secondary-accent)]/90 transition-colors whitespace-nowrap"
          >
            Add to Home Screen
          </button>
          <button
            on:click={handleDismiss}
            class="text-sm font-medium px-3 py-2 rounded-md text-[var(--color-secondary-text)] hover:bg-[var(--color-border)]/50 transition-colors whitespace-nowrap"
          >
            Maybe Later
          </button>
        </div>
      </div>
      
      <!-- Close button -->
      <button
        on:click={handleDismiss}
        class="flex-shrink-0 text-[var(--color-tertiary-text)] hover:text-[var(--color-secondary-text)] transition-colors"
        aria-label="Dismiss"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
        </svg>
      </button>
    </div>
  </div>
{/if}

<!-- iOS Instructions Modal -->
{#if showIOSInstructions}
  <div use:portal>
    <div 
      class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4"
      on:click={closeIOSInstructions}
      on:keydown={(e) => e.key === 'Escape' && closeIOSInstructions()}
      role="dialog"
      aria-modal="true"
      aria-labelledby="ios-instructions-title"
      tabindex="-1"
    >
      <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
      <div 
        class="bg-[var(--color-background-card)] rounded-xl max-w-sm w-full p-6 shadow-2xl"
        on:click|stopPropagation
        on:keydown|stopPropagation
        role="document"
      >
      <div class="flex justify-between items-start mb-4">
        <h2 id="ios-instructions-title" class="text-lg font-bold text-[var(--color-primary-text)]">
          Add to Home Screen
        </h2>
        <button
          on:click={closeIOSInstructions}
          class="text-[var(--color-tertiary-text)] hover:text-[var(--color-secondary-text)] transition-colors"
          aria-label="Close"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
          </svg>
        </button>
      </div>
      
      {#if platform === 'ios'}
        <div class="space-y-4">
          <div class="flex items-start gap-3">
            <div class="flex-shrink-0 w-7 h-7 rounded-full bg-[var(--color-secondary-accent)] text-white flex items-center justify-center text-sm font-bold">
              1
            </div>
            <div class="flex-1">
              <p class="text-sm text-[var(--color-primary-text)]">
                Tap the <strong>Share</strong> button
              </p>
              <div class="mt-1 flex items-center gap-1 text-[var(--color-secondary-text)]">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
                </svg>
                <span class="text-xs">(at the bottom of Safari)</span>
              </div>
            </div>
          </div>
          
          <div class="flex items-start gap-3">
            <div class="flex-shrink-0 w-7 h-7 rounded-full bg-[var(--color-secondary-accent)] text-white flex items-center justify-center text-sm font-bold">
              2
            </div>
            <p class="text-sm text-[var(--color-primary-text)] pt-1">
              Scroll down and tap <strong>"Add to Home Screen"</strong>
            </p>
          </div>
          
          <div class="flex items-start gap-3">
            <div class="flex-shrink-0 w-7 h-7 rounded-full bg-[var(--color-secondary-accent)] text-white flex items-center justify-center text-sm font-bold">
              3
            </div>
            <p class="text-sm text-[var(--color-primary-text)] pt-1">
              Tap <strong>"Add"</strong> to confirm
            </p>
          </div>
        </div>
      {:else}
        <!-- Generic/Android fallback instructions -->
        <div class="space-y-4">
          <p class="text-sm text-[var(--color-secondary-text)]">
            To add hn30 to your home screen:
          </p>
          <div class="flex items-start gap-3">
            <div class="flex-shrink-0 w-7 h-7 rounded-full bg-[var(--color-secondary-accent)] text-white flex items-center justify-center text-sm font-bold">
              1
            </div>
            <p class="text-sm text-[var(--color-primary-text)] pt-1">
              Tap the <strong>menu</strong> button (three dots)
            </p>
          </div>
          
          <div class="flex items-start gap-3">
            <div class="flex-shrink-0 w-7 h-7 rounded-full bg-[var(--color-secondary-accent)] text-white flex items-center justify-center text-sm font-bold">
              2
            </div>
            <p class="text-sm text-[var(--color-primary-text)] pt-1">
              Tap <strong>"Add to Home Screen"</strong> or <strong>"Install App"</strong>
            </p>
          </div>
        </div>
      {/if}
      
      <div class="mt-6">
        <button
          on:click={closeIOSInstructions}
          class="w-full bg-[var(--color-secondary-accent)] text-white font-medium py-2.5 px-4 rounded-md hover:bg-[var(--color-secondary-accent)]/90 transition-colors"
        >
          Got it
        </button>
      </div>
    </div>
  </div>
  </div>
{/if}
