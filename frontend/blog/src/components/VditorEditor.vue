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
  await loadHighlightTheme(codeThemeValue)
  const mdTheme = await getMarkdownTheme()
  
  // 使用加载的配置初始化 Vditor
  // 已手动设置window.VditorI18n，Vditor不会尝试动态加载i18n文件
  vditor = new Vditor(vditorRef.value, {
    height: props.height,
    mode: 'sv', // 分屏预览模式
    placeholder: '请输入文章内容，支持 Markdown 语法...',
    theme: 'classic',
    icon: 'material',
    typewriterMode: false,
    lang: 'zh_CN', // 设置语言为中文（已通过window.VditorI18n设置）
    cdn: '', // 禁用CDN，使用本地资源
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
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  width: 100%;
}

:deep(.vditor-toolbar) {
  background-color: #fff;
  border-bottom: 1px solid #e4e7ed;
  padding: 8px 12px;
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
  color: #24292e;
}
</style>
