<script setup>
import { cn } from '@/lib/utils'
import { DropdownMenuItem } from 'radix-vue'

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
  <DropdownMenuItem
    :disabled="disabled"
    :class="
      cn(
        'relative flex cursor-pointer select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none transition-colors focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50',
        props.class,
      )
    "
    v-bind="$attrs"
  >
    <slot />
  </DropdownMenuItem>
</template>
