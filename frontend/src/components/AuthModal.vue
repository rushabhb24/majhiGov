<script setup>
import { computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useUiStore } from '../stores/ui'
import { useI18n } from 'vue-i18n'
import AppDialog from './ui/AppDialog.vue'
import AppTabs from './ui/AppTabs.vue'
import AppInput from './ui/AppInput.vue'
import AppLabel from './ui/AppLabel.vue'
import AppSelect from './ui/AppSelect.vue'
import AppButton from './ui/AppButton.vue'

const authStore = useAuthStore()
const uiStore = useUiStore()
const { t } = useI18n()

const activeTab = computed({
  get: () => authStore.authTab,
  set: (val) => { authStore.authTab = val }
})

const tabs = computed(() => [
  { id: 'login', label: t('loginRegister').split(' / ')[0] },
  { id: 'register', label: t('loginRegister').split(' / ')[1] }
])
</script>

<template>
  <AppDialog
    :open="authStore.authModalOpen"
    @close="authStore.closeAuthModal"
    maxWidth="580px"
  >
    <!-- Tricolor Top Edge -->
    <div class="tricolor-bar tw-h-[4px] tw-absolute tw-top-0 tw-left-0 tw-w-full"></div>

    <div class="tw-mt-4">
      <AppTabs v-model="activeTab" :tabs="tabs">
        <template #default>
          
          <!-- Login Form -->
          <form v-if="activeTab === 'login'" @submit.prevent="authStore.loginUser()" class="tw-flex tw-flex-col tw-gap-4 tw-mt-4">
            <div>
              <AppLabel for="login-email">{{ t('emailLabel') }} *</AppLabel>
              <AppInput
                id="login-email"
                v-model="authStore.loginForm.email"
                type="email"
                placeholder="citizen@gov.in"
                required
              />
            </div>
            
            <div>
              <AppLabel for="login-password">{{ t('passwordLabel') }} *</AppLabel>
              <AppInput
                id="login-password"
                v-model="authStore.loginForm.password"
                type="password"
                placeholder="••••••••"
                required
              />
            </div>

            <AppButton
              type="submit"
              variant="primary"
              class="tw-mt-2 tw-w-full"
              :disabled="authStore.authSubmitting"
            >
              {{ authStore.authSubmitting ? t('submitting') : t('loginRegister').split(' / ')[0] }}
            </AppButton>
          </form>

          <!-- Register Form -->
          <form v-else @submit.prevent="authStore.registerUser()" class="tw-flex tw-flex-col tw-gap-4 tw-mt-4">
            <div class="tw-max-h-[60vh] tw-overflow-y-auto tw-pr-2 tw-flex tw-flex-col tw-gap-4">
              
              <div class="tw-grid tw-grid-cols-1 sm:tw-grid-cols-2 tw-gap-4">
                <div>
                  <AppLabel for="reg-name">{{ t('fullNameLabel') }} *</AppLabel>
                  <AppInput
                    id="reg-name"
                    v-model="authStore.regForm.full_name"
                    type="text"
                    placeholder="Ram Prasad"
                    required
                  />
                </div>
                <div>
                  <AppLabel for="reg-dob">{{ t('dobLabel') }} *</AppLabel>
                  <AppInput
                    id="reg-dob"
                    v-model="authStore.regForm.date_of_birth"
                    type="date"
                    required
                  />
                </div>
              </div>

              <div class="tw-grid tw-grid-cols-1 sm:tw-grid-cols-2 tw-gap-4">
                <div>
                  <AppLabel for="reg-email">{{ t('emailLabel') }} *</AppLabel>
                  <AppInput
                    id="reg-email"
                    v-model="authStore.regForm.email"
                    type="email"
                    placeholder="ram@gov.in"
                    required
                  />
                </div>
                <div>
                  <AppLabel for="reg-phone">{{ t('phoneLabel') }} *</AppLabel>
                  <AppInput
                    id="reg-phone"
                    v-model="authStore.regForm.phone"
                    type="tel"
                    placeholder="9876543210"
                    required
                  />
                </div>
              </div>

              <div>
                <AppLabel for="reg-password">{{ t('passwordLabel') }} *</AppLabel>
                <AppInput
                  id="reg-password"
                  v-model="authStore.regForm.password"
                  type="password"
                  placeholder="Create secure password"
                  required
                />
              </div>

              <div>
                <AppLabel for="reg-aadhaar">{{ t('aadhaarLabel') }} *</AppLabel>
                <AppInput
                  id="reg-aadhaar"
                  v-model="authStore.regForm.aadhaar"
                  type="text"
                  placeholder="e.g. 555566667777"
                  pattern="[0-9]{12}"
                  title="Aadhaar number must be exactly 12 digits"
                  required
                />
              </div>

              <hr class="tw-border-border/50 tw-my-2" />
              <h4 class="tw-font-heading tw-font-bold tw-text-xs tw-text-primary tw-m-0">
                {{ t('demographicDetails') }}
              </h4>

              <div class="tw-grid tw-grid-cols-2 tw-gap-4">
                <div>
                  <AppLabel for="reg-gender">{{ t('genderLabel') }} *</AppLabel>
                  <AppSelect
                    id="reg-gender"
                    v-model="authStore.regForm.gender"
                    required
                  >
                    <option value="Male">{{ t('maleOpt') }}</option>
                    <option value="Female">{{ t('femaleOpt') }}</option>
                    <option value="Other">{{ t('otherOpt') }}</option>
                  </AppSelect>
                </div>
                <div>
                  <AppLabel for="reg-caste">{{ t('casteLabel') }} *</AppLabel>
                  <AppSelect
                    id="reg-caste"
                    v-model="authStore.regForm.caste_category"
                    required
                  >
                    <option value="General">General / Open</option>
                    <option value="OBC">OBC</option>
                    <option value="SC">SC</option>
                    <option value="ST">ST</option>
                  </AppSelect>
                </div>
              </div>

              <div class="tw-grid tw-grid-cols-2 tw-gap-4">
                <div>
                  <AppLabel for="reg-state">{{ t('stateLabel') }} *</AppLabel>
                  <AppSelect
                    id="reg-state"
                    v-model="authStore.regForm.state"
                    required
                  >
                    <option value="Maharashtra">Maharashtra</option>
                    <option value="Gujarat">Gujarat</option>
                    <option value="Madhya Pradesh">Madhya Pradesh</option>
                    <option value="Karnataka">Karnataka</option>
                    <option value="Delhi">Delhi</option>
                    <option value="All">All India</option>
                  </AppSelect>
                </div>
                <div>
                  <AppLabel for="reg-district">{{ t('districtLabel') }} *</AppLabel>
                  <AppInput
                    id="reg-district"
                    v-model="authStore.regForm.district"
                    type="text"
                    placeholder="Pune"
                    required
                  />
                </div>
              </div>

              <div class="tw-grid tw-grid-cols-2 tw-gap-4">
                <div>
                  <AppLabel for="reg-income">{{ t('incomeLabel') }} *</AppLabel>
                  <AppInput
                    id="reg-income"
                    v-model="authStore.regForm.annual_income"
                    type="number"
                    required
                  />
                </div>
                <div>
                  <AppLabel for="reg-occupation">{{ t('occupationLabel') }} *</AppLabel>
                  <AppSelect
                    id="reg-occupation"
                    v-model="authStore.regForm.occupation"
                    required
                  >
                    <option value="Farmer">Farmer</option>
                    <option value="Student">Student</option>
                    <option value="Business Owner">Business Owner</option>
                    <option value="Unemployed">Unemployed</option>
                    <option value="Retired">Retired / Senior Citizen</option>
                  </AppSelect>
                </div>
              </div>

              <div class="tw-grid tw-grid-cols-2 tw-gap-4">
                <div>
                  <AppLabel for="reg-emptype">{{ t('employeeTypeLabel') }} *</AppLabel>
                  <AppSelect
                    id="reg-emptype"
                    v-model="authStore.regForm.employee_type"
                    required
                  >
                    <option value="Unemployed">Unemployed</option>
                    <option value="Private Employee">Private Sector</option>
                    <option value="Government Employee">Government Sector</option>
                    <option value="Self Employed">Self Employed</option>
                  </AppSelect>
                </div>
                <div>
                  <AppLabel for="reg-education">{{ t('educationLabel') }} *</AppLabel>
                  <AppSelect
                    id="reg-education"
                    v-model="authStore.regForm.education_level"
                    required
                  >
                    <option value="10th Pass">10th Standard or lower</option>
                    <option value="12th Pass">12th Standard</option>
                    <option value="Undergraduate">Undergraduate Degree</option>
                    <option value="Graduate">Graduate / Master Degree</option>
                    <option value="Post Graduate">Doctorate / Specialist</option>
                  </AppSelect>
                </div>
              </div>

              <div class="tw-flex tw-items-center tw-gap-2 tw-mt-2">
                <input 
                  type="checkbox" 
                  id="regIsDisabled" 
                  v-model="authStore.regForm.is_disabled" 
                  class="tw-w-5 tw-h-5 tw-cursor-pointer"
                />
                <label for="regIsDisabled" class="tw-cursor-pointer tw-text-sm tw-text-foreground">
                  Yes, I am differently-abled / Haan, mai divyang hu
                </label>
              </div>

            </div>

            <AppButton
              type="submit"
              variant="primary"
              class="tw-mt-2 tw-w-full"
              :disabled="authStore.authSubmitting"
            >
              {{ authStore.authSubmitting ? t('submitting') : t('loginRegister').split(' / ')[1] }}
            </AppButton>
          </form>

        </template>
      </AppTabs>
    </div>
  </AppDialog>
</template>
