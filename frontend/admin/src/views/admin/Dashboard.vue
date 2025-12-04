<template>
  <div class="dashboard">
    <h2>控制台</h2>
    
    <el-row :gutter="20">
      <el-col :xs="24" :sm="12" :md="6" v-for="item in stats" :key="item.title">
        <el-card class="stats-card">
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

    <el-row :gutter="20" class="mt-20">
      <el-col :span="24">
        <el-card>
          <template #header>
            <span>快速操作</span>
          </template>
          <el-space wrap>
            <el-button type="primary" @click="$router.push('/articles/create')">
              <el-icon><EditPen /></el-icon> 写文章
            </el-button>
            <el-button @click="$router.push('/works')">
              <el-icon><Picture /></el-icon> 添加作品
            </el-button>
            <el-button @click="$router.push('/categories')">
              <el-icon><Folder /></el-icon> 管理分类
            </el-button>
            <el-button @click="$router.push('/tags')">
              <el-icon><CollectionTag /></el-icon> 管理标签
            </el-button>
          </el-space>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Document, Picture, ChatDotRound, User, EditPen, Folder, CollectionTag } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import adminApi from '@/utils/adminApi'

const stats = ref([
  { title: '文章总数', value: 0, icon: Document, color: '#409eff' },
  { title: '作品总数', value: 0, icon: Picture, color: '#67c23a' },
  { title: '评论总数', value: 0, icon: ChatDotRound, color: '#e6a23c' },
  { title: '用户总数', value: 0, icon: User, color: '#f56c6c' }
])

// 获取统计数据
const fetchStats = async () => {
  try {
    // 获取文章数
    const articlesRes = await adminApi.get('/admin/articles', { params: { page: 1, page_size: 1 } })
    stats.value[0].value = articlesRes.data.total || 0

    // 获取作品数  
    const worksRes = await adminApi.get('/admin/works', { params: { page: 1, page_size: 1 } })
    stats.value[1].value = worksRes.data.total || 0

    // 获取评论数
    const commentsRes = await adminApi.get('/admin/comments', { params: { page: 1, page_size: 1 } })
    stats.value[2].value = commentsRes.data.total || 0

    // 获取用户数
    const usersRes = await adminApi.get('/admin/users', { params: { page: 1, page_size: 1 } })
    stats.value[3].value = usersRes.data.total || 0
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

onMounted(() => {
  fetchStats()
})
</script>

<style scoped>
.dashboard h2 {
  margin-bottom: 20px;
}

.stats-card {
  margin-bottom: 20px;
}

.stats-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.stats-icon {
  width: 60px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  color: white;
}

.stats-info h3 {
  margin: 0 0 5px 0;
  font-size: 28px;
  font-weight: bold;
}

.stats-info p {
  margin: 0;
  color: var(--text-secondary);
}

.mt-20 {
  margin-top: 20px;
}
</style>

