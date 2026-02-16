import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import type { AppState, SidebarState } from '@/types'
import { WindowSetBackgroundColour } from '../../wailsjs/runtime/runtime'

// Import api lazily to avoid circular dependency
let apiModule: typeof import('@/services/api') | null = null
async function getApi() {
  if (!apiModule) {
    apiModule = await import('@/services/api')
  }
  return apiModule.api
}

export const useAppStateStore = defineStore('appState', () => {
  // App state
  const sidebarOpen = ref(true)
  const sidebarWidth = ref(260)
  const layoutDirection = ref<'horizontal' | 'vertical'>('horizontal')
  const splitRatio = ref(50)
  const theme = ref<'light' | 'dark' | 'system'>('system')
  const activeEnvId = ref<number | null>(null)
  const requestTimeout = ref(30)
  const autoLocateSidebar = ref(true)
  const useSystemProxy = ref(true)
  const requestPanelTab = ref<'params' | 'headers' | 'body'>('params')
  const modalOpenCount = ref(0)
  
  // Window state
  const windowWidth = ref(1200)
  const windowHeight = ref(800)
  const windowX = ref<number | null>(null)
  const windowY = ref<number | null>(null)
  const windowMaximized = ref(false)
  
  // Sidebar expanded states
  const sidebarStates = ref<Map<string, boolean>>(new Map())
  
  // Sidebar active tab
  const sidebarTab = ref<'collections' | 'history'>('collections')
  
  // Currently highlighted request in sidebar (for tab-sidebar sync)
  const highlightedRequestId = ref<number | null>(null)
  
  // Loading state
  const loading = ref(true)

  // Track system dark mode preference reactively
  const systemIsDark = ref(
    typeof window !== 'undefined'
      ? window.matchMedia('(prefers-color-scheme: dark)').matches
      : false
  )

  // Current effective theme
  const effectiveTheme = computed(() => {
    if (theme.value === 'system') {
      return systemIsDark.value ? 'dark' : 'light'
    }
    return theme.value
  })

  const isModalOpen = computed(() => modalOpenCount.value > 0)

  // Apply theme to document and window background
  watch(effectiveTheme, (newTheme) => {
    if (newTheme === 'dark') {
      document.documentElement.classList.add('dark')
      WindowSetBackgroundColour(26, 26, 26, 255)
    } else {
      document.documentElement.classList.remove('dark')
      WindowSetBackgroundColour(255, 255, 255, 255)
    }
  }, { immediate: true })

  // Listen to system theme changes — update reactive ref
  if (typeof window !== 'undefined') {
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
      systemIsDark.value = e.matches
    })
  }

  // Toggle theme (light <-> dark only, for title bar button)
  function toggleTheme() {
    // Use effectiveTheme to determine current visual state
    theme.value = effectiveTheme.value === 'dark' ? 'light' : 'dark'
  }

  // Cycle theme (includes system, for settings)
  function cycleTheme() {
    if (theme.value === 'light') {
      theme.value = 'dark'
    } else if (theme.value === 'dark') {
      theme.value = 'system'
    } else {
      theme.value = 'light'
    }
  }

  // Toggle sidebar
  function toggleSidebar() {
    sidebarOpen.value = !sidebarOpen.value
  }

  // Toggle layout direction
  function toggleLayoutDirection() {
    layoutDirection.value = layoutDirection.value === 'horizontal' ? 'vertical' : 'horizontal'
  }

  function addModalOpen() {
    modalOpenCount.value += 1
  }

  function removeModalOpen() {
    modalOpenCount.value = Math.max(0, modalOpenCount.value - 1)
  }

  // Get sidebar item expanded state
  function isSidebarItemExpanded(itemType: string, itemId: number): boolean {
    return sidebarStates.value.get(`${itemType}-${itemId}`) ?? false
  }

  // Set sidebar item expanded state
  function setSidebarItemExpanded(itemType: string, itemId: number, expanded: boolean) {
    sidebarStates.value.set(`${itemType}-${itemId}`, expanded)
  }

  // Toggle sidebar item expanded state
  async function toggleSidebarItem(itemType: string, itemId: number) {
    const key = `${itemType}-${itemId}`
    const newValue = !sidebarStates.value.get(key)
    sidebarStates.value.set(key, newValue)
    
    // Persist to backend
    try {
      const api = await getApi()
      await api.setSidebarItemExpanded(itemType, itemId, newValue)
    } catch (error) {
      console.error('Failed to save sidebar state:', error)
    }
  }

  // Expand a sidebar item (convenience wrapper)
  async function expandSidebarItem(itemType: string, itemId: number) {
    const key = `${itemType}-${itemId}`
    if (!sidebarStates.value.get(key)) {
      sidebarStates.value.set(key, true)
      try {
        const api = await getApi()
        await api.setSidebarItemExpanded(itemType, itemId, true)
      } catch (error) {
        console.error('Failed to save sidebar state:', error)
      }
    }
  }

  // Load state from database
  function loadState(state: AppState, sidebarStateList: SidebarState[]) {
    sidebarOpen.value = state.sidebarOpen
    sidebarWidth.value = state.sidebarWidth
    layoutDirection.value = state.layoutDirection
    splitRatio.value = state.splitRatio
    theme.value = state.theme
    activeEnvId.value = state.activeEnvId
    requestTimeout.value = state.requestTimeout
    autoLocateSidebar.value = state.autoLocateSidebar
    useSystemProxy.value = state.useSystemProxy
    requestPanelTab.value = (state.requestPanelTab as 'params' | 'headers' | 'body') || 'params'
    
    // Load window state
    windowWidth.value = state.windowWidth || 1200
    windowHeight.value = state.windowHeight || 800
    windowX.value = state.windowX ?? null
    windowY.value = state.windowY ?? null
    windowMaximized.value = state.windowMaximized || false

    for (const item of sidebarStateList) {
      sidebarStates.value.set(`${item.itemType}-${item.itemId}`, item.expanded)
    }

    loading.value = false
  }

  // Get current state for saving
  function getStateForSave(): Partial<AppState> {
    return {
      sidebarOpen: sidebarOpen.value,
      sidebarWidth: sidebarWidth.value,
      layoutDirection: layoutDirection.value,
      splitRatio: splitRatio.value,
      theme: theme.value,
      activeEnvId: activeEnvId.value,
      requestTimeout: requestTimeout.value,
      autoLocateSidebar: autoLocateSidebar.value,
      useSystemProxy: useSystemProxy.value,
      requestPanelTab: requestPanelTab.value,
      windowWidth: windowWidth.value,
      windowHeight: windowHeight.value,
      windowX: windowX.value,
      windowY: windowY.value,
      windowMaximized: windowMaximized.value,
    }
  }

  return {
    sidebarOpen,
    sidebarWidth,
    layoutDirection,
    splitRatio,
    theme,
    activeEnvId,
    requestTimeout,
    autoLocateSidebar,
    useSystemProxy,
    requestPanelTab,
    modalOpenCount,
    isModalOpen,
    addModalOpen,
    removeModalOpen,
    windowWidth,
    windowHeight,
    windowX,
    windowY,
    windowMaximized,
    sidebarStates,
    sidebarTab,
    highlightedRequestId,
    loading,
    effectiveTheme,
    toggleTheme,
    cycleTheme,
    toggleSidebar,
    toggleLayoutDirection,
    isSidebarItemExpanded,
    setSidebarItemExpanded,
    toggleSidebarItem,
    expandSidebarItem,
    loadState,
    getStateForSave,
  }
})
