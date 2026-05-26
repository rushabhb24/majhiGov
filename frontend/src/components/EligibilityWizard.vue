<script setup>
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
  <div class="card wizard-panel">
    <h2 class="section-title">{{ t.eligibility }}</h2>
    
    <!-- Step Indicators -->
    <div class="steps-indicator">
      <div :class="['step-node', { active: step >= 1 }]">1</div>
      <div class="step-line"></div>
      <div :class="['step-node', { active: step >= 2 }]">2</div>
      <div class="step-line"></div>
      <div :class="['step-node', { active: step >= 3 }]">3</div>
    </div>

    <!-- STEP 1: Personal Details -->
    <div v-if="step === 1" class="step-content animate-fade">
      <h3 class="step-title">{{ t.personalProfile }}</h3>
      
      <div class="form-grid">
        <div class="form-group">
          <label class="form-label" for="age">{{ t.ageLabel }}</label>
          <input type="number" id="age" class="form-control" v-model.number="profile.age" min="1" max="120" />
        </div>
        
        <div class="form-group">
          <label class="form-label" for="gender">{{ t.genderLabel }}</label>
          <select id="gender" class="form-control" v-model="profile.gender">
            <option value="Male">{{ t.maleOpt }}</option>
            <option value="Female">{{ t.femaleOpt }}</option>
            <option value="Other">{{ t.otherOpt }}</option>
          </select>
        </div>

        <div class="form-group">
          <label class="form-label" for="state">{{ t.stateLabel }}</label>
          <select id="state" class="form-control" v-model="profile.state">
            <option value="Maharashtra">Maharashtra</option>
            <option value="Madhya Pradesh">Madhya Pradesh</option>
            <option value="Gujarat">Gujarat</option>
            <option value="Karnataka">Karnataka</option>
            <option value="Uttar Pradesh">Uttar Pradesh</option>
            <option value="Delhi">Delhi</option>
          </select>
        </div>

        <div class="form-group">
          <label class="form-label" for="caste">{{ t.casteLabel }}</label>
          <select id="caste" class="form-control" v-model="profile.caste">
            <option value="General">General / Open</option>
            <option value="OBC">OBC</option>
            <option value="SC">SC</option>
            <option value="ST">ST</option>
          </select>
        </div>
      </div>

      <div class="wizard-actions">
        <div></div>
        <button class="btn btn-primary" @click="emit('update:step', 2)">
          {{ t.next }}
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="5" y1="12" x2="19" y2="12"></line><polyline points="12 5 19 12 12 19"></polyline></svg>
        </button>
      </div>
    </div>

    <!-- STEP 2: Income and Occupation -->
    <div v-else-if="step === 2" class="step-content animate-fade">
      <h3 class="step-title">{{ t.incomeOccupation }}</h3>

      <div class="form-grid">
        <div class="form-group">
          <label class="form-label" for="income">{{ t.incomeLabel }}</label>
          <input type="number" id="income" class="form-control" v-model.number="profile.annual_income" />
        </div>

        <div class="form-group">
          <label class="form-label" for="occupation">{{ t.occupationLabel }}</label>
          <select id="occupation" class="form-control" v-model="profile.occupation">
            <option value="Farmer">Farmer (Kisan)</option>
            <option value="Student">Student (Vidyarthi)</option>
            <option value="Business Owner">Business Owner</option>
            <option value="Unemployed">Unemployed</option>
            <option value="Self-Employed">Self-Employed</option>
            <option value="Other">Other</option>
          </select>
        </div>

        <div class="form-group">
          <label class="form-label" for="emp-type">{{ t.employeeTypeLabel }}</label>
          <select id="emp-type" class="form-control" v-model="profile.employee_type">
            <option value="Unemployed">None / Unemployed</option>
            <option value="Private">Private Employee</option>
            <option value="Government">Government Employee</option>
            <option value="Self-Employed">Self-Employed</option>
          </select>
        </div>
      </div>

      <div class="wizard-actions">
        <button class="btn btn-secondary" @click="emit('update:step', 1)">
          {{ t.back }}
        </button>
        <button class="btn btn-primary" @click="emit('update:step', 3)">
          {{ t.next }}
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="5" y1="12" x2="19" y2="12"></line><polyline points="12 5 19 12 12 19"></polyline></svg>
        </button>
      </div>
    </div>

    <!-- STEP 3: Education & Disability -->
    <div v-else-if="step === 3" class="step-content animate-fade">
      <h3 class="step-title">{{ t.educationSpecial }}</h3>

      <div class="form-grid">
        <div class="form-group">
          <label class="form-label" for="education">{{ t.educationLabel }}</label>
          <select id="education" class="form-control" v-model="profile.education_level">
            <option value="None">None (Uneducated)</option>
            <option value="Primary">Primary Schooling</option>
            <option value="10th Pass">10th Class Pass</option>
            <option value="12th Pass">12th Class Pass</option>
            <option value="Graduate">College Graduate (Degree)</option>
            <option value="Post Graduate">Post Graduate</option>
          </select>
        </div>

        <!-- Custom vertical flex aligned checkbox block -->
        <div class="checkbox-align-wrapper">
          <input type="checkbox" id="disability" class="checkbox-control" v-model="profile.is_disabled" />
          <label class="checkbox-label mb-0" for="disability">
            {{ t.disabilityLabel }}
          </label>
        </div>
      </div>

      <div class="wizard-actions">
        <button class="btn btn-secondary" @click="emit('update:step', 2)">
          {{ t.back }}
        </button>
        <button class="btn btn-accent" @click="emit('submit')" :disabled="checking">
          {{ checking ? 'Calculating...' : t.calculate }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.animate-fade {
  animation: fadeIn 0.3s ease-out;
}
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(4px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Beautiful custom checkbox alignment rules */
.checkbox-align-wrapper {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  min-height: 48px;
  height: auto;
  align-self: end;
  margin-bottom: 20px;
  min-width: 0;
  padding-top: 14px;
}

.checkbox-control {
  width: 20px;
  height: 20px;
  cursor: pointer;
  accent-color: var(--clr-primary);
  flex-shrink: 0;
}

.checkbox-label {
  cursor: pointer;
  user-select: none;
  min-width: 0;
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--clr-text-main);
  line-height: 1.3;
  overflow-wrap: anywhere;
}
.mb-0 { margin-bottom: 0; }
</style>
