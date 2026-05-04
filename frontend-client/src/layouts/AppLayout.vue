<script setup>
import { useAuth } from '../composables/useAuth'

const { isLoggedIn, user, logout } = useAuth()
</script>

<template>
  <div class="min-h-screen flex flex-col bg-white">
    <header
      class="sticky top-0 z-50 flex h-14 items-center justify-between border-b border-gray-200 bg-white/95 px-4 backdrop-blur-sm"
    >
      <div class="flex items-center gap-4">
        <RouterLink to="/" class="text-lg font-semibold tracking-tight text-gray-900">
          Gate Crown
        </RouterLink>
        <RouterLink to="/catalog" class="text-sm font-medium text-gray-600 hover:text-gray-900">
          Catalog
        </RouterLink>
        <RouterLink v-if="isLoggedIn" to="/orders" class="text-sm font-medium text-gray-600 hover:text-gray-900">
          My Orders
        </RouterLink>
      </div>
      <div class="flex items-center gap-3">
        <template v-if="isLoggedIn">
          <span class="hidden text-sm text-gray-600 sm:inline">{{ user?.username }}</span>
          <RouterLink
            to="/welcome"
            class="text-sm font-medium text-emerald-600 hover:text-emerald-700"
          >
            Account
          </RouterLink>
          <button
            type="button"
            class="text-sm text-gray-500 hover:text-gray-800"
            @click="logout(); $router.push({ name: 'Home' })"
          >
            Sign out
          </button>
        </template>
        <template v-else>
          <RouterLink
            :to="{ name: 'Register' }"
            class="text-sm font-medium text-gray-600 hover:text-gray-900"
          >
            Register
          </RouterLink>
          <RouterLink
            :to="{ name: 'Login', query: { redirect: '/welcome' } }"
            class="text-sm font-medium text-emerald-600 hover:text-emerald-700"
          >
            Sign in
          </RouterLink>
        </template>
      </div>
    </header>

    <main class="mx-auto w-full max-w-4xl flex-1 px-4 py-4 md:py-6">
      <RouterView />
    </main>

    <footer class="mt-auto border-t border-gray-100 py-6 text-center text-xs text-gray-400">
      Gate Crown — client app (English)
    </footer>
  </div>
</template>
