import { initializeApp } from "firebase/app";
import { getMessaging, getToken, onMessage } from "firebase/messaging";

const firebaseConfig = {
  apiKey: "AIzaSyDmDoBAvQ-G8JGYlkxXKKHCME7rMsLvxjE",
  authDomain: "notification-preferences-46b33.firebaseapp.com",
  projectId: "notification-preferences-46b33",
  storageBucket: "notification-preferences-46b33.firebasestorage.app",
  messagingSenderId: "1049542857070",
  appId: "1:1049542857070:web:538ab9e3b78ddbb106f17f",
  measurementId: "G-K2M6ELZ5D5",
};

const app = initializeApp(firebaseConfig);
const messaging = getMessaging(app);

async function registerMessagingServiceWorker() {
  if (!("serviceWorker" in navigator)) {
    console.log("[FCM] service workers are not supported in this browser");
    return null;
  }

  console.log("[FCM] registering service worker", {
    controller: Boolean(navigator.serviceWorker.controller),
    visibilityState: document.visibilityState,
    notificationPermission: Notification.permission,
  });

  const registration = await navigator.serviceWorker.register(
    "/firebase-messaging-sw.js",
  );

  console.log("[FCM] service worker registered", registration?.scope);

  return registration;
}

// 🔵 Get token
export async function getDeviceToken() {
  console.log("[FCM] requesting notification permission");
  const permission = await Notification.requestPermission();
  console.log("[FCM] notification permission result", permission);

  if (permission !== "granted") {
    throw new Error("Notification permission denied");
  }

  const serviceWorkerRegistration = await registerMessagingServiceWorker();

  const token = await getToken(messaging, {
    vapidKey:
      "BO5YWWCcmLnH7qegP_elb4svtD1oQ3xSpU_vqQDhhNvfSGv78aWiw71ucWqO8XTmi6RQBwheu0l-xXYEFAjDbr0",
    serviceWorkerRegistration: serviceWorkerRegistration || undefined,
  });

  console.log("[FCM] device token acquired", token);

  return token;
}

// 🔵 Foreground listener
export function listenMessages() {
  console.log("[FCM] foreground listener registered");

  if (navigator.serviceWorker?.controller) {
    console.log("[FCM] service worker controller available");
  } else {
    console.log("[FCM] no active service worker controller yet");
  }

  if (navigator.permissions?.query) {
    navigator.permissions
      .query({ name: "notifications" })
      .then((status) => {
        console.log("[FCM] notifications permission probe", {
          state: status.state,
          onchange: typeof status.onchange,
        });
      })
      .catch((error) => {
        console.error("[FCM] notifications permission probe failed", error);
      });
  }

  return onMessage(messaging, (payload) => {
    console.log("[FCM] foreground message received", payload);

    const detail = {
      title: payload?.notification?.title || "New notification",
      body: payload?.notification?.body || "",
      raw: payload,
    };

    window.dispatchEvent(
      new CustomEvent("app:foreground-notification", {
        detail,
      }),
    );

    console.log("[FCM] notification permission state", Notification.permission);
    console.log("[FCM] document visibility state", document.visibilityState);

    if (Notification.permission === "granted") {
      try {
        navigator.serviceWorker
          ?.getRegistration()
          .then((registration) => {
            if (registration) {
              console.log(
                "[FCM] showing system notification via service worker",
              );
              return registration.showNotification(detail.title, {
                body: detail.body,
                tag: payload?.data?.notification_id || undefined,
                data: payload?.data || {},
              });
            }

            console.log(
              "[FCM] no service worker registration found, falling back to Notification API",
            );
            return new Notification(detail.title, {
              body: detail.body,
            });
          })
          .then(() => {
            console.log("[FCM] system notification shown in foreground");
          })
          .catch((error) => {
            console.error("[FCM] failed to show system notification", error);
          });
      } catch (error) {
        console.error("[FCM] failed to show system notification", error);
      }
    }
  });
}
