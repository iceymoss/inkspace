import { reactive, ref, watch } from 'vue'
import { defineStore } from 'pinia'

const STORAGE_KEY = 'inkspace_terminal_session_v1'
const MAX_OUTPUT_LINES = 200
const MAX_OUTPUT_LENGTH = 4000
const MAX_HISTORY = 100
const MAX_COMMAND_LENGTH = 500
const MAX_REFRESH_SIGNALS = 20
const DEFAULT_BOUNDS = Object.freeze({ x: 32, y: 96, width: 640, height: 420 })

function finiteNumber(value, fallback) {
  return Number.isFinite(value) ? Math.round(value) : fallback
}

function normalizeBounds(value) {
  return {
    x: Math.max(0, finiteNumber(value?.x, DEFAULT_BOUNDS.x)),
    y: Math.max(0, finiteNumber(value?.y, DEFAULT_BOUNDS.y)),
    width: Math.min(1200, Math.max(420, finiteNumber(value?.width, DEFAULT_BOUNDS.width))),
    height: Math.min(900, Math.max(280, finiteNumber(value?.height, DEFAULT_BOUNDS.height)))
  }
}

function normalizeOutput(value) {
  if (!Array.isArray(value)) return []
  return value.slice(-MAX_OUTPUT_LINES).map((entry, index) => ({
    id: finiteNumber(entry?.id, index + 1),
    type: ['command', 'error', 'success', 'system'].includes(entry?.type) ? entry.type : 'output',
    text: String(entry?.text ?? '').slice(0, MAX_OUTPUT_LENGTH)
  }))
}

function normalizeHistory(value) {
  if (!Array.isArray(value)) return []
  return value
    .map(command => String(command).trim().slice(0, MAX_COMMAND_LENGTH))
    .filter(Boolean)
    .slice(-MAX_HISTORY)
}

function normalizeConfirmation(value) {
  if (!value || typeof value !== 'object') return null
  const message = String(value.message ?? '').slice(0, 1000)
  if (!message) return null
  return { id: String(value.id ?? Date.now()).slice(0, 100), message }
}

function normalizeRefreshSignals(value) {
  if (!value || typeof value !== 'object' || Array.isArray(value)) return {}
  return Object.fromEntries(Object.entries(value).slice(0, MAX_REFRESH_SIGNALS).map(([key, count]) => [
    String(key).slice(0, 50),
    Math.max(0, finiteNumber(count, 0))
  ]))
}

function readSession() {
  try {
    const value = JSON.parse(window.sessionStorage.getItem(STORAGE_KEY))
    return value && typeof value === 'object' ? value : {}
  } catch {
    return {}
  }
}

function writeSession(value) {
  try {
    window.sessionStorage.setItem(STORAGE_KEY, JSON.stringify(value))
  } catch {
    // Session persistence is best-effort; the live terminal remains usable.
  }
}

export const useTerminalStore = defineStore('terminal', () => {
  const saved = readSession()
  const isOpen = ref(Boolean(saved.isOpen))
  const isMinimized = ref(Boolean(saved.isMinimized))
  const bounds = reactive(normalizeBounds(saved.bounds))
  const output = ref(normalizeOutput(saved.output))
  const commandHistory = ref(normalizeHistory(saved.commandHistory))
  const pendingConfirmation = ref(null)
  const isPending = ref(false)
  const refreshSignals = reactive(normalizeRefreshSignals(saved.refreshSignals))
  const sourceRect = ref(null)
  let nextOutputID = output.value.reduce((max, entry) => Math.max(max, entry.id), 0) + 1

  function open(rect = null) {
    sourceRect.value = rect ? {
      left: finiteNumber(rect.left, 0),
      top: finiteNumber(rect.top, 0),
      width: Math.max(1, finiteNumber(rect.width, 1)),
      height: Math.max(1, finiteNumber(rect.height, 1))
    } : null
    isOpen.value = true
    isMinimized.value = false
  }

  function close() {
    isOpen.value = false
    isMinimized.value = false
    sourceRect.value = null
  }

  function minimize() {
    if (isOpen.value) isMinimized.value = true
  }

  function restore() {
    if (isOpen.value) isMinimized.value = false
  }

  function setBounds(value) {
    Object.assign(bounds, normalizeBounds({ ...bounds, ...value }))
  }

  function appendOutput(text, type = 'output') {
    const entries = (Array.isArray(text) ? text : [text]).slice(-MAX_OUTPUT_LINES)
    const additions = entries.map(entry => ({
      id: nextOutputID++,
      type: ['command', 'error', 'success', 'system'].includes(type) ? type : 'output',
      text: String(entry ?? '').slice(0, MAX_OUTPUT_LENGTH)
    }))
    output.value = [...output.value, ...additions].slice(-MAX_OUTPUT_LINES)
  }

  function clearOutput() {
    output.value = []
  }

  function addCommand(command) {
    const normalized = String(command ?? '').trim().slice(0, MAX_COMMAND_LENGTH)
    if (!normalized) return
    commandHistory.value = [...commandHistory.value.filter(item => item !== normalized), normalized].slice(-MAX_HISTORY)
  }

  function setPendingConfirmation(confirmation) {
    pendingConfirmation.value = normalizeConfirmation(confirmation)
  }

  function clearPendingConfirmation() {
    pendingConfirmation.value = null
  }

  function setPending(value) {
    isPending.value = Boolean(value)
  }

  function signalRefresh(name = 'all') {
    const key = String(name).slice(0, 50) || 'all'
    if (!(key in refreshSignals) && Object.keys(refreshSignals).length >= MAX_REFRESH_SIGNALS) {
      delete refreshSignals[Object.keys(refreshSignals)[0]]
    }
    refreshSignals[key] = (refreshSignals[key] || 0) + 1
  }

  watch(
    [isOpen, isMinimized, bounds, output, commandHistory],
    () => writeSession({
      isOpen: isOpen.value,
      isMinimized: isMinimized.value,
      bounds: normalizeBounds(bounds),
      output: normalizeOutput(output.value),
      commandHistory: normalizeHistory(commandHistory.value)
    }),
    { deep: true }
  )

  return {
    isOpen,
    isMinimized,
    bounds,
    output,
    commandHistory,
    pendingConfirmation,
    isPending,
    refreshSignals,
    sourceRect,
    open,
    close,
    minimize,
    restore,
    setBounds,
    appendOutput,
    clearOutput,
    addCommand,
    setPendingConfirmation,
    clearPendingConfirmation,
    setPending,
    signalRefresh
  }
})

export { MAX_OUTPUT_LINES, MAX_HISTORY, STORAGE_KEY }
