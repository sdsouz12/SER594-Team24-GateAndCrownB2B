import { apiFetch } from './client'
import { useAuth } from '../composables/useAuth'

function authHeaders() {
  const { accessToken } = _getToken()
  return accessToken ? { Authorization: `Bearer ${accessToken}` } : {}
}

// Read token directly from localStorage to avoid composable import issues outside setup
function _getToken() {
  try {
    return { accessToken: localStorage.getItem('gate_crown_client_token') }
  } catch {
    return { accessToken: null }
  }
}

export async function listProducts(category = '') {
  const url = category ? `/api/catalog/products?category=${encodeURIComponent(category)}` : '/api/catalog/products'
  const res = await apiFetch(url)
  const data = await res.json().catch(() => ({}))
  if (!res.ok) throw new Error(data.error || 'Failed to fetch products')
  return data.data ?? []
}

export async function searchProducts(query, category = '') {
  const res = await apiFetch('/api/catalog/search', {
    method: 'POST',
    body: JSON.stringify({ query, category: category || undefined, limit: 10 }),
  })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) throw new Error(data.error || 'Search failed')
  return data.data ?? []
}

export async function askAssistant(userMessage, category = '') {
  const res = await apiFetch('/api/ai/assist', {
    method: 'POST',
    body: JSON.stringify({ user_message: userMessage, category: category || undefined }),
  })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) throw new Error(data.error || 'Assistant unavailable')
  return data.data
}
