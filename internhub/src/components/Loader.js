import React from 'react';
import './Loader.css';

export default function Loader() {
  return (
    <div className="loader-shell">
      <div className="loader-top">
        <div className="loader-line loader-line--lg" />
        <div className="loader-line loader-line--sm" />
      </div>
      <div className="loader-stack">
        {[0, 1, 2].map((item) => (
          <div key={item} className="loader-card">
            <div className="loader-avatar" />
            <div className="loader-text">
              <div className="loader-line loader-line--md" />
              <div className="loader-line loader-line--sm" />
              <div className="loader-line loader-line--xs" />
            </div>
          </div>
        ))}
      </div>
      <p>Загружаем стажировки и собираем дашборд…</p>
    </div>
  );
}
