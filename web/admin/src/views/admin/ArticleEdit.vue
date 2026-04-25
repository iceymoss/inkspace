<template>
  <div class="article-edit">
    <div class="page-header">
      <h2>{{ isEdit ? '编辑文章' : '新建文章' }}</h2>
      <Button variant="outline" @click="$router.back()">返回</Button>
    </div>

    <Card>
      <CardContent class="pt-6">
        <form class="space-y-4">
          <div class="space-y-2">
            <label class="text-sm font-medium">标题 <span class="text-destructive">*</span></label>
            <Input v-model="form.title" placeholder="请输入文章标题" />
            <p v-if="formErrors.title" class="text-sm text-destructive">{{ formErrors.title }}</p>
          </div>

          <div class="space-y-2">
            <label class="text-sm font-medium">摘要</label>
            <Textarea v-model="form.summary" :rows="3" placeholder="请输入文章摘要" />
          </div>

          <div class="space-y-2">
            <label class="text-sm font-medium">封面</label>
            <ImageCropUpload
              v-model="form.cover"
              preview-size="160px"
              placeholder="上传封面"
              tip="可自由裁切任意比例，最大5MB"
              :aspect-ratio="NaN"
            />
          </div>

          <div class="space-y-2">
            <label class="text-sm font-medium">分类 <span class="text-destructive">*</span></label>
            <Select :model-value="form.category_id ? String(form.category_id) : undefined" @update:model-value="form.category_id = $event ? Number($event) : null">
              <SelectTrigger>
                <SelectValue placeholder="请选择分类" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="category in categories" :key="category.id" :value="String(category.id)">
                  {{ category.name }}
                </SelectItem>
              </SelectContent>
            </Select>
            <p v-if="formErrors.category_id" class="text-sm text-destructive">{{ formErrors.category_id }}</p>
          </div>

          <div class="space-y-2">
            <label class="text-sm font-medium">标签</label>
            <div class="flex flex-wrap gap-2 mb-2" v-if="form.tag_ids.length">
              <Badge v-for="tagId in form.tag_ids" :key="tagId" variant="secondary" class="gap-1">
                {{ getTagName(tagId) }}
                <X class="h-3 w-3 cursor-pointer" @click="removeTag(tagId)" />
              </Badge>
            </div>
            <div class="flex gap-2">
              <Select @update:model-value="addTag($event)">
                <SelectTrigger class="w-[200px]">
                  <SelectValue placeholder="选择标签" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem v-for="tag in availableTags" :key="tag.id" :value="String(tag.id)">
                    {{ tag.name }}
                  </SelectItem>
                </SelectContent>
              </Select>
              <div class="flex gap-2">
                <Input v-model="newTagName" placeholder="创建新标签" class="w-[160px]" @keyup.enter="createNewTag" />
                <Button type="button" variant="outline" size="sm" @click="createNewTag" :disabled="!newTagName.trim()">添加</Button>
              </div>
            </div>
          </div>

          <div class="space-y-2 content-form-item">
            <label class="text-sm font-medium">内容 <span class="text-destructive">*</span></label>
            <VditorEditor v-model="form.content" height="600px" />
            <p v-if="formErrors.content" class="text-sm text-destructive">{{ formErrors.content }}</p>
          </div>

          <div class="space-y-2">
            <label class="text-sm font-medium">状态</label>
            <RadioGroup :model-value="String(form.status)" @update:model-value="form.status = Number($event)" class="flex gap-4">
              <div class="flex items-center space-x-2">
                <RadioGroupItem value="1" id="status-publish" />
                <label for="status-publish">发布</label>
              </div>
              <div class="flex items-center space-x-2">
                <RadioGroupItem value="0" id="status-draft" />
                <label for="status-draft">草稿</label>
              </div>
            </RadioGroup>
          </div>

          <div class="flex items-center gap-6">
            <div class="flex items-center space-x-2">
              <Switch :checked="form.is_top" @update:checked="form.is_top = $event" />
              <label class="text-sm font-medium">置顶</label>
            </div>
            <div class="flex items-center space-x-2">
              <Switch :checked="form.is_recommend" @update:checked="form.is_recommend = $event" />
              <label class="text-sm font-medium">推荐</label>
            </div>
          </div>

          <div class="flex gap-2 pt-4">
            <Button @click="handleSubmit" :disabled="loading">
              <Loader2 v-if="loading" class="h-4 w-4 mr-1 animate-spin" /> 保存
            </Button>
            <Button variant="outline" @click="$router.back()">取消</Button>
          </div>
        </form>
      </CardContent>
    </Card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { X, Loader2 } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select'
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'
import { Switch } from '@/components/ui/switch'
import { Badge } from '@/components/ui/badge'
import adminApi from '@/utils/adminApi'
import VditorEditor from '@/components/VditorEditor.vue'
import ImageCropUpload from '@/components/ImageCropUpload.vue'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const categories = ref([])
const tags = ref([])
const newTagName = ref('')

const formErrors = reactive({
  title: '',
  category_id: '',
  content: ''
})

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

const validateForm = () => {
  let valid = true
  formErrors.title = ''
  formErrors.category_id = ''
  formErrors.content = ''

  for (const rule of rules.title) {
    if (rule.required && !form.title) {
      formErrors.title = rule.message
      valid = false
      break
    }
  }
  for (const rule of rules.category_id) {
    if (rule.required && !form.category_id) {
      formErrors.category_id = rule.message
      valid = false
      break
    }
  }
  for (const rule of rules.content) {
    if (rule.required && !form.content) {
      formErrors.content = rule.message
      valid = false
      break
    }
  }
  return valid
}

const availableTags = computed(() => {
  return tags.value.filter(t => !form.tag_ids.includes(t.id))
})

const getTagName = (tagId) => {
  const tag = tags.value.find(t => t.id === tagId)
  return tag ? tag.name : tagId
}

const addTag = (value) => {
  if (value && !form.tag_ids.includes(Number(value))) {
    form.tag_ids.push(Number(value))
  }
}

const removeTag = (tagId) => {
  const index = form.tag_ids.indexOf(tagId)
  if (index > -1) {
    form.tag_ids.splice(index, 1)
  }
}

const createNewTag = async () => {
  const tagName = newTagName.value.trim()
  if (!tagName) return

  try {
    const response = await adminApi.post('/admin/tags', {
      name: tagName,
      slug: tagName.toLowerCase().replace(/\s+/g, '-'),
      color: '#409eff'
    })

    const newTag = response.data
    tags.value.push(newTag)
    form.tag_ids.push(newTag.id)
    newTagName.value = ''

    toast.success(`标签"${tagName}"创建成功`)
  } catch (error) {
    toast.error(`创建标签"${tagName}"失败`)
  }
}

const loadCategories = async () => {
  try {
    const response = await adminApi.get('/admin/categories', {
      params: {
        page: 1,
        page_size: 100
      }
    })
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

const handleTagChange = async (values) => {
  const newTags = values.filter(v => typeof v === 'string')

  if (newTags.length > 0) {
    for (const tagName of newTags) {
      try {
        const response = await adminApi.post('/admin/tags', {
          name: tagName,
          slug: tagName.toLowerCase().replace(/\s+/g, '-'),
          color: '#409eff'
        })

        const newTag = response.data
        tags.value.push(newTag)

        const index = form.tag_ids.indexOf(tagName)
        if (index > -1) {
          form.tag_ids[index] = newTag.id
        }

        toast.success(`标签"${tagName}"创建成功`)
      } catch (error) {
        toast.error(`创建标签"${tagName}"失败`)
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
    toast.error('加载文章失败')
  }
}

const handleSubmit = async () => {
  if (!validateForm()) return

  loading.value = true
  try {
    if (isEdit.value) {
      await adminApi.put(`/admin/articles/${route.params.id}`, form)
      toast.success('更新成功')
    } else {
      await adminApi.post('/admin/articles', form)
      toast.success('创建成功')
    }
    router.push('/articles')
  } catch (error) {
    toast.error('保存失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadCategories()
  loadTags()
  loadArticle()
})
</script>

<style scoped>
.article-edit {
  padding: var(--spacing-lg);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
}

.page-header h2 {
  font-size: var(--font-size-2xl);
  color: var(--color-text-primary);
  line-height: var(--line-height-tight);
}

.content-form-item {
  margin-bottom: var(--spacing-lg);
}
</style>
