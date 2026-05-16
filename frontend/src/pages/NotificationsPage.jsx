import { useEffect } from 'react';
import Navbar from '../components/Navbar.jsx';
import NotificationInbox from '../components/NotificationInbox.jsx';
import { useNotifications } from '../context/NotificationContext.jsx';

export default function NotificationsPage() {
  const { refreshNotifications } = useNotifications();

  useEffect(() => {
    refreshNotifications({ silent: true, markSeen: true });
  }, [refreshNotifications]);

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
            Use the bell icon to open the slide-out inbox. This page mirrors the same live feed for direct access.
          </p>
        </div>

        <NotificationInbox
          title="Notifications"
          subtitle="Auto-updating inbox and rich push activity"
          showRefresh={false}
        />
      </div>
    </>
  );
}
