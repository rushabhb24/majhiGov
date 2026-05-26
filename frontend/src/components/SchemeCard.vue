<script setup>
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

const emit = defineEmits(['toggleBookmark', 'openDetails', 'loginRequired']);

// Dynamic multi-language getters for schemes data
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
</script>

<template>
  <article class="card scheme-card">
    <div class="scheme-header">
      <!-- Multilingual Category Name -->
      <span class="scheme-category">
        {{ getCategoryName(scheme) }}
      </span>
      
      <!-- Government Level Badge (Central/State) -->
      <span :class="['level-badge', scheme.government_level]">
        {{ scheme.government_level === 'central' ? 'Central' : (scheme.state || 'State') }}
      </span>

      <!-- Bookmark Button -->
      <button 
        class="btn-bookmark" 
        @click.stop="emit('toggleBookmark', scheme.id)"
        :title="savedSchemeIds.includes(scheme.id) ? 'Remove Bookmark' : 'Save Scheme'"
      >
        <svg 
          xmlns="http://www.w3.org/2000/svg" 
          width="20" 
          height="20" 
          viewBox="0 0 24 24" 
          :fill="savedSchemeIds.includes(scheme.id) ? 'var(--clr-accent)' : 'none'" 
          :stroke="savedSchemeIds.includes(scheme.id) ? 'var(--clr-accent)' : 'currentColor'" 
          stroke-width="2" 
          stroke-linecap="round" 
          stroke-linejoin="round"
        >
          <path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"></path>
        </svg>
      </button>
    </div>
    
    <!-- Multilingual Title and Clamped Description -->
    <h3 class="scheme-title">{{ getSchemeTitle(scheme) }}</h3>
    <p class="scheme-desc text-muted">{{ getSchemeDesc(scheme) }}</p>
    
    <!-- Highlights info section -->
    <div class="scheme-info">
      <div class="info-item">
        <span class="info-label">{{ t.benefitsLabel }}:</span>
        <span class="info-value benefit-highlight">{{ scheme.benefits }}</span>
      </div>
      <div class="info-item">
        <span class="info-label">{{ t.lastDateLabel }}:</span>
        <span class="info-value date-tag">{{ scheme.application_end_date }}</span>
      </div>
    </div>

    <!-- Action Buttons with clean wrapping -->
    <div class="scheme-footer">
      <button class="btn btn-secondary flex-grow" @click="emit('openDetails', scheme)">
        {{ t.viewDetails }}
      </button>
      <a 
        v-if="isLoggedIn"
        :href="scheme.apply_link" 
        target="_blank" 
        class="btn btn-primary" 
        rel="noopener noreferrer"
      >
        {{ t.applyLink }} 
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="5" y1="12" x2="19" y2="12"></line><polyline points="12 5 19 12 12 19"></polyline></svg>
      </a>
      <button 
        v-else 
        class="btn btn-primary"
        @click.stop="emit('loginRequired', 'apply')"
      >
        🔒 {{ t.loginToApply || 'Login to Apply' }}
      </button>
    </div>
  </article>
</template>
