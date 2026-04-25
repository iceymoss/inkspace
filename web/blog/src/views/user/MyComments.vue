<template>
  <div class="my-comments">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>我的评论</span>
        </div>
      </template>

      <div class="comments-list">
        <div 
          v-for="comment in comments" 
          :key="comment.id" 
          class="comment-item"
        >
          <div class="comment-target">
            <el-link 
              v-if="comment.article"
              :href="`/blog/${comment.article.id}`"
              type="primary"
            >
              评论了文章：{{ comment.article.title }}
            </el-link>
            <el-link 
              v-else-if="comment.work"
              :href="`/works/${comment.work.id}`"
              type="primary"
            >
              评论了作品：{{ comment.work.title }}
            </el-link>
          </div>
          
          <div class="comment-content">
            {{ comment.content }}
          </div>
          
          <div class="comment-meta">
            <span class="comment-time">{{ formatDate(comment.created_at) }}</span>
            <el-button 
              size="small" 
              text 
              type="danger"
              @click="handleDelete(comment.id)"
            >
              删除
            </el-button>
          </div>
        </div>

        <el-empty v-if="comments.length === 0" description="暂无评论" />
      </div>

      <div class="pagination" v-if="total > 0">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          @current-change="loadComments"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/stores/user'
import api from '@/utils/api'

const userStore = useUserStore()

const comments = ref([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const loadComments = async () => {
  try {
    const response = await api.get('/comments', {
      params: {
        user_id: userStore.user?.id,
        page: currentPage.value,
        page_size: pageSize.value
      }
    })
    comments.value = response.data.list || []
    total.value = response.data.total || 0
  } catch (error) {
    console.error('Failed to load comments:', error)
  }
}

const handleDelete = async (commentId) => {
  try {
    await ElMessageBox.confirm('确定要删除这条评论吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await api.delete(`/comments/${commentId}`)
    ElMessage.success('删除成功')
    loadComments()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return d.toLocaleString('zh-CN')
}

onMounted(() => {
  loadComments()
})
</script>

<style scoped>
.my-comments {
  max-width: 900px;
}

.my-comments .el-card {
  box-shadow: var(--shadow-md);
  border-radius: var(--radius-lg);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: var(--font-size-lg);
  font-weight: 600;
  color: var(--theme-text-primary);
}

.comments-list {
  min-height: 400px;
}

.comment-item {
  padding: var(--spacing-lg);
  border-bottom: 1px solid var(--theme-border-light);
  transition: background-color var(--transition-base);
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-item:hover {
  background-color: var(--theme-bg-hover);
}

.comment-target {
  margin-bottom: var(--spacing-sm);
  font-size: var(--font-size-sm);
}

.comment-content {
  padding: var(--spacing-md);
  background: var(--theme-bg-secondary);
  border-radius: var(--radius-md);
  margin-bottom: var(--spacing-sm);
  line-height: var(--line-height-relaxed);
  color: var(--theme-text-primary);
}

.comment-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: var(--font-size-sm);
  color: var(--theme-text-tertiary);
}

.pagination {
  margin-top: var(--spacing-lg);
  display: flex;
  justify-content: center;
}
</style>
