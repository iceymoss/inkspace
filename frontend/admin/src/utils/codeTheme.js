// 代码主题工具函数（管理后台）
import adminApi from './adminApi'

let codeTheme = 'github' // 默认代码主题
let markdownTheme = 'light' // 默认 Markdown 主题
let themeLoaded = false

// 加载代码主题和 Markdown 主题配置
export const loadCodeTheme = async () => {
  if (themeLoaded) {
    return codeTheme
  }
  
  try {
    // 管理后台可以直接获取所有设置
    const response = await adminApi.get('/admin/settings')
    const settings = response.data || []
    const codeThemeSetting = settings.find(s => s.key === 'code_theme')
    if (codeThemeSetting && codeThemeSetting.value) {
      codeTheme = codeThemeSetting.value
    }
    const markdownThemeSetting = settings.find(s => s.key === 'markdown_theme')
    if (markdownThemeSetting && markdownThemeSetting.value) {
      markdownTheme = markdownThemeSetting.value
    }
    themeLoaded = true
    return codeTheme
  } catch (error) {
    console.warn('Failed to load code theme, using default:', error)
    // 如果失败，尝试从公开设置获取
    try {
      const publicResponse = await fetch('/api/settings/public')
      const publicSettings = await publicResponse.json()
      if (publicSettings.data) {
        if (publicSettings.data.code_theme) {
          codeTheme = publicSettings.data.code_theme
        }
        if (publicSettings.data.markdown_theme) {
          markdownTheme = publicSettings.data.markdown_theme
        }
      }
      themeLoaded = true
    } catch (e) {
      // 忽略错误，使用默认主题
      themeLoaded = true
    }
    return codeTheme
  }
}

// 获取 Markdown 主题
export const getMarkdownTheme = async () => {
  if (!themeLoaded) {
    await loadCodeTheme()
  }
  return markdownTheme
}

// 获取当前代码主题
export const getCodeTheme = () => {
  return codeTheme
}

// 设置代码主题（用于管理后台保存后刷新）
export const setCodeTheme = (theme) => {
  codeTheme = theme
  themeLoaded = true
}

// 设置 Markdown 主题（用于管理后台保存后刷新）
export const setMarkdownTheme = (theme) => {
  markdownTheme = theme
  themeLoaded = true
}

// 重置主题加载状态（用于刷新配置）
export const resetThemeCache = () => {
  themeLoaded = false
}

// 动态加载 highlight.js 主题 CSS
export const loadHighlightTheme = async (theme) => {
  if (!theme || theme === 'default') {
    return
  }
  
  // 检查是否已经加载
  const existingLink = document.getElementById(`hljs-theme-${theme}`)
  if (existingLink) {
    return Promise.resolve()
  }
  
  return new Promise((resolve, reject) => {
    try {
      // 动态导入 highlight.js 主题样式
      const link = document.createElement('link')
      link.id = `hljs-theme-${theme}`
      link.rel = 'stylesheet'
      link.href = `https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/styles/${theme}.min.css`
      
      // 等待 CSS 加载完成
      link.onload = () => {
        resolve()
      }
      link.onerror = () => {
        console.warn(`Failed to load highlight theme: ${theme}`)
        resolve() // 即使失败也继续，使用默认主题
      }
      
      document.head.appendChild(link)
    } catch (error) {
      console.warn(`Failed to load highlight theme: ${theme}`, error)
      resolve() // 即使失败也继续
    }
  })
}

