import React from 'react';
import { daysUntil, getSalaryStats, getTopTech } from '../utils';
import './DashboardView.css';

export default function DashboardView({ internships, liked, skipped, history, profileSkills }) {
  const likedItems = internships.filter((item) => liked.has(item.id));
  const skippedItems = internships.filter((item) => skipped.has(item.id));
  const stats = getSalaryStats(internships);
  const likedStats = getSalaryStats(likedItems);
  const topTech = getTopTech(internships);
  const topLikedTech = getTopTech(likedItems);
  const upcoming = internships
    .filter((item) => {
      const days = daysUntil(item.applicationDeadline);
      return days !== null && days >= 0 && days <= 14;
    })
    .sort((a, b) => (new Date(a.applicationDeadline) - new Date(b.applicationDeadline)))
    .slice(0, 5);

  const totalDecisions = liked.size + skipped.size;
  const likedPercent = totalDecisions ? Math.round((liked.size / totalDecisions) * 100) : 0;

  return (
    <div className="dashboard-view">
      <div className="dashboard-view__header">
        <h1 className="section-title">📊 Dashboard</h1>
        <p className="section-subtitle">Локальная аналитика по уже загруженным вакансиям: без дополнительных запросов и без изменений backend.</p>
      </div>

      <div className="dashboard-summary">
        <div className="summary-card glass-card">
          <span>Всего вакансий</span>
          <strong>{internships.length}</strong>
          <p>Доступно из backend</p>
        </div>
        <div className="summary-card glass-card">
          <span>В избранном</span>
          <strong>{liked.size}</strong>
          <p>{likedPercent}% от принятых решений</p>
        </div>
        <div className="summary-card glass-card">
          <span>Средняя зарплата</span>
          <strong>{stats.avg.toLocaleString('ru-RU')} ₽</strong>
          <p>Максимум: {stats.max.toLocaleString('ru-RU')} ₽</p>
        </div>
        <div className="summary-card glass-card">
          <span>Мой стек</span>
          <strong>{profileSkills.length}</strong>
          <p>{profileSkills.join(', ') || 'Не задан'}</p>
        </div>
      </div>

      <div className="dashboard-grid">
        <section className="dashboard-panel glass-card">
          <h2>Воронка решений</h2>
          <div className="funnel-ring" style={{ '--liked': `${likedPercent}%` }}>
            <div>
              <strong>{likedPercent}%</strong>
              <span>liked</span>
            </div>
          </div>
          <div className="funnel-legend">
            <span><i className="liked" /> Лайки: {liked.size}</span>
            <span><i className="skipped" /> Пропуски: {skipped.size}</span>
            <span><i className="neutral" /> Всего действий: {totalDecisions}</span>
          </div>
        </section>

        <section className="dashboard-panel glass-card">
          <h2>Топ технологий</h2>
          <div className="bars">
            {topTech.map((item) => (
              <div key={item.name} className="bar-row">
                <div className="bar-label">{item.name}</div>
                <div className="bar-track"><div style={{ width: `${(item.value / topTech[0].value) * 100}%` }} /></div>
                <strong>{item.value}</strong>
              </div>
            ))}
          </div>
        </section>

        <section className="dashboard-panel glass-card">
          <h2>Техи в избранном</h2>
          <div className="bars">
            {topLikedTech.length === 0 ? (
              <p className="dashboard-empty">Пока нет лайков — график появится автоматически.</p>
            ) : topLikedTech.map((item) => (
              <div key={item.name} className="bar-row">
                <div className="bar-label">{item.name}</div>
                <div className="bar-track secondary"><div style={{ width: `${(item.value / topLikedTech[0].value) * 100}%` }} /></div>
                <strong>{item.value}</strong>
              </div>
            ))}
          </div>
        </section>

        <section className="dashboard-panel glass-card">
          <h2>Ближайшие дедлайны</h2>
          <div className="upcoming-list">
            {upcoming.length === 0 ? (
              <p className="dashboard-empty">В ближайшие 14 дней дедлайнов не найдено.</p>
            ) : upcoming.map((item) => (
              <div key={item.id} className="upcoming-item">
                <div>
                  <strong>{item.positionName}</strong>
                  <p>{item.companyName}</p>
                </div>
                <span>{daysUntil(item.applicationDeadline)} дн.</span>
              </div>
            ))}
          </div>
        </section>
      </div>

      <div className="dashboard-note glass-card">
        <h3>Что считается на клиенте</h3>
        <p>
          Match score, фильтры, история пропусков, диаграммы, localStorage, уведомления, тема и PWA — всё это работает только во фронтенде и не требует изменений в твоём backend.
        </p>
      </div>
    </div>
  );
}
