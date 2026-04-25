<template>
  <div class="articles">
    <div class="page-header">
      <h2>文章管理</h2>
      <el-button type="primary" @click="$router.push('/admin/articles/create')">
        <el-icon><Plus /></el-icon> 新建文章
      </el-button>
    </div>

    <el-card>
      <el-table :data="articles" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="标题" min-width="200" />
        <el-table-column label="分类" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.category">{{ row.category.name }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '已发布' : '草稿' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="view_count" label="浏览" width="100" />
        <el-table-column label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="$router.push(`/admin/articles/${row.id}/edit`)">编辑</el-button>
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
          @current-change="loadArticles"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { toast } from 'vue-sonner'
import { Plus } from 'lucide-vue-next'
import adminApi from '@/utils/adminApi'
import dayjs from 'dayjs'

const articles = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD HH:mm')

const loadArticles = async () => {
  try {
    const response = await adminApi.get('/articles', {
      params: {
        page: currentPage.value,
        page_size: pageSize.value
      }
    })
    articles.value = response.data.list || []
    total.value = response.data.total || 0
  } catch (error) {
    toast.error('加载失败')
  }
}

const handleDelete = async (row) => {
  try {
    if (!confirm('确定要删除这篇文章吗？')) return
    await adminApi.delete(`/articles/${row.id}`)
    toast.success('删除成功')
    loadArticles()
  } catch (error) {
    toast.error('删除失败')
  }
}

onMounted(() => {
  loadArticles()
})
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>

