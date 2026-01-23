<template>
  <div class="p-4">
    <table class="w-full">
      <thead>
        <tr class="text-left text-xs" :class="effectiveTheme === 'dark' ? 'text-gray-500' : 'text-gray-400'">
          <th class="pb-2 font-medium">Header</th>
          <th class="pb-2 font-medium">Value</th>
        </tr>
      </thead>
      <tbody>
        <tr 
          v-for="(value, key) in headers" 
          :key="key"
          class="border-t text-sm"
          :class="effectiveTheme === 'dark' ? 'border-dark-border' : 'border-light-border'"
        >
          <td 
            class="py-2 pr-4 font-medium"
            :class="effectiveTheme === 'dark' ? 'text-gray-300' : 'text-gray-700'"
          >
            {{ key }}
          </td>
          <td 
            class="py-2 break-all"
            :class="effectiveTheme === 'dark' ? 'text-gray-400' : 'text-gray-600'"
          >
            {{ value }}
          </td>
        </tr>
      </tbody>
    </table>
    
    <div v-if="Object.keys(headers).length === 0" class="text-center py-8 text-gray-500">
      No headers
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAppStateStore } from '@/stores/appState'

defineProps<{
  headers: Record<string, string>
}>()

const appState = useAppStateStore()
const effectiveTheme = computed(() => appState.effectiveTheme)
</script>
