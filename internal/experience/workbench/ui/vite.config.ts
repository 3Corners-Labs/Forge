import {defineConfig} from 'vitest/config'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  build: {outDir: '../dist', emptyOutDir: true, sourcemap: true},
  server: {proxy: {'/api': 'http://127.0.0.1:7331'}},
  test: {environment: 'jsdom', include: ['src/**/*.test.ts']}
})
