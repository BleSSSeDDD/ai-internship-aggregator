import React, { useEffect, useMemo, useRef, useState } from 'react';
import './App.css';
import Navbar from './components/Navbar';
import SwipeView from './components/SwipeView';
import ListView from './components/ListView';
import LikedView from './components/LikedView';
import HistoryView from './components/HistoryView';
import DashboardView from './components/DashboardView';
import Toast from './components/Toast';
import Loader from './components/Loader';
import GlassModal from './components/GlassModal';
import ParticlesBackground from './components/ParticlesBackground';
import ConfettiLayer from './components/ConfettiLayer';
import { BASE_URL, fetchInternships } from './api';
import {
  notifyUpcomingDeadlines,
  readJson,
  readSet,
  requestDeadlinePermission,
  setToArray,
  STORAGE_KEYS,
  writeJson
} from './utils';

export default function App() {
  const [view, setView] = useState(() => readJson(STORAGE_KEYS.view, 'swipe'));
  const [liked, setLiked] = useState(() => readSet(STORAGE_KEYS.liked));
  const [skipped, setSkipped] = useState(() => readSet(STORAGE_KEYS.skipped));
  const [history, setHistory] = useState(() => readJson(STORAGE_KEYS.history, []));
  const [profileSkills, setProfileSkills] = useState(() => readJson(STORAGE_KEYS.profileSkills, ['React', 'Java', 'Spring']));
  const [theme, setTheme] = useState(() => readJson(STORAGE_KEYS.theme, 'dark'));
  const [toast, setToast] = useState(null);
  const [internships, setInternships] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [selectedJob, setSelectedJob] = useState(null);
  const [confettiBursts, setConfettiBursts] = useState([]);
  const [notificationsEnabled, setNotificationsEnabled] = useState(
    typeof Notification !== 'undefined' && Notification.permission === 'granted'
  );

  const toastTimerRef = useRef(null);

  useEffect(() => {
    document.body.setAttribute('data-theme', theme);
    writeJson(STORAGE_KEYS.theme, theme);
  }, [theme]);

  useEffect(() => {
    writeJson(STORAGE_KEYS.view, view);
  }, [view]);

  useEffect(() => {
    writeJson(STORAGE_KEYS.liked, setToArray(liked));
  }, [liked]);

  useEffect(() => {
    writeJson(STORAGE_KEYS.skipped, setToArray(skipped));
  }, [skipped]);

  useEffect(() => {
    writeJson(STORAGE_KEYS.history, history);
  }, [history]);

  useEffect(() => {
    writeJson(STORAGE_KEYS.profileSkills, profileSkills);
  }, [profileSkills]);

  useEffect(() => {
    const controller = new AbortController();
    let mounted = true;
    setLoading(true);
    setError(null);

    fetchInternships(controller.signal)
      .then((data) => {
        if (!mounted) return;
        setInternships(data);
        setLoading(false);
      })
      .catch((err) => {
        if (!mounted || err?.name === 'AbortError') return;
        setError(err.message);
        setLoading(false);
      });

    return () => {
      mounted = false;
      controller.abort();
    };
  }, []);

  useEffect(() => {
    if (!notificationsEnabled || !internships.length) return;
    const notified = readSet(STORAGE_KEYS.notified);

    const run = () => {
      notifyUpcomingDeadlines(internships, notified, (key) => {
        notified.add(key);
        writeJson(STORAGE_KEYS.notified, setToArray(notified));
      });
    };

    run();
    const interval = window.setInterval(run, 60000);
    return () => window.clearInterval(interval);
  }, [notificationsEnabled, internships]);

  useEffect(() => () => {
    if (toastTimerRef.current) clearTimeout(toastTimerRef.current);
  }, []);

  const showToast = (msg, type = '') => {
    if (toastTimerRef.current) clearTimeout(toastTimerRef.current);
    setToast({ msg, type });
    toastTimerRef.current = setTimeout(() => setToast(null), 2800);
  };

  const pushHistory = (entry) => {
    setHistory((prev) => [...prev.slice(-79), { ...entry, id: crypto.randomUUID?.() || `${Date.now()}-${Math.random()}`, at: new Date().toISOString() }]);
  };

  const toggleLike = (id, name) => {
    setSkipped((prev) => {
      if (!prev.has(id)) return prev;
      const next = new Set(prev);
      next.delete(id);
      return next;
    });

    setLiked((prev) => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
        pushHistory({ type: 'unlike', label: `Убрано из избранного: ${name}` });
        showToast(`Убрано: ${name}`, 'red');
      } else {
        next.add(id);
        pushHistory({ type: 'like', label: `Добавлено в избранное: ${name}` });
      }
      return next;
    });
  };

  const markSkipped = (id, name) => {
    setLiked((prev) => {
      if (!prev.has(id)) return prev;
      const next = new Set(prev);
      next.delete(id);
      return next;
    });

    setSkipped((prev) => {
      const next = new Set(prev);
      next.add(id);
      return next;
    });

    pushHistory({ type: 'skip', label: `Пропущено: ${name}` });
  };

  const restoreSkipped = (id) => {
    const intern = internships.find((item) => item.id === id);
    setSkipped((prev) => {
      const next = new Set(prev);
      next.delete(id);
      return next;
    });
    pushHistory({ type: 'restore', label: `Возвращено в стек: ${intern?.positionName || 'Вакансия'}` });
    showToast(`↩️ Возвращено: ${intern?.positionName || 'Вакансия'}`, 'green');
  };

  const onSuperLike = (intern) => {
    setLiked((prev) => {
      const next = new Set(prev);
      next.add(intern.id);
      return next;
    });
    setSkipped((prev) => {
      const next = new Set(prev);
      next.delete(intern.id);
      return next;
    });
    setConfettiBursts((prev) => [...prev.slice(-4), { id: Date.now() }]);
    pushHistory({ type: 'super', label: `Суперлайк: ${intern.positionName}` });
  };

  const handleEnableNotifications = async () => {
    const permission = await requestDeadlinePermission();
    if (permission === 'granted') {
      setNotificationsEnabled(true);
      showToast('🔔 Уведомления о дедлайнах включены', 'green');
    } else if (permission === 'unsupported') {
      showToast('Этот браузер не поддерживает Notification API', 'red');
    } else {
      showToast('Разрешение на уведомления не выдано', 'red');
    }
  };

  const quickStats = useMemo(() => ([
    ['Реальных вакансий', internships.length],
    ['Лайков', liked.size],
    ['Пропусков', skipped.size],
    ['Профильных скиллов', profileSkills.length]
  ]), [internships.length, liked.size, skipped.size, profileSkills.length]);

  if (loading) return <Loader />;

  if (error) {
    return (
      <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', minHeight: '100vh', gap: 16, padding: 24 }}>
        <div style={{ fontSize: 64 }}>🔌</div>
        <h2 style={{ fontSize: 22, fontWeight: 800 }}>Не удалось подключиться к серверу</h2>
        <p style={{ color: 'var(--muted)', textAlign: 'center' }}>
          Убедитесь, что backend запущен на <code style={{ color: 'var(--cyan)' }}>{BASE_URL}</code>
        </p>
        <p style={{ color: '#FF6B6B', fontSize: 13 }}>{error}</p>
        <button
          onClick={() => window.location.reload()}
          style={{ padding: '12px 28px', borderRadius: 12, background: 'var(--grad1)', color: '#fff', border: 'none', fontWeight: 700, cursor: 'pointer', fontSize: 15 }}
        >
          Повторить
        </button>
      </div>
    );
  }

  return (
    <div className="app-shell">
      <ParticlesBackground />
      <ConfettiLayer bursts={confettiBursts} />
      <Navbar
        view={view}
        setView={setView}
        likedCount={liked.size}
        skippedCount={skipped.size}
        theme={theme}
        toggleTheme={() => setTheme((prev) => prev === 'dark' ? 'light' : 'dark')}
        onEnableNotifications={handleEnableNotifications}
        notificationsEnabled={notificationsEnabled}
      />

      <div className="quick-actions">
        {quickStats.map(([label, value]) => (
          <div key={label} className="quick-pill">{label}<strong>{value}</strong></div>
        ))}
      </div>

      <main className="app-main">
        {view === 'swipe' && (
          <SwipeView
            internships={internships}
            liked={liked}
            skipped={skipped}
            toggleLike={toggleLike}
            markSkipped={markSkipped}
            showToast={showToast}
            onOpenDetails={setSelectedJob}
            profileSkills={profileSkills}
            onSuperLike={onSuperLike}
          />
        )}

        {view === 'list' && (
          <ListView
            internships={internships}
            liked={liked}
            skipped={skipped}
            toggleLike={toggleLike}
            markSkipped={markSkipped}
            onOpenDetails={setSelectedJob}
            profileSkills={profileSkills}
            setProfileSkills={setProfileSkills}
          />
        )}

        {view === 'liked' && (
          <LikedView
            internships={internships}
            liked={liked}
            showToast={showToast}
            toggleLike={toggleLike}
            onOpenDetails={setSelectedJob}
            profileSkills={profileSkills}
          />
        )}

        {view === 'history' && (
          <HistoryView
            internships={internships}
            skipped={skipped}
            history={history}
            restoreSkipped={restoreSkipped}
            onOpenDetails={setSelectedJob}
            toggleLike={toggleLike}
          />
        )}

        {view === 'dashboard' && (
          <DashboardView
            internships={internships}
            liked={liked}
            skipped={skipped}
            history={history}
            profileSkills={profileSkills}
          />
        )}
      </main>

      <div className="app-footer-note">
        Backend берёт на себя компанию, локацию, технологию и зарплату. Локально остаются matches, история, модалка, swipe, быстрый поиск и вся визуальная аналитика.
      </div>

      <GlassModal
        job={selectedJob}
        onClose={() => setSelectedJob(null)}
        liked={selectedJob ? liked.has(selectedJob.id) : false}
        toggleLike={toggleLike}
        profileSkills={profileSkills}
      />

      {toast && <Toast msg={toast.msg} type={toast.type} />}
    </div>
  );
}
