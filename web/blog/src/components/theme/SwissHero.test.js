import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import SwissHero from './SwissHero.vue'

const settings = {
  eyebrow: 'A1 / INDEX', location: '', coordinates: '', established: '',
  description: '文章 × 作品 × 摄影 × 知识库', primary_text: 'START READING', primary_link: '/blog',
  secondary_text: 'VIEW WORKS', secondary_link: '/works'
}

describe('SwissHero', () => {
  it('hides empty archive fields and renders real stats', () => {
    const wrapper = mount(SwissHero, {
      props: {
        settings,
        title: { before: 'WRITE & ', accent: 'ARCHIVE', after: '' },
        stats: { articleCount: 12, workCount: 4, publicDocCount: 8 },
        siteName: 'InkSpace'
      }
    })
    expect(wrapper.text()).not.toContain('SHANGHAI')
    expect(wrapper.text()).toContain('12+')
    expect(wrapper.text()).toContain('8+')
  })

  it('shows an unavailable state instead of a false zero', () => {
    const wrapper = mount(SwissHero, {
      props: {
        settings: { ...settings, location: '杭州', coordinates: '30.27°N 120.15°E', established: 'EST. 2020' },
        title: { before: 'ARCHIVE', accent: '', after: '' },
        stats: { articleCount: null, articleError: true, workCount: 0, publicDocCount: 0 },
        siteName: 'Example'
      }
    })
    expect(wrapper.text()).toContain('杭州')
    expect(wrapper.text()).toContain('公开文章暂不可用')
    expect(wrapper.text()).toContain('--')
  })
})
