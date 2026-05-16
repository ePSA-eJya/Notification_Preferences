import { useNavigate } from "react-router-dom";
import NotificationItem from "./NotificationItem.jsx";
import { useNotifications } from "../context/NotificationContext.jsx";

export default function NotificationInbox({
  title = "Notifications",
  subtitle = "Live updates from the app",
}) {
  const { notifications, loading, markAllRead, closeDrawer } =
    useNotifications();
  const navigate = useNavigate();

  const handleNotificationClick = async (notification) => {
    await markAllRead();

    if (notification.entity_type === "POST" && notification.entity_id) {
      closeDrawer();
      navigate(`/?postId=${notification.entity_id}`);
      return;
    }

    closeDrawer();
  };

  return (
    <section className="notification-inbox">
      <div className="notification-inbox-header">
        <div>
          <h2>{title}</h2>
          <p>{subtitle}</p>
        </div>
      </div>

      {loading ? (
        <div className="loading-spinner notification-inbox-loading">
          <div className="spinner"></div>
          Loading notifications...
        </div>
      ) : notifications.length === 0 ? (
        <div className="notification-inbox-empty">
          <div className="empty-state-icon">
            <i className="fas fa-bell"></i>
          </div>
          <div className="empty-state-text">
            No notifications yet. New activity will appear here instantly.
          </div>
        </div>
      ) : (
        <div className="notif-list">
          {notifications.map((notification) => (
            <NotificationItem
              key={notification.id}
              notification={notification}
              unread={!notification.channels?.in_app?.read_at}
              onClick={() => handleNotificationClick(notification)}
            />
          ))}
        </div>
      )}
    </section>
  );
}
