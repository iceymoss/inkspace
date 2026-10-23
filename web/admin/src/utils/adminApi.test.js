import { afterEach, describe, expect, it } from 'vitest'
import adminApi from './adminApi'

const originalAdapter = adminApi.defaults.adapter

afterEach(() => {
  adminApi.defaults.adapter = originalAdapter
})

describe('adminApi suppressed errors', () => {
  it('preserves a safe backend message without showing a global error', async () => {
    adminApi.defaults.adapter = async (config) => Promise.reject({
      config,
      response: {
        status: 403,
        data: { message: '初始化密钥错误或未配置' }
      }
    })

    await expect(adminApi.post('/admin/auth/bootstrap', {}, {
      skipAuth: true,
      suppressGlobalError: true
    })).rejects.toThrow('初始化密钥错误或未配置')
  })
})
