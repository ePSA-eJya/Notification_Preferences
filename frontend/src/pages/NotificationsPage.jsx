import { useState, useEffect } from 'react';
import { notificationsAPI } from '../services/api.js';
import Navbar from '../components/Navbar.jsx';
import NotificationItem from '../components/NotificationItem.jsx';

export default function NotificationsPage() {
  const [notifications, setNotifications] = useState([]);
  const [loading, setLoading] = useState(true);

  const load = async () => {
    try {
      const data = await notificationsAPI.getAll();
      setNotifications(data?.notifications || []);
    } catch (err) {
      console.error('Failed to load notifications:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  return (
    <>
      <Navbar title="Notifications" />
      <div className="page-container">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: 8 }}>
          <div>
            <p className="page-subtitle">Your notification inbox</p>
          </div>
          <button className="btn btn-secondary btn-sm" onClick={() => { setLoading(true); load(); }}>
            <i class="fa fa-refresh" aria-hidden="true"></i>
          </button>
        </div>

        {loading ? (
          <div className="loading-spinner">
            <div className="spinner"></div>
            Loading notifications...
          </div>
        ) : notifications.length === 0 ? (
          <div className="empty-state">
            <div className="empty-state-icon"><i className="fas fa-bell"></i></div>
            <div className="empty-state-text">
              No notifications yet. Interact with other users to see notifications here!
            </div>
          </div>
        ) : (
          <div className="notif-list">
            {notifications.map((n) => (
              <NotificationItem key={n.id} notification={n} />
            ))}
          </div>
        )}
      </div>
    </>
  );
}
