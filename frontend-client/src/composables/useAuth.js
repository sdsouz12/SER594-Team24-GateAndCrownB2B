import { ref, computed } from 'vue'
import { loginRequest } from '../api/auth'

const TOKEN_KEY = 'gate_crown_client_token'
const USER_KEY = 'gate_crown_client_user'

function readSession() {
  try {
    const token = localStorage.getItem(TOKEN_KEY)
    const raw = localStorage.getItem(USER_KEY)
    const user = raw ? JSON.parse(raw) : null
    return { token: token || null, user: user && typeof user === 'object' ? user : null }
  } catch {
    return { token: null, user: null }
  }
}

function writeSession(token, user) {
  localStorage.setItem(TOKEN_KEY, token)
  localStorage.setItem(USER_KEY, JSON.stringify(user))
}

function clearSession() {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(USER_KEY)
}

const initial = readSession()
const accessToken = ref(initial.token)
const user = ref(initial.user)

const isLoggedIn = computed(() => Boolean(accessToken.value))

export function useAuth() {
  async function login(username, password) {
    const data = await loginRequest(username.trim(), password)
    const profile = {
      userId: data.userId,
      username: data.username,
      role: data.role,
      organizationId: data.organizationId ?? null,
      expiresIn: data.expiresIn,
    }
    accessToken.value = data.accessToken
    user.value = profile
    try {
      writeSession(data.accessToken, profile)
    } catch {
      /* ignore */
    }
  }

  function logout() {
    accessToken.value = null
    user.value = null
    try {
      clearSession()
    } catch {
      /* ignore */
    }
  }

  return {
    isLoggedIn,
    user: computed(() => user.value),
    login,
    logout,
  }
}
