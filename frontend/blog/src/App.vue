<template>
  <RouterView />
</template>

<script setup>
import { onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { initTheme } from '@/utils/theme'

const userStore = useUserStore()

// 初始化主题
initTheme()

// 如果有token，尝试获取用户信息
onMounted(() => {
  if (userStore.isLoggedIn) {
    userStore.fetchProfile().catch(() => {
      // Token过期，清除
      userStore.logout()
    })
  }
})
</script>

<style>
@import '@/assets/main.css';

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

#app {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
