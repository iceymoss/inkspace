<template>
  <div class="favorites">
    <div class="container">
      <div class="page-header">
        <h1>我的收藏</h1>
        <p>共 {{ total }} 篇文章</p>
      </div>

      <el-row :gutter="20">
        <el-col :xs="24" :sm="12" :md="8" v-for="favorite in favorites" :key="favorite.id">
          <el-card class="favorite-card" shadow="hover">
            <div class="card-header">
              <div class="article-info" @click="goToArticle(favorite.article.id)">
                <img v-if="favorite.article.cover" :src="favorite.article.cover" class="article-cover" />
                <h3>{{ favorite.article.title }}</h3>
                <p class="article-summary">{{ favorite.article.summary }}</p>
              </div>
              <div class="card-actions">
                <el-button 
                  type="danger" 
                  size="small" 
                  :icon="Delete"
                  @click="handleRemoveFavorite(favorite)"
                >
                  取消收藏
                </el-button>
              </div>
            </div>
            
            <div class="article-meta">
              <span><el-icon><User /></el-icon> {{ favorite.article.author?.nickname }}</span>
              <span><el-icon><View /></el-icon> {{ favorite.article.view_count }}</span>
              <span><el-icon><Clock /></el-icon> {{ formatDate(favorite.created_at) }}</span>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-empty v-if="favorites.length === 0" description="还没有收藏任何文章">
        <el-button type="primary" @click="$router.push('/blog')">去逛逛</el-button>
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
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete } from '@element-plus/icons-vue'
import api from '@/utils/api'
import dayjs from 'dayjs'

const router = useRouter()

const favorites = ref([])
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD')

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

const handleRemoveFavorite = async (favorite) => {
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

const goToArticle = (articleId) => {
  router.push(`/blog/${articleId}`)
}

onMounted(() => {
  loadFavorites()
})
</script>

<style scoped>
.favorites {
  padding: 40px 0;
}

.page-header {
  text-align: center;
  margin-bottom: 40px;
}

.page-header h1 {
  margin-bottom: 10px;
}

.page-header p {
  color: var(--text-secondary);
}

.favorite-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.article-info {
  cursor: pointer;
}

.article-cover {
  width: 100%;
  height: 200px;
  object-fit: cover;
  border-radius: 4px;
  margin-bottom: 15px;
}

.favorite-card h3 {
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

.card-actions {
  display: flex;
  justify-content: flex-end;
}

.article-meta {
  display: flex;
  gap: 15px;
  color: var(--text-secondary);
  font-size: 14px;
  padding-top: 15px;
  border-top: 1px solid var(--border-lighter);
}

.article-meta span {
  display: flex;
  align-items: center;
  gap: 5px;
}

.pagination {
  margin-top: 30px;
  display: flex;
  justify-content: center;
}
</style>

