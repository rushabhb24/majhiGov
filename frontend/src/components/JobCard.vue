<script setup>
import { computed } from 'vue'
import AppCard from './ui/AppCard.vue'
import AppBadge from './ui/AppBadge.vue'
import AppButton from './ui/AppButton.vue'

const props = defineProps({
  job:             { type: Object, required: true },
  isPrivate:       { type: Boolean, default: false },
  isBookmarked:    { type: Boolean, default: false },
  isLoggedIn:      { type: Boolean, default: false },
  currentLanguage: { type: String, default: 'en' }
})

defineEmits(['applyClick', 'viewDetails', 'toggleBookmark', 'loginRequired'])

// Format currency range
const formattedSalary = computed(() => {
  if (!props.isPrivate) {
    return props.job.application_fee || 'Free / Exempted'
  }
  if (props.job.salary_min > 0 && props.job.salary_max > 0) {
    const minLPA = (props.job.salary_min / 100000).toFixed(1)
    const maxLPA = (props.job.salary_max / 100000).toFixed(1)
    return `₹${minLPA} - ₹${maxLPA} LPA`
  }
  if (props.job.stipend) {
    return props.job.stipend
  }
  if (props.job.prize_pool) {
    return `Prize: ${props.job.prize_pool}`
  }
  return 'Negotiable'
})

const getInitials = (name) => {
  if (!name) return 'C'
  const parts = name.split(' ')
  if (parts.length > 1) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return name.substring(0, 2).toUpperCase()
}

// Generate unique HSL colors for company initial avatars
const getAvatarColor = (name) => {
  let hash = 0
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash)
  }
  const h = Math.abs(hash % 360)
  return `hsl(${h}, 70%, 45%)`
}

const getDaysRemaining = (endDateStr) => {
  if (!endDateStr) return 99
  const end = new Date(endDateStr)
  const today = new Date()
  const diffTime = end - today
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  return diffDays
}
</script>

<template>
  <AppCard hoverable class="tw-flex tw-flex-col tw-h-full">
    
    <!-- Top row: Logo + Bookmarks -->
    <div class="tw-flex tw-justify-between tw-items-start tw-mb-4">
      <div class="tw-flex tw-items-center tw-gap-3">
        <!-- Logo circle -->
        <div v-if="job.company_logo_url" class="tw-w-12 tw-h-12 tw-rounded-xl tw-bg-white tw-border tw-border-border tw-flex tw-items-center tw-justify-center tw-overflow-hidden tw-flex-shrink-0">
          <img :src="job.company_logo_url" :alt="job.company_name || job.organization" class="tw-max-w-[85%] tw-max-h-[85%] tw-object-contain" />
        </div>
        <div v-else class="tw-w-12 tw-h-12 tw-rounded-xl tw-flex tw-items-center tw-justify-center tw-font-heading tw-font-black tw-text-base tw-text-white tw-flex-shrink-0" :style="{ backgroundColor: getAvatarColor(job.company_name || job.organization) }">
          {{ getInitials(job.company_name || job.organization) }}
        </div>

        <div class="tw-flex tw-flex-col">
          <h4 class="tw-font-heading tw-font-bold tw-text-sm tw-text-foreground tw-m-0 tw-line-clamp-1">
            {{ job.company_name || job.organization }}
          </h4>
          <div class="tw-mt-1">
            <AppBadge :tone="isPrivate ? 'info' : 'success'">
              {{ isPrivate ? (job.job_type || 'Private') : 'Govt Job' }}
            </AppBadge>
          </div>
        </div>
      </div>

      <button 
        class="tw-bg-transparent tw-border-none tw-cursor-pointer tw-p-1 tw-text-muted-foreground hover:tw-text-primary hover:tw-scale-110 tw-transition-all tw-outline-none"
        @click="$emit('toggleBookmark', job)"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" :fill="isBookmarked ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"/>
        </svg>
      </button>
    </div>

    <!-- Title and Tags -->
    <div class="tw-mb-4">
      <h3 class="tw-font-heading tw-font-bold tw-text-base tw-text-foreground tw-m-0 tw-leading-snug tw-line-clamp-2">
        {{ job.title }}
      </h3>
      <div class="tw-flex tw-flex-wrap tw-gap-1.5 tw-mt-2.5">
        <span v-if="job.work_mode" class="tw-bg-muted tw-text-muted-foreground tw-text-[10px] tw-px-2 tw-py-0.5 tw-rounded-md tw-font-semibold">{{ job.work_mode }}</span>
        <span v-if="job.employment_type" class="tw-bg-muted tw-text-muted-foreground tw-text-[10px] tw-px-2 tw-py-0.5 tw-rounded-md tw-font-semibold">{{ job.employment_type }}</span>
        <span class="tw-bg-muted tw-text-muted-foreground tw-text-[10px] tw-px-2 tw-py-0.5 tw-rounded-md tw-font-semibold">{{ job.education_qualification }}</span>
      </div>
    </div>

    <!-- Metadata Panel -->
    <div class="tw-grid tw-grid-cols-3 tw-gap-2 tw-bg-muted/40 tw-p-3 tw-rounded-xl tw-mb-4 tw-text-xs tw-text-muted-foreground">
      <div class="tw-flex tw-items-center tw-gap-1">
        <span>📍</span>
        <span class="tw-truncate">{{ job.location || 'India' }}</span>
      </div>
      <div class="tw-flex tw-items-center tw-gap-1">
        <span>💰</span>
        <span class="tw-truncate">{{ formattedSalary }}</span>
      </div>
      <div v-if="job.experience_required || job.experience_min !== undefined" class="tw-flex tw-items-center tw-gap-1">
        <span>⏱</span>
        <span class="tw-truncate">{{ isPrivate ? `${job.experience_min}-${job.experience_max} Yrs` : job.experience_required }}</span>
      </div>
    </div>

    <!-- AI Match Score -->
    <div v-if="job.ai_match_score || job.match_score" class="tw-mb-4 tw-flex tw-flex-col tw-gap-1.5">
      <div class="tw-flex tw-justify-between tw-text-xs tw-text-muted-foreground tw-font-semibold">
        <span>🤖 AI Profile Match</span>
        <span class="tw-text-primary">{{ job.ai_match_score || job.match_score }}%</span>
      </div>
      <div class="tw-h-1.5 tw-bg-muted tw-rounded-full tw-overflow-hidden">
        <div class="tw-h-full tw-bg-primary tw-rounded-full tw-transition-all" :style="{ width: (job.ai_match_score || job.match_score) + '%' }"></div>
      </div>
    </div>

    <!-- Footer Deadline & Actions -->
    <div class="tw-flex tw-items-center tw-justify-between tw-mt-auto tw-pt-3 tw-border-t tw-border-border/50">
      <div class="tw-flex tw-flex-col">
        <span class="tw-text-[10px] tw-text-muted-foreground">Apply by: {{ job.application_end_date }}</span>
        <span v-if="getDaysRemaining(job.application_end_date) <= 7" class="tw-text-[9px] tw-font-extrabold tw-text-destructive tw-uppercase tw-mt-0.5">
          Closing Soon
        </span>
      </div>
      <div class="tw-flex tw-gap-2">
        <AppButton variant="ghost" size="sm" @click="$emit('viewDetails', job)">
          Details
        </AppButton>
        <AppButton variant="primary" size="sm" @click="$emit('applyClick', job)">
          Apply →
        </AppButton>
      </div>
    </div>

  </AppCard>
</template>
