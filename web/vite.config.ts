import { fileURLToPath, URL } from 'node:url'
import { execSync } from 'node:child_process'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import vueDevTools from 'vite-plugin-vue-devtools'
import nightwatchPlugin from 'vite-plugin-nightwatch'
import Components from 'unplugin-vue-components/vite'
import { PrimeVueResolver } from '@primevue/auto-import-resolver'
import tailwindcss from '@tailwindcss/vite'

import pkg from './package.json'

function gitLastUpdate(): string {
  try {
    return execSync('git log -1 --format=%cI', { encoding: 'utf-8' }).trim()
  } catch {
    return new Date().toISOString()
  }
}

// Single source of truth for the app version is package.json.
const APP_VERSION = `v${pkg.version}`
const LAST_UPDATE = gitLastUpdate()

// https://vite.dev/config/
export default defineConfig({
  define: {
    __APP_VERSION__: JSON.stringify(APP_VERSION),
    __LAST_UPDATE__: JSON.stringify(LAST_UPDATE),
  },
  plugins: [
    tailwindcss(),
    vue(),
    vueJsx(),
    vueDevTools(),
    nightwatchPlugin(),
    Components({
      resolvers: [
        PrimeVueResolver()
      ]
    }),
 ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
      }
    }
  }
})
