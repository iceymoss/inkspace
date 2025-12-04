<template>
  <div class="my-comments">
    <el-card>
      <template #header>
        <div class="header">
          <span>我的评论</span>
        </div>
      </template>

      <!-- 评论列表 -->
      <el-table :data="comments" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="文章" min-width="200">
          <template #default="{ row }">
            <el-link 
              type="primary" 
              @click="handleViewArticle(row.article_id)"
              :underline="false"
            >
              {{ row.article?.title || `文章 #${row.article_id}` }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="content" label="评论内容" min-width="300" show-overflow-tooltip />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag 
              :type="getStatusType(row.status)"
              size="small"
            >
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="统计" width="120">
          <template #default="{ row }">
            <div class="stats-text">
              <span>❤️ {{ row.like_count || 0 }}</span>
              <span>💬 {{ row.reply_count || 0 }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleView(row)">查看</el-button>
            <el-button type="danger" size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="fetchComments"
        @current-change="fetchComments"
        class="mt-20"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/utils/api'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const comments = ref([])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 获取我的评论列表
const fetchComments = async () => {
  // 确保用户已登录
  if (!userStore.isLoggedIn || !userStore.user) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }

  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      user_id: userStore.user.id
    }

    const response = await api.get('/comments', { params })
    comments.value = response.data.list || []
    pagination.total = response.data.total || 0
  } catch (error) {
    ElMessage.error('获取评论列表失败')
  } finally {
    loading.value = false
  }
}

const getStatusType = (status) => {
  switch (status) {
    case 1:
      return 'success'
    case 0:
      return 'warning'
    case -1:
      return 'danger'
    default:
      return 'info'
  }
}

const getStatusText = (status) => {
  switch (status) {
    case 1:
      return '已通过'
    case 0:
      return '待审核'
    case -1:
      return '已拒绝'
    default:
      return '未知'
  }
}

const handleView = (row) => {
  router.push(`/blog/${row.article_id}`)
}

const handleViewArticle = (articleId) => {
  router.push(`/blog/${articleId}`)
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除这条评论吗？此操作不可恢复！`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'error'
      }
    )

    await api.delete(`/comments/${row.id}`)
    ElMessage.success('删除成功')
    fetchComments()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  fetchComments()
})
</script>

<style scoped>
.my-comments {
  max-width: 1400px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stats-text {
  display: flex;
  gap: 10px;
  font-size: 12px;
}

.mt-20 {
  margin-top: 20px;
}
</style>

