<script setup>
import { cn } from '@/lib/utils'
import { Inbox } from 'lucide-vue-next'

const props = defineProps({
  title: {
    type: String,
    default: '',
  },
  description: {
    type: String,
    default: '',
  },
  class: {
    type: String,
    default: '',
  },
})
</script>

<template>
  <div :class="cn('flex flex-col items-center justify-center py-12 text-center', props.class)">
    <slot name="icon">
      <Inbox class="h-12 w-12 text-muted-foreground" />
    </slot>
    <h3 v-if="title" class="mt-4 text-lg font-semibold">{{ title }}</h3>
    <p v-if="description" class="mt-2 text-sm text-muted-foreground">{{ description }}</p>
    <slot />
  </div>
</template>
