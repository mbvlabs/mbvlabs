import { fileURLToPath, URL } from 'node:url'
import { defineConfig, loadEnv, normalizePath, type Plugin, type ServerOptions } from 'vite'
import inertia from '@inertiajs/vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

function andurelReload(): Plugin {
  return {
    name: 'andurel-reload',
    apply: 'serve',
    configureServer(server) {
      // The runner touches this file only after the rebuilt app is listening.
      const file = process.env.ANDUREL_VITE_RELOAD_FILE
      if (!file) return
      const reloadFile = normalizePath(file)
      server.watcher.add(file)
      server.watcher.on('change', (path) => {
        if (normalizePath(path) === reloadFile) server.ws.send({ type: 'full-reload' })
      })
    },
  }
}

function developmentServer(mode: string): ServerOptions {
  const env = { ...loadEnv(mode, process.cwd(), ['HOST_', 'PROTOCOL', 'INERTIA_VITE_DEV_URL']), ...process.env }
  const viteURL = new URL(env.INERTIA_VITE_DEV_URL || 'http://localhost:5173/assets/dist')
  const primaryHost = env.HOST_PRIMARY || 'localhost:8080'
  const protocol = env.PROTOCOL || 'http'
  const origins = new Set([`${protocol}://${primaryHost}`])
  for (const [key, host] of Object.entries(env)) {
    if (key.startsWith('HOST_') && host) origins.add(`${protocol}://${host}`)
  }

  // Allow the default named hosts derived from HOST_PRIMARY as well.
  const escapedHost = primaryHost.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  const subdomainOrigin = new RegExp(`^${protocol}://[a-zA-Z0-9-]+\\.${escapedHost}$`)
  return {
    port: Number(viteURL.port || (viteURL.protocol === 'https:' ? 443 : 80)),
    strictPort: true,
    origin: viteURL.origin,
    cors: { origin: [...origins, subdomainOrigin] },
    // The Go runner reloads embedded root templates after a rebuild.
    watch: { ignored: ['**/resources/*.html'] },
  }
}

export default defineConfig(({ command, mode }) => ({
  plugins: [react(), inertia({ ssr: { entry: 'resources/js/ssr.tsx' } }), tailwindcss(), andurelReload()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./resources/js', import.meta.url)),
    },
  },
  base: command === 'build' ? './' : '/assets/dist/',
  // Bundle dependencies into the production SSR build so the embedded
  // assets/dist/ssr/ssr.js runs without node_modules. The dev server keeps them
  // external so Vite can serve SSR under `andurel run`.
  ssr: command === 'build' ? { noExternal: true } : undefined,
  build: {
    manifest: "vite/manifest.json",
    assetsDir: "",
    outDir: 'assets/dist',
    rollupOptions: {
      input: 'resources/js/app.tsx',
    },
  },
  server: command === 'serve' ? developmentServer(mode) : undefined,
}))
