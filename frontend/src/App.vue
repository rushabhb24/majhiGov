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

// Import Shared Components
import Header from './components/Header.vue'
import Hero from './components/Hero.vue'
import SchemeDetailsModal from './components/SchemeDetailsModal.vue'
import ToastBanner from './components/ToastBanner.vue'

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

// Handle tab/route navigation
function handleTabChange(tabName) {
  const routeMap = {
    explorer: '/',
    eligibility: '/eligibility',
    saved: '/saved',
    applications: '/applications',
    profile: '/profile'
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
  if (authStore.token) {
    authStore.fetchUserProfile()
  }
})
</script>

<template>
  <div :class="['app-wrapper', { 'rural-mode': uiStore.ruralMode }, uiStore.theme]">
    <!-- Header component (Logo, selects, tabs, togglers) -->
    <Header
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
      <Hero :t="tObj" @start-check="handleTabChange('eligibility')" />

      <!-- Router View - renders current route's component -->
      <router-view />
    </main>

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

    <!-- Beautiful Glassmorphic Auth Modal -->
    <Transition name="modal-fade">
      <div v-if="authStore.authModalOpen" class="modal-overlay" @click.self="authStore.closeAuthModal()">
        <div class="modal-content card" style="max-width: 580px; width: 100%;">
          <button class="btn-close-modal" @click="authStore.closeAuthModal()" title="Close Modal">×</button>
          
          <!-- Auth Tabs -->
          <div class="auth-tabs">
            <button 
              :class="['auth-tab-btn', { active: authStore.authTab === 'login' }]"
              @click="authStore.authTab = 'login'"
            >
              {{ t('loginRegister').split(' / ')[0] }}
            </button>
            <button 
              :class="['auth-tab-btn', { active: authStore.authTab === 'register' }]"
              @click="authStore.authTab = 'register'"
            >
              {{ t('loginRegister').split(' / ')[1] }}
            </button>
          </div>

          <!-- Login form -->
          <form v-if="authStore.authTab === 'login'" @submit.prevent="authStore.loginUser()" class="auth-form mt-4">
            <div class="form-group">
              <label class="form-label">{{ t('emailLabel') }} *</label>
              <input 
                v-model="authStore.loginForm.email" 
                type="email" 
                class="form-control" 
                placeholder="citizen@gov.in" 
                required 
              />
            </div>
            
            <div class="form-group">
              <label class="form-label">{{ t('passwordLabel') }} *</label>
              <input 
                v-model="authStore.loginForm.password" 
                type="password" 
                class="form-control" 
                placeholder="••••••••" 
                required 
              />
            </div>

            <button 
              type="submit" 
              class="btn btn-primary mt-4" 
              :disabled="authStore.authSubmitting"
            >
              {{ authStore.authSubmitting ? t('submitting') : t('loginRegister').split(' / ')[0] }}
            </button>
          </form>

          <!-- Register form -->
          <form v-else @submit.prevent="authStore.registerUser()" class="auth-form mt-4">
            <div class="auth-scroll-area">
              <div class="form-row">
                <div class="form-group">
                  <label class="form-label">{{ t('fullNameLabel') }} *</label>
                  <input 
                    v-model="authStore.regForm.full_name" 
                    type="text" 
                    class="form-control" 
                    placeholder="Ram Prasad" 
                    required 
                  />
                </div>
                <div class="form-group">
                  <label class="form-label">{{ t('dobLabel') }} *</label>
                  <input 
                    v-model="authStore.regForm.date_of_birth" 
                    type="date" 
                    class="form-control" 
                    required 
                  />
                </div>
              </div>

              <div class="form-row">
                <div class="form-group">
                  <label class="form-label">{{ t('emailLabel') }} *</label>
                  <input 
                    v-model="authStore.regForm.email" 
                    type="email" 
                    class="form-control" 
                    placeholder="ram@gov.in" 
                    required 
                  />
                </div>
                <div class="form-group">
                  <label class="form-label">{{ t('phoneLabel') }} *</label>
                  <input 
                    v-model="authStore.regForm.phone" 
                    type="tel" 
                    class="form-control" 
                    placeholder="9876543210" 
                    required 
                  />
                </div>
              </div>

              <div class="form-group">
                <label class="form-label">{{ t('passwordLabel') }} *</label>
                <input 
                  v-model="authStore.regForm.password" 
                  type="password" 
                  class="form-control" 
                  placeholder="Create secure password" 
                  required 
                  />
              </div>

              <hr class="divider mt-2" />
              <h4 class="form-section-title mt-2" style="font-size: 0.9rem; color: var(--clr-primary);">{{ t('demographicDetails') }}</h4>

              <div class="form-row mt-2">
                <div class="form-group">
                  <label class="form-label">{{ t('genderLabel') }} *</label>
                  <select v-model="authStore.regForm.gender" class="form-control" required>
                    <option value="Male">{{ t('maleOpt') }}</option>
                    <option value="Female">{{ t('femaleOpt') }}</option>
                    <option value="Other">{{ t('otherOpt') }}</option>
                  </select>
                </div>
                <div class="form-group">
                  <label class="form-label">{{ t('casteLabel') }} *</label>
                  <select v-model="authStore.regForm.caste_category" class="form-control" required>
                    <option value="General">General / Open</option>
                    <option value="OBC">OBC</option>
                    <option value="SC">SC</option>
                    <option value="ST">ST</option>
                  </select>
                </div>
              </div>

              <div class="form-row">
                <div class="form-group">
                  <label class="form-label">{{ t('stateLabel') }} *</label>
                  <select v-model="authStore.regForm.state" class="form-control" required>
                    <option value="Maharashtra">Maharashtra</option>
                    <option value="Gujarat">Gujarat</option>
                    <option value="Madhya Pradesh">Madhya Pradesh</option>
                    <option value="Karnataka">Karnataka</option>
                    <option value="Delhi">Delhi</option>
                    <option value="All">All India</option>
                  </select>
                </div>
                <div class="form-group">
                  <label class="form-label">{{ t('districtLabel') }} *</label>
                  <input 
                    v-model="authStore.regForm.district" 
                    type="text" 
                    class="form-control" 
                    placeholder="Pune" 
                    required 
                  />
                </div>
              </div>

              <div class="form-row">
                <div class="form-group">
                  <label class="form-label">{{ t('incomeLabel') }} *</label>
                  <input 
                    v-model="authStore.regForm.annual_income" 
                    type="number" 
                    class="form-control" 
                    required 
                  />
                </div>
                <div class="form-group">
                  <label class="form-label">{{ t('occupationLabel') }} *</label>
                  <select v-model="authStore.regForm.occupation" class="form-control" required>
                    <option value="Farmer">Farmer</option>
                    <option value="Student">Student</option>
                    <option value="Business Owner">Business Owner</option>
                    <option value="Unemployed">Unemployed</option>
                    <option value="Retired">Retired / Senior Citizen</option>
                  </select>
                </div>
              </div>

              <div class="form-row">
                <div class="form-group">
                  <label class="form-label">{{ t('employeeTypeLabel') }} *</label>
                  <select v-model="authStore.regForm.employee_type" class="form-control" required>
                    <option value="Unemployed">Unemployed</option>
                    <option value="Private Employee">Private Sector</option>
                    <option value="Government Employee">Government Sector</option>
                    <option value="Self Employed">Self Employed</option>
                  </select>
                </div>
                <div class="form-group">
                  <label class="form-label">{{ t('educationLabel') }} *</label>
                  <select v-model="authStore.regForm.education_level" class="form-control" required>
                    <option value="10th Pass">10th Standard or lower</option>
                    <option value="12th Pass">12th Standard</option>
                    <option value="Undergraduate">Undergraduate Degree</option>
                    <option value="Graduate">Graduate / Master Degree</option>
                    <option value="Post Graduate">Doctorate / Specialist</option>
                  </select>
                </div>
              </div>

              <div class="form-group">
                <label class="form-label">{{ t('disabilityLabel') }}</label>
                <div style="display: flex; align-items: center; gap: 8px;">
                  <input 
                    type="checkbox" 
                    id="regIsDisabled" 
                    v-model="authStore.regForm.is_disabled" 
                    style="width: 20px; height: 20px; cursor: pointer;"
                  />
                  <label for="regIsDisabled" style="cursor: pointer; font-size: 0.9rem;">Yes, I am differently-abled / Haan, mai divyang hu</label>
                </div>
              </div>
            </div>

            <button 
              type="submit" 
              class="btn btn-primary mt-4" 
              :disabled="authStore.authSubmitting"
            >
              {{ authStore.authSubmitting ? t('submitting') : t('loginRegister').split(' / ')[1] }}
            </button>
          </form>
        </div>
      </div>
    </Transition>
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
