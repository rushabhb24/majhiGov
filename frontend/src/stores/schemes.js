import { defineStore } from 'pinia'
import { ref } from 'vue'
import { schemesApi } from '../api/schemes.js'

export const useSchemeStore = defineStore('schemes', () => {
  // State
  const schemes = ref([])
  const loading = ref(false)
  const error = ref(null)
  const searchQuery = ref('')
  const selectedCategory = ref('All')
  const sortBy = ref('date_desc')
  const selectedScheme = ref(null)
  const detailModalOpen = ref(false)
  const categories = ['All', 'Farmers', 'Students', 'Women', 'Senior Citizens', 'Business Owners']

  // Actions
  async function fetchSchemes() {
    loading.value = true
    error.value = null
    try {
      const data = await schemesApi.fetchPublicSchemes({
        category: selectedCategory.value,
        search: searchQuery.value,
        sort_by: sortBy.value
      })
      schemes.value = data || []
    } catch (err) {
      console.error(err)
      error.value = 'Could not connect to Go backend.'
      schemes.value = []
    } finally {
      loading.value = false
    }
  }

  async function openDetails(scheme) {
    loading.value = true
    try {
      const data = await schemesApi.fetchSchemeDetails(scheme.id)
      selectedScheme.value = data
      detailModalOpen.value = true
    } catch (err) {
      console.error(err)
      // Fallback to loaded local list properties
      selectedScheme.value = scheme
      detailModalOpen.value = true
    } finally {
      loading.value = false
    }
  }

  function closeDetails() {
    detailModalOpen.value = false
    selectedScheme.value = null
  }

  return {
    schemes,
    loading,
    error,
    searchQuery,
    selectedCategory,
    sortBy,
    selectedScheme,
    detailModalOpen,
    categories,
    fetchSchemes,
    openDetails,
    closeDetails
  }
})
