<script setup>
import { onMounted, computed, ref } from 'vue'
import { useJobStore } from '../stores/jobs'
import { useAuthStore } from '../stores/auth'
import { useUiStore } from '../stores/ui'

const jobStore = useJobStore()
const authStore = useAuthStore()
const uiStore = useUiStore()

const localSearch = ref('')
const selectedQual = ref('All')

// Detailed Job Modal state
const isModalOpen = ref(false)
const selectedJob = ref(null)

onMounted(async () => {
  await fetchJobsList()
})

async function fetchJobsList() {
  jobStore.searchQuery = localSearch.value
  jobStore.selectedQualification = selectedQual.value
  await jobStore.fetchJobs()
}

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
function getDaysRemaining(endDateStr) {
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

function handleApplyAction(job) {
  const applyUrl = job.apply_link || job.official_website
  if (applyUrl) {
    window.open(applyUrl, '_blank', 'noopener,noreferrer')
    uiStore.showToast('Redirecting to official government job application portal...', 'success')
  } else {
    uiStore.showToast('Official apply link is currently offline.', 'warning')
  }
}
</script>

<template>
  <div class="jobs-explorer-container">
    
    <!-- HEADER INTRO TITLE -->
    <div class="jobs-header-hero">
      <h1 class="jobs-hero-title">💼 Sarkari Government Jobs Portal</h1>
      <p class="jobs-hero-subtitle">
        Explore active Central & State government vacancies, filter by educational eligibility, and apply directly via official portals.
      </p>
    </div>

    <!-- 🎯 PERSONALIZED MATCHING RECOMMENDATIONS CAROUSEL -->
    <div class="recommendations-section mt-4" v-if="authStore.isLoggedIn && personalizedRecommendations.length > 0">
      <div class="sec-header">
        <h2 class="sec-title-text">🎯 Personalized Job Recommendations for You</h2>
        <span class="recommendation-badge"><i class="ti ti-sparkles"></i> Eligibility Matched</span>
      </div>
      
      <div class="recommendations-slider-wrapper">
        <div class="recommendation-track">
          <div 
            class="recommendation-card" 
            v-for="job in personalizedRecommendations" 
            :key="job.id"
            @click="openJobDetails(job)"
          >
            <div class="rec-card-header">
              <span class="rec-org">{{ job.organization }}</span>
              <span class="rec-vacancies">{{ job.vacancies.toLocaleString() }} Vacancies</span>
            </div>
            <h3 class="rec-title">{{ uiStore.currentLanguage === 'hi' ? (job.title_hi || job.title) : (uiStore.currentLanguage === 'mr' ? (job.title_mr || job.title) : job.title) }}</h3>
            
            <div class="rec-footer mt-2">
              <span class="rec-qual"><i class="ti ti-school"></i> {{ job.education_qualification }}</span>
              <span class="rec-deadline text-danger" v-if="getDaysRemaining(job.application_end_date) <= 7">
                <i class="ti ti-alarm"></i> {{ getDaysRemaining(job.application_end_date) }} Days Left!
              </span>
              <span class="rec-deadline" v-else>
                <i class="ti ti-calendar"></i> Closes: {{ formatDate(job.application_end_date) }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- GUEST RECOMMENDATIONS SIGN IN CARD -->
    <div class="recommendations-guest-card mt-4" v-else-if="!authStore.isLoggedIn">
      <div class="guest-card-content">
        <div class="guest-icon"><i class="ti ti-lock"></i></div>
        <div>
          <h3 class="guest-card-title">Unlock Personalized Government Job Recommendations</h3>
          <p class="guest-card-text">Log in or create a profile with MajhiGov to instantly see government jobs matching your educational qualifications and profile characteristics.</p>
        </div>
        <button class="btn-guest-login" @click="authStore.openAuthModal('login')">
          <i class="ti ti-login"></i> Log In Now
        </button>
      </div>
    </div>

    <!-- MAIN GRID FILTER & EXPLORER -->
    <div class="main-explorer-layout mt-4">
      
      <!-- Filter Sidebar Row -->
      <div class="explorer-filters-row">
        <div class="jobs-search-box">
          <i class="ti ti-search search-icon"></i>
          <input 
            type="text" 
            v-model="localSearch" 
            @input="handleSearchChange"
            placeholder="Search jobs, organizations, department..." 
            class="jobs-search-input"
          />
        </div>

        <div class="filter-pills">
          <button 
            v-for="qual in ['All', '10th Pass', '12th Pass', 'Graduate']"
            :key="qual"
            :class="['filter-pill', { active: selectedQual === qual }]"
            @click="selectQualificationFilter(qual)"
          >
            {{ qual }}
          </button>
        </div>
      </div>

      <!-- Jobs Grid List -->
      <div class="jobs-grid mt-3">
        <div v-if="jobStore.loading" class="jobs-loader">
          <div class="jobs-spinner"></div>
          <p>Fetching active government job listings...</p>
        </div>

        <div v-else-if="jobStore.jobs.length === 0" class="no-jobs-card">
          <i class="ti ti-inbox no-jobs-icon"></i>
          <h3>No Government Jobs Found</h3>
          <p>We couldn't find any government job openings matching your search parameters. Try clearing your filters.</p>
        </div>

        <div 
          v-else
          class="job-card" 
          v-for="job in jobStore.jobs" 
          :key="job.id"
          @click="openJobDetails(job)"
        >
          <div class="job-card-top">
            <span class="job-org">{{ job.organization }}</span>
            <span class="job-department">{{ job.department }}</span>
            
            <!-- Deadline Alert Alert Badge -->
            <span class="closing-alert-badge" v-if="getDaysRemaining(job.application_end_date) <= 7">
              ⚠️ Closing soon!
            </span>
          </div>

          <h3 class="job-card-title">
            {{ uiStore.currentLanguage === 'hi' ? (job.title_hi || job.title) : (uiStore.currentLanguage === 'mr' ? (job.title_mr || job.title) : job.title) }}
          </h3>

          <div class="job-vacancies-row mt-2">
            <div class="vacancy-item">
              <span class="lbl">Total Vacancies</span>
              <span class="val font-semibold">{{ job.vacancies.toLocaleString() }} Positions</span>
            </div>
            <div class="vacancy-item">
              <span class="lbl">Education Level</span>
              <span class="val">{{ job.education_qualification }}</span>
            </div>
            <div class="vacancy-item">
              <span class="lbl">Application Fee</span>
              <span class="val">{{ job.application_fee }}</span>
            </div>
          </div>

          <div class="job-card-footer mt-3">
            <span class="end-date">
              <i class="ti ti-calendar-event"></i> Closes: {{ formatDate(job.application_end_date) }}
            </span>
            <button class="btn-job-action" @click.stop="openJobDetails(job)">
              View Details <i class="ti ti-arrow-right"></i>
            </button>
          </div>
        </div>
      </div>

    </div>

    <!-- JOB DETAILED VIEW OVERLAY MODAL -->
    <Transition name="modal-fade">
      <div class="admin-modal-overlay" v-if="isModalOpen" @click.self="closeJobDetails">
        <div class="admin-modal-box modal-jobs-override" style="max-width: 650px;">
          <button class="btn-close" @click="closeJobDetails">×</button>

          <div class="modal-jobs-header">
            <div class="modal-jobs-icon-box">
              <i class="ti ti-briefcase"></i>
            </div>
            <div>
              <h2 class="modal-job-title">{{ uiStore.currentLanguage === 'hi' ? (selectedJob?.title_hi || selectedJob?.title) : (uiStore.currentLanguage === 'mr' ? (selectedJob?.title_mr || selectedJob?.title) : selectedJob?.title) }}</h2>
              <span class="modal-job-org">{{ selectedJob?.organization }}</span>
            </div>
          </div>

          <div class="modal-jobs-body mt-3">
            <div class="job-details-grid">
              
              <!-- Left Column Core details -->
              <div class="job-meta-column">
                <div class="meta-section">
                  <div class="section-lbl">Department Branch</div>
                  <div class="section-val">{{ selectedJob?.department }}</div>
                </div>

                <div class="meta-section mt-2">
                  <div class="section-lbl">Qualifications Needed</div>
                  <div class="section-val font-semibold"><i class="ti ti-school"></i> {{ selectedJob?.education_qualification }}</div>
                </div>

                <div class="meta-section mt-2">
                  <div class="section-lbl">Total Job Vacancies</div>
                  <div class="section-val highlight-val">{{ selectedJob?.vacancies.toLocaleString() }} Positions Open</div>
                </div>

                <div class="meta-section mt-2">
                  <div class="section-lbl">Job Application Fee</div>
                  <div class="section-val">{{ selectedJob?.application_fee }}</div>
                </div>
              </div>

              <!-- Right Column Verification documents checklist -->
              <div class="job-docs-column">
                <div class="section-lbl">Required Documents Checklist</div>
                <div class="jobs-documents-checklist mt-1">
                  <div 
                    v-for="(doc, idx) in formatDocuments(selectedJob?.required_documents)" 
                    :key="idx" 
                    class="doc-checklist-item"
                  >
                    <i class="ti ti-circle-check checked-icon"></i>
                    <span>{{ doc }}</span>
                  </div>
                </div>

                <div class="meta-section mt-3">
                  <div class="section-lbl">Application Lifespans</div>
                  <div class="lifespan-card mt-1">
                    <div>Start Date: <strong>{{ formatDate(selectedJob?.application_start_date) }}</strong></div>
                    <div>Last Date: <strong class="text-danger">{{ formatDate(selectedJob?.application_end_date) }}</strong></div>
                  </div>
                </div>
              </div>

            </div>
          </div>

          <!-- Modal Action Footer Buttons -->
          <div class="modal-actions-footer mt-4">
            <button type="button" class="btn-cancel" @click="closeJobDetails">Close Details</button>
            <button type="button" class="btn-submit apply-link-btn" @click="handleApplyAction(selectedJob)">
              <i class="ti ti-external-link"></i> Apply on Official Portal
            </button>
          </div>

        </div>
      </div>
    </Transition>

  </div>
</template>

<style scoped>
.jobs-explorer-container {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px 20px;
  box-sizing: border-box;
}

.jobs-header-hero {
  background: linear-gradient(135deg, rgba(26, 58, 107, 0.08) 0%, rgba(249, 115, 22, 0.05) 100%);
  border: 0.5px solid var(--border);
  border-radius: 12px;
  padding: 30px 24px;
  box-sizing: border-box;
}

.jobs-hero-title {
  font-size: 26px;
  font-weight: 700;
  color: var(--text);
  margin: 0 0 8px 0;
}

.jobs-hero-subtitle {
  font-size: 14px;
  color: var(--text2);
  margin: 0;
  line-height: 1.5;
}

/* Recommendations Section */
.recommendations-section {
  background-color: var(--bg);
  border: 0.5px solid var(--border);
  border-radius: 12px;
  padding: 18px 20px;
  box-sizing: border-box;
}

.sec-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 14px;
}

.sec-title-text {
  font-size: 15px;
  font-weight: 600;
  color: var(--text);
  margin: 0;
}

.recommendation-badge {
  font-size: 10px;
  font-weight: 600;
  background-color: var(--primary-light);
  color: var(--primary);
  padding: 3px 8px;
  border-radius: 100px;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.recommendations-slider-wrapper {
  overflow-x: auto;
  padding-bottom: 8px;
}

.recommendation-track {
  display: flex;
  gap: 16px;
  width: max-content;
}

.recommendation-card {
  width: 280px;
  background: linear-gradient(135deg, var(--bg2) 0%, var(--bg) 100%);
  border: 0.5px solid var(--border);
  border-radius: 8px;
  padding: 12px 14px;
  box-sizing: border-box;
  cursor: pointer;
  transition: all 0.2s ease;
}

.recommendation-card:hover {
  transform: translateY(-2px);
  border-color: var(--primary);
  box-shadow: 0 4px 12px rgba(0,0,0,0.04);
}

.rec-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 10px;
  color: var(--text2);
}

.rec-org {
  font-weight: 500;
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.rec-vacancies {
  background-color: rgba(22, 163, 74, 0.1);
  color: #16a34a;
  padding: 1px 5px;
  border-radius: 4px;
  font-weight: 600;
}

.rec-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--text);
  margin: 8px 0;
  line-height: 1.3;
  height: 34px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.rec-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 11px;
  color: var(--text2);
}

/* Guest recommendation signup */
.recommendations-guest-card {
  background: linear-gradient(135deg, var(--bg2) 0%, var(--bg) 100%);
  border: 0.5px dashed var(--primary);
  border-radius: 12px;
  padding: 20px 24px;
  box-sizing: border-box;
}

.guest-card-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.guest-icon {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  background-color: var(--primary-light);
  color: var(--primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;
}

.guest-card-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text);
  margin: 0 0 4px 0;
}

.guest-card-text {
  font-size: 12px;
  color: var(--text2);
  margin: 0;
  line-height: 1.4;
}

.btn-guest-login {
  background-color: var(--primary);
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 8px 14px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  margin-left: auto;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.btn-guest-login:hover {
  opacity: 0.95;
}

/* Main Explorer layout */
.explorer-filters-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
  flex-wrap: wrap;
}

.jobs-search-box {
  position: relative;
  width: 100%;
  max-width: 400px;
}

.search-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text2);
}

.jobs-search-input {
  width: 100%;
  padding: 8px 12px 8px 36px;
  border: 0.5px solid var(--border);
  border-radius: 8px;
  background-color: var(--bg);
  color: var(--text);
  outline: none;
  font-size: 13px;
  font-family: inherit;
  box-sizing: border-box;
}

.jobs-search-input:focus {
  border-color: var(--primary);
}

.filter-pills {
  display: flex;
  gap: 8px;
}

.filter-pill {
  border: 0.5px solid var(--border);
  border-radius: 100px;
  background-color: var(--bg);
  color: var(--text2);
  padding: 6px 14px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  font-family: inherit;
}

.filter-pill.active {
  background-color: var(--primary);
  color: #fff;
  border-color: var(--primary);
}

/* Jobs Grid */
.jobs-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

@media (max-width: 768px) {
  .jobs-grid {
    grid-template-columns: 1fr;
  }
}

.job-card {
  background-color: var(--bg);
  border: 0.5px solid var(--border);
  border-radius: 12px;
  padding: 18px 20px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  cursor: pointer;
  transition: all 0.2s ease;
}

.job-card:hover {
  border-color: var(--primary);
  box-shadow: 0 4px 12px rgba(0,0,0,0.04);
}

.job-card-top {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  font-size: 11px;
}

.job-org {
  color: var(--primary);
  font-weight: 600;
}

.job-department {
  color: var(--text2);
}

.closing-alert-badge {
  background-color: var(--danger-bg);
  color: var(--danger);
  font-size: 9px;
  font-weight: 600;
  padding: 1px 6px;
  border-radius: 4px;
  margin-left: auto;
  animation: pulse 1.5s infinite;
}

@keyframes pulse {
  0% { transform: scale(1); }
  50% { transform: scale(1.05); }
  100% { transform: scale(1); }
}

.job-card-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--text);
  margin: 12px 0 8px 0;
  line-height: 1.4;
}

.job-vacancies-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
  background-color: var(--bg2);
  border: 0.5px solid var(--border);
  border-radius: 8px;
  padding: 8px 10px;
}

.vacancy-item {
  display: flex;
  flex-direction: column;
}

.vacancy-item .lbl {
  font-size: 9px;
  color: var(--text2);
  text-transform: uppercase;
}

.vacancy-item .val {
  font-size: 11px;
  color: var(--text);
  margin-top: 1px;
}

.job-card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: var(--text2);
}

.end-date {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.btn-job-action {
  background: transparent;
  border: none;
  color: var(--primary);
  font-weight: 600;
  font-size: 12px;
  cursor: pointer;
  padding: 0;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-family: inherit;
}

.btn-job-action:hover {
  text-decoration: underline;
}

/* Modal Jobs override details */
.modal-jobs-header {
  display: flex;
  align-items: center;
  gap: 16px;
  border-bottom: 0.5px solid var(--border);
  padding-bottom: 16px;
}

.modal-jobs-icon-box {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background-color: var(--primary-light);
  color: var(--primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.modal-job-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--text);
  margin: 0;
}

.modal-job-org {
  font-size: 12px;
  color: var(--text2);
  font-weight: 500;
}

.job-details-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

@media (max-width: 500px) {
  .job-details-grid {
    grid-template-columns: 1fr;
  }
}

.job-meta-column, .job-docs-column {
  display: flex;
  flex-direction: column;
}

.section-lbl {
  font-size: 10px;
  font-weight: 600;
  color: var(--text2);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 4px;
}

.section-val {
  font-size: 13px;
  color: var(--text);
}

.section-val.highlight-val {
  font-size: 14px;
  color: #16a34a;
  font-weight: 600;
}

.jobs-documents-checklist {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 180px;
  overflow-y: auto;
}

.doc-checklist-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: var(--text);
  background-color: var(--bg2);
  border: 0.5px solid var(--border);
  padding: 6px 10px;
  border-radius: 6px;
}

.checked-icon {
  color: #16a34a;
  font-size: 14px;
}

.lifespan-card {
  background-color: var(--bg2);
  border: 0.5px solid var(--border);
  border-radius: 6px;
  padding: 8px 10px;
  font-size: 11px;
  color: var(--text);
  line-height: 1.4;
}

.btn-submit.apply-link-btn {
  background-color: var(--accent);
}

.jobs-loader {
  grid-column: 1 / -1;
  text-align: center;
  padding: 40px;
}

.jobs-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--border);
  border-top-color: var(--primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto 10px auto;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.no-jobs-card {
  grid-column: 1 / -1;
  text-align: center;
  padding: 40px;
  border: 0.5px dashed var(--border);
  border-radius: 12px;
  background-color: var(--bg);
}

.no-jobs-icon {
  font-size: 40px;
  color: var(--text2);
}
</style>
