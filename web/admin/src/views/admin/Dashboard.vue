<template>
  <div class="dashboard">
    <h2>控制台</h2>
    
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-5">
      <Card v-for="item in stats" :key="item.title" class="stats-card">
        <CardContent class="p-5">
          <div class="stats-content">
            <div class="stats-icon" :style="{ backgroundColor: item.color }">
              <component :is="item.icon" class="w-8 h-8" />
            </div>
            <div class="stats-info">
              <h3>{{ item.value }}</h3>
              <p>{{ item.title }}</p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>

    <div class="mt-5">
      <Card>
        <CardHeader>
          <span>快速操作</span>
        </CardHeader>
        <CardContent>
          <div class="flex flex-wrap gap-3">
            <Button @click="$router.push('/articles/create')">
              <PencilLine class="mr-2 h-4 w-4" /> 写文章
            </Button>
            <Button variant="outline" @click="$router.push('/works')">
              <ImageIcon class="mr-2 h-4 w-4" /> 添加作品
            </Button>
            <Button variant="outline" @click="$router.push('/categories')">
              <Folder class="mr-2 h-4 w-4" /> 管理分类
            </Button>
            <Button variant="outline" @click="$router.push('/tags')">
              <Tag class="mr-2 h-4 w-4" /> 管理标签
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, markRaw } from 'vue'
import { FileText, Image as ImageIcon, MessageCircle, User as UserIcon, PencilLine, Folder, Tag } from 'lucide-vue-next'
import { Card, CardHeader, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import adminApi from '@/utils/adminApi'

const stats = ref([
  { title: '文章总数', value: 0, icon: markRaw(FileText), color: '#409eff' },
  { title: '作品总数', value: 0, icon: markRaw(ImageIcon), color: '#67c23a' },
  { title: '评论总数', value: 0, icon: markRaw(MessageCircle), color: '#e6a23c' },
  { title: '用户总数', value: 0, icon: markRaw(UserIcon), color: '#f56c6c' }
])

const fetchStats = async () => {
  try {
    const articlesRes = await adminApi.get('/admin/articles', { params: { page: 1, page_size: 1 } })
    stats.value[0].value = articlesRes.data.total || 0

    const worksRes = await adminApi.get('/admin/works', { params: { page: 1, page_size: 1 } })
    stats.value[1].value = worksRes.data.total || 0

    const commentsRes = await adminApi.get('/admin/comments', { params: { page: 1, page_size: 1 } })
    stats.value[2].value = commentsRes.data.total || 0

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
  @apply mb-4 text-2xl font-bold;
  color: var(--color-text-primary);
}

.stats-card {
  @apply mb-0 rounded-md;
  box-shadow: var(--shadow-sm);
  transition: box-shadow var(--transition-base);
}

.stats-card:hover {
  box-shadow: var(--shadow-md);
}

.stats-content {
  @apply flex items-center gap-4;
}

.stats-icon {
  @apply w-[60px] h-[60px] flex items-center justify-center rounded-md;
  color: var(--color-text-inverse);
}

.stats-info h3 {
  @apply m-0 mb-1 text-2xl font-bold;
  color: var(--color-text-primary);
}

.stats-info p {
  @apply m-0 text-sm;
  color: var(--color-text-secondary);
}
</style>
