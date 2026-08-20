<template>
  <!-- `html` must only come from an API response that applies the backend sanitizer. -->
  <div
    v-if="html"
    class="markdown-preview vditor-reset"
    :data-markdown-theme="markdownTheme"
    :data-editorial="editorial"
    v-html="html"
  />
  <div
    v-else
    ref="previewRef"
    class="markdown-preview vditor-reset"
    :data-markdown-theme="markdownTheme"
    :data-editorial="editorial"
  />
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import Vditor from 'vditor'
import 'vditor/dist/index.css'
import { useAppearanceStore } from '@/stores/appearance'
import { loadCodeTheme, loadHighlightTheme } from '@/utils/codeTheme'

const props = defineProps({
  html: {
    type: String,
    default: ''
  },
  content: {
    type: String,
    default: ''
  },
  editorial: {
    type: Boolean,
    default: false
  }
})

const previewRef = ref(null)
const appearanceStore = useAppearanceStore()
const markdownTheme = computed(() => appearanceStore.resolvedColorScheme)
let mounted = false
let renderVersion = 0
let codeTheme = 'github'

async function renderMarkdown() {
  const version = ++renderVersion
  if (!mounted || props.html) return

  await nextTick()
  const element = previewRef.value
  if (!element || version !== renderVersion) return
  element.replaceChildren()
  if (!props.content) return

  await Vditor.preview(element, props.content, {
    mode: markdownTheme.value,
    theme: {
      current: markdownTheme.value
    },
    markdown: {
      toc: true,
      mark: true,
      footnotes: true,
      autoSpace: true,
      sanitize: true
    },
    hljs: {
      lineNumber: false,
      style: codeTheme
    },
    speech: {
      enable: false
    },
    anchor: 0
  })
  if (!mounted || version !== renderVersion || previewRef.value !== element) element.replaceChildren()
}

watch([() => props.html, () => props.content, markdownTheme], renderMarkdown)

onMounted(async () => {
  mounted = true
  codeTheme = await loadCodeTheme()
  await loadHighlightTheme(codeTheme)
  await renderMarkdown()
})

onBeforeUnmount(() => {
  mounted = false
  renderVersion += 1
  previewRef.value?.replaceChildren()
})
</script>

<style scoped>
.markdown-preview {
  max-width: 100%;
  padding: 0;
  background: transparent;
  color: var(--theme-text-primary);
  font-size: 16px;
  line-height: 1.9;
  overflow-wrap: anywhere;
}

.markdown-preview[data-editorial='true'] {
  font-family: Georgia, 'Noto Serif SC', 'Songti SC', serif;
  font-size: clamp(16px, 2vw, 18px);
  line-height: 1.95;
}

.markdown-preview :deep(h1),
.markdown-preview :deep(h2),
.markdown-preview :deep(h3),
.markdown-preview :deep(h4) {
  margin: 1.7em 0 .7em;
  line-height: 1.35;
}

.markdown-preview :deep(h2) {
  padding-bottom: .35em;
  border-bottom: 1px solid var(--theme-border-light);
}

.markdown-preview :deep(p) {
  margin: 0 0 1.35em;
}

.markdown-preview :deep(ul),
.markdown-preview :deep(ol) {
  margin: 0 0 1.35em;
  padding-left: 1.6em;
}

.markdown-preview :deep(li) {
  margin: .4em 0;
}

.markdown-preview :deep(a) {
  color: var(--theme-primary);
  text-decoration-thickness: 1px;
  text-underline-offset: 3px;
}

.markdown-preview :deep(a:hover) {
  color: var(--theme-primary-hover);
}

.markdown-preview :deep(a:focus-visible) {
  outline: 2px solid var(--theme-primary);
  outline-offset: 3px;
}

.markdown-preview :deep(img) {
  display: block;
  max-width: 100%;
  height: auto;
  margin: 1.8em auto;
}

.markdown-preview :deep(blockquote) {
  margin: 1.8em 0;
  padding: 2px 0 2px 18px;
  border-left: 3px solid var(--theme-primary);
  color: var(--theme-text-secondary);
}

.markdown-preview :deep(pre) {
  max-width: 100%;
  margin: 1.8em 0;
  padding: 18px 20px;
  overflow-x: auto;
  border: 1px solid var(--theme-border-light);
  border-radius: 6px;
  background: var(--theme-bg-secondary);
  font-size: 13px;
  line-height: 1.7;
}

.markdown-preview :deep(code) {
  font-family: 'SFMono-Regular', Consolas, monospace;
}

.markdown-preview :deep(:not(pre) > code) {
  padding: 2px 5px;
  font-size: .9em;
}

.markdown-preview :deep(table) {
  display: block;
  width: max-content;
  max-width: 100%;
  margin: 1.8em 0;
  overflow-x: auto;
  border-collapse: collapse;
}

.markdown-preview :deep(th),
.markdown-preview :deep(td) {
  padding: 9px 12px;
  border: 1px solid var(--theme-border);
}

.markdown-preview :deep(th) {
  text-align: left;
}

.markdown-preview :deep(hr) {
  margin: 2.5em 0;
  border: 0;
  border-top: 1px solid var(--theme-border);
}

.markdown-preview :deep(input[type='checkbox']) {
  accent-color: var(--theme-primary);
}

@media (max-width: 600px) {
  .markdown-preview {
    font-size: 15px;
    line-height: 1.8;
  }

  .markdown-preview :deep(pre),
  .markdown-preview :deep(table) {
    -webkit-overflow-scrolling: touch;
  }
}
</style>
