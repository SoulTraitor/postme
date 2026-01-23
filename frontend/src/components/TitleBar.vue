<template>
  <div 
    class="h-10 flex items-center justify-between px-3 border-b select-none wails-drag"
    :class="[
      effectiveTheme === 'dark' 
        ? 'bg-dark-surface border-dark-border text-white' 
        : 'bg-light-surface border-light-border text-gray-900'
    ]"
  >
    <!-- Left section -->
    <div class="flex items-center gap-2">
      <!-- Sidebar toggle -->
      <button
        @click="appState.toggleSidebar"
        class="p-1.5 rounded-md transition-colors"
        :class="effectiveTheme === 'dark' ? 'hover:bg-dark-hover' : 'hover:bg-light-hover'"
        title="Toggle sidebar (Ctrl+B)"
      >
        <Bars3Icon class="w-5 h-5" />
      </button>
      
      <!-- App name -->
      <span class="font-semibold text-accent">PostMe</span>
    </div>
    
    <!-- Right section -->
    <div class="flex items-center gap-2">
      <!-- Environment selector -->
      <EnvironmentSelector />
      
      <!-- Theme toggle (light/dark only, system available in Settings) -->
      <button
        @click="appState.toggleTheme"
        class="p-1.5 rounded-md transition-colors"
        :class="effectiveTheme === 'dark' ? 'hover:bg-dark-hover' : 'hover:bg-light-hover'"
        :title="themeTooltip"
      >
        <SunIcon v-if="effectiveTheme === 'light'" class="w-5 h-5" />
        <MoonIcon v-else class="w-5 h-5" />
      </button>
      
      <!-- Settings -->
      <button
        @click="openSettings"
        class="p-1.5 rounded-md transition-colors"
        :class="effectiveTheme === 'dark' ? 'hover:bg-dark-hover' : 'hover:bg-light-hover'"
        title="Settings"
      >
        <Cog6ToothIcon class="w-5 h-5" />
      </button>
      
      <!-- Window controls (for custom title bar) -->
      <div class="flex items-center ml-4 -mr-3">
        <button
          @click="minimizeWindow"
          class="w-10 h-10 flex items-center justify-center hover:bg-gray-500/20"
        >
          <MinusIcon class="w-4 h-4" />
        </button>
        <button
          @click="maximizeWindow"
          class="w-10 h-10 flex items-center justify-center hover:bg-gray-500/20"
          :title="isMaximized ? 'Restore' : 'Maximize'"
        >
          <Square2StackIcon v-if="isMaximized" class="w-3 h-3" />
          <StopIcon v-else class="w-3 h-3" />
        </button>
        <button
          @click="closeWindow"
          class="w-10 h-10 flex items-center justify-center hover:bg-red-500"
        >
          <XMarkIcon class="w-4 h-4" />
        </button>
      </div>
    </div>
    
    <!-- Settings Modal -->
    <SettingsModal :isOpen="settingsOpen" @close="closeSettings" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useAppStateStore } from '@/stores/appState'
import { 
  Bars3Icon, 
  SunIcon, 
  MoonIcon, 
  Cog6ToothIcon,
  MinusIcon,
  StopIcon,
  XMarkIcon,
  Square2StackIcon,
} from '@heroicons/vue/24/outline'
import EnvironmentSelector from '@/components/EnvironmentSelector.vue'
import SettingsModal from '@/components/modals/SettingsModal.vue'

const appState = useAppStateStore()

const effectiveTheme = computed(() => appState.effectiveTheme)
const settingsOpen = ref(false)
const isMaximized = ref(false)

const themeTooltip = computed(() => {
  return appState.effectiveTheme === 'light' 
    ? 'Switch to dark mode' 
    : 'Switch to light mode'
})

function openSettings() {
  settingsOpen.value = true
}

function closeSettings() {
  settingsOpen.value = false
}

function minimizeWindow() {
  // @ts-ignore - Wails runtime
  if (window.runtime) {
    window.runtime.WindowMinimise()
  }
}

function maximizeWindow() {
  // @ts-ignore - Wails runtime
  if (window.runtime) {
    window.runtime.WindowToggleMaximise()
  }
}

function closeWindow() {
  // @ts-ignore - Wails runtime
  if (window.runtime) {
    window.runtime.Quit()
  }
}

// Listen for window state changes
onMounted(() => {
  // @ts-ignore - Wails runtime events
  if (window.runtime && window.runtime.EventsOn) {
    // @ts-ignore
    window.runtime.EventsOn('wails:window-maximised', () => {
      isMaximized.value = true
    })
    // @ts-ignore
    window.runtime.EventsOn('wails:window-restored', () => {
      isMaximized.value = false
    })
    // @ts-ignore
    window.runtime.EventsOn('wails:window-unmaximised', () => {
      isMaximized.value = false
    })
  }
})

onUnmounted(() => {
  // @ts-ignore - Wails runtime events
  if (window.runtime && window.runtime.EventsOff) {
    // @ts-ignore
    window.runtime.EventsOff('wails:window-maximised')
    // @ts-ignore
    window.runtime.EventsOff('wails:window-restored')
    // @ts-ignore
    window.runtime.EventsOff('wails:window-unmaximised')
  }
})
</script>
