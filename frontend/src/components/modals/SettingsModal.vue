<template>
  <TransitionRoot :show="isOpen" as="template">
    <Dialog as="div" class="relative z-50" @close="close">
      <TransitionChild
        enter="ease-out duration-300"
        enter-from="opacity-0"
        enter-to="opacity-100"
        leave="ease-in duration-200"
        leave-from="opacity-100"
        leave-to="opacity-0"
        as="template"
      >
        <div
          class="fixed inset-0 modal-backdrop backdrop-blur-sm"
          style="transition: opacity 300ms ease-out, backdrop-filter 300ms ease-out"
          aria-hidden="true"
        />
      </TransitionChild>

      <div class="fixed inset-0 overflow-y-auto">
        <div class="flex min-h-full items-center justify-center p-4">
          <TransitionChild
            enter="ease-out duration-300"
            enter-from="opacity-0 scale-90"
            enter-to="opacity-100 scale-100"
            leave="ease-in duration-200"
            leave-from="opacity-100 scale-100"
            leave-to="opacity-0 scale-95"
          >
            <DialogPanel 
              class="w-full max-w-lg rounded-lg p-6 shadow-xl"
              :class="effectiveTheme === 'dark' ? 'bg-dark-elevated' : 'bg-white'"
            >
              <DialogTitle 
                class="text-lg font-medium mb-6"
                :class="effectiveTheme === 'dark' ? 'text-white' : 'text-gray-900'"
              >
                Settings
              </DialogTitle>
              
              <div class="space-y-6">
                <!-- Request Timeout -->
                <div>
                  <label 
                    class="block text-sm font-medium mb-2"
                    :class="effectiveTheme === 'dark' ? 'text-gray-300' : 'text-gray-700'"
                  >
                    Request Timeout (seconds)
                  </label>
                  <input
                    v-model.number="localSettings.requestTimeout"
                    type="number"
                    min="1"
                    max="300"
                    class="w-full px-3 py-2 rounded-md border outline-none text-sm"
                    :class="[
                      effectiveTheme === 'dark'
                        ? 'bg-dark-surface border-dark-border text-white focus:border-accent'
                        : 'bg-white border-light-border text-gray-900 focus:border-accent'
                    ]"
                  />
                  <p class="mt-1 text-xs text-gray-500">
                    Maximum time to wait for a response (1-300 seconds)
                  </p>
                </div>
                
                <!-- Auto-locate in Sidebar -->
                <div class="flex items-center justify-between">
                  <div>
                    <label 
                      class="block text-sm font-medium"
                      :class="effectiveTheme === 'dark' ? 'text-gray-300' : 'text-gray-700'"
                    >
                      Auto-locate in Sidebar
                    </label>
                    <p class="text-xs text-gray-500">
                      Automatically expand and scroll to the active request in the sidebar
                    </p>
                  </div>
                  <button
                    @click="localSettings.autoLocateSidebar = !localSettings.autoLocateSidebar"
                    class="relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none"
                    :class="localSettings.autoLocateSidebar ? 'bg-accent' : (effectiveTheme === 'dark' ? 'bg-gray-600' : 'bg-gray-200')"
                  >
                    <span
                      class="pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out"
                      :class="localSettings.autoLocateSidebar ? 'translate-x-5' : 'translate-x-0'"
                    />
                  </button>
                </div>
                
                <!-- Use System Proxy -->
                <div class="flex items-center justify-between">
                  <div>
                    <label 
                      class="block text-sm font-medium"
                      :class="effectiveTheme === 'dark' ? 'text-gray-300' : 'text-gray-700'"
                    >
                      Use System Proxy
                    </label>
                    <p class="text-xs text-gray-500">
                      Use Windows system proxy settings for HTTP requests
                    </p>
                  </div>
                  <button
                    @click="localSettings.useSystemProxy = !localSettings.useSystemProxy"
                    class="relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none"
                    :class="localSettings.useSystemProxy ? 'bg-accent' : (effectiveTheme === 'dark' ? 'bg-gray-600' : 'bg-gray-200')"
                  >
                    <span
                      class="pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out"
                      :class="localSettings.useSystemProxy ? 'translate-x-5' : 'translate-x-0'"
                    />
                  </button>
                </div>
                
                <!-- Theme Selection -->
                <div>
                  <label 
                    class="block text-sm font-medium mb-2"
                    :class="effectiveTheme === 'dark' ? 'text-gray-300' : 'text-gray-700'"
                  >
                    Theme
                  </label>
                  <div class="flex gap-2">
                    <button
                      v-for="option in themeOptions"
                      :key="option.value"
                      @click="localSettings.theme = option.value"
                      class="flex-1 py-2 px-3 rounded-md text-sm font-medium transition-colors border"
                      :class="[
                        localSettings.theme === option.value
                          ? 'bg-accent text-white border-accent'
                          : (effectiveTheme === 'dark' 
                              ? 'bg-dark-surface border-dark-border text-gray-300 hover:bg-dark-hover' 
                              : 'bg-white border-light-border text-gray-700 hover:bg-gray-50')
                      ]"
                    >
                      {{ option.label }}
                    </button>
                  </div>
                </div>
                
                <!-- Layout Direction -->
                <div>
                  <label 
                    class="block text-sm font-medium mb-2"
                    :class="effectiveTheme === 'dark' ? 'text-gray-300' : 'text-gray-700'"
                  >
                    Layout Direction
                  </label>
                  <div class="flex gap-2">
                    <button
                      @click="localSettings.layoutDirection = 'horizontal'"
                      class="flex-1 py-2 px-3 rounded-md text-sm font-medium transition-colors border"
                      :class="[
                        localSettings.layoutDirection === 'horizontal'
                          ? 'bg-accent text-white border-accent'
                          : (effectiveTheme === 'dark' 
                              ? 'bg-dark-surface border-dark-border text-gray-300 hover:bg-dark-hover' 
                              : 'bg-white border-light-border text-gray-700 hover:bg-gray-50')
                      ]"
                    >
                      Horizontal
                    </button>
                    <button
                      @click="localSettings.layoutDirection = 'vertical'"
                      class="flex-1 py-2 px-3 rounded-md text-sm font-medium transition-colors border"
                      :class="[
                        localSettings.layoutDirection === 'vertical'
                          ? 'bg-accent text-white border-accent'
                          : (effectiveTheme === 'dark' 
                              ? 'bg-dark-surface border-dark-border text-gray-300 hover:bg-dark-hover' 
                              : 'bg-white border-light-border text-gray-700 hover:bg-gray-50')
                      ]"
                    >
                      Vertical
                    </button>
                  </div>
                </div>
              </div>
              
              <div class="flex justify-end gap-3 mt-8">
                <button
                  @click="close"
                  class="px-4 py-2 rounded-md font-medium transition-colors"
                  :class="effectiveTheme === 'dark' ? 'bg-dark-hover text-gray-300 hover:bg-dark-border' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'"
                >
                  Cancel
                </button>
                <button
                  @click="save"
                  class="px-4 py-2 rounded-md font-medium text-white bg-accent hover:bg-accent-hover transition-colors"
                >
                  Save
                </button>
              </div>
            </DialogPanel>
          </TransitionChild>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { Dialog, DialogPanel, DialogTitle, TransitionRoot, TransitionChild } from '@headlessui/vue'
import { useAppStateStore } from '@/stores/appState'
import { api } from '@/services/api'

const props = defineProps<{
  isOpen: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const appState = useAppStateStore()
const effectiveTheme = computed(() => appState.effectiveTheme)

const themeOptions = [
  { value: 'light' as const, label: 'Light' },
  { value: 'dark' as const, label: 'Dark' },
  { value: 'system' as const, label: 'System' },
]

const localSettings = reactive({
  requestTimeout: appState.requestTimeout,
  autoLocateSidebar: appState.autoLocateSidebar,
  useSystemProxy: appState.useSystemProxy,
  theme: appState.theme,
  layoutDirection: appState.layoutDirection,
})

// Reset local settings when modal opens
watch(() => props.isOpen, (isOpen, wasOpen) => {
  if (isOpen === wasOpen) return
  if (isOpen) {
    appState.addModalOpen()
    localSettings.requestTimeout = appState.requestTimeout
    localSettings.autoLocateSidebar = appState.autoLocateSidebar
    localSettings.useSystemProxy = appState.useSystemProxy
    localSettings.theme = appState.theme
    localSettings.layoutDirection = appState.layoutDirection
  } else {
    appState.removeModalOpen()
  }
})

function close() {
  emit('close')
}

async function save() {
  // Update store
  appState.requestTimeout = localSettings.requestTimeout
  appState.autoLocateSidebar = localSettings.autoLocateSidebar
  appState.useSystemProxy = localSettings.useSystemProxy
  appState.theme = localSettings.theme
  appState.layoutDirection = localSettings.layoutDirection
  
  // Persist to backend
  try {
    await api.updateAppState({
      requestTimeout: localSettings.requestTimeout,
      autoLocateSidebar: localSettings.autoLocateSidebar,
      useSystemProxy: localSettings.useSystemProxy,
      theme: localSettings.theme,
      layoutDirection: localSettings.layoutDirection,
    })

    // Update HTTP client proxy setting
    await api.setUseSystemProxy(localSettings.useSystemProxy)

    // Show success toast
    const toast = (window as any).$toast
    if (toast) {
      toast.success('Settings saved successfully')
    }
  } catch (error) {
    console.error('Failed to save settings:', error)

    // Show error toast
    const toast = (window as any).$toast
    if (toast) {
      toast.error('Failed to save settings')
    }
  }

  close()
}
</script>

