import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext.jsx';

export default function SignupPage() {
  const { signup } = useAuth();
  const navigate = useNavigate();

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [userHandle, setUserHandle] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [showSuccess, setShowSuccess] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await signup(email, password, userHandle);
      setShowSuccess(true);
      setTimeout(() => {
        navigate('/login');
      }, 2000);
    } catch (err) {
      setError(err.message || 'Signup failed');
      setLoading(false);
    }
  };

  return (
    <div className="auth-page">
      <div className="auth-card fade-in">
        <div className="card">
          <div className="auth-header">
            <div className="auth-logo"><i className="fas fa-bell"></i></div>
            <h1>Create Account</h1>
            <p>Join and customise your notification experience</p>
          </div>

          {error && <div className="form-error">{error}</div>}

          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label className="form-label" htmlFor="signup-handle">Username</label>
              <input
                id="signup-handle"
                className="form-input"
                type="text"
                placeholder="johndoe"
                value={userHandle}
                onChange={(e) => setUserHandle(e.target.value)}
                required
              />
            </div>

            <div className="form-group">
              <label className="form-label" htmlFor="signup-email">Email</label>
              <input
                id="signup-email"
                className="form-input"
                type="email"
                placeholder="you@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>

            <div className="form-group">
              <label className="form-label" htmlFor="signup-password">Password</label>
              <input
                id="signup-password"
                className="form-input"
                type="password"
                placeholder="min 6 characters"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                minLength={6}
              />
            </div>

            <button
              type="submit"
              className="btn btn-primary"
              disabled={loading}
              style={{ width: '100%' }}
            >
              {loading ? 'Creating account...' : 'Create Account'}
            </button>
          </form>

          <div className="auth-footer">
            Already have an account? <Link to="/login">Sign in</Link>
          </div>
        </div>
      </div>
      {showSuccess && (
        <div className="comments-modal-overlay">
          <div className="card fade-in" style={{ padding: '40px', width: '90%', maxWidth: '400px', display: 'flex', flexDirection: 'column', alignItems: 'center', textAlign: 'center', gap: '16px' }}>
            <div style={{ fontSize: '4rem', color: 'var(--success)' }}>
              <i className="fas fa-check-circle"></i>
            </div>
            <h2 style={{ fontSize: '1.5rem', fontWeight: '700' }}>Account Created!</h2>
            <p style={{ color: 'var(--text-secondary)' }}>Your account has been created successfully. Redirecting to login...</p>
          </div>
        </div>
      )}
    </div>
  );
}
