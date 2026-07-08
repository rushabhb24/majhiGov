<script setup>
import { computed } from 'vue'

const props = defineProps({
  tone: {
    type: String,
    default: 'neutral',
    validator: (v) => ['info', 'success', 'warn', 'danger', 'neutral'].includes(v)
  }
})

const classes = computed(() => {
  const base = 'tw-inline-flex tw-items-center tw-gap-1.5 tw-px-2.5 tw-py-0.5 tw-text-xs tw-font-semibold tw-rounded-full'
  
  const tones = {
    info: 'tw-bg-accent/15 tw-text-accent',
    success: 'tw-bg-success/15 tw-text-success',
    warn: 'tw-bg-warning/15 tw-text-warning',
    danger: 'tw-bg-destructive/15 tw-text-destructive',
    neutral: 'tw-bg-muted tw-text-muted-foreground'
  }
  
  return `${base} ${tones[props.tone]}`
})
</script>

<template>
  <span :class="classes">
    <slot />
  </span>
</template>
