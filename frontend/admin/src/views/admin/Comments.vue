<template>
  <div class="comments">
    <h2>评论管理</h2>

    <el-tabs v-model="activeTab" @tab-change="handleTabChange" style="margin-top: 20px;">
      <el-tab-pane label="文章评论" name="articles">
        <el-table :data="comments" style="width: 100%; margin-top: 20px;">
          <el-table-column prop="id" label="ID" width="80" />
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
          <el-table-column label="状态" width="120">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" size="small">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" width="180">
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
        <el-table :data="comments" style="width: 100%; margin-top: 20px;">
          <el-table-column prop="id" label="ID" width="80" />
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
          <el-table-column label="状态" width="120">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" size="small">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" width="180">
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
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import adminApi from '@/utils/adminApi'
import dayjs from 'dayjs'

const comments = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const activeTab = ref('articles') // 默认显示文章评论

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
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      show_all: true, // 管理后台显示所有状态的评论
      type: activeTab.value === 'articles' ? 'article' : activeTab.value === 'works' ? 'work' : '' // 根据标签页过滤类型
    }
    
    const response = await adminApi.get('/admin/comments', { params })
    comments.value = response.data.list || []
    total.value = response.data.total || 0
  } catch (error) {
    ElMessage.error('加载失败')
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
</script>

<style scoped>
.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>


