<script setup>
import { computed } from 'vue'
import { cn } from '@/lib/utils'
import { RadioGroupItem, RadioGroupIndicator } from 'radix-vue'
import { Circle } from 'lucide-vue-next'

const props = defineProps({
  value: {
    type: String,
    required: true,
  },
  disabled: {
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
    'aspect-square h-4 w-4 rounded-full border border-primary text-primary ring-offset-background focus:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50',
    props.class,
  ),
)
</script>

<template>
  <RadioGroupItem :value="value" :disabled="disabled" :class="classes">
    <RadioGroupIndicator class="flex items-center justify-center">
      <Circle class="h-2.5 w-2.5 fill-current text-current" />
    </RadioGroupIndicator>
  </RadioGroupItem>
</template>
