<template>
  <div class="flex flex-col h-full">
    <!-- Search and clear -->
    <div class="p-2 flex gap-2">
      <div class="flex-1 relative">
        <MagnifyingGlassIcon class="absolute left-2 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search..."
          class="w-full pl-8 pr-2 py-1.5 text-sm rounded-md border outline-none"
          :class="[
            effectiveTheme === 'dark'
              ? 'bg-dark-base border-dark-border text-white placeholder-gray-500 focus:border-accent'
              : 'bg-white border-light-border text-gray-900 placeholder-gray-400 focus:border-accent'
          ]"
        />
      </div>
      <button
        v-if="history.length > 0"
        @click="clearHistory"
        class="p-1.5 rounded-md transition-colors"
        :class="effectiveTheme === 'dark' ? 'hover:bg-dark-hover text-gray-400 hover:text-red-400' : 'hover:bg-light-hover text-gray-500 hover:text-red-500'"
        title="Clear all history"
      >
        <TrashIcon class="w-5 h-5" />
      </button>
    </div>
    
    <!-- History list -->
    <div class="flex-1 overflow-auto px-2 pb-2">
      <div v-if="filteredGroups.length === 0" class="text-center py-8 text-gray-500 text-sm">
        {{ searchQuery ? 'No results found' : 'No history yet' }}
      </div>
      
      <div v-for="group in filteredGroups" :key="group.date" class="mb-3">
        <div class="text-xs font-medium mb-1 px-2" :class="effectiveTheme === 'dark' ? 'text-gray-500' : 'text-gray-400'">
          {{ group.label }}
        </div>
        
        <div
          v-for="item in group.items"
          :key="item.id"
          class="flex items-center gap-2 px-2 py-1 rounded-md cursor-pointer group"
          :class="effectiveTheme === 'dark' ? 'hover:bg-dark-hover' : 'hover:bg-light-hover'"
          :title="item.url"
          @click="openHistory(item)"
        >
          <span class="text-xs text-gray-500 flex-shrink-0">
            {{ formatTime(item.createdAt) }}
          </span>
          <span 
            class="text-xs font-medium w-10 flex-shrink-0"
            :class="methodColor(item.method)"
          >
            {{ item.method }}
          </span>
          <span 
            class="flex-1 truncate text-sm"
            :class="effectiveTheme === 'dark' ? 'text-gray-300' : 'text-gray-700'"
          >
            {{ formatUrl(item.url) }}
          </span>
          <span 
            v-if="item.statusCode"
            class="text-xs font-medium flex-shrink-0"
            :class="statusColor(item.statusCode)"
          >
            {{ item.statusCode }}
          </span>
          <button
            @click.stop="deleteHistoryItem(item.id)"
            class="p-0.5 rounded opacity-0 group-hover:opacity-100 hover:bg-red-500/20 text-red-500 flex-shrink-0"
            title="Delete"
          >
            <XMarkIcon class="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { MagnifyingGlassIcon, TrashIcon, XMarkIcon } from '@heroicons/vue/24/outline'
import { useAppStateStore } from '@/stores/appState'
import { useHistoryStore } from '@/stores/history'
import { useTabsStore } from '@/stores/tabs'
import { api } from '@/services/api'
import type { History } from '@/types'

const appState = useAppStateStore()
const historyStore = useHistoryStore()
const tabsStore = useTabsStore()

const effectiveTheme = computed(() => appState.effectiveTheme)
const searchQuery = ref('')
const history = computed(() => historyStore.history)

const filteredGroups = computed(() => {
  if (!searchQuery.value) return historyStore.groupedHistory
  
  const query = searchQuery.value.toLowerCase()
  return historyStore.groupedHistory
    .map(group => ({
      ...group,
      items: group.items.filter(item =>
        item.url.toLowerCase().includes(query) ||
        item.method.toLowerCase().includes(query)
      )
    }))
    .filter(group => group.items.length > 0)
})

function formatTime(dateStr: string) {
  const date = new Date(dateStr)
  return date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', hour12: false })
}

function formatUrl(url: string) {
  try {
    const parsed = new URL(url)
    return parsed.pathname + parsed.search
  } catch {
    return url
  }
}

function methodColor(method: string) {
  switch (method.toUpperCase()) {
    case 'GET': return 'text-method-get'
    case 'POST': return 'text-method-post'
    case 'PUT': return 'text-method-put'
    case 'PATCH': return 'text-method-patch'
    case 'DELETE': return 'text-method-delete'
    default: return 'text-method-options'
  }
}

function statusColor(code: number) {
  if (code >= 200 && code < 300) return 'text-status-success'
  if (code >= 300 && code < 400) return 'text-status-redirect'
  if (code >= 400 && code < 500) return 'text-status-client-error'
  return 'text-status-server-error'
}

function openHistory(item: History) {
  let headers: any[] = []
  let params: any[] = []
  
  try {
    headers = JSON.parse(item.requestHeaders || '[]')
  } catch {}
  
  // Always open history as new unsaved tab
  tabsStore.openHistoryAsNewTab(
    `${item.method} ${formatUrl(item.url)}`,
    item.method,
    item.url,
    headers,
    params,
    item.requestBody || '',
    'json'
  )
}

async function deleteHistoryItem(id: number) {
  const modal = (window as any).$modal
  if (modal) {
    const confirmed = await modal.confirm({
      title: 'Delete History Item',
      message: 'Are you sure you want to delete this history item?',
      confirmText: 'Delete',
      danger: true,
    })
    
    if (!confirmed) return
  }
  
  try {
    await api.deleteHistoryItem(id)
    historyStore.deleteHistory(id)
  } catch (error) {
    console.error('Failed to delete history item:', error)
  }
}

async function clearHistory() {
  const modal = (window as any).$modal
  if (modal) {
    const confirmed = await modal.confirm({
      title: 'Clear All History',
      message: 'Are you sure you want to clear all history? This action cannot be undone.',
      confirmText: 'Clear All',
      danger: true,
    })
    
    if (!confirmed) return
  }
  
  try {
    await api.clearHistory()
    historyStore.clearHistory()
  } catch (error) {
    console.error('Failed to clear history:', error)
  }
}
</script>
