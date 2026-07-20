import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import api from '@/utils/api'

export const useWorkspaceStore = defineStore('workspace', () => {
  const workspaces = ref([])
  const currentWorkspace = ref(null)
  const catalogs = ref([])
  const docs = ref([])
  const loading = ref(false)

  const workspaceCount = computed(() => workspaces.value.length)

  async function fetchWorkspaces() {
    loading.value = true
    try {
      const response = await api.get('/workspaces')
      workspaces.value = Array.isArray(response.data) ? response.data : response.data?.list || []
      return workspaces.value
    } finally {
      loading.value = false
    }
  }

  async function createWorkspace(payload) {
    const response = await api.post('/workspaces', payload)
    workspaces.value.push(response.data)
    return response.data
  }

  async function updateWorkspace(id, payload) {
    const response = await api.put(`/workspaces/${id}`, payload)
    const workspace = response.data || { ...payload, id }
    const index = workspaces.value.findIndex(item => item.id === Number(id))
    if (index >= 0) workspaces.value[index] = workspace
    if (currentWorkspace.value?.id === Number(id)) currentWorkspace.value = workspace
    return workspace
  }

  async function deleteWorkspace(id) {
    await api.delete(`/workspaces/${id}`)
    workspaces.value = workspaces.value.filter(item => item.id !== Number(id))
    if (currentWorkspace.value?.id === Number(id)) resetCurrent()
  }

  async function fetchWorkspace(id) {
    const response = await api.get(`/workspaces/${id}`)
    currentWorkspace.value = response.data
    return response.data
  }

  async function fetchCatalogs(workspaceId) {
    const response = await api.get(`/workspaces/${workspaceId}/catalogs`)
    catalogs.value = response.data || []
    return catalogs.value
  }

  async function createCatalog(workspaceId, payload) {
    const response = await api.post(`/workspaces/${workspaceId}/catalogs`, payload)
    await fetchCatalogs(workspaceId)
    return response.data
  }

  async function renameCatalog(id, name, workspaceId) {
    await api.put(`/catalogs/${id}`, { name })
    await fetchCatalogs(workspaceId)
  }

  async function deleteCatalog(id, workspaceId) {
    await api.delete(`/catalogs/${id}`)
    await fetchCatalogs(workspaceId)
  }

  async function moveCatalog(id, payload, workspaceId, refresh = true) {
    await api.put(`/catalogs/${id}/move`, payload)
    if (refresh) await fetchCatalogs(workspaceId)
  }

  async function fetchDocs(workspaceId, catalogId = null) {
    const response = await api.get(`/workspaces/${workspaceId}/docs`, {
      params: catalogId === null ? {} : { catalog_id: catalogId }
    })
    const list = Array.isArray(response.data) ? response.data : response.data?.list || []
    docs.value = catalogId === null ? list.filter(item => item.catalog_id == null) : list
    return docs.value
  }

  async function searchDocs(workspaceId, query) {
    const response = await api.get(`/workspaces/${workspaceId}/search`, { params: { q: query } })
    docs.value = Array.isArray(response.data) ? response.data : response.data?.list || []
    return docs.value
  }

  async function createDoc(payload) {
    const response = await api.post('/docs', payload)
    return response.data
  }

  async function deleteDoc(id) {
    await api.delete(`/docs/${id}`)
    docs.value = docs.value.filter(item => item.id !== Number(id))
  }

  async function moveDoc(id, payload) {
    await api.put(`/docs/${id}/move`, payload)
  }

  function resetCurrent() {
    currentWorkspace.value = null
    catalogs.value = []
    docs.value = []
  }

  return {
    workspaces,
    currentWorkspace,
    catalogs,
    docs,
    loading,
    workspaceCount,
    fetchWorkspaces,
    createWorkspace,
    updateWorkspace,
    deleteWorkspace,
    fetchWorkspace,
    fetchCatalogs,
    createCatalog,
    renameCatalog,
    deleteCatalog,
    moveCatalog,
    fetchDocs,
    searchDocs,
    createDoc,
    deleteDoc,
    moveDoc,
    resetCurrent
  }
})
