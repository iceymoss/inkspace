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
        Authorization: `Bearer ${localStorage.getItem('user_token')}`
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
                [files[0].name]: `http://localhost:8081${response.data.url}`
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

// 暴露方法给父组件
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

:deep(.vditor-toolbar .vditor-tooltipped) {
  border-radius: 4px;
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

:deep(.vditor.vditor--fullscreen .vditor-sv) {
  padding: 0 !important;
}

:deep(.vditor.vditor--fullscreen .vditor-sv textarea) {
  padding-left: 32px !important;
  padding-right: 32px !important;
}

:deep(.vditor.vditor--fullscreen .vditor-sv__preview) {
  padding-left: 40px !important;
  padding-right: 40px !important;
}

:deep(.vditor-reset) {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', 'Helvetica Neue', Arial, sans-serif;
  font-size: 16px;
  line-height: 1.8;
  color: #24292e;
}

:deep(.vditor-reset h1) {
  font-size: 2em;
  margin-top: 0;
  margin-bottom: 16px;
  font-weight: 600;
  line-height: 1.25;
  padding-bottom: 0.3em;
  border-bottom: 1px solid #eaecef;
}

:deep(.vditor-reset h2) {
  font-size: 1.5em;
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
  line-height: 1.25;
  padding-bottom: 0.3em;
  border-bottom: 1px solid #eaecef;
}

:deep(.vditor-reset h3) {
  font-size: 1.25em;
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
  line-height: 1.25;
}

:deep(.vditor-reset p) {
  margin-top: 0;
  margin-bottom: 16px;
  line-height: 1.8;
}

:deep(.vditor-reset code) {
  font-family: 'SF Mono', 'Monaco', 'Consolas', 'Courier New', monospace;
  font-size: 85%;
  padding: 0.2em 0.4em;
  margin: 0;
  background-color: rgba(175, 184, 193, 0.2);
  border-radius: 6px;
}

:deep(.vditor-reset pre) {
  background-color: #f6f8fa;
  border-radius: 6px;
  padding: 16px;
  overflow: auto;
  font-size: 85%;
  line-height: 1.45;
  margin-bottom: 16px;
}

:deep(.vditor-reset pre code) {
  background-color: transparent;
  padding: 0;
  border-radius: 0;
}

:deep(.vditor-reset blockquote) {
  padding: 0 1em;
  color: #57606a;
  border-left: 0.25em solid #d0d7de;
  margin: 0 0 16px;
}

:deep(.vditor-reset a) {
  color: #0969da;
  text-decoration: none;
}

:deep(.vditor-reset a:hover) {
  text-decoration: underline;
}

:deep(.vditor-reset table) {
  border-spacing: 0;
  border-collapse: collapse;
  margin-bottom: 16px;
  width: 100%;
}

:deep(.vditor-reset table th,
.vditor-reset table td) {
  padding: 6px 13px;
  border: 1px solid #d0d7de;
}

:deep(.vditor-reset table th) {
  font-weight: 600;
  background-color: #f6f8fa;
}

:deep(.vditor-reset img) {
  max-width: 100%;
  box-sizing: content-box;
  background-color: #fff;
  border-radius: 6px;
}
</style>

