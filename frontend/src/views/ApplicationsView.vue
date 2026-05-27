<script setup>
import { onMounted } from 'vue'
import { useApplicationStore } from '../stores/applications'
import { useUiStore } from '../stores/ui'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const applicationStore = useApplicationStore()
const uiStore = useUiStore()
const { t } = useI18n()
const router = useRouter()

function getStatusIcon(status) {
  if (status === 'approved') return '✅'
  if (status === 'rejected') return '❌'
  return '⏳'
}

function getStepClass(appStatus, stepStatus) {
  const step = applicationStore.getStatusStep(appStatus)
  const stepNum = stepStatus === 'submitted' ? 1 : stepStatus === 'review' ? 2 : 3
  
  if (stepNum < step) return 'step-completed'
  if (stepNum === step) {
    if (appStatus === 'approved') return 'step-completed step-success'
    if (appStatus === 'rejected') return 'step-completed step-danger'
    return 'step-active'
  }
  return 'step-pending'
}

function getLineClass(appStatus, afterStep) {
  const step = applicationStore.getStatusStep(appStatus)
  return afterStep < step ? 'line-filled' : 'line-empty'
}

onMounted(() => {
  applicationStore.fetchApplications()
})
</script>

<template>
  <div class="tab-content animate-fade">
    <div class="card filter-panel">
      <div class="app-header-row">
        <div>
          <h2 class="section-title">{{ t('myApplications') }}</h2>
          <p class="text-muted">{{ t('applicationTimeline') }}</p>
        </div>
        <button 
          class="btn btn-secondary" 
          @click="applicationStore.refreshApplications()"
          :disabled="applicationStore.refreshing"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" :class="{ 'spin-anim': applicationStore.refreshing }"><polyline points="23 4 23 10 17 10"></polyline><polyline points="1 20 1 14 7 14"></polyline><path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"></path></svg>
          {{ applicationStore.refreshing ? t('submitting') : t('refreshStatus') || 'Refresh Status' }}
        </button>
      </div>
    </div>

    <div v-if="applicationStore.applications.length === 0" class="empty-state text-center mt-4 card">
      <div class="empty-bookmarks-art">📬</div>
      <h3>{{ t('noApplications') }}</h3>
      <p class="text-muted mt-2">{{ t('noApplicationsHint') || 'Find a scheme and apply through the official portal' }}</p>
      <button class="btn btn-primary mt-4" @click="router.push('/')">{{ t('exploreSchemes') || 'Explore Schemes' }}</button>
    </div>

    <!-- Application Tracking Cards -->
    <div v-else class="apps-list">
      <div 
        v-for="app in applicationStore.applications" 
        :key="app.id" 
        class="app-tracker-card card"
      >
        <!-- Card Header -->
        <div class="app-card-header">
          <div class="app-card-info">
            <span class="app-govt-badge" :class="app.government_level">
              {{ app.government_level === 'central' ? '🏛️ Central' : '🏢 State' }}
            </span>
            <h3 class="app-card-title">
              {{ uiStore.currentLanguage === 'mr' ? (app.title_mr || app.title) : (uiStore.currentLanguage === 'hi' ? (app.title_hi || app.title) : app.title) }}
            </h3>
            <div class="app-card-date">
              📅 {{ t('appliedOn') || 'Applied on' }}: {{ new Date(app.applied_at).toLocaleDateString(uiStore.currentLanguage === 'en' ? 'en-US' : (uiStore.currentLanguage === 'hi' ? 'hi-IN' : 'mr-IN'), { year: 'numeric', month: 'long', day: 'numeric' }) }}
            </div>
          </div>
          <div class="app-card-status">
            <span class="status-icon">{{ getStatusIcon(app.status) }}</span>
            <span 
              class="badge" 
              :class="{
                'badge-warning': app.status === 'pending',
                'badge-success': app.status === 'approved',
                'badge-danger': app.status === 'rejected'
              }"
            >
              {{ app.status === 'pending' ? t('statusPending') : (app.status === 'approved' ? t('statusApproved') : t('statusRejected')) }}
            </span>
          </div>
        </div>

        <!-- Status Progress Steps -->
        <div class="status-progress">
          <div class="progress-step" :class="getStepClass(app.status, 'submitted')">
            <div class="step-dot"></div>
            <span class="step-label">{{ t('stepSubmitted') || 'Submitted' }}</span>
          </div>
          <div class="progress-line" :class="getLineClass(app.status, 1)"></div>
          <div class="progress-step" :class="getStepClass(app.status, 'review')">
            <div class="step-dot"></div>
            <span class="step-label">{{ t('stepReview') || 'Under Review' }}</span>
          </div>
          <div class="progress-line" :class="getLineClass(app.status, 2)"></div>
          <div class="progress-step" :class="getStepClass(app.status, 'decision')">
            <div class="step-dot"></div>
            <span class="step-label">{{ app.status === 'rejected' ? (t('stepRejected') || 'Rejected') : (t('stepApproved') || 'Approved') }}</span>
          </div>
        </div>

        <!-- Citizen notes -->
        <div v-if="app.notes" class="app-card-notes">
          <strong>{{ t('notesTitle') }}</strong> {{ app.notes }}
        </div>

        <!-- Action: Re-visit official portal -->
        <div v-if="app.apply_link || app.official_website" class="app-card-actions">
          <a 
            :href="app.apply_link || app.official_website" 
            target="_blank" 
            rel="noopener noreferrer" 
            class="btn btn-secondary"
          >
            🔗 {{ t('visitOfficialPortal') || 'Visit Official Portal' }}
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path><polyline points="15 3 21 3 21 9"></polyline><line x1="10" y1="14" x2="21" y2="3"></line></svg>
          </a>
        </div>
      </div>
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
.empty-bookmarks-art {
  font-size: 3.5rem;
  margin-bottom: 12px;
}

.app-header-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  flex-wrap: wrap;
}

.apps-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
  margin-top: 24px;
}

.app-tracker-card {
  padding: 24px;
}

.app-tracker-card:hover {
  transform: translateY(-3px);
}

.app-card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  flex-wrap: wrap;
}

.app-card-info {
  flex: 1 1 300px;
  min-width: 0;
}

.app-govt-badge {
  display: inline-block;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 3px 10px;
  border-radius: var(--border-radius-full);
  background: var(--clr-secondary-light);
  color: var(--clr-secondary);
  margin-bottom: 8px;
}

.app-govt-badge.state {
  background: var(--clr-accent-light);
  color: hsl(35, 92%, 40%);
}

.app-card-title {
  font-family: var(--font-heading);
  font-weight: 700;
  font-size: 1.15rem;
  color: var(--clr-text-main);
  margin-bottom: 6px;
  line-height: 1.4;
}

.app-card-date {
  font-size: 0.82rem;
  color: var(--clr-text-muted);
}

.app-card-status {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.status-icon {
  font-size: 1.5rem;
}

.app-card-status .badge {
  padding: 6px 14px;
  font-size: 0.78rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.3px;
}

/* --- Status Progress Steps --- */
.status-progress {
  display: flex;
  align-items: center;
  gap: 0;
  margin: 24px 0 16px;
  padding: 0 8px;
}

.progress-step {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.step-dot {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: 3px solid var(--clr-border);
  background: var(--clr-surface-alt);
  transition: all 0.3s ease;
  position: relative;
}

.step-label {
  font-size: 0.72rem;
  font-weight: 600;
  color: var(--clr-text-muted);
  text-align: center;
  white-space: nowrap;
}

/* Step states */
.step-completed .step-dot {
  border-color: var(--clr-secondary);
  background: var(--clr-secondary);
  box-shadow: 0 0 0 4px var(--clr-secondary-light);
}

.step-completed .step-dot::after {
  content: '✓';
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 0.7rem;
  font-weight: 800;
}

.step-completed .step-label {
  color: var(--clr-secondary);
  font-weight: 700;
}

.step-active .step-dot {
  border-color: var(--clr-primary);
  background: var(--clr-primary-light);
  box-shadow: 0 0 0 4px var(--clr-primary-glow);
  animation: pulse-dot 2s infinite;
}

.step-active .step-label {
  color: var(--clr-primary);
  font-weight: 700;
}

.step-success .step-dot {
  border-color: var(--clr-success);
  background: var(--clr-success);
  box-shadow: 0 0 0 4px var(--clr-success-light);
}

.step-success .step-label {
  color: var(--clr-success);
}

.step-danger .step-dot {
  border-color: var(--clr-danger);
  background: var(--clr-danger);
  box-shadow: 0 0 0 4px var(--clr-danger-light);
}

.step-danger .step-dot::after {
  content: '✕';
}

.step-danger .step-label {
  color: var(--clr-danger);
}

@keyframes pulse-dot {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.15); }
}

/* Progress lines */
.progress-line {
  flex: 1;
  height: 3px;
  min-width: 20px;
  border-radius: 2px;
  transition: background 0.3s ease;
  margin: 0 4px;
  margin-bottom: 20px; /* offset for step-label height */
}

.line-filled {
  background: var(--clr-secondary);
}

.line-empty {
  background: var(--clr-border);
}

/* Notes and actions */
.app-card-notes {
  background: var(--clr-surface-alt);
  border-left: 3px solid var(--clr-primary);
  padding: 10px 14px;
  border-radius: var(--border-radius-sm);
  font-size: 0.85rem;
  color: var(--clr-text-muted);
  margin-top: 8px;
}

.app-card-actions {
  margin-top: 16px;
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.spin-anim {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* Mobile */
@media (max-width: 640px) {
  .app-card-header {
    flex-direction: column;
  }

  .status-progress {
    padding: 0;
  }

  .step-label {
    font-size: 0.65rem;
  }

  .step-dot {
    width: 24px;
    height: 24px;
  }

  .app-tracker-card {
    padding: 16px;
  }
}
</style>
