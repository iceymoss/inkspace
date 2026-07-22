export const VIRTUAL_ROOT = '/inkspace'

const TOP_LEVEL_PATHS = new Set(['index', 'blog', 'works', 'photos', 'wiki', 'users', 'about', 'links'])
const ACCOUNT_ROUTES = {
  articles: '/dashboard/articles',
  workspaces: '/dashboard/workspaces',
  'my-works': '/dashboard/works',
  favorites: '/favorites',
  notifications: '/dashboard/notifications',
  profile: '/profile/edit',
  appearance: '/dashboard/appearance'
}
const RESOURCE_EXTENSIONS = {
  blog: '.md',
  works: '.json',
  photos: '.jpg',
  wiki: '.md',
  users: ''
}

function isPositiveId(value) {
  return /^(?:[1-9]\d*)$/.test(String(value))
}

export function normalizeVirtualPath(path = '.', currentPath = VIRTUAL_ROOT) {
  if (typeof path !== 'string' || typeof currentPath !== 'string' || path.includes('\0')) return null

  const base = currentPath.startsWith(`${VIRTUAL_ROOT}/`) || currentPath === VIRTUAL_ROOT
    ? currentPath.split('/')
    : null
  if (!base) return null

  let parts
  if (path.startsWith('/')) {
    if (path !== VIRTUAL_ROOT && !path.startsWith(`${VIRTUAL_ROOT}/`)) return null
    parts = path.split('/')
  } else {
    parts = [...base, ...path.split('/')]
  }

  const normalized = []
  for (const part of parts) {
    if (!part || part === '.') continue
    if (part === '..') {
      if (normalized.length <= 1) return null
      normalized.pop()
      continue
    }
    normalized.push(part)
  }

  if (normalized[0] !== VIRTUAL_ROOT.slice(1)) return null
  return `/${normalized.join('/')}`
}

export function slugifyResourceName(value) {
  const slug = String(value ?? '')
    .normalize('NFKD')
    .replace(/[\u0300-\u036f]/g, '')
    .toLowerCase()
    .trim()
    .replace(/[^\p{Letter}\p{Number}]+/gu, '-')
    .replace(/^-+|-+$/g, '')
  return slug || 'untitled'
}

export function createResourceName(type, id, label) {
  if (!Object.hasOwn(RESOURCE_EXTENSIONS, type) || !isPositiveId(id)) return null
  return `${id}-${slugifyResourceName(label)}${RESOURCE_EXTENSIONS[type]}`
}

export function extractResourceId(value) {
  if (typeof value === 'number') return Number.isSafeInteger(value) && value > 0 ? value : null
  if (typeof value !== 'string') return null

  const name = value.replace(/\/+$/, '').split('/').pop()
  const match = name?.match(/^([1-9]\d*)(?:-|\.|$)/)
  if (!match) return null

  const id = Number(match[1])
  return Number.isSafeInteger(id) ? id : null
}

export function routeToVirtualPath(route) {
  const routePath = typeof route === 'string' ? route : route?.path
  if (typeof routePath !== 'string') return null
  const path = routePath.split(/[?#]/, 1)[0].replace(/\/+$/, '') || '/'

  if (path === '/') return `${VIRTUAL_ROOT}/index`
  if (path === '/login' || path === '/dashboard' || path.startsWith('/dashboard/') || path === '/favorites' || path === '/profile/edit') {
    return VIRTUAL_ROOT
  }
  if (path === '/user-search') return `${VIRTUAL_ROOT}/users`

  const topLevel = path.match(/^\/(blog|works|photos|wiki|about|links)$/)
  if (topLevel) return `${VIRTUAL_ROOT}/${topLevel[1]}`

  const detail = path.match(/^\/(blog|works|users)\/([1-9]\d*)$/)
  if (detail) {
    const extension = RESOURCE_EXTENSIONS[detail[1]]
    return `${VIRTUAL_ROOT}/${detail[1]}/${detail[2]}${extension}`
  }

  const workspace = path.match(/^\/wiki\/([1-9]\d*)$/)
  if (workspace) return `${VIRTUAL_ROOT}/wiki/${workspace[1]}`
  const doc = path.match(/^\/wiki\/docs\/([1-9]\d*)$/)
  if (doc) return `${VIRTUAL_ROOT}/wiki/docs/${doc[1]}.md`

  return null
}

export function virtualPathToRoute(path) {
  const normalized = normalizeVirtualPath(path)
  if (!normalized) return null
  if (normalized === VIRTUAL_ROOT) return '/dashboard'
  if (normalized === `${VIRTUAL_ROOT}/index`) return '/'

  const accountPath = normalized.slice(VIRTUAL_ROOT.length + 1)
  if (Object.hasOwn(ACCOUNT_ROUTES, accountPath)) return ACCOUNT_ROUTES[accountPath]

  const relative = normalized.slice(VIRTUAL_ROOT.length + 1)
  if (!relative || !TOP_LEVEL_PATHS.has(relative.split('/')[0])) return null
  if (relative === 'users') return '/user-search'
  if (['blog', 'works', 'photos', 'wiki', 'about', 'links'].includes(relative)) return `/${relative}`

  const [type, ...resourceParts] = relative.split('/')
  const id = extractResourceId(resourceParts.at(-1))
  if (!id) return null

  if (type === 'blog') return `/blog/${id}`
  if (type === 'works' || type === 'photos') return `/works/${id}`
  if (type === 'users' && resourceParts.length === 1) return `/users/${id}`
  if (type === 'wiki' && resourceParts[0] === 'docs') return `/wiki/docs/${id}`
  if (type === 'wiki' && resourceParts.length > 1) return `/wiki/docs/${id}`
  if (type === 'wiki') return `/wiki/${id}`
  return null
}
