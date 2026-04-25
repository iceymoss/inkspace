<script setup>
import { cn } from '@/lib/utils'

const props = defineProps({
  class: {
    type: String,
    default: '',
  },
})
</script>

<template>
  <li :class="cn(props.class)">
    <slot />
  </li>
</template>
