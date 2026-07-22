import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import CozyHero from './CozyHero.vue'

const settings = {
  eyebrow: '欢迎来坐坐',
  description: '真实内容说明',
  primary_text: '读读随笔',
  primary_link: '/blog',
  secondary_text: '看看照片',
  secondary_link: '/photos'
}

const mountHero = photos => mount(CozyHero, {
  props: { settings, title: { before: '慢慢', accent: '记录', after: '。' }, photos },
  global: { stubs: { RouterLink: { props: ['to'], template: '<a :href="to"><slot /></a>' } } }
})

describe('CozyHero', () => {
  it('renders only real photos without duplicating a short list', () => {
    const wrapper = mountHero([{ id: 7, title: '山间', cover: '/uploads/mountain.jpg', metadata: { location: '莫干山' } }])

    expect(wrapper.findAll('.cozy-polaroid')).toHaveLength(1)
    expect(wrapper.find('img').attributes('src')).toBe('/uploads/mountain.jpg')
    expect(wrapper.find('img').attributes('alt')).toBe('山间')
    expect(wrapper.text()).toContain('莫干山')
  })

  it('limits the stack to three and replaces failed images in place', async () => {
    const photos = Array.from({ length: 4 }, (_, index) => ({ id: index + 1, title: `照片${index + 1}`, images: [{ url: `/uploads/${index + 1}.jpg` }] }))
    const wrapper = mountHero(photos)

    expect(wrapper.findAll('.cozy-polaroid')).toHaveLength(3)
    await wrapper.find('img').trigger('error')
    expect(wrapper.findAll('.cozy-polaroid')).toHaveLength(3)
    expect(wrapper.find('.cozy-photo-placeholder').exists()).toBe(true)
  })

  it('shows an explicit empty state without fabricated images', () => {
    const wrapper = mountHero([])

    expect(wrapper.find('img').exists()).toBe(false)
    expect(wrapper.text()).toContain('照片墙还空着')
  })
})
