<template>
  <div class="main-layout">
    <header class="header">
      <div class="container">
        <div class="header-content">
          <div class="logo">
            <router-link to="/">InkSpace</router-link>
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
                    <el-dropdown-item command="searchUser">
                      <el-icon><User /></el-icon> 搜索用户
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
        <p>{{ siteSettings.site_copyright || '© 2024 InkSpace. All rights reserved.' }}</p>
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
    
    // 更新favicon
    if (siteSettings.value.site_logo) {
      updateFavicon(siteSettings.value.site_logo)
    }
  } catch (error) {
    console.error('Failed to load site settings:', error)
  }
}

// 更新favicon
const updateFavicon = (logoUrl) => {
  // 移除旧的favicon
  const oldFavicon = document.querySelector('link[rel="icon"]')
  if (oldFavicon) {
    oldFavicon.remove()
  }
  
  // 创建新的favicon
  const link = document.createElement('link')
  link.rel = 'icon'
  link.type = 'image/png'
  link.href = logoUrl
  document.head.appendChild(link)
}

const handleCommand = (command) => {
  if (command === 'logout') {
    userStore.logout()
    router.push('/')
  } else if (command === 'dashboard') {
    router.push('/dashboard')
  } else if (command === 'profile') {
    const currentUserId = userStore.user?.id
    if (!currentUserId) {
      console.error('User ID not found')
      return
    }
    const targetPath = `/users/${currentUserId}`
    const currentPath = router.currentRoute.value.path
    
    // 如果当前路由就是目标路由，使用 replace 强制刷新
    // 否则使用 push 正常跳转
    if (currentPath === targetPath) {
      // 先跳转到首页，然后立即跳转到目标路径，强制组件重新加载
      router.replace('/').then(() => {
        setTimeout(() => {
          router.replace(targetPath)
        }, 50)
      })
    } else {
      // 正常跳转，watch 会监听到路由变化并重新加载数据
      router.push(targetPath)
    }
  } else if (command === 'edit') {
    router.push('/profile/edit')
  } else if (command === 'favorites') {
    router.push('/favorites')
  } else if (command === 'searchUser') {
    router.push('/user-search')
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
  background: var(--theme-bg-card);
  box-shadow: 0 2px 8px var(--theme-shadow);
  border-bottom: 1px solid var(--theme-border);
  position: sticky;
  top: 0;
  z-index: 1000;
  color: var(--theme-text-primary);
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
}

.logo a {
  font-size: 28px;
  font-weight: 600;
  background: linear-gradient(135deg, #667eea 0%, var(--theme-primary) 30%, #764ba2 70%, #f093fb 100%);
  background-size: 200% 200%;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-decoration: none;
  letter-spacing: 2px;
  position: relative;
  display: inline-block;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  font-style: italic;
  text-shadow: 0 0 30px rgba(102, 126, 234, 0.3);
  animation: gradientShift 3s ease infinite;
}

.logo a:hover {
  transform: translateY(-2px) scale(1.03);
  filter: brightness(1.15);
  letter-spacing: 3px;
}

@keyframes gradientShift {
  0%, 100% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
}

.nav {
  display: flex;
  gap: 30px;
}

.nav a {
  color: var(--theme-text-secondary);
  font-size: 16px;
  transition: color 0.3s;
  text-decoration: none;
}

.nav a:hover,
.nav a.router-link-active {
  color: var(--theme-primary);
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
  color: var(--theme-text-primary);
}

.main-content {
  flex: 1;
  padding: 12px 0 40px;
}

.footer {
  background: #2c3e50; /* 默认灰黑色背景 */
  border-top: 1px solid #34495e;
  color: rgba(255, 255, 255, 0.8);
  padding: 30px 0;
  text-align: center;
  margin-top: auto;
}

/* 哀悼日主题 - 底部栏使用哀悼主题 */
body.theme-mourning .footer {
  background: var(--theme-bg-card) !important;
  border-top: 1px solid var(--theme-border) !important;
  color: var(--theme-text-primary) !important;
}

body.theme-mourning .footer p {
  color: var(--theme-text-secondary) !important;
}

body.theme-mourning .footer a {
  color: var(--theme-text-primary) !important;
}

body.theme-mourning .footer a:hover {
  color: var(--theme-primary) !important;
}

/* 其他主题保持灰黑色 */
body:not(.theme-mourning) .footer {
  background: #2c3e50 !important;
  border-top: 1px solid #34495e !important;
  color: rgba(255, 255, 255, 0.8) !important;
}

body:not(.theme-mourning) .footer p {
  color: rgba(255, 255, 255, 0.8) !important;
}

body:not(.theme-mourning) .footer a {
  color: rgba(255, 255, 255, 0.9) !important;
}

body:not(.theme-mourning) .footer a:hover {
  color: #ffffff !important;
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

