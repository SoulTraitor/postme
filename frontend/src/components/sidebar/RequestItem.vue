<template>
  <div
    class="flex items-center gap-2 px-2 py-1 rounded-md cursor-pointer"
    :class="[
      isHighlighted 
        ? 'bg-accent/20' 
        : (effectiveTheme === 'dark' ? 'hover:bg-dark-hover' : 'hover:bg-light-hover')
    ]"
    @click="$emit('click')"
    @dblclick="$emit('dblclick')"
  >
    <span 
      class="text-xs font-medium w-12 text-center flex-shrink-0"
      :class="methodColor"
    >
      {{ request.method }}
    </span>
    <span 
      class="flex-1 truncate text-sm"
      :class="effectiveTheme === 'dark' ? 'text-gray-300' : 'text-gray-700'"
    >
      {{ displayName }}
    </span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAppStateStore } from '@/stores/appState'
import type { Request } from '@/types'

const props = defineProps<{
  request: Request
}>()

defineEmits<{
  click: []
  dblclick: []
}>()

const appState = useAppStateStore()
const effectiveTheme = computed(() => appState.effectiveTheme)
const isHighlighted = computed(() => appState.highlightedRequestId === props.request.id)

const displayName = computed(() => {
  if (props.request.name) return props.request.name
  if (props.request.url) {
    try {
      const url = new URL(props.request.url)
      return url.pathname
    } catch {
      return props.request.url
    }
  }
  return 'Untitled'
})

const methodColor = computed(() => {
  switch (props.request.method.toUpperCase()) {
    case 'GET': return 'text-method-get'
    case 'POST': return 'text-method-post'
    case 'PUT': return 'text-method-put'
    case 'PATCH': return 'text-method-patch'
    case 'DELETE': return 'text-method-delete'
    default: return 'text-method-options'
  }
})
</script>
