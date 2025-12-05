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
      <el-table-column prop="view_count" label="浏览" width="100" />
      <el-table-column label="类型" width="100">
        <template #default="{ row }">
          <el-tag :type="row.type === 'photography' ? 'warning' : 'primary'">
            {{ row.type === 'photography' ? '📷 摄影' : '💻 项目' }}
          </el-tag>
        </template>
      </el-table-column>
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

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑作品' : '新建作品'" width="700px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="作品类型" prop="type">
          <el-radio-group v-model="form.type" @change="handleTypeChange">
            <el-radio label="project">开源项目</el-radio>
            <el-radio label="photography">摄影作品</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="输入作品标题" />
        </el-form-item>

        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="作品描述" />
        </el-form-item>

        <el-form-item label="封面图" prop="cover">
          <el-input v-model="form.cover" placeholder="封面图URL" />
          <el-text size="small" type="info">建议尺寸: 800x600</el-text>
        </el-form-item>

        <!-- 开源项目特有字段 -->
        <template v-if="form.type === 'project'">
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
          <el-divider content-position="left">摄影参数（选填）</el-divider>
          
          <el-form-item label="相机型号">
            <el-input v-model="photoMetadata.camera" placeholder="例如: Canon EOS R5" />
          </el-form-item>
          
          <el-form-item label="镜头">
            <el-input v-model="photoMetadata.lens" placeholder="例如: RF 24-70mm f/2.8" />
          </el-form-item>
          
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="焦段">
                <el-input v-model="photoMetadata.focal_length" placeholder="例如: 50mm" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="光圈">
                <el-input v-model="photoMetadata.aperture" placeholder="例如: f/2.8" />
              </el-form-item>
            </el-col>
          </el-row>
          
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="快门速度">
                <el-input v-model="photoMetadata.shutter_speed" placeholder="例如: 1/200s" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="ISO">
                <el-input v-model="photoMetadata.iso" placeholder="例如: 400" />
              </el-form-item>
            </el-col>
          </el-row>
          
          <el-form-item label="拍摄地点">
            <el-input v-model="photoMetadata.location" placeholder="例如: 杭州西湖" />
          </el-form-item>
          
          <el-form-item label="拍摄日期">
            <el-date-picker 
              v-model="photoMetadata.shooting_date" 
              type="date"
              placeholder="选择拍摄日期"
              value-format="YYYY-MM-DD"
              style="width: 100%"
            />
          </el-form-item>

          <el-alert 
            title="摄影作品说明" 
            type="info" 
            :closable="false"
            style="margin-bottom: 15px"
          >
            • 每天最多发布3张摄影作品<br>
            • 图片将保留原图质量，不会压缩<br>
            • 建议上传高质量JPG或PNG格式
          </el-alert>
        </template>

        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">发布</el-radio>
            <el-radio :label="0">草稿</el-radio>
          </el-radio-group>
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

const works = ref([])
const dialogVisible = ref(false)
const formRef = ref()
const loading = ref(false)
const editingId = ref(null)

const isEdit = computed(() => !!editingId.value)

const form = reactive({
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
  images: []
})

// 摄影作品元数据
const photoMetadata = reactive({
  camera: '',
  lens: '',
  focal_length: '',
  aperture: '',
  shutter_speed: '',
  iso: '',
  location: '',
  shooting_date: ''
})

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }]
}

const loadWorks = async () => {
  try {
    const response = await adminApi.get('/admin/works', {
      params: { page: 1, page_size: 100 }
    })
    works.value = response.data.list || []
  } catch (error) {
    ElMessage.error('加载失败')
  }
}

const handleTypeChange = () => {
  // 切换类型时清空对应的字段
  if (form.type === 'photography') {
    form.link = ''
    form.github_url = ''
    form.demo_url = ''
    form.tech_stack = ''
  } else if (form.type === 'project') {
    Object.assign(photoMetadata, {
      camera: '',
      lens: '',
      focal_length: '',
      aperture: '',
      shutter_speed: '',
      iso: '',
      location: '',
      shooting_date: ''
    })
  }
}

const showDialog = (work = null) => {
  if (work) {
    editingId.value = work.id
    Object.assign(form, work)
    
    // 加载摄影元数据
    if (work.metadata && work.type === 'photography') {
      Object.assign(photoMetadata, work.metadata)
    }
  } else {
    editingId.value = null
    Object.assign(form, {
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
      images: []
    })
    Object.assign(photoMetadata, {
      camera: '',
      lens: '',
      focal_length: '',
      aperture: '',
      shutter_speed: '',
      iso: '',
      location: '',
      shooting_date: ''
    })
  }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      // 准备提交数据
      const submitData = { ...form }
      
      // 如果是摄影作品，添加元数据
      if (form.type === 'photography') {
        submitData.metadata = photoMetadata
      } else {
        submitData.metadata = {}
      }

      if (isEdit.value) {
        await adminApi.put(`/admin/works/${editingId.value}`, submitData)
        ElMessage.success('更新成功')
      } else {
        await adminApi.post('/admin/works', submitData)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      loadWorks()
    } catch (error) {
      ElMessage.error('保存失败')
    } finally {
      loading.value = false
    }
  })
}

const handleToggleRecommend = async (work) => {
  try {
    await adminApi.put(`/admin/works/${work.id}/recommend`, {
      is_recommend: !work.is_recommend
    })
    ElMessage.success(work.is_recommend ? '已取消推荐' : '设置推荐成功')
    loadWorks()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleDelete = async (work) => {
  try {
    await ElMessageBox.confirm('确定要删除这个作品吗？', '提示', { type: 'warning' })
    await adminApi.delete(`/admin/works/${work.id}`)
    ElMessage.success('删除成功')
    loadWorks()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

onMounted(() => {
  loadWorks()
})
</script>

