// const options = ['ALL', 'FOLLOWERS', 'NONE'];
const preferenceOptions = {
  'follows': ['ALL', 'FOLLOWERS', 'NONE'],
  'posts': ['FOLLOWING', 'NONE'],
  'likes': ['ALL', 'FOLLOWERS', 'NONE'],
  'comments': ['ALL', 'FOLLOWERS', 'NONE'],
};

export default function PreferenceToggle({
  value,
  actionKey,
  onChange,
}) {
  const options = preferenceOptions[actionKey] || ['ALL'];

  return (
    <select
      className="form-select"
      value={value || 'ALL'}
      onChange={(e) => onChange(e.target.value)}
    >
      {options.map((opt) => (
        <option key={opt} value={opt}>
          {opt}
        </option>
      ))}
    </select>
  );
}
// export default function PreferenceToggle({ value, onChange }) {
//   return (
//     <select
//       className="form-select"
//       value={value || 'ALL'}
//       onChange={(e) => onChange(e.target.value)}
//     >
//       {options.map((opt) => (
//         <option key={opt} value={opt}>{opt}</option>
//       ))}
//     </select>
//   );
// }
