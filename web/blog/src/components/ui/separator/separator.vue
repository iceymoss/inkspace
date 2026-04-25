<script setup>
import { cn } from '@/lib/utils'

const props = defineProps({
  orientation: {
    type: String,
    default: 'horizontal',
  },
  decorative: {
    type: Boolean,
    default: false,
  },
  class: {
    type: String,
    default: '',
  },
})

const isHorizontal = props.orientation === 'horizontal'
</script>

<template>
  <component
    :is="decorative ? 'div' : 'span'"
    :role="decorative ? 'none' : 'separator'"
    :aria-orientation="decorative ? undefined : orientation"
    :class="cn(
      'shrink-0 bg-border',
      isHorizontal ? 'h-[1px] w-full' : 'w-[1px] h-full',
      props.class,
    )"
  />
</template>
