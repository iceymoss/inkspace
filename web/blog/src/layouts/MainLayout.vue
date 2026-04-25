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
  box-shadow: var(--shadow-sm);
  border-bottom: 1px solid var(--theme-border);
  position: sticky;
  top: 0;
  z-index: 1000;
  color: var(--theme-text-primary);
  transition: background-color var(--transition-slow), border-color var(--transition-slow);
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
}

.logo a {
  font-size: 28px;
  font-weight: 700;
  font-family: var(--font-serif);
  background: linear-gradient(135deg, var(--theme-primary) 0%, var(--theme-accent) 50%, var(--color-accent) 100%);
  background-size: 200% 200%;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-decoration: none;
  letter-spacing: 1px;
  position: relative;
  display: inline-block;
  transition: all var(--transition-base);
  animation: gradientShift 3s ease infinite;
  cursor: pointer;
}

.logo a:hover {
  letter-spacing: 2px;
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
  gap: var(--spacing-lg);
}

.nav a {
  color: var(--theme-text-secondary);
  font-size: var(--font-size-base);
  font-weight: 500;
  transition: color var(--transition-base);
  text-decoration: none;
  cursor: pointer;
  position: relative;
  padding: var(--spacing-xs) 0;
}

.nav a::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 0;
  height: 2px;
  background: var(--theme-primary);
  border-radius: var(--radius-full);
  transition: width var(--transition-base);
}

.nav a:hover::after,
.nav a.router-link-active::after {
  width: 100%;
}

.nav a:hover,
.nav a.router-link-active {
  color: var(--theme-primary);
}

.header-actions {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.user-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  cursor: pointer;
  transition: opacity var(--transition-base);
}

.user-info:hover {
  opacity: 0.8;
}

.username {
  color: var(--theme-text-primary);
  font-size: var(--font-size-sm);
  font-weight: 500;
}

.main-content {
  flex: 1;
  padding: var(--spacing-sm) 0 var(--spacing-xl);
}

.footer {
  background: var(--color-text-primary);
  border-top: 1px solid var(--theme-border);
  color: rgba(255, 255, 255, 0.8);
  padding: var(--spacing-lg) 0;
  text-align: center;
  margin-top: auto;
  transition: background-color var(--transition-slow), border-color var(--transition-slow);
}

.footer p {
  color: rgba(255, 255, 255, 0.7);
  font-size: var(--font-size-sm);
  line-height: var(--line-height-base);
  margin-bottom: var(--spacing-xs);
}

.footer a {
  color: rgba(255, 255, 255, 0.8);
  transition: color var(--transition-base);
  cursor: pointer;
}

.footer a:hover {
  color: #ffffff;
}

body.theme-mourning .footer {
  background: var(--theme-bg-card);
  border-top: 1px solid var(--theme-border);
  color: var(--theme-text-primary);
}

body.theme-mourning .footer p {
  color: var(--theme-text-secondary);
}

body.theme-mourning .footer a {
  color: var(--theme-text-primary);
}

body.theme-mourning .footer a:hover {
  color: var(--theme-primary);
}

@media (max-width: 768px) {
  .nav {
    gap: var(--spacing-md);
  }

  .nav a {
    font-size: var(--font-size-sm);
  }

  .logo a {
    font-size: 20px;
  }
}
</style>

