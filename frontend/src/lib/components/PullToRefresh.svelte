
<script>
  import { onMount, onDestroy } from 'svelte';
  import { invalidateAll } from '$app/navigation';
  import { isStandaloneMobilePWA } from '$lib/utils.js';

  let startY = 0;
  let currentY = 0;
  let pullDistance = 0;

  let isPulling = false;
  let isRefreshing = false;
  let isEnabled = false;

  const PULL_THRESHOLD = 80;
  const MAX_PULL = 160;

  // For smoother UI updates during touchmove
  let rafId = 0;
  let lastTouchY = 0;
  let contentEl;
  let containerEl;

  const listenerOptions = { passive: false };

  function getScrollTop() {
    return document.scrollingElement?.scrollTop ?? window.scrollY ?? 0;
  }

  function isAtTop() {
    return getScrollTop() <= 0;
  }

  function isInteractiveElement(el) {
    if (!el || el.nodeType !== 1) return false;
    const tag = el.tagName;
    return (
      tag === 'INPUT' ||
      tag === 'TEXTAREA' ||
      tag === 'SELECT' ||
      tag === 'BUTTON' ||
      el.isContentEditable === true
    );
  }

  function handleTouchStart(e) {
    if (!isEnabled || isRefreshing) return;
    if (e.touches.length !== 1) return;
    if (!isAtTop()) return;

    if (isInteractiveElement(e.target)) return;

    startY = e.touches[0].clientY;
    lastTouchY = startY;
    isPulling = true;
    
    // Disable transition during pull for immediate response
    if (contentEl) {
      contentEl.style.transition = 'none';
    }
  }

  function applyPullUpdate() {
    rafId = 0;
    if (!isPulling && !isRefreshing) return;

    // Calculate resistence/damping
    const rawDistance = lastTouchY - startY;
    // Logarithmic damping function roughly similar to iOS
    pullDistance = Math.min(Math.max(rawDistance * 0.5, 0), MAX_PULL);

    if (rawDistance < 0 || !isAtTop()) {
       resetPull();
       return;
    }

    // Direct DOM manipulation for performance
    if (contentEl) {
      contentEl.style.transform = `translateY(${pullDistance}px)`;
    }
    // Also height of the pull container - update directly to avoid state sync lag
    // Using variable reference if bound, or by ID? Let's use binding in markup mostly, 
    // but here we can force it if needed. However, since we removed transition on container,
    // svelte binding might be enough. Let's rely on Svelte binding for height, 
    // but ensure we update variable. `pullDistance` is updated. 
    // But we need to disable transition on container too.
  }

  function handleTouchMove(e) {
    if (!isEnabled || (!isPulling && !isRefreshing)) return;
    if (e.touches.length !== 1) return;

    lastTouchY = e.touches[0].clientY;
    
    // Only prevent browser scrolling if pulling down at the top
    const rawDistance = lastTouchY - startY;
    if (rawDistance > 0 && isAtTop() && e.cancelable) {
      // Prevent native scroll/refresh
      e.preventDefault();
      
      if (!isPulling) {
        // Recover if we started pulling while "scrolling up to top"
        isPulling = true;
        startY = lastTouchY; 
      }
    }

    if (!rafId) {
      rafId = requestAnimationFrame(applyPullUpdate);
    }
  }

  async function handleTouchEnd() {
    if (!isEnabled || !isPulling || isRefreshing) return;

    if (rafId) {
      cancelAnimationFrame(rafId);
      rafId = 0;
    }

    const wasOverThreshold = pullDistance >= PULL_THRESHOLD;

    // Immediately stop the "pulling" state so that CSS transition: height kicks in
    isPulling = false;

    // Ensure page content container has transition enabled for the snap-back
    if (contentEl) {
      contentEl.style.transition = 'transform 0.3s cubic-bezier(0.23, 1, 0.32, 1)';
    }

    if (wasOverThreshold) {
      isRefreshing = true;
      // Animate back to exactly the threshold height
      pullDistance = PULL_THRESHOLD;

      if (contentEl) {
        contentEl.style.transform = `translateY(${PULL_THRESHOLD}px)`;
      }

      try {
        await invalidateAll();
        // Artificial delay if request is too fast, so user sees the spinner
        await new Promise((resolve) => setTimeout(resolve, 400));
      } catch (err) {
        console.error('Refresh failed:', err);
      } finally {
        isRefreshing = false;
        resetPull();
      }
    } else {
      resetPull();
    }
  }

  function resetPull() {
    isPulling = false;
    pullDistance = 0;
    startY = 0;
    currentY = 0;
    lastTouchY = 0;
    
    if (contentEl) {
        // Ensure transition is on for the reset/snapback
        contentEl.style.transition = 'transform 0.3s cubic-bezier(0.23, 1, 0.32, 1)';
        contentEl.style.transform = '';
    }
  }

  onMount(() => {
    isEnabled = isStandaloneMobilePWA();

    if (isEnabled) {
      contentEl = document.getElementById('ptr-content');
      document.addEventListener('touchstart', handleTouchStart, listenerOptions);
      document.addEventListener('touchmove', handleTouchMove, listenerOptions);
      document.addEventListener('touchend', handleTouchEnd, listenerOptions);
      document.addEventListener('touchcancel', handleTouchEnd, listenerOptions);
    }
  });

  onDestroy(() => {
    if (rafId) cancelAnimationFrame(rafId);
    if (isEnabled) {
      document.removeEventListener('touchstart', handleTouchStart, listenerOptions);
      document.removeEventListener('touchmove', handleTouchMove, listenerOptions);
      document.removeEventListener('touchend', handleTouchEnd, listenerOptions);
      document.removeEventListener('touchcancel', handleTouchEnd, listenerOptions);
    }
  });

  $: isReady = pullDistance >= PULL_THRESHOLD;
</script>

{#if isEnabled}
  <div
    bind:this={containerEl}
    class="pull-to-refresh-container"
    class:pulling={isPulling}
    style="height: {pullDistance}px;"
    aria-hidden={!isPulling && !isRefreshing}
  >
    <div class="content">
      <div class="icon-wrapper" class:spinning={isRefreshing} style:transform={isRefreshing ? 'none' : `rotate(${Math.min((pullDistance/PULL_THRESHOLD) * 180, 180)}deg)`}>
        {#if isRefreshing}
          <!-- Spinner -->
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 12a9 9 0 1 1-6.219-8.56"></path>
          </svg>
        {:else}
           <!-- Arrow Down -->
           <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
             <line x1="12" y1="5" x2="12" y2="19"></line>
             <polyline points="19 12 12 19 5 12"></polyline>
           </svg>
        {/if}
      </div>
      <span class="label">
        {isRefreshing ? 'Refreshing...' : isReady ? 'Release to refresh' : 'Pull to refresh'}
      </span>
    </div>
  </div>
{/if}

<style>
  .pull-to-refresh-container {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    background-color: var(--color-primary-accent);
    color: white; /* Ensure good contrast on primary color */
    z-index: 1000;
    display: flex;
    justify-content: center;
    align-items: center;
    box-shadow: 0 4px 10px rgba(0,0,0,0.1);
    overflow: hidden;
    will-change: height;
    transition: height 0.3s cubic-bezier(0.23, 1, 0.32, 1); 
  }
  
  .pull-to-refresh-container.pulling {
    transition: none;
  }
  
  .content {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
  }

  .label {
    font-weight: 600;
    font-size: 0.95rem;
  }

  .icon-wrapper {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .spinning {
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }
</style>
