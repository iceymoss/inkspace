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
            <el-button type="primary" @click="$router.push('/admin/articles/create')">
              <el-icon><EditPen /></el-icon> 写文章
            </el-button>
            <el-button @click="$router.push('/admin/works')">
              <el-icon><Picture /></el-icon> 添加作品
            </el-button>
            <el-button @click="$router.push('/admin/categories')">
              <el-icon><Folder /></el-icon> 管理分类
            </el-button>
            <el-button @click="$router.push('/admin/tags')">
              <el-icon><CollectionTag /></el-icon> 管理标签
            </el-button>
          </el-space>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Document, Picture, ChatDotRound, User } from '@element-plus/icons-vue'

const stats = ref([
  { title: '文章总数', value: 0, icon: Document, color: '#409eff' },
  { title: '作品总数', value: 0, icon: Picture, color: '#67c23a' },
  { title: '评论总数', value: 0, icon: ChatDotRound, color: '#e6a23c' },
  { title: '访问总数', value: 0, icon: User, color: '#f56c6c' }
])
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

