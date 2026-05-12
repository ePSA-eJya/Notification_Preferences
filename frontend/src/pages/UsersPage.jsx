import { useState, useEffect } from 'react';
import { usersAPI } from '../services/api.js';
import { useAuth } from '../context/AuthContext.jsx';
import Navbar from '../components/Navbar.jsx';

export default function UsersPage() {
  const { user: currentUser } = useAuth();
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [followingSet, setFollowingSet] = useState(new Set());
  const [actionLoading, setActionLoading] = useState(null);

  useEffect(() => {
    // Fetch all users and current user's following list
    Promise.all([
      usersAPI.getAll(),
      usersAPI.getFollowing()
    ])
      .then(([users, following]) => {
        const list = Array.isArray(users) ? users : [];
        setUsers(list);
        
        // Convert following array to a Set of IDs for O(1) lookup
        const followingIds = new Set(following?.map(u => u.id) || []);
        setFollowingSet(followingIds);
      })
      .catch((err) => console.error('Failed to load users:', err))
      .finally(() => setLoading(false));
  }, []);

  const handleFollow = async (userId) => {
    setActionLoading(userId);
    try {
      await usersAPI.follow(userId);
      setFollowingSet((prev) => new Set([...prev, userId]));
    } catch (err) {
      // If already following, mark as following anyway
      if (err.message?.includes('already')) {
        setFollowingSet((prev) => new Set([...prev, userId]));
      }
      console.error('Follow failed:', err);
    } finally {
      setActionLoading(null);
    }
  };

  const handleUnfollow = async (userId) => {
    setActionLoading(userId);
    try {
      await usersAPI.unfollow(userId);
      setFollowingSet((prev) => {
        const next = new Set(prev);
        next.delete(userId);
        return next;
      });
    } catch (err) {
      console.error('Unfollow failed:', err);
    } finally {
      setActionLoading(null);
    }
  };

  return (
    <>
      <Navbar title="Discover Users" />
      <div className="page-container">
        <h1 className="page-title">Discover Users</h1>
        <p className="page-subtitle">Find and follow people to see their posts in your feed</p>

        {loading ? (
          <div className="loading-spinner">
            <div className="spinner"></div>
            Loading users...
          </div>
        ) : users.length === 0 ? (
          <div className="empty-state">
            <div className="empty-state-icon">👥</div>
            <div className="empty-state-text">No users found.</div>
          </div>
        ) : (
          <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
            {users
              .filter((u) => u.id !== currentUser?.id)
              .map((u) => {
                const isFollowing = followingSet.has(u.id);
                const isLoading = actionLoading === u.id;

                return (
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

                    {isFollowing ? (
                      <button
                        className="btn btn-secondary btn-sm"
                        onClick={() => handleUnfollow(u.id)}
                        disabled={isLoading}
                      >
                        {isLoading ? '...' : '✓ Following'}
                      </button>
                    ) : (
                      <button
                        className="btn btn-primary btn-sm"
                        onClick={() => handleFollow(u.id)}
                        disabled={isLoading}
                      >
                        {isLoading ? '...' : '+ Follow'}
                      </button>
                    )}
                  </div>
                );
              })}
          </div>
        )}
      </div>
    </>
  );
}
