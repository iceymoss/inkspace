import { describe, expect, it } from 'vitest'
import { availableThemeIds, getTheme, themeIds, themeRegistry } from './registry'

describe('theme registry', () => {
  it('opens all four implemented themes', () => {
    const cozy = getTheme('cozy')
    const swiss = getTheme('swiss')

    expect(cozy).toMatchObject({ status: 'available', stylesheet: 'cozy.css', defaultColorScheme: 'light' })
    expect(swiss).toMatchObject({ status: 'available', stylesheet: 'swiss.css', defaultColorScheme: 'light' })
    expect(availableThemeIds).toEqual(new Set(['magazine', 'terminal', 'cozy', 'swiss']))
    expect(themeIds).toEqual(new Set(themeRegistry.map(theme => theme.id)))
    expect(getTheme('unknown').id).toBe('magazine')
  })
})
