<script setup>
import { SelectRoot } from 'radix-vue'
</script>

<template>
  <SelectRoot v-bind="$attrs">
    <slot />
  </SelectRoot>
</template>
