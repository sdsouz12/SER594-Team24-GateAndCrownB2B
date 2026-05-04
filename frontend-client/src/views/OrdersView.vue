<script setup>
import { ref, onMounted } from 'vue'
import { getMyOrders } from '../api/orders'

const orders = ref([])
const loading = ref(true)
const error = ref('')

const statusLabel = {
  pending: 'Pending review',
  confirmed: 'Confirmed',
  processing: 'Processing',
  completed: 'Completed',
  cancelled: 'Cancelled',
}

const statusColor = {
  pending: 'status-pending',
  confirmed: 'status-confirmed',
  processing: 'status-processing',
  completed: 'status-completed',
  cancelled: 'status-cancelled',
}

function formatDate(d) {
  return new Date(d).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

onMounted(async () => {
  try {
    orders.value = await getMyOrders()
  } catch (err) {
    error.value = err.message || 'Failed to load orders'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="orders-page">
    <div class="orders-header">
      <h1 class="orders-title">My Orders</h1>
      <p class="orders-subtitle">Track all your product requests.</p>
    </div>

    <div v-if="loading" class="empty-state">Loading orders…</div>
    <div v-else-if="error" class="error-state">{{ error }}</div>

    <div v-else-if="orders.length === 0" class="empty-state">
      <div class="empty-icon">📋</div>
      <p class="empty-title">No orders yet</p>
      <p class="empty-sub">Go to the <RouterLink to="/catalog" class="link">Catalog</RouterLink> and place your first order.</p>
    </div>

    <div v-else class="orders-list">
      <div v-for="order in orders" :key="order.orderId" class="order-card">
        <div class="order-top">
          <div>
            <p class="order-product">{{ order.productName }}</p>
            <p class="order-meta">Qty: {{ order.quantity }} · Order #{{ order.orderId }}</p>
          </div>
          <span class="order-status" :class="statusColor[order.status] || 'status-pending'">
            {{ statusLabel[order.status] || order.status }}
          </span>
        </div>
        <p v-if="order.notes" class="order-notes">Note: {{ order.notes }}</p>
        <p class="order-date">{{ formatDate(order.createdAt) }}</p>
        <div v-if="order.status === 'pending'" class="order-notice">
          Your order has been received — our team will review it and get back to you shortly.
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.orders-page { padding: 1rem 0; }
.orders-header { margin-bottom: 1.5rem; }
.orders-title { font-size: 1.5rem; font-weight: 700; color: #111827; margin: 0 0 0.25rem; }
.orders-subtitle { font-size: 0.9375rem; color: #6b7280; margin: 0; }

.empty-state { text-align: center; padding: 4rem 1rem; color: #9ca3af; }
.empty-icon { font-size: 2.5rem; margin-bottom: 0.75rem; }
.empty-title { font-size: 1rem; font-weight: 600; color: #374151; margin: 0 0 0.4rem; }
.empty-sub { font-size: 0.9375rem; color: #6b7280; margin: 0; }
.link { color: #059669; font-weight: 600; text-decoration: none; }
.link:hover { text-decoration: underline; }
.error-state { color: #ef4444; text-align: center; padding: 2rem; }

.orders-list { display: flex; flex-direction: column; gap: 1rem; }

.order-card {
  background: #fff; border: 1px solid #e5e7eb; border-radius: 12px; padding: 1.25rem;
  transition: box-shadow 0.2s;
}
.order-card:hover { box-shadow: 0 4px 12px rgba(0,0,0,0.06); }

.order-top { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 0.5rem; }
.order-product { font-size: 1rem; font-weight: 600; color: #111827; margin: 0 0 0.2rem; }
.order-meta { font-size: 0.8125rem; color: #6b7280; margin: 0; }
.order-notes { font-size: 0.875rem; color: #4b5563; margin: 0.25rem 0; font-style: italic; }
.order-date { font-size: 0.8125rem; color: #9ca3af; margin: 0.25rem 0 0; }

.order-status {
  font-size: 0.75rem; font-weight: 600; padding: 0.25rem 0.6rem; border-radius: 9999px;
  white-space: nowrap; flex-shrink: 0; margin-left: 0.75rem;
}
.status-pending    { background: #fef3c7; color: #b45309; }
.status-confirmed  { background: #dbeafe; color: #1d4ed8; }
.status-processing { background: #ede9fe; color: #6d28d9; }
.status-completed  { background: #d1fae5; color: #047857; }
.status-cancelled  { background: #fee2e2; color: #b91c1c; }

.order-notice {
  margin-top: 0.75rem; padding: 0.6rem 0.875rem; background: #ecfdf5;
  border: 1px solid #a7f3d0; border-radius: 8px; font-size: 0.8125rem; color: #065f46; line-height: 1.5;
}
</style>
