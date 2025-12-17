<template>
  <div class="links">
    <div class="header">
      <h2>友情链接管理</h2>
      <el-button type="primary" @click="showDialog()"><el-icon><Plus /></el-icon> 新建友链</el-button>
    </div>

    <!-- 筛选区域 -->
    <div class="filter-bar">
      <el-form inline>
        <el-form-item label="关键字">
          <el-input
            v-model="filters.keyword"
            placeholder="名称 / 链接 / 描述"
            clearable
            @keyup.enter="handleSearch"
            style="width: 260px"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="filters.status"
            size="small"
            style="width: 140px"
            @change="handleFilterChange"
          >
            <el-option label="全部" value="all" />
            <el-option label="显示" value="1" />
            <el-option label="隐藏" value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table
      :data="links"
      v-loading="tableLoading"
      style="width: 100%; margin-top: 10px;"
      @sort-change="handleSortChange"
    >
      <el-table-column prop="id" label="ID" width="80" sortable="custom" />
      <el-table-column prop="name" label="名称" sortable="custom" />
      <el-table-column label="Logo" width="100">
        <template #default="{ row }">
          <el-avatar v-if="row.logo" :src="row.logo" :size="40" />
        </template>
      </el-table-column>
      <el-table-column prop="url" label="链接" />
      <el-table-column prop="description" label="描述" show-overflow-tooltip />
      <el-table-column prop="sort" label="排序" width="80" sortable="custom" />
      <el-table-column prop="status" label="状态" width="100" sortable="custom">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'">
            {{ row.status === 1 ? '显示' : '隐藏' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180" sortable="custom">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
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

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑友链' : '新建友链'" width="600px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="链接" prop="url">
          <el-input v-model="form.url" placeholder="https://example.com" />
          <div style="font-size: 12px; color: #909399; margin-top: 4px;">
            请输入完整的URL，例如：https://example.com
          </div>
        </el-form-item>
        <el-form-item label="Logo">
          <div style="display: flex; gap: 10px; align-items: flex-start;">
            <ImageCropUpload
              v-model="form.logo"
              :aspect-ratio="1"
              :output-width="200"
              :output-height="200"
              preview-size="100px"
              placeholder="点击上传Logo"
              tip="上传正方形Logo，系统会自动裁剪为200x200"
              :max-size="2"
            />
            <div style="flex: 1;">
              <el-input 
                v-model="form.logo" 
                placeholder="或直接输入Logo URL" 
                style="margin-bottom: 8px;"
              />
              <div style="font-size: 12px; color: #909399;">
                支持上传图片或直接输入图片URL。上传的图片将按照正方形比例自动裁剪为200x200
              </div>
            </div>
          </div>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="form.email" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">显示</el-radio>
            <el-radio :label="0">隐藏</el-radio>
          </el-radio-group>
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
import dayjs from 'dayjs'
import adminApi from '@/utils/adminApi'
import ImageCropUpload from '@/components/ImageCropUpload.vue'

const links = ref([])
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
  status: 'all' // 'all' | '1' | '0'
})

const sortField = ref('')
const sortOrder = ref('') // ascending / descending

const form = reactive({
  name: '',
  url: '',
  logo: '',
  description: '',
  email: '',
  sort: 0,
  status: 1
})

const rules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  url: [
    { required: true, message: '请输入链接', trigger: 'blur' },
    {
      type: 'url',
      message: '请输入有效的URL格式（例如：https://example.com）',
      trigger: 'blur'
    }
  ]
}

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

  if (filters.status && filters.status !== 'all') {
    params.status = filters.status
  }

  if (sortField.value && sortOrder.value) {
    const dir = sortOrder.value === 'ascending' ? 'asc' : 'desc'
    params.sort = `${sortField.value}_${dir}`
  }

  return params
}

const loadLinks = async () => {
  tableLoading.value = true
  try {
    const response = await adminApi.get('/admin/links', {
      params: buildQueryParams()
    })
    links.value = response.data?.list || []
    total.value = response.data?.total || 0
  } catch (error) {
    ElMessage.error('加载失败')
  } finally {
    tableLoading.value = false
  }
}

const showDialog = (link = null) => {
  if (link) {
    editingId.value = link.id
    Object.assign(form, link)
  } else {
    editingId.value = null
    Object.assign(form, {
      name: '',
      url: '',
      logo: '',
      description: '',
      email: '',
      sort: 0,
      status: 1
    })
  }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      if (isEdit.value) {
        await adminApi.put(`/admin/links/${editingId.value}`, form)
        ElMessage.success('更新成功')
      } else {
        await adminApi.post('/admin/links', form)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      loadLinks()
    } catch (error) {
      ElMessage.error('保存失败')
    } finally {
      loading.value = false
    }
  })
}

const handleDelete = async (link) => {
  try {
    await ElMessageBox.confirm('确定要删除这个友链吗？', '提示', { type: 'warning' })
    await adminApi.delete(`/admin/links/${link.id}`)
    ElMessage.success('删除成功')
    loadLinks()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

onMounted(() => {
  loadLinks()
})

// 搜索与筛选
const handleSearch = () => {
  currentPage.value = 1
  loadLinks()
}

const handleFilterChange = () => {
  currentPage.value = 1
  loadLinks()
}

const resetFilters = () => {
  filters.keyword = ''
  filters.status = 'all'
  sortField.value = ''
  sortOrder.value = ''
  currentPage.value = 1
  loadLinks()
}

// 排序
const handleSortChange = ({ prop, order }) => {
  sortField.value = order ? prop : ''
  sortOrder.value = order || ''
  loadLinks()
}

const handlePageChange = (page) => {
  currentPage.value = page
  loadLinks()
}
</script>

