<template>
  <div class="comments">
    <h2>评论管理</h2>

    <!-- 筛选区域 -->
    <div class="filter-bar">
      <el-form inline>
        <el-form-item label="关键字">
          <el-input
            v-model="filters.keyword"
            placeholder="内容 / 昵称 / 邮箱"
            clearable
            @keyup.enter="handleSearch"
            style="width: 260px"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="filters.status"
            placeholder="全部状态"
            clearable
            style="width: 140px"
            @change="handleFilterChange"
          >
            <el-option label="全部" :value="''" />
            <el-option label="待审核" value="0" />
            <el-option label="已通过" value="1" />
            <el-option label="已拒绝" value="-1" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-tabs v-model="activeTab" @tab-change="handleTabChange" style="margin-top: 10px;">
      <el-tab-pane label="文章评论" name="articles">
        <el-table
          :data="comments"
          style="width: 100%; margin-top: 20px;"
          v-loading="tableLoading"
          @sort-change="handleSortChange"
        >
          <el-table-column prop="id" label="ID" width="80" sortable="custom" />
          <el-table-column label="内容" min-width="200">
            <template #default="{ row }">
              {{ row.content }}
            </template>
          </el-table-column>
          <el-table-column label="评论者" width="150">
            <template #default="{ row }">
              {{ row.user?.nickname || row.nickname }}
            </template>
          </el-table-column>
          <el-table-column label="文章ID" width="100">
            <template #default="{ row }">
              {{ row.article_id || '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="like_count" label="点赞" width="90" sortable="custom" />
          <el-table-column prop="reply_count" label="回复" width="90" sortable="custom" />
          <el-table-column label="状态" width="120">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" size="small">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="创建时间" width="180" sortable="custom">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="280" fixed="right">
            <template #default="{ row }">
              <el-button v-if="row.status === 0" size="small" type="success" @click="approve(row)">
                通过
              </el-button>
              <el-button v-if="row.status === 0" size="small" type="danger" @click="reject(row)">
                拒绝
              </el-button>
              <el-button v-if="row.status === 1" size="small" type="warning" @click="reject(row)">
                拒绝
              </el-button>
              <el-button v-if="row.status === -1" size="small" type="success" @click="approve(row)">
                通过
              </el-button>
              <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="作品评论" name="works">
        <el-table
          :data="comments"
          style="width: 100%; margin-top: 20px;"
          v-loading="tableLoading"
          @sort-change="handleSortChange"
        >
          <el-table-column prop="id" label="ID" width="80" sortable="custom" />
          <el-table-column label="内容" min-width="200">
            <template #default="{ row }">
              {{ row.content }}
            </template>
          </el-table-column>
          <el-table-column label="评论者" width="150">
            <template #default="{ row }">
              {{ row.user?.nickname || row.nickname }}
            </template>
          </el-table-column>
          <el-table-column label="作品ID" width="100">
            <template #default="{ row }">
              {{ row.work_id || '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="like_count" label="点赞" width="90" sortable="custom" />
          <el-table-column prop="reply_count" label="回复" width="90" sortable="custom" />
          <el-table-column label="状态" width="120">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" size="small">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="创建时间" width="180" sortable="custom">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="280" fixed="right">
            <template #default="{ row }">
              <el-button v-if="row.status === 0" size="small" type="success" @click="approve(row)">
                通过
              </el-button>
              <el-button v-if="row.status === 0" size="small" type="danger" @click="reject(row)">
                拒绝
              </el-button>
              <el-button v-if="row.status === 1" size="small" type="warning" @click="reject(row)">
                拒绝
              </el-button>
              <el-button v-if="row.status === -1" size="small" type="success" @click="approve(row)">
                通过
              </el-button>
              <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="loadComments"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import adminApi from '@/utils/adminApi'
import dayjs from 'dayjs'

const comments = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const activeTab = ref('articles') // 默认显示文章评论
const tableLoading = ref(false)

const filters = reactive({
  keyword: '',
  status: '' // '', '0', '1', '-1'
})

const sortField = ref('')
const sortOrder = ref('') // ascending / descending

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD HH:mm')

// 获取状态文本
const getStatusText = (status) => {
  switch (status) {
    case 1:
      return '已通过'
    case -1:
      return '已拒绝'
    case 0:
    default:
      return '待审核'
  }
}

// 获取状态标签类型
const getStatusType = (status) => {
  switch (status) {
    case 1:
      return 'success'
    case -1:
      return 'danger'
    case 0:
    default:
      return 'warning'
  }
}

const loadComments = async () => {
  try {
    tableLoading.value = true
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      show_all: true, // 管理后台显示所有状态的评论
      type: activeTab.value === 'articles' ? 'article' : activeTab.value === 'works' ? 'work' : '' // 根据标签页过滤类型
    }

    if (filters.keyword && filters.keyword.trim()) {
      params.keyword = filters.keyword.trim()
    }

    if (filters.status !== '' && filters.status !== null && filters.status !== undefined) {
      params.status = Number(filters.status)
    }

    if (sortField.value && sortOrder.value) {
      const dir = sortOrder.value === 'ascending' ? 'asc' : 'desc'
      params.sort = `${sortField.value}_${dir}`
    }

    const response = await adminApi.get('/admin/comments', { params })
    comments.value = response.data.list || []
    total.value = response.data.total || 0
  } catch (error) {
    ElMessage.error('加载失败')
  } finally {
    tableLoading.value = false
  }
}

const handleTabChange = (tabName) => {
  // 切换标签页时，重置页码并重新加载
  currentPage.value = 1
  loadComments()
}

const approve = async (comment) => {
  try {
    await adminApi.put(`/admin/comments/${comment.id}/status`, { status: 1 })
    ElMessage.success('审核通过')
    loadComments()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const reject = async (comment) => {
  try {
    await ElMessageBox.confirm('确定要拒绝这条评论吗？', '提示', { type: 'warning' })
    await adminApi.put(`/admin/comments/${comment.id}/status`, { status: -1 })
    ElMessage.success('已拒绝')
    loadComments()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

const handleDelete = async (comment) => {
  try {
    await ElMessageBox.confirm('确定要删除这条评论吗？', '提示', { type: 'warning' })
    await adminApi.delete(`/admin/comments/${comment.id}`)
    ElMessage.success('删除成功')
    loadComments()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

onMounted(() => {
  loadComments()
})

// 搜索与筛选
const handleSearch = () => {
  currentPage.value = 1
  loadComments()
}

const handleFilterChange = () => {
  currentPage.value = 1
  loadComments()
}

const resetFilters = () => {
  filters.keyword = ''
  filters.status = ''
  sortField.value = ''
  sortOrder.value = ''
  currentPage.value = 1
  loadComments()
}

// 排序
const handleSortChange = ({ prop, order }) => {
  sortField.value = order ? prop : ''
  sortOrder.value = order || ''
  currentPage.value = 1
  loadComments()
}
</script>

<style scoped>
.comments {
  padding: 20px;
}

.filter-bar {
  margin-top: 10px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>


