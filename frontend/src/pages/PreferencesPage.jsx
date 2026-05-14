import { useState, useEffect } from 'react';
import { preferencesAPI } from '../services/api.js';
import Navbar from '../components/Navbar.jsx';
import PreferenceToggle from '../components/PreferenceToggle.jsx';

const actionTypes = [
  { key: 'likes', label: 'Likes', icon: 'fas fa-heart' },
  { key: 'comments', label: 'Comments', icon: 'fas fa-comment' },
  { key: 'follows', label: 'Follows', icon: 'fas fa-user' },
  { key: 'posts', label: 'Posts', icon: 'fas fa-edit' },
];

const channels = ['in_app', 'push', 'email'];
const channelLabels = { in_app: 'In-App', push: 'Push', email: 'Email' };

const defaultPrefs = {
  likes: { in_app: 'ALL', push: 'ALL', email: 'NONE' },
  comments: { in_app: 'ALL', push: 'ALL', email: 'NONE' },
  follows: { in_app: 'ALL', push: 'ALL', email: 'NONE' },
  posts: { in_app: 'ALL', push: 'ALL', email: 'NONE' },
};

export default function PreferencesPage() {
  const [prefs, setPrefs] = useState(defaultPrefs);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);

  useEffect(() => {
    preferencesAPI.get()
      .then((data) => {
        if (data?.preferences) {
          setPrefs(data.preferences);
        }
      })
      .catch((err) => console.error('Failed to load preferences:', err))
      .finally(() => setLoading(false));
  }, []);

  const handleChange = (action, channel, value) => {
    setPrefs((prev) => ({
      ...prev,
      [action]: {
        ...prev[action],
        [channel]: value,
      },
    }));
    setSaved(false);
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      await preferencesAPI.update(prefs);
      setSaved(true);
      setTimeout(() => setSaved(false), 3000);
    } catch (err) {
      console.error('Failed to save preferences:', err);
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <>
        <Navbar title="Preferences" />
        <div className="page-container">
          <div className="loading-spinner">
            <div className="spinner"></div>
            Loading preferences...
          </div>
        </div>
      </>
    );
  }

  return (
    <>
      <Navbar title="Preferences" />
      <div className="page-container">
        {/* <h1 className="page-title">Notification Preferences</h1> */}
        <p className="page-subtitle">
          Control how and when you receive notifications for each activity type
        </p>

        <div className="card" style={{ padding: 0, overflow: 'hidden' }}>
          <table className="pref-table">
            <thead>
              <tr>
                <th>Activity</th>
                {channels.map((ch) => (
                  <th key={ch}>{channelLabels[ch]}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {actionTypes.map((action) => (
                <tr key={action.key}>
                  <td>
                    <div className="pref-row-label">
                      <span className="pref-row-icon"><i className={action.icon}></i></span>
                      {action.label}
                    </div>
                  </td>
                  {channels.map((ch) => (
                    <td key={ch}>
                      <PreferenceToggle
                        value={prefs[action.key]?.[ch]}
                        onChange={(val) => handleChange(action.key, ch, val)}
                      />
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        <div style={{ marginTop: 24, display: 'flex', alignItems: 'center', gap: 16 }}>
          <button
            className="btn btn-primary"
            onClick={handleSave}
            disabled={saving}
          >
            {saving ? 'Saving...' : <> Save Preferences</>}
          </button>
          {saved && (
            <span className="badge badge-success fade-in"><i className="fas fa-check"></i> Saved successfully</span>
          )}
        </div>
      </div>
    </>
  );
}
