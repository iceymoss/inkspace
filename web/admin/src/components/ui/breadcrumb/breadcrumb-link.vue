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
  <a :class="cn('transition-colors hover:text-foreground', props.class)">
    <slot />
  </a>
</template>
