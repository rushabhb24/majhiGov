<script setup>
import { computed, ref } from 'vue'
import { useBookmarkStore } from '../stores/bookmarks'
import { useAuthStore } from '../stores/auth'
import { useSchemeStore } from '../stores/schemes'
import { useJobStore } from '../stores/jobs'
import { useApplicationStore } from '../stores/applications'
import { useUiStore } from '../stores/ui'
import SchemeCard from '../components/SchemeCard.vue'
import JobCard from '../components/JobCard.vue'
import AppButton from '../components/ui/AppButton.vue'
import AppBadge from '../components/ui/AppBadge.vue'
import AppDialog from '../components/ui/AppDialog.vue'
import AppLabel from '../components/ui/AppLabel.vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const bookmarkStore = useBookmarkStore()
const authStore = useAuthStore()
const schemeStore = useSchemeStore()
const jobStore = useJobStore()
const applicationStore = useApplicationStore()
const uiStore = useUiStore()
const { t, locale, messages } = useI18n()
const router = useRouter()

const tObj = computed(() => messages.value[locale.value] || {})

// Tabs: 'schemes' or 'jobs'
const activeSection = ref('schemes')

// Job Details Modal state
const isJobModalOpen = ref(false)
const selectedJob = ref(null)

function openJobDetails(job) {
  selectedJob.value = job
  isJobModalOpen.value = true
}

function closeJobDetails() {
  selectedJob.value = null
  isJobModalOpen.value = false
}

function handleApplyAction(scheme) {
  if (!authStore.isLoggedIn) {
    authStore.openAuthModal('login')
    return
  }
  applicationStore.openApplyModal(scheme)
}

function handleJobApply(job) {
  if (job.apply_link) {
    window.open(job.apply_link, '_blank')
  } else {
    uiStore.showToast('Official portal link not available.', 'info')
  }
}

function formatDate(dateStr) {
  if (!dateStr) return 'N/A'
  try {
    return new Date(dateStr).toLocaleDateString(locale.value === 'mr' ? 'mr-IN' : (locale.value === 'hi' ? 'hi-IN' : 'en-IN'), {
      day: 'numeric',
      month: 'short',
      year: 'numeric'
    })
  } catch (e) {
    return dateStr
  }
}

function formatDocuments(docs) {
  if (!docs) return []
  if (Array.isArray(docs)) return docs
  return docs.split(',').map(d => d.trim()).filter(Boolean)
}
</script>

<template>
  <div class="tw-max-w-7xl tw-mx-auto tw-px-4 tw-sm:px-6 tw-lg:px-8 tw-py-8">
    
    <!-- Header panel -->
    <div class="glass tw-p-6 tw-rounded-2xl tw-mb-6">
      <h2 class="tw-font-heading tw-font-bold tw-text-xl tw-text-foreground tw-m-0">
        {{ tObj.savedTitle || 'Bookmarks' }}
      </h2>
      <p class="tw-text-xs tw-text-muted-foreground tw-mt-1 tw-m-0">
        {{ tObj.savedSubtitle || 'View and manage items you have bookmarked for quick access.' }}
      </p>
    </div>

    <!-- Section Switcher Tabs -->
    <div class="tw-flex tw-gap-2 tw-bg-muted/40 tw-p-1 tw-rounded-xl tw-max-w-xs tw-mb-6">
      <button
        class="tw-flex-1 tw-py-2 tw-rounded-lg tw-text-xs tw-font-bold tw-font-heading tw-border-none tw-cursor-pointer tw-transition-all"
        :class="activeSection === 'schemes' ? 'tw-bg-card tw-text-foreground tw-shadow-sm' : 'tw-text-muted-foreground tw-bg-transparent hover:tw-text-foreground'"
        @click="activeSection = 'schemes'"
      >
        Schemes ({{ bookmarkStore.bookmarkedSchemes.length }})
      </button>
      <button
        class="tw-flex-1 tw-py-2 tw-rounded-lg tw-text-xs tw-font-bold tw-font-heading tw-border-none tw-cursor-pointer tw-transition-all"
        :class="activeSection === 'jobs' ? 'tw-bg-card tw-text-foreground tw-shadow-sm' : 'tw-text-muted-foreground tw-bg-transparent hover:tw-text-foreground'"
        @click="activeSection = 'jobs'"
      >
        Jobs ({{ bookmarkStore.bookmarkedJobs.length }})
      </button>
    </div>

    <!-- SCHEMES SECTION -->
    <div v-if="activeSection === 'schemes'">
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

    <!-- JOBS SECTION -->
    <div v-else>
      <!-- Empty State -->
      <div v-if="bookmarkStore.bookmarkedJobs.length === 0" class="glass tw-p-12 tw-rounded-2xl tw-text-center tw-flex tw-flex-col tw-items-center tw-gap-4">
        <div class="tw-text-5xl">💼</div>
        <h3 class="tw-font-heading tw-font-bold tw-text-base tw-text-foreground tw-m-0">
          No Saved Jobs
        </h3>
        <AppButton variant="primary" size="sm" @click="router.push('/jobs')">
          Explore Jobs
        </AppButton>
      </div>

      <!-- Jobs Grid List -->
      <div v-else class="tw-grid tw-grid-cols-1 md:tw-grid-cols-3 tw-gap-6">
        <JobCard
          v-for="job in bookmarkStore.bookmarkedJobs"
          :key="job.id"
          :job="job"
          :isPrivate="job.isPrivate"
          :currentLanguage="uiStore.currentLanguage"
          :isLoggedIn="authStore.isLoggedIn"
          :isBookmarked="true"
          @viewDetails="openJobDetails"
          @applyClick="handleJobApply"
          @toggleBookmark="bookmarkStore.toggleJobBookmark(job.id)"
        />
      </div>
    </div>

    <!-- JOB DETAILED VIEW MODAL -->
    <AppDialog
      :open="isJobModalOpen && !!selectedJob"
      @close="closeJobDetails"
      maxWidth="650px"
    >
      <div class="tricolor-bar tw-h-[4px] tw-absolute tw-top-0 tw-left-0 tw-w-full"></div>

      <div v-if="selectedJob" class="tw-flex tw-flex-col tw-gap-4 tw-mt-4">
        <div class="tw-flex tw-items-center tw-gap-4">
          <div class="tw-w-10 tw-h-10 tw-rounded-xl tw-bg-primary/10 tw-text-primary tw-flex tw-items-center tw-justify-center tw-text-xl">
            💼
          </div>
          <div>
            <h2 class="tw-font-heading tw-font-bold tw-text-lg tw-text-foreground tw-m-0">
              {{ selectedJob.title }}
            </h2>
            <span class="tw-text-xs tw-text-muted-foreground tw-font-semibold">
              {{ selectedJob.company_name || selectedJob.organization }}
            </span>
          </div>
        </div>

        <hr class="tw-border-border/50 tw-my-1" />

        <div class="tw-max-h-[55vh] tw-overflow-y-auto tw-pr-2 tw-flex tw-flex-col tw-gap-4">
          
          <!-- Private Job details -->
          <div v-if="selectedJob.isPrivate" class="tw-grid tw-grid-cols-1 sm:tw-grid-cols-2 tw-gap-4">
            <div class="tw-flex tw-flex-col tw-gap-3">
              <div v-if="selectedJob.industry">
                <AppLabel>Industry</AppLabel>
                <div class="tw-text-xs tw-text-foreground tw-font-medium">{{ selectedJob.industry }}</div>
              </div>
              <div v-if="selectedJob.location">
                <AppLabel>Location</AppLabel>
                <div class="tw-text-xs tw-text-foreground tw-font-medium">{{ selectedJob.location }} ({{ selectedJob.work_mode }})</div>
              </div>
              <div v-if="selectedJob.stipend">
                <AppLabel>Stipend</AppLabel>
                <div class="tw-text-xs tw-text-foreground tw-font-medium">{{ selectedJob.stipend }}</div>
              </div>
              <div v-else>
                <AppLabel>Salary Package</AppLabel>
                <div class="tw-text-xs tw-text-foreground tw-font-medium">₹{{ (selectedJob.salary_min / 100000).toFixed(1) }} - ₹{{ (selectedJob.salary_max / 100000).toFixed(1) }} LPA</div>
              </div>
            </div>
            
            <div class="tw-flex tw-flex-col tw-gap-3">
              <div>
                <AppLabel>Requirements</AppLabel>
                <div class="tw-text-xs tw-text-muted-foreground tw-line-height-[1.5]">{{ selectedJob.requirements }}</div>
              </div>
              <div v-if="selectedJob.skills && selectedJob.skills.length > 0">
                <AppLabel>Required Skills</AppLabel>
                <div class="tw-flex tw-flex-wrap tw-gap-1.5 tw-mt-1">
                  <AppBadge v-for="skill in selectedJob.skills" :key="skill" tone="neutral">
                    {{ skill }}
                  </AppBadge>
                </div>
              </div>
            </div>
          </div>

          <!-- Government Job details -->
          <div v-else class="tw-grid tw-grid-cols-1 sm:tw-grid-cols-2 tw-gap-4">
            <div class="tw-flex tw-flex-col tw-gap-3">
              <div>
                <AppLabel>Department / Branch</AppLabel>
                <div class="tw-text-xs tw-text-foreground tw-font-medium">{{ selectedJob.department }}</div>
              </div>
              <div>
                <AppLabel>Qualifications Needed</AppLabel>
                <div class="tw-text-xs tw-text-foreground tw-font-medium">🎓 {{ selectedJob.education_qualification }}</div>
              </div>
              <div>
                <AppLabel>Total Vacancies</AppLabel>
                <div class="tw-text-xs tw-text-emerald-600 dark:tw-text-emerald-400 tw-font-bold">{{ selectedJob.vacancies?.toLocaleString() }} Openings</div>
              </div>
            </div>

            <div class="tw-flex tw-flex-col tw-gap-3">
              <div>
                <AppLabel>Required Verification Documents</AppLabel>
                <div class="tw-flex tw-flex-col tw-gap-1.5 tw-mt-1">
                  <div 
                    v-for="(doc, idx) in formatDocuments(selectedJob.required_documents)" 
                    :key="idx" 
                    class="tw-flex tw-items-center tw-gap-2 tw-p-2 tw-bg-muted/40 tw-rounded-lg tw-text-[11px] tw-text-muted-foreground"
                  >
                    <span class="tw-text-emerald-500">✓</span>
                    <span>{{ doc }}</span>
                  </div>
                </div>
              </div>
              <div>
                <AppLabel>Application Calendar Dates</AppLabel>
                <div class="tw-p-2.5 tw-bg-muted/40 tw-rounded-lg tw-text-[11px] tw-text-muted-foreground tw-flex tw-flex-col tw-gap-1">
                  <div>Posted Date: <strong>{{ formatDate(selectedJob.created_at) }}</strong></div>
                  <div>Last Date: <strong class="tw-text-destructive">{{ formatDate(selectedJob.application_end_date) }}</strong></div>
                </div>
              </div>
            </div>
          </div>

        </div>

        <hr class="tw-border-border/50 tw-my-1" />

        <div class="tw-flex tw-justify-end tw-gap-3">
          <AppButton variant="outline" size="sm" @click="closeJobDetails">
            Close Details
          </AppButton>
          <AppButton variant="primary" size="sm" @click="handleJobApply(selectedJob)">
            Apply on Portal
          </AppButton>
        </div>
      </div>
    </AppDialog>

  </div>
</template>
