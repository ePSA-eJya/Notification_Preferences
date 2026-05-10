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

const statusConfig = {
  DELIVERED: { label: 'Delivered', cls: 'badge-success', icon: '✅' },
  SENT: { label: 'Sent', cls: 'badge-info', icon: '📤' },
  PENDING: { label: 'Pending', cls: 'badge-warning', icon: '⏳' },
  FAILED: { label: 'Failed', cls: 'badge-error', icon: '❌' },
  SKIPPED: { label: 'Skipped', cls: 'badge-muted', icon: '⏭' },
};

function ChannelBadge({ channel, status }) {
  const cfg = statusConfig[status] || statusConfig.PENDING;
  return (
    <span className={`badge ${cfg.cls}`} title={`${channel}: ${status}`}>
      {cfg.icon} {channel}
    </span>
  );
}

export default function NotificationItem({ notification }) {
  const { channels } = notification;

  return (
    <div className="notif-item fade-in">
      <div className="notif-icon" style={{ background: 'var(--accent-soft)', color: 'var(--accent)' }}>
        🔔
      </div>
      <div className="notif-body">
        <div className="notif-message">{notification.message}</div>
        <div className="notif-time">{timeAgo(notification.created_at)}</div>
        <div className="notif-channels">
          <ChannelBadge channel="InApp" status={channels?.in_app?.status} />
          <ChannelBadge channel="Push" status={channels?.push?.status} />
          <ChannelBadge channel="Email" status={channels?.email?.status} />
        </div>
      </div>
    </div>
  );
}
