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
            {{ doc.status === 1 ? '已发布' : '草稿' }}
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
      <div
        v-if="doc.kind === 'markdown' && doc.content_html"
        class="document-body"
        v-html="doc.content_html"
      />
      <pre
        v-else-if="isPlainText && doc.content"
        class="plain-content"
        :class="{ code: doc.kind === 'code' }"
      ><code>{{ doc.content }}</code></pre>
      <el-empty
        v-else
        description="这篇文档暂时没有正文"
        :image-size="72"
      />
    </article>
  </section>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import dayjs from 'dayjs'
import { ArrowLeft, EditPen } from '@element-plus/icons-vue'
import api from '@/utils/api'

const route = useRoute()
const router = useRouter()
const doc = ref(null)
const loading = ref(true)
const error = ref('')

const isPlainText = computed(() => doc.value?.kind === 'text' || doc.value?.kind === 'code')
const kindLabel = computed(() => {
  const labels = { markdown: 'Markdown', text: '纯文本', code: '代码' }
  const label = labels[doc.value?.kind] || doc.value?.kind || '文档'
  return doc.value?.language ? `${label} · ${doc.value.language}` : label
})
const formatDate = value => value ? dayjs(value).format('YYYY-MM-DD HH:mm') : '未记录'

async function loadDoc() {
  loading.value = true
  error.value = ''
  try {
    const response = await api.get(`/docs/${route.params.id}`)
    doc.value = response.data
  } catch (requestError) {
    doc.value = null
    error.value = requestError.response?.data?.message || '请检查网络连接后重试。'
  } finally {
    loading.value = false
  }
}

watch(() => route.params.id, (id, previousId) => {
  if (id && id !== previousId) loadDoc()
})
onMounted(loadDoc)
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
.document-body { font-size: 16px; line-height: 1.9; overflow-wrap: anywhere; }
.document-body :deep(h1), .document-body :deep(h2), .document-body :deep(h3), .document-body :deep(h4) { margin: 1.7em 0 .7em; line-height: 1.35; }
.document-body :deep(h2) { padding-bottom: .35em; border-bottom: 1px solid var(--theme-border-light); }
.document-body :deep(p) { margin: 0 0 1.35em; }
.document-body :deep(a) { color: var(--theme-primary); text-underline-offset: 3px; }
.document-body :deep(img) { display: block; max-width: 100%; height: auto; margin: 1.8em auto; }
.document-body :deep(blockquote) { margin: 1.8em 0; padding: 2px 0 2px 18px; border-left: 3px solid var(--theme-primary); color: var(--theme-text-secondary); }
.document-body :deep(pre), .plain-content { max-width: 100%; padding: 18px 20px; overflow: auto; border: 1px solid var(--theme-border-light); border-radius: 6px; background: var(--theme-bg-secondary); color: var(--theme-text-primary); font: 13px/1.7 'SFMono-Regular', Consolas, monospace; }
.document-body :deep(code) { font-family: 'SFMono-Regular', Consolas, monospace; }
.document-body :deep(table) { display: block; width: max-content; max-width: 100%; overflow-x: auto; border-collapse: collapse; }
.document-body :deep(th), .document-body :deep(td) { padding: 9px 12px; border: 1px solid var(--theme-border); }
.plain-content { min-height: 240px; margin: 0; white-space: pre-wrap; overflow-wrap: anywhere; }
.plain-content.code { white-space: pre; }
.header-topline :deep(.el-button:focus-visible), .document-body :deep(a:focus-visible) { outline: 2px solid var(--theme-primary); outline-offset: 3px; }
@media (max-width: 600px) { .doc-detail-page { padding-bottom: 20px; } .doc-card { min-height: calc(100vh - 100px); padding: 20px 16px 36px; border-right: 0; border-left: 0; } .title-row { align-items: flex-start; flex-direction: column-reverse; gap: 14px; } .doc-meta { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; } .document-body { font-size: 15px; line-height: 1.8; } .plain-content { margin-right: -16px; margin-left: -16px; border-right: 0; border-left: 0; border-radius: 0; } }
</style>
