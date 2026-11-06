// @vitest-environment node
import { afterAll, beforeAll, describe, expect, it } from 'vitest'
import { createServer as createHTTPServer } from 'node:http'
import { createServer as createViteServer } from 'vite'
import config from '../../vite.config.js'

describe('development public discovery proxies', () => {
  let backend
  let vite
  let origin
  const paths = ['/feed.xml', '/atom.xml', '/sitemap.xml', '/robots.txt', '/sitemaps/2.xml', '/api/settings/public']

  beforeAll(async () => {
    backend = createHTTPServer((req, res) => {
      res.setHeader('Content-Type', 'text/plain')
      res.end(`backend:${req.url}`)
    })
    await new Promise(resolve => backend.listen(0, '127.0.0.1', resolve))
    const target = `http://127.0.0.1:${backend.address().port}`
    const proxy = Object.fromEntries(Object.entries(config.server.proxy).map(([key, value]) => [key, { ...value, target }]))
    vite = await createViteServer({
      configFile: false,
      plugins: [],
      optimizeDeps: { noDiscovery: true, include: [] },
      server: { ...config.server, port: 0, host: '127.0.0.1', proxy }
    })
    await vite.listen()
    origin = `http://127.0.0.1:${vite.httpServer.address().port}`
  })

  afterAll(async () => {
    await vite?.close()
    if (backend) await new Promise(resolve => backend.close(resolve))
  })

  it.each(paths)('forwards %s (including query) instead of returning SPA HTML', async path => {
    const response = await fetch(`${origin}${path}?acceptance=1`)
    expect(response.status).toBe(200)
    expect(await response.text()).toBe(`backend:${path}?acceptance=1`)
  })

  it('uses the blog backend for runtime canonical configuration and discovery', () => {
    for (const key of ['/api', '/feed.xml', '/atom.xml', '/sitemap.xml', '/robots.txt', '/sitemaps/']) {
      expect(config.server.proxy[key]?.target).toBe('http://localhost:8081')
    }
  })
})
