<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuth } from '../composables/useAuth'

const router = useRouter()
const route = useRoute()
const { login: doLogin } = useAuth()

const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

const redirectTo = computed(() => {
  const r = route.query.redirect
  if (typeof r === 'string' && r.startsWith('/')) return r
  return '/welcome'
})

async function onSubmit(e) {
  e.preventDefault()
  error.value = ''
  loading.value = true
  try {
    await doLogin(username.value, password.value)
    await router.push(redirectTo.value)
  } catch (err) {
    error.value =
      err?.status === 401
        ? 'Invalid username or password.'
        : err instanceof Error && err.message
          ? err.message
          : 'Could not reach the server. Is the API running?'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-page">
    <div class="login-bg" aria-hidden="true">
      <div class="login-bg-gradient" />
      <div class="login-bg-grid" />
      <div class="login-bg-glow login-bg-glow-1" />
      <div class="login-bg-glow login-bg-glow-2" />
    </div>

    <div class="login-card">
      <div class="login-card-inner">
        <div class="login-top-link">
          <RouterLink :to="{ name: 'Home' }" class="login-back-link">Home</RouterLink>
        </div>
        <div class="login-header">
          <div class="login-logo" aria-hidden="true">G</div>
          <h1 class="login-title">Sign in</h1>
          <p class="login-subtitle">Use your account to continue.</p>
        </div>

        <form class="login-form" @submit="onSubmit">
          <div class="login-field">
            <label class="login-label" for="login-username">Username</label>
            <input
              id="login-username"
              v-model="username"
              type="text"
              autocomplete="username"
              class="login-input"
              placeholder="Username"
              required
            />
          </div>
          <div class="login-field">
            <label class="login-label" for="login-password">Password</label>
            <input
              id="login-password"
              v-model="password"
              type="password"
              autocomplete="current-password"
              class="login-input"
              placeholder="Password"
              required
            />
          </div>
          <p v-if="error" class="login-error" role="alert">{{ error }}</p>
          <button type="submit" class="login-submit" :disabled="loading">
            <span v-if="loading" class="login-spinner" aria-hidden="true" />
            <span v-else>Sign in</span>
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1.5rem;
  position: relative;
  overflow: hidden;
}

.login-bg {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.login-bg-gradient {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 40%, #0f172a 100%);
}

.login-bg-grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.03) 1px, transparent 1px);
  background-size: 48px 48px;
}

.login-bg-glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(100px);
  opacity: 0.4;
}

.login-bg-glow-1 {
  width: 400px;
  height: 400px;
  background: #059669;
  top: -120px;
  right: -120px;
}

.login-bg-glow-2 {
  width: 300px;
  height: 300px;
  background: #6366f1;
  bottom: -80px;
  left: -80px;
}

.login-card {
  position: relative;
  width: 100%;
  max-width: 400px;
}

.login-card-inner {
  background: rgba(255, 255, 255, 0.06);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 24px;
  padding: 2.5rem;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
}

.login-top-link {
  text-align: center;
  margin-bottom: 1rem;
}

.login-back-link {
  font-size: 0.875rem;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.65);
  text-decoration: none;
  transition: color 0.2s ease;
}

.login-back-link:hover {
  color: #fff;
}

.login-header {
  text-align: center;
  margin-bottom: 2rem;
}

.login-logo {
  width: 56px;
  height: 56px;
  margin: 0 auto 1.25rem;
  background: linear-gradient(135deg, #10b981, #059669);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.75rem;
  font-weight: 800;
  color: #fff;
  letter-spacing: -0.02em;
}

.login-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: #fff;
  letter-spacing: -0.02em;
  margin: 0 0 0.375rem;
}

.login-subtitle {
  font-size: 0.9375rem;
  color: rgba(255, 255, 255, 0.6);
  margin: 0;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.login-field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.login-label {
  font-size: 0.8125rem;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.8);
}

.login-input {
  width: 100%;
  padding: 0.875rem 1rem;
  font-size: 1rem;
  color: #fff;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 12px;
  outline: none;
  transition:
    border-color 0.2s,
    box-shadow 0.2s;
  box-sizing: border-box;
}

.login-input::placeholder {
  color: rgba(255, 255, 255, 0.35);
}

.login-input:focus {
  border-color: #10b981;
  box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.2);
}

.login-error {
  font-size: 0.875rem;
  color: #f87171;
  margin: -0.25rem 0 0;
  line-height: 1.4;
}

.login-submit {
  margin-top: 0.5rem;
  padding: 0.875rem 1.5rem;
  font-size: 1rem;
  font-weight: 600;
  color: #fff;
  background: linear-gradient(135deg, #10b981, #059669);
  border: none;
  border-radius: 12px;
  cursor: pointer;
  transition:
    transform 0.15s,
    box-shadow 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 48px;
}

.login-submit:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 10px 25px -5px rgba(16, 185, 129, 0.35);
}

.login-submit:disabled {
  opacity: 0.8;
  cursor: not-allowed;
}

.login-spinner {
  width: 22px;
  height: 22px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: login-spin 0.7s linear infinite;
}

@keyframes login-spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
