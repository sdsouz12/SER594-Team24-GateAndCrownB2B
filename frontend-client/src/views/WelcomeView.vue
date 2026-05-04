<script setup>
import { ref, onMounted } from 'vue'
import { useAuth } from '../composables/useAuth'
import { searchProducts, askAssistant } from '../api/catalog'
import { apiFetch } from '../api/client'

const { user } = useAuth()

// ── Pipeline status ───────────────────────────────────────────────────────────
const pipeline = ref(null)
const pipelineLoading = ref(false)

async function loadPipeline() {
  pipelineLoading.value = true
  try {
    const res = await apiFetch('/api/ai/pipeline/status')
    pipeline.value = await res.json()
  } catch {
    pipeline.value = { error: 'AI service not reachable' }
  } finally {
    pipelineLoading.value = false
  }
}

async function runIngest() {
  pipelineLoading.value = true
  try {
    await apiFetch('/api/ai/pipeline/ingest', {
      method: 'POST',
      body: '{}',
    })
    await loadPipeline()
  } finally {
    pipelineLoading.value = false
  }
}

// ── Vector Search ─────────────────────────────────────────────────────────────
const searchQuery = ref('heavy gate for warehouse entrance')
const searchResults = ref([])
const searchLoading = ref(false)
const searchError = ref('')
const searchDone = ref(false)

async function runSearch() {
  if (!searchQuery.value.trim()) return
  searchLoading.value = true
  searchError.value = ''
  searchDone.value = false
  try {
    searchResults.value = await searchProducts(searchQuery.value)
    searchDone.value = true
  } catch (err) {
    searchError.value = err.message
  } finally {
    searchLoading.value = false
  }
}

// ── AI Assistant ──────────────────────────────────────────────────────────────
const assistQuery = ref('I need a secure automated gate for a factory with 5 meter opening')
const assistResult = ref(null)
const assistLoading = ref(false)
const assistError = ref('')
const assistDone = ref(false)

async function runAssist() {
  assistLoading.value = true
  assistError.value = ''
  assistDone.value = false
  try {
    assistResult.value = await askAssistant(assistQuery.value)
    assistDone.value = true
  } catch (err) {
    assistError.value = err.message
  } finally {
    assistLoading.value = false
  }
}

// ── RAG ───────────────────────────────────────────────────────────────────────
const ragQuery = ref('What materials are available for crown panels?')
const ragResult = ref(null)
const ragLoading = ref(false)
const ragError = ref('')
const ragDone = ref(false)

async function runRag() {
  ragLoading.value = true
  ragError.value = ''
  ragDone.value = false
  try {
    const res = await apiFetch('/api/ai/rag', {
      method: 'POST',
      body: JSON.stringify({ question: ragQuery.value }),
    })
    ragResult.value = await res.json()
    ragDone.value = true
  } catch (err) {
    ragError.value = err.message
  } finally {
    ragLoading.value = false
  }
}

onMounted(loadPipeline)
</script>

<template>
  <div class="demo-page">

    <!-- Header -->
    <div class="demo-header">
      <div class="demo-badge">SER 594 — Final Milestone</div>
      <h1 class="demo-title">Gate & Crown B2B Platform</h1>
      <p class="demo-subtitle">
        Logged in as <strong>{{ user?.username }}</strong> ({{ user?.role }})
      </p>
    </div>

    <!-- Final milestone checklist -->
    <div class="checklist-card">
      <h2 class="section-title">Final Milestone — Feature Complete</h2>
      <ul class="checklist">
        <li class="check-item">
          <span class="check-icon">✅</span>
          <div>
            <strong>Authentication &amp; Registration</strong>
            <span class="check-note">JWT login, bcrypt register, role-based sessions (CLIENT / ADMIN / SUPERADMIN)</span>
          </div>
        </li>
        <li class="check-item">
          <span class="check-icon">✅</span>
          <div>
            <strong>AI Technique 1 — Vector Search</strong>
            <span class="check-note">BAAI/bge-small-en-v1.5 (ONNX) embeddings + pgvector cosine similarity</span>
          </div>
        </li>
        <li class="check-item">
          <span class="check-icon">✅</span>
          <div>
            <strong>AI Technique 2 — LLM Order Assistant</strong>
            <span class="check-note">Claude Haiku returns structured JSON product recommendations</span>
          </div>
        </li>
        <li class="check-item">
          <span class="check-icon">✅</span>
          <div>
            <strong>AI Technique 3 — RAG</strong>
            <span class="check-note">Top-5 catalog products retrieved as context → Claude generates grounded answers</span>
          </div>
        </li>
        <li class="check-item">
          <span class="check-icon">✅</span>
          <div>
            <strong>Order Management</strong>
            <span class="check-note">Users place orders from catalog; order list with status tracking</span>
          </div>
        </li>
        <li class="check-item">
          <span class="check-icon">✅</span>
          <div>
            <strong>Test Suite — 30 tests</strong>
            <span class="check-note">15 Go unit tests (auth service) + 15 Python unit tests (AI service); all mocked, no DB required</span>
          </div>
        </li>
        <li class="check-item">
          <span class="check-icon">✅</span>
          <div>
            <strong>Docker Deployment</strong>
            <span class="check-note">docker-compose up --build runs all 4 services (postgres/pgvector, backend, AI service, frontend)</span>
          </div>
        </li>
        <li class="check-item">
          <span class="check-icon">✅</span>
          <div>
            <strong>CI Pipeline</strong>
            <span class="check-note">GitHub Actions runs both test suites on every push to main</span>
          </div>
        </li>
      </ul>
    </div>

    <!-- ── Section 1: Data Pipeline ───────────────────────────────────────────── -->
    <div class="demo-section">
      <div class="section-header">
        <div class="section-num">1</div>
        <div>
          <h2 class="section-title">Data Pipeline</h2>
          <p class="section-desc">20 catalog products ingested into PostgreSQL with pgvector embeddings.</p>
        </div>
      </div>

      <div v-if="pipeline" class="pipeline-status">
        <div v-if="pipeline.error" class="status-error">{{ pipeline.error }}</div>
        <template v-else>
          <div class="stat-row">
            <div class="stat-box">
              <div class="stat-num">{{ pipeline.total }}</div>
              <div class="stat-label">Total products</div>
            </div>
            <div class="stat-box green">
              <div class="stat-num">{{ pipeline.embedded }}</div>
              <div class="stat-label">With embeddings</div>
            </div>
            <div class="stat-box" :class="pipeline.pending > 0 ? 'amber' : 'green'">
              <div class="stat-num">{{ pipeline.pending }}</div>
              <div class="stat-label">Pending</div>
            </div>
          </div>
          <div class="progress-bar-wrap">
            <div class="progress-bar" :style="{ width: pipeline.total ? (pipeline.embedded / pipeline.total * 100) + '%' : '0%' }"></div>
          </div>
          <p class="pipeline-note">
            {{ pipeline.embedded === pipeline.total ? '✅ All products embedded.' : '⚠️ Run ingest to embed remaining products.' }}
          </p>
          <button v-if="pipeline.pending > 0" class="btn-primary" :disabled="pipelineLoading" @click="runIngest">
            {{ pipelineLoading ? 'Ingesting…' : 'Run Ingest' }}
          </button>
        </template>
      </div>
      <div v-else class="loading-text">Loading pipeline status…</div>
    </div>

    <!-- ── Section 2: Vector Search ───────────────────────────────────────────── -->
    <div class="demo-section">
      <div class="section-header">
        <div class="section-num ai">2</div>
        <div>
          <h2 class="section-title">AI Technique 1 — Semantic Vector Search</h2>
          <p class="section-desc">Query embedded with BAAI/bge-small-en-v1.5 (ONNX), matched via pgvector cosine similarity.</p>
        </div>
      </div>
      <div class="input-row">
        <input v-model="searchQuery" class="demo-input" placeholder="e.g. heavy gate for warehouse entrance" @keyup.enter="runSearch" />
        <button class="btn-ai" :disabled="searchLoading" @click="runSearch">
          {{ searchLoading ? 'Searching…' : 'Search' }}
        </button>
      </div>
      <p v-if="searchError" class="error-text">{{ searchError }}</p>
      <div v-if="searchDone && searchResults.length" class="results-list">
        <div v-for="r in searchResults.slice(0, 5)" :key="r.product?.productId" class="result-item">
          <div class="result-top">
            <span class="result-name">{{ r.product?.name }}</span>
            <span class="result-score">{{ (r.similarity * 100).toFixed(1) }}% match</span>
          </div>
          <p class="result-desc">{{ r.product?.description }}</p>
          <span class="cat-chip">{{ r.product?.category }}</span>
        </div>
      </div>
      <div v-else-if="searchDone" class="empty-note">No results — run pipeline ingest first.</div>
    </div>

    <!-- ── Section 3: LLM Assistant ───────────────────────────────────────────── -->
    <div class="demo-section">
      <div class="section-header">
        <div class="section-num ai">3</div>
        <div>
          <h2 class="section-title">AI Technique 2 — LLM Order Assistant (Claude)</h2>
          <p class="section-desc">Structured prompting with Claude Haiku returns JSON order recommendations.</p>
        </div>
      </div>
      <div class="input-row">
        <input v-model="assistQuery" class="demo-input" placeholder="Describe your gate requirement…" @keyup.enter="runAssist" />
        <button class="btn-ai" :disabled="assistLoading" @click="runAssist">
          {{ assistLoading ? 'Asking Claude…' : 'Ask Claude' }}
        </button>
      </div>
      <p v-if="assistError" class="error-text">{{ assistError }}</p>
      <div v-if="assistDone && assistResult" class="assist-card">
        <p class="assist-rec">{{ assistResult.recommendation }}</p>
        <div class="assist-meta">
          <span v-if="assistResult.suggested_category" class="meta-tag">{{ assistResult.suggested_category }}</span>
          <span v-if="assistResult.estimated_price_range" class="price-tag">{{ assistResult.estimated_price_range }}</span>
        </div>
        <div v-if="assistResult.key_features?.length" class="feature-block">
          <p class="feature-label">Key features</p>
          <ul class="feature-list"><li v-for="f in assistResult.key_features" :key="f">{{ f }}</li></ul>
        </div>
        <div v-if="assistResult.next_steps?.length" class="feature-block">
          <p class="feature-label">Next steps</p>
          <ul class="feature-list"><li v-for="s in assistResult.next_steps" :key="s">{{ s }}</li></ul>
        </div>
        <div class="model-badge">claude-haiku-4-5-20251001</div>
      </div>
    </div>

    <!-- ── Section 4: RAG ─────────────────────────────────────────────────────── -->
    <div class="demo-section">
      <div class="section-header">
        <div class="section-num ai">4</div>
        <div>
          <h2 class="section-title">AI Technique 3 — Retrieval-Augmented Generation (RAG)</h2>
          <p class="section-desc">Top-5 relevant catalog products retrieved as context → Claude generates grounded answer.</p>
        </div>
      </div>
      <div class="input-row">
        <input v-model="ragQuery" class="demo-input" placeholder="Ask about catalog or pricing…" @keyup.enter="runRag" />
        <button class="btn-ai" :disabled="ragLoading" @click="runRag">
          {{ ragLoading ? 'Generating…' : 'Ask RAG' }}
        </button>
      </div>
      <p v-if="ragError" class="error-text">{{ ragError }}</p>
      <div v-if="ragDone && ragResult" class="rag-card">
        <p class="rag-answer">{{ ragResult.answer }}</p>
        <div class="rag-sources">
          <span class="source-label">Sources:</span>
          <span v-for="s in ragResult.sources" :key="s" class="source-chip">{{ s }}</span>
        </div>
        <div class="model-badge">claude-haiku-4-5-20251001</div>
      </div>
    </div>

    <!-- ── Section 5: Orders ──────────────────────────────────────────────────── -->
    <div class="demo-section">
      <div class="section-header">
        <div class="section-num">5</div>
        <div>
          <h2 class="section-title">Order Management</h2>
          <p class="section-desc">Browse the catalog, place orders, and track their status in My Orders.</p>
        </div>
      </div>
      <div class="action-links">
        <RouterLink to="/catalog" class="action-link primary">Go to Catalog →</RouterLink>
        <RouterLink to="/orders" class="action-link">My Orders →</RouterLink>
      </div>
    </div>

    <!-- ── Section 6: Tests & CI ──────────────────────────────────────────────── -->
    <div class="demo-section">
      <div class="section-header">
        <div class="section-num">6</div>
        <div>
          <h2 class="section-title">Tests &amp; CI</h2>
          <p class="section-desc">30 automated tests run on every push via GitHub Actions.</p>
        </div>
      </div>
      <div class="test-grid">
        <div class="test-box">
          <div class="test-count">15</div>
          <div class="test-label">Go unit tests</div>
          <div class="test-sub">JWT, login, register, updateMe</div>
          <code class="test-cmd">go test ./... -v</code>
        </div>
        <div class="test-box">
          <div class="test-count">15</div>
          <div class="test-label">Python unit tests</div>
          <div class="test-sub">search, assist, RAG, pipeline</div>
          <code class="test-cmd">pytest tests/ -v</code>
        </div>
      </div>
    </div>

  </div>
</template>

<style scoped>
.demo-page { display: flex; flex-direction: column; gap: 1.5rem; padding: 0.5rem 0 2rem; }

.demo-header { text-align: center; padding: 1.5rem 1rem; }
.demo-badge { display: inline-block; background: #059669; color: #fff; font-size: 0.75rem; font-weight: 700; letter-spacing: 0.08em; text-transform: uppercase; padding: 0.3rem 0.9rem; border-radius: 9999px; margin-bottom: 0.75rem; }
.demo-title { font-size: 1.75rem; font-weight: 800; color: #111827; margin: 0 0 0.5rem; }
.demo-subtitle { color: #6b7280; font-size: 0.9375rem; margin: 0; }

.checklist-card { background: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 16px; padding: 1.25rem 1.5rem; }
.checklist { list-style: none; padding: 0; margin: 0.75rem 0 0; display: flex; flex-direction: column; gap: 0.625rem; }
.check-item { display: flex; gap: 0.75rem; align-items: flex-start; }
.check-icon { font-size: 1rem; margin-top: 1px; flex-shrink: 0; }
.check-item strong { font-size: 0.9rem; color: #065f46; display: block; }
.check-note { font-size: 0.8125rem; color: #047857; }

.demo-section { background: #fff; border: 1px solid #e5e7eb; border-radius: 16px; padding: 1.5rem; display: flex; flex-direction: column; gap: 1rem; }
.section-header { display: flex; gap: 1rem; align-items: flex-start; }
.section-num { width: 36px; height: 36px; border-radius: 9999px; background: #e5e7eb; color: #374151; font-weight: 700; font-size: 1rem; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.section-num.ai { background: linear-gradient(135deg, #10b981, #059669); color: #fff; }
.section-title { font-size: 1rem; font-weight: 700; color: #111827; margin: 0 0 0.2rem; }
.section-desc { font-size: 0.8125rem; color: #6b7280; margin: 0; }

.pipeline-status { display: flex; flex-direction: column; gap: 0.75rem; }
.stat-row { display: flex; gap: 0.75rem; }
.stat-box { flex: 1; background: #f9fafb; border: 1px solid #e5e7eb; border-radius: 10px; padding: 0.75rem; text-align: center; }
.stat-box.green { background: #f0fdf4; border-color: #bbf7d0; }
.stat-box.amber { background: #fffbeb; border-color: #fde68a; }
.stat-num { font-size: 1.5rem; font-weight: 700; color: #111827; }
.stat-label { font-size: 0.75rem; color: #6b7280; margin-top: 0.2rem; }
.progress-bar-wrap { height: 8px; background: #e5e7eb; border-radius: 9999px; overflow: hidden; }
.progress-bar { height: 100%; background: linear-gradient(90deg, #10b981, #059669); border-radius: 9999px; transition: width 0.4s; }
.pipeline-note { font-size: 0.875rem; color: #374151; margin: 0; }
.status-error { color: #ef4444; font-size: 0.875rem; }
.loading-text { color: #9ca3af; font-size: 0.875rem; }

.input-row { display: flex; gap: 0.5rem; }
.demo-input { flex: 1; padding: 0.625rem 0.875rem; font-size: 0.9375rem; border: 1px solid #d1d5db; border-radius: 8px; outline: none; transition: border-color 0.2s; }
.demo-input:focus { border-color: #059669; box-shadow: 0 0 0 3px rgba(5,150,105,0.1); }
.btn-ai { padding: 0.625rem 1.25rem; background: linear-gradient(135deg, #10b981, #059669); color: #fff; border: none; border-radius: 8px; font-weight: 600; cursor: pointer; white-space: nowrap; }
.btn-ai:disabled { opacity: 0.7; cursor: not-allowed; }
.btn-primary { padding: 0.5rem 1.25rem; background: #059669; color: #fff; border: none; border-radius: 8px; font-weight: 600; cursor: pointer; align-self: flex-start; }
.error-text { font-size: 0.875rem; color: #ef4444; margin: 0; }
.empty-note { font-size: 0.875rem; color: #9ca3af; }

.results-list { display: flex; flex-direction: column; gap: 0.625rem; }
.result-item { background: #f9fafb; border: 1px solid #e5e7eb; border-radius: 10px; padding: 0.75rem 1rem; }
.result-top { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.25rem; }
.result-name { font-weight: 600; font-size: 0.9375rem; color: #111827; }
.result-score { font-size: 0.8125rem; font-weight: 700; color: #059669; }
.result-desc { font-size: 0.8125rem; color: #6b7280; margin: 0 0 0.5rem; }
.cat-chip { font-size: 0.75rem; background: #d1fae5; color: #065f46; padding: 0.15rem 0.5rem; border-radius: 9999px; font-weight: 600; }

.assist-card { background: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 12px; padding: 1rem; }
.assist-rec { font-size: 0.9375rem; color: #065f46; margin: 0 0 0.75rem; line-height: 1.6; }
.assist-meta { display: flex; gap: 0.5rem; flex-wrap: wrap; margin-bottom: 0.75rem; }
.meta-tag { background: #059669; color: #fff; font-size: 0.75rem; font-weight: 700; padding: 0.2rem 0.6rem; border-radius: 9999px; }
.price-tag { background: #ecfdf5; color: #059669; font-size: 0.875rem; font-weight: 700; padding: 0.2rem 0.6rem; border: 1px solid #a7f3d0; border-radius: 9999px; }
.feature-block { margin-bottom: 0.5rem; }
.feature-label { font-size: 0.75rem; font-weight: 700; color: #6b7280; text-transform: uppercase; letter-spacing: 0.05em; margin: 0 0 0.25rem; }
.feature-list { margin: 0 0 0 1rem; padding: 0; font-size: 0.875rem; color: #374151; line-height: 1.8; }
.model-badge { margin-top: 0.75rem; font-size: 0.75rem; color: #9ca3af; }

.rag-card { background: #eff6ff; border: 1px solid #bfdbfe; border-radius: 12px; padding: 1rem; }
.rag-answer { font-size: 0.9375rem; color: #1e3a5f; line-height: 1.7; margin: 0 0 0.75rem; }
.rag-sources { display: flex; flex-wrap: wrap; gap: 0.4rem; align-items: center; }
.source-label { font-size: 0.75rem; color: #6b7280; font-weight: 600; }
.source-chip { font-size: 0.75rem; background: #dbeafe; color: #1d4ed8; padding: 0.15rem 0.5rem; border-radius: 9999px; }

.action-links { display: flex; gap: 0.75rem; flex-wrap: wrap; }
.action-link { padding: 0.6rem 1.25rem; border-radius: 8px; font-size: 0.9375rem; font-weight: 600; text-decoration: none; background: #f3f4f6; color: #374151; border: 1px solid #e5e7eb; transition: background 0.15s; }
.action-link:hover { background: #e5e7eb; }
.action-link.primary { background: #059669; color: #fff; border-color: #059669; }
.action-link.primary:hover { background: #047857; }

.test-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
.test-box { background: #f9fafb; border: 1px solid #e5e7eb; border-radius: 12px; padding: 1rem; }
.test-count { font-size: 2rem; font-weight: 800; color: #059669; line-height: 1; }
.test-label { font-size: 0.9375rem; font-weight: 600; color: #111827; margin: 0.25rem 0 0.2rem; }
.test-sub { font-size: 0.8125rem; color: #6b7280; margin-bottom: 0.75rem; }
.test-cmd { display: block; background: #1f2937; color: #a7f3d0; font-size: 0.8125rem; padding: 0.4rem 0.75rem; border-radius: 6px; font-family: monospace; }
</style>
