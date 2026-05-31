<script setup>
import { ref, onMounted, computed } from 'vue'
import { jobsApi } from '../../api/jobs.js'
import { useAuthStore } from '../../stores/auth'
import { useUiStore } from '../../stores/ui'
import { API_BASE_URL } from '../../config.js'

const authStore = useAuthStore()
const uiStore = useUiStore()

const jobs = ref([])
const loading = ref(false)
const searchQuery = ref('')
const selectedQualFilter = ref('All')

// Modal Form states
const isModalOpen = ref(false)
const isEditing = ref(false)
const isTranslating = ref(false)
const activeJobId = ref(null)

const jobForm = ref({
  title: '', title_hi: '', title_mr: '',
  organization: '', department: '', vacancies: 10,
  education_qualification: 'Graduate',
  experience_required: 'None',
  required_documents: ['Aadhaar Card'],
  application_start_date: new Date().toISOString().split('T')[0],
  application_end_date: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
  official_website: 'https://gov.in',
  apply_link: 'https://gov.in',
  application_fee: 'General: ₹100, Reserved: Free',
  is_active: true
})

const documentInput = ref('')

onMounted(async () => {
  await fetchJobsList()
})

async function fetchJobsList() {
  loading.value = true
  try {
    const data = await jobsApi.adminFetchAllJobs()
    jobs.value = data || []
  } catch (err) {
    console.error(err)
    uiStore.showToast('Could not fetch jobs queue.', 'danger')
  } finally {
    loading.value = false
  }
}

// Compute Statistics Overview
const stats = computed(() => {
  const active = jobs.value.filter(j => j.is_active)
  const totalVacancies = active.reduce((sum, j) => sum + j.vacancies, 0)
  
  const today = new Date()
  const sevenDaysFromNow = new Date()
  sevenDaysFromNow.setDate(today.getDate() + 7)
  
  const expiring = active.filter(j => {
    const end = new Date(j.application_end_date)
    return end >= today && end <= sevenDaysFromNow
  }).length

  return [
    { num: jobs.value.length.toLocaleString(), label: 'Total Job Ads', icon: 'ti-id', class: 'blue' },
    { num: totalVacancies.toLocaleString(), label: 'Active Positions Open', icon: 'ti-briefcase', class: 'green' },
    { num: expiring.toLocaleString(), label: 'Closing Soon (7d)', icon: 'ti-alarm', class: 'orange' },
    { num: jobs.value.filter(j => !j.is_active).length.toLocaleString(), label: 'Suspended Listings', icon: 'ti-circle-x', class: 'red' }
  ]
})

const filteredJobs = computed(() => {
  let list = jobs.value
  
  if (selectedQualFilter.value !== 'All') {
    list = list.filter(j => j.education_qualification === selectedQualFilter.value)
  }

  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(j => 
      j.title.toLowerCase().includes(q) || 
      j.organization.toLowerCase().includes(q) || 
      j.department.toLowerCase().includes(q)
    )
  }

  return list
})

// Trigger auto-translations blur English Title
async function handleTitleBlur() {
  const text = jobForm.value.title
  if (!text || text.trim() === '') return
  
  isTranslating.value = true
  try {
    const hiRes = await fetch(`${API_BASE_URL}/api/translate?q=${encodeURIComponent(text)}&target=hi`, {
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    }).then(r => r.json())
    
    const mrRes = await fetch(`${API_BASE_URL}/api/translate?q=${encodeURIComponent(text)}&target=mr`, {
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    }).then(r => r.json())

    jobForm.value.title_hi = hiRes.translatedText || ''
    jobForm.value.title_mr = mrRes.translatedText || ''
  } catch (err) {
    console.error("Auto translation failed:", err)
  } finally {
    isTranslating.value = false
  }
}

function openCreateModal() {
  isEditing.value = false
  activeJobId.value = null
  jobForm.value = {
    title: '', title_hi: '', title_mr: '',
    organization: '', department: '', vacancies: 10,
    education_qualification: 'Graduate',
    experience_required: 'None',
    required_documents: ['Aadhaar Card'],
    application_start_date: new Date().toISOString().split('T')[0],
    application_end_date: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
    official_website: 'https://gov.in',
    apply_link: 'https://gov.in',
    application_fee: 'General: ₹100, Reserved: Free',
    is_active: true
  }
  isModalOpen.value = true
}

function openEditModal(job) {
  isEditing.value = true
  activeJobId.value = job.id
  jobForm.value = {
    title: job.title,
    title_hi: job.title_hi,
    title_mr: job.title_mr,
    organization: job.organization,
    department: job.department,
    vacancies: job.vacancies,
    education_qualification: job.education_qualification,
    experience_required: job.experience_required,
    required_documents: [...(job.required_documents || ['Aadhaar Card'])],
    application_start_date: job.application_start_date,
    application_end_date: job.application_end_date,
    official_website: job.official_website,
    apply_link: job.apply_link,
    application_fee: job.application_fee,
    is_active: job.is_active
  }
  isModalOpen.value = true
}

async function handleToggleActive(job) {
  try {
    const res = await jobsApi.adminDeleteJob(job.id)
    if (res.success) {
      uiStore.showToast(res.message, 'success')
      await fetchJobsList()
    }
  } catch (err) {
    console.error(err)
    uiStore.showToast('Failed to toggle status.', 'danger')
  }
}

function addDocument() {
  const doc = documentInput.value.trim()
  if (doc && !jobForm.value.required_documents.includes(doc)) {
    jobForm.value.required_documents.push(doc)
    documentInput.value = ''
  }
}

function removeDocument(idx) {
  jobForm.value.required_documents.splice(idx, 1)
}

async function submitForm() {
  try {
    let res
    if (isEditing.value) {
      res = await jobsApi.adminUpdateJob(activeJobId.value, jobForm.value)
    } else {
      res = await jobsApi.adminCreateJob(jobForm.value)
    }
    
    if (res.success) {
      uiStore.showToast(res.message, 'success')
      isModalOpen.value = false
      await fetchJobsList()
    }
  } catch (err) {
    console.error(err)
    uiStore.showToast('Could not save government job posting.', 'danger')
  }
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: 'numeric' })
}

defineExpose({
  openCreateModal
})
</script>

<template>
  <div class="admin-jobs-tab">
    
    <!-- Stats Grid Counters -->
    <div class="stats-grid">
      <div class="stat-card" v-for="s in stats" :key="s.label">
        <div :class="['stat-icon-box', s.class]">
          <i :class="['ti', s.icon]"></i>
        </div>
        <div class="stat-number">{{ s.num }}</div>
        <div class="stat-label">{{ s.label }}</div>
      </div>
    </div>

    <!-- Filter search and qualifications controls -->
    <div class="filter-panel mt-4">
      <div class="search-box">
        <i class="ti ti-search search-icon"></i>
        <input 
          type="text" 
          v-model="searchQuery" 
          placeholder="Search jobs, organizations, department..." 
          class="search-input"
        />
      </div>

      <div class="filter-controls">
        <select v-model="selectedQualFilter" class="filter-select">
          <option value="All">All Qualifications</option>
          <option value="10th Pass">10th Pass</option>
          <option value="12th Pass">12th Pass</option>
          <option value="Graduate">Graduate</option>
        </select>
        <button class="btn-create-direct" @click="openCreateModal">
          <i class="ti ti-plus"></i> Add Job Announcement
        </button>
      </div>
    </div>

    <!-- Table lists of job listings -->
    <div class="card mt-3">
      <div class="card-body p-0">
        <table class="data-table">
          <thead>
            <tr>
              <th>Job Information</th>
              <th>Dept & Branch</th>
              <th>Qualification</th>
              <th>Positions</th>
              <th>Fee Detail</th>
              <th>Closing Date</th>
              <th>Status</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="job in filteredJobs" :key="job.id" :class="{ 'inactive-row': !job.is_active }">
              <td>
                <div class="job-info-cell">
                  <div class="job-title-text font-semibold">{{ job.title }}</div>
                  <span class="job-org-sub">{{ job.organization }}</span>
                </div>
              </td>
              <td>{{ job.department }}</td>
              <td><span class="badge-qual">{{ job.education_qualification }}</span></td>
              <td class="font-semibold">{{ job.vacancies.toLocaleString() }}</td>
              <td class="fee-cell" :title="job.application_fee">{{ job.application_fee }}</td>
              <td class="date-cell">{{ formatDate(job.application_end_date) }}</td>
              <td>
                <span :class="['badge-status', job.is_active ? 'approved' : 'rejected']">
                  {{ job.is_active ? 'Active' : 'Suspended' }}
                </span>
              </td>
              <td>
                <div class="table-actions">
                  <button class="action-btn" title="Modify Job" @click="openEditModal(job)">
                    <i class="ti ti-pencil"></i>
                  </button>
                  <button 
                    :class="['action-btn', job.is_active ? 'reject-btn' : 'approve-btn']" 
                    :title="job.is_active ? 'Suspend Job' : 'Activate Job'" 
                    @click="handleToggleActive(job)"
                  >
                    <i :class="['ti', job.is_active ? 'ti-circle-x' : 'ti-circle-check']"></i>
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="filteredJobs.length === 0">
              <td colspan="8" class="no-records">
                <i class="ti ti-inbox no-records-icon"></i>
                <div class="no-records-text">No active job announcements found.</div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- CRUD MODAL DIALOG OVERLAY -->
    <div class="admin-modal-overlay" v-if="isModalOpen" @click.self="isModalOpen = false">
      <div class="admin-modal-box modal-jobs-dialog" style="max-width: 680px;">
        <button class="btn-close" @click="isModalOpen = false">×</button>
        <div class="modal-title-text">{{ isEditing ? 'Edit Government Job Posting' : 'Add New Government Job Posting' }}</div>

        <div class="modal-scrollable">
          <div v-if="isTranslating" class="translation-loader-banner">
            <i class="ti ti-loader rotate-spin"></i>
            <span>Generating translations into Hindi & Marathi...</span>
          </div>

          <form @submit.prevent="submitForm">
            
            <h4 class="form-sec-title">Core Job Advertisement Identity</h4>
            
            <div class="form-group">
              <label class="form-label">Job Title (English) *</label>
              <input 
                type="text" 
                class="form-input" 
                v-model="jobForm.title" 
                @blur="handleTitleBlur" 
                placeholder="e.g. RRB Assistant Loco Pilot Recruitment" 
                required 
              />
            </div>

            <div class="form-row">
              <div class="form-group">
                <label class="form-label">Hindi Translation Title</label>
                <input type="text" class="form-input" v-model="jobForm.title_hi" placeholder="आरआरबी सहायक लोको पायलट" />
              </div>
              <div class="form-group">
                <label class="form-label">Marathi Translation Title</label>
                <input type="text" class="form-input" v-model="jobForm.title_mr" placeholder="आरआरबी सहाय्यक लोको पायलट" />
              </div>
            </div>

            <div class="form-row mt-2">
              <div class="form-group">
                <label class="form-label">Conducting Organization *</label>
                <input type="text" class="form-input" v-model="jobForm.organization" placeholder="e.g. Railway Recruitment Board (RRB)" required />
              </div>
              <div class="form-group">
                <label class="form-label">Target Department *</label>
                <input type="text" class="form-input" v-model="jobForm.department" placeholder="e.g. Technical Department" required />
              </div>
            </div>

            <h4 class="form-sec-title mt-3">Eligibility & Positions Parameters</h4>

            <div class="form-row">
              <div class="form-group">
                <label class="form-label">Education Qualification Level *</label>
                <select class="form-input" v-model="jobForm.education_qualification">
                  <option value="10th Pass">10th Pass</option>
                  <option value="12th Pass">12th Pass</option>
                  <option value="Graduate">Graduate</option>
                </select>
              </div>
              <div class="form-group">
                <label class="form-label">Experience Required</label>
                <input type="text" class="form-input" v-model="jobForm.experience_required" placeholder="e.g. None or 2 Years Diploma" />
              </div>
            </div>

            <div class="form-row mt-2">
              <div class="form-group">
                <label class="form-label">Total Vacancies *</label>
                <input type="number" class="form-input" v-model="jobForm.vacancies" min="1" required />
              </div>
              <div class="form-group">
                <label class="form-label">Application Fee *</label>
                <input type="text" class="form-input" v-model="jobForm.application_fee" placeholder="e.g. Gen: ₹500, Reserved: Free" required />
              </div>
            </div>

            <div class="form-group mt-2">
              <label class="form-label">Documents Checklist Items *</label>
              <div class="tag-input-wrapper">
                <input 
                  type="text" 
                  class="form-input tag-input-box" 
                  v-model="documentInput" 
                  @keydown.enter.prevent="addDocument"
                  placeholder="Type document name and press Enter (Jaise: Caste Certificate)" 
                />
                <button type="button" class="btn-add-tag" @click="addDocument">Add</button>
              </div>

              <!-- Render Current tags checklist list -->
              <div class="tags-container mt-1" v-if="jobForm.required_documents.length > 0">
                <span class="tag-badge" v-for="(doc, idx) in jobForm.required_documents" :key="idx">
                  {{ doc }} <button type="button" class="tag-remove" @click="removeDocument(idx)">×</button>
                </span>
              </div>
            </div>

            <h4 class="form-sec-title mt-3">Gateways & Lifespans</h4>

            <div class="form-row">
              <div class="form-group">
                <label class="form-label">Start Date *</label>
                <input type="date" class="form-input" v-model="jobForm.application_start_date" required />
              </div>
              <div class="form-group">
                <label class="form-label">Last Date *</label>
                <input type="date" class="form-input" v-model="jobForm.application_end_date" required />
              </div>
            </div>

            <div class="form-row mt-2">
              <div class="form-group">
                <label class="form-label">Official Website Portal *</label>
                <input type="url" class="form-input" v-model="jobForm.official_website" placeholder="https://rrb.gov.in" required />
              </div>
              <div class="form-group">
                <label class="form-label">Apply Link Portal *</label>
                <input type="url" class="form-input" v-model="jobForm.apply_link" placeholder="https://rrb.gov.in/apply" required />
              </div>
            </div>

            <div class="form-group mt-3 flex-row-align">
              <input type="checkbox" id="job-active-cb" v-model="jobForm.is_active" class="jobs-cb" />
              <label for="job-active-cb" class="form-label m-0 cursor-pointer">Job advertisement active and open for applications</label>
            </div>

            <!-- Form Submit Buttons -->
            <button type="submit" class="submit-btn mt-3 font-semibold">
              <i class="ti ti-device-floppy"></i> Save Job Announcement
            </button>

          </form>
        </div>
      </div>
    </div>

  </div>
</template>

<style scoped>
.admin-jobs-tab {
  width: 100%;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
  width: 100%;
}

.stat-card {
  background-color: var(--bg);
  border: 0.5px solid var(--border);
  border-radius: 8px;
  padding: 14px 16px;
  box-sizing: border-box;
}

.stat-icon-box {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  margin-bottom: 10px;
}

.stat-icon-box.blue { background-color: var(--primary-light); color: var(--primary); }
.stat-icon-box.green { background-color: var(--success-bg); color: var(--success); }
.stat-icon-box.orange { background-color: var(--accent-light); color: var(--accent); }
.stat-icon-box.red { background-color: var(--danger-bg); color: var(--danger); }

.stat-number {
  font-size: 20px;
  font-weight: 700;
  color: var(--text);
}

.stat-label {
  font-size: 11px;
  color: var(--text2);
  margin-top: 2px;
}

.filter-panel {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.search-box {
  position: relative;
  flex: 1;
  max-width: 400px;
}

.search-icon {
  position: absolute;
  left: 10px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text2);
  font-size: 14px;
}

.search-input {
  width: 100%;
  padding: 7px 10px 7px 32px;
  border: 0.5px solid var(--border);
  border-radius: 6px;
  font-size: 13px;
  background-color: var(--bg);
  color: var(--text);
  outline: none;
  font-family: inherit;
  box-sizing: border-box;
}

.search-input:focus {
  border-color: var(--primary);
}

.filter-controls {
  display: flex;
  gap: 10px;
}

.filter-select {
  padding: 7px 10px;
  border: 0.5px solid var(--border);
  border-radius: 6px;
  font-size: 13px;
  background-color: var(--bg);
  color: var(--text);
  outline: none;
  font-family: inherit;
}

.btn-create-direct {
  background-color: var(--primary);
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 7px 14px;
  font-size: 12px;
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  font-family: inherit;
}

.btn-create-direct:hover {
  opacity: 0.95;
}

.card {
  background-color: var(--bg);
  border: 0.5px solid var(--border);
  border-radius: 12px;
  overflow: hidden;
  width: 100%;
}

.p-0 { padding: 0 !important; }

/* Table Styling */
.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th {
  background-color: var(--bg2);
  text-align: left;
  font-size: 11px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--text2);
  padding: 8px 12px;
  border-bottom: 0.5px solid var(--border);
}

.data-table td {
  padding: 10px 12px;
  border-bottom: 0.5px solid var(--border);
  font-size: 13px;
  vertical-align: middle;
}

.data-table tr:hover td {
  background-color: var(--bg2);
}

.inactive-row td {
  opacity: 0.6;
}

.job-info-cell {
  display: flex;
  flex-direction: column;
}

.job-title-text {
  color: var(--text);
  font-size: 13px;
  max-width: 250px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.job-org-sub {
  font-size: 11px;
  color: var(--text2);
}

.badge-qual {
  background-color: var(--primary-light);
  color: var(--primary);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
}

.fee-cell {
  color: var(--text2);
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.badge-status {
  display: inline-flex;
  align-items: center;
  padding: 2px 6px;
  border-radius: 100px;
  font-size: 11px;
  font-weight: 500;
}

.badge-status.approved { background-color: var(--success-bg); color: var(--success); }
.badge-status.rejected { background-color: var(--danger-bg); color: var(--danger); }

.table-actions {
  display: flex;
  gap: 6px;
}

.action-btn {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  border: 0.5px solid var(--border);
  background-color: var(--bg);
  color: var(--text2);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.15s ease;
  font-size: 14px;
  padding: 0;
  box-sizing: border-box;
}

.action-btn:hover {
  background-color: var(--bg2);
  color: var(--text);
}

.action-btn.approve-btn:hover {
  background-color: var(--success-bg);
  color: var(--success);
  border-color: var(--success);
}

.action-btn.reject-btn:hover {
  background-color: var(--danger-bg);
  color: var(--danger);
  border-color: var(--danger);
}

.no-records {
  text-align: center;
  padding: 40px !important;
}

.no-records-icon {
  font-size: 32px;
  color: var(--text2);
}

.no-records-text {
  margin-top: 8px;
  font-size: 13px;
  color: var(--text2);
}

/* Modals */
.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 12px;
}

.form-label {
  font-size: 12px;
  color: var(--text);
}

.form-input {
  padding: 8px 10px;
  border: 0.5px solid var(--border);
  border-radius: 6px;
  font-size: 13px;
  background-color: var(--bg);
  color: var(--text);
  outline: none;
  font-family: inherit;
  box-sizing: border-box;
  width: 100%;
}

.form-input:focus {
  border-color: var(--primary);
}

.form-row {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.form-sec-title {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  color: var(--primary);
  border-bottom: 0.5px solid var(--border);
  padding-bottom: 4px;
  margin: 14px 0 10px 0;
}

.flex-row-align {
  display: flex;
  align-items: center;
  gap: 8px;
}

.jobs-cb {
  width: 14px;
  height: 14px;
  cursor: pointer;
  accent-color: var(--primary);
}

.cursor-pointer {
  cursor: pointer;
}

/* Tags/Documents inputs styling */
.tag-input-wrapper {
  display: flex;
  gap: 8px;
}

.tag-input-box {
  flex-grow: 1;
}

.btn-add-tag {
  background-color: var(--primary-light);
  color: var(--primary);
  border: 0.5px solid var(--border);
  border-radius: 6px;
  padding: 0 14px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
}

.tags-container {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.tag-badge {
  background-color: var(--bg2);
  border: 0.5px solid var(--border);
  color: var(--text);
  padding: 3px 8px;
  border-radius: 4px;
  font-size: 11px;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.tag-remove {
  background: transparent;
  border: none;
  color: var(--text2);
  cursor: pointer;
  font-size: 12px;
  line-height: 1;
  padding: 0;
}

.tag-remove:hover {
  color: var(--danger);
}

.mt-4 { margin-top: 16px; }
.mt-3 { margin-top: 12px; }
.mt-2 { margin-top: 8px; }
.mt-1 { margin-top: 4px; }
.m-0 { margin: 0 !important; }
.font-semibold { font-weight: 600; }
</style>
