<template>
  <RouterView v-slot="{ Component }">
    <transition name="page-fade" mode="out-in">
      <component :is="Component" />
    </transition>
  </RouterView>
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
#app {
  min-height: 100vh;
}
</style>
