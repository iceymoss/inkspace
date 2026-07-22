<template>
  <div class="user-profile">
    <div class="container">
      <!-- 用户信息卡片 -->
      <el-card class="profile-card">
        <div class="profile-kicker">
          CONTRIBUTOR · PROFILE
        </div>
        <div class="profile-header">
          <el-avatar
            :size="120"
            :src="user?.avatar"
          />
          <div class="profile-info">
            <h2>{{ user?.nickname || user?.username }}</h2>
            <p class="bio">
              {{ user?.bio || '这个人很懒，什么都没写' }}
            </p>
            
            <div class="stats">
              <div
                class="stat-item clickable"
                @click="goToArticles"
              >
                <div class="stat-value">
                  {{ user?.article_count || 0 }}
                </div>
                <div class="stat-label">
                  文章
                </div>
              </div>
              <div
                class="stat-item clickable"
                @click="goToWorks"
              >
                <div class="stat-value">
                  {{ worksTotal || 0 }}
                </div>
                <div class="stat-label">
                  作品
                </div>
              </div>
              <!-- 只有自己的主页才能点击查看粉丝列表 -->
              <div 
                class="stat-item" 
                :class="{ clickable: isCurrentUser }" 
                @click="isCurrentUser && goToFollowers()"
              >
                <div class="stat-value">
                  {{ followStats?.follower_count || 0 }}
                </div>
                <div class="stat-label">
                  粉丝
                </div>
              </div>
              <!-- 只有自己的主页才能点击查看关注列表 -->
              <div 
                class="stat-item" 
                :class="{ clickable: isCurrentUser }" 
                @click="isCurrentUser && goToFollowing()"
              >
                <div class="stat-value">
                  {{ followStats?.following_count || 0 }}
                </div>
                <div class="stat-label">
                  关注
                </div>
              </div>
              <!-- 只有自己的主页才能点击查看收藏列表 -->
              <div 
                class="stat-item" 
                :class="{ clickable: isCurrentUser }" 
                @click="isCurrentUser && goToFavorites()"
              >
                <div class="stat-value">
                  {{ user?.favorite_count || 0 }}
                </div>
                <div class="stat-label">
                  收藏
                </div>
              </div>
            </div>

            <!-- 未登录用户：不显示任何操作按钮 -->
            <!-- 登录用户查看自己的主页：显示编辑资料按钮 -->
            <div
              v-if="isCurrentUser && userStore.isLoggedIn"
              class="profile-actions"
            >
              <el-button @click="$router.push('/profile/edit')">
                <el-icon><Edit /></el-icon> 编辑资料
              </el-button>
            </div>
            
            <!-- 登录用户查看他人的主页：显示关注/已关注按钮 -->
            <div
              v-else-if="!isCurrentUser && userStore.isLoggedIn"
              class="profile-actions"
            >
              <el-button 
                v-if="!followStats?.is_following"
                type="primary" 
                :loading="followLoading"
                @click="handleFollow"
              >
                <el-icon><Plus /></el-icon> 关注
              </el-button>
              <el-button 
                v-else
                :loading="followLoading"
                @click="handleUnfollow"
              >
                <el-icon><Check /></el-icon> 已关注
              </el-button>
              
              <el-tag
                v-if="followStats?.is_follower"
                type="info"
                size="large"
              >
                <el-icon><Star /></el-icon> 关注了你
              </el-tag>
            </div>
          </div>
        </div>
      </el-card>

      <!-- 标签页 -->
      <el-tabs
        v-model="activeTab"
        class="profile-tabs"
        @tab-change="onTabChange"
      >
        <el-tab-pane
          label="文章"
          name="articles"
        >
          <!-- 排序选择器 -->
          <div class="article-sort-header">
            <el-radio-group
              v-model="articleSortBy"
              size="default"
              @change="handleArticleSortChange"
            >
              <el-radio-button label="latest">
                最新排序
              </el-radio-button>
              <el-radio-button label="hot">
                最热排序
              </el-radio-button>
            </el-radio-group>
          </div>

          <!-- 用户文章列表（独立样式，不与博客列表页共用） -->
          <div class="user-article-list">
            <el-card
              v-for="article in articles"
              :key="article.id"
              class="user-article-item"
              shadow="hover"
              @click="$router.push(`/blog/${article.id}`)"
            >
              <div class="user-article-content">
                <img
                  v-if="article.cover"
                  :src="article.cover"
                  class="user-article-cover"
                >
                <div class="user-article-info">
                  <div class="user-article-header">
                    <h2>{{ article.title }}</h2>
                    <el-tag
                      v-if="article.is_top"
                      type="danger"
                      size="small"
                    >
                      置顶
                    </el-tag>
                  </div>
                  <p class="user-article-summary">
                    {{ article.summary }}
                  </p>
                  <div class="user-article-meta">
                    <el-tag
                      v-if="article.category"
                      size="small"
                    >
                      {{ article.category.name }}
                    </el-tag>
                    <span class="user-article-author">
                      <el-avatar
                        :size="20"
                        :src="article.author?.avatar"
                      />
                      <span>{{ article.author?.nickname || article.author?.username }}</span>
                    </span>
                    <span><el-icon><View /></el-icon> {{ article.view_count }}</span>
                    <span><el-icon><Clock /></el-icon> {{ formatDate(article.created_at) }}</span>
                    <!-- 书签样式的标签 -->
                    <div
                      v-if="article.tags && article.tags.length > 0"
                      class="user-article-bookmarks"
                    >
                      <span 
                        v-for="tag in article.tags.slice(0, 3)" 
                        :key="tag.id" 
                        class="user-bookmark-tag"
                        :style="{ backgroundColor: getTagColor(tag.name) }"
                        @click.stop="handleTagClick(tag.id)"
                      >
                        {{ tag.name }}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </el-card>
          </div>
          
          <div
            v-if="articlesTotal > 0"
            class="pagination"
          >
            <el-pagination
              v-model:current-page="articlesPage"
              :page-size="pageSize"
              :total="articlesTotal"
              layout="prev, pager, next"
              @current-change="loadArticles"
            />
          </div>
          
          <el-empty
            v-if="articles.length === 0"
            description="还没有发布文章"
          />
        </el-tab-pane>

        <!-- 作品标签页 -->
        <el-tab-pane
          label="作品"
          name="works"
        >
          <!-- 筛选条件 -->
          <div class="user-works-filters">
            <el-segmented
              v-model="workFilterType"
              :options="workTypeOptions"
              @change="handleWorkFilterChange"
            />
            <el-radio-group
              v-model="workSortBy"
              size="default"
              style="margin-left: 15px;"
              @change="handleWorkSortChange"
            >
              <el-radio-button label="latest">
                最新
              </el-radio-button>
              <el-radio-button label="hot">
                最热
              </el-radio-button>
            </el-radio-group>
          </div>
          
          <!-- 瀑布流布局 -->
          <div class="user-works-masonry-grid">
            <div 
              v-for="work in works" 
              :key="work.id" 
              class="user-works-masonry-item"
              @click="router.push(`/works/${work.id}`)"
            >
              <div class="user-works-image-container">
                <el-image 
                  :src="work.cover" 
                  :alt="work.title"
                  fit="cover"
                  class="user-works-image"
                  lazy
                />
                <div class="user-works-overlay">
                  <div class="user-works-overlay-content">
                    <div class="user-works-type-badge">
                      <el-tag
                        :type="work.type === 'photography' ? 'warning' : 'primary'"
                        size="small"
                      >
                        {{ work.type === 'photography' ? '📷' : '💻' }}
                      </el-tag>
                    </div>
                  </div>
                </div>
              </div>
              
              <div class="user-works-info">
                <h3 class="user-works-title">
                  {{ work.title }}
                </h3>
                <div class="user-works-meta">
                  <div class="user-works-author">
                    <el-avatar
                      :size="20"
                      :src="work.author?.avatar"
                    />
                    <span>{{ work.author?.nickname || work.author?.username }}</span>
                  </div>
                  <div class="user-works-stats">
                    <span><el-icon><View /></el-icon> {{ work.view_count }}</span>
                    <span v-if="work.like_count > 0">
                      <el-icon><Star /></el-icon> {{ work.like_count }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div
            v-if="worksTotal > 0"
            class="pagination"
          >
            <el-pagination
              v-model:current-page="worksPage"
              v-model:page-size="worksPageSize"
              :total="worksTotal"
              layout="prev, pager, next"
              @current-change="loadWorks"
            />
          </div>

          <el-empty
            v-if="works.length === 0"
            description="暂无作品"
          />
        </el-tab-pane>

        <!-- 只有自己的主页才显示收藏标签页 -->
        <el-tab-pane
          v-if="isCurrentUser"
          label="收藏"
          name="favorites"
        >
          <!-- 用户收藏列表（独立样式，不与博客列表页共用） -->
          <div class="user-favorite-list">
            <template
              v-for="favorite in favorites"
              :key="favorite.id"
            >
              <!-- 文章收藏 -->
              <el-card
                v-if="favorite && favorite.article"
                class="user-favorite-item"
                shadow="hover"
                @click="$router.push(`/blog/${favorite.article.id}`)"
              >
                <div class="user-favorite-content">
                  <img
                    v-if="favorite.article.cover"
                    :src="favorite.article.cover"
                    class="user-favorite-cover"
                  >
                  <div class="user-favorite-info">
                    <div class="user-favorite-header">
                      <h2>{{ favorite.article.title }}</h2>
                      <el-tag
                        type="info"
                        size="small"
                      >
                        文章
                      </el-tag>
                      <el-tag
                        v-if="favorite.article.is_top"
                        type="danger"
                        size="small"
                      >
                        置顶
                      </el-tag>
                    </div>
                    <p class="user-favorite-summary">
                      {{ favorite.article.summary }}
                    </p>
                    <div class="user-favorite-meta">
                      <el-tag
                        v-if="favorite.article.category"
                        size="small"
                      >
                        {{ favorite.article.category.name }}
                      </el-tag>
                      <span class="user-favorite-author">
                        <el-avatar
                          :size="20"
                          :src="favorite.article.author?.avatar"
                        />
                        <span>{{ favorite.article.author?.nickname || favorite.article.author?.username }}</span>
                      </span>
                      <span><el-icon><View /></el-icon> {{ favorite.article.view_count }}</span>
                      <span><el-icon><Clock /></el-icon> {{ formatDate(favorite.created_at) }}</span>
                      <!-- 书签样式的标签 -->
                      <div
                        v-if="favorite.article.tags && favorite.article.tags.length > 0"
                        class="user-favorite-bookmarks"
                      >
                        <span 
                          v-for="tag in favorite.article.tags.slice(0, 3)" 
                          :key="tag.id" 
                          class="user-favorite-bookmark-tag"
                          :style="{ backgroundColor: getTagColor(tag.name) }"
                          @click.stop="handleTagClick(tag.id)"
                        >
                          {{ tag.name }}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
              </el-card>
              
              <!-- 作品收藏 -->
              <el-card
                v-else-if="favorite && favorite.work"
                class="user-favorite-item user-favorite-work-item"
                shadow="hover"
                @click="$router.push(`/works/${favorite.work.id}`)"
              >
                <div class="user-favorite-content">
                  <div class="user-favorite-work-cover-container">
                    <el-image 
                      v-if="favorite.work.cover" 
                      :src="favorite.work.cover" 
                      class="user-favorite-cover"
                      fit="cover"
                    />
                    <el-tag 
                      :type="favorite.work.type === 'photography' ? 'warning' : 'primary'" 
                      size="small"
                      class="user-favorite-work-type-badge"
                    >
                      {{ favorite.work.type === 'photography' ? '📷 摄影' : '💻 项目' }}
                    </el-tag>
                  </div>
                  <div class="user-favorite-info">
                    <div class="user-favorite-header">
                      <h2>{{ favorite.work.title }}</h2>
                      <el-tag
                        type="success"
                        size="small"
                      >
                        作品
                      </el-tag>
                    </div>
                    <!-- 摄影作品显示描述，开源作品不显示 -->
                    <p
                      v-if="favorite.work.type === 'photography' && favorite.work.description"
                      class="user-favorite-summary"
                    >
                      {{ favorite.work.description }}
                    </p>
                    <div class="user-favorite-meta">
                      <span class="user-favorite-author">
                        <el-avatar
                          :size="20"
                          :src="favorite.work.author?.avatar"
                        />
                        <span>{{ favorite.work.author?.nickname || favorite.work.author?.username }}</span>
                      </span>
                      <span><el-icon><View /></el-icon> {{ favorite.work.view_count || 0 }}</span>
                      <span v-if="favorite.work.like_count > 0">
                        <el-icon><Star /></el-icon> {{ favorite.work.like_count }}
                      </span>
                      <span><el-icon><Clock /></el-icon> {{ formatDate(favorite.created_at) }}</span>
                    </div>
                  </div>
                </div>
              </el-card>
            </template>
          </div>
          
          <div
            v-if="favoritesTotal > 0"
            class="pagination"
          >
            <el-pagination
              v-model:current-page="favoritesPage"
              :page-size="12"
              :total="favoritesTotal"
              layout="prev, pager, next"
              @current-change="loadFavorites"
            />
          </div>
          
          <el-empty
            v-if="favorites.length === 0"
            description="还没有收藏"
          />
        </el-tab-pane>

        <!-- 只有自己的主页才显示关注标签页 -->
        <el-tab-pane
          v-if="isCurrentUser"
          :label="`关注 ${followStats?.following_count || 0}`"
          name="following"
        >
          <div class="user-list">
            <el-card 
              v-for="follow in following" 
              :key="follow.id" 
              class="user-card"
            >
              <div class="user-card-content">
                <div
                  class="user-info"
                  @click="goToUserProfile(follow.user?.id)"
                >
                  <el-avatar
                    :size="60"
                    :src="follow.user?.avatar"
                  />
                  <div class="user-details">
                    <h4>{{ follow.user?.nickname || follow.user?.username }}</h4>
                    <p>{{ follow.user?.bio || '这个人很懒，什么都没写' }}</p>
                    <div class="user-stats">
                      <span>文章 {{ follow.user?.article_count || 0 }}</span>
                      <span>粉丝 {{ follow.user?.follower_count || 0 }}</span>
                    </div>
                  </div>
                </div>
                <div
                  v-if="userStore.isLoggedIn && follow.user?.id !== userStore.user?.id"
                  class="user-actions"
                >
                  <el-button
                    :type="follow.is_following ? 'default' : 'primary'"
                    size="small"
                    :loading="followLoading"
                    @click.stop="handleFollowToggle(follow)"
                  >
                    {{ follow.is_following ? '已关注' : '关注' }}
                  </el-button>
                </div>
              </div>
            </el-card>
          </div>
          
          <div
            v-if="followingTotal > 0"
            class="pagination"
          >
            <el-pagination
              v-model:current-page="followingPage"
              :page-size="20"
              :total="followingTotal"
              layout="prev, pager, next"
              @current-change="loadFollowing"
            />
          </div>
          
          <el-empty
            v-if="following.length === 0"
            description="还没有关注任何人"
          />
        </el-tab-pane>

        <!-- 只有自己的主页才显示粉丝标签页 -->
        <el-tab-pane
          v-if="isCurrentUser"
          :label="`粉丝 ${followStats?.follower_count || 0}`"
          name="followers"
        >
          <div class="user-list">
            <el-card 
              v-for="follower in followers" 
              :key="follower.id" 
              class="user-card"
            >
              <div class="user-card-content">
                <div
                  class="user-info"
                  @click="goToUserProfile(follower.user?.id)"
                >
                  <el-avatar
                    :size="60"
                    :src="follower.user?.avatar"
                  />
                  <div class="user-details">
                    <h4>{{ follower.user?.nickname || follower.user?.username }}</h4>
                    <p>{{ follower.user?.bio || '这个人很懒，什么都没写' }}</p>
                    <div class="user-stats">
                      <span>文章 {{ follower.user?.article_count || 0 }}</span>
                      <span>粉丝 {{ follower.user?.follower_count || 0 }}</span>
                    </div>
                  </div>
                </div>
                <div
                  v-if="userStore.isLoggedIn && follower.user?.id !== userStore.user?.id"
                  class="user-actions"
                >
                  <el-button
                    :type="follower.is_following ? 'default' : 'primary'"
                    size="small"
                    :loading="followLoading"
                    @click.stop="handleFollowToggle(follower)"
                  >
                    {{ follower.is_following ? '已关注' : '关注' }}
                  </el-button>
                </div>
              </div>
            </el-card>
          </div>
          
          <div
            v-if="followersTotal > 0"
            class="pagination"
          >
            <el-pagination
              v-model:current-page="followersPage"
              :page-size="20"
              :total="followersTotal"
              layout="prev, pager, next"
              @current-change="loadFollowers"
            />
          </div>
          
          <el-empty
            v-if="followers.length === 0"
            description="还没有粉丝"
          />
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Plus, Check, Star, View, Clock } from '@element-plus/icons-vue'
import api from '@/utils/api'
import { useUserStore } from '@/stores/user'
import { useTerminalStore } from '@/stores/terminal'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const terminalStore = useTerminalStore()

const user = ref(null)
const followStats = ref(null)
const activeTab = ref('articles')
const followLoading = ref(false)

// 文章
const articles = ref([])
const articlesPage = ref(1)
const articlesTotal = ref(0)
const pageSize = ref(10)
const articleSortBy = ref('latest') // latest: 最新排序, hot: 最热排序

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

// 作品
const works = ref([])
const worksPage = ref(1)
const worksPageSize = ref(12)
const worksTotal = ref(0)
const workFilterType = ref('all')
const workSortBy = ref('latest') // latest: 最新, hot: 最热

const workTypeOptions = [
  { label: '全部', value: 'all' },
  { label: '💻 项目', value: 'project' },
  { label: '📷 摄影', value: 'photography' }
]

const userId = computed(() => {
  const id = route.params.id
  return id ? Number(id) : null
})
const isCurrentUser = computed(() => {
  // 严格比较：确保类型一致
  if (!userStore.user?.id || !userId.value) return false
  return Number(userStore.user.id) === Number(userId.value)
})

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD')

// 加载用户信息
const loadUser = async () => {
  try {
    // 查看自己的主页：使用 /api/profile（需要认证，返回完整信息）
    // 查看他人的主页：使用 /api/users/:id（公开API，只返回公开信息）
    if (isCurrentUser.value) {
      const response = await api.get('/profile')
      user.value = response.data
    } else {
      const response = await api.get(`/users/${userId.value}`)
      user.value = response.data
    }
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
    const params = {
      page: articlesPage.value,
      page_size: pageSize.value,
      sort_by: articleSortBy.value // latest 或 hot
    }
    const response = await api.get(`/users/${userId.value}/articles`, { params })
    articles.value = response.data.list || []
    articlesTotal.value = response.data.total || 0
  } catch (error) {
    console.error('Failed to load articles:', error)
  }
}

// 处理文章排序变化
const handleArticleSortChange = () => {
  articlesPage.value = 1
  loadArticles()
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

// 加载作品列表
const loadWorks = async () => {
  try {
    const params = {
      page: worksPage.value,
      page_size: worksPageSize.value,
      sort_by: workSortBy.value // latest 或 hot
    }
    
    if (workFilterType.value !== 'all') {
      params.type = workFilterType.value
    }
    
    const response = await api.get(`/users/${userId.value}/works`, { params })
    works.value = response.data.list || []
    worksTotal.value = response.data.total || 0
  } catch (error) {
    console.error('Failed to load works:', error)
  }
}

// 处理作品筛选变化
const handleWorkFilterChange = () => {
  worksPage.value = 1
  loadWorks()
}

// 处理作品排序变化
const handleWorkSortChange = () => {
  worksPage.value = 1
  loadWorks()
}

// 关注用户
const handleFollow = async () => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    return
  }

  // 不能关注自己
  if (userId.value === userStore.user?.id) {
    ElMessage.warning('不能关注自己')
    return
  }

  followLoading.value = true
  try {
    await api.post(`/users/${userId.value}/follow`)
    ElMessage.success('关注成功')
    // 更新关注统计
    await loadFollowStats()
    // 重新加载用户信息（更新粉丝数）
    await loadUser()
    // 如果当前在关注/粉丝标签页，重新加载列表
    if (activeTab.value === 'following' || activeTab.value === 'followers') {
      if (activeTab.value === 'following') {
        await loadFollowing()
      } else {
        await loadFollowers()
      }
    }
  } catch (error) {
    // 错误消息已经在响应拦截器中显示了，这里不需要再次显示
    // 只需要处理业务逻辑（如刷新状态）
    const errorMsg = error.response?.data?.message || error.message || '关注失败'
    // 如果错误是"已经关注过该用户"或"不能关注自己"，刷新关注状态
    if (errorMsg.includes('已经关注过') || errorMsg.includes('不能关注自己')) {
      await loadFollowStats()
      await loadUser()
    }
  } finally {
    followLoading.value = false
  }
}

// 取消关注
const handleUnfollow = async () => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    return
  }

  followLoading.value = true
  try {
    await api.delete(`/users/${userId.value}/follow`)
    ElMessage.success('取消关注成功')
    // 更新关注统计
    await loadFollowStats()
    // 重新加载用户信息（更新粉丝数）
    await loadUser()
    // 如果当前在关注/粉丝标签页，重新加载列表
    if (activeTab.value === 'following' || activeTab.value === 'followers') {
      if (activeTab.value === 'following') {
        await loadFollowing()
      } else {
        await loadFollowers()
      }
    }
  } catch (error) {
    // 错误消息已经在响应拦截器中显示了，这里不需要再次显示
    // 只需要处理业务逻辑（如刷新状态）
    const errorMsg = error.response?.data?.message || error.message || '操作失败'
    // 如果错误是"未关注该用户"，刷新状态
    if (errorMsg.includes('未关注')) {
      await loadFollowStats()
      await loadUser()
    }
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

// 跳转到作品列表
const goToWorks = () => {
  if (isCurrentUser.value) {
    // 本人 → 跳转到作品管理页面
    router.push('/dashboard/works')
  } else {
    // 他人 → 切换到作品标签页
    activeTab.value = 'works'
    setTimeout(() => {
      document.querySelector('.profile-tabs')?.scrollIntoView({ behavior: 'smooth' })
    }, 100)
  }
}

// 跳转到粉丝列表（切换到粉丝标签页）
const goToFollowers = () => {
  if (!isCurrentUser.value) {
    return
  }
  activeTab.value = 'followers'
  // 每次点击都重新加载数据，确保数据是最新的
  loadFollowers()
  // 同时更新关注统计
  loadFollowStats()
  setTimeout(() => {
    document.querySelector('.profile-tabs')?.scrollIntoView({ behavior: 'smooth' })
  }, 100)
}

// 跳转到关注列表（切换到关注标签页）
const goToFollowing = () => {
  if (!isCurrentUser.value) {
    return
  }
  activeTab.value = 'following'
  // 每次点击都重新加载数据，确保数据是最新的
  loadFollowing()
  // 同时更新关注统计
  loadFollowStats()
  setTimeout(() => {
    document.querySelector('.profile-tabs')?.scrollIntoView({ behavior: 'smooth' })
  }, 100)
}

// 跳转到收藏列表
const goToFavorites = () => {
  if (!isCurrentUser.value) {
    return
  }
  activeTab.value = 'favorites'
  // 每次点击都重新加载数据，确保数据是最新的
  loadFavorites()
  setTimeout(() => {
    document.querySelector('.profile-tabs')?.scrollIntoView({ behavior: 'smooth' })
  }, 100)
}

// 初始化数据加载
const initData = async () => {
  // 重置所有状态
  user.value = null
  followStats.value = null
  articles.value = []
  favorites.value = []
  following.value = []
  followers.value = []
  works.value = []
  articlesPage.value = 1
  favoritesPage.value = 1
  followingPage.value = 1
  followersPage.value = 1
  worksPage.value = 1
  worksTotal.value = 0
  activeTab.value = 'articles'
  
  // 加载数据
  await Promise.all([
    loadUser(),
    loadFollowStats(),
    loadArticles(),
    loadWorks() // 加载作品列表以获取作品总数
  ])
  
  // 如果URL中有tab参数，切换到对应的标签页（仅自己的主页）
  if (route.query.tab && isCurrentUser.value) {
    const tab = route.query.tab
    // 只允许访问自己的隐私标签页
    if (['following', 'followers', 'favorites'].includes(tab)) {
      activeTab.value = tab
      // 根据tab参数加载对应的数据，确保数据是最新的
      if (activeTab.value === 'following') {
        loadFollowing()
        loadFollowStats()
      } else if (activeTab.value === 'followers') {
        loadFollowers()
        loadFollowStats()
      } else if (activeTab.value === 'favorites') {
        loadFavorites()
      }
    }
  }
}

onMounted(() => {
  initData()
})

// 监听路由参数变化，当用户ID变化时重新加载数据
watch(
  () => route.params.id,
  (newId, oldId) => {
    // 只有当用户ID真正变化时才重新加载
    if (newId !== oldId) {
      initData()
    }
  },
  { immediate: false }
)

watch(
  () => terminalStore.refreshSignals[`user:${route.params.id}`],
  async (signal, previousSignal) => {
    if (!signal || signal === previousSignal) return
    await Promise.all([loadUser(), loadFollowStats()])
  }
)

// 监听tab切换
const onTabChange = (tabName) => {
  // 只有自己的主页才能访问隐私标签页
  if (!isCurrentUser.value) {
    // 如果切换到隐私标签页，强制切换回文章标签页
    if (['favorites', 'following', 'followers'].includes(tabName)) {
      activeTab.value = 'articles'
      ElMessage.warning('无权查看他人的隐私数据')
      return
    }
  }
  
  // 每次切换标签页都重新加载数据，确保数据是最新的
  if (tabName === 'favorites' && isCurrentUser.value) {
    loadFavorites()
  } else if (tabName === 'following' && isCurrentUser.value) {
    loadFollowing()
    // 切换关注列表时，更新关注统计
    loadFollowStats()
  } else if (tabName === 'followers' && isCurrentUser.value) {
    loadFollowers()
    // 切换粉丝列表时，更新关注统计
    loadFollowStats()
  } else if (tabName === 'works') {
    loadWorks()
  }
}

// 跳转到用户主页
const goToUserProfile = (targetUserId) => {
  if (targetUserId) {
    router.push(`/users/${targetUserId}`)
  }
}

// 处理标签点击
const handleTagClick = (tagId) => {
  router.push(`/blog?tag_id=${tagId}`)
}

// 根据标签名称生成颜色（用于书签样式）
const getTagColor = (tagName) => {
  const colors = [
    '#409EFF', '#67C23A', '#E6A23C', '#F56C6C', '#909399',
    '#409EFF', '#67C23A', '#E6A23C', '#F56C6C', '#909399'
  ]
  // 根据标签名称的hash值选择颜色
  let hash = 0
  for (let i = 0; i < tagName.length; i++) {
    hash = tagName.charCodeAt(i) + ((hash << 5) - hash)
  }
  return colors[Math.abs(hash) % colors.length]
}

// 在关注/粉丝列表中关注/取消关注
const handleFollowToggle = async (item) => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    return
  }
  
  const targetUserId = item.user?.id
  if (!targetUserId) {
    ElMessage.error('用户ID无效')
    return
  }
  
  // 不能关注自己
  if (targetUserId === userStore.user?.id) {
    ElMessage.warning('不能关注自己')
    return
  }
  
  followLoading.value = true
  try {
    if (item.is_following) {
      // 取消关注
      await api.delete(`/users/${targetUserId}/follow`)
      ElMessage.success('已取消关注')
      item.is_following = false
      // 更新关注统计
      await loadFollowStats()
      // 重新加载用户信息（更新关注数）
      await loadUser()
      // 如果当前在关注列表，重新加载列表（因为取消关注后，该用户会从关注列表中消失）
      if (activeTab.value === 'following') {
        await loadFollowing()
      } else if (activeTab.value === 'followers') {
        // 如果在粉丝列表，也需要重新加载（因为粉丝数可能变化）
        await loadFollowers()
      }
    } else {
      // 关注
      await api.post(`/users/${targetUserId}/follow`)
      ElMessage.success('关注成功')
      item.is_following = true
      // 更新关注统计
      await loadFollowStats()
      // 重新加载用户信息（更新关注数）
      await loadUser()
      // 如果在粉丝列表，重新加载（因为关注状态变化）
      if (activeTab.value === 'followers') {
        await loadFollowers()
      }
    }
  } catch (error) {
    // 错误消息已经在响应拦截器中显示了，这里不需要再次显示
    // 只需要处理业务逻辑（如刷新状态）
    const errorMsg = error.response?.data?.message || error.message || '操作失败'
    // 如果错误是"已经关注过该用户"，刷新状态
    if (errorMsg.includes('已经关注过')) {
      item.is_following = true
      await loadFollowStats()
      await loadUser()
    } else if (errorMsg.includes('未关注')) {
      item.is_following = false
      await loadFollowStats()
      await loadUser()
    }
  } finally {
    followLoading.value = false
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

/* 用户文章列表样式（独立样式，不与博客列表页共用） */
.article-sort-header {
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid var(--theme-border);
}

.user-article-list {
  margin-bottom: 30px;
  display: flex;
  flex-direction: column;
  gap: 8px !important;
}

.user-article-list :deep(.el-card) {
  margin-bottom: 0 !important;
}

.user-article-item {
  margin-bottom: 0 !important;
  cursor: pointer;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  height: 140px;
  overflow: hidden;
  position: relative;
  transition: all 0.3s;
}

.user-article-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px 0 rgba(0, 0, 0, 0.15);
}

.user-article-item :deep(.el-card__body) {
  height: 100%;
  padding: 14px;
  display: flex;
  flex-direction: column;
}

.user-article-content {
  display: flex;
  gap: 14px;
  height: 100%;
  align-items: center;
}

.user-article-cover {
  width: 140px;
  height: 112px;
  object-fit: cover;
  border-radius: 4px;
  flex-shrink: 0;
}

.user-article-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-width: 0;
  min-height: 0;
  overflow: hidden;
}

.user-article-header {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 6px;
}

.user-article-header h2 {
  font-size: 1.2rem;
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  line-height: 1.3;
  color: var(--theme-text-primary);
  word-break: break-word;
  max-height: 1.3em;
  flex-shrink: 0;
  font-weight: 500;
}

.user-article-summary {
  color: var(--theme-text-secondary);
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: 1.5;
  font-size: 13px;
  word-break: break-word;
  max-height: 3em;
  flex: 1;
  min-height: 0;
}

.user-article-meta {
  display: flex;
  gap: 12px;
  color: var(--theme-text-secondary);
  font-size: 13px;
  flex-wrap: wrap;
  align-items: center;
}

.user-article-meta span {
  display: flex;
  align-items: center;
  gap: 5px;
}

.user-article-author {
  display: flex;
  align-items: center;
  gap: 6px;
}

/* 书签样式标签 */
.user-article-bookmarks {
  display: flex;
  flex-direction: row;
  gap: 4px;
  flex-wrap: wrap;
  margin-left: auto;
}

.user-bookmark-tag {
  position: relative;
  padding: 3px 10px 3px 6px;
  color: white;
  font-size: 11px;
  font-weight: 500;
  line-height: 1.3;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
  white-space: nowrap;
  clip-path: polygon(0 0, calc(100% - 6px) 0, 100% 50%, calc(100% - 6px) 100%, 0 100%);
}

.user-bookmark-tag:hover {
  transform: translateX(-2px);
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.3);
  opacity: 0.9;
}

/* 用户收藏列表样式（独立样式，不与博客列表页共用） */
.user-favorite-list {
  margin-bottom: 30px;
  display: flex;
  flex-direction: column;
  gap: 8px !important;
}

.user-favorite-list :deep(.el-card) {
  margin-bottom: 0 !important;
}

.user-favorite-item {
  margin-bottom: 0 !important;
  cursor: pointer;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  height: 140px;
  overflow: hidden;
  position: relative;
}

.user-favorite-item :deep(.el-card__body) {
  height: 100%;
  padding: 14px;
  display: flex;
  flex-direction: column;
}

.user-favorite-content {
  display: flex;
  gap: 14px;
  height: 100%;
  align-items: center; /* 垂直居中对齐 */
}

.user-favorite-cover {
  width: 140px;
  height: 112px;
  object-fit: cover;
  border-radius: 4px;
  flex-shrink: 0;
}

.user-favorite-work-cover-container {
  position: relative;
  width: 140px;
  height: 112px;
  flex-shrink: 0;
  border-radius: 4px;
  overflow: hidden;
}

.user-favorite-work-cover-container .user-favorite-cover {
  width: 100%;
  height: 100%;
}

.user-favorite-work-type-badge {
  position: absolute;
  top: 8px;
  left: 8px;
  z-index: 1;
}

.user-favorite-work-item {
  /* 作品收藏项的特殊样式（如果需要） */
}

.user-favorite-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-width: 0; /* 允许flex收缩 */
  min-height: 0;
  overflow: hidden; /* 防止内容溢出 */
}

.user-favorite-header {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 6px;
}

.user-favorite-header h2 {
  font-size: 1.2rem;
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  line-height: 1.3;
  margin-bottom: 6px;
  color: var(--theme-text-primary);
  word-break: break-word;
  max-height: 1.3em; /* 确保只显示1行 */
  flex-shrink: 0; /* 防止标题被压缩 */
  font-weight: 500;
}

.user-favorite-summary {
  color: var(--theme-text-secondary);
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: 1.5; /* 使用标准行高 */
  font-size: 13px;
  word-break: break-word;
  max-height: 3em; /* 确保只显示2行 (1.5 * 2 = 3em) */
  flex: 1;
  min-height: 0; /* 允许flex收缩 */
}

.user-favorite-meta {
  display: flex;
  gap: 12px;
  color: var(--theme-text-secondary);
  font-size: 13px;
  flex-wrap: wrap;
  align-items: center;
}

.user-favorite-meta span {
  display: flex;
  align-items: center;
  gap: 5px;
}

.user-favorite-author {
  display: flex;
  align-items: center;
  gap: 6px;
}

/* 书签样式标签 */
.user-favorite-bookmarks {
  display: flex;
  flex-direction: row; /* 横向排列 */
  gap: 4px; /* 标签之间的间距 */
  flex-wrap: wrap;
  margin-left: auto; /* 自动推到右侧 */
}

.user-favorite-bookmark-tag {
  position: relative;
  padding: 3px 10px 3px 6px;
  color: white;
  font-size: 11px;
  font-weight: 500;
  line-height: 1.3;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
  white-space: nowrap; /* 防止文字换行 */
  /* 书签折角效果 */
  clip-path: polygon(0 0, calc(100% - 6px) 0, 100% 50%, calc(100% - 6px) 100%, 0 100%);
}

.user-favorite-bookmark-tag:hover {
  transform: translateX(-2px);
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.3);
  opacity: 0.9;
}

.user-list {
  display: grid;
  gap: 20px;
}

.user-card {
  transition: transform 0.3s;
}

.user-card:hover {
  transform: translateY(-3px);
}

.user-card-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
}

.user-info {
  display: flex;
  gap: 20px;
  align-items: center;
  flex: 1;
  cursor: pointer;
}

.user-stats {
  display: flex;
  gap: 15px;
  font-size: 13px;
  color: var(--text-secondary);
  margin-top: 5px;
}

.user-actions {
  flex-shrink: 0;
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

/* 作品列表样式（独立样式，不与作品列表页共用） */
.user-works-filters {
  display: flex;
  align-items: center;
  margin-bottom: 25px;
  padding: 15px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

/* 瀑布流布局 */
.user-works-masonry-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
  margin-bottom: 40px;
}

.user-works-masonry-item {
  cursor: pointer;
  background: white;
  border-radius: 6px;
  overflow: hidden;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.user-works-masonry-item:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.12);
}

.user-works-image-container {
  position: relative;
  width: 100%;
  overflow: hidden;
  aspect-ratio: 4 / 3;
}

.user-works-image {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: cover;
}

.user-works-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(to bottom, rgba(0,0,0,0.3) 0%, transparent 40%, transparent 60%, rgba(0,0,0,0.3) 100%);
  opacity: 0;
  transition: opacity 0.3s;
}

.user-works-masonry-item:hover .user-works-overlay {
  opacity: 1;
}

.user-works-overlay-content {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  padding: 15px;
}

.user-works-type-badge {
  display: flex;
  gap: 5px;
}

.user-works-info {
  padding: 15px;
}

.user-works-title {
  font-size: 1rem;
  margin: 0 0 10px 0;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: 1.4;
  color: var(--text-primary);
  font-weight: 500;
}

.user-works-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
}

.user-works-author {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.85rem;
  color: var(--text-secondary);
  flex: 1;
  min-width: 0;
}

.user-works-author span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-works-stats {
  display: flex;
  gap: 12px;
  font-size: 0.85rem;
  color: var(--text-secondary);
  flex-shrink: 0;
}

.user-works-stats span {
  display: flex;
  align-items: center;
  gap: 3px;
}

/* Magazine adaptation */
.user-profile { padding: 62px 0 80px; background: var(--theme-bg-primary); }
.user-profile .container { max-width: 1060px; padding: 0 32px; }
.profile-card { margin-bottom: 48px; border: 0; border-bottom: 1px solid var(--theme-border); border-radius: 0; box-shadow: none; background: transparent; }
.profile-card :deep(.el-card__body) { padding: 0 0 42px; }
.profile-kicker { margin-bottom: 22px; color: var(--theme-primary); font-family: Georgia, 'Songti SC', serif; font-size: 11px; letter-spacing: .26em; }
.profile-header { align-items: flex-start; }
.profile-header > :deep(.el-avatar) { border: 1px solid var(--theme-border); }
.profile-info h2 { margin-top: -8px; font-family: Georgia, 'Songti SC', 'Noto Serif SC', SimSun, serif; font-size: clamp(34px, 5vw, 48px); font-weight: 500; letter-spacing: .04em; }
.stats { gap: 0; border-top: 1px solid var(--theme-border); border-bottom: 1px solid var(--theme-border); }
.stat-item { min-width: 82px; padding: 14px 18px; border-right: 1px solid var(--theme-border); }
.stat-value { color: var(--theme-primary); font-family: Georgia, serif; font-weight: 500; }
.stat-item.clickable:hover { background: var(--theme-bg-secondary); transform: none; }
.profile-actions { margin-top: 20px; }
.profile-tabs :deep(.el-tabs__header) { margin-bottom: 34px; }
.profile-tabs :deep(.el-tabs__nav-wrap::after) { height: 1px; background: var(--theme-border); }
.profile-tabs :deep(.el-tabs__active-bar) { height: 1px; background: var(--theme-primary); }
.profile-tabs :deep(.el-tabs__item) { font-size: 13px; letter-spacing: .08em; }
.article-sort-header { border-bottom: 1px solid var(--theme-border); }
.user-article-list, .user-favorite-list { gap: 0 !important; padding: 10px 32px; border: 1px solid var(--theme-border); background: var(--theme-bg-card); }
.user-article-item, .user-favorite-item, .user-card { height: auto; min-height: 142px; border: 0 !important; border-bottom: 1px solid var(--theme-border) !important; border-radius: 0; box-shadow: none; transition: padding-left .25s ease; }
.user-article-item:hover, .user-favorite-item:hover, .user-card:hover { padding-left: 8px; box-shadow: none; transform: none; }
.user-article-item :deep(.el-card__body), .user-favorite-item :deep(.el-card__body) { padding: 26px 8px; }
.user-article-header h2, .user-favorite-header h2, .user-details h4, .user-works-title { font-family: Georgia, 'Songti SC', serif; font-weight: 500; letter-spacing: .03em; }
.user-article-cover, .user-favorite-cover, .user-favorite-work-cover-container { border-radius: 0; filter: saturate(.82) contrast(.96); }
.user-bookmark-tag, .user-favorite-bookmark-tag { padding: 3px 7px; background: transparent !important; border: 1px solid var(--theme-border); color: var(--theme-text-secondary); clip-path: none; box-shadow: none; }
.user-bookmark-tag:hover, .user-favorite-bookmark-tag:hover { border-color: var(--theme-primary); color: var(--theme-primary); box-shadow: none; }
.user-works-filters { padding: 0 0 20px; background: transparent; border-bottom: 1px solid var(--theme-border); border-radius: 0; box-shadow: none; }
.user-works-masonry-grid { grid-template-columns: repeat(3, 1fr); }
.user-works-masonry-item { background: transparent; border: 1px solid var(--theme-border); border-radius: 0; box-shadow: none; }
.user-works-masonry-item:hover { border-color: var(--theme-primary); box-shadow: none; transform: translateY(-4px); }
.user-works-image { filter: saturate(.82) contrast(.96); transition: transform .6s cubic-bezier(.2,.6,.2,1); }
.user-works-masonry-item:hover .user-works-image { transform: scale(1.045); }
.user-profile :deep(.el-button), .user-profile :deep(.el-tag), .user-profile :deep(.el-radio-button__inner), .user-profile :deep(.el-segmented) { border-radius: 1px; box-shadow: none; }
.stat-item.clickable:focus-visible, .user-article-item:focus-visible, .user-favorite-item:focus-visible, .user-works-masonry-item:focus-visible, .user-card:focus-visible { outline: 2px solid var(--theme-primary); outline-offset: 3px; }

@media (max-width: 900px) {
  .user-profile .container { padding: 0 24px; }
  .profile-header { flex-direction: column; }
  .stats { flex-wrap: wrap; }
  .user-works-masonry-grid { grid-template-columns: repeat(2, 1fr); }
}

@media (max-width: 560px) {
  .user-profile { padding: 38px 0 56px; }
  .user-profile .container { padding: 0 18px; }
  .profile-header { text-align: left; }
  .stats { justify-content: flex-start; }
  .stat-item { min-width: 33.333%; flex: 1; }
  .user-article-item, .user-favorite-item { height: auto; }
  .user-article-item :deep(.el-card__body), .user-favorite-item :deep(.el-card__body) { padding: 23px 2px; }
  .user-article-list, .user-favorite-list { padding: 4px 16px; }
  .user-article-content, .user-favorite-content { align-items: stretch; }
  .user-works-masonry-grid { grid-template-columns: 1fr; }
}

@media (max-width: 768px) {
  .profile-header {
    flex-direction: column;
    text-align: center;
  }

  .stats {
    justify-content: center;
  }

  .user-article-content {
    flex-direction: column;
  }

  .user-article-cover {
    width: 100%;
    height: 200px;
  }
  
  .user-article-bookmarks {
    margin-left: 0;
    width: 100%;
    margin-top: 8px;
  }

  .user-favorite-content {
    flex-direction: column;
  }

  .user-favorite-cover {
    width: 100%;
    height: 200px;
  }

  .user-favorite-work-cover-container {
    width: 100%;
    height: 200px;
  }
  
  .user-favorite-bookmarks {
    margin-left: 0;
    width: 100%;
    margin-top: 8px;
  }

  .user-works-masonry-grid {
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 15px;
  }
  
  .user-works-filters {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }
  
  .user-works-filters .el-radio-group {
    width: 100% !important;
    margin-left: 0 !important;
  }
  
  .user-works-meta {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
}
</style>

