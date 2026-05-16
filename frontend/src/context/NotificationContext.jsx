import { createContext, useCallback, useContext, useEffect, useMemo, useRef, useState } from 'react';
import { notificationsAPI } from '../services/api.js';

const NotificationContext = createContext(null);
const POLL_INTERVAL_MS = 15000;

function normalizeForegroundPayload(payload) {
  return {
    title: payload?.notification?.title || 'New notification',
    body: payload?.notification?.body || '',
  };
}

export function NotificationProvider({ children }) {
  const [notifications, setNotifications] = useState([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [hasNewNotifications, setHasNewNotifications] = useState(false);
  const [toasts, setToasts] = useState([]);

  const seenIdsRef = useRef(new Set());
  const initialLoadRef = useRef(true);
  const toastIdRef = useRef(0);
  const toastTimersRef = useRef(new Map());

  const dismissToast = useCallback((toastId) => {
    const timer = toastTimersRef.current.get(toastId);
    if (timer) {
      clearTimeout(timer);
      toastTimersRef.current.delete(toastId);
    }

    setToasts((current) => current.filter((toast) => toast.id !== toastId));
  }, []);

  const pushToast = useCallback((title, message) => {
    const id = ++toastIdRef.current;
    setToasts((current) => [...current.slice(-3), { id, title, message }]);

    const timer = window.setTimeout(() => {
      dismissToast(id);
    }, 4500);

    toastTimersRef.current.set(id, timer);
  }, [dismissToast]);

  const refreshNotifications = useCallback(async ({ silent = false, markSeen = false } = {}) => {
    if (!silent) {
      setRefreshing(true);
    }

    try {
      const data = await notificationsAPI.getAll(20);
      const nextNotifications = Array.isArray(data?.notifications) ? data.notifications : [];
      const nextIds = new Set(nextNotifications.map((notification) => notification.id));

      if (initialLoadRef.current) {
        seenIdsRef.current = nextIds;
        initialLoadRef.current = false;
        setNotifications(nextNotifications);
        setHasNewNotifications(false);
        return;
      }

      const previousIds = seenIdsRef.current;
      const newNotifications = nextNotifications.filter((notification) => !previousIds.has(notification.id));

      setNotifications(nextNotifications);
      seenIdsRef.current = nextIds;

      if (newNotifications.length > 0 && !markSeen) {
        setHasNewNotifications(true);
        newNotifications.slice(0, 3).forEach((notification) => {
          pushToast('New notification', notification.message || 'You have a new notification');
        });
      }

      if (markSeen) {
        setHasNewNotifications(false);
      }
    } catch (error) {
      console.error('Failed to refresh notifications:', error);
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  }, [pushToast]);

  const openDrawer = useCallback(async () => {
    setDrawerOpen(true);
    setHasNewNotifications(false);
    await refreshNotifications({ silent: true, markSeen: true });
  }, [refreshNotifications]);

  const closeDrawer = useCallback(() => {
    setDrawerOpen(false);
  }, []);

  useEffect(() => {
    refreshNotifications({ silent: false, markSeen: true });

    const intervalId = window.setInterval(() => {
      refreshNotifications({ silent: true });
    }, POLL_INTERVAL_MS);

    return () => {
      window.clearInterval(intervalId);
      toastTimersRef.current.forEach((timer) => clearTimeout(timer));
      toastTimersRef.current.clear();
    };
  }, [refreshNotifications]);

  useEffect(() => {
    const handleForegroundNotification = (event) => {
      const detail = normalizeForegroundPayload(event.detail);
      pushToast(detail.title, detail.body || detail.title);
      setHasNewNotifications(true);
    };

    window.addEventListener('app:foreground-notification', handleForegroundNotification);
    return () => {
      window.removeEventListener('app:foreground-notification', handleForegroundNotification);
    };
  }, [pushToast]);

  const value = useMemo(() => ({
    notifications,
    loading,
    refreshing,
    drawerOpen,
    hasNewNotifications,
    toasts,
    refreshNotifications,
    openDrawer,
    closeDrawer,
    dismissToast,
  }), [
    notifications,
    loading,
    refreshing,
    drawerOpen,
    hasNewNotifications,
    toasts,
    refreshNotifications,
    openDrawer,
    closeDrawer,
    dismissToast,
  ]);

  return (
    <NotificationContext.Provider value={value}>
      {children}
    </NotificationContext.Provider>
  );
}

export function useNotifications() {
  const context = useContext(NotificationContext);
  if (!context) {
    throw new Error('useNotifications must be used within NotificationProvider');
  }
  return context;
}