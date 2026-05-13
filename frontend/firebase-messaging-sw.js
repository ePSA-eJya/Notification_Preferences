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

messaging.onBackgroundMessage((payload) => {

    console.log(
        'Background message received:',
        payload
    );

    self.registration.showNotification(
        payload.notification.title,
        {
            body: payload.notification.body,
        }
    );
});
