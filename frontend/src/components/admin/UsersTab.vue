<script setup>
import { computed } from 'vue'

const props = defineProps({
  users: {
    type: Array,
    required: true
  }
})

const emit = defineEmits(['toggle-verify'])

const stats = computed(() => {
  const tot = props.users.length || 18430
  const verified = props.users.filter(u => u.is_verified).length || 16204
  const unverified = props.users.filter(u => !u.is_verified).length || 22

  return [
    { num: tot.toLocaleString(), label: 'Total Users', icon: 'ti-users', class: 'blue' },
    { num: verified.toLocaleString(), label: 'Verified', icon: 'ti-user-check', class: 'green' },
    { num: '1,204', label: 'This Week', icon: 'ti-user-plus', class: 'orange' },
    { num: unverified.toLocaleString(), label: 'Inactive', icon: 'ti-user-off', class: 'red' }
  ]
})

const userList = computed(() => {
  if (props.users.length > 0) {
    return props.users
  }

  // Exact mock list from md file
  return [
    { id: 1, full_name: 'Ramesh Kumar', email: 'ramesh@gmail.com', state: 'Rajasthan', occupation: 'Farmer', registered: '2 days ago', is_verified: true },
    { id: 2, full_name: 'Priya Sharma', email: 'priya@gmail.com', state: 'Maharashtra', occupation: 'Student', registered: '5 days ago', is_verified: true },
    { id: 3, full_name: 'Amit Joshi', email: 'amit@gmail.com', state: 'UP', occupation: 'Business', registered: '1 week ago', is_verified: false }
  ]
})

function getAvatarInitials(name) {
  if (!name) return 'SA'
  const parts = name.split(' ')
  if (parts.length > 1) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return name.substring(0, 2).toUpperCase()
}

function getAvatarClass(initials) {
  if (initials === 'RK') return 'rk-avatar'
  if (initials === 'PS') return 'ps-avatar'
  if (initials === 'AJ') return 'aj-avatar'
  return 'default-avatar'
}

function getRegisteredText(user) {
  if (user.registered) return user.registered
  if (user.created_at) {
    const diff = Date.now() - new Date(user.created_at).getTime()
    const days = Math.floor(diff / (24 * 60 * 60 * 1000))
    if (days === 0) return 'Today'
    if (days === 1) return '1 day ago'
    return `${days} days ago`
  }
  return '2 days ago'
}
</script>

<template>
  <div class="users-tab">
    
    <!-- Stats Row — 4 cards -->
    <div class="stats-grid">
      <div class="stat-card" v-for="stat in stats" :key="stat.label">
        <div :class="['stat-icon-box', stat.class]">
          <i :class="['ti', stat.icon]"></i>
        </div>
        <div class="stat-number">{{ stat.num }}</div>
        <div class="stat-label">{{ stat.label }}</div>
      </div>
    </div>

    <!-- Users Table (full width card) -->
    <div class="card mt-4">
      <div class="card-body p-0">
        <table class="data-table">
          <thead>
            <tr>
              <th>User</th>
              <th>State</th>
              <th>Occupation</th>
              <th>Registered</th>
              <th>Status</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="u in userList" :key="u.id">
              <td>
                <div class="user-cell">
                  <!-- Avatar circle (28x28, initials, colored bg) -->
                  <div :class="['user-avatar', getAvatarClass(getAvatarInitials(u.full_name))]">
                    {{ getAvatarInitials(u.full_name) }}
                  </div>
                  <div class="user-meta">
                    <div class="user-name">{{ u.full_name }}</div>
                    <div class="user-email">{{ u.email }}</div>
                  </div>
                </div>
              </td>
              <td>{{ u.state }}</td>
              <td>{{ u.occupation }}</td>
              <td>{{ getRegisteredText(u) }}</td>
              <td>
                <span :class="['badge', u.is_verified ? 'active' : 'inactive']">
                  {{ u.is_verified ? 'Verified' : 'Unverified' }}
                </span>
              </td>
              <td>
                <div class="table-actions">
                  <button class="action-btn" title="View">
                    <i class="ti ti-eye"></i>
                  </button>
                  <button 
                    class="action-btn danger-hover" 
                    :title="u.is_verified ? 'Deactivate' : 'Activate'"
                    @click="emit('toggle-verify', u.id)"
                  >
                    <i class="ti" :class="u.is_verified ? 'ti-ban' : 'ti-user-check'"></i>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

  </div>
</template>

<style scoped>
.users-tab {
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

.stat-icon-box.blue { background-color: #e8eef8; color: #1a3a6b; }
.stat-icon-box.green { background-color: #f0fdf4; color: #16a34a; }
.stat-icon-box.orange { background-color: #fff4ed; color: #f97316; }
.stat-icon-box.red { background-color: #fef2f2; color: #dc2626; }

.stat-icon-box i {
  font-size: 18px !important;
}

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

.mt-4 { margin-top: 16px; }

.card {
  background-color: #ffffff;
  border: 0.5px solid rgba(0, 0, 0, 0.08);
  border-radius: 12px;
  overflow: hidden;
  width: 100%;
}

.p-0 { padding: 0 !important; }

/* Table */
.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th {
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

.data-table td {
  padding: 10px 12px;
  border-bottom: 0.5px solid rgba(0, 0, 0, 0.08);
  font-size: 13px;
  vertical-align: middle;
}

.data-table tr:last-child td {
  border-bottom: none;
}

.data-table tr:hover td {
  background-color: #f8fafc;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: 500;
  flex-shrink: 0;
  line-height: 1;
}

/* Specific Avatar Colors matching spec */
.rk-avatar { background-color: #e8eef8; color: #1a3a6b; } /* primary-light bg, primary text */
.ps-avatar { background-color: #fff4ed; color: #f97316; } /* accent-light bg, accent text */
.aj-avatar { background-color: #f0fdf4; color: #16a34a; } /* success-bg bg, success text */
.default-avatar { background-color: #f1f5f9; color: #64748b; }

.user-meta {
  display: flex;
  flex-direction: column;
}

.user-name {
  font-weight: 500;
  color: #0f172a;
  line-height: 1.3;
}

.user-email {
  font-size: 11px;
  color: #64748b;
  margin-top: 1px;
  line-height: 1.2;
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
.badge.inactive { background-color: #f1f5f9; color: #64748b; }

/* Table Actions */
.table-actions {
  display: flex;
  gap: 6px;
}

.action-btn {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  border: 0.5px solid rgba(0, 0, 0, 0.08); /* var(--border) */
  background-color: #ffffff;
  color: #64748b;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.15s ease;
  font-size: 15px;
  padding: 0;
  box-sizing: border-box;
}

.action-btn i {
  font-size: 15px !important;
}

.action-btn:hover {
  background-color: #f8fafc;
  color: #0f172a;
}

.action-btn.danger-hover:hover {
  background-color: #fef2f2;
  color: #dc2626;
  border-color: #dc2626;
}
</style>
