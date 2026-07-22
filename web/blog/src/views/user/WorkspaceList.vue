<template>
  <section class="workspace-page">
    <header class="page-header">
      <div>
        <p class="eyebrow">
          个人知识空间
        </p>
        <h1>我的知识库</h1>
        <p class="subtitle">
          把灵感、资料和长期项目整理成独立空间。
        </p>
      </div>
      <el-button
        type="primary"
        size="large"
        @click="openCreate"
      >
        <el-icon><Plus /></el-icon>新建工作空间
      </el-button>
    </header>

    <div
      v-loading="store.loading"
      class="workspace-grid"
    >
      <article
        v-for="workspace in store.workspaces"
        :key="workspace.id"
        class="workspace-card"
        role="link"
        tabindex="0"
        @click="openWorkspace(workspace.id)"
        @keyup.enter="openWorkspace(workspace.id)"
        @keyup.space.prevent="openWorkspace(workspace.id)"
      >
        <div class="workspace-cover">
          <img
            v-if="isCoverImage(workspace.icon)"
            :src="workspace.icon"
            :alt="`${workspace.name}封面`"
          >
          <div
            v-else
            class="cover-placeholder"
            aria-hidden="true"
          >
            <span>{{ workspace.icon || workspace.name?.slice(0, 1) || '知' }}</span>
            <small>INKSPACE</small>
          </div>
        </div>
        <div class="card-top">
          <div class="card-title">
            <h2>{{ workspace.name }}</h2>
            <el-tag
              :type="workspace.is_public ? 'success' : 'info'"
              size="small"
              effect="plain"
            >
              {{ workspace.is_public ? '公开 Wiki' : '私有空间' }}
            </el-tag>
          </div>
          <span
            class="workspace-actions"
            @click.stop
            @keydown.stop
          >
            <el-dropdown
              trigger="click"
              @command="command => handleCommand(command, workspace)"
            >
              <el-button
                text
                circle
                aria-label="空间操作"
              ><el-icon><MoreFilled /></el-icon></el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="edit">编辑空间</el-dropdown-item>
                  <el-dropdown-item
                    command="delete"
                    divided
                  >删除空间</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </span>
        </div>
        <p>{{ workspace.description || '还没有简介，进入空间开始记录吧。' }}</p>
        <footer>
          <span class="doc-count"><strong>{{ workspace.doc_count || 0 }}</strong> 篇文档</span>
          <span>更新于 {{ formatDate(workspace.updated_at) }}</span>
          <el-link
            v-if="workspace.is_public"
            class="wiki-link"
            type="primary"
            :underline="false"
            @click.stop="openWiki(workspace.id)"
          >
            查看 Wiki
          </el-link>
          <el-icon class="enter-icon">
            <ArrowRight />
          </el-icon>
        </footer>
      </article>
    </div>

    <el-empty
      v-if="!store.loading && store.workspaceCount === 0"
      description="还没有工作空间"
    >
      <el-button
        type="primary"
        @click="openCreate"
      >
        创建第一个知识库
      </el-button>
    </el-empty>

    <el-dialog
      v-model="dialogVisible"
      :title="editingId ? '编辑工作空间' : '新建工作空间'"
      width="min(500px, 92vw)"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-position="top"
      >
        <el-form-item
          label="名称"
          prop="name"
        >
          <el-input
            v-model="form.name"
            maxlength="100"
            show-word-limit
            placeholder="例如：产品设计手册"
          />
        </el-form-item>
        <el-form-item label="简介">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="空间封面">
          <div class="cover-field">
            <ImageCropUpload
              v-model="form.icon"
              preview-size="100%"
              placeholder="上传封面图片"
              tip="支持最高 50MB 原图，将自动裁剪并压缩后上传"
              :max-size="5"
              :source-max-size="50"
              :aspect-ratio="16 / 9"
            />
            <el-button
              v-if="form.icon"
              text
              type="danger"
              @click="form.icon = ''"
            >
              移除封面
            </el-button>
          </div>
        </el-form-item>
        <el-form-item label="公开 Wiki">
          <div class="visibility-field">
            <el-switch
              v-model="form.is_public"
              active-text="公开"
              inactive-text="私有"
            />
            <el-alert
              v-if="form.is_public"
              title="公开后，任何人都可以访问此空间的 Wiki，并查看其中已发布的文档；草稿不会公开。"
              type="warning"
              :closable="false"
              show-icon
            />
            <p v-else>
              私有空间仅你自己可以管理，文档不会出现在公开 Wiki 中。
            </p>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">
          取消
        </el-button>
        <el-button
          type="primary"
          :loading="submitting"
          @click="submitWorkspace"
        >
          保存
        </el-button>
      </template>
    </el-dialog>
  </section>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowRight, MoreFilled, Plus } from '@element-plus/icons-vue'
import { useWorkspaceStore } from '@/stores/workspace'
import ImageCropUpload from '@/components/ImageCropUpload.vue'

const router = useRouter()
const store = useWorkspaceStore()
const dialogVisible = ref(false)
const submitting = ref(false)
const editingId = ref(null)
const formRef = ref(null)
const form = reactive({ name: '', description: '', icon: '', sort: 0, is_public: false })
const rules = { name: [{ required: true, message: '请输入空间名称', trigger: 'blur' }] }

const formatDate = value => value ? dayjs(value).format('YYYY-MM-DD') : '刚刚'
const isCoverImage = value => /^(https?:\/\/|\/uploads\/)/.test(value || '')
const openWorkspace = id => router.push(`/dashboard/workspaces/${id}`)
const openWiki = id => router.push(`/wiki/${id}`)

function openCreate() {
  editingId.value = null
  Object.assign(form, { name: '', description: '', icon: '', sort: 0, is_public: false })
  dialogVisible.value = true
}

function openEdit(workspace) {
  editingId.value = workspace.id
  Object.assign(form, {
    name: workspace.name,
    description: workspace.description || '',
    icon: isCoverImage(workspace.icon) ? workspace.icon : '',
    sort: workspace.sort || 0,
    is_public: Boolean(workspace.is_public)
  })
  dialogVisible.value = true
}

async function submitWorkspace() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    const payload = {
      name: form.name.trim(),
      description: form.description,
      icon: form.icon,
      sort: form.sort,
      is_public: form.is_public
    }
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
  try {
    await ElMessageBox.confirm(
      `删除“${workspace.name}”后，其中的目录和文档也会被删除。此操作不可恢复，是否继续？`,
      '删除工作空间',
      { type: 'warning', confirmButtonText: '确认删除', cancelButtonText: '取消' }
    )
  } catch {
    return
  }
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
.workspace-card { position: relative; min-height: 320px; padding: 0 26px 26px; overflow: hidden; background: var(--theme-bg-card); border: 1px solid var(--theme-border); border-radius: 0; box-shadow: none; cursor: pointer; transition: transform .25s, border-color .25s; }
.workspace-card::before { position: absolute; inset: 0 0 auto; height: 3px; background: var(--theme-primary); content: ''; opacity: 0; transform: scaleX(.45); transition: opacity .2s, transform .2s; }
.workspace-card:hover, .workspace-card:focus-visible { transform: translateY(-4px); border-color: var(--theme-primary); box-shadow: none; outline: none; }
.workspace-card:hover::before, .workspace-card:focus-visible::before { opacity: 1; transform: scaleX(1); }
.workspace-cover { height: 148px; margin: 0 -26px; overflow: hidden; background: color-mix(in srgb, var(--theme-primary) 8%, var(--theme-bg-secondary)); }
.workspace-cover img { width: 100%; height: 100%; object-fit: cover; transition: transform .35s ease; }
.workspace-card:hover .workspace-cover img { transform: scale(1.035); }
.cover-placeholder { display: flex; align-items: flex-end; justify-content: space-between; height: 100%; padding: 20px 22px; background: linear-gradient(135deg, color-mix(in srgb, var(--theme-primary) 82%, #18243a), color-mix(in srgb, var(--theme-primary) 30%, var(--theme-bg-secondary))); color: #fff; }
.cover-placeholder span { font-size: 44px; font-weight: 700; line-height: 1; opacity: .94; }
.cover-placeholder small { font-size: 10px; font-weight: 700; letter-spacing: 2px; opacity: .7; }
.card-top { display: flex; justify-content: space-between; align-items: center; gap: 12px; padding-top: 18px; }
.card-title { display: flex; align-items: center; gap: 8px; min-width: 0; }
.workspace-actions { flex: none; }
.workspace-card h2 { min-width: 0; margin: 0; overflow: hidden; font-size: 20px; text-overflow: ellipsis; white-space: nowrap; }
.workspace-card p { height: 48px; margin: 7px 0 0; color: var(--theme-text-tertiary); line-height: 1.7; overflow: hidden; }
.workspace-card footer { display: flex; align-items: center; gap: 12px; margin-top: 22px; padding-top: 14px; border-top: 1px solid var(--theme-border-light); color: var(--theme-text-tertiary); font-size: 12px; }
.workspace-card footer > span:nth-child(2) { margin-left: auto; }
.wiki-link { flex: none; font-size: 12px; }
.doc-count strong { color: var(--theme-text-primary); font-size: 15px; }
.enter-icon { color: var(--theme-primary); opacity: 0; transform: translateX(-5px); transition: opacity .2s, transform .2s; }
.workspace-card:hover .enter-icon, .workspace-card:focus-visible .enter-icon { opacity: 1; transform: translateX(0); }
.cover-field { width: 100%; }
.cover-field > .el-button { margin-top: 4px; padding-left: 0; }
.cover-field :deep(.image-crop-upload), .cover-field :deep(.upload-area) { width: 100%; }
.cover-field :deep(.upload-area) { min-height: 150px; border-color: var(--theme-border); border-radius: 10px; }
.cover-field :deep(.preview-image), .cover-field :deep(.preview-image img) { width: 100%; height: 150px; max-height: 150px !important; object-fit: cover; }
.cover-field :deep(.upload-placeholder) { width: 100% !important; min-height: 150px !important; background: var(--theme-bg-secondary); }
.visibility-field { display: flex; flex-direction: column; gap: 10px; width: 100%; }
.visibility-field p { margin: 0; color: var(--theme-text-tertiary); font-size: 13px; line-height: 1.6; }
@media (max-width: 640px) { .page-header { align-items: stretch; flex-direction: column; } .page-header .el-button { width: 100%; } .workspace-grid { grid-template-columns: 1fr; } }
</style>
