/**
 * File operations utility for bookmark import/export functionality
 * Handles JSON format for HN30-internal bookmark management
 */

/**
 * Export bookmarks to JSON format
 * @param {Array} bookmarks - Array of bookmark objects
 * @param {string} filename - Optional filename (defaults to timestamp-based name)
 * @returns {string} JSON string for download
 */
export function exportBookmarksToJSON(bookmarks, filename = null) {
  const exportData = {
    exportDate: new Date().toISOString(),
    format: 'hn-news-bookmarks-v1',
    totalBookmarks: bookmarks.length,
    bookmarks: bookmarks.map(bookmark => ({
      id: bookmark.id,
      title: bookmark.title,
      ogImage: bookmark.ogImage || '',
      url: bookmark.url,
      savedAt: bookmark.savedAt,
      postedAt: bookmark.postedAt || ''
    }))
  };

  const jsonString = JSON.stringify(exportData, null, 2);
  
  if (filename) {
    downloadFile(jsonString, filename, 'application/json');
  }
  
  return jsonString;
}

/**
 * Download file to user's device
 * @param {string} content - File content
 * @param {string} filename - Name for the downloaded file
 * @param {string} mimeType - MIME type for the file
 */
export function downloadFile(content, filename, mimeType = 'application/json') {
  // Create blob with the content
  const blob = new Blob([content], { type: mimeType });
  
  // Create download link
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  
  // Trigger download
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  
  // Clean up the object URL
  URL.revokeObjectURL(url);
}

/**
 * Read and parse a JSON file
 * @param {File} file - The file to read
 * @returns {Promise<Object>} Parsed JSON data
 */
export function readJSONFile(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    
    reader.onload = (event) => {
      try {
        const content = event.target.result;
        const data = JSON.parse(content);
        resolve(data);
      } catch (error) {
        reject(new Error(`Invalid JSON file: ${error.message}`));
      }
    };
    
    reader.onerror = () => {
      reject(new Error('Failed to read file'));
    };
    
    reader.readAsText(file);
  });
}

/**
 * Validate bookmark data structure
 * @param {Object} bookmark - Bookmark object to validate
 * @returns {Object} Validation result with isValid and errors
 */
export function validateBookmark(bookmark) {
  const errors = [];
  
  // Required fields
  if (!bookmark.id && bookmark.id !== 0) {
    errors.push('Missing required field: id');
  }
  
  if (!bookmark.title) {
    errors.push('Missing required field: title');
  }
  
  if (!bookmark.url) {
    errors.push('Missing required field: url');
  }
  
  // Validate URL format
  if (bookmark.url) {
    try {
      new URL(bookmark.url);
    } catch {
      errors.push('Invalid URL format');
    }
  }
  
  // Validate timestamps if present
  if (bookmark.savedAt && !isValidDate(bookmark.savedAt)) {
    errors.push('Invalid savedAt timestamp format');
  }
  
  if (bookmark.postedAt && !isValidDate(bookmark.postedAt)) {
    errors.push('Invalid postedAt timestamp format');
  }
  
  return {
    isValid: errors.length === 0,
    errors
  };
}

/**
 * Validate entire import data structure
 * @param {Object} data - Parsed JSON data
 * @returns {Object} Validation result with isValid and errors
 */
export function validateImportData(data) {
  const errors = [];
  
  // Check if data is an object
  if (typeof data !== 'object' || data === null) {
    errors.push('Data must be a valid JSON object');
    return { isValid: false, errors };
  }
  
  // Check for bookmarks array
  if (!Array.isArray(data.bookmarks)) {
    errors.push('Missing or invalid bookmarks array');
    return { isValid: false, errors };
  }
  
  // Validate each bookmark
  const bookmarkErrors = [];
  data.bookmarks.forEach((bookmark, index) => {
    const validation = validateBookmark(bookmark);
    if (!validation.isValid) {
      bookmarkErrors.push(`Bookmark ${index + 1}: ${validation.errors.join(', ')}`);
    }
  });
  
  errors.push(...bookmarkErrors);
  
  return {
    isValid: errors.length === 0,
    errors,
    totalBookmarks: data.bookmarks.length,
    validBookmarks: data.bookmarks.filter(bookmark => validateBookmark(bookmark).isValid)
  };
}

/**
 * Check if a string is a valid date
 * @param {string} dateString - Date string to validate
 * @returns {boolean} True if valid date
 */
function isValidDate(dateString) {
  const date = new Date(dateString);
  return date instanceof Date && !isNaN(date);
}

/**
 * Generate filename for export
 * @param {string} extension - File extension (without dot)
 * @returns {string} Generated filename
 */
export function generateExportFilename(extension = 'json') {
  const now = new Date();
  const timestamp = now.toISOString().slice(0, 19).replace(/[:.]/g, '-');
  return `hn-bookmarks-${timestamp}.${extension}`;
}