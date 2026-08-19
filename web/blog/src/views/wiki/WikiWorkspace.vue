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
            v-if="docs.length"
            class="tree-list"
          >
            <button
              type="button"
              :class="['root-row', { selected: selectedCatalogId === null }]"
              :aria-current="selectedCatalogId === null ? 'page' : undefined"
              @click="selectCatalog(null)"
            >
              <span aria-hidden="true">⌂</span>
              <span>根目录</span>
              <small>{{ rootDocCount }}</small>
            </button>
            <div
              v-for="row in visibleCatalogRows"
              :key="row.id"
              :class="['catalog-row', { selected: selectedCatalogId === row.id }]"
              :style="{ '--depth': row.depth }"
            >
              <button
                v-if="row.children?.length"
                class="catalog-toggle"
                type="button"
                :aria-expanded="!collapsedCatalogIds.has(row.id)"
                :aria-label="`${collapsedCatalogIds.has(row.id) ? '展开' : '收起'}${row.name}`"
                @click="toggleCatalog(row.id)"
              >
                <span aria-hidden="true">⌄</span>
              </button>
              <span
                v-else
                class="catalog-toggle-placeholder"
                aria-hidden="true"
              />
              <button
                class="catalog-select"
                type="button"
                :aria-current="selectedCatalogId === row.id ? 'page' : undefined"
                @click="selectCatalog(row.id)"
              >
                <span>{{ row.name }}</span>
                <small>{{ docCount(row.id) }}</small>
              </button>
            </div>
          </nav>
          <p
            v-else
            class="nav-empty"
          >
            暂无已公开文档
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
                {{ selectedCatalogName }}
              </h2>
            </div>
            <span>{{ selectedDocs.length }} 篇文档</span>
          </div>

          <div
            v-if="selectedDocs.length"
            id="workspace-document-list"
            class="document-index"
          >
            <router-link
              v-for="doc in selectedDocs"
              :key="doc.id"
              class="document-entry"
              :to="`/wiki/docs/${doc.id}`"
            >
              <span class="entry-branch" aria-hidden="true" />
              <div>
                <p>公开文档</p>
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
            <h2>此目录没有公开文档</h2>
            <p>请从左侧选择其他目录继续浏览。</p>
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
const collapsedCatalogIds = ref(new Set())
const selectedCatalogId = ref(null)

const docs = computed(() => [...(tree.value?.docs || [])].sort(compareNodes))
const docsByCatalog = computed(() => {
  const groups = new Map()
  for (const doc of docs.value) {
    const catalogId = doc.catalog_id ?? null
    if (!groups.has(catalogId)) groups.set(catalogId, [])
    groups.get(catalogId).push(doc)
  }
  return groups
})
const catalogRows = computed(() => {
  const rows = []
  const appendCatalogs = (catalogs, depth, ancestorIds = []) => {
    for (const catalog of [...catalogs].sort(compareNodes)) {
      rows.push({ ...catalog, depth, ancestorIds })
      const childAncestorIds = [...ancestorIds, catalog.id]
      appendCatalogs(catalog.children || [], depth + 1, childAncestorIds)
    }
  }
  appendCatalogs(tree.value?.catalogs || [], 0)
  return rows
})
const visibleCatalogRows = computed(() => catalogRows.value.filter(row => (
  !row.ancestorIds.some(id => collapsedCatalogIds.value.has(id))
)))
const selectedDocs = computed(() => docsByCatalog.value.get(selectedCatalogId.value) || [])
const selectedCatalogName = computed(() => {
  if (selectedCatalogId.value === null) return '根目录'
  return catalogRows.value.find(catalog => catalog.id === selectedCatalogId.value)?.name || '目录'
})
const rootDocCount = computed(() => docsByCatalog.value.get(null)?.length || 0)

const formatDate = value => value ? dayjs(value).format('YYYY.MM.DD') : '-'
const isImageIcon = value => /^(https?:\/\/|\/uploads\/)/.test(value || '')
const docCount = id => docsByCatalog.value.get(id)?.length || 0

function compareNodes(a, b) {
  return (a.sort || 0) - (b.sort || 0) || a.id - b.id
}

function toggleCatalog(id) {
  const next = new Set(collapsedCatalogIds.value)
  if (next.has(id)) next.delete(id)
  else next.add(id)
  collapsedCatalogIds.value = next
}

function selectCatalog(id) {
  selectedCatalogId.value = id
  mobileNavOpen.value = false
}

async function loadTree() {
  loading.value = true
  error.value = ''
  mobileNavOpen.value = false
  collapsedCatalogIds.value = new Set()
  selectedCatalogId.value = null
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
.workspace-page { min-height: calc(100vh - 60px); padding: 30px 0 80px; background: var(--bg-soft, var(--theme-bg-secondary)); color: var(--ink, var(--theme-text-primary)); }
.workspace-shell { width: min(1180px, calc(100% - 40px)); margin: 0 auto; }
.loading-shell { min-height: 640px; padding: 48px; border: 1px solid var(--hairline, var(--theme-border)); background: var(--surface, var(--theme-bg-card)); }
.workspace-hero { position: relative; padding: 28px 34px 32px; overflow: hidden; border: 1px solid var(--hairline, var(--theme-border)); border-bottom: 0; background: var(--surface, var(--theme-bg-card)); }
.workspace-hero::after { position: absolute; right: -34px; bottom: -58px; width: 190px; height: 190px; border: 34px solid color-mix(in srgb, var(--accent, var(--theme-primary)) 7%, transparent); border-radius: 50%; content: ''; pointer-events: none; }
.back-link { display: inline-flex; align-items: center; gap: 8px; margin-bottom: 24px; color: var(--sub, var(--theme-text-secondary)); font-size: 12px; font-weight: 700; text-decoration: none; }
.back-link:hover { color: var(--accent, var(--theme-primary)); }
.back-link:focus-visible, .tree-nav a:focus-visible, .document-entry:focus-visible, button:focus-visible { outline: 2px solid var(--accent, var(--theme-primary)); outline-offset: 3px; }
.hero-copy { display: grid; grid-template-columns: 82px minmax(0, 1fr); align-items: center; gap: 24px; }
.workspace-icon { display: grid; place-items: center; width: 82px; height: 82px; overflow: hidden; border: 1px solid color-mix(in srgb, var(--accent, var(--theme-primary)) 35%, var(--hairline, var(--theme-border))); border-radius: 18px; background: var(--accent-soft, var(--theme-bg-hover)); color: var(--accent, var(--theme-primary)); font-family: Georgia, 'Songti SC', serif; font-size: 32px; }
.workspace-icon img { width: 100%; height: 100%; object-fit: cover; }
.kicker, .content-heading p { margin: 0 0 7px; color: var(--accent, var(--theme-primary)); font-size: 9px; font-weight: 700; letter-spacing: .18em; }
.hero-copy h1 { max-width: 800px; margin: 0; font-family: Georgia, 'Songti SC', serif; font-size: clamp(34px, 5vw, 54px); font-weight: 500; letter-spacing: -.035em; line-height: 1.08; }
.description { max-width: 720px; margin: 10px 0 0; color: var(--sub, var(--theme-text-secondary)); font-size: 14px; line-height: 1.7; }
.reading-layout { display: grid; grid-template-columns: 284px minmax(0, 1fr); align-items: stretch; border: 1px solid var(--hairline, var(--theme-border)); background: var(--surface, var(--theme-bg-card)); }
.tree-nav { position: sticky; top: 76px; align-self: start; max-height: calc(100vh - 96px); padding: 28px 20px 34px; overflow-y: auto; border-right: 1px solid var(--hairline, var(--theme-border)); background: color-mix(in srgb, var(--bg-soft, var(--theme-bg-secondary)) 55%, var(--surface, var(--theme-bg-card))); }
.nav-heading { display: flex; align-items: center; justify-content: space-between; margin-bottom: 18px; padding: 0 6px 14px; border-bottom: 1px solid var(--hairline, var(--theme-border)); }
.nav-heading div { display: flex; flex-direction: column; gap: 3px; }
.nav-heading span { color: var(--accent, var(--theme-primary)); font-size: 9px; font-weight: 700; letter-spacing: .18em; }
.nav-heading strong { font-size: 16px; font-weight: 700; }
.nav-heading button { display: none; border: 0; background: transparent; color: inherit; font-size: 28px; cursor: pointer; }
.tree-list { display: flex; flex-direction: column; }
.root-row, .catalog-row { min-height: 36px; border-radius: 5px; }
.root-row { display: grid; grid-template-columns: 18px minmax(0, 1fr) auto; align-items: center; gap: 7px; width: 100%; padding: 7px 8px; border: 0; background: transparent; color: var(--ink, var(--theme-text-primary)); font: inherit; font-size: 12px; font-weight: 700; text-align: left; cursor: pointer; }
.root-row:hover, .root-row.selected, .catalog-row:hover, .catalog-row.selected { background: var(--surface, var(--theme-bg-card)); color: var(--accent, var(--theme-primary)); }
.root-row small, .catalog-select small { color: var(--sub, var(--theme-text-tertiary)); font-size: 10px; font-weight: 500; }
.catalog-row { display: grid; grid-template-columns: 18px minmax(0, 1fr); align-items: center; gap: 3px; margin-top: 3px; padding-left: calc(8px + var(--depth) * 16px); }
.catalog-toggle, .catalog-toggle-placeholder { display: inline-grid; width: 18px; height: 18px; place-items: center; }
.catalog-toggle { padding: 0; border: 0; background: transparent; color: var(--accent, var(--theme-primary)); font: inherit; font-size: 14px; line-height: 1; cursor: pointer; }
.catalog-toggle span { transition: transform .15s ease; }
.catalog-toggle[aria-expanded='false'] span { transform: rotate(-90deg); }
.catalog-select { display: flex; align-items: center; justify-content: space-between; min-width: 0; gap: 8px; padding: 7px 8px 7px 2px; border: 0; background: transparent; color: inherit; font: inherit; font-size: 12px; font-weight: 700; text-align: left; cursor: pointer; }
.catalog-select > span { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.nav-empty { color: var(--sub, var(--theme-text-tertiary)); font-size: 13px; }
.workspace-content { min-width: 0; padding: 30px 38px 52px; }
.content-heading { display: flex; align-items: end; justify-content: space-between; gap: 20px; padding-bottom: 18px; border-bottom: 2px solid var(--ink, var(--theme-text-primary)); }
.content-heading h2 { margin: 0; font-family: Georgia, 'Songti SC', serif; font-size: 29px; font-weight: 500; }
.content-heading > span { color: var(--sub, var(--theme-text-tertiary)); font-size: 12px; }
.mobile-nav-button { display: none; }
.document-index { padding-top: 10px; }
.document-entry { display: grid; grid-template-columns: 18px minmax(0, 1fr) auto 22px; align-items: center; gap: 14px; min-height: 82px; padding: 10px 12px; border-bottom: 1px solid var(--hairline, var(--theme-border)); color: inherit; text-decoration: none; transition: background .18s ease, transform .18s ease; }
.document-entry:hover { background: color-mix(in srgb, var(--accent, var(--theme-primary)) 6%, transparent); transform: translateX(3px); }
.entry-branch { width: 9px; height: 9px; border: 1px solid var(--accent, var(--theme-primary)); border-radius: 50%; }
.document-entry p { margin: 0 0 5px; color: var(--sub, var(--theme-text-tertiary)); font-size: 9px; font-weight: 700; letter-spacing: .08em; }
.document-entry h3 { margin: 0; font-family: Georgia, 'Songti SC', serif; font-size: clamp(18px, 2vw, 23px); font-weight: 500; line-height: 1.3; }
.document-entry time { color: var(--sub, var(--theme-text-tertiary)); font-size: 10px; }
.entry-arrow { color: var(--accent, var(--theme-primary)); font-size: 17px; }
.empty-state, .workspace-state { padding: 100px 20px; text-align: center; }
.empty-state > span, .workspace-state > span { display: grid; place-items: center; width: 56px; height: 56px; margin: 0 auto 20px; border: 1px solid var(--accent, var(--theme-primary)); border-radius: 50%; color: var(--accent, var(--theme-primary)); font: 26px Georgia, serif; }
.empty-state h2, .workspace-state h1 { margin: 0 0 10px; font-family: Georgia, 'Songti SC', serif; font-weight: 500; }
.empty-state p, .workspace-state p { color: var(--sub, var(--theme-text-secondary)); }
.state-actions { display: flex; justify-content: center; gap: 18px; margin-top: 24px; }
.state-actions button, .state-actions a { padding: 8px 0; border: 0; border-bottom: 1px solid currentColor; background: none; color: var(--accent, var(--theme-primary)); font: inherit; font-weight: 700; text-decoration: none; cursor: pointer; }
.nav-backdrop { display: none; }
@media (max-width: 900px) { .workspace-page { padding-top: 16px; } .workspace-shell { width: min(100% - 24px, 1180px); } .workspace-hero { padding: 22px 22px 26px; border-bottom: 1px solid var(--hairline, var(--theme-border)); } .back-link { margin-bottom: 20px; } .hero-copy { grid-template-columns: 64px minmax(0, 1fr); gap: 16px; } .workspace-icon { width: 64px; height: 64px; border-radius: 14px; font-size: 25px; } .hero-copy h1 { font-size: 38px; } .reading-layout { display: block; border-top: 0; } .tree-nav { position: fixed; inset: 0 auto 0 0; z-index: 31; width: min(340px, 88vw); max-height: none; padding: 25px 22px; border: 0; background: var(--surface, var(--theme-bg-card)); box-shadow: 12px 0 40px var(--shadow-color, var(--theme-shadow)); visibility: hidden; transform: translateX(-105%); transition: transform .2s ease, visibility .2s; } .tree-nav.open { visibility: visible; transform: translateX(0); } .nav-heading button { display: block; } .nav-backdrop { display: block; position: fixed; inset: 0; z-index: 30; border: 0; background: rgb(0 0 0 / 40%); } .workspace-content { padding: 24px 22px 40px; } .content-heading { align-items: center; } .content-heading > div { flex: 1; } .content-heading > span { display: none; } .mobile-nav-button { display: inline-flex; align-items: center; gap: 7px; padding: 8px 10px; border: 1px solid var(--hairline, var(--theme-border)); background: var(--surface, var(--theme-bg-card)); color: inherit; font: inherit; font-size: 12px; cursor: pointer; } }
@media (max-width: 560px) { .workspace-shell { width: calc(100% - 16px); } .workspace-hero { padding: 18px 16px 22px; } .hero-copy { grid-template-columns: 48px minmax(0, 1fr); gap: 12px; } .workspace-icon { width: 48px; height: 48px; border-radius: 10px; font-size: 20px; } .hero-copy h1 { font-size: 30px; } .description { grid-column: 1 / -1; font-size: 13px; } .workspace-content { padding: 20px 14px 32px; } .content-heading { gap: 12px; } .content-heading h2 { font-size: 23px; } .mobile-nav-button { padding: 7px 8px; } .document-entry { grid-template-columns: 12px minmax(0, 1fr) 18px; gap: 9px; min-height: 74px; padding-right: 6px; padding-left: 6px; } .document-entry time { display: none; } .entry-branch { width: 7px; height: 7px; } }
@media (prefers-reduced-motion: reduce) { .tree-nav, .document-entry, .catalog-toggle span { transition: none; } }
</style>
