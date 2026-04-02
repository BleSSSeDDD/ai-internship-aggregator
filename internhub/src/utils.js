export const STORAGE_KEYS = {
  liked: 'internhub:liked',
  skipped: 'internhub:skipped',
  history: 'internhub:history',
  theme: 'internhub:theme',
  listFilters: 'internhub:listFilters',
  profileSkills: 'internhub:profileSkills',
  notified: 'internhub:notifiedDeadlines',
  view: 'internhub:view'
};

export function readJson(key, fallback) {
  try {
    const raw = localStorage.getItem(key);
    return raw ? JSON.parse(raw) : fallback;
  } catch {
    return fallback;
  }
}

export function writeJson(key, value) {
  try {
    localStorage.setItem(key, JSON.stringify(value));
  } catch {}
}

export function readSet(key) {
  return new Set(readJson(key, []));
}

export function setToArray(setValue) {
  return Array.from(setValue || []);
}

export function normalizeText(value) {
  return String(value || '').toLowerCase().trim();
}

export function safeArray(value) {
  return Array.isArray(value) ? value : [];
}

export function parseDeadline(value) {
  if (!value) return null;
  const raw = String(value).trim();
  const native = new Date(raw);
  if (!Number.isNaN(native.getTime())) return native;

  const dotMatch = raw.match(/^(\d{1,2})[.\/-](\d{1,2})[.\/-](\d{2,4})$/);
  if (dotMatch) {
    const [, d, m, y] = dotMatch;
    const fullYear = y.length === 2 ? `20${y}` : y;
    const date = new Date(`${fullYear}-${m.padStart(2, '0')}-${d.padStart(2, '0')}T00:00:00`);
    if (!Number.isNaN(date.getTime())) return date;
  }
  return null;
}

export function formatDeadline(value) {
  const date = value instanceof Date ? value : parseDeadline(value);
  if (!date) return String(value || 'Не указан');
  return date.toLocaleDateString('ru-RU', { day: '2-digit', month: 'short', year: 'numeric' });
}

export function daysUntil(value) {
  const date = value instanceof Date ? value : parseDeadline(value);
  if (!date) return null;
  const start = new Date();
  start.setHours(0, 0, 0, 0);
  const target = new Date(date);
  target.setHours(0, 0, 0, 0);
  return Math.ceil((target - start) / 86400000);
}

export function deadlineMeta(value) {
  const days = daysUntil(value);
  if (days === null) return { label: 'Без дедлайна', tone: 'neutral' };
  if (days < 0) return { label: 'Дедлайн прошёл', tone: 'danger' };
  if (days === 0) return { label: 'Сегодня', tone: 'danger' };
  if (days <= 3) return { label: `${days} дн. осталось`, tone: 'warn' };
  return { label: `${days} дн. осталось`, tone: 'ok' };
}

export function calculateMatchScore(techStack, profileSkills) {
  const jobSkills = safeArray(techStack).map(normalizeText).filter(Boolean);
  const profile = safeArray(profileSkills).map(normalizeText).filter(Boolean);
  if (!jobSkills.length || !profile.length) return 0;
  const matched = jobSkills.filter((skill) => profile.includes(skill)).length;
  return Math.round((matched / jobSkills.length) * 100);
}

export function getTopTech(items, limit = 6) {
  const counter = new Map();
  safeArray(items).forEach((item) => {
    safeArray(item.techStack).forEach((tech) => {
      counter.set(tech, (counter.get(tech) || 0) + 1);
    });
  });
  return Array.from(counter.entries())
    .sort((a, b) => b[1] - a[1])
    .slice(0, limit)
    .map(([name, value]) => ({ name, value }));
}

export function getSalaryStats(items) {
  const salaries = safeArray(items).map((item) => Number(item.minSalary || 0)).filter((n) => n > 0);
  if (!salaries.length) return { avg: 0, max: 0, min: 0 };
  const total = salaries.reduce((sum, val) => sum + val, 0);
  return {
    avg: Math.round(total / salaries.length),
    min: Math.min(...salaries),
    max: Math.max(...salaries)
  };
}

export function requestDeadlinePermission() {
  if (typeof Notification === 'undefined') return Promise.resolve('unsupported');
  if (Notification.permission === 'granted') return Promise.resolve('granted');
  return Notification.requestPermission();
}

export function notifyUpcomingDeadlines(items, alreadyNotified, onNotified) {
  if (typeof Notification === 'undefined' || Notification.permission !== 'granted') return;
  safeArray(items).forEach((item) => {
    const days = daysUntil(item.applicationDeadline);
    if (days === null || days < 0 || days > 3) return;
    const key = `${item.id}:${days}`;
    if (alreadyNotified.has(key)) return;

    const body = days === 0
      ? `Сегодня дедлайн у ${item.companyName}`
      : `До дедлайна ${days} дн. · ${item.companyName}`;

    try {
      new Notification(`⏰ ${item.positionName}`, { body, tag: key });
      onNotified?.(key);
    } catch {}
  });
}

export function playFeedback(type = 'like') {
  try {
    if (navigator.vibrate) {
      navigator.vibrate(type === 'super' ? [20, 35, 20] : type === 'skip' ? [18] : [25]);
    }
  } catch {}

  try {
    const AudioContextClass = window.AudioContext || window.webkitAudioContext;
    if (!AudioContextClass) return;
    const ctx = new AudioContextClass();
    const oscillator = ctx.createOscillator();
    const gain = ctx.createGain();
    oscillator.connect(gain);
    gain.connect(ctx.destination);
    oscillator.type = 'sine';
    oscillator.frequency.value = type === 'super' ? 880 : type === 'skip' ? 220 : 660;
    gain.gain.setValueAtTime(0.0001, ctx.currentTime);
    gain.gain.exponentialRampToValueAtTime(0.06, ctx.currentTime + 0.01);
    gain.gain.exponentialRampToValueAtTime(0.0001, ctx.currentTime + 0.12);
    oscillator.start();
    oscillator.stop(ctx.currentTime + 0.13);
    oscillator.onended = () => ctx.close();
  } catch {}
}
