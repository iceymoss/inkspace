<template>
  <div class="terminal-hero">
    <div class="terminal-hero-copy">
      <span class="terminal-status">{{ settings.status }}</span>
      <h1>
        {{ title.before }}<em v-if="title.accent">{{ title.accent }}</em>{{ title.after }}
      </h1>
      <p>{{ settings.description }}</p>
      <div class="terminal-actions">
        <a
          class="terminal-primary"
          :href="safeLink(settings.primary_link)"
          @click="$emit('navigate', $event, settings.primary_link)"
        >{{ settings.primary_text }}</a>
        <a
          class="terminal-secondary"
          :href="safeLink(settings.secondary_link)"
          @click="$emit('navigate', $event, settings.secondary_link)"
        >{{ settings.secondary_text }}</a>
      </div>
    </div>

    <div
      v-if="active"
      class="terminal-window terminal-placeholder"
      aria-live="polite"
    >
      <span>terminal session active</span>
      <small>Use the floating window to continue.</small>
    </div>
    <div
      v-else
      ref="terminalWindow"
      class="terminal-window"
      role="button"
      tabindex="0"
      aria-label="Open interactive terminal"
      @click="activateTerminal"
      @keydown.enter.prevent="activateTerminal"
      @keydown.space.prevent="activateTerminal"
    >
      <div class="terminal-bar">
        <i class="dot red" /><i class="dot amber" /><i class="dot green" />
        <span>visitor@inkspace.log: ~</span>
      </div>
      <div class="terminal-body">
        <div><b>❯</b> <strong>whoami</strong></div>
        <p>{{ settings.eyebrow }}</p>
        <div><b>❯</b> <strong>cat content.json</strong></div>
        <p>{ articles: <i>{{ stats.articleCount }}</i>, works: <i>{{ stats.workCount }}</i>, categories: <i>{{ stats.categoryCount }}</i> }</p>
        <div><b>❯</b> <strong>{{ typedCommand }}</strong><span class="terminal-caret" /></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onBeforeUnmount, onMounted, ref } from 'vue'

defineProps({
  settings: { type: Object, required: true },
  title: { type: Object, required: true },
  stats: { type: Object, required: true },
  active: { type: Boolean, default: false }
})

const emit = defineEmits(['navigate', 'activate'])

const terminalWindow = ref(null)
const commands = ['ls photos/', 'git log --all', 'open wiki/', 'npm run build']
const typedCommand = ref(commands[0])
let timer
let commandIndex = 0
let position = commands[0].length
let deleting = true
let reducedMotion

const safeLink = link => /^https?:\/\//.test(link) || link?.startsWith('/') ? link : '#'

function activateTerminal() {
  emit('activate', terminalWindow.value?.getBoundingClientRect())
}

function tick() {
  const command = commands[commandIndex]
  typedCommand.value = command.slice(0, position)
  if (!deleting && position < command.length) {
    position += 1
    timer = window.setTimeout(tick, 90)
    return
  }
  if (!deleting) {
    deleting = true
    timer = window.setTimeout(tick, 1500)
    return
  }
  if (position > 0) {
    position -= 1
    timer = window.setTimeout(tick, 35)
    return
  }
  deleting = false
  commandIndex = (commandIndex + 1) % commands.length
  timer = window.setTimeout(tick, 450)
}

onMounted(() => {
  reducedMotion = window.matchMedia?.('(prefers-reduced-motion: reduce)')
  reducedMotion?.addEventListener?.('change', handleReducedMotionChange)
  if (reducedMotion?.matches) return
  timer = window.setTimeout(tick, 1500)
})

function handleReducedMotionChange(event) {
  window.clearTimeout(timer)
  typedCommand.value = commands[0]
  if (!event.matches) timer = window.setTimeout(tick, 1500)
}

onBeforeUnmount(() => {
  window.clearTimeout(timer)
  reducedMotion?.removeEventListener?.('change', handleReducedMotionChange)
})
</script>

<style scoped>
.terminal-hero { display: grid; grid-template-columns: 1.05fr .95fr; gap: 56px; align-items: center; padding: 60px 0 66px; }
.terminal-status { display: inline-flex; align-items: center; gap: 8px; margin-bottom: 25px; padding: 4px 14px; border: 1px solid color-mix(in srgb, var(--green) 35%, transparent); border-radius: 99px; background: color-mix(in srgb, var(--green) 12%, transparent); color: var(--green); font: 12px var(--terminal-mono); }
.terminal-status::before { width: 7px; height: 7px; border-radius: 50%; background: var(--green); box-shadow: 0 0 8px var(--green); content: ''; }
h1 { margin: 0; color: var(--bright); font-size: clamp(34px, 5vw, 52px); line-height: 1.3; }
h1 em { color: var(--accent); font-style: normal; }
.terminal-hero-copy > p { max-width: 35em; margin: 20px 0 0; color: var(--sub); font-size: 15px; line-height: 1.8; }
.terminal-actions { display: flex; gap: 14px; margin-top: 34px; font: 14px var(--terminal-mono); }
.terminal-actions a { padding: 10px 24px; border: 1px solid var(--line); border-radius: 10px; color: var(--ink); text-decoration: none; }
.terminal-actions .terminal-primary { border-color: var(--accent); background: var(--accent); color: var(--terminal-button-ink); font-weight: 600; }
.terminal-actions a:hover { border-color: var(--accent); color: var(--accent); transform: translateY(-1px); }
.terminal-actions .terminal-primary:hover { color: var(--terminal-button-ink); filter: brightness(1.1); }
.terminal-window { overflow: hidden; border: 1px solid var(--line); border-radius: 14px; background: var(--panel); box-shadow: var(--terminal-shadow); font-family: var(--terminal-mono); font-size: 13px; }
.terminal-window[role="button"] { cursor: pointer; transition: border-color .2s, transform .2s; }
.terminal-window[role="button"]:hover { border-color: var(--accent); transform: translateY(-2px); }
.terminal-window[role="button"]:focus-visible { outline: 2px solid var(--accent); outline-offset: 4px; }
.terminal-placeholder { display: grid; min-height: 190px; align-content: center; gap: 7px; border-style: dashed; color: var(--green); text-align: center; }
.terminal-placeholder::before { color: var(--accent); font-size: 24px; content: '❯_'; }
.terminal-placeholder small { color: var(--sub); }
.terminal-bar { display: flex; align-items: center; gap: 7px; padding: 12px 16px; border-bottom: 1px solid var(--line); }
.terminal-bar > span { margin: 0 auto; color: var(--sub); font-size: 12px; }
.dot { width: 11px; height: 11px; border-radius: 50%; }.dot.red { background: #f26d6d; }.dot.amber { background: #f2c46d; }.dot.green { background: #6dd087; }
.terminal-body { padding: 20px 22px 25px; color: var(--bright); line-height: 2.05; overflow-wrap: anywhere; }
.terminal-body div b { color: var(--green); }.terminal-body p { margin: 0; color: var(--sub); }.terminal-body p i { color: var(--amber); font-style: normal; }
.terminal-caret { display: inline-block; width: 8px; height: 15px; margin-left: 3px; background: var(--accent); vertical-align: -2px; animation: terminal-blink 1.1s steps(1) infinite; }
@keyframes terminal-blink { 50% { opacity: 0; } }
@media (max-width: 900px) { .terminal-hero { grid-template-columns: 1fr; gap: 40px; padding-top: 48px; } }
@media (max-width: 560px) { .terminal-hero { padding: 34px 0 46px; }.terminal-actions { flex-wrap: wrap; }.terminal-actions a { flex: 1; padding-inline: 14px; text-align: center; }.terminal-body { padding-inline: 16px; } }
@media (prefers-reduced-motion: reduce) { .terminal-caret { animation: none; }.terminal-window[role="button"] { transition: none; }.terminal-window[role="button"]:hover { transform: none; } }
</style>
