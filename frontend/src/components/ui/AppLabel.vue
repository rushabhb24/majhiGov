<script setup>
defineProps({
  for: {
    type: String,
    default: ''
  }
})
</script>

<template>
  <label
    :for="for"
    class="tw-block tw-text-xs tw-font-semibold tw-text-muted-foreground tw-mb-1 tw-font-heading"
  >
    <slot />
  </label>
</template>
