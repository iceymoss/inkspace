<script setup>
import { cn } from '@/lib/utils'
import { AccordionItem } from 'radix-vue'

const props = defineProps({
  value: {
    type: String,
    required: true,
  },
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
  <AccordionItem :value="value" :disabled="disabled" :class="cn('border-b', props.class)">
    <slot />
  </AccordionItem>
</template>
