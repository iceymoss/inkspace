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
  <ul :class="cn('flex flex-row items-center gap-1', props.class)">
    <slot />
  </ul>
</template>
