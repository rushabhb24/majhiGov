<script setup>
import { computed, watch, onMounted } from 'vue'
import { useSchemeStore } from '../stores/schemes'
import { useBookmarkStore } from '../stores/bookmarks'
import { useAuthStore } from '../stores/auth'
import { useApplicationStore } from '../stores/applications'
import { useUiStore } from '../stores/ui'
import SchemeExplorer from '../components/SchemeExplorer.vue'
import { useI18n } from 'vue-i18n'

const schemeStore = useSchemeStore()
const bookmarkStore = useBookmarkStore()
const authStore = useAuthStore()
const applicationStore = useApplicationStore()
const uiStore = useUiStore()
const { t, locale, messages } = useI18n()

// Translation object for child components that use t.key property access
const tObj = computed(() => messages.value[locale.value] || {})

function handleApplyAction(scheme) {
  if (!authStore.isLoggedIn) {
    authStore.openAuthModal('login')
    uiStore.showToast(t('loginRequiredToast'), 'info')
    return
  }
  applicationStore.openApplyModal(scheme)
}

// Search debouncer
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
})
</script>

<template>
  <SchemeExplorer
    v-model:selectedCategory="schemeStore.selectedCategory"
    v-model:sortBy="schemeStore.sortBy"
    v-model:searchQu="schemeStore.searchQuery"
    :schemes="schemeStore.schemes"
    :loading="schemeStore.loading"
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
  />
</template>
