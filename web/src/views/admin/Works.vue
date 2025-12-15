<template>
  <div class="works">
    <h2>作品管理</h2>
    <el-button type="primary" @click="showDialog()"><el-icon><Plus /></el-icon> 新建作品</el-button>

    <el-table :data="works" style="width: 100%; margin-top: 20px;">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="title" label="标题" />
      <el-table-column label="封面" width="120">
        <template #default="{ row }">
          <el-image :src="row.cover" style="width: 80px; height: 60px;" fit="cover" />
        </template>
      </el-table-column>
      <el-table-column prop="view_count" label="浏览" width="100" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'">
            {{ row.status === 1 ? '已发布' : '草稿' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="showDialog(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑作品' : '新建作品'" width="600px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="封面" prop="cover">
          <el-input v-model="form.cover" />
        </el-form-item>
        <el-form-item label="链接" prop="link">
          <el-input v-model="form.link" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">发布</el-radio>
            <el-radio :label="0">草稿</el-radio>
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
import adminApi from '@/utils/adminApi'

const works = ref([])
const dialogVisible = ref(false)
const formRef = ref()
const loading = ref(false)
const editingId = ref(null)

const isEdit = computed(() => !!editingId.value)

const form = reactive({
  title: '',
  description: '',
  cover: '',
  link: '',
  status: 1,
  images: []
})

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }]
}

const loadWorks = async () => {
  try {
    const response = await adminApi.get('/admin/works', {
      params: { page: 1, page_size: 100 }
    })
    works.value = response.data.list || []
  } catch (error) {
    ElMessage.error('加载失败')
  }
}

const showDialog = (work = null) => {
  if (work) {
    editingId.value = work.id
    Object.assign(form, work)
  } else {
    editingId.value = null
    Object.assign(form, {
      title: '',
      description: '',
      cover: '',
      link: '',
      status: 1,
      images: []
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
        await adminApi.put(`/admin/works/${editingId.value}`, form)
        ElMessage.success('更新成功')
      } else {
        await adminApi.post('/admin/works', form)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      loadWorks()
    } catch (error) {
      ElMessage.error('保存失败')
    } finally {
      loading.value = false
    }
  })
}

const handleDelete = async (work) => {
  try {
    await ElMessageBox.confirm('确定要删除这个作品吗？', '提示', { type: 'warning' })
    await adminApi.delete(`/admin/works/${work.id}`)
    ElMessage.success('删除成功')
    loadWorks()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

onMounted(() => {
  loadWorks()
})
</script>

