import { describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import BlogDetail from './BlogDetail.vue'
import api from '@/utils/api'

vi.mock('@/utils/api', () => ({ default: { get: vi.fn() } }))
vi.mock('vditor', () => ({ default: { preview: vi.fn() } }))
vi.mock('@/utils/codeTheme', () => ({ loadCodeTheme: vi.fn().mockResolvedValue('github'), loadHighlightTheme: vi.fn().mockResolvedValue() }))
vi.mock('vue-router', () => ({ useRoute: () => ({ params: { id: '1' } }), useRouter: () => ({ push: vi.fn() }) }))
vi.mock('@/stores/user', () => ({ useUserStore: () => ({ isLoggedIn: false, user: null }) }))
vi.mock('@/stores/appearance', () => ({ useAppearanceStore: () => ({ resolvedColorScheme: 'light' }) }))
vi.mock('@/stores/terminal', () => ({ useTerminalStore: () => ({ refreshSignals: {} }) }))

describe('BlogDetail', () => {
  it('displays the detail reading estimate in the title metadata', async () => {
    api.get.mockImplementation(async path => {
      if (path === '/articles/1') return { data: { id: 1, title: '测试文章', content: '正文', reading_time: 2, tags: [], view_count: 0 } }
      if (path === '/settings/public') return { data: { article_comment_enabled: '0' } }
      return { data: {} }
    })

    const wrapper = mount(BlogDetail, {
      attachTo: document.body,
      global: { stubs: {
        'el-card': { template: '<div><slot /></div>' },
        'el-tag': true, 'el-avatar': true, 'el-icon': true, 'el-button': true,
        'el-input': true, 'el-form': true, 'el-link': true, 'el-alert': true
      } }
    })
    await flushPromises()

    expect(wrapper.find('.article-title-row h1').text()).toBe('测试文章')
    expect(wrapper.find('.article-meta').text()).toContain('约 2 分钟阅读')
    expect(api.get).toHaveBeenCalledWith('/articles/1')
    wrapper.unmount()
  })
})
