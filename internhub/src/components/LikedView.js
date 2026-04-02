import React from 'react';
import { getColor, getEmoji } from '../data';
import { calculateMatchScore, deadlineMeta } from '../utils';
import './LikedView.css';

export default function LikedView({ internships, liked, showToast, toggleLike, onOpenDetails, profileSkills }) {
  const items = internships.filter((item) => liked.has(item.id));
  const avgSalary = items.length
    ? Math.round(items.reduce((acc, item) => acc + Number(item.minSalary || 0), 0) / items.length)
    : 0;

  return (
    <div className="liked-view">
      <div className="liked-view__header">
        <h1 className="section-title">❤️ Matches</h1>
        <p className="section-subtitle">Стажировки, которые прошли свайп вправо или были лайкнуты в списке. Система matches хранится локально и не теряется после перезапуска.</p>
      </div>

      <div className="liked-stats">
        <div className="stat-pill">Всего <span>{items.length}</span></div>
        <div className="stat-pill">Средняя зарплата <span>{avgSalary.toLocaleString('ru-RU')} ₽</span></div>
      </div>

      <div className="liked-list">
        {items.length === 0 ? (
          <div className="liked-empty glass-card">
            <div className="empty-emoji">🫙</div>
            <h2>Пока пусто</h2>
            <p>Свайпайте вправо или лайкайте вакансии в списке — они сразу появятся здесь.</p>
          </div>
        ) : items.map((intern) => {
          const color = getColor(intern.companyName);
          const emoji = getEmoji(intern.companyName);
          const deadline = deadlineMeta(intern.applicationDeadline);
          const match = calculateMatchScore(intern.techStack, profileSkills);

          return (
            <div key={intern.id} className="liked-item glass-card">
              <div className="liked-logo" style={{ background: color }}>{emoji}</div>

              <div className="liked-info">
                <div className="liked-info__top">
                  <div>
                    <h3>{intern.positionName}</h3>
                    <p>{intern.companyName} · {intern.location} · от {(intern.minSalary || 0).toLocaleString('ru-RU')} ₽/мес</p>
                  </div>
                  <span className="liked-match">{match}%</span>
                </div>

                <div className="liked-info__meta">
                  <span className={`status-badge ${deadline.tone}`}>{deadline.label}</span>
                  <span className="status-badge neutral">Технологий: {(intern.techStack || []).length}</span>
                </div>
              </div>

              <div className="liked-actions">
                <button className="btn-ghost" onClick={() => onOpenDetails(intern)}>Подробнее</button>
                <button className="btn-ghost" onClick={() => toggleLike(intern.id, intern.positionName)}>Убрать</button>
                <button className="btn-apply" onClick={() => showToast(`🚀 Отклик в ${intern.companyName} отправлен!`, 'green')}>
                  Откликнуться
                </button>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
