<template>
  <div
    class="main-layout"
    :class="{ 'nav-is-open': navOpen }"
  >
    <header class="header">
      <div class="container">
        <div class="header-content">
          <div class="masthead">
            <div class="logo">
              <router-link to="/">
                <template v-if="isTerminal"><span class="terminal-dollar">$</span>inkspace.log<i class="brand-caret" /></template>
                <template v-else>Ink<span>Space</span></template>
              </router-link>
            </div>
            <span
              v-if="!isTerminal"
              class="issue-mark serif"
              aria-hidden="true"
            >{{ issueMark }}</span>
          </div>
          <nav
            id="public-navigation"
            class="nav"
            :class="{ 'is-open': navOpen }"
            aria-label="主导航"
          >
            <router-link to="/">
              {{ navLabel('首页', 'home') }}
            </router-link>
            <router-link to="/blog">
              {{ navLabel('博客', 'blog') }}
            </router-link>
            <router-link to="/works">
              {{ navLabel('作品', 'projects') }}
            </router-link>
            <router-link to="/photos">
              {{ navLabel('摄影', 'photos') }}
            </router-link>
            <router-link to="/wiki">
              {{ navLabel('知识库', 'wiki') }}
            </router-link>
            <router-link to="/links">
              {{ navLabel('友链', 'links') }}
            </router-link>
            <router-link to="/about">
              {{ navLabel('关于', 'about') }}
            </router-link>
          </nav>
          <div class="header-actions">
            <el-dropdown
              v-if="!userStore.isLoggedIn"
              trigger="click"
              @command="setGuestScheme"
            >
              <button
                class="scheme-toggle"
                type="button"
                aria-label="切换明暗模式"
                title="切换明暗模式"
              >
                <Monitor v-if="appearance.savedPreference.color_scheme === 'system'" />
                <Sunny v-else-if="appearance.savedPreference.color_scheme === 'light'" />
                <Moon v-else />
              </button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="system">
                    跟随系统
                  </el-dropdown-item>
                  <el-dropdown-item command="light">
                    浅色模式
                  </el-dropdown-item>
                  <el-dropdown-item command="dark">
                    深色模式
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <template v-if="userStore.isLoggedIn">
              <NotificationDropdown />
              <el-dropdown @command="handleCommand">
                <span class="user-info">
                  <el-avatar
                    :size="32"
                    :src="userStore.user?.avatar"
                  />
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
                    <el-dropdown-item
                      command="logout"
                      divided
                    >
                      <el-icon><SwitchButton /></el-icon> 退出登录
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </template>
            <template v-else>
              <el-button
                type="primary"
                @click="$router.push('/login')"
              >
                登录
              </el-button>
            </template>
            <button
              class="nav-toggle"
              type="button"
              :aria-expanded="navOpen"
              aria-controls="public-navigation"
              :aria-label="navOpen ? '关闭导航' : '打开导航'"
              @click="navOpen = !navOpen"
            >
              <Close v-if="navOpen" />
              <Menu v-else />
            </button>
          </div>
        </div>
      </div>
    </header>

    <main class="main-content">
      <div
        v-if="!isTerminal"
        class="publication-signature serif"
        aria-hidden="true"
      >
        记录 · 观察 · <b>创造</b> · 分享
      </div>
      <RouterView />
    </main>

    <footer class="footer">
      <div class="container">
        <div class="footer-identity">
          <strong :class="{ serif: !isTerminal }">{{ isTerminal ? '$inkspace.log' : 'InkSpace' }}</strong>
          <p>{{ siteSettings.site_description || '以文字收存观察，让思想缓慢生长。' }}</p>
        </div>
        <div class="footer-meta">
          <p v-if="isTerminal" class="system-status"><i /> all systems normal</p>
          <p>{{ siteSettings.site_copyright || '© 2024 InkSpace. All rights reserved.' }}</p>
          <p v-if="siteSettings.site_icp">
            <a
              :href="`https://beian.miit.gov.cn/`"
              target="_blank"
              rel="noopener"
            >
              {{ siteSettings.site_icp }}
            </a>
          </p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { computed, ref, onMounted, watch } from 'vue'
import { useUserStore } from '@/stores/user'
import { useAppearanceStore } from '@/stores/appearance'
import { useRouter } from 'vue-router'
import { User, Edit, Collection, SwitchButton, Odometer, Monitor, Moon, Sunny, Menu, Close } from '@element-plus/icons-vue'
import NotificationDropdown from '@/components/NotificationDropdown.vue'
import api from '@/utils/api'

const userStore = useUserStore()
const appearance = useAppearanceStore()
const router = useRouter()
const siteSettings = ref({})
const navOpen = ref(false)
const issueMark = computed(() => {
  try {
    return JSON.parse(siteSettings.value.home_hero || '{}').issue || 'VOL. 01'
  } catch {
    return 'VOL. 01'
  }
})
const isTerminal = computed(() => appearance.activePreference.ui_theme === 'terminal')
const navLabel = (fallback, terminal) => isTerminal.value ? terminal : fallback

watch(() => router.currentRoute.value.fullPath, () => {
  navOpen.value = false
})

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

const setGuestScheme = (colorScheme) => {
  appearance.saveGuestPreference({
    ui_theme: appearance.savedPreference.ui_theme,
    color_scheme: colorScheme
  })
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
  background: color-mix(in srgb, var(--theme-bg-primary) 88%, transparent);
  backdrop-filter: blur(10px);
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

.masthead {
  display: flex;
  align-items: baseline;
  gap: 18px;
  flex: 0 0 auto;
}

.logo a {
  color: var(--theme-text-primary);
  font-family: Georgia, 'Songti SC', 'Noto Serif SC', SimSun, serif;
  font-size: 22px;
  font-weight: 600;
  text-decoration: none;
  letter-spacing: .1em;
  display: inline-block;
  white-space: nowrap;
}

.logo a span {
  color: var(--theme-primary);
}

.issue-mark {
  padding-left: 18px;
  border-left: 1px solid var(--theme-border);
  color: var(--theme-primary);
  font-size: 10px;
  letter-spacing: .2em;
}

.nav {
  display: flex;
  align-items: center;
  gap: clamp(16px, 2vw, 30px);
}

.nav a {
  color: var(--theme-text-secondary);
  padding: 4px 0;
  border-bottom: 1px solid transparent;
  font-size: 14px;
  transition: color .2s ease, border-color .2s ease;
  text-decoration: none;
}

.nav a:hover,
.nav a.router-link-exact-active {
  border-color: var(--theme-primary);
  color: var(--theme-text-primary);
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 15px;
}

.scheme-toggle {
  display: grid;
  width: 34px;
  height: 34px;
  padding: 7px;
  color: var(--theme-text-secondary);
  cursor: pointer;
  border: 1px solid var(--theme-border);
  border-radius: 50%;
  background: transparent;
  transition: color .25s ease, border-color .25s ease, transform .25s ease;
}

.scheme-toggle:hover {
  color: var(--theme-primary);
  border-color: var(--theme-primary);
  transform: rotate(15deg);
}

.nav-toggle {
  display: none;
  place-items: center;
  width: 36px;
  height: 36px;
  padding: 8px;
  border: 1px solid var(--theme-border);
  border-radius: 50%;
  background: transparent;
  color: var(--theme-text-primary);
  cursor: pointer;
}

.nav-toggle svg {
  width: 100%;
  height: 100%;
}

.scheme-toggle svg {
  width: 100%;
  height: 100%;
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
  position: relative;
  flex: 1;
  padding: 24px 0 56px;
}

.publication-signature {
  position: fixed;
  z-index: 2;
  top: 144px;
  right: max(18px, calc((100vw - 1200px) / 2 - 62px));
  height: 284px;
  padding-left: 18px;
  border-left: 1px solid var(--theme-border);
  color: var(--theme-text-secondary);
  font-size: 13px;
  letter-spacing: .5em;
  writing-mode: vertical-rl;
  pointer-events: none;
}

.publication-signature b {
  color: var(--theme-primary);
  font-weight: 500;
}

.footer {
  background: var(--theme-bg-primary);
  border-top: 1px solid var(--theme-border);
  color: var(--theme-text-secondary);
  padding: 48px 0 56px;
  margin-top: auto;
}

.footer .container {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 32px;
}

.footer-identity strong {
  color: var(--theme-text-primary);
  font-size: 28px;
  font-weight: 500;
  letter-spacing: .1em;
}

.footer-identity strong span {
  color: var(--theme-primary);
}

.footer p {
  margin: 8px 0 0;
  color: var(--theme-text-secondary);
  font-size: 12px;
  line-height: 1.7;
}

.footer-meta {
  text-align: right;
}

.footer a {
  color: inherit;
}

.footer a:hover {
  color: var(--theme-primary);
}

@media (max-width: 900px) {
  .nav {
    position: absolute;
    top: 64px;
    left: 0;
    right: 0;
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 0;
    padding: 16px max(24px, calc((100vw - 1060px) / 2 + 32px)) 24px;
    border-bottom: 1px solid var(--theme-border);
    background: var(--theme-bg-primary);
    opacity: 0;
    visibility: hidden;
    transform: translateY(-8px);
    transition: opacity .2s ease, transform .2s ease, visibility .2s;
  }

  .nav a {
    padding: 11px 0;
    border-bottom-color: var(--theme-border);
  }

  .nav.is-open {
    opacity: 1;
    visibility: visible;
    transform: none;
  }

  .nav-toggle {
    display: grid;
  }

  .publication-signature {
    display: none;
  }
}

@media (max-width: 560px) {
  .header-content {
    height: 58px;
  }

  .issue-mark,
  .username {
    display: none;
  }

  .header-actions {
    gap: 8px;
  }

  .nav {
    top: 58px;
    grid-template-columns: 1fr;
    padding-inline: 20px;
  }

  .main-content {
    padding-top: 12px;
  }

  .footer .container {
    display: block;
  }

  .footer-meta {
    margin-top: 24px;
    padding-top: 18px;
    border-top: 1px solid var(--theme-border);
    text-align: left;
  }
}
</style>

