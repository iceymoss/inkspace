<template>
  <div class="user-dashboard">
    <el-card>
      <template #header>
        <span>欢迎回来，{{ userStore.user?.nickname || userStore.user?.username }}！</span>
      </template>

      <el-row :gutter="20">
        <el-col :xs="24" :sm="12" :md="6" v-for="item in stats" :key="item.title">
          <el-card class="stats-card" :class="{ clickable: item.clickable }" @click="handleStatClick(item.action)">
            <div class="stats-content">
              <div class="stats-icon" :style="{ backgroundColor: item.color }">
                <el-icon :size="32"><component :is="item.icon" /></el-icon>
              </div>
              <div class="stats-info">
                <h3>{{ item.value }}</h3>
                <p>{{ item.title }}</p>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-card>

    <el-row :gutter="20" class="mt-20">
      <el-col :span="24">
        <el-card>
          <template #header>
            <span>快速操作</span>
          </template>
          <el-space wrap>
            <el-button type="primary" @click="$router.push('/dashboard/articles/create')">
              <el-icon><EditPen /></el-icon> 写文章
            </el-button>
            <el-button @click="$router.push('/dashboard/articles')">
              <el-icon><Document /></el-icon> 我的文章
            </el-button>
            <el-button @click="$router.push('/favorites')">
              <el-icon><Collection /></el-icon> 我的收藏
            </el-button>
            <el-button @click="$router.push('/profile/edit')">
              <el-icon><User /></el-icon> 个人设置
            </el-button>
          </el-space>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { Document, Collection, ChatDotRound, User, EditPen } from '@element-plus/icons-vue'
import api from '@/utils/api'

const router = useRouter()
const userStore = useUserStore()

const stats = ref([
  { title: '我的文章', value: 0, icon: Document, color: '#409eff', clickable: true, action: 'articles' },
  { title: '我的收藏', value: 0, icon: Collection, color: '#67c23a', clickable: true, action: 'favorites' },
  { title: '我的评论', value: 0, icon: ChatDotRound, color: '#e6a23c', clickable: true, action: 'comments' },
  { title: '粉丝数', value: 0, icon: User, color: '#f56c6c', clickable: true, action: 'followers' }
])

const fetchStats = async () => {
  try {
    const profile = await api.get('/profile')
    stats.value[0].value = profile.data.article_count || 0
    stats.value[1].value = profile.data.favorite_count || 0
    stats.value[2].value = profile.data.comment_count || 0
    stats.value[3].value = profile.data.follower_count || 0
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

// 处理统计卡片点击
const handleStatClick = (action) => {
  if (!action) return
  
  switch (action) {
    case 'articles':
      router.push('/dashboard/articles')
      break
    case 'favorites':
      router.push('/favorites')
      break
    case 'comments':
      router.push('/dashboard/comments')
      break
    case 'followers':
      router.push(`/users/${userStore.user.id}/follows?tab=followers`)
      break
  }
}

onMounted(() => {
  fetchStats()
})
</script>

<style scoped>
.user-dashboard {
  max-width: 1200px;
}

.user-dashboard .el-card {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.stats-card {
  margin-bottom: 20px;
  transition: all 0.3s;
}

.stats-card.clickable {
  cursor: pointer;
}

.stats-card.clickable:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.stats-content {
  display: flex;
  align-items: center;
  gap: 15px;
}

.stats-icon {
  width: 60px;
  height: 60px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stats-info h3 {
  margin: 0;
  font-size: 24px;
  font-weight: bold;
}

.stats-info p {
  margin: 5px 0 0 0;
  color: #909399;
  font-size: 14px;
}

.mt-20 {
  margin-top: 20px;
}
</style>

