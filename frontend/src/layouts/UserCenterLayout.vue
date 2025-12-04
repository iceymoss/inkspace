<template>
  <div class="user-center-layout">
    <el-container>
      <el-aside width="200px" class="sidebar">
        <div class="logo">
          <router-link to="/">My Site</router-link>
        </div>
        <el-menu
          :default-active="activeMenu"
          :router="true"
          class="menu"
        >
          <el-menu-item index="/dashboard">
            <el-icon><Odometer /></el-icon>
            <span>我的主页</span>
          </el-menu-item>
          <el-menu-item index="/dashboard/articles">
            <el-icon><Document /></el-icon>
            <span>我的文章</span>
          </el-menu-item>
          <el-menu-item index="/dashboard/works">
            <el-icon><Picture /></el-icon>
            <span>我的作品</span>
          </el-menu-item>
          <el-menu-item index="/favorites">
            <el-icon><Collection /></el-icon>
            <span>我的收藏</span>
          </el-menu-item>
          <el-menu-item index="/notifications">
            <el-icon><Bell /></el-icon>
            <span>我的通知</span>
          </el-menu-item>
          <el-menu-item index="/profile/edit">
            <el-icon><User /></el-icon>
            <span>个人设置</span>
          </el-menu-item>
        </el-menu>
      </el-aside>

      <el-container>
        <el-header class="header">
          <div class="header-content">
            <el-breadcrumb separator="/">
              <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
              <el-breadcrumb-item>{{ breadcrumbTitle }}</el-breadcrumb-item>
            </el-breadcrumb>
            <div class="header-actions">
              <el-button text @click="$router.push('/')">返回网站</el-button>
              <el-button v-if="userStore.isAdmin" type="primary" @click="$router.push('/admin')">
                管理后台
              </el-button>
              <el-dropdown @command="handleCommand">
                <span class="user-info">
                  <el-avatar :size="32" :src="userStore.user?.avatar" />
                  <span>{{ userStore.user?.nickname || userStore.user?.username }}</span>
                </span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="logout">
                      <el-icon><SwitchButton /></el-icon> 退出登录
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>
        </el-header>

        <el-main class="main">
          <RouterView />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import {
  Odometer,
  Document,
  Picture,
  Collection,
  Bell,
  User,
  SwitchButton
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)

const breadcrumbTitle = computed(() => {
  const titles = {
    '/dashboard': '我的主页',
    '/dashboard/articles': '我的文章',
    '/dashboard/articles/create': '写文章',
    '/dashboard/works': '我的作品',
    '/favorites': '我的收藏',
    '/notifications': '我的通知',
    '/profile/edit': '个人设置'
  }
  // 检查文章编辑路径
  if (route.path.includes('/dashboard/articles/') && route.path.includes('/edit')) {
    return '编辑文章'
  }
  return titles[route.path] || '用户中心'
})

const handleCommand = (command) => {
  if (command === 'logout') {
    userStore.logout()
    router.push('/')
  }
}
</script>

<style scoped>
.user-center-layout {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.sidebar {
  background-color: #fff;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: bold;
  border-bottom: 1px solid #e4e7ed;
}

.logo a {
  color: #409eff;
  text-decoration: none;
}

.menu {
  border-right: none;
}

.header {
  background-color: #fff;
  border-bottom: 1px solid #e4e7ed;
  padding: 0 20px;
}

.header-content {
  height: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
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

.main {
  padding: 20px;
}
</style>

