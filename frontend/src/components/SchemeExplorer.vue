<script setup>
import SchemeCard from './SchemeCard.vue';
import SkeletonCard from './SkeletonCard.vue';
import EmptyState from './EmptyState.vue';

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
  loadingMore: {
    type: Boolean,
    default: false
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

// Map visual icons and custom classes for category grid representation
const categoryIcons = {
  'All': 'IN',
  'Farmers': '🌾',
  'Students': '🎓',
  'Women': '👩',
  'Senior Citizens': '👵',
  'Health': '🏥',
  'Housing': '🏠',
  'Business Owners': '💼',
  'Business': '💼',
  'Divyang': '♿'
};

function getCategoryLabel(cat) {
  if (cat === 'All') return 'All';
  if (cat === 'Farmers') return 'Farmers';
  if (cat === 'Students') return 'Students';
  if (cat === 'Women') return 'Women';
  if (cat === 'Senior Citizens') return 'Senior Citizens';
  if (cat === 'Business Owners') return 'Business';
  if (cat === 'Business') return 'Business';
  if (cat === 'Health') return 'Health';
  if (cat === 'Housing') return 'Housing';
  if (cat === 'Divyang') return 'Divyang';
  return cat;
}
</script>

<template>
  <div id="schemes-section" class="tw-max-w-7xl tw-mx-auto tw-px-4 tw-sm:px-6 tw-lg:px-8 tw-py-8">
    
    <!-- Top Search & Sort Action Container -->
    <div class="tw-flex tw-flex-col sm:tw-flex-row tw-gap-4 tw-items-center tw-justify-between">
      
      <!-- Search Input Box -->
      <div class="glass tw-flex-1 tw-flex tw-items-center tw-px-4 tw-py-2 tw-rounded-full tw-w-full">
        <svg class="tw-flex-shrink-0 tw-mr-3" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
        <input 
          type="text" 
          class="tw-border-none tw-bg-transparent tw-outline-none tw-text-sm tw-text-foreground tw-w-full tw-font-body" 
          placeholder="Search schemes by name or keyword..." 
          :value="searchQu"
          @input="emit('update:searchQu', $event.target.value)"
        />
      </div>

      <!-- Sort Toggle Buttons -->
      <div class="tw-flex tw-bg-muted/60 tw-p-1 tw-rounded-full tw-flex-shrink-0">
        <button 
          class="tw-border-none tw-outline-none tw-text-xs tw-font-bold tw-font-heading tw-px-4 tw-py-2 tw-rounded-full tw-cursor-pointer tw-transition-all"
          :class="sortBy === 'date_desc' ? 'tw-bg-card tw-text-foreground tw-shadow-sm' : 'tw-text-muted-foreground tw-bg-transparent hover:tw-text-foreground'"
          @click="emit('update:sortBy', 'date_desc')"
        >
          Recent
        </button>
        <button 
          class="tw-border-none tw-outline-none tw-text-xs tw-font-bold tw-font-heading tw-px-4 tw-py-2 tw-rounded-full tw-cursor-pointer tw-transition-all"
          :class="sortBy === 'title_asc' ? 'tw-bg-card tw-text-foreground tw-shadow-sm' : 'tw-text-muted-foreground tw-bg-transparent hover:tw-text-foreground'"
          @click="emit('update:sortBy', 'title_asc')"
        >
          A-Z
        </button>
      </div>

    </div>

    <!-- Category Rounded Avatars Scroller Strip -->
    <div class="tw-flex tw-gap-4 tw-overflow-x-auto tw-pb-3 tw-mt-6 tw-scrollbar-none">
      <div 
        v-for="cat in ['All', 'Farmers', 'Students', 'Women', 'Senior Citizens', 'Health', 'Housing', 'Business', 'Divyang']"
        :key="cat"
        class="glass tw-flex-shrink-0 tw-w-24 tw-flex tw-flex-col tw-items-center tw-gap-2.5 tw-cursor-pointer tw-rounded-2xl tw-py-3 tw-px-1 tw-transition-all tw-duration-250"
        :class="selectedCategory === cat || (cat === 'Business' && selectedCategory === 'Business Owners') ? 'tw-border-primary/80 tw-bg-primary/5 tw-scale-105' : 'hover:tw-border-primary/30'"
        @click="emit('update:selectedCategory', cat === 'Business' ? 'Business Owners' : cat)"
      >
        <div 
          class="tw-w-12 tw-h-12 tw-rounded-full tw-flex tw-items-center tw-justify-center tw-transition-all"
          :class="selectedCategory === cat || (cat === 'Business' && selectedCategory === 'Business Owners') ? 'tw-bg-primary tw-text-white' : 'tw-bg-muted'"
        >
          <span v-if="cat === 'All'" class="tw-font-heading tw-font-extrabold tw-text-xs" :class="selectedCategory === cat ? 'tw-text-white' : 'tw-text-primary'">IN</span>
          <span v-else class="tw-text-xl">{{ categoryIcons[cat] || '🌾' }}</span>
        </div>
        <span 
          class="tw-text-[10px] tw-font-bold tw-text-center tw-truncate tw-w-full tw-font-heading"
          :class="selectedCategory === cat || (cat === 'Business' && selectedCategory === 'Business Owners') ? 'tw-text-primary' : 'tw-text-muted-foreground'"
        >
          {{ getCategoryLabel(cat) }}
        </span>
      </div>
    </div>

    <!-- Loading Shimmer State -->
    <div v-if="loading && schemes.length === 0" class="tw-grid tw-grid-cols-1 md:tw-grid-cols-3 tw-gap-6 tw-mt-6">
      <SkeletonCard v-for="i in 6" :key="i" />
    </div>

    <!-- Error State -->
    <div v-else-if="error && schemes.length === 0" class="glass tw-p-8 tw-rounded-2xl tw-text-center tw-mt-6">
      <h3 class="tw-font-heading tw-font-bold tw-text-lg tw-text-foreground tw-m-0">Database Connection Error</h3>
      <p class="tw-text-sm tw-text-muted-foreground tw-mt-2">{{ error }}</p>
      <AppButton variant="primary" size="sm" class="tw-mt-4" @click="emit('retry')">
        Retry
      </AppButton>
    </div>

    <!-- Empty State -->
    <EmptyState
      v-else-if="schemes.length === 0"
      title="No Schemes Found"
      description="Try adjusting your filters or search terms."
      icon="search"
      actionLabel="Reset Search"
      @action="emit('update:searchQu', ''); emit('update:selectedCategory', 'All')"
    />

    <!-- Schemes Grid List -->
    <div v-else class="tw-grid tw-grid-cols-1 md:tw-grid-cols-3 tw-gap-6 tw-mt-6">
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

    <!-- Loading More Shimmer Skeletons -->
    <div v-if="loadingMore" class="tw-grid tw-grid-cols-1 md:tw-grid-cols-3 tw-gap-6 tw-mt-6">
      <SkeletonCard v-for="i in 3" :key="'more-' + i" />
    </div>

  </div>
</template>
