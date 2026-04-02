import React from 'react';
import { calculateMatchScore, deadlineMeta, formatDeadline, safeArray } from '../utils';
import './GlassModal.css';

export default function GlassModal({ job, onClose, liked, toggleLike, profileSkills }) {
  if (!job) return null;
  const score = calculateMatchScore(job.techStack, profileSkills);
  const deadline = deadlineMeta(job.applicationDeadline);

  return (
    <div className="modal-backdrop" onClick={onClose}>
      <div className="job-modal" onClick={(e) => e.stopPropagation()}>
        <button className="job-modal__close" onClick={onClose}>✕</button>
        <div className="job-modal__header">
          <div>
            <h2>{job.positionName}</h2>
            <p>{job.companyName} · {job.location}</p>
          </div>
          <div className="job-modal__actions">
            <span className={`status-badge ${deadline.tone}`}>{deadline.label}</span>
            <button className={`job-modal__like ${liked ? 'liked' : ''}`} onClick={() => toggleLike(job.id, job.positionName)}>
              {liked ? '❤️ В избранном' : '🤍 В избранное'}
            </button>
          </div>
        </div>

        <div className="job-modal__hero">
          <div className="job-modal__hero-card">
            <span>Match score</span>
            <strong>{score}%</strong>
          </div>
          <div className="job-modal__hero-card">
            <span>Зарплата</span>
            <strong>{(job.minSalary || 0).toLocaleString('ru-RU')} ₽</strong>
          </div>
          <div className="job-modal__hero-card">
            <span>Дедлайн</span>
            <strong>{formatDeadline(job.applicationDeadline)}</strong>
          </div>
        </div>

        <div className="job-modal__grid">
          <section>
            <h3>Описание</h3>
            <p>{job.description || 'Описание пока не указано.'}</p>
          </section>
          <section>
            <h3>Стек</h3>
            <div className="job-modal__tags">
              {safeArray(job.techStack).map((tech) => <span key={tech}>{tech}</span>)}
            </div>
          </section>
          <section>
            <h3>Требования</h3>
            <p>{job.experienceRequirements || 'Без дополнительных требований.'}</p>
          </section>
          <section>
            <h3>Детали</h3>
            <ul>
              <li><strong>Период:</strong> {job.internshipDates || 'Не указан'}</li>
              <li><strong>Процесс:</strong> {job.selectionProcess || 'Не указан'}</li>
              <li><strong>Контакт:</strong> {job.contactInfo || 'Не указан'}</li>
            </ul>
          </section>
        </div>
      </div>
    </div>
  );
}
