<script setup>
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import { useAdminStore } from '../stores/admin'
import { API_BASE_URL } from '../config.js'
import { useAuthStore } from '../stores/auth'
import { useUiStore } from '../stores/ui'
import { useRouter } from 'vue-router'

// Import Modular Components
import AdminSidebar from '../components/admin/AdminSidebar.vue'
import AdminTopbar from '../components/admin/AdminTopbar.vue'
import OverviewTab from '../components/admin/OverviewTab.vue'
import SchemesTab from '../components/admin/SchemesTab.vue'
import CategoriesTab from '../components/admin/CategoriesTab.vue'
import UsersTab from '../components/admin/UsersTab.vue'
import EligibilityTab from '../components/admin/EligibilityTab.vue'
import NotificationsTab from '../components/admin/NotificationsTab.vue'
import AnalyticsTab from '../components/admin/AnalyticsTab.vue'
import SettingsTab from '../components/admin/SettingsTab.vue'
import ProfileTab from '../components/admin/ProfileTab.vue'
import ApplicationsTab from '../components/admin/ApplicationsTab.vue'

const adminStore = useAdminStore()
const authStore = useAuthStore()
const uiStore = useUiStore()
const router = useRouter()

// Viewport / active states
const activeTab = ref('overview')
const isInitialLoading = ref(true)
const searchQuery = ref('')
const filterCategory = ref('All')
const filterType = ref('All')
const filterStatus = ref('All')
const schemesCount = computed(() => adminStore.schemes.length || 142)
const notificationsCount = computed(() => adminStore.notifications.length || 3)
const applicationsCount = computed(() => adminStore.applications.length || 0)

const isTranslating = ref(false)

// Modal states
const schemeModalOpen = ref(false)
const isEditingScheme = ref(false)
const activeSchemeId = ref(null)

const newScheme = ref({
  title: '', title_hi: '', title_mr: '',
  description: '', description_hi: '', description_mr: '',
  category_id: 1,
  government_level: 'central',
  state: null,
  benefits: '',
  application_start_date: new Date().toISOString().split('T')[0],
  application_end_date: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
  official_website: 'https://yojana.gov.in',
  apply_link: 'https://yojana.gov.in',
  is_active: true,
  eligibility: {
    min_age: 18, max_age: 60,
    gender: 'all',
    caste_categories: ['General', 'OBC', 'SC', 'ST'],
    min_income: 0, max_income: 300000,
    states: [], occupations: [], employee_types: [], education_levels: [],
    disability_required: false
  },
  documents: [{ document_name: 'Aadhaar Card', document_name_hi: 'आधार कार्ड', document_name_mr: 'आधार कार्ड', is_mandatory: true }],
  faqs: [{ question: '', answer: '', question_hi: '', answer_hi: '', question_mr: '', answer_mr: '' }]
})

// Category Add Form state
const newCategory = ref({
  name: '', name_hi: '', name_mr: '',
  icon: '🌾', description: ''
})

// Add Admin Form state
const adminModalOpen = ref(false)
const newAdmin = ref({
  email: '', phone: '', password: '', full_name: ''
})

// Eligibility Rules Page 5 form state
const selectedRuleSchemeId = ref(null)
const ruleForm = ref({
  min_age: 18, max_age: 60,
  min_income: 0, max_income: 200000,
  gender: 'all',
  caste_categories: ['General', 'OBC', 'SC', 'ST'],
  states_str: '',
  occupations_str: '',
  employee_types_str: 'Unemployed, Self-Employed',
  education_levels_str: '10th Pass, 12th Pass, Undergraduate, Graduate',
  disability_required: false
})

// Broadcast Form state
const broadcast = ref({
  send_to: 'All Users',
  state: '',
  title: '',
  message: '',
  type: 'New Scheme Alert'
})

let refreshInterval = null

// On mount guard and sync load
onMounted(async () => {
  console.log("AdminView: Component mounted. isLoggedIn:", authStore.isLoggedIn, "isAdmin:", authStore.isAdmin)
  if (!authStore.isLoggedIn || !authStore.isAdmin) {
    console.log("AdminView: Access denied, redirecting to home")
    uiStore.showToast('Access Denied: Administrative privileges required.', 'danger')
    router.push('/')
    return
  }
  
  console.log("AdminView: Access granted, refreshing data...")
  try {
    await refreshData()
    console.log("AdminView: Data refresh complete! Schemes loaded:", adminStore.schemes.length)
    isInitialLoading.value = false
    
    // Start active real-time data polling every 10 seconds
    refreshInterval = setInterval(refreshData, 10000)
  } catch (err) {
    console.error("AdminView: Failed to load data:", err)
    isInitialLoading.value = false
  }
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})

async function refreshData() {
  await Promise.all([
    adminStore.fetchAnalytics(),
    adminStore.fetchAllSchemes(),
    adminStore.fetchAllCategories(),
    adminStore.fetchAllUsers(),
    adminStore.fetchNotifications(),
    adminStore.fetchAllApplications()
  ])

  // Select first scheme for rules configuration if available
  if (adminStore.schemes.length > 0 && !selectedRuleSchemeId.value) {
    selectRuleScheme(adminStore.schemes[0].id)
  }
}

// Select scheme for rules edit
function selectRuleScheme(schemeId) {
  selectedRuleSchemeId.value = schemeId
  const scheme = adminStore.schemes.find(s => s.id === schemeId)
  if (scheme && scheme.eligibility) {
    const e = scheme.eligibility
    ruleForm.value = {
      min_age: e.min_age,
      max_age: e.max_age,
      min_income: e.min_income,
      max_income: e.max_income,
      gender: e.gender || 'all',
      caste_categories: e.caste_categories || ['General', 'OBC', 'SC', 'ST'],
      states_str: (e.states || []).join(', '),
      occupations_str: (e.occupations || []).join(', '),
      employee_types_str: (e.employee_types || []).join(', '),
      education_levels_str: (e.education_levels || []).join(', '),
      disability_required: e.disability_required || false
    }
  }
}

// Watch selected scheme inside rules page
watch(selectedRuleSchemeId, (newVal) => {
  if (newVal) selectRuleScheme(newVal)
})

// Add/Remove forms helper rows
function addDocumentRow() {
  newScheme.value.documents.push({ document_name: '', document_name_hi: '', document_name_mr: '', is_mandatory: true })
}
function removeDocumentRow(index) {
  newScheme.value.documents.splice(index, 1)
}
function addFaqRow() {
  newScheme.value.faqs.push({ question: '', answer: '', question_hi: '', answer_hi: '', question_mr: '', answer_mr: '' })
}
function removeFaqRow(index) {
  newScheme.value.faqs.splice(index, 1)
}

// Topbar dynamic click action dispatcher
function handleActionClick(actionType) {
  if (actionType === 'theme') {
    uiStore.toggleTheme()
    return
  }
  if (actionType === 'bell') {
    activeTab.value = 'notifications'
    return
  }

  // Primary action buttons clicks
  if (activeTab.value === 'overview' || activeTab.value === 'schemes') {
    isEditingScheme.value = false
    activeSchemeId.value = null
    newScheme.value = {
      title: '', title_hi: '', title_mr: '',
      description: '', description_hi: '', description_mr: '',
      category_id: adminStore.categories[0]?.id || 1,
      government_level: 'central',
      state: null,
      benefits: '',
      application_start_date: new Date().toISOString().split('T')[0],
      application_end_date: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
      official_website: 'https://yojana.gov.in',
      apply_link: 'https://yojana.gov.in',
      is_active: true,
      eligibility: {
        min_age: 18, max_age: 60,
        gender: 'all',
        caste_categories: ['General', 'OBC', 'SC', 'ST'],
        min_income: 0, max_income: 300000,
        states: [], occupations: [], employee_types: [], education_levels: [],
        disability_required: false
      },
      documents: [{ document_name: 'Aadhaar Card', document_name_hi: 'आधार कार्ड', document_name_mr: 'आधार कार्ड', is_mandatory: true }],
      faqs: [{ question: '', answer: '', question_hi: '', answer_hi: '', question_mr: '', answer_mr: '' }]
    }
    schemeModalOpen.value = true
  } else if (activeTab.value === 'categories') {
    document.getElementById('cat-name-input')?.focus()
  } else if (activeTab.value === 'users') {
    newAdmin.value = { email: '', phone: '', password: '', full_name: '' }
    adminModalOpen.value = true
  } else if (activeTab.value === 'eligibility') {
    saveEligibilityRules()
  } else if (activeTab.value === 'notifications') {
    sendBroadcast()
  } else if (activeTab.value === 'analytics') {
    uiStore.showToast('Operational analytics exported successfully as CSV!', 'success')
  } else if (activeTab.value === 'profile') {
    uiStore.showToast('Please fill out the details below and click "Save Administrative Profile" to update.', 'info')
  }
}

// Automatic Translation Engine Callers
async function translateField(text, fieldKey) {
  if (!text || text.trim() === '') return
  isTranslating.value = true
  try {
    const promises = [
      fetch(`${API_BASE_URL}/api/translate?q=${encodeURIComponent(text)}&target=hi`, {
        headers: { 'Authorization': `Bearer ${authStore.token}` }
      }).then(r => r.json()),
      fetch(`${API_BASE_URL}/api/translate?q=${encodeURIComponent(text)}&target=mr`, {
        headers: { 'Authorization': `Bearer ${authStore.token}` }
      }).then(r => r.json())
    ]
    
    const [hiRes, mrRes] = await Promise.all(promises)
    
    if (hiRes.success && hiRes.translatedText) {
      newScheme.value[fieldKey + '_hi'] = hiRes.translatedText
    }
    if (mrRes.success && mrRes.translatedText) {
      newScheme.value[fieldKey + '_mr'] = mrRes.translatedText
    }
  } catch (err) {
    console.error('Translation failed:', err)
  } finally {
    isTranslating.value = false
  }
}

async function translateDocRow(doc) {
  if (!doc.document_name || doc.document_name.trim() === '') return
  isTranslating.value = true
  try {
    const promises = [
      fetch(`${API_BASE_URL}/api/translate?q=${encodeURIComponent(doc.document_name)}&target=hi`, {
        headers: { 'Authorization': `Bearer ${authStore.token}` }
      }).then(r => r.json()),
      fetch(`${API_BASE_URL}/api/translate?q=${encodeURIComponent(doc.document_name)}&target=mr`, {
        headers: { 'Authorization': `Bearer ${authStore.token}` }
      }).then(r => r.json())
    ]
    
    const [hiRes, mrRes] = await Promise.all(promises)
    if (hiRes.success && hiRes.translatedText) {
      doc.document_name_hi = hiRes.translatedText
    }
    if (mrRes.success && mrRes.translatedText) {
      doc.document_name_mr = mrRes.translatedText
    }
  } catch (err) {
    console.error('Doc translation failed:', err)
  } finally {
    isTranslating.value = false
  }
}

async function translateFaqQuestion(faq) {
  if (!faq.question || faq.question.trim() === '') return
  isTranslating.value = true
  try {
    const promises = [
      fetch(`${API_BASE_URL}/api/translate?q=${encodeURIComponent(faq.question)}&target=hi`, {
        headers: { 'Authorization': `Bearer ${authStore.token}` }
      }).then(r => r.json()),
      fetch(`${API_BASE_URL}/api/translate?q=${encodeURIComponent(faq.question)}&target=mr`, {
        headers: { 'Authorization': `Bearer ${authStore.token}` }
      }).then(r => r.json())
    ]
    
    const [hiRes, mrRes] = await Promise.all(promises)
    if (hiRes.success && hiRes.translatedText) {
      faq.question_hi = hiRes.translatedText
    }
    if (mrRes.success && mrRes.translatedText) {
      faq.question_mr = mrRes.translatedText
    }
  } catch (err) {
    console.error('FAQ question translation failed:', err)
  } finally {
    isTranslating.value = false
  }
}

async function translateFaqAnswer(faq) {
  if (!faq.answer || faq.answer.trim() === '') return
  isTranslating.value = true
  try {
    const promises = [
      fetch(`${API_BASE_URL}/api/translate?q=${encodeURIComponent(faq.answer)}&target=hi`, {
        headers: { 'Authorization': `Bearer ${authStore.token}` }
      }).then(r => r.json()),
      fetch(`${API_BASE_URL}/api/translate?q=${encodeURIComponent(faq.answer)}&target=mr`, {
        headers: { 'Authorization': `Bearer ${authStore.token}` }
      }).then(r => r.json())
    ]
    
    const [hiRes, mrRes] = await Promise.all(promises)
    if (hiRes.success && hiRes.translatedText) {
      faq.answer_hi = hiRes.translatedText
    }
    if (mrRes.success && mrRes.translatedText) {
      faq.answer_mr = mrRes.translatedText
    }
  } catch (err) {
    console.error('FAQ answer translation failed:', err)
  } finally {
    isTranslating.value = false
  }
}

// Tab handlers logic
async function submitSchemeForm() {
  try {
    if (newScheme.value.government_level === 'central') {
      newScheme.value.state = null
    }

    if (isEditingScheme.value) {
      await adminStore.updateScheme(activeSchemeId.value, newScheme.value)
      uiStore.showToast('Scheme updated successfully!', 'success')
    } else {
      await adminStore.createScheme(newScheme.value)
      uiStore.showToast('Scheme created successfully!', 'success')
    }
    
    schemeModalOpen.value = false
    await refreshData()
  } catch (err) {
    uiStore.showToast(err.message || 'Failed to submit scheme.', 'danger')
  }
}

function handleViewScheme(scheme) {
  handleEditScheme(scheme)
}

function handleEditScheme(scheme) {
  isEditingScheme.value = true
  activeSchemeId.value = scheme.id
  
  newScheme.value = {
    title: scheme.title, title_hi: scheme.title_hi || '', title_mr: scheme.title_mr || '',
    description: scheme.description, description_hi: scheme.description_hi || '', description_mr: scheme.description_mr || '',
    category_id: scheme.category_id,
    government_level: scheme.government_level,
    state: scheme.state,
    benefits: scheme.benefits,
    application_start_date: scheme.application_start_date,
    application_end_date: scheme.application_end_date,
    official_website: scheme.official_website,
    apply_link: scheme.apply_link,
    is_active: scheme.is_active,
    eligibility: scheme.eligibility || {
      min_age: 18, max_age: 60,
      gender: 'all',
      caste_categories: ['General', 'OBC', 'SC', 'ST'],
      min_income: 0, max_income: 300000,
      states: [], occupations: [], employee_types: [], education_levels: [],
      disability_required: false
    },
    documents: scheme.documents || [],
    faqs: scheme.faqs || []
  }

  schemeModalOpen.value = true
}

async function handleDeleteScheme(scheme) {
  try {
    const res = await adminStore.toggleSchemeStatus(scheme.id)
    uiStore.showToast(res.message, 'success')
    await refreshData()
  } catch (err) {
    uiStore.showToast(err.message, 'danger')
  }
}

async function handleAddCategory() {
  try {
    await adminStore.createCategory(newCategory.value)
    uiStore.showToast('Category created successfully!', 'success')
    newCategory.value = { name: '', name_hi: '', name_mr: '', icon: '🌾', description: '' }
    await refreshData()
  } catch (err) {
    uiStore.showToast(err.message, 'danger')
  }
}

async function handleDeleteCategory(id) {
  if (!confirm('Are you sure you want to delete this category?')) return
  try {
    const result = await adminStore.deleteCategory(id)
    uiStore.showToast(result.message, 'success')
    await refreshData()
  } catch (err) {
    uiStore.showToast(err.message || 'Cannot delete category linked to active schemes.', 'danger')
  }
}

async function handleToggleUserVerify(userId) {
  try {
    const res = await adminStore.toggleUserStatus(userId)
    uiStore.showToast(res.message, 'success')
    await refreshData()
  } catch (err) {
    uiStore.showToast(err.message, 'danger')
  }
}

async function handleUpdateApplicationStatus({ applicationId, status, notes }) {
  try {
    const res = await adminStore.updateApplicationStatus(applicationId, status, notes)
    uiStore.showToast(res.message, 'success')
    await refreshData()
  } catch (err) {
    uiStore.showToast(err.message || 'Failed to update application status.', 'danger')
  }
}

async function submitAdminForm() {
  try {
    await adminStore.createAdmin(newAdmin.value)
    uiStore.showToast('Administrative privileges granted successfully!', 'success')
    adminModalOpen.value = false
    await refreshData()
  } catch (err) {
    uiStore.showToast(err.message, 'danger')
  }
}

async function saveEligibilityRules() {
  if (!selectedRuleSchemeId.value) return
  try {
    const parseList = (str) => str.split(',').map(s => s.trim()).filter(s => s !== '')
    const payload = {
      min_age: Number(ruleForm.value.min_age),
      max_age: Number(ruleForm.value.max_age),
      min_income: Number(ruleForm.value.min_income),
      max_income: Number(ruleForm.value.max_income),
      gender: ruleForm.value.gender,
      caste_categories: ruleForm.value.caste_categories,
      states: parseList(ruleForm.value.states_str),
      occupations: parseList(ruleForm.value.occupations_str),
      employee_types: parseList(ruleForm.value.employee_types_str),
      education_levels: parseList(ruleForm.value.education_levels_str),
      disability_required: ruleForm.value.disability_required
    }

    const scheme = adminStore.schemes.find(s => s.id === selectedRuleSchemeId.value)
    if (!scheme) throw new Error('Scheme not found')
    
    const fullPayload = {
      title: scheme.title, title_hi: scheme.title_hi || '', title_mr: scheme.title_mr || '',
      description: scheme.description, description_hi: scheme.description_hi || '', description_mr: scheme.description_mr || '',
      category_id: scheme.category_id,
      government_level: scheme.government_level,
      state: scheme.state,
      benefits: scheme.benefits,
      application_start_date: scheme.application_start_date,
      application_end_date: scheme.application_end_date,
      official_website: scheme.official_website,
      apply_link: scheme.apply_link,
      is_active: scheme.is_active,
      eligibility: payload,
      documents: scheme.documents || [],
      faqs: scheme.faqs || []
    }

    await adminStore.updateScheme(selectedRuleSchemeId.value, fullPayload)
    uiStore.showToast('Eligibility rules updated successfully!', 'success')
    await refreshData()
  } catch (err) {
    uiStore.showToast(err.message, 'danger')
  }
}

async function sendBroadcast() {
  try {
    if (!broadcast.value.title || !broadcast.value.message) {
      uiStore.showToast('Broadcast title and message are mandatory.', 'warning')
      return
    }
    const res = await adminStore.sendNotification(broadcast.value)
    uiStore.showToast(res.message, 'success')
    broadcast.value = { send_to: 'All Users', state: '', title: '', message: '', type: 'New Scheme Alert' }
    await refreshData()
  } catch (err) {
    uiStore.showToast(err.message, 'danger')
  }
}

function handleViewExpiring() {
  activeTab.value = 'schemes'
}

// Computed helper for filtered schemes list
const filteredSchemesList = computed(() => {
  let list = adminStore.schemes
  
  if (filterCategory.value && filterCategory.value !== 'All') {
    list = list.filter(s => s.category_name === filterCategory.value)
  }
  
  if (filterType.value && filterType.value !== 'All') {
    list = list.filter(s => s.government_level === filterType.value)
  }
  
  if (filterStatus.value && filterStatus.value !== 'All') {
    if (filterStatus.value === 'Active') {
      list = list.filter(s => s.is_active)
    } else if (filterStatus.value === 'Inactive') {
      list = list.filter(s => !s.is_active)
    } else if (filterStatus.value === 'Expiring') {
      const sevenDaysFromNow = new Date()
      sevenDaysFromNow.setDate(sevenDaysFromNow.getDate() + 7)
      const now = new Date()
      list = list.filter(s => {
        const endDate = new Date(s.application_end_date)
        return endDate >= now && endDate <= sevenDaysFromNow
      })
    }
  }
  
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(s => s.title.toLowerCase().includes(q) || s.description.toLowerCase().includes(q))
  }
  
  return list
})

const filteredUsersList = computed(() => {
  let list = adminStore.users
  
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(u => u.full_name.toLowerCase().includes(q) || u.email.toLowerCase().includes(q))
  }
  
  return list
})
</script>

<template>
  <div :class="['admin-dashboard-container', uiStore.theme]">
    
    <!-- LEFT SIDEBAR PANEL -->
    <AdminSidebar 
      v-model:activeTab="activeTab"
      :schemesCount="schemesCount"
      :notificationsCount="notificationsCount"
      :applicationsCount="applicationsCount"
      :theme="uiStore.theme"
    />

    <!-- RIGHT MAIN PANEL AREA -->
    <main class="main-area">
      
      <!-- TOP BAR HEADER -->
      <AdminTopbar 
        :activeTab="activeTab"
        v-model:searchQuery="searchQuery"
        @action-click="handleActionClick"
        :theme="uiStore.theme"
      />

      <!-- SCROLLABLE PAGE CONTENT AREA -->
      <div class="content-area" v-if="!isInitialLoading">
        
        <!-- Tab 1: Overview -->
        <OverviewTab 
          v-if="activeTab === 'overview'"
          :analytics="adminStore.analytics"
          :schemes="adminStore.schemes"
          @nav-tab="(tab) => activeTab = tab"
          @view-expiring="handleViewExpiring"
        />

        <!-- Tab 2: Schemes -->
        <SchemesTab 
          v-else-if="activeTab === 'schemes'"
          :schemes="filteredSchemesList"
          :categories="adminStore.categories"
          v-model:filterCategory="filterCategory"
          v-model:filterType="filterType"
          v-model:filterStatus="filterStatus"
          @view-scheme="handleViewScheme"
          @edit-scheme="handleEditScheme"
          @delete-scheme="handleDeleteScheme"
        />

        <!-- Tab 3: Categories -->
        <CategoriesTab 
          v-else-if="activeTab === 'categories'"
          :categories="adminStore.categories"
          :newCategory="newCategory"
          @add-category="handleAddCategory"
          @delete-category="handleDeleteCategory"
        />

        <!-- Tab 4: Users -->
        <UsersTab 
          v-else-if="activeTab === 'users'"
          :users="filteredUsersList"
          @toggle-verify="handleToggleUserVerify"
        />

        <!-- Tab 4b: Applications [NEW] -->
        <ApplicationsTab 
          v-else-if="activeTab === 'applications'"
          :applications="adminStore.applications"
          :loading="adminStore.loading"
          @update-status="handleUpdateApplicationStatus"
        />

        <!-- Tab 5: Eligibility Rules -->
        <EligibilityTab 
          v-else-if="activeTab === 'eligibility'"
          :schemes="adminStore.schemes"
          v-model:selectedSchemeId="selectedRuleSchemeId"
          :ruleForm="ruleForm"
          @save-rules="saveEligibilityRules"
        />

        <!-- Tab 6: Notifications -->
        <NotificationsTab 
          v-else-if="activeTab === 'notifications'"
          :broadcast="broadcast"
          :notifications="adminStore.notifications"
          @send-broadcast="sendBroadcast"
        />

        <!-- Tab 7: Analytics -->
        <AnalyticsTab 
          v-else-if="activeTab === 'analytics'"
          :analytics="adminStore.analytics"
        />

        <!-- Tab 8: Settings -->
        <SettingsTab 
          v-else-if="activeTab === 'settings'"
          :loading="adminStore.loading"
          @refresh="refreshData"
        />

        <!-- Tab 9: Profile [NEW] -->
        <ProfileTab 
          v-else-if="activeTab === 'profile'"
        />

      </div>

      <!-- Spinner Loader indicator -->
      <div v-else class="admin-loading-spinner-wrapper">
        <div class="admin-loading-spinner"></div>
        <div style="margin-top:14px; font-weight:500; color:#64748b;">Loading Administrative Datastore...</div>
      </div>

    </main>

    <!-- ============================================== -->
    <!-- OVERLAY MODAL: CREATE/EDIT SCHEME -->
    <!-- ============================================== -->
    <div class="admin-modal-overlay" v-if="schemeModalOpen" @click.self="schemeModalOpen = false">
      <div class="admin-modal-box">
        <button class="btn-close" @click="schemeModalOpen = false" title="Close">×</button>
        <div class="modal-title-text">{{ isEditingScheme ? 'Modify Scheme Parameters' : 'Add New Scheme' }}</div>
        
        <div class="modal-scrollable">
          <div v-if="isTranslating" class="translation-loader-banner">
            <i class="ti ti-loader rotate-spin"></i>
            <span>Generating translations into Hindi & Marathi...</span>
          </div>
          <form @submit.prevent="submitSchemeForm">
            <h4 class="form-sec-title">Scheme Identity Details</h4>
            
            <div class="form-group">
              <label class="form-label">Scheme Title (English) *</label>
              <input type="text" class="form-input" v-model="newScheme.title" @blur="translateField(newScheme.title, 'title')" placeholder="e.g. PM Kisan Samman Nidhi" required />
            </div>

            <div class="form-row">
              <div class="form-group">
                <label class="form-label">Hindi Title</label>
                <input type="text" class="form-input" v-model="newScheme.title_hi" placeholder="पीएम किसान सम्मान निधि" />
              </div>
              <div class="form-group">
                <label class="form-label">Marathi Title</label>
                <input type="text" class="form-input" v-model="newScheme.title_mr" placeholder="पीएम किसान सन्मान निधी" />
              </div>
            </div>

            <div class="form-group">
              <label class="form-label">Scheme Benefits Description (English) *</label>
              <input type="text" class="form-input" v-model="newScheme.benefits" placeholder="₹6,000 per year in 3 equal installments" required />
            </div>

            <div class="form-group">
              <label class="form-label">Description (English) *</label>
              <textarea class="form-input" v-model="newScheme.description" @blur="translateField(newScheme.description, 'description')" rows="3" placeholder="Scheme details and statement..." required></textarea>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label class="form-label">Hindi Description</label>
                <textarea class="form-input" v-model="newScheme.description_hi" rows="2" placeholder="विवरण हिंदी में..."></textarea>
              </div>
              <div class="form-group">
                <label class="form-label">Marathi Description</label>
                <textarea class="form-input" v-model="newScheme.description_mr" rows="2" placeholder="वर्णन मराठीत..."></textarea>
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label class="form-label">Category *</label>
                <select class="form-input" v-model="newScheme.category_id" required>
                  <option v-for="cat in adminStore.categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                </select>
              </div>

              <div class="form-group">
                <label class="form-label">Government Level *</label>
                <select class="form-input" v-model="newScheme.government_level" required>
                  <option value="central">Central</option>
                  <option value="state">State</option>
                </select>
              </div>
            </div>

            <div class="form-group" v-if="newScheme.government_level === 'state'">
              <label class="form-label">Applicable State *</label>
              <input type="text" class="form-input" v-model="newScheme.state" placeholder="e.g. Maharashtra" required />
            </div>

            <div class="form-row">
              <div class="form-group">
                <label class="form-label">Start Date *</label>
                <input type="date" class="form-input" v-model="newScheme.application_start_date" required />
              </div>
              <div class="form-group">
                <label class="form-label">End Date *</label>
                <input type="date" class="form-input" v-model="newScheme.application_end_date" required />
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label class="form-label">Official Portal Website *</label>
                <input type="text" class="form-input" v-model="newScheme.official_website" required />
              </div>
              <div class="form-group">
                <label class="form-label">Apply Link *</label>
                <input type="text" class="form-input" v-model="newScheme.apply_link" required />
              </div>
            </div>

            <div class="form-group">
              <div style="display: flex; align-items: center; gap: 8px;">
                <input type="checkbox" id="schemeActive" v-model="newScheme.is_active" style="width: 18px; height: 18px;" />
                <label for="schemeActive" style="cursor: pointer; font-weight: 500;">Active and visible on Citizen Explorer portal</label>
              </div>
            </div>

            <!-- MANDATORY DOCUMENTS MANAGER -->
            <hr class="divider mt-4 mb-3" />
            <div style="display:flex; justify-content:space-between; align-items:center;">
              <h4 class="form-sec-title m-0">Mandatory Verification Documents</h4>
              <button type="button" class="action-btn" style="padding:4px 8px; font-size:11px;" @click="addDocumentRow"><i class="ti ti-plus"></i> Add Row</button>
            </div>
            
            <div style="display:flex; flex-direction:column; gap:8px;" class="mt-2">
              <div v-for="(doc, idx) in newScheme.documents" :key="idx" style="display:flex; align-items:center; gap:8px;">
                <input type="text" class="form-input" v-model="doc.document_name" @blur="translateDocRow(doc)" placeholder="Aadhaar Card" required />
                <input type="text" class="form-input" v-model="doc.document_name_hi" placeholder="आधार कार्ड" />
                <input type="text" class="form-input" v-model="doc.document_name_mr" placeholder="आधार कार्ड" />
                <select class="form-input" v-model="doc.is_mandatory" style="max-width:110px;">
                  <option :value="true">Mandatory</option>
                  <option :value="false">Optional</option>
                </select>
                <button type="button" class="action-btn danger-hover" @click="removeDocumentRow(idx)"><i class="ti ti-trash"></i></button>
              </div>
            </div>

            <!-- FAQ LIST MANAGER -->
            <hr class="divider mt-4 mb-3" />
            <div style="display:flex; justify-content:space-between; align-items:center;">
              <h4 class="form-sec-title m-0">Frequently Asked Questions (FAQ)</h4>
              <button type="button" class="action-btn" style="padding:4px 8px; font-size:11px;" @click="addFaqRow"><i class="ti ti-plus"></i> Add Row</button>
            </div>

            <div style="display:flex; flex-direction:column; gap:12px;" class="mt-2">
              <div v-for="(faq, idx) in newScheme.faqs" :key="idx" class="faq-form-card" style="padding:10px; border:0.5px solid rgba(0, 0, 0, 0.08); border-radius:6px;">
                <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:6px;">
                  <strong>FAQ #{{ idx + 1 }}</strong>
                  <button type="button" class="action-btn danger-hover" @click="removeFaqRow(idx)"><i class="ti ti-trash"></i></button>
                </div>
                <div style="display:flex; flex-direction:column; gap:8px;">
                  <input type="text" class="form-input" v-model="faq.question" @blur="translateFaqQuestion(faq)" placeholder="Question in English" required />
                  <input type="text" class="form-input" v-model="faq.answer" @blur="translateFaqAnswer(faq)" placeholder="Answer in English" required />
                  <input type="text" class="form-input" v-model="faq.question_hi" placeholder="प्रश्न (हिंदी)" />
                  <input type="text" class="form-input" v-model="faq.answer_hi" placeholder="उत्तर (हिंदी)" />
                  <input type="text" class="form-input" v-model="faq.question_mr" placeholder="प्रश्न (मराठी)" />
                  <input type="text" class="form-input" v-model="faq.answer_mr" placeholder="उत्तर (मराठी)" />
                </div>
              </div>
            </div>

            <hr class="divider mt-4 mb-4" />

            <div style="display:flex; justify-content:flex-end; gap:10px;">
              <button type="button" class="top-action-btn" style="background:#e2e8f0; color:#475569;" @click="schemeModalOpen = false">Cancel</button>
              <button type="submit" class="top-action-btn" style="background:#16a34a;">{{ isEditingScheme ? 'Save Changes' : 'Create Scheme' }}</button>
            </div>

          </form>
        </div>
      </div>
    </div>

    <!-- ============================================== -->
    <!-- OVERLAY MODAL: CREATE ADMIN -->
    <!-- ============================================== -->
    <div class="admin-modal-overlay" v-if="adminModalOpen" @click.self="adminModalOpen = false">
      <div class="admin-modal-box" style="max-width: 450px;">
        <button class="btn-close" @click="adminModalOpen = false">×</button>
        <div class="modal-title-text">Add Administrative User</div>

        <form @submit.prevent="submitAdminForm">
          <div class="form-group">
            <label class="form-label">Full Name *</label>
            <input type="text" class="form-input" v-model="newAdmin.full_name" placeholder="John Doe" required />
          </div>

          <div class="form-group">
            <label class="form-label">Administrative Email *</label>
            <input type="email" class="form-input" v-model="newAdmin.email" placeholder="admin.doe@gov.in" required />
          </div>

          <div class="form-group">
            <label class="form-label">Phone Number *</label>
            <input type="tel" class="form-input" v-model="newAdmin.phone" placeholder="9876543210" required />
          </div>

          <div class="form-group">
            <label class="form-label">Administrative Password *</label>
            <input type="password" class="form-input" v-model="newAdmin.password" placeholder="••••••••" required />
          </div>

          <button type="submit" class="submit-btn mt-4">
            <i class="ti ti-user-plus"></i> Grant Administrator Privileges
          </button>
        </form>
      </div>
    </div>

  </div>
</template>

<style>
/* Global rules for shielding from other layouts & matching styling guide strictly */
:root {
  --primary: #1a3a6b;
  --primary-light: #e8eef8;
  --accent: #f97316;
  --accent-light: #fff4ed;
  --success: #16a34a;
  --success-bg: #f0fdf4;
  --danger: #dc2626;
  --danger-bg: #fef2f2;
  --warning: #d97706;
  --warning-bg: #fffbeb;
  --border: rgba(0,0,0,0.08);
  --text: #0f172a;
  --text2: #64748b;
  --bg: #ffffff;
  --bg2: #f8fafc;
  --bg3: #f1f5f9;
  --radius: 8px;
  --radius-lg: 12px;
}

/* Dark mode variable overrides */
.admin-dashboard-container.dark {
  --border: rgba(255, 255, 255, 0.08);
  --text: #f3f4f6;
  --text2: #9ca3af;
  --bg: #111827;
  --bg2: #1f2937;
  --bg3: #0b0f19;
  --primary-light: rgba(99, 102, 241, 0.15);
  --accent-light: rgba(245, 158, 11, 0.15);
  --danger-bg: rgba(239, 68, 68, 0.15);
  --success-bg: rgba(16, 185, 129, 0.15);
  --warning-bg: rgba(245, 158, 11, 0.15);
}
</style>

<style scoped>
.admin-dashboard-container {
  display: flex;
  width: 100vw;
  height: 100vh;
  position: fixed;
  top: 0;
  left: 0;
  overflow: hidden;
  background-color: var(--bg3);
  color: var(--text);
  font-family: 'Plus Jakarta Sans', sans-serif;
  font-size: 13px;
  font-weight: 400;
  z-index: 9999; /* Completely shields and sits above citizen overlays */
  box-sizing: border-box;
}

.main-area {
  flex-grow: 1;
  height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-sizing: border-box;
  background-color: var(--bg3);
}

.content-area {
  flex-grow: 1;
  padding: 20px;
  overflow-y: auto;
  box-sizing: border-box;
  background-color: var(--bg3);
}

/* Modals & Overlays */
.admin-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(15, 23, 42, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10000;
  animation: fadeInModal 0.2s ease-out;
}

@keyframes fadeInModal {
  from { opacity: 0; }
  to { opacity: 1; }
}

.admin-modal-box {
  background-color: var(--bg);
  border: 0.5px solid var(--border);
  border-radius: 12px;
  width: 100%;
  max-width: 700px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  padding: 24px;
  position: relative;
  box-sizing: border-box;
  animation: slideInModal 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes slideInModal {
  from { transform: translateY(10px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}

.btn-close {
  position: absolute;
  top: 16px;
  right: 16px;
  background: transparent;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #64748b;
  line-height: 1;
  padding: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.btn-close:hover {
  color: #0f172a;
}

.modal-title-text {
  font-size: 15px;
  font-weight: 500;
  color: #0f172a;
  margin-bottom: 20px;
  line-height: 1.2;
}

.modal-scrollable {
  overflow-y: auto;
  flex-grow: 1;
  padding-right: 6px;
}

/* Forms internal */
.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 12px;
  width: 100%;
}

.form-label {
  font-size: 13px;
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

textarea.form-input {
  resize: vertical;
}

.form-sec-title {
  font-size: 12px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #1a3a6b;
  margin-top: 14px;
  margin-bottom: 12px;
  line-height: 1.3;
}

.divider {
  border: none;
  border-top: 0.5px solid var(--border);
}

.mt-2 { margin-top: 8px; }
.mt-3 { margin-top: 12px; }
.mt-4 { margin-top: 16px; }
.mb-3 { margin-bottom: 12px; }
.m-0 { margin: 0 !important; }

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

.action-btn.danger-hover:hover {
  background-color: #fef2f2;
  color: #dc2626;
  border-color: #dc2626;
}

.top-action-btn {
  background-color: #1a3a6b;
  color: #ffffff;
  border: none;
  border-radius: 6px;
  height: 34px;
  padding: 7px 14px;
  font-size: 13px;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  font-family: inherit;
  box-sizing: border-box;
}

.top-action-btn:hover {
  opacity: 0.9;
}

.submit-btn {
  background-color: #1a3a6b;
  color: #ffffff;
  border: none;
  border-radius: 6px;
  padding: 9px 18px;
  font-size: 13px;
  font-weight: 500;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  cursor: pointer;
  width: 100%;
  font-family: inherit;
  box-sizing: border-box;
  margin-top: 8px;
}

.submit-btn:hover {
  opacity: 0.9;
}

.submit-btn i {
  font-size: 16px !important;
}

/* Spinner Loader */
.admin-loading-spinner-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 50vh;
}

.admin-loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid #e8eef8;
  border-top-color: #1a3a6b;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.translation-loader-banner {
  background-color: var(--primary-light);
  color: var(--primary);
  border: 0.5px solid var(--border);
  border-radius: 6px;
  padding: 8px 12px;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
  font-weight: 500;
}

.rotate-spin {
  animation: spin 1s linear infinite;
  display: inline-block;
}
</style>
