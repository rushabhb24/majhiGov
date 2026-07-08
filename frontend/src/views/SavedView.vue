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
import AppButton from '../components/ui/AppButton.vue'

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
  <div class="tw-max-w-7xl tw-mx-auto tw-px-4 tw-sm:px-6 tw-lg:px-8 tw-py-8">
    
    <!-- Header panel -->
    <div class="glass tw-p-6 tw-rounded-2xl tw-mb-6">
      <h2 class="tw-font-heading tw-font-bold tw-text-xl tw-text-foreground tw-m-0">
        {{ tObj.savedTitle || 'Saved Schemes' }}
      </h2>
      <p class="tw-text-xs tw-text-muted-foreground tw-mt-1 tw-m-0">
        {{ tObj.savedSubtitle || 'View and manage schemes you have bookmarked for quick access.' }}
      </p>
    </div>

    <!-- Empty State -->
    <div v-if="bookmarkStore.bookmarkedSchemes.length === 0" class="glass tw-p-12 tw-rounded-2xl tw-text-center tw-flex tw-flex-col tw-items-center tw-gap-4">
      <div class="tw-text-5xl">🔖</div>
      <h3 class="tw-font-heading tw-font-bold tw-text-base tw-text-foreground tw-m-0">
        {{ tObj.noSaved || 'No Saved Schemes' }}
      </h3>
      <AppButton variant="primary" size="sm" @click="router.push('/')">
        {{ tObj.exploreSchemes || 'Explore Schemes' }}
      </AppButton>
    </div>

    <!-- Schemes Grid List -->
    <div v-else class="tw-grid tw-grid-cols-1 md:tw-grid-cols-3 tw-gap-6">
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
