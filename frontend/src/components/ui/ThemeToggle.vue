<script setup>
import { useUiStore } from '../../stores/ui'
const uiStore = useUiStore()
</script>

<template>
  <button 
    class="tw-w-9 tw-h-9 tw-flex tw-items-center tw-justify-center tw-rounded-full tw-bg-muted/80 hover:tw-bg-muted tw-text-muted-foreground hover:tw-text-foreground tw-transition-colors tw-border-none tw-cursor-pointer tw-outline-none focus-visible:tw-ring-2 focus-visible:tw-ring-ring"
    @click="uiStore.toggleTheme"
    :title="uiStore.theme === 'dark' ? 'Light Mode' : 'Dark Mode'"
  >
    <!-- Sun icon -->
    <svg v-if="uiStore.theme === 'dark'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="5"></circle><line x1="12" y1="1" x2="12" y2="3"></line><line x1="12" y1="21" x2="12" y2="23"></line><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line><line x1="1" y1="12" x2="3" y2="12"></line><line x1="21" y1="12" x2="23" y2="12"></line><line x1="4.22" y1="19.77" x2="5.64" y2="18.36"></line><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line></svg>
    <!-- Moon icon -->
    <svg v-else xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path></svg>
  </button>
</template>
