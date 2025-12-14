import { writable, get } from 'svelte/store';
import { bookmarksDB } from '$lib/db';
import { toast } from '$lib/stores/toast';
import { exportBookmarksToJSON, readJSONFile, validateImportData, downloadFile, generateExportFilename } from '$lib/utils/fileOperations';


const createBookmarkStore = () => {
  const { subscribe, set } = writable([]);

  /**
   * Initializes the store by loading all bookmarks from IndexedDB.
   * This should be called once when the application starts.
   */
  async function init() {
    if (typeof window !== 'undefined') {
      const bookmarks = await bookmarksDB.getAll();
      set(bookmarks);
    }
  }

  /**
   * Toggles a bookmark's state.
   * It checks if the story is already bookmarked and either adds or removes it.
   * @param {object} story - The story object to bookmark/unbookmark.
   */
  async function toggle(story) {
    const currentBookmarks = get({ subscribe });
    const isBookmarked = currentBookmarks.some(b => b.id === story.id);

    if (isBookmarked) {
      await bookmarksDB.delete(story.id);
      set(currentBookmarks.filter(b => b.id !== story.id));
      toast.show('Bookmark removed', 'info');
    } else {
      const newBookmark = {
        id: story.id,
        title: story.title,
        ogImage: story.ogImage,
        url: story.url,
        savedAt: new Date().toISOString(),
        postedAt: story.time ? new Date(story.time * 1000).toISOString() : ''
      };
      await bookmarksDB.put(newBookmark);
      set([...currentBookmarks, newBookmark]);
      toast.show('Bookmark saved!', 'success');
    }
  }

  async function remove(id) {
    const currentBookmarks = get({ subscribe });
    const exists = currentBookmarks.some(b => b.id === id);
    if (!exists) return;
    await bookmarksDB.delete(id);
    set(currentBookmarks.filter(b => b.id !== id));
    toast.show('Bookmark removed', 'info');
  }

  async function add(story) {
    const currentBookmarks = get({ subscribe });
    if (currentBookmarks.some(b => b.id === story.id)) return;
    const newBookmark = {
      id: story.id,
      title: story.title,
      ogImage: story.ogImage,
      url: story.url,
      savedAt: new Date().toISOString(),
      postedAt: story.time ? new Date(story.time * 1000).toISOString() : ''
    };
    await bookmarksDB.put(newBookmark);
    set([...currentBookmarks, newBookmark]);
    toast.show('Bookmark saved!', 'success');
  }

  function isBookmarked(id) {
    const currentBookmarks = get({ subscribe });
    return currentBookmarks.some(b => b.id === id);
  }

  /**
   * Export bookmarks to JSON file
   * @param {Object} options - Export options
   * @param {string} options.filename - Custom filename (optional)
   */
  async function exportBookmarks(options = {}) {
    try {
      const currentBookmarks = get({ subscribe });
      const filename = options.filename || generateExportFilename('json');
      
      exportBookmarksToJSON(currentBookmarks, filename);
      toast.show(`Successfully exported ${currentBookmarks.length} bookmarks`, 'success');
    } catch (error) {
      console.error('Export failed:', error);
      toast.show('Failed to export bookmarks', 'error');
    }
  }

  /**
   * Import bookmarks from JSON file
   * @param {File} file - File to import
   * @param {Object} options - Import options
   * @param {boolean} options.skipDuplicates - Skip bookmarks that already exist (default: true)
   */
  async function importBookmarks(file, options = {}) {
    try {
      const skipDuplicates = options.skipDuplicates !== false; // Default to true
      
      // Read and parse the file
      const importData = await readJSONFile(file);
      
      // Validate the data structure
      const validation = validateImportData(importData);
      
      if (!validation.isValid) {
        throw new Error(`Invalid file format: ${validation.errors.join(', ')}`);
      }
      
      const currentBookmarks = get({ subscribe });
      const existingIds = new Set(currentBookmarks.map(b => b.id));
      
      // Filter bookmarks based on duplicate handling preference
      let bookmarksToImport = validation.validBookmarks;
      if (skipDuplicates) {
        bookmarksToImport = bookmarksToImport.filter(bookmark => !existingIds.has(bookmark.id));
      }
      
      // Batch import valid bookmarks
      let importedCount = 0;
      let skippedCount = validation.validBookmarks.length - bookmarksToImport.length;
      
      for (const bookmark of bookmarksToImport) {
        try {
          await bookmarksDB.put(bookmark);
          importedCount++;
        } catch (error) {
          console.error('Failed to import bookmark:', bookmark.id, error);
        }
      }
      
      // Update the store with new bookmarks
      const updatedBookmarks = [...currentBookmarks, ...bookmarksToImport];
      set(updatedBookmarks);
      
      // Show success message
      let message = `Successfully imported ${importedCount} bookmarks`;
      if (skippedCount > 0) {
        message += ` (${skippedCount} skipped as duplicates)`;
      }
      toast.show(message, 'success');
      
      return {
        imported: importedCount,
        skipped: skippedCount,
        total: validation.validBookmarks.length
      };
      
    } catch (error) {
      console.error('Import failed:', error);
      const errorMessage = error.message || 'Failed to import bookmarks';
      toast.show(errorMessage, 'error');
      throw error;
    }
  }

  return {
    subscribe,
    init,
    toggle,
    add,
    remove,
    isBookmarked,
    exportBookmarks,
    importBookmarks,
  };
};

export const bookmarks = createBookmarkStore();

// Named helper to match existing component imports
export const toggleBookmark = (story) => bookmarks.toggle(story);
