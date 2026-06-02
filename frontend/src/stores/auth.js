import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '../api/auth.js'

function parseJwt(token) {
  try {
    const base64Url = token.split('.')[1]
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
    const jsonPayload = decodeURIComponent(window.atob(base64).split('').map(function(c) {
      return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
    }).join(''))
    return JSON.parse(jsonPayload)
  } catch (e) {
    return null
  }
}

export const useAuthStore = defineStore('auth', () => {
  // State
  const token = ref(localStorage.getItem('yojana_auth_token') || null)
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

  // Getters
  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => {
    if (!token.value) return false
    const claims = parseJwt(token.value)
    return claims ? !!claims.is_admin : false
  })

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
      console.error(err)
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
      if (data.success && data.token) {
        const isUserAdmin = data.profile ? !!data.profile.is_admin : false

        if (!isAdminLogin && isUserAdmin) {
          throw new Error('Access Denied: Administrative credentials must log in via the Admin Console.')
        }

        if (isAdminLogin && !isUserAdmin) {
          throw new Error('Access Denied: Administrative privileges required.')
        }

        token.value = data.token
        userProfile.value = data.profile
        localStorage.setItem('yojana_auth_token', data.token)

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
        }

        uiStore.showToast('Welcome back! Logged in successfully.', 'success')
        authModalOpen.value = false
        loginForm.value.password = ''
      } else {
        throw new Error('Authentication response was invalid')
      }
    } catch (err) {
      console.error(err)
      uiStore.showToast(err.message || 'Authentication failed.', 'danger')
      throw err
    } finally {
      authSubmitting.value = false
    }
  }

  async function fetchUserProfile() {
    if (!token.value) return
    try {
      const data = await authApi.fetchProfile()
      if (data.success && data.profile) {
        userProfile.value = data.profile

        // Prefill eligibility
        const { useEligibilityStore } = await import('./eligibility.js')
        const eligibilityStore = useEligibilityStore()
        eligibilityStore.prefillFromProfile(data.profile)

        // Fetch bookmarks and applications
        const { useBookmarkStore } = await import('./bookmarks.js')
        const bookmarkStore = useBookmarkStore()
        bookmarkStore.fetchSavedSchemes()

        const { useApplicationStore } = await import('./applications.js')
        const applicationStore = useApplicationStore()
        applicationStore.fetchApplications()
      }
    } catch (err) {
      console.error('Session restoration failed:', err)
      logoutUser()
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

        // Re-prefill eligibility with updated profile
        const { useEligibilityStore } = await import('./eligibility.js')
        const eligibilityStore = useEligibilityStore()
        eligibilityStore.prefillFromProfile(result.profile)
      }
    } catch (err) {
      console.error(err)
      uiStore.showToast(err.message || 'Failed to update profile', 'danger')
      throw err
    }
  }

  function logoutUser() {
    token.value = null
    userProfile.value = null
    localStorage.removeItem('yojana_auth_token')
    localStorage.removeItem('yojana_saved_ids')

    // Clear dependent stores (import dynamically)
    import('./bookmarks.js').then(({ useBookmarkStore }) => {
      try { useBookmarkStore().clearBookmarks() } catch (e) { /* ok */ }
    })
    import('./applications.js').then(({ useApplicationStore }) => {
      try { useApplicationStore().clearApplications() } catch (e) { /* ok */ }
    })
  }

  function openAuthModal(tab = 'login') {
    authModalOpen.value = true
    authTab.value = tab
  }

  function closeAuthModal() {
    authModalOpen.value = false
  }

  return {
    token,
    userProfile,
    authModalOpen,
    authTab,
    regForm,
    loginForm,
    authSubmitting,
    isLoggedIn,
    isAdmin,
    registerUser,
    loginUser,
    fetchUserProfile,
    updateProfile,
    logoutUser,
    openAuthModal,
    closeAuthModal
  }
})
