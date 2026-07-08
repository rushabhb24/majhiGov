<script setup>
import { ref, onMounted, computed } from 'vue'
import { API_BASE_URL } from '../config.js'

const companies = ref([])
const loading = ref(true)
const searchQ = ref('')
const selectedIndustry = ref('All')

onMounted(async () => {
  try {
    const resp = await fetch(`${API_BASE_URL}/api/companies`, { credentials: 'include' })
    if (resp.ok) {
      companies.value = await resp.json()
    }
  } catch (err) {
    console.error("Failed to load companies:", err)
  } finally {
    loading.value = false
  }
})

const industries = computed(() => {
  const list = new Set(['All'])
  companies.value.forEach(c => {
    if (c.industry) list.add(c.industry)
  })
  return Array.from(list)
})

const filteredCompanies = computed(() => {
  return companies.value.filter(c => {
    const matchesSearch = c.name.toLowerCase().includes(searchQ.value.toLowerCase()) ||
                          c.description.toLowerCase().includes(searchQ.value.toLowerCase()) ||
                          c.location.toLowerCase().includes(searchQ.value.toLowerCase())
    const matchesIndustry = selectedIndustry.value === 'All' || c.industry === selectedIndustry.value
    return matchesSearch && matchesIndustry
  })
})

function getInitials(name) {
  if (!name) return 'C'
  const parts = name.split(' ')
  if (parts.length > 1) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return name.substring(0, 2).toUpperCase()
}

// Generate unique HSL colors for company initial avatars
function getAvatarColor(name) {
  let hash = 0
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash)
  }
  const h = Math.abs(hash % 360)
  return `hsl(${h}, 70%, 45%)`
}
</script>

<template>
  <div class="companies-view">
    <!-- Hero Header -->
    <header class="companies-hero">
      <div class="blob blob-1"></div>
      <div class="blob blob-2"></div>
      <div class="hero-content">
        <h1 class="hero-title">Top hiring <span class="gradient-text">Companies</span></h1>
        <p class="hero-subtitle">Explore verified organizations listing private jobs, internships, and walk-in drives on MajhiGov.</p>
        <div class="search-box">
          <svg class="search-icon" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
          <input v-model="searchQ" type="text" placeholder="Search companies by name or location..." class="search-input" />
        </div>
      </div>
    </header>

    <!-- Industry Filter Bar -->
    <div class="filters-bar" v-if="industries.length > 1">
      <div 
        v-for="ind in industries" 
        :key="ind"
        :class="['filter-pill', { active: selectedIndustry === ind }]"
        @click="selectedIndustry = ind"
      >
        {{ ind }}
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="companies-grid mt-6">
      <div v-for="i in 6" :key="i" class="company-skeleton-card">
        <div class="sk-header">
          <div class="sk-avatar"></div>
          <div class="sk-meta">
            <div class="sk-line w-24"></div>
            <div class="sk-line w-16"></div>
          </div>
        </div>
        <div class="sk-desc"></div>
        <div class="sk-footer"></div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="filteredCompanies.length === 0" class="empty-state-wrapper mt-6">
      <h3>No companies found 🔍</h3>
      <p>Try modifying your search query or industry filters.</p>
    </div>

    <!-- Company Grid -->
    <div v-else class="companies-grid mt-6">
      <div v-for="comp in filteredCompanies" :key="comp.id" class="company-card glass-card">
        <div class="card-header">
          <div v-if="comp.logo_url" class="company-logo-img">
            <img :src="comp.logo_url" :alt="comp.name" />
          </div>
          <div v-else class="company-logo-initials" :style="{ backgroundColor: getAvatarColor(comp.name) }">
            {{ getInitials(comp.name) }}
          </div>

          <div class="header-details">
            <h3 class="company-name">
              {{ comp.name }}
              <span v-if="comp.is_verified" class="verify-badge" title="Verified Employer">✓</span>
            </h3>
            <span class="industry-tag">{{ comp.industry }}</span>
          </div>
        </div>

        <p class="company-desc">{{ comp.description }}</p>

        <div class="meta-details">
          <div class="meta-item">
            <span>📍 {{ comp.location }}</span>
          </div>
          <div class="meta-item">
            <span>👥 {{ comp.company_size }}</span>
          </div>
        </div>

        <div class="card-footer">
          <span class="job-badge" v-if="comp.job_count > 0">{{ comp.job_count }} Active Jobs</span>
          <span class="job-badge empty" v-else>No Active Jobs</span>
          
          <button class="view-btn btn btn-primary btn-sm" @click="$router.push(`/jobs?company=${comp.name}`)">
            View Jobs
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.companies-view {
  max-width: var(--max-width);
  margin: 0 auto;
  padding: 0 24px 48px;
  box-sizing: border-box;
}

.companies-hero {
  position: relative;
  overflow: hidden;
  border-radius: var(--border-radius-lg);
  background: linear-gradient(135deg, #1e293b 0%, #0f172a 100%);
  padding: 64px 32px;
  margin-top: 24px;
  margin-bottom: 32px;
  text-align: center;
  border: 1px solid var(--clr-border);
}

.blob {
  position: absolute;
  border-radius: 50%;
  pointer-events: none;
  filter: blur(60px);
  opacity: 0.15;
}
.blob-1 {
  width: 300px; height: 300px;
  background: #6366f1;
  top: -80px; left: -50px;
}
.blob-2 {
  width: 250px; height: 250px;
  background: #06b6d4;
  bottom: -40px; right: -50px;
}

.hero-content {
  position: relative;
  z-index: 1;
}

.hero-title {
  font-family: var(--font-heading);
  font-size: 2.75rem;
  font-weight: 800;
  color: white;
  margin-bottom: 12px;
}

.hero-subtitle {
  color: #94a3b8;
  max-width: 600px;
  margin: 0 auto 24px;
  line-height: 1.6;
  font-size: 1rem;
}

.search-box {
  max-width: 500px;
  margin: 0 auto;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 9999px;
  display: flex;
  align-items: center;
  padding: 4px 18px;
  height: 48px;
}

.search-icon {
  color: #94a3b8;
  margin-right: 12px;
  flex-shrink: 0;
}

.search-input {
  border: none;
  outline: none;
  background: transparent;
  width: 100%;
  color: white;
  font-size: 0.95rem;
}

.search-input::placeholder {
  color: #64748b;
}

/* Industry Filters */
.filters-bar {
  display: flex;
  gap: 12px;
  overflow-x: auto;
  padding-bottom: 8px;
  margin-bottom: 32px;
  scrollbar-width: none;
}

.filters-bar::-webkit-scrollbar {
  display: none;
}

.filter-pill {
  flex-shrink: 0;
  padding: 8px 18px;
  border-radius: 9999px;
  background: var(--clr-surface-alt);
  border: 1px solid var(--clr-border);
  color: var(--clr-text-main);
  font-weight: 600;
  font-size: 0.82rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.filter-pill:hover {
  background: var(--clr-surface);
  border-color: var(--clr-primary);
}

.filter-pill.active {
  background: var(--clr-primary);
  color: white;
  border-color: var(--clr-primary);
}

/* Companies Grid */
.companies-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 24px;
}

.company-card {
  padding: 24px;
  display: flex;
  flex-direction: column;
  height: 100%;
  box-sizing: border-box;
}

.card-header {
  display: flex;
  gap: 16px;
  align-items: center;
  margin-bottom: 16px;
}

.company-logo-img {
  width: 54px;
  height: 54px;
  border-radius: 12px;
  overflow: hidden;
  background: white;
  border: 1px solid var(--clr-border);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.company-logo-img img {
  max-width: 80%;
  max-height: 80%;
  object-fit: contain;
}

.company-logo-initials {
  width: 54px;
  height: 54px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-family: var(--font-heading);
  font-weight: 800;
  font-size: 1.2rem;
  flex-shrink: 0;
}

.header-details {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.company-name {
  font-family: var(--font-heading);
  font-size: 1.15rem;
  font-weight: 700;
  color: var(--clr-text-main);
  margin: 0;
  display: flex;
  align-items: center;
  gap: 6px;
}

.verify-badge {
  background: #3b82f6;
  color: white;
  font-size: 0.65rem;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: 800;
}

.industry-tag {
  font-size: 0.78rem;
  color: var(--clr-text-muted);
}

.company-desc {
  font-size: 0.88rem;
  color: var(--clr-text-muted);
  line-height: 1.6;
  margin: 0 0 20px;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  height: 4.8em;
}

.meta-details {
  display: flex;
  gap: 16px;
  margin-top: auto;
  padding-top: 12px;
  border-top: 1px solid var(--clr-border);
  margin-bottom: 18px;
}

.meta-item {
  font-size: 0.8rem;
  color: var(--clr-text-muted);
  font-weight: 500;
}

.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.job-badge {
  font-size: 0.75rem;
  padding: 4px 10px;
  background: var(--clr-success-light);
  color: var(--clr-success);
  border-radius: 9999px;
  font-weight: 700;
}

.job-badge.empty {
  background: var(--clr-surface-alt);
  color: var(--clr-text-muted);
}

.view-btn {
  font-size: 0.82rem;
  padding: 8px 16px;
}

/* Skeletons */
.company-skeleton-card {
  background: var(--clr-surface);
  border: 1px solid var(--clr-border);
  border-radius: var(--border-radius-md);
  padding: 24px;
  height: 250px;
  display: flex;
  flex-direction: column;
}

.sk-header {
  display: flex;
  gap: 16px;
  align-items: center;
  margin-bottom: 20px;
}

.sk-avatar {
  width: 54px;
  height: 54px;
  border-radius: 12px;
  background: var(--clr-surface-alt);
  animation: shimmer 1.5s infinite;
}

.sk-meta {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
}

.sk-line {
  height: 12px;
  background: var(--clr-surface-alt);
  border-radius: 4px;
  animation: shimmer 1.5s infinite;
}

.w-24 { width: 120px; }
.w-16 { width: 80px; }

.sk-desc {
  height: 48px;
  background: var(--clr-surface-alt);
  border-radius: 8px;
  margin-bottom: 20px;
  animation: shimmer 1.5s infinite;
}

.sk-footer {
  height: 32px;
  background: var(--clr-surface-alt);
  border-radius: 8px;
  margin-top: auto;
  animation: shimmer 1.5s infinite;
}

@keyframes shimmer {
  0% { opacity: 0.6; }
  50% { opacity: 1; }
  100% { opacity: 0.6; }
}

.empty-state-wrapper {
  text-align: center;
  padding: 48px;
  background: var(--clr-surface-alt);
  border-radius: var(--border-radius-md);
  color: var(--clr-text-muted);
}
</style>
