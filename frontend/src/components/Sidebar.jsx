import { NavLink, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext.jsx';

const navItems = [
  { path: '/', label: 'Feed', icon: '📰' },
  { path: '/users', label: 'Discover', icon: '👥' },
  { path: '/notifications', label: 'Notifications', icon: '🔔' },
  { path: '/preferences', label: 'Preferences', icon: '⚙️' },
];

export default function Sidebar() {
  const { user, logout } = useAuth();
  const location = useLocation();

  return (
    <aside className="sidebar">
      <div className="sidebar-brand">
        <span className="sidebar-brand-icon">🔔</span>
        <span className="sidebar-brand-text">NotifPref</span>
      </div>

      <nav className="sidebar-nav">
        {navItems.map((item) => (
          <NavLink
            key={item.path}
            to={item.path}
            className={`sidebar-link ${location.pathname === item.path ? 'active' : ''}`}
          >
            <span className="sidebar-link-icon">{item.icon}</span>
            {item.label}
          </NavLink>
        ))}
      </nav>

      <div className="sidebar-footer">
        <div className="sidebar-user">
          <div className="sidebar-user-avatar">
            {user?.user_handle?.[0]?.toUpperCase() || '?'}
          </div>
          <div className="sidebar-user-info">
            <div className="sidebar-user-name">{user?.user_handle || 'User'}</div>
            <div className="sidebar-user-email">{user?.email || ''}</div>
          </div>
        </div>
        <button className="btn btn-ghost btn-sm" onClick={logout} style={{ width: '100%', marginTop: 8 }}>
          🚪 Logout
        </button>
      </div>
    </aside>
  );
}
