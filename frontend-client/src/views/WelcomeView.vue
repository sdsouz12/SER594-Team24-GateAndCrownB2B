<script setup>
import { useAuth } from '../composables/useAuth'

const { user } = useAuth()
</script>

<template>
  <div class="rounded-2xl border border-emerald-100 bg-emerald-50/60 px-6 py-10 text-center shadow-sm">
    <h1 class="text-2xl font-semibold text-gray-900">You’re signed in</h1>
    <p class="mt-3 text-gray-700">
      Hello<span v-if="user?.username">, <strong>{{ user.username }}</strong></span> — you have successfully logged in.
    </p>
    <p class="mt-4 text-sm text-gray-500">
      This is a placeholder screen. More features will be added step by step.
    </p>
    <RouterLink
      to="/"
      class="mt-8 inline-block rounded-xl bg-emerald-600 px-5 py-2.5 text-sm font-medium text-white hover:bg-emerald-700"
    >
      Back to home
    </RouterLink>
  </div>
</template>
