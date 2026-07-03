import { defineStore } from 'pinia'
import { ref } from 'vue'
import { schemesApi } from '../api/schemes.js'

const PAGE_SIZE = 5 // Infinite scroll loads 5 records at a time

export const useSchemeStore = defineStore('schemes', () => {
  // State
  const schemes = ref([])
  const loading = ref(false)
  const loadingMore = ref(false)
  const error = ref(null)
  const searchQuery = ref('')
  const selectedCategory = ref('All')
  const sortBy = ref('date_desc')
  const selectedScheme = ref(null)
  const detailModalOpen = ref(false)
  const categories = ['All', 'Farmers', 'Students', 'Women', 'Senior Citizens', 'Business Owners']

  // Pagination state
  const currentPage = ref(1)
  const totalSchemes = ref(0)
  const hasNextPage = ref(false)

  // Actions
  /** Initial fetch — resets the list to page 1 */
  async function fetchSchemes() {
    loading.value = true
    error.value = null
    currentPage.value = 1
    schemes.value = []
    try {
      const resp = await schemesApi.fetchPublicSchemes({
        category: selectedCategory.value,
        search: searchQuery.value,
        sort_by: sortBy.value,
        page: 1,
        limit: PAGE_SIZE
      })
      schemes.value = resp.data || []
      totalSchemes.value = resp.meta?.total ?? 0
      hasNextPage.value = resp.meta?.hasNext ?? false
    } catch (err) {
      error.value = 'Could not connect to backend.'
      schemes.value = []
    } finally {
      loading.value = false
    }
  }

  /** Load next page (called by IntersectionObserver at scroll bottom) */
  async function loadMoreSchemes() {
    if (loadingMore.value || !hasNextPage.value) return
    loadingMore.value = true
    try {
      const nextPage = currentPage.value + 1
      const resp = await schemesApi.fetchPublicSchemes({
        category: selectedCategory.value,
        search: searchQuery.value,
        sort_by: sortBy.value,
        page: nextPage,
        limit: PAGE_SIZE
      })
      schemes.value = [...schemes.value, ...(resp.data || [])]
      currentPage.value = nextPage
      hasNextPage.value = resp.meta?.hasNext ?? false
      totalSchemes.value = resp.meta?.total ?? totalSchemes.value
    } catch (err) {
      // fail silently on load-more
    } finally {
      loadingMore.value = false
    }
  }

  async function openDetails(scheme) {
    loading.value = true
    try {
      const data = await schemesApi.fetchSchemeDetails(scheme.id)
      selectedScheme.value = data
      detailModalOpen.value = true
    } catch (err) {
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
    loadingMore,
    error,
    searchQuery,
    selectedCategory,
    sortBy,
    selectedScheme,
    detailModalOpen,
    categories,
    currentPage,
    totalSchemes,
    hasNextPage,
    fetchSchemes,
    loadMoreSchemes,
    openDetails,
    closeDetails
  }
})
