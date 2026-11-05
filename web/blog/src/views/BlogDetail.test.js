import { describe, expect, it, vi, beforeEach } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { createMemoryHistory, createRouter } from 'vue-router'
import { defineComponent } from 'vue'
import { createPinia, setActivePinia } from 'pinia'
import { useTerminalStore } from '@/stores/terminal'
import { createTerminalExecutor } from '@/utils/terminal/executor'
import BlogDetail from './BlogDetail.vue'
import MainLayout from '@/layouts/MainLayout.vue'
import api from '@/utils/api'
import { clearArticleMetadata, setSiteDefaults } from '@/utils/metadata'

vi.mock('@/utils/api', () => ({ default: { get: vi.fn(), post: vi.fn(), delete: vi.fn() } }))
vi.mock('vditor', () => ({ default: { preview: vi.fn() } }))
vi.mock('@/utils/codeTheme', () => ({ loadCodeTheme: vi.fn().mockResolvedValue('github'), loadHighlightTheme: vi.fn().mockResolvedValue() }))
vi.mock('@/stores/user', () => ({ useUserStore: () => ({ isLoggedIn: true, user: { id: 9 } }) }))
vi.mock('@/stores/appearance', () => ({ useAppearanceStore: () => ({ resolvedColorScheme: 'light', activePreference: { ui_theme: 'magazine' } }) }))

const article = (id, slug) => ({ id, slug, title: slug, content: '正文', reading_time: 2, tags: [], view_count: 0,
  seo: { title: `${slug} | Public site`, headline: slug, description: '纯文本摘要', canonical: `https://canonical.example/blog/${slug}`, card: 'summary',
    json_ld: { '@context': 'https://schema.org', '@type': 'BlogPosting', headline: slug, url: `https://canonical.example/blog/${slug}` } }
})

async function renderPage(path = '/blog/first-post', withLayout = false) {
  const children = [
    { path: '/blog/:slug', component: BlogDetail },
    { path: '/', component: { template: '<div>Home</div>' } },
    ...['/blog', '/works', '/photos', '/wiki', '/links', '/about'].map(path => ({ path, component: { template: '<div>Page</div>' } }))
  ]
  const router = createRouter({ history: createMemoryHistory(), routes: [
    ...(withLayout ? [{ path: '/', component: MainLayout, children }] : children)
  ] })
  await router.push(path)
  await router.isReady()
  const wrapper = mount(defineComponent({ template: '<router-view />' }), {
    attachTo: document.body,
    global: { plugins: [router], stubs: {
      'el-card': { template: '<div><slot /></div>' },
      'el-tag': true, 'el-avatar': true, 'el-icon': true,
      'el-button': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
      'el-input': true, 'el-form': true, 'el-link': true, 'el-alert': true,
      'el-dropdown': true, 'el-dropdown-menu': true, 'el-dropdown-item': true, NotificationDropdown: true
    } }
  })
  await flushPromises()
  return { wrapper, router }
}

describe('BlogDetail slug routing', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    clearArticleMetadata()
    setSiteDefaults({ site_name: 'Public site', site_description: '默认描述', public_site_url: 'https://canonical.example' })
    api.get.mockImplementation(async path => {
      if (path === '/articles/slug/first-post') return { data: article(41, 'first-post') }
      if (path === '/articles/slug/second-post') return { data: article(72, 'second-post') }
      if (path === '/settings/public') return { data: { article_comment_enabled: '1', site_name: 'Public site', site_description: '默认描述', public_site_url: 'https://canonical.example' } }
      if (path === '/comments') return { data: { list: [], total: 0 } }
      return { data: {} }
    })
    api.post.mockResolvedValue({ data: {} })
    api.delete.mockResolvedValue({ data: {} })
  })

  it('consumes terminal refreshes once across loading and revisiting the real route', async () => {
    const { wrapper, router } = await renderPage()
    const terminal = useTerminalStore()
    const executor = createTerminalExecutor({ router, route: router.currentRoute.value, api,
      userStore: { isLoggedIn: true }, appearanceStore: {}, terminalStore: terminal })
    const original = api.get.getMockImplementation()
    const pending = []
    api.get.mockImplementation(path => path === '/articles/slug/first-post'
      ? new Promise(resolve => pending.push(resolve)) : original(path))
    const requests = () => api.get.mock.calls.filter(([path]) => path === '/articles/slug/first-post').length
    try {
      for (let signal = 1; signal <= 2; signal++) {
        await executor.execute(`${signal === 1 ? 'favorite' : 'unfavorite'} article 41`)
        await executor.confirm({ ...terminal.pendingConfirmation, accepted: true })
        await flushPromises()
        expect(pending).toHaveLength(1)
        pending.shift()({ data: article(41, 'first-post') })
        await flushPromises()
        expect(requests()).toBe(signal + 1)
        expect(pending).toHaveLength(0)
        expect(document.title).toBe('first-post | Public site')
        expect(document.querySelectorAll('script[type="application/ld+json"]')).toHaveLength(1)
      }
      await router.push('/blog/second-post')
      await flushPromises()
      await router.push('/blog/first-post')
      await flushPromises()
      expect(pending).toHaveLength(1)
      pending.shift()({ data: article(41, 'first-post') })
      await flushPromises()
      expect(requests()).toBe(4)
      expect(pending).toHaveLength(0)
      expect(document.querySelector('link[rel="canonical"]').href).toBe('https://canonical.example/blog/first-post')
    } finally {
      wrapper.unmount()
    }
  })

  it('loads the real slug route and uses returned numeric IDs for interactions', async () => {
    const { wrapper, router } = await renderPage()
    expect(wrapper.find('.article-title-row h1').text()).toBe('first-post')
    expect(document.title).toBe('first-post | Public site')
    expect(document.querySelector('link[rel="canonical"]').href).toBe('https://canonical.example/blog/first-post')
    expect(wrapper.find('.article-meta').text()).toContain('约 2 分钟阅读')
    expect(api.get).toHaveBeenCalledWith('/articles/slug/first-post')
    expect(api.get).toHaveBeenCalledWith('/articles/41/is-liked')
    expect(api.get).toHaveBeenCalledWith('/articles/41/is-favorited')
    expect(api.get).toHaveBeenCalledWith('/comments', { params: { article_id: 41, page: 1, page_size: 10 } })

    await router.push('/blog/second-post')
    await flushPromises()
    expect(wrapper.find('h1').text()).toBe('second-post')
    expect(document.title).toBe('second-post | Public site')
    expect(document.querySelectorAll('script[type="application/ld+json"]')).toHaveLength(1)
    expect(JSON.parse(document.querySelector('script[type="application/ld+json"]').textContent).headline).toBe('second-post')
    expect(document.querySelector('meta[property="og:image"]')).toBeNull()
    expect(api.get).toHaveBeenCalledWith('/articles/slug/second-post')
    expect(api.get).toHaveBeenCalledWith('/articles/72/is-liked')
    const favorite = wrapper.findAll('button').find(button => button.text().includes('收藏'))
    await favorite.trigger('click')
    await flushPromises()
    expect(api.post).toHaveBeenCalledWith('/articles/72/favorite')
    await router.push('/')
    await flushPromises()
    expect(document.title).toBe('Public site')
    expect(document.querySelector('link[rel="canonical"]')).toBeNull()
    expect(document.querySelector('script[type="application/ld+json"]')).toBeNull()
    expect(document.querySelector('meta[name="description"]').content).toBe('默认描述')
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
    expect(document.title).toBe('second-post | Public site')
    expect(api.get).not.toHaveBeenCalledWith('/articles/41/is-liked')
    wrapper.unmount()
  })

  it('loads articles before settings finish, preserves their title, and clears on 404', async () => {
    let resolveSettings
    const original = api.get.getMockImplementation()
    api.get.mockImplementation(path => path === '/settings/public'
      ? new Promise(resolve => { resolveSettings = resolve })
      : path === '/articles/slug/missing' ? Promise.reject(new Error('404')) : original(path))
    const { wrapper, router } = await renderPage()
    expect(document.title).toBe('first-post | Public site')
    resolveSettings({ data: { site_name: 'New site', site_description: '新默认', public_site_url: 'https://canonical.example' } })
    await flushPromises()
    expect(document.title).toBe('first-post | Public site')
    await router.push('/blog/missing')
    await flushPromises()
    expect(document.title).toBe('New site')
    expect(document.querySelector('link[rel="canonical"]')).toBeNull()
    expect(document.querySelector('meta[property="og:title"]')).toBeNull()
    wrapper.unmount()
  })

  it('late MainLayout settings cannot overwrite an article head and become defaults after leaving', async () => {
    const settingsRequests = []
    const original = api.get.getMockImplementation()
    api.get.mockImplementation(path => path === '/settings/public'
      ? new Promise(resolve => { settingsRequests.push(resolve) }) : original(path))
    const { wrapper, router } = await renderPage('/blog/first-post', true)
    expect(settingsRequests).toHaveLength(2)
    expect(document.title).toBe('first-post | Public site')
    settingsRequests.forEach(resolve => resolve({ data: { site_name: 'Late layout', site_description: 'Late defaults', public_site_url: 'https://canonical.example' } }))
    await flushPromises()
    expect(document.title).toBe('first-post | Public site')
    expect(document.querySelector('meta[property="og:title"]').content).toBe('first-post')
    await router.push('/')
    await flushPromises()
    expect(document.title).toBe('Late layout')
    expect(document.querySelector('meta[name="description"]').content).toBe('Late defaults')
    expect(document.querySelectorAll('link[data-inkspace-feed]')).toHaveLength(2)
    wrapper.unmount()
  })

  it('does not restore metadata from a fetch completing after the reader leaves', async () => {
    let resolveArticle
    const original = api.get.getMockImplementation()
    api.get.mockImplementation(path => path === '/articles/slug/first-post'
      ? new Promise(resolve => { resolveArticle = resolve }) : original(path))
    const { wrapper, router } = await renderPage()
    await router.push('/')
    await flushPromises()
    resolveArticle({ data: article(41, 'first-post') })
    await flushPromises()
    expect(document.title).toBe('Public site')
    expect(document.querySelector('link[rel="canonical"]')).toBeNull()
    expect(document.querySelector('script[type="application/ld+json"]')).toBeNull()
    wrapper.unmount()
  })
})
