<script setup>
import { ref } from 'vue'
import AppCard from '../components/ui/AppCard.vue'
import AppBadge from '../components/ui/AppBadge.vue'
import AppInput from '../components/ui/AppInput.vue'

const searchQuery = ref('')
const selectedCategory = ref('All')

const categories = [
  { name: 'All', icon: '🎯' },
  { name: 'Resume Writing', icon: '📄' },
  { name: 'Interview Tips', icon: '🎤' },
  { name: 'Career Roadmaps', icon: '🗺️' },
  { name: 'Certifications', icon: '📜' },
  { name: 'Online Courses', icon: '💻' },
  { name: 'Salary Guide', icon: '📊' }
]

const articles = [
  {
    title: 'How to Crack UPSC Interview in 2026',
    desc: 'UPSC interview guidelines, structured answering patterns, and mental preparation keys to scoring 180+ marks.',
    category: 'Interview Tips',
    readTime: '8 min read',
    date: 'July 10, 2026'
  },
  {
    title: 'Top 10 Government Jobs for Engineering Graduates',
    desc: 'An exhaustive analysis of central and state level engineering vacancies in ISRO, DRDO, PSUs, and state bodies.',
    category: 'Career Roadmaps',
    readTime: '6 min read',
    date: 'July 08, 2026'
  },
  {
    title: 'Complete Guide to Resume Writing for Freshers',
    desc: 'Step-by-step drafting rules, formatting layout checklists, and critical ATS-optimized keywords for starting careers.',
    category: 'Resume Writing',
    readTime: '5 min read',
    date: 'July 05, 2026'
  },
  {
    title: 'Maharashtrian Government Job Schemes You Must Know',
    desc: 'Overview of MahaDBT, MPSC opportunities, and state-backed skill development programs available for youth.',
    category: 'Career Roadmaps',
    readTime: '7 min read',
    date: 'June 28, 2026'
  },
  {
    title: 'How to Prepare for SSC CGL Examination',
    desc: 'An actionable study timetable, topic weightage guide, and recommendation checklist for preparation materials.',
    category: 'Online Courses',
    readTime: '10 min read',
    date: 'June 25, 2026'
  },
  {
    title: 'Top 5 Free Online Certification Courses for Job Seekers',
    desc: 'Curated list of premium accredited courses from Google, AWS, Microsoft, and top educational foundations.',
    category: 'Certifications',
    readTime: '4 min read',
    date: 'June 18, 2026'
  }
]

const certifications = [
  { name: 'AWS Cloud Practitioner', provider: 'Amazon Web Services', duration: '20 hours', cost: 'Free learning' },
  { name: 'Google Data Analytics Certificate', provider: 'Coursera / Google', duration: '6 months', cost: 'Financial Aid available' },
  { name: 'Project Management Professional (PMP)', provider: 'PMI', duration: '35 hours training', cost: 'Paid exam' },
  { name: 'GATE (Graduate Aptitude Test in Engineering)', provider: 'IITs / IISc', duration: '1 year prep recommended', cost: 'Standard registration' },
  { name: 'UPSC Civil Services Foundation', provider: 'Standard Academy', duration: '1-2 years prep', cost: 'Free standard notes' }
]

const salaryGuide = [
  { sector: 'IT/Software Development', entry: '₹3.5L - ₹8L', mid: '₹8L - ₹18L', senior: '₹18L - ₹35L' },
  { sector: 'Civil Services / Group A Govt', entry: '₹7.2L - ₹10L', mid: '₹12L - ₹18L', senior: '₹20L - ₹30L' },
  { sector: 'Banking & Financial Services', entry: '₹3.0L - ₹6.5L', mid: '₹7L - ₹12L', senior: '₹12L - ₹22L' },
  { sector: 'Secondary / High School Teaching', entry: '₹2.8L - ₹5L', mid: '₹5L - ₹8L', senior: '₹8L - ₹14L' },
  { sector: 'Healthcare & Nursing', entry: '₹2.5L - ₹4.5L', mid: '₹5L - ₹9L', senior: '₹9L - ₹15L' }
]

function filteredArticles() {
  return articles.filter(a => {
    const matchesSearch = a.title.toLowerCase().includes(searchQuery.value.toLowerCase()) || 
                          a.desc.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesCat = selectedCategory.value === 'All' || a.category === selectedCategory.value
    return matchesSearch && matchesCat
  })
}
</script>

<template>
  <div class="tw-max-w-7xl tw-mx-auto tw-px-4 tw-sm:px-6 tw-lg:px-8 tw-py-8">
    
    <!-- Hero Header -->
    <div class="glass tw-p-6 tw-rounded-2xl tw-mb-6">
      <h1 class="tw-font-heading tw-font-bold tw-text-2xl tw-text-foreground tw-m-0">
        Career Resources Hub
      </h1>
      <p class="tw-text-xs tw-text-muted-foreground tw-mt-2 tw-m-0 tw-line-height-[1.5]">
        Equip yourself with professional templates, study material guidelines, and salary indicators to accelerate your career growth.
      </p>

      <!-- Search Box -->
      <div class="glass tw-flex tw-items-center tw-px-4 tw-py-2 tw-rounded-full tw-mt-4 tw-w-full sm:tw-max-w-md">
        <svg class="tw-flex-shrink-0 tw-mr-3" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="Search guides, articles, templates..." 
          class="tw-border-none tw-bg-transparent tw-outline-none tw-text-sm tw-text-foreground tw-w-full" 
        />
      </div>
    </div>

    <!-- Categories Strip -->
    <div class="tw-flex tw-gap-2 tw-overflow-x-auto tw-pb-3 tw-mb-8 tw-scrollbar-none">
      <button 
        v-for="cat in categories" 
        :key="cat.name" 
        class="tw-border-none tw-outline-none tw-px-4 tw-py-2 tw-rounded-full tw-text-xs tw-font-bold tw-font-heading tw-cursor-pointer tw-flex tw-items-center tw-gap-1.5 tw-transition-colors"
        :class="selectedCategory === cat.name ? 'tw-bg-primary tw-text-white' : 'tw-bg-muted/80 tw-text-muted-foreground hover:tw-text-foreground'"
        @click="selectedCategory = cat.name"
      >
        <span>{{ cat.icon }}</span>
        <span>{{ cat.name }}</span>
      </button>
    </div>

    <!-- Featured Articles -->
    <section class="tw-mb-8">
      <div class="tw-mb-6">
        <h2 class="tw-font-heading tw-font-bold tw-text-lg tw-text-foreground tw-m-0">Featured Resources &amp; Guidebooks</h2>
        <p class="tw-text-xs tw-text-muted-foreground tw-mt-1 tw-m-0">Handpicked strategies and insights written by seasoned industry leaders and career developers.</p>
      </div>

      <div v-if="filteredArticles().length === 0" class="glass tw-p-8 tw-rounded-2xl tw-text-center tw-text-xs tw-text-muted-foreground">
        No articles found matching the current search parameters.
      </div>

      <div v-else class="tw-grid tw-grid-cols-1 md:tw-grid-cols-3 tw-gap-6">
        <AppCard v-for="art in filteredArticles()" :key="art.title" hoverable class="tw-flex tw-flex-col tw-h-full">
          <div>
            <AppBadge tone="info">{{ art.category }}</AppBadge>
          </div>
          <h3 class="tw-font-heading tw-font-bold tw-text-sm tw-text-foreground tw-my-3 tw-leading-snug tw-m-0">
            {{ art.title }}
          </h3>
          <p class="tw-text-xs tw-text-muted-foreground tw-line-height-[1.5] tw-mb-4 tw-flex-grow">
            {{ art.desc }}
          </p>
          <div class="tw-flex tw-items-center tw-gap-2 tw-text-[10px] tw-text-muted-foreground tw-mt-auto">
            <span>{{ art.readTime }}</span>
            <span>•</span>
            <span>{{ art.date }}</span>
          </div>
        </AppCard>
      </div>
    </section>

    <!-- Certifications and Salary Guides Row -->
    <div class="tw-grid tw-grid-cols-1 lg:tw-grid-cols-12 tw-gap-6 tw-mt-8">
      
      <!-- Certifications -->
      <section class="lg:tw-col-span-5 glass tw-p-6 tw-rounded-2xl tw-flex tw-flex-col">
        <h3 class="tw-font-heading tw-font-bold tw-text-sm tw-text-foreground tw-m-0">Popular Credentials &amp; Standards</h3>
        <p class="tw-text-xs tw-text-muted-foreground tw-mt-1 tw-mb-4">Enhance your employability credentials by exploring standard accredited pathways.</p>
        
        <div class="tw-flex tw-flex-col tw-gap-4">
          <div 
            v-for="cert in certifications" 
            :key="cert.name" 
            class="tw-flex tw-justify-between tw-items-center tw-pb-3 tw-border-b tw-border-border/50 last:tw-border-b-0"
          >
            <div class="tw-flex tw-flex-col">
              <span class="tw-font-heading tw-font-bold tw-text-xs tw-text-foreground">{{ cert.name }}</span>
              <span class="tw-text-[10px] tw-text-muted-foreground">{{ cert.provider }}</span>
            </div>
            <div class="tw-flex tw-flex-col tw-items-end tw-gap-1">
              <span class="tw-text-[10px] tw-text-muted-foreground">{{ cert.duration }}</span>
              <AppBadge tone="success">{{ cert.cost }}</AppBadge>
            </div>
          </div>
        </div>
      </section>

      <!-- Salary Matrix -->
      <section class="lg:tw-col-span-7 glass tw-p-6 tw-rounded-2xl tw-flex tw-flex-col tw-overflow-x-auto">
        <h3 class="tw-font-heading tw-font-bold tw-text-sm tw-text-foreground tw-m-0">Indian Annual Salary Indicators</h3>
        <p class="tw-text-xs tw-text-muted-foreground tw-mt-1 tw-mb-4">Approximate yearly salary parameters across major private and public hiring sectors.</p>
        
        <table class="tw-w-full tw-border-collapse tw-text-left tw-text-xs">
          <thead>
            <tr class="tw-border-b tw-border-border">
              <th class="tw-pb-2 tw-font-heading tw-font-bold tw-text-muted-foreground">Sector</th>
              <th class="tw-pb-2 tw-font-heading tw-font-bold tw-text-muted-foreground">Entry Level</th>
              <th class="tw-pb-2 tw-font-heading tw-font-bold tw-text-muted-foreground">Mid-Career</th>
              <th class="tw-pb-2 tw-font-heading tw-font-bold tw-text-muted-foreground">Senior / Lead</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="row in salaryGuide" :key="row.sector" class="tw-border-b tw-border-border/50 last:tw-border-none">
              <td class="tw-py-3 tw-font-bold tw-text-foreground">{{ row.sector }}</td>
              <td class="tw-py-3 tw-text-muted-foreground">{{ row.entry }}</td>
              <td class="tw-py-3 tw-text-muted-foreground">{{ row.mid }}</td>
              <td class="tw-py-3 tw-text-muted-foreground">{{ row.senior }}</td>
            </tr>
          </tbody>
        </table>
      </section>

    </div>

  </div>
</template>
