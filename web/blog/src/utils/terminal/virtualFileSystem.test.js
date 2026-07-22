import { describe, expect, it } from 'vitest'
import {
  createResourceName,
  extractResourceId,
  normalizeVirtualPath,
  routeToVirtualPath,
  slugifyResourceName,
  virtualPathToRoute
} from './virtualFileSystem'

describe('terminal virtual file system', () => {
  it('normalizes paths without allowing escape from /inkspace', () => {
    expect(normalizeVirtualPath('../blog/./12-post.md', '/inkspace/index')).toBe('/inkspace/blog/12-post.md')
    expect(normalizeVirtualPath('/inkspace/works/../blog')).toBe('/inkspace/blog')
    expect(normalizeVirtualPath('../../outside', '/inkspace/blog')).toBeNull()
    expect(normalizeVirtualPath('/etc/passwd')).toBeNull()
  })

  it.each([
    ['/', '/inkspace/index'],
    ['/blog', '/inkspace/blog'],
    ['/works', '/inkspace/works'],
    ['/photos', '/inkspace/photos'],
    ['/wiki', '/inkspace/wiki'],
    ['/user-search', '/inkspace/users'],
    ['/about', '/inkspace/about'],
    ['/links', '/inkspace/links'],
    ['/blog/42', '/inkspace/blog/42.md'],
    ['/works/7', '/inkspace/works/7.json'],
    ['/users/9', '/inkspace/users/9'],
    ['/wiki/3', '/inkspace/wiki/3'],
    ['/wiki/docs/18', '/inkspace/wiki/docs/18.md']
  ])('maps route %s to virtual path %s', (route, path) => {
    expect(routeToVirtualPath(route)).toBe(path)
    expect(virtualPathToRoute(path)).toBe(route)
  })

  it('maps account routes to the inkspace user directory', () => {
    expect(routeToVirtualPath('/dashboard')).toBe('/inkspace')
    expect(routeToVirtualPath('/dashboard/articles')).toBe('/inkspace/articles')
    expect(routeToVirtualPath('/dashboard/workspaces')).toBe('/inkspace/workspaces')
    expect(routeToVirtualPath('/dashboard/works')).toBe('/inkspace/my-works')
    expect(routeToVirtualPath('/favorites')).toBe('/inkspace/favorites')
    expect(routeToVirtualPath('/dashboard/notifications')).toBe('/inkspace/notifications')
    expect(routeToVirtualPath('/profile/edit')).toBe('/inkspace/profile')
    expect(routeToVirtualPath('/dashboard/appearance')).toBe('/inkspace/appearance')
    expect(routeToVirtualPath('/login')).toBe('/inkspace')
    expect(virtualPathToRoute('/inkspace')).toBe('/dashboard')
    expect(virtualPathToRoute('/inkspace/articles')).toBe('/dashboard/articles')
    expect(virtualPathToRoute('/inkspace/workspaces')).toBe('/dashboard/workspaces')
    expect(virtualPathToRoute('/inkspace/my-works')).toBe('/dashboard/works')
    expect(virtualPathToRoute('/inkspace/appearance')).toBe('/dashboard/appearance')
  })

  it('maps readable resource paths back to ID-based routes', () => {
    expect(virtualPathToRoute('/inkspace/blog/42-vue-终端.md')).toBe('/blog/42')
    expect(virtualPathToRoute('/inkspace/photos/8-summer-trip.jpg')).toBe('/works/8')
    expect(virtualPathToRoute('/inkspace/users/9-alice')).toBe('/users/9')
    expect(virtualPathToRoute('/inkspace/wiki/3-frontend/18-routing.md')).toBe('/wiki/docs/18')
  })

  it('creates readable names and extracts positive IDs only', () => {
    expect(slugifyResourceName('  Vue 3: 终端基础  ')).toBe('vue-3-终端基础')
    expect(createResourceName('blog', 12, 'Hello, InkSpace!')).toBe('12-hello-inkspace.md')
    expect(createResourceName('works', 0, 'invalid')).toBeNull()
    expect(extractResourceId('/inkspace/blog/12-hello.md')).toBe(12)
    expect(extractResourceId('0-invalid.md')).toBeNull()
    expect(extractResourceId('-4-invalid.md')).toBeNull()
    expect(extractResourceId('12abc.md')).toBeNull()
  })

  it('rejects unknown routes and virtual roots', () => {
    expect(routeToVirtualPath('/blog/not-an-id')).toBeNull()
    expect(virtualPathToRoute('/inkspace/private/1-file')).toBeNull()
  })
})
