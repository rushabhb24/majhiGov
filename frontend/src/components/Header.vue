<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import ThemeToggle from './ui/ThemeToggle.vue'
import LanguageSwitcher from './ui/LanguageSwitcher.vue'
import RuralModeToggle from './ui/RuralModeToggle.vue'
import AppButton from './ui/AppButton.vue'

const authStore = useAuthStore()
const showNotifications = ref(false)
const mobileMenuOpen = ref(false)

const props = defineProps({
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
})

const emit = defineEmits([
  'update:activeTab',
  'update:currentLanguage',
  'update:ruralMode',
  'update:theme',
  'loginClick',
  'logout'
])

function navigateTo(tab) {
  emit('update:activeTab', tab)
  mobileMenuOpen.value = false
}

function handleLogout() {
  emit('logout')
  mobileMenuOpen.value = false
}
</script>

<template>
  <header class="glass tw-sticky tw-top-0 tw-z-[999] tw-w-full">
    <!-- Tricolor Top Bar -->
    <div class="tricolor-bar tw-h-[4px] tw-w-full"></div>

    <div class="tw-max-w-7xl tw-mx-auto tw-px-4 tw-sm:px-6 tw-lg:px-8">
      <div class="tw-flex tw-items-center tw-justify-between tw-h-16">
        
        <!-- Left: Logo -->
        <div class="tw-flex-shrink-0 tw-flex tw-items-center tw-cursor-pointer" @click="navigateTo('explorer')">
          <div class="tw-w-10 tw-h-10 tw-rounded-xl tw-bg-primary tw-text-white tw-font-heading tw-font-extrabold tw-flex tw-items-center tw-justify-center tw-text-xl tw-shadow-glow">
            M
          </div>
          <div class="tw-ml-3">
            <span class="tw-block tw-font-heading tw-font-extrabold tw-text-lg tw-tracking-tight tw-text-foreground">
              MajhiGov
            </span>
            <span class="tw-block tw-text-[9px] tw-font-bold tw-text-primary tw-tracking-widest tw-uppercase">
              Yojana &amp; Careers
            </span>
          </div>
        </div>

        <!-- Center: Navigation links (Desktop) -->
        <nav class="tw-hidden md:tw-flex tw-space-x-1.5 tw-mx-6">
          <button
            v-for="nav in [
              { id: 'explorer', label: t.explorer || 'Explorer' },
              { id: 'eligibility', label: t.eligibility || 'Eligibility' },
              { id: 'jobs', label: t.govtJobsNav || 'Jobs' },
              { id: 'ai-assistant', label: t.aiAssistant || 'AI' },
              { id: 'career-resources', label: t.careerResources || 'Resources' }
            ]"
            :key="nav.id"
            @click="navigateTo(nav.id)"
            class="tw-px-4 tw-py-2 tw-rounded-full tw-text-xs tw-font-heading tw-font-bold tw-transition-all tw-border-none tw-cursor-pointer"
            :class="activeTab === nav.id ? 'tw-bg-primary tw-text-black tw-shadow-[0_0_14px_rgba(249,115,22,0.6)] tw-font-black' : 'tw-text-slate-600 dark:tw-text-slate-300 hover:tw-text-primary hover:dark:tw-text-primary tw-bg-transparent'"
          >
            {{ nav.label }}
          </button>

          <!-- Saved with badge count -->
          <button
            @click="navigateTo('saved')"
            class="tw-px-4 tw-py-2 tw-rounded-full tw-text-xs tw-font-heading tw-font-bold tw-transition-all tw-border-none tw-cursor-pointer tw-flex tw-items-center tw-gap-1.5"
            :class="activeTab === 'saved' ? 'tw-bg-primary tw-text-black tw-shadow-[0_0_14px_rgba(249,115,22,0.6)] tw-font-black' : 'tw-text-slate-600 dark:tw-text-slate-300 hover:tw-text-primary hover:dark:tw-text-primary tw-bg-transparent'"
          >
            <span>{{ t.saved || 'Saved' }}</span>
            <span 
              class="tw-text-[10px] tw-px-1.5 tw-py-0.5 tw-rounded-full"
              :class="activeTab === 'saved' ? 'tw-bg-black tw-text-primary' : 'tw-bg-primary tw-text-white'"
            >
              {{ savedCount }}
            </span>
          </button>

          <!-- Citizen Applications -->
          <button
            v-if="user"
            @click="navigateTo('applications')"
            class="tw-px-4 tw-py-2 tw-rounded-full tw-text-xs tw-font-heading tw-font-bold tw-transition-all tw-border-none tw-cursor-pointer"
            :class="activeTab === 'applications' ? 'tw-bg-primary tw-text-black tw-shadow-[0_0_14px_rgba(249,115,22,0.6)] tw-font-black' : 'tw-text-slate-600 dark:tw-text-slate-300 hover:tw-text-primary hover:dark:tw-text-primary tw-bg-transparent'"
          >
            {{ t.myApplications || 'Applications' }}
          </button>
        </nav>

        <!-- Right Tools: Switchers, Notifications, Profile -->
        <div class="tw-hidden md:tw-flex tw-items-center tw-gap-3">
          
          <RuralModeToggle />
          <ThemeToggle />
          <LanguageSwitcher />

          <!-- Notifications Bell dropdown -->
          <div v-if="user" class="tw-relative">
            <button 
              class="tw-w-9 tw-h-9 tw-flex tw-items-center tw-justify-center tw-rounded-full tw-bg-muted/80 hover:tw-bg-muted tw-text-muted-foreground hover:tw-text-foreground tw-transition-colors tw-border-none tw-cursor-pointer tw-relative"
              @click="showNotifications = !showNotifications"
              title="Notifications"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path><path d="M13.73 21a2 2 0 0 1-3.46 0"></path></svg>
              <span v-if="authStore.unreadCount > 0" class="tw-absolute tw-top-0.5 tw-right-0.5 tw-flex tw-h-2 tw-w-2">
                <span class="tw-animate-ping tw-absolute tw-inline-flex tw-h-full tw-w-full tw-rounded-full tw-bg-destructive tw-opacity-75"></span>
                <span class="tw-relative tw-inline-flex tw-rounded-full tw-h-2 tw-w-2 tw-bg-destructive"></span>
              </span>
            </button>

            <!-- Dropdown box -->
            <div 
              v-if="showNotifications" 
              class="glass tw-absolute tw-right-0 tw-mt-2 tw-w-80 tw-rounded-2xl tw-p-4 tw-z-[99] tw-max-h-96 tw-overflow-y-auto"
            >
              <div class="tw-flex tw-justify-between tw-items-center tw-mb-3">
                <h4 class="tw-font-heading tw-font-bold tw-text-sm tw-text-foreground tw-m-0">Notifications</h4>
                <button 
                  v-if="authStore.unreadCount > 0" 
                  @click="authStore.markNotificationsRead()" 
                  class="tw-bg-transparent tw-border-none tw-text-primary tw-font-bold tw-text-xs tw-cursor-pointer hover:tw-underline"
                >
                  Mark all read
                </button>
              </div>
              <div v-if="authStore.notifications.length === 0" class="tw-text-center tw-text-muted-foreground tw-text-xs tw-py-6">
                No notifications
              </div>
              <div v-else class="tw-flex tw-flex-col tw-gap-2">
                <div 
                  v-for="(notif, idx) in authStore.notifications" 
                  :key="idx" 
                  class="tw-p-2.5 tw-rounded-xl tw-border tw-border-border"
                  :class="notif.is_read ? 'tw-bg-transparent' : 'tw-bg-primary/5'"
                >
                  <div class="tw-font-semibold tw-text-xs tw-text-foreground">{{ notif.title }}</div>
                  <div class="tw-text-[11px] tw-text-muted-foreground tw-mt-1">{{ notif.message }}</div>
                </div>
              </div>
            </div>
          </div>

          <!-- User Badge / Profile / Login Button -->
          <div v-if="user" class="tw-flex tw-items-center tw-gap-2 tw-bg-muted/80 tw-pl-3 tw-pr-1 tw-py-1 tw-rounded-full">
            <span class="tw-text-xs tw-font-bold tw-text-foreground tw-cursor-pointer" @click="navigateTo('profile')">
              {{ user.full_name.split(' ')[0] }}
            </span>
            <button 
              class="tw-w-7 tw-h-7 tw-rounded-full tw-bg-primary tw-text-white tw-flex tw-items-center tw-justify-center tw-border-none tw-cursor-pointer"
              @click="navigateTo('profile')"
              title="View Profile"
            >
              👤
            </button>
            <button 
              class="tw-w-7 tw-h-7 tw-rounded-full tw-bg-transparent tw-text-muted-foreground hover:tw-text-destructive tw-flex tw-items-center tw-justify-center tw-border-none tw-cursor-pointer tw-transition-colors"
              @click="handleLogout"
              :title="t.logout || 'Logout'"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path><polyline points="16 17 21 12 16 7"></polyline><line x1="21" y1="12" x2="9" y2="12"></line></svg>
            </button>
          </div>

          <AppButton 
            v-else 
            variant="primary" 
            size="sm" 
            @click="emit('loginClick')"
          >
            {{ t.loginRegister || 'Login' }}
          </AppButton>
        </div>

        <!-- Hamburger Icon Button (Mobile) -->
        <button 
          class="md:tw-hidden tw-w-10 tw-h-10 tw-flex tw-items-center tw-justify-center tw-rounded-xl tw-bg-muted/80 hover:tw-bg-muted tw-border-none tw-cursor-pointer tw-text-foreground"
          @click="mobileMenuOpen = !mobileMenuOpen"
          :aria-expanded="mobileMenuOpen"
        >
          <svg v-if="!mobileMenuOpen" xmlns="http://www.w3.org/2000/svg" class="tw-h-6 tw-w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
          <svg v-else xmlns="http://www.w3.org/2000/svg" class="tw-h-6 tw-w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>

      </div>
    </div>

    <!-- Mobile Navigation Drawer -->
    <div v-if="mobileMenuOpen" class="md:tw-hidden tw-bg-card/98 tw-backdrop-blur-md tw-shadow-2xl tw-border-t tw-border-border tw-absolute tw-w-full">
      <nav class="tw-flex tw-flex-col tw-gap-1 tw-p-4">
        <button
          v-for="nav in [
            { id: 'explorer', label: t.explorer || 'Explorer' },
            { id: 'eligibility', label: t.eligibility || 'Eligibility' },
            { id: 'jobs', label: t.govtJobsNav || 'Jobs' },
            { id: 'ai-assistant', label: t.aiAssistant || 'AI' },
            { id: 'career-resources', label: t.careerResources || 'Resources' }
          ]"
          :key="nav.id"
          @click="navigateTo(nav.id)"
          class="tw-w-full tw-text-left tw-px-4 tw-py-2.5 tw-rounded-xl tw-text-sm tw-font-heading tw-font-bold tw-border-none tw-cursor-pointer"
          :class="activeTab === nav.id ? 'tw-bg-primary tw-text-black tw-shadow-[0_0_12px_rgba(249,115,22,0.5)]' : 'tw-text-slate-700 dark:tw-text-slate-200 hover:tw-text-primary hover:dark:tw-text-primary tw-bg-transparent'"
        >
          {{ nav.label }}
        </button>

        <button
          @click="navigateTo('saved')"
          class="tw-w-full tw-text-left tw-px-4 tw-py-2.5 tw-rounded-xl tw-text-sm tw-font-heading tw-font-bold tw-border-none tw-cursor-pointer tw-flex tw-items-center tw-justify-between"
          :class="activeTab === 'saved' ? 'tw-bg-primary tw-text-black tw-shadow-[0_0_12px_rgba(249,115,22,0.5)]' : 'tw-text-slate-700 dark:tw-text-slate-200 hover:tw-text-primary hover:dark:tw-text-primary tw-bg-transparent'"
        >
          <span>{{ t.saved || 'Saved' }}</span>
          <span 
            class="tw-text-[10px] tw-px-2 tw-py-0.5 tw-rounded-full"
            :class="activeTab === 'saved' ? 'tw-bg-black tw-text-primary' : 'tw-bg-primary tw-text-white'"
          >
            {{ savedCount }}
          </span>
        </button>

        <button
          v-if="user"
          @click="navigateTo('applications')"
          class="tw-w-full tw-text-left tw-px-4 tw-py-2.5 tw-rounded-xl tw-text-sm tw-font-heading tw-font-bold tw-border-none tw-cursor-pointer"
          :class="activeTab === 'applications' ? 'tw-bg-primary tw-text-black tw-shadow-[0_0_12px_rgba(249,115,22,0.5)]' : 'tw-text-slate-700 dark:tw-text-slate-200 hover:tw-text-primary hover:dark:tw-text-primary tw-bg-transparent'"
        >
          {{ t.myApplications || 'Applications' }}
        </button>

        <!-- Divider -->
        <div class="tw-h-[1px] tw-bg-border/50 tw-my-2"></div>

        <!-- Toggles & Selects Row -->
        <div class="tw-flex tw-items-center tw-justify-between tw-px-4 tw-py-2">
          <div class="tw-flex tw-items-center tw-gap-3">
            <RuralModeToggle />
            <ThemeToggle />
            <LanguageSwitcher />
          </div>

          <div v-if="user" class="tw-flex tw-items-center tw-gap-2">
            <span class="tw-text-xs tw-font-bold tw-text-foreground" @click="navigateTo('profile')">
              {{ user.full_name.split(' ')[0] }}
            </span>
            <button 
              class="tw-w-8 tw-h-8 tw-rounded-full tw-bg-primary tw-text-white tw-flex tw-items-center tw-justify-center tw-border-none tw-cursor-pointer"
              @click="navigateTo('profile')"
            >
              👤
            </button>
            <button 
              class="tw-w-8 tw-h-8 tw-rounded-full tw-bg-muted tw-text-muted-foreground hover:tw-text-destructive tw-flex tw-items-center tw-justify-center tw-border-none tw-cursor-pointer"
              @click="handleLogout"
            >
              ✕
            </button>
          </div>

          <AppButton 
            v-else 
            variant="primary" 
            size="sm" 
            @click="emit('loginClick')"
          >
            {{ t.loginRegister || 'Login' }}
          </AppButton>
        </div>
      </nav>
    </div>

  </header>
</template>

<style scoped>
/* Glassmorphism overrides */
.glass {
  border-bottom: 1px solid hsl(var(--border) / 0.8);
}
</style>
