import React from 'react';
import './Navbar.css';

const tabs = [
  ['swipe', '🔥 Свайп'],
  ['list', '🔍 Список'],
  ['liked', '❤️ Matches'],
  ['history', '🕓 История'],
  ['dashboard', '📊 Дашборд']
];

export default function Navbar({
  view,
  setView,
  likedCount,
  skippedCount,
  theme,
  toggleTheme,
  onEnableNotifications,
  notificationsEnabled
}) {
  return (
    <nav className="navbar">
      <div className="navbar__logo">intern<span>hub</span><span className="navbar__dot" /></div>

      <div className="navbar__tabs">
        {tabs.map(([key, label]) => (
          <button key={key} className={`navbar__tab ${view === key ? 'active' : ''}`} onClick={() => setView(key)}>
            {label}
            {key === 'liked' && likedCount > 0 && <span className="navbar__badge">{likedCount}</span>}
            {key === 'history' && skippedCount > 0 && <span className="navbar__badge muted">{skippedCount}</span>}
          </button>
        ))}
      </div>

      <div className="navbar__right">
        <button className={`navbar__icon ${notificationsEnabled ? 'active' : ''}`} onClick={onEnableNotifications} title="Уведомления о дедлайнах">
          🔔
        </button>
        <button className="navbar__icon" onClick={toggleTheme} title="Переключить тему">
          {theme === 'dark' ? '🌙' : '☀️'}
        </button>
        <div className="navbar__avatar">AI</div>
      </div>
    </nav>
  );
}
