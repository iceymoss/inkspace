import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    port: 3002,  // 管理后台前端端口
    proxy: {
      '/api': {
        target: 'http://localhost:8083',  // 管理后台服务
        changeOrigin: true
      }
    }
  }
})

