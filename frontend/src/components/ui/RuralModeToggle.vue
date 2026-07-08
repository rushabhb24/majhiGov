<script setup>
import { useUiStore } from '../../stores/ui'
const uiStore = useUiStore()
</script>

<template>
  <button 
    class="tw-w-9 tw-h-9 tw-flex tw-items-center tw-justify-center tw-rounded-full tw-bg-muted/80 hover:tw-bg-muted tw-transition-colors tw-border-none tw-cursor-pointer tw-outline-none focus-visible:tw-ring-2 focus-visible:tw-ring-ring"
    :class="uiStore.ruralMode ? 'tw-text-primary tw-bg-primary/10' : 'tw-text-muted-foreground hover:tw-text-foreground'"
    @click="uiStore.toggleRuralMode"
    :title="uiStore.ruralMode ? 'Normal Mode' : 'Rural Mode'"
  >
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path></svg>
  </button>
</template>
