import { ElLoading } from 'element-plus'
import api from './api'

/**
 * 预加载作品数据并跳转到详情页
 * @param {number} workId - 作品ID
 * @param {object} routerInstance - Vue Router 实例（从 useRouter() 获取）
 */
export async function navigateToWorkDetail(workId, routerInstance) {
  // 显示全屏加载动画
  const loadingInstance = ElLoading.service({
    lock: true,
    text: '加载中...',
    background: 'rgba(0, 0, 0, 0.7)'
  })

  try {
    // 预加载作品数据
    const response = await api.get(`/works/${workId}`)
    const workData = response.data

    // 将数据存储到 sessionStorage，供详情页使用
    sessionStorage.setItem(`preloaded_work_${workId}`, JSON.stringify(workData))

    // 关闭加载动画
    loadingInstance.close()

    // 跳转到详情页
    routerInstance.push(`/works/${workId}`)
  } catch (error) {
    // 关闭加载动画
    loadingInstance.close()
    // 即使加载失败，也跳转到详情页（详情页会处理错误）
    routerInstance.push(`/works/${workId}`)
  }
}

