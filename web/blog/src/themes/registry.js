export const themeRegistry = Object.freeze([
  {
    id: 'magazine',
    name: '屿刊',
    subtitle: '极简杂志风',
    description: '纸张、衬线与发丝线，让内容像一本持续更新的个人刊物。',
    status: 'available',
    defaultColorScheme: 'system',
    design: 'style-a-magazine',
    stylesheet: 'magazine.css'
  },
  {
    id: 'terminal',
    name: 'inkspace.log',
    subtitle: '暗色科技感',
    description: '以终端、日志流和冰蓝状态光构成的开发者工作台。',
    status: 'available',
    defaultColorScheme: 'dark',
    design: 'style-b-terminal',
    stylesheet: 'terminal.css'
  },
  {
    id: 'cozy',
    name: '小屿的角落',
    subtitle: '温暖手作感',
    description: '拍立得、手绘线和温暖纸张组成的一间内容小屋。',
    status: 'coming_soon',
    defaultColorScheme: 'light',
    design: 'style-c-cozy',
    stylesheet: null
  },
  {
    id: 'swiss',
    name: 'CHEN YU®',
    subtitle: '瑞士网格风',
    description: '外露网格、精确编号与克莱因蓝构成的视觉系统。',
    status: 'coming_soon',
    defaultColorScheme: 'light',
    design: 'style-d-swiss',
    stylesheet: null
  }
])

export const themeIds = new Set(themeRegistry.map((theme) => theme.id))
export const availableThemeIds = new Set(
  themeRegistry.filter((theme) => theme.status === 'available').map((theme) => theme.id)
)

export function getTheme(themeId) {
  return themeRegistry.find((theme) => theme.id === themeId) || themeRegistry[0]
}
