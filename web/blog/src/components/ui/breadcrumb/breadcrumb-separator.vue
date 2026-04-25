<script setup>
import { cn } from '@/lib/utils'
import { ChevronRight } from 'lucide-vue-next'

const props = defineProps({
  class: {
    type: String,
    default: '',
  },
})
</script>

<template>
  <span :class="cn('[&>svg]:size-3.5', props.class)" role="presentation" aria-hidden="true">
    <ChevronRight />
  </span>
</template>
