<template>
  <div class="article-edit">
    <div class="edit-container">
      <div class="edit-header">
        <h1>{{ isEdit ? '编辑文章' : '写文章' }}</h1>
        <el-button @click="$router.back()" plain>返回</el-button>
      </div>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="0"
        class="article-form"
      >
        <el-form-item prop="title">
          <el-input
            v-model="form.title"
            placeholder="请输入文章标题"
            maxlength="100"
            show-word-limit
            size="large"
            class="title-input"
          />
        </el-form-item>

        <div class="form-row">
          <el-form-item prop="category_id" class="form-item-inline">
            <template #label>
              <span class="form-label">分类</span>
            </template>
            <el-select v-model="form.category_id" placeholder="请选择分类" style="width: 200px">
              <el-option
                v-for="cat in categories"
                :key="cat.id"
                :label="cat.name"
                :value="cat.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item prop="tag_ids" class="form-item-inline">
            <template #label>
              <span class="form-label">标签</span>
            </template>
            <el-select 
              v-model="form.tag_ids" 
              multiple 
              filterable 
              allow-create
              default-first-option
              placeholder="选择或创建标签"
              @change="handleTagChange"
              style="width: 300px"
            >
              <el-option
                v-for="tag in tags"
                :key="tag.id"
                :label="tag.name"
                :value="tag.id"
              />
            </el-select>
          </el-form-item>
        </div>

        <el-form-item prop="summary" class="form-item-block">
          <template #label>
            <span class="form-label-block">摘要</span>
          </template>
          <el-input
            v-model="form.summary"
            type="textarea"
            :rows="4"
            placeholder="文章摘要，将显示在列表页..."
            maxlength="500"
            show-word-limit
            class="summary-input"
          />
        </el-form-item>

        <el-form-item prop="cover" class="form-item-block">
          <template #label>
            <span class="form-label-block">封面图</span>
          </template>
          <ImageCropUpload 
            v-model="form.cover" 
            preview-size="160px"
            placeholder="上传封面"
            tip="可自由裁切任意比例，最大5MB"
            :aspect-ratio="NaN"
          />
        </el-form-item>

        <el-form-item prop="content" class="content-form-item" required>
          <VditorEditor v-model="form.content" height="600px" />
        </el-form-item>


        <div class="form-actions">
          <div class="action-left">
            <el-radio-group v-model="form.status" size="default">
              <el-radio-button :label="0">保存草稿</el-radio-button>
              <el-radio-button :label="1">立即发布</el-radio-button>
            </el-radio-group>
          </div>
          <div class="action-right">
            <el-button @click="$router.back()" size="large">取消</el-button>
            <el-button type="primary" @click="handleSubmit" :loading="loading" size="large">
              {{ isEdit ? '保存修改' : (form.status === 1 ? '发布文章' : '保存草稿') }}
            </el-button>
          </div>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import api from '@/utils/api'
import VditorEditor from '@/components/VditorEditor.vue'
import ImageCropUpload from '@/components/ImageCropUpload.vue'

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
  cover: '',
  content: '',
  status: 0  // 0: draft, 1: published
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
      cover: article.cover || '',
      content: article.content || '',
      status: article.status || 0  // 0: draft, 1: published
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
      let articleId = route.params.id
      
      // 准备提交数据
      const submitData = {
        title: form.title,
        content: form.content,
        summary: form.summary,
        cover: form.cover,
        category_id: form.category_id,
        tag_ids: form.tag_ids,
        status: form.status,
        is_top: form.is_top || false,
        is_recommend: form.is_recommend || false
      }
      
      // 打印提交的数据，方便调试
      console.log('Submitting article data:', submitData)
      
      if (isEdit.value) {
        await api.put(`/articles/${articleId}`, submitData)
        ElMessage.success('保存成功')
      } else {
        const response = await api.post('/articles', submitData)
        articleId = response.data.id
        ElMessage.success('发布成功')
      }
      
      // 跳转到文章详情页
      router.push(`/blog/${articleId}`)
    } catch (error) {
      console.error('Submit error:', error)
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
  background: #f5f7fa;
  min-height: calc(100vh - 60px);
  padding: 20px 0;
}

.edit-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 20px;
}

.edit-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.edit-header h1 {
  font-size: 24px;
  font-weight: 600;
  color: #1f2329;
  margin: 0;
}

.article-form {
  background: #fff;
  border-radius: 8px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  width: 100%;
}

.form-label {
  font-size: 14px;
  font-weight: 600;
  color: #495057;
  margin-right: 12px;
  min-width: 50px;
  display: inline-block;
  line-height: 32px;
}

.title-input {
  margin-bottom: 16px;
}

.title-input :deep(.el-input__inner) {
  font-size: 18px;
  font-weight: 500;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  padding: 12px 16px;
}

.title-input :deep(.el-input__inner:focus) {
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.1);
}

.form-row {
  display: flex;
  gap: 24px;
  margin-bottom: 18px;
  flex-wrap: wrap;
  align-items: center;
}

.form-item-inline {
  margin-bottom: 0 !important;
  display: flex;
  align-items: center;
  flex: 0 0 auto;
}

.form-item-inline :deep(.el-form-item__label) {
  padding: 0;
  min-width: 50px;
  margin-right: 12px;
  line-height: 32px;
}

.form-item-inline :deep(.el-form-item__content) {
  margin-left: 0 !important;
}

.form-item-block {
  margin-bottom: 18px;
}

.form-label-block {
  font-size: 14px;
  font-weight: 600;
  color: #495057;
  display: block;
  margin-bottom: 8px;
  line-height: 1.5;
}

.summary-input :deep(.el-textarea__inner) {
  border-radius: 6px;
  border: 1px solid #e4e7ed;
  font-family: inherit;
}

.summary-input :deep(.el-textarea__inner:focus) {
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.1);
}

.cover-input :deep(.el-input__inner) {
  border-radius: 6px;
  border: 1px solid #e4e7ed;
}

.cover-input :deep(.el-input__inner:focus) {
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.1);
}

.form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid #e4e7ed;
}

.action-left {
  flex: 1;
}

.action-right {
  display: flex;
  gap: 12px;
}

:deep(.el-select) {
  border-radius: 6px;
}

:deep(.el-select .el-input__inner) {
  border-radius: 6px;
  border: 1px solid #e4e7ed;
}

:deep(.el-select .el-input__inner:focus) {
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.1);
}

:deep(.el-form-item__label) {
  font-size: 14px;
  font-weight: 600;
  color: #495057;
}

.content-form-item {
  margin-bottom: 24px;
}

:deep(.content-form-item .el-form-item__content) {
  line-height: normal;
  width: 100%;
  max-width: 100%;
}

:deep(.content-form-item .vditor-container) {
  width: 100%;
}

:deep(.vditor) {
  width: 100% !important;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
}
</style>

