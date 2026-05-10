import { useState, useEffect } from 'react';
import { feedAPI } from '../services/api.js';
import { useAuth } from '../context/AuthContext.jsx';
import Navbar from '../components/Navbar.jsx';
import PostCard from '../components/PostCard.jsx';

export default function DashboardPage() {
  const { user } = useAuth();
  const [posts, setPosts] = useState([]);
  const [content, setContent] = useState('');
  const [loading, setLoading] = useState(true);
  const [posting, setPosting] = useState(false);

  const loadFeed = async () => {
    try {
      const data = await feedAPI.getFeed();
      setPosts(Array.isArray(data) ? data : []);
    } catch (err) {
      console.error('Failed to load feed:', err);
      setPosts([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadFeed();
  }, []);

  const handleCreatePost = async (e) => {
    e.preventDefault();
    if (!content.trim() || posting) return;
    setPosting(true);
    try {
      const newPost = await feedAPI.createPost(content.trim());
      setPosts((prev) => [newPost, ...prev]);
      setContent('');
    } catch (err) {
      console.error('Post creation failed:', err);
    } finally {
      setPosting(false);
    }
  };

  return (
    <>
      <Navbar title="Feed" />
      <div className="page-container">
        <h1 className="page-title">Your Feed</h1>
        <p className="page-subtitle">See what people you follow are posting</p>

        {/* Create Post */}
        <form className="create-post card" onSubmit={handleCreatePost}>
          <textarea
            className="form-input"
            placeholder="What's on your mind?"
            value={content}
            onChange={(e) => setContent(e.target.value)}
            rows={3}
          />
          <button
            type="submit"
            className="btn btn-primary"
            disabled={posting || !content.trim()}
          >
            {posting ? 'Posting...' : '📝 Publish Post'}
          </button>
        </form>

        {/* Feed */}
        {loading ? (
          <div className="loading-spinner">
            <div className="spinner"></div>
            Loading feed...
          </div>
        ) : posts.length === 0 ? (
          <div className="empty-state">
            <div className="empty-state-icon">📰</div>
            <div className="empty-state-text">
              No posts yet. Follow some users or create your first post!
            </div>
          </div>
        ) : (
          <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
            {posts.map((post) => (
              <PostCard key={post.id} post={post} currentUserId={user?.id} />
            ))}
          </div>
        )}
      </div>
    </>
  );
}
