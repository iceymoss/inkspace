<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <!-- 用户信息卡片 -->
      <el-col :span="24">
        <el-card class="user-card">
          <div class="user-info">
            <el-avatar :size="80" :src="userStore.user?.avatar" />
            <div class="user-details">
              <h2>{{ userStore.user?.nickname || userStore.user?.username }}</h2>
              <p class="bio">{{ userStore.user?.bio || '这个人很懒，什么都没留下' }}</p>
              <div class="stats">
                <div class="stat-item clickable" @click="$router.push('/dashboard/articles')">
                  <span class="stat-value">{{ userStore.user?.article_count || 0 }}</span>
                  <span class="stat-label">文章</span>
                </div>
                <div class="stat-item clickable" @click="$router.push('/dashboard/works')">
                  <span class="stat-value">{{ works.length }}</span>
                  <span class="stat-label">作品</span>
                </div>
                <div class="stat-item">
                  <span class="stat-value">{{ userStore.user?.comment_count || 0 }}</span>
                  <span class="stat-label">评论</span>
                </div>
                <div class="stat-item clickable" @click="$router.push('/favorites')">
                  <span class="stat-value">{{ userStore.user?.favorite_count || 0 }}</span>
                  <span class="stat-label">收藏</span>
                </div>
                <div class="stat-item clickable" @click="goToFollowers">
                  <span class="stat-value">{{ userStore.user?.follower_count || 0 }}</span>
                  <span class="stat-label">粉丝</span>
                </div>
                <div class="stat-item clickable" @click="goToFollowing">
                  <span class="stat-value">{{ userStore.user?.following_count || 0 }}</span>
                  <span class="stat-label">关注</span>
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <!-- 快捷入口 -->
      <el-col :span="24">
        <el-card>
          <template #header>
            <span>快捷入口</span>
          </template>
          <div class="quick-links">
            <el-button @click="$router.push('/dashboard/articles')">
              <el-icon><Document /></el-icon>
              我的文章
            </el-button>
            <el-button @click="$router.push('/dashboard/works')">
              <el-icon><Picture /></el-icon>
              我的作品
            </el-button>
            <el-button @click="$router.push('/dashboard/comments')">
              <el-icon><ChatDotRound /></el-icon>
              我的评论
            </el-button>
            <el-button @click="$router.push('/favorites')">
              <el-icon><Star /></el-icon>
              我的收藏
            </el-button>
            <el-button @click="$router.push('/dashboard/notifications')">
              <el-icon><Bell /></el-icon>
              消息通知
              <el-badge :value="unreadCount" :hidden="unreadCount === 0" class="badge-inline" />
            </el-button>
            <el-button @click="$router.push('/profile/edit')">
              <el-icon><Setting /></el-icon>
              个人设置
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <!-- 最近作品 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>最近作品</span>
              <el-link type="primary" @click="$router.push('/dashboard/works')">查看全部</el-link>
            </div>
          </template>
          <div class="recent-items">
            <div 
              v-for="work in works" 
              :key="work.id" 
              class="recent-item"
              @click="$router.push(`/works/${work.id}`)"
            >
              <el-image :src="work.cover" fit="cover" style="width: 60px; height: 60px; border-radius: 4px;" />
              <div class="item-info">
                <div class="item-title">{{ work.title }}</div>
                <div class="item-meta">
                  <span><el-icon><View /></el-icon> {{ work.view_count }}</span>
                  <span><el-icon><ChatDotRound /></el-icon> {{ work.comment_count }}</span>
                </div>
              </div>
            </div>
            <el-empty v-if="works.length === 0" description="暂无作品" :image-size="60" />
          </div>
        </el-card>
      </el-col>

      <!-- 最近文章 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>最近文章</span>
              <el-link type="primary" @click="$router.push('/dashboard/articles')">查看全部</el-link>
            </div>
          </template>
          <div class="recent-items">
            <div 
              v-for="article in articles" 
              :key="article.id" 
              class="recent-item"
              @click="$router.push(`/blog/${article.id}`)"
            >
              <el-image 
                v-if="article.cover"
                :src="article.cover" 
                fit="cover" 
                style="width: 60px; height: 60px; border-radius: 4px;" 
              />
              <div class="item-info">
                <div class="item-title">{{ article.title }}</div>
                <div class="item-meta">
                  <span><el-icon><View /></el-icon> {{ article.view_count }}</span>
                  <span><el-icon><ChatDotRound /></el-icon> {{ article.comment_count }}</span>
                </div>
              </div>
            </div>
            <el-empty v-if="articles.length === 0" description="暂无文章" :image-size="60" />
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Document, Picture, ChatDotRound, Star, Bell, Setting, View } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import api from '@/utils/api'

const router = useRouter()
const userStore = useUserStore()

// 跳转到粉丝列表（跳转到用户主页的粉丝标签页）
const goToFollowers = () => {
  router.push(`/users/${userStore.user?.id}?tab=followers`)
}

// 跳转到关注列表（跳转到用户主页的关注标签页）
const goToFollowing = () => {
  router.push(`/users/${userStore.user?.id}?tab=following`)
}

const works = ref([])
const articles = ref([])
const unreadCount = ref(0)

const loadRecentWorks = async () => {
  try {
    const response = await api.get('/works/my', {
      params: {
        page: 1,
        page_size: 5
      }
    })
    works.value = response.data.list || []
  } catch (error) {
    console.error('Failed to load works:', error)
  }
}

const loadRecentArticles = async () => {
  try {
    const response = await api.get('/articles', {
      params: {
        author_id: userStore.user?.id,
        page: 1,
        page_size: 5
      }
    })
    articles.value = response.data.list || []
  } catch (error) {
    console.error('Failed to load articles:', error)
  }
}

const loadUnreadCount = async () => {
  try {
    const response = await api.get('/notifications/unread-count')
    unreadCount.value = response.data.count || 0
  } catch (error) {
    console.error('Failed to load unread count:', error)
  }
}

onMounted(() => {
  loadRecentWorks()
  loadRecentArticles()
  loadUnreadCount()
})
</script>

<style scoped>
.dashboard {
  max-width: 1200px;
}

.user-card {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.user-info {
  display: flex;
  gap: 30px;
  align-items: center;
}

.user-details {
  flex: 1;
}

.user-details h2 {
  margin: 0 0 10px 0;
  font-size: 1.8rem;
}

.bio {
  color: var(--text-secondary);
  margin-bottom: 20px;
}

.stats {
  display: flex;
  gap: 40px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-item.clickable {
  cursor: pointer;
  transition: all 0.3s;
}

.stat-item.clickable:hover {
  transform: translateY(-2px);
}

.stat-item.clickable:hover .stat-value {
  color: var(--theme-primary);
}

.stat-value {
  font-size: 1.8rem;
  font-weight: 600;
  color: var(--theme-primary);
}

.stat-label {
  font-size: 0.9rem;
  color: var(--text-secondary);
  margin-top: 5px;
}

.quick-links {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
}

.quick-links .el-button {
  flex: 1;
  min-width: 140px;
  height: 50px;
  font-size: 15px;
}

.badge-inline {
  margin-left: 5px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.recent-items {
  display: flex;
  flex-direction: column;
  gap: 15px;
  min-height: 200px;
}

.recent-item {
  display: flex;
  gap: 15px;
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.recent-item:hover {
  background-color: var(--theme-bg-secondary);
}

.item-info {
  flex: 1;
  min-width: 0;
}

.item-title {
  font-weight: 500;
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.item-meta {
  display: flex;
  gap: 15px;
  font-size: 0.9rem;
  color: var(--text-secondary);
}

.item-meta span {
  display: flex;
  align-items: center;
  gap: 5px;
}

@media (max-width: 768px) {
  .user-info {
    flex-direction: column;
    text-align: center;
  }
  
  .stats {
    justify-content: center;
  }
  
  .quick-links {
    flex-direction: column;
  }
  
  .quick-links .el-button {
    width: 100%;
  }
}
</style>
