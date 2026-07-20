<template>
  <section v-loading="pageLoading" class="workspace-detail">
    <header class="workspace-header">
      <div class="workspace-heading">
        <el-button text circle aria-label="返回知识库" @click="$router.push('/dashboard/workspaces')">
          <el-icon><ArrowLeft /></el-icon>
        </el-button>
        <span class="workspace-icon">{{ store.currentWorkspace?.icon || '知' }}</span>
        <div>
          <h1>{{ store.currentWorkspace?.name || '工作空间' }}</h1>
          <p>{{ store.currentWorkspace?.description || '整理目录，沉淀你的知识。' }}</p>
        </div>
      </div>
      <el-button type="primary" @click="createDoc"><el-icon><Plus /></el-icon>新建文档</el-button>
    </header>

    <div class="workspace-shell">
      <aside :class="['catalog-panel', { open: catalogOpen }]">
        <div class="panel-title">
          <strong>目录</strong>
          <el-button text circle title="新建根目录" @click="openCatalogCreate(null)"><el-icon><FolderAdd /></el-icon></el-button>
        </div>
        <button :class="['root-row', { active: selectedCatalogId === null }]" @click="selectCatalog(null)">
          <el-icon><Files /></el-icon><span>根目录</span>
        </button>
        <el-tree
          ref="treeRef"
          class="catalog-tree"
          node-key="id"
          :data="store.catalogs"
          :props="treeProps"
          default-expand-all
          draggable
          :expand-on-click-node="false"
          @node-click="node => selectCatalog(node.id)"
          @node-drop="handleCatalogDrop"
        >
          <template #default="{ data }">
            <div :class="['tree-node', { active: selectedCatalogId === data.id }]">
              <span class="node-label"><el-icon><Folder /></el-icon>{{ data.name }}</span>
              <el-dropdown trigger="click" @command="command => handleCatalogCommand(command, data)" @click.stop>
                <el-button text circle size="small"><el-icon><MoreFilled /></el-icon></el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="create">新建子目录</el-dropdown-item>
                    <el-dropdown-item command="rename">重命名</el-dropdown-item>
                    <el-dropdown-item command="delete" divided>删除</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </template>
        </el-tree>
      </aside>

      <main class="docs-panel">
        <div class="docs-toolbar">
          <div class="title-line">
            <el-button class="mobile-catalog-button" text circle @click="catalogOpen = !catalogOpen"><el-icon><Menu /></el-icon></el-button>
            <div>
              <h2>{{ listTitle }}</h2>
              <span>{{ store.docs.length }} 篇文档</span>
            </div>
          </div>
          <el-input
            v-model="searchQuery"
            clearable
            class="search-input"
            placeholder="搜索标题或内容"
            @keyup.enter="runSearch"
            @clear="clearSearch"
          >
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
        </div>

        <div v-loading="docsLoading" class="doc-list">
          <article
            v-for="(doc, index) in store.docs"
            :key="doc.id"
            class="doc-row"
            :draggable="!searching"
            @dragstart="startDocDrag(index)"
            @dragover.prevent
            @drop.stop="dropDoc(index)"
            @click="editDoc(doc.id)"
          >
            <div class="doc-mark"><el-icon><Document /></el-icon></div>
            <div class="doc-main">
              <div class="doc-title">
                <strong>{{ doc.title }}</strong>
                <el-tag :type="doc.status === 1 ? 'success' : 'info'" size="small">
                  {{ doc.status === 1 ? '已发布' : '草稿' }}
                </el-tag>
              </div>
              <p>{{ doc.summary || doc.excerpt || contentExcerpt(doc.content) || '暂无内容' }}</p>
              <span>{{ doc.word_count || 0 }} 字 · 更新于 {{ formatTime(doc.updated_at) }}</span>
            </div>
            <el-dropdown trigger="click" @command="command => handleDocCommand(command, doc)" @click.stop>
              <el-button text circle><el-icon><MoreFilled /></el-icon></el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="edit">编辑</el-dropdown-item>
                  <el-dropdown-item command="move">移动到目录</el-dropdown-item>
                  <el-dropdown-item command="delete" divided>删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
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
          <el-input v-model="catalogDialog.name" maxlength="100" autofocus @keyup.enter="submitCatalog" />
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
  ArrowLeft, Document, Files, Folder, FolderAdd, Menu, MoreFilled, Plus, Search
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
const formatTime = value => value ? dayjs(value).format('YYYY-MM-DD HH:mm') : '刚刚'
const contentExcerpt = content => (content || '').replace(/[#>*`\[\]()_-]/g, '').slice(0, 100)

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

const editDoc = id => router.push(`/dashboard/docs/${id}/edit`)

function openCatalogCreate(parentId) {
  Object.assign(catalogDialog, { visible: true, mode: 'create', id: null, parentId, name: '' })
}

function openCatalogRename(catalog) {
  Object.assign(catalogDialog, { visible: true, mode: 'rename', id: catalog.id, parentId: null, name: catalog.name })
}

async function submitCatalog() {
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
  if (deletedIds.includes(selectedCatalogId.value)) await selectCatalog(null)
  else await loadDocs()
  ElMessage.success('目录已删除')
}

async function handleCatalogDrop(draggingNode, dropNode, dropType) {
  const parentNode = dropType === 'inner' ? dropNode : dropNode.parent
  const parentId = parentNode.level === 0 ? null : parentNode.data.id
  const siblings = parentNode.childNodes
  try {
    await Promise.all(siblings.map((node, sort) => store.moveCatalog(
      node.data.id,
      { parent_id: parentId, sort },
      workspaceId,
      false
    )))
    await store.fetchCatalogs(workspaceId)
    ElMessage.success('目录位置已更新')
  } catch {
    await store.fetchCatalogs(workspaceId)
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
    await loadDocs()
  } catch {
    await loadDocs()
  }
}

async function handleDocCommand(command, doc) {
  if (command === 'edit') return editDoc(doc.id)
  if (command === 'move') {
    Object.assign(moveDialog, { visible: true, docId: doc.id, catalogId: doc.catalog_id || 0 })
    return
  }
  await ElMessageBox.confirm(`确定删除文档“${doc.title}”吗？`, '删除文档', {
    type: 'warning', confirmButtonText: '确认删除', cancelButtonText: '取消'
  })
  await store.deleteDoc(doc.id)
  ElMessage.success('文档已删除')
}

async function submitDocMove() {
  moveDialog.loading = true
  try {
    await store.moveDoc(moveDialog.docId, { catalog_id: moveDialog.catalogId || null, sort: 0 })
    moveDialog.visible = false
    await loadDocs()
    ElMessage.success('文档已移动')
  } finally {
    moveDialog.loading = false
  }
}

onMounted(async () => {
  try {
    await Promise.all([store.fetchWorkspace(workspaceId), store.fetchCatalogs(workspaceId)])
    await loadDocs()
  } catch {
    router.replace('/dashboard/workspaces')
  } finally {
    pageLoading.value = false
  }
})
</script>

<style scoped>
.workspace-detail { max-width: 1380px; min-height: calc(100vh - 100px); margin: 0 auto; color: var(--theme-text-primary); }
.workspace-header { display: flex; align-items: center; justify-content: space-between; gap: 20px; margin-bottom: 18px; }
.workspace-heading { display: flex; align-items: center; gap: 13px; min-width: 0; }
.workspace-heading h1 { margin: 0; font-size: 25px; }
.workspace-heading p { margin: 2px 0 0; color: var(--theme-text-tertiary); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.workspace-icon { display: grid; place-items: center; flex: none; width: 44px; height: 44px; border-radius: 12px; background: color-mix(in srgb, var(--theme-primary) 15%, var(--theme-bg-card)); color: var(--theme-primary); font-size: 20px; font-weight: 700; }
.workspace-shell { display: grid; grid-template-columns: 260px minmax(0, 1fr); min-height: calc(100vh - 180px); background: var(--theme-bg-card); border: 1px solid var(--theme-border); border-radius: 14px; overflow: hidden; box-shadow: 0 8px 30px var(--theme-shadow); }
.catalog-panel { padding: 18px 12px; border-right: 1px solid var(--theme-border); background: color-mix(in srgb, var(--theme-bg-secondary) 65%, var(--theme-bg-card)); overflow: auto; }
.panel-title { display: flex; align-items: center; justify-content: space-between; padding: 0 8px 10px; }
.root-row { display: flex; align-items: center; gap: 8px; width: 100%; padding: 8px 10px; border: 0; border-radius: 7px; background: transparent; color: var(--theme-text-secondary); cursor: pointer; }
.root-row.active, .root-row:hover { background: color-mix(in srgb, var(--theme-primary) 12%, transparent); color: var(--theme-primary); }
.catalog-tree { background: transparent; color: var(--theme-text-primary); }
.tree-node { display: flex; align-items: center; justify-content: space-between; flex: 1; min-width: 0; padding-right: 2px; border-radius: 6px; }
.node-label { display: flex; align-items: center; gap: 6px; min-width: 0; overflow: hidden; text-overflow: ellipsis; }
.tree-node.active { color: var(--theme-primary); font-weight: 600; }
:deep(.el-tree-node__content) { height: 38px; border-radius: 7px; }
:deep(.el-tree-node__content:hover) { background: var(--theme-bg-hover); }
.docs-panel { min-width: 0; padding: 22px 28px; }
.docs-toolbar { display: flex; align-items: center; justify-content: space-between; gap: 20px; padding-bottom: 18px; border-bottom: 1px solid var(--theme-border-light); }
.title-line { display: flex; align-items: center; gap: 8px; }
.title-line h2 { margin: 0; font-size: 21px; }
.title-line span { color: var(--theme-text-tertiary); font-size: 12px; }
.search-input { width: min(320px, 45%); }
.mobile-catalog-button { display: none; }
.doc-list { min-height: 120px; }
.doc-row { display: flex; align-items: center; gap: 15px; padding: 18px 8px; border-bottom: 1px solid var(--theme-border-light); cursor: pointer; transition: background .2s; }
.doc-row:hover { background: var(--theme-bg-hover); }
.doc-mark { display: grid; place-items: center; flex: none; width: 40px; height: 40px; border-radius: 10px; background: color-mix(in srgb, var(--theme-primary) 10%, var(--theme-bg-card)); color: var(--theme-primary); }
.doc-main { flex: 1; min-width: 0; }
.doc-title { display: flex; align-items: center; gap: 9px; }
.doc-main p { margin: 5px 0; color: var(--theme-text-tertiary); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.doc-main > span { color: var(--theme-text-tertiary); font-size: 12px; }
@media (max-width: 900px) { .workspace-shell { grid-template-columns: 220px minmax(0, 1fr); } .docs-panel { padding: 18px; } }
@media (max-width: 700px) { .workspace-header { align-items: flex-start; } .workspace-heading p { white-space: normal; } .workspace-shell { display: block; position: relative; } .catalog-panel { display: none; position: absolute; inset: 0 auto 0 0; z-index: 20; width: min(280px, 85vw); box-shadow: 8px 0 24px var(--theme-shadow); } .catalog-panel.open { display: block; } .mobile-catalog-button { display: inline-flex; } .docs-panel { padding: 14px; } .docs-toolbar { align-items: stretch; flex-direction: column; } .search-input { width: 100%; } .doc-row { gap: 10px; } .doc-main p { display: none; } }
</style>
