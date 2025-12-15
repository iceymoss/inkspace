<template>
  <RouterView />
</template>

<script setup>
import { onMounted } from 'vue'
import { useAdminStore } from '@/stores/admin'

const adminStore = useAdminStore()

// 如果有token，尝试获取管理员信息
onMounted(() => {
  if (adminStore.isLoggedIn) {
    adminStore.fetchProfile().catch(() => {
      // Token过期，清除
      adminStore.logout()
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
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

:root {
  --primary-color: #409eff;
  --text-primary: #303133;
  --text-secondary: #909399;
  --border-color: #dcdfe6;
}
</style>
