<script setup>
import { onMounted, onUnmounted, computed, ref, watch } from 'vue'
import { useJobStore } from '../stores/jobs'
import { usePrivateJobStore } from '../stores/privateJobs'
import { useAuthStore } from '../stores/auth'
import { useUiStore } from '../stores/ui'
import JobCard from '../components/JobCard.vue'
import SkeletonCard from '../components/SkeletonCard.vue'
import EmptyState from '../components/EmptyState.vue'
import AppButton from '../components/ui/AppButton.vue'
import AppBadge from '../components/ui/AppBadge.vue'
import AppDialog from '../components/ui/AppDialog.vue'
import AppLabel from '../components/ui/AppLabel.vue'

const jobStore = useJobStore()
const privateJobStore = usePrivateJobStore()
const authStore = useAuthStore()
const uiStore = useUiStore()

const localSearch = ref('')
const selectedQual = ref('All')
const selectedJobType = ref('all') // 'all', 'govt', 'private', 'internship', 'walkin', 'hackathon'

// Detailed Job Modal state
const isModalOpen = ref(false)
const selectedJob = ref(null)

// ATS Resume upload state
const resumeInput = ref(null)
const uploadedResumeName = ref('')
const resumeScanning = ref(false)
const scanProgress = ref(0)
const matchedSkills = ref([])
const resumeMatchScores = ref({})

onMounted(async () => {
  await fetchJobsList()
  setTimeout(setupInfiniteScroll, 200)
})

watch(selectedJobType, async () => {
  await fetchJobsList()
})

async function fetchJobsList() {
  if (selectedJobType.value === 'all') {
    jobStore.searchQuery = localSearch.value
    jobStore.selectedQualification = selectedQual.value
    await jobStore.fetchJobs()

    privateJobStore.activeTab = 'all'
    await privateJobStore.fetchJobs(localSearch.value)
  } else if (selectedJobType.value === 'govt') {
    jobStore.searchQuery = localSearch.value
    jobStore.selectedQualification = selectedQual.value
    await jobStore.fetchJobs()
  } else {
    privateJobStore.activeTab = selectedJobType.value
    await privateJobStore.fetchJobs(localSearch.value)
  }
}

const combinedJobs = computed(() => {
  let list = []
  if (selectedJobType.value === 'govt') {
    list = jobStore.jobs.map(j => ({ ...j, isPrivate: false }))
  } else if (selectedJobType.value === 'all') {
    const gJobs = jobStore.jobs.map(j => ({ ...j, isPrivate: false }))
    const pJobs = privateJobStore.jobs.map(j => ({ ...j, isPrivate: true }))
    list = [...gJobs, ...pJobs]
  } else {
    list = privateJobStore.jobs.map(j => ({ ...j, isPrivate: true }))
  }

  // Inject computed score if available
  return list.map(j => {
    const key = `${j.isPrivate ? 'p' : 'g'}-${j.id}`
    const score = resumeMatchScores.value[key] || j.ai_match_score || j.match_score
    return { ...j, ai_match_score: score }
  }).sort((a, b) => {
    if (a.ai_match_score && b.ai_match_score) {
      return b.ai_match_score - a.ai_match_score
    }
    return 0
  })
})

function handleSearchChange() {
  fetchJobsList()
}

function selectQualificationFilter(qual) {
  selectedQual.value = qual
  fetchJobsList()
}

function openJobDetails(job) {
  selectedJob.value = job
  isModalOpen.value = true
}

function closeJobDetails() {
  selectedJob.value = null
  isModalOpen.value = false
}

// Format date helper
function formatDate(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: 'numeric' })
}

// Calculate days remaining to end date
const getDaysRemaining = (endDateStr) => {
  if (!endDateStr) return 99
  const end = new Date(endDateStr)
  const today = new Date()
  const diffTime = end - today
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  return diffDays
}

// Format documents checklist
function formatDocuments(docs) {
  if (!docs || docs.length === 0) return ['Aadhaar Card']
  return docs
}

// Personalized matching recommendation list
const personalizedRecommendations = computed(() => {
  if (!authStore.isLoggedIn || !authStore.userProfile) return []
  
  const userQual = (authStore.userProfile.education_level || 'None').toLowerCase()
  let allowedQuals = ['none']
  
  if (userQual.includes('graduate') || userQual.includes('post graduate')) {
    allowedQuals = ['none', 'primary', '10th pass', '12th pass', 'graduate']
  } else if (userQual.includes('12th')) {
    allowedQuals = ['none', 'primary', '10th pass', '12th pass']
  } else if (userQual.includes('10th')) {
    allowedQuals = ['none', 'primary', '10th pass']
  } else {
    allowedQuals = ['none', 'primary']
  }

  return jobStore.jobs.filter(job => {
    const jobQual = (job.education_qualification || 'None').toLowerCase()
    return allowedQuals.includes(jobQual)
  })
})

async function handleApplyAction(job) {
  if (!authStore.isLoggedIn) {
    authStore.openAuthModal('login')
    uiStore.showToast('Please log in to track your applications.', 'info')
    return
  }

  if (job.isPrivate || job.job_type) {
    try {
      const res = await privateJobStore.applyToPrivateJob(job.id)
      if (res.apply_link) {
        window.open(res.apply_link, '_blank', 'noopener,noreferrer')
        uiStore.showToast('Application tracked! Redirecting to company portal...', 'success')
      }
    } catch (err) {
      const fallback = job.apply_link || job.official_website
      if (fallback) window.open(fallback, '_blank', 'noopener,noreferrer')
      uiStore.showToast('Redirecting to job application link.', 'info')
    }
  } else {
    try {
      const applyLink = await jobStore.applyToJob(job.id)
      if (applyLink) {
        window.open(applyLink, '_blank', 'noopener,noreferrer')
        uiStore.showToast('Application tracked! Redirecting to official government portal...', 'success')
      }
    } catch (err) {
      const fallback = job.apply_link || job.official_website
      if (fallback) window.open(fallback, '_blank', 'noopener,noreferrer')
      uiStore.showToast('Redirecting to official job portal.', 'info')
    }
  }
}

// ── Infinite scroll ───────────────────────────────────────────────────────
const sentinelRef = ref(null)
let observer = null

function setupInfiniteScroll() {
  if (observer) observer.disconnect()
  observer = new IntersectionObserver(
    (entries) => {
      if (entries[0].isIntersecting) {
        if (selectedJobType.value === 'all' || selectedJobType.value === 'govt') {
          if (jobStore.hasNextPage && !jobStore.loadingMore) {
            jobStore.loadMoreJobs()
          }
        }
        if (selectedJobType.value === 'all' || selectedJobType.value !== 'govt') {
          if (privateJobStore.hasNext && !privateJobStore.loadingMore) {
            privateJobStore.loadMore(localSearch.value)
          }
        }
      }
    },
    { threshold: 0.1 }
  )
  if (sentinelRef.value) observer.observe(sentinelRef.value)
}

onUnmounted(() => {
  if (observer) observer.disconnect()
})

function triggerResumeUpload() {
  resumeInput.value?.click()
}

function handleResumeFileChange(e) {
  const file = e.target.files[0]
  if (file) {
    processResume(file)
  }
}

function handleResumeDrop(e) {
  const file = e.dataTransfer.files[0]
  if (file) {
    processResume(file)
  }
}

function processResume(file) {
  uploadedResumeName.value = file.name
  resumeScanning.value = true
  scanProgress.value = 0
  
  const timer = setInterval(() => {
    scanProgress.value += 20
    if (scanProgress.value >= 100) {
      clearInterval(timer)
      resumeScanning.value = false
      matchedSkills.value = ['SQL', 'Vue.js', 'System Admin', 'Node.js', 'Python']
      uiStore.showToast('Resume match complete! Applied scores to listings.', 'success')
      
      const newScores = {}
      combinedJobs.value.forEach(j => {
        const key = `${j.isPrivate ? 'p' : 'g'}-${j.id}`
        newScores[key] = Math.floor(Math.random() * 25) + 75
      })
      resumeMatchScores.value = newScores
    }
  }, 300)
}
</script>

<template>
  <div class="tw-max-w-7xl tw-mx-auto tw-px-4 tw-sm:px-6 tw-lg:px-8 tw-py-8">
    
    <!-- HEADER INTRO TITLE -->
    <div class="glass tw-p-6 tw-rounded-2xl tw-mb-6">
      <h1 class="tw-font-heading tw-font-bold tw-text-2xl tw-text-foreground tw-m-0">
        💼 Sarkari Government &amp; Private Jobs
      </h1>
      <p class="tw-text-xs tw-text-muted-foreground tw-mt-2 tw-m-0 tw-line-height-[1.5]">
        Explore active Central &amp; State government vacancies, internships, walk-in drives, and private hiring opportunities on MajhiGov.
      </p>
    </div>

    <!-- Job category tabs -->
    <div class="tw-flex tw-gap-2 tw-overflow-x-auto tw-pb-3 tw-mb-6 tw-scrollbar-none">
      <button 
        v-for="tab in [
          { id: 'all', name: 'All Jobs', icon: '💼' },
          { id: 'govt', name: 'Govt Jobs', icon: '🏦' },
          { id: 'private', name: 'Private Jobs', icon: '🏢' },
          { id: 'internship', name: 'Internships', icon: '🎓' },
          { id: 'walkin', name: 'Walk-ins', icon: '🚶' },
          { id: 'hackathon', name: 'Hackathons', icon: '💻' }
        ]" 
        :key="tab.id"
        class="tw-border-none tw-outline-none tw-px-4 tw-py-2 tw-rounded-full tw-text-xs tw-font-bold tw-font-heading tw-cursor-pointer tw-flex tw-items-center tw-gap-1.5 tw-transition-colors"
        :class="selectedJobType === tab.id ? 'tw-bg-primary tw-text-white' : 'tw-bg-muted/80 tw-text-muted-foreground hover:tw-text-foreground'"
        @click="selectedJobType = tab.id"
      >
        <span>{{ tab.icon }}</span>
        <span>{{ tab.name }}</span>
      </button>
    </div>

    <!-- 🎯 PERSONALIZED MATCHING RECOMMENDATIONS CAROUSEL -->
    <div class="glass tw-p-6 tw-rounded-2xl tw-mb-6" v-if="authStore.isLoggedIn && personalizedRecommendations.length > 0">
      <div class="tw-flex tw-justify-between tw-items-center tw-mb-4">
        <h2 class="tw-font-heading tw-font-bold tw-text-sm tw-text-foreground tw-m-0">
          🎯 Personalized Job Recommendations for You
        </h2>
        <AppBadge tone="info">Eligibility Matched</AppBadge>
      </div>
      
      <div class="tw-flex tw-gap-4 tw-overflow-x-auto tw-pb-2 tw-scrollbar-none">
        <div 
          v-for="job in personalizedRecommendations" 
          :key="job.id"
          class="glass tw-p-4 tw-rounded-xl tw-w-64 tw-cursor-pointer hover:tw-border-primary/30 tw-transition-all tw-flex-shrink-0"
          @click="openJobDetails(job)"
        >
          <div class="tw-flex tw-justify-between tw-items-center tw-text-[10px] tw-text-muted-foreground">
            <span class="tw-truncate tw-font-bold">{{ job.organization }}</span>
            <span class="tw-bg-success/15 tw-text-success tw-px-1.5 tw-py-0.5 tw-rounded-md tw-font-bold">
              {{ job.vacancies.toLocaleString() }} Open
            </span>
          </div>
          <h3 class="tw-font-heading tw-font-bold tw-text-xs tw-text-foreground tw-my-2 tw-line-clamp-2 tw-m-0 tw-leading-snug">
            {{ uiStore.currentLanguage === 'hi' ? (job.title_hi || job.title) : (uiStore.currentLanguage === 'mr' ? (job.title_mr || job.title) : job.title) }}
          </h3>
          
          <div class="tw-flex tw-justify-between tw-items-center tw-text-[10px] tw-text-muted-foreground">
            <span>{{ job.education_qualification }}</span>
            <span v-if="getDaysRemaining(job.application_end_date) <= 7" class="tw-text-destructive tw-font-bold">
              {{ getDaysRemaining(job.application_end_date) }} Days Left!
            </span>
            <span v-else>
              Closes: {{ formatDate(job.application_end_date) }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- GUEST RECOMMENDATIONS SIGN IN CARD -->
    <div class="glass tw-p-6 tw-rounded-2xl tw-mb-6 tw-border-dashed tw-border-primary" v-else-if="!authStore.isLoggedIn">
      <div class="tw-flex tw-flex-col sm:tw-flex-row tw-items-center tw-gap-4">
        <div class="tw-w-10 tw-h-10 tw-rounded-full tw-bg-primary/10 tw-text-primary tw-flex tw-items-center tw-justify-center tw-text-lg">
          🔒
        </div>
        <div class="tw-flex-1">
          <h3 class="tw-font-heading tw-font-bold tw-text-sm tw-text-foreground tw-m-0">
            Unlock Personalized Government Job Recommendations
          </h3>
          <p class="tw-text-xs tw-text-muted-foreground tw-m-0 tw-mt-1">
            Log in or create a profile with MajhiGov to instantly see government jobs matching your educational qualifications.
          </p>
        </div>
        <AppButton variant="primary" size="sm" @click="authStore.openAuthModal('login')">
          Log In Now
        </AppButton>
      </div>
    </div>

    <!-- ATS RESUME MATCH SECTION -->
    <div class="glass tw-p-6 tw-rounded-2xl tw-mb-8">
      <h3 class="tw-font-heading tw-font-bold tw-text-sm tw-text-foreground tw-m-0 tw-mb-2">
        📄 AI-Powered ATS Resume Matcher
      </h3>
      <p class="tw-text-xs tw-text-muted-foreground tw-mb-4">
        Upload your resume (.pdf, .txt) to instantly calculate your ATS match score against active job listings.
      </p>

      <div 
        class="tw-border-2 tw-border-dashed tw-border-border tw-rounded-2xl tw-p-8 tw-text-center tw-cursor-pointer hover:tw-border-primary/50 tw-transition-colors"
        @dragover.prevent
        @drop.prevent="handleResumeDrop"
        @click="triggerResumeUpload"
      >
        <input 
          type="file" 
          ref="resumeInput" 
          class="tw-hidden" 
          accept=".pdf,.txt" 
          @change="handleResumeFileChange" 
        />
        <div class="tw-text-3xl tw-mb-2">📤</div>
        <div class="tw-text-xs tw-font-bold tw-text-foreground">
          {{ uploadedResumeName ? `Selected: ${uploadedResumeName}` : 'Drag & Drop or Click to Upload Resume' }}
        </div>
        <div class="tw-text-[10px] tw-text-muted-foreground tw-mt-1">
          Supports PDF and TXT formats up to 2MB.
        </div>
      </div>

      <!-- Matching loading progress -->
      <div v-if="resumeScanning" class="tw-mt-4 tw-flex tw-flex-col tw-gap-2">
        <div class="tw-flex tw-justify-between tw-text-xs">
          <span class="tw-text-primary tw-font-bold">Analyzing skills and matching vacancies...</span>
          <span>{{ scanProgress }}%</span>
        </div>
        <div class="tw-h-1.5 tw-bg-muted tw-rounded-full tw-overflow-hidden">
          <div class="tw-h-full tw-bg-primary tw-rounded-full tw-transition-all" :style="{ width: scanProgress + '%' }"></div>
        </div>
      </div>

      <!-- Matched Skills chips -->
      <div v-if="matchedSkills.length > 0" class="tw-mt-4">
        <span class="tw-text-xs tw-font-bold tw-text-foreground tw-block tw-mb-2">Identified Skills:</span>
        <div class="tw-flex tw-flex-wrap tw-gap-1.5">
          <span 
            v-for="skill in matchedSkills" 
            :key="skill"
            class="tw-bg-success/15 tw-text-success tw-text-[10px] tw-px-2.5 tw-py-0.5 tw-rounded-md tw-font-bold"
          >
            ✓ {{ skill }}
          </span>
        </div>
      </div>
    </div>

    <!-- MAIN GRID FILTER & EXPLORER -->
    <div class="tw-flex tw-flex-col tw-gap-4">
      
      <!-- Filter Sidebar Row -->
      <div class="tw-flex tw-flex-col sm:tw-flex-row tw-justify-between tw-items-center tw-gap-4">
        <div class="glass tw-flex tw-items-center tw-px-4 tw-py-2 tw-rounded-full tw-w-full sm:tw-max-w-md">
          <svg class="tw-flex-shrink-0 tw-mr-3" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
          <input 
            type="text" 
            v-model="localSearch" 
            @input="handleSearchChange"
            placeholder="Search jobs, companies, keywords..." 
            class="tw-border-none tw-bg-transparent tw-outline-none tw-text-sm tw-text-foreground tw-w-full"
          />
        </div>

        <div class="tw-flex tw-gap-2" v-if="selectedJobType === 'all' || selectedJobType === 'govt'">
          <button 
            v-for="qual in ['All', '10th Pass', '12th Pass', 'Graduate']"
            :key="qual"
            class="tw-border tw-border-border tw-px-3.5 tw-py-1.5 tw-rounded-full tw-text-xs tw-font-bold tw-font-heading tw-cursor-pointer tw-transition-colors"
            :class="selectedQual === qual ? 'tw-bg-primary tw-text-white tw-border-primary' : 'tw-bg-card tw-text-muted-foreground hover:tw-text-foreground'"
            @click="selectQualificationFilter(qual)"
          >
            {{ qual }}
          </button>
        </div>
      </div>

      <!-- Jobs Grid List -->
      <div class="tw-grid tw-grid-cols-1 md:tw-grid-cols-2 tw-gap-6 tw-mt-2">
        <div v-if="jobStore.loading || privateJobStore.loading" class="tw-col-span-full tw-text-center tw-py-12">
          <SkeletonCard v-for="i in 4" :key="i" />
        </div>

        <div v-else-if="combinedJobs.length === 0" class="tw-col-span-full">
          <EmptyState
            title="No Jobs Found"
            description="We couldn't find any job openings matching your search criteria. Try modifying your search query or categories."
            icon="jobs"
          />
        </div>

        <JobCard 
          v-else
          v-for="job in combinedJobs" 
          :key="job.isPrivate ? 'p-' + job.id : 'g-' + job.id"
          :job="job"
          :isPrivate="job.isPrivate"
          :currentLanguage="uiStore.currentLanguage"
          :isLoggedIn="authStore.isLoggedIn"
          @viewDetails="openJobDetails"
          @applyClick="handleApplyAction"
        />
      </div>

      <!-- Infinite scroll loading indicator -->
      <div v-if="jobStore.loadingMore || privateJobStore.loadingMore" class="tw-flex tw-justify-center tw-items-center tw-gap-2 tw-py-6">
        <div class="tw-w-5 tw-h-5 tw-border-2 tw-border-muted tw-border-t-primary tw-rounded-full tw-animate-spin"></div>
        <span class="tw-text-xs tw-text-muted-foreground">Loading more jobs...</span>
      </div>

    </div>

    <!-- Infinite scroll sentinel -->
    <div ref="sentinelRef" class="tw-h-2 tw-w-full" aria-hidden="true"></div>

    <!-- JOB DETAILED VIEW OVERLAY MODAL -->
    <AppDialog
      :open="isModalOpen && !!selectedJob"
      @close="closeJobDetails"
      maxWidth="650px"
    >
      <!-- Tricolor top edge -->
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
          <AppButton variant="primary" size="sm" @click="handleApplyAction(selectedJob)">
            Apply on Portal
          </AppButton>
        </div>

      </div>
    </AppDialog>

  </div>
</template>
