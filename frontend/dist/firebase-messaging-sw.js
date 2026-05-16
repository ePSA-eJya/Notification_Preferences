importScripts(
    'https://www.gstatic.com/firebasejs/10.12.2/firebase-app-compat.js'
);

importScripts(
    'https://www.gstatic.com/firebasejs/10.12.2/firebase-messaging-compat.js'
);

firebase.initializeApp({
    apiKey: "AIzaSyDmDoBAvQ-G8JGYlkxXKKHCME7rMsLvxjE",
    authDomain: "notification-preferences-46b33.firebaseapp.com",
    projectId: "notification-preferences-46b33",
    storageBucket: "notification-preferences-46b33.firebasestorage.app",
    messagingSenderId: "1049542857070",
    appId: "1:1049542857070:web:538ab9e3b78ddbb106f17f",
    measurementId: "G-K2M6ELZ5D5"
});

const messaging = firebase.messaging();

self.addEventListener('install', (event) => {
    console.log('Service worker install event');
    self.skipWaiting();
});

self.addEventListener('activate', (event) => {
    console.log('Service worker activate event');
    event.waitUntil(self.clients.claim());
});

self.addEventListener('notificationclick', (event) => {
    console.log('Notification click event', event.notification?.data || null);
    event.notification.close();

    event.waitUntil(
        self.clients.matchAll({
            type: 'window',
            includeUncontrolled: true,
        }).then((clients) => {
            console.log('Notification click matched clients:', clients.length);

            if (clients.length > 0) {
                const client = clients[0];
                console.log('Focusing existing client:', client.url);
                return client.focus();
            }

            console.log('Opening new client from notification click');
            return self.clients.openWindow('/');
        })
    );
});

messaging.onBackgroundMessage((payload) => {

    console.log(
        'Background message received:',
        payload
    );
    console.log('Background notification permission state:', Notification.permission);
    console.log('Background showNotification available:', typeof self.registration?.showNotification);

    self.clients.matchAll({
        type: 'window',
        includeUncontrolled: true,
    }).then((clients) => {
        console.log('Background message matched clients:', clients.length);

        clients.forEach((client) => {
            console.log('Posting background message to client:', client.url);
            client.postMessage({
                type: 'FCM_BACKGROUND_MESSAGE',
                payload,
            });
        });
    });

    try {
        self.registration.showNotification(
            payload.notification.title,
            {
                body: payload.notification.body,
                data: payload.data || {},
            }
        ).then(() => {
            console.log('Background notification shown');
        }).catch((error) => {
            console.error('Background notification showNotification failed', error);
        });
    } catch (error) {
        console.error('Background notification threw synchronously', error);
    }
});
