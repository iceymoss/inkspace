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
  <div :class="cn('relative w-full overflow-auto', props.class)">
    <table class="w-full caption-bottom text-sm">
      <slot />
    </table>
  </div>
</template>
