<template>
  <main class="workspace-page">
    <div
      v-if="loading"
      class="workspace-shell loading-shell"
      aria-busy="true"
    >
      <el-skeleton
        :rows="11"
        animated
      />
    </div>

    <section
      v-else-if="error"
      class="workspace-state"
      role="alert"
    >
      <span>!</span>
      <h1>知识库暂时无法打开</h1>
      <p>{{ error }}</p>
      <div class="state-actions">
        <button
          type="button"
          @click="loadTree"
        >
          重新加载
        </button>
        <router-link to="/wiki">
          返回知识刊物
        </router-link>
      </div>
    </section>

    <div
      v-else-if="tree"
      class="workspace-shell"
    >
      <header class="workspace-hero">
        <router-link
          class="back-link"
          to="/wiki"
        >
          <span aria-hidden="true">←</span> 全部知识库
        </router-link>
        <div class="hero-copy">
          <span
            class="workspace-icon"
            aria-hidden="true"
          >
            <img
              v-if="isImageIcon(tree.workspace.icon)"
              :src="tree.workspace.icon"
              alt=""
            >
            <template v-else>{{ tree.workspace.icon || tree.workspace.name?.slice(0, 1) || '知' }}</template>
          </span>
          <div>
            <p class="kicker">
              PUBLIC KNOWLEDGE · {{ docs.length }} ARTICLES
            </p>
            <h1>{{ tree.workspace.name }}</h1>
            <p class="description">
              {{ tree.workspace.description || '一个持续整理、公开阅读的知识空间。' }}
            </p>
          </div>
        </div>
      </header>

      <button
        v-if="mobileNavOpen"
        class="nav-backdrop"
        type="button"
        aria-label="关闭目录"
        @click="mobileNavOpen = false"
      />

      <div class="reading-layout">
        <aside
          :class="['tree-nav', { open: mobileNavOpen }]"
          aria-label="知识库目录"
        >
          <div class="nav-heading">
            <div><span>CONTENTS</span><strong>阅读目录</strong></div>
            <button
              type="button"
              aria-label="关闭目录"
              @click="mobileNavOpen = false"
            >
              ×
            </button>
          </div>
          <nav
            v-if="navigationRows.length"
            class="tree-list"
          >
            <template
              v-for="row in navigationRows"
              :key="`${row.type}-${row.id}`"
            >
              <div
                v-if="row.type === 'catalog'"
                class="catalog-row"
                :style="{ '--depth': row.depth }"
              >
                <span class="catalog-rule" />{{ row.name }}
              </div>
              <router-link
                v-else
                class="doc-row"
                :style="{ '--depth': row.depth }"
                :to="`/wiki/docs/${row.id}`"
                @click="mobileNavOpen = false"
              >
                <span>{{ row.title }}</span><span aria-hidden="true">↗</span>
              </router-link>
            </template>
          </nav>
          <p
            v-else
            class="nav-empty"
          >
            暂无已发布文档
          </p>
        </aside>

        <section
          class="workspace-content"
          aria-labelledby="documents-heading"
        >
          <div class="content-heading">
            <button
              class="mobile-nav-button"
              type="button"
              :aria-expanded="mobileNavOpen"
              aria-controls="workspace-document-list"
              @click="mobileNavOpen = true"
            >
              <span aria-hidden="true">☰</span> 阅读目录
            </button>
            <div>
              <p>EDITOR'S INDEX</p>
              <h2 id="documents-heading">
                全部文章
              </h2>
            </div>
            <span>{{ docs.length }} 篇公开文档</span>
          </div>

          <div
            v-if="docs.length"
            id="workspace-document-list"
            class="document-index"
          >
            <router-link
              v-for="(doc, index) in docs"
              :key="doc.id"
              class="document-entry"
              :to="`/wiki/docs/${doc.id}`"
            >
              <span class="entry-number">{{ String(index + 1).padStart(2, '0') }}</span>
              <div>
                <p>{{ catalogName(doc.catalog_id) }}</p>
                <h3>{{ doc.title || '无标题文档' }}</h3>
              </div>
              <time :datetime="doc.updated_at">{{ formatDate(doc.updated_at) }}</time>
              <span
                class="entry-arrow"
                aria-hidden="true"
              >→</span>
            </router-link>
          </div>

          <div
            v-else
            class="empty-state"
          >
            <span>∅</span>
            <h2>这里还没有公开文章</h2>
            <p>内容发布后，会按目录顺序陈列在这里。</p>
          </div>
        </section>
      </div>
    </div>
  </main>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import dayjs from 'dayjs'
import api from '@/utils/api'

const route = useRoute()
const tree = ref(null)
const loading = ref(true)
const error = ref('')
const mobileNavOpen = ref(false)

const docs = computed(() => [...(tree.value?.docs || [])].sort(compareNodes))
const catalogNames = computed(() => {
  const names = new Map()
  const visit = catalogs => catalogs.forEach(catalog => {
    names.set(catalog.id, catalog.name)
    visit(catalog.children || [])
  })
  visit(tree.value?.catalogs || [])
  return names
})

const navigationRows = computed(() => {
  if (!tree.value) return []
  const docsByCatalog = new Map()
  for (const doc of tree.value.docs || []) {
    const catalogId = doc.catalog_id ?? null
    if (!docsByCatalog.has(catalogId)) docsByCatalog.set(catalogId, [])
    docsByCatalog.get(catalogId).push(doc)
  }
  docsByCatalog.forEach(group => group.sort(compareNodes))
  const rows = (docsByCatalog.get(null) || []).map(doc => ({ ...doc, type: 'doc', depth: 0 }))
  const appendCatalogs = (catalogs, depth) => {
    for (const catalog of [...catalogs].sort(compareNodes)) {
      rows.push({ ...catalog, type: 'catalog', depth })
      rows.push(...(docsByCatalog.get(catalog.id) || []).map(doc => ({ ...doc, type: 'doc', depth: depth + 1 })))
      appendCatalogs(catalog.children || [], depth + 1)
    }
  }
  appendCatalogs(tree.value.catalogs || [], 0)
  return rows
})

const formatDate = value => value ? dayjs(value).format('YYYY.MM.DD') : '-'
const isImageIcon = value => /^(https?:\/\/|\/uploads\/)/.test(value || '')
const catalogName = id => id == null ? '根目录' : catalogNames.value.get(id) || '未分类'

function compareNodes(a, b) {
  return (a.sort || 0) - (b.sort || 0) || a.id - b.id
}

async function loadTree() {
  loading.value = true
  error.value = ''
  mobileNavOpen.value = false
  try {
    const response = await api.get(`/wiki/workspaces/${route.params.workspaceId}/tree`)
    tree.value = response.data
  } catch (requestError) {
    tree.value = null
    error.value = requestError.response?.data?.message || '请检查网络连接后重试。'
  } finally {
    loading.value = false
  }
}

watch(() => route.params.workspaceId, (id, previousId) => {
  if (id && id !== previousId) loadTree()
})
onMounted(loadTree)
</script>

<style scoped>
.workspace-page { min-height: calc(100vh - 60px); background: var(--bg, var(--theme-bg-primary)); color: var(--ink, var(--theme-text-primary)); }
.workspace-shell { width: min(1120px, calc(100% - 40px)); margin: 0 auto; padding: 46px 0 80px; }
.loading-shell { max-width: 1000px; padding-top: 90px; }
.workspace-hero { padding-bottom: 36px; border-bottom: 3px double var(--ink, var(--theme-text-primary)); }
.back-link { display: inline-flex; gap: 8px; margin-bottom: 44px; color: var(--sub, var(--theme-text-secondary)); font-size: 12px; font-weight: 700; letter-spacing: .04em; text-decoration: none; }
.back-link:hover { color: var(--accent, var(--theme-primary)); }
.back-link:focus-visible, .tree-nav a:focus-visible, .document-entry:focus-visible, button:focus-visible { outline: 2px solid var(--accent, var(--theme-primary)); outline-offset: 3px; }
.hero-copy { display: grid; grid-template-columns: 112px minmax(0, 1fr); align-items: start; gap: 30px; }
.workspace-icon { display: grid; place-items: center; width: 112px; height: 112px; overflow: hidden; border: 1px solid var(--hairline, var(--theme-border)); border-radius: 50%; background: var(--accent-soft, var(--theme-bg-hover)); font-family: Georgia, 'Songti SC', serif; font-size: 42px; }
.workspace-icon img { width: 100%; height: 100%; object-fit: cover; }
.kicker, .content-heading p { margin: 4px 0 10px; color: var(--accent, var(--theme-primary)); font-size: 10px; font-weight: 700; letter-spacing: .18em; }
.hero-copy h1 { max-width: 900px; margin: 0; font-family: Georgia, 'Songti SC', serif; font-size: clamp(46px, 7vw, 84px); font-weight: 500; letter-spacing: -.05em; line-height: 1; }
.description { max-width: 720px; margin: 18px 0 0; color: var(--sub, var(--theme-text-secondary)); font-family: Georgia, 'Songti SC', serif; font-size: 17px; line-height: 1.75; }
.reading-layout { display: grid; grid-template-columns: 300px minmax(0, 1fr); align-items: start; }
.tree-nav { position: sticky; top: 76px; max-height: calc(100vh - 95px); padding: 30px 26px 30px 0; overflow-y: auto; border-right: 1px solid var(--hairline, var(--theme-border)); }
.nav-heading { display: flex; align-items: center; justify-content: space-between; margin-bottom: 22px; }
.nav-heading div { display: flex; flex-direction: column; gap: 3px; }
.nav-heading span { color: var(--accent, var(--theme-primary)); font-size: 9px; font-weight: 700; letter-spacing: .18em; }
.nav-heading strong { font-family: Georgia, 'Songti SC', serif; font-size: 19px; font-weight: 500; }
.nav-heading button { display: none; border: 0; background: transparent; color: inherit; font-size: 28px; cursor: pointer; }
.tree-list { display: flex; flex-direction: column; }
.catalog-row { display: flex; align-items: center; gap: 8px; margin-top: 18px; padding: 7px 0 7px calc(var(--depth) * 14px); color: var(--ink, var(--theme-text-primary)); font-family: Georgia, 'Songti SC', serif; font-size: 14px; font-weight: 700; }
.catalog-rule { width: 14px; height: 1px; background: var(--accent, var(--theme-primary)); }
.doc-row { display: flex; align-items: flex-start; justify-content: space-between; gap: 10px; padding: 8px 2px 8px calc(var(--depth) * 14px); border-bottom: 1px dotted var(--hairline, var(--theme-border)); color: var(--sub, var(--theme-text-secondary)); font-size: 13px; line-height: 1.45; text-decoration: none; }
.doc-row:hover { color: var(--accent, var(--theme-primary)); }
.nav-empty { color: var(--sub, var(--theme-text-tertiary)); font-size: 13px; }
.workspace-content { min-width: 0; padding: 32px 0 0 42px; }
.content-heading { display: flex; align-items: end; justify-content: space-between; gap: 20px; padding-bottom: 18px; border-bottom: 1px solid var(--ink, var(--theme-text-primary)); }
.content-heading p { margin-top: 0; }
.content-heading h2 { margin: 0; font-family: Georgia, 'Songti SC', serif; font-size: 35px; font-weight: 500; }
.content-heading > span { color: var(--sub, var(--theme-text-tertiary)); font-size: 12px; }
.mobile-nav-button { display: none; }
.document-entry { display: grid; grid-template-columns: 48px minmax(0, 1fr) auto 22px; align-items: center; gap: 18px; min-height: 106px; border-bottom: 1px solid var(--hairline, var(--theme-border)); color: inherit; text-decoration: none; transition: padding .2s, background .2s; }
.document-entry:hover { padding-right: 12px; padding-left: 12px; background: var(--accent-soft, var(--theme-bg-hover)); }
.entry-number { color: var(--accent, var(--theme-primary)); font-family: Georgia, serif; font-size: 13px; font-style: italic; }
.document-entry p { margin: 0 0 6px; color: var(--sub, var(--theme-text-tertiary)); font-size: 10px; font-weight: 700; letter-spacing: .08em; }
.document-entry h3 { margin: 0; font-family: Georgia, 'Songti SC', serif; font-size: clamp(20px, 2.2vw, 27px); font-weight: 500; line-height: 1.3; }
.document-entry time { color: var(--sub, var(--theme-text-tertiary)); font-size: 11px; }
.entry-arrow { color: var(--accent, var(--theme-primary)); font-size: 18px; }
.empty-state, .workspace-state { padding: 100px 20px; text-align: center; }
.empty-state > span, .workspace-state > span { display: grid; place-items: center; width: 56px; height: 56px; margin: 0 auto 20px; border: 1px solid var(--accent, var(--theme-primary)); border-radius: 50%; color: var(--accent, var(--theme-primary)); font: 26px Georgia, serif; }
.empty-state h2, .workspace-state h1 { margin: 0 0 10px; font-family: Georgia, 'Songti SC', serif; font-weight: 500; }
.empty-state p, .workspace-state p { color: var(--sub, var(--theme-text-secondary)); }
.state-actions { display: flex; justify-content: center; gap: 18px; margin-top: 24px; }
.state-actions button, .state-actions a { padding: 8px 0; border: 0; border-bottom: 1px solid currentColor; background: none; color: var(--accent, var(--theme-primary)); font: inherit; font-weight: 700; text-decoration: none; cursor: pointer; }
.nav-backdrop { display: none; }
@media (max-width: 900px) { .workspace-shell { width: min(100% - 24px, 1280px); padding-top: 28px; } .back-link { margin-bottom: 28px; } .hero-copy { grid-template-columns: 70px minmax(0, 1fr); gap: 18px; } .workspace-icon { width: 70px; height: 70px; font-size: 28px; } .hero-copy h1 { font-size: 44px; } .description { font-size: 15px; } .reading-layout { display: block; } .tree-nav { position: fixed; inset: 0 auto 0 0; z-index: 31; width: min(340px, 88vw); max-height: none; padding: 25px 22px; border: 0; background: var(--surface, var(--theme-bg-card)); box-shadow: 12px 0 40px var(--shadow-color, var(--theme-shadow)); visibility: hidden; transform: translateX(-105%); transition: transform .2s ease, visibility .2s; } .tree-nav.open { visibility: visible; transform: translateX(0); } .nav-heading button { display: block; } .nav-backdrop { display: block; position: fixed; inset: 0; z-index: 30; border: 0; background: rgb(0 0 0 / 40%); } .workspace-content { padding: 26px 0 0; } .content-heading { align-items: center; } .content-heading > div { flex: 1; } .content-heading > span { display: none; } .mobile-nav-button { display: inline-flex; align-items: center; gap: 7px; padding: 8px 10px; border: 1px solid var(--hairline, var(--theme-border)); background: var(--surface, var(--theme-bg-card)); color: inherit; font: inherit; font-size: 12px; cursor: pointer; } .document-entry { grid-template-columns: 35px minmax(0, 1fr) 18px; gap: 10px; min-height: 94px; } .document-entry time { display: none; } }
@media (prefers-reduced-motion: reduce) { .tree-nav, .document-entry { transition: none; } }
</style>
