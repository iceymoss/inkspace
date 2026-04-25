<template>
  <div class="work-edit">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ isEdit ? '编辑作品' : '创建作品' }}</span>
          <el-button text @click="handleCancel">返回</el-button>
        </div>
      </template>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
        style="max-width: 1200px"
      >
        <!-- 创建模式：显示类型选择器 -->
        <el-form-item v-if="!isEdit" label="作品类型" prop="type">
          <el-radio-group v-model="form.type" @change="handleTypeChange">
            <el-radio label="project">💻 开源项目</el-radio>
            <el-radio label="photography">📷 摄影作品</el-radio>
          </el-radio-group>
          <div class="form-tip" v-if="form.type === 'photography'">
            摄影作品每天最多发布3个相册，已用：{{ quotaUsed }}/3<br>
            照片限制：{{ photoLimit }}张/相册
          </div>
        </el-form-item>
        
        <!-- 编辑模式：只显示类型标签 -->
        <el-form-item v-else label="作品类型">
          <el-tag :type="form.type === 'photography' ? 'warning' : 'primary'" size="large">
            {{ form.type === 'photography' ? '📷 摄影作品' : '💻 开源项目' }}
          </el-tag>
        </el-form-item>

        <el-form-item label="作品标题" prop="title">
          <el-input
            v-model="form.title"
            placeholder="输入作品标题"
            maxlength="100"
            show-word-limit
          />
        </el-form-item>

        <el-form-item label="作品描述" prop="description">
          <!-- 只有在非编辑模式或编辑模式下权限验证通过后才渲染VditorEditor -->
          <VditorEditor 
            v-if="!loading && (!isEdit || (isEdit && !permissionError))"
            v-model="form.description" 
            height="400px"
          />
          <div v-else style="height: 400px; display: flex; align-items: center; justify-content: center;">
            <Loader2 class="h-5 w-5 animate-spin" />
          </div>
        </el-form-item>

        <!-- 开源项目字段 -->
        <template v-if="form.type === 'project'">
          <el-divider content-position="left">项目信息</el-divider>

          <el-form-item label="封面图">
            <el-upload
              class="cover-uploader"
              :action="uploadImageUrl"
              :headers="uploadHeaders"
              :show-file-list="false"
              :on-success="handleCoverSuccess"
              :before-upload="beforeUpload"
              accept="image/*"
            >
              <el-image
                v-if="form.cover"
                :src="form.cover"
                fit="cover"
                class="cover-image"
              />
              <el-icon v-else class="cover-uploader-icon"><Plus /></el-icon>
            </el-upload>
          </el-form-item>

          <el-form-item label="项目链接">
            <el-input v-model="form.link" placeholder="项目主页URL" />
          </el-form-item>

          <el-form-item label="GitHub">
            <el-input v-model="form.github_url" placeholder="https://github.com/..." />
          </el-form-item>

          <el-form-item label="在线演示">
            <el-input v-model="form.demo_url" placeholder="演示地址" />
          </el-form-item>

          <el-form-item label="技术栈">
            <el-input v-model="form.tech_stack" placeholder="Go, Vue, MySQL" />
          </el-form-item>
        </template>

        <!-- 摄影作品字段 -->
        <template v-if="form.type === 'photography'">
          <el-divider content-position="left">相册信息</el-divider>

          <el-form-item label="拍摄地点">
            <el-input v-model="albumMetadata.location" placeholder="例如: 杭州西湖" />
          </el-form-item>

          <el-form-item label="拍摄日期">
            <el-date-picker
              v-model="albumMetadata.shooting_date"
              type="date"
              placeholder="选择日期"
              value-format="YYYY-MM-DD"
              style="width: 100%"
            />
          </el-form-item>

          <el-divider content-position="left">照片管理（{{ photos.length }}/{{ photoLimit }}）</el-divider>

          <el-form-item label="上传照片">
            <el-upload
              :action="uploadPhotoUrl"
              :headers="uploadHeaders"
              :limit="photoLimit"
              list-type="picture-card"
              :on-success="handlePhotoSuccess"
              :on-remove="handlePhotoRemove"
              :before-upload="beforePhotoUpload"
              :file-list="photoFileList"
              accept="image/jpeg,image/jpg,image/png"
              multiple
            >
              <el-icon><Plus /></el-icon>
              <template #tip>
                <div class="el-upload__tip">
                  支持 JPG/PNG，单张最大20MB，保留原图质量
                </div>
              </template>
            </el-upload>
          </el-form-item>

          <!-- 照片列表及参数编辑 -->
          <el-form-item label="照片参数" v-if="photos.length > 0">
            <div class="photos-params">
              <el-collapse v-model="activePhotoIndex" accordion>
                <el-collapse-item 
                  v-for="(photo, index) in photos" 
                  :key="index"
                  :name="index"
                >
                  <template #title>
                    <div class="photo-header">
                      <el-image :src="photo.url" class="photo-thumb" fit="cover" />
                      <span>照片 {{ index + 1 }}</span>
                      <el-tag v-if="index === 0" type="success" size="small" class="cover-tag">封面</el-tag>
                      <div style="flex: 1"></div>
                      <el-button 
                        size="small" 
                        type="danger" 
                        text
                        @click.stop="removePhoto(index)"
                      >
                        删除
                      </el-button>
                    </div>
                  </template>
                  
                  <el-form label-width="100px" class="photo-params-form">
                    <el-form-item label="照片描述">
                      <el-input 
                        v-model="photo.description" 
                        placeholder="这张照片的描述"
                        maxlength="200"
                      />
                    </el-form-item>

                    <el-row :gutter="20">
                      <el-col :span="12">
                        <el-form-item label="相机">
                          <el-input v-model="photo.metadata.camera" placeholder="Canon EOS R5" />
                        </el-form-item>
                      </el-col>
                      <el-col :span="12">
                        <el-form-item label="镜头">
                          <el-input v-model="photo.metadata.lens" placeholder="RF 24-70mm f/2.8" />
                        </el-form-item>
                      </el-col>
                    </el-row>

                    <el-row :gutter="20">
                      <el-col :span="12">
                        <el-form-item label="焦段">
                          <el-input v-model="photo.metadata.focal_length" placeholder="50mm" />
                        </el-form-item>
                      </el-col>
                      <el-col :span="12">
                        <el-form-item label="光圈">
                          <el-input v-model="photo.metadata.aperture" placeholder="f/2.8" />
                        </el-form-item>
                      </el-col>
                    </el-row>

                    <el-row :gutter="20">
                      <el-col :span="12">
                        <el-form-item label="快门">
                          <el-input v-model="photo.metadata.shutter_speed" placeholder="1/200s" />
                        </el-form-item>
                      </el-col>
                      <el-col :span="12">
                        <el-form-item label="ISO">
                          <el-input v-model="photo.metadata.iso" placeholder="400" />
                        </el-form-item>
                      </el-col>
                    </el-row>
                  </el-form>
                </el-collapse-item>
              </el-collapse>
            </div>
          </el-form-item>

          <el-alert
            title="摄影作品说明"
            type="info"
            :closable="false"
            class="alert-spacing"
          >
            • 每个相册包含多张照片（普通用户最多10张，管理员最多50张）<br>
            • 每天最多发布3个摄影相册<br>
            • 图片将保留原图质量，建议上传高质量JPG或PNG<br>
            • 第一张照片将作为相册封面
          </el-alert>
        </template>

        <el-form-item label="发布状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">立即发布</el-radio>
            <el-radio :label="0">保存为草稿</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            {{ isEdit ? '保存修改' : '发布作品' }}
          </el-button>
          <el-button @click="handleCancel">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <Dialog :open="showConfirmDialog" @update:open="onConfirmDialogUpdateOpen">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>确认</DialogTitle>
          <DialogDescription>{{ confirmDialogMessage }}</DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" @click="cancelConfirmDialog">取消</Button>
          <Button @click="confirmDialogCallback?.()">确认</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter, onBeforeRouteLeave } from 'vue-router'
import { toast } from 'vue-sonner'
import { Plus, Loader2 } from 'lucide-vue-next'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { useUserStore } from '@/stores/user'
import api from '@/utils/api'
import VditorEditor from '@/components/VditorEditor.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const formRef = ref()
const submitting = ref(false)
const loading = ref(true)
const permissionError = ref(false)
const quotaUsed = ref(0)
const activePhotoIndex = ref(0)
const hasUnsavedChanges = ref(false)
const originalData = ref(null)
const isSaved = ref(false)

const showConfirmDialog = ref(false)
const confirmDialogMessage = ref('')
const confirmDialogCallback = ref(null)
let _confirmDialogReject = null

const confirmDialog = (message) => {
  return new Promise((resolve, reject) => {
    confirmDialogMessage.value = message
    _confirmDialogReject = reject
    confirmDialogCallback.value = () => {
      _confirmDialogReject = null
      showConfirmDialog.value = false
      resolve()
    }
    showConfirmDialog.value = true
  })
}

const onConfirmDialogUpdateOpen = (open) => {
  showConfirmDialog.value = open
  if (!open && _confirmDialogReject) {
    const rejectFn = _confirmDialogReject
    _confirmDialogReject = null
    rejectFn('cancel')
  }
}

const cancelConfirmDialog = () => {
  if (_confirmDialogReject) {
    const rejectFn = _confirmDialogReject
    _confirmDialogReject = null
    showConfirmDialog.value = false
    rejectFn('cancel')
  }
}

const isEdit = computed(() => !!route.params.id)

// 根据用户角色确定照片限制
const photoLimit = computed(() => {
  return userStore.user?.role === 'admin' ? 50 : 10
})

const form = reactive({
  title: '',
  type: 'project',
  description: '',
  cover: '',
  link: '',
  github_url: '',
  demo_url: '',
  tech_stack: '',
  status: 1
})

// 相册元数据
const albumMetadata = reactive({
  location: '',
  shooting_date: ''
})

// 照片数组
const photos = ref([])

// 用于 el-upload 的文件列表
const photoFileList = ref([])

const rules = {
  title: [{ required: true, message: '请输入作品标题', trigger: 'blur' }],
  type: [{ required: true, message: '请选择作品类型', trigger: 'change' }]
}

// 使用相对路径，通过 Nginx 代理到后端
const uploadImageUrl = import.meta.env.VITE_API_URL 
  ? `${import.meta.env.VITE_API_URL}/api/upload/image`
  : '/api/upload/image'
const uploadPhotoUrl = import.meta.env.VITE_API_URL
  ? `${import.meta.env.VITE_API_URL}/api/upload/photo`
  : '/api/upload/photo'
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${userStore.token}`
}))

const handleTypeChange = () => {
  // 切换类型时清空对应字段
  if (form.type === 'photography') {
    form.link = ''
    form.github_url = ''
    form.demo_url = ''
    form.tech_stack = ''
  } else {
    photos.value = []
    photoFileList.value = []
    Object.assign(albumMetadata, {
      location: '',
      shooting_date: ''
    })
  }
}

const handleCoverSuccess = (response) => {
  if (response.code === 0) {
    form.cover = response.data.url
    toast.success('封面上传成功')
  }
}

const beforeUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt10M = file.size / 1024 / 1024 < 10

  if (!isImage) {
    toast.error('只能上传图片')
    return false
  }
  if (!isLt10M) {
    toast.error('图片大小不能超过 10MB')
    return false
  }
  return true
}

const beforePhotoUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt20M = file.size / 1024 / 1024 < 20

  if (!isImage) {
    toast.error('只能上传图片')
    return false
  }
  if (!isLt20M) {
    toast.error('照片大小不能超过 20MB')
    return false
  }
  if (photos.value.length >= photoLimit.value) {
    toast.warning(`照片数量已达上限（${photoLimit.value}张）`)
    return false
  }
  return true
}

const handlePhotoSuccess = (response, file) => {
  if (response.code === 0) {
    // 添加新照片到数组
    photos.value.push({
      url: response.data.url,
      description: '',
      metadata: {
        camera: '',
        lens: '',
        focal_length: '',
        aperture: '',
        shutter_speed: '',
        iso: ''
      }
    })
    
    // 如果是第一张照片，设置为封面
    if (photos.value.length === 1) {
      form.cover = response.data.url
    }
    
    toast.success('照片上传成功')
  }
}

const handlePhotoRemove = (file) => {
  // 从 photos 数组中移除对应的照片
  const index = photoFileList.value.findIndex(f => f.uid === file.uid)
  if (index !== -1 && index < photos.value.length) {
    photos.value.splice(index, 1)
    
    // 如果删除的是封面，更新封面为第一张照片
    if (photos.value.length > 0 && form.cover === file.url) {
      form.cover = photos.value[0].url
    }
  }
}

const removePhoto = (index) => {
  photos.value.splice(index, 1)
  photoFileList.value.splice(index, 1)
  
  // 更新封面
  if (photos.value.length > 0 && index === 0) {
    form.cover = photos.value[0].url
  }
}

const loadQuota = async () => {
  try {
    const response = await api.get('/works/quota')
    quotaUsed.value = response.data.used || 0
  } catch (error) {
    console.error('Failed to load quota:', error)
  }
}

// 检查表单是否有未保存的更改
const checkUnsavedChanges = () => {
  if (!isEdit.value || !originalData.value) {
    // 创建模式下，检查是否有任何输入
    return form.title || form.description || form.cover || 
           (form.type === 'project' && (form.link || form.github_url || form.demo_url || form.tech_stack)) ||
           (form.type === 'photography' && (photos.value.length > 0 || albumMetadata.location || albumMetadata.shooting_date))
  }
  
  // 编辑模式下，与原始数据对比
  const original = originalData.value
  
  // 检查基础字段
  if (form.title !== original.title ||
      form.description !== original.description ||
      form.cover !== original.cover ||
      form.status !== original.status) {
    return true
  }
  
  // 检查项目类型字段
  if (form.type === 'project') {
    if (form.link !== original.link ||
        form.github_url !== original.github_url ||
        form.demo_url !== original.demo_url ||
        form.tech_stack !== original.tech_stack) {
      return true
    }
  }
  
  // 检查摄影类型字段
  if (form.type === 'photography') {
    if (albumMetadata.location !== (original.metadata?.location || '') ||
        albumMetadata.shooting_date !== (original.metadata?.shooting_date || '')) {
      return true
    }
    
    // 检查照片数组
    if (JSON.stringify(photos.value) !== JSON.stringify(original.images || [])) {
      return true
    }
  }
  
  return false
}

const loadWork = async () => {
  // 重置保存状态和权限错误状态
  isSaved.value = false
  permissionError.value = false
  
  if (!isEdit.value) {
    // 创建模式下，检查是否登录
    if (!userStore.isLoggedIn) {
      toast.warning('请先登录')
      router.push('/login')
      return
    }
    hasUnsavedChanges.value = false
    originalData.value = null
    loading.value = false // 创建模式不需要加载数据，可以直接渲染VditorEditor
    return
  }

  // 编辑模式下，先检查是否登录
  if (!userStore.isLoggedIn) {
    toast.warning('请先登录')
    router.push('/login')
    return
  }

  // 确保在权限检查完成之前，loading保持为true，防止VditorEditor被渲染
  loading.value = true
  permissionError.value = false
  submitting.value = true
  try {
    // 使用编辑专用API，后端会进行权限检查
    const response = await api.get(`/works/${route.params.id}/edit`)
    const work = response.data

    Object.assign(form, {
      title: work.title,
      type: work.type || 'project',
      description: work.description,
      cover: work.cover,
      link: work.link,
      github_url: work.github_url,
      demo_url: work.demo_url,
      tech_stack: work.tech_stack,
      status: work.status
    })

    // 加载摄影相册数据
    if (work.type === 'photography') {
      photos.value = work.images || []
      
      // 加载相册元数据
      if (work.metadata) {
        Object.assign(albumMetadata, {
          location: work.metadata.location || '',
          shooting_date: work.metadata.shooting_date || ''
        })
      }
      
      // 构建文件列表用于显示
      photoFileList.value = photos.value.map((photo, index) => ({
        uid: index,
        name: `photo-${index}`,
        url: photo.url,
        status: 'success'
      }))
    }
    
    // 保存原始数据用于对比
    originalData.value = {
      title: work.title,
      type: work.type || 'project',
      description: work.description,
      cover: work.cover,
      link: work.link,
      github_url: work.github_url,
      demo_url: work.demo_url,
      tech_stack: work.tech_stack,
      status: work.status,
      images: work.images || [],
      metadata: work.metadata || {}
    }
    
    hasUnsavedChanges.value = false
    loading.value = false // 加载成功，允许渲染VditorEditor
  } catch (error) {
    // 标记权限错误，防止渲染VditorEditor
    permissionError.value = true
    loading.value = true // 保持loading为true，确保不渲染VditorEditor
    
    // 处理各种错误情况
    const status = error.response?.status
    if (status === 403) {
      toast.error('无权限编辑此作品')
    } else if (status === 404) {
      toast.error('作品不存在')
    } else if (status === 401) {
      toast.error('请先登录')
      router.push('/login')
      return
    } else {
      toast.error('加载作品失败')
    }
    
    // 标记为已保存，避免路由守卫拦截
    isSaved.value = true
    // 立即跳转，不等待任何异步操作，避免VditorEditor被渲染
    router.push('/dashboard/works')
    return // 立即返回，不执行 finally 中的代码
  } finally {
    submitting.value = false
  }
}

const handleSubmit = async () => {
  await formRef.value.validate(async (valid) => {
    if (!valid) return

    // 验证摄影作品照片数量
    if (form.type === 'photography') {
      if (photos.value.length === 0) {
        toast.warning('请至少上传1张照片')
        return
      }
      if (photos.value.length > photoLimit.value) {
        toast.warning(`照片数量超过限制（最多${photoLimit.value}张）`)
        return
      }
      
      // 检查配额
      if (!isEdit.value && quotaUsed.value >= 3) {
        toast.warning('今日摄影作品发布数量已达上限（3个相册/天）')
        return
      }
    }

    submitting.value = true
    try {
      const submitData = {
        title: form.title,
        type: form.type,
        description: form.description,
        cover: form.cover,
        status: form.status,
        images: [],
        metadata: {}
      }

      if (form.type === 'project') {
        // 项目类型：简单字符串数组（保持向后兼容）
        submitData.images = []
        submitData.link = form.link
        submitData.github_url = form.github_url
        submitData.demo_url = form.demo_url
        submitData.tech_stack = form.tech_stack
      } else if (form.type === 'photography') {
        // 摄影类型：PhotoItem 对象数组
        submitData.images = photos.value
        submitData.cover = photos.value[0]?.url || ''
        submitData.metadata = {
          ...albumMetadata,
          photo_count: photos.value.length
        }
      }

      if (isEdit.value) {
        await api.put(`/works/${route.params.id}`, submitData)
        toast.success('更新成功')
      } else {
        await api.post('/works', submitData)
        toast.success('创建成功')
      }
      
      // 清除未保存标记，标记为已保存
      hasUnsavedChanges.value = false
      originalData.value = null
      isSaved.value = true
      
      router.push('/dashboard/works')
    } catch (error) {
      toast.error(error.response?.data?.message || '保存失败')
    } finally {
      submitting.value = false
    }
  })
}

// 处理取消/返回按钮
const handleCancel = async () => {
  if (checkUnsavedChanges()) {
    try {
      await confirmDialog(
        '您有未保存的更改，确定要离开吗？'
      )
      router.back()
    } catch {
      // 用户取消，不执行任何操作
    }
  } else {
    router.back()
  }
}

// 路由守卫：离开页面前的提示
onBeforeRouteLeave((to, from, next) => {
  // 如果已保存成功，直接允许离开
  if (isSaved.value) {
    next()
    return
  }
  
  // 如果正在提交，直接允许离开
  if (submitting.value) {
    next()
    return
  }
  
  // 如果跳转到作品列表页（可能是权限错误导致的跳转），直接允许离开
  if (to.path === '/dashboard/works') {
    next()
    return
  }
  
  if (checkUnsavedChanges()) {
    confirmDialog(
      '您有未保存的更改，确定要离开吗？'
    ).then(() => {
      next()
    }).catch(() => {
      next(false)
    })
  } else {
    next()
  }
})

// 监听页面刷新/关闭
const handleBeforeUnload = (e) => {
  // 如果已保存或正在提交，不提示
  if (isSaved.value || submitting.value) {
    return
  }
  
  if (checkUnsavedChanges()) {
    e.preventDefault()
    e.returnValue = '您有未保存的更改，确定要离开吗？'
    return e.returnValue
  }
}

// 监听表单变化
watch(
  [() => form.title, () => form.description, () => form.cover, () => form.status,
   () => form.link, () => form.github_url, () => form.demo_url, () => form.tech_stack,
   () => albumMetadata.location, () => albumMetadata.shooting_date,
   () => photos.value],
  () => {
    hasUnsavedChanges.value = checkUnsavedChanges()
  },
  { deep: true }
)

onMounted(() => {
  loadQuota()
  loadWork()
  // 添加页面刷新/关闭监听
  window.addEventListener('beforeunload', handleBeforeUnload)
})

onUnmounted(() => {
  // 移除页面刷新/关闭监听
  window.removeEventListener('beforeunload', handleBeforeUnload)
})
</script>

<style scoped>
.work-edit {
  max-width: 1200px;
  background: var(--theme-bg-secondary);
  min-height: calc(100vh - 60px);
  padding: var(--spacing-lg) 0;
}

.work-edit .el-card {
  box-shadow: var(--shadow-md);
  border-radius: var(--radius-md);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header span {
  font-size: var(--font-size-2xl);
  font-weight: 600;
  color: var(--theme-text-primary);
}

.form-tip {
  font-size: var(--font-size-xs);
  color: var(--theme-text-tertiary);
  margin-top: 5px;
  line-height: var(--line-height-base);
}

.cover-uploader {
  width: 200px;
  height: 150px;
  border: 1px dashed var(--theme-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: border-color var(--transition-base);
}

.cover-uploader:hover {
  border-color: var(--theme-primary);
}

.cover-image {
  width: 200px;
  height: 150px;
  border-radius: var(--radius-sm);
}

.cover-uploader-icon {
  font-size: var(--font-size-2xl);
  color: var(--theme-text-tertiary);
}

.photos-params {
  width: 100%;
}

.photo-header {
  display: flex;
  align-items: center;
  width: 100%;
}

.photo-thumb {
  width: 60px;
  height: 60px;
  margin-right: var(--spacing-sm);
  border-radius: var(--radius-sm);
}

.cover-tag {
  margin-left: var(--spacing-sm);
}

.photo-params-form {
  padding: var(--spacing-sm);
}

.alert-spacing {
  margin-bottom: var(--spacing-lg);
}

:deep(.el-divider__text) {
  color: var(--theme-text-secondary);
  font-size: var(--font-size-sm);
  font-weight: 600;
}

:deep(.el-form-item__label) {
  color: var(--theme-text-secondary);
  font-weight: 600;
}

:deep(.el-radio) {
  cursor: pointer;
}

:deep(.el-upload--picture-card) {
  border-radius: var(--radius-sm);
  border-color: var(--theme-border);
  transition: border-color var(--transition-base);
}

:deep(.el-upload--picture-card:hover) {
  border-color: var(--theme-primary);
}
</style>
