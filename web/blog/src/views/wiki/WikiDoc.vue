<template>
  <main class="doc-page">
    <div
      v-if="loading"
      class="doc-paper loading-paper"
      aria-busy="true"
      aria-label="正在加载文档"
    >
      <el-skeleton
        :rows="12"
        animated
      />
    </div>

    <section
      v-else-if="error"
      class="doc-state"
      role="alert"
    >
      <span>!</span>
      <h1>文章暂时无法打开</h1>
      <p>{{ error }}</p>
      <div class="state-actions">
        <button
          type="button"
          @click="loadDoc"
        >
          重新加载
        </button>
        <router-link to="/wiki">
          返回知识刊物
        </router-link>
      </div>
    </section>

    <article
      v-else-if="doc"
      class="doc-paper"
    >
      <nav
        class="doc-nav"
        aria-label="返回导航"
      >
        <router-link :to="`/wiki/${doc.workspace_id}`">
          <span aria-hidden="true">←</span> 返回知识库
        </router-link>
        <router-link to="/wiki">
          知识刊物
        </router-link>
      </nav>

      <header class="doc-header">
        <p class="kicker">
          PUBLIC DOCUMENT · NO. {{ doc.id }}
        </p>
        <h1>{{ doc.title || '无标题文档' }}</h1>
        <dl class="metadata">
          <div>
            <dt>发布</dt>
            <dd><time :datetime="doc.published_at">{{ formatDate(doc.published_at) }}</time></dd>
          </div>
          <div>
            <dt>修订</dt>
            <dd><time :datetime="doc.updated_at">{{ formatDate(doc.updated_at) }}</time></dd>
          </div>
          <div>
            <dt>阅读</dt>
            <dd>{{ doc.view_count ?? 0 }} 次</dd>
          </div>
        </dl>
      </header>

      <div class="editorial-rule">
        <span>INKSPACE</span>
      </div>
      <MarkdownPreview
        v-if="doc.content_html"
        class="document-body"
        :html="doc.content_html"
        editorial
      />
      <div
        v-else
        class="content-empty"
      >
        <span>∅</span>
        <p>这篇文章暂时没有正文。</p>
      </div>

      <footer class="doc-footer">
        <p>END OF DOCUMENT</p>
        <router-link :to="`/wiki/${doc.workspace_id}`">
          返回目录继续阅读 <span aria-hidden="true">→</span>
        </router-link>
      </footer>
    </article>
  </main>
</template>

<script setup>
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import dayjs from 'dayjs'
import api from '@/utils/api'
import MarkdownPreview from '@/components/docs/MarkdownPreview.vue'

const route = useRoute()
const doc = ref(null)
const loading = ref(true)
const error = ref('')
const formatDate = value => value ? dayjs(value).format('YYYY年MM月DD日') : '未记录'

async function loadDoc() {
  loading.value = true
  error.value = ''
  try {
    const response = await api.get(`/wiki/docs/${route.params.id}`)
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
.doc-page { min-height: calc(100vh - 60px); padding: clamp(20px, 5vw, 64px) 20px 90px; background: var(--bg-soft, var(--theme-bg-secondary)); color: var(--ink, var(--theme-text-primary)); }
.doc-paper { width: min(940px, 100%); margin: 0 auto; padding: clamp(38px, 8vw, 86px) clamp(24px, 8vw, 94px); border: 1px solid var(--hairline, var(--theme-border)); border-top: 3px double var(--ink, var(--theme-text-primary)); background: var(--surface, var(--theme-bg-card)); box-shadow: none; }
.loading-paper { min-height: 620px; }
.doc-nav { display: flex; justify-content: space-between; padding-bottom: 16px; border-bottom: 1px solid var(--hairline, var(--theme-border)); }
.doc-nav a { color: var(--sub, var(--theme-text-secondary)); font-size: 11px; font-weight: 700; letter-spacing: .05em; text-decoration: none; }
.doc-nav a:hover, .doc-footer a:hover { color: var(--accent-hover, var(--theme-primary-hover)); }
.doc-nav a:focus-visible, .doc-footer a:focus-visible, button:focus-visible { outline: 2px solid var(--accent, var(--theme-primary)); outline-offset: 4px; }
.doc-header { padding: clamp(44px, 8vw, 76px) 0 34px; }
.kicker { margin: 0 0 20px; color: var(--accent, var(--theme-primary)); font-size: 10px; font-weight: 700; letter-spacing: .2em; }
.doc-header h1 { margin: 0; font-family: Georgia, 'Songti SC', serif; font-size: clamp(40px, 7vw, 72px); font-weight: 500; letter-spacing: -.045em; line-height: 1.12; text-wrap: balance; }
.metadata { display: flex; flex-wrap: wrap; gap: 16px 38px; margin: 34px 0 0; }
.metadata div { display: grid; grid-template-columns: auto auto; gap: 8px; }
.metadata dt { color: var(--accent, var(--theme-primary)); font-size: 10px; font-weight: 700; letter-spacing: .1em; }
.metadata dd { margin: 0; color: var(--sub, var(--theme-text-secondary)); font-size: 11px; }
.editorial-rule { display: flex; align-items: center; gap: 14px; color: var(--sub, var(--theme-text-tertiary)); font-size: 8px; font-weight: 700; letter-spacing: .2em; }
.editorial-rule::before, .editorial-rule::after { height: 1px; background: var(--ink, var(--theme-text-primary)); content: ''; }
.editorial-rule::before { width: 54px; }
.editorial-rule::after { flex: 1; }
.document-body { padding-top: clamp(36px, 6vw, 58px); }
.content-empty { padding: 75px 0; color: var(--sub, var(--theme-text-tertiary)); text-align: center; }
.content-empty span { font: 30px Georgia, serif; }
.doc-footer { display: flex; align-items: center; justify-content: space-between; gap: 20px; margin-top: 70px; padding-top: 18px; border-top: 3px double var(--ink, var(--theme-text-primary)); }
.doc-footer p { margin: 0; color: var(--sub, var(--theme-text-tertiary)); font-size: 9px; font-weight: 700; letter-spacing: .18em; }
.doc-footer a { color: var(--accent, var(--theme-primary)); font-size: 12px; font-weight: 700; text-decoration: none; }
.doc-state { max-width: 720px; margin: 0 auto; padding: 100px 20px; text-align: center; }
.doc-state > span { display: grid; place-items: center; width: 56px; height: 56px; margin: 0 auto 20px; border: 1px solid var(--accent, var(--theme-primary)); border-radius: 50%; color: var(--accent, var(--theme-primary)); font: 26px Georgia, serif; }
.doc-state h1 { margin: 0 0 10px; font-family: Georgia, 'Songti SC', serif; font-weight: 500; }
.doc-state p { color: var(--sub, var(--theme-text-secondary)); }
.state-actions { display: flex; justify-content: center; gap: 18px; margin-top: 24px; }
.state-actions button, .state-actions a { padding: 8px 0; border: 0; border-bottom: 1px solid currentColor; background: none; color: var(--accent, var(--theme-primary)); font: inherit; font-weight: 700; text-decoration: none; cursor: pointer; }
@media (max-width: 560px) { .doc-page { padding: 10px 8px 40px; } .doc-paper { padding: 32px 16px 40px; box-shadow: none; } .doc-nav a:last-child { display: none; } .doc-header h1 { font-size: 39px; } .metadata { gap: 10px 20px; } .metadata div { display: block; } .metadata dd { margin-top: 3px; } .doc-footer { align-items: flex-start; flex-direction: column; } }
</style>
