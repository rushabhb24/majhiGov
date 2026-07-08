<script setup>
defineProps({
  modelValue: {
    type: [String, Number, Boolean],
    default: ''
  },
  options: {
    type: Array,
    default: () => [] // Array of { value, label } or simple strings
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

defineEmits(['update:modelValue', 'change'])
</script>

<template>
  <select
    :id="id"
    :value="modelValue"
    :disabled="disabled"
    :required="required"
    @change="$emit('update:modelValue', $event.target.value); $emit('change', $event)"
    class="tw-w-full tw-px-3 tw-py-2 tw-text-sm tw-bg-background tw-text-foreground tw-border tw-border-border tw-rounded-lg tw-transition-all tw-duration-150 tw-outline-none focus:tw-border-primary focus-visible:tw-ring-2 focus-visible:tw-ring-ring focus-visible:tw-ring-offset-2 disabled:tw-opacity-50"
  >
    <slot>
      <option
        v-for="opt in options"
        :key="typeof opt === 'object' ? opt.value : opt"
        :value="typeof opt === 'object' ? opt.value : opt"
      >
        {{ typeof opt === 'object' ? opt.label : opt }}
      </option>
    </slot>
  </select>
</template>
