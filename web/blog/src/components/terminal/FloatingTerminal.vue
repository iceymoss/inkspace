<template>
  <Teleport to="body">
    <section
      v-if="terminal.isOpen"
      ref="terminalElement"
      class="floating-terminal"
      :class="{ 'is-minimized': terminal.isMinimized, 'is-entering': enterTransform }"
      :style="terminalStyle"
      role="dialog"
      aria-label="InkSpace terminal"
    >
      <header class="terminal-titlebar" @pointerdown="startDrag">
        <span class="terminal-lights" aria-hidden="true"><i /><i /><i /></span>
        <strong>visitor@inkspace:{{ currentVirtualPath }}</strong>
        <span class="terminal-controls">
          <button
            type="button"
            :aria-label="terminal.isMinimized ? 'Restore terminal' : 'Minimize terminal'"
            @pointerdown.stop
            @click="terminal.isMinimized ? restore() : terminal.minimize()"
          >{{ terminal.isMinimized ? '□' : '−' }}</button>
          <button type="button" aria-label="Close terminal" @pointerdown.stop @click="terminal.close()">×</button>
        </span>
      </header>

      <template v-if="!terminal.isMinimized">
        <div ref="outputElement" class="terminal-output" aria-live="polite" aria-atomic="false">
          <p v-if="!terminal.output.length" class="terminal-welcome">InkSpace shell ready. Type <b>help</b> to begin.</p>
          <p v-for="entry in terminal.output" :key="entry.id" :class="`line-${entry.type}`">
            <span v-if="entry.type === 'command'" aria-hidden="true">❯ </span>{{ entry.text }}
          </p>
          <p v-if="terminal.pendingConfirmation" class="terminal-confirmation">
            {{ terminal.pendingConfirmation.message }} <b>[yes/no]</b>
          </p>
        </div>

        <form class="terminal-input-row" @submit.prevent="submitCommand">
          <label for="global-terminal-input" aria-label="Terminal command">❯</label>
          <input
            id="global-terminal-input"
            ref="inputElement"
            v-model="command"
            type="text"
            autocomplete="off"
            autocapitalize="off"
            spellcheck="false"
            maxlength="500"
            :readonly="terminal.isPending"
            :aria-busy="terminal.isPending"
            aria-describedby="terminal-shortcuts"
            @keydown="handleInputKeydown"
          >
          <span id="terminal-shortcuts" class="visually-hidden">Use up and down arrows for command history. Control L clears output.</span>
        </form>
        <button
          type="button"
          class="terminal-resize"
          aria-label="Resize terminal"
          @pointerdown.stop="startResize"
        />
      </template>
    </section>
  </Teleport>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useTerminalStore } from '@/stores/terminal'
import { parseCommand, tokenizeCommand } from '@/utils/terminal/parser'
import { routeToVirtualPath } from '@/utils/terminal/virtualFileSystem'

const emit = defineEmits(['submit', 'confirm'])
const terminal = useTerminalStore()
const route = useRoute()
const terminalElement = ref(null)
const outputElement = ref(null)
const inputElement = ref(null)
const command = ref('')
const historyIndex = ref(-1)
const historyDraft = ref('')
const enterTransform = ref('')
let interaction = null
let animationFrame = 0
const commandNames = ['help', 'status', 'pwd', 'ls', 'cd', 'cat', 'grep', 'open', 'search', 'list', 'filter', 'sort', 'page', 'reset', 'theme', 'scheme', 'scroll', 'focus', 'history', 'clear', 'minimize', 'close', 'like', 'unlike', 'favorite', 'unfavorite', 'follow', 'unfollow']
const currentVirtualPath = computed(() => routeToVirtualPath(route) || '/inkspace/index')
const staticCompletions = {
  open: ['home', 'index', 'blog', 'works', 'photos', 'wiki', 'users', 'about', 'links', 'login', 'dashboard', 'article', 'work', 'user', 'workspace', 'doc'],
  search: ['blog', 'works', 'photos', 'users', 'wiki'],
  list: ['articles', 'works', 'photos', 'users', 'wiki'],
  theme: ['magazine', 'terminal'],
  scheme: ['system', 'light', 'dark'],
  scroll: ['top', 'bottom'],
  focus: ['search'],
  reset: ['filters'],
  like: ['article', 'work'],
  unlike: ['article', 'work'],
  favorite: ['article', 'work'],
  unfavorite: ['article', 'work'],
  follow: ['user'],
  unfollow: ['user']
}

const terminalStyle = computed(() => {
  const bounds = terminal.bounds
  const style = terminal.isMinimized
    ? { left: `${bounds.x}px`, top: `${bounds.y}px`, width: '260px' }
    : { left: `${bounds.x}px`, top: `${bounds.y}px`, width: `${bounds.width}px`, height: `${bounds.height}px` }
  if (enterTransform.value) style.transform = enterTransform.value
  return style
})

function viewportSize() {
  return { width: window.innerWidth, height: window.innerHeight }
}

function clampBounds(candidate = terminal.bounds) {
  const viewport = viewportSize()
  if (viewport.width <= 560) return
  const width = Math.min(Math.max(420, candidate.width), Math.max(420, viewport.width - 16))
  const height = Math.min(Math.max(280, candidate.height), Math.max(280, viewport.height - 16))
  terminal.setBounds({
    width,
    height,
    x: Math.min(Math.max(8, candidate.x), Math.max(8, viewport.width - width - 8)),
    y: Math.min(Math.max(8, candidate.y), Math.max(8, viewport.height - height - 8))
  })
}

function focusInput() {
  if (!terminal.isMinimized) nextTick(() => inputElement.value?.focus())
}

function animateFromSource() {
  const source = terminal.sourceRect
  if (!source || viewportSize().width <= 560) return
  const bounds = terminal.bounds
  enterTransform.value = `translate(${source.left - bounds.x}px, ${source.top - bounds.y}px) scale(${source.width / bounds.width}, ${source.height / bounds.height})`
  animationFrame = window.requestAnimationFrame(() => {
    enterTransform.value = ''
    terminal.sourceRect = null
  })
}

function restore() {
  terminal.restore()
  clampBounds()
  focusInput()
}

function startDrag(event) {
  if (terminal.isMinimized || viewportSize().width <= 560 || event.button !== 0) return
  event.preventDefault()
  interaction = { type: 'drag', startX: event.clientX, startY: event.clientY, bounds: { ...terminal.bounds } }
  event.currentTarget.setPointerCapture?.(event.pointerId)
}

function startResize(event) {
  if (viewportSize().width <= 560 || event.button !== 0) return
  event.preventDefault()
  interaction = { type: 'resize', startX: event.clientX, startY: event.clientY, bounds: { ...terminal.bounds } }
  event.currentTarget.setPointerCapture?.(event.pointerId)
}

function handlePointerMove(event) {
  if (!interaction) return
  const deltaX = event.clientX - interaction.startX
  const deltaY = event.clientY - interaction.startY
  if (interaction.type === 'drag') {
    clampBounds({ ...interaction.bounds, x: interaction.bounds.x + deltaX, y: interaction.bounds.y + deltaY })
  } else {
    clampBounds({ ...interaction.bounds, width: interaction.bounds.width + deltaX, height: interaction.bounds.height + deltaY })
  }
}

function stopInteraction() {
  interaction = null
}

function showHistory(direction) {
  const history = terminal.commandHistory
  if (!history.length) return
  if (historyIndex.value === -1) historyDraft.value = command.value
  historyIndex.value = Math.min(history.length - 1, Math.max(-1, historyIndex.value + direction))
  command.value = historyIndex.value === -1 ? historyDraft.value : history[history.length - 1 - historyIndex.value]
}

function commonPrefix(values) {
  if (!values.length) return ''
  return values.reduce((prefix, value) => {
    let index = 0
    while (index < prefix.length && index < value.length && prefix[index] === value[index]) index += 1
    return prefix.slice(0, index)
  })
}

function listedResourceNames() {
  return terminal.output
    .flatMap(entry => String(entry.text).split(/\s+/))
    .filter(value => /^(?:[1-9]\d*)-[^/\s]+(?:\.md|\.json|\.jpg|\/)$/.test(value))
    .slice(-50)
}

function virtualPathCandidates() {
  if (currentVirtualPath.value === '/inkspace') {
    return ['index/', 'articles/', 'workspaces/', 'my-works/', 'favorites/', 'notifications/', 'profile/', 'appearance/']
  }
  if (currentVirtualPath.value === '/inkspace/index') {
    return ['../', '../blog/', '../works/', '../photos/', '../wiki/', '../users/', '../about/', '../links/']
  }
  return ['../', './', ...listedResourceNames()]
}

function completionCandidates(tokens, argumentIndex) {
  const name = tokens[0]?.toLowerCase()
  if (argumentIndex === 0) return commandNames
  if (['ls', 'cd', 'cat'].includes(name) && argumentIndex === 1) return virtualPathCandidates()
  if (name === 'grep' && argumentIndex === 2) return virtualPathCandidates()
  if (name === 'filter') {
    if (argumentIndex === 1) return ['blog', 'works']
    if (argumentIndex === 2) return tokens[1] === 'blog' ? ['category', 'tag', 'rank'] : ['type']
    if (argumentIndex === 3 && tokens[2] === 'rank') return ['hot', 'week', 'month', 'year']
    if (argumentIndex === 3 && tokens[2] === 'type') return ['project', 'photography']
  }
  if (name === 'sort') {
    if (argumentIndex === 1) return ['blog', 'works', 'photos']
    if (argumentIndex === 2) return tokens[1] === 'blog'
      ? ['time', 'hot', 'view_count', 'like_count', 'comment_count']
      : ['latest', 'popular', 'view', 'likes']
  }
  return argumentIndex === 1 ? staticCompletions[name] || [] : []
}

function completeCommand(event) {
  event.preventDefault()
  const trailingSpace = /\s$/.test(command.value)
  let tokens
  try {
    tokens = tokenizeCommand(command.value)
  } catch {
    return
  }
  const argumentIndex = trailingSpace ? tokens.length : Math.max(0, tokens.length - 1)
  const partial = trailingSpace ? '' : tokens[argumentIndex] || ''
  const matches = completionCandidates(tokens, argumentIndex).filter(value => value.startsWith(partial))
  if (!matches.length) return

  const completion = matches.length === 1 ? matches[0] : commonPrefix(matches)
  if (completion && completion !== partial) {
    const completed = [...tokens]
    if (trailingSpace) completed.push(completion)
    else completed[argumentIndex] = completion
    command.value = `${completed.join(' ')}${matches.length === 1 ? ' ' : ''}`
  } else if (matches.length > 1) {
    terminal.appendOutput(matches.join('  '), 'system')
  }
  nextTick(() => inputElement.value?.setSelectionRange(command.value.length, command.value.length))
}

function handleInputKeydown(event) {
  if (event.key === 'ArrowUp') {
    event.preventDefault()
    showHistory(1)
  } else if (event.key === 'ArrowDown') {
    event.preventDefault()
    showHistory(-1)
  } else if (event.key.toLowerCase() === 'l' && event.ctrlKey) {
    event.preventDefault()
    terminal.clearOutput()
  } else if (event.key === 'Escape') {
    if (terminal.pendingConfirmation) {
      const confirmation = terminal.pendingConfirmation
      terminal.clearPendingConfirmation()
      emit('confirm', { ...confirmation, accepted: false })
    } else terminal.minimize()
  } else if (event.key === 'Tab') {
    completeCommand(event)
  }
}

function submitCommand() {
  const value = command.value.trim()
  if (!value || terminal.isPending) return
  command.value = ''
  historyIndex.value = -1
  historyDraft.value = ''
  terminal.addCommand(value)
  terminal.appendOutput(value, 'command')

  if (terminal.pendingConfirmation) {
    const confirmation = terminal.pendingConfirmation
    terminal.clearPendingConfirmation()
    if (/^yes$/i.test(value)) {
      emit('confirm', { ...confirmation, accepted: true })
      return
    }
    emit('confirm', { ...confirmation, accepted: false })
    if (/^no$/i.test(value)) return
  }

  try {
    const parsed = parseCommand(value)
    if (parsed.name === 'clear') terminal.clearOutput()
    else if (parsed.name === 'help') emit('submit', value)
    else if (parsed.name === 'history') {
      terminal.appendOutput(terminal.commandHistory.map((item, index) => `${index + 1}  ${item}`), 'system')
    } else emit('submit', value)
  } catch (error) {
    terminal.appendOutput(error.message, 'error')
  }
  focusInput()
}

watch(() => [terminal.isOpen, terminal.isMinimized], ([open, minimized], previous = []) => {
  if (open && !minimized) {
    clampBounds()
    focusInput()
    if (!previous[0]) nextTick(animateFromSource)
  }
})

watch(() => terminal.output.length, () => nextTick(() => {
  if (outputElement.value) outputElement.value.scrollTop = outputElement.value.scrollHeight
}))

watch(() => terminal.isPending, pending => {
  if (!pending && terminal.isOpen && !terminal.isMinimized) focusInput()
})

onMounted(() => {
  window.addEventListener('pointermove', handlePointerMove)
  window.addEventListener('pointerup', stopInteraction)
  window.addEventListener('pointercancel', stopInteraction)
  window.addEventListener('resize', clampBounds)
  if (terminal.isOpen) {
    clampBounds()
    focusInput()
  }
})

onBeforeUnmount(() => {
  window.cancelAnimationFrame(animationFrame)
  window.removeEventListener('pointermove', handlePointerMove)
  window.removeEventListener('pointerup', stopInteraction)
  window.removeEventListener('pointercancel', stopInteraction)
  window.removeEventListener('resize', clampBounds)
})
</script>

<style scoped>
.floating-terminal {
  --ft-bg: var(--panel, var(--theme-bg-card, #101622));
  --ft-bg-soft: var(--panel-2, var(--theme-bg-secondary, #182131));
  --ft-text: var(--ink, var(--theme-text-primary, #d8e2f0));
  --ft-bright: var(--bright, var(--theme-text-primary, #f4f7fb));
  --ft-muted: var(--sub, var(--theme-text-secondary, #8794a8));
  --ft-line: var(--line, var(--theme-border, #344156));
  --ft-accent: var(--accent, var(--theme-primary, #4b9eff));
  --ft-green: var(--green, #57c985);
  position: fixed;
  z-index: 1800;
  display: flex;
  overflow: hidden;
  flex-direction: column;
  min-width: 420px;
  min-height: 280px;
  border: 1px solid var(--ft-line);
  border-radius: 12px;
  background: var(--ft-bg);
  box-shadow: var(--terminal-shadow, 0 22px 70px rgba(0, 0, 0, .32));
  color: var(--ft-text);
  font: 13px/1.65 var(--terminal-mono, 'SFMono-Regular', Consolas, monospace);
  transform-origin: top left;
  transition: transform .24s ease-out, opacity .24s ease-out;
}
.floating-terminal.is-entering { opacity: .35; }
.floating-terminal.is-minimized { min-width: 0; min-height: 0; height: auto; }
.terminal-titlebar { display: flex; flex: none; align-items: center; gap: 12px; height: 42px; padding: 0 10px 0 13px; border-bottom: 1px solid var(--ft-line); background: var(--ft-bg-soft); cursor: move; touch-action: none; user-select: none; }
.is-minimized .terminal-titlebar { border-bottom: 0; cursor: default; }
.terminal-titlebar strong { overflow: hidden; flex: 1; color: var(--ft-muted); font-size: 11px; font-weight: 500; text-align: center; text-overflow: ellipsis; white-space: nowrap; }
.terminal-lights { display: flex; gap: 6px; }
.terminal-lights i { width: 9px; height: 9px; border-radius: 50%; background: #ef6b6b; }
.terminal-lights i:nth-child(2) { background: #e9b84f; }
.terminal-lights i:nth-child(3) { background: #58bd75; }
.terminal-controls { display: flex; gap: 3px; }
.terminal-controls button { display: grid; width: 26px; height: 26px; padding: 0; border: 0; border-radius: 5px; background: transparent; color: var(--ft-muted); font: 18px/1 inherit; cursor: pointer; place-items: center; }
.terminal-controls button:hover, .terminal-controls button:focus-visible { background: color-mix(in srgb, var(--ft-accent) 15%, transparent); color: var(--ft-bright); outline: none; }
.terminal-output { overflow: auto; flex: 1; padding: 17px 19px 8px; overscroll-behavior: contain; scrollbar-color: var(--ft-line) transparent; }
.terminal-output p { margin: 0 0 5px; color: var(--ft-text); white-space: pre-wrap; overflow-wrap: anywhere; }
.terminal-output .terminal-welcome, .terminal-output .line-system { color: var(--ft-muted); }
.terminal-output .line-command { color: var(--ft-bright); }
.terminal-output .line-command span { color: var(--ft-green); }
.terminal-output .line-error { color: #ee7777; }
.terminal-output .line-success { color: var(--ft-green); }
.terminal-confirmation { padding: 8px 10px; border-left: 2px solid var(--ft-accent); background: color-mix(in srgb, var(--ft-accent) 9%, transparent); }
.terminal-input-row { display: flex; flex: none; align-items: center; gap: 9px; min-height: 46px; padding: 7px 19px 10px; border-top: 1px solid color-mix(in srgb, var(--ft-line) 65%, transparent); }
.terminal-input-row label { color: var(--ft-green); font-weight: 700; }
.terminal-input-row input { min-width: 0; flex: 1; border: 0; outline: 0; background: transparent; color: var(--ft-bright); caret-color: var(--ft-accent); font: inherit; }
.terminal-resize { position: absolute; right: 0; bottom: 0; width: 22px; height: 22px; padding: 0; border: 0; background: linear-gradient(135deg, transparent 50%, var(--ft-line) 51% 58%, transparent 59% 67%, var(--ft-line) 68% 75%, transparent 76%); cursor: nwse-resize; touch-action: none; }
.visually-hidden { position: absolute; overflow: hidden; width: 1px; height: 1px; padding: 0; border: 0; margin: -1px; clip: rect(0, 0, 0, 0); white-space: nowrap; }
@media (max-width: 560px) {
  .floating-terminal, .floating-terminal.is-minimized { inset: auto 0 0 !important; width: 100% !important; height: min(86dvh, 680px) !important; min-width: 0; min-height: 220px; border-width: 1px 0 0; border-radius: 16px 16px 0 0; transform: none !important; }
  .floating-terminal.is-minimized { height: 44px !important; min-height: 0; }
  .terminal-titlebar { height: 44px; cursor: default; }
  .terminal-output { padding-inline: 15px; }
  .terminal-input-row { padding-inline: 15px; padding-bottom: max(10px, env(safe-area-inset-bottom)); }
  .terminal-resize { display: none; }
}
@media (prefers-reduced-motion: reduce) {
  .floating-terminal { transition: none; }
}
</style>
