import { beforeEach, describe, expect, it, vi } from 'vitest'
import { reactive } from 'vue'
import { createTerminalExecutor } from './executor'

function setup(path = '/') {
  const route = reactive({ path, query: {}, params: {} })
  const router = {
    push: vi.fn(async target => {
      route.path = typeof target === 'string' ? target : target.path
      route.query = typeof target === 'string' ? {} : target.query || {}
    })
  }
  const api = {
    get: vi.fn(),
    post: vi.fn().mockResolvedValue({ data: {} }),
    delete: vi.fn().mockResolvedValue({ data: {} })
  }
  const terminalStore = {
    isPending: false,
    output: [],
    appendOutput(text, type = 'output') {
      const lines = Array.isArray(text) ? text : [text]
      this.output.push(...lines.map(line => ({ text: String(line), type })))
    },
    setPending(value) { this.isPending = value },
    setPendingConfirmation: vi.fn(),
    signalRefresh: vi.fn(),
    minimize: vi.fn(),
    close: vi.fn()
  }
  const userStore = { isLoggedIn: false }
  const appearanceStore = {
    activePreference: { ui_theme: 'terminal', color_scheme: 'system' },
    resolvedColorScheme: 'dark',
    saveGuestPreference: vi.fn(),
    save: vi.fn()
  }
  return {
    route, router, api, terminalStore, userStore, appearanceStore,
    executor: createTerminalExecutor({ route, router, api, terminalStore, userStore, appearanceStore })
  }
}

describe('terminal executor', () => {
  beforeEach(() => vi.clearAllMocks())

  it('lists the public content shown on the homepage', async () => {
    const context = setup('/')
    context.api.get
      .mockResolvedValueOnce({ data: [{ id: 12, title: 'Vue Notes' }] })
      .mockResolvedValueOnce({ data: [{ id: 7, title: 'InkSpace CLI' }] })
      .mockResolvedValueOnce({ data: { list: [{ id: 8, title: 'Night Sky' }] } })
      .mockResolvedValueOnce({ data: { list: [{ id: 3, name: 'Frontend' }] } })
    await context.executor.execute('ls')

    expect(context.api.get).toHaveBeenCalledTimes(4)
    expect(context.terminalStore.output.map(item => item.text)).toEqual(expect.arrayContaining([
      'blog/12-vue-notes.md',
      'works/7-inkspace-cli.json',
      'photos/8-night-sky.jpg',
      'wiki/3-frontend/'
    ]))
  })

  it('changes virtual directories through the router', async () => {
    const context = setup('/')
    await context.executor.execute('cd blog')

    expect(context.router.push).toHaveBeenCalledWith('/blog')
    expect(context.route.path).toBe('/blog')
  })

  it('enters public directories exactly as listed from index', async () => {
    const context = setup('/')

    await context.executor.execute('cd works')

    expect(context.router.push).toHaveBeenCalledWith('/works')
    expect(context.terminalStore.output.at(-1).text).toBe('cwd /inkspace/works')
  })

  it('treats /inkspace as the user directory and opens login for guests', async () => {
    const context = setup('/')
    await context.executor.execute('cd ../')

    expect(context.router.push).toHaveBeenCalledWith('/login')
    expect(context.terminalStore.output.at(-1).text).toContain('opened login')
  })

  it('opens the dashboard when an authenticated user enters /inkspace', async () => {
    const context = setup('/')
    context.userStore.isLoggedIn = true
    await context.executor.execute('cd ../')

    expect(context.router.push).toHaveBeenCalledWith('/dashboard')
  })

  it('lists account directories from the authenticated user directory', async () => {
    const context = setup('/dashboard')
    context.userStore.isLoggedIn = true
    await context.executor.execute('ls')

    expect(context.terminalStore.output.at(-1).text).toContain('articles/')
    expect(context.terminalStore.output.at(-1).text).toContain('appearance/')
  })

  it('keeps the selected account directory after navigation', async () => {
    const context = setup('/dashboard')
    context.userStore.isLoggedIn = true

    await context.executor.execute('cd my-works')

    expect(context.router.push).toHaveBeenCalledWith('/dashboard/works')
    expect(context.terminalStore.output.at(-1).text).toBe('cwd /inkspace/my-works')
    await context.executor.execute('ls')
    expect(context.terminalStore.output.at(-1).text).toBe('my-works/')
    await context.executor.execute('pwd')
    expect(context.terminalStore.output.at(-1).text).toBe('/inkspace/my-works')
  })

  it('lists real articles through the existing API', async () => {
    const context = setup('/blog')
    context.api.get.mockResolvedValue({ data: { list: [{ id: 12, title: 'Vue Notes' }] } })
    await context.executor.execute('ls')

    expect(context.api.get).toHaveBeenCalledWith('/articles', {
      skipAuth: true,
      silentError: true,
      params: { page: 1, page_size: 10 }
    })
    expect(context.terminalStore.output.at(-1).text).toContain('12-vue-notes.md')
  })

  it('navigates searches with URL query and calls the query API', async () => {
    const context = setup('/')
    context.api.get.mockResolvedValue({ data: { list: [] } })
    await context.executor.execute('grep "vue router" ../blog')

    expect(context.router.push).toHaveBeenCalledWith({ path: '/blog', query: { keyword: 'vue router' } })
    expect(context.api.get).toHaveBeenCalledWith('/articles', {
      skipAuth: true,
      silentError: true,
      params: { keyword: 'vue router', page: 1, page_size: 10 }
    })
  })

  it('requires login before preparing a write action', async () => {
    const context = setup('/blog/12')
    await context.executor.execute('favorite article 12')

    expect(context.api.post).not.toHaveBeenCalled()
    expect(context.terminalStore.output.at(-1).text).toContain('login required')
  })

  it('confirms a favorite before using the existing endpoint', async () => {
    const context = setup('/blog/12')
    context.userStore.isLoggedIn = true
    await context.executor.execute('favorite article 12')
    const confirmation = context.terminalStore.setPendingConfirmation.mock.calls[0][0]

    expect(context.api.post).not.toHaveBeenCalled()
    await context.executor.confirm({ ...confirmation, accepted: true })

    expect(context.api.post).toHaveBeenCalledWith('/articles/12/favorite', undefined, { silentError: true })
    expect(context.terminalStore.signalRefresh).toHaveBeenCalledWith('article:12')
  })
})
