import NotificationInbox from "./NotificationInbox.jsx";
import { useNotifications } from "../context/NotificationContext.jsx";

export default function NotificationDrawer() {
  const { drawerOpen, closeDrawer } = useNotifications();

  if (!drawerOpen) {
    return null;
  }

  return (
    <div className="notification-drawer-overlay" onClick={closeDrawer}>
      <aside
        className="notification-drawer"
        onClick={(event) => event.stopPropagation()}
      >
        <button
          className="notification-drawer-close"
          type="button"
          onClick={closeDrawer}
          aria-label="Close notifications"
        >
          <i className="fa fa-times" aria-hidden="true"></i>
        </button>
        <NotificationInbox
          title="Notifications"
          subtitle="Your live inbox and push activity"
        />
      </aside>
    </div>
  );
}
