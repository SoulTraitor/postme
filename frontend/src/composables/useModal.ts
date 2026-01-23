import { ref } from 'vue'

interface UseModalOptions {
  title?: string
  message?: string
  confirmText?: string
  cancelText?: string
  danger?: boolean
}

export function useModal() {
  const isOpen = ref(false)
  const resolvePromise = ref<((value: boolean) => void) | null>(null)

  function open(): Promise<boolean> {
    isOpen.value = true
    return new Promise(resolve => {
      resolvePromise.value = resolve
    })
  }

  function confirm() {
    isOpen.value = false
    resolvePromise.value?.(true)
    resolvePromise.value = null
  }

  function cancel() {
    isOpen.value = false
    resolvePromise.value?.(false)
    resolvePromise.value = null
  }

  return {
    isOpen,
    open,
    confirm,
    cancel,
  }
}
