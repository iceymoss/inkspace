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
            <router-link to="/about">关于</router-link>
            <router-link v-if="userStore.isAdmin" to="/admin">管理</router-link>
          </nav>
          <div class="header-actions">
            <template v-if="userStore.isLoggedIn">
              <el-dropdown @command="handleCommand">
                <span class="user-info">
                  <el-avatar :size="32" :src="userStore.user?.avatar" />
                  <span class="username">{{ userStore.user?.nickname || userStore.user?.username }}</span>
                </span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="profile">个人信息</el-dropdown-item>
                    <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
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
        <p>&copy; 2024 My Site. All rights reserved.</p>
        <p>Powered by Go + Gin + Vue</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { useUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'

const userStore = useUserStore()
const router = useRouter()

// Load user profile if logged in
if (userStore.isLoggedIn && !userStore.user) {
  userStore.fetchProfile()
}

const handleCommand = (command) => {
  if (command === 'logout') {
    userStore.logout()
    router.push('/')
  } else if (command === 'profile') {
    // Navigate to profile page
  }
}
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

