<script setup>
import { ref } from 'vue';
import AppDialog from './ui/AppDialog.vue'
import AppButton from './ui/AppButton.vue'
import AppBadge from './ui/AppBadge.vue'

const props = defineProps({
  scheme: {
    type: Object,
    default: null
  },
  currentLanguage: {
    type: String,
    required: true
  },
  savedSchemeIds: {
    type: Array,
    required: true
  },
  open: {
    type: Boolean,
    required: true
  },
  t: {
    type: Object,
    required: true
  },
  isLoggedIn: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['close', 'toggleBookmark', 'loginRequired', 'applyClick']);

const activeFaqIndex = ref(null);

function toggleFaq(index) {
  if (activeFaqIndex.value === index) {
    activeFaqIndex.value = null;
  } else {
    activeFaqIndex.value = index;
  }
}

// Localized getters
function getSchemeTitle(scheme) {
  if (!scheme) return '';
  if (props.currentLanguage === 'mr' && scheme.title_mr) return scheme.title_mr;
  if (props.currentLanguage === 'hi' && scheme.title_hi) return scheme.title_hi;
  return scheme.title;
}

function getSchemeDesc(scheme) {
  if (!scheme) return '';
  if (props.currentLanguage === 'mr' && scheme.description_mr) return scheme.description_mr;
  if (props.currentLanguage === 'hi' && scheme.description_hi) return scheme.description_hi;
  return scheme.description;
}

function getDocName(doc) {
  if (props.currentLanguage === 'mr' && doc.document_name_mr) return doc.document_name_mr;
  if (props.currentLanguage === 'hi' && doc.document_name_hi) return doc.document_name_hi;
  return doc.document_name;
}

function getFaqQuestion(faq) {
  if (props.currentLanguage === 'mr' && faq.question_mr) return faq.question_mr;
  if (props.currentLanguage === 'hi' && faq.question_hi) return faq.question_hi;
  return faq.question;
}

function getFaqAnswer(faq) {
  if (props.currentLanguage === 'mr' && faq.answer_mr) return faq.answer_mr;
  if (props.currentLanguage === 'hi' && faq.answer_hi) return faq.answer_hi;
  return faq.answer;
}

function getCategoryName(scheme) {
  if (!scheme) return '';
  const cat = props.currentLanguage === 'mr' ? scheme.category_name_mr : (props.currentLanguage === 'hi' ? scheme.category_name_hi : scheme.category_name);
  return cat || scheme.category_name;
}
</script>

<template>
  <AppDialog
    :open="open && !!scheme"
    @close="emit('close')"
    maxWidth="650px"
  >
    <!-- Tricolor Top Edge -->
    <div class="tricolor-bar tw-h-[4px] tw-absolute tw-top-0 tw-left-0 tw-w-full"></div>

    <div v-if="scheme" class="tw-flex tw-flex-col tw-gap-4 tw-mt-4">
      
      <!-- Category Badge -->
      <div>
        <AppBadge tone="info">
          {{ getCategoryName(scheme) }}
        </AppBadge>
      </div>

      <!-- Title -->
      <h2 class="tw-font-heading tw-font-bold tw-text-xl tw-text-foreground tw-m-0">
        {{ getSchemeTitle(scheme) }}
      </h2>

      <hr class="tw-border-border/50 tw-my-1" />

      <div class="tw-max-h-[60vh] tw-overflow-y-auto tw-pr-2 tw-flex tw-flex-col tw-gap-5">
        
        <!-- Overview Section -->
        <section class="tw-flex tw-flex-col tw-gap-2">
          <h3 class="tw-font-heading tw-font-bold tw-text-sm tw-text-foreground tw-m-0">
            📖 {{ t.modalOverview || 'Overview' }}
          </h3>
          <p class="tw-text-xs tw-text-muted-foreground tw-line-height-[1.6] tw-m-0">
            {{ getSchemeDesc(scheme) }}
          </p>
        </section>

        <!-- Benefits Section -->
        <section class="tw-flex tw-flex-col tw-gap-2 tw-bg-emerald-500/10 tw-border tw-border-emerald-500/25 tw-rounded-2xl tw-p-4">
          <h3 class="tw-font-heading tw-font-bold tw-text-sm tw-text-emerald-600 dark:tw-text-emerald-400 tw-m-0">
            💰 {{ t.modalBenefits || 'Benefits' }}
          </h3>
          <p class="tw-text-xs tw-text-emerald-700 dark:tw-text-emerald-300 tw-font-semibold tw-m-0 tw-line-height-[1.5]">
            {{ scheme.benefits }}
          </p>
        </section>

        <!-- Required Documents Section -->
        <section class="tw-flex tw-flex-col tw-gap-2">
          <h3 class="tw-font-heading tw-font-bold tw-text-sm tw-text-foreground tw-m-0">
            📄 {{ t.modalDocs || 'Required Documents' }}
          </h3>
          <div class="tw-flex tw-flex-col tw-gap-2 tw-mt-1">
            <div 
              v-for="doc in scheme.documents" 
              :key="doc.id" 
              class="tw-flex tw-items-center tw-justify-between tw-p-2.5 tw-bg-muted/40 tw-border tw-border-border tw-rounded-xl tw-text-xs"
            >
              <div class="tw-flex tw-items-center tw-gap-2">
                <span class="tw-text-base">📁</span>
                <span class="tw-font-medium tw-text-foreground">{{ getDocName(doc) }}</span>
              </div>
              <span 
                class="tw-text-[10px] tw-font-bold tw-px-2 tw-py-0.5 tw-rounded-full"
                :class="doc.is_mandatory ? 'tw-bg-destructive/15 tw-text-destructive' : 'tw-bg-success/15 tw-text-success'"
              >
                {{ doc.is_mandatory ? (t.mandatoryBadge || 'Mandatory') : (t.optionalBadge || 'Optional') }}
              </span>
            </div>
          </div>
        </section>

        <!-- Collapsible FAQs accordion -->
        <section v-if="scheme.faqs && scheme.faqs.length > 0" class="tw-flex tw-flex-col tw-gap-2">
          <h3 class="tw-font-heading tw-font-bold tw-text-sm tw-text-foreground tw-m-0">
            ❓ {{ t.modalFaqs || 'Frequently Asked Questions' }}
          </h3>
          <div class="tw-flex tw-flex-col tw-gap-2 tw-mt-1">
            <div 
              v-for="(faq, idx) in scheme.faqs" 
              :key="faq.id" 
              class="tw-border tw-border-border tw-rounded-xl tw-overflow-hidden"
            >
              <button 
                class="tw-w-full tw-flex tw-items-center tw-justify-between tw-p-3 tw-bg-muted/40 hover:tw-bg-muted tw-transition-colors tw-border-none tw-cursor-pointer tw-text-left tw-font-heading tw-font-bold tw-text-xs tw-text-foreground"
                @click="toggleFaq(idx)"
              >
                <span>{{ getFaqQuestion(faq) }}</span>
                <span class="tw-text-xs tw-text-muted-foreground tw-transition-transform" :class="{ 'tw-rotate-180': activeFaqIndex === idx }">▼</span>
              </button>
              
              <div 
                v-if="activeFaqIndex === idx" 
                class="tw-p-3 tw-bg-card tw-border-t tw-border-border tw-text-xs tw-text-muted-foreground tw-line-height-[1.5]"
              >
                {{ getFaqAnswer(faq) }}
              </div>
            </div>
          </div>
        </section>

      </div>

      <hr class="tw-border-border/50 tw-my-1" />

      <!-- Footer Actions -->
      <div class="tw-flex tw-justify-end tw-gap-3">
        <AppButton variant="outline" size="sm" @click="emit('toggleBookmark', scheme.id)">
          {{ savedSchemeIds.includes(scheme.id) ? (t.removeSavedBtn || 'Remove Bookmark') : (t.saveSchemeBtn || 'Save Scheme') }}
        </AppButton>
        
        <AppButton 
          v-if="isLoggedIn" 
          variant="primary" 
          size="sm" 
          class="tw-flex tw-items-center tw-gap-1.5"
          @click="emit('applyClick', scheme)"
        >
          <span>{{ t.applyOnOfficialPortal || 'Apply on Official Portal' }}</span>
          <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path><polyline points="15 3 21 3 21 9"></polyline><line x1="10" y1="14" x2="21" y2="3"></line></svg>
        </AppButton>
        
        <AppButton 
          v-else 
          variant="primary" 
          size="sm" 
          @click="emit('loginRequired', 'apply')"
        >
          🔒 {{ t.loginToApply || 'Login to Apply' }}
        </AppButton>
      </div>

    </div>
  </AppDialog>
</template>
