<script setup>
import { useAuth } from '../composables/useAuth'

const { isLoggedIn } = useAuth()

const services = [
  { id: 6, slug: 'frame', label: 'Frame gate', icon: 'frame', color: 'violet' },
  { id: 7, slug: 'crown', label: 'Crown', icon: 'crown', color: 'teal' },
  { id: 2, slug: 'cutting', label: 'Cutting', icon: 'cut', color: 'amber' },
  { id: 3, slug: 'fittings', label: 'Fittings', icon: 'fittings', color: 'indigo' },
  { id: 4, slug: 'mdf', label: 'MDF', icon: 'mdf', color: 'emerald' },
  { id: 5, slug: 'glass', label: 'Glass', icon: 'glass', color: 'sky' },
]

/** Any service tap → login if guest; if already logged in → welcome (placeholder for future per-service routes). */
function serviceTo() {
  if (isLoggedIn.value) {
    return { name: 'Welcome' }
  }
  return { name: 'Login', query: { redirect: '/welcome' } }
}

const colorClasses = {
  amber: { bg: 'bg-amber-50', circle: 'bg-amber-100', icon: 'text-amber-600' },
  indigo: { bg: 'bg-indigo-50', circle: 'bg-indigo-100', icon: 'text-indigo-600' },
  emerald: { bg: 'bg-emerald-50', circle: 'bg-emerald-100', icon: 'text-emerald-600' },
  sky: { bg: 'bg-sky-50', circle: 'bg-sky-100', icon: 'text-sky-600' },
  violet: { bg: 'bg-violet-50', circle: 'bg-violet-100', icon: 'text-violet-600' },
  teal: { bg: 'bg-teal-50', circle: 'bg-teal-100', icon: 'text-teal-600' },
}
</script>

<template>
  <div class="py-1">
    <h1
      class="mt-4 mb-4 text-xl font-semibold text-gray-800 md:mt-6 md:mb-6 sr-only md:not-sr-only"
    >
      Choose a service
    </h1>
    <ul class="grid grid-cols-3 gap-2 sm:gap-3" role="list">
      <li v-for="service in services" :key="service.id" class="group">
        <RouterLink
          :to="serviceTo()"
          class="block overflow-hidden rounded-xl bg-gray-100 transition-colors hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:ring-offset-2"
        >
          <div
            class="flex aspect-square items-center justify-center transition-colors group-hover:opacity-95"
            :class="colorClasses[service.color].bg"
            :aria-label="service.label"
          >
            <span
              class="flex h-12 w-12 items-center justify-center rounded-xl p-1 transition-all duration-200 group-hover:scale-105 group-hover:shadow-md sm:h-16 sm:w-16"
              :class="[colorClasses[service.color].circle, colorClasses[service.color].icon]"
            >
              <svg
                v-if="service.icon === 'cut'"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 64 64"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="h-full w-full"
              >
                <circle cx="32" cy="28" r="18" fill="currentColor" fill-opacity="0.15" />
                <circle cx="32" cy="28" r="14" fill="none" />
                <path
                  d="M32 28 L46 28 M32 28 L44.5 32.5 M32 28 L44.5 23.5 M32 28 L40 40 M32 28 L40 16 M32 28 L32 42 M32 28 L32 14 M32 28 L24 40 M32 28 L24 16 M32 28 L19.5 32.5 M32 28 L19.5 23.5 M32 28 L18 28"
                />
                <circle cx="32" cy="28" r="5" fill="currentColor" />
                <rect x="8" y="48" width="48" height="8" rx="2" fill="currentColor" fill-opacity="0.25" />
                <line x1="32" y1="46" x2="32" y2="52" stroke-width="2.5" />
              </svg>
              <svg
                v-else-if="service.icon === 'fittings'"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 64 64"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="h-full w-full"
              >
                <rect
                  x="12"
                  y="26"
                  width="40"
                  height="12"
                  rx="6"
                  fill="currentColor"
                  fill-opacity="0.2"
                  stroke="currentColor"
                />
                <circle cx="20" cy="32" r="4" fill="currentColor" />
                <circle cx="44" cy="32" r="4" fill="currentColor" />
                <path
                  d="M8 20 L8 44 L16 40 L16 24 Z"
                  fill="currentColor"
                  fill-opacity="0.3"
                  stroke="currentColor"
                />
                <path
                  d="M56 20 L56 44 L48 40 L48 24 Z"
                  fill="currentColor"
                  fill-opacity="0.3"
                  stroke="currentColor"
                />
              </svg>
              <svg
                v-else-if="service.icon === 'mdf'"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 64 64"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="h-full w-full"
              >
                <path d="M12 20 L12 44 L52 44 L52 20 Z" fill="currentColor" fill-opacity="0.2" />
                <path d="M16 24 L16 40 L48 40 L48 24 Z" fill="currentColor" fill-opacity="0.15" />
                <path d="M20 28 L20 36 L44 36 L44 28 Z" fill="currentColor" fill-opacity="0.1" />
                <line x1="12" y1="24" x2="52" y2="24" />
                <line x1="12" y1="32" x2="52" y2="32" />
                <line x1="12" y1="40" x2="52" y2="40" />
                <path d="M52 20 L60 24 L60 48 L52 44 Z" fill="currentColor" fill-opacity="0.25" />
              </svg>
              <svg
                v-else-if="service.icon === 'glass'"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 64 64"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="h-full w-full"
              >
                <rect x="10" y="8" width="44" height="48" rx="4" fill="currentColor" fill-opacity="0.08" />
                <path d="M10 8 L54 8 L54 56 L10 56 Z" />
                <path d="M10 8 L54 56" stroke-width="1.5" stroke-opacity="0.6" />
                <ellipse cx="32" cy="32" rx="12" ry="14" fill="currentColor" fill-opacity="0.12" />
              </svg>
              <svg
                v-else-if="service.icon === 'frame'"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 64 64"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="h-full w-full"
              >
                <rect x="6" y="6" width="52" height="52" rx="2" fill="currentColor" fill-opacity="0.1" />
                <rect x="14" y="14" width="36" height="36" rx="1" fill="none" />
                <rect x="6" y="6" width="52" height="52" rx="2" />
                <rect x="18" y="18" width="28" height="28" fill="currentColor" fill-opacity="0.05" />
              </svg>
              <svg
                v-else-if="service.icon === 'crown'"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 64 64"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="h-full w-full"
              >
                <path
                  d="M8 44 L12 28 L24 36 L32 20 L40 36 L52 28 L56 44 L8 44 Z"
                  fill="currentColor"
                  fill-opacity="0.15"
                />
                <path d="M8 44 L12 28 L24 36 L32 20 L40 36 L52 28 L56 44" />
                <path d="M12 28 L32 20 L52 28" stroke-width="1.5" />
              </svg>
            </span>
          </div>
          <div
            class="bg-gray-900 px-1.5 py-1.5 text-center text-xs font-medium text-white sm:px-2 sm:py-2 sm:text-sm"
          >
            {{ service.label }}
          </div>
        </RouterLink>
      </li>
    </ul>
  </div>
</template>
