<script setup>
import SchemeCard from './SchemeCard.vue';

const props = defineProps({
  schemes: {
    type: Array,
    required: true
  },
  recommendedSchemes: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    required: true
  },
  error: {
    type: String,
    default: null
  },
  currentLanguage: {
    type: String,
    required: true
  },
  savedSchemeIds: {
    type: Array,
    required: true
  },
  categories: {
    type: Array,
    required: true
  },
  selectedCategory: {
    type: String,
    required: true
  },
  sortBy: {
    type: String,
    required: true
  },
  searchQu: {
    type: String,
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

const emit = defineEmits([
  'update:selectedCategory',
  'update:sortBy',
  'update:searchQu',
  'toggleBookmark',
  'openDetails',
  'retry',
  'loginRequired',
  'applyClick'
]);

function getCategoryName(categoryName) {
  if (categoryName === 'Farmers') return props.t.farmerCategory;
  if (categoryName === 'Students') return props.t.studentCategory;
  if (categoryName === 'Women') return props.t.womenCategory;
  if (categoryName === 'Senior Citizens') return props.t.seniorCategory;
  if (categoryName === 'Business Owners') return props.t.businessCategory;
  if (categoryName === 'All') return props.t.allCategory;
  return categoryName;
}
</script>

<template>
  <div class="tab-content">
    <!-- Search and Filters Panel -->
    <div class="filter-panel card">
      <h2 class="section-title">{{ t.explorer }}</h2>
      
      <div class="search-filter-row">
        <!-- Search field -->
        <div class="form-group flex-grow">
          <label class="form-label" for="search-input">{{ t.searchLabel }}</label>
          <div class="search-input-wrapper">
            <svg class="search-icon" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
            <input 
              type="text" 
              id="search-input" 
              class="form-control padded-left" 
              :placeholder="t.searchPlaceholder" 
              :value="searchQu"
              @input="emit('update:searchQu', $event.target.value)"
            />
          </div>
        </div>

        <!-- Sort By dropdown -->
        <div class="form-group sort-field">
          <label class="form-label" for="sort-select">{{ t.sortByLabel }}</label>
          <select 
            id="sort-select" 
            class="form-control" 
            :value="sortBy"
            @change="emit('update:sortBy', $event.target.value)"
          >
            <option value="date_desc">{{ t.recentFirst }}</option>
            <option value="title_asc">{{ t.titleAlphabetical }}</option>
          </select>
        </div>
      </div>

      <!-- Category Chips -->
      <div class="category-wrapper">
        <span class="category-label">{{ t.chooseCategory }}</span>
        <div class="category-chips">
          <button 
            v-for="cat in categories" 
            :key="cat"
            :class="['chip', { active: selectedCategory === cat }]"
            @click="emit('update:selectedCategory', cat)"
          >
            {{ getCategoryName(cat) }}
          </button>
        </div>
      </div>
    </div>

    <!-- 🎯 Personalized Recommendations Carousel (Sleek Glassmorphism) -->
    <Transition name="fade-slide">
      <div v-if="isLoggedIn && recommendedSchemes.length > 0" class="recommendations-container card mt-4">
        <div class="recommendations-header">
          <h3 class="recommendations-title">
            <span class="pulse-sparkle">🎯</span>
            {{ t.recommendedForYou || 'Recommended Schemes for You' }}
          </h3>
          <span class="recommendations-count-badge">{{ recommendedSchemes.length }} Schemes matched</span>
        </div>
        <div class="recommendations-carousel mt-3">
          <div 
            v-for="scheme in recommendedSchemes" 
            :key="'rec-' + scheme.id" 
            class="rec-card card"
          >
            <span class="rec-badge">Best Match</span>
            <h4 class="rec-title">{{ currentLanguage === 'mr' ? (scheme.title_mr || scheme.title) : (currentLanguage === 'hi' ? (scheme.title_hi || scheme.title) : scheme.title) }}</h4>
            <p class="rec-benefits">{{ scheme.benefits }}</p>
            <div class="rec-footer mt-4">
              <button class="btn btn-secondary btn-sm" @click="emit('openDetails', scheme)">
                Learn More
              </button>
              <button class="btn btn-primary btn-sm" @click="emit('applyClick', scheme)">
                Apply
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Schemes Rendering States -->
    <div v-if="loading && schemes.length === 0" class="loading-state text-center mt-4 card">
      <div class="spinner"></div>
      <p class="mt-4">{{ t.loading }}</p>
    </div>

    <div v-else-if="error && schemes.length === 0" class="error-state text-center mt-4 card danger-card">
      <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="error-icon"><polygon points="7.86 2 16.14 2 22 7.86 22 16.14 16.14 22 7.86 22 2 16.14 2 7.86 7.86 2"></polygon><line x1="12" y1="9" x2="12" y2="13"></line><line x1="12" y1="17" x2="12.01" y2="17"></line></svg>
      <h3 class="mt-4">Database Connection Error</h3>
      <p>{{ error }}</p>
      <button class="btn btn-primary mt-4" @click="emit('retry')">{{ t.retry }}</button>
    </div>

    <div v-else-if="schemes.length === 0" class="empty-state text-center mt-4 card">
      <h3>Koi Scheme nahi mili! 🔍</h3>
      <p class="text-muted">Apna category ya search term badal kar try karein.</p>
    </div>

    <!-- Active schemes grid layout -->
    <div v-else class="schemes-grid mt-4">
      <SchemeCard 
        v-for="scheme in schemes" 
        :key="scheme.id"
        :scheme="scheme"
        :current-language="currentLanguage"
        :saved-scheme-ids="savedSchemeIds"
        :t="t"
        :is-logged-in="isLoggedIn"
        @toggle-bookmark="emit('toggleBookmark', $event)"
        @open-details="emit('openDetails', $event)"
        @login-required="emit('loginRequired', $event)"
        @apply-click="emit('applyClick', $event)"
      />
    </div>
  </div>
</template>

<style scoped>
.danger-card {
  border-color: var(--clr-danger);
  color: var(--clr-danger);
}
.error-icon {
  color: var(--clr-danger);
  margin: 0 auto;
}

/* Recommendations sliding carousel styling */
.recommendations-container {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.08) 0%, rgba(168, 85, 247, 0.05) 100%);
  border: 1px solid rgba(99, 102, 241, 0.2);
  padding: 20px;
  border-radius: var(--border-radius-md);
  overflow: hidden;
  position: relative;
  animation: slideInRec 0.5s cubic-bezier(0.16, 1, 0.3, 1) both;
}

@keyframes slideInRec {
  from { transform: translateY(12px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}

.recommendations-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.recommendations-title {
  font-family: var(--font-heading);
  font-size: 1.15rem;
  font-weight: 700;
  color: var(--clr-text-main);
  display: flex;
  align-items: center;
  gap: 8px;
}

.pulse-sparkle {
  animation: pulse-ring 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse-ring {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.15); }
}

.recommendations-count-badge {
  background: var(--clr-primary-light);
  color: var(--clr-primary);
  font-size: 0.75rem;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: var(--border-radius-full);
  border: 1px solid rgba(99, 102, 241, 0.15);
}

.recommendations-carousel {
  display: flex;
  gap: 16px;
  overflow-x: auto;
  padding: 4px 4px 12px;
  scroll-behavior: smooth;
  -webkit-overflow-scrolling: touch;
}

/* Custom Scrollbar for carousel */
.recommendations-carousel::-webkit-scrollbar {
  height: 6px;
}
.recommendations-carousel::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.03);
  border-radius: var(--border-radius-full);
}
.recommendations-carousel::-webkit-scrollbar-thumb {
  background: rgba(99, 102, 241, 0.15);
  border-radius: var(--border-radius-full);
}
.recommendations-carousel::-webkit-scrollbar-thumb:hover {
  background: rgba(99, 102, 241, 0.3);
}

.rec-card {
  flex: 0 0 280px;
  background: var(--clr-surface);
  border: 1px solid var(--clr-border);
  padding: 16px;
  border-radius: var(--border-radius-md);
  box-shadow: 0 4px 12px rgba(0,0,0,0.02);
  transition: all var(--transition-normal);
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.rec-card:hover {
  transform: translateY(-3px);
  border-color: var(--clr-primary);
  box-shadow: 0 6px 16px var(--clr-primary-light);
}

.rec-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  background: linear-gradient(135deg, var(--clr-secondary) 0%, #10b981 100%);
  color: #fff;
  font-size: 0.68rem;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: var(--border-radius-full);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.rec-title {
  font-size: 0.92rem;
  font-weight: 700;
  color: var(--clr-text-main);
  line-height: 1.3;
  margin-top: 12px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  height: 2.6em;
}

.rec-benefits {
  font-size: 0.8rem;
  color: var(--clr-secondary);
  font-weight: 600;
  margin-top: 8px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  height: 2.4em;
  background: rgba(16, 185, 129, 0.05);
  padding: 4px 8px;
  border-radius: var(--border-radius-sm);
}

.rec-footer {
  display: flex;
  gap: 8px;
  width: 100%;
}

.rec-footer .btn-sm {
  flex: 1;
  padding: 6px 12px;
  font-size: 0.78rem;
  border-radius: var(--border-radius-sm);
}
</style>
