import { apiFetch } from './client'

/**
 * @param {string} username
 * @param {string} password
 */
export async function loginRequest(username, password) {
  const res = await apiFetch('/api/auth/login', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) {
    const err = new Error(typeof data.error === 'string' ? data.error : 'Login failed')
    err.status = res.status
    err.data = data
    throw err
  }
  return data
}
