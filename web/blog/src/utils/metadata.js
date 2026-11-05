// Article metadata is a backend projection: canonical origin, excerpt, resources
// and JSON-LD are identical to initial HTML, not re-derived from browser Host.
let site = { name: 'InkSpace', description: '', url: '' }
let article = null

const initialSite = document.getElementById('inkspace-site')
if (initialSite) {
  try { site = { ...site, ...JSON.parse(initialSite.textContent) } } catch { /* Invalid bootstrap data is not trusted. */ }
}

const append = (tag, attributes, text) => {
  const node = document.createElement(tag)
  node.setAttribute('data-inkspace-meta', '')
  for (const [key, value] of Object.entries(attributes)) node.setAttribute(key, value)
  if (text !== undefined) node.textContent = text
  document.head.appendChild(node)
}

const meta = (key, content, property = false) => append('meta', { [property ? 'property' : 'name']: key, content })

const render = () => {
  document.head.querySelectorAll('[data-inkspace-meta]').forEach(node => node.remove())
  // Adopt the single static title used by Vite before Go head injection.
  document.head.querySelectorAll('title').forEach(node => node.remove())
  append('title', {}, article?.title || site.name)
  meta('description', article?.description || site.description)
  if (!article) return
  append('link', { rel: 'canonical', href: article.canonical })
  meta('og:type', 'article', true)
  meta('og:title', article.headline, true)
  meta('og:description', article.description, true)
  meta('og:url', article.canonical, true)
  meta('twitter:card', article.card)
  meta('twitter:title', article.headline)
  meta('twitter:description', article.description)
  if (article.image) {
    meta('og:image', article.image, true)
    meta('twitter:image', article.image)
  }
  // DOM textContent (never innerHTML) keeps script-like user text inert.
  append('script', { type: 'application/ld+json' }, JSON.stringify(article.json_ld))
}

/** Update defaults without allowing late settings to overwrite an article head. */
export const setSiteDefaults = settings => {
  site = {
    name: settings.site_name?.trim() || 'InkSpace',
    description: settings.site_description || '',
    url: settings.public_site_url || site.url
  }
  if (site.url) {
    document.head.querySelectorAll('[data-inkspace-feed]').forEach(node => node.remove())
    for (const [path, type, title] of [['/feed.xml', 'application/rss+xml', 'RSS'], ['/atom.xml', 'application/atom+xml', 'Atom']]) {
      const node = document.createElement('link')
      node.setAttribute('data-inkspace-feed', '')
      node.rel = 'alternate'
      node.type = type
      node.title = title
      node.href = site.url + path
      document.head.appendChild(node)
    }
  }
  // Preserve Go's initial article head while the first slug API call is pending.
  if (article || !document.querySelector('link[data-inkspace-meta][rel="canonical"]')) render()
}

/** Apply only the current slug API response, after the route's stale guard. */
export const setArticleMetadata = metadata => {
  article = metadata || null
  render()
}

/** Restore site defaults when leaving an article, loading another, or failing. */
export const clearArticleMetadata = () => {
  article = null
  render()
}
