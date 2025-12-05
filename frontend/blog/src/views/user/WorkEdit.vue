<template>
  <div class="work-edit">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ isEdit ? '编辑作品' : '创建作品' }}</span>
          <el-button text @click="$router.back()">返回</el-button>
        </div>
      </template>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
        style="max-width: 1200px"
      >
        <el-form-item label="作品类型" prop="type">
          <el-radio-group v-model="form.type" @change="handleTypeChange">
            <el-radio label="project">💻 开源项目</el-radio>
            <el-radio label="photography">📷 摄影作品</el-radio>
          </el-radio-group>
          <div class="form-tip" v-if="form.type === 'photography'">
            摄影作品每天最多发布3个相册，已用：{{ quotaUsed }}/3<br>
            照片限制：{{ photoLimit }}张/相册
          </div>
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
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="4"
            placeholder="作品描述"
            maxlength="500"
            show-word-limit
          />
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
                      <el-image :src="photo.url" style="width: 60px; height: 60px; margin-right: 10px" fit="cover" />
                      <span>照片 {{ index + 1 }}</span>
                      <el-tag v-if="index === 0" type="success" size="small" style="margin-left: 10px">封面</el-tag>
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
            style="margin-bottom: 20px"
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
import { Plus } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import api from '@/utils/api'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const formRef = ref()
const submitting = ref(false)
const quotaUsed = ref(0)
const activePhotoIndex = ref(0)

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

const uploadImageUrl = `${import.meta.env.VITE_API_URL || 'http://localhost:8081'}/api/upload/image`
const uploadPhotoUrl = `${import.meta.env.VITE_API_URL || 'http://localhost:8081'}/api/upload/photo`
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
    ElMessage.success('封面上传成功')
  }
}

const beforeUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt10M = file.size / 1024 / 1024 < 10

  if (!isImage) {
    ElMessage.error('只能上传图片')
    return false
  }
  if (!isLt10M) {
    ElMessage.error('图片大小不能超过 10MB')
    return false
  }
  return true
}

const beforePhotoUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt20M = file.size / 1024 / 1024 < 20

  if (!isImage) {
    ElMessage.error('只能上传图片')
    return false
  }
  if (!isLt20M) {
    ElMessage.error('照片大小不能超过 20MB')
    return false
  }
  if (photos.value.length >= photoLimit.value) {
    ElMessage.warning(`照片数量已达上限（${photoLimit.value}张）`)
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
    
    ElMessage.success('照片上传成功')
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

const loadWork = async () => {
  if (!isEdit.value) return

  try {
    const response = await api.get(`/works/${route.params.id}`)
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
  } catch (error) {
    ElMessage.error('加载作品失败')
  }
}

const handleSubmit = async () => {
  await formRef.value.validate(async (valid) => {
    if (!valid) return

    // 验证摄影作品照片数量
    if (form.type === 'photography') {
      if (photos.value.length === 0) {
        ElMessage.warning('请至少上传1张照片')
        return
      }
      if (photos.value.length > photoLimit.value) {
        ElMessage.warning(`照片数量超过限制（最多${photoLimit.value}张）`)
        return
      }
      
      // 检查配额
      if (!isEdit.value && quotaUsed.value >= 3) {
        ElMessage.warning('今日摄影作品发布数量已达上限（3个相册/天）')
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
        ElMessage.success('更新成功')
      } else {
        await api.post('/works', submitData)
        ElMessage.success('创建成功')
      }
      
      router.push('/dashboard/works')
    } catch (error) {
      ElMessage.error(error.response?.data?.message || '保存失败')
    } finally {
      submitting.value = false
    }
  })
}

onMounted(() => {
  loadQuota()
  loadWork()
})
</script>

<style scoped>
.work-edit {
  max-width: 1200px;
}

.work-edit .el-card {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}

.cover-uploader {
  width: 200px;
  height: 150px;
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.cover-uploader:hover {
  border-color: #409eff;
}

.cover-image {
  width: 200px;
  height: 150px;
}

.cover-uploader-icon {
  font-size: 28px;
  color: #8c939d;
}

.photos-params {
  width: 100%;
}

.photo-header {
  display: flex;
  align-items: center;
  width: 100%;
}
</style>
