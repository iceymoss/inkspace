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
          <ImageCropUpload 
            v-model="form.cover" 
            preview-size="160px"
            placeholder="上传封面"
            tip="可自由裁切任意比例，最大5MB"
            :aspect-ratio="NaN"
          />
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
          <el-select 
            v-model="form.tag_ids" 
            multiple 
            filterable 
            allow-create
            default-first-option
            placeholder="选择或创建标签"
            @change="handleTagChange"
          >
            <el-option
              v-for="tag in tags"
              :key="tag.id"
              :label="tag.name"
              :value="tag.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="内容" prop="content" class="content-form-item">
          <VditorEditor v-model="form.content" height="600px" />
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
import VditorEditor from '@/components/VditorEditor.vue'
import ImageCropUpload from '@/components/ImageCropUpload.vue'

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
  category_id: [{ required: true, message: '请选择分类', trigger: 'change' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }]
}

const loadCategories = async () => {
  try {
    const response = await adminApi.get('/admin/categories', {
      params: {
        page: 1,
        page_size: 100
      }
    })
    // 兼容分页结构和旧的数组结构
    categories.value = response.data?.list || response.data || []
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

// 处理标签变化（支持创建新标签）
const handleTagChange = async (values) => {
  // 检查是否有新标签（字符串类型的值）
  const newTags = values.filter(v => typeof v === 'string')
  
  if (newTags.length > 0) {
    for (const tagName of newTags) {
      try {
        // 创建新标签
        const response = await adminApi.post('/admin/tags', {
          name: tagName,
          slug: tagName.toLowerCase().replace(/\s+/g, '-'),
          color: '#409eff'
        })
        
        // 添加到标签列表
        const newTag = response.data
        tags.value.push(newTag)
        
        // 替换form中的字符串为ID
        const index = form.tag_ids.indexOf(tagName)
        if (index > -1) {
          form.tag_ids[index] = newTag.id
        }
        
        ElMessage.success(`标签"${tagName}"创建成功`)
      } catch (error) {
        ElMessage.error(`创建标签"${tagName}"失败`)
        // 移除失败的标签
        const index = form.tag_ids.indexOf(tagName)
        if (index > -1) {
          form.tag_ids.splice(index, 1)
        }
      }
    }
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
.article-edit {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.content-form-item {
  margin-bottom: 24px;
}

.content-form-item :deep(.el-form-item__label) {
  display: none;
}

.content-form-item :deep(.el-form-item__content) {
  margin-left: 0 !important;
}

:deep(.el-form-item__label) {
  font-weight: 600;
}
</style>

