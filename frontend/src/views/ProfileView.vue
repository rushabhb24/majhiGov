<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useUiStore } from '../stores/ui'
import { useI18n } from 'vue-i18n'

const authStore = useAuthStore()
const uiStore = useUiStore()
const { t } = useI18n()

const isEditing = ref(false)
const saving = ref(false)

// Editable form data (clone of profile)
const editForm = ref({})

function startEditing() {
  const p = authStore.userProfile
  editForm.value = {
    full_name: p.full_name || '',
    date_of_birth: p.date_of_birth ? p.date_of_birth.split('T')[0] : '1995-01-01',
    gender: p.gender || 'Male',
    state: p.state || 'Maharashtra',
    district: p.district || 'Pune',
    caste_category: p.caste_category || 'General',
    annual_income: p.annual_income || 150000,
    occupation: p.occupation || 'Unemployed',
    employee_type: p.employee_type || 'Unemployed',
    education_level: p.education_level || 'Graduate',
    is_disabled: p.is_disabled || false
  }
  isEditing.value = true
}

function cancelEditing() {
  isEditing.value = false
}

async function saveProfile() {
  saving.value = true
  try {
    await authStore.updateProfile(editForm.value)
    isEditing.value = false
  } catch (err) {
    // Toast already shown by store
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  if (!authStore.userProfile) {
    authStore.fetchUserProfile()
  }
})
</script>

<template>
  <div class="tab-content animate-fade">
    <div class="card filter-panel">
      <div class="profile-header-row">
        <div>
          <h2 class="section-title">{{ t('myProfile') }}</h2>
          <p class="text-muted">{{ t('demographicDetails') }}</p>
        </div>
        <div class="profile-actions">
          <button v-if="!isEditing" class="btn btn-primary" @click="startEditing">
            ✏️ {{ t('editProfile') }}
          </button>
          <template v-else>
            <button class="btn btn-secondary" @click="cancelEditing" :disabled="saving">
              {{ t('cancelEdit') }}
            </button>
            <button class="btn btn-primary" @click="saveProfile" :disabled="saving">
              {{ saving ? t('submitting') : t('saveProfile') }}
            </button>
          </template>
        </div>
      </div>
    </div>

    <!-- Profile Display / Edit Form -->
    <div v-if="authStore.userProfile" class="profile-card card mt-4">
      <!-- View Mode -->
      <div v-if="!isEditing" class="profile-grid">
        <div class="profile-section">
          <h3 class="profile-section-title">{{ t('personalInfo') }}</h3>
          <div class="profile-field">
            <span class="field-label">{{ t('fullNameLabel') }}</span>
            <span class="field-value">{{ authStore.userProfile.full_name }}</span>
          </div>
          <div class="profile-field">
            <span class="field-label">{{ t('dobLabel') }}</span>
            <span class="field-value">{{ authStore.userProfile.date_of_birth ? new Date(authStore.userProfile.date_of_birth).toLocaleDateString() : '-' }}</span>
          </div>
          <div class="profile-field">
            <span class="field-label">{{ t('genderLabel') }}</span>
            <span class="field-value">{{ authStore.userProfile.gender }}</span>
          </div>
          <div class="profile-field">
            <span class="field-label">{{ t('stateLabel') }}</span>
            <span class="field-value">{{ authStore.userProfile.state }}</span>
          </div>
          <div class="profile-field">
            <span class="field-label">{{ t('districtLabel') }}</span>
            <span class="field-value">{{ authStore.userProfile.district }}</span>
          </div>
        </div>

        <div class="profile-section">
          <h3 class="profile-section-title">{{ t('demographicDetails') }}</h3>
          <div class="profile-field">
            <span class="field-label">{{ t('casteLabel') }}</span>
            <span class="field-value">{{ authStore.userProfile.caste_category }}</span>
          </div>
          <div class="profile-field">
            <span class="field-label">{{ t('incomeLabel') }}</span>
            <span class="field-value">₹{{ Number(authStore.userProfile.annual_income).toLocaleString() }}</span>
          </div>
          <div class="profile-field">
            <span class="field-label">{{ t('occupationLabel') }}</span>
            <span class="field-value">{{ authStore.userProfile.occupation }}</span>
          </div>
          <div class="profile-field">
            <span class="field-label">{{ t('employeeTypeLabel') }}</span>
            <span class="field-value">{{ authStore.userProfile.employee_type }}</span>
          </div>
          <div class="profile-field">
            <span class="field-label">{{ t('educationLabel') }}</span>
            <span class="field-value">{{ authStore.userProfile.education_level }}</span>
          </div>
          <div class="profile-field">
            <span class="field-label">{{ t('disabilityLabel') }}</span>
            <span class="field-value">{{ authStore.userProfile.is_disabled ? '✅ Yes' : '❌ No' }}</span>
          </div>
        </div>
      </div>

      <!-- Edit Mode -->
      <form v-else @submit.prevent="saveProfile" class="profile-edit-form">
        <div class="profile-section">
          <h3 class="profile-section-title">{{ t('personalInfo') }}</h3>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">{{ t('fullNameLabel') }} *</label>
              <input v-model="editForm.full_name" type="text" class="form-control" required />
            </div>
            <div class="form-group">
              <label class="form-label">{{ t('dobLabel') }} *</label>
              <input v-model="editForm.date_of_birth" type="date" class="form-control" required />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">{{ t('genderLabel') }} *</label>
              <select v-model="editForm.gender" class="form-control" required>
                <option value="Male">{{ t('maleOpt') }}</option>
                <option value="Female">{{ t('femaleOpt') }}</option>
                <option value="Other">{{ t('otherOpt') }}</option>
              </select>
            </div>
            <div class="form-group">
              <label class="form-label">{{ t('casteLabel') }} *</label>
              <select v-model="editForm.caste_category" class="form-control" required>
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
              <select v-model="editForm.state" class="form-control" required>
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
              <input v-model="editForm.district" type="text" class="form-control" required />
            </div>
          </div>
        </div>

        <hr class="divider" />

        <div class="profile-section">
          <h3 class="profile-section-title">{{ t('demographicDetails') }}</h3>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">{{ t('incomeLabel') }} *</label>
              <input v-model="editForm.annual_income" type="number" class="form-control" required />
            </div>
            <div class="form-group">
              <label class="form-label">{{ t('occupationLabel') }} *</label>
              <select v-model="editForm.occupation" class="form-control" required>
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
              <select v-model="editForm.employee_type" class="form-control" required>
                <option value="Unemployed">Unemployed</option>
                <option value="Private Employee">Private Sector</option>
                <option value="Government Employee">Government Sector</option>
                <option value="Self Employed">Self Employed</option>
              </select>
            </div>
            <div class="form-group">
              <label class="form-label">{{ t('educationLabel') }} *</label>
              <select v-model="editForm.education_level" class="form-control" required>
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
                id="editIsDisabled"
                v-model="editForm.is_disabled"
                style="width: 20px; height: 20px; cursor: pointer;"
              />
              <label for="editIsDisabled" style="cursor: pointer; font-size: 0.9rem;">Yes, I am differently-abled</label>
            </div>
          </div>
        </div>
      </form>
    </div>

    <!-- Loading state -->
    <div v-else class="loading-state text-center mt-4 card">
      <div class="spinner"></div>
      <p class="mt-4">{{ t('loading') }}</p>
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

.profile-header-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  flex-wrap: wrap;
}

.profile-actions {
  display: flex;
  gap: 8px;
}

.profile-card {
  padding: 28px;
}

.profile-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 32px;
}

.profile-section-title {
  font-family: var(--font-heading);
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--clr-primary);
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 2px solid var(--clr-border);
}

.profile-field {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid var(--clr-border);
}

.profile-field:last-child {
  border-bottom: none;
}

.field-label {
  font-size: 0.88rem;
  color: var(--clr-text-muted);
  font-weight: 500;
}

.field-value {
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--clr-text-main);
}

.profile-edit-form {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.profile-edit-form .profile-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>
