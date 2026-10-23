import { mount, flushPromises } from '@vue/test-utils'
import { defineComponent, h, nextTick } from 'vue'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import Bootstrap from './Bootstrap.vue'
import adminApi from '@/utils/adminApi'

vi.mock('@/utils/adminApi', () => ({ default: { get: vi.fn(), post: vi.fn() } }))
const router = vi.hoisted(() => ({ push: vi.fn() }))
vi.mock('vue-router', () => ({ useRouter: () => router }))

const elementStubs = {
  'el-card': defineComponent({ setup(_, { slots }) { return () => h('section', slots.default?.()) } }),
  'el-form': defineComponent({
    props: ['model', 'rules'],
    emits: ['submit'],
    setup(props, { expose, slots }) {
      expose({
        async validate() {
          for (const [field, rules] of Object.entries(props.rules)) {
            const value = props.model[field]
            for (const rule of rules) {
              if (rule.required && !value) throw new Error('invalid')
              if (rule.min && value.length < rule.min) throw new Error('invalid')
              if (rule.max && value.length > rule.max) throw new Error('invalid')
              if (rule.type === 'email' && !/^[^@\s]+@[^@\s]+\.[^@\s]+$/.test(value)) throw new Error('invalid')
              if (rule.validator) {
                await new Promise((resolve, reject) => rule.validator({}, value, (error) => error ? reject(error) : resolve()))
              }
            }
          }
        }
      })
      return () => h('form', slots.default?.())
    }
  }),
  'el-form-item': defineComponent({ setup(_, { slots }) { return () => h('label', slots.default?.()) } }),
  'el-input': defineComponent({
    props: ['modelValue', 'type'], emits: ['update:modelValue'],
    setup(props, { emit }) {
      return () => h('input', {
        type: props.type || 'text', value: props.modelValue,
        onInput: (event) => emit('update:modelValue', event.target.value)
      })
    }
  }),
  'el-button': defineComponent({ props: ['nativeType'], setup(_, { slots }) { return () => h('button', { type: 'submit' }, slots.default?.()) } })
}

function mountBootstrap() {
  return mount(Bootstrap, { global: { stubs: elementStubs } })
}

describe('Bootstrap', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
    adminApi.get.mockResolvedValue({ data: { state: 'pending', secret_configured: true } })
  })

  it('shows maintenance feedback if bootstrap status cannot be read without exposing the error', async () => {
    adminApi.get.mockRejectedValue(new Error('database unavailable'))
    const wrapper = mountBootstrap()
    await flushPromises()
    expect(wrapper.text()).toContain('系统暂不可用')
    expect(wrapper.text()).not.toContain('database unavailable')
  })

  it('does not offer the form when the server has already initialized', async () => {
    adminApi.get.mockResolvedValue({ data: { state: 'initialized', secret_configured: true } })
    const wrapper = mountBootstrap()
    await flushPromises()
    expect(wrapper.text()).toContain('系统已完成初始化')
    expect(wrapper.find('form').exists()).toBe(false)
  })

  it('rejects mismatched confirmation and submits validated credentials', async () => {
    adminApi.post.mockResolvedValue({ data: {} })
    const wrapper = mountBootstrap()
    await flushPromises()
    const inputs = wrapper.findAll('input')
    for (const [index, value] of ['setup-secret', 'first-admin', 'admin@example.com', 'secure-pass-123', 'different-pass'].entries()) {
      await inputs[index].setValue(value)
    }

    await wrapper.vm.submit()
    await flushPromises()
    expect(adminApi.post).not.toHaveBeenCalled()
    expect(wrapper.text()).toContain('两次输入的密码不一致')

    await inputs[4].setValue('secure-pass-123')
    await wrapper.vm.submit()
    await flushPromises()
    expect(adminApi.post).toHaveBeenCalledWith('/admin/auth/bootstrap', {
      secret: 'setup-secret', username: 'first-admin', email: 'admin@example.com',
      password: 'secure-pass-123', confirm_password: 'secure-pass-123'
    }, { skipAuth: true, suppressGlobalError: true })
    expect(router.push).toHaveBeenCalledWith('/login')
  })

  it('does not submit an empty form', async () => {
    const wrapper = mountBootstrap()
    await flushPromises()

    await wrapper.vm.submit()
    await flushPromises()

    expect(adminApi.post).not.toHaveBeenCalled()
    expect(router.push).not.toHaveBeenCalled()
  })

  it('shows safe failure feedback without exposing raw API errors', async () => {
    const wrapper = mountBootstrap()
    await flushPromises()
    const inputs = wrapper.findAll('input')
    for (const [index, value] of ['setup-secret', 'first-admin', 'admin@example.com', 'secure-pass-123', 'secure-pass-123'].entries()) {
      await inputs[index].setValue(value)
    }
    adminApi.post.mockRejectedValue(new Error('internal SQL error'))
    await wrapper.vm.submit()
    await flushPromises()
    expect(wrapper.text()).toContain('提交失败，请稍后重试')
    expect(wrapper.text()).not.toContain('internal SQL error')
    expect(router.push).not.toHaveBeenCalled()
  })

  it('shows unavailable messaging when the initialization secret is not configured', async () => {
    adminApi.get.mockResolvedValue({ data: { state: 'pending', secret_configured: false } })
    const wrapper = mountBootstrap()
    await flushPromises()
    expect(wrapper.text()).toContain('初始化暂不可用')
    expect(wrapper.find('form').exists()).toBe(false)
  })

  it('redirects to login if another request completed initialization', async () => {
    adminApi.post.mockRejectedValue(new Error('系统已初始化，请登录'))
    const wrapper = mountBootstrap()
    await flushPromises()
    const inputs = wrapper.findAll('input')
    for (const [index, value] of ['setup-secret', 'first-admin', 'admin@example.com', 'secure-pass-123', 'secure-pass-123'].entries()) {
      await inputs[index].setValue(value)
    }

    await wrapper.vm.submit()
    await flushPromises()

    expect(router.push).toHaveBeenCalledWith('/login')
  })

  it('fails closed when the server returns an unknown state', async () => {
    adminApi.get.mockResolvedValue({ data: { state: 'unexpected' } })
    const wrapper = mountBootstrap()
    await flushPromises()
    expect(wrapper.text()).toContain('系统暂不可用')
    expect(wrapper.find('form').exists()).toBe(false)
  })
})
