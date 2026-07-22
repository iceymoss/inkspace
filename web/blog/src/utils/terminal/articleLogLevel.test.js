import { describe, expect, it } from 'vitest'
import { getArticleLogLevel } from './articleLogLevel'

describe('article log level', () => {
  it('always marks top articles as FEAT', () => {
    expect(getArticleLogLevel({ id: 1, is_top: true, tags: [{ name: 'warning' }] })).toBe('FEAT')
  })

  it.each([
    ['观点', 'NOTE'],
    ['OPINION', 'NOTE'],
    ['故障复盘', 'WARN'],
    ['Warning', 'WARN']
  ])('uses semantic tag %s as a %s override', (name, level) => {
    expect(getArticleLogLevel({ id: 4, tags: [{ id: 2, name }] })).toBe(level)
  })

  it('uses WARN when both semantic override types are present', () => {
    expect(getArticleLogLevel({ tags: [{ name: '随笔' }, { name: '踩坑' }] })).toBe('WARN')
  })

  it('is stable across refreshes and tag ordering', () => {
    const article = {
      id: 27,
      category: { id: 3, name: 'Frontend' },
      tags: [{ id: 8, name: 'Vue' }, { id: 2, name: 'JavaScript' }]
    }
    const reordered = { ...article, tags: [...article.tags].reverse() }
    const level = getArticleLogLevel(article)

    expect(['INFO', 'NOTE', 'WARN']).toContain(level)
    expect(getArticleLogLevel(article)).toBe(level)
    expect(getArticleLogLevel(reordered)).toBe(level)
  })
})
