<template>
  <div class="my-works">
    <el-card>
      <template #header>
        <div class="header">
          <span>我的作品</span>
          <el-button type="primary" @click="$router.push('/dashboard/works/create')">
            <el-icon><Plus /></el-icon> 创建作品
          </el-button>
        </div>
      </template>

      <!-- 筛选栏 -->
      <el-form :inline="true" class="search-form">
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="全部类型" clearable style="width: 150px">
            <el-option label="开源项目" value="project" />
            <el-option label="摄影作品" value="photography" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部状态" clearable style="width: 120px">
            <el-option label="草稿" :value="0" />
            <el-option label="已发布" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
        <el-form-item>
          <el-text type="info" size="small">
            今日摄影作品配额：{{ quotaUsed }}/3
          </el-text>
        </el-form-item>
      </el-form>

      <!-- 作品列表 -->
      <el-table :data="works" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="封面" width="120">
          <template #default="{ row }">
            <el-image :src="row.cover" style="width: 80px; height: 60px;" fit="cover" />
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" min-width="200" />
        <el-table-column label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === 'photography' ? 'warning' : 'primary'" size="small">
              {{ row.type === 'photography' ? '📷 摄影' : '💻 项目' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="200">
          <template #default="{ row }">
            <div>
              <el-tag :type="getStatusType(row.status)" size="small">
                {{ getStatusText(row.status) }}
              </el-tag>
              <div v-if="row.audit_message" style="margin-top: var(--spacing-xs);">
                <el-text type="info" size="small" style="font-size: var(--font-size-xs);">
                  {{ row.audit_message }}
                </el-text>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="统计" width="120">
          <template #default="{ row }">
            <div class="stats-text">
              <span>👁 {{ row.view_count || 0 }}</span>
              <span>💬 {{ row.comment_count || 0 }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleView(row)">查看</el-button>
            <el-button size="small" type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination" v-if="total > 0">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="loadWorks"
        />
      </div>

      <el-empty v-if="works.length === 0 && !loading" description="还没有发布作品">
        <el-button type="primary" @click="$router.push('/dashboard/works/create')">
          创建第一个作品
        </el-button>
      </el-empty>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import api from '@/utils/api'
import dayjs from 'dayjs'
import { navigateToWorkDetail } from '@/utils/workNavigation'

const router = useRouter()

const works = ref([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const quotaUsed = ref(0)

const searchForm = reactive({
  type: '',
  status: null
})

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD HH:mm')

// 加载配额信息
const loadQuota = async () => {
  try {
    const response = await api.get('/works/quota')
    quotaUsed.value = response.data.used || 0
  } catch (error) {
    console.error('Failed to load quota:', error)
  }
}

// 加载作品列表
const loadWorks = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    
    if (searchForm.type) params.type = searchForm.type
    if (searchForm.status !== null && searchForm.status !== '') params.status = searchForm.status

    const response = await api.get('/works/my', { params })
    works.value = response.data.list || []
    total.value = response.data.total || 0
  } catch (error) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  loadWorks()
}

const handleReset = () => {
  searchForm.type = ''
  searchForm.status = null
  currentPage.value = 1
  loadWorks()
}

const handleView = (row) => {
  // 使用预加载导航
  navigateToWorkDetail(row.id, router)
}

const handleEdit = (row) => {
  router.push(`/dashboard/works/${row.id}/edit`)
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除这个作品吗？', '提示', {
      type: 'warning'
    })
    await api.delete(`/works/${row.id}`)
    ElMessage.success('删除成功')
    loadWorks()
    loadQuota()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 获取状态文本
const getStatusText = (status) => {
  switch (status) {
    case 1:
      return '已发布'
    case 0:
      return '草稿'
    case 2:
      return '审核中'
    case 3:
      return '审核不通过'
    default:
      return '未知'
  }
}

// 获取状态标签类型
const getStatusType = (status) => {
  switch (status) {
    case 1:
      return 'success'
    case 0:
      return 'info'
    case 2:
      return 'warning'
    case 3:
      return 'danger'
    default:
      return 'info'
  }
}

onMounted(() => {
  loadWorks()
  loadQuota()
})
</script>

<style scoped>
.my-works {
  max-width: 1400px;
}

.my-works .el-card {
  box-shadow: var(--shadow-md);
  border-radius: var(--radius-lg);
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: var(--font-size-lg);
  font-weight: 600;
  color: var(--theme-text-primary);
}

.search-form {
  margin-bottom: var(--spacing-lg);
}

.stats-text {
  display: flex;
  gap: var(--spacing-sm);
  font-size: var(--font-size-base);
  color: var(--theme-text-secondary);
}

.pagination {
  margin-top: var(--spacing-lg);
  display: flex;
  justify-content: center;
}
</style>

