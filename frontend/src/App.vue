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
          <RequestPanel
            class="min-w-0 min-h-0"
            :style="{
              flex: `0 0 ${appState.splitRatio}%`
            }"
          />

          <!-- Resizer -->
          <div
            class="resizer relative group"
            :class="[
              appState.layoutDirection === 'vertical' ? 'resizer-horizontal' : 'resizer-vertical',
              isResizing ? 'resizing' : ''
            ]"
            @mousedown="startResize"
            @dblclick="resetSplitRatio"
            :title="isResizing ? `${splitRatioPercent}%` : 'Double-click to reset'"
          >
            <!-- Ratio tooltip during drag -->
            <div
              v-if="isResizing"
              class="absolute z-10 px-2 py-1 bg-accent text-white text-xs rounded-md pointer-events-none whitespace-nowrap"
              :class="appState.layoutDirection === 'vertical' ? '-top-8 left-1/2 -translate-x-1/2' : '-left-16 top-1/2 -translate-y-1/2'"
            >
              {{ splitRatioPercent }}%
            </div>
          </div>

          <!-- Response Panel -->
          <ResponsePanel
            class="min-w-0 min-h-0"
            :style="{
              flex: `1 1 ${100 - appState.splitRatio}%`
            }"
          />
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
const splitRatioPercent = computed(() => Math.round(appState.splitRatio))
let startX = 0
let startY = 0
let startRatio = 50

function startResize(e: MouseEvent) {
  isResizing.value = true
  startX = e.clientX
  startY = e.clientY
  startRatio = appState.splitRatio
  document.addEventListener('mousemove', onResize)
  document.addEventListener('mouseup', stopResize)
  e.preventDefault()
}

function resetSplitRatio() {
  appState.splitRatio = 50
}

function onResize(e: MouseEvent) {
  if (!isResizing.value) return

  const container = document.querySelector('.flex-1.flex.overflow-hidden') as HTMLElement
  if (!container) return

  const rect = container.getBoundingClientRect()

  if (appState.layoutDirection === 'horizontal') {
    // Calculate delta from start position
    const deltaX = e.clientX - startX
    const deltaRatio = (deltaX / rect.width) * 100
    const newRatio = startRatio + deltaRatio
    appState.splitRatio = Math.max(20, Math.min(80, newRatio))
  } else {
    // Calculate delta from start position
    const deltaY = e.clientY - startY
    const deltaRatio = (deltaY / rect.height) * 100
    const newRatio = startRatio + deltaRatio
    appState.splitRatio = Math.max(20, Math.min(80, newRatio))
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
  const startTime = performance.now()
  console.log('[Performance] loadData start')

  try {
    // Priority 1: Load app state and tabs first (critical for UI)
    const t1 = performance.now()
    const [state, sidebarStates, sessions] = await Promise.all([
      api.getAppState(),
      api.getSidebarState(),
      api.getTabSessions(),
    ])
    console.log(`[Performance] Priority 1 loaded in ${(performance.now() - t1).toFixed(2)}ms - sessions count: ${sessions.length}`)
    appState.loadState(state, sidebarStates)

    // Priority 2: Initialize tabs immediately for instant UI
    const t2 = performance.now()
    if (sessions.length > 0) {
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
        originalState: null,
      }))
      tabsStore.init(tabs)
      const activeSession = sessions.find(s => s.isActive)
      if (activeSession) {
        tabsStore.setActiveTab(activeSession.tabId)
      }
      console.log(`[Performance] Priority 2 tabs initialized in ${(performance.now() - t2).toFixed(2)}ms - ${tabs.length} tabs`)
    } else {
      tabsStore.init()
      console.log('[Performance] No saved tabs, created empty tab')
    }

    // Priority 3: Load other data in background (non-blocking, sequential to avoid DB locks)
    const t3 = performance.now()
    ;(async () => {
      try {
        // Apply proxy setting first (non-DB operation)
        await api.setUseSystemProxy(state.useSystemProxy)

        // Load critical UI data (collections for sidebar)
        const tree = await api.getCollectionTree()
        collectionStore.setTree(tree)
        console.log(`[Performance] Collections loaded in ${(performance.now() - t3).toFixed(2)}ms`)

        // Load environments (needed for requests)
        const [envs, globalVars] = await Promise.all([
          api.getEnvironments(),
          api.getGlobalVariables(),
        ])
        environmentStore.setEnvironments(envs)
        environmentStore.setGlobalVariables(globalVars)
        environmentStore.setActiveEnv(state.activeEnvId)
        console.log(`[Performance] Environments loaded in ${(performance.now() - t3).toFixed(2)}ms`)

        // Load history last (least important)
        const history = await api.getHistory()
        historyStore.setHistory(history)
        console.log(`[Performance] All background data loaded in ${(performance.now() - t3).toFixed(2)}ms`)
      } catch (err) {
        console.error('Failed to load background data:', err)
      }
    })()

    console.log(`[Performance] loadData completed in ${(performance.now() - startTime).toFixed(2)}ms (tabs should be visible now)`)
    console.log('[Performance] Background data loading...')

    // Priority 4: Load original states in background (for dirty detection)
    if (sessions.length > 0) {
      Promise.all(sessions.map(async (s) => {
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
    } catch (error: any) {
      // Retry once if database is locked
      if (error?.message?.includes('locked') || error?.message?.includes('BUSY')) {
        console.warn('Database locked, retrying in 500ms...')
        setTimeout(async () => {
          try {
            await api.updateAppState(appState.getStateForSave())
          } catch (retryError) {
            console.error('Failed to save app state after retry:', retryError)
          }
        }, 500)
      } else {
        console.error('Failed to save app state:', error)
      }
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
      // Build all sessions first
      const sessions = tabsStore.tabs.map((tab, i) => ({
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
      }))

      // Clear and save atomically - if save fails, clear won't happen
      await api.clearTabSessions()

      // Save all sessions - if any fails, log but continue
      for (const session of sessions) {
        try {
          await api.saveTabSession(session)
        } catch (err) {
          console.error('Failed to save individual tab session:', err, session)
          // Continue saving other tabs
        }
      }

      console.log(`[Tabs] Saved ${sessions.length} tab sessions`)
    } catch (error) {
      console.error('Failed to save tab sessions:', error)
      // Critical: if clear fails, tabs might be lost on next restart
      console.error('[CRITICAL] Tab sessions may be corrupted!')
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
      // Wait for animation to complete
      setTimeout(async () => {
        try {
          // @ts-ignore
          const isMax = await window.runtime.WindowIsMaximised()
          console.log('[Window] Maximised event - actual state:', isMax)
          appState.windowMaximized = isMax
          immediateSave()
        } catch (err) {
          console.error('[Window] Error in maximised event:', err)
          appState.windowMaximized = true
          immediateSave()
        }
      }, 300)
    })
    // @ts-ignore
    window.runtime.EventsOn('wails:window-restored', async () => {
      // Wait longer for window to fully settle
      setTimeout(async () => {
        try {
          // @ts-ignore
          const isMax = await window.runtime.WindowIsMaximised()
          console.log('[Window] Restored event - actual state:', isMax)
          appState.windowMaximized = isMax
          updateWindowState(true)
        } catch (err) {
          console.error('[Window] Error in restored event:', err)
        }
      }, 400)
    })
    // @ts-ignore
    window.runtime.EventsOn('wails:window-unmaximised', async () => {
      // Wait for window to settle before saving position
      setTimeout(async () => {
        try {
          // @ts-ignore
          const isMax = await window.runtime.WindowIsMaximised()
          console.log('[Window] Unmaximised event - actual state:', isMax)
          appState.windowMaximized = isMax
          updateWindowState(true)
        } catch (err) {
          console.error('[Window] Error in unmaximised event:', err)
        }
      }, 300)
    })

    // Backup: Check state when window gains focus (catches missed events)
    window.addEventListener('focus', async () => {
      try {
        // @ts-ignore
        if (window.runtime && window.runtime.WindowIsMaximised) {
          // @ts-ignore
          const isMax = await window.runtime.WindowIsMaximised()
          if (appState.windowMaximized !== isMax) {
            console.warn('[Window] Focus check - state mismatch corrected:', { was: appState.windowMaximized, now: isMax })
            appState.windowMaximized = isMax
          }
        }
      } catch {
        // Ignore errors
      }
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
  transition: background 0.2s, box-shadow 0.2s;
  position: relative;
}

.resizer:hover {
  background: #d97706;
}

.resizer.resizing {
  background: #d97706;
}

.resizer.resizing.resizer-vertical {
  box-shadow: 0 0 12px 2px rgba(217, 119, 6, 0.5);
}

.resizer.resizing.resizer-horizontal {
  box-shadow: 0 0 12px 2px rgba(217, 119, 6, 0.5);
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
