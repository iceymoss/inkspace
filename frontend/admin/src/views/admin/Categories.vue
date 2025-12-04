<template>
  <div class="categories">
    <h2>分类管理</h2>
    <el-button type="primary" @click="showDialog()"><el-icon><Plus /></el-icon> 新建分类</el-button>

    <el-table :data="categories" style="width: 100%; margin-top: 20px;">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="slug" label="别名" />
      <el-table-column prop="description" label="描述" />
      <el-table-column label="操作" width="150">
        <template #default="{ row }">
          <el-button size="small" @click="showDialog(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑分类' : '新建分类'" width="500px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="别名" prop="slug">
          <el-input v-model="form.slug" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" />
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

const categories = ref([])
const dialogVisible = ref(false)
const formRef = ref()
const loading = ref(false)
const editingId = ref(null)

const isEdit = computed(() => !!editingId.value)

const form = reactive({
  name: '',
  slug: '',
  description: ''
})

const rules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }]
}

const loadCategories = async () => {
  try {
    const response = await adminApi.get('/admin/categories')
    categories.value = response.data || []
  } catch (error) {
    ElMessage.error('加载失败')
  }
}

const showDialog = (category = null) => {
  if (category) {
    editingId.value = category.id
    Object.assign(form, category)
  } else {
    editingId.value = null
    Object.assign(form, { name: '', slug: '', description: '' })
  }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      if (isEdit.value) {
        await adminApi.put(`/admin/categories/${editingId.value}`, form)
        ElMessage.success('更新成功')
      } else {
        await adminApi.post('/admin/categories', form)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      loadCategories()
    } catch (error) {
      ElMessage.error('保存失败')
    } finally {
      loading.value = false
    }
  })
}

const handleDelete = async (category) => {
  try {
    await ElMessageBox.confirm('确定要删除这个分类吗？', '提示', { type: 'warning' })
    await adminApi.delete(`/admin/categories/${category.id}`)
    ElMessage.success('删除成功')
    loadCategories()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

onMounted(() => {
  loadCategories()
})
</script>

