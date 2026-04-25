import api from './api'

const themes = {
  day: {
    name: '白天',
    cssVars: {
      '--theme-bg-primary': '#fafafa',
      '--theme-bg-secondary': '#f4f4f5',
      '--theme-bg-card': '#ffffff',
      '--theme-bg-hover': '#f4f4f5',
      '--theme-text-primary': '#09090b',
      '--theme-text-secondary': '#52525b',
      '--theme-text-tertiary': '#71717a',
      '--theme-border': '#d4d4d8',
      '--theme-border-light': '#e4e4e7',
      '--theme-primary': '#18181b',
      '--theme-primary-hover': '#27272a',
      '--theme-shadow': 'rgba(0, 0, 0, 0.1)',
      '--theme-content-bg': '#ffffff',
      '--theme-hero-gradient': 'linear-gradient(135deg, #18181b 0%, #3f3f46 100%)',
      '--theme-accent': '#ec4899'
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

export function applyTheme(themeName) {
  const theme = themes[themeName] || themes.day
  const root = document.documentElement

  Object.keys(theme.cssVars).forEach(key => {
    root.style.setProperty(key, theme.cssVars[key])
  })

  document.body.className = document.body.className.replace(/theme-\w+/g, '')
  document.body.classList.add(`theme-${themeName}`)

  localStorage.setItem('site_theme', themeName)
}

export async function loadTheme() {
  try {
    const response = await api.get('/settings/public')
    const themeName = response.data?.site_theme || 'day'

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
    const savedTheme = localStorage.getItem('site_theme')
    if (savedTheme && themes[savedTheme]) {
      applyTheme(savedTheme)
      return savedTheme
    }
    applyTheme('day')
    return 'day'
  }
}

function adjustBrightness(color, percent) {
  const num = parseInt(color.replace('#', ''), 16)
  const amt = Math.round(2.55 * percent)
  const R = Math.min(255, Math.max(0, (num >> 16) + amt))
  const G = Math.min(255, Math.max(0, ((num >> 8) & 0x00FF) + amt))
  const B = Math.min(255, Math.max(0, (num & 0x0000FF) + amt))
  return '#' + (0x1000000 + R * 0x10000 + G * 0x100 + B).toString(16).slice(1)
}

function generateGradientFromColor(primaryColor) {
  if (!primaryColor) {
    return 'linear-gradient(135deg, #ff4444 0%, #cc0000 100%)'
  }

  const startColor = adjustBrightness(primaryColor, 10)
  const endColor = adjustBrightness(primaryColor, -20)

  return `linear-gradient(135deg, ${startColor} 0%, ${endColor} 100%)`
}

function applyCustomTheme(themeName, customVars) {
  const root = document.documentElement

  Object.keys(customVars).forEach(key => {
    root.style.setProperty(key, customVars[key])
  })

  document.body.className = document.body.className.replace(/theme-\w+/g, '')
  document.body.classList.add(`theme-${themeName}`)

  localStorage.setItem('site_theme', themeName)
}

export function initTheme() {
  const savedTheme = localStorage.getItem('site_theme')
  if (savedTheme && themes[savedTheme]) {
    applyTheme(savedTheme)
  }

  loadTheme()
}

export default {
  themes,
  applyTheme,
  loadTheme,
  initTheme
}
