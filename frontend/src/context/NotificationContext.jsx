import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import { notificationsAPI } from "../services/api.js";

const NotificationContext = createContext(null);

function normalizeForegroundPayload(payload) {
  return {
    title: payload?.notification?.title || "New notification",
    body: payload?.notification?.body || "",
  };
}

function normalizeIncomingNotification(
  payload,
  fallbackTitle = "New notification",
) {
  const data = payload?.data || payload?.raw?.data || {};
  const notificationId =
    data.notification_id || payload?.notification_id || `push-${Date.now()}`;
  const now = new Date().toISOString();

  return {
    id: notificationId,
    recipient_id: data.recipient_id || "",
    event_id: data.notification_id || notificationId,
    entity_id: data.entity_id || "",
    entity_type: data.entity_type || "",
    message: data.message || payload?.notification?.body || fallbackTitle,
    channels: {
      in_app: {
        status: "SENT",
        read_at: null,
      },
      push: {
        status: "DELIVERED",
      },
      email: {
        status: "SKIPPED",
      },
    },
    created_at: now,
  };
}

export function NotificationProvider({ children }) {
  const [notifications, setNotifications] = useState([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [hasNewNotifications, setHasNewNotifications] = useState(false);
  const [toasts, setToasts] = useState([]);

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

  const pushToast = useCallback(
    (title, message) => {
      const id = ++toastIdRef.current;
      setToasts((current) => [...current.slice(-3), { id, title, message }]);

      const timer = window.setTimeout(() => {
        dismissToast(id);
      }, 4500);

      toastTimersRef.current.set(id, timer);
    },
    [dismissToast],
  );

  const refreshNotifications = useCallback(async ({ silent = false } = {}) => {
    console.log("[Notifications] refreshNotifications called", { silent });

    if (!silent) {
      setRefreshing(true);
    }

    try {
      const data = await notificationsAPI.getAll(20);
      const nextNotifications = Array.isArray(data?.notifications)
        ? data.notifications
        : [];

      console.log("[Notifications] refreshNotifications received", {
        count: nextNotifications.length,
        unreadCount: nextNotifications.filter(
          (notification) => !notification.channels?.in_app?.read_at,
        ).length,
      });

      setNotifications(nextNotifications);

      setHasNewNotifications(
        nextNotifications.some(
          (notification) => !notification.channels?.in_app?.read_at,
        ),
      );
    } catch (error) {
      console.error("Failed to refresh notifications:", error);
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  }, []);

  const upsertNotificationFromPush = useCallback((payload) => {
    const nextNotification = normalizeIncomingNotification(
      payload,
      payload?.notification?.title || "New notification",
    );

    setNotifications((current) => {
      const existingIndex = current.findIndex(
        (notification) => notification.id === nextNotification.id,
      );

      if (existingIndex >= 0) {
        const cloned = [...current];
        cloned[existingIndex] = {
          ...cloned[existingIndex],
          ...nextNotification,
          channels: {
            ...cloned[existingIndex].channels,
            ...nextNotification.channels,
            in_app: {
              ...cloned[existingIndex].channels?.in_app,
              ...nextNotification.channels.in_app,
            },
          },
        };
        return cloned;
      }

      return [nextNotification, ...current].slice(0, 20);
    });

    setHasNewNotifications(true);
  }, []);

  const markAllRead = useCallback(async () => {
    try {
      await notificationsAPI.markAllRead();
      setNotifications((current) =>
        current.map((notification) => ({
          ...notification,
          channels: {
            ...notification.channels,
            in_app: {
              ...notification.channels?.in_app,
              read_at:
                notification.channels?.in_app?.read_at ||
                new Date().toISOString(),
            },
          },
        })),
      );
      setHasNewNotifications(false);
    } catch (error) {
      console.error("Failed to mark notifications as read:", error);
    }
  }, []);

  const openDrawer = useCallback(async () => {
    setDrawerOpen(true);
    await markAllRead();
  }, [markAllRead]);

  const closeDrawer = useCallback(() => {
    setDrawerOpen(false);
  }, []);

  useEffect(() => {
    refreshNotifications({ silent: false });

    return () => {
      toastTimersRef.current.forEach((timer) => clearTimeout(timer));
      toastTimersRef.current.clear();
    };
  }, [refreshNotifications]);

  useEffect(() => {
    const handleForegroundNotification = async (event) => {
      console.log(
        "[Notifications] foreground notification event",
        event.detail,
      );
      const detail = normalizeForegroundPayload(event.detail);
      upsertNotificationFromPush(event.detail?.raw || event.detail);
      pushToast(detail.title, detail.body || detail.title);
      await refreshNotifications({ silent: true });
    };

    window.addEventListener(
      "app:foreground-notification",
      handleForegroundNotification,
    );
    return () => {
      window.removeEventListener(
        "app:foreground-notification",
        handleForegroundNotification,
      );
    };
  }, [pushToast, refreshNotifications, upsertNotificationFromPush]);

  useEffect(() => {
    const handleServiceWorkerMessage = async (event) => {
      if (event?.data?.type !== "FCM_BACKGROUND_MESSAGE") {
        return;
      }

      console.log(
        "[Notifications] service worker message received",
        event.data,
      );

      upsertNotificationFromPush(event.data?.payload || event.data);

      await refreshNotifications({ silent: true });
    };

    navigator.serviceWorker?.addEventListener(
      "message",
      handleServiceWorkerMessage,
    );

    return () => {
      navigator.serviceWorker?.removeEventListener(
        "message",
        handleServiceWorkerMessage,
      );
    };
  }, [refreshNotifications, upsertNotificationFromPush]);

  useEffect(() => {
    const handleVisibilityChange = () => {
      if (document.visibilityState === "visible") {
        console.log(
          "[Notifications] visibilitychange -> visible, refreshing inbox",
        );
        refreshNotifications({ silent: true });
      }
    };

    const handleWindowFocus = () => {
      console.log("[Notifications] window focus, refreshing inbox");
      refreshNotifications({ silent: true });
    };

    document.addEventListener("visibilitychange", handleVisibilityChange);
    window.addEventListener("focus", handleWindowFocus);

    return () => {
      document.removeEventListener("visibilitychange", handleVisibilityChange);
      window.removeEventListener("focus", handleWindowFocus);
    };
  }, [refreshNotifications]);

  const value = useMemo(
    () => ({
      notifications,
      loading,
      refreshing,
      drawerOpen,
      hasNewNotifications,
      toasts,
      refreshNotifications,
      markAllRead,
      openDrawer,
      closeDrawer,
      dismissToast,
    }),
    [
      notifications,
      loading,
      refreshing,
      drawerOpen,
      hasNewNotifications,
      toasts,
      refreshNotifications,
      markAllRead,
      openDrawer,
      closeDrawer,
      dismissToast,
    ],
  );

  return (
    <NotificationContext.Provider value={value}>
      {children}
    </NotificationContext.Provider>
  );
}

export function useNotifications() {
  const context = useContext(NotificationContext);
  if (!context) {
    throw new Error(
      "useNotifications must be used within NotificationProvider",
    );
  }
  return context;
}
