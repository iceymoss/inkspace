<template>
  <main class="wiki-index">
    <div class="index-shell">
      <header class="masthead">
        <div>
          <p class="kicker">
            INKSPACE · PUBLIC LIBRARY
          </p>
          <h1>知识刊物</h1>
        </div>
        <p class="declaration">
          从公开工作空间中，阅读被持续整理、校订与更新的知识。
        </p>
      </header>

      <section
        class="issue-heading"
        aria-labelledby="workspace-heading"
      >
        <h2 id="workspace-heading">
          本期目录
        </h2>
        <span
          v-if="!loading && !error"
          aria-live="polite"
        >共 {{ total }} 个知识库</span>
      </section>

      <div
        v-if="loading"
        class="workspace-grid"
        aria-busy="true"
        aria-label="正在加载知识库"
      >
        <article
          v-for="index in 6"
          :key="index"
          class="workspace-card skeleton-card"
        >
          <el-skeleton animated>
            <template #template>
              <el-skeleton-item
                variant="circle"
                class="skeleton-icon"
              />
              <el-skeleton-item
                variant="h1"
                style="width: 62%"
              />
              <el-skeleton-item variant="text" />
              <el-skeleton-item
                variant="text"
                style="width: 78%"
              />
            </template>
          </el-skeleton>
        </article>
      </div>

      <section
        v-else-if="error"
        class="state-panel"
        role="alert"
      >
        <span class="state-mark">!</span>
        <h2>目录暂时无法打开</h2>
        <p>{{ error }}</p>
        <button
          class="text-action"
          type="button"
          @click="loadWorkspaces"
        >
          重新加载
        </button>
      </section>

      <section
        v-else-if="workspaces.length === 0"
        class="state-panel"
      >
        <span class="state-mark">∅</span>
        <h2>本期尚未收录内容</h2>
        <p>公开知识库发布后，会在这里形成可阅读的目录。</p>
      </section>

      <section
        v-else
        class="workspace-grid"
        aria-label="公开知识库列表"
      >
        <article
          v-for="(workspace, index) in workspaces"
          :key="workspace.id"
          class="workspace-card"
        >
          <router-link
            class="workspace-main"
            :to="`/wiki/${workspace.id}`"
          >
            <div class="workspace-cover">
              <img
                v-if="isImageIcon(workspace.icon)"
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
            <div class="workspace-copy">
              <div class="card-topline">
                <span>VOL. {{ String((currentPage - 1) * pageSize + index + 1).padStart(2, '0') }}</span>
                <span>{{ workspace.doc_count }} 篇</span>
              </div>
              <h3>{{ workspace.name }}</h3>
              <p>{{ workspace.description || '一个持续生长的公开知识空间。' }}</p>
            </div>
          </router-link>
          <router-link
            v-if="workspace.author_id"
            class="workspace-author"
            :to="`/users/${workspace.author_id}`"
            :aria-label="`查看${workspace.author_name || '作者'}的个人主页`"
          >
            <el-avatar
              :size="30"
              :src="workspace.author_avatar"
            >
              {{ workspace.author_name?.slice(0, 1) || '作' }}
            </el-avatar>
            <span>{{ workspace.author_name || '匿名作者' }}</span>
          </router-link>
          <footer>
            <time :datetime="workspace.updated_at">更新于 {{ formatDate(workspace.updated_at) }}</time>
            <router-link
              class="read-link"
              :to="`/wiki/${workspace.id}`"
            >进入阅读 <span aria-hidden="true">→</span></router-link>
          </footer>
        </article>
      </section>

      <nav
        v-if="total > pageSize"
        class="pagination"
        aria-label="知识库分页"
      >
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          @current-change="handlePageChange"
        />
      </nav>
    </div>
  </main>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import dayjs from 'dayjs'
import api from '@/utils/api'

const workspaces = ref([])
const loading = ref(true)
const error = ref('')
const currentPage = ref(1)
const pageSize = 9
const total = ref(0)

const formatDate = value => value ? dayjs(value).format('YYYY.MM.DD') : '-'
const isImageIcon = value => /^(https?:\/\/|\/uploads\/)/.test(value || '')

async function loadWorkspaces() {
  loading.value = true
  error.value = ''
  try {
    const response = await api.get('/wiki/workspaces', {
      params: { page: currentPage.value, page_size: pageSize }
    })
    workspaces.value = response.data?.list || []
    total.value = response.data?.total || 0
  } catch (requestError) {
    workspaces.value = []
    error.value = requestError.response?.data?.message || '请检查网络连接后重试。'
  } finally {
    loading.value = false
  }
}

function handlePageChange() {
  loadWorkspaces()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(loadWorkspaces)
</script>

<style scoped>
.wiki-index { min-height: calc(100vh - 60px); background: var(--bg, var(--theme-bg-primary)); color: var(--ink, var(--theme-text-primary)); }
.index-shell { width: min(1180px, calc(100% - 40px)); margin: 0 auto; padding: clamp(42px, 7vw, 92px) 0 80px; }
.masthead { display: grid; grid-template-columns: minmax(0, 1.35fr) minmax(260px, .65fr); align-items: end; gap: 48px; padding-bottom: 30px; border-bottom: 3px double var(--ink, var(--theme-text-primary)); }
.kicker { margin: 0 0 12px; color: var(--accent, var(--theme-primary)); font-size: 11px; font-weight: 700; letter-spacing: .2em; }
.masthead h1 { margin: 0; font-family: Georgia, 'Songti SC', serif; font-size: clamp(54px, 9vw, 104px); font-weight: 500; letter-spacing: -.065em; line-height: .88; }
.declaration { max-width: 390px; margin: 0 0 3px; color: var(--sub, var(--theme-text-secondary)); font-family: Georgia, 'Songti SC', serif; font-size: 17px; line-height: 1.85; }
.issue-heading { display: flex; align-items: baseline; justify-content: space-between; margin: 28px 0 18px; padding-bottom: 10px; border-bottom: 1px solid var(--hairline, var(--theme-border)); }
.issue-heading h2 { margin: 0; font-size: 21px; font-weight: 500; }
.issue-heading span { color: var(--sub, var(--theme-text-tertiary)); font-size: 12px; }
.workspace-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 20px; }
.workspace-card { display: flex; min-height: 390px; overflow: hidden; border: 1px solid var(--hairline, var(--theme-border)); background: var(--surface, var(--theme-bg-card)); flex-direction: column; transition: border-color .25s ease, transform .25s ease; }
.workspace-card:hover { border-color: var(--accent, var(--theme-primary)); transform: translateY(-4px); }
.workspace-main { color: inherit; text-decoration: none; }
.workspace-main:focus-visible, .workspace-author:focus-visible, .read-link:focus-visible { outline: 2px solid var(--accent, var(--theme-primary)); outline-offset: -2px; }
.workspace-cover { height: 150px; overflow: hidden; background: var(--bg-soft, var(--theme-bg-secondary)); }
.workspace-cover img { width: 100%; height: 100%; object-fit: cover; transition: transform .35s ease; }
.workspace-card:hover .workspace-cover img { transform: scale(1.035); }
.cover-placeholder { display: flex; align-items: flex-end; justify-content: space-between; height: 100%; padding: 20px 22px; background: linear-gradient(135deg, color-mix(in srgb, var(--accent, var(--theme-primary)) 82%, #18243a), color-mix(in srgb, var(--accent, var(--theme-primary)) 30%, var(--bg-soft, var(--theme-bg-secondary)))); color: #fff; }
.cover-placeholder span { font-family: Georgia, 'Songti SC', serif; font-size: 44px; line-height: 1; }
.cover-placeholder small { font-size: 10px; font-weight: 700; letter-spacing: .2em; opacity: .72; }
.workspace-copy { padding: 18px 26px 0; }
.card-topline { display: flex; justify-content: space-between; color: var(--sub, var(--theme-text-tertiary)); font-size: 10px; font-weight: 700; letter-spacing: .12em; }
.workspace-card h3 { margin: 12px 0 8px; overflow: hidden; font-family: Georgia, 'Songti SC', serif; font-size: clamp(22px, 2.5vw, 28px); font-weight: 500; line-height: 1.25; text-overflow: ellipsis; white-space: nowrap; }
.workspace-card p { display: -webkit-box; margin: 0; overflow: hidden; color: var(--sub, var(--theme-text-secondary)); font-size: 14px; line-height: 1.65; -webkit-box-orient: vertical; -webkit-line-clamp: 2; }
.workspace-author { display: inline-flex; align-items: center; align-self: flex-start; gap: 9px; min-width: 0; margin: 17px 26px 0; color: var(--ink, var(--theme-text-primary)); font-size: 12px; text-decoration: none; }
.workspace-author span { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.workspace-author:hover { color: var(--accent, var(--theme-primary)); }
.workspace-card footer { display: flex; align-items: center; justify-content: space-between; margin: auto 26px 0; padding: 14px 0 22px; border-top: 1px solid var(--hairline, var(--theme-border)); color: var(--sub, var(--theme-text-tertiary)); font-size: 11px; }
.read-link { color: var(--accent, var(--theme-primary)); font-weight: 700; text-decoration: none; }
.skeleton-card { pointer-events: none; }
.skeleton-icon { display: block; width: 100%; height: 150px; margin-bottom: 24px; border-radius: 0; }
.state-panel { padding: 90px 20px; border-top: 1px solid var(--hairline, var(--theme-border)); border-bottom: 1px solid var(--hairline, var(--theme-border)); text-align: center; }
.state-mark { display: grid; place-items: center; width: 54px; height: 54px; margin: 0 auto 22px; border: 1px solid var(--accent, var(--theme-primary)); border-radius: 50%; color: var(--accent, var(--theme-primary)); font-family: Georgia, serif; font-size: 25px; }
.state-panel h2 { margin: 0 0 10px; font-size: 27px; }
.state-panel p { margin: 0 0 24px; color: var(--sub, var(--theme-text-secondary)); }
.text-action { padding: 8px 0; border: 0; border-bottom: 1px solid currentColor; background: transparent; color: var(--accent, var(--theme-primary)); font: inherit; font-weight: 700; cursor: pointer; }
.text-action:focus-visible { outline: 2px solid var(--accent, var(--theme-primary)); outline-offset: 4px; }
.pagination { display: flex; justify-content: center; margin-top: 36px; }
@media (max-width: 900px) { .workspace-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); } .masthead { grid-template-columns: 1fr; gap: 24px; } }
@media (max-width: 560px) { .index-shell { width: min(100% - 24px, 1180px); padding-top: 34px; } .masthead h1 { font-size: 58px; } .declaration { font-size: 15px; } .workspace-grid { grid-template-columns: 1fr; } .workspace-card { min-height: 380px; } }
@media (prefers-reduced-motion: reduce) { .workspace-card, .workspace-cover img { transition: none; } html { scroll-behavior: auto; } }
</style>
