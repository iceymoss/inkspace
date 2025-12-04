<template>
  <div class="main-layout">
    <header class="header">
      <div class="container">
        <div class="header-content">
          <div class="logo">
            <router-link to="/">My Site</router-link>
          </div>
          <nav class="nav">
            <router-link to="/">首页</router-link>
            <router-link to="/blog">博客</router-link>
            <router-link to="/works">作品</router-link>
            <router-link to="/links">友链</router-link>
            <router-link to="/about">关于</router-link>
          </nav>
          <div class="header-actions">
            <template v-if="userStore.isLoggedIn">
              <NotificationDropdown />
              <el-dropdown @command="handleCommand">
                <span class="user-info">
                  <el-avatar :size="32" :src="userStore.user?.avatar" />
                  <span class="username">{{ userStore.user?.nickname || userStore.user?.username }}</span>
                </span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="profile">
                      <el-icon><User /></el-icon> 个人主页
                    </el-dropdown-item>
                    <el-dropdown-item command="dashboard">
                      <el-icon><Odometer /></el-icon> 用户中心
                    </el-dropdown-item>
                    <el-dropdown-item command="edit">
                      <el-icon><Edit /></el-icon> 编辑资料
                    </el-dropdown-item>
                    <el-dropdown-item command="favorites">
                      <el-icon><Collection /></el-icon> 我的收藏
                    </el-dropdown-item>
                    <el-dropdown-item command="logout" divided>
                      <el-icon><SwitchButton /></el-icon> 退出登录
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </template>
            <template v-else>
              <el-button type="primary" @click="$router.push('/login')">登录</el-button>
            </template>
          </div>
        </div>
      </div>
    </header>

    <main class="main-content">
      <RouterView />
    </main>

    <footer class="footer">
      <div class="container">
        <p>{{ siteSettings.site_copyright || '© 2024 My Site. All rights reserved.' }}</p>
        <p>{{ siteSettings.site_description || 'Powered by Go + Gin + Vue' }}</p>
        <p v-if="siteSettings.site_icp">
          <a :href="`https://beian.miit.gov.cn/`" target="_blank" rel="noopener">
            {{ siteSettings.site_icp }}
          </a>
        </p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'
import { User, Edit, Collection, SwitchButton, Odometer } from '@element-plus/icons-vue'
import NotificationDropdown from '@/components/NotificationDropdown.vue'
import api from '@/utils/api'

const userStore = useUserStore()
const router = useRouter()
const siteSettings = ref({})

// Load user profile if logged in
if (userStore.isLoggedIn && !userStore.user) {
  userStore.fetchProfile()
}

// Load public settings
const loadSiteSettings = async () => {
  try {
    const response = await api.get('/settings/public')
    siteSettings.value = response.data || {}
    
    // 更新页面标题
    if (siteSettings.value.site_name) {
      document.title = siteSettings.value.site_name
    }
  } catch (error) {
    console.error('Failed to load site settings:', error)
  }
}

const handleCommand = (command) => {
  if (command === 'logout') {
    userStore.logout()
    router.push('/')
  } else if (command === 'dashboard') {
    router.push('/dashboard')
  } else if (command === 'profile') {
    router.push(`/users/${userStore.user.id}`)
  } else if (command === 'edit') {
    router.push('/profile/edit')
  } else if (command === 'favorites') {
    router.push('/favorites')
  }
}

onMounted(() => {
  loadSiteSettings()
})
</script>

<style scoped>
.main-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 1000;
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
}

.logo a {
  font-size: 24px;
  font-weight: bold;
  color: var(--primary-color);
}

.nav {
  display: flex;
  gap: 30px;
}

.nav a {
  color: var(--text-regular);
  font-size: 16px;
  transition: color 0.3s;
}

.nav a:hover,
.nav a.router-link-active {
  color: var(--primary-color);
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 15px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.username {
  color: var(--text-primary);
}

.main-content {
  flex: 1;
  padding: 40px 0;
}

.footer {
  background: #2c3e50;
  color: white;
  padding: 30px 0;
  text-align: center;
}

.footer p {
  margin: 5px 0;
  color: rgba(255, 255, 255, 0.8);
}

@media (max-width: 768px) {
  .nav {
    gap: 15px;
  }

  .nav a {
    font-size: 14px;
  }

  .logo a {
    font-size: 20px;
  }
}
</style>

