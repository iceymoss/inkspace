<template>
  <div class="works">
    <div class="header">
      <h2>作品管理</h2>
      <el-button type="primary" @click="showDialog()"><Plus class="h-4 w-4" /> 新建作品</el-button>
    </div>

    <!-- 筛选区域 -->
    <div class="filter-bar">
      <el-form inline>
        <el-form-item label="关键字">
          <el-input
            v-model="keyword"
            placeholder="标题 / 描述"
            clearable
            style="width: 240px"
            @keyup.enter="handleSearch"
          />
        </el-form-item>

        <el-form-item label="状态">
          <el-radio-group v-model="statusFilter" @change="handleStatusFilterChange" size="default">
            <el-radio-button :label="null">全部</el-radio-button>
            <el-radio-button :label="1">审核通过</el-radio-button>
            <el-radio-button :label="2">待审核</el-radio-button>
            <el-radio-button :label="3">审核不通过</el-radio-button>
            <el-radio-button :label="0">草稿</el-radio-button>
          </el-radio-group>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table
      :data="works"
      class="works-table"
      @sort-change="handleSortChange"
    >
      <el-table-column prop="id" label="ID" width="80" sortable="custom" />
      <el-table-column prop="title" label="标题" sortable="custom" />
      <el-table-column label="封面" width="120">
        <template #default="{ row }">
          <el-image :src="row.cover" style="width: 80px; height: 60px;" fit="cover" />
        </template>
      </el-table-column>
      <el-table-column prop="type" label="类型" width="120" sortable="custom">
        <template #default="{ row }">
          <el-tag :type="row.type === 'photography' ? 'warning' : 'primary'">
            {{ row.type === 'photography' ? '📷 摄影' : '💻 项目' }}
          </el-tag>
          <div v-if="row.type === 'photography' && row.metadata?.photo_count" class="cell-detail">
            {{ row.metadata.photo_count }} 张
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="view_count" label="浏览" width="100" sortable="custom" />
      <el-table-column prop="status" label="状态" width="120" sortable="custom">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">
            {{ getStatusText(row.status) }}
          </el-tag>
          <div v-if="row.audit_message" class="audit-hint" :title="row.audit_message">
            {{ row.audit_message }}
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="is_recommend" label="推荐" width="90" align="center" sortable="custom">
        <template #default="{ row }">
          <el-tag :type="row.is_recommend ? 'warning' : ''">
            {{ row.is_recommend ? '★ 推荐' : '-' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180" sortable="custom">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column prop="updated_at" label="处理时间" width="180" sortable="custom">
        <template #default="{ row }">
          {{ formatDate(row.updated_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="350" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="showDialog(row)">编辑</el-button>
          <el-button 
            size="small" 
            :type="row.is_recommend ? 'warning' : 'default'"
            @click="handleToggleRecommend(row)"
          >
            {{ row.is_recommend ? '取消推荐' : '推荐' }}
          </el-button>
          <el-button 
            v-if="row.status === 2" 
            size="small" 
            type="success" 
            @click="handleApprove(row)"
          >
            通过
          </el-button>
          <el-button 
            v-if="row.status === 2 || row.status === 1" 
            size="small" 
            type="warning" 
            @click="handleReject(row)"
          >
            拒绝
          </el-button>
          <el-button 
            v-if="row.status === 3" 
            size="small" 
            type="success" 
            @click="handleApprove(row)"
          >
            通过
          </el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog 
      v-model="dialogVisible" 
      :title="isEdit ? '编辑作品' : '新建作品'" 
      width="900px"
      :close-on-click-modal="false"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="作品类型" prop="type">
          <el-radio-group v-model="form.type" @change="handleTypeChange">
            <el-radio label="project">💻 开源项目</el-radio>
            <el-radio label="photography">📷 摄影作品（相册）</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="输入作品标题" />
        </el-form-item>

        <el-form-item label="作品描述" prop="description">
          <VditorEditor v-model="form.description" height="400px" />
        </el-form-item>

        <!-- 开源项目特有字段 -->
        <template v-if="form.type === 'project'">
          <el-form-item label="封面图" prop="cover">
            <!-- 0 表示不限制比例，支持任意裁切 -->
            <ImageCropUpload 
              v-model="form.cover"
              preview-size="160px"
              :aspect-ratio="0"
              :output-width="1600"
              :output-height="1600"
              placeholder="点击上传封面图"
              tip="支持任意比例裁切，最大 5MB"
              :max-size="5"
            />
          </el-form-item>

          <el-form-item label="项目链接">
            <el-input v-model="form.link" placeholder="项目主页URL" />
          </el-form-item>
          
          <el-form-item label="GitHub">
            <el-input v-model="form.github_url" placeholder="GitHub 仓库URL" />
          </el-form-item>
          
          <el-form-item label="在线演示">
            <el-input v-model="form.demo_url" placeholder="演示地址" />
          </el-form-item>
          
          <el-form-item label="技术栈">
            <el-input v-model="form.tech_stack" placeholder="例如: Go, Vue, MySQL" />
          </el-form-item>
        </template>

        <!-- 摄影作品特有字段 -->
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

          <el-divider content-position="left">照片管理（{{ photos.length }} 张）</el-divider>

          <el-form-item label="添加照片">
            <el-button @click="addPhoto" :disabled="photos.length >= 50">
              <Plus class="h-4 w-4" /> 添加照片
            </el-button>
            <el-text size="small" type="info" style="margin-left: 10px;">
              最多 50 张照片
            </el-text>
          </el-form-item>

          <!-- 照片列表 -->
          <div class="photos-list" v-if="photos.length > 0">
            <el-collapse v-model="activePhotoIndex" accordion>
              <el-collapse-item 
                v-for="(photo, index) in photos" 
                :key="index"
                :name="index"
              >
                <template #title>
                  <div class="photo-header">
                    <el-image 
                      v-if="photo.url"
                      :src="photo.url" 
                      style="width: 60px; height: 60px; margin-right: 10px" 
                      fit="cover" 
                    />
                    <span>照片 {{ index + 1 }}</span>
                    <el-tag v-if="index === 0" type="success" size="small" style="margin-left: 10px">
                      封面
                    </el-tag>
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
                
                <el-form label-width="100px" style="padding: 10px">
                  <el-form-item label="照片" required>
                    <!-- aspect-ratio=0 支持任意裁切比例 -->
                    <ImageCropUpload 
                      v-model="photo.url"
                      preview-size="160px"
                      :aspect-ratio="0"
                      :output-width="1600"
                      :output-height="1600"
                      placeholder="点击上传照片"
                      tip="支持任意比例裁切，最大 5MB"
                      :max-size="5"
                    />
                  </el-form-item>

                  <el-form-item label="照片描述">
                    <el-input 
                      v-model="photo.description" 
                      placeholder="这张照片的描述"
                      maxlength="200"
                    />
                  </el-form-item>

                  <el-divider content-position="left">拍摄参数（选填）</el-divider>

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
        </template>

        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">发布</el-radio>
            <el-radio :label="0">草稿</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="推荐" prop="is_recommend">
          <el-switch v-model="form.is_recommend" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <Dialog v-model:open="confirmDialogVisible">
      <DialogContent class="sm:max-w-[400px]">
        <DialogHeader>
          <DialogTitle>{{ confirmDialogConfig.title }}</DialogTitle>
        </DialogHeader>
        <p class="text-sm text-muted-foreground">{{ confirmDialogConfig.message }}</p>
        <DialogFooter>
          <Button variant="outline" @click="handleConfirmCancel">取消</Button>
          <Button variant="destructive" @click="handleConfirmOk">确定</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <Dialog v-model:open="promptDialogVisible">
      <DialogContent class="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>{{ promptDialogConfig.title }}</DialogTitle>
        </DialogHeader>
        <p class="text-sm text-muted-foreground">{{ promptDialogConfig.message }}</p>
        <el-input
          v-model="promptDialogConfig.inputValue"
          type="textarea"
          :rows="4"
          :placeholder="promptDialogConfig.placeholder"
        />
        <DialogFooter>
          <Button variant="outline" @click="handlePromptCancel">取消</Button>
          <Button @click="handlePromptOk">确定</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { toast } from 'vue-sonner'
import { Plus } from 'lucide-vue-next'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import adminApi from '@/utils/adminApi'
import dayjs from 'dayjs'
import ImageCropUpload from '@/components/ImageCropUpload.vue'
import VditorEditor from '@/components/VditorEditor.vue'

const works = ref([])
const dialogVisible = ref(false)
const formRef = ref()
const isEdit = ref(false)
const activePhotoIndex = ref(0)
const statusFilter = ref(null) // 状态筛选：null=全部, 0=草稿, 1=已发布, 2=待审核, 3=审核不通过
const sortField = ref(null)
const sortOrder = ref(null) // 'ascending' | 'descending' | null
const keyword = ref('')

const confirmDialogVisible = ref(false)
const confirmDialogConfig = reactive({
  title: '提示',
  message: '',
  onConfirm: null,
  onCancel: null,
})

const confirmDialog = (message, title = '提示') => {
  return new Promise((resolve, reject) => {
    confirmDialogConfig.title = title
    confirmDialogConfig.message = message
    confirmDialogConfig.onConfirm = resolve
    confirmDialogConfig.onCancel = () => reject('cancel')
    confirmDialogVisible.value = true
  })
}

const handleConfirmOk = () => {
  confirmDialogVisible.value = false
  confirmDialogConfig.onConfirm?.()
}

const handleConfirmCancel = () => {
  confirmDialogVisible.value = false
  confirmDialogConfig.onCancel?.()
}

const promptDialogVisible = ref(false)
const promptDialogConfig = reactive({
  title: '',
  message: '',
  placeholder: '',
  inputValue: '',
  required: false,
  requiredMessage: '',
  onConfirm: null,
  onCancel: null,
})

const promptDialog = (message, title, options = {}) => {
  return new Promise((resolve, reject) => {
    promptDialogConfig.title = title
    promptDialogConfig.message = message
    promptDialogConfig.placeholder = options.inputPlaceholder || ''
    promptDialogConfig.inputValue = ''
    promptDialogConfig.required = !!options.inputValidator
    promptDialogConfig.requiredMessage = ''
    if (options.inputValidator) {
      const testResult = options.inputValidator('')
      if (typeof testResult === 'string') {
        promptDialogConfig.requiredMessage = testResult
      }
    }
    promptDialogConfig.onConfirm = () => resolve({ value: promptDialogConfig.inputValue })
    promptDialogConfig.onCancel = () => reject('cancel')
    promptDialogVisible.value = true
  })
}

const handlePromptOk = () => {
  if (promptDialogConfig.required && !promptDialogConfig.inputValue.trim()) {
    toast.error(promptDialogConfig.requiredMessage || '此项为必填')
    return
  }
  promptDialogVisible.value = false
  promptDialogConfig.onConfirm?.()
}

const handlePromptCancel = () => {
  promptDialogVisible.value = false
  promptDialogConfig.onCancel?.()
}

const form = reactive({
  id: null,
  title: '',
  type: 'project',
  description: '',
  cover: '',
  link: '',
  github_url: '',
  demo_url: '',
  tech_stack: '',
  status: 1,
  is_recommend: false,
  sort: 0
})

// 相册元数据
const albumMetadata = reactive({
  location: '',
  shooting_date: ''
})

// 照片数组
const photos = ref([])

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }]
}

const formatDate = (date) => {
  if (!date) return ''
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

const loadWorks = async () => {
  try {
    const params = {}
    if (statusFilter.value !== null) {
      params.status = statusFilter.value
    }
    if (keyword.value) {
      params.keyword = keyword.value
    }
    if (sortField.value && sortOrder.value) {
      const dir = sortOrder.value === 'ascending' ? 'asc' : 'desc'
      params.sort = `${sortField.value}_${dir}`
    }
    const response = await adminApi.get('/admin/works', { params })
    works.value = response.data.list || []
  } catch (error) {
    toast.error('加载失败')
  }
}

const handleStatusFilterChange = () => {
  loadWorks()
}

const handleSortChange = ({ prop, order }) => {
  // Element Plus: order 为 'ascending' / 'descending' / null
  sortField.value = order ? prop : null
  sortOrder.value = order
  loadWorks()
}

const handleSearch = () => {
  loadWorks()
}

const resetFilters = () => {
  keyword.value = ''
  statusFilter.value = null
  sortField.value = null
  sortOrder.value = null
  loadWorks()
}

const showDialog = (work = null) => {
  isEdit.value = !!work
  
  if (work) {
    Object.assign(form, {
      id: work.id,
      title: work.title,
      type: work.type || 'project',
      description: work.description,
      cover: work.cover,
      link: work.link,
      github_url: work.github_url,
      demo_url: work.demo_url,
      tech_stack: work.tech_stack,
      status: work.status,
      is_recommend: work.is_recommend,
      sort: work.sort || 0
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
    }
  } else {
    resetForm()
  }
  
  dialogVisible.value = true
}

const resetForm = () => {
  Object.assign(form, {
    id: null,
    title: '',
    type: 'project',
    description: '',
    cover: '',
    link: '',
    github_url: '',
    demo_url: '',
    tech_stack: '',
    status: 1,
    is_recommend: false,
    sort: 0
  })
  Object.assign(albumMetadata, {
    location: '',
    shooting_date: ''
  })
  photos.value = []
}

const handleTypeChange = () => {
  // 切换类型时清空对应字段
  if (form.type === 'photography') {
    form.link = ''
    form.github_url = ''
    form.demo_url = ''
    form.tech_stack = ''
  } else {
    photos.value = []
    Object.assign(albumMetadata, {
      location: '',
      shooting_date: ''
    })
  }
}

const addPhoto = () => {
  if (photos.value.length >= 50) {
    toast.warning('照片数量已达上限（50张）')
    return
  }
  
  photos.value.push({
    url: '',
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
  
  // 展开新添加的照片
  activePhotoIndex.value = photos.value.length - 1
}

const removePhoto = (index) => {
  photos.value.splice(index, 1)
}

const handleSubmit = async () => {
  await formRef.value.validate(async (valid) => {
    if (!valid) return

    try {
      const submitData = {
        title: form.title,
        type: form.type,
        description: form.description,
        status: form.status,
        is_recommend: form.is_recommend,
        sort: form.sort,
        images: [],
        metadata: {}
      }

      if (form.type === 'project') {
        // 项目类型
        submitData.cover = form.cover
        submitData.link = form.link
        submitData.github_url = form.github_url
        submitData.demo_url = form.demo_url
        submitData.tech_stack = form.tech_stack
        submitData.images = [] // 空数组
      } else if (form.type === 'photography') {
        // 摄影类型
        if (photos.value.length === 0) {
          toast.warning('请至少添加1张照片')
          return
        }
        
        // 验证所有照片都有URL
        const hasEmptyUrl = photos.value.some(p => !p.url)
        if (hasEmptyUrl) {
          toast.warning('请填写所有照片的URL')
          return
        }
        
        submitData.images = photos.value
        submitData.cover = photos.value[0]?.url || ''
        submitData.metadata = {
          ...albumMetadata,
          photo_count: photos.value.length
        }
      }

      if (isEdit.value) {
        await adminApi.put(`/admin/works/${form.id}`, submitData)
        toast.success('更新成功')
      } else {
        await adminApi.post('/admin/works', submitData)
        toast.success('创建成功')
      }
      
      dialogVisible.value = false
      loadWorks()
    } catch (error) {
      toast.error(error.response?.data?.message || '操作失败')
    }
  })
}

const handleToggleRecommend = async (work) => {
  try {
    await adminApi.put(`/admin/works/${work.id}/recommend`, {
      is_recommend: !work.is_recommend
    })
    toast.success('设置成功')
    loadWorks()
  } catch (error) {
    toast.error('操作失败')
  }
}

const handleDelete = async (work) => {
  try {
    await confirmDialog('确定要删除这个作品吗？', '提示')
    
    await adminApi.delete(`/admin/works/${work.id}`)
    toast.success('删除成功')
    loadWorks()
  } catch (error) {
    if (error !== 'cancel') {
      toast.error('删除失败')
    }
  }
}

// 获取状态文本
const getStatusText = (status) => {
  switch (status) {
    case 1:
      return '已发布'
    case 0:
      return '草稿'
    case 2:
      return '待审核'
    case 3:
      return '审核不通过'
    default:
      return '未知'
  }
}

// 获取状态标签类型
const getStatusType = (status) => {
  switch (status) {
    case 1:
      return 'success'
    case 0:
      return 'info'
    case 2:
      return 'warning'
    case 3:
      return 'danger'
    default:
      return 'info'
  }
}

// 审核通过
const handleApprove = async (work) => {
  try {
    const { value: auditMessage } = await promptDialog(
      '请输入审核消息（可选，用于记录审核通过的原因）',
      '审核通过',
      {
        inputPlaceholder: '请输入审核消息，例如：内容符合规范，审核通过...',
      }
    )
    
    await adminApi.put(`/admin/works/${work.id}/status`, {
      status: 1,
      audit_message: auditMessage || ''
    })
    toast.success('审核通过')
    loadWorks()
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') {
      toast.error('操作失败')
    }
  }
}

const handleReject = async (work) => {
  try {
    const { value: auditMessage } = await promptDialog(
      '请输入拒绝原因（必填，用于告知作者审核不通过的原因）',
      '审核拒绝',
      {
        inputPlaceholder: '请输入拒绝原因，例如：内容不符合规范、包含违规信息等...',
        inputValidator: (value) => {
          if (!value || value.trim() === '') {
            return '拒绝原因不能为空'
          }
          return true
        }
      }
    )
    
    await adminApi.put(`/admin/works/${work.id}/status`, {
      status: 3,
      audit_message: auditMessage || ''
    })
    toast.success('已拒绝')
    loadWorks()
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') {
      toast.error('操作失败')
    }
  }
}

loadWorks()
</script>

<style scoped>
.works {
  padding: var(--spacing-lg);
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
}

.filter-bar {
  padding: var(--spacing-md);
  margin-top: var(--spacing-lg);
  margin-bottom: var(--spacing-lg);
  background: var(--color-bg-card);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
}

.works-table {
  width: 100%;
  margin-top: var(--spacing-lg);
}

.photos-list {
  margin-top: var(--spacing-lg);
}

.photo-header {
  display: flex;
  align-items: center;
  width: 100%;
}

.cell-detail {
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
  margin-top: var(--spacing-xs);
}

.audit-hint {
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
  margin-top: var(--spacing-xs);
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
