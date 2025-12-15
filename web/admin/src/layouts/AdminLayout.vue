<template>
  <div class="admin-layout">
    <el-container>
      <el-aside width="200px" class="sidebar">
        <div class="logo">
          <h3>后台管理</h3>
        </div>
        <el-menu
          :default-active="activeMenu"
          router
          background-color="#304156"
          text-color="#bfcbd9"
          active-text-color="#409eff"
        >
          <el-menu-item index="/">
            <el-icon><House /></el-icon>
            <span>控制台</span>
          </el-menu-item>
          <el-menu-item index="/articles">
            <el-icon><Document /></el-icon>
            <span>文章管理</span>
          </el-menu-item>
          <el-menu-item index="/works">
            <el-icon><Picture /></el-icon>
            <span>作品管理</span>
          </el-menu-item>
          <el-menu-item index="/categories">
            <el-icon><Folder /></el-icon>
            <span>分类管理</span>
          </el-menu-item>
          <el-menu-item index="/tags">
            <el-icon><CollectionTag /></el-icon>
            <span>标签管理</span>
          </el-menu-item>
          <el-menu-item index="/comments">
            <el-icon><ChatDotRound /></el-icon>
            <span>评论管理</span>
          </el-menu-item>
          <el-menu-item index="/links">
            <el-icon><Link /></el-icon>
            <span>友链管理</span>
          </el-menu-item>
          <el-menu-item index="/settings">
            <el-icon><Setting /></el-icon>
            <span>系统配置</span>
          </el-menu-item>
          <el-menu-item index="/users">
            <el-icon><User /></el-icon>
            <span>用户管理</span>
          </el-menu-item>
          <el-menu-item index="/ads">
            <el-icon><Promotion /></el-icon>
            <span>广告管理</span>
          </el-menu-item>
        </el-menu>
      </el-aside>

      <el-container>
        <el-header class="header">
          <div class="header-content">
            <el-breadcrumb separator="/">
              <el-breadcrumb-item :to="{ path: '/' }">管理后台</el-breadcrumb-item>
              <el-breadcrumb-item>{{ breadcrumbTitle }}</el-breadcrumb-item>
            </el-breadcrumb>
            <div class="header-actions">
              <el-dropdown @command="handleCommand">
                <span class="user-info">
                  <el-avatar :size="32" :src="adminStore.admin?.avatar" />
                  <span>{{ adminStore.admin?.nickname || adminStore.admin?.username }}</span>
                </span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="logout">退出登录</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>
        </el-header>

        <el-main class="main-content">
          <RouterView />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAdminStore } from '@/stores/admin'
import adminApi from '@/utils/adminApi'
import {
  House,
  Document,
  Picture,
  Folder,
  CollectionTag,
  ChatDotRound,
  Link,
  Setting,
  User,
  Promotion
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const adminStore = useAdminStore()

const activeMenu = computed(() => route.path)

const breadcrumbTitle = computed(() => {
  const titles = {
    '/': '控制台',
    '/articles': '文章管理',
    '/works': '作品管理',
    '/categories': '分类管理',
    '/tags': '标签管理',
    '/comments': '评论管理',
    '/links': '友链管理',
    '/settings': '系统配置',
    '/users': '用户管理',
    '/ads': '广告管理'
  }
  return titles[route.path] || '管理'
})

const handleCommand = (command) => {
  if (command === 'logout') {
    adminStore.logout()
    router.push('/login')
  }
}

// 加载网站设置并更新favicon
const loadSiteSettings = async () => {
  try {
    const response = await adminApi.get('/admin/settings')
    const settings = response.data || []
    const siteLogoSetting = settings.find(s => s.key === 'site_logo')
    
    if (siteLogoSetting && siteLogoSetting.value) {
      updateFavicon(siteLogoSetting.value)
    }
    
    // 更新页面标题
    const siteNameSetting = settings.find(s => s.key === 'site_name')
    if (siteNameSetting && siteNameSetting.value) {
      document.title = `${siteNameSetting.value} - 管理后台`
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

onMounted(() => {
  loadSiteSettings()
})
</script>

<style scoped>
.admin-layout {
  height: 100vh;
}

.el-container {
  height: 100%;
}

.sidebar {
  background-color: #304156;
  height: 100vh;
}

.logo {
  padding: 20px;
  text-align: center;
  color: white;
  border-bottom: 1px solid #1f2d3d;
}

.logo h3 {
  margin: 0;
  color: white;
}

.el-menu {
  border-right: none;
}

.header {
  background: white;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 20px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.main-content {
  background: #f0f2f5;
  padding: 20px;
}
</style>

