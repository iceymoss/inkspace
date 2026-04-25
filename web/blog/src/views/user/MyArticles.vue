<template>
  <div class="my-articles">
    <el-card>
      <template #header>
        <div class="header">
          <span>我的文章</span>
          <el-button type="primary" @click="$router.push('/dashboard/articles/create')">
            <el-icon><Plus /></el-icon> 写文章
          </el-button>
        </div>
      </template>

      <!-- Tab 标签 -->
      <el-tabs v-model="activeTab" @tab-change="handleTabChange" class="status-tabs">
        <el-tab-pane label="全部" name="all" />
        <el-tab-pane label="已发布" name="published" />
        <el-tab-pane label="私有" name="private" />
        <el-tab-pane label="草稿" name="draft" />
      </el-tabs>

      <!-- 搜索栏 -->
      <el-form :inline="true" class="search-form">
        <el-form-item label="标题">
          <el-input v-model="searchForm.title" placeholder="请输入标题" clearable @clear="handleSearch" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 文章列表 -->
      <el-table :data="articles" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="标题" min-width="200" />
        <el-table-column label="分类" width="120">
          <template #default="{ row }">
            {{ row.category?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="统计" width="150">
          <template #default="{ row }">
            <div class="stats-text">
              <span>👁️ {{ row.view_count }}</span>
              <span>❤️ {{ row.like_count }}</span>
              <span>💬 {{ row.comment_count }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleView(row)">查看</el-button>
            <el-button type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button 
              v-if="row.status === 0" 
              type="success" 
              size="small" 
              @click="handlePublish(row)"
            >
              发布
            </el-button>
            <el-button 
              v-if="row.status === 1" 
              type="warning" 
              size="small" 
              @click="handleMakePrivate(row)"
            >
              设为私有
            </el-button>
            <el-button 
              v-if="row.status === 2" 
              type="success" 
              size="small" 
              @click="handleMakePublic(row)"
            >
              设为公开
            </el-button>
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
        @size-change="fetchArticles"
        @current-change="fetchArticles"
        class="mt-20"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import api from '@/utils/api'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const articles = ref([])
const activeTab = ref('all') // all, published, private, draft

const searchForm = reactive({
  title: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 获取我的文章列表
const fetchArticles = async () => {
  // 确保用户已登录
  if (!userStore.isLoggedIn || !userStore.user) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }

  loading.value = true
  try {
    // 根据 tab 设置 status
    let status = null
    if (activeTab.value === 'published') {
      status = 1
    } else if (activeTab.value === 'private') {
      status = 2
    } else if (activeTab.value === 'draft') {
      status = 0
    }
    // activeTab.value === 'all' 时，status 为 null，显示所有状态

    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      author_id: userStore.user.id, // 只获取当前用户的文章
      ...searchForm
    }
    
    // 如果 status 不为 null，添加到参数中
    if (status !== null) {
      params.status = status
    }
    
    // 清理空值参数
    Object.keys(params).forEach(key => {
      if (params[key] === '' || params[key] === null || params[key] === undefined) {
        delete params[key]
      }
    })

    const response = await api.get('/articles', { params })
    articles.value = response.data.list || []
    pagination.total = response.data.total || 0
  } catch (error) {
    ElMessage.error('获取文章列表失败')
  } finally {
    loading.value = false
  }
}

// Tab 切换
const handleTabChange = () => {
  pagination.page = 1
  fetchArticles()
}

const handleSearch = () => {
  pagination.page = 1
  fetchArticles()
}

const handleReset = () => {
  searchForm.title = ''
  activeTab.value = 'all'
  pagination.page = 1
  fetchArticles()
}

// 获取状态文本
const getStatusText = (status) => {
  switch (status) {
    case 0:
      return '草稿'
    case 1:
      return '已发布'
    case 2:
      return '私有'
    default:
      return '未知'
  }
}

// 获取状态标签类型
const getStatusType = (status) => {
  switch (status) {
    case 0:
      return 'info'
    case 1:
      return 'success'
    case 2:
      return 'warning'
    default:
      return ''
  }
}

// 发布草稿
const handlePublish = async (row) => {
  try {
    await ElMessageBox.confirm('确定要发布这篇文章吗？', '提示', {
      type: 'warning'
    })
    
    // 使用编辑API获取完整的文章详情（确保权限检查）
    const detailResponse = await api.get(`/articles/${row.id}/edit`)
    const article = detailResponse.data
    
    // 更新status字段
    await api.put(`/articles/${row.id}`, {
      title: article.title,
      content: article.content,
      summary: article.summary || '',
      cover: article.cover || '',
      category_id: article.category_id || null,
      tag_ids: article.tags?.map(t => t.id) || [],
      status: 1,
      is_top: article.is_top || false,
      is_recommend: article.is_recommend || false
    })
    
    ElMessage.success('发布成功')
    fetchArticles()
  } catch (error) {
    if (error !== 'cancel') {
      const errorMessage = error.response?.data?.message || error.message || '发布失败'
      ElMessage.error(errorMessage)
    }
  }
}

// 设为私有
const handleMakePrivate = async (row) => {
  try {
    await ElMessageBox.confirm('确定要将这篇文章设为私有吗？设为私有后，只有您可以查看。', '提示', {
      type: 'warning'
    })
    
    // 使用编辑API获取完整的文章详情（确保权限检查）
    const detailResponse = await api.get(`/articles/${row.id}/edit`)
    const article = detailResponse.data
    
    // 更新status字段
    await api.put(`/articles/${row.id}`, {
      title: article.title,
      content: article.content,
      summary: article.summary || '',
      cover: article.cover || '',
      category_id: article.category_id || null,
      tag_ids: article.tags?.map(t => t.id) || [],
      status: 2,
      is_top: article.is_top || false,
      is_recommend: article.is_recommend || false
    })
    
    ElMessage.success('已设为私有')
    fetchArticles()
  } catch (error) {
    if (error !== 'cancel') {
      const errorMessage = error.response?.data?.message || error.message || '操作失败'
      ElMessage.error(errorMessage)
    }
  }
}

// 设为公开
const handleMakePublic = async (row) => {
  try {
    await ElMessageBox.confirm('确定要将这篇文章设为公开吗？设为公开后，所有人都可以查看。', '提示', {
      type: 'warning'
    })
    
    // 使用编辑API获取完整的文章详情（确保权限检查）
    const detailResponse = await api.get(`/articles/${row.id}/edit`)
    const article = detailResponse.data
    
    // 更新status字段
    await api.put(`/articles/${row.id}`, {
      title: article.title,
      content: article.content,
      summary: article.summary || '',
      cover: article.cover || '',
      category_id: article.category_id || null,
      tag_ids: article.tags?.map(t => t.id) || [],
      status: 1,
      is_top: article.is_top || false,
      is_recommend: article.is_recommend || false
    })
    
    ElMessage.success('已设为公开')
    fetchArticles()
  } catch (error) {
    if (error !== 'cancel') {
      const errorMessage = error.response?.data?.message || error.message || '操作失败'
      ElMessage.error(errorMessage)
    }
  }
}

const handleView = (row) => {
  router.push(`/blog/${row.id}`)
}

const handleEdit = (row) => {
  router.push(`/dashboard/articles/${row.id}/edit`)
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除文章《${row.title}》吗？此操作不可恢复！`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'error'
      }
    )

    await api.delete(`/articles/${row.id}`)
    ElMessage.success('删除成功')
    fetchArticles()
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
  fetchArticles()
})
</script>

<style scoped>
.my-articles {
  max-width: 1400px;
}

.my-articles .el-card {
  box-shadow: var(--shadow-md);
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header span {
  font-size: var(--font-size-xl);
  font-weight: 600;
  color: var(--theme-text-primary);
}

.status-tabs {
  margin-bottom: var(--spacing-lg);
}

.search-form {
  margin-bottom: var(--spacing-lg);
}

.stats-text {
  display: flex;
  gap: var(--spacing-sm);
  font-size: var(--font-size-xs);
  color: var(--theme-text-secondary);
}

.mt-20 {
  margin-top: var(--spacing-lg);
}
</style>

