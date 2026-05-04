<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { listProducts, searchProducts, askAssistant } from '../api/catalog'
import { createOrder } from '../api/orders'
import { useAuth } from '../composables/useAuth'

const { isLoggedIn } = useAuth()

// ── Order modal ───────────────────────────────────────────────────────────────
const orderModal = ref(null)   // { product }
const orderQty = ref(1)
const orderNotes = ref('')
const orderLoading = ref(false)
const orderSuccess = ref(null) // success message string
const orderError = ref('')

function openOrder(product) {
  orderModal.value = { product }
  orderQty.value = 1
  orderNotes.value = ''
  orderSuccess.value = null
  orderError.value = ''
}

function closeOrder() {
  orderModal.value = null
}

async function submitOrder() {
  orderLoading.value = true
  orderError.value = ''
  try {
    const resp = await createOrder({
      productId: orderModal.value.product.productId,
      quantity: orderQty.value,
      notes: orderNotes.value,
    })
    orderSuccess.value = resp.message
  } catch (err) {
    orderError.value = err.message || 'Failed to place order'
  } finally {
    orderLoading.value = false
  }
}

const categories = [
  { value: '', label: 'All' },
  { value: 'frame', label: 'Frame Gates' },
  { value: 'crown', label: 'Crown Panels' },
  { value: 'mdf', label: 'MDF' },
  { value: 'glass', label: 'Glass' },
  { value: 'fittings', label: 'Fittings' },
  { value: 'cutting', label: 'Cutting' },
]

const selectedCategory = ref('')
const searchQuery = ref('')
const products = ref([])
const searchResults = ref([])
const loadingProducts = ref(false)
const loadingSearch = ref(false)
const searchError = ref('')

// AI Assistant
const assistQuery = ref('')
const assistResult = ref(null)
const loadingAssist = ref(false)
const assistError = ref('')

const isSearchMode = computed(() => searchResults.value.length > 0 || (searchQuery.value.trim().length > 0 && !loadingSearch.value))

async function loadProducts() {
  loadingProducts.value = true
  try {
    products.value = await listProducts(selectedCategory.value)
  } catch {
    products.value = []
  } finally {
    loadingProducts.value = false
  }
}

async function runSearch() {
  const q = searchQuery.value.trim()
  if (!q) {
    searchResults.value = []
    return
  }
  searchError.value = ''
  loadingSearch.value = true
  try {
    searchResults.value = await searchProducts(q, selectedCategory.value)
  } catch (err) {
    searchError.value = err.message || 'Search failed'
    searchResults.value = []
  } finally {
    loadingSearch.value = false
  }
}

async function runAssist() {
  const q = assistQuery.value.trim()
  if (!q) return
  assistError.value = ''
  assistResult.value = null
  loadingAssist.value = true
  try {
    assistResult.value = await askAssistant(q, selectedCategory.value)
  } catch (err) {
    assistError.value = err.message || 'AI assistant unavailable'
  } finally {
    loadingAssist.value = false
  }
}

function clearSearch() {
  searchQuery.value = ''
  searchResults.value = ''
  searchError.value = ''
}

watch(selectedCategory, () => {
  clearSearch()
  loadProducts()
})

onMounted(loadProducts)

const displayedProducts = computed(() =>
  isSearchMode.value
    ? searchResults.value.map(r => r.product ?? r)
    : products.value
)

const categoryColorMap = {
  frame: 'bg-violet-100 text-violet-700',
  crown: 'bg-teal-100 text-teal-700',
  mdf: 'bg-emerald-100 text-emerald-700',
  glass: 'bg-sky-100 text-sky-700',
  fittings: 'bg-indigo-100 text-indigo-700',
  cutting: 'bg-amber-100 text-amber-700',
}
function categoryColor(cat) {
  return categoryColorMap[cat] || 'bg-gray-100 text-gray-700'
}
</script>

<template>
  <div class="catalog-page">
    <!-- Header -->
    <div class="catalog-header">
      <h1 class="catalog-title">Product Catalog</h1>
      <p class="catalog-subtitle">Browse gates, crowns, and accessories — or search semantically.</p>
    </div>

    <!-- Category tabs -->
    <div class="category-tabs">
      <button
        v-for="cat in categories"
        :key="cat.value"
        class="cat-tab"
        :class="{ active: selectedCategory === cat.value }"
        @click="selectedCategory = cat.value"
      >
        {{ cat.label }}
      </button>
    </div>

    <!-- AI Semantic Search -->
    <div class="search-box">
      <div class="search-row">
        <input
          v-model="searchQuery"
          type="text"
          class="search-input"
          placeholder="Search by meaning — e.g. 'heavy gate for warehouse entrance'"
          @keyup.enter="runSearch"
        />
        <button class="search-btn" :disabled="loadingSearch" @click="runSearch">
          <span v-if="loadingSearch" class="btn-spinner" />
          <span v-else>Search</span>
        </button>
        <button v-if="searchQuery" class="clear-btn" @click="clearSearch">Clear</button>
      </div>
      <p class="search-hint">Powered by vector embeddings — understands natural language queries.</p>
      <p v-if="searchError" class="search-error">{{ searchError }}</p>
    </div>

    <!-- Product grid -->
    <div v-if="loadingProducts && !isSearchMode" class="loading-state">
      Loading products…
    </div>
    <div v-else-if="displayedProducts.length === 0" class="empty-state">
      <span v-if="isSearchMode">No matching products found. Try a different query.</span>
      <span v-else>No products in this category.</span>
    </div>
    <div v-else class="product-grid">
      <div v-for="product in displayedProducts" :key="product.productId" class="product-card">
        <div class="product-top">
          <span class="product-badge" :class="categoryColor(product.category)">{{ product.category }}</span>
          <span v-if="product.priceRange" class="product-price">{{ product.priceRange }}</span>
        </div>
        <h3 class="product-name">{{ product.name }}</h3>
        <p class="product-desc">{{ product.description }}</p>
        <div class="product-meta">
          <span v-if="product.material" class="meta-chip">{{ product.material }}</span>
          <span v-if="product.dimensions" class="meta-chip">{{ product.dimensions }}</span>
        </div>
        <button v-if="isLoggedIn" class="order-btn" @click="openOrder(product)">Place Order</button>
        <RouterLink v-else :to="{ name: 'Login' }" class="order-btn order-btn-login">Sign in to order</RouterLink>
      </div>
    </div>

    <!-- Order modal -->
    <div v-if="orderModal" class="modal-overlay" @click.self="closeOrder">
      <div class="modal-box">
        <div class="modal-header">
          <h3 class="modal-title">Place Order</h3>
          <button class="modal-close" @click="closeOrder">✕</button>
        </div>

        <div v-if="orderSuccess" class="order-success">
          <div class="success-icon">✓</div>
          <p class="success-title">Order sent!</p>
          <p class="success-msg">{{ orderSuccess }}</p>
          <button class="order-btn" style="margin-top:1rem;" @click="closeOrder">Close</button>
        </div>

        <template v-else>
          <p class="modal-product-name">{{ orderModal.product.name }}</p>
          <p v-if="orderModal.product.priceRange" class="modal-price">{{ orderModal.product.priceRange }}</p>

          <div class="modal-field">
            <label class="modal-label">Quantity</label>
            <input v-model.number="orderQty" type="number" min="1" class="modal-input" />
          </div>
          <div class="modal-field">
            <label class="modal-label">Notes <span class="optional">(optional)</span></label>
            <textarea v-model="orderNotes" rows="3" class="modal-input" placeholder="Dimensions, colour, delivery requirements…" />
          </div>

          <p v-if="orderError" class="order-error">{{ orderError }}</p>

          <div class="modal-actions">
            <button class="cancel-btn" @click="closeOrder">Cancel</button>
            <button class="order-btn" :disabled="orderLoading || orderQty < 1" @click="submitOrder">
              <span v-if="orderLoading" class="btn-spinner" />
              <span v-else>Send Order</span>
            </button>
          </div>
        </template>
      </div>
    </div>

    <!-- AI Order Assistant -->
    <div class="assistant-section">
      <div class="assistant-header">
        <div class="assistant-icon">AI</div>
        <div>
          <h2 class="assistant-title">Order Assistant</h2>
          <p class="assistant-subtitle">Describe what you need and Claude will configure a recommendation.</p>
        </div>
      </div>
      <div class="assistant-input-row">
        <input
          v-model="assistQuery"
          type="text"
          class="search-input"
          placeholder="e.g. 'I need a secure automated gate for a factory with 5m opening'"
          @keyup.enter="runAssist"
        />
        <button class="search-btn" :disabled="loadingAssist || !assistQuery.trim()" @click="runAssist">
          <span v-if="loadingAssist" class="btn-spinner" />
          <span v-else>Ask AI</span>
        </button>
      </div>
      <p v-if="assistError" class="search-error">{{ assistError }}</p>
      <div v-if="assistResult" class="assist-result">
        <p class="assist-rec">{{ assistResult.recommendation }}</p>
        <div class="assist-grid">
          <div v-if="assistResult.suggested_category" class="assist-item">
            <span class="assist-label">Category</span>
            <span class="product-badge" :class="categoryColor(assistResult.suggested_category)">{{ assistResult.suggested_category }}</span>
          </div>
          <div v-if="assistResult.estimated_price_range" class="assist-item">
            <span class="assist-label">Est. price</span>
            <span class="assist-value">{{ assistResult.estimated_price_range }}</span>
          </div>
        </div>
        <div v-if="assistResult.key_features?.length" class="assist-features">
          <span class="assist-label">Key features</span>
          <ul class="feature-list">
            <li v-for="f in assistResult.key_features" :key="f">{{ f }}</li>
          </ul>
        </div>
        <div v-if="assistResult.next_steps?.length" class="assist-features">
          <span class="assist-label">Next steps</span>
          <ul class="feature-list">
            <li v-for="s in assistResult.next_steps" :key="s">{{ s }}</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.catalog-page { padding: 1rem 0; }
.catalog-header { margin-bottom: 1.5rem; }
.catalog-title { font-size: 1.5rem; font-weight: 700; color: #111827; margin: 0 0 0.25rem; }
.catalog-subtitle { font-size: 0.9375rem; color: #6b7280; margin: 0; }

.category-tabs { display: flex; gap: 0.5rem; flex-wrap: wrap; margin-bottom: 1.25rem; }
.cat-tab {
  padding: 0.375rem 0.875rem; border-radius: 9999px; font-size: 0.875rem; font-weight: 500;
  background: #f3f4f6; color: #374151; border: 1px solid transparent; cursor: pointer; transition: all 0.15s;
}
.cat-tab:hover { background: #e5e7eb; }
.cat-tab.active { background: #059669; color: #fff; border-color: #059669; }

.search-box { background: #f9fafb; border: 1px solid #e5e7eb; border-radius: 12px; padding: 1rem; margin-bottom: 1.5rem; }
.search-row { display: flex; gap: 0.5rem; margin-bottom: 0.5rem; }
.search-input {
  flex: 1; padding: 0.625rem 0.875rem; font-size: 0.9375rem; border: 1px solid #d1d5db;
  border-radius: 8px; outline: none; transition: border-color 0.2s;
}
.search-input:focus { border-color: #059669; box-shadow: 0 0 0 3px rgba(5,150,105,0.1); }
.search-btn {
  padding: 0.625rem 1.25rem; background: #059669; color: #fff; border: none; border-radius: 8px;
  font-weight: 600; cursor: pointer; transition: background 0.15s; display: flex; align-items: center; gap: 0.5rem; min-width: 80px; justify-content: center;
}
.search-btn:hover:not(:disabled) { background: #047857; }
.search-btn:disabled { opacity: 0.7; cursor: not-allowed; }
.clear-btn { padding: 0.625rem 0.875rem; background: #f3f4f6; border: 1px solid #d1d5db; border-radius: 8px; cursor: pointer; font-size: 0.875rem; color: #6b7280; }
.clear-btn:hover { background: #e5e7eb; }
.search-hint { font-size: 0.8125rem; color: #9ca3af; margin: 0; }
.search-error { font-size: 0.875rem; color: #ef4444; margin-top: 0.5rem; }

.loading-state, .empty-state { text-align: center; color: #9ca3af; padding: 3rem 1rem; font-size: 0.9375rem; }

.product-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(240px, 1fr)); gap: 1rem; margin-bottom: 2rem; }
.product-card {
  background: #fff; border: 1px solid #e5e7eb; border-radius: 12px; padding: 1rem;
  transition: box-shadow 0.2s, transform 0.15s;
}
.product-card:hover { box-shadow: 0 4px 12px rgba(0,0,0,0.08); transform: translateY(-2px); }
.product-top { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.625rem; }
.product-badge { font-size: 0.75rem; font-weight: 600; padding: 0.2rem 0.5rem; border-radius: 9999px; }
.product-price { font-size: 0.8125rem; font-weight: 600; color: #059669; }
.product-name { font-size: 0.9375rem; font-weight: 600; color: #111827; margin: 0 0 0.375rem; }
.product-desc { font-size: 0.8125rem; color: #6b7280; margin: 0 0 0.75rem; line-height: 1.5; }
.product-meta { display: flex; gap: 0.375rem; flex-wrap: wrap; }
.meta-chip { font-size: 0.75rem; background: #f3f4f6; color: #4b5563; padding: 0.2rem 0.5rem; border-radius: 6px; }

.assistant-section { background: linear-gradient(135deg, #ecfdf5, #f0fdf4); border: 1px solid #a7f3d0; border-radius: 16px; padding: 1.5rem; }
.assistant-header { display: flex; gap: 1rem; align-items: flex-start; margin-bottom: 1rem; }
.assistant-icon { width: 44px; height: 44px; background: linear-gradient(135deg, #10b981, #059669); border-radius: 10px; display: flex; align-items: center; justify-content: center; color: #fff; font-weight: 700; font-size: 0.875rem; flex-shrink: 0; }
.assistant-title { font-size: 1rem; font-weight: 700; color: #064e3b; margin: 0 0 0.2rem; }
.assistant-subtitle { font-size: 0.875rem; color: #065f46; margin: 0; }
.assistant-input-row { display: flex; gap: 0.5rem; }

.assist-result { margin-top: 1rem; background: #fff; border: 1px solid #d1fae5; border-radius: 10px; padding: 1rem; }
.assist-rec { font-size: 0.9375rem; color: #065f46; margin: 0 0 0.75rem; line-height: 1.6; }
.assist-grid { display: flex; gap: 1rem; flex-wrap: wrap; margin-bottom: 0.75rem; }
.assist-item { display: flex; flex-direction: column; gap: 0.25rem; }
.assist-label { font-size: 0.75rem; font-weight: 600; color: #6b7280; text-transform: uppercase; letter-spacing: 0.05em; }
.assist-value { font-size: 0.9375rem; font-weight: 600; color: #059669; }
.assist-features { margin-bottom: 0.75rem; }
.feature-list { margin: 0.25rem 0 0 1rem; padding: 0; font-size: 0.875rem; color: #374151; line-height: 1.8; }

.btn-spinner { width: 18px; height: 18px; border: 2px solid rgba(255,255,255,0.4); border-top-color: #fff; border-radius: 50%; animation: spin 0.7s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.order-btn {
  display: flex; align-items: center; justify-content: center; gap: 0.4rem;
  margin-top: 0.75rem; width: 100%; padding: 0.5rem 1rem;
  background: #059669; color: #fff; border: none; border-radius: 8px;
  font-size: 0.875rem; font-weight: 600; cursor: pointer; transition: background 0.15s;
  text-decoration: none;
}
.order-btn:hover:not(:disabled) { background: #047857; }
.order-btn:disabled { opacity: 0.7; cursor: not-allowed; }
.order-btn-login { background: #6b7280; }
.order-btn-login:hover { background: #4b5563; }

.modal-overlay {
  position: fixed; inset: 0; background: rgba(0,0,0,0.45); z-index: 1000;
  display: flex; align-items: center; justify-content: center; padding: 1rem;
}
.modal-box {
  background: #fff; border-radius: 16px; padding: 1.5rem; width: 100%; max-width: 420px;
  box-shadow: 0 20px 50px rgba(0,0,0,0.2);
}
.modal-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem; }
.modal-title { font-size: 1.125rem; font-weight: 700; color: #111827; margin: 0; }
.modal-close { background: none; border: none; font-size: 1.125rem; color: #9ca3af; cursor: pointer; padding: 0.25rem; }
.modal-close:hover { color: #374151; }
.modal-product-name { font-size: 0.9375rem; font-weight: 600; color: #059669; margin: 0 0 0.2rem; }
.modal-price { font-size: 0.8125rem; color: #6b7280; margin: 0 0 1rem; }
.modal-field { margin-bottom: 0.875rem; }
.modal-label { display: block; font-size: 0.8125rem; font-weight: 600; color: #374151; margin-bottom: 0.35rem; }
.optional { font-weight: 400; color: #9ca3af; }
.modal-input {
  width: 100%; padding: 0.5rem 0.75rem; font-size: 0.9375rem; border: 1px solid #d1d5db;
  border-radius: 8px; outline: none; box-sizing: border-box; font-family: inherit; resize: vertical;
  transition: border-color 0.2s;
}
.modal-input:focus { border-color: #059669; box-shadow: 0 0 0 3px rgba(5,150,105,0.1); }
.modal-actions { display: flex; gap: 0.5rem; margin-top: 1rem; }
.modal-actions .order-btn { margin-top: 0; flex: 1; }
.cancel-btn {
  flex: 1; padding: 0.5rem 1rem; background: #f3f4f6; color: #374151; border: 1px solid #d1d5db;
  border-radius: 8px; font-size: 0.875rem; font-weight: 600; cursor: pointer;
}
.cancel-btn:hover { background: #e5e7eb; }
.order-error { font-size: 0.875rem; color: #ef4444; margin: 0.25rem 0 0; }

.order-success { text-align: center; padding: 1rem 0; }
.success-icon { width: 52px; height: 52px; background: #d1fae5; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 1.5rem; color: #059669; margin: 0 auto 1rem; }
.success-title { font-size: 1.125rem; font-weight: 700; color: #111827; margin: 0 0 0.5rem; }
.success-msg { font-size: 0.9375rem; color: #4b5563; line-height: 1.6; margin: 0; }

/* Category badge colors */
.bg-violet-100 { background: #ede9fe; }
.text-violet-700 { color: #6d28d9; }
.bg-teal-100 { background: #ccfbf1; }
.text-teal-700 { color: #0f766e; }
.bg-emerald-100 { background: #d1fae5; }
.text-emerald-700 { color: #047857; }
.bg-sky-100 { background: #e0f2fe; }
.text-sky-700 { color: #0369a1; }
.bg-indigo-100 { background: #e0e7ff; }
.text-indigo-700 { color: #4338ca; }
.bg-amber-100 { background: #fef3c7; }
.text-amber-700 { color: #b45309; }
.bg-gray-100 { background: #f3f4f6; }
.text-gray-700 { color: #374151; }
</style>
