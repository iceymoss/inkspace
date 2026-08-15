<template>
  <section class="pdf-preview">
    <iframe
      :src="src"
      title="PDF 文档预览"
    />
    <p>浏览器无法显示 PDF 时，请下载原文件查看。</p>
  </section>
</template>

<script setup>
defineProps({ src: { type: String, required: true } })
</script>

<style scoped>
.pdf-preview { border: 1px solid var(--theme-border-light); background: var(--theme-bg-secondary); }
.pdf-preview iframe { display: block; width: 100%; min-height: min(76vh, 860px); border: 0; background: #fff; }
.pdf-preview p { margin: 0; padding: 9px 14px; color: var(--theme-text-tertiary); font-size: 11px; text-align: center; }
@media (max-width: 600px) { .pdf-preview { margin-right: -16px; margin-left: -16px; border-right: 0; border-left: 0; } .pdf-preview iframe { min-height: 68vh; } }
</style>
