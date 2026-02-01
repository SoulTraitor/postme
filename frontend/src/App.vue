<template>
  <div class="h-screen flex flex-col" :class="effectiveTheme === 'dark' ? 'dark' : ''">
    <!-- Title Bar -->
    <TitleBar />
    
    <!-- Main Content -->
    <div class="flex-1 flex overflow-hidden">
      <!-- Sidebar -->
      <Sidebar v-if="appState.sidebarOpen" />
      
      <!-- Right section with tabs and panels -->
      <div class="flex-1 flex flex-col overflow-hidden min-w-0">
        <!-- Tab Bar -->
        <TabBar />
        
        <!-- Main Panel -->
        <div 
          class="flex-1 flex overflow-hidden"
          :class="[
            appState.layoutDirection === 'vertical' ? 'flex-col' : 'flex-row',
            effectiveTheme === 'dark' ? 'bg-dark-border' : 'bg-light-border'
          ]"
        >
          <!-- Request Panel -->
          <RequestPanel class="flex-1 min-w-0 min-h-0" />
          
          <!-- Resizer -->
          <div 
            class="resizer"
            :class="appState.layoutDirection === 'vertical' ? 'resizer-horizontal' : 'resizer-vertical'"
            @mousedown="startResize"
          />
          
          <!-- Response Panel -->
          <ResponsePanel class="flex-1 min-w-0 min-h-0" />
        </div>
      </div>
    </div>
    
    <!-- Toasts -->
    <ToastContainer />
    
    <!-- Modals -->
    <ModalContainer />
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref, computed, watch } from 'vue'
import { useAppStateStore } from '@/stores/appState'
import { useTabsStore } from '@/stores/tabs'
import { useCollectionStore } from '@/stores/collection'
import { useEnvironmentStore } from '@/stores/environment'
import { useHistoryStore } from '@/stores/history'
import { api } from '@/services/api'
import { emitKeyboardAction } from '@/composables/useKeyboardActions'
import TitleBar from '@/components/TitleBar.vue'
import TabBar from '@/components/tabs/TabBar.vue'
import Sidebar from '@/components/sidebar/Sidebar.vue'
import RequestPanel from '@/components/request/RequestPanel.vue'
import ResponsePanel from '@/components/response/ResponsePanel.vue'
import ToastContainer from '@/components/common/ToastContainer.vue'
import ModalContainer from '@/components/modals/ModalContainer.vue'

const appState = useAppStateStore()
const tabsStore = useTabsStore()
const collectionStore = useCollectionStore()
const environmentStore = useEnvironmentStore()
const historyStore = useHistoryStore()

const effectiveTheme = computed(() => appState.effectiveTheme)

// Sync dark class to <html> element for CSS selectors to work with teleported modals
watch(effectiveTheme, (theme) => {
  if (theme === 'dark') {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}, { immediate: true })

// Resize handling
const isResizing = ref(false)

function startResize(e: MouseEvent) {
  isResizing.value = true
  document.addEventListener('mousemove', onResize)
  document.addEventListener('mouseup', stopResize)
}

function onResize(e: MouseEvent) {
  if (!isResizing.value) return
  
  const container = document.querySelector('.flex-1.flex.overflow-hidden') as HTMLElement
  if (!container) return
  
  const rect = container.getBoundingClientRect()
  
  if (appState.layoutDirection === 'horizontal') {
    const ratio = ((e.clientX - rect.left) / rect.width) * 100
    appState.splitRatio = Math.max(20, Math.min(80, ratio))
  } else {
    const ratio = ((e.clientY - rect.top) / rect.height) * 100
    appState.splitRatio = Math.max(20, Math.min(80, ratio))
  }
}

function stopResize() {
  isResizing.value = false
  document.removeEventListener('mousemove', onResize)
  document.removeEventListener('mouseup', stopResize)
}

// Keyboard shortcuts
function handleKeydown(e: KeyboardEvent) {
  // Ctrl+B - Toggle sidebar
  if (e.ctrlKey && e.key === 'b') {
    e.preventDefault()
    appState.toggleSidebar()
  }
  
  // Ctrl+\ - Toggle layout direction
  if (e.ctrlKey && e.key === '\\') {
    e.preventDefault()
    appState.toggleLayoutDirection()
  }
  
  // Ctrl+T - New tab
  if (e.ctrlKey && e.key === 't') {
    e.preventDefault()
    tabsStore.addTab()
  }
  
  // Ctrl+W - Close tab
  if (e.ctrlKey && e.key === 'w') {
    e.preventDefault()
    if (tabsStore.activeTabId) {
      closeActiveTabWithConfirmation()
    }
  }
  
  // Ctrl+S - Save request
  if (e.ctrlKey && e.key === 's') {
    e.preventDefault()
    emitKeyboardAction('save')
  }
  
  // Ctrl+Enter - Send request
  if (e.ctrlKey && e.key === 'Enter') {
    e.preventDefault()
    emitKeyboardAction('send')
  }
}

// Close active tab with unsaved changes confirmation
async function closeActiveTabWithConfirmation() {
  const tab = tabsStore.activeTab
  if (!tab) return
  
  if (tab.isDirty) {
    const modal = (window as any).$modal
    if (modal) {
      const confirmed = await modal.confirm({
        title: 'Unsaved Changes',
        message: 'This tab has unsaved changes. Are you sure you want to close it?',
        confirmText: 'Close Without Saving',
        danger: true,
      })
      
      if (!confirmed) return
    }
  }
  
  tabsStore.closeTab(tab.id)
}

// Load all data from backend
async function loadData() {
  try {
    // Load app state and sidebar state
    const [state, sidebarStates] = await Promise.all([
      api.getAppState(),
      api.getSidebarState(),
    ])
    appState.loadState(state, sidebarStates)
    
    // Apply system proxy setting to HTTP client
    await api.setUseSystemProxy(state.useSystemProxy)

    // Load collections, environments, history in parallel
    const [tree, envs, globalVars, history] = await Promise.all([
      api.getCollectionTree(),
      api.getEnvironments(),
      api.getGlobalVariables(),
      api.getHistory(),
    ])

    collectionStore.setTree(tree)
    environmentStore.setEnvironments(envs)
    environmentStore.setGlobalVariables(globalVars)
    environmentStore.setActiveEnv(state.activeEnvId)
    historyStore.setHistory(history)

    // Load tab sessions - show UI immediately, load original state in background
    const sessions = await api.getTabSessions()
    if (sessions.length > 0) {
      // Build tabs without original state first (fast UI display)
      const tabs = sessions.map(s => ({
        id: s.tabId,
        requestId: s.requestId,
        title: s.title,
        method: s.method,
        url: s.url,
        headers: s.headers,
        params: s.params,
        body: s.body,
        bodyType: s.bodyType,
        isDirty: s.isDirty,
        isPreview: false,
        originalState: null, // Load in background
      }))

      // Initialize tabs immediately for fast UI
      tabsStore.init(tabs)
      const activeSession = sessions.find(s => s.isActive)
      if (activeSession) {
        tabsStore.setActiveTab(activeSession.tabId)
      }

      // Load original states in background (for dirty detection)
      Promise.all(sessions.map(async (s, index) => {
        if (s.requestId) {
          try {
            const originalRequest = await api.getRequest(s.requestId)
            const tab = tabsStore.getTab(s.tabId)
            if (tab) {
              tab.originalState = {
                method: originalRequest.method,
                url: originalRequest.url,
                headers: [...originalRequest.headers],
                params: [...originalRequest.params],
                body: originalRequest.body,
                bodyType: originalRequest.bodyType,
              }
              // Recalculate dirty state with original data
              tab.isDirty = s.isDirty
            }
          } catch {
            // Request deleted, use session data as original
            const tab = tabsStore.getTab(s.tabId)
            if (tab) {
              tab.originalState = {
                method: s.method,
                url: s.url,
                headers: [...s.headers],
                params: [...s.params],
                body: s.body,
                bodyType: s.bodyType,
              }
            }
          }
        }
      })).catch(err => {
        console.error('Failed to load original states:', err)
      })
    } else {
      tabsStore.init()
    }
  } catch (error) {
    console.error('Failed to load data:', error)
    // Initialize with defaults on error
    tabsStore.init()
  }
}

// Save app state periodically and on changes
let saveTimeout: number | null = null
function debouncedSave() {
  if (saveTimeout) {
    clearTimeout(saveTimeout)
  }
  saveTimeout = window.setTimeout(async () => {
    try {
      await api.updateAppState(appState.getStateForSave())
    } catch (error) {
      console.error('Failed to save app state:', error)
    }
  }, 1000)
}

// Immediate save for critical state changes (window maximize/restore)
async function immediateSave() {
  if (saveTimeout) {
    clearTimeout(saveTimeout)
    saveTimeout = null
  }
  try {
    await api.updateAppState(appState.getStateForSave())
  } catch (error) {
    console.error('Failed to save app state:', error)
  }
}

// Watch for state changes
watch(
  () => [
    appState.sidebarOpen,
    appState.sidebarWidth,
    appState.layoutDirection,
    appState.splitRatio,
    appState.theme,
    appState.activeEnvId,
    appState.requestPanelTab,
  ],
  () => {
    debouncedSave()
  },
  { deep: true }
)

// Watch for active tab changes to sync with sidebar
watch(
  () => tabsStore.activeTabId,
  () => {
    const tab = tabsStore.activeTab
    if (!tab) {
      appState.highlightedRequestId = null
      return
    }
    
    if (!tab.requestId) {
      // New unsaved request, no highlight
      appState.highlightedRequestId = null
      return
    }
    
    // Always update highlighted request id when switching tabs
    appState.highlightedRequestId = tab.requestId
    
    // If auto-locate is enabled, also expand parents
    if (appState.autoLocateSidebar) {
      const path = collectionStore.findRequestPath(tab.requestId)
      if (path) {
        appState.setSidebarItemExpanded('collection', path.collectionId, true)
        if (path.folderId) {
          appState.setSidebarItemExpanded('folder', path.folderId, true)
        }
      }
    }
  },
  { immediate: true }
)

// Save tab sessions periodically and on changes
let tabSaveTimeout: number | null = null
function debouncedSaveTabSessions() {
  if (tabSaveTimeout) {
    clearTimeout(tabSaveTimeout)
  }
  tabSaveTimeout = window.setTimeout(async () => {
    try {
      // Clear existing sessions first
      await api.clearTabSessions()
      
      // Save all current tabs
      for (let i = 0; i < tabsStore.tabs.length; i++) {
        const tab = tabsStore.tabs[i]
        await api.saveTabSession({
          tabId: tab.id,
          requestId: tab.requestId ?? undefined,
          title: tab.title,
          sortOrder: i,
          isActive: tab.id === tabsStore.activeTabId,
          isDirty: tab.isDirty,
          method: tab.method,
          url: tab.url,
          headers: tab.headers,
          params: tab.params,
          body: tab.body,
          bodyType: tab.bodyType,
        })
      }
    } catch (error) {
      console.error('Failed to save tab sessions:', error)
    }
  }, 500)
}

// Watch for tab changes to save sessions
watch(
  () => [tabsStore.tabs, tabsStore.activeTabId],
  () => {
    debouncedSaveTabSessions()
  },
  { deep: true }
)

// Track window state changes
async function updateWindowState(saveImmediately = false) {
  try {
    // @ts-ignore - Wails runtime
    if (window.runtime) {
      // @ts-ignore
      const isMax = await window.runtime.WindowIsMaximised()
      appState.windowMaximized = isMax
      
      if (!isMax) {
        // Only save size/position when not maximized
        // @ts-ignore
        const size = await window.runtime.WindowGetSize()
        // @ts-ignore
        const pos = await window.runtime.WindowGetPosition()
        
        appState.windowWidth = size.w
        appState.windowHeight = size.h
        appState.windowX = pos.x
        appState.windowY = pos.y
      }
      
      if (saveImmediately) {
        await immediateSave()
      } else {
        debouncedSave()
      }
    }
  } catch (error) {
    console.error('Failed to get window state:', error)
  }
}

// Debounced window resize handler
let resizeTimeout: number | null = null
async function handleWindowResize() {
  // Don't update size while maximized
  // @ts-ignore - Wails runtime
  if (window.runtime) {
    try {
      // @ts-ignore
      const isMax = await window.runtime.WindowIsMaximised()
      if (isMax) return
    } catch {
      // Ignore
    }
  }
  
  if (resizeTimeout) {
    clearTimeout(resizeTimeout)
  }
  resizeTimeout = window.setTimeout(() => {
    updateWindowState()
  }, 500)
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
  loadData()
  
  // Set up window state tracking
  // @ts-ignore - Wails runtime
  if (window.runtime && window.runtime.EventsOn) {
    // @ts-ignore
    window.runtime.EventsOn('wails:window-maximised', async () => {
      // Verify actual state
      try {
        // @ts-ignore
        const isMax = await window.runtime.WindowIsMaximised()
        appState.windowMaximized = isMax
      } catch {
        appState.windowMaximized = true
      }
      immediateSave()
    })
    // @ts-ignore
    window.runtime.EventsOn('wails:window-restored', async () => {
      // Don't update maximized state here - let maximised/unmaximised events handle it
      // Just update window position/size
      updateWindowState(true)
    })
    // @ts-ignore
    window.runtime.EventsOn('wails:window-unmaximised', async () => {
      // Verify actual state
      try {
        // @ts-ignore
        const isMax = await window.runtime.WindowIsMaximised()
        appState.windowMaximized = isMax
      } catch {
        appState.windowMaximized = false
      }
      updateWindowState(true)
    })
  }
  
  // Get initial window state (don't save, just sync)
  updateWindowState(false)
  
  // Track resize events
  window.addEventListener('resize', handleWindowResize)
  
  // Save immediately before window closes
  window.addEventListener('beforeunload', () => {
    immediateSave()
  })
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
  window.removeEventListener('resize', handleWindowResize)
  if (saveTimeout) {
    clearTimeout(saveTimeout)
  }
  if (tabSaveTimeout) {
    clearTimeout(tabSaveTimeout)
  }
  if (resizeTimeout) {
    clearTimeout(resizeTimeout)
  }
})
</script>

<style scoped>
.resizer {
  flex-shrink: 0;
  background: transparent;
  transition: background 0.2s;
}

.resizer:hover {
  background: rgb(var(--color-accent));
}

.resizer-vertical {
  width: 4px;
  cursor: col-resize;
}

.resizer-horizontal {
  height: 4px;
  cursor: row-resize;
}

.dark {
  --color-bg-base: 26 26 26;
  --color-bg-surface: 38 38 38;
  --color-bg-elevated: 51 51 51;
  --color-text-primary: 245 245 245;
  --color-text-secondary: 163 163 163;
  --color-border: 64 64 64;
}
</style>
