import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'
import { writeFileSync } from 'node:fs'
import { resolve } from 'node:path'

const outDir = process.env.VITE_OUT_DIR || '../../internal/webassets/blog/dist'

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
    port: 3001,  // 博客前端端口
    proxy: {
      '/api': {
        target: 'http://localhost:8081',  // 用户服务
        changeOrigin: true
      },
      '/uploads': {
        target: 'http://localhost:8081',  // 静态文件服务
        changeOrigin: true
      }
    }
  }
})
