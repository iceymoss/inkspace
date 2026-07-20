<template>
  <section class="share-page">
    <div v-if="loading" class="share-loading">
      <el-skeleton :rows="8" animated />
    </div>

    <article v-else-if="doc" class="shared-document">
      <header>
        <div class="share-badge"><el-icon><Link /></el-icon>共享文档</div>
        <h1>{{ doc.title }}</h1>
        <div class="doc-meta">
          <span v-if="doc.owner_name || doc.author_name">由 {{ doc.owner_name || doc.author_name }} 分享</span>
          <span>更新于 {{ formatDate(doc.updated_at) }}</span>
          <span v-if="doc.word_count">{{ doc.word_count }} 字</span>
          <span v-if="doc.view_count !== undefined">阅读 {{ doc.view_count }}</span>
        </div>
      </header>
      <div class="document-rule"></div>
      <div class="markdown-body" v-html="doc.content_html"></div>
    </article>

    <div v-else class="invalid-state">
      <div class="invalid-icon"><el-icon><Warning /></el-icon></div>
      <h1>{{ errorState.title }}</h1>
      <p>{{ errorState.description }}</p>
      <el-button type="primary" plain @click="$router.push('/')">返回 InkSpace 首页</el-button>
    </div>
  </section>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import dayjs from 'dayjs'
import { Link, Warning } from '@element-plus/icons-vue'
import api from '@/utils/api'

const route = useRoute()
const loading = ref(true)
const doc = ref(null)
const errorState = ref({ title: '链接已失效', description: '该分享链接当前无法访问。' })
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

onMounted(async () => {
  try {
    const response = await api.get(`/share/${route.params.token}`)
    doc.value = response.data
  } catch (error) {
    errorState.value = resolveError(error)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.share-page { min-height: 65vh; padding: 38px 20px 70px; background: var(--theme-bg-secondary); }
.share-loading, .shared-document, .invalid-state { max-width: 860px; margin: 0 auto; }
.share-loading, .shared-document { padding: clamp(24px, 6vw, 64px); background: var(--theme-content-bg); border: 1px solid var(--theme-border); border-radius: 14px; box-shadow: 0 12px 40px var(--theme-shadow); }
.share-badge { display: inline-flex; align-items: center; gap: 6px; padding: 5px 10px; border-radius: 99px; background: color-mix(in srgb, var(--theme-primary) 11%, transparent); color: var(--theme-primary); font-size: 12px; font-weight: 600; }
.shared-document h1 { margin: 20px 0 10px; font-size: clamp(30px, 5vw, 45px); line-height: 1.25; }
.doc-meta { display: flex; flex-wrap: wrap; gap: 8px 18px; color: var(--theme-text-tertiary); font-size: 13px; }
.document-rule { width: 55px; height: 3px; margin: 28px 0 36px; border-radius: 99px; background: var(--theme-primary); }
.markdown-body { color: var(--theme-text-primary); font-size: 16px; line-height: 1.85; overflow-wrap: anywhere; }
.markdown-body :deep(h1), .markdown-body :deep(h2), .markdown-body :deep(h3) { margin: 1.6em 0 .7em; }
.markdown-body :deep(p) { color: var(--theme-text-secondary); }
.markdown-body :deep(img) { max-width: 100%; border-radius: 8px; }
.markdown-body :deep(pre) { padding: 16px; overflow: auto; border: 1px solid var(--theme-border); border-radius: 8px; background: var(--theme-bg-secondary); }
.markdown-body :deep(code) { font-family: 'SFMono-Regular', Consolas, monospace; }
.markdown-body :deep(blockquote) { margin-left: 0; padding-left: 18px; border-left: 4px solid var(--theme-primary); color: var(--theme-text-tertiary); }
.markdown-body :deep(table) { display: block; max-width: 100%; border-collapse: collapse; overflow-x: auto; }
.markdown-body :deep(th), .markdown-body :deep(td) { padding: 8px 12px; border: 1px solid var(--theme-border); }
.invalid-state { padding: 90px 20px; text-align: center; }
.invalid-icon { display: grid; place-items: center; width: 72px; height: 72px; margin: 0 auto 20px; border-radius: 22px; background: color-mix(in srgb, var(--theme-text-tertiary) 12%, transparent); color: var(--theme-text-tertiary); font-size: 34px; }
.invalid-state h1 { margin: 0 0 10px; font-size: 28px; }
.invalid-state p { margin: 0 0 24px; color: var(--theme-text-tertiary); }
@media (max-width: 600px) { .share-page { padding: 14px 10px 40px; } .shared-document { border-radius: 10px; } }
</style>
