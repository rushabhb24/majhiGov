<script setup>
defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: ''
  },
  rows: {
    type: [Number, String],
    default: 3
  },
  disabled: {
    type: Boolean,
    default: false
  },
  required: {
    type: Boolean,
    default: false
  },
  id: {
    type: String,
    default: ''
  }
})

defineEmits(['update:modelValue', 'input'])
</script>

<template>
  <textarea
    :id="id"
    :value="modelValue"
    :placeholder="placeholder"
    :rows="rows"
    :disabled="disabled"
    :required="required"
    @input="$emit('update:modelValue', $event.target.value); $emit('input', $event)"
    class="tw-w-full tw-px-3 tw-py-2 tw-text-sm tw-bg-background tw-text-foreground tw-border tw-border-border tw-rounded-lg tw-transition-all tw-duration-150 tw-outline-none focus:tw-border-primary focus-visible:tw-ring-2 focus-visible:tw-ring-ring focus-visible:tw-ring-offset-2 disabled:tw-opacity-50 tw-resize-y"
  ></textarea>
</template>
