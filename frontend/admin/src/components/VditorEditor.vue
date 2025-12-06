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
      mode: mdTheme || 'light', // 使用配置的 Markdown 主题
      delay: 500, // 增加延迟，确保预览同步稳定
      hljs: {
        style: codeThemeValue || 'github', // 使用配置的代码主题
        lineNumber: true,
        enable: true,
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
                [files[0].name]: `http://localhost:8083${response.data.url}`
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
      // Vditor 初始化完成后设置初始值
      // 使用 nextTick 确保 Vditor 完全初始化
      nextTick(() => {
        if (props.modelValue && vditor) {
          vditor.setValue(props.modelValue)
        }
        // 确保 sv 模式下的预览区域可见并正确渲染
        setTimeout(() => {
          if (vditor && vditor.vditor) {
            const svElement = vditor.vditor.querySelector('.vditor-sv')
            const previewElement = vditor.vditor.querySelector('.vditor-sv__preview')
            if (svElement && previewElement) {
              // 确保预览区域可见
              previewElement.style.display = 'block'
              previewElement.style.visibility = 'visible'
              previewElement.style.width = '50%'
              // 触发预览更新
              if (vditor.getValue) {
                const content = vditor.getValue()
                if (content) {
                  vditor.setValue(content)
                }
              }
            }
          }
        }, 100)
      })
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
  display: flex;
  flex-direction: row;
}

:deep(.vditor-sv textarea) {
  flex: 1;
  min-width: 0;
  width: 50%;
}

:deep(.vditor-sv .vditor-sv__preview) {
  flex: 1;
  min-width: 0;
  width: 50%;
  display: block !important;
  visibility: visible !important;
  border-left: 1px solid #e4e7ed;
  overflow-y: auto;
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

:deep(.vditor-sv .vditor-sv__preview) {
  padding: 20px 32px;
}

:deep(.vditor.vditor--fullscreen .vditor-sv) {
  padding: 0 20px;
}

:deep(.vditor.vditor--fullscreen .vditor-content) {
  padding: 0 20px;
}

:deep(.vditor-reset) {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', 'Helvetica Neue', Arial, sans-serif;
  font-size: 16px;
  line-height: 1.8;
  color: #24292e;
}
</style>

