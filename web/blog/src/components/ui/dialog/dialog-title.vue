<script setup>
import { cn } from '@/lib/utils'
import { DialogTitle } from 'radix-vue'

const props = defineProps({
  class: {
    type: String,
    default: '',
  },
})
</script>

<template>
  <DialogTitle :class="cn('text-lg font-semibold leading-none tracking-tight', props.class)">
    <slot />
  </DialogTitle>
</template>
