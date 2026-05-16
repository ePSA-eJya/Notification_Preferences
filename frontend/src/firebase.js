import { initializeApp } from "firebase/app";
import {
    getMessaging,
    getToken,
    onMessage
} from "firebase/messaging";

const firebaseConfig = {
    apiKey: "AIzaSyDmDoBAvQ-G8JGYlkxXKKHCME7rMsLvxjE",
    authDomain: "notification-preferences-46b33.firebaseapp.com",
    projectId: "notification-preferences-46b33",
    storageBucket: "notification-preferences-46b33.firebasestorage.app",
    messagingSenderId: "1049542857070",
    appId: "1:1049542857070:web:538ab9e3b78ddbb106f17f",
    measurementId: "G-K2M6ELZ5D5"
};

const app = initializeApp(firebaseConfig);
const messaging = getMessaging(app);

// 🔵 Get token
export async function getDeviceToken() {
    const permission = await Notification.requestPermission();

    if (permission !== "granted") {
        throw new Error("Notification permission denied");
    }

    const token = await getToken(messaging, {
        vapidKey: "BO5YWWCcmLnH7qegP_elb4svtD1oQ3xSpU_vqQDhhNvfSGv78aWiw71ucWqO8XTmi6RQBwheu0l-xXYEFAjDbr0"
    });

    return token;
}

// 🔵 Foreground listener
export function listenMessages() {
    onMessage(messaging, (payload) => {
        const detail = {
            title: payload?.notification?.title || 'New notification',
            body: payload?.notification?.body || '',
            raw: payload,
        };

        window.dispatchEvent(new CustomEvent('app:foreground-notification', {
            detail,
        }));

        if (document.hidden && Notification.permission === 'granted') {
            new Notification(detail.title, {
                body: detail.body,
            });
        }
    });
}