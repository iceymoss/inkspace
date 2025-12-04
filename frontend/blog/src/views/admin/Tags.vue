<template>
  <div class="tags">
    <h2>标签管理</h2>
    <el-button type="primary" @click="showDialog()"><el-icon><Plus /></el-icon> 新建标签</el-button>

    <el-table :data="tags" style="width: 100%; margin-top: 20px;">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="slug" label="别名" />
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

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑标签' : '新建标签'" width="500px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="别名" prop="slug">
          <el-input v-model="form.slug" />
        </el-form-item>
        <el-form-item label="颜色">
          <el-input v-model="form.color" placeholder="#409eff" />
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

const tags = ref([])
const dialogVisible = ref(false)
const formRef = ref()
const loading = ref(false)
const editingId = ref(null)

const isEdit = computed(() => !!editingId.value)

const form = reactive({
  name: '',
  slug: '',
  color: ''
})

const rules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }]
}

const loadTags = async () => {
  try {
    const response = await adminApi.get('/tags')
    tags.value = response.data || []
  } catch (error) {
    ElMessage.error('加载失败')
  }
}

const showDialog = (tag = null) => {
  if (tag) {
    editingId.value = tag.id
    Object.assign(form, tag)
  } else {
    editingId.value = null
    Object.assign(form, { name: '', slug: '', color: '' })
  }
  dialogVisible.value = true
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
</script>

