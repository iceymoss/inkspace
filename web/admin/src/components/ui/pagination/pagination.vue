<script setup>
import { cn } from '@/lib/utils'

const props = defineProps({
  class: {
    type: String,
    default: '',
  },
})
</script>

<template>
  <nav :class="cn('mx-auto flex w-full justify-center', props.class)" role="navigation" aria-label="pagination">
    <slot />
  </nav>
</template>
