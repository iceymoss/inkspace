<template>
  <div
    ref="editorElement"
    class="code-text-editor"
  />
</template>

<script setup>
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { basicSetup } from 'codemirror'
import { EditorState } from '@codemirror/state'
import { EditorView } from '@codemirror/view'
import { StreamLanguage } from '@codemirror/language'

const props = defineProps({
  modelValue: { type: String, default: '' },
  language: { type: String, default: '' },
  readOnly: { type: Boolean, default: false }
})
const emit = defineEmits(['update:modelValue'])
const editorElement = ref(null)
let view

async function languageExtension(language) {
  switch (language) {
    case 'json': case 'jsonl': return (await import('@codemirror/lang-json')).json()
    case 'yaml': return (await import('@codemirror/lang-yaml')).yaml()
    case 'javascript': case 'jsx': return (await import('@codemirror/lang-javascript')).javascript({ jsx: language === 'jsx' })
    case 'typescript': case 'tsx': return (await import('@codemirror/lang-javascript')).javascript({ typescript: true, jsx: language === 'tsx' })
    case 'vue': return (await import('@codemirror/lang-vue')).vue()
    case 'python': return (await import('@codemirror/lang-python')).python()
    case 'html': return (await import('@codemirror/lang-html')).html()
    case 'xml': return (await import('@codemirror/lang-xml')).xml()
    case 'css': case 'scss': case 'less': return (await import('@codemirror/lang-css')).css()
    case 'java': return (await import('@codemirror/lang-java')).java()
    case 'sql': return (await import('@codemirror/lang-sql')).sql()
    case 'c': case 'cpp': return (await import('@codemirror/lang-cpp')).cpp()
    case 'rust': return (await import('@codemirror/lang-rust')).rust()
    case 'php': return (await import('@codemirror/lang-php')).php()
    case 'go': return StreamLanguage.define((await import('@codemirror/legacy-modes/mode/go')).go)
    case 'shell': return StreamLanguage.define((await import('@codemirror/legacy-modes/mode/shell')).shell)
    case 'ruby': return StreamLanguage.define((await import('@codemirror/legacy-modes/mode/ruby')).ruby)
    case 'csharp': case 'kotlin': case 'scala': return StreamLanguage.define((await import('@codemirror/legacy-modes/mode/clike'))[language])
    case 'swift': return StreamLanguage.define((await import('@codemirror/legacy-modes/mode/swift')).swift)
    default: return []
  }
}

onMounted(async () => {
  const syntax = await languageExtension(props.language)
  view = new EditorView({
    parent: editorElement.value,
    state: EditorState.create({
      doc: props.modelValue,
      extensions: [
        basicSetup,
        syntax,
        EditorState.readOnly.of(props.readOnly),
        EditorView.editable.of(!props.readOnly),
        EditorView.lineWrapping,
        EditorView.updateListener.of(update => {
          if (update.docChanged) emit('update:modelValue', update.state.doc.toString())
        }),
        EditorView.theme({
          '&': props.readOnly
            ? { maxHeight: '70vh', minHeight: '240px', backgroundColor: 'var(--theme-bg-card)', color: 'var(--theme-text-primary)' }
            : { height: 'calc(100vh - 245px)', minHeight: '420px', backgroundColor: 'var(--theme-bg-card)', color: 'var(--theme-text-primary)' },
          '.cm-scroller': { overflow: 'auto', fontFamily: "'SFMono-Regular', Consolas, monospace", fontSize: '13px', lineHeight: '1.7' },
          '.cm-gutters': { backgroundColor: 'var(--theme-bg-secondary)', color: 'var(--theme-text-tertiary)', borderRight: '1px solid var(--theme-border)' },
          '.cm-activeLine, .cm-activeLineGutter': { backgroundColor: 'color-mix(in srgb, var(--theme-primary) 7%, transparent)' },
          '&.cm-focused': { outline: props.readOnly ? 'none' : '1px solid var(--theme-primary)' }
        })
      ]
    })
  })
})

watch(() => props.modelValue, value => {
  if (!view || value === view.state.doc.toString()) return
  view.dispatch({ changes: { from: 0, to: view.state.doc.length, insert: value } })
})

onBeforeUnmount(() => view?.destroy())
</script>

<style scoped>
.code-text-editor { border: 1px solid var(--theme-border); }
</style>
