import React, { memo, useEffect, useMemo, useRef } from 'react';
import { getColor, getEmoji } from '../data';
import { calculateMatchScore, deadlineMeta, safeArray } from '../utils';
import './SwipeCard.css';

const SWIPE_THRESHOLD = 110;

function SwipeCard({ intern, isTop, flying, onSwipe, stackIndex, onOpenDetails, profileSkills }) {
  const cardRef = useRef(null);
  const dragRef = useRef({
    active: false,
    pointerId: null,
    startX: 0,
    startY: 0,
    x: 0,
    y: 0,
    moved: false,
    frame: 0
  });

  const color = useMemo(() => getColor(intern.companyName), [intern.companyName]);
  const emoji = useMemo(() => getEmoji(intern.companyName), [intern.companyName]);
  const deadline = useMemo(() => deadlineMeta(intern.applicationDeadline), [intern.applicationDeadline]);
  const matchScore = useMemo(() => calculateMatchScore(intern.techStack, profileSkills), [intern.techStack, profileSkills]);
  const baseScale = 1 - stackIndex * 0.04;
  const baseTranslate = stackIndex * 12;
  const brightness = 1 - stackIndex * 0.08;
  const flyClass = flying ? `fly-${flying}` : '';

  useEffect(() => {
    const node = cardRef.current;
    if (!node) return;
    node.style.setProperty('--stack-y', `${baseTranslate}px`);
    node.style.setProperty('--stack-scale', `${baseScale}`);
    node.style.setProperty('--drag-x', '0px');
    node.style.setProperty('--drag-y', '0px');
    node.style.setProperty('--drag-rotate', '0deg');
    node.style.setProperty('--like-opacity', '0');
    node.style.setProperty('--nope-opacity', '0');
  }, [baseScale, baseTranslate, intern.id]);

  useEffect(() => () => {
    if (dragRef.current.frame) cancelAnimationFrame(dragRef.current.frame);
  }, []);

  const applyFrame = () => {
    dragRef.current.frame = 0;
    const node = cardRef.current;
    if (!node) return;
    const { x, y } = dragRef.current;
    const clampedX = Math.max(-220, Math.min(220, x));
    node.style.setProperty('--drag-x', `${clampedX}px`);
    node.style.setProperty('--drag-y', `${y}px`);
    node.style.setProperty('--drag-rotate', `${clampedX * 0.05}deg`);
    node.style.setProperty('--like-opacity', String(clampedX > 20 ? Math.min(clampedX / 80, 1) : 0));
    node.style.setProperty('--nope-opacity', String(clampedX < -20 ? Math.min(Math.abs(clampedX) / 80, 1) : 0));
  };

  const scheduleFrame = () => {
    if (dragRef.current.frame) return;
    dragRef.current.frame = requestAnimationFrame(applyFrame);
  };

  const resetDragVisual = () => {
    dragRef.current.x = 0;
    dragRef.current.y = 0;
    dragRef.current.moved = false;
    const node = cardRef.current;
    if (!node) return;
    node.classList.remove('dragging');
    node.style.setProperty('--drag-x', '0px');
    node.style.setProperty('--drag-y', '0px');
    node.style.setProperty('--drag-rotate', '0deg');
    node.style.setProperty('--like-opacity', '0');
    node.style.setProperty('--nope-opacity', '0');
  };

  const finishGesture = () => {
    const dx = dragRef.current.x;
    dragRef.current.active = false;
    dragRef.current.pointerId = null;

    if (Math.abs(dx) > SWIPE_THRESHOLD) {
      onSwipe(dx > 0 ? 'right' : 'left');
      return;
    }

    resetDragVisual();
  };

  const handlePointerDown = (event) => {
    if (!isTop) return;
    if (event.target.closest('button, a, input, textarea, select')) return;

    const node = cardRef.current;
    if (!node) return;

    dragRef.current.active = true;
    dragRef.current.pointerId = event.pointerId;
    dragRef.current.startX = event.clientX;
    dragRef.current.startY = event.clientY;
    dragRef.current.x = 0;
    dragRef.current.y = 0;
    dragRef.current.moved = false;

    node.setPointerCapture?.(event.pointerId);
    node.classList.add('dragging');
  };

  const handlePointerMove = (event) => {
    if (!isTop || !dragRef.current.active || dragRef.current.pointerId !== event.pointerId) return;

    dragRef.current.x = event.clientX - dragRef.current.startX;
    dragRef.current.y = (event.clientY - dragRef.current.startY) * 0.55;
    dragRef.current.moved = dragRef.current.moved || Math.abs(dragRef.current.x) > 6 || Math.abs(dragRef.current.y) > 6;
    scheduleFrame();
  };

  const handlePointerUp = (event) => {
    if (!dragRef.current.active || dragRef.current.pointerId !== event.pointerId) return;
    cardRef.current?.releasePointerCapture?.(event.pointerId);
    finishGesture();
  };

  const handlePointerCancel = (event) => {
    if (!dragRef.current.active || dragRef.current.pointerId !== event.pointerId) return;
    cardRef.current?.releasePointerCapture?.(event.pointerId);
    dragRef.current.active = false;
    dragRef.current.pointerId = null;
    resetDragVisual();
  };

  return (
    <div
      ref={cardRef}
      className={`swipe-card ${flyClass}`}
      style={{
        filter: `brightness(${brightness})`,
        zIndex: 10 - stackIndex,
        pointerEvents: isTop ? 'all' : 'none',
        cursor: isTop ? 'grab' : 'default'
      }}
      onPointerDown={handlePointerDown}
      onPointerMove={handlePointerMove}
      onPointerUp={handlePointerUp}
      onPointerCancel={handlePointerCancel}
    >
      <div className="card-band" style={{ background: color }} />
      <span className="badge like-badge">ЛАЙК 💚</span>
      <span className="badge nope-badge">ПРОПУСК ✕</span>

      <div className="card-content">
        <div className="company-row">
          <div className="company-logo" style={{ background: color }}>{emoji}</div>
          <div>
            <div className="company-name">{intern.companyName}</div>
            <div className="position-name">{intern.positionName}</div>
          </div>
          <button className="card-details-btn" onClick={(e) => { e.stopPropagation(); onOpenDetails(intern); }}>
            Подробнее
          </button>
        </div>

        <div className="card-divider" />

        <div className="card-meta">
          <div className="meta-item">
            <div className="meta-label">💰 Зарплата</div>
            <div className="meta-value meta-salary">{(intern.minSalary || 0).toLocaleString('ru-RU')} ₽</div>
          </div>
          <div className="meta-item">
            <div className="meta-label">📍 Локация</div>
            <div className="meta-value">{intern.location || 'Не указана'}</div>
          </div>
          <div className="meta-item">
            <div className="meta-label">📅 Период</div>
            <div className="meta-value">{intern.internshipDates || 'Не указан'}</div>
          </div>
          <div className="meta-item">
            <div className="meta-label">⏰ Дедлайн</div>
            <div className={`meta-value meta-deadline meta-${deadline.tone}`}>{deadline.label}</div>
          </div>
        </div>

        <div className="match-row">
          <div className="match-pill">
            <span>Match score</span>
            <strong>{matchScore}%</strong>
          </div>
          <div className="match-bar">
            <div style={{ width: `${matchScore}%` }} />
          </div>
        </div>

        <p className="card-desc">{intern.description || 'Описание пока не добавлено.'}</p>

        <div className="tech-stack">
          {safeArray(intern.techStack).slice(0, 6).map((tech) => (
            <span key={tech} className="tech-tag">{tech}</span>
          ))}
        </div>
      </div>
    </div>
  );
}

export default memo(SwipeCard);
