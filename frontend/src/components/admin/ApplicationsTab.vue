<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  applications: {
    type: Array,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update-status'])

// Filter States
const searchQuery = ref('')
const statusFilter = ref('All')
const levelFilter = ref('All')

// Modal States for Review
const selectedApp = ref(null)
const isModalOpen = ref(false)
const actionType = ref('') // 'approve', 'reject', 'view'
const remarksText = ref('')

const stats = computed(() => {
  const tot = props.applications.length
  const pending = props.applications.filter(a => a.status === 'pending').length
  const approved = props.applications.filter(a => a.status === 'approved').length
  const rejected = props.applications.filter(a => a.status === 'rejected').length

  return [
    { num: tot.toLocaleString(), label: 'Total Applications', icon: 'ti-checklist', class: 'blue' },
    { num: pending.toLocaleString(), label: 'Pending Approvals', icon: 'ti-clock', class: 'orange' },
    { num: approved.toLocaleString(), label: 'Approved Applications', icon: 'ti-checkbox', class: 'green' },
    { num: rejected.toLocaleString(), label: 'Rejected Applications', icon: 'ti-circle-x', class: 'red' }
  ]
})

const filteredApplications = computed(() => {
  let list = props.applications

  // Search filter
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(a => 
      (a.full_name || '').toLowerCase().includes(q) || 
      (a.scheme_title || '').toLowerCase().includes(q) || 
      (a.email || '').toLowerCase().includes(q)
    )
  }

  // Status Filter
  if (statusFilter.value !== 'All') {
    list = list.filter(a => a.status === statusFilter.value.toLowerCase())
  }

  // Level Filter
  if (levelFilter.value !== 'All') {
    list = list.filter(a => a.government_level === levelFilter.value.toLowerCase())
  }

  return list
})

function openReviewModal(app, type) {
  selectedApp.value = app
  actionType.value = type
  remarksText.value = app.notes || ''
  isModalOpen.value = true
}

function closeModal() {
  selectedApp.value = null
  isModalOpen.value = false
  remarksText.value = ''
}

function submitAction() {
  if (!selectedApp.value) return
  emit('update-status', {
    applicationId: selectedApp.value.id,
    status: actionType.value === 'approve' ? 'approved' : 'rejected',
    notes: remarksText.value
  })
  closeModal()
}

function getAvatarInitials(name) {
  if (!name) return 'U'
  const parts = name.split(' ')
  if (parts.length > 1) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return name.substring(0, 2).toUpperCase()
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: 'numeric' })
}
</script>

<template>
  <div class="applications-tab">
    
    <!-- Stats Grid Row -->
    <div class="stats-grid">
      <div class="stat-card" v-for="s in stats" :key="s.label">
        <div :class="['stat-icon-box', s.class]">
          <i :class="['ti', s.icon]"></i>
        </div>
        <div class="stat-number">{{ s.num }}</div>
        <div class="stat-label">{{ s.label }}</div>
      </div>
    </div>

    <!-- Filter & Search Panel -->
    <div class="filter-panel mt-4">
      <div class="search-box">
        <i class="ti ti-search search-icon"></i>
        <input 
          type="text" 
          v-model="searchQuery" 
          placeholder="Search citizen name, scheme title, email..." 
          class="search-input"
        />
      </div>

      <div class="filter-controls">
        <select v-model="statusFilter" class="filter-select">
          <option value="All">All Statuses</option>
          <option value="Pending">Pending</option>
          <option value="Approved">Approved</option>
          <option value="Rejected">Rejected</option>
        </select>

        <select v-model="levelFilter" class="filter-select">
          <option value="All">All Levels</option>
          <option value="Central">Central Schemes</option>
          <option value="State">State Schemes</option>
        </select>
      </div>
    </div>

    <!-- Data Table Container -->
    <div class="card mt-3">
      <div class="card-body p-0">
        <table class="data-table">
          <thead>
            <tr>
              <th>Citizen Details</th>
              <th>Scheme Title</th>
              <th>Date Applied</th>
              <th>Status</th>
              <th>Remarks / Notes</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="app in filteredApplications" :key="app.id">
              <td>
                <div class="citizen-info">
                  <div class="avatar-circle">
                    {{ getAvatarInitials(app.full_name) }}
                  </div>
                  <div>
                    <div class="citizen-name">{{ app.full_name || 'Citizen User' }}</div>
                    <div class="citizen-meta">
                      <span><i class="ti ti-mail"></i> {{ app.email }}</span>
                      <span class="ml-2"><i class="ti ti-phone"></i> {{ app.phone }}</span>
                    </div>
                  </div>
                </div>
              </td>
              <td>
                <div class="scheme-title-cell">{{ app.scheme_title }}</div>
                <span :class="['level-badge', app.government_level === 'central' ? 'central' : 'state']">
                  {{ app.government_level === 'central' ? 'Central' : 'State' }}
                </span>
              </td>
              <td class="date-cell">{{ formatDate(app.applied_at) }}</td>
              <td>
                <span :class="['badge-status', app.status]">
                  <i class="ti" :class="app.status === 'approved' ? 'ti-circle-check' : (app.status === 'rejected' ? 'ti-circle-x' : 'ti-clock')"></i>
                  {{ app.status }}
                </span>
              </td>
              <td class="remarks-cell" :title="app.notes">
                {{ app.notes || '-' }}
              </td>
              <td>
                <div class="table-actions">
                  <button 
                    class="action-btn" 
                    title="View Details" 
                    @click="openReviewModal(app, 'view')"
                  >
                    <i class="ti ti-eye"></i>
                  </button>
                  <button 
                    v-if="app.status === 'pending'"
                    class="action-btn approve-btn" 
                    title="Approve Application" 
                    @click="openReviewModal(app, 'approve')"
                  >
                    <i class="ti ti-check"></i>
                  </button>
                  <button 
                    v-if="app.status === 'pending'"
                    class="action-btn reject-btn" 
                    title="Reject Application" 
                    @click="openReviewModal(app, 'reject')"
                  >
                    <i class="ti ti-x"></i>
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="filteredApplications.length === 0">
              <td colspan="6" class="no-records">
                <i class="ti ti-inbox no-records-icon"></i>
                <div class="no-records-text">No active applications match your filter selection.</div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Dynamic Review / View Details Modal -->
    <div class="admin-modal-overlay" v-if="isModalOpen" @click.self="closeModal">
      <div class="admin-modal-box" style="max-width: 550px;">
        <button class="btn-close" @click="closeModal">×</button>
        
        <div class="modal-header-section">
          <div class="modal-indicator" :class="actionType">
            <i class="ti" :class="actionType === 'approve' ? 'ti-circle-check' : (actionType === 'reject' ? 'ti-circle-x' : 'ti-eye')"></i>
          </div>
          <div class="modal-title-text">
            <span v-if="actionType === 'view'">Citizen Application Details</span>
            <span v-else-if="actionType === 'approve'">Approve Citizen Application</span>
            <span v-else-if="actionType === 'reject'">Reject Citizen Application</span>
          </div>
        </div>

        <div class="modal-review-body mt-3">
          <!-- Citizen Overview Cards -->
          <div class="review-sec">
            <div class="sec-label">Applicant Citizen Details</div>
            <div class="citizen-detail-grid">
              <div class="detail-item">
                <span class="lbl">Full Name</span>
                <span class="val font-semibold">{{ selectedApp?.full_name }}</span>
              </div>
              <div class="detail-item">
                <span class="lbl">Email Address</span>
                <span class="val">{{ selectedApp?.email }}</span>
              </div>
              <div class="detail-item">
                <span class="lbl">Mobile Number</span>
                <span class="val">{{ selectedApp?.phone }}</span>
              </div>
              <div class="detail-item">
                <span class="lbl">User Database ID</span>
                <span class="val">#{{ selectedApp?.user_id }}</span>
              </div>
            </div>
          </div>

          <div class="review-sec mt-3">
            <div class="sec-label">Applied Scheme Details</div>
            <div class="scheme-details-card">
              <div class="scheme-details-title">{{ selectedApp?.scheme_title }}</div>
              <div class="scheme-details-meta">
                <span :class="['level-badge', selectedApp?.government_level === 'central' ? 'central' : 'state']">
                  {{ selectedApp?.government_level === 'central' ? 'Central Scheme' : 'State Scheme' }}
                </span>
                <span class="date"><i class="ti ti-calendar-event"></i> Applied On: {{ formatDate(selectedApp?.applied_at) }}</span>
              </div>
            </div>
          </div>

          <!-- Application Remarks Review / Submission -->
          <div class="review-sec mt-3">
            <div class="sec-label">
              <span v-if="actionType === 'view'">Administrative Remarks</span>
              <span v-else>Provide Decision Remarks (Optional)</span>
            </div>
            
            <div v-if="actionType === 'view'" class="remarks-box-readonly">
              <div class="status-indicator">
                Status: <span :class="['status-val', selectedApp?.status]">{{ selectedApp?.status }}</span>
              </div>
              <p class="remarks-desc">{{ selectedApp?.notes || 'No administrative remarks submitted.' }}</p>
            </div>
            
            <div v-else class="form-group">
              <textarea 
                v-model="remarksText"
                class="form-input" 
                rows="3" 
                placeholder="e.g. Verified Aadhaar and income criteria successfully. Approved."
              ></textarea>
            </div>
          </div>
        </div>

        <!-- Action Footer Buttons -->
        <div class="modal-actions-footer mt-4">
          <button type="button" class="btn-cancel" @click="closeModal">Cancel</button>
          
          <button 
            v-if="actionType === 'approve'" 
            type="button" 
            class="btn-submit approve"
            @click="submitAction"
          >
            <i class="ti ti-circle-check"></i> Finalize Approval
          </button>
          
          <button 
            v-else-if="actionType === 'reject'" 
            type="button" 
            class="btn-submit reject"
            @click="submitAction"
          >
            <i class="ti ti-circle-x"></i> Reject Application
          </button>
        </div>
      </div>
    </div>

  </div>
</template>

<style scoped>
.applications-tab {
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
.stat-icon-box.orange { background-color: var(--accent-light); color: var(--accent); }
.stat-icon-box.green { background-color: var(--success-bg); color: var(--success); }
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

.citizen-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.avatar-circle {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background-color: var(--primary-light);
  color: var(--primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 12px;
}

.citizen-name {
  font-weight: 500;
  color: var(--text);
}

.citizen-meta {
  font-size: 10px;
  color: var(--text2);
  margin-top: 1px;
  display: flex;
  gap: 8px;
}

.citizen-meta span {
  display: inline-flex;
  align-items: center;
  gap: 2px;
}

.scheme-title-cell {
  font-weight: 500;
  color: var(--text);
  max-width: 220px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.level-badge {
  display: inline-block;
  font-size: 9px;
  font-weight: 500;
  padding: 1px 5px;
  border-radius: 4px;
  margin-top: 2px;
  text-transform: uppercase;
}

.level-badge.central { background-color: var(--primary-light); color: var(--primary); }
.level-badge.state { background-color: var(--accent-light); color: var(--accent); }

.badge-status {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 8px;
  border-radius: 100px;
  font-size: 11px;
  font-weight: 500;
  text-transform: capitalize;
}

.badge-status.pending { background-color: var(--warning-bg); color: var(--warning); }
.badge-status.approved { background-color: var(--success-bg); color: var(--success); }
.badge-status.rejected { background-color: var(--danger-bg); color: var(--danger); }

.remarks-cell {
  color: var(--text2);
  font-size: 12px;
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

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

/* Modal Headers */
.modal-header-section {
  display: flex;
  align-items: center;
  gap: 12px;
  border-bottom: 0.5px solid var(--border);
  padding-bottom: 12px;
}

.modal-indicator {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
}

.modal-indicator.view { background-color: var(--primary-light); color: var(--primary); }
.modal-indicator.approve { background-color: var(--success-bg); color: var(--success); }
.modal-indicator.reject { background-color: var(--danger-bg); color: var(--danger); }

.review-sec {
  border-bottom: 0.5px solid var(--border);
  padding-bottom: 14px;
}

.review-sec:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.sec-label {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  color: var(--text2);
  margin-bottom: 8px;
}

.citizen-detail-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
}

.detail-item {
  display: flex;
  flex-direction: column;
}

.detail-item .lbl {
  font-size: 10px;
  color: var(--text2);
}

.detail-item .val {
  font-size: 13px;
  color: var(--text);
  margin-top: 1px;
}

.scheme-details-card {
  background-color: var(--bg2);
  border: 0.5px solid var(--border);
  border-radius: 8px;
  padding: 10px 12px;
}

.scheme-details-title {
  font-weight: 600;
  font-size: 13px;
  color: var(--text);
}

.scheme-details-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 6px;
  font-size: 11px;
  color: var(--text2);
}

.scheme-details-meta span.date {
  display: inline-flex;
  align-items: center;
  gap: 2px;
}

.remarks-box-readonly {
  background-color: var(--bg2);
  border: 0.5px solid var(--border);
  border-radius: 8px;
  padding: 10px 12px;
}

.remarks-box-readonly .status-indicator {
  font-size: 11px;
  font-weight: 500;
  color: var(--text);
}

.remarks-box-readonly .status-val {
  text-transform: uppercase;
  font-size: 10px;
  font-weight: 600;
  padding: 1px 5px;
  border-radius: 4px;
}

.remarks-box-readonly .status-val.pending { background-color: var(--warning-bg); color: var(--warning); }
.remarks-box-readonly .status-val.approved { background-color: var(--success-bg); color: var(--success); }
.remarks-box-readonly .status-val.rejected { background-color: var(--danger-bg); color: var(--danger); }

.remarks-box-readonly .remarks-desc {
  font-size: 12px;
  color: var(--text2);
  margin-top: 6px;
  line-height: 1.4;
}

.modal-actions-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.btn-cancel {
  padding: 6px 12px;
  border-radius: 6px;
  border: 0.5px solid var(--border);
  background-color: var(--bg);
  color: var(--text2);
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
}

.btn-cancel:hover {
  background-color: var(--bg2);
  color: var(--text);
}

.btn-submit {
  padding: 6px 12px;
  border-radius: 6px;
  border: none;
  color: #fff;
  font-size: 12px;
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
}

.btn-submit.approve { background-color: var(--success); }
.btn-submit.approve:hover { filter: brightness(0.95); }
.btn-submit.reject { background-color: var(--danger); }
.btn-submit.reject:hover { filter: brightness(0.95); }

.ml-2 { margin-left: 8px; }
.mt-4 { margin-top: 16px; }
.mt-3 { margin-top: 12px; }
.mt-2 { margin-top: 8px; }
.font-semibold { font-weight: 600; }
</style>
