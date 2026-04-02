import React, { useCallback, useMemo, useState } from 'react';
import SwipeCard from './SwipeCard';
import { playFeedback } from '../utils';
import './SwipeView.css';

export default function SwipeView({
  internships,
  liked,
  skipped,
  toggleLike,
  markSkipped,
  showToast,
  onOpenDetails,
  profileSkills,
  onSuperLike
}) {
  const [flying, setFlying] = useState(null);

  const remaining = useMemo(
    () => internships.filter((intern) => !liked.has(intern.id) && !skipped.has(intern.id)),
    [internships, liked, skipped]
  );

  const progressDone = internships.length - remaining.length;

  const commitSwipe = useCallback((dir) => {
    if (!remaining.length) return;
    const intern = remaining[0];
    setFlying(dir);
    playFeedback(dir === 'super' ? 'super' : dir === 'left' ? 'skip' : 'like');

    window.setTimeout(() => {
      setFlying(null);
      if (dir === 'right') {
        toggleLike(intern.id, intern.positionName);
        showToast(`💚 Лайк: ${intern.positionName}`, 'green');
      } else if (dir === 'super') {
        onSuperLike(intern);
        showToast(`⭐ Суперлайк: ${intern.positionName}`, 'star');
      } else {
        markSkipped(intern.id, intern.positionName);
        showToast(`✕ Пропущено: ${intern.companyName}`, 'red');
      }
    }, 260);
  }, [remaining, toggleLike, showToast, onSuperLike, markSkipped]);

  if (!internships.length) {
    return (
      <div className="swipe-view">
        <div className="swipe-view__empty" style={{ position: 'static', height: 'auto', marginTop: 80 }}>
          <div className="empty-emoji">📭</div>
          <h2>Вакансий пока нет</h2>
          <p>Бекенд вернул пустой список</p>
        </div>
      </div>
    );
  }

  return (
    <div className="swipe-view">
      <div className="swipe-view__header">
        <h1>Найди стажировку мечты</h1>
        <p>Перетаскивание теперь обрабатывается плавнее: drag не триггерит постоянные React-перерисовки и ощущается заметно мягче.</p>
      </div>

      <div className="swipe-view__analytics">
        <div className="analytics-pill">Просмотрено <strong>{progressDone}</strong></div>
        <div className="analytics-pill">Осталось <strong>{remaining.length}</strong></div>
        <div className="analytics-pill">Matches <strong>{liked.size}</strong></div>
      </div>

      <div className="swipe-view__progress">
        {internships.map((intern) => {
          const isDone = liked.has(intern.id) || skipped.has(intern.id);
          const isCurrent = remaining[0]?.id === intern.id;
          return <div key={intern.id} className={`prog-dot ${isDone ? 'done' : isCurrent ? 'current' : ''}`} />;
        })}
      </div>

      <div className="swipe-view__stack">
        {remaining.length === 0 ? (
          <div className="swipe-view__empty">
            <div className="empty-emoji">🎉</div>
            <h2>Все просмотрено!</h2>
            <p>Перейдите в историю или в matches — там уже собрана ваша воронка.</p>
          </div>
        ) : (
          remaining.slice(0, 3).reverse().map((intern, index, arr) => (
            <SwipeCard
              key={intern.id}
              intern={intern}
              isTop={index === arr.length - 1}
              flying={index === arr.length - 1 ? flying : null}
              onSwipe={commitSwipe}
              stackIndex={arr.length - 1 - index}
              onOpenDetails={onOpenDetails}
              profileSkills={profileSkills}
            />
          ))
        )}
      </div>

      <div className="swipe-view__buttons">
        <button className="btn-swipe btn-nope" onClick={() => commitSwipe('left')} title="Пропустить">✕</button>
        <button className="btn-swipe btn-like" onClick={() => commitSwipe('right')} title="Добавить в matches">♥</button>
        <button className="btn-swipe btn-super" onClick={() => commitSwipe('super')} title="Суперлайк">⭐</button>
      </div>

      <div className="swipe-view__hint">
        <span>← В историю</span>
        <span>Открыть карточку — подробнее</span>
        <span>В matches →</span>
      </div>
    </div>
  );
}
