import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useSchemeStore } from './schemes.js'
import { bookmarksApi } from '../api/bookmarks.js'

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
      const data = await bookmarksApi.fetchSavedSchemes()
      if (Array.isArray(data)) {
        savedSchemeIds.value = data
        localStorage.setItem('yojana_saved_ids', JSON.stringify(data))
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
      await bookmarksApi.toggleSavedScheme(schemeId)

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
