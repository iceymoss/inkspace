import { describe, expect, it, vi, beforeEach } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { createMemoryHistory, createRouter } from 'vue-router'
import { defineComponent } from 'vue'
import BlogDetail from './BlogDetail.vue'
import api from '@/utils/api'

vi.mock('@/utils/api', () => ({ default: { get: vi.fn(), post: vi.fn(), delete: vi.fn() } }))
vi.mock('vditor', () => ({ default: { preview: vi.fn() } }))
vi.mock('@/utils/codeTheme', () => ({ loadCodeTheme: vi.fn().mockResolvedValue('github'), loadHighlightTheme: vi.fn().mockResolvedValue() }))
vi.mock('@/stores/user', () => ({ useUserStore: () => ({ isLoggedIn: true, user: { id: 9 } }) }))
vi.mock('@/stores/appearance', () => ({ useAppearanceStore: () => ({ resolvedColorScheme: 'light' }) }))
vi.mock('@/stores/terminal', () => ({ useTerminalStore: () => ({ refreshSignals: {} }) }))

const article = (id, slug) => ({ id, slug, title: slug, content: '正文', reading_time: 2, tags: [], view_count: 0 })

async function renderPage(path = '/blog/first-post') {
  const router = createRouter({ history: createMemoryHistory(), routes: [
    { path: '/blog/:slug', component: BlogDetail },
    { path: '/', component: { template: '<div>Home</div>' } }
  ] })
  await router.push(path)
  await router.isReady()
  const wrapper = mount(defineComponent({ template: '<router-view />' }), {
    attachTo: document.body,
    global: { plugins: [router], stubs: {
      'el-card': { template: '<div><slot /></div>' },
      'el-tag': true, 'el-avatar': true, 'el-icon': true,
      'el-button': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
      'el-input': true, 'el-form': true, 'el-link': true, 'el-alert': true
    } }
  })
  await flushPromises()
  return { wrapper, router }
}

describe('BlogDetail slug routing', () => {
  beforeEach(() => {
    api.get.mockImplementation(async path => {
      if (path === '/articles/slug/first-post') return { data: article(41, 'first-post') }
      if (path === '/articles/slug/second-post') return { data: article(72, 'second-post') }
      if (path === '/settings/public') return { data: { article_comment_enabled: '1' } }
      if (path === '/comments') return { data: { list: [], total: 0 } }
      return { data: {} }
    })
    api.post.mockResolvedValue({ data: {} })
  })

  it('loads the real slug route and uses returned numeric IDs for interactions', async () => {
    const { wrapper, router } = await renderPage()
    expect(wrapper.find('.article-title-row h1').text()).toBe('first-post')
    expect(wrapper.find('.article-meta').text()).toContain('约 2 分钟阅读')
    expect(api.get).toHaveBeenCalledWith('/articles/slug/first-post')
    expect(api.get).toHaveBeenCalledWith('/articles/41/is-liked')
    expect(api.get).toHaveBeenCalledWith('/articles/41/is-favorited')
    expect(api.get).toHaveBeenCalledWith('/comments', { params: { article_id: 41, page: 1, page_size: 10 } })

    await router.push('/blog/second-post')
    await flushPromises()
    expect(wrapper.find('h1').text()).toBe('second-post')
    expect(api.get).toHaveBeenCalledWith('/articles/slug/second-post')
    expect(api.get).toHaveBeenCalledWith('/articles/72/is-liked')
    const favorite = wrapper.findAll('button').find(button => button.text().includes('收藏'))
    await favorite.trigger('click')
    await flushPromises()
    expect(api.post).toHaveBeenCalledWith('/articles/72/favorite')
    wrapper.unmount()
  })

  it('does not let a late article response overwrite a newer slug', async () => {
    let resolveFirst
    const original = api.get.getMockImplementation()
    api.get.mockImplementation(path => path === '/articles/slug/first-post'
      ? new Promise(resolve => { resolveFirst = resolve }) : original(path))
    const { wrapper, router } = await renderPage()
    await router.push('/blog/second-post')
    await flushPromises()
    resolveFirst({ data: article(41, 'first-post') })
    await flushPromises()
    expect(wrapper.find('h1').text()).toBe('second-post')
    expect(api.get).not.toHaveBeenCalledWith('/articles/41/is-liked')
    wrapper.unmount()
  })
})
