<script setup>
import { ref } from 'vue';

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
  <Transition name="modal-fade">
    <div v-if="open && scheme" class="modal-overlay" @click.self="emit('close')">
      <div class="modal-content card">
        <button class="btn-close-modal" @click="emit('close')" title="Close Modal">×</button>
        
        <span class="modal-category">
          {{ getCategoryName(scheme) }}
        </span>
        <h2 class="modal-title">{{ getSchemeTitle(scheme) }}</h2>
        
        <hr class="divider mt-4" />

        <div class="modal-scroll-area">
          <!-- Overview Section -->
          <section class="modal-section mt-4">
            <h3>📖 {{ t.modalOverview }}</h3>
            <p class="mt-4">{{ getSchemeDesc(scheme) }}</p>
          </section>

          <!-- Benefits Section -->
          <section class="modal-section mt-4 benefits-section">
            <h3>💰 {{ t.modalBenefits }}</h3>
            <p class="mt-4 benefit-text">{{ scheme.benefits }}</p>
          </section>

          <!-- Required Documents Section with Clamping/Alignment -->
          <section class="modal-section mt-4">
            <h3>📄 {{ t.modalDocs }}</h3>
            <ul class="docs-list mt-4">
              <li v-for="doc in scheme.documents" :key="doc.id">
                <span class="doc-icon">📁</span> 
                <span class="doc-text">{{ getDocName(doc) }}</span>
                <!-- Mandatory/Optional Badges -->
                <span :class="['doc-badge', doc.is_mandatory ? 'mandatory' : 'optional']">
                  {{ doc.is_mandatory ? t.mandatoryBadge : t.optionalBadge }}
                </span>
              </li>
            </ul>
          </section>

          <!-- Collapsible FAQs with rotate-arrow hooks -->
          <section v-if="scheme.faqs && scheme.faqs.length > 0" class="modal-section mt-4">
            <h3>❓ {{ t.modalFaqs }}</h3>
            <div class="faqs-accordion mt-4">
              <div 
                v-for="(faq, idx) in scheme.faqs" 
                :key="faq.id" 
                :class="['faq-panel', { active: activeFaqIndex === idx }]"
              >
                <button class="faq-header" @click="toggleFaq(idx)">
                  <span>{{ getFaqQuestion(faq) }}</span>
                  <span class="faq-arrow">▼</span>
                </button>
                <Transition name="accordion">
                  <div v-show="activeFaqIndex === idx" class="faq-body">
                    <p>{{ getFaqAnswer(faq) }}</p>
                  </div>
                </Transition>
              </div>
            </div>
          </section>
        </div>

        <!-- Footer actions -->
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="emit('toggleBookmark', scheme.id)">
            {{ savedSchemeIds.includes(scheme.id) ? t.removeSavedBtn : t.saveSchemeBtn }}
          </button>
          <button 
            v-if="isLoggedIn"
            class="btn btn-primary text-center" 
            @click="emit('applyClick', scheme)"
          >
            🔗 {{ t.applyOnOfficialPortal || 'Apply on Official Portal' }}
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path><polyline points="15 3 21 3 21 9"></polyline><line x1="10" y1="14" x2="21" y2="3"></line></svg>
          </button>
          <button 
            v-else 
            class="btn btn-primary text-center"
            @click="emit('loginRequired', 'apply')"
          >
            🔒 {{ t.loginToApply || 'Login to Apply' }}
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.3s ease;
}
.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

.modal-fade-enter-active .modal-content {
  animation: zoomElastic 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.modal-fade-leave-active .modal-content {
  animation: zoomElastic 0.3s cubic-bezier(0.34, 1.56, 0.64, 1) reverse;
}

/* Accordion slide animation */
.accordion-enter-active,
.accordion-leave-active {
  transition: max-height 0.3s ease-out, opacity 0.2s ease;
  max-height: 600px;
  overflow: hidden;
}
.accordion-enter-from,
.accordion-leave-to {
  max-height: 0;
  opacity: 0;
}

.doc-text {
  flex: 1 1 140px;
  min-width: 0;
  word-break: break-word;
  overflow-wrap: anywhere;
  padding-right: 8px;
}
</style>
