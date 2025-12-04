<template>
  <div class="article-edit">
    <div class="page-header">
      <h2>{{ isEdit ? '编辑文章' : '新建文章' }}</h2>
      <el-button @click="$router.back()">返回</el-button>
    </div>

    <el-card>
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入文章标题" />
        </el-form-item>

        <el-form-item label="摘要" prop="summary">
          <el-input v-model="form.summary" type="textarea" :rows="3" placeholder="请输入文章摘要" />
        </el-form-item>

        <el-form-item label="封面" prop="cover">
          <el-input v-model="form.cover" placeholder="请输入封面图片URL" />
        </el-form-item>

        <el-form-item label="分类" prop="category_id">
          <el-select v-model="form.category_id" placeholder="请选择分类">
            <el-option
              v-for="category in categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="标签" prop="tag_ids">
          <el-select v-model="form.tag_ids" multiple placeholder="请选择标签">
            <el-option
              v-for="tag in tags"
              :key="tag.id"
              :label="tag.name"
              :value="tag.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="内容" prop="content">
          <el-input v-model="form.content" type="textarea" :rows="20" placeholder="请输入文章内容（支持Markdown）" />
        </el-form-item>

        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">发布</el-radio>
            <el-radio :label="0">草稿</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="置顶">
          <el-switch v-model="form.is_top" />
        </el-form-item>

        <el-form-item label="推荐">
          <el-switch v-model="form.is_recommend" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="loading">保存</el-button>
          <el-button @click="$router.back()">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import adminApi from '@/utils/adminApi'

const route = useRoute()
const router = useRouter()

const formRef = ref()
const loading = ref(false)
const categories = ref([])
const tags = ref([])

const isEdit = computed(() => !!route.params.id)

const form = reactive({
  title: '',
  summary: '',
  cover: '',
  category_id: null,
  tag_ids: [],
  content: '',
  status: 1,
  is_top: false,
  is_recommend: false
})

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }]
}

const loadCategories = async () => {
  try {
    const response = await adminApi.get('/admin/categories')
    categories.value = response.data || []
  } catch (error) {
    console.error('Failed to load categories:', error)
  }
}

const loadTags = async () => {
  try {
    const response = await adminApi.get('/admin/tags')
    tags.value = response.data || []
  } catch (error) {
    console.error('Failed to load tags:', error)
  }
}

const loadArticle = async () => {
  if (!isEdit.value) return
  try {
    const response = await adminApi.get(`/admin/articles/${route.params.id}`)
    const article = response.data
    Object.assign(form, {
      title: article.title,
      summary: article.summary,
      cover: article.cover,
      category_id: article.category_id,
      tag_ids: article.tags?.map(t => t.id) || [],
      content: article.content,
      status: article.status,
      is_top: article.is_top,
      is_recommend: article.is_recommend
    })
  } catch (error) {
    ElMessage.error('加载文章失败')
  }
}

const handleSubmit = async () => {
  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      if (isEdit.value) {
        await adminApi.put(`/admin/articles/${route.params.id}`, form)
        ElMessage.success('更新成功')
      } else {
        await adminApi.post('/admin/articles', form)
        ElMessage.success('创建成功')
      }
      router.push('/articles')
    } catch (error) {
      ElMessage.error('保存失败')
    } finally {
      loading.value = false
    }
  })
}

onMounted(() => {
  loadCategories()
  loadTags()
  loadArticle()
})
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
</style>

