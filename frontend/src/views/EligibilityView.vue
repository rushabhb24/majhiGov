<script setup>
import { computed } from 'vue'
import { useEligibilityStore } from '../stores/eligibility'
import { useSchemeStore } from '../stores/schemes'
import { useUiStore } from '../stores/ui'
import EligibilityWizard from '../components/EligibilityWizard.vue'
import EligibilityResults from '../components/EligibilityResults.vue'
import { useI18n } from 'vue-i18n'

const eligibilityStore = useEligibilityStore()
const schemeStore = useSchemeStore()
const uiStore = useUiStore()
const { t, locale, messages } = useI18n()

// Translation object for child components that use t.key property access
const tObj = computed(() => messages.value[locale.value] || {})
</script>

<template>
  <div class="tab-content">
    <div class="grid-layout">
      <EligibilityWizard
        v-model:step="eligibilityStore.step"
        :profile="eligibilityStore.profile"
        :t="tObj"
        :checking="eligibilityStore.checking"
        @submit="eligibilityStore.submitEligibility"
      />

      <EligibilityResults
        :checked="eligibilityStore.checked"
        :results="eligibilityStore.results"
        :current-language="uiStore.currentLanguage"
        :t="tObj"
        @open-details="schemeStore.openDetails"
      />
    </div>
  </div>
</template>
