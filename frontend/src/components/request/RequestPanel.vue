<template>
  <div 
    class="flex flex-col h-full overflow-hidden"
    :class="effectiveTheme === 'dark' ? 'bg-dark-base' : 'bg-light-base'"
  >
    <!-- URL bar -->
    <div class="p-4 border-b" :class="effectiveTheme === 'dark' ? 'border-dark-border' : 'border-light-border'">
      <div class="flex gap-2">
        <!-- Method selector -->
        <MethodSelect 
          :modelValue="activeTab?.method || 'GET'"
          @update:modelValue="updateMethod"
        />
        
        <!-- URL input -->
        <UrlInput 
          :modelValue="activeTab?.url || ''"
          @update:modelValue="updateUrl"
          class="flex-1"
        />
        
        <!-- Save button -->
        <button
          @click="openSaveModal"
          class="px-4 py-2 rounded-md font-medium transition-colors"
          :class="[
            effectiveTheme === 'dark' 
              ? 'bg-dark-hover text-gray-300 hover:bg-dark-border' 
              : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
          ]"
          :title="activeTab?.requestId ? 'Update request (Ctrl+S)' : 'Save request (Ctrl+S)'"
        >
          {{ activeTab?.requestId ? 'Update' : 'Save' }}
        </button>
        
        <!-- Send button -->
        <button
          v-if="!isLoading"
          @click="sendRequest"
          class="px-6 py-2 rounded-md font-medium text-white bg-accent hover:bg-accent-hover transition-colors"
          :disabled="!activeTab?.url"
        >
          Send
        </button>
        <button
          v-else
          @click="cancelRequest"
          class="px-6 py-2 rounded-md font-medium text-white bg-red-500 hover:bg-red-600 transition-colors"
        >
          Cancel
        </button>
      </div>
    </div>
    
    <!-- Request details tabs -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- Tab headers -->
      <div class="flex border-b" :class="effectiveTheme === 'dark' ? 'border-dark-border' : 'border-light-border'">
        <button
          v-for="tab in requestTabs"
          :key="tab.id"
          @click="activeRequestTab = tab.id"
          class="px-4 py-2 text-sm font-medium transition-colors relative"
          :class="[
            activeRequestTab === tab.id
              ? 'text-accent'
              : (effectiveTheme === 'dark' ? 'text-gray-400 hover:text-white' : 'text-gray-500 hover:text-gray-900')
          ]"
        >
          {{ tab.label }}
          <span v-if="tab.count > 0" class="ml-1 text-xs text-gray-500">({{ tab.count }})</span>
          <div
            v-if="activeRequestTab === tab.id"
            class="absolute bottom-0 left-0 right-0 h-0.5 bg-accent"
          />
        </button>
      </div>
      
      <!-- Tab content -->
      <div class="flex-1 overflow-auto p-4">
        <ParamsEditor 
          v-if="activeRequestTab === 'params'"
          :params="activeTab?.params || []"
          @update:params="updateParams"
        />
        <HeadersEditor 
          v-else-if="activeRequestTab === 'headers'"
          :headers="activeTab?.headers || []"
          @update:headers="updateHeaders"
        />
        <BodyEditor 
          v-else-if="activeRequestTab === 'body'"
          :body="activeTab?.body || ''"
          :bodyType="activeTab?.bodyType || 'none'"
          @update:body="updateBody"
          @update:bodyType="updateBodyType"
        />
      </div>
    </div>
    
    <!-- Save Request Modal -->
    <SaveRequestModal 
      :isOpen="saveModalOpen" 
      :tabId="activeTab?.id"
      @close="saveModalOpen = false"
      @saved="onRequestSaved"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useAppStateStore } from '@/stores/appState'
import { useTabsStore } from '@/stores/tabs'
import { useResponseStore } from '@/stores/response'
import { useEnvironmentStore } from '@/stores/environment'
import { useCollectionStore } from '@/stores/collection'
import { useHistoryStore } from '@/stores/history'
import { api } from '@/services/api'
import { onKeyboardAction } from '@/composables/useKeyboardActions'
import type { KeyValue } from '@/types'
import MethodSelect from './MethodSelect.vue'
import UrlInput from './UrlInput.vue'
import ParamsEditor from './ParamsEditor.vue'
import HeadersEditor from './HeadersEditor.vue'
import BodyEditor from './BodyEditor.vue'
import SaveRequestModal from '@/components/modals/SaveRequestModal.vue'

const appState = useAppStateStore()
const tabsStore = useTabsStore()
const responseStore = useResponseStore()
const environmentStore = useEnvironmentStore()
const collectionStore = useCollectionStore()
const historyStore = useHistoryStore()

const effectiveTheme = computed(() => appState.effectiveTheme)
const activeTab = computed(() => tabsStore.activeTab)
const activeRequestTab = computed({
  get: () => appState.requestPanelTab,
  set: (value) => { appState.requestPanelTab = value }
})
const saveModalOpen = ref(false)

const isLoading = computed(() => {
  if (!activeTab.value) return false
  const state = responseStore.getResponse(activeTab.value.id)
  return state.status === 'loading'
})

const requestTabs = computed(() => [
  { id: 'params' as const, label: 'Params', count: activeTab.value?.params.filter(p => p.enabled).length || 0 },
  { id: 'headers' as const, label: 'Headers', count: activeTab.value?.headers.filter(h => h.enabled).length || 0 },
  { id: 'body' as const, label: 'Body', count: activeTab.value?.body ? 1 : 0 },
])

function updateMethod(method: string) {
  if (activeTab.value) {
    tabsStore.updateTab(activeTab.value.id, { method })
  }
}

function updateUrl(url: string) {
  if (activeTab.value) {
    tabsStore.updateTab(activeTab.value.id, { url })
  }
}

function updateParams(params: KeyValue[]) {
  if (activeTab.value) {
    tabsStore.updateTab(activeTab.value.id, { params })
  }
}

function updateHeaders(headers: KeyValue[]) {
  if (activeTab.value) {
    tabsStore.updateTab(activeTab.value.id, { headers })
  }
}

function updateBody(body: string) {
  if (activeTab.value) {
    tabsStore.updateTab(activeTab.value.id, { body })
  }
}

function updateBodyType(bodyType: string) {
  if (activeTab.value) {
    const tab = activeTab.value
    const currentHeaders = [...tab.headers]
    
    // Determine the Content-Type for this body type
    let contentType: string | null = null
    if (bodyType === 'json') {
      contentType = 'application/json'
    } else if (bodyType === 'xml') {
      contentType = 'application/xml'
    } else if (bodyType === 'text') {
      contentType = 'text/plain'
    } else if (bodyType === 'form-data') {
      contentType = 'multipart/form-data'
    } else if (bodyType === 'x-www-form-urlencoded') {
      contentType = 'application/x-www-form-urlencoded'
    }
    
    // Find existing Content-Type header (case-insensitive)
    const contentTypeIndex = currentHeaders.findIndex(
      h => h.key.toLowerCase() === 'content-type'
    )
    
    if (bodyType === 'none') {
      // Remove Content-Type if body is none
      if (contentTypeIndex !== -1) {
        currentHeaders.splice(contentTypeIndex, 1)
      }
    } else if (contentType) {
      if (contentTypeIndex !== -1) {
        // Update existing Content-Type header
        currentHeaders[contentTypeIndex] = {
          ...currentHeaders[contentTypeIndex],
          value: contentType,
          enabled: true,
        }
      } else {
        // Add new Content-Type header at the beginning
        currentHeaders.unshift({
          key: 'Content-Type',
          value: contentType,
          enabled: true,
        })
      }
    }
    
    // Don't clear body - preserve it to avoid accidental data loss
    tabsStore.updateTab(tab.id, { bodyType, headers: currentHeaders })
  }
}

async function sendRequest() {
  if (!activeTab.value?.url) return
  
  const tab = activeTab.value
  responseStore.setLoading(tab.id)
  
  try {
    // Build URL with params and replace variables
    let url = environmentStore.replaceVariables(tab.url)
    const urlObj = new URL(url)
    for (const param of tab.params) {
      if (param.enabled && param.key) {
        urlObj.searchParams.set(
          environmentStore.replaceVariables(param.key),
          environmentStore.replaceVariables(param.value)
        )
      }
    }
    
    // Replace variables in headers
    const headers: KeyValue[] = tab.headers
      .filter(h => h.enabled && h.key)
      .map(h => ({
        key: environmentStore.replaceVariables(h.key),
        value: environmentStore.replaceVariables(h.value),
        enabled: true,
      }))
    
    // Replace variables in body
    const body = environmentStore.replaceVariables(tab.body)
    
    // Execute request via Wails backend
    const response = await api.executeRequest({
      tabId: tab.id,
      method: tab.method,
      url: urlObj.toString(),
      headers,
      body,
      bodyType: tab.bodyType,
      timeout: appState.requestTimeout,
    })
    
    responseStore.setSuccess(tab.id, response)
    
    // Save to history
    try {
      const historyItem = await api.addHistory({
        requestId: tab.requestId ?? undefined,
        method: tab.method,
        url: urlObj.toString(),
        requestHeaders: JSON.stringify(headers),
        requestBody: body,
        statusCode: response.statusCode,
        responseHeaders: JSON.stringify(response.headers),
        responseBody: response.body,
        durationMs: response.duration,
      })
      historyStore.addHistory(historyItem)
    } catch (err) {
      console.error('Failed to save history:', err)
    }
  } catch (error: any) {
    const errorMessage = error?.message || String(error) || 'Request failed'
    if (errorMessage.includes('cancelled') || errorMessage.includes('canceled')) {
      responseStore.setCancelled(tab.id)
    } else if (errorMessage.includes('timeout')) {
      responseStore.setTimeout(tab.id, appState.requestTimeout)
    } else {
      responseStore.setError(tab.id, errorMessage)
    }
  }
}

async function cancelRequest() {
  if (activeTab.value) {
    try {
      await api.cancelRequest(activeTab.value.id)
    } catch (error) {
      console.error('Failed to cancel request:', error)
    }
    responseStore.setCancelled(activeTab.value.id)
  }
}

async function openSaveModal() {
  if (!activeTab.value) return
  
  // If this is an existing request, just update it
  if (activeTab.value.requestId) {
    await updateExistingRequest()
  } else {
    // Open save modal for new requests
    saveModalOpen.value = true
  }
}

async function updateExistingRequest() {
  const tab = activeTab.value
  if (!tab || !tab.requestId) return
  
  try {
    const existingRequest = collectionStore.getRequest(tab.requestId)
    if (!existingRequest) return
    
    const updatedRequest = {
      ...existingRequest,
      method: tab.method,
      url: tab.url,
      headers: tab.headers,
      params: tab.params,
      body: tab.body,
      bodyType: tab.bodyType,
    }
    
    await api.updateRequest(updatedRequest)
    collectionStore.updateRequest(updatedRequest)
    // Use markSaved to properly update originalState and clear dirty flag
    tabsStore.markSaved(tab.id, tab.requestId, tab.title)
  } catch (error) {
    console.error('Failed to update request:', error)
  }
}

function onRequestSaved(_requestId: number) {
  // Request was saved successfully, modal will close
}

// Keyboard action listeners
let unsubSave: (() => void) | null = null
let unsubSend: (() => void) | null = null

onMounted(() => {
  unsubSave = onKeyboardAction('save', openSaveModal)
  unsubSend = onKeyboardAction('send', sendRequest)
})

onUnmounted(() => {
  unsubSave?.()
  unsubSend?.()
})
</script>
