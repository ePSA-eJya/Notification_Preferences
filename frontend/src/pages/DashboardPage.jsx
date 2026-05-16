import { useState, useEffect, useRef } from 'react';
import { feedAPI } from '../services/api.js';
import { useAuth } from '../context/AuthContext.jsx';
import Navbar from '../components/Navbar.jsx';
import PostCard from '../components/PostCard.jsx';

export default function DashboardPage() {
  const { user } = useAuth();
  const [posts, setPosts] = useState([]);
  const [showComposer, setShowComposer] = useState(false);
  const [content, setContent] = useState('');
  const [selectedMedia, setSelectedMedia] = useState([]);
  const [mediaPreviews, setMediaPreviews] = useState([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [showNewPosts, setShowNewPosts] = useState(false);
  const [posting, setPosting] = useState(false);
  const fileInputRef = useRef(null);
  const latestPostIdRef = useRef(null);

  const resetComposer = () => {
    setContent('');
    setSelectedMedia([]);
    setMediaPreviews([]);
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const loadFeed = async ({ silent = false, markAsSeen = false } = {}) => {
    if (!silent) {
      setRefreshing(true);
    }
    try {
      const data = await feedAPI.getFeed();
      const nextPosts = Array.isArray(data) ? data : [];
      const nextTopPostId = nextPosts[0]?.id || null;

      if (latestPostIdRef.current && nextTopPostId && nextTopPostId !== latestPostIdRef.current && !markAsSeen) {
        setShowNewPosts(true);
      }

      setPosts(nextPosts);
      latestPostIdRef.current = nextTopPostId;

      if (markAsSeen) {
        setShowNewPosts(false);
      }
    } catch (err) {
      console.error('Failed to load feed:', err);
      setPosts([]);
    } finally {
      setLoading(false);
      if (!silent) {
        setRefreshing(false);
      }
    }
  };

  useEffect(() => {
    loadFeed({ silent: false, markAsSeen: true });

    const intervalId = setInterval(() => {
      loadFeed({ silent: true });
    }, 30000);

    return () => clearInterval(intervalId);
  }, []);

  const handleMediaSelect = (e) => {
    const files = Array.from(e.target.files || []);
    const validFiles = [];

    files.forEach((file) => {
      // Validate file type
      if (!file.type.startsWith('image/') && !file.type.startsWith('video/')) {
        alert(`${file.name} is not a valid image or video`);
        return;
      }

      // Validate file size (max 50MB)
      if (file.size > 50 * 1024 * 1024) {
        alert(`${file.name} is too large (max 50MB)`);
        return;
      }

      validFiles.push(file);

      // Create preview
      const reader = new FileReader();
      reader.onload = (event) => {
        setMediaPreviews(prev => [...prev, {
          type: file.type.startsWith('image/') ? 'image' : 'video',
          url: event.target.result,
          filename: file.name,
        }]);
      };
      reader.readAsDataURL(file);
    });

    setSelectedMedia((prev) => [...prev, ...validFiles]);
    e.target.value = '';
  };

  const removeMedia = (index) => {
    setSelectedMedia(prev => prev.filter((_, i) => i !== index));
    setMediaPreviews(prev => prev.filter((_, i) => i !== index));
  };

  const handleCreatePost = async (e) => {
    e.preventDefault();
    if ((!content.trim() && selectedMedia.length === 0) || posting) return;
    setPosting(true);
    try {
      const newPost = await feedAPI.createPost(content.trim(), selectedMedia);
      setPosts((prev) => [newPost, ...prev]);
      latestPostIdRef.current = newPost.id || latestPostIdRef.current;
      setShowNewPosts(false);
      resetComposer();
      setShowComposer(false);
    } catch (err) {
      console.error('Post creation failed:', err);
      alert('Failed to create post: ' + err.message);
    } finally {
      setPosting(false);
    }
  };

  const handleOpenComposer = () => {
    setShowComposer(true);
  };

  const handleCloseComposer = () => {
    if (posting) return;
    resetComposer();
    setShowComposer(false);
  };

  const handleShowNewPosts = async () => {
    await loadFeed({ silent: false, markAsSeen: true });
  };

  return (
    <>
      <Navbar title="My Feed" />
      <div className="page-container">
        {showNewPosts && (
          <button
            type="button"
            className="new-posts-pill"
            onClick={handleShowNewPosts}
            disabled={refreshing}
          >
            <i className="fa fa-refresh" aria-hidden="true"></i>
            {refreshing ? 'Updating...' : 'New Posts'}
          </button>
        )}

        {/* Feed */}
        {loading ? (
          <div className="loading-spinner">
            <div className="spinner"></div>
            Loading feed...
          </div>
        ) : posts.length === 0 ? (
          <div className="empty-state">
            <div className="empty-state-icon"><i className="fas fa-newspaper"></i></div>
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

        <button
          type="button"
          className="fab-btn"
          onClick={handleOpenComposer}
          aria-label="Create post"
        >
          <i className="fa fa-pencil" aria-hidden="true"></i>
        </button>

        {showComposer && (
          <div className="composer-modal-overlay" onClick={handleCloseComposer}>
            <div className="comments-modal card composer-modal" onClick={(event) => event.stopPropagation()}>
              <div className="comments-modal-header">
                <h3>Create Post</h3>
                <button className="btn btn-ghost btn-sm" type="button" onClick={handleCloseComposer} disabled={posting}>
                  Close
                </button>
              </div>

              <form className="composer-form" onSubmit={handleCreatePost}>
                <textarea
                  className="form-input composer-textarea"
                  placeholder="What's on your mind?"
                  value={content}
                  onChange={(e) => setContent(e.target.value)}
                  rows={5}
                />

                <div className="composer-actions-row">
                  <input
                    ref={fileInputRef}
                    type="file"
                    multiple
                    accept="image/*,video/*"
                    onChange={handleMediaSelect}
                    className="composer-file-input"
                  />
                  <button
                    type="button"
                    onClick={() => fileInputRef.current?.click()}
                    className="btn btn-secondary"
                    disabled={posting}
                  >
                    <i className="fa fa-image" aria-hidden="true"></i>
                    Add Media
                  </button>
                </div>

                {mediaPreviews.length > 0 && (
                  <div className="composer-preview-grid">
                    {mediaPreviews.map((preview, index) => (
                      <div key={`${preview.filename}-${index}`} className="composer-preview-item">
                        {preview.type === 'image' ? (
                          <img src={preview.url} alt="preview" />
                        ) : (
                          <video src={preview.url} />
                        )}
                        <button type="button" className="composer-preview-remove" onClick={() => removeMedia(index)}>
                          ✕
                        </button>
                      </div>
                    ))}
                  </div>
                )}

                <div className="composer-footer">
                  <span className="composer-hint">You can post with text, media, or both.</span>
                  <button
                    type="submit"
                    className="btn btn-primary"
                    disabled={posting || (!content.trim() && selectedMedia.length === 0)}
                  >
                    {posting ? 'Posting...' : 'Publish Post'}
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}
      </div>
    </>
  );
}
