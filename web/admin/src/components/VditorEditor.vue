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
  await nextTick()
  
  const codeThemeValue = await loadCodeTheme()
  await loadHighlightTheme(codeThemeValue)
  
  const mdTheme = await getMarkdownTheme()
  
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
      link.onerror = () => resolve()
      document.head.appendChild(link)
    })
  }
  
  await loadVditorContentTheme(mdTheme || 'light')
  await new Promise(resolve => setTimeout(resolve, 100))
  
  vditor = new Vditor(vditorRef.value, {
    height: props.height,
    mode: 'sv',
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
        style: codeThemeValue || 'github',
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
      max: 5 * 1024 * 1024,
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
  border-radius: var(--radius);
  width: 100%;
  min-width: 0;
}

:deep(.vditor) {
  border: 1px solid hsl(var(--border));
  border-radius: var(--radius);
  width: 100%;
}

:deep(.vditor-toolbar) {
  background-color: hsl(var(--background));
  border-bottom: 1px solid hsl(var(--border));
  padding: 0.25rem 0.5rem;
}

:deep(.vditor-sv) {
  font-family: var(--font-mono);
  font-size: var(--font-size-base);
  line-height: var(--line-height-relaxed);
}

:deep(.vditor-sv .vditor-sv__preview) {
  padding: var(--spacing-lg) var(--spacing-xl);
}

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
  color: hsl(var(--foreground));
}
</style>
