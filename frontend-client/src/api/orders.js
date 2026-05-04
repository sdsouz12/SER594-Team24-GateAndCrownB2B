import { apiFetch } from './client'

export async function createOrder({ productId, quantity, notes }) {
  const res = await apiFetch('/api/orders', {
    method: 'POST',
    body: JSON.stringify({ productId, quantity, notes: notes || '' }),
  })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) throw new Error(data.error || 'Failed to place order')
  return data.data
}

export async function getMyOrders() {
  const res = await apiFetch('/api/orders')
  const data = await res.json().catch(() => ({}))
  if (!res.ok) throw new Error(data.error || 'Failed to fetch orders')
  return data.data ?? []
}
