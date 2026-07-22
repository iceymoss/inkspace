import { parseCommand } from './parser'
import { availableThemeIds } from '@/themes/registry'
import {
  VIRTUAL_ROOT,
  createResourceName,
  extractResourceId,
  normalizeVirtualPath,
  routeToVirtualPath,
  virtualPathToRoute
} from './virtualFileSystem'

const TOP_LEVEL = ['index/', 'blog/', 'works/', 'photos/', 'wiki/', 'users/', 'about/', 'links/']
const ACCOUNT_LEVEL = ['index/', 'articles/', 'workspaces/', 'my-works/', 'favorites/', 'notifications/', 'profile/', 'appearance/']
const NAVIGATION = {
  home: '/', index: '/', blog: '/blog', works: '/works', photos: '/photos', wiki: '/wiki',
  users: '/user-search', about: '/about', links: '/links', login: '/login', dashboard: '/dashboard'
}
const WRITE_COMMANDS = new Set(['like', 'unlike', 'favorite', 'unfavorite', 'follow', 'unfollow'])
const HELP = [
  'filesystem: ls [path], cd <path>, pwd, cat <path>, grep <word> <path>',
  'navigation: open home|blog|works|photos|wiki|about|links|article <id>|work <id>|user <id>',
  'query: list articles|works|photos|users|wiki [keyword], search blog|works|photos|users <keyword>',
  `appearance: theme ${[...availableThemeIds].join('|')}, scheme system|light|dark, status`,
  'actions: like|unlike article|work <id>, favorite|unfavorite article|work <id>, follow|unfollow user <id>',
  'window: clear, history, minimize, close; write actions require yes/no'
]
const PUBLIC_REQUEST = Object.freeze({ skipAuth: true, silentError: true })
const TERMINAL_REQUEST = Object.freeze({ silentError: true })

function positiveID(value) {
  const id = extractResourceId(value)
  if (!id) throw new Error('resource id must be a positive integer or virtual resource path')
  return id
}

function routePath(route) {
  return routeToVirtualPath(route) || `${VIRTUAL_ROOT}/index`
}

function resolvePath(rawPath, cwd) {
  const relativePath = String(rawPath).replace(/^(?:\.\/)+/, '')
  const firstPart = relativePath.split('/')[0]
  if (cwd === `${VIRTUAL_ROOT}/index` && TOP_LEVEL.includes(`${firstPart}/`)) {
    return normalizeVirtualPath(relativePath, VIRTUAL_ROOT)
  }
  return normalizeVirtualPath(rawPath, cwd)
}

function textContent(html) {
  if (!html) return ''
  const documentNode = new DOMParser().parseFromString(String(html), 'text/html')
  return (documentNode.body.textContent || '').replace(/\s+/g, ' ').trim().slice(0, 1200)
}

function resourceLines(type, list) {
  return list.slice(0, 20).map(item => {
    const label = item.title || item.name || item.nickname || item.username || 'untitled'
    const name = createResourceName(type, item.id, label)
    return `${name || item.id}  ${label}`
  })
}

function prefixedResourceLines(type, list) {
  return list.slice(0, 20).map(item => `${type}/${createResourceName(type, item.id, item.title || item.name)}`)
}

function queryWithoutEmpty(query) {
  return Object.fromEntries(Object.entries(query).filter(([, value]) => value !== '' && value != null))
}

export function createTerminalExecutor({ router, route, api, userStore, appearanceStore, terminalStore }) {
  const pendingActions = new Map()

  const write = (text, type = 'output') => terminalStore.appendOutput(text, type)

  async function navigate(path, query) {
    await router.push(query ? { path, query: queryWithoutEmpty(query) } : path)
    write(`cwd ${routePath(route)}`, 'success')
  }

  async function listPath(rawPath = '.') {
    const cwd = routePath(route)
    const path = resolvePath(rawPath, cwd)
    if (!path) throw new Error('path escapes /inkspace or is invalid')
    if (path === VIRTUAL_ROOT) {
      write(userStore.isLoggedIn ? ACCOUNT_LEVEL.join('  ') : 'login required; use open login or cd index')
      return
    }
    if (path === `${VIRTUAL_ROOT}/index`) {
      const [articles, works, photos, workspaces] = await Promise.all([
        api.get('/articles/hot?limit=6', PUBLIC_REQUEST),
        api.get('/works/hot?limit=4', PUBLIC_REQUEST),
        api.get('/works', { ...PUBLIC_REQUEST, params: { type: 'photography', status: 1, page: 1, page_size: 3 } }),
        api.get('/wiki/workspaces', { ...PUBLIC_REQUEST, params: { page: 1, page_size: 4 } })
      ])
      write([
        ...prefixedResourceLines('blog', articles.data || []),
        ...prefixedResourceLines('works', works.data || []),
        ...prefixedResourceLines('photos', photos.data?.list || []),
        ...(workspaces.data?.list || []).map(item => `wiki/${item.id}-${createResourceName('wiki', item.id, item.name).replace(/^\d+-|\.md$/g, '')}/`)
      ])
      return
    }

    const relative = path.slice(VIRTUAL_ROOT.length + 1)
    const parts = relative.split('/')
    const [type] = parts
    if (type === 'blog') {
      if (parts.length !== 1) throw new Error(`not a directory: ${path}`)
      const response = await api.get('/articles', { ...PUBLIC_REQUEST, params: { page: 1, page_size: 10 } })
      write(resourceLines('blog', response.data?.list || []))
    } else if (type === 'works' || type === 'photos') {
      if (parts.length !== 1) throw new Error(`not a directory: ${path}`)
      const response = await api.get('/works', { ...PUBLIC_REQUEST, params: { page: 1, page_size: 10, ...(type === 'photos' ? { type: 'photography' } : {}) } })
      write(resourceLines(type, response.data?.list || []))
    } else if (type === 'wiki') {
      if (parts.length > 2) throw new Error(`not a directory: ${path}`)
      const workspaceID = parts[1] && positiveID(parts[1])
      if (workspaceID) {
        const response = await api.get(`/wiki/workspaces/${workspaceID}/tree`, PUBLIC_REQUEST)
        const catalogs = (response.data?.catalogs || []).map(item => `${item.id}-${item.name}/`)
        const docs = (response.data?.docs || []).map(item => createResourceName('wiki', item.id, item.title))
        write([...catalogs, ...docs])
      } else {
        const response = await api.get('/wiki/workspaces', { ...PUBLIC_REQUEST, params: { page: 1, page_size: 10 } })
        write((response.data?.list || []).map(item => `${item.id}-${item.name}/  ${item.doc_count || 0} docs`))
      }
    } else if (type === 'users') {
      if (parts.length !== 1) throw new Error(`not a directory: ${path}`)
      write('usage: list users <keyword> or grep <keyword> users/')
    } else if (TOP_LEVEL.includes(`${type}/`) || ACCOUNT_LEVEL.includes(`${type}/`)) {
      write(`${type}/`)
    } else throw new Error(`directory not found: ${path}`)
  }

  async function catPath(rawPath) {
    if (!rawPath) throw new Error('usage: cat <virtual-resource-path>')
    const path = resolvePath(rawPath, routePath(route))
    if (!path) throw new Error('invalid virtual path')
    const relative = path.slice(VIRTUAL_ROOT.length + 1)
    const [type] = relative.split('/')
    const id = positiveID(relative)
    if (type === 'blog') {
      const response = await api.get(`/articles/${id}`, PUBLIC_REQUEST)
      const item = response.data || {}
      if (Number(item.status) !== 1) throw new Error('resource is not publicly published')
      write([`# ${item.title || `article ${id}`}`, `${item.author?.nickname || item.author?.username || 'anonymous'} · ${item.created_at || ''}`, item.summary || textContent(item.content_html)])
      await navigate(`/blog/${id}`)
    } else if (type === 'works' || type === 'photos') {
      const response = await api.get(`/works/${id}`, PUBLIC_REQUEST)
      const item = response.data || {}
      if (Number(item.status) !== 1) throw new Error('resource is not publicly published')
      write([`# ${item.title || `work ${id}`}`, `${item.type || 'work'} · ${item.author?.nickname || item.author?.username || 'anonymous'}`, item.description || 'no description'])
      await navigate(`/works/${id}`)
    } else if (type === 'wiki') {
      const response = await api.get(`/wiki/docs/${id}`, PUBLIC_REQUEST)
      const item = response.data || {}
      write([`# ${item.title || `doc ${id}`}`, textContent(item.content_html) || 'empty document'])
      await navigate(`/wiki/docs/${id}`)
    } else throw new Error('cat supports blog, works, photos and wiki document paths')
  }

  async function search(type, keyword, navigatePage = true, page = 1) {
    if (!keyword && navigatePage) throw new Error(`usage: search ${type} <keyword>`)
    if (type === 'blog' || type === 'articles') {
      if (navigatePage) await navigate('/blog', { keyword })
      const response = await api.get('/articles', { ...PUBLIC_REQUEST, params: { ...(keyword ? { keyword } : {}), page, page_size: 10 } })
      write(resourceLines('blog', response.data?.list || []))
    } else if (type === 'works' || type === 'photos') {
      if (navigatePage) await navigate(type === 'photos' ? '/photos' : '/works', { keyword })
      const response = await api.get('/works', { ...PUBLIC_REQUEST, params: { ...(keyword ? { keyword } : {}), page, page_size: 10, ...(type === 'photos' ? { type: 'photography' } : {}) } })
      write(resourceLines(type, response.data?.list || []))
    } else if (type === 'users') {
      if (navigatePage) await navigate('/user-search', { keyword })
      const response = await api.get('/users/search', { ...PUBLIC_REQUEST, params: { keyword, limit: 10 } })
      write(resourceLines('users', [...(response.data?.top || []), ...(response.data?.users || [])]))
    } else if (type === 'wiki') {
      const response = await api.get('/wiki/workspaces', { ...PUBLIC_REQUEST, params: { page, page_size: 10 } })
      const matches = (response.data?.list || []).filter(item => `${item.name} ${item.description || ''}`.toLowerCase().includes(keyword.toLowerCase()))
      write(matches.map(item => `${item.id}-${item.name}/  ${item.doc_count || 0} docs`))
      if (navigatePage) await navigate('/wiki')
    } else throw new Error('search supports blog, works, photos, users and wiki workspace names')
  }

  async function open(args) {
    const [target, value] = args
    if (!target) throw new Error('usage: open <destination> [id]')
    if (NAVIGATION[target]) {
      if (target === 'dashboard' && !userStore.isLoggedIn) throw new Error('login required; use open login')
      await navigate(NAVIGATION[target])
      return
    }
    const id = positiveID(value)
    const path = target === 'article' ? `/blog/${id}`
      : target === 'work' ? `/works/${id}`
        : target === 'user' ? `/users/${id}`
          : target === 'workspace' ? `/wiki/${id}`
            : target === 'doc' ? `/wiki/docs/${id}` : null
    if (!path) throw new Error(`unknown destination: ${target}`)
    await navigate(path)
  }

  async function cd(rawPath) {
    if (!rawPath) throw new Error('usage: cd <virtual-path>')
    const path = resolvePath(rawPath, routePath(route))
    if (/\.(?:md|json|jpg)$/i.test(path || '')) throw new Error(`not a directory: ${rawPath}`)
    const destination = path === VIRTUAL_ROOT
      ? (userStore.isLoggedIn ? '/dashboard' : '/login')
      : path && virtualPathToRoute(path)
    if (!destination) throw new Error(`not a navigable directory: ${rawPath}`)
    await navigate(destination)
    if (path === VIRTUAL_ROOT && !userStore.isLoggedIn) {
      write('authentication required for /inkspace; opened login', 'system')
    }
  }

  async function setAppearance(kind, value) {
    const preference = { ...appearanceStore.activePreference }
    if (kind === 'theme') {
      if (!availableThemeIds.has(value)) throw new Error(`theme must be one of: ${[...availableThemeIds].join(', ')}`)
      preference.ui_theme = value
    } else {
      if (!['system', 'light', 'dark'].includes(value)) throw new Error('scheme must be system, light or dark')
      preference.color_scheme = value
    }
    if (userStore.isLoggedIn) await appearanceStore.save(preference)
    else appearanceStore.preview(preference)
    write(`${kind} set to ${value}`, 'success')
  }

  async function prepareWrite(name, args) {
    if (!userStore.isLoggedIn) throw new Error('login required; use open login')
    let [type, target] = args
    if (args.length === 1 && String(type).includes('/')) {
      const path = normalizeVirtualPath(type, routePath(route))
      if (!path) throw new Error('invalid virtual resource path')
      const relative = path.slice(VIRTUAL_ROOT.length + 1)
      const section = relative.split('/')[0]
      type = section === 'blog' ? 'article' : ['works', 'photos'].includes(section) ? 'work' : section === 'users' ? 'user' : ''
      target = path
    }
    if ((name === 'follow' || name === 'unfollow') && type !== 'user') {
      target = type
      type = 'user'
    }
    if (!['article', 'work', 'user'].includes(type)) throw new Error(`usage: ${name} article|work|user <id>`)
    const id = positiveID(target)
    if ((name === 'follow' || name === 'unfollow') && type !== 'user') throw new Error(`${name} only supports users`)
    if ((name === 'favorite' || name === 'unfavorite') && type === 'user') throw new Error(`${name} does not support users`)
    if ((name === 'like' || name === 'unlike') && type === 'user') throw new Error(`${name} does not support users`)

    if (name === 'like' || name === 'unlike') {
      const statusURL = type === 'article' ? `/articles/${id}/is-liked` : `/works/${id}/liked`
      const response = await api.get(statusURL, TERMINAL_REQUEST)
      const liked = Boolean(response.data?.liked || response.data?.is_liked)
      const desired = name === 'like'
      if (liked === desired) {
        write(`${type} ${id} is already ${desired ? 'liked' : 'not liked'}`, 'system')
        return
      }
    }

    const confirmationID = `${Date.now()}-${name}-${type}-${id}`
    pendingActions.set(confirmationID, { name, type, id })
    terminalStore.setPendingConfirmation({ id: confirmationID, message: `${name} ${type} ${id}?` })
  }

  async function runWrite({ name, type, id }) {
    let method = 'post'
    let url
    if (name === 'like' || name === 'unlike') {
      const statusURL = type === 'article' ? `/articles/${id}/is-liked` : `/works/${id}/liked`
      const status = await api.get(statusURL, TERMINAL_REQUEST)
      const liked = Boolean(status.data?.liked || status.data?.is_liked)
      const desired = name === 'like'
      if (liked === desired) {
        write(`${type} ${id} already has the requested like state`, 'system')
        return
      }
      url = `/${type === 'article' ? 'articles' : 'works'}/${id}/like`
    }
    else if (name === 'favorite' || name === 'unfavorite') {
      method = name === 'favorite' ? 'post' : 'delete'
      url = `/${type === 'article' ? 'articles' : 'works'}/${id}/favorite`
    } else {
      method = name === 'follow' ? 'post' : 'delete'
      url = `/users/${id}/follow`
    }
    if (method === 'delete') await api.delete(url, TERMINAL_REQUEST)
    else await api.post(url, undefined, TERMINAL_REQUEST)
    terminalStore.signalRefresh(`${type}:${id}`)
    write(`${name} ${type} ${id}: done`, 'success')
  }

  async function execute(input) {
    if (terminalStore.isPending) {
      write('another command is still running', 'system')
      return
    }
    const { name, args } = parseCommand(input)
    terminalStore.setPending(true)
    try {
      if (name === 'help') write(HELP, 'system')
      else if (name === 'pwd') write(routePath(route))
      else if (name === 'ls') await listPath(args[0] || '.')
      else if (name === 'cd') await cd(args[0])
      else if (name === 'cat') await catPath(args[0])
      else if (name === 'grep') {
        const path = resolvePath(args[1] || '.', routePath(route))
        if (!path) throw new Error('invalid grep path')
        const type = path.slice(VIRTUAL_ROOT.length + 1).split('/')[0]
        await search(type, args[0])
      }
      else if (name === 'open') await open(args)
      else if (name === 'search') await search(args[0], args.slice(1).join(' '))
      else if (name === 'list') {
        const type = args[0]
        if (!type) throw new Error('usage: list articles|works|photos|users|wiki [keyword|page]')
        if (type === 'wiki') await search(type, '', false, args[1] ? positiveID(args[1]) : 1)
        else await search(type, args.slice(1).join(' '), false)
      }
      else if (name === 'filter') {
        const [section, key, value] = args
        const allowed = section === 'blog'
          ? { category: 'category_id', tag: 'tag_id', rank: 'rank_type' }
          : section === 'works' ? { type: 'type' } : null
        if (!allowed?.[key] || !value) throw new Error('usage: filter blog category|tag|rank <value> or filter works type <value>')
        if (['category', 'tag'].includes(key) && !/^[1-9]\d*$/.test(value)) throw new Error(`${key} must be a positive integer`)
        if (key === 'rank' && !['hot', 'week', 'month', 'year'].includes(value)) throw new Error('invalid blog rank')
        if (key === 'type' && !['project', 'photography'].includes(value)) throw new Error('invalid work type')
        const path = section === 'blog' ? '/blog' : '/works'
        await navigate(path, { [allowed[key]]: value, page: 1 })
      } else if (name === 'sort') {
        const [section, value] = args
        if (!['blog', 'works', 'photos'].includes(section) || !value) throw new Error('usage: sort blog|works|photos <mode>')
        const key = section === 'blog' ? 'sort_by' : 'sort'
        const aliases = { latest: 'time', popular: 'hot', likes: 'like' }
        const normalized = aliases[value] || value
        const valid = section === 'blog'
          ? ['time', 'hot', 'view_count', 'like_count', 'comment_count'].includes(normalized)
          : ['time', 'hot', 'view', 'like'].includes(normalized)
        if (!valid) throw new Error('invalid sort mode')
        await navigate(section === 'blog' ? '/blog' : `/${section}`, { [key]: normalized, page: 1 })
      } else if (name === 'page') {
        if (!/^[1-9]\d*$/.test(args[0] || '')) throw new Error('page must be a positive integer')
        if (!['/blog', '/works', '/photos', '/wiki'].includes(route.path)) throw new Error('page is only available on list pages')
        await navigate(route.path, { ...route.query, page: Number(args[0]) })
      } else if (name === 'reset' && args[0] === 'filters') {
        if (!['/blog', '/works', '/photos', '/wiki', '/user-search'].includes(route.path)) throw new Error('this page has no URL filters')
        await navigate(route.path)
      }
      else if (name === 'theme' || name === 'scheme') await setAppearance(name, args[0])
      else if (name === 'status') write([
        `path: ${routePath(route)}`,
        `auth: ${userStore.isLoggedIn ? 'authenticated' : 'guest'}`,
        `theme: ${appearanceStore.activePreference.ui_theme}`,
        `scheme: ${appearanceStore.activePreference.color_scheme} (${appearanceStore.resolvedColorScheme})`
      ])
      else if (name === 'scroll') {
        if (!['top', 'bottom'].includes(args[0])) throw new Error('usage: scroll top|bottom')
        window.scrollTo({ top: args[0] === 'top' ? 0 : document.documentElement.scrollHeight, behavior: 'smooth' })
      } else if (name === 'minimize') terminalStore.minimize()
      else if (name === 'close') terminalStore.close()
      else if (name === 'focus' && args[0] === 'search') {
        const searchInput = document.querySelector('.search-input input, .works-search-input input, .search-bar input')
        if (!searchInput) throw new Error('this page has no registered search input')
        searchInput.focus()
        write('search focused', 'success')
      }
      else if (WRITE_COMMANDS.has(name)) await prepareWrite(name, args)
      else throw new Error(`command not found: ${name || '(empty)'}. Type help.`)
    } catch (error) {
      write(error.message || 'command failed', 'error')
    } finally {
      terminalStore.setPending(false)
    }
  }

  async function confirm({ id, accepted }) {
    if (!accepted) {
      pendingActions.delete(id)
      write('operation cancelled', 'system')
      return
    }
    const action = pendingActions.get(id)
    if (!action) {
      write('confirmation expired; run the command again', 'error')
      return
    }
    if (terminalStore.isPending) {
      write('another command is still running', 'system')
      return
    }
    pendingActions.delete(id)
    terminalStore.setPending(true)
    try {
      await runWrite(action)
    } catch (error) {
      write(error.message || 'operation failed', 'error')
    } finally {
      terminalStore.setPending(false)
    }
  }

  return { execute, confirm, help: HELP }
}
