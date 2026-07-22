import { beforeEach, describe, expect, it } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { MAX_HISTORY, MAX_OUTPUT_LINES, STORAGE_KEY, useTerminalStore } from './terminal'

describe('terminal store', () => {
  beforeEach(() => {
    window.sessionStorage.clear()
    setActivePinia(createPinia())
  })

  it('persists normalized session state', async () => {
    const store = useTerminalStore()
    store.open({ left: 10, top: 20, width: 300, height: 200 })
    store.setBounds({ x: -10, width: 10, height: 9999 })
    store.appendOutput('ready', 'success')
    store.addCommand('help')
    store.setPendingConfirmation({ id: 'delete', message: 'Continue?' })
    store.signalRefresh('articles')

    await Promise.resolve()
    const saved = JSON.parse(window.sessionStorage.getItem(STORAGE_KEY))
    expect(saved.isOpen).toBe(true)
    expect(saved.bounds).toEqual({ x: 0, y: 96, width: 420, height: 900 })
    expect(saved.output[0]).toMatchObject({ text: 'ready', type: 'success' })
    expect(saved.commandHistory).toEqual(['help'])
    expect(saved.pendingConfirmation).toBeUndefined()
    expect(saved.refreshSignals).toBeUndefined()
  })

  it('caps output and de-duplicates bounded history', () => {
    const store = useTerminalStore()
    for (let index = 0; index < MAX_OUTPUT_LINES + 5; index += 1) store.appendOutput(`line ${index}`)
    for (let index = 0; index < MAX_HISTORY + 5; index += 1) store.addCommand(`command ${index}`)
    store.addCommand('command 50')

    expect(store.output).toHaveLength(MAX_OUTPUT_LINES)
    expect(store.output[0].text).toBe('line 5')
    expect(store.commandHistory).toHaveLength(MAX_HISTORY)
    expect(store.commandHistory.at(-1)).toBe('command 50')
    expect(store.commandHistory.filter(item => item === 'command 50')).toHaveLength(1)
  })

  it('recovers safely from invalid persisted data', () => {
    window.sessionStorage.setItem(STORAGE_KEY, '{invalid')
    setActivePinia(createPinia())

    const store = useTerminalStore()

    expect(store.isOpen).toBe(false)
    expect(store.output).toEqual([])
    expect(store.bounds.width).toBe(640)
  })

  it('never restores a pending confirmation', () => {
    window.sessionStorage.setItem(STORAGE_KEY, JSON.stringify({
      pendingConfirmation: { id: 'stale', message: 'Continue?' }
    }))
    setActivePinia(createPinia())

    expect(useTerminalStore().pendingConfirmation).toBeNull()
  })
})
