<template>
  <div class="settings">
    <h2>系统配置</h2>

    <el-tabs v-model="activeGroup">
      <el-tab-pane label="网站信息" name="site">
        <el-form :model="siteSettings" label-width="120px">
          <el-form-item label="网站名称">
            <el-input v-model="siteSettings.site_name" />
          </el-form-item>
          <el-form-item label="网站描述">
            <el-input v-model="siteSettings.site_description" type="textarea" :rows="3" />
          </el-form-item>
          <el-form-item label="关键词">
            <el-input v-model="siteSettings.site_keywords" placeholder="关键词用逗号分隔" />
          </el-form-item>
          <el-form-item label="备案号">
            <el-input v-model="siteSettings.site_icp" />
          </el-form-item>
          <el-form-item label="版权信息">
            <el-input v-model="siteSettings.site_copyright" />
          </el-form-item>
          <el-form-item label="网站Logo">
            <el-input v-model="siteSettings.site_logo" placeholder="Logo URL" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="saveSiteSettings" :loading="saving">保存</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="功能设置" name="feature">
        <el-form :model="featureSettings" label-width="120px">
          <el-form-item label="开放注册">
            <el-switch v-model="featureSettings.register_enabled" />
          </el-form-item>
          <el-form-item label="开放评论">
            <el-switch v-model="featureSettings.comment_enabled" />
          </el-form-item>
          <el-form-item label="评论审核">
            <el-switch v-model="featureSettings.comment_audit" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="saveFeatureSettings" :loading="saving">保存</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="所有配置" name="all">
        <el-button type="primary" @click="showEditDialog()" style="margin-bottom: 20px;">
          <el-icon><Plus /></el-icon> 新建配置
        </el-button>

        <el-table :data="allSettings" style="width: 100%;">
          <el-table-column prop="key" label="配置键" width="200" />
          <el-table-column prop="value" label="配置值" show-overflow-tooltip />
          <el-table-column prop="type" label="类型" width="100" />
          <el-table-column prop="group" label="分组" width="100" />
          <el-table-column label="公开" width="80">
            <template #default="{ row }">
              <el-tag :type="row.is_public ? 'success' : 'info'" size="small">
                {{ row.is_public ? '是' : '否' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="description" label="描述" show-overflow-tooltip />
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="{ row }">
              <el-button size="small" @click="showEditDialog(row)">编辑</el-button>
              <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <el-dialog v-model="editDialogVisible" title="编辑配置" width="600px">
      <el-form :model="editForm" :rules="editRules" ref="editFormRef" label-width="100px">
        <el-form-item label="配置键" prop="key">
          <el-input v-model="editForm.key" :disabled="isEditMode" />
        </el-form-item>
        <el-form-item label="配置值" prop="value">
          <el-input v-model="editForm.value" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="editForm.type">
            <el-option label="字符串" value="string" />
            <el-option label="整数" value="int" />
            <el-option label="布尔值" value="bool" />
            <el-option label="JSON" value="json" />
          </el-select>
        </el-form-item>
        <el-form-item label="分组">
          <el-input v-model="editForm.group" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="editForm.description" />
        </el-form-item>
        <el-form-item label="公开">
          <el-switch v-model="editForm.is_public" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitEdit" :loading="editLoading">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import adminApi from '@/utils/adminApi'

const activeGroup = ref('site')
const allSettings = ref([])
const saving = ref(false)
const editDialogVisible = ref(false)
const editFormRef = ref()
const editLoading = ref(false)
const isEditMode = ref(false)

const siteSettings = reactive({
  site_name: '',
  site_description: '',
  site_keywords: '',
  site_icp: '',
  site_copyright: '',
  site_logo: ''
})

const featureSettings = reactive({
  register_enabled: true,
  comment_enabled: true,
  comment_audit: false
})

const editForm = reactive({
  key: '',
  value: '',
  type: 'string',
  group: 'general',
  description: '',
  is_public: false
})

const editRules = {
  key: [{ required: true, message: '请输入配置键', trigger: 'blur' }]
}

const loadAllSettings = async () => {
  try {
    const response = await adminApi.get('/admin/settings')
    allSettings.value = response.data || []
    
    // 分组加载到不同的设置对象
    allSettings.value.forEach(setting => {
      if (setting.group === 'site') {
        siteSettings[setting.key] = setting.value
      } else if (setting.group === 'feature') {
        featureSettings[setting.key] = setting.value === '1' || setting.value === 'true'
      }
    })
  } catch (error) {
    ElMessage.error('加载失败')
  }
}

const saveSiteSettings = async () => {
  saving.value = true
  try {
    await adminApi.put('/admin/settings/batch', siteSettings)
    ElMessage.success('保存成功')
    loadAllSettings()
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const saveFeatureSettings = async () => {
  saving.value = true
  try {
    // 转换布尔值为字符串
    const settings = {}
    Object.keys(featureSettings).forEach(key => {
      settings[key] = featureSettings[key] ? '1' : '0'
    })
    
    await adminApi.put('/admin/settings/batch', settings)
    ElMessage.success('保存成功')
    loadAllSettings()
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const showEditDialog = (setting = null) => {
  if (setting) {
    isEditMode.value = true
    Object.assign(editForm, setting)
  } else {
    isEditMode.value = false
    Object.assign(editForm, {
      key: '',
      value: '',
      type: 'string',
      group: 'general',
      description: '',
      is_public: false
    })
  }
  editDialogVisible.value = true
}

const submitEdit = async () => {
  await editFormRef.value.validate(async (valid) => {
    if (!valid) return

    editLoading.value = true
    try {
      await adminApi.put('/admin/settings', editForm)
      ElMessage.success('保存成功')
      editDialogVisible.value = false
      loadAllSettings()
    } catch (error) {
      ElMessage.error('保存失败')
    } finally {
      editLoading.value = false
    }
  })
}

const handleDelete = async (setting) => {
  try {
    await ElMessageBox.confirm('确定要删除这个配置吗？', '提示', { type: 'warning' })
    await adminApi.delete(`/admin/settings/${setting.key}`)
    ElMessage.success('删除成功')
    loadAllSettings()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

onMounted(() => {
  loadAllSettings()
})
</script>

