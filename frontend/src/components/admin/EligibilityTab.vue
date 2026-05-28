<script setup>
import { watch } from 'vue'

const props = defineProps({
  schemes: {
    type: Array,
    required: true
  },
  selectedSchemeId: {
    type: [Number, String],
    default: null
  },
  ruleForm: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['update:selectedSchemeId', 'save-rules'])

function handleSchemeChange(e) {
  emit('update:selectedSchemeId', Number(e.target.value))
}
</script>

<template>
  <div class="eligibility-tab">
    
    <!-- Single full-width card -->
    <div class="card">
      <div class="card-header">
        <div class="card-title">Set Eligibility Rules for a Scheme</div>
      </div>
      
      <div class="card-body">
        
        <!-- Select Scheme dropdown -->
        <div class="form-group">
          <label class="form-label font-medium">Select Scheme</label>
          <select 
            class="form-input font-medium" 
            :value="selectedSchemeId" 
            @change="handleSchemeChange"
          >
            <option v-for="s in schemes" :key="s.id" :value="s.id">{{ s.title }}</option>
          </select>
        </div>

        <hr class="divider mt-4 mb-4" />

        <form @submit.prevent="emit('save-rules')" v-if="selectedSchemeId">
          
          <!-- Two-column row: Min Age / Max Age -->
          <div class="form-row">
            <div class="form-group">
              <label class="form-label">Min Age</label>
              <input 
                type="number" 
                class="form-input" 
                v-model="ruleForm.min_age" 
                placeholder="18" 
                required 
              />
            </div>
            <div class="form-group">
              <label class="form-label">Max Age (0 = no limit)</label>
              <input 
                type="number" 
                class="form-input" 
                v-model="ruleForm.max_age" 
                placeholder="60 (0 = no limit)" 
                required 
              />
            </div>
          </div>

          <!-- Two-column row: Min Annual Income / Max Annual Income -->
          <div class="form-row">
            <div class="form-group">
              <label class="form-label">Min Annual Income (₹)</label>
              <input 
                type="number" 
                class="form-input" 
                v-model="ruleForm.min_income" 
                placeholder="0" 
                required 
              />
            </div>
            <div class="form-group">
              <label class="form-label">Max Annual Income (₹)</label>
              <input 
                type="number" 
                class="form-input" 
                v-model="ruleForm.max_income" 
                placeholder="200000" 
                required 
              />
            </div>
          </div>

          <!-- Two-column row: Gender / Caste -->
          <div class="form-row">
            <div class="form-group">
              <label class="form-label">Gender</label>
              <select class="form-input" v-model="ruleForm.gender" required>
                <option value="all">All</option>
                <option value="male">Male</option>
                <option value="female">Female</option>
                <option value="other">Other</option>
              </select>
            </div>
            <div class="form-group">
              <label class="form-label">Caste Categories (SC / ST / OBC / General)</label>
              <select class="form-input" v-model="ruleForm.caste_categories" multiple style="height: 60px; padding: 4px;">
                <option value="General">General</option>
                <option value="OBC">OBC</option>
                <option value="SC">SC</option>
                <option value="ST">ST</option>
              </select>
            </div>
          </div>

          <!-- Applicable States -->
          <div class="form-group">
            <label class="form-label">Applicable States (leave empty = All India)</label>
            <input 
              type="text" 
              class="form-input" 
              v-model="ruleForm.states_str" 
              placeholder="e.g. Maharashtra, Rajasthan, UP" 
            />
          </div>

          <!-- Occupations -->
          <div class="form-group">
            <label class="form-label">Occupations (leave empty = All)</label>
            <input 
              type="text" 
              class="form-input" 
              v-model="ruleForm.occupations_str" 
              placeholder="e.g. farmer, student, unemployed" 
            />
          </div>

          <!-- Disability Toggle -->
          <div class="form-group mt-2">
            <div style="display: flex; align-items: center; gap: 8px;">
              <input 
                type="checkbox" 
                id="disabilityCheck" 
                v-model="ruleForm.disability_required" 
                style="width: 18px; height: 18px; cursor: pointer;" 
              />
              <label for="disabilityCheck" style="cursor: pointer; font-weight: 500;">
                Differently-abled criteria required
              </label>
            </div>
          </div>

          <!-- Submit Button -->
          <button type="submit" class="submit-btn mt-3">
            <i class="ti ti-device-floppy"></i>
            <span>Save Eligibility Rules</span>
          </button>

        </form>

        <div v-else class="empty-state">
          No active schemes available to configure rules.
        </div>

      </div>
    </div>

  </div>
</template>

<style scoped>
.eligibility-tab {
  width: 100%;
}

.card {
  background-color: #ffffff;
  border: 0.5px solid rgba(0, 0, 0, 0.08);
  border-radius: 12px;
  overflow: hidden;
}

.card-header {
  padding: 14px 16px;
  border-bottom: 0.5px solid rgba(0, 0, 0, 0.08);
}

.card-title {
  font-size: 13px;
  font-weight: 500;
  color: #0f172a;
}

.card-body {
  padding: 16px;
}

/* Forms */
.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 12px;
}

.form-label {
  font-size: 13px;
  color: #0f172a;
}

.font-medium {
  font-weight: 500;
}

.form-input {
  padding: 8px 10px;
  border: 0.5px solid rgba(0, 0, 0, 0.08);
  border-radius: 6px;
  font-size: 13px;
  background-color: #ffffff;
  color: #0f172a;
  outline: none;
  font-family: inherit;
  box-sizing: border-box;
  width: 100%;
}

.form-input:focus {
  border-color: #1a3a6b;
}

.form-row {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.divider {
  border: none;
  border-top: 0.5px solid rgba(0, 0, 0, 0.08);
}

.mt-2 { margin-top: 8px; }
.mt-3 { margin-top: 12px; }
.mt-4 { margin-top: 16px; }
.mb-4 { margin-bottom: 16px; }

.submit-btn {
  background-color: #1a3a6b; /* var(--primary) */
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
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.empty-state {
  text-align: center;
  padding: 30px;
  color: #64748b;
}
</style>
