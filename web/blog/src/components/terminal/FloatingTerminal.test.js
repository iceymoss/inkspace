import { afterEach, beforeEach, describe, expect, it } from 'vitest'
import { DOMWrapper, mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createMemoryHistory, createRouter } from 'vue-router'
import FloatingTerminal from './FloatingTerminal.vue'
import { useTerminalStore } from '@/stores/terminal'

describe('FloatingTerminal', () => {
  let wrapper

  beforeEach(async () => {
    window.sessionStorage.clear()
    const pinia = createPinia()
    setActivePinia(pinia)
    const router = createRouter({ history: createMemoryHistory(), routes: [{ path: '/', component: { template: '<div />' } }] })
    await router.push('/')
    await router.isReady()
    wrapper = mount(FloatingTerminal, { global: { plugins: [pinia, router] } })
  })

  afterEach(() => {
    wrapper.unmount()
    document.body.innerHTML = ''
  })

  function terminalDOM() {
    return new DOMWrapper(document.body)
  }

  it('opens, focuses, and runs the help built-in', async () => {
    const store = useTerminalStore()
    store.open()
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    const dom = terminalDOM()
    const input = dom.find('input')

    expect(document.activeElement).toBe(input.element)
    await input.setValue('help')
    await dom.find('form').trigger('submit')

    expect(store.commandHistory).toEqual(['help'])
    expect(wrapper.emitted('submit')?.[0]).toEqual(['help'])
  })

  it('recalls commands with arrow keys and clears output with Control L', async () => {
    const store = useTerminalStore()
    store.addCommand('first')
    store.addCommand('second')
    store.appendOutput('existing output')
    store.open()
    await wrapper.vm.$nextTick()
    const input = terminalDOM().find('input')

    await input.trigger('keydown', { key: 'ArrowUp' })
    expect(input.element.value).toBe('second')
    await input.trigger('keydown', { key: 'ArrowUp' })
    expect(input.element.value).toBe('first')
    await input.trigger('keydown', { key: 'l', ctrlKey: true })
    expect(store.output).toEqual([])
  })

  it('accepts only full yes and executes a replacement command after cancellation', async () => {
    const store = useTerminalStore()
    store.open()
    store.setPendingConfirmation({ id: 'favorite-12', message: 'favorite article 12?' })
    await wrapper.vm.$nextTick()
    const dom = terminalDOM()
    const input = dom.find('input')

    await input.setValue('y')
    await dom.find('form').trigger('submit')
    expect(wrapper.emitted('confirm')?.[0][0]).toMatchObject({ id: 'favorite-12', accepted: false })
    expect(wrapper.emitted('submit')?.[0]).toEqual(['y'])

    store.setPendingConfirmation({ id: 'favorite-12-again', message: 'favorite article 12?' })
    await input.setValue('yes')
    await dom.find('form').trigger('submit')
    expect(wrapper.emitted('confirm')?.[1][0]).toMatchObject({ id: 'favorite-12-again', accepted: true })
    expect(wrapper.emitted('submit')).toHaveLength(1)
  })

  it('keeps input focus while pending and restores it when a command finishes', async () => {
    const store = useTerminalStore()
    store.open()
    await wrapper.vm.$nextTick()
    const input = terminalDOM().find('input')

    store.setPending(true)
    await wrapper.vm.$nextTick()
    expect(input.attributes('readonly')).toBeDefined()
    expect(document.activeElement).toBe(input.element)

    store.setPending(false)
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    expect(document.activeElement).toBe(input.element)
  })

  it('completes commands, subcommands, and virtual directories with Tab', async () => {
    const store = useTerminalStore()
    store.open()
    await wrapper.vm.$nextTick()
    const input = terminalDOM().find('input')

    await input.setValue('op')
    await input.trigger('keydown', { key: 'Tab' })
    expect(input.element.value).toBe('open ')

    await input.setValue('open b')
    await input.trigger('keydown', { key: 'Tab' })
    expect(input.element.value).toBe('open blog ')

    await input.setValue('cd ../b')
    await input.trigger('keydown', { key: 'Tab' })
    expect(input.element.value).toBe('cd ../blog/ ')

    await input.setValue('theme c')
    await input.trigger('keydown', { key: 'Tab' })
    expect(input.element.value).toBe('theme cozy ')

    await input.setValue('theme s')
    await input.trigger('keydown', { key: 'Tab' })
    expect(input.element.value).toBe('theme swiss ')
    expect(document.activeElement).toBe(input.element)
  })

  it('prints ambiguous Tab candidates without moving focus', async () => {
    const store = useTerminalStore()
    store.open()
    await wrapper.vm.$nextTick()
    const input = terminalDOM().find('input')

    await input.setValue('f')
    await input.trigger('keydown', { key: 'Tab' })

    expect(store.output.at(-1)?.text).toContain('favorite')
    expect(store.output.at(-1)?.text).toContain('follow')
    expect(document.activeElement).toBe(input.element)
  })
})
