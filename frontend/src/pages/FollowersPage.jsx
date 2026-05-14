import { useState, useEffect } from 'react';
import { usersAPI } from '../services/api.js';
import Navbar from '../components/Navbar.jsx';

export default function FollowersPage() {
  const [followers, setFollowers] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    usersAPI.getFollowers()
      .then((data) => {
        setFollowers(Array.isArray(data) ? data : []);
      })
      .catch((err) => console.error('Failed to load followers:', err))
      .finally(() => setLoading(false));
  }, []);

  return (
    <>
      <Navbar title="My Followers" />
      <div className="page-container">
        {/* <h1 className="page-title">My Followers</h1> */}
        <p className="page-subtitle">People who follow you and receive your post notifications</p>

        {loading ? (
          <div className="loading-spinner">
            <div className="spinner"></div>
            Loading followers...
          </div>
        ) : followers.length === 0 ? (
          <div className="empty-state">
            <div className="empty-state-icon"><i className="fas fa-users"></i></div>
            <div className="empty-state-text">You don't have any followers yet.</div>
          </div>
        ) : (
          <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
            {followers.map((u) => (
              <div key={u.id} className="user-card fade-in">
                <div className="user-info">
                  <div className="user-avatar">
                    {u.user_handle?.[0]?.toUpperCase() || '?'}
                  </div>
                  <div>
                    <div className="user-name">@{u.user_handle}</div>
                    <div className="user-email">{u.email}</div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </>
  );
}
