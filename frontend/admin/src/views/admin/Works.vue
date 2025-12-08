<template>
  <div class="works">
    <h2>作品管理</h2>
    <el-button type="primary" @click="showDialog()"><el-icon><Plus /></el-icon> 新建作品</el-button>

    <el-table :data="works" style="width: 100%; margin-top: 20px;">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="title" label="标题" />
      <el-table-column label="封面" width="120">
        <template #default="{ row }">
          <el-image :src="row.cover" style="width: 80px; height: 60px;" fit="cover" />
        </template>
      </el-table-column>
      <el-table-column label="类型" width="120">
        <template #default="{ row }">
          <el-tag :type="row.type === 'photography' ? 'warning' : 'primary'">
            {{ row.type === 'photography' ? '📷 摄影' : '💻 项目' }}
          </el-tag>
          <div v-if="row.type === 'photography' && row.metadata?.photo_count" style="font-size: 12px; color: #909399; margin-top: 4px;">
            {{ row.metadata.photo_count }} 张
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="view_count" label="浏览" width="100" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'">
            {{ row.status === 1 ? '已发布' : '草稿' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="推荐" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="row.is_recommend ? 'warning' : ''">
            {{ row.is_recommend ? '★ 推荐' : '-' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="250" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="showDialog(row)">编辑</el-button>
          <el-button 
            size="small" 
            :type="row.is_recommend ? 'warning' : 'default'"
            @click="handleToggleRecommend(row)"
          >
            {{ row.is_recommend ? '取消推荐' : '推荐' }}
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

        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="作品描述" />
        </el-form-item>

        <!-- 开源项目特有字段 -->
        <template v-if="form.type === 'project'">
          <el-form-item label="封面图" prop="cover">
            <el-input v-model="form.cover" placeholder="封面图URL" />
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
              <el-icon><Plus /></el-icon> 添加照片
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
                  <el-form-item label="照片URL" required>
                    <el-input v-model="photo.url" placeholder="照片URL" />
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
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import adminApi from '@/utils/adminApi'

const works = ref([])
const dialogVisible = ref(false)
const formRef = ref()
const isEdit = ref(false)
const activePhotoIndex = ref(0)

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

const loadWorks = async () => {
  try {
    const response = await adminApi.get('/admin/works')
    works.value = response.data.list || []
  } catch (error) {
    ElMessage.error('加载失败')
  }
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
    ElMessage.warning('照片数量已达上限（50张）')
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
          ElMessage.warning('请至少添加1张照片')
          return
        }
        
        // 验证所有照片都有URL
        const hasEmptyUrl = photos.value.some(p => !p.url)
        if (hasEmptyUrl) {
          ElMessage.warning('请填写所有照片的URL')
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
        ElMessage.success('更新成功')
      } else {
        await adminApi.post('/admin/works', submitData)
        ElMessage.success('创建成功')
      }
      
      dialogVisible.value = false
      loadWorks()
    } catch (error) {
      ElMessage.error(error.response?.data?.message || '操作失败')
    }
  })
}

const handleToggleRecommend = async (work) => {
  try {
    await adminApi.put(`/admin/works/${work.id}/recommend`, {
      is_recommend: !work.is_recommend
    })
    ElMessage.success('设置成功')
    loadWorks()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleDelete = async (work) => {
  try {
    await ElMessageBox.confirm('确定要删除这个作品吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await adminApi.delete(`/admin/works/${work.id}`)
    ElMessage.success('删除成功')
    loadWorks()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

loadWorks()
</script>

<style scoped>
.works {
  padding: 20px;
}

.photos-list {
  margin-top: 20px;
}

.photo-header {
  display: flex;
  align-items: center;
  width: 100%;
}
</style>
