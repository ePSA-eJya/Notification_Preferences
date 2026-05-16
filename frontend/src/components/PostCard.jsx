import { useState, useEffect } from 'react';
import { feedAPI } from '../services/api.js';

function timeAgo(dateStr) {
  const diff = Date.now() - new Date(dateStr).getTime();
  const mins = Math.floor(diff / 60000);
  if (mins < 1) return 'just now';
  if (mins < 60) return `${mins}m ago`;
  const hrs = Math.floor(mins / 60);
  if (hrs < 24) return `${hrs}h ago`;
  const days = Math.floor(hrs / 24);
  return `${days}d ago`;
}

export default function PostCard({ post, currentUserId }) {
  const [liked, setLiked] = useState(false);
  const [likeLoading, setLikeLoading] = useState(true);
  const [showCommentsModal, setShowCommentsModal] = useState(false);
  const [comments, setComments] = useState(Array.isArray(post.comments) ? post.comments : []);
  const [commentsLoading, setCommentsLoading] = useState(false);
  const [commentText, setCommentText] = useState('');
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    setComments(Array.isArray(post.comments) ? post.comments : []);
  }, [post.id, post.comments]);

  useEffect(() => {
    if (!showCommentsModal) return;

    let active = true;
    setCommentsLoading(true);

    feedAPI.getPostComments(post.id)
      .then((data) => {
        if (!active) return;
        setComments(Array.isArray(data?.comments) ? data.comments : []);
      })
      .catch((err) => {
        console.error('Failed to load comments:', err);
      })
      .finally(() => {
        if (active) setCommentsLoading(false);
      });

    return () => {
      active = false;
    };
  }, [showCommentsModal, post.id]);

  // Check if post is already liked by current user
  useEffect(() => {
    feedAPI.isPostLiked(post.id)
      .then((data) => {
        setLiked(data?.liked || false);
      })
      .catch((err) => console.error('Failed to check like status:', err))
      .finally(() => setLikeLoading(false));
  }, [post.id]);

  const handleLike = async () => {
    if (liked) {
      // Unlike
      setLiked(false);
      try {
        await feedAPI.unlikePost(post.id);
      } catch (err) {
        setLiked(true);
        console.error('Unlike failed:', err);
      }
    } else {
      // Like
      setLiked(true);
      try {
        await feedAPI.likePost(post.id);
      } catch (err) {
        setLiked(false);
        console.error('Like failed:', err);
      }
    }
  };

  const handleComment = async () => {
    if (!commentText.trim() || submitting) return;
    setSubmitting(true);
    const text = commentText.trim();
    try {
      await feedAPI.commentOnPost(post.id, text);
      const data = await feedAPI.getPostComments(post.id);
      setComments(Array.isArray(data?.comments) ? data.comments : []);
      setCommentText('');
    } catch (err) {
      console.error('Comment failed:', err);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="post-card fade-in">
      <div className="post-header">
        <div className="post-avatar">{post.user_handle?.[0]?.toUpperCase() || '?'}</div>
        <div>
          <div className="post-author">@{post.user_handle || 'unknown'}</div>
          <div className="post-time">{timeAgo(post.created_at)}</div>
        </div>
      </div>

      <div className="post-content">{post.content}</div>

      {/* Media Display */}
      {post.media_urls && post.media_urls.length > 0 && (
        <div style={{ marginTop: '12px', display: 'grid', gridTemplateColumns: post.media_urls.length === 1 ? '1fr' : 'repeat(auto-fit, minmax(200px, 1fr))', gap: '8px' }}>
          {post.media_urls.map((mediaUrl, index) => {
            const isVideo = mediaUrl.endsWith('.mp4') || mediaUrl.endsWith('.webm') || mediaUrl.endsWith('.mov') || mediaUrl.endsWith('.avi');
            console.log(`Rendering media ${index}:`, mediaUrl, 'isVideo:', isVideo);
            return (
              <div key={index} style={{ borderRadius: '8px', overflow: 'hidden', backgroundColor: '#f0f0f0' }}>
                {isVideo ? (
                  <video 
                    src={mediaUrl} 
                    controls 
                    onError={(e) => console.error('Video failed to load:', mediaUrl, e)}
                    style={{ width: '100%', height: 'auto', maxHeight: '400px', objectFit: 'cover' }} 
                  />
                ) : (
                  <img 
                    src={mediaUrl} 
                    alt="post media" 
                    onError={(e) => console.error('Image failed to load:', mediaUrl, e)}
                    style={{ width: '100%', height: 'auto', maxHeight: '400px', objectFit: 'cover' }} 
                  />
                )}
              </div>
            );
          })}
        </div>
      )}

      <div className="post-actions">
        <button
          className={`post-action-btn ${liked ? 'active' : ''}`}
          onClick={handleLike}
          disabled={likeLoading}
        >
          {liked ? <i class="fa fa-heart" style={{ color: 'red', fontsize: '18px' }} aria-hidden="true"></i> : <i class="fa fa-heart" style={{ color: '#6f6c6cff' }} aria-hidden="true"></i>} {liked ? 'Unlike' : 'Like'}
        </button>
        <button
          className="post-action-btn"
          onClick={() => setShowCommentsModal(true)}
        >
          <i class="fa fa-comment" style={{ color: '#6f6c6cff' }} aria-hidden="true"></i>
          Comments {comments.length > 0 ? `(${comments.length})` : ''}
        </button>

      </div>

      {showCommentsModal && (
        <div className="comments-modal-overlay" onClick={() => setShowCommentsModal(false)}>
          <div className="comments-modal card" onClick={(e) => e.stopPropagation()}>
            <div className="comments-modal-header">
              <h3>Comments</h3>
              <button
                className="btn btn-ghost btn-sm"
                type="button"
                onClick={() => setShowCommentsModal(false)}
              >
                Close
              </button>
            </div>

            <div className="comments-modal-list">
              {commentsLoading ? (
                <div className="comments-empty">Loading comments...</div>
              ) : comments.length === 0 ? (
                <div className="comments-empty">No comments yet. Be the first one to comment.</div>
              ) : (
                comments.map((comment) => (
                  <div className="comment-item" key={comment.id || `${comment.user_id}-${comment.created_at}`}>
                    <div className="comment-item-header">
                      <span className="comment-author">@{comment.user_handle || 'unknown'}</span>
                      <span className="comment-time">{timeAgo(comment.created_at)}</span>
                    </div>
                    <div className="comment-text">{comment.text}</div>
                  </div>
                ))
              )}
            </div>

            <div className="comment-input-row">
              <input
                className="form-input"
                placeholder="Write a comment..."
                value={commentText}
                onChange={(e) => setCommentText(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && handleComment()}
              />
              <button
                className="btn btn-primary btn-sm"
                onClick={handleComment}
                disabled={submitting || !commentText.trim()}
              >
                {submitting ? '...' : 'Send'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
