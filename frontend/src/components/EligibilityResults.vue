<script setup>
import AppCard from './ui/AppCard.vue'
import AppBadge from './ui/AppBadge.vue'
import AppButton from './ui/AppButton.vue'

const props = defineProps({
  checked: {
    type: Boolean,
    required: true
  },
  results: {
    type: Object,
    default: null
  },
  currentLanguage: {
    type: String,
    required: true
  },
  t: {
    type: Object,
    required: true
  }
});

const emit = defineEmits(['openDetails']);

function getSchemeTitle(scheme) {
  if (props.currentLanguage === 'mr' && scheme.title_mr) return scheme.title_mr;
  if (props.currentLanguage === 'hi' && scheme.title_hi) return scheme.title_hi;
  return scheme.title;
}

function filterDisqualifyingReasons(reasons) {
  return reasons.filter(r => 
    r.includes('not') || 
    r.includes('exceeds') || 
    r.includes('outside') || 
    r.includes('specifically') || 
    r.includes('Required') || 
    r.includes('नाही') || // Marathi check for negatives
    r.includes('अपात्र') || 
    r.includes('अधिक') || 
    r.includes('नहीं') || // Hindi check for negatives
    r.includes('बाहर')
  );
}
</script>

<template>
  <AppCard class="tw-max-w-2xl tw-mx-auto tw-p-8 tw-mt-6">
    <h2 class="tw-font-heading tw-font-bold tw-text-xl tw-text-foreground tw-mb-6 tw-m-0 text-center">
      {{ t.resultsTitle || 'Eligibility Results' }}
    </h2>
    
    <!-- Intro State (Not Checked) -->
    <div v-if="!checked" class="tw-text-center tw-py-12 tw-flex tw-flex-col tw-items-center tw-gap-4">
      <div class="tw-text-5xl">🛡️</div>
      <h3 class="tw-font-heading tw-font-bold tw-text-base tw-text-foreground tw-m-0">
        {{ t.notCheckedIntro || 'Check Your Entitlements' }}
      </h3>
      <p class="tw-text-sm tw-text-muted-foreground tw-m-0">
        {{ t.noDetailsChecked || 'Fill out the wizard demographic questions to view schemes matching your background.' }}
      </p>
    </div>

    <!-- Results State -->
    <div v-else-if="results" class="tw-flex tw-flex-col tw-gap-8">
      
      <!-- 1. Eligible Schemes -->
      <div class="tw-flex tw-flex-col tw-gap-4">
        <h3 class="tw-font-heading tw-font-bold tw-text-sm tw-text-success tw-flex tw-items-center tw-gap-2 tw-m-0">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline></svg>
          {{ t.eligibleTitle || 'Eligible Schemes' }} ({{ results.eligible ? results.eligible.length : 0 }})
        </h3>

        <div v-if="!results.eligible || results.eligible.length === 0" class="tw-p-4 tw-rounded-xl tw-bg-muted/50 tw-border tw-border-border tw-text-xs tw-text-muted-foreground tw-text-center">
          {{ t.noEligible || 'No schemes match your eligibility profile.' }}
        </div>

        <div v-else class="tw-flex tw-flex-col tw-gap-4">
          <div 
            v-for="status in results.eligible" 
            :key="status.scheme.id" 
            class="tw-p-5 tw-border tw-border-success/20 tw-bg-success/5 tw-rounded-2xl tw-flex tw-flex-col tw-gap-3"
          >
            <div class="tw-flex tw-justify-between tw-items-start">
              <h4 class="tw-font-heading tw-font-bold tw-text-sm tw-text-foreground tw-m-0">
                {{ getSchemeTitle(status.scheme) }}
              </h4>
              <AppBadge tone="success">Eligible</AppBadge>
            </div>
            
            <p class="tw-text-xs tw-text-success tw-font-semibold tw-m-0">
              Benefit: {{ status.scheme.benefits }}
            </p>
            
            <div class="tw-p-3 tw-bg-card tw-rounded-xl tw-border tw-border-border tw-text-xs">
              <strong class="tw-block tw-font-heading tw-text-foreground tw-mb-1.5">
                {{ t.reasonsLabel || 'Matching Criteria:' }}
              </strong>
              <ul class="tw-list-none tw-p-0 tw-m-0 tw-flex tw-flex-col tw-gap-1 tw-text-muted-foreground">
                <li v-for="(reason, index) in status.reasons" :key="index" class="tw-flex tw-items-start tw-gap-1.5">
                  <span class="tw-text-success">✓</span>
                  <span>{{ reason }}</span>
                </li>
              </ul>
            </div>

            <div class="tw-flex tw-gap-3">
              <AppButton variant="outline" size="sm" class="tw-flex-1" @click="emit('openDetails', status.scheme)">
                {{ t.viewDetails }}
              </AppButton>
              <a :href="status.scheme.apply_link" target="_blank" class="tw-flex-1 tw-no-underline" rel="noopener noreferrer">
                <AppButton variant="primary" size="sm" class="tw-w-full">
                  Apply Now
                </AppButton>
              </a>
            </div>
          </div>
        </div>
      </div>

      <!-- 2. Not Eligible Schemes -->
      <div class="tw-flex tw-flex-col tw-gap-4">
        <h3 class="tw-font-heading tw-font-bold tw-text-sm tw-text-destructive tw-flex tw-items-center tw-gap-2 tw-m-0">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="15" y1="9" x2="9" y2="15"></line><line x1="9" y1="9" x2="15" y2="15"></line></svg>
          {{ t.notEligibleTitle || 'Not Eligible Schemes' }} ({{ results.not_eligible ? results.not_eligible.length : 0 }})
        </h3>

        <div v-if="!results.not_eligible || results.not_eligible.length === 0" class="tw-p-4 tw-rounded-xl tw-bg-muted/50 tw-border tw-border-border tw-text-xs tw-text-muted-foreground tw-text-center">
          No disqualifying schemes.
        </div>

        <div v-else class="tw-flex tw-flex-col tw-gap-4">
          <div 
            v-for="status in results.not_eligible" 
            :key="status.scheme.id" 
            class="tw-p-5 tw-border tw-border-border tw-bg-muted/30 tw-rounded-2xl tw-flex tw-flex-col tw-gap-3 tw-opacity-80"
          >
            <div class="tw-flex tw-justify-between tw-items-start">
              <h4 class="tw-font-heading tw-font-bold tw-text-sm tw-text-foreground tw-m-0">
                {{ getSchemeTitle(status.scheme) }}
              </h4>
              <AppBadge tone="danger">Not Eligible</AppBadge>
            </div>
            
            <div class="tw-p-3 tw-bg-card tw-rounded-xl tw-border tw-border-border tw-text-xs">
              <strong class="tw-block tw-font-heading tw-text-foreground tw-mb-1.5">
                {{ t.notEligReasonsLabel || 'Disqualifying Reasons:' }}
              </strong>
              <ul class="tw-list-none tw-p-0 tw-m-0 tw-flex tw-flex-col tw-gap-1 tw-text-muted-foreground">
                <li v-for="(reason, index) in filterDisqualifyingReasons(status.reasons)" :key="index" class="tw-flex tw-items-start tw-gap-1.5">
                  <span class="tw-text-destructive">✗</span>
                  <span>{{ reason }}</span>
                </li>
              </ul>
            </div>

            <div class="tw-flex">
              <AppButton variant="outline" size="sm" class="tw-w-full" @click="emit('openDetails', status.scheme)">
                {{ t.viewDetails }}
              </AppButton>
            </div>
          </div>
        </div>
      </div>

    </div>
  </AppCard>
</template>
