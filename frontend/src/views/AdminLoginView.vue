<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useUiStore } from '../stores/ui'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const uiStore = useUiStore()
const router = useRouter()

const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function handleAdminLogin() {
  if (!email.value || !password.value) {
    error.value = 'Please enter both email and password.'
    return
  }

  loading.value = true
  error.value = ''

  try {
    // 1. Authenticate credentials (restricting to administrative logins)
    authStore.loginForm.email = email.value
    authStore.loginForm.password = password.value
    await authStore.loginUser(true)

    // 2. Redirect directly to administrative dashboard
    if (authStore.isLoggedIn && authStore.isAdmin) {
      uiStore.showToast('Welcome back to the Administrative Datastore!', 'success')
      router.push('/admin/dashboard')
    } else {
      error.value = 'Access Denied: Administrative privileges required.'
      authStore.logoutUser()
    }
  } catch (err) {
    console.error(err)
    error.value = err.message || 'Invalid admin credentials'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="admin-login-wrapper">
    <div class="login-card">
      
      <!-- Brand Logo / Identity Header -->
      <div class="brand-header">
        <div class="logo-box">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="9"></rect><rect x="14" y="3" width="7" height="5"></rect><rect x="14" y="12" width="7" height="9"></rect><rect x="3" y="16" width="7" height="5"></rect></svg>
        </div>
        <h1 class="brand-title">MajhiGov Portal</h1>
        <p class="brand-subtitle">Administrative gateway access control</p>
      </div>

      <!-- Error Alert Message -->
      <div class="alert danger-alert" v-if="error">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
        <span>{{ error }}</span>
      </div>

      <!-- Login Form -->
      <form @submit.prevent="handleAdminLogin" class="login-form">
        
        <div class="form-group">
          <label class="form-label" for="admin-email">Administrative email</label>
          <input 
            id="admin-email"
            type="email" 
            class="form-control" 
            v-model="email" 
            placeholder="e.g. admin@gov.in"
            required 
            autocomplete="email"
          />
        </div>

        <div class="form-group mt-3">
          <label class="form-label" for="admin-password">Administrative password</label>
          <input 
            id="admin-password" 
            type="password" 
            class="form-control" 
            v-model="password" 
            placeholder="••••••••"
            required 
            autocomplete="current-password"
          />
        </div>

        <!-- Submit Button -->
        <button type="submit" class="submit-btn" :disabled="loading">
          <span v-if="loading" class="spinner"></span>
          <span>{{ loading ? 'Authenticating credentials...' : 'Authenticate and enter' }}</span>
        </button>

      </form>

      <!-- Portal Footer Links -->
      <div class="login-footer">
        <a href="#/" class="back-link">
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="19" y1="12" x2="5" y2="12"></line><polyline points="12 19 5 12 12 5"></polyline></svg>
          <span>Return to main citizen portal</span>
        </a>
      </div>

    </div>
  </div>
</template>

<style scoped>
.admin-login-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  width: 100vw;
  position: fixed;
  top: 0;
  left: 0;
  background-color: #f1f5f9; /* var(--bg3) in light theme */
  z-index: 9999;
  font-family: 'Plus Jakarta Sans', 'Inter', sans-serif;
}

.login-card {
  background-color: #ffffff;
  border: 0.5px solid rgba(0, 0, 0, 0.08);
  border-radius: 12px;
  width: 100%;
  max-width: 400px;
  padding: 32px;
  box-sizing: border-box;
}

.brand-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  margin-bottom: 24px;
}

.logo-box {
  width: 42px;
  height: 42px;
  background-color: #f97316; /* var(--accent) */
  color: #ffffff;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 12px;
}

.logo-box svg {
  color: #ffffff !important;
}

.brand-title {
  font-family: 'Plus Jakarta Sans', sans-serif;
  font-weight: 500;
  font-size: 19px;
  color: #0f172a;
  margin: 0;
  line-height: 1.2;
}

.brand-subtitle {
  font-size: 12px;
  color: #64748b;
  margin-top: 4px;
  line-height: 1.3;
}

/* Alert styles */
.alert {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  border-radius: 6px;
  font-size: 12px;
  margin-bottom: 18px;
  box-sizing: border-box;
}

.danger-alert {
  background-color: #fef2f2;
  border: 0.5px solid rgba(220, 38, 38, 0.16);
  color: #dc2626;
}

.danger-alert svg {
  color: #dc2626 !important;
  flex-shrink: 0;
}

/* Form inputs */
.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
  width: 100%;
}

.mt-3 {
  margin-top: 14px;
}

.form-label {
  font-size: 12px;
  font-weight: 500;
  color: #334155;
  text-align: left;
}

.form-control {
  width: 100%;
  height: 38px;
  padding: 0 12px;
  border: 0.5px solid rgba(0, 0, 0, 0.08);
  border-radius: 6px;
  font-size: 13px;
  font-family: 'Inter', sans-serif;
  background-color: #ffffff;
  color: #0f172a;
  outline: none;
  box-sizing: border-box;
}

.form-control:focus {
  border-color: #1a3a6b;
  box-shadow: 0 0 0 2px rgba(26, 58, 107, 0.06);
}

.form-control::placeholder {
  color: #94a3b8;
}

/* Submit Button */
.submit-btn {
  width: 100%;
  height: 40px;
  background-color: #1a3a6b; /* var(--primary) */
  color: #ffffff;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-top: 24px;
  font-family: inherit;
  box-sizing: border-box;
  transition: opacity 0.15s ease;
}

.submit-btn:hover {
  opacity: 0.95;
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #ffffff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Footer link */
.login-footer {
  margin-top: 24px;
  text-align: center;
  border-top: 0.5px solid rgba(0, 0, 0, 0.08);
  padding-top: 16px;
}

.back-link {
  font-size: 12px;
  color: #64748b;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  transition: color 0.15s ease;
}

.back-link:hover {
  color: #1a3a6b;
}

.back-link svg {
  color: #64748b !important;
}

.back-link:hover svg {
  color: #1a3a6b !important;
}
</style>
