<script setup>
import { computed } from 'vue'
import { useAuthStore } from '../../stores/auth'

defineProps({
  activeTab: {
    type: String,
    required: true
  },
  schemesCount: {
    type: Number,
    default: 142
  },
  notificationsCount: {
    type: Number,
    default: 3
  },
  applicationsCount: {
    type: Number,
    default: 0
  },
  theme: {
    type: String,
    default: 'dark'
  }
})

const emit = defineEmits(['update:activeTab'])
const authStore = useAuthStore()

// Avatar configurations
const avatarPresets = [
  { url: 'avatar-dir-blue', bg: '#e8eef8', color: '#1a3a6b', icon: 'ti-user-check' },
  { url: 'avatar-adm-amber', bg: '#fff4ed', color: '#f97316', icon: 'ti-shield-check' },
  { url: 'avatar-off-emerald', bg: '#f0fdf4', color: '#16a34a', icon: 'ti-medal' },
  { url: 'avatar-mgr-rose', bg: '#fef2f2', color: '#dc2626', icon: 'ti-briefcase' },
  { url: 'avatar-exe-violet', bg: '#faf5ff', color: '#7c3aed', icon: 'ti-award' },
  { url: 'avatar-sup-teal', bg: '#f0fdfa', color: '#0d9488', icon: 'ti-id' },
  { url: 'avatar-spe-cyan', bg: '#ecfeff', color: '#0891b2', icon: 'ti-user-plus' },
  { url: 'avatar-coo-indigo', bg: '#eef2ff', color: '#4f46e5', icon: 'ti-settings' }
]

const adminProfile = computed(() => authStore.userProfile || {})
const adminName = computed(() => adminProfile.value.full_name || 'Super Admin')
const adminInitials = computed(() => {
  const name = adminName.value
  return name ? name.split(' ').map(n => n[0]).join('').slice(0, 2).toUpperCase() : 'SA'
})
const currentPreset = computed(() => {
  const url = adminProfile.value.avatar_url
  return avatarPresets.find(p => p.url === url) || null
})

function selectTab(tab) {
  emit('update:activeTab', tab)
}
</script>

<template>
  <aside :class="['sidebar', theme]">
    <!-- Brand Area (top of sidebar) -->
    <div class="brand-area">
      <div class="logo-box">
        <i class="ti ti-building-bank"></i>
      </div>
      <div class="brand-text">
        <div class="app-name">Yojana Portal</div>
        <div class="panel-subtitle">Admin Panel</div>
      </div>
    </div>

    <!-- Navigation Menu -->
    <div class="nav-menu">
      <div class="nav-section-label">PORTAL</div>
      <div 
        class="sb-item return-portal-btn" 
        @click="$router.push('/')"
        id="nav-return-portal"
        title="Go back to the Main citizen portal"
      >
        <i class="ti ti-arrow-back-up"></i>
        <span>Return to Portal</span>
      </div>

      <div class="nav-section-label">MAIN</div>
      <div 
        :class="['sb-item', { active: activeTab === 'overview' }]" 
        @click="selectTab('overview')"
        id="nav-overview"
      >
        <i class="ti ti-layout-dashboard"></i>
        <span>Overview</span>
      </div>
      <div 
        :class="['sb-item', { active: activeTab === 'schemes' }]" 
        @click="selectTab('schemes')"
        id="nav-schemes"
      >
        <i class="ti ti-files"></i>
        <span>Schemes</span>
        <span class="sb-badge bg-accent">{{ schemesCount }}</span>
      </div>
      <div 
        :class="['sb-item', { active: activeTab === 'categories' }]" 
        @click="selectTab('categories')"
        id="nav-categories"
      >
        <i class="ti ti-grid-dots"></i>
        <span>Categories</span>
      </div>

      <div class="nav-section-label">MANAGEMENT</div>
      <div 
        :class="['sb-item', { active: activeTab === 'users' }]" 
        @click="selectTab('users')"
        id="nav-users"
      >
        <i class="ti ti-users"></i>
        <span>Users</span>
      </div>
      <div 
        :class="['sb-item', { active: activeTab === 'applications' }]" 
        @click="selectTab('applications')"
        id="nav-applications"
      >
        <i class="ti ti-briefcase"></i>
        <span>Applications</span>
        <span class="sb-badge bg-accent" v-if="applicationsCount > 0">{{ applicationsCount }}</span>
      </div>
      <div 
        :class="['sb-item', { active: activeTab === 'eligibility' }]" 
        @click="selectTab('eligibility')"
        id="nav-eligibility"
      >
        <i class="ti ti-checklist"></i>
        <span>Eligibility Rules</span>
      </div>
      <div 
        :class="['sb-item', { active: activeTab === 'notifications' }]" 
        @click="selectTab('notifications')"
        id="nav-notifications"
      >
        <i class="ti ti-bell"></i>
        <span>Notifications</span>
        <span class="sb-badge bg-accent">{{ notificationsCount }}</span>
      </div>

      <div class="nav-section-label">SYSTEM</div>
      <div 
        :class="['sb-item', { active: activeTab === 'analytics' }]" 
        @click="selectTab('analytics')"
        id="nav-analytics"
      >
        <i class="ti ti-chart-bar"></i>
        <span>Analytics</span>
      </div>
      <div 
        :class="['sb-item', { active: activeTab === 'profile' }]" 
        @click="selectTab('profile')"
        id="nav-profile"
      >
        <i class="ti ti-user"></i>
        <span>Admin Profile</span>
      </div>
      <div 
        :class="['sb-item', { active: activeTab === 'settings' }]" 
        @click="selectTab('settings')"
        id="nav-settings"
      >
        <i class="ti ti-settings"></i>
        <span>Settings</span>
      </div>
    </div>

    <!-- Footer (bottom of sidebar) -->
    <div class="sidebar-footer">
      <div 
        v-if="currentPreset" 
        class="admin-avatar"
        :style="{ backgroundColor: currentPreset.bg, color: currentPreset.color }"
      >
        <i class="ti" :class="currentPreset.icon" style="font-size: 14px !important;"></i>
      </div>
      <div 
        v-else-if="adminProfile.avatar_url && adminProfile.avatar_url.startsWith('http')"
        class="admin-avatar img-avatar"
        :style="{ backgroundImage: `url(${adminProfile.avatar_url})` }"
      >
      </div>
      <div v-else class="admin-avatar">
        {{ adminInitials }}
      </div>

      <div class="admin-details">
        <div class="admin-name">{{ adminName }}</div>
        <div class="admin-role">Administrator</div>
      </div>
      <button class="btn-logout" @click="authStore.logoutUser(); $router.push('/admin-dashboard')" title="Exit Dashboard">
        <i class="ti ti-logout"></i>
      </button>
    </div>
  </aside>
</template>

<style scoped>
.sidebar {
  width: 220px;
  height: 100vh;
  background-color: #1a3a6b; /* var(--primary) */
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  z-index: 10;
  user-select: none;
}

.brand-area {
  padding: 14px 16px;
  display: flex;
  align-items: center;
  gap: 10px;
  border-bottom: 0.5px solid rgba(255, 255, 255, 0.1);
  height: 56px;
  flex-shrink: 0;
  box-sizing: border-box;
}

.logo-box {
  width: 32px;
  height: 32px;
  background-color: #f97316; /* var(--accent) */
  border-radius: 8px; /* var(--radius) */
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 18px;
}

.logo-box i {
  color: #fff !important;
  font-size: 18px !important;
}

.brand-text {
  display: flex;
  flex-direction: column;
}

.app-name {
  color: #fff;
  font-size: 13px;
  font-weight: 500;
  line-height: 1.2;
}

.panel-subtitle {
  color: rgba(255, 255, 255, 0.5);
  font-size: 11px;
  line-height: 1.2;
}

.nav-menu {
  flex-grow: 1;
  padding: 16px 12px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.nav-section-label {
  color: rgba(255, 255, 255, 0.4);
  font-size: 10px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  padding: 8px 10px 4px;
}

.sb-item {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 36px;
  padding: 0 10px;
  border-radius: 6px;
  color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  transition: all 0.15s ease;
  font-weight: 500;
  font-size: 13px;
  box-sizing: border-box;
}

.sb-item i {
  font-size: 16px !important;
  width: 18px;
  text-align: center;
  color: rgba(255, 255, 255, 0.7) !important;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.sb-item:hover {
  background-color: rgba(255, 255, 255, 0.08);
  color: #fff;
}

.sb-item:hover i {
  color: #fff !important;
}

.sb-item.active {
  background-color: rgba(255, 255, 255, 0.15);
  color: #fff;
}

.sb-item.active i {
  color: #f97316 !important; /* var(--accent) */
}

.return-portal-btn {
  border: 0.5px solid rgba(249, 115, 22, 0.2);
  background-color: rgba(249, 115, 22, 0.05);
}

.return-portal-btn i {
  color: #f97316 !important;
}

.return-portal-btn:hover {
  background-color: rgba(249, 115, 22, 0.15) !important;
  color: #fff;
}

.sb-badge {
  color: #fff;
  font-size: 10px;
  font-weight: 500;
  padding: 2px 6px;
  border-radius: 10px;
  margin-left: auto;
  line-height: 1;
}

.bg-accent {
  background-color: #f97316;
}

.sidebar-footer {
  padding: 14px 16px;
  display: flex;
  align-items: center;
  gap: 10px;
  border-top: 0.5px solid rgba(255, 255, 255, 0.1);
  color: #fff;
  height: 58px;
  flex-shrink: 0;
  box-sizing: border-box;
}

.admin-avatar {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 500;
  color: #fff;
  flex-shrink: 0;
  border: 0.5px solid rgba(255, 255, 255, 0.1);
}

.admin-avatar.img-avatar {
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
}

.admin-details {
  flex-grow: 1;
  min-width: 0;
}

.admin-name {
  font-size: 12px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.8);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: 1.2;
}

.admin-role {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
  line-height: 1.2;
}

.btn-logout {
  background: transparent;
  border: none;
  color: rgba(255, 255, 255, 0.4);
  cursor: pointer;
  font-size: 14px;
  padding: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.btn-logout:hover {
  color: #fff;
}

/* ======================================================== */
/* PREMIUM LIGHT SIDEBAR STYLING OVERRIDES                   */
/* ======================================================== */
.sidebar.light {
  background-color: #ffffff;
  border-right: 0.5px solid rgba(0, 0, 0, 0.08);
}

.sidebar.light .brand-area {
  border-bottom: 0.5px solid rgba(0, 0, 0, 0.08);
}

.sidebar.light .app-name {
  color: #0f172a;
}

.sidebar.light .panel-subtitle {
  color: #64748b;
}

.sidebar.light .nav-section-label {
  color: #94a3b8;
}

.sidebar.light .sb-item {
  color: #475569;
}

.sidebar.light .sb-item i {
  color: #64748b !important;
}

.sidebar.light .sb-item:hover {
  background-color: rgba(0, 0, 0, 0.04);
  color: #0f172a;
}

.sidebar.light .sb-item:hover i {
  color: #0f172a !important;
}

.sidebar.light .sb-item.active {
  background-color: rgba(26, 58, 107, 0.06);
  color: #1a3a6b;
}

.sidebar.light .sb-item.active i {
  color: #f97316 !important; /* Keep accent for active icon */
}

.sidebar.light .return-portal-btn {
  border: 0.5px solid rgba(249, 115, 22, 0.3);
  background-color: rgba(249, 115, 22, 0.06);
}

.sidebar.light .return-portal-btn:hover {
  background-color: rgba(249, 115, 22, 0.12) !important;
  color: #ea580c;
}

.sidebar.light .sidebar-footer {
  border-top: 0.5px solid rgba(0, 0, 0, 0.08);
  color: #0f172a;
}

.sidebar.light .admin-avatar {
  background-color: rgba(0, 0, 0, 0.05);
  color: #475569;
  border: 0.5px solid rgba(0, 0, 0, 0.08);
}

.sidebar.light .admin-name {
  color: #0f172a;
}

.sidebar.light .admin-role {
  color: #64748b;
}

.sidebar.light .btn-logout {
  color: #64748b;
}

.sidebar.light .btn-logout:hover {
  color: #0f172a;
}
</style>

