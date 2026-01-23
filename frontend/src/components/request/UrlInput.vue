<template>
  <input
    type="text"
    :value="modelValue"
    @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    @keydown.enter="$emit('submit')"
    placeholder="Enter URL or paste text"
    class="w-full px-3 py-2 rounded-md border outline-none text-sm"
    :class="[
      effectiveTheme === 'dark'
        ? 'bg-dark-surface border-dark-border text-white placeholder-gray-500 focus:border-accent'
        : 'bg-light-surface border-light-border text-gray-900 placeholder-gray-400 focus:border-accent'
    ]"
  />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAppStateStore } from '@/stores/appState'

defineProps<{
  modelValue: string
}>()

defineEmits<{
  'update:modelValue': [value: string]
  'submit': []
}>()

const appState = useAppStateStore()
const effectiveTheme = computed(() => appState.effectiveTheme)
</script>
