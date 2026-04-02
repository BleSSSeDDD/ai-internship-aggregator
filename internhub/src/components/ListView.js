import React, { useEffect, useMemo, useRef, useState } from 'react';
import { getColor, getEmoji } from '../data';
import { fetchTechOptions, searchInternships } from '../api';
import { calculateMatchScore, deadlineMeta, normalizeText, readJson, STORAGE_KEYS, writeJson } from '../utils';
import './ListView.css';

const defaultFilters = {
  companyName: '',
  location: '',
  tech: '',
  minSalary: 0,
  quickText: '',
  scope: 'all',
  sortBy: 'match',
  advancedOpen: false
};

export default function ListView({ internships, liked, skipped, toggleLike, markSkipped, onOpenDetails, profileSkills, setProfileSkills }) {
  const [filters, setFilters] = useState(() => readJson(STORAGE_KEYS.listFilters, defaultFilters));
  const [skillsInput, setSkillsInput] = useState(profileSkills.join(', '));
  const [techOptions, setTechOptions] = useState([]);
  const [backendResults, setBackendResults] = useState(internships);
  const [backendLoading, setBackendLoading] = useState(false);
  const [backendError, setBackendError] = useState('');
  const requestIdRef = useRef(0);
  const abortRef = useRef(null);

  useEffect(() => {
    writeJson(STORAGE_KEYS.listFilters, filters);
  }, [filters]);

  useEffect(() => {
    setSkillsInput(profileSkills.join(', '));
  }, [profileSkills]);

  useEffect(() => {
    setBackendResults(internships);
  }, [internships]);

  useEffect(() => {
    const controller = new AbortController();
    fetchTechOptions(controller.signal)
      .then(setTechOptions)
      .catch(() => setTechOptions([]));
    return () => controller.abort();
  }, []);

  useEffect(() => () => {
    abortRef.current?.abort();
  }, []);

  const runBackendSearch = async () => {
    abortRef.current?.abort();
    const controller = new AbortController();
    abortRef.current = controller;
    const requestId = requestIdRef.current + 1;
    requestIdRef.current = requestId;
    setBackendLoading(true);
    setBackendError('');

    try {
      const data = await searchInternships({
        companyName: filters.companyName.trim(),
        location: filters.location.trim(),
        tech: filters.tech ? [filters.tech] : [],
        minSalary: Number(filters.minSalary || 0)
      }, controller.signal);

      if (requestIdRef.current !== requestId) return;
      setBackendResults(data);
    } catch (error) {
      if (requestIdRef.current !== requestId || error?.name === 'AbortError') return;
      setBackendError(error.message || 'Не удалось выполнить поиск');
    } finally {
      if (requestIdRef.current === requestId) {
        setBackendLoading(false);
        if (abortRef.current === controller) {
          abortRef.current = null;
        }
      }
    }
  };

  const resetBackendSearch = () => {
    requestIdRef.current += 1;
    abortRef.current?.abort();
    abortRef.current = null;
    setBackendLoading(false);
    setBackendError('');
    setBackendResults(internships);
    setFilters((prev) => ({
      ...prev,
      companyName: '',
      location: '',
      tech: '',
      minSalary: 0,
      quickText: ''
    }));
  };

  const filtered = useMemo(() => {
    const quickText = normalizeText(filters.quickText);
    const scoped = backendResults.filter((item) => {
      if (filters.scope === 'liked') return liked.has(item.id);
      if (filters.scope === 'skipped') return skipped.has(item.id);
      if (filters.scope === 'fresh') return !liked.has(item.id) && !skipped.has(item.id);
      return true;
    });

    const searched = scoped.filter((item) => {
      if (!quickText) return true;
      const haystack = [
        item.positionName,
        item.companyName,
        item.location,
        item.description,
        ...(item.techStack || [])
      ].map(normalizeText).join(' ');
      return haystack.includes(quickText);
    });

    return searched.slice().sort((a, b) => {
      if (filters.sortBy === 'salary') return Number(b.minSalary || 0) - Number(a.minSalary || 0);
      if (filters.sortBy === 'company') return String(a.companyName || '').localeCompare(String(b.companyName || ''), 'ru');
      if (filters.sortBy === 'title') return String(a.positionName || '').localeCompare(String(b.positionName || ''), 'ru');
      return calculateMatchScore(b.techStack, profileSkills) - calculateMatchScore(a.techStack, profileSkills);
    });
  }, [backendResults, filters.quickText, filters.scope, filters.sortBy, liked, skipped, profileSkills]);

  const activeBackendFilters = useMemo(() => ([
    filters.companyName ? `Компания: ${filters.companyName}` : '',
    filters.location ? `Локация: ${filters.location}` : '',
    filters.tech ? `Технология: ${filters.tech}` : '',
    Number(filters.minSalary || 0) > 0 ? `От ${Number(filters.minSalary).toLocaleString('ru-RU')} ₽` : ''
  ].filter(Boolean)), [filters.companyName, filters.location, filters.tech, filters.minSalary]);

  const saveSkills = () => {
    const next = skillsInput.split(',').map((value) => value.trim()).filter(Boolean);
    setProfileSkills(next);
  };

  return (
    <div className="list-view">
      <div className="list-view__top">
        <div>
          <h1 className="section-title">Все стажировки</h1>
          <p className="section-subtitle">Backend получает только то, что умеет фильтровать: компанию, локацию, технологию и мин. зарплату. Всё остальное — уже локально и без лишних запросов.</p>
        </div>
      </div>

      <div className="search-shell glass-card">
        <div className="search-shell__row search-shell__row--main">
          <div className="search-field search-field--grow">
            <label>Компания</label>
            <input
              type="text"
              placeholder="Яндекс, VK, Т-Банк..."
              value={filters.companyName}
              onChange={(e) => setFilters((prev) => ({ ...prev, companyName: e.target.value }))}
            />
          </div>

          <div className="search-field search-field--select">
            <label>Выбор вакансий</label>
            <select value={filters.scope} onChange={(e) => setFilters((prev) => ({ ...prev, scope: e.target.value }))}>
              <option value="all">Все результаты</option>
              <option value="fresh">Только неразобранные</option>
              <option value="liked">Только matches</option>
              <option value="skipped">Только история</option>
            </select>
          </div>

          <div className="search-field search-field--range">
            <label>Мин. зарплата</label>
            <input
              type="range"
              min="0"
              max="250000"
              step="10000"
              value={filters.minSalary}
              onChange={(e) => setFilters((prev) => ({ ...prev, minSalary: Number(e.target.value) }))}
            />
            <div className="range-value">{Number(filters.minSalary || 0).toLocaleString('ru-RU')} ₽</div>
          </div>

          <div className="search-actions">
            <button className="btn-primary" onClick={runBackendSearch} disabled={backendLoading}>
              {backendLoading ? 'Ищу...' : 'Искать в backend'}
            </button>
            <button className="btn-ghost" onClick={() => setFilters((prev) => ({ ...prev, advancedOpen: !prev.advancedOpen }))}>
              {filters.advancedOpen ? 'Скрыть поиск' : 'Другие параметры'}
            </button>
          </div>
        </div>

        {filters.advancedOpen && (
          <div className="search-shell__row search-shell__row--advanced">
            <div className="search-field">
              <label>Локация (backend)</label>
              <input
                type="text"
                placeholder="Москва, Remote..."
                value={filters.location}
                onChange={(e) => setFilters((prev) => ({ ...prev, location: e.target.value }))}
              />
            </div>

            <div className="search-field">
              <label>Технология (backend)</label>
              <select value={filters.tech} onChange={(e) => setFilters((prev) => ({ ...prev, tech: e.target.value }))}>
                <option value="">Любая</option>
                {techOptions.map((tech) => <option key={tech} value={tech}>{tech}</option>)}
              </select>
            </div>

            <div className="search-field search-field--grow">
              <label>Быстрый локальный поиск</label>
              <input
                type="text"
                placeholder="Название, описание, стек, город..."
                value={filters.quickText}
                onChange={(e) => setFilters((prev) => ({ ...prev, quickText: e.target.value }))}
              />
            </div>

            <div className="search-field">
              <label>Сортировка</label>
              <select value={filters.sortBy} onChange={(e) => setFilters((prev) => ({ ...prev, sortBy: e.target.value }))}>
                <option value="match">По match score</option>
                <option value="salary">По зарплате</option>
                <option value="company">По компании</option>
                <option value="title">По названию</option>
              </select>
            </div>

            <div className="search-field search-field--skills">
              <label>Мой стек</label>
              <div className="skills-row">
                <input value={skillsInput} onChange={(e) => setSkillsInput(e.target.value)} placeholder="React, Java, Spring, Kafka" />
                <button onClick={saveSkills}>Сохранить</button>
              </div>
            </div>
          </div>
        )}

        <div className="search-shell__footer">
          <div className="active-filter-list">
            {activeBackendFilters.length === 0 ? (
              <span className="active-filter-chip muted">Backend-фильтры не заданы</span>
            ) : activeBackendFilters.map((label) => (
              <span key={label} className="active-filter-chip">{label}</span>
            ))}
            {filters.quickText && <span className="active-filter-chip cyan">Локально: {filters.quickText}</span>}
          </div>

          <div className="search-shell__footer-actions">
            <span className="list-counter">Показано <strong>{filtered.length}</strong> из <strong>{backendResults.length}</strong></span>
            <button className="btn-link" onClick={resetBackendSearch}>Сбросить всё</button>
          </div>
        </div>

        {backendError && <div className="search-error">{backendError}</div>}
      </div>

      <div className="intern-grid">
        {filtered.length === 0 && <p className="empty-copy">Ничего не найдено — попробуйте ослабить backend-фильтры или очистить локальный поиск.</p>}
        {filtered.map((intern) => {
          const color = getColor(intern.companyName);
          const emoji = getEmoji(intern.companyName);
          const isLiked = liked.has(intern.id);
          const isSkipped = skipped.has(intern.id);
          const match = calculateMatchScore(intern.techStack, profileSkills);
          const deadline = deadlineMeta(intern.applicationDeadline);

          return (
            <div key={intern.id} className="intern-card glass-card">
              <div className="ic-head">
                <div className="ic-logo" style={{ background: color }}>{emoji}</div>
                <div>
                  <div className="ic-title">{intern.positionName}</div>
                  <div className="ic-company">{intern.companyName}</div>
                </div>
                <span className="ic-match">{match}%</span>
              </div>

              <div className="ic-meta-row">
                <span>📍 {intern.location || 'Не указана'}</span>
                <span>💰 от {(intern.minSalary || 0).toLocaleString('ru-RU')} ₽</span>
              </div>
              <p className="ic-desc">{intern.description || 'Описание пока не добавлено.'}</p>

              <div className="ic-tags">
                {(intern.techStack || []).slice(0, 5).map((tech) => <span key={tech} className="ic-tag">{tech}</span>)}
              </div>

              <div className="ic-foot">
                <div className={`status-badge ${deadline.tone}`}>{deadline.label}</div>

                <div className="ic-actions">
                  <button className="ic-action-btn" onClick={() => onOpenDetails(intern)}>Подробнее</button>
                  <button className={`ic-like-btn ${isLiked ? 'liked' : ''}`} onClick={() => toggleLike(intern.id, intern.positionName)}>
                    {isLiked ? '❤️' : '🤍'}
                  </button>
                  <button className={`ic-like-btn ${isSkipped ? 'skipped' : ''}`} onClick={() => markSkipped(intern.id, intern.positionName)}>
                    ⏭
                  </button>
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
