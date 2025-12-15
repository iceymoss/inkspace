import api from './api'

// 主题配置 - 基于美学配色方案
const themes = {
  day: {
    name: '白天',
    cssVars: {
      '--theme-bg-primary': '#ffffff',
      '--theme-bg-secondary': '#f5f7fa',
      '--theme-bg-card': '#ffffff',
      '--theme-bg-hover': '#fafbfc',
      '--theme-text-primary': '#303133',
      '--theme-text-secondary': '#606266',
      '--theme-text-tertiary': '#909399',
      '--theme-border': '#dcdfe6',
      '--theme-border-light': '#ebeef5',
      '--theme-primary': '#409eff',
      '--theme-primary-hover': '#66b1ff',
      '--theme-shadow': 'rgba(0, 0, 0, 0.1)',
      '--theme-content-bg': '#ffffff',
      '--theme-hero-gradient': 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
      '--theme-accent': '#409eff'
    }
  },
  night: {
    name: '黑夜',
    cssVars: {
      '--theme-bg-primary': '#0d1117',
      '--theme-bg-secondary': '#161b22',
      '--theme-bg-card': '#21262d',
      '--theme-bg-hover': '#30363d',
      '--theme-text-primary': '#c9d1d9',
      '--theme-text-secondary': '#8b949e',
      '--theme-text-tertiary': '#6e7681',
      '--theme-border': '#30363d',
      '--theme-border-light': '#21262d',
      '--theme-primary': '#58a6ff',
      '--theme-primary-hover': '#79c0ff',
      '--theme-shadow': 'rgba(0, 0, 0, 0.5)',
      '--theme-content-bg': '#1c2128',
      '--theme-hero-gradient': 'linear-gradient(135deg, #1f2937 0%, #111827 100%)',
      '--theme-accent': '#58a6ff'
    }
  },
  holiday: {
    name: '节假日',
    cssVars: {
      '--theme-bg-primary': '#fff5f5',
      '--theme-bg-secondary': '#ffe8e8',
      '--theme-bg-card': '#ffffff',
      '--theme-bg-hover': '#fff0f0',
      '--theme-text-primary': '#8b1a1a',
      '--theme-text-secondary': '#a52a2a',
      '--theme-text-tertiary': '#b22222',
      '--theme-border': '#ffb3b3',
      '--theme-border-light': '#ffcccc',
      '--theme-primary': '#ff3333',
      '--theme-primary-hover': '#ff5555',
      '--theme-shadow': 'rgba(255, 51, 51, 0.25)',
      '--theme-content-bg': '#fffafa',
      '--theme-hero-gradient': 'linear-gradient(135deg, #ff4444 0%, #cc0000 100%)',
      '--theme-accent': '#ff3333'
    }
  },
  mourning: {
    name: '哀悼日',
    cssVars: {
      '--theme-bg-primary': '#2d2d2d',
      '--theme-bg-secondary': '#1f1f1f',
      '--theme-bg-card': '#3a3a3a',
      '--theme-bg-hover': '#454545',
      '--theme-text-primary': '#d4d4d4',
      '--theme-text-secondary': '#b8b8b8',
      '--theme-text-tertiary': '#9c9c9c',
      '--theme-border': '#4d4d4d',
      '--theme-border-light': '#3a3a3a',
      '--theme-primary': '#999999',
      '--theme-primary-hover': '#aaaaaa',
      '--theme-shadow': 'rgba(0, 0, 0, 0.5)',
      '--theme-content-bg': '#404040',
      '--theme-hero-gradient': 'linear-gradient(135deg, #3a3a3a 0%, #2d2d2d 100%)',
      '--theme-accent': '#999999'
    }
  }
}

// 应用主题
export function applyTheme(themeName) {
  const theme = themes[themeName] || themes.day
  const root = document.documentElement
  
  // 应用CSS变量
  Object.keys(theme.cssVars).forEach(key => {
    root.style.setProperty(key, theme.cssVars[key])
  })
  
  // 添加主题类名到body
  document.body.className = document.body.className.replace(/theme-\w+/g, '')
  document.body.classList.add(`theme-${themeName}`)
  
  // 保存到localStorage
  localStorage.setItem('site_theme', themeName)
}

// 获取主题设置
export async function loadTheme() {
  try {
    const response = await api.get('/settings/public')
    const themeName = response.data?.site_theme || 'day'
    
    // 如果是节假日主题，应用自定义颜色
    if (themeName === 'holiday' && response.data) {
      const holidayPrimary = response.data.holiday_primary || themes.holiday.cssVars['--theme-primary']
      const holidayTheme = {
        ...themes.holiday,
        cssVars: {
          ...themes.holiday.cssVars,
          '--theme-bg-primary': response.data.holiday_bg_primary || themes.holiday.cssVars['--theme-bg-primary'],
          '--theme-bg-secondary': response.data.holiday_bg_secondary || themes.holiday.cssVars['--theme-bg-secondary'],
          '--theme-text-primary': response.data.holiday_text_primary || themes.holiday.cssVars['--theme-text-primary'],
          '--theme-primary': holidayPrimary,
          '--theme-primary-hover': adjustBrightness(holidayPrimary, 20),
          // 根据主色调动态生成渐变背景
          '--theme-hero-gradient': generateGradientFromColor(holidayPrimary)
        }
      }
      applyCustomTheme('holiday', holidayTheme.cssVars)
    } else {
      applyTheme(themeName)
    }
    return themeName
  } catch (error) {
    console.error('Failed to load theme:', error)
    // 如果API失败，尝试从localStorage读取
    const savedTheme = localStorage.getItem('site_theme')
    if (savedTheme && themes[savedTheme]) {
      applyTheme(savedTheme)
      return savedTheme
    }
    // 默认使用白天主题
    applyTheme('day')
    return 'day'
  }
}

// 调整颜色亮度（用于生成hover颜色）
function adjustBrightness(color, percent) {
  const num = parseInt(color.replace('#', ''), 16)
  const amt = Math.round(2.55 * percent)
  const R = Math.min(255, Math.max(0, (num >> 16) + amt))
  const G = Math.min(255, Math.max(0, ((num >> 8) & 0x00FF) + amt))
  const B = Math.min(255, Math.max(0, (num & 0x0000FF) + amt))
  return '#' + (0x1000000 + R * 0x10000 + G * 0x100 + B).toString(16).slice(1)
}

// 根据主色调生成渐变背景
function generateGradientFromColor(primaryColor) {
  if (!primaryColor) {
    return 'linear-gradient(135deg, #ff4444 0%, #cc0000 100%)'
  }
  
  // 将主色调稍微调亮作为起始色
  const startColor = adjustBrightness(primaryColor, 10)
  // 将主色调稍微调暗作为结束色
  const endColor = adjustBrightness(primaryColor, -20)
  
  return `linear-gradient(135deg, ${startColor} 0%, ${endColor} 100%)`
}

// 应用自定义主题
function applyCustomTheme(themeName, customVars) {
  const root = document.documentElement
  
  // 应用自定义CSS变量
  Object.keys(customVars).forEach(key => {
    root.style.setProperty(key, customVars[key])
  })
  
  // 添加主题类名到body
  document.body.className = document.body.className.replace(/theme-\w+/g, '')
  document.body.classList.add(`theme-${themeName}`)
  
  // 保存到localStorage
  localStorage.setItem('site_theme', themeName)
}

// 初始化主题（在应用启动时调用）
export function initTheme() {
  // 先从localStorage读取，避免闪烁
  const savedTheme = localStorage.getItem('site_theme')
  if (savedTheme && themes[savedTheme]) {
    applyTheme(savedTheme)
  }
  
  // 然后从服务器加载最新设置
  loadTheme()
}

export default {
  themes,
  applyTheme,
  loadTheme,
  initTheme
}

