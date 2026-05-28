<script setup>
import { computed } from 'vue'

const props = defineProps({
  activeTab: {
    type: String,
    required: true
  },
  searchQuery: {
    type: String,
    default: ''
  },
  theme: {
    type: String,
    default: 'dark'
  }
})

const emit = defineEmits(['update:searchQuery', 'action-click'])

const pageTitle = computed(() => {
  const titles = {
    overview: 'Overview',
    schemes: 'Schemes',
    categories: 'Categories',
    users: 'Users',
    eligibility: 'Eligibility Rules',
    notifications: 'Notifications',
    analytics: 'Analytics',
    settings: 'Settings',
    profile: 'Admin Profile'
  }
  return titles[props.activeTab] || 'Dashboard'
})

const actionLabel = computed(() => {
  const actions = {
    overview: 'Add Scheme',
    schemes: 'Add New Scheme',
    categories: 'Add Category',
    users: 'Add Admin',
    eligibility: 'Save Rules',
    notifications: 'Send Notification',
    analytics: 'Export Report',
    settings: 'Save Settings'
  }
  return actions[props.activeTab] || 'Save Changes'
})

const showSearch = computed(() => {
  return props.activeTab === 'schemes' || props.activeTab === 'users'
})

const isSaveIcon = computed(() => {
  return props.activeTab === 'eligibility' || props.activeTab === 'settings'
})
</script>

<template>
  <header class="topbar">
    <div class="page-title">{{ pageTitle }}</div>

    <div class="topbar-actions">
      <!-- Search Bar (conditional) -->
      <div class="search-bar" v-if="showSearch">
        <i class="ti ti-search"></i>
        <input 
          type="text" 
          :value="searchQuery" 
          @input="emit('update:searchQuery', $event.target.value)"
          placeholder="Search schemes, users..." 
        />
      </div>

      <!-- Theme Toggle Button -->
      <button class="theme-toggle-btn" @click="emit('action-click', 'theme')" :title="theme === 'dark' ? 'Switch to Light Mode' : 'Switch to Dark Mode'">
        <i class="ti" :class="theme === 'dark' ? 'ti-sun' : 'ti-moon'"></i>
      </button>

      <!-- Bell Button -->
      <button class="bell-btn" @click="emit('action-click', 'bell')">
        <i class="ti ti-bell"></i>
        <span class="bell-dot"></span>
      </button>

      <!-- Primary Action Button -->
      <button class="top-action-btn" @click="emit('action-click', 'primary')">
        <i class="ti" :class="isSaveIcon ? 'ti-device-floppy' : 'ti-plus'"></i>
        <span>{{ actionLabel }}</span>
      </button>
    </div>
  </header>
</template>

<style scoped>
.topbar {
  height: 56px;
  background-color: var(--bg);
  border-bottom: 0.5px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  flex-shrink: 0;
  box-sizing: border-box;
}

.page-title {
  font-size: 15px;
  font-weight: 500;
  color: var(--text);
}

.topbar-actions {
  display: flex;
  align-items: center;
  gap: 14px;
}

.search-bar {
  width: 200px;
  height: 32px;
  background-color: var(--bg2);
  border: 0.5px solid var(--border);
  border-radius: 6px;
  display: flex;
  align-items: center;
  padding: 6px 10px;
  gap: 8px;
  box-sizing: border-box;
}

.search-bar i {
  color: var(--text2);
  font-size: 15px !important;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.search-bar input {
  background: transparent;
  border: none;
  outline: none;
  font-family: inherit;
  font-size: 13px;
  color: var(--text);
  width: 100%;
  padding: 0;
}

.search-bar input::placeholder {
  color: var(--text2);
}

.theme-toggle-btn {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: 0.5px solid var(--border);
  background: var(--bg);
  color: var(--text2);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  padding: 0;
  box-sizing: border-box;
}

.theme-toggle-btn:hover {
  background-color: var(--bg2);
  color: var(--text);
}

.theme-toggle-btn i {
  font-size: 16px !important;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.bell-btn {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: 0.5px solid var(--border);
  background: var(--bg);
  color: var(--text2);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  position: relative;
  padding: 0;
  box-sizing: border-box;
}

.bell-btn:hover {
  background-color: var(--bg2);
  color: var(--text);
}

.bell-btn i {
  font-size: 16px !important;
}

.bell-dot {
  width: 7px;
  height: 7px;
  background-color: var(--accent);
  border-radius: 50%;
  position: absolute;
  top: 3px;
  right: 3px;
}

.top-action-btn {
  background-color: #1a3a6b; /* var(--primary) */
  color: #ffffff;
  border: none;
  border-radius: 6px; /* var(--radius) */
  height: 34px;
  padding: 7px 14px;
  font-size: 13px;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  font-family: inherit;
  box-sizing: border-box;
}

.top-action-btn:hover {
  opacity: 0.9;
}

.top-action-btn i {
  font-size: 16px !important;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}
</style>
