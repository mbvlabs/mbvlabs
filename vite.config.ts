import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig(({ command }) => ({
  plugins: [vue(), tailwindcss()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./resources/js', import.meta.url)),
    },
  },
  base: command === 'build' ? './' : '/assets/dist/',
  build: {
    manifest: "vite/manifest.json",
    assetsDir: "",
    outDir: 'assets/dist',
    chunkSizeWarningLimit: 600,
    rollupOptions: {
      input: 'resources/js/app.ts',
      onLog(level, log, handler) {
        if (log.code === 'INVALID_ANNOTATION' && log.id?.includes('/@vueuse/core/')) return
        handler(level, log)
      },
    },
  },
  server: {
    port: 5173,
    strictPort: true,
  },
}))
