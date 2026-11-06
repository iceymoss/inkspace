import { beforeEach, describe, expect, it } from 'vitest'
import { clearArticleMetadata, setArticleMetadata, setSiteDefaults } from './metadata'

describe('metadata DOM ownership', () => {
  beforeEach(() => {
    document.head.innerHTML = '<title data-inkspace-meta>Server</title><meta data-inkspace-meta property="og:title" content="Server"><script data-inkspace-meta type="application/ld+json">{}</script><script src="/assets/app.js"></script>'
    clearArticleMetadata()
    setSiteDefaults({ site_name: '站点', site_description: '默认', public_site_url: 'https://canonical.example' })
  })

  it('adopts server nodes and inserts hostile text without parsing it as HTML', () => {
    const headline = '中文 </script><script>alert("x")</script> & \\'
    const metadata = {
      title: `${headline} | 站点`, headline, description: '纯文本 & "摘要"',
      canonical: 'https://canonical.example/blog/post', image: 'https://cdn.example/cover.png', card: 'summary_large_image',
      json_ld: { '@type': 'BlogPosting', headline }
    }
    setArticleMetadata(metadata)
    setArticleMetadata(metadata)
    expect(document.title).toBe(`${headline} | 站点`)
    expect(document.querySelectorAll('script[type="application/ld+json"]')).toHaveLength(1)
    expect(JSON.parse(document.querySelector('script[type="application/ld+json"]').textContent).headline).toBe(headline)
    expect(document.querySelectorAll('script:not([type="application/ld+json"])')).toHaveLength(1)
    expect(document.querySelectorAll('meta[property="og:title"]')).toHaveLength(1)
    expect(document.querySelector('meta[name="twitter:image"]').content).toBe('https://cdn.example/cover.png')
    clearArticleMetadata()
    expect(document.title).toBe('站点')
    expect(document.querySelector('meta[property="og:image"]')).toBeNull()
    expect(document.querySelector('script[src="/assets/app.js"]')).not.toBeNull()
    expect(document.querySelectorAll('link[data-inkspace-feed]')).toHaveLength(2)
  })
})
