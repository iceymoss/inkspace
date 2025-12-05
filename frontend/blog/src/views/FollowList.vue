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
            @click="goToUserProfile(item.following_id)"
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
                v-if="item.following_id !== userStore.user?.id"
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
            @click="goToUserProfile(item.follower_id)"
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
                v-if="item.follower_id !== userStore.user?.id"
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
import { ElMessage } from 'element-plus'
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
    ElMessage.error('获取用户信息失败')
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
    ElMessage.error('获取关注列表失败')
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
    ElMessage.error('获取粉丝列表失败')
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
  const targetUserId = item.following_id || item.follower_id
  
  try {
    if (item.is_following) {
      await api.delete(`/users/${targetUserId}/follow`)
      ElMessage.success('已取消关注')
      item.is_following = false
    } else {
      await api.post(`/users/${targetUserId}/follow`)
      ElMessage.success('关注成功')
      item.is_following = true
    }
  } catch (error) {
    ElMessage.error('操作失败')
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
  padding: 20px;
  background-color: #f5f7fa;
  min-height: 100vh;
}

.follow-list .el-card {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
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
  margin: 0 0 15px 0;
}

.user-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.user-item {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.user-item:hover {
  background-color: #f5f7fa;
  transform: translateX(5px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.user-info {
  flex: 1;
}

.user-info h4 {
  margin: 0 0 5px 0;
  font-size: 16px;
}

.user-info .bio {
  margin: 0 0 8px 0;
  color: #909399;
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-info .stats {
  display: flex;
  gap: 15px;
  font-size: 13px;
  color: #606266;
}

.user-actions {
  flex-shrink: 0;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

:deep(.el-tabs__nav-wrap::after) {
  display: none;
}
</style>

