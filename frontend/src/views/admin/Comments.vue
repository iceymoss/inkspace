<template>
  <div class="comments">
    <h2>评论管理</h2>

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
          {{ row.article_id }}
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'warning'">
            {{ row.status === 1 ? '已通过' : '待审核' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button v-if="row.status === 0" size="small" type="success" @click="approve(row)">
            通过
          </el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

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
import api from '@/utils/api'
import dayjs from 'dayjs'

const comments = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD HH:mm')

const loadComments = async () => {
  try {
    const response = await api.get('/comments', {
      params: {
        page: currentPage.value,
        page_size: pageSize.value
      }
    })
    comments.value = response.data.list || []
    total.value = response.data.total || 0
  } catch (error) {
    ElMessage.error('加载失败')
  }
}

const approve = async (comment) => {
  try {
    await api.put(`/admin/comments/${comment.id}/status`, { status: 1 })
    ElMessage.success('已通过')
    loadComments()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleDelete = async (comment) => {
  try {
    await ElMessageBox.confirm('确定要删除这条评论吗？', '提示', { type: 'warning' })
    await api.delete(`/comments/${comment.id}`)
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

