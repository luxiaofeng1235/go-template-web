import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  base: '/admin/',
  build: {
    outDir: '../public/admin',
    emptyOutDir: true,
    rollupOptions: {
      output: {
        manualChunks: {
          'element-plus': ['element-plus'],
          'vue-vendor': ['vue', 'vue-router']
        }
      }
    }
  },
  server: {
    proxy: {
      '/admin/api': {
        target: 'http://127.0.0.1:8006',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/admin\/api/, '/admin')
      }
    }
  }
})