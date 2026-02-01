<template>
  <div 
    class="flex flex-col h-full overflow-hidden"
    :class="effectiveTheme === 'dark' ? 'bg-dark-base' : 'bg-light-base'"
  >
    <!-- Response state: Idle -->
    <div v-if="responseState.status === 'idle'" class="flex-1 flex flex-col items-center justify-center text-gray-500">
      <PaperAirplaneIcon class="w-12 h-12 mb-4 opacity-50" />
      <p>Click Send to make a request</p>
    </div>
    
    <!-- Response state: Loading -->
    <div v-else-if="responseState.status === 'loading'" class="flex-1 flex flex-col items-center justify-center text-gray-500">
      <div class="w-8 h-8 border-2 border-accent border-t-transparent rounded-full animate-spin mb-4" />
      <p>Sending request...</p>
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
        <span 
          class="font-medium"
          :class="statusColor"
        >
          {{ responseState.response.statusCode }} {{ responseState.response.status.split(' ').slice(1).join(' ') }}
        </span>
        
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
  ArrowsRightLeftIcon
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
