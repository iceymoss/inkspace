<template>
  <RouterView />
  <FloatingTerminal @submit="handleTerminalSubmit" @confirm="handleTerminalConfirm" />
</template>

<script setup>
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useTerminalStore } from '@/stores/terminal'
import { useAppearanceStore } from '@/stores/appearance'
import { createTerminalExecutor } from '@/utils/terminal/executor'
import api from '@/utils/api'
import FloatingTerminal from '@/components/terminal/FloatingTerminal.vue'

const userStore = useUserStore()
const terminalStore = useTerminalStore()
const appearanceStore = useAppearanceStore()
const router = useRouter()
const route = useRoute()
const terminalExecutor = createTerminalExecutor({ router, route, api, userStore, appearanceStore, terminalStore })

function handleTerminalSubmit(command) {
  terminalExecutor.execute(command)
}

function handleTerminalConfirm(confirmation) {
  terminalExecutor.confirm(confirmation)
}

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
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

#app {
  font-family: var(--theme-font-body, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
