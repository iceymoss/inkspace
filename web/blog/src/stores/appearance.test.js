import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import api from '@/utils/api'
import { useAppearanceStore, bootstrapCachedAppearance } from './appearance'
import { useUserStore } from './user'

vi.mock('@/utils/api', () => ({
  default: {
    get: vi.fn(),
    put: vi.fn()
  }
}))

const matchMediaListeners = new Set()
let systemDark = false

function installMatchMedia() {
  window.matchMedia = vi.fn(() => ({
    get matches() {
      return systemDark
    },
    addEventListener: (_event, listener) => matchMediaListeners.add(listener),
    removeEventListener: (_event, listener) => matchMediaListeners.delete(listener)
  }))
}

function emitSystemSchemeChange() {
  matchMediaListeners.forEach((listener) => listener({ matches: systemDark }))
}

describe('appearance store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    window.localStorage.clear()
    document.documentElement.removeAttribute('data-ui-theme')
    document.documentElement.removeAttribute('data-theme')
    document.documentElement.removeAttribute('data-site-theme')
    document.documentElement.removeAttribute('style')
    document.body.className = ''
    systemDark = false
    matchMediaListeners.clear()
    installMatchMedia()
    vi.clearAllMocks()
    api.get.mockResolvedValue({ data: {} })
  })

  it('falls back to magazine and system for an unknown guest preference', () => {
    window.localStorage.setItem('inkspace_guest_appearance_v1', JSON.stringify({
      ui_theme: 'unknown',
      color_scheme: 'sepia'
    }))

    bootstrapCachedAppearance()

    expect(document.documentElement.dataset.uiTheme).toBe('magazine')
    expect(document.documentElement.dataset.theme).toBe('light')
  })

  it('bootstraps a cached terminal preference', () => {
    window.localStorage.setItem('inkspace_guest_appearance_v1', JSON.stringify({
      ui_theme: 'terminal',
      color_scheme: 'dark'
    }))

    bootstrapCachedAppearance()

    expect(document.documentElement.dataset.uiTheme).toBe('terminal')
    expect(document.documentElement.dataset.theme).toBe('dark')
  })

  it('bootstraps a cached cozy preference', () => {
    window.localStorage.setItem('inkspace_guest_appearance_v1', JSON.stringify({ ui_theme: 'cozy', color_scheme: 'light' }))

    bootstrapCachedAppearance()

    expect(document.documentElement.dataset.uiTheme).toBe('cozy')
    expect(document.documentElement.dataset.theme).toBe('light')
  })

  it('previews a preference without saving and restores it on cancel', () => {
    const store = useAppearanceStore()

    store.preview({ ui_theme: 'terminal', color_scheme: 'dark' })
    expect(store.isPreviewing).toBe(true)
    expect(document.documentElement.dataset.uiTheme).toBe('terminal')
    expect(document.documentElement.dataset.theme).toBe('dark')
    expect(api.put).not.toHaveBeenCalled()

    store.cancelPreview()
    expect(store.isPreviewing).toBe(false)
    expect(document.documentElement.dataset.uiTheme).toBe('magazine')
    expect(document.documentElement.dataset.theme).toBe('light')
  })

  it('rolls back a failed account save', async () => {
    const userStore = useUserStore()
    userStore.token = 'token-without-readable-payload'
    userStore.user = { id: 7 }
    const store = useAppearanceStore()
    api.put.mockRejectedValue(new Error('network unavailable'))

    await expect(store.save({ ui_theme: 'terminal', color_scheme: 'dark' })).rejects.toThrow()

    expect(store.savedPreference.ui_theme).toBe('magazine')
    expect(store.savedPreference.color_scheme).toBe('system')
    expect(store.isPreviewing).toBe(false)
    expect(document.documentElement.dataset.theme).toBe('light')
  })

  it('saves terminal to the current account cache', async () => {
    const userStore = useUserStore()
    userStore.token = 'token-without-readable-payload'
    userStore.user = { id: 7 }
    const store = useAppearanceStore()
    api.put.mockResolvedValue({ data: { ui_theme: 'terminal', color_scheme: 'dark' } })

    await store.save({ ui_theme: 'terminal', color_scheme: 'dark' })

    expect(store.savedPreference).toEqual({ ui_theme: 'terminal', color_scheme: 'dark' })
    expect(JSON.parse(window.localStorage.getItem('inkspace_user_appearance_v1:7'))).toEqual({
      ui_theme: 'terminal', color_scheme: 'dark'
    })
  })

  it('saves cozy and Swiss preferences', async () => {
    const userStore = useUserStore()
    userStore.token = 'token-without-readable-payload'
    userStore.user = { id: 7 }
    const store = useAppearanceStore()
    api.put
      .mockResolvedValueOnce({ data: { ui_theme: 'cozy', color_scheme: 'light' } })
      .mockResolvedValueOnce({ data: { ui_theme: 'swiss', color_scheme: 'dark' } })

    await store.save({ ui_theme: 'cozy', color_scheme: 'light' })

    expect(store.savedPreference).toEqual({ ui_theme: 'cozy', color_scheme: 'light' })
    expect(store.preview({ ui_theme: 'swiss', color_scheme: 'dark' })).toBe(true)
    await store.save({ ui_theme: 'swiss', color_scheme: 'dark' })
    expect(store.savedPreference).toEqual({ ui_theme: 'swiss', color_scheme: 'dark' })
    expect(api.put).toHaveBeenCalledTimes(2)
  })

  it('keeps a site override independent from terminal preference', async () => {
    api.get.mockResolvedValue({ data: { site_theme: 'holiday' } })
    const store = useAppearanceStore()
    store.saveGuestPreference({ ui_theme: 'terminal', color_scheme: 'dark' })

    await store.refreshSiteOverride(true)

    expect(document.documentElement.dataset.uiTheme).toBe('terminal')
    expect(document.documentElement.dataset.siteTheme).toBe('holiday')
    expect(store.savedPreference.ui_theme).toBe('terminal')
  })

  it('applies the public default appearance only to a guest without cache', async () => {
    api.get.mockResolvedValue({
      data: { default_guest_ui_theme: 'swiss', default_guest_color_scheme: 'dark' }
    })
    const store = useAppearanceStore()

    store.initialize()
    await vi.waitFor(() => expect(store.savedPreference).toEqual({ ui_theme: 'swiss', color_scheme: 'dark' }))

    expect(document.documentElement.dataset.uiTheme).toBe('swiss')
    expect(JSON.parse(window.localStorage.getItem('inkspace_guest_appearance_v1'))).toEqual({
      ui_theme: 'swiss', color_scheme: 'dark'
    })
    store.dispose()
  })

  it('does not overwrite an existing guest preference with the public default', async () => {
    window.localStorage.setItem('inkspace_guest_appearance_v1', JSON.stringify({
      ui_theme: 'cozy', color_scheme: 'light'
    }))
    api.get.mockResolvedValue({
      data: { default_guest_ui_theme: 'swiss', default_guest_color_scheme: 'dark' }
    })
    const store = useAppearanceStore()

    store.initialize()
    await vi.waitFor(() => expect(api.get).toHaveBeenCalledWith('/settings/public'))

    expect(store.savedPreference).toEqual({ ui_theme: 'cozy', color_scheme: 'light' })
    expect(document.documentElement.dataset.uiTheme).toBe('cozy')
    store.dispose()
  })

  it('loads the public default after logout when no guest cache exists', async () => {
    const userStore = useUserStore()
    userStore.user = { id: 7 }
    userStore.token = 'account-token'
    api.get.mockImplementation((url) => Promise.resolve({
      data: url === '/profile/appearance'
        ? { ui_theme: 'terminal', color_scheme: 'light' }
        : { default_guest_ui_theme: 'cozy', default_guest_color_scheme: 'dark' }
    }))
    const store = useAppearanceStore()
    store.initialize()
    await vi.waitFor(() => expect(store.savedPreference.ui_theme).toBe('terminal'))

    userStore.logout()
    await vi.waitFor(() => expect(store.savedPreference).toEqual({ ui_theme: 'cozy', color_scheme: 'dark' }))

    expect(JSON.parse(window.localStorage.getItem('inkspace_guest_appearance_v1'))).toEqual({
      ui_theme: 'cozy', color_scheme: 'dark'
    })
    store.dispose()
  })

  it('normalizes invalid public default appearance values', async () => {
    api.get.mockResolvedValue({
      data: { default_guest_ui_theme: 'unknown', default_guest_color_scheme: 'sepia' }
    })
    const store = useAppearanceStore()

    store.initialize()
    await vi.waitFor(() => expect(window.localStorage.getItem('inkspace_guest_appearance_v1')).not.toBeNull())

    expect(store.savedPreference).toEqual({ ui_theme: 'magazine', color_scheme: 'system' })
    store.dispose()
  })

  it('does not apply an account save after the session changes', async () => {
    window.localStorage.setItem('inkspace_guest_appearance_v1', JSON.stringify({
      ui_theme: 'magazine',
      color_scheme: 'dark'
    }))
    const userStore = useUserStore()
    userStore.token = 'account-token'
    userStore.user = { id: 7 }
    const store = useAppearanceStore()
    store.initialize()

    let resolveSave
    api.put.mockReturnValue(new Promise((resolve) => { resolveSave = resolve }))
    const savePromise = store.save({ ui_theme: 'magazine', color_scheme: 'light' })
    userStore.logout()
    await vi.waitFor(() => expect(store.savedPreference.color_scheme).toBe('dark'))

    resolveSave({ data: { ui_theme: 'magazine', color_scheme: 'light' } })
    await savePromise

    expect(store.savedPreference.color_scheme).toBe('dark')
    expect(JSON.parse(window.localStorage.getItem('inkspace_guest_appearance_v1')).color_scheme).toBe('dark')
    expect(window.localStorage.getItem('inkspace_user_appearance_v1:7')).toBeNull()
    store.dispose()
  })

  it('keeps guest and account caches isolated across session changes', async () => {
    window.localStorage.setItem('inkspace_guest_appearance_v1', JSON.stringify({
      ui_theme: 'magazine',
      color_scheme: 'dark'
    }))
    const userStore = useUserStore()
    const store = useAppearanceStore()
    store.initialize()

    api.get.mockImplementation((url) => Promise.resolve({
      data: url === '/profile/appearance'
        ? { ui_theme: 'magazine', color_scheme: 'light' }
        : {}
    }))
    userStore.user = { id: 9 }
    userStore.token = 'account-token'
    await vi.waitFor(() => expect(store.savedPreference.color_scheme).toBe('light'))

    expect(JSON.parse(window.localStorage.getItem('inkspace_guest_appearance_v1')).color_scheme).toBe('dark')
    expect(JSON.parse(window.localStorage.getItem('inkspace_user_appearance_v1:9')).color_scheme).toBe('light')

    userStore.logout()
    await vi.waitFor(() => expect(store.savedPreference.color_scheme).toBe('dark'))
    store.dispose()
  })

  it('reacts to system changes only while following the system', () => {
    const store = useAppearanceStore()
    store.initialize()

    store.saveGuestPreference({ ui_theme: 'terminal', color_scheme: 'system' })

    systemDark = true
    emitSystemSchemeChange()
    expect(document.documentElement.dataset.theme).toBe('dark')

    store.saveGuestPreference({ ui_theme: 'terminal', color_scheme: 'light' })
    systemDark = false
    emitSystemSchemeChange()
    expect(document.documentElement.dataset.theme).toBe('light')
    systemDark = true
    emitSystemSchemeChange()
    expect(document.documentElement.dataset.theme).toBe('light')
    store.dispose()
  })

  it('removes legacy theme state during bootstrap', () => {
    window.localStorage.setItem('site_theme', 'night')
    document.body.classList.add('theme-night')
    document.documentElement.style.setProperty('--theme-primary', '#000')

    bootstrapCachedAppearance()

    expect(window.localStorage.getItem('site_theme')).toBeNull()
    expect(document.body.classList.contains('theme-night')).toBe(false)
    expect(document.documentElement.style.getPropertyValue('--theme-primary')).toBe('')
  })
})
