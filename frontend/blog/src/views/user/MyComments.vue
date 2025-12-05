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

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.comments-list {
  min-height: 400px;
}

.comment-item {
  padding: 20px;
  border-bottom: 1px solid #ebeef5;
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-target {
  margin-bottom: 10px;
  font-size: 0.9rem;
}

.comment-content {
  padding: 15px;
  background: #f5f7fa;
  border-radius: 8px;
  margin-bottom: 10px;
  line-height: 1.6;
}

.comment-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.9rem;
  color: #909399;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
