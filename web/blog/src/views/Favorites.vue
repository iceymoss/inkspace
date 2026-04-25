<template>
  <div class="favorites">
    <div class="container">
      <div class="page-header">
        <h1>我的收藏</h1>
        <p>共 {{ total }} 项收藏</p>
      </div>

      <!-- 分类标签 -->
      <div class="category-tabs">
        <div class="tab-group">
          <button
            v-for="tab in categoryTabs"
            :key="tab.value"
            :class="['tab-btn', { active: activeCategory === tab.value }]"
            @click="activeCategory = tab.value; handleCategoryChange()"
          >
            {{ tab.label }}
          </button>
        </div>
      </div>

      <!-- 收藏列表 -->
      <div class="favorites-grid">
        <div class="favorite-card-wrapper" v-for="favorite in filteredFavorites" :key="favorite.id">
          <Card class="favorite-card">
            <!-- 文章收藏 -->
            <CardContent v-if="favorite.type === 'article' && favorite.article">
              <div class="card-content" @click="goToArticle(favorite.article.id)">
                <img 
                  v-if="favorite.article.cover" 
                  :src="favorite.article.cover" 
                  class="cover-image" 
                />
                <div v-else class="cover-placeholder">
                  <FileText :size="48" />
                </div>
                <div class="content-info">
                  <h3 class="title">{{ favorite.article.title }}</h3>
                  <p class="summary">{{ favorite.article.summary || '暂无摘要' }}</p>
                </div>
              </div>
              <div class="card-meta">
                <span><UserIcon :size="14" /> {{ favorite.article.author?.nickname || '未知' }}</span>
                <span><Eye :size="14" /> {{ favorite.article.view_count || 0 }}</span>
                <span><Clock :size="14" /> {{ formatDate(favorite.created_at) }}</span>
              </div>
              <div class="card-actions">
                <Button 
                  variant="destructive" 
                  size="sm"
                  @click.stop="handleRemoveArticleFavorite(favorite)"
                >
                  <Trash2 :size="14" class="icon-left" /> 取消收藏
                </Button>
              </div>
            </CardContent>

            <!-- 作品收藏 -->
            <CardContent v-else-if="favorite.type === 'work' && favorite.work">
              <div class="card-content" @click="goToWork(favorite.work.id)">
                <img 
                  v-if="favorite.work.cover" 
                  :src="favorite.work.cover" 
                  class="cover-image" 
                />
                <div v-else class="cover-placeholder">
                  <ImageIcon :size="48" />
                </div>
                <div class="content-info">
                  <h3 class="title">{{ favorite.work.title }}</h3>
                  <p class="summary">{{ favorite.work.description || '暂无描述' }}</p>
                  <Badge v-if="favorite.work.type === 'project'" variant="default">开源项目</Badge>
                  <Badge v-else variant="secondary">摄影作品</Badge>
                </div>
              </div>
              <div class="card-meta">
                <span><UserIcon :size="14" /> {{ favorite.work.author?.nickname || '未知' }}</span>
                <span><Eye :size="14" /> {{ favorite.work.view_count || 0 }}</span>
                <span><Clock :size="14" /> {{ formatDate(favorite.created_at) }}</span>
              </div>
              <div class="card-actions">
                <Button 
                  variant="destructive" 
                  size="sm"
                  @click.stop="handleRemoveWorkFavorite(favorite)"
                >
                  <Trash2 :size="14" class="icon-left" /> 取消收藏
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>

      <div v-if="filteredFavorites.length === 0" class="empty-state">
        <p class="empty-description">{{ getEmptyDescription() }}</p>
        <Button variant="default" @click="goToBrowse">{{ activeCategory === 'article' ? '去浏览文章' : activeCategory === 'work' ? '去浏览作品' : '去逛逛' }}</Button>
      </div>

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

    <!-- Confirmation Dialog -->
    <Dialog :open="showConfirmDialog" @update:open="onConfirmDialogUpdateOpen">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>确认</DialogTitle>
          <DialogDescription>{{ confirmDialogMessage }}</DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" @click="cancelConfirmDialog">取消</Button>
          <Button @click="confirmDialogCallback?.()">确认</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { Trash2, User as UserIcon, Eye, Clock, FileText, Image as ImageIcon } from 'lucide-vue-next'
import api from '@/utils/api'
import dayjs from 'dayjs'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from '@/components/ui/dialog'

const router = useRouter()

const favorites = ref([])
const activeCategory = ref('all')
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)

const categoryTabs = [
  { label: '全部', value: 'all' },
  { label: '文章', value: 'article' },
  { label: '作品', value: 'work' }
]

const showConfirmDialog = ref(false)
const confirmDialogMessage = ref('')
const confirmDialogCallback = ref(null)
let _confirmDialogReject = null

const confirmDialog = (message) => {
  return new Promise((resolve, reject) => {
    confirmDialogMessage.value = message
    _confirmDialogReject = reject
    confirmDialogCallback.value = () => {
      _confirmDialogReject = null
      showConfirmDialog.value = false
      resolve()
    }
    showConfirmDialog.value = true
  })
}

const onConfirmDialogUpdateOpen = (open) => {
  showConfirmDialog.value = open
  if (!open && _confirmDialogReject) {
    const rejectFn = _confirmDialogReject
    _confirmDialogReject = null
    rejectFn('cancel')
  }
}

const cancelConfirmDialog = () => {
  if (_confirmDialogReject) {
    const rejectFn = _confirmDialogReject
    _confirmDialogReject = null
    showConfirmDialog.value = false
    rejectFn('cancel')
  }
}

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
    toast.error('加载失败')
  }
}

const handleRemoveArticleFavorite = async (favorite) => {
  try {
    await confirmDialog('确定要取消收藏这篇文章吗？')
    
    await api.delete(`/articles/${favorite.article_id}/favorite`)
    toast.success('取消收藏成功')
    loadFavorites()
  } catch (error) {
    if (error !== 'cancel') {
      toast.error('操作失败')
    }
  }
}

const handleRemoveWorkFavorite = async (favorite) => {
  try {
    await confirmDialog('确定要取消收藏这个作品吗？')
    
    await api.delete(`/works/${favorite.work_id}/favorite`)
    toast.success('取消收藏成功')
    loadFavorites()
  } catch (error) {
    if (error !== 'cancel') {
      toast.error('操作失败')
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

.tab-group {
  display: inline-flex;
  border-radius: var(--radius-md);
  overflow: hidden;
  border: 1px solid var(--theme-border);
}

.tab-btn {
  padding: var(--spacing-sm) var(--spacing-lg);
  background: var(--theme-bg-card);
  border: none;
  cursor: pointer;
  font-size: var(--font-size-sm);
  color: var(--theme-text-secondary);
  transition: all var(--transition-fast);
}

.tab-btn:not(:last-child) {
  border-right: 1px solid var(--theme-border);
}

.tab-btn.active {
  background: var(--theme-primary);
  color: var(--color-text-inverse);
}

.tab-btn:hover:not(.active) {
  background: var(--theme-bg-hover);
}

.favorites-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  margin-bottom: var(--spacing-lg);
}

@media (max-width: 768px) {
  .favorites-grid {
    grid-template-columns: 1fr;
  }
}

@media (min-width: 769px) and (max-width: 992px) {
  .favorites-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

.favorite-card-wrapper {
  display: flex;
}

.favorite-card {
  width: 100%;
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

.icon-left {
  margin-right: 4px;
}

.empty-state {
  text-align: center;
  padding: var(--spacing-xl) 0;
}

.empty-description {
  color: var(--theme-text-secondary);
  margin-bottom: var(--spacing-md);
}

.pagination {
  margin-top: var(--spacing-lg);
  display: flex;
  justify-content: center;
}
</style>
