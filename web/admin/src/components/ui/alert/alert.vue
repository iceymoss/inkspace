<script setup>
import { computed } from 'vue'
import { cn } from '@/lib/utils'
import { cva } from 'class-variance-authority'

const props = defineProps({
  variant: {
    type: String,
    default: 'default',
  },
  class: {
    type: String,
    default: '',
  },
})

const alertVariants = cva(
  'relative w-full rounded-lg border p-4 [&>svg+div]:translate-y-[-3px] [&>svg]:absolute [&>svg]:left-4 [&>svg]:top-4 [&>svg]:text-foreground [&>svg]:pointer-events-none',
  {
    variants: {
      variant: {
        default: 'bg-background text-foreground',
        destructive:
          'border-destructive/50 text-destructive dark:border-destructive [&>svg]:text-destructive',
      },
    },
    defaultVariants: {
      variant: 'default',
    },
  },
)

const classes = computed(() => cn(alertVariants({ variant: props.variant }), props.class))
</script>

<template>
  <div :class="classes" role="alert">
    <slot />
  </div>
</template>
