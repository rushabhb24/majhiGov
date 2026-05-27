<script setup>
import { computed } from 'vue'
import { useBookmarkStore } from '../stores/bookmarks'
import { useAuthStore } from '../stores/auth'
import { useSchemeStore } from '../stores/schemes'
import { useApplicationStore } from '../stores/applications'
import { useUiStore } from '../stores/ui'
import SchemeCard from '../components/SchemeCard.vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const bookmarkStore = useBookmarkStore()
const authStore = useAuthStore()
const schemeStore = useSchemeStore()
const applicationStore = useApplicationStore()
const uiStore = useUiStore()
const { t, locale, messages } = useI18n()
const router = useRouter()

// Translation object for child components that use t.key property access
const tObj = computed(() => messages.value[locale.value] || {})

function handleApplyAction(scheme) {
  if (!authStore.isLoggedIn) {
    authStore.openAuthModal('login')
    return
  }
  applicationStore.openApplyModal(scheme)
}
</script>

<template>
  <div class="tab-content animate-fade">
    <div class="card filter-panel">
      <h2 class="section-title">{{ tObj.savedTitle }}</h2>
      <p class="text-muted">{{ tObj.savedSubtitle }}</p>
    </div>

    <div v-if="bookmarkStore.bookmarkedSchemes.length === 0" class="empty-state text-center mt-4 card">
      <div class="empty-bookmarks-art">🔖</div>
      <h3>{{ tObj.noSaved }}</h3>
      <button class="btn btn-primary mt-4" @click="router.push('/')">{{ tObj.exploreSchemes || 'Explore Schemes' }}</button>
    </div>

    <div v-else class="schemes-grid mt-4">
      <SchemeCard
        v-for="scheme in bookmarkStore.bookmarkedSchemes"
        :key="scheme.id"
        :scheme="scheme"
        :current-language="uiStore.currentLanguage"
        :saved-scheme-ids="bookmarkStore.savedSchemeIds"
        :t="tObj"
        :is-logged-in="authStore.isLoggedIn"
        @toggle-bookmark="bookmarkStore.toggleBookmark"
        @open-details="schemeStore.openDetails"
        @login-required="authStore.openAuthModal('login')"
        @apply-click="handleApplyAction"
      />
    </div>
  </div>
</template>

<style scoped>
.animate-fade {
  animation: fadeIn 0.4s ease-out;
}
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}
.empty-bookmarks-art {
  font-size: 3.5rem;
  margin-bottom: 12px;
  filter: drop-shadow(0 6px 10px rgba(0,0,0,0.05));
}
</style>
