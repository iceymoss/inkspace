<template>
  <div ref="vditorRef" class="vditor-container"></div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import Vditor from 'vditor'
import 'vditor/dist/index.css'
import 'vditor/dist/js/i18n/zh_CN.js'
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
  
  // 检查 ref 是否存在
  if (!vditorRef.value) {
    console.error('VditorEditor: vditorRef.value is null')
    return
  }
  
  // 先同步加载主题配置，确保在初始化 Vditor 前配置已准备好
  const codeThemeValue = await loadCodeTheme()
  await loadHighlightTheme(codeThemeValue)
  const mdTheme = await getMarkdownTheme()
  
  // 使用加载的配置初始化 Vditor
  // 中文资源随前端一起打包，避免运行时依赖 CDN。
  try {
    vditor = new Vditor(vditorRef.value, {
    height: props.height,
    mode: 'sv', // 分屏预览模式
    placeholder: '请输入文章内容，支持 Markdown 语法...',
    theme: 'classic',
    icon: 'material',
    typewriterMode: false,
    lang: 'zh_CN',
    // 不设置 cdn，使用 Vditor 默认的 CDN（unpkg.com）
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
      url: '/api/upload/markdown-image',
      max: 5 * 1024 * 1024, // 5MB
      accept: 'image/*',
      fieldName: 'file',
      // 不需要Authorization header，这是公开API
      format(files, responseText) {
        const response = JSON.parse(responseText)
        if (response.code === 0) {
          return JSON.stringify({
            msg: '',
            code: 0,
            data: {
              errFiles: [],
              succMap: {
                // 直接使用相对路径（博客系统）
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
  } catch (error) {
    console.error('VditorEditor: Failed to initialize Vditor', error)
  }
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
  border-radius: 8px;
  width: 100%;
  min-width: 0;
}

:deep(.vditor) {
  border: 1px solid var(--theme-border);
  border-radius: 8px;
  width: 100%;
  background: var(--theme-bg-card);
  color: var(--theme-text-primary);
}

:deep(.vditor-toolbar) {
  background-color: var(--theme-bg-card);
  border-bottom: 1px solid var(--theme-border);
  padding: 8px 12px;
}

:deep(.vditor-toolbar__item button) {
  color: var(--theme-text-secondary);
}

:deep(.vditor-toolbar__item button:hover),
:deep(.vditor-toolbar__item button.vditor-menu--current) {
  background: var(--theme-bg-hover);
  color: var(--theme-primary);
}

:deep(.vditor-content),
:deep(.vditor-sv),
:deep(.vditor-sv textarea),
:deep(.vditor-sv__preview) {
  background: var(--theme-bg-card);
  color: var(--theme-text-primary);
}

:deep(.vditor-sv textarea) {
  border-right-color: var(--theme-border);
  caret-color: var(--theme-primary);
}

:deep(.vditor-counter) {
  background: color-mix(in srgb, var(--theme-bg-card) 92%, transparent);
  color: var(--theme-text-tertiary);
}

:deep(.vditor-sv) {
  font-family: 'SF Mono', 'Monaco', 'Inconsolata', 'Fira Code', 'Menlo', 'Consolas', 'Source Code Pro', monospace;
  font-size: 15px;
  line-height: 1.8;
}

:deep(.vditor-sv .vditor-sv__preview) {
  padding: 20px 32px;
}

/* 全屏模式下的内边距 */
:deep(.vditor.vditor--fullscreen) {
  padding: 0 32px !important;
}

:deep(.vditor.vditor--fullscreen .vditor-toolbar) {
  padding-left: 32px !important;
  padding-right: 32px !important;
}

:deep(.vditor.vditor--fullscreen .vditor-content) {
  padding-left: 32px !important;
  padding-right: 32px !important;
}

:deep(.vditor.vditor--fullscreen .vditor-sv textarea) {
  padding-left: 32px !important;
  padding-right: 32px !important;
}

:deep(.vditor.vditor--fullscreen .vditor-sv__preview) {
  padding-left: 40px !important;
  padding-right: 40px !important;
}

:deep(.vditor.vditor--fullscreen .vditor-sv) {
  padding: 0 !important;
}

:deep(.vditor-sv .vditor-sv__preview) {
  padding: 20px 32px;
}

:deep(.vditor-reset) {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', 'Helvetica Neue', Arial, sans-serif;
  font-size: 16px;
  line-height: 1.8;
  color: var(--theme-text-primary);
}

:deep(.vditor-reset p),
:deep(.vditor-reset li) {
  color: var(--theme-text-secondary);
}

:deep(.vditor-reset blockquote) {
  color: var(--theme-text-tertiary);
  border-left-color: var(--theme-primary);
}

@media (max-width: 700px) {
  :deep(.vditor-toolbar) { padding: 6px; }
  :deep(.vditor-sv .vditor-sv__preview) { padding: 16px; }
}
</style>
