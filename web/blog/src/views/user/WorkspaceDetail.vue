<template>
  <section v-loading="pageLoading" class="workspace-detail">
    <div class="workspace-shell">
      <header class="workspace-header">
        <div class="workspace-heading">
          <el-button text circle aria-label="返回知识库" @click="$router.push('/dashboard/workspaces')">
            <el-icon><ArrowLeft /></el-icon>
          </el-button>
          <span class="workspace-cover-thumb">
            <img
              v-if="isCoverImage(store.currentWorkspace?.icon)"
              :src="store.currentWorkspace.icon"
              :alt="`${store.currentWorkspace.name}封面`"
            />
            <template v-else>{{ store.currentWorkspace?.icon || store.currentWorkspace?.name?.slice(0, 1) || '知' }}</template>
          </span>
          <div>
            <h1>{{ store.currentWorkspace?.name || '工作空间' }}</h1>
            <p>{{ store.currentWorkspace?.description || '整理目录，沉淀你的知识。' }}</p>
          </div>
        </div>
        <el-button type="primary" @click="createDoc"><el-icon><Plus /></el-icon>新建文档</el-button>
      </header>
      <button
        v-if="catalogOpen"
        class="catalog-backdrop"
        aria-label="关闭目录"
        @click="catalogOpen = false"
      ></button>
      <aside :class="['catalog-panel', { open: catalogOpen }]">
        <div class="panel-title">
          <div>
            <strong>知识树</strong>
            <small>{{ store.treeDocs.length }} 篇文档</small>
          </div>
          <div class="panel-actions">
            <el-button text circle title="新建根目录" aria-label="新建根目录" @click="openCatalogCreate(null)"><el-icon><FolderAdd /></el-icon></el-button>
            <el-button class="catalog-close" text circle aria-label="关闭目录" @click="catalogOpen = false"><el-icon><Close /></el-icon></el-button>
          </div>
        </div>
        <button :class="['root-row', { active: selectedCatalogId === null }]" @click="selectCatalog(null)">
          <el-icon><Files /></el-icon><span>根目录</span>
        </button>
        <el-tree
          ref="treeRef"
          class="catalog-tree"
          node-key="id"
          :data="knowledgeTree"
          :props="treeProps"
          default-expand-all
          draggable
          :allow-drag="allowTreeDrag"
          :allow-drop="allowTreeDrop"
          :expand-on-click-node="false"
          @node-click="handleTreeNodeClick"
          @node-drop="handleCatalogDrop"
        >
          <template #default="{ data }">
            <div :class="['tree-node', data.type, { active: data.type === 'catalog' && selectedCatalogId === data.catalogId }]">
              <span class="node-label">
                <el-icon><Folder v-if="data.type === 'catalog'" /><Document v-else /></el-icon>
                <span>{{ data.name }}</span>
              </span>
              <el-dropdown v-if="data.type === 'catalog'" trigger="click" @command="command => handleCatalogCommand(command, data.catalog)" @click.stop>
                <el-button text circle size="small" :aria-label="`${data.name}目录操作`"><el-icon><MoreFilled /></el-icon></el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="create">新建子目录</el-dropdown-item>
                    <el-dropdown-item command="rename">重命名</el-dropdown-item>
                    <el-dropdown-item command="delete" divided>删除</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
              <span v-else-if="data.doc.article_id" class="blog-link-dot" title="已发布到博客"></span>
            </div>
          </template>
        </el-tree>
      </aside>

      <main class="docs-panel">
        <div class="docs-toolbar">
          <div class="title-line">
            <el-button
              class="mobile-catalog-button"
              text
              circle
              aria-label="打开目录"
              :aria-expanded="catalogOpen"
              @click="catalogOpen = !catalogOpen"
            ><el-icon><Menu /></el-icon></el-button>
            <div>
              <h2>{{ listTitle }}</h2>
              <span>{{ store.docs.length }} 篇文档{{ searching ? ' · 搜索范围：整个知识库' : '' }}</span>
            </div>
          </div>
          <el-input
            v-model="searchQuery"
            clearable
            class="search-input"
            placeholder="搜索整个知识库"
            @keyup.enter="runSearch"
            @clear="clearSearch"
          >
            <template #prefix><el-icon><Search /></el-icon></template>
            <template #append>
              <el-button aria-label="搜索" @click="runSearch"><el-icon><Search /></el-icon></el-button>
            </template>
          </el-input>
        </div>

        <div v-loading="docsLoading" class="doc-list">
          <article
            v-for="(doc, index) in store.docs"
            :key="doc.id"
            :class="['doc-row', { dragging: draggedDocIndex === index }]"
            role="link"
            tabindex="0"
            @dragover.prevent
            @drop.stop="dropDoc(index)"
            @click="editDoc(doc.id)"
            @keyup.enter="editDoc(doc.id)"
            @keyup.space.prevent="editDoc(doc.id)"
          >
            <span
              v-if="!searching"
              class="drag-handle"
              draggable="true"
              title="拖拽排序"
              @dragstart.stop="startDocDrag(index)"
              @dragend="draggedDocIndex = null"
              @click.stop
            ><el-icon><Rank /></el-icon></span>
            <div class="doc-mark"><el-icon><Document /></el-icon></div>
            <div class="doc-main">
              <div class="doc-title">
                <strong>{{ doc.title }}</strong>
                <el-tag v-if="doc.article_id" type="success" size="small">已发布到博客</el-tag>
              </div>
              <p>{{ doc.summary || doc.excerpt || contentExcerpt(doc.content) || '暂无内容' }}</p>
              <span>{{ doc.word_count || 0 }} 字 · 更新于 {{ formatTime(doc.updated_at) }}</span>
            </div>
            <span class="doc-actions" @click.stop @keydown.stop>
              <el-dropdown trigger="click" @command="command => handleDocCommand(command, doc)">
                <el-button text circle :aria-label="`${doc.title}文档操作`"><el-icon><MoreFilled /></el-icon></el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="publish">{{ doc.article_id ? '更新到博客' : '发布到博客' }}</el-dropdown-item>
                    <el-dropdown-item command="edit">编辑</el-dropdown-item>
                    <el-dropdown-item command="share">分享</el-dropdown-item>
                    <el-dropdown-item command="delete" divided>删除</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </span>
          </article>
        </div>

        <el-empty v-if="!docsLoading && store.docs.length === 0" :description="searching ? '没有找到匹配的文档' : '这里还没有文档'">
          <el-button v-if="!searching" type="primary" @click="createDoc">新建文档</el-button>
          <el-button v-else @click="clearSearch">清除搜索</el-button>
        </el-empty>
      </main>
    </div>

    <el-dialog v-model="catalogDialog.visible" :title="catalogDialog.mode === 'rename' ? '重命名目录' : '新建目录'" width="min(440px, 92vw)">
      <el-form @submit.prevent="submitCatalog">
        <el-form-item label="目录名称">
          <el-input v-model="catalogDialog.name" maxlength="100" autofocus />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="catalogDialog.visible = false">取消</el-button>
        <el-button type="primary" :loading="catalogDialog.loading" @click="submitCatalog">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="moveDialog.visible" title="移动文档" width="min(440px, 92vw)">
      <el-tree-select
        v-model="moveDialog.catalogId"
        :data="catalogOptions"
        :props="treeProps"
        node-key="id"
        check-strictly
        default-expand-all
        style="width: 100%"
      />
      <template #footer>
        <el-button @click="moveDialog.visible = false">取消</el-button>
        <el-button type="primary" :loading="moveDialog.loading" @click="submitDocMove">移动</el-button>
      </template>
    </el-dialog>
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  ArrowLeft, Close, Document, Files, Folder, FolderAdd, Menu, MoreFilled, Plus, Rank, Search
} from '@element-plus/icons-vue'
import { useWorkspaceStore } from '@/stores/workspace'

const route = useRoute()
const router = useRouter()
const store = useWorkspaceStore()
const workspaceId = Number(route.params.id)
const pageLoading = ref(true)
const docsLoading = ref(false)
const catalogOpen = ref(false)
const selectedCatalogId = ref(null)
const selectedCatalogName = ref('根目录')
const searchQuery = ref('')
const searching = ref(false)
const treeRef = ref(null)
const treeProps = { label: 'name', children: 'children' }
const catalogDialog = reactive({ visible: false, loading: false, mode: 'create', id: null, parentId: null, name: '' })
const moveDialog = reactive({ visible: false, loading: false, docId: null, catalogId: 0 })
const draggedDocIndex = ref(null)

const listTitle = computed(() => searching.value ? `“${searchQuery.value}”的搜索结果` : selectedCatalogName.value)
const catalogOptions = computed(() => [{ id: 0, name: '根目录', children: store.catalogs }])
const knowledgeTree = computed(() => buildKnowledgeTree(store.catalogs, store.treeDocs))
const formatTime = value => value ? dayjs(value).format('YYYY-MM-DD HH:mm') : '刚刚'
const isCoverImage = value => /^(https?:\/\/|\/uploads\/)/.test(value || '')
const contentExcerpt = content => (content || '').replace(/[#>*`\[\]()_-]/g, '').slice(0, 100)

function buildKnowledgeTree(catalogs, docs) {
  const docsByCatalog = new Map()
  for (const doc of docs) {
    const catalogId = doc.catalog_id ?? null
    if (!docsByCatalog.has(catalogId)) docsByCatalog.set(catalogId, [])
    docsByCatalog.get(catalogId).push(doc)
  }
  const docNodes = catalogId => (docsByCatalog.get(catalogId) || []).map(doc => ({
    id: `doc-${doc.id}`,
    type: 'doc',
    name: doc.title || '无标题文档',
    doc,
    children: []
  }))
  const catalogNodes = nodes => nodes.map(catalog => ({
    id: `catalog-${catalog.id}`,
    type: 'catalog',
    name: catalog.name,
    catalogId: catalog.id,
    catalog,
    children: [...catalogNodes(catalog.children || []), ...docNodes(catalog.id)]
  }))
  return [...catalogNodes(catalogs), ...docNodes(null)]
}

function handleTreeNodeClick(node) {
  if (node.type === 'doc') {
    editDoc(node.doc.id)
    return
  }
  selectCatalog(node.catalogId)
}

const allowTreeDrag = node => node.data.type === 'catalog'
const allowTreeDrop = (draggingNode, dropNode) => dropNode.data.type === 'catalog'

async function loadDocs() {
  docsLoading.value = true
  try {
    await store.fetchDocs(workspaceId, selectedCatalogId.value)
  } finally {
    docsLoading.value = false
  }
}

function findCatalog(nodes, id) {
  for (const node of nodes) {
    if (node.id === id) return node
    const found = findCatalog(node.children || [], id)
    if (found) return found
  }
  return null
}

async function selectCatalog(id) {
  selectedCatalogId.value = id
  selectedCatalogName.value = id === null ? '根目录' : findCatalog(store.catalogs, id)?.name || '目录'
  searching.value = false
  searchQuery.value = ''
  catalogOpen.value = false
  await loadDocs()
}

async function runSearch() {
  const query = searchQuery.value.trim()
  if (!query) return clearSearch()
  searching.value = true
  docsLoading.value = true
  try {
    await store.searchDocs(workspaceId, query)
  } finally {
    docsLoading.value = false
  }
}

async function clearSearch() {
  searchQuery.value = ''
  searching.value = false
  await loadDocs()
}

async function createDoc() {
  const response = await store.createDoc({
    workspace_id: workspaceId,
    catalog_id: selectedCatalogId.value,
    title: '无标题文档',
    content: ''
  })
  router.push(`/dashboard/docs/${response.id}/edit`)
}

const editDoc = (id, action) => router.push({
  path: `/dashboard/docs/${id}/edit`,
  query: action ? { action } : undefined
})

function openCatalogCreate(parentId) {
  Object.assign(catalogDialog, { visible: true, mode: 'create', id: null, parentId, name: '' })
}

function openCatalogRename(catalog) {
  Object.assign(catalogDialog, { visible: true, mode: 'rename', id: catalog.id, parentId: null, name: catalog.name })
}

async function submitCatalog() {
  if (catalogDialog.loading) return
  const name = catalogDialog.name.trim()
  if (!name) return ElMessage.warning('请输入目录名称')
  catalogDialog.loading = true
  try {
    if (catalogDialog.mode === 'rename') {
      await store.renameCatalog(catalogDialog.id, name, workspaceId)
      if (selectedCatalogId.value === catalogDialog.id) selectedCatalogName.value = name
      ElMessage.success('目录已重命名')
    } else {
      await store.createCatalog(workspaceId, { parent_id: catalogDialog.parentId, name, sort: 0 })
      ElMessage.success('目录已创建')
    }
    catalogDialog.visible = false
  } finally {
    catalogDialog.loading = false
  }
}

async function handleCatalogCommand(command, catalog) {
  if (command === 'create') return openCatalogCreate(catalog.id)
  if (command === 'rename') return openCatalogRename(catalog)
  await ElMessageBox.confirm(
    `删除“${catalog.name}”会同时删除其子目录和文档。此操作不可恢复，是否继续？`,
    '删除目录',
    { type: 'warning', confirmButtonText: '确认删除', cancelButtonText: '取消' }
  )
  const deletedIds = collectCatalogIds(catalog)
  await store.deleteCatalog(catalog.id, workspaceId)
  await store.fetchTreeDocs(workspaceId)
  if (deletedIds.includes(selectedCatalogId.value)) await selectCatalog(null)
  else await loadDocs()
  ElMessage.success('目录已删除')
}

async function handleCatalogDrop(draggingNode, dropNode, dropType) {
  const parentNode = dropType === 'inner' ? dropNode : dropNode.parent
  const parentId = parentNode.level === 0 ? null : parentNode.data.catalogId
  const siblings = parentNode.childNodes.filter(node => node.data.type === 'catalog')
  try {
    await Promise.all(siblings.map((node, sort) => store.moveCatalog(
      node.data.catalogId,
      { parent_id: parentId, sort },
      workspaceId,
      false
    )))
    await store.fetchCatalogs(workspaceId)
    ElMessage.success('目录位置已更新')
  } catch {
    await store.fetchCatalogs(workspaceId)
    ElMessage.error('目录位置保存失败，已恢复原顺序')
  }
}

function collectCatalogIds(catalog) {
  return [catalog.id, ...(catalog.children || []).flatMap(collectCatalogIds)]
}

function startDocDrag(index) {
  draggedDocIndex.value = index
}

async function dropDoc(targetIndex) {
  const sourceIndex = draggedDocIndex.value
  draggedDocIndex.value = null
  if (sourceIndex === null || sourceIndex === targetIndex || searching.value) return
  const reordered = [...store.docs]
  const [moved] = reordered.splice(sourceIndex, 1)
  reordered.splice(targetIndex, 0, moved)
  store.docs = reordered
  try {
    await Promise.all(reordered.map((doc, sort) => store.moveDoc(doc.id, {
      catalog_id: selectedCatalogId.value,
      sort
    })))
    await Promise.all([loadDocs(), store.fetchTreeDocs(workspaceId)])
  } catch {
    await Promise.all([loadDocs(), store.fetchTreeDocs(workspaceId)])
    ElMessage.error('文档排序保存失败，已恢复原顺序')
  }
}

async function handleDocCommand(command, doc) {
  if (command === 'publish') return editDoc(doc.id, 'publish')
  if (command === 'edit') return editDoc(doc.id)
  if (command === 'share') return editDoc(doc.id, 'share')
  try {
    await ElMessageBox.confirm(`确定删除文档“${doc.title}”吗？`, '删除文档', {
      type: 'warning', confirmButtonText: '确认删除', cancelButtonText: '取消'
    })
  } catch {
    return
  }
  await store.deleteDoc(doc.id)
  await store.fetchTreeDocs(workspaceId)
  ElMessage.success('文档已删除')
}

async function submitDocMove() {
  moveDialog.loading = true
  try {
    await store.moveDoc(moveDialog.docId, { catalog_id: moveDialog.catalogId || null, sort: 0 })
    moveDialog.visible = false
    await Promise.all([loadDocs(), store.fetchTreeDocs(workspaceId)])
    ElMessage.success('文档已移动')
  } finally {
    moveDialog.loading = false
  }
}

onMounted(async () => {
  try {
    await Promise.all([
      store.fetchWorkspace(workspaceId),
      store.fetchCatalogs(workspaceId),
      store.fetchTreeDocs(workspaceId)
    ])
    await loadDocs()
  } catch {
    router.replace('/dashboard/workspaces')
  } finally {
    pageLoading.value = false
  }
})
</script>

<style scoped>
.workspace-detail { max-width: 1380px; min-height: calc(100vh - 100px); margin: -4px auto 0; color: var(--theme-text-primary); }
.workspace-header { display: flex; grid-column: 1 / -1; align-items: center; justify-content: space-between; gap: 20px; padding: 16px 20px; border-bottom: 1px solid var(--theme-border); background: color-mix(in srgb, var(--theme-primary) 3%, var(--theme-bg-card)); }
.workspace-heading { display: flex; align-items: center; gap: 13px; min-width: 0; }
.workspace-heading h1 { margin: 0; font-size: 25px; }
.workspace-heading p { margin: 2px 0 0; color: var(--theme-text-tertiary); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.workspace-cover-thumb { display: grid; place-items: center; flex: none; width: 64px; height: 44px; overflow: hidden; border: 1px solid var(--theme-border); border-radius: 10px; background: linear-gradient(135deg, color-mix(in srgb, var(--theme-primary) 80%, #18243a), color-mix(in srgb, var(--theme-primary) 25%, var(--theme-bg-card))); color: #fff; font-size: 20px; font-weight: 700; box-shadow: 0 4px 12px var(--theme-shadow); }
.workspace-cover-thumb img { width: 100%; height: 100%; object-fit: cover; }
.workspace-shell { display: grid; grid-template-columns: 260px minmax(0, 1fr); grid-template-rows: auto minmax(0, 1fr); min-height: calc(100vh - 112px); background: var(--theme-bg-card); border: 1px solid var(--theme-border); border-radius: 14px; overflow: hidden; box-shadow: 0 8px 30px var(--theme-shadow); }
.catalog-backdrop { display: none; }
.catalog-panel { grid-column: 1; grid-row: 2; padding: 18px 12px; border-right: 1px solid var(--theme-border); background: color-mix(in srgb, var(--theme-bg-secondary) 65%, var(--theme-bg-card)); overflow: auto; }
.panel-title { display: flex; align-items: center; justify-content: space-between; padding: 0 8px 10px; }
.panel-title > div:first-child { display: flex; flex-direction: column; }
.panel-title small { margin-top: 1px; color: var(--theme-text-tertiary); font-size: 11px; font-weight: 400; }
.panel-actions { display: flex; align-items: center; }
.catalog-close { display: none; }
.root-row { display: flex; align-items: center; gap: 8px; width: 100%; padding: 8px 10px; border: 0; border-radius: 7px; background: transparent; color: var(--theme-text-secondary); cursor: pointer; }
.root-row.active, .root-row:hover { background: color-mix(in srgb, var(--theme-primary) 12%, transparent); color: var(--theme-primary); }
.catalog-tree { background: transparent; color: var(--theme-text-primary); }
.tree-node { display: flex; align-items: center; justify-content: space-between; flex: 1; min-width: 0; padding-right: 2px; border-radius: 6px; }
.node-label { display: flex; align-items: center; gap: 6px; min-width: 0; overflow: hidden; }
.node-label span { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.tree-node.active { color: var(--theme-primary); font-weight: 600; }
.tree-node.doc { padding-right: 9px; color: var(--theme-text-secondary); }
.tree-node.doc .node-label .el-icon { color: var(--theme-text-tertiary); }
.blog-link-dot { flex: none; width: 6px; height: 6px; border-radius: 50%; background: var(--el-color-success); opacity: .9; }
:deep(.el-tree-node__content) { height: 38px; border-radius: 7px; }
:deep(.el-tree-node__content:hover) { background: var(--theme-bg-hover); }
.docs-panel { grid-column: 2; grid-row: 2; min-width: 0; padding: 22px 28px; }
.docs-toolbar { display: flex; align-items: center; justify-content: space-between; gap: 20px; padding-bottom: 18px; border-bottom: 1px solid var(--theme-border-light); }
.title-line { display: flex; align-items: center; gap: 8px; }
.title-line h2 { margin: 0; font-size: 21px; }
.title-line span { color: var(--theme-text-tertiary); font-size: 12px; }
.search-input { width: min(320px, 45%); }
.mobile-catalog-button { display: none; }
.doc-list { min-height: 120px; }
.doc-row { display: flex; align-items: center; gap: 13px; padding: 18px 8px; border-bottom: 1px solid var(--theme-border-light); cursor: pointer; transition: background .2s, opacity .2s; }
.doc-row:hover, .doc-row:focus-visible { background: var(--theme-bg-hover); outline: none; }
.doc-row.dragging { opacity: .45; }
.drag-handle { display: grid; place-items: center; flex: none; width: 24px; height: 34px; color: var(--theme-text-tertiary); cursor: grab; opacity: .35; transition: color .2s, opacity .2s; }
.drag-handle:active { cursor: grabbing; }
.doc-row:hover .drag-handle { color: var(--theme-primary); opacity: 1; }
.doc-mark { display: grid; place-items: center; flex: none; width: 40px; height: 40px; border-radius: 10px; background: color-mix(in srgb, var(--theme-primary) 10%, var(--theme-bg-card)); color: var(--theme-primary); }
.doc-main { flex: 1; min-width: 0; }
.doc-title { display: flex; align-items: center; gap: 9px; }
.doc-main p { margin: 5px 0; color: var(--theme-text-tertiary); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.doc-main > span { color: var(--theme-text-tertiary); font-size: 12px; }
.doc-actions { flex: none; }
@media (max-width: 900px) { .workspace-shell { grid-template-columns: 220px minmax(0, 1fr); } .docs-panel { padding: 18px; } }
@media (max-width: 700px) { .workspace-header { align-items: flex-start; padding: 14px; } .workspace-heading { gap: 9px; } .workspace-heading p { white-space: normal; } .workspace-cover-thumb { width: 52px; } .workspace-header > .el-button { flex: none; } .workspace-shell { display: block; position: relative; min-height: calc(100vh - 90px); overflow: hidden; } .catalog-backdrop { display: block; position: absolute; inset: 0; z-index: 19; padding: 0; border: 0; background: rgb(0 0 0 / 38%); } .catalog-panel { display: block; position: absolute; inset: 0 auto 0 0; z-index: 20; width: min(290px, 86vw); box-shadow: 8px 0 24px var(--theme-shadow); transform: translateX(-105%); visibility: hidden; transition: transform .2s ease, visibility .2s; } .catalog-panel.open { transform: translateX(0); visibility: visible; } .catalog-close { display: inline-flex; } .mobile-catalog-button { display: inline-flex; } .docs-panel { padding: 14px; } .docs-toolbar { align-items: stretch; flex-direction: column; } .search-input { width: 100%; } .doc-row { gap: 10px; } .drag-handle { display: none; } .doc-main p { display: none; } }
</style>
