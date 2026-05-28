<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../../stores/auth'
import { useUiStore } from '../../stores/ui'

const authStore = useAuthStore()
const uiStore = useUiStore()

const profileForm = ref({
  full_name: '',
  email: '',
  phone: '',
  avatar_url: '',
  password: '',
  confirm_password: ''
})

const isSubmitting = ref(false)

// Pre-seeded premium administrative illustration avatars
const avatarPresets = [
  { name: 'Director Blue', url: 'avatar-dir-blue', bg: '#e8eef8', color: '#1a3a6b', icon: 'ti-user-check' },
  { name: 'Admin Amber', url: 'avatar-adm-amber', bg: '#fff4ed', color: '#f97316', icon: 'ti-shield-check' },
  { name: 'Officer Emerald', url: 'avatar-off-emerald', bg: '#f0fdf4', color: '#16a34a', icon: 'ti-medal' },
  { name: 'Manager Rose', url: 'avatar-mgr-rose', bg: '#fef2f2', color: '#dc2626', icon: 'ti-briefcase' },
  { name: 'Executive Violet', url: 'avatar-exe-violet', bg: '#faf5ff', color: '#7c3aed', icon: 'ti-award' },
  { name: 'Superviser Teal', url: 'avatar-sup-teal', bg: '#f0fdfa', color: '#0d9488', icon: 'ti-id' },
  { name: 'Specialist Cyan', url: 'avatar-spe-cyan', bg: '#ecfeff', color: '#0891b2', icon: 'ti-user-plus' },
  { name: 'Coordinator Indigo', url: 'avatar-coo-indigo', bg: '#eef2ff', color: '#4f46e5', icon: 'ti-settings' }
]

onMounted(() => {
  if (authStore.userProfile) {
    const p = authStore.userProfile
    profileForm.value.full_name = p.full_name || ''
    profileForm.value.email = p.email || ''
    profileForm.value.phone = p.phone || ''
    profileForm.value.avatar_url = p.avatar_url || 'avatar-dir-blue'
  }
})

function selectAvatar(url) {
  profileForm.value.avatar_url = url
}

async function saveAdminProfile() {
  if (profileForm.value.password) {
    if (profileForm.value.password !== profileForm.value.confirm_password) {
      uiStore.showToast('Passwords do not match.', 'danger')
      return
    }
    if (profileForm.value.password.length < 6) {
      uiStore.showToast('Password must be at least 6 characters long.', 'danger')
      return
    }
  }

  isSubmitting.value = true
  try {
    const payload = {
      full_name: profileForm.value.full_name,
      email: profileForm.value.email,
      phone: profileForm.value.phone,
      avatar_url: profileForm.value.avatar_url,
      // Pass other required fields from existing profile (matching models.UserProfile constraints)
      date_of_birth: authStore.userProfile?.date_of_birth || '1990-01-01',
      gender: authStore.userProfile?.gender || 'Male',
      state: authStore.userProfile?.state || 'Maharashtra',
      district: authStore.userProfile?.district || 'Mumbai',
      caste_category: authStore.userProfile?.caste_category || 'General',
      annual_income: authStore.userProfile?.annual_income || 0,
      occupation: authStore.userProfile?.occupation || 'Other',
      employee_type: authStore.userProfile?.employee_type || 'Government Employee',
      education_level: authStore.userProfile?.education_level || 'Graduate',
      is_disabled: authStore.userProfile?.is_disabled || false
    }

    if (profileForm.value.password) {
      payload.password = profileForm.value.password
    }

    await authStore.updateProfile(payload)
    profileForm.value.password = ''
    profileForm.value.confirm_password = ''
  } catch (err) {
    console.error('Failed to update admin profile:', err)
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="profile-tab">
    <div class="profile-grid">
      <!-- LEFT SECTION: Avatar Selector Card -->
      <div class="card shadow-sm">
        <div class="card-header">
          <div class="card-title">Administrative Avatar Selection</div>
        </div>
        <div class="card-body">
          <div class="avatar-preview-area">
            <div class="large-avatar-box">
              <template v-if="profileForm.avatar_url && avatarPresets.find(a => a.url === profileForm.avatar_url)">
                <div 
                  class="preview-badge" 
                  :style="{ 
                    backgroundColor: avatarPresets.find(a => a.url === profileForm.avatar_url).bg,
                    color: avatarPresets.find(a => a.url === profileForm.avatar_url).color
                  }"
                >
                  <i class="ti" :class="avatarPresets.find(a => a.url === profileForm.avatar_url).icon"></i>
                </div>
              </template>
              <template v-else>
                <div class="preview-badge def-badge">
                  <i class="ti ti-user"></i>
                </div>
              </template>
            </div>
            <div class="avatar-meta">
              <div class="avatar-name-display">{{ avatarPresets.find(a => a.url === profileForm.avatar_url)?.name || 'Default Avatar' }}</div>
              <div class="avatar-desc-display">Choose an illustration to identify yourself across the administration panel.</div>
            </div>
          </div>

          <hr class="divider mt-3 mb-3" />

          <div class="avatar-presets-grid">
            <div 
              v-for="preset in avatarPresets" 
              :key="preset.url"
              :class="['preset-card', { active: profileForm.avatar_url === preset.url }]"
              @click="selectAvatar(preset.url)"
              :style="{ backgroundColor: preset.bg, color: preset.color }"
              :title="preset.name"
            >
              <i class="ti" :class="preset.icon"></i>
              <div class="check-indicator" v-if="profileForm.avatar_url === preset.url">
                <i class="ti ti-check"></i>
              </div>
            </div>
          </div>

          <div class="form-group mt-3">
            <label class="form-label">Custom Profile Photo URL (Optional)</label>
            <input 
              type="text" 
              class="form-input" 
              v-model="profileForm.avatar_url" 
              placeholder="e.g. https://images.unsplash.com/... or base64" 
            />
          </div>
        </div>
      </div>

      <!-- RIGHT SECTION: Personal Details Card -->
      <div class="card shadow-sm">
        <div class="card-header">
          <div class="card-title">Administrative Credentials & Details</div>
        </div>
        <div class="card-body">
          <form @submit.prevent="saveAdminProfile" class="profile-form">
            <div class="form-group">
              <label class="form-label">Super Admin Full Name *</label>
              <input 
                type="text" 
                class="form-input" 
                v-model="profileForm.full_name" 
                placeholder="Super Admin"
                required 
              />
            </div>

            <div class="form-row">
              <div class="form-group">
                <label class="form-label">Administrative Email *</label>
                <input 
                  type="email" 
                  class="form-input" 
                  v-model="profileForm.email" 
                  placeholder="admin@gov.in"
                  required 
                />
              </div>
              <div class="form-group">
                <label class="form-label">Phone Number *</label>
                <input 
                  type="tel" 
                  class="form-input" 
                  v-model="profileForm.phone" 
                  placeholder="9999999999"
                  required 
                />
              </div>
            </div>

            <hr class="divider mt-3 mb-3" />
            <h4 class="section-sub-title">Update Security Password</h4>
            <p class="section-sub-desc">Leave fields blank if you do not wish to update your administrative password.</p>

            <div class="form-row mt-2">
              <div class="form-group">
                <label class="form-label">New Password</label>
                <input 
                  type="password" 
                  class="form-input" 
                  v-model="profileForm.password" 
                  placeholder="••••••••" 
                />
              </div>
              <div class="form-group">
                <label class="form-label">Confirm New Password</label>
                <input 
                  type="password" 
                  class="form-input" 
                  v-model="profileForm.confirm_password" 
                  placeholder="••••••••" 
                />
              </div>
            </div>

            <div class="actions-row mt-4">
              <button 
                type="submit" 
                class="save-btn" 
                :disabled="isSubmitting"
              >
                <i class="ti" :class="isSubmitting ? 'ti-refresh spin-anim' : 'ti-device-floppy'"></i>
                <span>{{ isSubmitting ? 'Saving Changes...' : 'Save Administrative Profile' }}</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.profile-tab {
  width: 100%;
}

.profile-grid {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 16px;
  align-items: start;
}

.card {
  background-color: var(--bg);
  border: 0.5px solid var(--border);
  border-radius: 12px;
  overflow: hidden;
  box-sizing: border-box;
}

.card-header {
  padding: 14px 16px;
  border-bottom: 0.5px solid var(--border);
}

.card-title {
  font-size: 13px;
  font-weight: 500;
  color: var(--text);
}

.card-body {
  padding: 16px;
}

/* Avatar Preview Section */
.avatar-preview-area {
  display: flex;
  align-items: center;
  gap: 14px;
}

.large-avatar-box {
  width: 60px;
  height: 60px;
  flex-shrink: 0;
}

.preview-badge {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: 500;
  border: 0.5px solid var(--border);
}

.preview-badge.def-badge {
  background-color: var(--bg3);
  color: var(--text2);
}

.avatar-meta {
  flex-grow: 1;
}

.avatar-name-display {
  font-size: 13px;
  font-weight: 500;
  color: var(--text);
}

.avatar-desc-display {
  font-size: 11px;
  color: var(--text2);
  line-height: 1.3;
  margin-top: 2px;
}

/* Avatar presets grid */
.avatar-presets-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
}

.preset-card {
  height: 52px;
  border-radius: 8px;
  border: 0.5px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  cursor: pointer;
  position: relative;
  transition: all 0.15s ease;
}

.preset-card:hover {
  transform: translateY(-2px);
}

.preset-card.active {
  border-color: var(--accent);
  box-shadow: 0 0 0 1.5px var(--accent);
}

.check-indicator {
  position: absolute;
  top: -4px;
  right: -4px;
  width: 16px;
  height: 16px;
  background-color: var(--accent);
  border-radius: 50%;
  color: #fff;
  font-size: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.check-indicator i {
  font-size: 9px !important;
  color: #fff !important;
}

/* Forms */
.profile-form {
  display: flex;
  flex-direction: column;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 12px;
  width: 100%;
  box-sizing: border-box;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.form-label {
  font-size: 13px;
  color: var(--text);
  font-weight: 500;
}

.form-input {
  padding: 8px 10px;
  border: 0.5px solid var(--border);
  border-radius: 6px;
  font-size: 13px;
  background-color: var(--bg);
  color: var(--text);
  outline: none;
  font-family: inherit;
  box-sizing: border-box;
  width: 100%;
}

.form-input:focus {
  border-color: var(--primary);
}

.divider {
  border: none;
  border-top: 0.5px solid var(--border);
}

.mt-2 { margin-top: 8px; }
.mt-3 { margin-top: 12px; }
.mt-4 { margin-top: 16px; }
.mb-3 { margin-bottom: 12px; }

.section-sub-title {
  font-size: 13px;
  font-weight: 500;
  color: var(--text);
  margin: 0;
}

.section-sub-desc {
  font-size: 11px;
  color: var(--text2);
  margin: 2px 0 0 0;
}

.save-btn {
  background-color: var(--primary);
  color: var(--clr-text-light);
  border: none;
  border-radius: 6px;
  padding: 8px 16px;
  font-size: 13px;
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  font-family: inherit;
  box-sizing: border-box;
}

.save-btn:hover {
  opacity: 0.9;
}

.save-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.save-btn i {
  font-size: 16px !important;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.spin-anim {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
