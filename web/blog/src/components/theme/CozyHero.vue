<template>
  <div class="cozy-hero">
    <div class="cozy-hero-copy">
      <span class="cozy-greeting">{{ settings.eyebrow }}</span>
      <h1>
        {{ title.before }}<em v-if="title.accent">{{ title.accent }}</em>{{ title.after }}
      </h1>
      <p>{{ settings.description }}</p>
      <div class="cozy-actions">
        <a
          class="cozy-primary"
          :href="safeLink(settings.primary_link)"
          @click="$emit('navigate', $event, settings.primary_link)"
        >{{ settings.primary_text }}</a>
        <a
          class="cozy-secondary"
          :href="safeLink(settings.secondary_link)"
          @click="$emit('navigate', $event, settings.secondary_link)"
        >{{ settings.secondary_text }} →</a>
      </div>
    </div>

    <div class="cozy-stack" :class="`has-${displayPhotos.length}`" aria-label="最近发布的摄影作品">
      <div v-if="loading" class="cozy-hero-state">正在整理照片…</div>
      <div v-else-if="error" class="cozy-hero-state is-error">
        <span>照片暂时没有挂好。</span>
        <button type="button" @click="$emit('retry')">重新试试</button>
      </div>
      <div v-else-if="displayPhotos.length === 0" class="cozy-hero-state">
        <span>照片墙还空着，先去别处坐坐。</span>
        <a href="/photos" @click="$emit('navigate', $event, '/photos')">看看照片页</a>
      </div>
      <router-link
        v-for="(photo, index) in displayPhotos"
        v-else
        :key="photo.id"
        class="cozy-polaroid"
        :class="`polaroid-${index + 1}`"
        :to="`/works/${photo.id}`"
      >
        <i class="cozy-tape" aria-hidden="true" />
        <div class="cozy-polaroid-media">
          <img
            v-if="photoImage(photo) && !failedImages.has(photo.id)"
            :src="photoImage(photo)"
            :alt="photo.title"
            @error="markImageFailed(photo.id)"
          >
          <div v-else class="cozy-photo-placeholder" aria-hidden="true">
            <strong>{{ photo.title?.slice(0, 1) || '照' }}</strong>
            <span>暂无照片</span>
          </div>
        </div>
        <span class="cozy-caption">{{ photo.title }}</span>
        <small v-if="photo.metadata?.location">{{ photo.metadata.location }}</small>
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { computed, reactive } from 'vue'

const props = defineProps({
  settings: { type: Object, required: true },
  title: { type: Object, required: true },
  photos: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
  error: { type: Boolean, default: false }
})

defineEmits(['navigate', 'retry'])

const failedImages = reactive(new Set())
const displayPhotos = computed(() => props.photos.slice(0, 3))
const safeLink = link => /^https?:\/\//.test(link) || link?.startsWith('/') ? link : '#'
const photoImage = photo => photo.cover || (typeof photo.images?.[0] === 'string' ? photo.images[0] : photo.images?.[0]?.url) || ''
const markImageFailed = id => failedImages.add(id)
</script>

<style scoped>
.cozy-hero { display: grid; grid-template-columns: 1.08fr .92fr; align-items: center; gap: 46px; min-height: 500px; padding: 62px 0 70px; }
.cozy-greeting { display: inline-block; margin-bottom: 22px; padding: 6px 18px; border-radius: 99px; background: var(--moss-soft); color: var(--moss); font-size: 13px; transform: rotate(-1.5deg); }
h1 { margin: 0; color: var(--ink); font-size: clamp(32px, 5vw, 48px); font-weight: 800; line-height: 1.4; }
h1 em { position: relative; color: var(--mustard); font-style: normal; }
h1 em::after { position: absolute; right: 0; bottom: -7px; left: 0; height: 8px; background: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 120 10'%3E%3Cpath d='M0 7 Q15 1 30 7 T60 7 T90 7 T120 7' fill='none' stroke='%23D9A441' stroke-width='3' stroke-linecap='round'/%3E%3C/svg%3E") center / 100% 100% no-repeat; content: ''; }
.cozy-hero-copy > p { max-width: 34em; margin: 24px 0 0; color: var(--sub); font-size: 15.5px; line-height: 1.85; }
.cozy-actions { display: flex; align-items: center; gap: 18px; margin-top: 34px; }
.cozy-actions a { min-height: 44px; text-decoration: none; }
.cozy-primary { display: inline-flex; align-items: center; padding: 11px 28px; border-radius: 99px; background: var(--mustard); box-shadow: 0 4px 0 color-mix(in srgb, var(--mustard) 60%, #7a5a16); color: var(--mustard-ink); font-weight: 700; transition: transform .15s, box-shadow .15s; }
.cozy-primary:hover { box-shadow: 0 2px 0 color-mix(in srgb, var(--mustard) 60%, #7a5a16); transform: translateY(2px); }
.cozy-secondary { display: inline-flex; align-items: center; border-bottom: 2px dotted var(--sub); color: var(--sub); }
.cozy-secondary:hover { border-color: var(--moss); color: var(--moss); }
.cozy-stack { position: relative; min-height: 385px; }
.cozy-polaroid { position: absolute; width: min(220px, 58%); padding: 11px 11px 38px; border-radius: 4px; background: var(--polaroid); box-shadow: var(--shadow); color: var(--ink); text-decoration: none; transition: transform .35s cubic-bezier(.3, 1.4, .5, 1); }
.cozy-polaroid:hover, .cozy-polaroid:focus-visible { z-index: 9; transform: rotate(0) scale(1.05); }
.polaroid-1 { top: 9%; left: 2%; z-index: 1; transform: rotate(-6deg); }
.polaroid-2 { top: 1%; right: 0; z-index: 2; transform: rotate(3deg); }
.polaroid-3 { top: 39%; left: 24%; z-index: 3; transform: rotate(-2deg); }
.has-1 .polaroid-1 { top: 12%; left: 50%; transform: translateX(-50%) rotate(-3deg); }
.has-2 .polaroid-1 { left: 8%; }.has-2 .polaroid-2 { top: 28%; right: 7%; }
.cozy-tape { position: absolute; z-index: 2; top: -12px; left: 50%; width: 74px; height: 25px; border-radius: 2px; background: color-mix(in srgb, var(--mustard) 45%, var(--card)); opacity: .86; transform: translateX(-50%) rotate(-3deg); }
.cozy-polaroid-media { display: grid; aspect-ratio: 1; overflow: hidden; place-items: center; border-radius: 2px; background: var(--card-soft); }
.cozy-polaroid-media img { width: 100%; height: 100%; object-fit: cover; }
.cozy-photo-placeholder { display: grid; color: var(--sub); text-align: center; }.cozy-photo-placeholder strong { color: var(--moss); font-size: 42px; }.cozy-photo-placeholder span { font-size: 11px; }
.cozy-caption { position: absolute; right: 8px; bottom: 14px; left: 8px; overflow: hidden; font-size: 12px; text-align: center; text-overflow: ellipsis; white-space: nowrap; }
.cozy-polaroid small { position: absolute; right: 10px; bottom: 2px; left: 10px; overflow: hidden; color: var(--sub); font-size: 9px; text-align: center; text-overflow: ellipsis; white-space: nowrap; }
.cozy-hero-state { position: absolute; inset: 70px 20px; display: grid; align-content: center; justify-items: center; gap: 12px; padding: 30px; border: 2px dashed var(--line); border-radius: 20px 14px 18px 16px; background: var(--card); color: var(--sub); text-align: center; transform: rotate(1deg); }
.cozy-hero-state :is(a, button) { padding: 7px 14px; border: 1px dashed var(--moss); border-radius: 99px; background: transparent; color: var(--moss); cursor: pointer; text-decoration: none; }.cozy-hero-state.is-error { color: var(--danger); }
@media (max-width: 900px) { .cozy-hero { grid-template-columns: 1fr; }.cozy-stack { width: 100%; max-width: 430px; margin: 0 auto; } }
@media (max-width: 560px) { .cozy-hero { gap: 30px; min-height: 0; padding: 38px 0 46px; }.cozy-actions { align-items: stretch; flex-direction: column; }.cozy-actions a { justify-content: center; }.cozy-stack { min-height: 315px; }.cozy-polaroid { width: min(190px, 58%); }.polaroid-3 { top: 38%; left: 21%; }.cozy-hero-state { inset: 35px 0; } }
@media (prefers-reduced-motion: reduce) { .cozy-greeting, .cozy-polaroid, .cozy-polaroid:hover, .cozy-polaroid:focus-visible, .cozy-hero-state { transition: none; transform: none; }.has-1 .polaroid-1 { left: 30%; transform: none; } }
</style>
