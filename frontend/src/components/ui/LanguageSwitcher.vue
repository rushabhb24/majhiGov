<script setup>
import { useUiStore } from '../../stores/ui'
const uiStore = useUiStore()
</script>

<template>
  <div class="tw-relative tw-inline-block">
    <select
      :value="uiStore.currentLanguage"
      @change="uiStore.setLanguage($event.target.value)"
      class="tw-appearance-none tw-px-3 tw-pr-7 tw-py-1.5 tw-text-xs tw-font-bold tw-bg-muted/80 hover:tw-bg-muted tw-text-foreground tw-border-none tw-rounded-full tw-transition-colors tw-cursor-pointer tw-outline-none focus-visible:tw-ring-2 focus-visible:tw-ring-ring"
    >
      <option value="en">EN</option>
      <option value="hi">हि</option>
      <option value="mr">मरा</option>
    </select>
    <!-- Chevron -->
    <div class="tw-pointer-events-none tw-absolute tw-inset-y-0 tw-right-2 tw-flex tw-items-center tw-text-muted-foreground">
      <svg class="tw-fill-current tw-h-3 tw-w-3" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"><path d="M9.293 12.95l.707.707L15.657 8l-1.414-1.414L10 10.828 5.757 6.586 4.343 8z"/></svg>
    </div>
  </div>
</template>
