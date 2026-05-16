function timeAgo(dateStr) {
  const diff = Date.now() - new Date(dateStr).getTime();
  const mins = Math.floor(diff / 60000);
  if (mins < 1) return "just now";
  if (mins < 60) return `${mins}m ago`;
  const hrs = Math.floor(mins / 60);
  if (hrs < 24) return `${hrs}h ago`;
  const days = Math.floor(hrs / 24);
  return `${days}d ago`;
}

const statusConfig = {
  DELIVERED: { label: "Delivered", cls: "badge-success" },
  SENT: { label: "Sent", cls: "badge-info" },
  PENDING: { label: "Pending", cls: "badge-warning" },
  FAILED: { label: "Failed", cls: "badge-error" },
  SKIPPED: { label: "Skipped", cls: "badge-muted" },
};

function ChannelBadge({ channel, status }) {
  const cfg = statusConfig[status] || statusConfig.PENDING;
  const isEnabled = status && status !== "SKIPPED";
  const colorClass = isEnabled ? "badge-dark" : "badge-light";
  return (
    <span className={`badge ${colorClass}`} title={`${channel}: ${status}`}>
      <i className={cfg.icon}></i> {channel}
    </span>
  );
}

export default function NotificationItem({
  notification,
  unread = false,
  onClick,
}) {
  const { channels } = notification;
  const clickable = typeof onClick === "function";

  return (
    <div
      className={`notif-item fade-in ${unread ? "unread" : ""} ${clickable ? "clickable" : ""}`}
      role={clickable ? "button" : undefined}
      tabIndex={clickable ? 0 : undefined}
      onClick={onClick}
      onKeyDown={(event) => {
        if (!clickable) return;
        if (event.key === "Enter" || event.key === " ") {
          event.preventDefault();
          onClick();
        }
      }}
    >
      <div
        className="notif-icon"
        style={{ background: "var(--accent-soft)", color: "var(--accent)" }}
      >
        <i className="fas fa-bell"></i>
      </div>
      <div className="notif-body">
        <div className="notif-message">{notification.message}</div>
        <div className="notif-channels">
          <ChannelBadge channel="InApp" status={channels?.in_app?.status} />
          <ChannelBadge channel="Push" status={channels?.push?.status} />
          <ChannelBadge channel="Email" status={channels?.email?.status} />
        </div>
      </div>
      <div
        className="notif-time"
        style={{ marginTop: 0, whiteSpace: "nowrap" }}
      >
        {timeAgo(notification.created_at)}
      </div>
    </div>
  );
}
