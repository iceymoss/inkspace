// 代码主题工具函数
import api from './api'

let codeTheme = 'github' // 默认代码主题
let markdownTheme = 'light' // 默认 Markdown 主题
let themeLoaded = false

// 加载代码主题和 Markdown 主题配置
export const loadCodeTheme = async () => {
  if (themeLoaded) {
    return codeTheme
  }
  
  try {
    const response = await api.get('/settings/public')
    const settings = response.data || {}
    if (settings.code_theme) {
      codeTheme = settings.code_theme
      console.log('加载代码主题配置:', codeTheme)
    }
    if (settings.markdown_theme) {
      markdownTheme = settings.markdown_theme
      console.log('加载Markdown主题配置:', markdownTheme)
    }
    themeLoaded = true
    return codeTheme
  } catch (error) {
    console.warn('Failed to load code theme, using default:', error)
    themeLoaded = true // 即使失败也标记为已加载，避免重复请求
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
    // 如果样式表已存在，检查是否已完全加载
    if (existingLink.sheet) {
      // 样式表已加载，等待一小段时间确保样式应用
      return new Promise((resolve) => {
        setTimeout(() => resolve(), 100)
      })
    }
    // 如果样式表存在但还没完全加载，等待加载完成
    return new Promise((resolve) => {
      if (existingLink.complete || existingLink.readyState === 'complete') {
        setTimeout(() => resolve(), 100)
      } else {
        existingLink.onload = () => {
          setTimeout(() => resolve(), 100)
        }
        existingLink.onerror = () => resolve()
      }
    })
  }
  
  // 移除之前加载的其他主题（如果有）
  const oldLinks = document.querySelectorAll('link[id^="hljs-theme-"]')
  oldLinks.forEach(link => {
    if (link.id !== `hljs-theme-${theme}`) {
      link.remove()
    }
  })
  
  return new Promise((resolve, reject) => {
    try {
      // 动态导入 highlight.js 主题样式
      const link = document.createElement('link')
      link.id = `hljs-theme-${theme}`
      link.rel = 'stylesheet'
      link.href = `https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/styles/${theme}.min.css`
      
      // 等待 CSS 加载完成
      link.onload = () => {
        console.log(`Highlight.js theme loaded: ${theme}`)
        // 确保样式应用，等待一小段时间让浏览器处理
        setTimeout(() => {
          resolve()
        }, 100)
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

