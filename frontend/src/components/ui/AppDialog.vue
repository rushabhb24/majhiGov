<script setup>
import { onMounted, onUnmounted, watch } from 'vue'

const props = defineProps({
  open: {
    type: Boolean,
    required: true
  },
  maxWidth: {
    type: String,
    default: '580px'
  }
})

const emit = defineEmits(['close'])

const handleKeydown = (e) => {
  if (e.key === 'Escape' && props.open) {
    emit('close')
  }
}

watch(() => props.open, (newVal) => {
  if (newVal) {
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = ''
  }
}, { immediate: true })

onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
  document.body.style.overflow = ''
})
</script>

<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div
        v-if="open"
        class="tw-fixed tw-inset-0 tw-z-[9999] tw-flex tw-items-center tw-justify-center tw-p-4 tw-bg-black/75 tw-backdrop-blur-md"
        @click.self="$emit('close')"
      >
        <div
          class="tw-bg-card tw-text-card-foreground tw-border tw-border-solid tw-border-border tw-w-full tw-rounded-2xl tw-p-6 tw-relative tw-shadow-2xl tw-max-h-[90vh] tw-overflow-y-auto tw-flex tw-flex-col"
          :style="{ maxWidth: maxWidth }"
        >
          <!-- Close button -->
          <button
            class="tw-absolute tw-top-4 tw-right-4 tw-w-8 tw-height-8 tw-flex tw-items-center tw-justify-center tw-rounded-full tw-bg-muted tw-text-muted-foreground hover:tw-bg-muted-foreground/20 tw-transition-colors tw-border-none tw-cursor-pointer"
            @click="$emit('close')"
            aria-label="Close modal"
          >
            ✕
          </button>
          
          <slot />
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
}
.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
</style>
