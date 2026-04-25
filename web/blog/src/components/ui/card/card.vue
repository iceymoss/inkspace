<script setup>
import { computed } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps({
  class: {
    type: String,
    default: '',
  },
})

const mergedClass = computed(() =>
  cn('rounded-lg border bg-card text-card-foreground shadow-sm', props.class),
)
</script>

<template>
  <div :class="mergedClass" style="box-shadow: inset 0 1px 0 rgba(255,255,255,.1)">
    <slot />
  </div>
</template>
