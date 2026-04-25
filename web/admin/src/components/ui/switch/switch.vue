<script setup>
import { computed } from 'vue'
import { cn } from '@/lib/utils'
import { SwitchRoot, SwitchThumb } from 'radix-vue'

const props = defineProps({
  class: {
    type: String,
    default: '',
  },
  defaultChecked: {
    type: Boolean,
    default: undefined,
  },
  checked: {
    type: Boolean,
    default: undefined,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  id: {
    type: String,
    default: undefined,
  },
})

const emit = defineEmits(['update:checked'])

const rootClass = computed(() =>
  cn(
    'peer inline-flex h-[24px] w-[44px] shrink-0 cursor-pointer items-center rounded-full border-2 border-transparent transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background disabled:cursor-not-allowed disabled:opacity-50 data-[state=checked]:bg-primary data-[state=unchecked]:bg-input',
    props.class,
  ),
)
</script>

<template>
  <SwitchRoot
    :class="rootClass"
    :default-checked="defaultChecked"
    :checked="checked"
    :disabled="disabled"
    :id="id"
    @update:checked="emit('update:checked', $event)"
  >
    <SwitchThumb
      class="pointer-events-none block h-5 w-5 rounded-full bg-background shadow-lg ring-0 transition-transform data-[state=checked]:translate-x-5 data-[state=unchecked]:translate-x-0"
    />
  </SwitchRoot>
</template>
