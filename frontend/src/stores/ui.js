import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useUiStore = defineStore('ui', () => {
  // State
  const theme = ref(localStorage.getItem('yojana_theme') || 'dark')
  const ruralMode = ref(localStorage.getItem('yojana_rural') === 'true')
  const currentLanguage = ref(localStorage.getItem('yojana_lang') || 'en')
  const toast = ref({ show: false, message: '', type: 'success' })

  // Actions
  function showToast(message, type = 'success') {
    toast.value = { show: true, message, type }
    setTimeout(() => {
      toast.value.show = false
    }, 3000)
  }

  function toggleTheme() {
    theme.value = theme.value === 'dark' ? 'light' : 'dark'
    localStorage.setItem('yojana_theme', theme.value)
  }

  function setTheme(val) {
    theme.value = val
    localStorage.setItem('yojana_theme', val)
  }

  function toggleRuralMode() {
    ruralMode.value = !ruralMode.value
    localStorage.setItem('yojana_rural', ruralMode.value.toString())
  }

  function setRuralMode(val) {
    ruralMode.value = val
    localStorage.setItem('yojana_rural', val.toString())
  }

  function setLanguage(lang) {
    currentLanguage.value = lang
    localStorage.setItem('yojana_lang', lang)
  }

  return {
    theme,
    ruralMode,
    currentLanguage,
    toast,
    showToast,
    toggleTheme,
    setTheme,
    toggleRuralMode,
    setRuralMode,
    setLanguage
  }
})
