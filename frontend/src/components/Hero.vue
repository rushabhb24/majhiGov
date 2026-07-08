<script setup>
defineProps({ t: Object })
const emit = defineEmits(['start-check'])
</script>

<template>
  <section class="hero-gradient tw-relative tw-overflow-hidden tw-rounded-3xl tw-text-center tw-px-6 tw-pt-16 tw-pb-0 tw-mb-10 tw-shadow-glow">
    <!-- Floating decorative blobs -->
    <div class="blob-1 tw-absolute tw-w-96 tw-h-96 tw-rounded-full tw-bg-white/10 tw-blur-3xl tw-top-[-10%] tw-left-[-10%] tw-pointer-events-none"></div>
    <div class="blob-2 tw-absolute tw-w-96 tw-h-96 tw-rounded-full tw-bg-white/10 tw-blur-3xl tw-bottom-[-10%] tw-right-[-10%] tw-pointer-events-none"></div>

    <div class="tw-relative tw-z-10 tw-max-w-3xl tw-mx-auto">
      <div class="tw-inline-flex tw-items-center tw-gap-2 tw-bg-white/15 tw-backdrop-blur-md tw-border tw-border-white/25 tw-text-white tw-px-4 tw-py-1.5 tw-rounded-full tw-text-xs tw-font-bold tw-tracking-wide tw-mb-6">
        <span class="tw-w-2 tw-h-2 tw-rounded-full tw-bg-emerald-400 tw-animate-pulse"></span>
        AI-Powered MajhiGov Portal
      </div>

      <h1 class="tw-font-display tw-text-4xl sm:tw-text-5xl tw-font-black tw-text-white tw-leading-tight tw-tracking-tight tw-mb-4">
        Discover Government Schemes &amp; Careers
      </h1>

      <p class="tw-text-base sm:tw-text-lg tw-text-white/90 tw-mb-8 tw-max-w-2xl tw-mx-auto tw-leading-relaxed">
        Unlock your entitlements and career match points. Find schemes for Farmers, Students, Women, and Senior Citizens using AI matching.
      </p>

      <!-- Actions -->
      <div class="tw-flex tw-justify-center tw-gap-4 tw-flex-wrap tw-mb-12">
        <button 
          class="tw-px-6 tw-py-3 tw-rounded-full tw-bg-white tw-text-primary tw-font-heading tw-font-bold tw-text-sm tw-border-none tw-cursor-pointer tw-transition-all hover:tw-translate-y-[-2px] hover:tw-shadow-xl"
          @click="emit('start-check')"
        >
          Check Eligibility
        </button>
        <button 
          class="tw-px-6 tw-py-3 tw-rounded-full tw-bg-white/15 tw-text-white tw-font-heading tw-font-bold tw-text-sm tw-border-2 tw-border-white/40 tw-cursor-pointer tw-transition-all hover:tw-bg-white/20 hover:tw-translate-y-[-2px] hover:tw-shadow-xl"
          @click="$router.push('/')"
        >
          Browse Schemes
        </button>
      </div>
    </div>

    <!-- Stat Strip -->
    <div class="glass tw-border-b-0 tw-rounded-t-2xl tw-mx-auto tw-max-w-5xl tw-px-6 tw-py-5 tw-grid tw-grid-cols-2 sm:tw-grid-cols-4 tw-gap-4">
      <div class="tw-text-center">
        <div class="tw-font-display tw-text-2xl tw-font-black tw-text-white">10,000+</div>
        <div class="tw-text-[10px] tw-font-bold tw-text-white/80 tw-uppercase tw-tracking-wider tw-mt-1">Schemes</div>
      </div>
      <div class="tw-text-center">
        <div class="tw-font-display tw-text-2xl tw-font-black tw-text-white">5,000+</div>
        <div class="tw-text-[10px] tw-font-bold tw-text-white/80 tw-uppercase tw-tracking-wider tw-mt-1">Vacancies</div>
      </div>
      <div class="tw-text-center">
        <div class="tw-font-display tw-text-2xl tw-font-black tw-text-white">50,000+</div>
        <div class="tw-text-[10px] tw-font-bold tw-text-white/80 tw-uppercase tw-tracking-wider tw-mt-1">Citizens Helped</div>
      </div>
      <div class="tw-text-center">
        <div class="tw-font-display tw-text-2xl tw-font-black tw-text-white">95%</div>
        <div class="tw-text-[10px] tw-font-bold tw-text-white/80 tw-uppercase tw-tracking-wider tw-mt-1">Success Rate</div>
      </div>
    </div>
  </section>
</template>
