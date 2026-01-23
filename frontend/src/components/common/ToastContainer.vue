<template>
  <Teleport to="body">
    <div class="fixed top-4 right-4 z-50 space-y-2">
      <TransitionGroup
        enter-active-class="transition duration-200 ease-out"
        enter-from-class="transform translate-x-full opacity-0"
        enter-to-class="transform translate-x-0 opacity-100"
        leave-active-class="transition duration-150 ease-in"
        leave-from-class="transform translate-x-0 opacity-100"
        leave-to-class="transform translate-x-full opacity-0"
      >
        <div
          v-for="toast in toasts"
          :key="toast.id"
          class="flex items-center gap-2 px-4 py-3 rounded-lg shadow-lg min-w-[200px]"
          :class="toastClass(toast.type)"
        >
          <component :is="toastIcon(toast.type)" class="w-5 h-5 flex-shrink-0" />
          <span class="flex-1 text-sm">{{ toast.message }}</span>
          <button
            @click="removeToast(toast.id)"
            class="p-0.5 rounded hover:bg-black/10"
          >
            <XMarkIcon class="w-4 h-4" />
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { 
  CheckCircleIcon, 
  ExclamationCircleIcon, 
  ExclamationTriangleIcon, 
  InformationCircleIcon,
  XMarkIcon
} from '@heroicons/vue/24/solid'
import type { Toast, ToastType } from '@/types'

const toasts = ref<Toast[]>([])

function toastClass(type: ToastType) {
  switch (type) {
    case 'success': return 'bg-green-500 text-white'
    case 'error': return 'bg-red-500 text-white'
    case 'warning': return 'bg-yellow-500 text-white'
    case 'info': return 'bg-blue-500 text-white'
  }
}

function toastIcon(type: ToastType) {
  switch (type) {
    case 'success': return CheckCircleIcon
    case 'error': return ExclamationCircleIcon
    case 'warning': return ExclamationTriangleIcon
    case 'info': return InformationCircleIcon
  }
}

function removeToast(id: string) {
  const index = toasts.value.findIndex(t => t.id === id)
  if (index !== -1) {
    toasts.value.splice(index, 1)
  }
}

function addToast(type: ToastType, message: string, duration = 3000) {
  const id = `toast-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
  toasts.value.push({ id, type, message })
  
  if (duration > 0) {
    setTimeout(() => removeToast(id), duration)
  }
  
  return id
}

// Expose for global access
const toastAPI = {
  success: (message: string) => addToast('success', message),
  error: (message: string) => addToast('error', message),
  warning: (message: string) => addToast('warning', message),
  info: (message: string) => addToast('info', message),
}

// Make available globally
onMounted(() => {
  (window as any).$toast = toastAPI
})

onUnmounted(() => {
  delete (window as any).$toast
})

defineExpose(toastAPI)
</script>
