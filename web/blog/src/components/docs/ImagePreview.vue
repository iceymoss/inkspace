<template>
  <figure class="image-preview">
    <img
      :src="src"
      :alt="fileName"
    >
    <figcaption>{{ fileName }}</figcaption>
  </figure>
</template>

<script setup>
defineProps({
  src: { type: String, required: true },
  fileName: { type: String, default: '图片预览' }
})
</script>

<style scoped>
.image-preview { margin: 0; border: 1px solid var(--theme-border-light); background: repeating-conic-gradient(color-mix(in srgb, var(--theme-border) 45%, transparent) 0 25%, transparent 0 50%) 50% / 20px 20px; text-align: center; }
.image-preview img { display: block; max-width: 100%; max-height: 76vh; margin: 0 auto; object-fit: contain; }
.image-preview figcaption { padding: 9px 14px; border-top: 1px solid var(--theme-border-light); background: var(--theme-bg-card); color: var(--theme-text-tertiary); font-size: 11px; }
@media (max-width: 600px) { .image-preview { margin-right: -16px; margin-left: -16px; border-right: 0; border-left: 0; } }
</style>
