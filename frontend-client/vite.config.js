import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const proxyTarget = env.API_PROXY_TARGET || 'http://localhost:8080'
  const apiProxy = {
    '/api': { target: proxyTarget, changeOrigin: true },
  }

  return {
    plugins: [vue()],
    server: { proxy: apiProxy },
    preview: { proxy: apiProxy },
  }
})
