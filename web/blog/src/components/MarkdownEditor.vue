<template>
  <div class="markdown-editor" :class="{ fullscreen: isFullscreen }">
    <div class="editor-toolbar">
      <div class="flex items-center gap-1">
        <Button variant="ghost" size="sm" @click="insertBold">
          <BoldIcon class="mr-1" /> 粗体
        </Button>
        <Button variant="ghost" size="sm" @click="insertItalic">
          <ItalicIcon class="mr-1" /> 斜体
        </Button>
        <Button variant="ghost" size="sm" @click="insertHeading">
          <Heading class="mr-1" /> 标题
        </Button>
        <Button variant="ghost" size="sm" @click="insertCode">
          <CodeIcon class="mr-1" /> 代码
        </Button>
        <Button variant="ghost" size="sm" @click="insertLink">
          <LinkIcon class="h-4 w-4 mr-1" /> 链接
        </Button>
        <Button variant="ghost" size="sm" @click="insertImage">
          <ImageIcon class="h-4 w-4 mr-1" /> 图片
        </Button>
      </div>

      <div class="ml-2">
        <Button variant="ghost" size="sm" @click="toggleFullscreen">
          <Maximize v-if="!isFullscreen" class="h-4 w-4 mr-1" />
          <X v-else class="h-4 w-4 mr-1" />
          {{ isFullscreen ? '退出' : '全屏' }}
        </Button>
      </div>
    </div>

    <div class="editor-content">
      <div class="editor-pane">
        <div class="pane-header">
          <span>Markdown 编辑</span>
          <span class="word-count">{{ wordCount }} 字</span>
        </div>
        <textarea
          ref="textareaRef"
          v-model="content"
          class="editor-textarea"
          placeholder="请输入文章内容，支持Markdown语法..."
          @input="handleInput"
          @scroll="handleEditorScroll"
        />
      </div>

      <div class="preview-pane">
        <div class="pane-header">
          <span>预览效果</span>
        </div>
        <div
          ref="previewRef"
          class="preview-content markdown-body"
          v-html="renderedHtml"
          @scroll="handlePreviewScroll"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import {
  Maximize,
  X,
  Image as ImageIcon,
  Link as LinkIcon
} from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'

const BoldIcon = { template: '<b>B</b>' }
const ItalicIcon = { template: '<i>I</i>' }
const Heading = { template: '<span>H</span>' }
const CodeIcon = { template: '<span>&lt;/&gt;</span>' }

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

const textareaRef = ref(null)
const previewRef = ref(null)
const content = ref(props.modelValue)
const isFullscreen = ref(false)
let isEditorScrolling = false
let isPreviewScrolling = false

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  highlight: function (str, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return hljs.highlight(str, { language: lang }).value
      } catch (__) {}
    }
    return ''
  }
})

const wordCount = computed(() => {
  return content.value.replace(/\s/g, '').length
})

const renderedHtml = computed(() => {
  if (!content.value) {
    return '<p class="placeholder">预览区域将显示渲染后的内容...</p>'
  }
  return md.render(content.value)
})

watch(content, (newVal) => {
  emit('update:modelValue', newVal)
})

watch(() => props.modelValue, (newVal) => {
  if (newVal !== content.value) {
    content.value = newVal
  }
})

const handleInput = () => {
  emit('update:modelValue', content.value)
}

const getSelection = () => {
  const textarea = textareaRef.value
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selectedText = content.value.substring(start, end)
  return { start, end, selectedText }
}

const insertText = (before, after = '', placeholder = '') => {
  const textarea = textareaRef.value
  const { start, end, selectedText } = getSelection()
  const text = selectedText || placeholder
  const newText = before + text + after

  content.value = content.value.substring(0, start) + newText + content.value.substring(end)

  nextTick(() => {
    textarea.focus()
    if (!selectedText && placeholder) {
      textarea.setSelectionRange(start + before.length, start + before.length + placeholder.length)
    } else {
      textarea.setSelectionRange(start + before.length + text.length + after.length, start + before.length + text.length + after.length)
    }
  })
}

const insertBold = () => insertText('**', '**', '粗体文字')
const insertItalic = () => insertText('*', '*', '斜体文字')
const insertHeading = () => insertText('\n## ', '', '标题')
const insertCode = () => insertText('`', '`', '代码')
const insertLink = () => insertText('[', '](https://example.com)', '链接文字')
const insertImage = () => insertText('![', '](https://example.com/image.jpg)', '图片描述')

const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value
  if (isFullscreen.value) {
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = ''
  }
}

const handleEditorScroll = () => {
  if (isPreviewScrolling) {
    isPreviewScrolling = false
    return
  }

  isEditorScrolling = true
  const editor = textareaRef.value
  const preview = previewRef.value

  if (editor && preview) {
    const scrollRatio = editor.scrollTop / (editor.scrollHeight - editor.clientHeight)
    preview.scrollTop = scrollRatio * (preview.scrollHeight - preview.clientHeight)
  }

  setTimeout(() => {
    isEditorScrolling = false
  }, 100)
}

const handlePreviewScroll = () => {
  if (isEditorScrolling) {
    isEditorScrolling = false
    return
  }

  isPreviewScrolling = true
  const editor = textareaRef.value
  const preview = previewRef.value

  if (editor && preview) {
    const scrollRatio = preview.scrollTop / (preview.scrollHeight - preview.clientHeight)
    editor.scrollTop = scrollRatio * (editor.scrollHeight - editor.clientHeight)
  }

  setTimeout(() => {
    isPreviewScrolling = false
  }, 100)
}
</script>

<style scoped>
.markdown-editor {
  border: 1px solid var(--theme-border);
  border-radius: var(--radius-sm);
  overflow: hidden;
  background: var(--theme-bg-primary);
}

.markdown-editor.fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 9999;
  border-radius: 0;
}

.editor-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-sm) var(--spacing-md);
  background: var(--theme-bg-primary);
  border-bottom: 1px solid var(--theme-border);
  box-shadow: var(--shadow-sm);
}

.editor-content {
  display: flex;
  height: 600px;
}

.fullscreen .editor-content {
  height: calc(100vh - 52px);
}

.editor-pane,
.preview-pane {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.editor-pane {
  background: var(--theme-bg-secondary);
}

.pane-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-sm) var(--spacing-md);
  background: var(--theme-bg-secondary);
  border-bottom: 1px solid var(--theme-border);
  font-size: var(--font-size-sm);
  color: var(--theme-text-secondary);
  font-weight: 600;
  letter-spacing: 0.3px;
}

.word-count {
  font-size: var(--font-size-xs);
  color: var(--theme-text-tertiary);
}

.editor-textarea {
  flex: 1;
  padding: var(--spacing-lg) var(--spacing-xl);
  border: none;
  outline: none;
  resize: none;
  font-family: 'SF Mono', 'Monaco', 'Inconsolata', 'Fira Code', 'Menlo', 'Consolas', 'Courier New', monospace;
  font-size: var(--font-size-base);
  line-height: var(--line-height-relaxed);
  color: var(--theme-text-primary);
  background: var(--theme-bg-secondary);
  overflow-y: auto;
  transition: background-color var(--transition-base);
}

.editor-textarea:focus {
  background: var(--theme-bg-primary);
}

.preview-content {
  flex: 1;
  padding: var(--spacing-lg) var(--spacing-xl);
  overflow-y: auto;
  background: var(--theme-bg-primary);
  border-left: 1px solid var(--theme-border);
}

.markdown-body {
  font-size: var(--font-size-lg);
  line-height: var(--line-height-relaxed);
  color: var(--theme-text-primary);
}

.markdown-body .placeholder {
  color: var(--theme-text-tertiary);
  font-style: italic;
}

.markdown-body h1,
.markdown-body h2,
.markdown-body h3,
.markdown-body h4 {
  margin-top: var(--spacing-lg);
  margin-bottom: var(--spacing-md);
  font-weight: 600;
  line-height: var(--line-height-tight);
}

.markdown-body h1 {
  font-size: 2em;
  border-bottom: 1px solid var(--theme-border-light);
  padding-bottom: 0.3em;
}

.markdown-body h2 {
  font-size: 1.5em;
  border-bottom: 1px solid var(--theme-border-light);
  padding-bottom: 0.3em;
}

.markdown-body h3 {
  font-size: 1.25em;
}

.markdown-body p {
  margin-bottom: var(--spacing-md);
}

.markdown-body code {
  padding: 0.2em 0.4em;
  margin: 0;
  font-size: 85%;
  background-color: var(--theme-bg-hover);
  border-radius: var(--radius-sm);
  font-family: 'Monaco', 'Menlo', 'Consolas', 'Courier New', monospace;
}

.markdown-body pre {
  padding: var(--spacing-md);
  overflow: auto;
  font-size: 85%;
  line-height: 1.45;
  background-color: var(--theme-bg-secondary);
  border-radius: var(--radius-md);
  margin-bottom: var(--spacing-md);
}

.markdown-body pre code {
  display: inline;
  padding: 0;
  margin: 0;
  overflow: visible;
  line-height: inherit;
  background-color: transparent;
  border: 0;
}

.markdown-body blockquote {
  padding: 0 1em;
  color: var(--theme-text-tertiary);
  border-left: 0.25em solid var(--theme-border-light);
  margin-bottom: var(--spacing-md);
}

.markdown-body ul,
.markdown-body ol {
  padding-left: 2em;
  margin-bottom: var(--spacing-md);
}

.markdown-body img {
  max-width: 100%;
  border-radius: var(--radius-sm);
  margin: var(--spacing-md) 0;
}

.markdown-body a {
  color: var(--theme-primary);
  text-decoration: none;
  transition: color var(--transition-fast);
}

.markdown-body a:hover {
  text-decoration: underline;
}

.markdown-body table {
  border-collapse: collapse;
  width: 100%;
  margin-bottom: var(--spacing-md);
}

.markdown-body table th,
.markdown-body table td {
  padding: 6px 13px;
  border: 1px solid var(--theme-border-light);
}

.markdown-body table th {
  font-weight: 600;
  background-color: var(--theme-bg-secondary);
}
</style>
