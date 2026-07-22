import { computed, ref, watch } from 'vue'
import { defineStore } from 'pinia'
import api from '@/utils/api'
import { useUserStore } from '@/stores/user'
import { availableThemeIds } from '@/themes/registry'

const DEFAULT_PREFERENCE = Object.freeze({ ui_theme: 'magazine', color_scheme: 'system' })
const GUEST_CACHE_KEY = 'inkspace_guest_appearance_v1'
const USER_CACHE_PREFIX = 'inkspace_user_appearance_v1:'
const MIGRATION_KEY = 'inkspace_appearance_migrated_v1'
const COLOR_SCHEMES = new Set(['system', 'light', 'dark'])
const SPECIAL_THEMES = new Set(['holiday', 'mourning'])
const SPECIAL_REFRESH_INTERVAL = 5 * 60 * 1000

function normalizePreference(value) {
  return {
    ui_theme: availableThemeIds.has(value?.ui_theme) ? value.ui_theme : DEFAULT_PREFERENCE.ui_theme,
    color_scheme: COLOR_SCHEMES.has(value?.color_scheme)
      ? value.color_scheme
      : DEFAULT_PREFERENCE.color_scheme
  }
}

function readCache(key) {
  try {
    return normalizePreference(JSON.parse(storageGet(key)))
  } catch {
    return { ...DEFAULT_PREFERENCE }
  }
}

function writeCache(key, preference) {
  storageSet(key, JSON.stringify(normalizePreference(preference)))
}

function storageGet(key) {
  try {
    return window.localStorage.getItem(key)
  } catch {
    return null
  }
}

function storageSet(key, value) {
  try {
    window.localStorage.setItem(key, value)
    return true
  } catch {
    return false
  }
}

function storageRemove(key) {
  try {
    window.localStorage.removeItem(key)
  } catch {
    // Appearance caching is best-effort; the in-memory preference remains valid.
  }
}

function tokenUserId(token) {
  try {
    const payload = token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')
    return JSON.parse(decodeURIComponent(escape(atob(payload)))).user_id || null
  } catch {
    return null
  }
}

function resolveScheme(colorScheme) {
  if (colorScheme !== 'system') return colorScheme
  return window.matchMedia?.('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

function applyDocumentPreference(preference) {
  const normalized = normalizePreference(preference)
  document.documentElement.dataset.uiTheme = normalized.ui_theme
  document.documentElement.dataset.theme = resolveScheme(normalized.color_scheme)
  return normalized
}

function removeLegacyTheme() {
  if (storageGet(MIGRATION_KEY)) return
  storageRemove('site_theme')
  document.body.classList.forEach((className) => {
    if (className.startsWith('theme-')) document.body.classList.remove(className)
  })
  Array.from(document.documentElement.style).forEach((property) => {
    if (property.startsWith('--theme-')) document.documentElement.style.removeProperty(property)
  })
  storageSet(MIGRATION_KEY, '1')
}

function cachedPreferenceForCurrentSession() {
  const token = storageGet('token') || ''
  const userId = tokenUserId(token)
  return readCache(token && userId ? `${USER_CACHE_PREFIX}${userId}` : GUEST_CACHE_KEY)
}

export function bootstrapCachedAppearance() {
  removeLegacyTheme()
  return applyDocumentPreference(cachedPreferenceForCurrentSession())
}

export const useAppearanceStore = defineStore('appearance', () => {
  const savedPreference = ref(cachedPreferenceForCurrentSession())
  const previewPreference = ref(null)
  const siteOverride = ref(null)
  const initialized = ref(false)
  const syncing = ref(false)
  const saving = ref(false)
  const lastSpecialRefresh = ref(0)
  const resolvedColorScheme = ref(resolveScheme(savedPreference.value.color_scheme))

  const activePreference = computed(() => previewPreference.value || savedPreference.value)
  const isPreviewing = computed(() => previewPreference.value !== null)

  let mediaQuery = null
  let stopUserWatch = null
  let requestVersion = 0

  function apply(preference = activePreference.value) {
    const normalized = applyDocumentPreference(preference)
    resolvedColorScheme.value = resolveScheme(normalized.color_scheme)
  }

  function currentUserId() {
    const userStore = useUserStore()
    return userStore.user?.id || tokenUserId(userStore.token)
  }

  function currentCacheKey() {
    const userStore = useUserStore()
    const userId = currentUserId()
    return userStore.token && userId ? `${USER_CACHE_PREFIX}${userId}` : GUEST_CACHE_KEY
  }

  function useCachedContext() {
    savedPreference.value = readCache(currentCacheKey())
    previewPreference.value = null
    apply(savedPreference.value)
  }

  async function loadAccountPreference() {
    const userStore = useUserStore()
    if (!userStore.token || !currentUserId()) return
    const version = ++requestVersion
    syncing.value = true
    try {
      const response = await api.get('/profile/appearance')
      if (version !== requestVersion || !userStore.token) return
      savedPreference.value = normalizePreference(response.data)
      previewPreference.value = null
      writeCache(currentCacheKey(), savedPreference.value)
      apply(savedPreference.value)
    } catch (error) {
      console.error('Failed to load appearance preference:', error)
    } finally {
      if (version === requestVersion) syncing.value = false
    }
  }

  function handleSessionChange() {
    requestVersion += 1
    useCachedContext()
    loadAccountPreference()
    const userStore = useUserStore()
    if (!userStore.token && storageGet(GUEST_CACHE_KEY) === null) refreshSiteOverride(true)
  }

  function preview(preference) {
    if (!availableThemeIds.has(preference?.ui_theme)) return false
    const normalized = normalizePreference(preference)
    previewPreference.value = normalized
    apply(normalized)
    return true
  }

  function cancelPreview() {
    previewPreference.value = null
    apply(savedPreference.value)
  }

  async function save(preference = activePreference.value) {
    const userStore = useUserStore()
    if (!userStore.token || !availableThemeIds.has(preference?.ui_theme)) {
      throw new Error('当前主题不可保存')
    }
    const candidate = normalizePreference(preference)

    const rollback = { ...savedPreference.value }
    const sessionToken = userStore.token
    const sessionUserID = currentUserId()
    const sessionCacheKey = currentCacheKey()
    const sameSession = () => userStore.token === sessionToken && currentUserId() === sessionUserID
    preview(candidate)
    saving.value = true
    try {
      const response = await api.put('/profile/appearance', candidate)
      const saved = normalizePreference(response.data || candidate)
      if (!sameSession()) return saved
      savedPreference.value = saved
      previewPreference.value = null
      writeCache(sessionCacheKey, savedPreference.value)
      apply(savedPreference.value)
      return savedPreference.value
    } catch (error) {
      if (sameSession()) {
        savedPreference.value = rollback
        previewPreference.value = null
        apply(rollback)
      }
      throw error
    } finally {
      saving.value = false
    }
  }

  function saveGuestPreference(preference) {
    if (!availableThemeIds.has(preference?.ui_theme)) return false
    const normalized = normalizePreference(preference)
    savedPreference.value = normalized
    previewPreference.value = null
    writeCache(GUEST_CACHE_KEY, normalized)
    apply(normalized)
    return true
  }

  async function refreshSiteOverride(force = false) {
    const now = Date.now()
    if (!force && now - lastSpecialRefresh.value < SPECIAL_REFRESH_INTERVAL) return
    try {
      const response = await api.get('/settings/public')
      const settings = response.data || {}
      const userStore = useUserStore()
      if (!userStore.token && storageGet(GUEST_CACHE_KEY) === null) {
        const guestDefault = normalizePreference({
          ui_theme: settings.default_guest_ui_theme,
          color_scheme: settings.default_guest_color_scheme
        })
        savedPreference.value = guestDefault
        previewPreference.value = null
        writeCache(GUEST_CACHE_KEY, guestDefault)
        apply(guestDefault)
      }
      siteOverride.value = SPECIAL_THEMES.has(settings.site_theme) ? settings.site_theme : null
      const root = document.documentElement
      if (siteOverride.value) root.dataset.siteTheme = siteOverride.value
      else delete root.dataset.siteTheme
      root.style.setProperty('--site-holiday-bg', settings.holiday_bg_primary || '#fff5f2')
      root.style.setProperty('--site-holiday-bg-soft', settings.holiday_bg_secondary || '#ffe9e4')
      root.style.setProperty('--site-holiday-ink', settings.holiday_text_primary || '#741f1f')
      root.style.setProperty('--site-holiday-accent', settings.holiday_primary || '#b52a32')
      lastSpecialRefresh.value = now
    } catch (error) {
      console.error('Failed to load site appearance override:', error)
    }
  }

  function handleSystemSchemeChange() {
    if (activePreference.value.color_scheme === 'system') apply()
  }

  function handleWindowFocus() {
    refreshSiteOverride()
  }

  function initialize() {
    if (initialized.value) return
    initialized.value = true
    apply(savedPreference.value)

    mediaQuery = window.matchMedia?.('(prefers-color-scheme: dark)')
    mediaQuery?.addEventListener?.('change', handleSystemSchemeChange)
    window.addEventListener('focus', handleWindowFocus)

    const userStore = useUserStore()
    stopUserWatch = watch(
      () => [userStore.token, userStore.user?.id],
      handleSessionChange
    )
    if (userStore.token) loadAccountPreference()
    refreshSiteOverride(true)
  }

  function dispose() {
    mediaQuery?.removeEventListener?.('change', handleSystemSchemeChange)
    window.removeEventListener('focus', handleWindowFocus)
    stopUserWatch?.()
    initialized.value = false
  }

  return {
    savedPreference,
    previewPreference,
    activePreference,
    resolvedColorScheme,
    siteOverride,
    initialized,
    syncing,
    saving,
    isPreviewing,
    initialize,
    dispose,
    preview,
    cancelPreview,
    save,
    saveGuestPreference,
    refreshSiteOverride
  }
})
