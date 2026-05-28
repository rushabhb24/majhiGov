<script setup>
import { ref } from 'vue';
import { useAuthStore } from '../stores/auth';

const authStore = useAuthStore();

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
        <div class="logo-text">MajhiGov <span class="accent-text">Portal</span></div>
      </div>

      <!-- Unified Nav + Settings Row (Desktop) -->
      <div class="nav-row nav-desktop">
        <!-- Navigation Links -->
        <nav class="nav-menu">
          <div 
            :class="['nav-link', { active: activeTab === 'explorer' }]" 
            @click="navigateTo('explorer')"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
            <span>{{ t.explorer }}</span>
          </div>
          <div 
            :class="['nav-link', { active: activeTab === 'eligibility' }]" 
            @click="navigateTo('eligibility')"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 11 12 14 22 4"></polyline><path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"></path></svg>
            <span>{{ t.eligibility }}</span>
          </div>
          <div 
            :class="['nav-link', { active: activeTab === 'saved' }]" 
            @click="navigateTo('saved')"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"></path></svg>
            <span>{{ t.saved }}</span>
            <span v-if="savedCount > 0" class="badge">{{ savedCount }}</span>
          </div>
          <div 
            v-if="user"
            :class="['nav-link', { active: activeTab === 'applications' }]" 
            @click="navigateTo('applications')"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line><polyline points="10 9 9 9 8 9"></polyline></svg>
            <span>{{ t.myApplications || 'Applications' }}</span>
          </div>
          <div 
            v-if="user"
            :class="['nav-link', { active: activeTab === 'profile' }]" 
            @click="navigateTo('profile')"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
            <span>{{ t.myProfile || 'Profile' }}</span>
          </div>
          <div 
            v-if="user && authStore.isAdmin"
            :class="['nav-link', { active: activeTab === 'admin' }]" 
            @click="navigateTo('admin')"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="9"></rect><rect x="14" y="3" width="7" height="5"></rect><rect x="14" y="12" width="7" height="9"></rect><rect x="3" y="16" width="7" height="5"></rect></svg>
            <span>Admin</span>
          </div>
        </nav>

        <!-- Divider -->
        <div class="nav-divider"></div>

        <!-- Settings inline with nav -->
        <div class="nav-settings">
          <!-- Language Selector -->
          <select 
            class="form-control select-lang" 
            :value="currentLanguage"
            @change="emit('update:currentLanguage', $event.target.value)"
          >
            <option value="en">EN</option>
            <option value="hi">हि</option>
            <option value="mr">मरा</option>
          </select>

          <!-- Theme Toggle -->
          <button 
            class="btn-icon" 
            @click="emit('update:theme', theme === 'dark' ? 'light' : 'dark')"
            :title="theme === 'dark' ? 'Light Mode' : 'Dark Mode'"
          >
            <svg v-if="theme === 'dark'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="5"></circle><line x1="12" y1="1" x2="12" y2="3"></line><line x1="12" y1="21" x2="12" y2="23"></line><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line><line x1="1" y1="12" x2="3" y2="12"></line><line x1="21" y1="12" x2="23" y2="12"></line><line x1="4.22" y1="19.77" x2="5.64" y2="18.36"></line><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line></svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path></svg>
          </button>

          <!-- Rural Mode Toggle -->
          <button 
            :class="['btn-icon', { active: ruralMode }]" 
            @click="emit('update:ruralMode', !ruralMode)"
            :title="ruralMode ? (t.normalMode || 'Normal Mode') : (t.ruralMode || 'Rural Mode')"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path></svg>
          </button>

          <!-- Login / User Badge -->
          <div v-if="user" class="user-profile-badge">
            <div class="user-avatar" :title="user.full_name">
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
            </div>
            <span class="user-name-text">{{ user.full_name.split(' ')[0] }}</span>
            <button class="btn-logout" @click="handleLogout" :title="t.logout || 'Logout'">
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path><polyline points="16 17 21 12 16 7"></polyline><line x1="21" y1="12" x2="9" y2="12"></line></svg>
            </button>
          </div>
          <button v-else class="btn-login" @click="emit('loginClick')">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"></path><polyline points="10 17 15 12 10 7"></polyline><line x1="15" y1="12" x2="3" y2="12"></line></svg>
            <span>{{ t.loginRegister || 'Login' }}</span>
          </button>
        </div>
      </div>

      <!-- Hamburger Button (mobile only) -->
      <button 
        class="btn-hamburger" 
        @click="mobileMenuOpen = !mobileMenuOpen"
        :class="{ 'is-active': mobileMenuOpen }"
      >
        <span class="hamburger-line"></span>
        <span class="hamburger-line"></span>
        <span class="hamburger-line"></span>
      </button>
    </div>

    <!-- Mobile Slide-down Drawer -->
    <Transition name="slide-drawer">
      <div v-if="mobileMenuOpen" class="mobile-drawer">
        <nav class="mobile-nav">
          <div :class="['mobile-nav-link', { active: activeTab === 'explorer' }]" @click="navigateTo('explorer')">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
            <span>{{ t.explorer }}</span>
          </div>
          <div :class="['mobile-nav-link', { active: activeTab === 'eligibility' }]" @click="navigateTo('eligibility')">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 11 12 14 22 4"></polyline><path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"></path></svg>
            <span>{{ t.eligibility }}</span>
          </div>
          <div :class="['mobile-nav-link', { active: activeTab === 'saved' }]" @click="navigateTo('saved')">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"></path></svg>
            <span>{{ t.saved }}</span>
            <span v-if="savedCount > 0" class="badge" style="margin-left:auto;">{{ savedCount }}</span>
          </div>
          <div v-if="user" :class="['mobile-nav-link', { active: activeTab === 'applications' }]" @click="navigateTo('applications')">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line><polyline points="10 9 9 9 8 9"></polyline></svg>
            <span>{{ t.myApplications || 'Applications' }}</span>
          </div>
          <div v-if="user" :class="['mobile-nav-link', { active: activeTab === 'profile' }]" @click="navigateTo('profile')">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
            <span>{{ t.myProfile || 'Profile' }}</span>
          </div>
          <div v-if="user && authStore.isAdmin" :class="['mobile-nav-link', { active: activeTab === 'admin' }]" @click="navigateTo('admin')">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="9"></rect><rect x="14" y="3" width="7" height="5"></rect><rect x="14" y="12" width="7" height="9"></rect><rect x="3" y="16" width="7" height="5"></rect></svg>
            <span>Admin Panel</span>
          </div>

          <!-- Mobile Settings Section -->
          <div class="mobile-settings-divider"></div>
          <div class="mobile-settings-row">
            <select 
              class="form-control select-lang" 
              :value="currentLanguage"
              @change="emit('update:currentLanguage', $event.target.value)"
              style="flex: 1;"
            >
              <option value="en">English</option>
              <option value="hi">हिंदी (Hindi)</option>
              <option value="mr">मराठी (Marathi)</option>
            </select>
            <button class="btn-icon" @click="emit('update:theme', theme === 'dark' ? 'light' : 'dark')">
              <svg v-if="theme === 'dark'" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="5"></circle><line x1="12" y1="1" x2="12" y2="3"></line><line x1="12" y1="21" x2="12" y2="23"></line><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line><line x1="1" y1="12" x2="3" y2="12"></line><line x1="21" y1="12" x2="23" y2="12"></line><line x1="4.22" y1="19.77" x2="5.64" y2="18.36"></line><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line></svg>
              <svg v-else xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path></svg>
            </button>
            <button :class="['btn-icon', { active: ruralMode }]" @click="emit('update:ruralMode', !ruralMode)">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path></svg>
            </button>
          </div>
          <!-- Mobile login/logout -->
          <div v-if="user" class="mobile-nav-link" @click="handleLogout">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path><polyline points="16 17 21 12 16 7"></polyline><line x1="21" y1="12" x2="9" y2="12"></line></svg>
            <span>{{ t.logout || 'Logout' }} ({{ user.full_name.split(' ')[0] }})</span>
          </div>
          <div v-else class="mobile-nav-link" @click="emit('loginClick'); mobileMenuOpen = false">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"></path><polyline points="10 17 15 12 10 7"></polyline><line x1="15" y1="12" x2="3" y2="12"></line></svg>
            <span>{{ t.loginRegister || 'Login / Register' }}</span>
          </div>
        </nav>
      </div>
    </Transition>
  </header>
</template>

<style scoped>
/* Unified nav + settings row */
.nav-row {
  display: flex;
  align-items: center;
  gap: 0;
  flex: 1 1 auto;
  min-width: 0;
  background: var(--clr-surface-alt);
  border-radius: var(--border-radius-full);
  border: 1px solid var(--clr-border);
  padding: 3px;
}

.nav-menu {
  display: flex;
  gap: 2px;
  flex-wrap: nowrap;
  min-width: 0;
  overflow: hidden;
}

.nav-divider {
  width: 1px;
  height: 24px;
  background: var(--clr-border);
  margin: 0 6px;
  flex-shrink: 0;
}

.nav-settings {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
  padding-right: 2px;
}

/* Icon button (theme, rural) */
.btn-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: transparent;
  border: none;
  color: var(--clr-text-muted);
  cursor: pointer;
  transition: all 0.15s ease;
  flex-shrink: 0;
}

.btn-icon:hover {
  color: var(--clr-primary);
  background: var(--clr-primary-light);
}

.btn-icon.active {
  color: var(--clr-secondary);
  background: var(--clr-secondary-light);
}

/* Compact language selector */
.nav-settings .select-lang {
  padding: 5px 26px 5px 8px !important;
  font-size: 0.78rem;
  font-weight: 700;
  border-radius: var(--border-radius-full) !important;
  max-width: 60px;
  height: 32px;
  background-color: transparent !important;
  border-color: transparent !important;
}

.nav-settings .select-lang:focus {
  border-color: var(--clr-primary) !important;
  background-color: var(--clr-surface) !important;
}

/* Compact user badge inside nav row */
.nav-settings .user-profile-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 2px 4px 2px 8px;
  background: var(--clr-primary-light);
  border: none;
  border-radius: var(--border-radius-full);
}

.nav-settings .user-avatar {
  width: 24px;
  height: 24px;
}

.nav-settings .user-name-text {
  font-size: 0.78rem;
  max-width: 70px;
}

.nav-settings .btn-logout {
  width: 24px;
  height: 24px;
}

/* Compact login button */
.nav-settings .btn-login {
  padding: 5px 12px;
  font-size: 0.78rem;
  gap: 5px;
}

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

.mobile-settings-divider {
  height: 1px;
  background: var(--clr-border);
  margin: 8px 0;
}

.mobile-settings-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 16px;
}

.mobile-settings-row .select-lang {
  padding: 8px 32px 8px 12px !important;
  font-size: 0.85rem;
  border-radius: var(--border-radius-full) !important;
}

.mobile-settings-row .btn-icon {
  width: 40px;
  height: 40px;
  background: var(--clr-surface-alt);
  border: 1px solid var(--clr-border);
  border-radius: 50%;
}

/* Slide drawer animation */
.slide-drawer-enter-active,
.slide-drawer-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  max-height: 600px;
  overflow: hidden;
}
.slide-drawer-enter-from,
.slide-drawer-leave-to {
  max-height: 0;
  opacity: 0;
}

/* --- Mobile: show hamburger, hide unified nav row --- */
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
}

/* Tablet: allow nav row to scroll if needed */
@media (max-width: 1100px) and (min-width: 769px) {
  .nav-menu {
    overflow-x: auto;
    scrollbar-width: none;
  }
  .nav-menu::-webkit-scrollbar {
    display: none;
  }
  .nav-link span {
    font-size: 0.78rem;
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
}
</style>
