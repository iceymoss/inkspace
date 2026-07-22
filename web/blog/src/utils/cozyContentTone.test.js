import { describe, expect, it } from 'vitest'
import { getCozyContentTone, getCozyRotation } from './cozyContentTone'

describe('cozy content presentation', () => {
  it.each([
    [{ id: 1, category: { name: '设计思考' } }, 'moss'],
    [{ id: 2, tags: [{ name: '前端工程' }] }, 'sky'],
    [{ id: 3, type: 'photography' }, 'blush']
  ])('maps known semantics to a stable tone', (resource, tone) => {
    expect(getCozyContentTone(resource)).toBe(tone)
  })

  it('uses deterministic bounded fallbacks', () => {
    const resource = { id: 41, title: '未分类内容' }
    expect(getCozyContentTone(resource)).toBe(getCozyContentTone(resource))
    expect(['moss', 'sky', 'blush']).toContain(getCozyContentTone(resource))
    expect(getCozyRotation(resource, 2)).toBe(getCozyRotation(resource, 2))
    expect(Math.abs(getCozyRotation(resource, 2))).toBeLessThanOrEqual(7)
  })
})
