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

/**
 * Returns true if the app is running in an "installed" context.
 *
 * - On Chromium: checks display-mode media query (standalone).
 * - On iOS Safari: checks navigator.standalone (non-standard but common).
 * - Optional: detects Trusted Web Activity (TWA) via referrer.
 *
 * References:
 * - MDN: display-mode / matchMedia standalone detection [1](https://developer.mozilla.org/en-US/docs/Web/Progressive_web_apps/How_to/Create_a_standalone_app)
 * - web.dev: display-mode & TWA referrer detection patterns [2](https://web.dev/learn/pwa/detection/)
 */
export function isInstalledStandalonePWA() {
  if (typeof window === 'undefined') return false;

  // Optional: if your app is also shipped as a TWA wrapper
  const isTWA =
    typeof document !== 'undefined' &&
    typeof document.referrer === 'string' &&
    document.referrer.startsWith('android-app://'); // web.dev pattern [2](https://web.dev/learn/pwa/detection/)

  const mql = (query) => window.matchMedia?.(query)?.matches === true;

  // Manifest display is "standalone" (your case), so this is the main check.
  const isStandaloneDisplayMode = mql('(display-mode: standalone)'); // MDN [1](https://developer.mozilla.org/en-US/docs/Web/Progressive_web_apps/How_to/Create_a_standalone_app)

  // iOS "Added to Home Screen" / standalone mode
  const isIOSStandalone = window.navigator?.standalone === true;

  return Boolean(isTWA || isStandaloneDisplayMode || isIOSStandalone);
}

/**
 * Returns true for "mobile-like" devices, including iPadOS.
 *
 * We avoid pure UA sniffing where possible and prefer capability checks:
 * - pointer: coarse AND hover: none are good indicators of touch-first devices
 * - maxTouchPoints helps detect touch devices (and iPadOS masquerading as Mac)
 */
export function isMobileOrIPad() {
  if (typeof window === 'undefined') return false;

  const coarsePointer = window.matchMedia?.('(pointer: coarse)')?.matches === true;
  const noHover = window.matchMedia?.('(hover: none)')?.matches === true;
  const touchPoints = navigator.maxTouchPoints ?? 0;

  // iPadOS (often reports platform "MacIntel" but has touch points)
  const isIPadOS = /MacIntel/.test(navigator.platform) && touchPoints > 1;

  // Typical mobile/touch device
  const isTouchMobileLike = coarsePointer && noHover && touchPoints > 0;

  return Boolean(isTouchMobileLike || isIPadOS);
}


function devForcePullToRefresh() {
  // Only allow override in dev
  if (!import.meta.env.DEV) return false;

  // Enable via URL: ?__pwa=1
  const qs = new URLSearchParams(location.search);
  if (qs.get('__pwa') === '1') return true;

  // Or via localStorage: localStorage.__pwa = "1"
  if (localStorage.getItem('__pwa') === '1') return true;

  return false;
}


/**
 * Your final gate:
 * True only when:
 * - installed standalone PWA context
 * - AND mobile/iPad
 *
 * This matches your goal: "only for installed mobile PWAs; browsers are an exception".
 */
export function isStandaloneMobilePWA() {

  if (typeof window === 'undefined') return false;

  if (devForcePullToRefresh()) return true;

  return isInstalledStandalonePWA() && isMobileOrIPad();
}

