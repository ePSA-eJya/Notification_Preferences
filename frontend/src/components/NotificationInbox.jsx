import NotificationItem from './NotificationItem.jsx';
import { useNotifications } from '../context/NotificationContext.jsx';

export default function NotificationInbox({ title = 'Notifications', subtitle = 'Live updates from the app', showRefresh = true }) {
  const {
    notifications,
    loading,
    refreshing,
    hasNewNotifications,
    refreshNotifications,
  } = useNotifications();

  const handleRefresh = async () => {
    await refreshNotifications({ silent: false, markSeen: true });
  };

  return (
    <section className="notification-inbox">
      <div className="notification-inbox-header">
        <div>
          <h2>{title}</h2>
          <p>{subtitle}</p>
        </div>
        {showRefresh && (
          <button className="btn btn-secondary btn-sm" type="button" onClick={handleRefresh} disabled={loading || refreshing}>
            <i className="fa fa-refresh" aria-hidden="true"></i>
          </button>
        )}
      </div>

      <div className="notification-inbox-meta">
        <span>{hasNewNotifications ? 'New activity available' : `${notifications.length} notifications`}</span>
        <span>{refreshing ? 'Refreshing...' : 'Auto-updates on'}</span>
      </div>

      {loading ? (
        <div className="loading-spinner notification-inbox-loading">
          <div className="spinner"></div>
          Loading notifications...
        </div>
      ) : notifications.length === 0 ? (
        <div className="notification-inbox-empty">
          <div className="empty-state-icon"><i className="fas fa-bell"></i></div>
          <div className="empty-state-text">No notifications yet. New activity will appear here instantly.</div>
        </div>
      ) : (
        <div className="notif-list">
          {notifications.map((notification) => (
            <NotificationItem key={notification.id} notification={notification} />
          ))}
        </div>
      )}
    </section>
  );
}