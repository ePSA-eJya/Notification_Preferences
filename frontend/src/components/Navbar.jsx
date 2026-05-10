export default function Navbar({ title }) {
  return (
    <header className="navbar">
      <h2 className="navbar-title">{title}</h2>
      <div className="navbar-actions">
        {/* Future: search bar, notification bell count */}
      </div>
    </header>
  );
}
