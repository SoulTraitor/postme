<template>
  <div
    class="group flex items-center gap-1 px-3 py-1.5 border-r cursor-pointer min-w-0 max-w-[200px] relative"
    :class="[
      isActive 
        ? (effectiveTheme === 'dark' ? 'bg-dark-surface' : 'bg-white')
        : (effectiveTheme === 'dark' ? 'hover:bg-dark-hover' : 'hover:bg-light-hover'),
      effectiveTheme === 'dark' ? 'border-dark-border' : 'border-light-border'
    ]"
    @click="$emit('click')"
    @dblclick="$emit('dblclick')"
  >
    <!-- Active indicator bar -->
    <div 
      v-if="isActive"
      class="absolute bottom-0 left-0 right-0 h-0.5 bg-accent"
    />
    <!-- Dirty indicator -->
    <span 
      v-if="tab.isDirty" 
      class="w-2 h-2 rounded-full bg-accent flex-shrink-0"
    />
    
    <!-- Method badge -->
    <span 
      class="text-xs font-medium flex-shrink-0"
      :class="methodColor"
    >
      {{ tab.method }}
    </span>
    
    <!-- Title -->
    <span 
      class="truncate text-sm flex-1 min-w-0"
      :class="[
        tab.isPreview ? 'italic' : '',
        effectiveTheme === 'dark' ? 'text-gray-300' : 'text-gray-700'
      ]"
    >
      {{ displayTitle }}
    </span>
    
    <!-- Close button -->
    <button
      @click.stop="$emit('close')"
      class="p-0.5 rounded opacity-0 group-hover:opacity-100 transition-opacity flex-shrink-0"
      :class="effectiveTheme === 'dark' ? 'hover:bg-dark-hover' : 'hover:bg-light-hover'"
    >
      <XMarkIcon class="w-3.5 h-3.5" />
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { XMarkIcon } from '@heroicons/vue/24/outline'
import { useAppStateStore } from '@/stores/appState'
import type { Tab } from '@/types'

const props = defineProps<{
  tab: Tab
  isActive: boolean
}>()

defineEmits<{
  click: []
  close: []
  dblclick: []
}>()

const appState = useAppStateStore()
const effectiveTheme = computed(() => appState.effectiveTheme)

const displayTitle = computed(() => {
  if (props.tab.title && props.tab.title !== 'Untitled') {
    return props.tab.title
  }
  if (props.tab.url) {
    try {
      const url = new URL(props.tab.url)
      return url.pathname || props.tab.url
    } catch {
      return props.tab.url || 'Untitled'
    }
  }
  return 'Untitled'
})

const methodColor = computed(() => {
  switch (props.tab.method.toUpperCase()) {
    case 'GET': return 'text-method-get'
    case 'POST': return 'text-method-post'
    case 'PUT': return 'text-method-put'
    case 'PATCH': return 'text-method-patch'
    case 'DELETE': return 'text-method-delete'
    default: return 'text-method-options'
  }
})
</script>
