<script setup>
import { DropdownMenuRoot } from 'radix-vue'
</script>

<template>
  <DropdownMenuRoot v-bind="$attrs">
    <slot />
  </DropdownMenuRoot>
</template>
