const DEFAULT_BASE_URL = (typeof import.meta !== 'undefined' && import.meta.env && (import.meta.env.VITE_API_URL || import.meta.env.REACT_APP_API_URL)) || 'http://localhost:8082';

export const BASE_URL = DEFAULT_BASE_URL.replace(/\/$/, '');

function buildUrl(path, params) {
  const url = new URL(`${BASE_URL}${path}`);
  Object.entries(params || {}).forEach(([key, value]) => {
    if (value == null || value === '' || (Array.isArray(value) && !value.length)) return;
    if (Array.isArray(value)) {
      value.forEach((entry) => url.searchParams.append(key, entry));
      return;
    }
    url.searchParams.set(key, String(value));
  });
  return url.toString();
}

function makeFallbackId(item) {
  const raw = [
    item?.companyName || 'company',
    item?.positionName || 'position',
    item?.location || '',
    item?.sourceUrl || '',
    item?.applicationDeadline || ''
  ].join('::');

  return raw
    .toLowerCase()
    .replace(/[^a-zа-я0-9]+/gi, '-')
    .replace(/^-+|-+$/g, '') || 'internship-fallback';
}

function normalizeInternship(item) {
  return {
    id: item?.id ?? makeFallbackId(item),
    positionName: item?.positionName || 'Без названия',
    companyName: item?.companyName || 'Неизвестная компания',
    techStack: Array.isArray(item?.techStack)
      ? item.techStack.filter(Boolean)
      : typeof item?.techStack === 'string'
        ? item.techStack.split(',').map((entry) => entry.trim()).filter(Boolean)
        : [],
    minSalary: Number(item?.minSalary || 0),
    location: item?.location || '',
    internshipDates: item?.internshipDates || '',
    selectionProcess: item?.selectionProcess || '',
    description: item?.description || '',
    applicationDeadline: item?.applicationDeadline || '',
    contactInfo: item?.contactInfo || '',
    experienceRequirements: item?.experienceRequirements || '',
    sourceUrl: item?.sourceUrl || '',
    sourceSite: item?.sourceSite || ''
  };
}

async function fetchJson(url, signal) {
  const res = await fetch(url, { signal });
  if (!res.ok) {
    const text = await res.text().catch(() => '');
    throw new Error(`HTTP ${res.status}${text ? ` · ${text}` : ''}`);
  }
  return res.json();
}

async function fetchPaged(path, params = {}, signal) {
  const first = await fetchJson(buildUrl(path, { page: 0, size: 100, ...params }), signal);

  if (Array.isArray(first)) {
    return first.map(normalizeInternship);
  }

  const firstContent = Array.isArray(first?.content) ? first.content.map(normalizeInternship) : [];
  const totalPages = Number(first?.totalPages ?? 1);

  if (totalPages <= 1) return firstContent;

  const rest = await Promise.all(
    Array.from({ length: totalPages - 1 }, (_, idx) =>
      fetchJson(buildUrl(path, { page: idx + 1, size: 100, ...params }), signal)
    )
  );

  return [
    ...firstContent,
    ...rest.flatMap((page) => (Array.isArray(page) ? page : page?.content || []).map(normalizeInternship))
  ];
}

export function fetchInternships(signal) {
  return fetchPaged('/api/internship/all', {}, signal);
}

export function searchInternships(filters = {}, signal) {
  const params = {
    page: 0,
    size: 100,
    companyName: filters.companyName,
    location: filters.location,
    minSalary: filters.minSalary && Number(filters.minSalary) > 0 ? Number(filters.minSalary) : undefined,
    tech: Array.isArray(filters.tech) ? filters.tech : filters.tech ? [filters.tech] : undefined
  };


  if (!filters.companyName && !filters.location && !filters.minSalary && !(Array.isArray(filters.tech) ? filters.tech.length : filters.tech)) {
    return fetchInternships(signal);
  }

  return fetchPaged('/api/internship/', params, signal);
}

export async function fetchTechOptions(signal) {
  const data = await fetchJson(buildUrl('/api/internship/tech', {}), signal);
  const options = Array.isArray(data) ? data : [];
  return options
    .map((entry) => {
      if (typeof entry === 'string') return entry;
      return entry?.technology || entry?.tech || entry?.name || '';
    })
    .map((entry) => String(entry || '').trim())
    .filter(Boolean)
    .sort((a, b) => a.localeCompare(b, 'ru'));
}
