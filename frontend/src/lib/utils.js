/**
 * Generates a hex color code from a string.
 * @param {string} str The input string.
 * @returns {string} A 6-character hex color code.
 */
export function stringToColor(str) {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash);
  }
  let color = '';
  for (let i = 0; i < 3; i++) {
    const value = (hash >> (i * 8)) & 0xFF;
    color += ('00' + value.toString(16)).substr(-2);
  }
  return color;
}

/**
 * Generates a placeholder URL from a string.
 * @param {string} str The input string.
 * @param {number} width The width of the placeholder.
 * @param {number} height The height of the placeholder.
 * @returns {string} A URL to a placeholder image.
 */
export function generatePlaceholder(str, width, height) {
    const color1 = stringToColor(str);
    const color2 = stringToColor(str.split('').reverse().join(''));
    return `https://colorr.me/i/g/${color1}-${color2}?w=${width}&h=${height}`;
}

/**
 * Extracts the domain from a URL.
 * @param {string} url The input URL.
 * @returns {string} The domain name.
 */
export function getDomain(url) {
  if (!url) return '';
  try {
    const domain = new URL(url).hostname;
    return domain.replace(/^www\./, '');
  } catch (e) {
    return '';
  }
}

/**
 * Formats a unix timestamp into a human-readable string.
 * @param {number} unixTimestamp The unix timestamp.
 * @returns {string} A human-readable time ago string.
 */
export function timeAgo(ts) {
  const seconds = Math.floor((new Date() - new Date(ts * 1000)) / 1000);
  let interval = seconds / 31536000;
  if (interval > 1) return Math.floor(interval) + " years ago";
  interval = seconds / 2592000;
  if (interval > 1) return Math.floor(interval) + " months ago";
  interval = seconds / 86400;
  if (interval > 1) return Math.floor(interval) + " days ago";
  interval = seconds / 3600;
  if (interval > 1) return Math.floor(interval) + " hours ago";
  interval = seconds / 60;
  if (interval > 1) return Math.floor(interval) + " minutes ago";
  return Math.floor(seconds) + " seconds ago";
}

export function proxyImageUrl(url) {
  if (!url) return '';
  return `/api/image-proxy?url=${encodeURIComponent(url)}`;
}
