import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '../api/auth.js'

export const useAuthStore = defineStore('auth', () => {
  // State — no token stored in JS memory (auth is managed via httpOnly cookie on the server)
  const userProfile = ref(null)
  const authModalOpen = ref(false)
  const authTab = ref('login')
  const regForm = ref({
    email: '',
    phone: '',
    password: '',
    full_name: '',
    date_of_birth: '1995-01-01',
    gender: 'Male',
    state: 'Maharashtra',
    district: 'Pune',
    caste_category: 'General',
    annual_income: 150000,
    occupation: 'Unemployed',
    employee_type: 'Unemployed',
    education_level: 'Graduate',
    is_disabled: false,
    aadhaar: ''
  })
  const loginForm = ref({ email: '', password: '' })
  const authSubmitting = ref(false)

  // Getters — derive auth state from profile (not from a stored token)
  const isLoggedIn = computed(() => !!userProfile.value)
  const isAdmin = computed(() => !!userProfile.value?.is_admin)

  // Actions
  async function registerUser() {
    const { useUiStore } = await import('./ui.js')
    const uiStore = useUiStore()

    authSubmitting.value = true
    try {
      const payload = {
        email: regForm.value.email,
        phone: regForm.value.phone,
        password: regForm.value.password,
        full_name: regForm.value.full_name,
        date_of_birth: regForm.value.date_of_birth,
        gender: regForm.value.gender,
        state: regForm.value.state,
        district: regForm.value.district,
        caste_category: regForm.value.caste_category,
        annual_income: Number(regForm.value.annual_income),
        occupation: regForm.value.occupation,
        employee_type: regForm.value.employee_type,
        education_level: regForm.value.education_level,
        is_disabled: regForm.value.is_disabled,
        aadhaar: regForm.value.aadhaar
      }

      await authApi.register(payload)
      uiStore.showToast('Registration successful! You can now log in.', 'success')
      authTab.value = 'login'
      loginForm.value.email = regForm.value.email
    } catch (err) {
      uiStore.showToast(err.message || 'Authentication failed.', 'danger')
    } finally {
      authSubmitting.value = false
    }
  }

  async function loginUser(isAdminLogin = false) {
    const { useUiStore } = await import('./ui.js')
    const uiStore = useUiStore()

    authSubmitting.value = true
    try {
      const data = await authApi.login(loginForm.value)
      if (data.success && data.profile) {
        const isUserAdmin = !!data.profile.is_admin

        if (!isAdminLogin && isUserAdmin) {
          throw new Error('Access Denied: Administrative credentials must log in via the Admin Console.')
        }

        if (isAdminLogin && !isUserAdmin) {
          throw new Error('Access Denied: Administrative privileges required.')
        }

        userProfile.value = data.profile

        // Prefill eligibility from profile
        const { useEligibilityStore } = await import('./eligibility.js')
        const eligibilityStore = useEligibilityStore()
        eligibilityStore.prefillFromProfile(data.profile)

        // Sync bookmarks and applications (only for regular citizens)
        if (!isUserAdmin) {
          const { useBookmarkStore } = await import('./bookmarks.js')
          const bookmarkStore = useBookmarkStore()
          await bookmarkStore.fetchSavedSchemes()

          const { useApplicationStore } = await import('./applications.js')
          const applicationStore = useApplicationStore()
          await applicationStore.fetchApplications()

          // Connect WebSocket for real-time notifications
          connectNotifications()
        }

        uiStore.showToast('Welcome back! Logged in successfully.', 'success')
        authModalOpen.value = false
        loginForm.value.password = ''
      } else {
        throw new Error('Authentication response was invalid')
      }
    } catch (err) {
      uiStore.showToast(err.message || 'Authentication failed.', 'danger')
      throw err
    } finally {
      authSubmitting.value = false
    }
  }

  async function fetchUserProfile() {
    try {
      const data = await authApi.fetchProfile()
      if (data.success && data.profile) {
        userProfile.value = data.profile

        const { useEligibilityStore } = await import('./eligibility.js')
        const eligibilityStore = useEligibilityStore()
        eligibilityStore.prefillFromProfile(data.profile)

        if (!data.profile.is_admin) {
          const { useBookmarkStore } = await import('./bookmarks.js')
          const bookmarkStore = useBookmarkStore()
          bookmarkStore.fetchSavedSchemes()

          const { useApplicationStore } = await import('./applications.js')
          const applicationStore = useApplicationStore()
          applicationStore.fetchApplications()

          connectNotifications()
        }
      }
    } catch (err) {
      // Profile fetch failed — user is not authenticated
      userProfile.value = null
    }
  }

  async function updateProfile(data) {
    const { useUiStore } = await import('./ui.js')
    const uiStore = useUiStore()

    try {
      const result = await authApi.updateProfile(data)
      if (result.success && result.profile) {
        userProfile.value = result.profile
        uiStore.showToast('Profile updated successfully!', 'success')

        const { useEligibilityStore } = await import('./eligibility.js')
        const eligibilityStore = useEligibilityStore()
        eligibilityStore.prefillFromProfile(result.profile)
      }
    } catch (err) {
      uiStore.showToast(err.message || 'Failed to update profile', 'danger')
      throw err
    }
  }

  async function logoutUser() {
    try {
      // Tell the server to clear the httpOnly cookie
      await authApi.logout()
    } catch (e) {
      // Still clear local state even if server call fails
    }
    userProfile.value = null
    disconnectNotifications()

    // Clear all legacy localStorage leftovers
    localStorage.removeItem('yojana_auth_token')
    localStorage.removeItem('yojana_saved_ids')

    import('./bookmarks.js').then(({ useBookmarkStore }) => {
      try { useBookmarkStore().clearBookmarks() } catch (e) { /* ok */ }
    })
    import('./applications.js').then(({ useApplicationStore }) => {
      try { useApplicationStore().clearApplications() } catch (e) { /* ok */ }
    })
  }

  // ─── WebSocket / Notifications ────────────────────────────────────────────

  const notifications = ref([])
  const unreadCount = computed(() => notifications.value.filter(n => !n.is_read).length)
  let wsConnection = null

  function connectNotifications() {
    if (wsConnection) return // Already connected
    const wsUrl = `${import.meta.env.VITE_WS_URL || 'ws://localhost:8080'}/ws`
    try {
      wsConnection = new WebSocket(wsUrl)
      wsConnection.onmessage = (event) => {
        try {
          const msg = JSON.parse(event.data)
          notifications.value.unshift({ ...msg, is_read: false })
        } catch (e) { /* ignore parse errors */ }
      }
      wsConnection.onclose = () => {
        wsConnection = null
        // Auto-reconnect after 5 seconds if user is still logged in
        if (userProfile.value) {
          setTimeout(connectNotifications, 5000)
        }
      }
      wsConnection.onerror = () => {
        wsConnection = null
      }
    } catch (e) {
      // WebSocket not available (e.g., HTTP-only dev mode)
    }
  }

  function disconnectNotifications() {
    if (wsConnection) {
      wsConnection.close()
      wsConnection = null
    }
    notifications.value = []
  }

  async function fetchNotifications() {
    if (!isLoggedIn.value) return
    try {
      const data = await authApi.fetchNotifications()
      if (Array.isArray(data)) {
        notifications.value = data
      }
    } catch (e) { /* fail silently */ }
  }

  async function markNotificationsRead() {
    try {
      await authApi.markNotificationsRead()
      notifications.value = notifications.value.map(n => ({ ...n, is_read: true }))
    } catch (e) { /* fail silently */ }
  }

  function openAuthModal(tab = 'login') {
    authModalOpen.value = true
    authTab.value = tab
  }

  function closeAuthModal() {
    authModalOpen.value = false
  }

  return {
    userProfile,
    authModalOpen,
    authTab,
    regForm,
    loginForm,
    authSubmitting,
    isLoggedIn,
    isAdmin,
    notifications,
    unreadCount,
    registerUser,
    loginUser,
    fetchUserProfile,
    updateProfile,
    logoutUser,
    fetchNotifications,
    markNotificationsRead,
    connectNotifications,
    openAuthModal,
    closeAuthModal
  }
})
