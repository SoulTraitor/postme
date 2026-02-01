import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Tab, KeyValue } from '@/types'

// Generate unique ID
function generateId(): string {
  return `tab-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
}

// Create a new empty tab
function createEmptyTab(): Tab {
  return {
    id: generateId(),
    requestId: null,
    title: 'Untitled',
    method: 'GET',
    url: '',
    headers: [],
    params: [],
    body: '',
    bodyType: 'none',
    isDirty: false,
    isPreview: false,
    originalState: null,
  }
}

// HTTP methods that don't support request body
const methodsWithoutBody = ['GET', 'HEAD', 'OPTIONS', 'TRACE']

// Check if a tab is dirty compared to its original state
function computeDirty(tab: Tab): boolean {
  if (!tab.originalState) {
    // New tab, dirty if any content was entered
    return !!(tab.url || tab.body || tab.headers.length > 0 || tab.params.length > 0)
  }
  
  const orig = tab.originalState
  
  // Compare method
  if (tab.method !== orig.method) return true
  
  // Compare URL
  if (tab.url !== orig.url) return true
  
  // Compare params
  if (JSON.stringify(tab.params) !== JSON.stringify(orig.params)) return true
  
  // Compare headers
  if (JSON.stringify(tab.headers) !== JSON.stringify(orig.headers)) return true
  
  // Compare body and bodyType
  // Always check regardless of method, since the UI allows editing body for all methods
  if (tab.body !== orig.body) return true
  if (tab.bodyType !== orig.bodyType) return true

  return false
}

export const useTabsStore = defineStore('tabs', () => {
  const tabs = ref<Tab[]>([])
  const activeTabId = ref<string | null>(null)

  const activeTab = computed(() => {
    return tabs.value.find(t => t.id === activeTabId.value) || null
  })

  const activeTabIndex = computed(() => {
    return tabs.value.findIndex(t => t.id === activeTabId.value)
  })

  // Initialize with an empty tab if none exists
  function init(savedTabs?: Tab[]) {
    if (savedTabs && savedTabs.length > 0) {
      tabs.value = savedTabs
      const activeOne = savedTabs.find(t => t.isPreview === false) || savedTabs[0]
      activeTabId.value = activeOne.id
    } else {
      const newTab = createEmptyTab()
      tabs.value = [newTab]
      activeTabId.value = newTab.id
    }
  }

  // Add a new empty tab at the end
  function addTab() {
    const newTab = createEmptyTab()
    tabs.value.push(newTab)
    activeTabId.value = newTab.id
    return newTab
  }

  // Open a request in a new tab or switch to existing
  function openRequest(requestId: number, title: string, method: string, url: string, headers: KeyValue[], params: KeyValue[], body: string, bodyType: string) {
    // Check if already open
    const existing = tabs.value.find(t => t.requestId === requestId)
    if (existing) {
      activeTabId.value = existing.id
      return existing
    }

    const originalState = { method, url, headers: [...headers], params: [...params], body, bodyType }

    // Check if current tab is preview and reuse it
    const previewTab = tabs.value.find(t => t.isPreview)
    if (previewTab) {
      previewTab.requestId = requestId
      previewTab.title = title
      previewTab.method = method
      previewTab.url = url
      previewTab.headers = headers
      previewTab.params = params
      previewTab.body = body
      previewTab.bodyType = bodyType
      previewTab.isDirty = false
      previewTab.isPreview = false
      previewTab.originalState = originalState
      activeTabId.value = previewTab.id
      return previewTab
    }

    // Create new tab
    const newTab: Tab = {
      id: generateId(),
      requestId,
      title,
      method,
      url,
      headers,
      params,
      body,
      bodyType,
      isDirty: false,
      isPreview: false,
      originalState,
    }
    tabs.value.push(newTab)
    activeTabId.value = newTab.id
    return newTab
  }

  // Preview a request (single-click)
  function previewRequest(requestId: number, title: string, method: string, url: string, headers: KeyValue[], params: KeyValue[], body: string, bodyType: string) {
    // Check if already open (pinned or preview)
    const existing = tabs.value.find(t => t.requestId === requestId)
    if (existing) {
      activeTabId.value = existing.id
      return existing
    }

    const originalState = { method, url, headers: [...headers], params: [...params], body, bodyType }

    // Find or create preview tab
    let previewTab = tabs.value.find(t => t.isPreview)
    if (previewTab) {
      previewTab.requestId = requestId
      previewTab.title = title
      previewTab.method = method
      previewTab.url = url
      previewTab.headers = headers
      previewTab.params = params
      previewTab.body = body
      previewTab.bodyType = bodyType
      previewTab.isDirty = false
      previewTab.originalState = originalState
    } else {
      previewTab = {
        id: generateId(),
        requestId,
        title,
        method,
        url,
        headers,
        params,
        body,
        bodyType,
        isDirty: false,
        isPreview: true,
        originalState,
      }
      tabs.value.push(previewTab)
    }
    activeTabId.value = previewTab.id
    return previewTab
  }

  // Convert preview to permanent
  function pinTab(tabId: string) {
    const tab = tabs.value.find(t => t.id === tabId)
    if (tab) {
      tab.isPreview = false
    }
  }

  // Close a tab
  function closeTab(tabId: string) {
    const index = tabs.value.findIndex(t => t.id === tabId)
    if (index === -1) return

    tabs.value.splice(index, 1)

    if (tabs.value.length === 0) {
      // Add a new empty tab
      const newTab = createEmptyTab()
      tabs.value.push(newTab)
      activeTabId.value = newTab.id
    } else if (activeTabId.value === tabId) {
      // Switch to adjacent tab
      const newIndex = Math.min(index, tabs.value.length - 1)
      activeTabId.value = tabs.value[newIndex].id
    }
  }

  // Set active tab
  function setActiveTab(tabId: string) {
    const tab = tabs.value.find(t => t.id === tabId)
    if (tab) {
      activeTabId.value = tabId
    }
  }

  // Update tab content
  function updateTab(tabId: string, updates: Partial<Tab>) {
    const tab = tabs.value.find(t => t.id === tabId)
    if (tab) {
      Object.assign(tab, updates)
      // Recalculate dirty state properly
      tab.isDirty = computeDirty(tab)
    }
  }

  // Mark tab as saved
  function markSaved(tabId: string, requestId: number, title: string) {
    const tab = tabs.value.find(t => t.id === tabId)
    if (tab) {
      tab.requestId = requestId
      tab.title = title
      tab.isDirty = false
      tab.isPreview = false
      // Update original state to match current saved state
      tab.originalState = {
        method: tab.method,
        url: tab.url,
        headers: [...tab.headers],
        params: [...tab.params],
        body: tab.body,
        bodyType: tab.bodyType,
      }
    }
  }

  // Get tab by ID
  function getTab(tabId: string) {
    return tabs.value.find(t => t.id === tabId)
  }

  // Reorder tabs
  function reorderTabs(fromIndex: number, toIndex: number) {
    if (fromIndex === toIndex || fromIndex < 0 || toIndex < 0) return
    if (fromIndex >= tabs.value.length || toIndex >= tabs.value.length) return
    
    const [moved] = tabs.value.splice(fromIndex, 1)
    tabs.value.splice(toIndex, 0, moved)
  }

  // Update tab title when a request is renamed
  function updateTabTitleByRequestId(requestId: number, newTitle: string) {
    const tab = tabs.value.find(t => t.requestId === requestId)
    if (tab) {
      tab.title = newTitle
    }
  }

  // Close tab when a request is deleted
  function closeTabByRequestId(requestId: number) {
    const tab = tabs.value.find(t => t.requestId === requestId)
    if (tab) {
      closeTab(tab.id)
    }
  }

  // Open history item as new unsaved tab (always creates new tab)
  function openHistoryAsNewTab(title: string, method: string, url: string, headers: KeyValue[], params: KeyValue[], body: string, bodyType: string) {
    // Create new tab without requestId (unsaved)
    const newTab: Tab = {
      id: generateId(),
      requestId: null,
      title,
      method,
      url,
      headers,
      params,
      body,
      bodyType,
      isDirty: false,
      isPreview: false,
      originalState: null, // No original state = new unsaved request
    }
    tabs.value.push(newTab)
    activeTabId.value = newTab.id
    return newTab
  }

  // Duplicate an existing tab
  function duplicateTab(tabId: string) {
    const tab = tabs.value.find(t => t.id === tabId)
    if (!tab) return null

    const newTab: Tab = {
      id: generateId(),
      requestId: null, // New tab is unsaved
      title: `${tab.title} (Copy)`,
      method: tab.method,
      url: tab.url,
      headers: JSON.parse(JSON.stringify(tab.headers)),
      params: JSON.parse(JSON.stringify(tab.params)),
      body: tab.body,
      bodyType: tab.bodyType,
      isDirty: false,
      isPreview: false,
      originalState: null, // No original state = new unsaved request
    }
    
    // Insert after the current tab
    const index = tabs.value.findIndex(t => t.id === tabId)
    tabs.value.splice(index + 1, 0, newTab)
    activeTabId.value = newTab.id
    return newTab
  }

  return {
    tabs,
    activeTabId,
    activeTab,
    activeTabIndex,
    init,
    addTab,
    openRequest,
    previewRequest,
    pinTab,
    closeTab,
    setActiveTab,
    updateTab,
    markSaved,
    getTab,
    reorderTabs,
    updateTabTitleByRequestId,
    closeTabByRequestId,
    openHistoryAsNewTab,
    duplicateTab,
  }
})
