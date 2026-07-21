import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'
import { writeFileSync } from 'node:fs'
import { resolve } from 'node:path'

const outDir = process.env.VITE_OUT_DIR || '../../internal/webassets/admin/dist'

export default defineConfig({
  plugins: [
    vue(),
    {
      name: 'preserve-embed-directory',
      closeBundle() {
        writeFileSync(resolve(process.cwd(), outDir, '.gitkeep'), '')
      }
    }
  ],
  build: {
    outDir,
    emptyOutDir: true
  },
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
      },
      '/uploads': {
        target: 'http://localhost:8083',  // 静态文件代理
        changeOrigin: true
      }
    }
  }
})
