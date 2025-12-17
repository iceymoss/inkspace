<template>
  <div class="tags">
    <div class="header">
      <h2>标签管理</h2>
      <el-button type="primary" @click="showDialog()"><el-icon><Plus /></el-icon> 新建标签</el-button>
    </div>

    <div class="filter-bar">
      <el-form inline>
        <el-form-item label="关键字">
          <el-input
            v-model="filters.keyword"
            placeholder="名称 / 别名"
            clearable
            @keyup.enter="handleSearch"
            style="width: 220px"
          />
        </el-form-item>
        <el-form-item label="标签类型">
          <el-select
            v-model="filters.scope"
            placeholder="全部类型"
            style="width: 160px"
            @change="handleFilterChange"
          >
            <el-option label="全部" value="all" />
            <el-option label="系统标签" value="system" />
            <el-option label="用户标签" value="user" />
          </el-select>
        </el-form-item>
        <el-form-item label="文章数">
          <el-select
            v-model="filters.hasArticles"
            placeholder="全部"
            clearable
            style="width: 140px"
            @change="handleFilterChange"
          >
            <el-option label="全部" :value="''" />
            <el-option label="有文章" value="1" />
            <el-option label="无文章" value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table
      :data="tags"
      v-loading="tableLoading"
      style="width: 100%; margin-top: 10px;"
      @sort-change="handleSortChange"
    >
      <el-table-column prop="id" label="ID" width="80" sortable="custom" />
      <el-table-column prop="name" label="名称" sortable="custom" />
      <el-table-column prop="slug" label="别名" />
      <el-table-column prop="article_count" label="文章数" width="100" sortable="custom" />
      <el-table-column prop="created_at" label="创建时间" width="180" sortable="custom">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="颜色" width="120">
        <template #default="{ row }">
          <el-tag :color="row.color">{{ row.name }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150">
        <template #default="{ row }">
          <el-button size="small" @click="showDialog(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        layout="prev, pager, next, jumper, ->, total"
        @current-change="handlePageChange"
      />
    </div>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑标签' : '新建标签'" width="500px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="别名" prop="slug">
          <el-input v-model="form.slug" />
        </el-form-item>
        <el-form-item label="颜色">
          <div style="display: flex; align-items: center; gap: 10px;">
            <el-color-picker v-model="form.color" />
            <el-input 
              v-model="form.color" 
              placeholder="#409eff" 
              style="width: 200px;"
              @change="handleColorChange"
            />
          </div>
          <div style="margin-top: 8px; color: #909399; font-size: 12px;">
            点击颜色选择器选择颜色，或直接在输入框中输入颜色值（如：#409eff）
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="loading">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import adminApi from '@/utils/adminApi'
import dayjs from 'dayjs'

const tags = ref([])
const dialogVisible = ref(false)
const formRef = ref()
const loading = ref(false)
const tableLoading = ref(false)
const editingId = ref(null)

const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const isEdit = computed(() => !!editingId.value)

const filters = reactive({
  keyword: '',
  scope: 'all', // all/system/user
  hasArticles: '' // '' | '1' | '0'
})

const form = reactive({
  name: '',
  slug: '',
  color: ''
})

const rules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }]
}

const sortField = ref('')
const sortOrder = ref('') // ascending / descending

const formatDate = (date) => {
  if (!date) return ''
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

const buildQueryParams = () => {
  const params = {
    page: currentPage.value,
    page_size: pageSize.value
  }

  if (filters.keyword && filters.keyword.trim()) {
    params.keyword = filters.keyword.trim()
  }

  if (filters.scope && filters.scope !== 'all') {
    params.scope = filters.scope
  }

  if (filters.hasArticles !== '' && filters.hasArticles !== null && filters.hasArticles !== undefined) {
    params.has_articles = filters.hasArticles
  }

  if (sortField.value && sortOrder.value) {
    const dir = sortOrder.value === 'ascending' ? 'asc' : 'desc'
    params.sort = `${sortField.value}_${dir}`
  }

  return params
}

const loadTags = async () => {
  tableLoading.value = true
  try {
    const response = await adminApi.get('/admin/tags', {
      params: buildQueryParams()
    })
    tags.value = response.data?.list || []
    total.value = response.data?.total || 0
  } catch (error) {
    ElMessage.error('加载失败')
  } finally {
    tableLoading.value = false
  }
}

const showDialog = (tag = null) => {
  if (tag) {
    editingId.value = tag.id
    Object.assign(form, tag)
    // 确保颜色值存在，如果没有则设置默认值
    if (!form.color) {
      form.color = '#409eff'
    }
  } else {
    editingId.value = null
    Object.assign(form, { name: '', slug: '', color: '#409eff' })
  }
  dialogVisible.value = true
}

const handleColorChange = (value) => {
  // 确保颜色值格式正确（以#开头）
  if (value && !value.startsWith('#')) {
    form.color = '#' + value
  }
}

const handleSubmit = async () => {
  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      if (isEdit.value) {
        await adminApi.put(`/admin/tags/${editingId.value}`, form)
        ElMessage.success('更新成功')
      } else {
        await adminApi.post('/admin/tags', form)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      loadTags()
    } catch (error) {
      ElMessage.error('保存失败')
    } finally {
      loading.value = false
    }
  })
}

const handleDelete = async (tag) => {
  try {
    await ElMessageBox.confirm('确定要删除这个标签吗？', '提示', { type: 'warning' })
    await adminApi.delete(`/admin/tags/${tag.id}`)
    ElMessage.success('删除成功')
    loadTags()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

onMounted(() => {
  loadTags()
})

// 搜索与筛选
const handleSearch = () => {
  currentPage.value = 1
  loadTags()
}

const handleFilterChange = () => {
  currentPage.value = 1
  loadTags()
}

const resetFilters = () => {
  filters.keyword = ''
  filters.scope = 'all'
  filters.hasArticles = ''
  sortField.value = ''
  sortOrder.value = ''
  currentPage.value = 1
  loadTags()
}

// 排序
const handleSortChange = ({ prop, order }) => {
  sortField.value = order ? prop : ''
  sortOrder.value = order || ''
  currentPage.value = 1
  loadTags()
}

const handlePageChange = (page) => {
  currentPage.value = page
  loadTags()
}
</script>

<style scoped>
.tags {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.filter-bar {
  margin-bottom: 10px;
}

.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>

