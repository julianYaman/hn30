<script>
  import { theme } from '$lib/stores/theme.js';
  import { bookmarks } from '$lib/stores/bookmarks.js';
  import { notifications } from '$lib/stores/notifications.js';
  import SettingsModal from './SettingsModal.svelte';
  import { fly } from 'svelte/transition';
  import { onMount } from 'svelte';

  let menuOpen = false;
  let settingsOpen = false;
  let isMobile = false;

  function toggleTheme() {
    theme.setTheme($theme === 'light' ? 'dark' : 'light');
  }

  onMount(() => {
    notifications.initialize();
    
    const mql = window.matchMedia('(max-width: 767px)');
    isMobile = mql.matches;
    const listener = (e) => {
      isMobile = e.matches;
      if (!isMobile) {
        menuOpen = false; // Close menu when switching to desktop view
      }
    };
    mql.addEventListener('change', listener);
    
    return () => {
      mql.removeEventListener('change', listener);
    }
  });
</script>

<header class="bg-[var(--color-background-card)] shadow-sm sticky top-0 z-20">
  <div class="max-w-7xl mx-auto px-4 py-5 flex justify-between items-center">
    <a href="/" class="flex items-center space-x-2">
      <span class="bg-[var(--color-primary-accent)] text-white font-bold text-xl rounded-md px-2 py-1">hn30</span>
      <h1 class="text-2xl font-bold text-[var(--color-primary-text)]">tech news</h1>
    </a>

    <!-- Desktop Nav -->
    <nav class="hidden md:flex items-center space-x-6">
      <a href="https://yaman.pro" class="text-[var(--color-secondary-text)] hover:text-[var(--color-primary-accent)] transition-colors" target="_blank" rel="noopener noreferrer">yaman.pro</a>
      <a href="https://yaman.pro/blog" class="text-[var(--color-secondary-text)] hover:text-[var(--color-primary-accent)] transition-colors" target="_blank" rel="noopener noreferrer">Blog</a>
      <a href="https://github.com/julianyaman/hn-tech-news" class="text-[var(--color-secondary-text)] hover:text-[var(--color-primary-accent)] transition-colors" target="_blank" rel="noopener noreferrer" title="GitHub">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 2C6.477 2 2 6.484 2 12.012c0 4.418 2.865 8.166 6.839 9.489.5.092.682-.217.682-.482 0-.237-.009-.868-.013-1.703-2.782.605-3.369-1.342-3.369-1.342-.454-1.157-1.11-1.465-1.11-1.465-.908-.62.069-.608.069-.608 1.004.07 1.532 1.032 1.532 1.032.892 1.53 2.341 1.089 2.91.833.091-.647.35-1.089.636-1.341-2.221-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.254-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.025A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.295 2.748-1.025 2.748-1.025.546 1.378.202 2.396.099 2.65.64.7 1.028 1.595 1.028 2.688 0 3.847-2.337 4.695-4.566 4.944.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .267.18.577.688.479C19.138 20.175 22 16.427 22 12.012 22 6.484 17.523 2 12 2z"/>
        </svg>
      </a>
      <a href="/bookmarks" class="text-[var(--color-secondary-text)] hover:text-[var(--color-primary-accent)] transition-colors" title="Bookmarks">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
        </svg>
      </a>
      <button on:click={() => settingsOpen = true} class="text-[var(--color-secondary-text)] hover:text-[var(--color-primary-accent)] transition-colors" title="Settings">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
        </svg>
      </button>
      <button on:click={toggleTheme} class="text-[var(--color-secondary-text)] hover:text-[var(--color-primary-accent)] transition-colors" title="Toggle Theme">
        {#if $theme === 'light'}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" /></svg>
        {:else}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" /></svg>
        {/if}
      </button>
    </nav>

    <!-- Mobile Menu Button -->
    <div class="md:hidden">
      <button on:click={() => menuOpen = true} aria-label="Open menu" class="text-[var(--color-primary-text)] p-2 -mr-2">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path></svg>
      </button>
    </div>
  </div>
</header>

<!-- Mobile Menu Panel -->
{#if menuOpen && isMobile}
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div transition:fly={{ duration: 200, opacity: 0 }} on:click={() => menuOpen = false} class="fixed inset-0 bg-white/30 backdrop-blur-sm z-30" role="button" tabindex="0"></div>

  <div transition:fly={{ duration: 300, x: '100%' }} class="fixed top-0 right-0 h-full w-64 bg-[var(--color-background-card)] shadow-lg z-40">
    <div class="p-5 flex justify-end">
        <button on:click={() => menuOpen = false} aria-label="Close menu" class="text-[var(--color-primary-text)] p-2 -mr-2">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
        </button>
    </div>
    <nav class="flex flex-col text-lg">
      <a href="https://yaman.pro" on:click={() => menuOpen = false} class="p-4 hover:bg-[var(--color-background-dark-sections)]" target="_blank" rel="noopener noreferrer">yaman.pro</a>
      <a href="https://yaman.pro/blog" on:click={() => menuOpen = false} class="p-4 hover:bg-[var(--color-background-dark-sections)]" target="_blank" rel="noopener noreferrer">Blog</a>
      <a href="https://github.com/julianyaman/hn-tech-news" on:click={() => menuOpen = false} class="p-4 hover:bg-[var(--color-background-dark-sections)]" target="_blank" rel="noopener noreferrer">GitHub</a>
      <a href="/bookmarks" on:click={() => menuOpen = false} class="p-4 hover:bg-[var(--color-background-dark-sections)]">Bookmarks</a>
      <button on:click={() => { settingsOpen = true; menuOpen = false; }} class="p-4 text-left hover:bg-[var(--color-background-dark-sections)] w-full">
        Settings
      </button>
      <button on:click={() => { toggleTheme(); menuOpen = false; }} class="p-4 text-left hover:bg-[var(--color-background-dark-sections)] w-full">
        Toggle Theme
      </button>
    </nav>
  </div>
{/if}

<SettingsModal bind:isOpen={settingsOpen} />
