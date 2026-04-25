<script setup>
import { TabsRoot } from 'radix-vue'
</script>

<template>
  <TabsRoot v-bind="$attrs">
    <slot />
  </TabsRoot>
</template>
