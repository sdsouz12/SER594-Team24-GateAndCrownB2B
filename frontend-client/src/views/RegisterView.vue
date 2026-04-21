<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { registerRequest } from '../api/auth'

const router = useRouter()

const username = ref('')
const fullName = ref('')
const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')
const success = ref(false)

async function onSubmit(e) {
  e.preventDefault()
  error.value = ''
  loading.value = true
  try {
    await registerRequest({
      username: username.value.trim(),
      password: password.value,
      fullName: fullName.value.trim(),
      email: email.value.trim() || undefined,
    })
    success.value = true
    setTimeout(() => router.push({ name: 'Login' }), 1500)
  } catch (err) {
    error.value =
      err?.status === 409
        ? (err?.data?.error || 'Username or email already taken.')
        : err instanceof Error && err.message
          ? err.message
          : 'Registration failed. Please try again.'
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
          <RouterLink :to="{ name: 'Login' }" class="login-back-link">← Back to Sign in</RouterLink>
        </div>
        <div class="login-header">
          <div class="login-logo" aria-hidden="true">G</div>
          <h1 class="login-title">Create account</h1>
          <p class="login-subtitle">Join Gate & Crown B2B platform.</p>
        </div>

        <div v-if="success" class="success-msg">
          Account created! Redirecting to login…
        </div>

        <form v-else class="login-form" @submit="onSubmit">
          <div class="login-field">
            <label class="login-label" for="reg-fullname">Full name</label>
            <input id="reg-fullname" v-model="fullName" type="text" class="login-input"
              placeholder="John Doe" required />
          </div>
          <div class="login-field">
            <label class="login-label" for="reg-username">Username</label>
            <input id="reg-username" v-model="username" type="text" autocomplete="username"
              class="login-input" placeholder="john_doe" required minlength="3" />
          </div>
          <div class="login-field">
            <label class="login-label" for="reg-email">Email <span class="optional">(optional)</span></label>
            <input id="reg-email" v-model="email" type="email" autocomplete="email"
              class="login-input" placeholder="john@example.com" />
          </div>
          <div class="login-field">
            <label class="login-label" for="reg-password">Password</label>
            <input id="reg-password" v-model="password" type="password" autocomplete="new-password"
              class="login-input" placeholder="Min. 8 characters" required minlength="8" />
          </div>
          <p v-if="error" class="login-error" role="alert">{{ error }}</p>
          <button type="submit" class="login-submit" :disabled="loading">
            <span v-if="loading" class="login-spinner" aria-hidden="true" />
            <span v-else>Create account</span>
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
.login-bg { position: absolute; inset: 0; pointer-events: none; }
.login-bg-gradient { position: absolute; inset: 0; background: linear-gradient(135deg, #0f172a 0%, #1e293b 40%, #0f172a 100%); }
.login-bg-grid {
  position: absolute; inset: 0;
  background-image: linear-gradient(rgba(255,255,255,0.03) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.03) 1px, transparent 1px);
  background-size: 48px 48px;
}
.login-bg-glow { position: absolute; border-radius: 50%; filter: blur(100px); opacity: 0.4; }
.login-bg-glow-1 { width: 400px; height: 400px; background: #059669; top: -120px; right: -120px; }
.login-bg-glow-2 { width: 300px; height: 300px; background: #6366f1; bottom: -80px; left: -80px; }
.login-card { position: relative; width: 100%; max-width: 400px; }
.login-card-inner {
  background: rgba(255,255,255,0.06); backdrop-filter: blur(20px); -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255,255,255,0.1); border-radius: 24px; padding: 2.5rem;
  box-shadow: 0 25px 50px -12px rgba(0,0,0,0.5);
}
.login-top-link { text-align: center; margin-bottom: 1rem; }
.login-back-link { font-size: 0.875rem; font-weight: 500; color: rgba(255,255,255,0.65); text-decoration: none; transition: color 0.2s; }
.login-back-link:hover { color: #fff; }
.login-header { text-align: center; margin-bottom: 2rem; }
.login-logo {
  width: 56px; height: 56px; margin: 0 auto 1.25rem;
  background: linear-gradient(135deg, #10b981, #059669); border-radius: 16px;
  display: flex; align-items: center; justify-content: center;
  font-size: 1.75rem; font-weight: 800; color: #fff; letter-spacing: -0.02em;
}
.login-title { font-size: 1.75rem; font-weight: 700; color: #fff; letter-spacing: -0.02em; margin: 0 0 0.375rem; }
.login-subtitle { font-size: 0.9375rem; color: rgba(255,255,255,0.6); margin: 0; }
.login-form { display: flex; flex-direction: column; gap: 1.25rem; }
.login-field { display: flex; flex-direction: column; gap: 0.5rem; }
.login-label { font-size: 0.8125rem; font-weight: 500; color: rgba(255,255,255,0.8); }
.optional { font-weight: 400; color: rgba(255,255,255,0.4); }
.login-input {
  width: 100%; padding: 0.875rem 1rem; font-size: 1rem; color: #fff;
  background: rgba(255,255,255,0.06); border: 1px solid rgba(255,255,255,0.12);
  border-radius: 12px; outline: none; transition: border-color 0.2s, box-shadow 0.2s; box-sizing: border-box;
}
.login-input::placeholder { color: rgba(255,255,255,0.35); }
.login-input:focus { border-color: #10b981; box-shadow: 0 0 0 3px rgba(16,185,129,0.2); }
.login-error { font-size: 0.875rem; color: #f87171; margin: -0.25rem 0 0; line-height: 1.4; }
.success-msg { padding: 1rem; background: rgba(16,185,129,0.15); border: 1px solid rgba(16,185,129,0.3); border-radius: 12px; color: #6ee7b7; text-align: center; font-size: 0.9375rem; }
.login-submit {
  margin-top: 0.5rem; padding: 0.875rem 1.5rem; font-size: 1rem; font-weight: 600; color: #fff;
  background: linear-gradient(135deg, #10b981, #059669); border: none; border-radius: 12px;
  cursor: pointer; transition: transform 0.15s, box-shadow 0.2s;
  display: flex; align-items: center; justify-content: center; min-height: 48px;
}
.login-submit:hover:not(:disabled) { transform: translateY(-1px); box-shadow: 0 10px 25px -5px rgba(16,185,129,0.35); }
.login-submit:disabled { opacity: 0.8; cursor: not-allowed; }
.login-spinner {
  width: 22px; height: 22px; border: 2px solid rgba(255,255,255,0.3);
  border-top-color: #fff; border-radius: 50%; animation: spin 0.7s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }
</style>
