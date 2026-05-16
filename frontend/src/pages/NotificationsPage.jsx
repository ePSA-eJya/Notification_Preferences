import { useEffect } from "react";
import Navbar from "../components/Navbar.jsx";
import NotificationInbox from "../components/NotificationInbox.jsx";
import { useNotifications } from "../context/NotificationContext.jsx";

export default function NotificationsPage() {
  const { markAllRead } = useNotifications();

  useEffect(() => {
    markAllRead();
  }, [markAllRead]);

  return (
    <>
      <Navbar title="Notifications" />
      <div className="page-container notification-page-shell">
        <div className="notification-page-intro">
          <div>
            <p className="page-subtitle">Your live inbox</p>
            <h2>Notifications are now anchored to the top bar</h2>
          </div>
          <p className="notification-page-copy">
            Use the bell icon to open the slide-out inbox. This page mirrors the
            same live feed for direct access.
          </p>
        </div>

        <NotificationInbox
          title="Notifications"
          subtitle="Rich push activity and in-app inbox"
        />
      </div>
    </>
  );
}
