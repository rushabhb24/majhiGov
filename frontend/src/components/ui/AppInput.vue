<script setup>
defineProps({
  modelValue: {
    type: [String, Number],
    default: ''
  },
  type: {
    type: String,
    default: 'text'
  },
  placeholder: {
    type: String,
    default: ''
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

defineEmits(['update:modelValue', 'input', 'change'])
</script>

<template>
  <input
    :id="id"
    :type="type"
    :value="modelValue"
    :placeholder="placeholder"
    :disabled="disabled"
    :required="required"
    @input="$emit('update:modelValue', $event.target.value); $emit('input', $event)"
    @change="$emit('change', $event)"
    class="tw-w-full tw-px-3 tw-py-2 tw-text-sm tw-bg-background tw-text-foreground tw-border tw-border-border tw-rounded-lg tw-transition-all tw-duration-150 tw-outline-none focus:tw-border-primary focus-visible:tw-ring-2 focus-visible:tw-ring-ring focus-visible:tw-ring-offset-2 disabled:tw-opacity-50"
  />
</template>
