import { useState } from 'react';
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
  const [commenting, setCommenting] = useState(false);
  const [commentText, setCommentText] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const initial = (post.user_id || '?')[0]?.toUpperCase();

  const handleLike = async () => {
    if (liked) return;
    try {
      await feedAPI.likePost(post.id);
      setLiked(true);
    } catch (err) {
      console.error('Like failed:', err);
    }
  };

  const handleComment = async () => {
    if (!commentText.trim() || submitting) return;
    setSubmitting(true);
    try {
      await feedAPI.commentOnPost(post.id, commentText.trim());
      setCommentText('');
      setCommenting(false);
    } catch (err) {
      console.error('Comment failed:', err);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="post-card fade-in">
      <div className="post-header">
        <div className="post-avatar">{initial}</div>
        <div>
          <div className="post-author">{post.user_id?.slice(0, 8) || 'Unknown'}</div>
          <div className="post-time">{timeAgo(post.created_at)}</div>
        </div>
      </div>

      <div className="post-content">{post.content}</div>

      <div className="post-actions">
        <button
          className={`post-action-btn ${liked ? 'active' : ''}`}
          onClick={handleLike}
        >
          {liked ? '❤️' : '🤍'} Like
        </button>
        <button
          className="post-action-btn"
          onClick={() => setCommenting(!commenting)}
        >
          💬 Comment
        </button>
      </div>

      {commenting && (
        <div className="comment-input-row fade-in">
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
      )}
    </div>
  );
}
