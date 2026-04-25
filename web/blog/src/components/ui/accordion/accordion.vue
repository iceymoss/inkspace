<script setup>
import { cn } from '@/lib/utils'
import { AccordionRoot } from 'radix-vue'

const props = defineProps({
  type: {
    type: String,
    default: 'single',
  },
  collapsible: {
    type: Boolean,
    default: true,
  },
  class: {
    type: String,
    default: '',
  },
})
</script>

<template>
  <AccordionRoot
    :type="type"
    :collapsible="collapsible"
    :class="cn(props.class)"
  >
    <slot />
  </AccordionRoot>
</template>
