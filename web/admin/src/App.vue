<template>
  <RouterView />
</template>

<script setup>
import { onMounted } from 'vue'
import { useAdminStore } from '@/stores/admin'

const adminStore = useAdminStore()

onMounted(() => {
  if (adminStore.isLoggedIn) {
    adminStore.fetchProfile().catch(() => {
      adminStore.logout()
    })
  }
})
</script>
