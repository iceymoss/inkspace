<script setup>
import { computed } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps({
  isActive: {
    type: Boolean,
    default: false,
  },
  class: {
    type: String,
    default: '',
  },
})

const classes = computed(() =>
  cn(
    'inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 h-9 w-9 hover:bg-accent hover:text-accent-foreground cursor-pointer',
    props.isActive && 'bg-primary text-primary-foreground',
    props.class,
  ),
)
</script>

<template>
  <button :class="classes">
    <slot />
  </button>
</template>
