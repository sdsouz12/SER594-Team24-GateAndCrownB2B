/**
 * @param {string} path
 */
export function apiUrl(path) {
  const base = String(import.meta.env.VITE_API_BASE_URL || '').replace(/\/$/, '')
  const p = path.startsWith('/') ? path : `/${path}`
  return `${base}${p}`
}

export async function apiFetch(path, options = {}) {
  const headers = new Headers(options.headers)
  if (
    options.body != null &&
    typeof options.body === 'string' &&
    !headers.has('Content-Type')
  ) {
    headers.set('Content-Type', 'application/json')
  }
  try {
    const token = localStorage.getItem('gate_crown_client_token')
    if (token && !headers.has('Authorization')) {
      headers.set('Authorization', `Bearer ${token}`)
    }
  } catch {}
  return fetch(apiUrl(path), { ...options, headers })
}
