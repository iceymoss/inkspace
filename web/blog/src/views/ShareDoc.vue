<template>
  <section class="share-page">
    <div
      v-if="loading"
      class="share-loading"
    >
      <el-skeleton
        :rows="8"
        animated
      />
    </div>

    <article
      v-else-if="doc"
      class="shared-document"
    >
      <header>
        <div class="share-badge">
          <el-icon><Link /></el-icon>SHARED DOCUMENT · 共享文档
        </div>
        <h1>{{ doc.title }}</h1>
        <div class="doc-meta">
          <span><el-icon><Calendar /></el-icon>更新于 {{ formatDate(doc.updated_at) }}</span>
          <span v-if="doc.word_count"><el-icon><EditPen /></el-icon>{{ doc.word_count }} 字</span>
          <span v-if="doc.view_count !== undefined"><el-icon><View /></el-icon>阅读 {{ doc.view_count }}</span>
        </div>
      </header>
      <div class="document-rule" />
      <div
        class="markdown-body"
        :data-markdown-theme="markdownTheme"
        v-html="doc.content_html"
      />
    </article>

    <div
      v-else
      class="invalid-state"
    >
      <div class="invalid-icon">
        <el-icon><Warning /></el-icon>
      </div>
      <h1>{{ errorState.title }}</h1>
      <p>{{ errorState.description }}</p>
      <div class="invalid-actions">
        <el-button
          v-if="canRetry"
          type="primary"
          @click="fetchSharedDoc"
        >
          重新加载
        </el-button>
        <el-button
          :type="canRetry ? 'default' : 'primary'"
          plain
          @click="$router.push('/')"
        >
          返回 InkSpace 首页
        </el-button>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import dayjs from 'dayjs'
import { Calendar, EditPen, Link, View, Warning } from '@element-plus/icons-vue'
import api from '@/utils/api'
import { loadCodeTheme, loadHighlightTheme } from '@/utils/codeTheme'
import { useAppearanceStore } from '@/stores/appearance'

const route = useRoute()
const appearanceStore = useAppearanceStore()
const loading = ref(true)
const doc = ref(null)
const canRetry = ref(false)
const errorState = ref({ title: '链接已失效', description: '该分享链接当前无法访问。' })
const markdownTheme = computed(() => appearanceStore.resolvedColorScheme)
const formatDate = value => value ? dayjs(value).format('YYYY年MM月DD日 HH:mm') : '-'

function resolveError(error) {
  const status = error.response?.status
  const message = error.response?.data?.message || error.message || ''
  if (message.includes('过期')) {
    return { title: '链接已过期', description: '这个分享链接已超过作者设置的有效期。' }
  }
  if (message.includes('禁用') || message.includes('关闭') || status === 403) {
    return { title: '链接已被关闭', description: '作者已停止通过这个链接分享文档。' }
  }
  if (status === 404 || message.includes('不存在')) {
    return { title: '链接不存在', description: '请检查地址是否完整，或联系分享者获取新链接。' }
  }
  return { title: '暂时无法访问', description: '加载共享文档时出现问题，请稍后重试。' }
}

async function fetchSharedDoc() {
  loading.value = true
  doc.value = null
  canRetry.value = false
  try {
    const response = await api.get(`/share/${route.params.token}`)
    doc.value = response.data
  } catch (error) {
    errorState.value = resolveError(error)
    canRetry.value = errorState.value.title === '暂时无法访问'
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  const codeTheme = await loadCodeTheme()
  await loadHighlightTheme(codeTheme)
  await fetchSharedDoc()
})
</script>

<style scoped>
.share-page { min-height: 70vh; padding: 62px 20px 80px; background: var(--theme-bg-primary); }
.share-loading, .shared-document, .invalid-state { max-width: 860px; margin: 0 auto; }
.share-loading, .shared-document { padding: clamp(36px, 7vw, 76px) clamp(20px, 7vw, 72px); background: var(--theme-bg-card); border: 1px solid var(--theme-border); border-radius: 0; box-shadow: none; }
.share-badge { display: inline-flex; align-items: center; gap: 7px; padding: 0; color: var(--theme-primary); font-family: Georgia, 'Songti SC', serif; font-size: 11px; font-weight: 400; letter-spacing: .2em; }
.shared-document h1 { margin: 20px 0 10px; font-family: Georgia, 'Songti SC', 'Noto Serif SC', SimSun, serif; font-size: clamp(34px, 6vw, 52px); font-weight: 500; letter-spacing: .03em; line-height: 1.3; }
.doc-meta { display: flex; flex-wrap: wrap; gap: 8px 18px; color: var(--theme-text-tertiary); font-size: 13px; }
.doc-meta span { display: inline-flex; align-items: center; gap: 5px; }
.document-rule { width: 100%; height: 1px; margin: 28px 0 36px; background: var(--theme-border); }
.markdown-body { font-size: 16px; line-height: 1.85; overflow-wrap: anywhere; }
.markdown-body :deep(h1), .markdown-body :deep(h2), .markdown-body :deep(h3) { margin: 1.6em 0 .7em; font-family: Georgia, 'Songti SC', serif; font-weight: 500; }
.markdown-body :deep(p) { margin: 0 0 1.25em; }
.markdown-body :deep(a) { color: var(--theme-primary); text-decoration-thickness: 1px; text-underline-offset: 3px; }
.markdown-body :deep(a:hover) { color: var(--theme-primary-hover); }
.markdown-body :deep(ul), .markdown-body :deep(ol) { margin: 0 0 1.25em; padding-left: 1.6em; }
.markdown-body :deep(li) { margin: .35em 0; }
.markdown-body :deep(img) { max-width: 100%; border-radius: 0; }
.markdown-body :deep(pre) { padding: 16px; overflow: auto; border-radius: 0; }
.markdown-body :deep(code) { font-family: 'SFMono-Regular', Consolas, monospace; }
.markdown-body :deep(:not(pre) > code) { padding: 2px 6px; border-radius: 0; font-size: .9em; }
.markdown-body :deep(blockquote) { margin-left: 0; padding-left: 18px; }
.markdown-body :deep(table) { display: block; max-width: 100%; border-collapse: collapse; overflow-x: auto; }
.markdown-body :deep(th), .markdown-body :deep(td) { padding: 8px 12px; border: 1px solid var(--theme-border); }
.markdown-body :deep(th) { text-align: left; }
.markdown-body :deep(hr) { margin: 2em 0; border: 0; border-top: 1px solid var(--theme-border); }
.markdown-body :deep(input[type='checkbox']) { margin-right: 7px; accent-color: var(--theme-primary); }
.invalid-state { padding: 90px 20px; text-align: center; }
.invalid-icon { display: grid; place-items: center; width: 72px; height: 72px; margin: 0 auto 20px; border: 1px solid var(--theme-border); border-radius: 0; color: var(--theme-text-tertiary); font-size: 34px; }
.invalid-state h1 { margin: 0 0 10px; font-family: Georgia, 'Songti SC', serif; font-size: 32px; font-weight: 500; }
.invalid-state p { margin: 0 0 24px; color: var(--theme-text-tertiary); }
.invalid-actions { display: flex; justify-content: center; gap: 10px; }
.share-page :deep(.el-button) { border-radius: 1px; box-shadow: none; }
.markdown-body :deep(a):focus-visible { outline: 2px solid var(--theme-primary); outline-offset: 3px; }
@media (max-width: 900px) { .share-page { padding-right: 24px; padding-left: 24px; } }
@media (max-width: 560px) { .share-page { padding: 24px 8px 40px; } .shared-document { padding: 32px 18px 40px; } .shared-document h1 { font-size: 32px; } .document-rule { margin: 22px 0 28px; } .invalid-actions { flex-direction: column; } }
</style>
