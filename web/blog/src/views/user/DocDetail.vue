<template>
  <section class="doc-detail-page">
    <div
      v-if="loading"
      class="doc-card loading-card"
      aria-busy="true"
      aria-label="正在加载文档"
    >
      <el-skeleton
        :rows="10"
        animated
      />
    </div>

    <el-result
      v-else-if="error"
      icon="error"
      title="文档暂时无法打开"
      :sub-title="error"
    >
      <template #extra>
        <el-button @click="loadDoc">
          重新加载
        </el-button>
        <el-button @click="router.push('/dashboard/workspaces')">
          返回知识库
        </el-button>
      </template>
    </el-result>

    <article
      v-else-if="doc"
      class="doc-card"
    >
      <header class="doc-header">
        <div class="header-topline">
          <el-button
            text
            class="back-button"
            @click="router.push(`/dashboard/workspaces/${doc.workspace_id}`)"
          >
            <el-icon><ArrowLeft /></el-icon>
            返回空间
          </el-button>
          <el-button
            v-if="doc.capabilities?.can_edit"
            type="primary"
            @click="router.push(`/dashboard/docs/${doc.id}/edit`)"
          >
            <el-icon><EditPen /></el-icon>
            编辑文档
          </el-button>
          <el-button
            v-if="doc.capabilities?.can_download"
            :loading="downloading"
            @click="downloadDoc"
          >
            <el-icon><Download /></el-icon>
            下载
          </el-button>
        </div>

        <div class="title-row">
          <div>
            <p class="doc-kicker">
              DOCUMENT / {{ doc.id }}
            </p>
            <h1>{{ doc.title || '无标题文档' }}</h1>
          </div>
          <el-tag
            :type="doc.status === 1 ? 'success' : 'info'"
            effect="plain"
          >
            {{ doc.status === 1 ? '已公开' : '私有' }}
          </el-tag>
        </div>

        <dl class="doc-meta">
          <div>
            <dt>更新时间</dt>
            <dd><time :datetime="doc.updated_at">{{ formatDate(doc.updated_at) }}</time></dd>
          </div>
          <div>
            <dt>文档类型</dt>
            <dd>{{ kindLabel }}</dd>
          </div>
          <div>
            <dt>字数</dt>
            <dd>{{ doc.word_count ?? 0 }} 字</dd>
          </div>
        </dl>
      </header>

      <div class="content-rule" />
      <MarkdownPreview
        v-if="doc.kind === 'markdown' && doc.content_html"
        :html="doc.content_html"
      />
      <StructuredDataPreview
        v-else-if="isStructuredData && doc.content"
        :content="doc.content"
        :language="doc.language"
      />
      <CodeTextEditor
        v-else-if="isPlainText && doc.content"
        :model-value="doc.content"
        :language="doc.language"
        read-only
      />
      <div
        v-else-if="previewLoading"
        class="preview-loading"
      >
        <el-skeleton :rows="8" animated />
      </div>
      <PdfPreview
        v-else-if="previewUrl && doc.attachment?.mime_type === 'application/pdf'"
        :src="previewUrl"
      />
      <ImagePreview
        v-else-if="previewUrl && isRasterImage"
        :src="previewUrl"
        :file-name="doc.attachment.file_name"
      />
      <DocxPreview
        v-else-if="previewBlob && isDocx"
        :blob="previewBlob"
        @error="handlePreviewError"
      />
      <section
        v-else-if="doc.kind === 'file' && doc.attachment"
        class="file-fallback"
      >
        <div class="file-mark">
          {{ fileExtension }}
        </div>
        <div class="file-copy">
          <p class="file-label">
            FILE ATTACHMENT
          </p>
          <h2>{{ doc.attachment.file_name }}</h2>
          <p>{{ formatFileSize(doc.attachment.file_size) }} · {{ doc.attachment.mime_type }}</p>
          <el-tag
            :type="previewTagType"
            effect="plain"
          >
            {{ previewStatusLabel }}
          </el-tag>
        </div>
        <el-button
          type="primary"
          :loading="downloading"
          @click="downloadDoc"
        >
          <el-icon><Download /></el-icon>
          下载原文件
        </el-button>
      </section>
      <el-empty
        v-else
        description="这篇文档暂时没有正文"
        :image-size="72"
      />
    </article>
  </section>
</template>

<script setup>
import { computed, defineAsyncComponent, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import dayjs from 'dayjs'
import { ArrowLeft, Download, EditPen } from '@element-plus/icons-vue'
import api from '@/utils/api'
import MarkdownPreview from '@/components/docs/MarkdownPreview.vue'
import ImagePreview from '@/components/docs/ImagePreview.vue'
import PdfPreview from '@/components/docs/PdfPreview.vue'
import DocxPreview from '@/components/docs/DocxPreview.vue'

const CodeTextEditor = defineAsyncComponent(() => import('@/components/docs/CodeTextEditor.vue'))
const StructuredDataPreview = defineAsyncComponent(() => import('@/components/docs/StructuredDataPreview.vue'))

const route = useRoute()
const router = useRouter()
const doc = ref(null)
const loading = ref(true)
const downloading = ref(false)
const previewLoading = ref(false)
const previewUrl = ref('')
const previewBlob = ref(null)
const error = ref('')

const isPlainText = computed(() => doc.value?.kind === 'text' || doc.value?.kind === 'code')
const isStructuredData = computed(() => ['json', 'jsonl', 'yaml', 'csv'].includes(doc.value?.language))
const isRasterImage = computed(() => ['image/jpeg', 'image/png', 'image/gif', 'image/webp'].includes(doc.value?.attachment?.mime_type))
const isDocx = computed(() => doc.value?.attachment?.mime_type === 'application/vnd.openxmlformats-officedocument.wordprocessingml.document')
const kindLabel = computed(() => {
	const labels = { markdown: 'Markdown', text: '纯文本', code: '代码', file: '文件' }
  const label = labels[doc.value?.kind] || doc.value?.kind || '文档'
  return doc.value?.language ? `${label} · ${doc.value.language}` : label
})
const fileExtension = computed(() => (doc.value?.attachment?.extension || 'FILE').replace('.', '').toUpperCase())
const previewStatusLabel = computed(() => {
  const labels = { ready: '可直接查看', pending: '等待生成安全预览', running: '正在生成安全预览', failed: '预览生成失败', none: '仅支持下载' }
  return labels[doc.value?.attachment?.preview_status] || '暂无在线预览'
})
const previewTagType = computed(() => ({ ready: 'success', failed: 'danger' }[doc.value?.attachment?.preview_status] || 'info'))
const formatDate = value => value ? dayjs(value).format('YYYY-MM-DD HH:mm') : '未记录'
const formatFileSize = size => {
  if (!Number.isFinite(size) || size < 1024) return `${size || 0} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KiB`
  return `${(size / 1024 / 1024).toFixed(1)} MiB`
}

async function downloadDoc() {
  downloading.value = true
  try {
    const response = await api.get(`/docs/${doc.value.id}/download`, { responseType: 'blob' })
	const href = URL.createObjectURL(response.data)
    const link = document.createElement('a')
    link.href = href
    link.download = doc.value.attachment?.file_name || doc.value.title
    link.click()
    URL.revokeObjectURL(href)
  } catch (requestError) {
    error.value = requestError.response?.data?.message || '文件下载失败，请稍后重试。'
  } finally {
    downloading.value = false
  }
}

async function loadDoc() {
	revokePreviewUrl()
  loading.value = true
  error.value = ''
  try {
    const response = await api.get(`/docs/${route.params.id}`)
    doc.value = response.data
	await loadPreview()
  } catch (requestError) {
    doc.value = null
    error.value = requestError.response?.data?.message || '请检查网络连接后重试。'
  } finally {
    loading.value = false
  }
}

async function loadPreview() {
	const attachment = doc.value?.attachment
	if (!attachment || attachment.preview_status !== 'ready' || (attachment.mime_type !== 'application/pdf' && !isRasterImage.value && !isDocx.value)) return
	previewLoading.value = true
	try {
		const response = await api.get(`/docs/${doc.value.id}/preview`, { responseType: 'blob', silentError: true })
		previewBlob.value = response.data
		if (!isDocx.value) previewUrl.value = URL.createObjectURL(response.data)
	} catch {
		previewUrl.value = ''
		previewBlob.value = null
	} finally {
		previewLoading.value = false
	}
}

function revokePreviewUrl() {
	if (previewUrl.value) URL.revokeObjectURL(previewUrl.value)
	previewUrl.value = ''
	previewBlob.value = null
}

function handlePreviewError() {
	revokePreviewUrl()
}

watch(() => route.params.id, (id, previousId) => {
  if (id && id !== previousId) loadDoc()
})
onMounted(loadDoc)
onBeforeUnmount(revokePreviewUrl)
</script>

<style scoped>
.doc-detail-page { min-height: calc(100vh - 100px); padding: 4px 0 48px; color: var(--theme-text-primary); }
.doc-card { width: min(980px, 100%); min-height: 520px; margin: 0 auto; padding: clamp(26px, 5vw, 58px) clamp(20px, 7vw, 76px); border: 1px solid var(--theme-border); border-top: 3px solid var(--theme-primary); background: var(--theme-bg-card); }
.loading-card { padding-top: 54px; }
.header-topline { display: flex; align-items: center; justify-content: space-between; gap: 16px; }
.back-button { margin-left: -12px; }
.title-row { display: flex; align-items: flex-start; justify-content: space-between; gap: 24px; margin-top: clamp(34px, 7vw, 68px); }
.doc-kicker { margin: 0 0 12px; color: var(--theme-primary); font-size: 11px; font-weight: 700; letter-spacing: .16em; }
.title-row h1 { margin: 0; font-size: clamp(32px, 5vw, 52px); line-height: 1.15; overflow-wrap: anywhere; }
.doc-meta { display: flex; flex-wrap: wrap; gap: 18px 44px; margin: 32px 0 0; }
.doc-meta div { min-width: 120px; }
.doc-meta dt { margin-bottom: 5px; color: var(--theme-text-tertiary); font-size: 11px; font-weight: 600; letter-spacing: .08em; }
.doc-meta dd { margin: 0; color: var(--theme-text-secondary); font-size: 13px; }
.content-rule { height: 1px; margin: 36px 0 clamp(28px, 5vw, 48px); background: var(--theme-border); }
.preview-loading { min-height: 360px; padding: 28px; border: 1px solid var(--theme-border-light); background: var(--theme-bg-secondary); }
.file-fallback { display: grid; grid-template-columns: auto 1fr auto; align-items: center; gap: 24px; padding: clamp(24px, 5vw, 42px); border: 1px solid var(--theme-border-light); background: var(--theme-bg-secondary); }
.file-mark { display: grid; width: 82px; height: 104px; place-items: center; border: 1px solid var(--theme-primary); background: var(--theme-bg-card); color: var(--theme-primary); font: 700 13px/1 'SFMono-Regular', Consolas, monospace; letter-spacing: .08em; }
.file-copy { min-width: 0; }
.file-copy h2 { margin: 5px 0 8px; overflow-wrap: anywhere; font-size: 22px; }
.file-copy p { margin: 0 0 12px; color: var(--theme-text-secondary); font-size: 13px; }
.file-copy .file-label { margin: 0; color: var(--theme-primary); font-size: 10px; font-weight: 700; letter-spacing: .14em; }
.header-topline :deep(.el-button:focus-visible) { outline: 2px solid var(--theme-primary); outline-offset: 3px; }
@media (max-width: 600px) { .doc-detail-page { padding-bottom: 20px; } .doc-card { min-height: calc(100vh - 100px); padding: 20px 16px 36px; border-right: 0; border-left: 0; } .header-topline { flex-wrap: wrap; } .title-row { align-items: flex-start; flex-direction: column-reverse; gap: 14px; } .doc-meta { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; } .plain-content { margin-right: -16px; margin-left: -16px; border-right: 0; border-left: 0; border-radius: 0; } .file-fallback { grid-template-columns: auto 1fr; margin-right: -16px; margin-left: -16px; padding: 24px 16px; border-right: 0; border-left: 0; } .file-fallback .el-button { grid-column: 1 / -1; width: 100%; } }
</style>
