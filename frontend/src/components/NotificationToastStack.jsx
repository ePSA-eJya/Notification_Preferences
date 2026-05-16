import { useNotifications } from "../context/NotificationContext.jsx";

export default function NotificationToastStack() {
  const { toasts, dismissToast } = useNotifications();

  if (toasts.length === 0) {
    return null;
  }

  return (
    <div
      className="notification-toast-stack"
      aria-live="polite"
      aria-atomic="true"
    >
      {toasts.map((toast) => (
        <div key={toast.id} className="notification-toast card-glass">
          <div className="notification-toast-icon">
            <i className="fas fa-bell"></i>
          </div>
          <div className="notification-toast-body">
            <div className="notification-toast-title">{toast.title}</div>
            <div className="notification-toast-message">{toast.message}</div>
          </div>
          <button
            type="button"
            className="notification-toast-close"
            onClick={() => dismissToast(toast.id)}
            aria-label="Dismiss notification"
          >
            <i className="fa fa-times" aria-hidden="true"></i>
          </button>
        </div>
      ))}
    </div>
  );
}
