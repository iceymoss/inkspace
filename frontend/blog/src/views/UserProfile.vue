<template>
  <div class="user-profile">
    <div class="container">
      <!-- 用户信息卡片 -->
      <el-card class="profile-card">
        <div class="profile-header">
          <el-avatar :size="120" :src="user?.avatar" />
          <div class="profile-info">
            <h2>{{ user?.nickname || user?.username }}</h2>
            <p class="bio">{{ user?.bio || '这个人很懒，什么都没写' }}</p>
            
            <div class="stats">
              <div class="stat-item clickable" @click="goToArticles">
                <div class="stat-value">{{ user?.article_count || 0 }}</div>
                <div class="stat-label">文章</div>
              </div>
              <div class="stat-item clickable" @click="goToFollowers">
                <div class="stat-value">{{ followStats?.follower_count || 0 }}</div>
                <div class="stat-label">粉丝</div>
              </div>
              <div class="stat-item clickable" @click="goToFollowing">
                <div class="stat-value">{{ followStats?.following_count || 0 }}</div>
                <div class="stat-label">关注</div>
              </div>
              <div class="stat-item clickable" @click="goToFavorites">
                <div class="stat-value">{{ user?.favorite_count || 0 }}</div>
                <div class="stat-label">收藏</div>
              </div>
            </div>

            <div class="profile-actions" v-if="!isCurrentUser">
              <el-button 
                v-if="!followStats?.is_following"
                type="primary" 
                @click="handleFollow"
                :loading="followLoading"
              >
                <el-icon><Plus /></el-icon> 关注
              </el-button>
              <el-button 
                v-else
                @click="handleUnfollow"
                :loading="followLoading"
              >
                <el-icon><Check /></el-icon> 已关注
              </el-button>
              
              <el-tag v-if="followStats?.is_follower" type="info" size="large">
                <el-icon><Star /></el-icon> 关注了你
              </el-tag>
            </div>
            
            <div class="profile-actions" v-else>
              <el-button @click="$router.push('/profile/edit')">
                <el-icon><Edit /></el-icon> 编辑资料
              </el-button>
            </div>
          </div>
        </div>
      </el-card>

      <!-- 标签页 -->
      <el-tabs v-model="activeTab" class="profile-tabs">
        <el-tab-pane label="文章" name="articles">
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12" :md="8" v-for="article in articles" :key="article.id">
              <el-card class="article-card" shadow="hover" @click="$router.push(`/blog/${article.id}`)">
                <img v-if="article.cover" :src="article.cover" class="article-cover" />
                <h3>{{ article.title }}</h3>
                <p class="article-summary">{{ article.summary }}</p>
                <div class="article-meta">
                  <span><el-icon><View /></el-icon> {{ article.view_count }}</span>
                  <span><el-icon><Clock /></el-icon> {{ formatDate(article.created_at) }}</span>
                </div>
              </el-card>
            </el-col>
          </el-row>
          
          <div class="pagination" v-if="articlesTotal > 0">
            <el-pagination
              v-model:current-page="articlesPage"
              :page-size="12"
              :total="articlesTotal"
              layout="prev, pager, next"
              @current-change="loadArticles"
            />
          </div>
          
          <el-empty v-if="articles.length === 0" description="还没有发布文章" />
        </el-tab-pane>

        <el-tab-pane label="收藏" name="favorites">
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12" :md="8" v-for="favorite in favorites" :key="favorite.id">
              <el-card 
                v-if="favorite.article" 
                class="article-card" 
                shadow="hover" 
                @click="$router.push(`/blog/${favorite.article.id}`)"
              >
                <img v-if="favorite.article.cover" :src="favorite.article.cover" class="article-cover" />
                <h3>{{ favorite.article.title }}</h3>
                <p class="article-summary">{{ favorite.article.summary }}</p>
                <div class="article-meta">
                  <span><el-icon><View /></el-icon> {{ favorite.article.view_count }}</span>
                  <span><el-icon><Clock /></el-icon> {{ formatDate(favorite.created_at) }}</span>
                </div>
              </el-card>
            </el-col>
          </el-row>
          
          <div class="pagination" v-if="favoritesTotal > 0">
            <el-pagination
              v-model:current-page="favoritesPage"
              :page-size="12"
              :total="favoritesTotal"
              layout="prev, pager, next"
              @current-change="loadFavorites"
            />
          </div>
          
          <el-empty v-if="favorites.length === 0" description="还没有收藏" />
        </el-tab-pane>

        <el-tab-pane :label="`关注 ${followStats?.following_count || 0}`" name="following">
          <div class="user-list">
            <el-card 
              v-for="follow in following" 
              :key="follow.id" 
              class="user-card"
              @click="$router.push(`/users/${follow.user.id}`)"
            >
              <div class="user-info">
                <el-avatar :size="60" :src="follow.user?.avatar" />
                <div class="user-details">
                  <h4>{{ follow.user?.nickname || follow.user?.username }}</h4>
                  <p>{{ follow.user?.bio || '这个人很懒，什么都没写' }}</p>
                </div>
              </div>
            </el-card>
          </div>
          
          <div class="pagination" v-if="followingTotal > 0">
            <el-pagination
              v-model:current-page="followingPage"
              :page-size="20"
              :total="followingTotal"
              layout="prev, pager, next"
              @current-change="loadFollowing"
            />
          </div>
          
          <el-empty v-if="following.length === 0" description="还没有关注任何人" />
        </el-tab-pane>

        <el-tab-pane :label="`粉丝 ${followStats?.follower_count || 0}`" name="followers">
          <div class="user-list">
            <el-card 
              v-for="follower in followers" 
              :key="follower.id" 
              class="user-card"
              @click="$router.push(`/users/${follower.user.id}`)"
            >
              <div class="user-info">
                <el-avatar :size="60" :src="follower.user?.avatar" />
                <div class="user-details">
                  <h4>{{ follower.user?.nickname || follower.user?.username }}</h4>
                  <p>{{ follower.user?.bio || '这个人很懒，什么都没写' }}</p>
                </div>
              </div>
            </el-card>
          </div>
          
          <div class="pagination" v-if="followersTotal > 0">
            <el-pagination
              v-model:current-page="followersPage"
              :page-size="20"
              :total="followersTotal"
              layout="prev, pager, next"
              @current-change="loadFollowers"
            />
          </div>
          
          <el-empty v-if="followers.length === 0" description="还没有粉丝" />
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Plus, Check, Star } from '@element-plus/icons-vue'
import api from '@/utils/api'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const user = ref(null)
const followStats = ref(null)
const activeTab = ref('articles')
const followLoading = ref(false)

// 文章
const articles = ref([])
const articlesPage = ref(1)
const articlesTotal = ref(0)

// 收藏
const favorites = ref([])
const favoritesPage = ref(1)
const favoritesTotal = ref(0)

// 关注
const following = ref([])
const followingPage = ref(1)
const followingTotal = ref(0)

// 粉丝
const followers = ref([])
const followersPage = ref(1)
const followersTotal = ref(0)

const userId = computed(() => parseInt(route.params.id))
const isCurrentUser = computed(() => userStore.user?.id === userId.value)

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD')

// 加载用户信息
const loadUser = async () => {
  try {
    const response = await api.get(`/users/${userId.value}`)
    user.value = response.data
  } catch (error) {
    ElMessage.error('用户加载失败')
  }
}

// 加载关注统计
const loadFollowStats = async () => {
  try {
    const response = await api.get(`/users/${userId.value}/follow-stats`)
    followStats.value = response.data
  } catch (error) {
    console.error('Failed to load follow stats:', error)
  }
}

// 加载用户文章
const loadArticles = async () => {
  try {
    const response = await api.get(`/users/${userId.value}/articles`, {
      params: {
        page: articlesPage.value,
        page_size: 12
      }
    })
    articles.value = response.data.list || []
    articlesTotal.value = response.data.total || 0
  } catch (error) {
    console.error('Failed to load articles:', error)
  }
}

// 加载收藏列表
const loadFavorites = async () => {
  try {
    const response = await api.get(`/users/${userId.value}/favorites`, {
      params: {
        page: favoritesPage.value,
        page_size: 12
      }
    })
    favorites.value = response.data.list || []
    favoritesTotal.value = response.data.total || 0
  } catch (error) {
    console.error('Failed to load favorites:', error)
  }
}

// 加载关注列表
const loadFollowing = async () => {
  try {
    const response = await api.get(`/users/${userId.value}/following`, {
      params: {
        page: followingPage.value,
        page_size: 20
      }
    })
    following.value = response.data.list || []
    followingTotal.value = response.data.total || 0
  } catch (error) {
    console.error('Failed to load following:', error)
  }
}

// 加载粉丝列表
const loadFollowers = async () => {
  try {
    const response = await api.get(`/users/${userId.value}/followers`, {
      params: {
        page: followersPage.value,
        page_size: 20
      }
    })
    followers.value = response.data.list || []
    followersTotal.value = response.data.total || 0
  } catch (error) {
    console.error('Failed to load followers:', error)
  }
}

// 关注用户
const handleFollow = async () => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    return
  }

  followLoading.value = true
  try {
    await api.post(`/users/${userId.value}/follow`)
    ElMessage.success('关注成功')
    await loadFollowStats()
  } catch (error) {
    ElMessage.error('关注失败')
  } finally {
    followLoading.value = false
  }
}

// 取消关注
const handleUnfollow = async () => {
  followLoading.value = true
  try {
    await api.delete(`/users/${userId.value}/follow`)
    ElMessage.success('取消关注成功')
    await loadFollowStats()
  } catch (error) {
    ElMessage.error('操作失败')
  } finally {
    followLoading.value = false
  }
}

// 跳转到文章列表
const goToArticles = () => {
  if (isCurrentUser.value) {
    // 本人 → 跳转到文章管理页面
    router.push('/dashboard/articles')
  } else {
    // 他人 → 切换到文章标签页
    activeTab.value = 'articles'
    setTimeout(() => {
      document.querySelector('.profile-tabs')?.scrollIntoView({ behavior: 'smooth' })
    }, 100)
  }
}

// 跳转到粉丝列表
const goToFollowers = () => {
  router.push(`/users/${userId.value}/follows?tab=followers`)
}

// 跳转到关注列表
const goToFollowing = () => {
  router.push(`/users/${userId.value}/follows?tab=following`)
}

// 跳转到收藏列表
const goToFavorites = () => {
  if (!isCurrentUser.value) {
    ElMessage.warning('无法查看其他用户的收藏')
    return
  }
  // 跳转到收藏管理页面
  router.push('/favorites')
}

onMounted(async () => {
  await Promise.all([
    loadUser(),
    loadFollowStats(),
    loadArticles()
  ])
})

// 监听tab切换
const onTabChange = () => {
  if (activeTab.value === 'favorites' && favorites.value.length === 0) {
    loadFavorites()
  } else if (activeTab.value === 'following' && following.value.length === 0) {
    loadFollowing()
  } else if (activeTab.value === 'followers' && followers.value.length === 0) {
    loadFollowers()
  }
}
</script>

<style scoped>
.user-profile {
  padding: 40px 0;
  background-color: var(--theme-bg-secondary);
  min-height: 100vh;
}

.profile-card {
  margin-bottom: 30px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.article-card,
.work-card {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.profile-header {
  display: flex;
  gap: 30px;
  align-items: center;
}

.profile-info {
  flex: 1;
}

.profile-info h2 {
  margin-bottom: 10px;
}

.bio {
  color: var(--text-secondary);
  margin-bottom: 20px;
}

.stats {
  display: flex;
  gap: 40px;
  margin-bottom: 20px;
}

.stat-item {
  text-align: center;
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
  font-size: 24px;
  font-weight: bold;
  color: var(--primary-color);
}

.stat-label {
  color: var(--text-secondary);
  margin-top: 5px;
}

.profile-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

.profile-tabs {
  margin-top: 20px;
}

.article-card {
  cursor: pointer;
  margin-bottom: 20px;
  height: 100%;
}

.article-cover {
  width: 100%;
  height: 200px;
  object-fit: cover;
  border-radius: 4px;
  margin-bottom: 15px;
}

.article-card h3 {
  margin-bottom: 10px;
  font-size: 1.2rem;
}

.article-summary {
  color: var(--text-secondary);
  margin-bottom: 15px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.article-meta {
  display: flex;
  gap: 15px;
  color: var(--text-secondary);
  font-size: 14px;
}

.article-meta span {
  display: flex;
  align-items: center;
  gap: 5px;
}

.user-list {
  display: grid;
  gap: 20px;
}

.user-card {
  cursor: pointer;
  transition: transform 0.3s;
}

.user-card:hover {
  transform: translateY(-3px);
}

.user-info {
  display: flex;
  gap: 20px;
  align-items: center;
}

.user-details {
  flex: 1;
}

.user-details h4 {
  margin-bottom: 5px;
}

.user-details p {
  color: var(--text-secondary);
  margin: 0;
}

.pagination {
  margin-top: 30px;
  display: flex;
  justify-content: center;
}

@media (max-width: 768px) {
  .profile-header {
    flex-direction: column;
    text-align: center;
  }

  .stats {
    justify-content: center;
  }
}
</style>

