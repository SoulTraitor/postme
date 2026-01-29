<template>
  <Teleport to="body">
    <!-- Confirm Modal -->
    <TransitionRoot :show="confirmModal.open" as="template">
      <Dialog as="div" class="relative z-[60]" @close="confirmModal.onCancel?.()">
        <TransitionChild
          enter="ease-out duration-200"
          enter-from="opacity-0"
          enter-to="opacity-100"
          leave="ease-in duration-150"
          leave-from="opacity-100"
          leave-to="opacity-0"
        >
          <div class="fixed inset-0 modal-backdrop" aria-hidden="true" />
        </TransitionChild>

        <div class="fixed inset-0 overflow-y-auto">
          <div class="flex min-h-full items-center justify-center p-4">
            <TransitionChild
              enter="ease-out duration-200"
              enter-from="opacity-0 scale-95"
              enter-to="opacity-100 scale-100"
              leave="ease-in duration-150"
              leave-from="opacity-100 scale-100"
              leave-to="opacity-0 scale-95"
            >
              <DialogPanel 
                class="w-full max-w-md rounded-lg p-6 shadow-xl"
                :class="effectiveTheme === 'dark' ? 'bg-dark-elevated' : 'bg-white'"
              >
                <DialogTitle 
                  class="text-lg font-medium mb-4"
                  :class="effectiveTheme === 'dark' ? 'text-white' : 'text-gray-900'"
                >
                  {{ confirmModal.title }}
                </DialogTitle>
                
                <p 
                  class="mb-6"
                  :class="effectiveTheme === 'dark' ? 'text-gray-300' : 'text-gray-600'"
                >
                  {{ confirmModal.message }}
                </p>
                
                <div class="flex justify-end gap-3">
                  <button
                    ref="cancelButtonRef"
                    @click="confirmModal.onCancel?.()"
                    class="px-4 py-2 rounded-md font-medium transition-colors"
                    :class="effectiveTheme === 'dark' ? 'bg-dark-hover text-gray-300 hover:bg-dark-border' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'"
                  >
                    Cancel
                  </button>
                  <button
                    @click="confirmModal.onConfirm?.()"
                    class="px-4 py-2 rounded-md font-medium text-white transition-colors"
                    :class="confirmModal.danger ? 'bg-red-500 hover:bg-red-600' : 'bg-accent hover:bg-accent-hover'"
                  >
                    {{ confirmModal.confirmText || 'Confirm' }}
                  </button>
                </div>
              </DialogPanel>
            </TransitionChild>
          </div>
        </div>
      </Dialog>
    </TransitionRoot>
    
    <!-- Input Modal -->
    <TransitionRoot :show="inputModal.open" as="template">
      <Dialog as="div" class="relative z-[60]" @close="inputModal.onCancel?.()">
        <TransitionChild
          enter="ease-out duration-200"
          enter-from="opacity-0"
          enter-to="opacity-100"
          leave="ease-in duration-150"
          leave-from="opacity-100"
          leave-to="opacity-0"
        >
          <div class="fixed inset-0 modal-backdrop" aria-hidden="true" />
        </TransitionChild>

        <div class="fixed inset-0 overflow-y-auto">
          <div class="flex min-h-full items-center justify-center p-4">
            <TransitionChild
              enter="ease-out duration-200"
              enter-from="opacity-0 scale-95"
              enter-to="opacity-100 scale-100"
              leave="ease-in duration-150"
              leave-from="opacity-100 scale-100"
              leave-to="opacity-0 scale-95"
            >
              <DialogPanel 
                class="w-full max-w-md rounded-lg p-6 shadow-xl"
                :class="effectiveTheme === 'dark' ? 'bg-dark-elevated' : 'bg-white'"
              >
                <DialogTitle 
                  class="text-lg font-medium mb-4"
                  :class="effectiveTheme === 'dark' ? 'text-white' : 'text-gray-900'"
                >
                  {{ inputModal.title }}
                </DialogTitle>
                
                <input
                  v-model="inputModal.value"
                  type="text"
                  :placeholder="inputModal.placeholder"
                  @keydown.enter="inputModal.onConfirm?.(inputModal.value)"
                  class="w-full px-3 py-2 rounded-md border outline-none text-sm mb-6"
                  :class="[
                    effectiveTheme === 'dark'
                      ? 'bg-dark-surface border-dark-border text-white placeholder-gray-500 focus:border-accent'
                      : 'bg-white border-light-border text-gray-900 placeholder-gray-400 focus:border-accent'
                  ]"
                />
                
                <div class="flex justify-end gap-3">
                  <button
                    @click="inputModal.onCancel?.()"
                    class="px-4 py-2 rounded-md font-medium transition-colors"
                    :class="effectiveTheme === 'dark' ? 'bg-dark-hover text-gray-300 hover:bg-dark-border' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'"
                  >
                    Cancel
                  </button>
                  <button
                    @click="inputModal.onConfirm?.(inputModal.value)"
                    class="px-4 py-2 rounded-md font-medium text-white bg-accent hover:bg-accent-hover transition-colors"
                  >
                    {{ inputModal.confirmText || 'Save' }}
                  </button>
                </div>
              </DialogPanel>
            </TransitionChild>
          </div>
        </div>
      </Dialog>
    </TransitionRoot>
    
    <!-- Select Modal -->
    <TransitionRoot :show="selectModal.open" as="template">
      <Dialog as="div" class="relative z-[60]" @close="selectModal.onCancel?.()">
        <TransitionChild
          enter="ease-out duration-200"
          enter-from="opacity-0"
          enter-to="opacity-100"
          leave="ease-in duration-150"
          leave-from="opacity-100"
          leave-to="opacity-0"
        >
          <div class="fixed inset-0 modal-backdrop" aria-hidden="true" />
        </TransitionChild>

        <div class="fixed inset-0 overflow-y-auto">
          <div class="flex min-h-full items-center justify-center p-4">
            <TransitionChild
              enter="ease-out duration-200"
              enter-from="opacity-0 scale-95"
              enter-to="opacity-100 scale-100"
              leave="ease-in duration-150"
              leave-from="opacity-100 scale-100"
              leave-to="opacity-0 scale-95"
            >
              <DialogPanel 
                class="w-full max-w-md rounded-lg p-6 shadow-xl"
                :class="effectiveTheme === 'dark' ? 'bg-dark-elevated' : 'bg-white'"
              >
                <DialogTitle 
                  class="text-lg font-medium mb-4"
                  :class="effectiveTheme === 'dark' ? 'text-white' : 'text-gray-900'"
                >
                  {{ selectModal.title }}
                </DialogTitle>
                
                <div class="space-y-2 max-h-64 overflow-auto mb-6">
                  <button
                    v-for="option in selectModal.options"
                    :key="option.value"
                    @click="selectModal.onSelect?.(option.value)"
                    class="w-full text-left px-4 py-2 rounded-md text-sm transition-colors"
                    :class="[
                      effectiveTheme === 'dark'
                        ? 'hover:bg-dark-hover text-gray-200'
                        : 'hover:bg-light-hover text-gray-800'
                    ]"
                  >
                    {{ option.label }}
                  </button>
                </div>
                
                <div class="flex justify-end">
                  <button
                    @click="selectModal.onCancel?.()"
                    class="px-4 py-2 rounded-md font-medium transition-colors"
                    :class="effectiveTheme === 'dark' ? 'bg-dark-hover text-gray-300 hover:bg-dark-border' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'"
                  >
                    Cancel
                  </button>
                </div>
              </DialogPanel>
            </TransitionChild>
          </div>
        </div>
      </Dialog>
    </TransitionRoot>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, watch } from 'vue'
import { Dialog, DialogPanel, DialogTitle, TransitionRoot, TransitionChild } from '@headlessui/vue'
import { useAppStateStore } from '@/stores/appState'

const appState = useAppStateStore()
const effectiveTheme = computed(() => appState.effectiveTheme)

const cancelButtonRef = ref(null)

const confirmModal = reactive({
  open: false,
  title: '',
  message: '',
  confirmText: 'Confirm',
  danger: false,
  onConfirm: null as (() => void) | null,
  onCancel: null as (() => void) | null,
})

const inputModal = reactive({
  open: false,
  title: '',
  value: '',
  placeholder: '',
  confirmText: 'Save',
  onConfirm: null as ((value: string) => void) | null,
  onCancel: null as (() => void) | null,
})

const selectModal = reactive({
  open: false,
  title: '',
  options: [] as Array<{ value: any; label: string }>,
  onSelect: null as ((value: any) => void) | null,
  onCancel: null as (() => void) | null,
})

function showConfirm(options: {
  title: string
  message: string
  confirmText?: string
  danger?: boolean
}): Promise<boolean> {
  return new Promise(resolve => {
    confirmModal.title = options.title
    confirmModal.message = options.message
    confirmModal.confirmText = options.confirmText || 'Confirm'
    confirmModal.danger = options.danger || false
    confirmModal.onConfirm = () => {
      confirmModal.open = false
      resolve(true)
    }
    confirmModal.onCancel = () => {
      confirmModal.open = false
      resolve(false)
    }
    confirmModal.open = true
  })
}

function showInput(options: {
  title: string
  value?: string
  placeholder?: string
  confirmText?: string
}): Promise<string | null> {
  return new Promise(resolve => {
    inputModal.title = options.title
    inputModal.value = options.value || ''
    inputModal.placeholder = options.placeholder || ''
    inputModal.confirmText = options.confirmText || 'Save'
    inputModal.onConfirm = (value: string) => {
      inputModal.open = false
      resolve(value)
    }
    inputModal.onCancel = () => {
      inputModal.open = false
      resolve(null)
    }
    inputModal.open = true
  })
}

function showSelect<T>(options: {
  title: string
  options: Array<{ value: T; label: string }>
}): Promise<T | null> {
  return new Promise(resolve => {
    selectModal.title = options.title
    selectModal.options = options.options
    selectModal.onSelect = (value: T) => {
      selectModal.open = false
      resolve(value)
    }
    selectModal.onCancel = () => {
      selectModal.open = false
      resolve(null)
    }
    selectModal.open = true
  })
}

const modalAPI = {
  confirm: showConfirm,
  input: showInput,
  select: showSelect,
}

const anyModalOpen = computed(() => confirmModal.open || inputModal.open || selectModal.open)

watch(anyModalOpen, (isOpen, wasOpen) => {
  if (isOpen === wasOpen) return
  if (isOpen) {
    appState.addModalOpen()
  } else {
    appState.removeModalOpen()
  }
})

onMounted(() => {
  (window as any).$modal = modalAPI
})

onUnmounted(() => {
  delete (window as any).$modal
})

defineExpose(modalAPI)
</script>
