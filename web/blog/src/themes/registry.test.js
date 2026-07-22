import { describe, expect, it } from 'vitest'
import { availableThemeIds, getTheme, themeIds, themeRegistry } from './registry'

describe('theme registry', () => {
  it('opens cozy while keeping swiss unavailable', () => {
    const cozy = getTheme('cozy')

    expect(cozy).toMatchObject({ status: 'available', stylesheet: 'cozy.css', defaultColorScheme: 'light' })
    expect(availableThemeIds).toEqual(new Set(['magazine', 'terminal', 'cozy']))
    expect(themeIds).toEqual(new Set(themeRegistry.map(theme => theme.id)))
    expect(availableThemeIds.has('swiss')).toBe(false)
    expect(getTheme('unknown').id).toBe('magazine')
  })
})
