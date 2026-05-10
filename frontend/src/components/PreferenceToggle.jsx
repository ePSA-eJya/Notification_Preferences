const options = ['ALL', 'FOLLOWERS', 'NONE'];

export default function PreferenceToggle({ value, onChange }) {
  return (
    <select
      className="form-select"
      value={value || 'ALL'}
      onChange={(e) => onChange(e.target.value)}
    >
      {options.map((opt) => (
        <option key={opt} value={opt}>{opt}</option>
      ))}
    </select>
  );
}
