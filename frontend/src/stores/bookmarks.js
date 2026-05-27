import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { API_BASE_URL } from '../config.js'
import { useSchemeStore } from './schemes.js'

export const useBookmarkStore = defineStore('bookmarks', () => {
  // State
  const savedSchemeIds = ref([])

  // Getters
  const bookmarkedSchemes = computed(() => {
    const schemeStore = useSchemeStore()
    return schemeStore.schemes.filter(s => savedSchemeIds.value.includes(s.id))
  })

  // Actions
  async function fetchSavedSchemes() {
    const { useAuthStore } = await import('./auth.js')
    const authStore = useAuthStore()
    if (!authStore.token) return

    try {
      const response = await fetch(`${API_BASE_URL}/api/user/saved`, {
        method: 'GET',
        headers: { 'Authorization': `Bearer ${authStore.token}` }
      })
      if (response.ok) {
        const data = await response.json()
        if (Array.isArray(data)) {
          savedSchemeIds.value = data
          localStorage.setItem('yojana_saved_ids', JSON.stringify(data))
        }
      }
    } catch (err) {
      console.error('Failed to fetch bookmarks from server:', err)
    }
  }

  async function toggleBookmark(schemeId) {
    const { useAuthStore } = await import('./auth.js')
    const authStore = useAuthStore()
    const { useUiStore } = await import('./ui.js')
    const uiStore = useUiStore()

    if (!authStore.token) {
      authStore.openAuthModal('login')
      uiStore.showToast('Please login to save schemes!', 'info')
      return
    }

    try {
      const response = await fetch(`${API_BASE_URL}/api/user/saved`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${authStore.token}`
        },
        body: JSON.stringify({ scheme_id: Number(schemeId) })
      })

      if (!response.ok) {
        throw new Error('Failed to toggle bookmark on server')
      }

      const index = savedSchemeIds.value.indexOf(schemeId)
      if (index === -1) {
        savedSchemeIds.value.push(schemeId)
        uiStore.showToast('Scheme saved to your profile!', 'success')
      } else {
        savedSchemeIds.value.splice(index, 1)
        uiStore.showToast('Scheme removed from your profile.', 'info')
      }
      localStorage.setItem('yojana_saved_ids', JSON.stringify(savedSchemeIds.value))
    } catch (err) {
      console.error('Error toggling bookmark:', err)
      uiStore.showToast('Failed to sync bookmark with account.', 'danger')
    }
  }

  function loadBookmarks() {
    const data = localStorage.getItem('yojana_saved_ids')
    if (data) {
      savedSchemeIds.value = JSON.parse(data)
    }
  }

  function clearBookmarks() {
    savedSchemeIds.value = []
  }

  return {
    savedSchemeIds,
    bookmarkedSchemes,
    fetchSavedSchemes,
    toggleBookmark,
    loadBookmarks,
    clearBookmarks
  }
})
