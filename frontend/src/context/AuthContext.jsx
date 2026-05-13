import { createContext, useContext, useState, useEffect } from 'react';
import { authAPI, usersAPI } from '../services/api.js';

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [token, setToken] = useState(() => localStorage.getItem('token'));
  const [loading, setLoading] = useState(true);

  // On mount, if token exists, fetch user info
  useEffect(() => {
    if (token) {
      usersAPI.getMe()
        .then((data) => setUser(data))
        .catch(() => {
          // Token expired or invalid
          localStorage.removeItem('token');
          setToken(null);
        })
        .finally(() => setLoading(false));
    } else {
      setLoading(false);
    }
  }, [token]);

  const login = async (email, password, deviceToken) => {
    const data = await authAPI.signin(email, password, deviceToken);
    localStorage.setItem('token', data.token);
    setToken(data.token);
    setUser(data.user);
    return data;
  };

  const signup = async (email, password, userHandle, deviceToken) => {
    const data = await authAPI.signup(email, password, userHandle, deviceToken);
    return data;
  };

  const logout = () => {
    localStorage.removeItem('token');
    setToken(null);
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, token, loading, login, signup, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}
