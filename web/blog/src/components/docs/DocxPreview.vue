<template>
  <section class="docx-preview">
    <div
      ref="container"
      class="docx-pages"
    />
  </section>
</template>

<script setup>
import { onMounted, ref } from 'vue'

const props = defineProps({ blob: { type: Blob, required: true } })
const emit = defineEmits(['error'])
const container = ref(null)

onMounted(async () => {
  try {
    const { renderAsync } = await import('docx-preview')
    await renderAsync(props.blob, container.value, null, {
      breakPages: true,
      ignoreLastRenderedPageBreak: false,
      renderAltChunks: false,
      renderComments: false,
      useBase64URL: true
    })
    for (const link of container.value.querySelectorAll('a[href]')) {
      let protocol = ''
      try { protocol = new URL(link.href, window.location.origin).protocol } catch { /* invalid links stay disabled */ }
      if (!['http:', 'https:', 'mailto:'].includes(protocol)) link.removeAttribute('href')
      link.rel = 'noopener noreferrer'
      link.target = '_blank'
    }
  } catch (error) {
    emit('error', error)
  }
})
</script>

<style scoped>
.docx-preview { max-height: 78vh; overflow: auto; border: 1px solid var(--theme-border-light); background: #d8dce2; }
.docx-pages { min-width: 680px; padding: 24px; }
.docx-pages :deep(.docx-wrapper) { padding: 0; background: transparent; }
.docx-pages :deep(.docx) { margin: 0 auto 20px; box-shadow: 0 8px 30px rgb(15 23 42 / 18%); }
@media (max-width: 600px) { .docx-preview { margin-right: -16px; margin-left: -16px; border-right: 0; border-left: 0; } }
</style>
