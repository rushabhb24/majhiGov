<script setup>
import { ref } from 'vue';

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

const mobileMenuOpen = ref(false);

function navigateTo(tab) {
  emit('update:activeTab', tab);
  mobileMenuOpen.value = false;
}

function handleLogout() {
  emit('logout');
  mobileMenuOpen.value = false;
}
</script>

<template>
  <header class="header">
    <div class="header-container">
      <!-- Branding Logo -->
      <div class="logo" @click="navigateTo('explorer')">
        <div class="logo-icon">M</div>
        <div>MajhiGov <span class="accent-text">Portal</span></div>
      </div>

      <!-- Desktop Navigation Tabs (hidden on mobile) -->
      <nav class="nav-menu nav-desktop">
        <div 
          :class="['nav-link', { active: activeTab === 'explorer' }]" 
          @click="navigateTo('explorer')"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
          <span>{{ t.explorer }}</span>
        </div>
        <div 
          :class="['nav-link', { active: activeTab === 'eligibility' }]" 
          @click="navigateTo('eligibility')"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 11 12 14 22 4"></polyline><path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"></path></svg>
          <span>{{ t.eligibility }}</span>
        </div>
        <div 
          :class="['nav-link', { active: activeTab === 'saved' }]" 
          @click="navigateTo('saved')"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"></path></svg>
          <span>{{ t.saved }}</span>
          <span v-if="savedCount > 0" class="badge">{{ savedCount }}</span>
        </div>
        <div 
          v-if="user"
          :class="['nav-link', { active: activeTab === 'applications' }]" 
          @click="navigateTo('applications')"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line><polyline points="10 9 9 9 8 9"></polyline></svg>
          <span>{{ t.myApplications || 'My Applications' }}</span>
        </div>
        <div 
          v-if="user"
          :class="['nav-link', { active: activeTab === 'profile' }]" 
          @click="navigateTo('profile')"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
          <span>{{ t.myProfile || 'My Profile' }}</span>
        </div>
      </nav>

      <!-- Settings Controls -->
      <div class="settings-bar">
        <!-- Language Selector -->
        <div class="lang-selector">
          <select 
            class="form-control select-lang" 
            :value="currentLanguage"
            @change="emit('update:currentLanguage', $event.target.value)"
          >
            <option value="en">English</option>
            <option value="hi">हिंदी</option>
            <option value="mr">मराठी</option>
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

        <!-- Rural Accessibility Toggle (hidden on very small screens) -->
        <button 
          :class="['btn-toggle', 'hide-mobile-xs', { active: ruralMode }]" 
          @click="emit('update:ruralMode', !ruralMode)"
          :title="t.ruralMode"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path></svg>
        </button>

        <!-- User Authentication Control -->
        <div v-if="user" class="user-profile-badge">
          <div class="user-avatar" :title="user.full_name">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
          </div>
          <span class="user-name-text">{{ user.full_name.split(' ')[0] }}</span>
          <button class="btn-logout" @click="handleLogout" :title="t.logout || 'Logout'">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path><polyline points="16 17 21 12 16 7"></polyline><line x1="21" y1="12" x2="9" y2="12"></line></svg>
          </button>
        </div>
        <button v-else class="btn-login" @click="emit('loginClick')">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"></path><polyline points="10 17 15 12 10 7"></polyline><line x1="15" y1="12" x2="3" y2="12"></line></svg>
          <span class="hide-mobile-xs">{{ t.loginRegister || 'Login' }}</span>
        </button>

        <!-- Hamburger Menu Button (visible on mobile only) -->
        <button 
          class="btn-hamburger" 
          @click="mobileMenuOpen = !mobileMenuOpen"
          :title="mobileMenuOpen ? 'Close menu' : 'Open menu'"
          :class="{ 'is-active': mobileMenuOpen }"
        >
          <span class="hamburger-line"></span>
          <span class="hamburger-line"></span>
          <span class="hamburger-line"></span>
        </button>
      </div>
    </div>

    <!-- Mobile Slide-down Navigation Drawer -->
    <Transition name="slide-drawer">
      <div v-if="mobileMenuOpen" class="mobile-drawer">
        <nav class="mobile-nav">
          <div 
            :class="['mobile-nav-link', { active: activeTab === 'explorer' }]"
            @click="navigateTo('explorer')"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
            <span>{{ t.explorer }}</span>
          </div>
          <div 
            :class="['mobile-nav-link', { active: activeTab === 'eligibility' }]"
            @click="navigateTo('eligibility')"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 11 12 14 22 4"></polyline><path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"></path></svg>
            <span>{{ t.eligibility }}</span>
          </div>
          <div 
            :class="['mobile-nav-link', { active: activeTab === 'saved' }]"
            @click="navigateTo('saved')"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"></path></svg>
            <span>{{ t.saved }}</span>
            <span v-if="savedCount > 0" class="badge" style="margin-left:auto;">{{ savedCount }}</span>
          </div>
          <div 
            v-if="user"
            :class="['mobile-nav-link', { active: activeTab === 'applications' }]"
            @click="navigateTo('applications')"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line><polyline points="10 9 9 9 8 9"></polyline></svg>
            <span>{{ t.myApplications || 'My Applications' }}</span>
          </div>
          <div 
            v-if="user"
            :class="['mobile-nav-link', { active: activeTab === 'profile' }]"
            @click="navigateTo('profile')"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
            <span>{{ t.myProfile || 'My Profile' }}</span>
          </div>
          <!-- Rural toggle for mobile -->
          <div class="mobile-nav-link" @click="emit('update:ruralMode', !ruralMode); mobileMenuOpen = false">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path></svg>
            <span>{{ ruralMode ? (t.normalMode || 'Normal Mode') : (t.ruralMode || 'Rural Mode') }}</span>
          </div>
        </nav>
      </div>
    </Transition>
  </header>
</template>

<style scoped>
/* Hamburger button — hidden on desktop */
.btn-hamburger {
  display: none;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 4px;
  width: 38px;
  height: 38px;
  padding: 8px;
  background: var(--clr-surface-alt);
  border: 1px solid var(--clr-border);
  border-radius: var(--border-radius-sm);
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-hamburger:hover {
  border-color: var(--clr-primary);
}

.hamburger-line {
  display: block;
  width: 18px;
  height: 2px;
  background-color: var(--clr-text-main);
  border-radius: 1px;
  transition: all 0.3s ease;
}

.btn-hamburger.is-active .hamburger-line:nth-child(1) {
  transform: translateY(6px) rotate(45deg);
}
.btn-hamburger.is-active .hamburger-line:nth-child(2) {
  opacity: 0;
}
.btn-hamburger.is-active .hamburger-line:nth-child(3) {
  transform: translateY(-6px) rotate(-45deg);
}

/* Mobile Drawer */
.mobile-drawer {
  display: none;
}

/* Mobile Nav Links */
.mobile-nav {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 8px 16px 16px;
}

.mobile-nav-link {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  font-family: var(--font-heading);
  font-weight: 600;
  font-size: 0.95rem;
  color: var(--clr-text-muted);
  border-radius: var(--border-radius-md);
  cursor: pointer;
  transition: all 0.2s ease;
}

.mobile-nav-link:hover {
  background: var(--clr-surface-alt);
  color: var(--clr-text-main);
}

.mobile-nav-link.active {
  background: linear-gradient(135deg, var(--clr-primary) 0%, hsl(265, 100%, 60%) 100%);
  color: var(--clr-text-light);
  box-shadow: 0 4px 12px var(--clr-primary-glow);
}

/* Slide drawer animation */
.slide-drawer-enter-active,
.slide-drawer-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  max-height: 500px;
  overflow: hidden;
}
.slide-drawer-enter-from,
.slide-drawer-leave-to {
  max-height: 0;
  opacity: 0;
}

/* --- Mobile breakpoint: show hamburger, hide desktop nav --- */
@media (max-width: 768px) {
  .nav-desktop {
    display: none !important;
  }

  .btn-hamburger {
    display: flex;
  }

  .mobile-drawer {
    display: block;
    background: var(--clr-surface);
    border-top: 1px solid var(--clr-border);
  }

  .hide-mobile-xs {
    display: none;
  }
}

/* Very small screens — compact user badge */
@media (max-width: 420px) {
  .user-name-text {
    display: none;
  }

  .user-profile-badge {
    gap: 4px;
    padding: 4px;
  }

  .select-lang {
    max-width: 72px !important;
    font-size: 0.78rem !important;
    padding: 6px 28px 6px 8px !important;
  }
}
</style>
