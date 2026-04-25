import api from './api'

let loadingEl = null
const startLoading = () => {
  loadingEl = document.createElement('div')
  loadingEl.className = 'fixed inset-0 z-[9999] flex items-center justify-center bg-black/20'
  loadingEl.innerHTML = '<div class="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>'
  document.body.appendChild(loadingEl)
}
const closeLoading = () => {
  if (loadingEl) {
    loadingEl.remove()
    loadingEl = null
  }
}

export async function navigateToWorkDetail(workId, routerInstance) {
  startLoading()

  try {
    const response = await api.get(`/works/${workId}`)
    const workData = response.data

    sessionStorage.setItem(`preloaded_work_${workId}`, JSON.stringify(workData))

    closeLoading()

    routerInstance.push(`/works/${workId}`)
  } catch (error) {
    closeLoading()
    routerInstance.push(`/works/${workId}`)
  }
}
