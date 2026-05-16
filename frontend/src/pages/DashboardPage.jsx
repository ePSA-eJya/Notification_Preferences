import { useState, useEffect, useRef } from 'react';
import { feedAPI } from '../services/api.js';
import { useAuth } from '../context/AuthContext.jsx';
import Navbar from '../components/Navbar.jsx';
import PostCard from '../components/PostCard.jsx';

export default function DashboardPage() {
  const { user } = useAuth();
  const [posts, setPosts] = useState([]);
  const [content, setContent] = useState('');
  const [selectedMedia, setSelectedMedia] = useState([]);
  const [mediaPreviews, setMediaPreviews] = useState([]);
  const [loading, setLoading] = useState(true);
  const [posting, setPosting] = useState(false);
  const fileInputRef = useRef(null);

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

  const handleMediaSelect = (e) => {
    const files = Array.from(e.target.files || []);
    const validFiles = [];
    const newPreviews = [];

    files.forEach(file => {
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
        newPreviews.push({
          type: file.type.startsWith('image/') ? 'image' : 'video',
          url: event.target.result,
          filename: file.name,
        });
        setMediaPreviews(prev => [...prev, {
          type: file.type.startsWith('image/') ? 'image' : 'video',
          url: event.target.result,
          filename: file.name,
        }]);
      };
      reader.readAsDataURL(file);
    });

    setSelectedMedia(validFiles);
  };

  const removeMedia = (index) => {
    setSelectedMedia(prev => prev.filter((_, i) => i !== index));
    setMediaPreviews(prev => prev.filter((_, i) => i !== index));
  };

  const handleCreatePost = async (e) => {
    e.preventDefault();
    if (!content.trim() && selectedMedia.length === 0 || posting) return;
    setPosting(true);
    try {
      const newPost = await feedAPI.createPost(content.trim(), selectedMedia);
      console.log('Post created:', newPost);
      console.log('Media URLs:', newPost.media_urls);
      setPosts((prev) => [newPost, ...prev]);
      setContent('');
      setSelectedMedia([]);
      setMediaPreviews([]);
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
    } catch (err) {
      console.error('Post creation failed:', err);
      alert('Failed to create post: ' + err.message);
    } finally {
      setPosting(false);
    }
  };

  return (
    <>
      <Navbar title="My Feed" />
      <div className="page-container">
        {/* Create Post */}
        <form className="create-post card" onSubmit={handleCreatePost}>
          <textarea
            className="form-input"
            placeholder="What's on your mind?"
            value={content}
            onChange={(e) => setContent(e.target.value)}
            rows={3}
          />
          
          {/* Media Upload Section */}
          <div style={{ marginTop: '12px' }}>
            <input
              ref={fileInputRef}
              type="file"
              multiple
              accept="image/*,video/*"
              onChange={handleMediaSelect}
              style={{ display: 'none' }}
              id="media-input"
            />
            <label htmlFor="media-input" style={{ cursor: 'pointer', marginRight: '8px' }}>
              <button
                type="button"
                onClick={() => fileInputRef.current?.click()}
                className="btn btn-secondary"
                style={{ marginRight: '8px' }}
              >
                <i className="fas fa-image"></i> Add Media
              </button>
            </label>
          </div>

          {/* Media Previews */}
          {mediaPreviews.length > 0 && (
            <div style={{ marginTop: '12px', display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(80px, 1fr))', gap: '8px' }}>
              {mediaPreviews.map((preview, index) => (
                <div key={index} style={{ position: 'relative', borderRadius: '4px', overflow: 'hidden', backgroundColor: '#f0f0f0' }}>
                  {preview.type === 'image' ? (
                    <img src={preview.url} alt="preview" style={{ width: '100%', height: '80px', objectFit: 'cover' }} />
                  ) : (
                    <video src={preview.url} style={{ width: '100%', height: '80px', objectFit: 'cover' }} />
                  )}
                  <button
                    type="button"
                    onClick={() => removeMedia(index)}
                    style={{
                      position: 'absolute',
                      top: '0',
                      right: '0',
                      background: 'rgba(0,0,0,0.7)',
                      color: 'white',
                      border: 'none',
                      borderRadius: '0',
                      cursor: 'pointer',
                      width: '24px',
                      height: '24px',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      fontSize: '14px',
                    }}
                  >
                    ✕
                  </button>
                </div>
              ))}
            </div>
          )}

          <button
            type="submit"
            className="btn btn-primary"
            disabled={posting || (!content.trim() && selectedMedia.length === 0)}
            style={{ marginTop: '12px' }}
          >
            {posting ? 'Posting...' : 'Publish Post'}
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
      </div>
    </>
  );
}
