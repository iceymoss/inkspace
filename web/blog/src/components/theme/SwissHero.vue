<template>
  <section class="swiss-hero" aria-labelledby="swiss-hero-title">
    <div class="swiss-hero-main">
      <div class="swiss-coordinates">
        <span class="is-primary">{{ settings.eyebrow }}</span>
        <span v-if="settings.location">{{ settings.location }}</span>
        <span v-if="settings.coordinates">{{ settings.coordinates }}</span>
        <span v-if="settings.established">{{ settings.established }}</span>
      </div>
      <h1 id="swiss-hero-title">
        {{ title.before }}<em v-if="title.accent">{{ title.accent }}</em>{{ title.after }}
      </h1>
      <p>{{ settings.description }}</p>
    </div>
    <div class="swiss-hero-stats" aria-label="站点公开内容统计">
      <div v-for="stat in statItems" :key="stat.key" class="swiss-stat-cell">
        <strong>{{ stat.error ? '--' : stat.value ?? '--' }}<em v-if="!stat.error && stat.value != null">+</em></strong>
        <span>{{ stat.error ? `${stat.label}暂不可用` : stat.label }}</span>
      </div>
    </div>
    <div class="swiss-hero-foot">
      <span>Blog — Works — Photography — Knowledge Base — {{ siteName }}</span>
      <div class="swiss-hero-actions">
        <a :href="safeLink(settings.secondary_link)" @click="$emit('navigate', $event, settings.secondary_link)">{{ settings.secondary_text }}</a>
        <a class="is-primary" :href="safeLink(settings.primary_link)" @click="$emit('navigate', $event, settings.primary_link)">{{ settings.primary_text }}</a>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  settings: { type: Object, required: true },
  title: { type: Object, required: true },
  stats: { type: Object, required: true },
  siteName: { type: String, default: 'InkSpace' }
})

defineEmits(['navigate'])

const statItems = computed(() => [
  { key: 'articles', label: '公开文章', value: props.stats.articleCount, error: props.stats.articleError },
  { key: 'works', label: '已发布作品', value: props.stats.workCount, error: props.stats.workError },
  { key: 'docs', label: '公开知识文档', value: props.stats.publicDocCount, error: props.stats.docError }
])
const safeLink = link => /^https?:\/\//.test(link) || link?.startsWith('/') ? link : '#'
</script>
