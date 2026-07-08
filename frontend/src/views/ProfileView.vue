<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useUiStore } from '../stores/ui'
import { useI18n } from 'vue-i18n'
import AppCard from '../components/ui/AppCard.vue'
import AppBadge from '../components/ui/AppBadge.vue'
import AppButton from '../components/ui/AppButton.vue'
import AppInput from '../components/ui/AppInput.vue'
import AppLabel from '../components/ui/AppLabel.vue'
import AppSelect from '../components/ui/AppSelect.vue'

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
    is_disabled: p.is_disabled || false,
    aadhaar: p.aadhaar || ''
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

function getInitials(name) {
  if (!name) return 'C'
  const parts = name.split(' ')
  if (parts.length > 1) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return name.substring(0, 2).toUpperCase()
}
</script>

<template>
  <div class="tw-max-w-4xl tw-mx-auto tw-px-4 tw-sm:px-6 tw-lg:px-8 tw-py-8">
    
    <!-- Header panel -->
    <div class="glass tw-p-6 tw-rounded-2xl tw-mb-6">
      <div class="tw-flex tw-justify-between tw-items-center tw-flex-wrap tw-gap-4">
        <div class="tw-flex tw-items-center tw-gap-4">
          <!-- Monogram avatar -->
          <div class="tw-w-14 tw-h-14 tw-rounded-2xl tw-bg-primary tw-text-white tw-font-heading tw-font-black tw-text-xl tw-flex tw-items-center tw-justify-center tw-shadow-glow">
            {{ authStore.userProfile ? getInitials(authStore.userProfile.full_name) : 'C' }}
          </div>
          <div>
            <h2 class="tw-font-heading tw-font-bold tw-text-xl tw-text-foreground tw-m-0">
              {{ authStore.userProfile ? authStore.userProfile.full_name : t('myProfile') }}
            </h2>
            <p class="tw-text-xs tw-text-muted-foreground tw-mt-1 tw-m-0">
              {{ t('demographicDetails') || 'Your citizen verification demographics.' }}
            </p>
          </div>
        </div>
        
        <div class="tw-flex tw-gap-2">
          <AppButton v-if="!isEditing" variant="primary" size="sm" @click="startEditing">
            ✏️ {{ t('editProfile') || 'Edit Profile' }}
          </AppButton>
          <template v-else>
            <AppButton variant="secondary" size="sm" @click="cancelEditing" :disabled="saving">
              {{ t('cancelEdit') || 'Cancel' }}
            </AppButton>
            <AppButton variant="primary" size="sm" @click="saveProfile" :disabled="saving">
              {{ saving ? t('submitting') : (t('saveProfile') || 'Save Profile') }}
            </AppButton>
          </template>
        </div>
      </div>
    </div>

    <!-- Profile Display / Edit Form -->
    <div v-if="authStore.userProfile">
      
      <!-- View Mode -->
      <div v-if="!isEditing" class="tw-grid tw-grid-cols-1 md:tw-grid-cols-2 tw-gap-6">
        
        <!-- Left: Personal info -->
        <AppCard class="tw-flex tw-flex-col tw-gap-4">
          <h3 class="tw-font-heading tw-font-bold tw-text-sm tw-text-primary tw-m-0 tw-pb-2 tw-border-b tw-border-border/50">
            {{ t('personalInfo') || 'Personal Information' }}
          </h3>
          
          <div class="tw-flex tw-flex-col tw-gap-3 tw-text-xs">
            <div class="tw-flex tw-justify-between tw-py-2 tw-border-b tw-border-border/30">
              <span class="tw-text-muted-foreground">{{ t('fullNameLabel') }}</span>
              <strong class="tw-text-foreground">{{ authStore.userProfile.full_name }}</strong>
            </div>
            <div class="tw-flex tw-justify-between tw-py-2 tw-border-b tw-border-border/30">
              <span class="tw-text-muted-foreground">{{ t('dobLabel') }}</span>
              <strong class="tw-text-foreground">{{ authStore.userProfile.date_of_birth ? new Date(authStore.userProfile.date_of_birth).toLocaleDateString() : '-' }}</strong>
            </div>
            <div class="tw-flex tw-justify-between tw-py-2 tw-border-b tw-border-border/30">
              <span class="tw-text-muted-foreground">{{ t('genderLabel') }}</span>
              <strong class="tw-text-foreground">{{ authStore.userProfile.gender }}</strong>
            </div>
            <div class="tw-flex tw-justify-between tw-py-2 tw-border-b tw-border-border/30">
              <span class="tw-text-muted-foreground">{{ t('stateLabel') }}</span>
              <strong class="tw-text-foreground">{{ authStore.userProfile.state }}</strong>
            </div>
            <div class="tw-flex tw-justify-between tw-py-2 tw-border-b tw-border-border/30">
              <span class="tw-text-muted-foreground">{{ t('districtLabel') }}</span>
              <strong class="tw-text-foreground">{{ authStore.userProfile.district }}</strong>
            </div>
            <div class="tw-flex tw-justify-between tw-py-2">
              <span class="tw-text-muted-foreground">{{ t('aadhaarLabel') }}</span>
              <strong class="tw-text-foreground">{{ authStore.userProfile.aadhaar || 'Not Provided' }}</strong>
            </div>
          </div>
        </AppCard>

        <!-- Right: Demographic details -->
        <AppCard class="tw-flex tw-flex-col tw-gap-4">
          <h3 class="tw-font-heading tw-font-bold tw-text-sm tw-text-primary tw-m-0 tw-pb-2 tw-border-b tw-border-border/50">
            {{ t('demographicDetails') || 'Demographics & Socioeconomics' }}
          </h3>
          
          <div class="tw-flex tw-flex-col tw-gap-3 tw-text-xs">
            <div class="tw-flex tw-justify-between tw-py-2 tw-border-b tw-border-border/30">
              <span class="tw-text-muted-foreground">{{ t('casteLabel') }}</span>
              <strong class="tw-text-foreground">{{ authStore.userProfile.caste_category }}</strong>
            </div>
            <div class="tw-flex tw-justify-between tw-py-2 tw-border-b tw-border-border/30">
              <span class="tw-text-muted-foreground">{{ t('incomeLabel') }}</span>
              <strong class="tw-text-foreground">₹{{ Number(authStore.userProfile.annual_income).toLocaleString() }}</strong>
            </div>
            <div class="tw-flex tw-justify-between tw-py-2 tw-border-b tw-border-border/30">
              <span class="tw-text-muted-foreground">{{ t('occupationLabel') }}</span>
              <strong class="tw-text-foreground">{{ authStore.userProfile.occupation }}</strong>
            </div>
            <div class="tw-flex tw-justify-between tw-py-2 tw-border-b tw-border-border/30">
              <span class="tw-text-muted-foreground">{{ t('employeeTypeLabel') }}</span>
              <strong class="tw-text-foreground">{{ authStore.userProfile.employee_type }}</strong>
            </div>
            <div class="tw-flex tw-justify-between tw-py-2 tw-border-b tw-border-border/30">
              <span class="tw-text-muted-foreground">{{ t('educationLabel') }}</span>
              <strong class="tw-text-foreground">{{ authStore.userProfile.education_level }}</strong>
            </div>
            <div class="tw-flex tw-justify-between tw-py-2">
              <span class="tw-text-muted-foreground">{{ t('disabilityLabel') || 'Differently-Abled' }}</span>
              <strong class="tw-text-foreground">{{ authStore.userProfile.is_disabled ? '✅ Yes' : '❌ No' }}</strong>
            </div>
          </div>
        </AppCard>

      </div>

      <!-- Edit Mode Form -->
      <form v-else @submit.prevent="saveProfile" class="tw-flex tw-flex-col tw-gap-6">
        
        <AppCard class="tw-flex tw-flex-col tw-gap-4">
          <h3 class="tw-font-heading tw-font-bold tw-text-sm tw-text-primary tw-m-0 tw-pb-2 tw-border-b tw-border-border/50">
            {{ t('personalInfo') }}
          </h3>

          <div class="tw-grid tw-grid-cols-1 sm:tw-grid-cols-2 tw-gap-4">
            <div>
              <AppLabel>{{ t('fullNameLabel') }} *</AppLabel>
              <AppInput v-model="editForm.full_name" type="text" required />
            </div>
            <div>
              <AppLabel>{{ t('dobLabel') }} *</AppLabel>
              <AppInput v-model="editForm.date_of_birth" type="date" required />
            </div>
            <div>
              <AppLabel>{{ t('genderLabel') }} *</AppLabel>
              <AppSelect v-model="editForm.gender" required>
                <option value="Male">{{ t('maleOpt') }}</option>
                <option value="Female">{{ t('femaleOpt') }}</option>
                <option value="Other">{{ t('otherOpt') }}</option>
              </AppSelect>
            </div>
            <div>
              <AppLabel>{{ t('casteLabel') }} *</AppLabel>
              <AppSelect v-model="editForm.caste_category" required>
                <option value="General">General / Open</option>
                <option value="OBC">OBC</option>
                <option value="SC">SC</option>
                <option value="ST">ST</option>
              </AppSelect>
            </div>
            <div>
              <AppLabel>{{ t('stateLabel') }} *</AppLabel>
              <AppSelect v-model="editForm.state" required>
                <option value="Maharashtra">Maharashtra</option>
                <option value="Gujarat">Gujarat</option>
                <option value="Madhya Pradesh">Madhya Pradesh</option>
                <option value="Karnataka">Karnataka</option>
                <option value="Delhi">Delhi</option>
                <option value="All">All India</option>
              </AppSelect>
            </div>
            <div>
              <AppLabel>{{ t('districtLabel') }} *</AppLabel>
              <AppInput v-model="editForm.district" type="text" required />
            </div>
          </div>

          <div class="tw-grid tw-grid-cols-1 tw-gap-4">
            <div>
              <AppLabel>{{ t('aadhaarLabel') }} *</AppLabel>
              <AppInput v-model="editForm.aadhaar" type="text" pattern="[0-9]{12}" title="Aadhaar number must be exactly 12 digits" required />
            </div>
          </div>
        </AppCard>

        <AppCard class="tw-flex tw-flex-col tw-gap-4">
          <h3 class="tw-font-heading tw-font-bold tw-text-sm tw-text-primary tw-m-0 tw-pb-2 tw-border-b tw-border-border/50">
            {{ t('demographicDetails') }}
          </h3>

          <div class="tw-grid tw-grid-cols-1 sm:tw-grid-cols-2 tw-gap-4">
            <div>
              <AppLabel>{{ t('incomeLabel') }} *</AppLabel>
              <AppInput v-model.number="editForm.annual_income" type="number" required />
            </div>
            <div>
              <AppLabel>{{ t('occupationLabel') }} *</AppLabel>
              <AppSelect v-model="editForm.occupation" required>
                <option value="Farmer">Farmer</option>
                <option value="Student">Student</option>
                <option value="Business Owner">Business Owner</option>
                <option value="Unemployed">Unemployed</option>
                <option value="Retired">Retired / Senior Citizen</option>
              </AppSelect>
            </div>
            <div>
              <AppLabel>{{ t('employeeTypeLabel') }} *</AppLabel>
              <AppSelect v-model="editForm.employee_type" required>
                <option value="Unemployed">Unemployed</option>
                <option value="Private Employee">Private Sector</option>
                <option value="Government Employee">Government Sector</option>
                <option value="Self Employed">Self Employed</option>
              </AppSelect>
            </div>
            <div>
              <AppLabel>{{ t('educationLabel') }} *</AppLabel>
              <AppSelect v-model="editForm.education_level" required>
                <option value="10th Pass">10th Standard or lower</option>
                <option value="12th Pass">12th Standard</option>
                <option value="Undergraduate">Undergraduate Degree</option>
                <option value="Graduate">Graduate / Master Degree</option>
                <option value="Post Graduate">Doctorate / Specialist</option>
              </AppSelect>
            </div>
          </div>

          <div class="tw-flex tw-items-center tw-gap-2.5 tw-mt-2">
            <input
              type="checkbox"
              id="editIsDisabled"
              v-model="editForm.is_disabled"
              class="tw-w-5 tw-h-5 tw-cursor-pointer"
            />
            <label for="editIsDisabled" class="tw-text-xs tw-font-bold tw-text-foreground tw-cursor-pointer">
              Yes, I am differently-abled / Haan, mai divyang hu
            </label>
          </div>
        </AppCard>

      </form>

    </div>

    <!-- Loading State -->
    <div v-else class="tw-text-center tw-py-12 tw-flex tw-flex-col tw-items-center tw-gap-3">
      <div class="tw-w-8 tw-h-8 tw-border-3 tw-border-muted tw-border-t-primary tw-rounded-full tw-animate-spin"></div>
      <p class="tw-text-xs tw-text-muted-foreground">{{ t('loading') || 'Loading Profile...' }}</p>
    </div>

  </div>
</template>
