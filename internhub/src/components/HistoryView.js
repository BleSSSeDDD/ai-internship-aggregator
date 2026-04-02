import React from 'react';
import { getColor, getEmoji } from '../data';
import { deadlineMeta, formatDeadline } from '../utils';
import './HistoryView.css';

export default function HistoryView({ internships, skipped, history, restoreSkipped, onOpenDetails, toggleLike }) {
  const skippedItems = internships.filter((item) => skipped.has(item.id));
  const recent = history.slice().reverse().slice(0, 12);

  return (
    <div className="history-view">
      <div className="history-view__header">
        <h1 className="section-title">🕓 История и пропуски</h1>
        <p className="section-subtitle">Пропущенные вакансии можно вернуть в свайпы. История действий тоже хранится локально.</p>
      </div>

      <div className="history-grid">
        <section className="history-card glass-card">
          <h2>Пропущенные вакансии</h2>
          {skippedItems.length === 0 ? (
            <p className="history-empty">Пока ничего не пропускали.</p>
          ) : skippedItems.map((intern) => {
            const color = getColor(intern.companyName);
            const emoji = getEmoji(intern.companyName);
            const deadline = deadlineMeta(intern.applicationDeadline);

            return (
              <div key={intern.id} className="history-item">
                <div className="history-logo" style={{ background: color }}>{emoji}</div>
                <div className="history-info">
                  <strong>{intern.positionName}</strong>
                  <span>{intern.companyName} · {formatDeadline(intern.applicationDeadline)}</span>
                  <span className={`status-badge ${deadline.tone}`}>{deadline.label}</span>
                </div>
                <div className="history-actions">
                  <button onClick={() => restoreSkipped(intern.id)}>Вернуть</button>
                  <button onClick={() => onOpenDetails(intern)}>Подробнее</button>
                  <button onClick={() => toggleLike(intern.id, intern.positionName)}>Лайк</button>
                </div>
              </div>
            );
          })}
        </section>

        <section className="history-card glass-card">
          <h2>Лента действий</h2>
          {recent.length === 0 ? (
            <p className="history-empty">История пока пустая.</p>
          ) : (
            <div className="history-timeline">
              {recent.map((entry) => (
                <div key={entry.id} className="timeline-item">
                  <div className={`timeline-dot ${entry.type}`} />
                  <div>
                    <strong>{entry.label}</strong>
                    <p>{new Date(entry.at).toLocaleString('ru-RU')}</p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </section>
      </div>
    </div>
  );
}
