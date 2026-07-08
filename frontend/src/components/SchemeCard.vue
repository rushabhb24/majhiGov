<script setup>
import { computed } from 'vue'
import AppButton from './ui/AppButton.vue'
import AppBadge from './ui/AppBadge.vue'
import AppCard from './ui/AppCard.vue'

const props = defineProps({
  scheme: {
    type: Object,
    required: true
  },
  currentLanguage: {
    type: String,
    required: true
  },
  savedSchemeIds: {
    type: Array,
    required: true
  },
  t: {
    type: Object,
    required: true
  },
  isLoggedIn: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['toggleBookmark', 'openDetails', 'loginRequired', 'applyClick']);

function getSchemeTitle(scheme) {
  if (props.currentLanguage === 'mr' && scheme.title_mr) return scheme.title_mr;
  if (props.currentLanguage === 'hi' && scheme.title_hi) return scheme.title_hi;
  return scheme.title;
}

function getSchemeDesc(scheme) {
  if (props.currentLanguage === 'mr' && scheme.description_mr) return scheme.description_mr;
  if (props.currentLanguage === 'hi' && scheme.description_hi) return scheme.description_hi;
  return scheme.description;
}

function getCategoryName(scheme) {
  const cat = props.currentLanguage === 'mr' ? scheme.category_name_mr : (props.currentLanguage === 'hi' ? scheme.category_name_hi : scheme.category_name);
  return cat || scheme.category_name;
}

function formatDate(dateStr) {
  if (!dateStr) return '31 Dec 2026'
  const d = new Date(dateStr)
  const day = d.getDate()
  const month = d.toLocaleString('en-US', { month: 'short' })
  const year = d.getFullYear()
  return `${day} ${month} ${year}`
}
</script>

<template>
  <AppCard hoverable class="tw-flex tw-flex-col tw-h-full">
    <!-- Header: Category + Bookmark -->
    <div class="tw-flex tw-justify-between tw-items-center tw-mb-4">
      <AppBadge tone="info">
        {{ getCategoryName(scheme) }}
      </AppBadge>
      
      <button 
        class="tw-bg-transparent tw-border-none tw-cursor-pointer tw-p-1 tw-flex tw-items-center tw-justify-center tw-text-slate-400 dark:tw-text-slate-300 hover:tw-text-primary hover:tw-scale-110 tw-transition-all tw-outline-none"
        @click.stop="emit('toggleBookmark', scheme.id)"
        :title="savedSchemeIds.includes(scheme.id) ? 'Remove Bookmark' : 'Save Scheme'"
      >
        <svg 
          xmlns="http://www.w3.org/2000/svg" 
          width="18" 
          height="18" 
          viewBox="0 0 24 24" 
          :fill="savedSchemeIds.includes(scheme.id) ? 'currentColor' : 'none'" 
          stroke="currentColor" 
          stroke-width="2" 
          stroke-linecap="round" 
          stroke-linejoin="round"
        >
          <path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"></path>
        </svg>
      </button>
    </div>

    <!-- Title and Description -->
    <h3 class="tw-font-heading tw-font-bold tw-text-base tw-text-foreground tw-mb-2 tw-leading-snug tw-m-0">
      {{ getSchemeTitle(scheme) }}
    </h3>
    <p class="tw-text-xs tw-text-muted-foreground tw-leading-relaxed tw-mb-4 tw-line-clamp-3">
      {{ getSchemeDesc(scheme) }}
    </p>

    <!-- Green Benefit Banner -->
    <div class="tw-bg-emerald-500/10 tw-border tw-border-emerald-500/20 tw-rounded-xl tw-p-3 tw-mb-4 tw-flex tw-gap-2 tw-items-start tw-text-xs tw-leading-snug">
      <span class="tw-text-emerald-600 dark:tw-text-emerald-400 tw-font-black">Benefit:</span>
      <span class="tw-text-emerald-700 dark:tw-text-emerald-300 tw-font-semibold">{{ scheme.benefits }}</span>
    </div>

    <!-- Calendar Apply Date -->
    <div class="tw-flex tw-items-center tw-gap-2 tw-text-xs tw-text-muted-foreground tw-mt-auto tw-mb-4">
      <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
      <span class="tw-font-medium">Apply by {{ formatDate(scheme.application_end_date) }}</span>
    </div>

    <!-- Footer Action Buttons -->
    <div class="tw-flex tw-gap-3 tw-w-full">
      <AppButton variant="outline" size="sm" class="tw-flex-1" @click="emit('openDetails', scheme)">
        {{ t.viewDetails }}
      </AppButton>
      <AppButton variant="primary" size="sm" class="tw-flex-1 tw-flex tw-items-center tw-justify-center tw-gap-1.5" @click.stop="emit('applyClick', scheme)">
        Apply
        <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path><polyline points="15 3 21 3 21 9"></polyline><line x1="10" y1="14" x2="21" y2="3"></line></svg>
      </AppButton>
    </div>
  </AppCard>
</template>
