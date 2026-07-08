<script setup>
import { computed } from 'vue'

const props = defineProps({
  variant: {
    type: String,
    default: 'primary',
    validator: (v) => ['primary', 'secondary', 'ghost', 'outline', 'destructive'].includes(v)
  },
  size: {
    type: String,
    default: 'md',
    validator: (v) => ['sm', 'md', 'lg'].includes(v)
  },
  tag: {
    type: String,
    default: 'button'
  },
  disabled: {
    type: Boolean,
    default: false
  }
})

defineEmits(['click'])

const classes = computed(() => {
  const base = 'tw-inline-flex tw-items-center tw-justify-center tw-font-heading tw-font-semibold tw-transition-all tw-duration-200 tw-rounded-lg tw-border tw-border-transparent tw-outline-none focus-visible:tw-ring-2 focus-visible:tw-ring-ring focus-visible:tw-ring-offset-2 tw-cursor-pointer disabled:tw-opacity-50 disabled:tw-cursor-not-allowed'
  
  const variants = {
    primary: 'tw-bg-primary tw-text-primary-foreground hover:tw-bg-primary/90 tw-border-transparent tw-shadow-[0_4px_12px_rgba(249,115,22,0.15)] dark:tw-shadow-[0_0_12px_rgba(249,115,22,0.35)] hover:dark:tw-shadow-[0_0_18px_rgba(249,115,22,0.55)]',
    secondary: 'tw-bg-muted tw-text-muted-foreground hover:tw-bg-muted/80 tw-border-transparent',
    ghost: 'tw-bg-transparent tw-text-foreground hover:tw-bg-muted tw-border-transparent',
    outline: 'tw-bg-transparent tw-text-foreground tw-border-border hover:tw-border-primary/60 hover:tw-text-primary hover:tw-bg-primary/5 dark:hover:tw-shadow-[0_0_10px_rgba(249,115,22,0.3)]',
    destructive: 'tw-bg-destructive tw-text-primary-foreground hover:tw-bg-destructive/90 tw-border-transparent'
  }

  const sizes = {
    sm: 'tw-px-3 tw-py-1.5 tw-text-xs',
    md: 'tw-px-4 tw-py-2 tw-text-sm',
    lg: 'tw-px-6 tw-py-3 tw-text-base'
  }

  return `${base} ${variants[props.variant]} ${sizes[props.size]}`
})
</script>

<template>
  <component
    :is="tag"
    :class="classes"
    :disabled="disabled"
    @click="$emit('click', $event)"
  >
    <slot />
  </component>
</template>
