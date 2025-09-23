import { writable } from 'svelte/store';

/**
 * Creates a custom Svelte store that holds a Set of URLs.
 * This is designed to be a session-wide cache to track loaded images.
 */
function createUrlSetStore() {
	const { subscribe, update } = writable(new Set());

	return {
		subscribe,
		/**
		 * Adds a URL to the set. The `update` method ensures that all
		 * subscribers are properly notified of the change.
		 * @param {string} url The image URL to add to the cache.
		 */
		add: (url) => {
			update((set) => {
				// Only update if the URL is not already in the set to avoid unnecessary updates.
				if (!set.has(url)) {
					set.add(url);
				}
				return set;
			});
		}
	};
}

export const imageCache = createUrlSetStore();
