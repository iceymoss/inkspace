<template>
  <div class="ads-management">
    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <!-- 广告位置管理 -->
      <el-tab-pane label="广告位置" name="positions">
        <div class="toolbar">
          <el-button type="primary" @click="handleCreatePosition">
            <el-icon><Plus /></el-icon>
            新增位置
          </el-button>
        </div>
        <el-table :data="positions" v-loading="loading">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="name" label="位置名称" />
          <el-table-column prop="code" label="位置代码" />
          <el-table-column prop="max_count" label="最大数量" width="100" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'">
                {{ row.status === 1 ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="sort" label="排序" width="100" />
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link @click="handleEditPosition(row)">编辑</el-button>
              <el-button type="danger" link @click="handleDeletePosition(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
        <el-pagination
          v-model:current-page="positionPage"
          v-model:page-size="positionPageSize"
          :total="positionTotal"
          layout="total, prev, pager, next"
          @current-change="loadPositions"
        />
      </el-tab-pane>

      <!-- 广告信息管理 -->
      <el-tab-pane label="广告信息" name="advertisements">
        <div class="toolbar">
          <el-button type="primary" @click="handleCreateAdvertisement">
            <el-icon><Plus /></el-icon>
            新增广告
          </el-button>
        </div>
        <el-table :data="advertisements" v-loading="loading">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="title" label="标题" />
          <el-table-column label="图片" width="120">
            <template #default="{ row }">
              <el-image
                v-if="row.image"
                :src="row.image"
                style="width: 80px; height: 60px"
                fit="cover"
                :preview-src-list="[row.image]"
              />
            </template>
          </el-table-column>
          <el-table-column prop="link" label="链接" show-overflow-tooltip />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'">
                {{ row.status === 1 ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="click_count" label="点击" width="80" />
          <el-table-column prop="view_count" label="展示" width="80" />
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link @click="handleEditAdvertisement(row)">编辑</el-button>
              <el-button type="danger" link @click="handleDeleteAdvertisement(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
        <el-pagination
          v-model:current-page="adPage"
          v-model:page-size="adPageSize"
          :total="adTotal"
          layout="total, prev, pager, next"
          @current-change="loadAdvertisements"
        />
      </el-tab-pane>

      <!-- 广告投放管理 -->
      <el-tab-pane label="广告投放" name="placements">
        <div class="toolbar">
          <el-button type="primary" @click="handleCreatePlacement">
            <el-icon><Plus /></el-icon>
            新增投放
          </el-button>
          <el-select
            v-model="selectedPositionId"
            placeholder="筛选位置"
            clearable
            style="width: 200px; margin-left: 10px"
            @change="loadPlacements"
          >
            <el-option
              v-for="pos in positions"
              :key="pos.id"
              :label="pos.name"
              :value="pos.id"
            />
          </el-select>
        </div>
        <el-table :data="placements" v-loading="loading">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column label="位置" width="150">
            <template #default="{ row }">
              {{ row.position?.name || '-' }}
            </template>
          </el-table-column>
          <el-table-column label="广告" width="150">
            <template #default="{ row }">
              {{ row.advertisement?.title || '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="start_time" label="开始时间" width="180" />
          <el-table-column prop="end_time" label="结束时间" width="180" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'">
                {{ row.status === 1 ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="sort" label="排序" width="100" />
          <el-table-column prop="click_count" label="点击" width="80" />
          <el-table-column prop="view_count" label="展示" width="80" />
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link @click="handleEditPlacement(row)">编辑</el-button>
              <el-button type="danger" link @click="handleDeletePlacement(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
        <el-pagination
          v-model:current-page="placementPage"
          v-model:page-size="placementPageSize"
          :total="placementTotal"
          layout="total, prev, pager, next"
          @current-change="loadPlacements"
        />
      </el-tab-pane>
    </el-tabs>

    <!-- 广告位置对话框 -->
    <el-dialog
      v-model="positionDialogVisible"
      :title="positionForm.id ? '编辑位置' : '新增位置'"
      width="600px"
    >
      <el-form :model="positionForm" label-width="100px">
        <el-form-item label="位置名称" required>
          <el-input v-model="positionForm.name" />
        </el-form-item>
        <el-form-item label="位置代码" required>
          <el-input v-model="positionForm.code" placeholder="如：blog_sidebar_1" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="positionForm.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="最大投放数量">
          <el-input-number v-model="positionForm.max_count" :min="1" :max="100" />
          <div style="margin-top: 8px; color: #909399; font-size: 12px;">
            该位置最多可以同时投放多少个广告
          </div>
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="positionForm.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="positionForm.sort" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="positionDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSavePosition">保存</el-button>
      </template>
    </el-dialog>

    <!-- 广告信息对话框 -->
    <el-dialog
      v-model="adDialogVisible"
      :title="adForm.id ? '编辑广告信息' : '新增广告信息'"
      width="600px"
    >
      <el-form :model="adForm" label-width="100px">
        <el-form-item label="标题" required>
          <el-input v-model="adForm.title" />
        </el-form-item>
        <el-form-item label="图片">
          <el-input v-model="adForm.image" placeholder="图片URL" />
          <el-button
            type="primary"
            style="margin-top: 10px"
            @click="handleUploadImage('ad')"
          >
            上传图片
          </el-button>
          <el-image
            v-if="adForm.image"
            :src="adForm.image"
            style="width: 200px; height: 150px; margin-top: 10px"
            fit="cover"
          />
        </el-form-item>
        <el-form-item label="链接">
          <el-input v-model="adForm.link" placeholder="广告链接" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="adForm.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="adForm.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="adDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveAdvertisement">保存</el-button>
      </template>
    </el-dialog>

    <!-- 广告投放对话框 -->
    <el-dialog
      v-model="placementDialogVisible"
      :title="placementForm.id ? '编辑投放' : '新增投放'"
      width="600px"
    >
      <el-form :model="placementForm" label-width="100px">
        <el-form-item label="广告位置" required>
          <el-select v-model="placementForm.position_id" style="width: 100%">
            <el-option
              v-for="pos in positions"
              :key="pos.id"
              :label="pos.name"
              :value="pos.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="广告信息" required>
          <el-select v-model="placementForm.advertisement_id" style="width: 100%">
            <el-option
              v-for="ad in advertisements"
              :key="ad.id"
              :label="ad.title"
              :value="ad.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="开始时间">
          <el-date-picker
            v-model="placementForm.start_time"
            type="datetime"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="结束时间">
          <el-date-picker
            v-model="placementForm.end_time"
            type="datetime"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="placementForm.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="placementForm.sort" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="placementDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSavePlacement">保存</el-button>
      </template>
    </el-dialog>

    <!-- 图片上传对话框 -->
    <el-dialog v-model="uploadDialogVisible" title="上传图片" width="500px">
      <el-upload
        class="upload-demo"
        :action="uploadAction"
        :headers="uploadHeaders"
        :on-success="handleUploadSuccess"
        :before-upload="beforeUpload"
        drag
      >
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
        <template #tip>
          <div class="el-upload__tip">只能上传jpg/png文件，且不超过2MB</div>
        </template>
      </el-upload>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, UploadFilled } from '@element-plus/icons-vue'
import adminApi from '@/utils/adminApi'

const activeTab = ref('positions')
const loading = ref(false)

// 广告位置
const positions = ref([])
const positionPage = ref(1)
const positionPageSize = ref(20)
const positionTotal = ref(0)
const positionDialogVisible = ref(false)
const positionForm = ref({
  id: null,
  name: '',
  code: '',
  description: '',
  max_count: 4,
  status: 1,
  sort: 0
})

// 广告信息
const advertisements = ref([])
const adPage = ref(1)
const adPageSize = ref(20)
const adTotal = ref(0)
const adDialogVisible = ref(false)
const adForm = ref({
  id: null,
  title: '',
  image: '',
  link: '',
  description: '',
  status: 1
})

// 广告投放
const placements = ref([])
const placementPage = ref(1)
const placementPageSize = ref(20)
const placementTotal = ref(0)
const selectedPositionId = ref(null)
const placementDialogVisible = ref(false)
const placementForm = ref({
  id: null,
  position_id: null,
  advertisement_id: null,
  start_time: null,
  end_time: null,
  status: 1,
  sort: 0
})

// 图片上传
const uploadDialogVisible = ref(false)
const uploadType = ref('')
const uploadAction = ref('/api/upload/image')
const uploadHeaders = ref({})

// 加载数据
const loadPositions = async () => {
  try {
    loading.value = true
    const res = await adminApi.get('/admin/ad-positions', {
      params: { page: positionPage.value, page_size: positionPageSize.value }
    })
    positions.value = res.data.list || []
    positionTotal.value = res.data.total || 0
  } catch (error) {
    ElMessage.error('加载位置列表失败')
  } finally {
    loading.value = false
  }
}

const loadAdvertisements = async () => {
  try {
    loading.value = true
    const res = await adminApi.get('/admin/advertisements', {
      params: { page: adPage.value, page_size: adPageSize.value }
    })
    advertisements.value = res.data.list || []
    adTotal.value = res.data.total || 0
  } catch (error) {
    ElMessage.error('加载广告信息列表失败')
  } finally {
    loading.value = false
  }
}

const loadPlacements = async () => {
  try {
    loading.value = true
    const params = { page: placementPage.value, page_size: placementPageSize.value }
    if (selectedPositionId.value) {
      params.position_id = selectedPositionId.value
    }
    const res = await adminApi.get('/admin/ad-placements', { params })
    placements.value = res.data.list || []
    placementTotal.value = res.data.total || 0
  } catch (error) {
    ElMessage.error('加载投放列表失败')
  } finally {
    loading.value = false
  }
}

// 位置管理
const handleCreatePosition = () => {
  positionForm.value = {
    id: null,
    name: '',
    code: '',
    description: '',
    max_count: 4,
    status: 1,
    sort: 0
  }
  positionDialogVisible.value = true
}

const handleEditPosition = (row) => {
  positionForm.value = { ...row }
  positionDialogVisible.value = true
}

const handleSavePosition = async () => {
  try {
    if (positionForm.value.id) {
      await adminApi.put(`/admin/ad-positions/${positionForm.value.id}`, positionForm.value)
      ElMessage.success('更新成功')
    } else {
      await adminApi.post('/admin/ad-positions', positionForm.value)
      ElMessage.success('创建成功')
    }
    positionDialogVisible.value = false
    loadPositions()
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '操作失败')
  }
}

const handleDeletePosition = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除这个位置吗？', '提示', {
      type: 'warning'
    })
    await adminApi.delete(`/admin/ad-positions/${row.id}`)
    ElMessage.success('删除成功')
    loadPositions()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 广告管理
const handleCreateAdvertisement = () => {
  adForm.value = {
    id: null,
    title: '',
    image: '',
    link: '',
    description: '',
    status: 1
  }
  adDialogVisible.value = true
}

const handleEditAdvertisement = (row) => {
  adForm.value = { ...row }
  adDialogVisible.value = true
}

const handleSaveAdvertisement = async () => {
  try {
    if (adForm.value.id) {
      await adminApi.put(`/admin/advertisements/${adForm.value.id}`, adForm.value)
      ElMessage.success('更新成功')
    } else {
      await adminApi.post('/admin/advertisements', adForm.value)
      ElMessage.success('创建成功')
    }
    adDialogVisible.value = false
    loadAdvertisements()
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '操作失败')
  }
}

const handleDeleteAdvertisement = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除这个广告吗？', '提示', {
      type: 'warning'
    })
    await adminApi.delete(`/admin/advertisements/${row.id}`)
    ElMessage.success('删除成功')
    loadAdvertisements()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 投放管理
const handleCreatePlacement = () => {
  placementForm.value = {
    id: null,
    position_id: null,
    advertisement_id: null,
    start_time: null,
    end_time: null,
    status: 1,
    sort: 0
  }
  placementDialogVisible.value = true
}

const handleEditPlacement = (row) => {
  placementForm.value = {
    ...row,
    start_time: row.start_time || null,
    end_time: row.end_time || null
  }
  placementDialogVisible.value = true
}

const handleSavePlacement = async () => {
  try {
    const data = {
      ...placementForm.value,
      start_time: placementForm.value.start_time || null,
      end_time: placementForm.value.end_time || null
    }
    if (placementForm.value.id) {
      await adminApi.put(`/admin/ad-placements/${placementForm.value.id}`, data)
      ElMessage.success('更新成功')
    } else {
      await adminApi.post('/admin/ad-placements', data)
      ElMessage.success('创建成功')
    }
    placementDialogVisible.value = false
    loadPlacements()
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '操作失败')
  }
}

const handleDeletePlacement = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除这个投放吗？', '提示', {
      type: 'warning'
    })
    await adminApi.delete(`/admin/ad-placements/${row.id}`)
    ElMessage.success('删除成功')
    loadPlacements()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 图片上传
const handleUploadImage = (type) => {
  uploadType.value = type
  uploadHeaders.value = {
        Authorization: `Bearer ${localStorage.getItem('admin_token')}`
      }
  uploadDialogVisible.value = true
}

const beforeUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt2M = file.size / 1024 / 1024 < 2
  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过2MB!')
    return false
  }
  return true
}

const handleUploadSuccess = (response) => {
  if (response.code === 0 && response.data) {
    if (uploadType.value === 'ad') {
      adForm.value.image = response.data.url
    }
    ElMessage.success('上传成功')
    uploadDialogVisible.value = false
  } else {
    ElMessage.error('上传失败')
  }
}

const handleTabChange = (tab) => {
  if (tab === 'positions') {
    loadPositions()
  } else if (tab === 'advertisements') {
    loadAdvertisements()
  } else if (tab === 'placements') {
    loadPlacements()
  }
}

onMounted(() => {
  loadPositions()
  loadAdvertisements()
  loadPlacements()
})
</script>

<style scoped>
.ads-management {
  padding: 20px;
}

.toolbar {
  margin-bottom: 20px;
}

.el-pagination {
  margin-top: 20px;
  justify-content: flex-end;
}
</style>

