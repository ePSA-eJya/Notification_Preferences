import { Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from './context/AuthContext.jsx';
import DashboardPage from './pages/DashboardPage.jsx';
import LoginPage from './pages/LoginPage.jsx';
import SignupPage from './pages/SignupPage.jsx';
import NotificationsPage from './pages/NotificationsPage.jsx';
import PreferencesPage from './pages/PreferencesPage.jsx';
import UsersPage from './pages/UsersPage.jsx';
import ProtectedRoute from './components/ProtectedRoute.jsx';
import Layout from './components/Layout.jsx';
import FollowersPage from './pages/FollowersPage.jsx';
import FollowingPage from './pages/FollowingPage.jsx';
import { useEffect } from 'react';
import { getDeviceToken, listenMessages } from './firebase.js';

export default function App() {
  const { token, loading } = useAuth();
  useEffect(() => {
    // getDeviceToken().then(token => {
    //   if (token) {
    //     console.log("Send this token to backend:", token);
    //   }
    // });

    listenMessages();
  }, []);
  if (loading) {
    return (
      <div className="loading-spinner" style={{ height: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <div className="spinner"></div>
      </div>
    );
  }


  return (
    <Routes>
      {/* Public Routes */}
      <Route path="/login" element={token ? <Navigate to="/" replace /> : <LoginPage />} />
      <Route path="/signup" element={token ? <Navigate to="/" replace /> : <SignupPage />} />

      {/* Protected Routes wrapped with Layout */}
      <Route element={<Layout />}>
        <Route
          path="/"
          element={
            <ProtectedRoute>
              <DashboardPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/notifications"
          element={
            <ProtectedRoute>
              <NotificationsPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/preferences"
          element={
            <ProtectedRoute>
              <PreferencesPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/users"
          element={
            <ProtectedRoute>
              <UsersPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/followers"
          element={
            <ProtectedRoute>
              <FollowersPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/following"
          element={
            <ProtectedRoute>
              <FollowingPage />
            </ProtectedRoute>
          }
        />
      </Route>

      {/* Fallback */}
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
