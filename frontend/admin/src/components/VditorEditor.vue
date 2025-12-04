<template>
  <div ref="vditorRef" class="vditor-container"></div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import Vditor from 'vditor'
import 'vditor/dist/index.css'

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

onMounted(() => {
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
        style: 'github',
        lineNumber: true,
      },
      markdown: {
        toc: true,
        mark: true,
        footnotes: true,
        autoSpace: true, // 自动空格
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
  if (vditor && newVal !== vditor.getValue()) {
    vditor.setValue(newVal || '')
  }
})

defineExpose({
  getValue: () => vditor?.getValue(),
  setValue: (value) => vditor?.setValue(value),
  focus: () => vditor?.focus(),
})
</script>

<style scoped>
.vditor-container {
  border-radius: 8px;
  overflow: hidden;
  width: 100%;
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

