<script setup>
import SchemeCard from './SchemeCard.vue';

const props = defineProps({
  schemes: {
    type: Array,
    required: true
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
</style>
