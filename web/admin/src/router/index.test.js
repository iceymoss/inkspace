import { beforeEach, describe, expect, it, vi } from 'vitest'
import adminApi from '@/utils/adminApi'
import router from './index'

vi.mock('@/stores/admin', () => ({
  useAdminStore: vi.fn(() => ({ isLoggedIn: false }))
}))
vi.mock('@/utils/adminApi', () => ({
  default: { get: vi.fn() }
}))

describe('admin bootstrap navigation', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    window.history.replaceState({}, '', '/')
  })

  it('routes anonymous visits to login through bootstrap when initialization is pending', async () => {
    adminApi.get.mockResolvedValue({ data: { state: 'pending', secret_configured: true } })

    await router.push('/login')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/bootstrap')
    expect(adminApi.get).toHaveBeenCalledWith('/admin/auth/bootstrap/status', {
      skipAuth: true,
      suppressGlobalError: true
    })
  })

  it('routes status failures to the bootstrap maintenance page', async () => {
    adminApi.get.mockRejectedValue(new Error('database unavailable'))

    await router.push('/login')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/bootstrap')
  })

  it('routes unknown status values to the bootstrap maintenance page', async () => {
    adminApi.get.mockResolvedValue({ data: { state: 'unexpected' } })

    await router.push('/login')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/bootstrap')
  })

  it('does not trust a stale local login while initialization is pending', async () => {
    const { useAdminStore } = await import('@/stores/admin')
    useAdminStore.mockReturnValueOnce({ isLoggedIn: true })
    adminApi.get.mockResolvedValue({ data: { state: 'pending', secret_configured: true } })

    await router.push('/')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/bootstrap')
  })

  it('routes initialized visitors away from bootstrap to login', async () => {
    adminApi.get.mockResolvedValue({ data: { state: 'initialized' } })

    await router.push('/login')
    await router.isReady()
    await router.push('/bootstrap')

    expect(router.currentRoute.value.path).toBe('/login')
  })

  it('allows an authenticated administrator into the initialized backend', async () => {
    const { useAdminStore } = await import('@/stores/admin')
    useAdminStore.mockReturnValueOnce({ isLoggedIn: true })
    adminApi.get.mockResolvedValue({ data: { state: 'initialized' } })

    await router.push('/')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/')
  })
})
