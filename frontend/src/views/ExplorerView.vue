<script setup>
import { computed, watch, onMounted, onUnmounted, ref } from 'vue'
import { useSchemeStore } from '../stores/schemes'
import { useBookmarkStore } from '../stores/bookmarks'
import { useAuthStore } from '../stores/auth'
import { useApplicationStore } from '../stores/applications'
import { useUiStore } from '../stores/ui'
import { useEligibilityStore } from '../stores/eligibility'
import SchemeExplorer from '../components/SchemeExplorer.vue'
import { useI18n } from 'vue-i18n'

const schemeStore = useSchemeStore()
const bookmarkStore = useBookmarkStore()
const authStore = useAuthStore()
const applicationStore = useApplicationStore()
const uiStore = useUiStore()
const eligibilityStore = useEligibilityStore()
const { t, locale, messages } = useI18n()

const tObj = computed(() => messages.value[locale.value] || {})

const recommendedSchemes = computed(() => {
  if (!authStore.isLoggedIn) return []
  return eligibilityStore.results?.eligible ? eligibilityStore.results.eligible.map(es => es.scheme) : []
})

function handleApplyAction(scheme) {
  if (!authStore.isLoggedIn) {
    authStore.openAuthModal('login')
    uiStore.showToast(t('loginRequiredToast'), 'info')
    return
  }
  applicationStore.openApplyModal(scheme)
}

// ── Infinite scroll ────────────────────────────────────────────────────────
const sentinelRef = ref(null)
let observer = null

function setupInfiniteScroll() {
  if (observer) observer.disconnect()
  observer = new IntersectionObserver(
    (entries) => {
      if (entries[0].isIntersecting && schemeStore.hasNextPage && !schemeStore.loadingMore) {
        schemeStore.loadMoreSchemes()
      }
    },
    { threshold: 0.1 }
  )
  if (sentinelRef.value) observer.observe(sentinelRef.value)
}

// Search debouncer — resets to page 1
let searchTimeout
watch(() => schemeStore.searchQuery, () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    schemeStore.fetchSchemes()
  }, 400)
})

watch([() => schemeStore.selectedCategory, () => schemeStore.sortBy], () => {
  schemeStore.fetchSchemes()
})

onMounted(() => {
  schemeStore.fetchSchemes()
  // Sentinel element is rendered by SchemeExplorer; set up after next tick
  setTimeout(setupInfiniteScroll, 200)
})

onUnmounted(() => {
  if (observer) observer.disconnect()
})
</script>

<template>
  <SchemeExplorer
    v-model:selectedCategory="schemeStore.selectedCategory"
    v-model:sortBy="schemeStore.sortBy"
    v-model:searchQu="schemeStore.searchQuery"
    :schemes="schemeStore.schemes"
    :recommended-schemes="recommendedSchemes"
    :loading="schemeStore.loading"
    :loading-more="schemeStore.loadingMore"
    :has-next-page="schemeStore.hasNextPage"
    :error="schemeStore.error"
    :current-language="uiStore.currentLanguage"
    :saved-scheme-ids="bookmarkStore.savedSchemeIds"
    :categories="schemeStore.categories"
    :t="tObj"
    :is-logged-in="authStore.isLoggedIn"
    @toggle-bookmark="bookmarkStore.toggleBookmark"
    @open-details="schemeStore.openDetails"
    @retry="schemeStore.fetchSchemes"
    @login-required="authStore.openAuthModal('login')"
    @apply-click="handleApplyAction"
    @load-more="schemeStore.loadMoreSchemes"
  />
  <!-- Infinite scroll sentinel -->
  <div ref="sentinelRef" class="scroll-sentinel" aria-hidden="true"></div>
</template>

<style scoped>
.scroll-sentinel {
  height: 1px;
  width: 100%;
}
</style>
