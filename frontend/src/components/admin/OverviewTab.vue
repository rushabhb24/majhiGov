<script setup>
import { computed } from 'vue'

const props = defineProps({
  analytics: {
    type: Object,
    default: () => ({})
  },
  schemes: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['nav-tab', 'view-expiring'])

const stats = computed(() => {
  const sa = props.analytics || {}
  const totalSch = sa.total_schemes || 142
  const totalUsr = sa.total_users || 18430
  const totalApp = sa.total_applications || 24
  const pendingApp = sa.pending_applications || 7

  return [
    { num: totalSch, label: 'Total Schemes', desc: '+8 this month', icon: 'ti-files', class: 'blue' },
    { num: totalUsr.toLocaleString(), label: 'Registered Users', desc: '+1,204 this week', icon: 'ti-users', class: 'green' },
    { num: totalApp, label: 'Total Applications', desc: 'Real database count', icon: 'ti-checklist', class: 'orange' },
    { num: pendingApp, label: 'Pending Approvals', desc: 'Awaiting review', icon: 'ti-clock', class: 'red' }
  ]
})

const categoryProgress = computed(() => {
  // exact mock values and styling from md file
  return [
    { emoji: '🌾', name: 'Farmers', fill: 85, count: 34 },
    { emoji: '🎓', name: 'Students', fill: 70, count: 28 },
    { emoji: '👩', name: 'Women', fill: 55, count: 22 },
    { emoji: '💼', name: 'Business', fill: 45, count: 18 },
    { emoji: '👴', name: 'Senior Citizens', fill: 35, count: 14 },
    { emoji: '♿', name: 'Disabled', fill: 22, count: 9 }
  ]
})

const recentActivity = computed(() => {
  const sa = props.analytics || {}
  const list = sa.recent_activity || []
  if (list.length === 0) {
    return [
      { dot: 'green', text: 'New scheme PM Vishwakarma Yojana added to Business category', time: '2 hours ago' },
      { dot: 'orange', text: 'Scheme Ladli Behna deadline updated to 31 March 2025', time: '5 hours ago' },
      { dot: 'blue', text: 'User Ramesh Kumar registered from Rajasthan', time: '1 day ago' },
      { dot: 'green', text: 'Category Senior Citizens updated with 2 new schemes', time: '2 days ago' }
    ]
  }
  return list.map(act => ({
    dot: act.type === 'scheme' ? 'green' : (act.type === 'user' ? 'blue' : 'orange'),
    text: act.text,
    time: act.time_ago || 'Just now'
  }))
})

const recentApplications = computed(() => {
  const sa = props.analytics || {}
  return sa.recent_applications || [
    { id: 1, full_name: "Ramesh Kumar", scheme_name: "PM Kisan Samman Nidhi", status: "approved" },
    { id: 2, full_name: "Priya Sharma", scheme_name: "NSP Post Matric Scholarship", status: "pending" },
    { id: 3, full_name: "Amit Joshi", scheme_name: "PM Mudra Loan", status: "rejected" },
    { id: 4, full_name: "Sunita Patil", scheme_name: "Lado Deviprasad Scheme", status: "pending" }
  ]
})
</script>

<template>
  <div class="overview-tab">
    
    <!-- Stats Row — 4 cards in a grid -->
    <div class="stats-grid">
      <div class="stat-card" v-for="s in stats" :key="s.label">
        <div :class="['stat-icon-box', s.class]">
          <i :class="['ti', s.icon]"></i>
        </div>
        <div class="stat-number">{{ s.num }}</div>
        <div class="stat-label">{{ s.label }}</div>
        <div class="stat-desc" :class="s.class === 'red' ? 'text-danger' : 'text-success'">
          <i class="ti ti-trending-up" v-if="s.class === 'blue' || s.class === 'green'"></i>
          <i class="ti ti-check" v-if="s.class === 'orange'"></i>
          <i class="ti ti-alert-circle" v-if="s.class === 'red'"></i>
          <span>{{ s.desc }}</span>
        </div>
      </div>
    </div>

    <!-- Warning Alert Banner -->
    <div class="warning-alert-banner">
      <div class="alert-content">
        <i class="ti ti-alert-triangle"></i>
        <span>Review operational applications and statistics. Ensure timely compliance updates.</span>
      </div>
      <span class="alert-link" @click="emit('nav-tab', 'users')">Manage Users</span>
    </div>

    <!-- Two-Column Section (equal width) -->
    <div class="two-col-grid mt-4">
      
      <!-- LEFT CARD — "Recent Applications" -->
      <div class="card">
        <div class="card-header">
          <div class="card-title">Recent Applications</div>
          <span class="card-header-link" @click="emit('nav-tab', 'users')">
            <span>View all</span>
            <i class="ti ti-arrow-right"></i>
          </span>
        </div>
        <div class="card-body p-0">
          <table class="overview-table">
            <thead>
              <tr>
                <th>Applicant</th>
                <th>Scheme</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="app in recentApplications" :key="app.id">
                <td>
                  <div class="scheme-title">{{ app.full_name }}</div>
                  <div class="scheme-subtitle">Applicant</div>
                </td>
                <td>
                  <div class="scheme-title" style="max-width:180px; overflow:hidden; text-overflow:ellipsis; white-space:nowrap;">{{ app.scheme_name }}</div>
                </td>
                <td>
                  <span :class="['badge', app.status === 'approved' ? 'active' : (app.status === 'rejected' ? 'expiring' : 'central')]">
                    {{ app.status }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- RIGHT CARD — "Schemes by Category" -->
      <div class="card">
        <div class="card-header">
          <div class="card-title">Schemes by Category</div>
        </div>
        <div class="card-body">
          <div class="category-row" v-for="cat in categoryProgress" :key="cat.name">
            <span class="cat-label-text">{{ cat.emoji }} {{ cat.name }}</span>
            <span class="arrow-indicator">→</span>
            <div class="progress-bar-container">
              <div class="progress-bar-fill" :style="{ width: cat.fill + '%' }"></div>
            </div>
            <span class="cat-count-text">{{ cat.count }}</span>
          </div>
        </div>
      </div>

    </div>

    <!-- Activity Feed Card (full width) -->
    <div class="card mt-4">
      <div class="card-header">
        <div class="card-title">Recent Activity</div>
      </div>
      <div class="card-body p-0">
        <div class="activity-feed">
          <div class="activity-item" v-for="(act, idx) in recentActivity" :key="idx">
            <span :class="['dot-indicator', act.dot]"></span>
            <div class="activity-content">
              <div class="activity-text">{{ act.text }}</div>
              <div class="activity-time">{{ act.time }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<style scoped>
.overview-tab {
  width: 100%;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
  width: 100%;
}

.stat-card {
  background-color: #ffffff; /* var(--bg) */
  border: 0.5px solid rgba(0, 0, 0, 0.08); /* var(--border) */
  border-radius: 8px; /* var(--radius) */
  padding: 14px 16px;
  box-sizing: border-box;
}

.stat-icon-box {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
}

.stat-icon-box i {
  font-size: 18px !important;
}

.stat-icon-box.blue { background-color: #e8eef8; color: #1a3a6b; }
.stat-icon-box.green { background-color: #f0fdf4; color: #16a34a; }
.stat-icon-box.orange { background-color: #fff4ed; color: #f97316; }
.stat-icon-box.red { background-color: #fef2f2; color: #dc2626; }

.stat-number {
  font-size: 22px;
  font-weight: 500;
  color: #0f172a;
  margin-top: 10px;
}

.stat-label {
  font-size: 12px;
  color: #64748b;
  margin-top: 2px;
}

.stat-desc {
  font-size: 11px;
  font-weight: 500;
  margin-top: 8px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.stat-desc i {
  font-size: 12px !important;
}

.text-success { color: #16a34a; }
.text-danger { color: #dc2626; }

/* Warning Alert Banner */
.warning-alert-banner {
  background-color: #fffbeb; /* var(--warning-bg) */
  border: 0.5px solid #fbbf24;
  border-radius: 8px;
  padding: 10px 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 12px;
  box-sizing: border-box;
}

.alert-content {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
  color: #0f172a;
}

.alert-content i {
  color: #d97706 !important;
  font-size: 16px !important;
}

.alert-link {
  font-size: 13px;
  font-weight: 500;
  text-decoration: underline;
  cursor: pointer;
  color: #1a3a6b;
}

/* Grids */
.two-col-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  align-items: start;
}

.mt-4 { margin-top: 16px; }

.card {
  background-color: #ffffff;
  border: 0.5px solid rgba(0, 0, 0, 0.08);
  border-radius: 12px; /* var(--radius-lg) */
  overflow: hidden;
  width: 100%;
  box-sizing: border-box;
}

.card-header {
  padding: 14px 16px;
  border-bottom: 0.5px solid rgba(0, 0, 0, 0.08);
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-sizing: border-box;
}

.card-title {
  font-size: 13px;
  font-weight: 500;
  color: #0f172a;
}

.card-header-link {
  color: #1a3a6b; /* var(--primary) */
  font-weight: 500;
  font-size: 13px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
}

.card-header-link:hover {
  text-decoration: underline;
}

.card-header-link i {
  font-size: 13px !important;
}

.card-body {
  padding: 16px;
}

.p-0 { padding: 0 !important; }

/* Overview Table */
.overview-table {
  width: 100%;
  border-collapse: collapse;
}

.overview-table th {
  background-color: #f8fafc; /* var(--bg2) */
  text-align: left;
  font-size: 11px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #64748b;
  padding: 8px 12px;
  border-bottom: 0.5px solid rgba(0, 0, 0, 0.08);
}

.overview-table td {
  padding: 10px 12px;
  border-bottom: 0.5px solid rgba(0, 0, 0, 0.08);
  font-size: 13px;
}

.overview-table tr:last-child td {
  border-bottom: none;
}

.overview-table tr:hover td {
  background-color: #f8fafc;
}

.scheme-title {
  font-weight: 500;
  color: #0f172a;
}

.scheme-subtitle {
  font-size: 11px;
  color: #64748b;
  margin-top: 1px;
}

/* Badges */
.badge {
  display: inline-flex;
  align-items: center;
  padding: 3px 8px;
  border-radius: 100px;
  font-size: 11px;
  font-weight: 500;
  line-height: 1;
}

.badge.active { background-color: #f0fdf4; color: #16a34a; }
.badge.expiring { background-color: #fffbeb; color: #d97706; }
.badge.central { background-color: #e8eef8; color: #1a3a6b; }
.badge.state { background-color: #fff4ed; color: #f97316; }

/* Category fill list */
.category-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
}

.category-row:last-child {
  margin-bottom: 0;
}

.cat-label-text {
  width: 100px;
  font-weight: 500;
  font-size: 13px;
  color: #0f172a;
}

.arrow-indicator {
  color: #64748b;
  font-weight: 400;
  width: 12px;
}

.progress-bar-container {
  flex-grow: 1;
  height: 4px;
  background-color: #f8fafc; /* var(--bg2) */
  border-radius: 2px;
  overflow: hidden;
}

.progress-bar-fill {
  height: 100%;
  background-color: #1a3a6b; /* var(--primary) */
  border-radius: 2px;
}

.cat-count-text {
  width: 25px;
  text-align: right;
  font-weight: 500;
  font-size: 13px;
  color: #0f172a;
}

/* Activity Feed */
.activity-feed {
  display: flex;
  flex-direction: column;
}

.activity-item {
  display: flex;
  gap: 12px;
  padding: 12px 16px;
  border-bottom: 0.5px solid rgba(0, 0, 0, 0.08);
  align-items: flex-start;
  box-sizing: border-box;
}

.activity-item:last-child {
  border-bottom: none;
}

.dot-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  margin-top: 5px;
  flex-shrink: 0;
}

.dot-indicator.green { background-color: #16a34a; }
.dot-indicator.orange { background-color: #d97706; }
.dot-indicator.blue { background-color: #1a3a6b; }

.activity-content {
  flex-grow: 1;
}

.activity-text {
  font-size: 13px;
  color: #0f172a;
}

.activity-time {
  font-size: 11px;
  color: #64748b;
  margin-top: 2px;
}
</style>
