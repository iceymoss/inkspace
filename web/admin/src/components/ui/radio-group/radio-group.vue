<script setup>
import { cn } from '@/lib/utils'
import { RadioGroupRoot } from 'radix-vue'

const props = defineProps({
  class: {
    type: String,
    default: '',
  },
  modelValue: {
    type: String,
    default: undefined,
  },
  defaultValue: {
    type: String,
    default: undefined,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue'])
</script>

<template>
  <RadioGroupRoot
    :class="cn('grid gap-2', props.class)"
    :model-value="modelValue"
    :default-value="defaultValue"
    :disabled="disabled"
    @update:model-value="emit('update:modelValue', $event)"
  >
    <slot />
  </RadioGroupRoot>
</template>
