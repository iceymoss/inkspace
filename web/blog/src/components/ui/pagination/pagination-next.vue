<script setup>
import { cn } from '@/lib/utils'
import { ChevronRight } from 'lucide-vue-next'

const props = defineProps({
  class: {
    type: String,
    default: '',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})
</script>

<template>
  <button
    :class="cn(
      'inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 h-9 px-3 hover:bg-accent hover:text-accent-foreground cursor-pointer gap-1',
      props.class,
    )"
    :disabled="disabled"
  >
    <slot />
    <ChevronRight class="h-4 w-4" />
  </button>
</template>
