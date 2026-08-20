<template>
  <section
    v-loading="pageLoading"
    class="workspace-detail"
  >
    <div class="workspace-shell">
      <header class="workspace-header">
        <div class="workspace-heading">
          <el-button
            text
            circle
            aria-label="返回知识库"
            @click="$router.push('/dashboard/workspaces')"
          >
            <el-icon><ArrowLeft /></el-icon>
          </el-button>
          <span class="workspace-cover-thumb">
            <img
              v-if="isCoverImage(store.currentWorkspace?.icon)"
              :src="store.currentWorkspace.icon"
              :alt="`${store.currentWorkspace.name}封面`"
            >
            <template v-else>{{ store.currentWorkspace?.icon || store.currentWorkspace?.name?.slice(0, 1) || '知' }}</template>
          </span>
          <div>
            <div class="workspace-title-line">
              <h1>{{ store.currentWorkspace?.name || '工作空间' }}</h1>
              <el-tag
                :type="store.currentWorkspace?.is_public ? 'success' : 'info'"
                size="small"
                effect="plain"
              >
                {{ store.currentWorkspace?.is_public ? '公开 Wiki' : '私有空间' }}
              </el-tag>
            </div>
            <p>{{ store.currentWorkspace?.description || '整理目录，沉淀你的知识。' }}</p>
          </div>
        </div>
        <div class="workspace-header-actions">
          <el-button
            v-if="store.currentWorkspace?.is_public"
            @click="router.push(`/wiki/${workspaceId}`)"
          >
            查看 Wiki
          </el-button>
          <el-button
            v-if="canEditWorkspace"
            :loading="uploadState.active"
            @click="chooseFile"
          >
            <el-icon><Upload /></el-icon>上传文件
          </el-button>
          <el-button
            v-if="canEditWorkspace"
            type="primary"
            @click="createDoc"
          >
            <el-icon><Plus /></el-icon>新建文件
          </el-button>
          <input
            ref="fileInputRef"
            class="file-input"
            type="file"
            :accept="acceptedFileTypes"
            @change="uploadFile"
          >
        </div>
      </header>
      <button
        v-if="catalogOpen"
        class="catalog-backdrop"
        aria-label="关闭目录"
        @click="catalogOpen = false"
      />
      <aside :class="['catalog-panel', { open: catalogOpen }]">
        <div class="panel-title">
          <div>
            <strong>知识树</strong>
            <small>{{ store.treeDocs.length }} 篇文档</small>
          </div>
          <div class="panel-actions">
            <el-button
              v-if="canEditWorkspace"
              text
              circle
              title="新建根目录"
              aria-label="新建根目录"
              @click="openCatalogCreate(null)"
            >
              <el-icon><FolderAdd /></el-icon>
            </el-button>
            <el-button
              class="catalog-close"
              text
              circle
              aria-label="关闭目录"
              @click="catalogOpen = false"
            >
              <el-icon><Close /></el-icon>
            </el-button>
          </div>
        </div>
        <button
          :class="['root-row', { active: selectedCatalogId === null }]"
          @click="selectCatalog(null)"
        >
          <el-icon><Files /></el-icon><span>根目录</span>
        </button>
        <el-tree
          ref="treeRef"
          class="catalog-tree"
          node-key="id"
          :data="knowledgeTree"
          :props="treeProps"
          default-expand-all
          :draggable="canEditWorkspace"
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
              <el-dropdown
                v-if="data.type === 'catalog' && canEditWorkspace"
                trigger="click"
                @command="command => handleCatalogCommand(command, data.catalog)"
                @click.stop
              >
                <el-button
                  text
                  circle
                  size="small"
                  :aria-label="`${data.name}目录操作`"
                >
                  <el-icon><MoreFilled /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="create">
                      新建子目录
                    </el-dropdown-item>
                    <el-dropdown-item command="rename">
                      重命名
                    </el-dropdown-item>
                    <el-dropdown-item
                      command="delete"
                      divided
                    >
                      删除
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
              <span
                v-else-if="data.doc.article_id"
                class="blog-link-dot"
                title="已发布到博客"
              />
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
            >
              <el-icon><Menu /></el-icon>
            </el-button>
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
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
            <template #append>
              <el-button
                aria-label="搜索"
                @click="runSearch"
              >
                <el-icon><Search /></el-icon>
              </el-button>
            </template>
          </el-input>
        </div>

        <div
          v-if="uploadState.active"
          class="upload-strip"
          role="status"
          aria-live="polite"
        >
          <div class="upload-file-mark">
            {{ uploadState.extension }}
          </div>
          <div>
            <strong>正在上传 {{ uploadState.fileName }}</strong>
            <span>文件上传完成后会进入只读详情，预览转换将在后台继续。</span>
          </div>
          <el-progress
            :percentage="uploadState.progress"
            :stroke-width="6"
          />
        </div>

        <div
          v-loading="docsLoading"
          class="doc-list"
        >
          <article
            v-for="(doc, index) in store.docs"
            :key="doc.id"
            :class="['doc-row', { dragging: draggedDocIndex === index }]"
            role="link"
            tabindex="0"
            @dragover.prevent
            @drop.stop="dropDoc(index)"
            @click="viewDoc(doc.id)"
            @keyup.enter="viewDoc(doc.id)"
            @keyup.space.prevent="viewDoc(doc.id)"
          >
            <span
              v-if="!searching && canEditWorkspace"
              class="drag-handle"
              draggable="true"
              title="拖拽排序"
              @dragstart.stop="startDocDrag(index)"
              @dragend="draggedDocIndex = null"
              @click.stop
            ><el-icon><Rank /></el-icon></span>
            <div class="doc-mark">
              <el-icon><Document /></el-icon>
            </div>
            <div class="doc-main">
              <div class="doc-title">
                <strong>{{ doc.title }}</strong>
                <el-tag
                  :type="doc.status === 1 ? 'success' : 'info'"
                  size="small"
                  effect="plain"
                >
                  {{ doc.status === 1 ? '已公开' : '私有' }}
                </el-tag>
                <el-tag
                  v-if="doc.article_id"
                  type="success"
                  size="small"
                >
                  已发布到博客
                </el-tag>
              </div>
              <p>{{ doc.summary || doc.excerpt || contentExcerpt(doc.content) || '暂无内容' }}</p>
              <span>{{ doc.word_count || 0 }} 字 · 更新于 {{ formatTime(doc.updated_at) }}</span>
            </div>
            <span
              class="doc-actions"
              @click.stop
              @keydown.stop
            >
              <el-dropdown
                v-if="canEditWorkspace"
                trigger="click"
                @command="command => handleDocCommand(command, doc)"
              >
                <el-button
                  text
                  circle
                  :aria-label="`${doc.title}文档操作`"
                ><el-icon><MoreFilled /></el-icon></el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item
                      v-if="isMarkdownDocument(doc)"
                      command="wiki-publish"
                    >{{ doc.status === 1 ? '重新公开' : '公开' }}</el-dropdown-item>
                    <el-dropdown-item
                      v-if="isMarkdownDocument(doc)"
                      command="publish"
                    >{{ doc.article_id ? '更新到博客' : '发布到博客' }}</el-dropdown-item>
                    <el-dropdown-item
                      v-if="doc.editable"
                      command="edit"
                    >编辑</el-dropdown-item>
                    <el-dropdown-item
                      v-if="isMarkdownDocument(doc)"
                      command="share"
                    >分享</el-dropdown-item>
                    <el-dropdown-item
                      command="delete"
                      divided
                    >删除</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </span>
          </article>
        </div>

        <el-empty
          v-if="!docsLoading && store.docs.length === 0"
          :description="searching ? '没有找到匹配的文档' : '这里还没有文档'"
        >
          <el-button
            v-if="!searching && canEditWorkspace"
            type="primary"
            @click="createDoc"
          >
            新建文件
          </el-button>
          <el-button
            v-if="!searching && canEditWorkspace"
            @click="chooseFile"
          >
            上传文件
          </el-button>
          <el-button
            v-else
            @click="clearSearch"
          >
            清除搜索
          </el-button>
        </el-empty>
      </main>
    </div>

    <el-dialog
      v-model="createDialog.visible"
      title="新建文本文件"
      width="min(520px, 94vw)"
    >
      <el-form @submit.prevent="submitCreateDoc">
        <el-form-item label="文件名">
          <el-input
            v-model="createDialog.fileName"
            maxlength="200"
            autofocus
            placeholder="例如 README.md、styles.css、main.rs"
            @keyup.enter="submitCreateDoc"
          />
        </el-form-item>
        <div class="file-type-result" :class="{ invalid: createDialog.fileName && !createFileType }">
          <template v-if="createFileType">
            <strong>{{ createFileType.label }}</strong>
            <span>kind={{ createFileType.kind }} · language={{ createFileType.language }}</span>
          </template>
          <template v-else-if="createDialog.fileName">
            该扩展名不能直接创建，请使用“上传文件”保留原件。
          </template>
          <template v-else>
            输入文件名后将自动选择编辑器和语法高亮。
          </template>
        </div>
        <div class="file-presets">
          <button
            v-for="name in createFilePresets"
            :key="name"
            type="button"
            @click="createDialog.fileName = name"
          >
            {{ name }}
          </button>
        </div>
      </el-form>
      <template #footer>
        <el-button @click="createDialog.visible = false">取消</el-button>
        <el-button
          type="primary"
          :disabled="!createFileType"
          :loading="createDialog.loading"
          @click="submitCreateDoc"
        >
          创建并编辑
        </el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="catalogDialog.visible"
      :title="catalogDialog.mode === 'rename' ? '重命名目录' : '新建目录'"
      width="min(440px, 92vw)"
    >
      <el-form @submit.prevent="submitCatalog">
        <el-form-item label="目录名称">
          <el-input
            v-model="catalogDialog.name"
            maxlength="100"
            autofocus
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="catalogDialog.visible = false">
          取消
        </el-button>
        <el-button
          type="primary"
          :loading="catalogDialog.loading"
          @click="submitCatalog"
        >
          确定
        </el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="moveDialog.visible"
      title="移动文档"
      width="min(440px, 92vw)"
    >
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
        <el-button @click="moveDialog.visible = false">
          取消
        </el-button>
        <el-button
          type="primary"
          :loading="moveDialog.loading"
          @click="submitDocMove"
        >
          移动
        </el-button>
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
  ArrowLeft, Close, Document, Files, Folder, FolderAdd, Menu, MoreFilled, Plus, Rank, Search, Upload
} from '@element-plus/icons-vue'
import { useWorkspaceStore } from '@/stores/workspace'
import api from '@/utils/api'
import { inferCreatableFileType, textFileAccept } from '@/utils/knowledgeFileTypes'

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
const fileInputRef = ref(null)
const treeProps = { label: 'name', children: 'children' }
const catalogDialog = reactive({ visible: false, loading: false, mode: 'create', id: null, parentId: null, name: '' })
const moveDialog = reactive({ visible: false, loading: false, docId: null, catalogId: 0 })
const draggedDocIndex = ref(null)
const uploadState = reactive({ active: false, progress: 0, fileName: '', extension: 'FILE' })
const createDialog = reactive({ visible: false, loading: false, fileName: 'README.md' })
const createFilePresets = ['README.md', 'notes.txt', 'styles.css', 'app.js', 'main.go', 'worker.rs', 'config.json', 'config.yaml', 'data.csv']

const acceptedFileTypes = `${textFileAccept},.pdf,.jpg,.jpeg,.png,.gif,.webp,.svg,.doc,.docx,.odt,.rtf,.xls,.xlsx,.ods,.ppt,.pptx,.odp,.zip,.rar,.7z,.tar,.gz,.tgz`

const listTitle = computed(() => searching.value ? `“${searchQuery.value}”的搜索结果` : selectedCatalogName.value)
const canEditWorkspace = computed(() => store.currentWorkspace?.capabilities?.can_edit === true)
const catalogOptions = computed(() => [{ id: 0, name: '根目录', children: store.catalogs }])
const knowledgeTree = computed(() => buildKnowledgeTree(store.catalogs, store.treeDocs))
const createFileType = computed(() => inferCreatableFileType(createDialog.fileName))
const formatTime = value => value ? dayjs(value).format('YYYY-MM-DD HH:mm') : '刚刚'
const isCoverImage = value => /^(https?:\/\/|\/uploads\/)/.test(value || '')
const contentExcerpt = content => (content || '').replace(/[#>*`\[\]()_-]/g, '').slice(0, 100)
const isMarkdownDocument = doc => (doc?.kind || 'markdown') === 'markdown'

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
    viewDoc(node.doc.id)
    return
  }
  selectCatalog(node.catalogId)
}

const allowTreeDrag = node => canEditWorkspace.value && node.data.type === 'catalog'
const allowTreeDrop = (draggingNode, dropNode) => canEditWorkspace.value && dropNode.data.type === 'catalog'

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
	Object.assign(createDialog, { visible: true, loading: false, fileName: 'README.md' })
}

async function submitCreateDoc() {
	const fileName = createDialog.fileName.trim()
	const fileType = inferCreatableFileType(fileName)
	if (!fileType || createDialog.loading) return
	createDialog.loading = true
	try {
		const response = await store.createDoc({
			workspace_id: workspaceId,
			catalog_id: selectedCatalogId.value,
			title: fileName,
			file_name: fileName,
			content: '',
			kind: fileType.kind,
			language: fileType.language
		})
		createDialog.visible = false
		await store.fetchTreeDocs(workspaceId)
		router.push(`/dashboard/docs/${response.id}/edit`)
	} finally {
		createDialog.loading = false
	}
}

function chooseFile() {
  if (!uploadState.active) fileInputRef.value?.click()
}

async function uploadFile(event) {
  const file = event.target.files?.[0]
  event.target.value = ''
  if (!file) return
  if (file.size > 100 * 1024 * 1024) {
    ElMessage.warning('单个文件不能超过 100 MiB')
    return
  }

  const form = new FormData()
  form.append('file', file)
  if (selectedCatalogId.value !== null) form.append('catalog_id', String(selectedCatalogId.value))
  if (globalThis.crypto?.randomUUID) form.append('request_id', globalThis.crypto.randomUUID())

  Object.assign(uploadState, {
    active: true,
    progress: 0,
    fileName: file.name,
    extension: file.name.includes('.') ? file.name.split('.').pop().toUpperCase() : 'FILE'
  })
  try {
    const response = await api.post(`/workspaces/${workspaceId}/files`, form, {
      onUploadProgress: ({ loaded, total }) => {
        uploadState.progress = total ? Math.min(100, Math.round(loaded * 100 / total)) : 0
      }
    })
    uploadState.progress = 100
    await Promise.all([loadDocs(), store.fetchTreeDocs(workspaceId)])
    ElMessage.success('文件已上传')
    router.push(`/dashboard/docs/${response.data.id}`)
  } finally {
    uploadState.active = false
  }
}

const editDoc = (id, action) => router.push({
  path: `/dashboard/docs/${id}/edit`,
  query: action ? { action } : undefined
})
const viewDoc = id => router.push(`/dashboard/docs/${id}`)

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
  if (command === 'wiki-publish') return editDoc(doc.id, 'wiki-publish')
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
.workspace-title-line, .workspace-header-actions { display: flex; align-items: center; gap: 9px; }
.file-input { position: absolute; width: 1px; height: 1px; overflow: hidden; clip-path: inset(50%); white-space: nowrap; }
.workspace-heading p { margin: 2px 0 0; color: var(--theme-text-tertiary); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.workspace-cover-thumb { display: grid; place-items: center; flex: none; width: 64px; height: 44px; overflow: hidden; border: 1px solid var(--theme-border); border-radius: 10px; background: linear-gradient(135deg, color-mix(in srgb, var(--theme-primary) 80%, #18243a), color-mix(in srgb, var(--theme-primary) 25%, var(--theme-bg-card))); color: #fff; font-size: 20px; font-weight: 700; box-shadow: 0 4px 12px var(--theme-shadow); }
.workspace-cover-thumb img { width: 100%; height: 100%; object-fit: cover; }
.workspace-shell { display: grid; grid-template-columns: 260px minmax(0, 1fr); grid-template-rows: auto minmax(0, 1fr); min-height: calc(100vh - 112px); background: transparent; border: 0; border-top: 1px solid var(--theme-border); border-bottom: 1px solid var(--theme-border); border-radius: 0; overflow: hidden; box-shadow: none; }
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
.upload-strip { display: grid; grid-template-columns: auto minmax(0, 1fr) minmax(120px, 220px); align-items: center; gap: 14px; margin: 16px 0 4px; padding: 12px 14px; border: 1px solid color-mix(in srgb, var(--theme-primary) 35%, var(--theme-border)); background: color-mix(in srgb, var(--theme-primary) 6%, var(--theme-bg-card)); }
.upload-strip > div:nth-child(2) { display: flex; min-width: 0; flex-direction: column; }
.upload-strip strong { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-size: 13px; }
.upload-strip span { margin-top: 2px; color: var(--theme-text-tertiary); font-size: 11px; }
.upload-file-mark { display: grid; width: 42px; height: 48px; place-items: center; border: 1px solid var(--theme-primary); color: var(--theme-primary); font: 700 9px/1 'SFMono-Regular', Consolas, monospace; letter-spacing: .05em; }
.file-type-result { display: flex; min-height: 54px; flex-direction: column; justify-content: center; padding: 10px 12px; border-left: 3px solid var(--theme-primary); background: color-mix(in srgb, var(--theme-primary) 7%, var(--theme-bg-secondary)); color: var(--theme-text-secondary); font-size: 12px; }
.file-type-result strong { color: var(--theme-text-primary); font-size: 14px; }
.file-type-result span { margin-top: 2px; color: var(--theme-text-tertiary); font: 10px/1.5 'SFMono-Regular', Consolas, monospace; }
.file-type-result.invalid { border-left-color: var(--el-color-danger); color: var(--el-color-danger); }
.file-presets { display: flex; flex-wrap: wrap; gap: 7px; margin-top: 14px; }
.file-presets button { padding: 5px 8px; border: 1px solid var(--theme-border); background: var(--theme-bg-card); color: var(--theme-text-secondary); font: 11px/1.2 'SFMono-Regular', Consolas, monospace; cursor: pointer; }
.file-presets button:hover, .file-presets button:focus-visible { border-color: var(--theme-primary); color: var(--theme-primary); outline: none; }
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
@media (max-width: 700px) { .workspace-header { align-items: flex-start; padding: 14px; } .workspace-heading { gap: 9px; } .workspace-heading p { white-space: normal; } .workspace-cover-thumb { width: 52px; } .workspace-header-actions { align-items: stretch; flex-direction: column; flex: none; } .workspace-header-actions .el-button { margin-left: 0; } .workspace-shell { display: block; position: relative; min-height: calc(100vh - 90px); overflow: hidden; } .catalog-backdrop { display: block; position: absolute; inset: 0; z-index: 19; padding: 0; border: 0; background: rgb(0 0 0 / 38%); } .catalog-panel { display: block; position: absolute; inset: 0 auto 0 0; z-index: 20; width: min(290px, 86vw); box-shadow: 8px 0 24px var(--theme-shadow); transform: translateX(-105%); visibility: hidden; transition: transform .2s ease, visibility .2s; } .catalog-panel.open { transform: translateX(0); visibility: visible; } .catalog-close { display: inline-flex; } .mobile-catalog-button { display: inline-flex; } .docs-panel { padding: 14px; } .docs-toolbar { align-items: stretch; flex-direction: column; } .upload-strip { grid-template-columns: auto minmax(0, 1fr); } .upload-strip .el-progress { grid-column: 1 / -1; } .search-input { width: 100%; } .doc-row { gap: 10px; } .drag-handle { display: none; } .doc-main p { display: none; } }
</style>
