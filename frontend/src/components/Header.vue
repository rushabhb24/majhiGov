<script setup>
defineProps({
  activeTab: {
    type: String,
    required: true
  },
  currentLanguage: {
    type: String,
    required: true
  },
  ruralMode: {
    type: Boolean,
    required: true
  },
  theme: {
    type: String,
    required: true
  },
  savedCount: {
    type: Number,
    default: 0
  },
  t: {
    type: Object,
    required: true
  },
  user: {
    type: Object,
    default: null
  }
});

const emit = defineEmits([
  'update:activeTab',
  'update:currentLanguage',
  'update:ruralMode',
  'update:theme',
  'loginClick',
  'logout'
]);
</script>

<template>
  <header class="header">
    <div class="header-container">
      <!-- Branding Logo -->
      <div class="logo" @click="emit('update:activeTab', 'explorer')">
        <div class="logo-icon">M</div>
        <div>MajhiGov <span class="accent-text">Portal</span></div>
      </div>

      <!-- Navigation Tabs -->
      <nav class="nav-menu">
        <div 
          :class="['nav-link', { active: activeTab === 'explorer' }]" 
          @click="emit('update:activeTab', 'explorer')"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
          <span>{{ t.explorer }}</span>
        </div>
        <div 
          :class="['nav-link', { active: activeTab === 'eligibility' }]" 
          @click="emit('update:activeTab', 'eligibility')"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 11 12 14 22 4"></polyline><path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"></path></svg>
          <span>{{ t.eligibility }}</span>
        </div>
        <div 
          :class="['nav-link', { active: activeTab === 'saved' }]" 
          @click="emit('update:activeTab', 'saved')"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"></path></svg>
          <span>{{ t.saved }}</span>
          <span v-if="savedCount > 0" class="badge">{{ savedCount }}</span>
        </div>
      </nav>

      <!-- Settings Controls -->
      <div class="settings-bar">
        <!-- Language Selector -->
        <div class="lang-selector mr-4">
          <select 
            class="form-control select-lang" 
            :value="currentLanguage"
            @change="emit('update:currentLanguage', $event.target.value)"
          >
            <option value="en">English</option>
            <option value="hi">हिंदी (Hindi)</option>
            <option value="mr">मराठी (Marathi)</option>
          </select>
        </div>

        <!-- Theme Toggle -->
        <button 
          class="btn-theme" 
          @click="emit('update:theme', theme === 'dark' ? 'light' : 'dark')"
          :title="theme === 'dark' ? 'Switch to Light Mode' : 'Switch to Dark Mode'"
        >
          <svg v-if="theme === 'dark'" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="5"></circle><line x1="12" y1="1" x2="12" y2="3"></line><line x1="12" y1="21" x2="12" y2="23"></line><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line><line x1="1" y1="12" x2="3" y2="12"></line><line x1="21" y1="12" x2="23" y2="12"></line><line x1="4.22" y1="19.77" x2="5.64" y2="18.36"></line><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line></svg>
          <svg v-else xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path></svg>
        </button>

        <!-- Rural Accessibility Toggle -->
        <button 
          :class="['btn-toggle', { active: ruralMode }]" 
          @click="emit('update:ruralMode', !ruralMode)"
          :title="t.ruralMode"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path></svg>
        </button>

        <!-- User Authentication Control (Option A) -->
        <div v-if="user" class="user-profile-badge">
          <div class="user-avatar" :title="user.full_name">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
          </div>
          <span class="user-name-text">{{ user.full_name.split(' ')[0] }}</span>
          <button class="btn-logout" @click="emit('logout')" :title="t.logout || 'Logout'">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path><polyline points="16 17 21 12 16 7"></polyline><line x1="21" y1="12" x2="9" y2="12"></line></svg>
          </button>
        </div>
        <button v-else class="btn-login" @click="emit('loginClick')">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"></path><polyline points="10 17 15 12 10 7"></polyline><line x1="15" y1="12" x2="3" y2="12"></line></svg>
          <span>{{ t.loginRegister || 'Login' }}</span>
        </button>
      </div>
    </div>
  </header>
</template>
