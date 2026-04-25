<template>
  <div class="favorites">
    <div class="container">
      <div class="page-header">
        <h1>我的收藏</h1>
        <p>共 {{ total }} 项收藏</p>
      </div>

      <!-- 分类标签 -->
      <div class="category-tabs">
        <el-radio-group v-model="activeCategory" @change="handleCategoryChange">
          <el-radio-button label="all">全部</el-radio-button>
          <el-radio-button label="article">文章</el-radio-button>
          <el-radio-button label="work">作品</el-radio-button>
        </el-radio-group>
      </div>

      <!-- 收藏列表 -->
      <el-row :gutter="20" class="favorites-grid">
        <el-col :xs="24" :sm="12" :md="8" v-for="favorite in filteredFavorites" :key="favorite.id">
          <el-card class="favorite-card" shadow="hover">
            <!-- 文章收藏 -->
            <template v-if="favorite.type === 'article' && favorite.article">
              <div class="card-content" @click="goToArticle(favorite.article.id)">
                <img 
                  v-if="favorite.article.cover" 
                  :src="favorite.article.cover" 
                  class="cover-image" 
                />
                <div v-else class="cover-placeholder">
                  <el-icon><Document /></el-icon>
                </div>
                <div class="content-info">
                  <h3 class="title">{{ favorite.article.title }}</h3>
                  <p class="summary">{{ favorite.article.summary || '暂无摘要' }}</p>
                </div>
              </div>
              <div class="card-meta">
                <span><el-icon><User /></el-icon> {{ favorite.article.author?.nickname || '未知' }}</span>
                <span><el-icon><View /></el-icon> {{ favorite.article.view_count || 0 }}</span>
                <span><el-icon><Clock /></el-icon> {{ formatDate(favorite.created_at) }}</span>
              </div>
              <div class="card-actions">
                <el-button 
                  type="danger" 
                  size="small" 
                  :icon="Delete"
                  @click.stop="handleRemoveArticleFavorite(favorite)"
                >
                  取消收藏
                </el-button>
              </div>
            </template>

            <!-- 作品收藏 -->
            <template v-else-if="favorite.type === 'work' && favorite.work">
              <div class="card-content" @click="goToWork(favorite.work.id)">
                <img 
                  v-if="favorite.work.cover" 
                  :src="favorite.work.cover" 
                  class="cover-image" 
                />
                <div v-else class="cover-placeholder">
                  <el-icon><Picture /></el-icon>
                </div>
                <div class="content-info">
                  <h3 class="title">{{ favorite.work.title }}</h3>
                  <p class="summary">{{ favorite.work.description || '暂无描述' }}</p>
                  <el-tag v-if="favorite.work.type === 'project'" type="primary" size="small">开源项目</el-tag>
                  <el-tag v-else type="warning" size="small">摄影作品</el-tag>
                </div>
              </div>
              <div class="card-meta">
                <span><el-icon><User /></el-icon> {{ favorite.work.author?.nickname || '未知' }}</span>
                <span><el-icon><View /></el-icon> {{ favorite.work.view_count || 0 }}</span>
                <span><el-icon><Clock /></el-icon> {{ formatDate(favorite.created_at) }}</span>
              </div>
              <div class="card-actions">
                <el-button 
                  type="danger" 
                  size="small" 
                  :icon="Delete"
                  @click.stop="handleRemoveWorkFavorite(favorite)"
                >
                  取消收藏
                </el-button>
              </div>
            </template>
          </el-card>
        </el-col>
      </el-row>

      <el-empty v-if="filteredFavorites.length === 0" :description="getEmptyDescription()">
        <el-button type="primary" @click="goToBrowse">{{ activeCategory === 'article' ? '去浏览文章' : activeCategory === 'work' ? '去浏览作品' : '去逛逛' }}</el-button>
      </el-empty>

      <div class="pagination" v-if="total > 0">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="loadFavorites"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete, User, View, Clock, Document, Picture } from '@element-plus/icons-vue'
import api from '@/utils/api'
import dayjs from 'dayjs'

const router = useRouter()

const favorites = ref([])
const activeCategory = ref('all')
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD')

// 根据分类过滤收藏
const filteredFavorites = computed(() => {
  if (activeCategory.value === 'all') {
    return favorites.value
  }
  return favorites.value.filter(f => f.type === activeCategory.value)
})

const getEmptyDescription = () => {
  if (activeCategory.value === 'article') {
    return '还没有收藏任何文章'
  } else if (activeCategory.value === 'work') {
    return '还没有收藏任何作品'
  }
  return '还没有任何收藏'
}

const goToBrowse = () => {
  if (activeCategory.value === 'article') {
    router.push('/blog')
  } else if (activeCategory.value === 'work') {
    router.push('/works')
  } else {
    router.push('/')
  }
}

const handleCategoryChange = () => {
  currentPage.value = 1
  // 不需要重新加载，因为数据已经在内存中
}

const loadFavorites = async () => {
  try {
    const response = await api.get('/favorites', {
      params: {
        page: currentPage.value,
        page_size: pageSize.value
      }
    })
    favorites.value = response.data.list || []
    total.value = response.data.total || 0
  } catch (error) {
    ElMessage.error('加载失败')
  }
}

const handleRemoveArticleFavorite = async (favorite) => {
  try {
    await ElMessageBox.confirm('确定要取消收藏这篇文章吗？', '提示', {
      type: 'warning'
    })
    
    await api.delete(`/articles/${favorite.article_id}/favorite`)
    ElMessage.success('取消收藏成功')
    loadFavorites()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

const handleRemoveWorkFavorite = async (favorite) => {
  try {
    await ElMessageBox.confirm('确定要取消收藏这个作品吗？', '提示', {
      type: 'warning'
    })
    
    await api.delete(`/works/${favorite.work_id}/favorite`)
    ElMessage.success('取消收藏成功')
    loadFavorites()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

const goToArticle = (articleId) => {
  router.push(`/blog/${articleId}`)
}

const goToWork = (workId) => {
  router.push(`/works/${workId}`)
}

onMounted(() => {
  loadFavorites()
})
</script>

<style scoped>
.favorites {
  padding: var(--spacing-xl) 0;
  background-color: var(--theme-bg-secondary);
  min-height: 100vh;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 var(--spacing-md);
}

.page-header {
  text-align: center;
  margin-bottom: var(--spacing-lg);
}

.page-header h1 {
  margin-bottom: var(--spacing-sm);
  font-size: var(--font-size-3xl);
  color: var(--theme-text-primary);
}

.page-header p {
  color: var(--theme-text-secondary);
}

.category-tabs {
  display: flex;
  justify-content: center;
  margin-bottom: var(--spacing-lg);
}

.category-tabs :deep(.el-radio-button__inner) {
  cursor: pointer;
  transition: all var(--transition-fast);
}

.favorites-grid {
  margin-bottom: var(--spacing-lg);
}

.favorite-card {
  height: 100%;
  display: flex;
  flex-direction: column;
  box-shadow: var(--shadow-md);
  border-radius: var(--radius-lg);
  transition: transform var(--transition-slow), box-shadow var(--transition-slow);
}

.favorite-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-lg);
}

.card-content {
  cursor: pointer;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.cover-image {
  width: 100%;
  height: 200px;
  object-fit: cover;
  border-radius: var(--radius-sm);
  margin-bottom: var(--spacing-md);
}

.cover-placeholder {
  width: 100%;
  height: 200px;
  background-color: var(--theme-bg-hover);
  border-radius: var(--radius-sm);
  margin-bottom: var(--spacing-md);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--theme-text-tertiary);
  font-size: 48px;
}

.content-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.title {
  margin: 0 0 var(--spacing-sm) 0;
  font-size: var(--font-size-lg);
  font-weight: 600;
  color: var(--theme-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: var(--line-height-base);
  min-height: 3em;
}

.summary {
  color: var(--theme-text-secondary);
  margin: 0 0 var(--spacing-sm) 0;
  font-size: var(--font-size-sm);
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: var(--line-height-base);
  min-height: 3em;
  flex: 1;
}

.card-meta {
  display: flex;
  gap: var(--spacing-md);
  color: var(--theme-text-secondary);
  font-size: var(--font-size-sm);
  padding: var(--spacing-md) 0;
  border-top: 1px solid var(--theme-border-light);
  flex-wrap: wrap;
}

.card-meta span {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.card-actions {
  display: flex;
  justify-content: flex-end;
  padding-top: var(--spacing-sm);
}

.card-actions .el-button {
  cursor: pointer;
  transition: all var(--transition-fast);
}

.pagination {
  margin-top: var(--spacing-lg);
  display: flex;
  justify-content: center;
}

:deep(.el-col) {
  display: flex;
}

:deep(.favorite-card) {
  width: 100%;
}
</style>
