import { defineStore } from 'pinia'
import { ref } from 'vue'
import { API_BASE_URL } from '../config.js'

export const useEligibilityStore = defineStore('eligibility', () => {
  // State
  const profile = ref({
    age: 25,
    gender: 'Male',
    state: 'Maharashtra',
    caste: 'General',
    annual_income: 180000,
    occupation: 'Unemployed',
    employee_type: 'Unemployed',
    education_level: '12th Pass',
    is_disabled: false
  })
  const checking = ref(false)
  const results = ref(null)
  const checked = ref(false)
  const step = ref(1)

  // Actions
  async function submitEligibility() {
    const { useUiStore } = await import('./ui.js')
    const uiStore = useUiStore()

    checking.value = true
    results.value = null
    try {
      const response = await fetch(`${API_BASE_URL}/api/eligibility-check`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(profile.value)
      })
      if (!response.ok) throw new Error('Eligibility check failed.')
      const data = await response.json()
      results.value = data
      checked.value = true
      uiStore.showToast('Eligibility calculated successfully!', 'success')
    } catch (err) {
      console.error(err)
      uiStore.showToast('Could not connect to Go backend.', 'danger')
    } finally {
      checking.value = false
    }
  }

  async function prefillFromProfile(userProfile) {
    if (!userProfile) return

    // Calculate age from DOB
    let calculatedAge = 25
    if (userProfile.date_of_birth) {
      const dob = new Date(userProfile.date_of_birth)
      const diffMs = Date.now() - dob.getTime()
      const ageDate = new Date(diffMs)
      calculatedAge = Math.abs(ageDate.getUTCFullYear() - 1970)
    }

    profile.value = {
      age: calculatedAge,
      gender: userProfile.gender || 'Male',
      state: userProfile.state || 'Maharashtra',
      caste: userProfile.caste_category || 'General',
      annual_income: userProfile.annual_income || 150000,
      occupation: userProfile.occupation || 'Unemployed',
      employee_type: userProfile.employee_type || 'Unemployed',
      education_level: userProfile.education_level || 'Graduate',
      is_disabled: userProfile.is_disabled || false
    }

    const { useUiStore } = await import('./ui.js')
    const uiStore = useUiStore()
    uiStore.showToast('Account logged in & eligibility profile prefilled!', 'info')
  }

  function reset() {
    profile.value = {
      age: 25,
      gender: 'Male',
      state: 'Maharashtra',
      caste: 'General',
      annual_income: 180000,
      occupation: 'Unemployed',
      employee_type: 'Unemployed',
      education_level: '12th Pass',
      is_disabled: false
    }
    checking.value = false
    results.value = null
    checked.value = false
    step.value = 1
  }

  return {
    profile,
    checking,
    results,
    checked,
    step,
    submitEligibility,
    prefillFromProfile,
    reset
  }
})
