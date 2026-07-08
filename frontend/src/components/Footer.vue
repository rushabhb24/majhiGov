<script setup>
import { ref } from 'vue'
import { API_BASE_URL } from '../config.js'

const email = ref('')
const subscribing = ref(false)
const subscribeMsg = ref('')

async function handleSubscribe() {
  if (!email.value) return
  subscribing.value = true
  subscribeMsg.value = ''
  try {
    await fetch(`${API_BASE_URL}/api/newsletter/subscribe`, {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: email.value })
    })
    subscribeMsg.value = 'Subscribed! Thank you 🎉'
    email.value = ''
  } catch {
    subscribeMsg.value = 'Could not subscribe. Try again.'
  } finally {
    subscribing.value = false
  }
}
</script>

<template>
  <footer class="tw-relative tw-bg-navy tw-text-slate-300 tw-mt-16 tw-border-t tw-border-border/10">
    <!-- Tricolor Top Bar -->
    <div class="tricolor-bar tw-h-[4px] tw-w-full"></div>

    <div class="tw-max-w-7xl tw-mx-auto tw-px-4 tw-sm:px-6 tw-lg:px-8 tw-py-12">
      <!-- Main Columns Grid -->
      <div class="tw-grid tw-grid-cols-1 md:tw-grid-cols-12 tw-gap-8 tw-mb-12">
        
        <!-- Column 1: Brand (About) -->
        <div class="md:tw-col-span-4 tw-flex tw-flex-col tw-gap-4">
          <div class="tw-flex tw-items-center">
            <div class="tw-w-9 tw-h-9 tw-rounded-xl tw-bg-primary tw-text-white tw-font-heading tw-font-extrabold tw-flex tw-items-center tw-justify-center tw-text-lg">
              M
            </div>
            <span class="tw-ml-3 tw-font-heading tw-font-extrabold tw-text-xl tw-text-white">MajhiGov</span>
          </div>
          <p class="tw-text-sm tw-text-slate-400 tw-line-height-[1.6]">
            India's most trusted AI-powered portal for discovering government schemes, job announcements, and career assistance. Empowering every citizen with glassmorphism ease.
          </p>
          <div class="tw-flex tw-gap-2">
            <!-- Twitter/X -->
            <a href="#" class="tw-w-8 tw-h-8 tw-rounded-lg tw-bg-white/5 tw-border tw-border-white/10 tw-flex tw-items-center tw-justify-center tw-text-slate-400 hover:tw-bg-primary/20 hover:tw-text-primary tw-transition-all" aria-label="Twitter">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-4.714-6.231-5.401 6.231H2.748l7.73-8.835L1.254 2.25H8.08l4.259 5.63zm-1.161 17.52h1.833L7.084 4.126H5.117z"/></svg>
            </a>
            <!-- LinkedIn -->
            <a href="#" class="tw-w-8 tw-h-8 tw-rounded-lg tw-bg-white/5 tw-border tw-border-white/10 tw-flex tw-items-center tw-justify-center tw-text-slate-400 hover:tw-bg-primary/20 hover:tw-text-primary tw-transition-all" aria-label="LinkedIn">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433a2.062 2.062 0 0 1-2.063-2.065 2.064 2.064 0 1 1 2.063 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z"/></svg>
            </a>
          </div>
        </div>

        <!-- Column 2: Quick Links -->
        <div class="md:tw-col-span-2 tw-flex tw-flex-col tw-gap-3">
          <h4 class="tw-font-heading tw-font-bold tw-text-sm tw-text-white tw-tracking-wide">Quick Links</h4>
          <ul class="tw-list-none tw-p-0 tw-m-0 tw-flex tw-flex-col tw-gap-2 tw-text-sm">
            <li><a href="#/" class="tw-text-slate-400 hover:tw-text-primary tw-transition-colors tw-no-underline">Explorer</a></li>
            <li><a href="#/eligibility" class="tw-text-slate-400 hover:tw-text-primary tw-transition-colors tw-no-underline">Eligibility Wizard</a></li>
            <li><a href="#/jobs" class="tw-text-slate-400 hover:tw-text-primary tw-transition-colors tw-no-underline">Sarkari Jobs</a></li>
            <li><a href="#/ai-assistant" class="tw-text-slate-400 hover:tw-text-primary tw-transition-colors tw-no-underline">AI Assistant</a></li>
            <li><a href="#/career-resources" class="tw-text-slate-400 hover:tw-text-primary tw-transition-colors tw-no-underline">Resources</a></li>
          </ul>
        </div>

        <!-- Column 3: Contact & Support -->
        <div class="md:tw-col-span-2 tw-flex tw-flex-col tw-gap-3">
          <h4 class="tw-font-heading tw-font-bold tw-text-sm tw-text-white tw-tracking-wide">Support</h4>
          <ul class="tw-list-none tw-p-0 tw-m-0 tw-flex tw-flex-col tw-gap-2 tw-text-sm">
            <li><a href="#" class="tw-text-slate-400 hover:tw-text-primary tw-transition-colors tw-no-underline">About Us</a></li>
            <li><a href="#" class="tw-text-slate-400 hover:tw-text-primary tw-transition-colors tw-no-underline">Contact Us</a></li>
            <li><a href="#" class="tw-text-slate-400 hover:tw-text-primary tw-transition-colors tw-no-underline">Privacy Policy</a></li>
            <li><a href="#" class="tw-text-slate-400 hover:tw-text-primary tw-transition-colors tw-no-underline">Terms of Service</a></li>
          </ul>
        </div>

        <!-- Column 4: Newsletter -->
        <div class="md:tw-col-span-4 tw-flex tw-flex-col tw-gap-3">
          <h4 class="tw-font-heading tw-font-bold tw-text-sm tw-text-white tw-tracking-wide">Stay Updated</h4>
          <p class="tw-text-xs tw-text-slate-400 tw-line-height-[1.5]">
            Subscribe to our newsletter to receive updates on newly launched central/state schemes and hot vacancies.
          </p>
          <form class="tw-flex tw-gap-2" @submit.prevent="handleSubscribe">
            <input
              v-model="email"
              type="email"
              placeholder="Enter your email"
              class="tw-flex-1 tw-px-3 tw-py-2 tw-text-xs tw-bg-white/5 tw-text-white tw-border tw-border-white/10 tw-rounded-lg tw-outline-none focus:tw-border-primary tw-transition-colors"
              required
            />
            <button
              type="submit"
              class="tw-px-4 tw-py-2 tw-rounded-lg tw-bg-primary tw-text-white tw-font-heading tw-font-bold tw-text-xs tw-border-none tw-cursor-pointer tw-transition-transform hover:tw-scale-105"
              :disabled="subscribing"
            >
              {{ subscribing ? '...' : 'Subscribe' }}
            </button>
          </form>
          <p v-if="subscribeMsg" class="tw-text-xs tw-text-primary tw-m-0">{{ subscribeMsg }}</p>
        </div>

      </div>

      <!-- Bottom Bar -->
      <div class="tw-border-t tw-border-white/5 tw-pt-6 tw-flex tw-flex-col sm:tw-flex-row tw-items-center tw-justify-between tw-gap-4 tw-text-xs tw-text-slate-500">
        <span>© 2026 MajhiGov. All rights reserved.</span>
        <span>Made with ❤️ for citizens of India</span>
      </div>

    </div>
  </footer>
</template>
