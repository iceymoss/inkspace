<template>
  <div class="follow-list">
    <div class="container">
      <el-card>
        <template #header>
          <div class="header">
            <div>
              <h2>{{ user?.nickname || user?.username }}</h2>
              <el-tabs v-model="activeTab" @tab-change="handleTabChange">
                <el-tab-pane :label="`关注 ${followingTotal}`" name="following" />
                <el-tab-pane :label="`粉丝 ${followersTotal}`" name="followers" />
              </el-tabs>
            </div>
            <el-button @click="$router.back()">返回</el-button>
          </div>
        </template>

        <!-- 关注列表 -->
        <div v-if="activeTab === 'following'" class="user-list">
          <div 
            v-for="item in followingList" 
            :key="item.id" 
            class="user-item"
            @click="goToUserProfile(item.user?.id)"
          >
            <el-avatar :size="60" :src="item.user?.avatar" />
            <div class="user-info">
              <h4>{{ item.user?.nickname || item.user?.username }}</h4>
              <p class="bio">{{ item.user?.bio || '这个人很懒，什么都没写' }}</p>
              <div class="stats">
                <span>文章 {{ item.user?.article_count || 0 }}</span>
                <span>粉丝 {{ item.user?.follower_count || 0 }}</span>
              </div>
            </div>
            <div class="user-actions" @click.stop>
              <el-button
                v-if="item.user?.id !== userStore.user?.id && userStore.isLoggedIn"
                :type="item.is_following ? 'default' : 'primary'"
                size="small"
                @click="handleFollowToggle(item)"
              >
                {{ item.is_following ? '已关注' : '关注' }}
              </el-button>
            </div>
          </div>

          <el-empty v-if="followingList.length === 0" description="还没有关注任何人" />

          <el-pagination
            v-if="followingTotal > 0"
            v-model:current-page="followingPage"
            :page-size="20"
            :total="followingTotal"
            layout="total, prev, pager, next"
            @current-change="loadFollowing"
            class="pagination"
          />
        </div>

        <!-- 粉丝列表 -->
        <div v-if="activeTab === 'followers'" class="user-list">
          <div 
            v-for="item in followersList" 
            :key="item.id" 
            class="user-item"
            @click="goToUserProfile(item.user?.id)"
          >
            <el-avatar :size="60" :src="item.user?.avatar" />
            <div class="user-info">
              <h4>{{ item.user?.nickname || item.user?.username }}</h4>
              <p class="bio">{{ item.user?.bio || '这个人很懒，什么都没写' }}</p>
              <div class="stats">
                <span>文章 {{ item.user?.article_count || 0 }}</span>
                <span>粉丝 {{ item.user?.follower_count || 0 }}</span>
              </div>
            </div>
            <div class="user-actions" @click.stop>
              <el-button
                v-if="item.user?.id !== userStore.user?.id && userStore.isLoggedIn"
                :type="item.is_following ? 'default' : 'primary'"
                size="small"
                @click="handleFollowToggle(item)"
              >
                {{ item.is_following ? '已关注' : '关注' }}
              </el-button>
            </div>
          </div>

          <el-empty v-if="followersList.length === 0" description="还没有粉丝" />

          <el-pagination
            v-if="followersTotal > 0"
            v-model:current-page="followersPage"
            :page-size="20"
            :total="followersTotal"
            layout="total, prev, pager, next"
            @current-change="loadFollowers"
            class="pagination"
          />
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import api from '@/utils/api'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const user = ref(null)
const activeTab = ref(route.query.tab || 'following')

// 关注列表
const followingList = ref([])
const followingPage = ref(1)
const followingTotal = ref(0)

// 粉丝列表
const followersList = ref([])
const followersPage = ref(1)
const followersTotal = ref(0)

const userId = computed(() => parseInt(route.params.id))

// 加载用户信息
const loadUser = async () => {
  try {
    const response = await api.get(`/users/${userId.value}`)
    user.value = response.data
  } catch (error) {
    toast.error('获取用户信息失败')
  }
}

// 加载关注列表
const loadFollowing = async () => {
  try {
    const response = await api.get(`/users/${userId.value}/following`, {
      params: { page: followingPage.value, page_size: 20 }
    })
    followingList.value = response.data.list || []
    followingTotal.value = response.data.total || 0
  } catch (error) {
    toast.error('获取关注列表失败')
  }
}

// 加载粉丝列表
const loadFollowers = async () => {
  try {
    const response = await api.get(`/users/${userId.value}/followers`, {
      params: { page: followersPage.value, page_size: 20 }
    })
    followersList.value = response.data.list || []
    followersTotal.value = response.data.total || 0
  } catch (error) {
    toast.error('获取粉丝列表失败')
  }
}

// 切换标签
const handleTabChange = (tab) => {
  router.replace({ query: { tab } })
  if (tab === 'following') {
    loadFollowing()
  } else {
    loadFollowers()
  }
}

// 跳转到用户主页
const goToUserProfile = (uid) => {
  router.push(`/users/${uid}`)
}

// 关注/取消关注
const handleFollowToggle = async (item) => {
  if (!userStore.isLoggedIn) {
    toast.warning('请先登录')
    return
  }
  
  const targetUserId = item.user?.id
  if (!targetUserId) {
    toast.error('用户ID无效')
    return
  }
  
  try {
    if (item.is_following) {
      await api.delete(`/users/${targetUserId}/follow`)
      toast.success('已取消关注')
      item.is_following = false
    } else {
      await api.post(`/users/${targetUserId}/follow`)
      toast.success('关注成功')
      item.is_following = true
    }
  } catch (error) {
    toast.error(error.response?.data?.message || '操作失败')
  }
}

onMounted(() => {
  loadUser()
  if (activeTab.value === 'following') {
    loadFollowing()
  } else {
    loadFollowers()
  }
})
</script>

<style scoped>
.follow-list {
  padding: var(--spacing-lg);
  background-color: var(--theme-bg-secondary);
  min-height: 100vh;
}

.follow-list .el-card {
  box-shadow: var(--shadow-md);
  border-radius: var(--radius-lg);
}

.container {
  max-width: 800px;
  margin: 0 auto;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.header h2 {
  margin: 0 0 var(--spacing-md) 0;
  font-size: var(--font-size-2xl);
  color: var(--theme-text-primary);
}

.header .el-button {
  cursor: pointer;
  transition: color var(--transition-fast);
}

.user-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.user-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-md);
  border: 1px solid var(--theme-border-light);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--transition-slow);
}

.user-item:hover {
  background-color: var(--theme-bg-hover);
  transform: translateX(5px);
  box-shadow: var(--shadow-sm);
}

.user-info {
  flex: 1;
}

.user-info h4 {
  margin: 0 0 var(--spacing-xs) 0;
  font-size: var(--font-size-base);
  color: var(--theme-text-primary);
}

.user-info .bio {
  margin: 0 0 var(--spacing-sm) 0;
  color: var(--theme-text-tertiary);
  font-size: var(--font-size-sm);
  line-height: var(--line-height-base);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-info .stats {
  display: flex;
  gap: var(--spacing-md);
  font-size: var(--font-size-sm);
  color: var(--theme-text-secondary);
}

.user-actions {
  flex-shrink: 0;
}

.user-actions .el-button {
  cursor: pointer;
  transition: all var(--transition-fast);
}

.pagination {
  margin-top: var(--spacing-lg);
  display: flex;
  justify-content: center;
}

:deep(.el-tabs__nav-wrap::after) {
  display: none;
}
</style>

