<template>
  <div class="article-edit">
    <el-card>
      <template #header>
        <div class="header">
          <span>{{ isEdit ? '编辑文章' : '写文章' }}</span>
          <el-button @click="$router.back()">返回</el-button>
        </div>
      </template>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="文章标题" prop="title">
          <el-input
            v-model="form.title"
            placeholder="请输入文章标题"
            maxlength="100"
            show-word-limit
          />
        </el-form-item>

        <el-form-item label="分类" prop="category_id">
          <el-select v-model="form.category_id" placeholder="请选择分类">
            <el-option
              v-for="cat in categories"
              :key="cat.id"
              :label="cat.name"
              :value="cat.id"
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
            placeholder="请选择或输入新标签"
            @change="handleTagChange"
          >
            <el-option
              v-for="tag in tags"
              :key="tag.id"
              :label="tag.name"
              :value="tag.id"
            />
          </el-select>
          <div class="form-tip">可以输入新标签名称并按回车创建</div>
        </el-form-item>

        <el-form-item label="摘要" prop="summary">
          <el-input
            v-model="form.summary"
            type="textarea"
            :rows="3"
            placeholder="请输入文章摘要（可选）"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>

        <el-form-item label="封面图" prop="cover_image">
          <el-input v-model="form.cover_image" placeholder="请输入封面图URL（可选）" />
        </el-form-item>

        <el-form-item label="内容" prop="content" required>
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="20"
            placeholder="请输入文章内容（支持Markdown）"
          />
          <div class="content-tip">支持Markdown语法</div>
        </el-form-item>

        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio label="draft">草稿</el-radio>
            <el-radio label="published">发布</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="loading">
            {{ isEdit ? '保存' : '发布' }}
          </el-button>
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
import api from '@/utils/api'

const route = useRoute()
const router = useRouter()
const formRef = ref(null)
const loading = ref(false)
const categories = ref([])
const tags = ref([])

const isEdit = computed(() => !!route.params.id)

const form = reactive({
  title: '',
  category_id: null,
  tag_ids: [],
  summary: '',
  cover_image: '',
  content: '',
  status: 'draft'
})

const rules = {
  title: [
    { required: true, message: '请输入文章标题', trigger: 'blur' },
    { min: 1, max: 100, message: '标题长度1-100个字符', trigger: 'blur' }
  ],
  category_id: [
    { required: true, message: '请选择分类', trigger: 'change' }
  ],
  content: [
    { required: true, message: '请输入文章内容', trigger: 'blur' }
  ]
}

// 获取分类列表
const fetchCategories = async () => {
  try {
    const response = await api.get('/categories')
    categories.value = response.data || []
  } catch (error) {
    ElMessage.error('获取分类列表失败')
  }
}

// 获取标签列表
const fetchTags = async () => {
  try {
    const response = await api.get('/tags')
    tags.value = response.data || []
  } catch (error) {
    ElMessage.error('获取标签列表失败')
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
        const response = await api.post('/tags', {
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

// 获取文章详情
const fetchArticle = async () => {
  if (!isEdit.value) return
  
  try {
    const response = await api.get(`/articles/${route.params.id}`)
    const article = response.data
    
    Object.assign(form, {
      title: article.title || '',
      category_id: article.category_id || null,
      tag_ids: article.tags?.map(t => t.id) || [],
      summary: article.summary || '',
      cover_image: article.cover_image || '',
      content: article.content || '',
      status: article.status || 'draft'
    })
  } catch (error) {
    ElMessage.error('获取文章详情失败')
    router.back()
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      if (isEdit.value) {
        await api.put(`/articles/${route.params.id}`, form)
        ElMessage.success('保存成功')
      } else {
        await api.post('/articles', form)
        ElMessage.success('发布成功')
      }
      
      router.push('/dashboard/articles')
    } catch (error) {
      ElMessage.error(error.message || '操作失败')
    } finally {
      loading.value = false
    }
  })
}

onMounted(() => {
  fetchCategories()
  fetchTags()
  fetchArticle()
})
</script>

<style scoped>
.article-edit {
  max-width: 1000px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.content-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}

:deep(.el-textarea__inner) {
  font-family: 'Courier New', Courier, monospace;
}
</style>

