<template>
  <section v-loading="initialLoading" class="doc-edit-page">
    <header class="editor-header">
      <div class="header-left">
        <el-button text circle aria-label="返回空间" @click="goBack"><el-icon><ArrowLeft /></el-icon></el-button>
        <div class="doc-location">
          <span>{{ workspaceName }}</span>
          <small>{{ saveStateText }}</small>
        </div>
      </div>
      <div class="header-actions">
        <el-tag :type="form.status === 1 ? 'success' : 'info'">{{ form.status === 1 ? '已发布' : '草稿' }}</el-tag>
        <el-button @click="openVersions"><el-icon><Clock /></el-icon><span class="button-text">版本</span></el-button>
        <el-button @click="openShares"><el-icon><Share /></el-icon><span class="button-text">分享</span></el-button>
        <el-button :loading="saving" @click="manualSave">保存</el-button>
        <el-button v-if="form.status === 1" :loading="publishing" @click="changeStatus(0)">转为草稿</el-button>
        <el-button v-else type="primary" :loading="publishing" @click="changeStatus(1)">发布</el-button>
      </div>
    </header>

    <div class="editor-surface">
      <el-input
        v-model="form.title"
        class="title-input"
        maxlength="200"
        placeholder="无标题文档"
        @input="markTitleDirty"
      />
      <VditorEditor v-if="editorReady" v-model="form.content" height="calc(100vh - 245px)" @update:model-value="markContentDirty" />
    </div>

    <el-drawer v-model="versionDrawer" title="版本历史" size="min(520px, 92vw)" @open="fetchVersions">
      <div v-loading="versionsLoading" class="version-list">
        <article
          v-for="version in versions"
          :key="version.id || version.version"
          :class="['version-item', { selected: selectedVersion?.version === version.version }]"
          @click="fetchVersionDetail(version.version)"
        >
          <div>
            <strong>v{{ version.version }}</strong>
            <el-tag v-if="version.remark" size="small" effect="plain">{{ version.remark }}</el-tag>
          </div>
          <span>{{ formatDate(version.created_at) }}</span>
        </article>
        <el-empty v-if="!versionsLoading && versions.length === 0" description="还没有版本快照" />
      </div>
      <section v-if="selectedVersion" class="version-detail">
        <div class="version-detail-header">
          <div><strong>v{{ selectedVersion.version }} · {{ selectedVersion.title || '无标题文档' }}</strong></div>
          <el-button type="primary" plain :loading="rollingBack" @click="rollbackVersion(selectedVersion.version)">回滚到此版本</el-button>
        </div>
        <pre>{{ selectedVersion.content || '此版本没有内容' }}</pre>
      </section>
    </el-drawer>

    <el-dialog v-model="shareDialog" title="分享文档" width="min(720px, 94vw)" @open="fetchShares">
      <div class="share-create">
        <el-radio-group v-model="shareMode">
          <el-radio-button label="permanent">永久有效</el-radio-button>
          <el-radio-button label="expiry">设置有效期</el-radio-button>
        </el-radio-group>
        <el-date-picker
          v-if="shareMode === 'expiry'"
          v-model="shareExpiry"
          type="datetime"
          placeholder="选择过期时间"
          :disabled-date="disablePastDate"
          style="width: 220px"
        />
        <el-button type="primary" :loading="creatingShare" @click="createShare">创建链接</el-button>
      </div>
      <el-divider />
      <div v-loading="sharesLoading" class="share-list">
        <article v-for="link in shares" :key="link.id" class="share-row">
          <div class="share-main">
            <div class="share-url">
              <el-tag :type="shareStatus(link).type" size="small">{{ shareStatus(link).label }}</el-tag>
              <span>{{ shareUrl(link.token) }}</span>
            </div>
            <small>
              {{ link.expires_at ? `有效期至 ${formatDate(link.expires_at)}` : '永久有效' }}
              · 访问 {{ link.view_count || 0 }} 次
            </small>
          </div>
          <div class="share-actions">
            <el-button text type="primary" @click="copyShare(link.token)">复制</el-button>
            <el-button text @click="toggleShare(link)">{{ link.enabled ? '禁用' : '启用' }}</el-button>
            <el-button text type="danger" @click="deleteShare(link)">删除</el-button>
          </div>
        </article>
        <el-empty v-if="!sharesLoading && shares.length === 0" description="还没有分享链接" :image-size="70" />
      </div>
    </el-dialog>
  </section>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import { onBeforeRouteLeave, useRoute, useRouter } from 'vue-router'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft, Clock, Share } from '@element-plus/icons-vue'
import api from '@/utils/api'
import VditorEditor from '@/components/VditorEditor.vue'

const route = useRoute()
const router = useRouter()
const docId = Number(route.params.id)
const form = reactive({ title: '', content: '', status: 0, workspace_id: null, catalog_id: null })
const workspaceName = ref('知识库文档')
const initialLoading = ref(true)
const editorReady = ref(false)
const saving = ref(false)
const publishing = ref(false)
const dirty = ref(false)
const titleDirty = ref(false)
const contentDirty = ref(false)
const lastSavedAt = ref(null)
const autosaveFailed = ref(false)
const versionDrawer = ref(false)
const versions = ref([])
const versionsLoading = ref(false)
const selectedVersion = ref(null)
const rollingBack = ref(false)
const shareDialog = ref(false)
const shares = ref([])
const sharesLoading = ref(false)
const creatingShare = ref(false)
const shareMode = ref('permanent')
const shareExpiry = ref(null)
let autosaveTimer = null
let allowLeave = false

const saveStateText = computed(() => {
  if (saving.value) return '正在保存...'
  if (autosaveFailed.value) return '自动保存失败，请手动保存'
  if (dirty.value) return '有未保存更改'
  if (lastSavedAt.value) return `已保存 ${dayjs(lastSavedAt.value).format('HH:mm')}`
  return '已加载'
})

const markTitleDirty = () => {
  if (!initialLoading.value) {
    dirty.value = true
    titleDirty.value = true
    autosaveFailed.value = false
  }
}

const markContentDirty = () => {
  if (!initialLoading.value) {
    dirty.value = true
    contentDirty.value = true
    autosaveFailed.value = false
  }
}

async function fetchDoc() {
  const response = await api.get(`/docs/${docId}/edit`)
  const doc = response.data
  Object.assign(form, {
    title: doc.title || '',
    content: doc.content || '',
    status: doc.status || 0,
    workspace_id: doc.workspace_id,
    catalog_id: doc.catalog_id ?? null
  })
  workspaceName.value = doc.workspace?.name || doc.workspace_name || '知识库文档'
  lastSavedAt.value = doc.updated_at || new Date()
  await nextTick()
  editorReady.value = true
}

async function autosave() {
  if (!contentDirty.value || saving.value || publishing.value || !docId) return
  const savedContent = form.content
  try {
    await api.put(`/docs/${docId}/autosave`, { content: savedContent })
    if (form.content === savedContent) contentDirty.value = false
    dirty.value = titleDirty.value
    if (contentDirty.value) dirty.value = true
    autosaveFailed.value = false
    lastSavedAt.value = new Date()
  } catch {
    autosaveFailed.value = true
  }
}

async function manualSave(showMessage = true) {
  const title = form.title.trim()
  if (!title) return ElMessage.warning('请输入文档标题')
  const savedContent = form.content
  saving.value = true
  try {
    const response = await api.put(`/docs/${docId}`, { title, content: savedContent })
    if (form.title.trim() === title) {
      form.title = title
      titleDirty.value = false
    }
    if (response.data?.status !== undefined) form.status = response.data.status
    if (form.content === savedContent) contentDirty.value = false
    dirty.value = titleDirty.value || contentDirty.value
    autosaveFailed.value = false
    lastSavedAt.value = new Date()
    if (showMessage) ElMessage.success('文档已保存，并生成版本快照')
    return true
  } finally {
    saving.value = false
  }
}

async function changeStatus(status) {
  if (dirty.value || !form.title.trim()) {
    const saved = await manualSave(false)
    if (!saved) return
  }
  publishing.value = true
  const publishedTitle = form.title
  const publishedContent = form.content
  try {
    const response = await api.post(`/docs/${docId}/publish`, { status })
    form.status = response.data?.status ?? status
    if (form.title === publishedTitle && form.content === publishedContent) {
      dirty.value = false
      titleDirty.value = false
      contentDirty.value = false
    }
    lastSavedAt.value = new Date()
    ElMessage.success(status === 1 ? '文档已发布' : '文档已转为草稿')
  } finally {
    publishing.value = false
  }
}

function openVersions() {
  versionDrawer.value = true
}

async function fetchVersions() {
  versionsLoading.value = true
  selectedVersion.value = null
  try {
    const response = await api.get(`/docs/${docId}/versions`)
    versions.value = Array.isArray(response.data) ? response.data : response.data?.list || []
  } finally {
    versionsLoading.value = false
  }
}

async function fetchVersionDetail(version) {
  const response = await api.get(`/docs/${docId}/versions/${version}`)
  selectedVersion.value = response.data
}

async function rollbackVersion(version) {
  await ElMessageBox.confirm(`回滚到 v${version} 后，当前内容会被替换并生成新的版本快照。`, '确认回滚', {
    type: 'warning', confirmButtonText: '确认回滚', cancelButtonText: '取消'
  })
  rollingBack.value = true
  try {
    await api.post(`/docs/${docId}/versions/${version}/rollback`)
    editorReady.value = false
    await fetchDoc()
    dirty.value = false
    titleDirty.value = false
    contentDirty.value = false
    selectedVersion.value = null
    await fetchVersions()
    ElMessage.success(`已回滚到 v${version}`)
  } finally {
    rollingBack.value = false
  }
}

function openShares() {
  shareDialog.value = true
}

async function fetchShares() {
  sharesLoading.value = true
  try {
    const response = await api.get(`/docs/${docId}/shares`)
    shares.value = Array.isArray(response.data) ? response.data : response.data?.list || []
  } finally {
    sharesLoading.value = false
  }
}

async function createShare() {
  if (shareMode.value === 'expiry' && !shareExpiry.value) return ElMessage.warning('请选择过期时间')
  if (shareMode.value === 'expiry' && dayjs(shareExpiry.value).isBefore(dayjs())) return ElMessage.warning('过期时间必须晚于现在')
  creatingShare.value = true
  try {
    await api.post(`/docs/${docId}/shares`, {
      permanent: shareMode.value === 'permanent',
      expires_at: shareMode.value === 'permanent' ? null : new Date(shareExpiry.value).toISOString()
    })
    shareExpiry.value = null
    await fetchShares()
    ElMessage.success('分享链接已创建')
  } finally {
    creatingShare.value = false
  }
}

async function toggleShare(link) {
  await api.put(`/shares/${link.id}`, { enabled: !link.enabled })
  await fetchShares()
  ElMessage.success(link.enabled ? '分享链接已禁用' : '分享链接已启用')
}

async function deleteShare(link) {
  await ElMessageBox.confirm('删除后该分享地址将永久失效，是否继续？', '删除分享链接', {
    type: 'warning', confirmButtonText: '确认删除', cancelButtonText: '取消'
  })
  await api.delete(`/shares/${link.id}`)
  await fetchShares()
  ElMessage.success('分享链接已删除')
}

const shareUrl = token => `${window.location.origin}/share/${token}`
async function copyShare(token) {
  try {
    await navigator.clipboard.writeText(shareUrl(token))
    ElMessage.success('链接已复制')
  } catch {
    ElMessage.error('复制失败，请手动选择链接')
  }
}

function shareStatus(link) {
  if (!link.enabled) return { label: '已禁用', type: 'info' }
  if (link.expires_at && dayjs(link.expires_at).isBefore(dayjs())) return { label: '已过期', type: 'danger' }
  return { label: '生效中', type: 'success' }
}

const disablePastDate = date => dayjs(date).endOf('day').isBefore(dayjs())
const formatDate = value => value ? dayjs(value).format('YYYY-MM-DD HH:mm') : '-'
const goBack = () => router.push(form.workspace_id ? `/dashboard/workspaces/${form.workspace_id}` : '/dashboard/workspaces')

function beforeUnload(event) {
  if (!dirty.value) return
  event.preventDefault()
  event.returnValue = ''
}

onBeforeRouteLeave(async () => {
  if (!dirty.value || allowLeave) return true
  try {
    await ElMessageBox.confirm('文档还有未保存更改，确定离开吗？', '离开编辑页', {
      type: 'warning', confirmButtonText: '直接离开', cancelButtonText: '继续编辑'
    })
    allowLeave = true
    return true
  } catch {
    return false
  }
})

onMounted(async () => {
  try {
    await fetchDoc()
    autosaveTimer = window.setInterval(autosave, 30000)
    window.addEventListener('beforeunload', beforeUnload)
  } catch {
    allowLeave = true
    router.replace('/dashboard/workspaces')
  } finally {
    initialLoading.value = false
  }
})

onBeforeUnmount(() => {
  if (autosaveTimer) window.clearInterval(autosaveTimer)
  window.removeEventListener('beforeunload', beforeUnload)
})
</script>

<style scoped>
.doc-edit-page { min-height: calc(100vh - 100px); color: var(--theme-text-primary); }
.editor-header { display: flex; align-items: center; justify-content: space-between; gap: 20px; margin-bottom: 14px; }
.header-left, .header-actions { display: flex; align-items: center; gap: 9px; }
.doc-location { display: flex; flex-direction: column; }
.doc-location span { font-weight: 600; }
.doc-location small { color: var(--theme-text-tertiary); }
.editor-surface { padding: 22px; background: var(--theme-bg-card); border: 1px solid var(--theme-border); border-radius: 13px; box-shadow: 0 8px 30px var(--theme-shadow); }
.title-input { margin-bottom: 16px; }
.title-input :deep(.el-input__wrapper) { padding: 5px 2px; background: transparent; box-shadow: none !important; }
.title-input :deep(.el-input__inner) { height: 48px; background: transparent !important; font-size: clamp(24px, 3vw, 34px); font-weight: 700; }
.version-list { display: flex; flex-direction: column; gap: 8px; }
.version-item { padding: 13px; border: 1px solid var(--theme-border); border-radius: 9px; cursor: pointer; }
.version-item:hover, .version-item.selected { border-color: var(--theme-primary); background: color-mix(in srgb, var(--theme-primary) 7%, var(--theme-bg-card)); }
.version-item > div { display: flex; align-items: center; gap: 8px; }
.version-item > span { display: block; margin-top: 5px; color: var(--theme-text-tertiary); font-size: 12px; }
.version-detail { margin-top: 18px; padding-top: 18px; border-top: 1px solid var(--theme-border); }
.version-detail-header { display: flex; justify-content: space-between; align-items: center; gap: 10px; }
.version-detail pre { max-height: 40vh; padding: 15px; overflow: auto; white-space: pre-wrap; background: var(--theme-bg-secondary); border-radius: 8px; color: var(--theme-text-secondary); font: 13px/1.7 monospace; }
.share-create { display: flex; align-items: center; gap: 12px; flex-wrap: wrap; }
.share-list { min-height: 100px; }
.share-row { display: flex; align-items: center; gap: 12px; padding: 14px 0; border-bottom: 1px solid var(--theme-border-light); }
.share-main { flex: 1; min-width: 0; }
.share-url { display: flex; align-items: center; gap: 8px; }
.share-url span { min-width: 0; overflow: hidden; color: var(--theme-text-secondary); text-overflow: ellipsis; white-space: nowrap; }
.share-main small { display: block; margin-top: 6px; color: var(--theme-text-tertiary); }
.share-actions { display: flex; flex: none; }
@media (max-width: 760px) { .editor-header { align-items: flex-start; flex-direction: column; } .header-actions { width: 100%; overflow-x: auto; padding-bottom: 3px; } .editor-surface { padding: 12px; } .button-text { display: none; } .share-row { align-items: flex-start; flex-direction: column; } .share-actions { width: 100%; justify-content: flex-end; } .version-detail-header { align-items: flex-start; flex-direction: column; } }
</style>
