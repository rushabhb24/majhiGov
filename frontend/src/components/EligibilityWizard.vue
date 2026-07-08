<script setup>
import AppCard from './ui/AppCard.vue'
import AppInput from './ui/AppInput.vue'
import AppSelect from './ui/AppSelect.vue'
import AppButton from './ui/AppButton.vue'
import AppLabel from './ui/AppLabel.vue'

defineProps({
  profile: {
    type: Object,
    required: true
  },
  step: {
    type: Number,
    required: true
  },
  t: {
    type: Object,
    required: true
  },
  checking: {
    type: Boolean,
    required: true
  }
});

const emit = defineEmits(['update:step', 'submit']);
</script>

<template>
  <AppCard class="tw-max-w-2xl tw-mx-auto tw-p-8">
    <h2 class="tw-font-heading tw-font-bold tw-text-xl tw-text-foreground tw-mb-6 tw-m-0 text-center">
      {{ t.eligibility || 'Smart Eligibility Wizard' }}
    </h2>
    
    <!-- Step Indicators Stepper -->
    <div class="tw-flex tw-items-center tw-justify-between tw-mb-8 tw-relative tw-max-w-md tw-mx-auto">
      
      <!-- Line behind -->
      <div class="tw-absolute tw-h-[2px] tw-bg-border tw-w-full tw-top-1/2 -tw-translate-y-1/2 tw-z-0"></div>
      
      <!-- Step 1 node -->
      <div 
        class="tw-w-8 tw-h-8 tw-rounded-full tw-flex tw-items-center tw-justify-center tw-font-heading tw-font-bold tw-text-xs tw-z-10 tw-transition-colors"
        :class="step >= 1 ? 'tw-bg-primary tw-text-white' : 'tw-bg-muted tw-text-muted-foreground'"
      >
        <span v-if="step > 1">✓</span>
        <span v-else>1</span>
      </div>

      <!-- Step 2 node -->
      <div 
        class="tw-w-8 tw-h-8 tw-rounded-full tw-flex tw-items-center tw-justify-center tw-font-heading tw-font-bold tw-text-xs tw-z-10 tw-transition-colors"
        :class="step >= 2 ? 'tw-bg-primary tw-text-white' : 'tw-bg-muted tw-text-muted-foreground'"
      >
        <span v-if="step > 2">✓</span>
        <span v-else>2</span>
      </div>

      <!-- Step 3 node -->
      <div 
        class="tw-w-8 tw-h-8 tw-rounded-full tw-flex tw-items-center tw-justify-center tw-font-heading tw-font-bold tw-text-xs tw-z-10 tw-transition-colors"
        :class="step >= 3 ? 'tw-bg-primary tw-text-white' : 'tw-bg-muted tw-text-muted-foreground'"
      >
        <span>3</span>
      </div>

    </div>

    <!-- STEP 1: Personal Details -->
    <div v-if="step === 1" class="tw-flex tw-flex-col tw-gap-4">
      <h3 class="tw-font-heading tw-font-bold tw-text-sm tw-text-primary tw-m-0">
        {{ t.personalProfile || 'Personal Profile' }}
      </h3>
      
      <div class="tw-grid tw-grid-cols-1 sm:tw-grid-cols-2 tw-gap-4">
        <div>
          <AppLabel for="age">{{ t.ageLabel || 'Age' }}</AppLabel>
          <AppInput
            id="age"
            type="number"
            v-model.number="profile.age"
            min="1"
            max="120"
          />
        </div>
        
        <div>
          <AppLabel for="gender">{{ t.genderLabel || 'Gender' }}</AppLabel>
          <AppSelect id="gender" v-model="profile.gender">
            <option value="Male">{{ t.maleOpt || 'Male' }}</option>
            <option value="Female">{{ t.femaleOpt || 'Female' }}</option>
            <option value="Other">{{ t.otherOpt || 'Other' }}</option>
          </AppSelect>
        </div>

        <div>
          <AppLabel for="state">{{ t.stateLabel || 'State' }}</AppLabel>
          <AppSelect id="state" v-model="profile.state">
            <option value="Maharashtra">Maharashtra</option>
            <option value="Madhya Pradesh">Madhya Pradesh</option>
            <option value="Gujarat">Gujarat</option>
            <option value="Karnataka">Karnataka</option>
            <option value="Uttar Pradesh">Uttar Pradesh</option>
            <option value="Delhi">Delhi</option>
          </AppSelect>
        </div>

        <div>
          <AppLabel for="caste">{{ t.casteLabel || 'Caste Category' }}</AppLabel>
          <AppSelect id="caste" v-model="profile.caste">
            <option value="General">General / Open</option>
            <option value="OBC">OBC</option>
            <option value="SC">SC</option>
            <option value="ST">ST</option>
          </AppSelect>
        </div>
      </div>

      <div class="tw-flex tw-justify-between tw-mt-4">
        <div></div>
        <AppButton 
          variant="primary" 
          size="sm" 
          class="tw-flex tw-items-center tw-gap-1.5"
          @click="emit('update:step', 2)"
        >
          <span>{{ t.next || 'Next' }}</span>
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="5" y1="12" x2="19" y2="12"></line><polyline points="12 5 19 12 12 19"></polyline></svg>
        </AppButton>
      </div>
    </div>

    <!-- STEP 2: Income and Occupation -->
    <div v-else-if="step === 2" class="tw-flex tw-flex-col tw-gap-4">
      <h3 class="tw-font-heading tw-font-bold tw-text-sm tw-text-primary tw-m-0">
        {{ t.incomeOccupation || 'Income & Occupation' }}
      </h3>

      <div class="tw-grid tw-grid-cols-1 sm:tw-grid-cols-3 tw-gap-4">
        <div>
          <AppLabel for="income">{{ t.incomeLabel || 'Annual Income (₹)' }}</AppLabel>
          <AppInput
            id="income"
            type="number"
            v-model.number="profile.annual_income"
          />
        </div>

        <div>
          <AppLabel for="occupation">{{ t.occupationLabel || 'Occupation' }}</AppLabel>
          <AppSelect id="occupation" v-model="profile.occupation">
            <option value="Farmer">Farmer (Kisan)</option>
            <option value="Student">Student (Vidyarthi)</option>
            <option value="Business Owner">Business Owner</option>
            <option value="Unemployed">Unemployed</option>
            <option value="Self-Employed">Self-Employed</option>
            <option value="Other">Other</option>
          </AppSelect>
        </div>

        <div>
          <AppLabel for="emp-type">{{ t.employeeTypeLabel || 'Employment Type' }}</AppLabel>
          <AppSelect id="emp-type" v-model="profile.employee_type">
            <option value="Unemployed">None / Unemployed</option>
            <option value="Private">Private Employee</option>
            <option value="Government">Government Employee</option>
            <option value="Self-Employed">Self-Employed</option>
          </AppSelect>
        </div>
      </div>

      <div class="tw-flex tw-justify-between tw-mt-4">
        <AppButton variant="outline" size="sm" @click="emit('update:step', 1)">
          {{ t.back || 'Back' }}
        </AppButton>
        <AppButton 
          variant="primary" 
          size="sm" 
          class="tw-flex tw-items-center tw-gap-1.5"
          @click="emit('update:step', 3)"
        >
          <span>{{ t.next || 'Next' }}</span>
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="5" y1="12" x2="19" y2="12"></line><polyline points="12 5 19 12 12 19"></polyline></svg>
        </AppButton>
      </div>
    </div>

    <!-- STEP 3: Education & Disability -->
    <div v-else-if="step === 3" class="tw-flex tw-flex-col tw-gap-4">
      <h3 class="tw-font-heading tw-font-bold tw-text-sm tw-text-primary tw-m-0">
        {{ t.educationSpecial || 'Education & Special Criteria' }}
      </h3>

      <div class="tw-grid tw-grid-cols-1 sm:tw-grid-cols-2 tw-gap-4">
        <div>
          <AppLabel for="education">{{ t.educationLabel || 'Education Level' }}</AppLabel>
          <AppSelect id="education" v-model="profile.education_level">
            <option value="None">None (Uneducated)</option>
            <option value="Primary">Primary Schooling</option>
            <option value="10th Pass">10th Class Pass</option>
            <option value="12th Pass">12th Class Pass</option>
            <option value="Graduate">College Graduate (Degree)</option>
            <option value="Post Graduate">Post Graduate</option>
          </AppSelect>
        </div>

        <div class="tw-flex tw-items-center tw-gap-2.5 tw-mt-6">
          <input 
            type="checkbox" 
            id="disability" 
            v-model="profile.is_disabled" 
            class="tw-w-5 tw-h-5 tw-cursor-pointer"
          />
          <label for="disability" class="tw-text-xs tw-font-bold tw-text-foreground tw-cursor-pointer">
            {{ t.disabilityLabel || 'Differently-Abled (Divyangjan)' }}
          </label>
        </div>
      </div>

      <div class="tw-flex tw-justify-between tw-mt-4">
        <AppButton variant="outline" size="sm" @click="emit('update:step', 2)">
          {{ t.back || 'Back' }}
        </AppButton>
        <AppButton 
          variant="primary" 
          size="sm" 
          @click="emit('submit')" 
          :disabled="checking"
        >
          {{ checking ? 'Calculating...' : (t.calculate || 'Check Eligibility') }}
        </AppButton>
      </div>
    </div>
  </AppCard>
</template>
