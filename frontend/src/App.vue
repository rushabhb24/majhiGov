<script setup>
import { computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

// Import Stores
import { useAuthStore } from './stores/auth'
import { useSchemeStore } from './stores/schemes'
import { useBookmarkStore } from './stores/bookmarks'
import { useApplicationStore } from './stores/applications'
import { useUiStore } from './stores/ui'

import Header from './components/Header.vue'
import Hero from './components/Hero.vue'
import Footer from './components/Footer.vue'
import SchemeDetailsModal from './components/SchemeDetailsModal.vue'
import ToastBanner from './components/ToastBanner.vue'
import AuthModal from './components/AuthModal.vue'

// Initialize stores
const authStore = useAuthStore()
const schemeStore = useSchemeStore()
const bookmarkStore = useBookmarkStore()
const applicationStore = useApplicationStore()
const uiStore = useUiStore()

// i18n — t() for function calls in this template, tObj for child component props
const { t, locale, messages } = useI18n()
const tObj = computed(() => messages.value[locale.value] || {})
const router = useRouter()

// Sync i18n locale with store
watch(() => uiStore.currentLanguage, (newLang) => {
  locale.value = newLang
}, { immediate: true })

// Watch theme to set dark class on <html>
watch(() => uiStore.theme, (newTheme) => {
  if (newTheme === 'dark') {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}, { immediate: true })

// Handle tab/route navigation
function handleTabChange(tabName) {
  const routeMap = {
    explorer: '/',
    eligibility: '/eligibility',
    saved: '/saved',
    jobs: '/jobs',
    'ai-assistant': '/ai-assistant',
    'career-resources': '/career-resources',
    applications: '/applications',
    profile: '/profile',
    admin: '/admin/dashboard'
  }
  const path = routeMap[tabName] || '/'
  router.push(path)
}

// Handle apply action — opens official link synchronously, then records application
function handleApplyAction(scheme) {
  if (!authStore.isLoggedIn) {
    authStore.openAuthModal('login')
    uiStore.showToast(t('loginRequiredToast'), 'info')
    return
  }
  // Open official portal IMMEDIATELY (synchronous user gesture — won't be popup-blocked)
  const applyUrl = scheme.apply_link || scheme.official_website
  if (applyUrl) {
    window.open(applyUrl, '_blank', 'noopener,noreferrer')
  } else {
    uiStore.showToast(t('noOfficialLink') || 'Official apply link not available for this scheme.', 'warning')
  }
  // Record in background (async — doesn't need popup)
  applicationStore.applyViaOfficialLink(scheme)
}

// Lifecycle
onMounted(() => {
  bookmarkStore.loadBookmarks()
  // Unconditionally attempt to restore session via httpOnly cookie
  authStore.fetchUserProfile()
})
</script>

<template>
  <div :class="['app-wrapper', uiStore.theme, { 'rural-mode': uiStore.ruralMode }]">
    <!-- Header component (Logo, selects, tabs, togglers) -->
    <Header
      v-if="!$route.path.startsWith('/admin')"
      :activeTab="$route.name"
      @update:activeTab="handleTabChange"
      :currentLanguage="uiStore.currentLanguage"
      @update:currentLanguage="(v) => { uiStore.setLanguage(v) }"
      :ruralMode="uiStore.ruralMode"
      @update:ruralMode="(v) => { uiStore.setRuralMode(v) }"
      :theme="uiStore.theme"
      @update:theme="(v) => { uiStore.setTheme(v) }"
      :saved-count="bookmarkStore.savedSchemeIds.length"
      :t="tObj"
      :user="authStore.userProfile"
      @loginClick="authStore.openAuthModal('login')"
      @logout="authStore.logoutUser(); router.push('/')"
    />

    <!-- Main Viewport Shell -->
    <main class="main-container">
      
      <!-- Premium Hero Headline banner -->
      <Hero v-if="$route.path === '/'" :t="tObj" @start-check="handleTabChange('eligibility')" />

      <!-- Router View - renders current route's component -->
      <router-view />
    </main>

    <!-- Site Footer -->
    <Footer v-if="!$route.path.startsWith('/admin')" />

    <!-- Details relational modal overlay (Acc FAQ + Docs lists) -->
    <SchemeDetailsModal
      :scheme="schemeStore.selectedScheme"
      :current-language="uiStore.currentLanguage"
      :saved-scheme-ids="bookmarkStore.savedSchemeIds"
      :open="schemeStore.detailModalOpen"
      :t="tObj"
      :is-logged-in="authStore.isLoggedIn"
      @close="schemeStore.closeDetails"
      @toggle-bookmark="bookmarkStore.toggleBookmark"
      @login-required="authStore.openAuthModal('login')"
      @apply-click="(s) => { schemeStore.closeDetails(); handleApplyAction(s); }"
    />


    <!-- Frosted Notification banner alerts -->
    <ToastBanner 
      :show="uiStore.toast.show"
      :message="uiStore.toast.message"
      :type="uiStore.toast.type"
    />

    <!-- Beautiful Glassmorphic Application Form Modal -->
    <AppDialog
      :open="applicationStore.applyModalOpen && !!applicationStore.applyingScheme"
      @close="applicationStore.closeApplyModal()"
      maxWidth="600px"
    >
      <div v-if="applicationStore.applyingScheme" class="tw-flex tw-flex-col tw-gap-4">
        <h2 class="tw-font-heading tw-font-bold tw-text-xl tw-text-foreground tw-m-0">
          {{ t('applyFormTitle') || 'Submit Government Application Form' }}
        </h2>
        <h4 class="tw-font-heading tw-font-bold tw-text-sm tw-text-primary tw-m-0">
          {{ uiStore.currentLanguage === 'mr' ? (applicationStore.applyingScheme.title_mr || applicationStore.applyingScheme.title) : (uiStore.currentLanguage === 'hi' ? (applicationStore.applyingScheme.title_hi || applicationStore.applyingScheme.title) : applicationStore.applyingScheme.title) }}
        </h4>

        <hr class="tw-border-border/50 tw-my-1" />

        <p class="tw-text-xs tw-text-muted-foreground tw-line-height-[1.5] tw-m-0">
          {{ t('requiredDocumentsInfo') || 'Review required documents before submitting your request. Make sure you possess all mandatory credentials.' }}
        </p>

        <!-- Mandatory Documents checklist -->
        <div class="tw-p-4 tw-rounded-xl tw-bg-muted/50 tw-border tw-border-border">
          <h5 class="tw-font-heading tw-font-bold tw-text-xs tw-text-foreground tw-mb-2 tw-mt-0">
            📋 {{ t('modalDocs') || 'Required Documents' }}
          </h5>
          <div 
            v-for="doc in applicationStore.applyingScheme.documents" 
            :key="doc.id" 
            class="tw-flex tw-items-start tw-gap-2 tw-mt-2 tw-text-xs tw-text-muted-foreground"
          >
            <span class="tw-text-primary">🟢</span>
            <div>
              <strong class="tw-text-foreground">{{ uiStore.currentLanguage === 'mr' ? doc.document_name_mr : (uiStore.currentLanguage === 'hi' ? doc.document_name_hi : doc.document_name) }}</strong>
              <span class="tw-text-[10px] tw-px-1.5 tw-py-0.5 tw-rounded-full tw-ml-2" :class="doc.is_mandatory ? 'tw-bg-destructive/15 tw-text-destructive' : 'tw-bg-success/15 tw-text-success'">
                {{ doc.is_mandatory ? t('mandatoryBadge') || 'Mandatory' : t('optionalBadge') || 'Optional' }}
              </span>
            </div>
          </div>
        </div>

        <!-- Notes Input -->
        <div>
          <AppLabel for="apply-notes-input">{{ t('notesLabel') || 'Demographic Notes & Supporting Statement' }}</AppLabel>
          <AppTextarea
            id="apply-notes-input"
            v-model="applicationStore.applyNotes"
            rows="4"
            :placeholder="t('notesPlaceholder') || 'Enter any additional notes or details...'"
          />
        </div>

        <div class="tw-flex tw-justify-end tw-gap-3 tw-mt-2">
          <AppButton variant="outline" size="sm" @click="applicationStore.closeApplyModal()">
            {{ t('back') || 'Back' }}
          </AppButton>
          <AppButton
            variant="primary"
            size="sm"
            @click="applicationStore.submitApplication()"
            :disabled="applicationStore.applySubmitting"
          >
            {{ applicationStore.applySubmitting ? (t('submittingApp') || 'Submitting...') : (t('submitAppBtn') || 'Submit Application') }}
          </AppButton>
        </div>
      </div>
    </AppDialog>

    <!-- Teleported Auth Modal component -->
    <AuthModal />
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
