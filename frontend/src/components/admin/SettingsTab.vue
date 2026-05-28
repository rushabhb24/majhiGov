<script setup>
import { API_BASE_URL } from '../../config.js'

defineProps({
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['refresh'])
</script>

<template>
  <div class="settings-tab">
    
    <div class="card" style="max-width: 600px;">
      <div class="card-header">
        <div class="card-title">Administrative Settings</div>
      </div>
      
      <div class="card-body">
        <div class="settings-wrapper">
          
          <div class="form-group">
            <label class="form-label">API Base Sync URL</label>
            <input 
              type="text" 
              class="form-input readonly-input" 
              :value="API_BASE_URL" 
              readonly 
            />
          </div>

          <div class="form-group">
            <label class="form-label">Database Connection Host</label>
            <input 
              type="text" 
              class="form-input readonly-input" 
              value="PostgreSQL (localhost:5432)" 
              readonly 
            />
          </div>

          <hr class="divider mt-3 mb-3" />

          <div class="actions-row">
            <button 
              type="button" 
              class="refresh-btn" 
              @click="emit('refresh')"
              :disabled="loading"
            >
              <i class="ti ti-refresh" :class="{ 'spin-anim': loading }"></i>
              <span>{{ loading ? 'Synchronizing...' : 'Force Database Sync' }}</span>
            </button>
          </div>

        </div>
      </div>
    </div>

  </div>
</template>

<style scoped>
.settings-tab {
  width: 100%;
}

.card {
  background-color: var(--bg);
  border: 0.5px solid var(--border);
  border-radius: 12px;
  overflow: hidden;
}

.card-header {
  padding: 14px 16px;
  border-bottom: 0.5px solid var(--border);
}

.card-title {
  font-size: 13px;
  font-weight: 500;
  color: var(--text);
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

.readonly-input {
  background-color: var(--bg2);
  color: var(--text2);
}

.divider {
  border: none;
  border-top: 0.5px solid var(--border);
}

.mt-3 { margin-top: 12px; }
.mb-3 { margin-bottom: 12px; }

.refresh-btn {
  background-color: var(--primary);
  color: var(--clr-text-light);
  border: none;
  border-radius: 6px;
  padding: 8px 16px;
  font-size: 13px;
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  font-family: inherit;
  box-sizing: border-box;
}

.refresh-btn:hover {
  opacity: 0.9;
}

.refresh-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.refresh-btn i {
  font-size: 16px !important;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.spin-anim {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
