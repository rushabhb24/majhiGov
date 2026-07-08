<script setup>
defineProps({
  modelValue: {
    type: String,
    required: true
  },
  tabs: {
    type: Array,
    required: true // Array of { id, label, icon }
  }
})

defineEmits(['update:modelValue'])
</script>

<template>
  <div class="tw-flex tw-flex-col tw-w-full">
    <!-- Tabs List -->
    <div class="tw-flex tw-gap-1 tw-bg-muted/60 tw-p-1 tw-rounded-xl tw-w-full tw-overflow-x-auto">
      <button
        v-for="t in tabs"
        :key="t.id"
        class="tw-flex-1 tw-flex tw-items-center tw-justify-center tw-gap-2 tw-py-2 tw-px-3 tw-text-sm tw-font-semibold tw-font-heading tw-rounded-lg tw-transition-all tw-border-none tw-cursor-pointer tw-outline-none focus-visible:tw-ring-2 focus-visible:tw-ring-ring"
        :class="modelValue === t.id ? 'tw-bg-card tw-text-foreground tw-shadow-sm' : 'tw-text-muted-foreground hover:tw-text-foreground tw-bg-transparent'"
        @click="$emit('update:modelValue', t.id)"
      >
        <span v-if="t.icon" class="tw-text-base">{{ t.icon }}</span>
        <span>{{ t.label }}</span>
      </button>
    </div>
    <!-- Tabs Content -->
    <div class="tw-mt-4">
      <slot :activeTab="modelValue" />
    </div>
  </div>
</template>
