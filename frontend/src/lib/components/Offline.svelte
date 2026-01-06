<script>
  
  let isRetrying = false;
  
  async function handleRetry() {
    isRetrying = true;
    
    // Simple retry logic - try to fetch a small resource or just wait a moment
    // The online store will automatically update when navigator.onLine changes
   	window.location.reload();
    
    // Wait a moment then reset retrying state
    setTimeout(() => {
      isRetrying = false;
    }, 1000);
  }
</script>

<div class="min-h-screen flex items-center justify-center px-4">
  <div class="max-w-md w-full text-center">
    <!-- Offline Icon -->
    <div class="mb-8 flex justify-center">
      <div class="relative">
        <!-- WiFi Icon with X overlay -->
        <svg 
          class="w-24 h-24 text-[var(--color-secondary-text)]" 
          fill="none" 
          stroke="currentColor" 
          viewBox="0 0 24 24"
        >
          <!-- WiFi arcs -->
          <path 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            stroke-width="2" 
            d="M8.111 16.404a5.5 5.5 0 017.778 0M12 20h.01m-7.08-7.071c3.904-3.905 10.236-3.905 14.141 0M1.394 9.393c5.857-5.857 15.355-5.857 21.213 0"
          />
        </svg>
        
        <!-- Red X overlay -->
        <div class="absolute inset-0 flex items-center justify-center">
          <svg 
            class="w-10 h-10 text-red-500" 
            fill="none" 
            stroke="currentColor" 
            viewBox="0 0 24 24"
          >
            <path 
              stroke-linecap="round" 
              stroke-linejoin="round" 
              stroke-width="3" 
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </div>
      </div>
    </div>
    
    <!-- Title and Description -->
    <div class="mb-8">
      <h1 class="text-3xl md:text-4xl font-bold text-[var(--color-primary-text)] mb-4">
        You're Offline
      </h1>
      <p class="text-lg text-[var(--color-secondary-text)] leading-relaxed">
        It looks like you've lost your internet connection. Please check your network settings and try again.
      </p>
    </div>
    
    <!-- Retry Button -->
    <button
      on:click={handleRetry}
      disabled={isRetrying}
      class="inline-flex items-center px-6 py-3 bg-[var(--color-primary-accent)] hover:bg-[var(--color-primary-accent)]/90 disabled:bg-[var(--color-tertiary-text)] text-white font-semibold rounded-lg transition-all duration-200 transform hover:scale-105 disabled:scale-100 focus:outline-none focus:ring-2 focus:ring-[var(--color-primary-accent)] focus:ring-offset-2 focus:ring-offset-[var(--color-background-dark-sections)]"
    >
      {#if isRetrying}
        <!-- Spinner -->
        <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        Checking Connection...
      {:else}
        <!-- Refresh Icon -->
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
        </svg>
        Retry Connection
      {/if}
    </button>
    
    <!-- Additional Help Text -->
    <p class="mt-6 text-sm text-[var(--color-tertiary-text)]">
      Make sure your device is connected to a stable internet connection
    </p>
  </div>
</div>
