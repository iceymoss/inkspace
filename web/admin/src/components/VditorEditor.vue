<template>
  <div ref="vditorRef" class="vditor-container"></div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import Vditor from 'vditor'
import 'vditor/dist/index.css'
import { loadCodeTheme, loadHighlightTheme, getMarkdownTheme } from '@/utils/codeTheme'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  height: {
    type: String,
    default: '600px'
  }
})

const emit = defineEmits(['update:modelValue'])

const vditorRef = ref(null)
let vditor = null

onMounted(async () => {
  // 等待 DOM 准备好
  await nextTick()
  
  // 先同步加载主题配置，确保在初始化 Vditor 前配置已准备好
  const codeThemeValue = await loadCodeTheme()
  // 确保 highlight.js 主题 CSS 完全加载完成
  await loadHighlightTheme(codeThemeValue)
  
  const mdTheme = await getMarkdownTheme()
  
  // 加载 Vditor 内容主题 CSS（如果需要）
  const loadVditorContentTheme = (theme) => {
    return new Promise((resolve) => {
      const themeId = `vditor-content-theme-${theme}`
      if (document.getElementById(themeId)) {
        resolve()
        return
      }
      
      const link = document.createElement('link')
      link.id = themeId
      link.rel = 'stylesheet'
      link.href = `https://unpkg.com/vditor@3.10.4/dist/css/content-theme/${theme}.css`
      link.onload = () => resolve()
      link.onerror = () => resolve() // 即使失败也继续
      document.head.appendChild(link)
    })
  }
  
  // 加载 Vditor 内容主题
  await loadVditorContentTheme(mdTheme || 'light')
  // 额外等待一小段时间，确保 CSS 样式已应用
  await new Promise(resolve => setTimeout(resolve, 100))
  
  // 使用加载的配置初始化 Vditor
  vditor = new Vditor(vditorRef.value, {
    height: props.height,
    mode: 'sv', // 分屏预览模式
    placeholder: '请输入文章内容，支持 Markdown 语法...',
    theme: 'classic',
    icon: 'material',
    typewriterMode: false,
    toolbarConfig: {
      pin: true,
    },
    cache: {
      enable: false,
    },
    counter: {
      enable: true,
      type: 'markdown',
    },
    preview: {
      delay: 500,
      hljs: {
        style: codeThemeValue || 'github', // 使用配置的代码主题
        lineNumber: true,
      },
      markdown: {
        toc: true,
        mark: true,
        footnotes: true,
        autoSpace: true,
      },
    },
    upload: {
      url: '/api/upload/image',
      max: 5 * 1024 * 1024, // 5MB
      accept: 'image/*',
      fieldName: 'file',
      headers: {
        Authorization: `Bearer ${localStorage.getItem('admin_token')}`
      },
      format(files, responseText) {
        const response = JSON.parse(responseText)
        if (response.code === 0) {
          return JSON.stringify({
            msg: '',
            code: 0,
            data: {
              errFiles: [],
              succMap: {
                [files[0].name]: response.data.url
              }
            }
          })
        }
        return JSON.stringify({
          msg: response.message || '上传失败',
          code: 1,
          data: {
            errFiles: [files[0].name],
            succMap: {}
          }
        })
      }
    },
    after: () => {
      vditor.setValue(props.modelValue || '')
    },
    input: (value) => {
      emit('update:modelValue', value)
    },
  })
})

onBeforeUnmount(() => {
  if (vditor) {
    vditor.destroy()
    vditor = null
  }
})

watch(() => props.modelValue, (newVal) => {
  if (vditor) {
    const currentValue = vditor.getValue()
    // 只有当值真正改变时才更新，避免循环更新
    if (newVal !== currentValue) {
      vditor.setValue(newVal || '')
    }
  }
}, { immediate: false })

defineExpose({
  getValue: () => vditor?.getValue(),
  setValue: (value) => vditor?.setValue(value),
  focus: () => vditor?.focus(),
})
</script>

<style scoped>
.vditor-container {
  border-radius: var(--radius-md);
  width: 100%;
  min-width: 0;
}

:deep(.vditor) {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  width: 100%;
}

:deep(.vditor-toolbar) {
  background-color: var(--color-bg-primary);
  border-bottom: 1px solid var(--color-border);
  padding: var(--spacing-sm) var(--spacing-md);
}

:deep(.vditor-sv) {
  font-family: var(--font-mono);
  font-size: var(--font-size-base);
  line-height: var(--line-height-relaxed);
}

:deep(.vditor-sv .vditor-sv__preview) {
  padding: var(--spacing-lg) var(--spacing-xl);
}

/* 全屏模式下的内边距 */
:deep(.vditor.vditor--fullscreen) {
  padding: 0 var(--spacing-xl) !important;
}

:deep(.vditor.vditor--fullscreen .vditor-toolbar) {
  padding-left: var(--spacing-xl) !important;
  padding-right: var(--spacing-xl) !important;
}

:deep(.vditor.vditor--fullscreen .vditor-content) {
  padding-left: var(--spacing-xl) !important;
  padding-right: var(--spacing-xl) !important;
}

:deep(.vditor.vditor--fullscreen .vditor-sv textarea) {
  padding-left: var(--spacing-xl) !important;
  padding-right: var(--spacing-xl) !important;
}

:deep(.vditor.vditor--fullscreen .vditor-sv__preview) {
  padding-left: 40px !important;
  padding-right: 40px !important;
}

:deep(.vditor.vditor--fullscreen .vditor-sv) {
  padding: 0 !important;
}

:deep(.vditor-sv .vditor-sv__preview) {
  padding: var(--spacing-lg) var(--spacing-xl);
}

:deep(.vditor-reset) {
  font-family: var(--font-sans);
  font-size: var(--font-size-base);
  line-height: var(--line-height-relaxed);
  color: var(--color-text-primary);
}
</style>

