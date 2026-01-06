import { writable } from 'svelte/store';

const preferenceKey = 'notifications_preference';

const initialState = {
  wantsNotifications: false,
  enabled: false,
  loading: false,
  loadingMessage: null,
  isInitialized: false,
  error: null,
  permission: 'default',
  cookiesAccepted: false,
  message: null,
  messageType: null
};

const createNotificationsStore = () => {
  const { subscribe, update } = writable(initialState);

  const checkCookies = () => {
    if (typeof window === 'undefined') return false;
    return localStorage.getItem('privacy_accepted') === 'true';
  };

  const getStoredPreference = () => {
    if (typeof window === 'undefined') return false;
    return localStorage.getItem(preferenceKey) === 'true';
  };

  const persistPreference = (value) => {
    if (typeof window === 'undefined') return;
    localStorage.setItem(preferenceKey, value ? 'true' : 'false');
  };

  const syncState = async () => {
    if (typeof window === 'undefined') return;

    const cookiesAccepted = checkCookies();
    const permission = window.Notification?.permission || 'default';
    const isSubscribed = window.OneSignal?.User?.PushSubscription?.optedIn === true;
    const storedPref = cookiesAccepted ? getStoredPreference() : false;
    
    // Use OneSignal as source of truth, but only reconcile if:
    // 1. OneSignal is fully initialized AND
    // 2. There's a clear mismatch (not during initial setup)
    let wantsNotifications = storedPref;
    
    // Only reconcile if OneSignal state is definitive (user is actually subscribed)
    // This prevents overwriting during the brief moment after optIn but before OneSignal updates
    if (window.OneSignal && isSubscribed && !storedPref) {
      // User is subscribed in OneSignal but localStorage says no - trust OneSignal
      wantsNotifications = true;
      persistPreference(true);
    } else if (window.OneSignal && !isSubscribed && storedPref && permission !== 'default') {
      // User is NOT subscribed in OneSignal but localStorage says yes - trust OneSignal
      // But only if permission is not 'default' (meaning the flow has completed)
      wantsNotifications = false;
      persistPreference(false);
    }

    update((state) => ({
      ...state,
      cookiesAccepted,
      permission,
      wantsNotifications,
      enabled: permission === 'granted' && isSubscribed,
      loading: false,
      loadingMessage: null
    }));
  };

  const requireCookies = () => {
    const cookiesAccepted = checkCookies();
    if (cookiesAccepted) return true;

    update((state) => ({
      ...state,
      cookiesAccepted: false,
      wantsNotifications: false,
      loading: false,
      loadingMessage: null,
      error: 'Please accept cookies to enable notifications.',
      message: null,
      messageType: null
    }));
    return false;
  };

  return {
    subscribe,

    initialize: async () => {
      if (typeof window === 'undefined') return;

      // Check browser compatibility
      if (!window.Notification) {
        update((state) => ({
          ...state,
          isInitialized: true,
          error: 'Notifications are not supported in your browser.'
        }));
        return;
      }

      await syncState();

      const OS = window.OneSignal || window.OneSignalDeferred;
      if (!OS) {
        update((state) => ({ ...state, isInitialized: true }));
        return;
      }

      OS.push(() => {
        syncState();

        window.OneSignal?.Notifications?.addEventListener('permissionChange', (event) => {
          update((state) => ({
            ...state,
            permission: event?.to || window.Notification?.permission || 'default',
            error:
              event?.to === 'denied'
                ? 'Notifications are blocked in your browser settings.'
                : state.error
          }));
        });

        window.OneSignal?.User?.PushSubscription?.addEventListener('change', () => {
          syncState();
        });

        update((state) => ({ ...state, isInitialized: true }));
      });
    },

    setNotificationsEnabled: async (shouldEnable) => {
      if (typeof window === 'undefined') return;
      if (!requireCookies()) return;

      update((state) => ({
        ...state,
        wantsNotifications: shouldEnable,
        loading: true,
        loadingMessage: shouldEnable ? 'Enabling notifications…' : 'Disabling notifications…',
        error: null,
        message: null,
        messageType: null
      }));

      try {
        const oneSignal = window.OneSignal;
        if (!oneSignal) throw new Error('Notification service unavailable.');

        // Helper to add timeout to promises
        const withTimeout = (promise, ms = 10000) => {
          return Promise.race([
            promise,
            new Promise((_, reject) => 
              setTimeout(() => reject(new Error('Operation timed out')), ms)
            )
          ]);
        };

        if (shouldEnable) {
          const permission = await withTimeout(
            oneSignal.Notifications.requestPermission(),
            10000
          );

          if (permission === 'denied') {
            persistPreference(false);
            update((state) => ({
              ...state,
              permission: 'denied',
              wantsNotifications: false,
              error: 'Notifications are blocked in your browser settings.'
            }));
            return;
          }

          if (permission === 'granted') {
            await withTimeout(
              oneSignal.User.PushSubscription.optIn(),
              10000
            );
          }
        } else {
          await withTimeout(
            oneSignal.User.PushSubscription.optOut(),
            10000
          );
        }

        // Persist preference AFTER OneSignal operation succeeds
        persistPreference(shouldEnable);
        
        // Get the actual subscription state from OneSignal
        const isSubscribed = window.OneSignal?.User?.PushSubscription?.optedIn === true;
        const permission = window.Notification?.permission || 'default';
        
        // Update state immediately with the new values
        // Don't call syncState() here - let OneSignal's change event handle it
        update((state) => ({
          ...state,
          wantsNotifications: shouldEnable,
          enabled: permission === 'granted' && isSubscribed,
          permission,
          messageType: 'success',
          message: shouldEnable ? 'Notifications enabled.' : 'Notifications disabled.',
          error: null,
          loading: false,
          loadingMessage: null
        }));

        // Auto-clear success message after 4 seconds
        setTimeout(() => {
          update((state) => ({
            ...state,
            message: null,
            messageType: null
          }));
        }, 4000);
      } catch (error) {
        console.error('Notification error:', error);
        persistPreference(false);
        update((state) => ({
          ...state,
          wantsNotifications: false,
          loading: false,
          loadingMessage: null,
          error: error?.message || 'Unable to update notifications right now.'
        }));
      }
    },

    refresh: syncState,

    clearMessage: () => {
      update((state) => ({
        ...state,
        message: null,
        messageType: null,
        error: null
      }));
    }
  };
};

export const notifications = createNotificationsStore();
