import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Response, ResponseState } from '@/types'

export const useResponseStore = defineStore('response', () => {
  // Map of tabId to response state
  const responses = ref<Map<string, ResponseState>>(new Map())

  function getResponse(tabId: string): ResponseState {
    return responses.value.get(tabId) || { status: 'idle' }
  }

  function setLoading(tabId: string) {
    responses.value.set(tabId, { status: 'loading' })
  }

  function setSuccess(tabId: string, response: Response) {
    responses.value.set(tabId, { status: 'success', response })
  }

  function setError(tabId: string, message: string) {
    responses.value.set(tabId, { status: 'error', message })
  }

  function setCancelled(tabId: string) {
    responses.value.set(tabId, { status: 'cancelled' })
  }

  function setTimeout(tabId: string, seconds: number) {
    responses.value.set(tabId, { status: 'timeout', seconds })
  }

  function setIdle(tabId: string) {
    responses.value.set(tabId, { status: 'idle' })
  }

  function clearResponse(tabId: string) {
    responses.value.delete(tabId)
  }

  function clearAll() {
    responses.value.clear()
  }

  return {
    responses,
    getResponse,
    setLoading,
    setSuccess,
    setError,
    setCancelled,
    setTimeout,
    setIdle,
    clearResponse,
    clearAll,
  }
})
