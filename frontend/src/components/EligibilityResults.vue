<script setup>
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
  <div class="card results-panel">
    <h2 class="section-title">{{ t.resultsTitle }}</h2>
    
    <!-- Intro State (Not Checked) -->
    <div v-if="!checked" class="results-intro text-center">
      <div class="intro-art">🛡️</div>
      <h3>{{ t.notCheckedIntro }}</h3>
      <p class="text-muted">{{ t.noDetailsChecked }}</p>
    </div>

    <!-- Results State -->
    <div v-else-if="results" class="results-content">
      
      <!-- 1. Eligible Schemes -->
      <div class="results-section">
        <h3 class="results-sub-title success-text">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline></svg>
          {{ t.eligibleTitle }} - {{ results.eligible ? results.eligible.length : 0 }}
        </h3>

        <div v-if="!results.eligible || results.eligible.length === 0" class="empty-results-box">
          {{ t.noEligible }}
        </div>

        <div v-else class="results-list">
          <div 
            v-for="status in results.eligible" 
            :key="status.scheme.id" 
            class="result-item eligible-border"
          >
            <div class="result-header">
              <h4>{{ getSchemeTitle(status.scheme) }}</h4>
              <span class="badge badge-success">Eligible</span>
            </div>
            <p class="result-benefits benefit-highlight mt-4">{{ status.scheme.benefits }}</p>
            
            <div class="reasons-box mt-4">
              <strong>{{ t.reasonsLabel }}</strong>
              <ul>
                <li v-for="(reason, index) in status.reasons" :key="index">
                  ✓ {{ reason }}
                </li>
              </ul>
            </div>

            <div class="result-actions mt-4">
              <button class="btn btn-secondary" @click="emit('openDetails', status.scheme)">
                {{ t.viewDetails }}
              </button>
              <a :href="status.scheme.apply_link" target="_blank" class="btn btn-primary" rel="noopener noreferrer">
                Apply Now
              </a>
            </div>
          </div>
        </div>
      </div>

      <!-- 2. Not Eligible Schemes -->
      <div class="results-section mt-4">
        <h3 class="results-sub-title danger-text">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="15" y1="9" x2="9" y2="15"></line><line x1="9" y1="9" x2="15" y2="15"></line></svg>
          {{ t.notEligibleTitle }} - {{ results.not_eligible ? results.not_eligible.length : 0 }}
        </h3>

        <div v-if="!results.not_eligible || results.not_eligible.length === 0" class="empty-results-box">
          No disqualifying schemes.
        </div>

        <div v-else class="results-list">
          <div 
            v-for="status in results.not_eligible" 
            :key="status.scheme.id" 
            class="result-item not-eligible-border muted-card"
          >
            <div class="result-header">
              <h4>{{ getSchemeTitle(status.scheme) }}</h4>
              <span class="badge badge-danger">Not Eligible</span>
            </div>
            
            <div class="reasons-box mt-4">
              <strong>{{ t.notEligReasonsLabel }}</strong>
              <ul>
                <li v-for="(reason, index) in filterDisqualifyingReasons(status.reasons)" :key="index" class="text-danger">
                  ✗ {{ reason }}
                </li>
              </ul>
            </div>

            <div class="result-actions mt-4">
              <button class="btn btn-secondary" @click="emit('openDetails', status.scheme)">
                {{ t.viewDetails }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
