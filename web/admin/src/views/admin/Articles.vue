<template>
  <div class="articles">
    <div class="page-header">
      <h2>文章管理</h2>
      <el-button type="primary" @click="$router.push('/articles/create')">
        <el-icon><Plus /></el-icon> 新建文章
      </el-button>
    </div>

    <el-card>
      <!-- 筛选区域 -->
      <div class="filter-bar">
        <el-form inline>
          <el-form-item label="关键字">
            <el-input
              v-model="filters.keyword"
              placeholder="标题 / 内容"
              clearable
              @keyup.enter="handleSearch"
              style="width: 220px"
            />
          </el-form-item>

          <el-form-item label="分类">
            <el-select
              v-model="filters.categoryId"
              placeholder="全部分类"
              clearable
              filterable
              style="width: 180px"
              @change="handleFilterChange"
            >
              <el-option
                v-for="cat in categories"
                :key="cat.id"
                :label="cat.name"
                :value="cat.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="状态">
            <el-select
              v-model="filters.status"
              placeholder="全部状态"
              clearable
              style="width: 140px"
              @change="handleFilterChange"
            >
              <el-option :value="1" label="已发布" />
              <el-option :value="0" label="草稿" />
              <el-option :value="2" label="私有" />
            </el-select>
          </el-form-item>

          <el-form-item label="排序">
            <el-select
              v-model="filters.sortBy"
              placeholder="默认排序"
              clearable
              style="width: 160px"
              @change="handleSortChange"
            >
              <el-option label="默认（置顶 + 最新）" value="" />
              <el-option label="最新发布" value="time" />
              <el-option label="最热排序" value="hot" />
              <el-option label="最多浏览" value="view_count" />
              <el-option label="最多点赞" value="like_count" />
              <el-option label="最多评论" value="comment_count" />
            </el-select>
          </el-form-item>

          <el-form-item>
            <el-button type="primary" @click="handleSearch">查询</el-button>
            <el-button @click="resetFilters">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

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
        <el-table-column label="推荐" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="row.is_recommend ? 'warning' : ''">
              {{ row.is_recommend ? '★ 推荐' : '-' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="view_count" label="浏览" width="100" />
        <el-table-column label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="350" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="$router.push(`/articles/${row.id}`)">查看</el-button>
            <el-button size="small" type="primary" @click="$router.push(`/articles/${row.id}/edit`)">编辑</el-button>
            <el-button 
              size="small" 
              :type="row.is_recommend ? 'warning' : 'default'"
              @click="handleToggleRecommend(row)"
            >
              {{ row.is_recommend ? '取消推荐' : '推荐' }}
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
          @current-change="loadArticles"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import adminApi from '@/utils/adminApi'
import dayjs from 'dayjs'

const articles = ref([])
const categories = ref([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const filters = ref({
  keyword: '',
  categoryId: null,
  status: null,
  sortBy: '',
  sortOrder: 'desc'
})

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD HH:mm')

const loadArticles = async () => {
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }

    if (filters.value.keyword) {
      params.keyword = filters.value.keyword
    }
    if (filters.value.categoryId) {
      params.category_id = filters.value.categoryId
    }
    if (filters.value.status !== null && filters.value.status !== undefined && filters.value.status !== '') {
      params.status = filters.value.status
    }
    if (filters.value.sortBy) {
      params.sort_by = filters.value.sortBy
      params.sort_order = filters.value.sortOrder || 'desc'
    }

    const response = await adminApi.get('/admin/articles', {
      params
    })
    articles.value = response.data.list || []
    total.value = response.data.total || 0
  } catch (error) {
    ElMessage.error('加载失败')
  }
}

const loadCategories = async () => {
  try {
    const res = await adminApi.get('/admin/categories', {
      params: { page: 1, page_size: 100 }
    })
    categories.value = res.data?.list || []
  } catch (error) {
    // 分类加载失败不影响文章列表，静默失败即可
    console.error('加载分类失败', error)
  }
}

const handleSearch = () => {
  currentPage.value = 1
  loadArticles()
}

const handleFilterChange = () => {
  currentPage.value = 1
  loadArticles()
}

const handleSortChange = () => {
  currentPage.value = 1
  loadArticles()
}

const resetFilters = () => {
  filters.value = {
    keyword: '',
    categoryId: null,
    status: null,
    sortBy: '',
    sortOrder: 'desc'
  }
  currentPage.value = 1
  loadArticles()
}

const handleToggleRecommend = async (row) => {
  try {
    await adminApi.put(`/admin/articles/${row.id}/recommend`, {
      is_recommend: !row.is_recommend
    })
    ElMessage.success(row.is_recommend ? '已取消推荐' : '设置推荐成功')
    loadArticles()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除这篇文章吗？', '提示', {
      type: 'warning'
    })
    await adminApi.delete(`/admin/articles/${row.id}`)
    ElMessage.success('删除成功')
    loadArticles()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

onMounted(() => {
  loadCategories()
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

.filter-bar {
  margin-bottom: 16px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>

