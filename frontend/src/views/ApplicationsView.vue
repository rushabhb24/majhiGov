<script setup>
import { onMounted } from 'vue'
import { useApplicationStore } from '../stores/applications'
import { useUiStore } from '../stores/ui'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import AppCard from '../components/ui/AppCard.vue'
import AppBadge from '../components/ui/AppBadge.vue'
import AppButton from '../components/ui/AppButton.vue'

const applicationStore = useApplicationStore()
const uiStore = useUiStore()
const { t } = useI18n()
const router = useRouter()

function getStatusIcon(status) {
  if (status === 'approved') return '✅'
  if (status === 'rejected') return '❌'
  return '⏳'
}

function getStepClass(appStatus, stepStatus) {
  const step = applicationStore.getStatusStep(appStatus)
  const stepNum = stepStatus === 'submitted' ? 1 : stepStatus === 'review' ? 2 : 3
  
  if (stepNum < step) return 'step-completed'
  if (stepNum === step) {
    if (appStatus === 'approved') return 'step-completed step-success'
    if (appStatus === 'rejected') return 'step-completed step-danger'
    return 'step-active'
  }
  return 'step-pending'
}

function getLineClass(appStatus, afterStep) {
  const step = applicationStore.getStatusStep(appStatus)
  return afterStep < step ? 'line-filled' : 'line-empty'
}

onMounted(() => {
  applicationStore.fetchApplications()
})
</script>

<template>
  <div class="tw-max-w-4xl tw-mx-auto tw-px-4 tw-sm:px-6 tw-lg:px-8 tw-py-8">
    
    <!-- Header panel -->
    <div class="glass tw-p-6 tw-rounded-2xl tw-mb-6">
      <div class="tw-flex tw-justify-between tw-items-center tw-flex-wrap tw-gap-4">
        <div>
          <h2 class="tw-font-heading tw-font-bold tw-text-xl tw-text-foreground tw-m-0">
            {{ t('myApplications') || 'My Applications' }}
          </h2>
          <p class="tw-text-xs tw-text-muted-foreground tw-mt-1 tw-m-0">
            {{ t('applicationTimeline') || 'Track stages and real-time approval status for submitted schemes.' }}
          </p>
        </div>
        <AppButton 
          variant="outline" 
          size="sm"
          class="tw-flex tw-items-center tw-gap-1.5"
          @click="applicationStore.refreshApplications()"
          :disabled="applicationStore.refreshing"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" :class="{ 'tw-animate-spin': applicationStore.refreshing }"><polyline points="23 4 23 10 17 10"></polyline><polyline points="1 20 1 14 7 14"></polyline><path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"></path></svg>
          <span>{{ applicationStore.refreshing ? t('submitting') : (t('refreshStatus') || 'Refresh') }}</span>
        </AppButton>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="applicationStore.applications.length === 0" class="glass tw-p-12 tw-rounded-2xl tw-text-center tw-flex tw-flex-col tw-items-center tw-gap-4">
      <div class="tw-text-5xl">📬</div>
      <h3 class="tw-font-heading tw-font-bold tw-text-base tw-text-foreground tw-m-0">
        {{ t('noApplications') || 'No Applications Found' }}
      </h3>
      <p class="tw-text-sm tw-text-muted-foreground tw-m-0">
        {{ t('noApplicationsHint') || 'Find a scheme and apply through the official portal' }}
      </p>
      <AppButton variant="primary" size="sm" @click="router.push('/')">
        {{ t('exploreSchemes') || 'Explore Schemes' }}
      </AppButton>
    </div>

    <!-- Application Tracking Cards -->
    <div v-else class="tw-flex tw-flex-col tw-gap-6">
      <AppCard 
        v-for="app in applicationStore.applications" 
        :key="app.id" 
        class="tw-flex tw-flex-col tw-gap-4"
      >
        <!-- Card Header -->
        <div class="tw-flex tw-justify-between tw-items-start tw-flex-wrap tw-gap-4">
          <div class="tw-flex-1 tw-min-w-0">
            <AppBadge tone="info">
              {{ app.government_level === 'central' ? '🏛️ Central' : '🏢 State' }}
            </AppBadge>
            <h3 class="tw-font-heading tw-font-bold tw-text-base tw-text-foreground tw-my-2 tw-m-0">
              {{ uiStore.currentLanguage === 'mr' ? (app.title_mr || app.title) : (uiStore.currentLanguage === 'hi' ? (app.title_hi || app.title) : app.title) }}
            </h3>
            <div class="tw-text-[11px] tw-text-muted-foreground">
              📅 {{ t('appliedOn') || 'Applied on' }}: {{ new Date(app.applied_at).toLocaleDateString(uiStore.currentLanguage === 'en' ? 'en-US' : (uiStore.currentLanguage === 'hi' ? 'hi-IN' : 'mr-IN'), { year: 'numeric', month: 'long', day: 'numeric' }) }}
            </div>
          </div>
          <div class="tw-flex tw-items-center tw-gap-2">
            <span class="tw-text-lg">{{ getStatusIcon(app.status) }}</span>
            <AppBadge 
              :tone="app.status === 'pending' ? 'warn' : (app.status === 'approved' ? 'success' : 'danger')"
            >
              {{ app.status === 'pending' ? t('statusPending') : (app.status === 'approved' ? t('statusApproved') : t('statusRejected')) }}
            </AppBadge>
          </div>
        </div>

        <!-- Status Progress Steps Stepper -->
        <div class="tw-flex tw-items-center tw-gap-1 tw-my-4 tw-px-2">
          
          <!-- Step 1: Submitted -->
          <div class="tw-flex tw-flex-col tw-items-center tw-gap-1.5" :class="getStepClass(app.status, 'submitted')">
            <div class="step-indicator tw-w-6 tw-h-6 tw-rounded-full tw-bg-muted tw-flex tw-items-center tw-justify-center tw-text-[10px] tw-font-bold">
              1
            </div>
            <span class="tw-text-[10px] tw-font-bold tw-text-muted-foreground">{{ t('stepSubmitted') || 'Submitted' }}</span>
          </div>

          <div class="tw-flex-1 tw-h-[2px] tw-mb-4 tw-rounded-full" :class="getLineClass(app.status, 1) === 'line-filled' ? 'tw-bg-primary' : 'tw-bg-border'"></div>

          <!-- Step 2: Under Review -->
          <div class="tw-flex tw-flex-col tw-items-center tw-gap-1.5" :class="getStepClass(app.status, 'review')">
            <div class="step-indicator tw-w-6 tw-h-6 tw-rounded-full tw-bg-muted tw-flex tw-items-center tw-justify-center tw-text-[10px] tw-font-bold">
              2
            </div>
            <span class="tw-text-[10px] tw-font-bold tw-text-muted-foreground">{{ t('stepReview') || 'Under Review' }}</span>
          </div>

          <div class="tw-flex-1 tw-h-[2px] tw-mb-4 tw-rounded-full" :class="getLineClass(app.status, 2) === 'line-filled' ? 'tw-bg-primary' : 'tw-bg-border'"></div>

          <!-- Step 3: Approved / Rejected -->
          <div class="tw-flex tw-flex-col tw-items-center tw-gap-1.5" :class="getStepClass(app.status, 'decision')">
            <div class="step-indicator tw-w-6 tw-h-6 tw-rounded-full tw-bg-muted tw-flex tw-items-center tw-justify-center tw-text-[10px] tw-font-bold">
              3
            </div>
            <span class="tw-text-[10px] tw-font-bold tw-text-muted-foreground">{{ app.status === 'rejected' ? (t('stepRejected') || 'Rejected') : (t('stepApproved') || 'Approved') }}</span>
          </div>

        </div>

        <!-- Citizen notes -->
        <div v-if="app.notes" class="tw-bg-muted/40 tw-border-l-4 tw-border-primary tw-p-3 tw-rounded-r-xl tw-text-xs tw-text-muted-foreground">
          <strong class="tw-text-foreground">{{ t('notesTitle') || 'Notes:' }}</strong> {{ app.notes }}
        </div>

        <!-- Re-visit official portal link -->
        <div v-if="app.apply_link || app.official_website" class="tw-flex tw-mt-2">
          <a 
            :href="app.apply_link || app.official_website" 
            target="_blank" 
            rel="noopener noreferrer" 
            class="tw-no-underline"
          >
            <AppButton variant="outline" size="sm" class="tw-flex tw-items-center tw-gap-1.5">
              <span>{{ t('visitOfficialPortal') || 'Visit Official Portal' }}</span>
              <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path><polyline points="15 3 21 3 21 9"></polyline><line x1="10" y1="14" x2="21" y2="3"></line></svg>
            </AppButton>
          </a>
        </div>
      </AppCard>
    </div>

  </div>
</template>

<style scoped>
/* Stepper node overrides */
.step-completed .step-indicator {
  background-color: hsl(var(--primary));
  color: white;
  box-shadow: 0 0 0 3px hsl(var(--primary) / 0.15);
}
.step-active .step-indicator {
  background-color: hsl(var(--accent));
  color: white;
  box-shadow: 0 0 0 3px hsl(var(--accent) / 0.15);
  animation: pulse-step 2s infinite;
}
.step-completed.step-success .step-indicator {
  background-color: hsl(var(--success));
}
.step-completed.step-danger .step-indicator {
  background-color: hsl(var(--destructive));
}

@keyframes pulse-step {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.08); }
}
</style>
