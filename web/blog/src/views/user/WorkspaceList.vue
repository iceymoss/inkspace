<template>
  <section class="workspace-page">
    <header class="page-header">
      <div>
        <p class="eyebrow">PRIVATE LIBRARY</p>
        <h1>我的知识库</h1>
        <p class="subtitle">把灵感、资料和长期项目整理成独立空间。</p>
      </div>
      <el-button type="primary" size="large" @click="openCreate">
        <el-icon><Plus /></el-icon>新建工作空间
      </el-button>
    </header>

    <div v-loading="store.loading" class="workspace-grid">
      <article
        v-for="workspace in store.workspaces"
        :key="workspace.id"
        class="workspace-card"
        @click="openWorkspace(workspace.id)"
      >
        <div class="card-top">
          <span class="workspace-icon">{{ workspace.icon || '知' }}</span>
          <el-dropdown trigger="click" @command="command => handleCommand(command, workspace)" @click.stop>
            <el-button text circle aria-label="空间操作"><el-icon><MoreFilled /></el-icon></el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="edit">编辑空间</el-dropdown-item>
                <el-dropdown-item command="delete" divided>删除空间</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
        <h2>{{ workspace.name }}</h2>
        <p>{{ workspace.description || '还没有简介，进入空间开始记录吧。' }}</p>
        <footer>
          <span>{{ workspace.doc_count || 0 }} 篇文档</span>
          <span>更新于 {{ formatDate(workspace.updated_at) }}</span>
        </footer>
      </article>
    </div>

    <el-empty v-if="!store.loading && store.workspaceCount === 0" description="还没有工作空间">
      <el-button type="primary" @click="openCreate">创建第一个知识库</el-button>
    </el-empty>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑工作空间' : '新建工作空间'" width="min(500px, 92vw)">
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" maxlength="100" show-word-limit placeholder="例如：产品设计手册" />
        </el-form-item>
        <el-form-item label="简介">
          <el-input v-model="form.description" type="textarea" :rows="3" maxlength="500" show-word-limit />
        </el-form-item>
        <el-form-item label="图标">
          <el-input v-model="form.icon" maxlength="20" placeholder="可输入一个 emoji 或文字" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitWorkspace">保存</el-button>
      </template>
    </el-dialog>
  </section>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'
import { MoreFilled, Plus } from '@element-plus/icons-vue'
import { useWorkspaceStore } from '@/stores/workspace'

const router = useRouter()
const store = useWorkspaceStore()
const dialogVisible = ref(false)
const submitting = ref(false)
const editingId = ref(null)
const formRef = ref(null)
const form = reactive({ name: '', description: '', icon: '', sort: 0 })
const rules = { name: [{ required: true, message: '请输入空间名称', trigger: 'blur' }] }

const formatDate = value => value ? dayjs(value).format('YYYY-MM-DD') : '刚刚'
const openWorkspace = id => router.push(`/dashboard/workspaces/${id}`)

function openCreate() {
  editingId.value = null
  Object.assign(form, { name: '', description: '', icon: '', sort: 0 })
  dialogVisible.value = true
}

function openEdit(workspace) {
  editingId.value = workspace.id
  Object.assign(form, {
    name: workspace.name,
    description: workspace.description || '',
    icon: workspace.icon || '',
    sort: workspace.sort || 0
  })
  dialogVisible.value = true
}

async function submitWorkspace() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    const payload = { name: form.name.trim(), description: form.description, icon: form.icon, sort: form.sort }
    if (editingId.value) await store.updateWorkspace(editingId.value, payload)
    else await store.createWorkspace(payload)
    ElMessage.success(editingId.value ? '空间已更新' : '空间已创建')
    dialogVisible.value = false
    await store.fetchWorkspaces()
  } finally {
    submitting.value = false
  }
}

async function handleCommand(command, workspace) {
  if (command === 'edit') return openEdit(workspace)
  await ElMessageBox.confirm(
    `删除“${workspace.name}”后，其中的目录和文档也会被删除。此操作不可恢复，是否继续？`,
    '删除工作空间',
    { type: 'warning', confirmButtonText: '确认删除', cancelButtonText: '取消' }
  )
  await store.deleteWorkspace(workspace.id)
  ElMessage.success('工作空间已删除')
}

onMounted(() => store.fetchWorkspaces())
</script>

<style scoped>
.workspace-page { max-width: 1180px; margin: 0 auto; color: var(--theme-text-primary); }
.page-header { display: flex; align-items: flex-end; justify-content: space-between; gap: 24px; margin: 8px 0 30px; }
.eyebrow { margin: 0 0 4px; color: var(--theme-primary); font-size: 12px; font-weight: 700; letter-spacing: 2px; }
h1 { margin: 0; font-size: clamp(28px, 4vw, 40px); }
.subtitle { margin: 6px 0 0; color: var(--theme-text-tertiary); }
.workspace-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 18px; min-height: 120px; }
.workspace-card { min-height: 210px; padding: 22px; background: var(--theme-bg-card); border: 1px solid var(--theme-border); border-radius: 14px; box-shadow: 0 8px 30px var(--theme-shadow); cursor: pointer; transition: transform .2s, border-color .2s; }
.workspace-card:hover { transform: translateY(-3px); border-color: var(--theme-primary); }
.card-top { display: flex; justify-content: space-between; align-items: center; }
.workspace-icon { display: grid; place-items: center; width: 48px; height: 48px; border-radius: 13px; background: color-mix(in srgb, var(--theme-primary) 14%, var(--theme-bg-card)); color: var(--theme-primary); font-size: 22px; font-weight: 700; }
.workspace-card h2 { margin: 20px 0 7px; font-size: 20px; }
.workspace-card p { height: 48px; margin: 0; color: var(--theme-text-tertiary); line-height: 1.7; overflow: hidden; }
.workspace-card footer { display: flex; justify-content: space-between; margin-top: 22px; padding-top: 14px; border-top: 1px solid var(--theme-border-light); color: var(--theme-text-tertiary); font-size: 12px; }
@media (max-width: 640px) { .page-header { align-items: stretch; flex-direction: column; } .page-header .el-button { width: 100%; } .workspace-grid { grid-template-columns: 1fr; } }
</style>
