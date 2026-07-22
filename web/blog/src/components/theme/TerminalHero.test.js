import { afterEach, describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import TerminalHero from './TerminalHero.vue'

const props = {
  settings: {
    status: 'system online',
    eyebrow: '持续记录',
    description: '真实站点描述',
    primary_text: 'tail -f blog',
    primary_link: '/blog',
    secondary_text: 'ls projects/',
    secondary_link: '/works'
  },
  title: { before: '记录', accent: '此刻', after: '。' },
  stats: { articleCount: 3, workCount: 2, categoryCount: 1 }
}

afterEach(() => {
  vi.useRealTimers()
  vi.restoreAllMocks()
})

describe('TerminalHero', () => {
  it('renders real configuration and remains static with reduced motion', () => {
    vi.spyOn(window, 'matchMedia').mockReturnValue({ matches: true, addEventListener: vi.fn(), removeEventListener: vi.fn() })
    const setTimeoutSpy = vi.spyOn(window, 'setTimeout')
    const wrapper = mount(TerminalHero, { props })

    expect(wrapper.text()).toContain('真实站点描述')
    expect(wrapper.text()).toContain('articles: 3')
    expect(wrapper.text()).toContain('ls photos/')
    expect(wrapper.text()).toContain('tail -f blog')
    expect(wrapper.text()).toContain('ls projects/')
    expect(setTimeoutSpy).not.toHaveBeenCalled()
  })

  it('clears its typing timer when unmounted', () => {
    vi.useFakeTimers()
    vi.spyOn(window, 'matchMedia').mockReturnValue({ matches: false, addEventListener: vi.fn(), removeEventListener: vi.fn() })
    const clearTimeoutSpy = vi.spyOn(window, 'clearTimeout')
    const wrapper = mount(TerminalHero, { props })

    wrapper.unmount()

    expect(clearTimeoutSpy).toHaveBeenCalled()
  })
})
