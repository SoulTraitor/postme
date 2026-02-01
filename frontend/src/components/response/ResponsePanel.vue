<template>
  <div 
    class="flex flex-col h-full overflow-hidden"
    :class="effectiveTheme === 'dark' ? 'bg-dark-base' : 'bg-light-base'"
  >
    <!-- Response state: Idle -->
    <div v-if="responseState.status === 'idle'" class="flex-1 flex flex-col items-center justify-center text-gray-500">
      <div class="relative">
        <div class="absolute inset-0 bg-accent/5 blur-2xl rounded-full"></div>
        <PaperAirplaneIcon class="w-16 h-16 mb-6 opacity-40 relative" />
      </div>
      <p class="text-lg font-medium mb-2" :class="effectiveTheme === 'dark' ? 'text-gray-300' : 'text-gray-700'">
        Ready to send
      </p>
      <p class="text-sm mb-4">Click Send to make a request</p>
      <div class="flex items-center gap-2 text-xs px-3 py-1.5 rounded-md" :class="effectiveTheme === 'dark' ? 'bg-dark-surface' : 'bg-gray-100'">
        <kbd class="px-2 py-0.5 rounded font-mono" :class="effectiveTheme === 'dark' ? 'bg-dark-hover' : 'bg-white'">Ctrl</kbd>
        <span>+</span>
        <kbd class="px-2 py-0.5 rounded font-mono" :class="effectiveTheme === 'dark' ? 'bg-dark-hover' : 'bg-white'">Enter</kbd>
        <span class="ml-1">to send</span>
      </div>
    </div>
    
    <!-- Response state: Loading -->
    <div v-else-if="responseState.status === 'loading'" class="flex-1 flex flex-col items-center justify-center text-gray-500">
      <div class="relative mb-6">
        <!-- Outer pulse ring -->
        <div class="absolute inset-0 w-12 h-12 -m-2 rounded-full bg-accent/20 animate-ping"></div>
        <!-- Spinner -->
        <div class="relative w-8 h-8 border-2 border-accent border-t-transparent rounded-full animate-spin"></div>
      </div>
      <p class="text-lg font-medium mb-2" :class="effectiveTheme === 'dark' ? 'text-gray-300' : 'text-gray-700'">
        Sending request
      </p>
      <p class="text-sm">Please wait...</p>
    </div>
    
    <!-- Response state: Cancelled -->
    <div v-else-if="responseState.status === 'cancelled'" class="flex-1 flex flex-col items-center justify-center text-gray-500">
      <NoSymbolIcon class="w-12 h-12 mb-4 opacity-50" />
      <p>Request was cancelled</p>
    </div>
    
    <!-- Response state: Timeout -->
    <div v-else-if="responseState.status === 'timeout'" class="flex-1 flex flex-col items-center justify-center text-gray-500">
      <ClockIcon class="w-12 h-12 mb-4 opacity-50" />
      <p>Request timed out ({{ responseState.seconds }}s)</p>
    </div>
    
    <!-- Response state: Error -->
    <div v-else-if="responseState.status === 'error'" class="flex-1 flex flex-col items-center justify-center text-red-500">
      <ExclamationCircleIcon class="w-12 h-12 mb-4" />
      <p>{{ responseState.message }}</p>
    </div>
    
    <!-- Response state: Success -->
    <template v-else-if="responseState.status === 'success'">
      <!-- Status bar -->
      <div 
        class="flex items-center gap-4 px-4 py-2 border-b"
        :class="effectiveTheme === 'dark' ? 'border-dark-border' : 'border-light-border'"
      >
        <!-- Status code -->
        <div class="flex items-center gap-2">
          <span
            class="px-3 py-1 rounded-full font-medium text-sm flex items-center gap-1.5"
            :class="statusBadgeClass"
          >
            <component :is="statusIcon" v-if="statusIcon" class="w-4 h-4" />
            {{ responseState.response.statusCode }}
          </span>
          <span
            class="text-sm"
            :class="effectiveTheme === 'dark' ? 'text-gray-400' : 'text-gray-600'"
          >
            {{ responseState.response.status.split(' ').slice(1).join(' ') }}
          </span>
        </div>
        
        <!-- Duration -->
        <span class="status-value text-sm text-gray-500">
          {{ responseState.response.duration }}ms
        </span>

        <!-- Size -->
        <span class="status-value text-sm text-gray-500">
          {{ formatSize(responseState.response.size) }}
        </span>
        
        <!-- Layout toggle -->
        <button
          @click="appState.toggleLayoutDirection"
          class="ml-auto p-1.5 rounded-md transition-colors"
          :class="effectiveTheme === 'dark' ? 'hover:bg-dark-hover text-gray-400' : 'hover:bg-light-hover text-gray-500'"
          :title="appState.layoutDirection === 'horizontal' ? 'Switch to vertical layout' : 'Switch to horizontal layout'"
        >
          <ArrowsUpDownIcon v-if="appState.layoutDirection === 'horizontal'" class="w-4 h-4" />
          <ArrowsRightLeftIcon v-else class="w-4 h-4" />
        </button>
      </div>
      
      <!-- Response tabs -->
      <div class="flex-1 flex flex-col overflow-hidden">
        <!-- Tab headers -->
        <div class="flex border-b" :class="effectiveTheme === 'dark' ? 'border-dark-border' : 'border-light-border'">
          <button
            v-for="tab in responseTabs"
            :key="tab.id"
            @click="activeResponseTab = tab.id"
            class="px-4 py-2 text-sm font-medium transition-colors relative"
            :class="[
              activeResponseTab === tab.id
                ? 'text-accent'
                : (effectiveTheme === 'dark' ? 'text-gray-400 hover:text-white' : 'text-gray-500 hover:text-gray-900')
            ]"
          >
            {{ tab.label }}
            <div
              v-if="activeResponseTab === tab.id"
              class="absolute bottom-0 left-0 right-0 h-0.5 bg-accent"
            />
          </button>
        </div>
        
        <!-- Tab content -->
        <div class="flex-1 overflow-auto">
          <ResponseBody 
            v-if="activeResponseTab === 'body'"
            :body="responseState.response.body"
            :contentType="responseState.response.headers['Content-Type'] || responseState.response.headers['content-type'] || ''"
          />
          <ResponseHeaders 
            v-else-if="activeResponseTab === 'headers'"
            :headers="responseState.response.headers"
          />
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  PaperAirplaneIcon,
  NoSymbolIcon,
  ClockIcon,
  ExclamationCircleIcon,
  ArrowsUpDownIcon,
  ArrowsRightLeftIcon,
  CheckCircleIcon,
  XCircleIcon
} from '@heroicons/vue/24/outline'
import { useAppStateStore } from '@/stores/appState'
import { useTabsStore } from '@/stores/tabs'
import { useResponseStore } from '@/stores/response'
import ResponseBody from './ResponseBody.vue'
import ResponseHeaders from './ResponseHeaders.vue'

const appState = useAppStateStore()
const tabsStore = useTabsStore()
const responseStore = useResponseStore()

const effectiveTheme = computed(() => appState.effectiveTheme)
const activeTab = computed(() => tabsStore.activeTab)
const activeResponseTab = ref<'body' | 'headers'>('body')

const responseState = computed(() => {
  if (!activeTab.value) return { status: 'idle' as const }
  return responseStore.getResponse(activeTab.value.id)
})

const statusColor = computed(() => {
  if (responseState.value.status !== 'success') return ''
  const code = responseState.value.response.statusCode
  if (code >= 200 && code < 300) return 'text-status-success'
  if (code >= 300 && code < 400) return 'text-status-redirect'
  if (code >= 400 && code < 500) return 'text-status-client-error'
  return 'text-status-server-error'
})

const statusBadgeClass = computed(() => {
  if (responseState.value.status !== 'success') return ''
  const code = responseState.value.response.statusCode
  if (code >= 200 && code < 300) {
    return 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
  }
  if (code >= 300 && code < 400) {
    return 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'
  }
  if (code >= 400 && code < 500) {
    return 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400'
  }
  return 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400'
})

const statusIcon = computed(() => {
  if (responseState.value.status !== 'success') return null
  const code = responseState.value.response.statusCode
  if (code >= 200 && code < 300) return CheckCircleIcon
  if (code >= 400) return XCircleIcon
  return null
})

const responseTabs = [
  { id: 'body' as const, label: 'Body' },
  { id: 'headers' as const, label: 'Headers' },
]

function formatSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}
</script>
