<script setup>
defineProps({
  hoverable: {
    type: Boolean,
    default: false
  }
})
</script>

<template>
  <div
    class="glass tw-rounded-2xl tw-p-5 tw-transition-all tw-duration-300"
    :class="{ 'hover:tw-translate-y-[-4px] hover:tw-shadow-xl hover:tw-border-primary/30': hoverable }"
  >
    <slot />
  </div>
</template>
