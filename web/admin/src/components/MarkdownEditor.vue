<!-- 复制博客前端的MarkdownEditor组件 -->
<template>
  <div class="markdown-editor" :class="{ fullscreen: isFullscreen }">
    <!-- 工具栏 -->
    <div class="editor-toolbar">
      <el-button-group>
        <el-button size="small" @click="insertBold">
          <b>B</b> 粗体
        </el-button>
        <el-button size="small" @click="insertItalic">
          <i>I</i> 斜体
        </el-button>
        <el-button size="small" @click="insertHeading">
          <span>H</span> 标题
        </el-button>
        <el-button size="small" @click="insertCode">
          <span>&lt;/&gt;</span> 代码
        </el-button>
        <el-button size="small" @click="insertLink">
          <el-icon><Link /></el-icon> 链接
        </el-button>
        <el-button size="small" @click="insertImage">
          <el-icon><Picture /></el-icon> 图片
        </el-button>
      </el-button-group>
      
      <el-button-group class="ml-10">
        <el-button size="small" @click="toggleFullscreen">
          <el-icon v-if="!isFullscreen"><FullScreen /></el-icon>
          <el-icon v-else><Close /></el-icon>
          {{ isFullscreen ? '退出' : '全屏' }}
        </el-button>
      </el-button-group>
    </div>

    <!-- 编辑区域 -->
    <div class="editor-content">
      <!-- 左侧编辑器 -->
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

      <!-- 右侧预览 -->
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
import { FullScreen, Close, Picture, Link } from '@element-plus/icons-vue'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'

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
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  overflow: hidden;
  background: var(--color-bg-primary);
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
  background: var(--color-bg-primary);
  border-bottom: 1px solid var(--color-border);
  box-shadow: var(--shadow-sm);
}

.ml-10 {
  margin-left: var(--spacing-sm);
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
  background: var(--color-bg-secondary);
}

.pane-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-sm) var(--spacing-md);
  background: var(--color-bg-hover);
  border-bottom: 1px solid var(--color-border);
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  font-weight: 600;
  letter-spacing: 0.3px;
}

.word-count {
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
}

.editor-textarea {
  flex: 1;
  padding: var(--spacing-lg);
  border: none;
  outline: none;
  resize: none;
  font-family: var(--font-mono);
  font-size: var(--font-size-base);
  line-height: var(--line-height-relaxed);
  color: var(--color-text-primary);
  background: var(--color-bg-secondary);
  overflow-y: auto;
  transition: background var(--transition-fast);
}

.editor-textarea:focus {
  background: var(--color-bg-primary);
}

.preview-content {
  flex: 1;
  padding: var(--spacing-lg) var(--spacing-xl);
  overflow-y: auto;
  background: var(--color-bg-primary);
  border-left: 1px solid var(--color-border);
}

.markdown-body {
  font-size: var(--font-size-base);
  line-height: var(--line-height-relaxed);
  color: var(--color-text-primary);
}

.markdown-body .placeholder {
  color: var(--color-text-tertiary);
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
  border-bottom: 1px solid var(--color-border-lighter);
  padding-bottom: 0.3em;
}

.markdown-body h2 {
  font-size: 1.5em;
  border-bottom: 1px solid var(--color-border-lighter);
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
  background-color: var(--color-bg-hover);
  border-radius: var(--radius-sm);
  font-family: var(--font-mono);
}

.markdown-body pre {
  padding: var(--spacing-md);
  overflow: auto;
  font-size: 85%;
  line-height: 1.45;
  background-color: var(--color-bg-secondary);
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
  color: var(--color-text-tertiary);
  border-left: 0.25em solid var(--color-border-lighter);
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
  color: var(--color-primary);
  text-decoration: none;
  transition: color var(--transition-fast);
}

.markdown-body a:hover {
  color: var(--color-primary-hover);
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
  border: 1px solid var(--color-border-lighter);
}

.markdown-body table th {
  font-weight: 600;
  background-color: var(--color-bg-secondary);
}
</style>

