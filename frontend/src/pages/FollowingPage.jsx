import { useState, useEffect } from 'react';
import { usersAPI } from '../services/api.js';
import Navbar from '../components/Navbar.jsx';

export default function FollowingPage() {
  const [following, setFollowing] = useState([]);
  const [loading, setLoading] = useState(true);
  const [actionLoading, setActionLoading] = useState(null);

  useEffect(() => {
    loadFollowing();
  }, []);

  const loadFollowing = () => {
    setLoading(true);
    usersAPI.getFollowing()
      .then((data) => {
        setFollowing(Array.isArray(data) ? data : []);
      })
      .catch((err) => console.error('Failed to load following:', err))
      .finally(() => setLoading(false));
  };

  const handleUnfollow = async (userId) => {
    setActionLoading(userId);
    try {
      await usersAPI.unfollow(userId);
      setFollowing((prev) => prev.filter((u) => u.id !== userId));
    } catch (err) {
      console.error('Unfollow failed:', err);
    } finally {
      setActionLoading(null);
    }
  };

  return (
    <>
      <Navbar title="Following" />
      <div className="page-container">
        <h1 className="page-title">Following</h1>
        <p className="page-subtitle">People you follow and receive notifications from</p>

        {loading ? (
          <div className="loading-spinner">
            <div className="spinner"></div>
            Loading following list...
          </div>
        ) : following.length === 0 ? (
          <div className="empty-state">
            <div className="empty-state-icon">👥</div>
            <div className="empty-state-text">You are not following anyone yet.</div>
          </div>
        ) : (
          <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
            {following.map((u) => (
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
                <button
                  className="btn btn-secondary btn-sm"
                  onClick={() => handleUnfollow(u.id)}
                  disabled={actionLoading === u.id}
                >
                  {actionLoading === u.id ? '...' : 'Unfollow'}
                </button>
              </div>
            ))}
          </div>
        )}
      </div>
    </>
  );
}
