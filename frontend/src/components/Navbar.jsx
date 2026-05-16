import NotificationDrawer from "./NotificationDrawer.jsx";
import NotificationToastStack from "./NotificationToastStack.jsx";
import { useNotifications } from "../context/NotificationContext.jsx";

export default function Navbar({ title }) {
  const { notifications, hasNewNotifications, openDrawer } = useNotifications();

  return (
    <>
      <header className="navbar">
        <h1 className="navbar-title">{title}</h1>
        <div className="navbar-actions">
          <button
            type="button"
            className={`navbar-notification-btn ${hasNewNotifications ? "has-new" : ""}`}
            onClick={openDrawer}
            aria-label="Open notifications"
          >
            <i className="fas fa-bell"></i>
            <span className="navbar-notification-count">
              {notifications.length}
            </span>
            {hasNewNotifications && (
              <span className="navbar-notification-dot" />
            )}
          </button>
        </div>
      </header>
      <NotificationDrawer />
      <NotificationToastStack />
    </>
  );
}
