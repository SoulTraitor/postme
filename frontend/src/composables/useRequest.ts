import { computed } from 'vue'
import { useTabsStore } from '@/stores/tabs'
import { useResponseStore } from '@/stores/response'
import { useAppStateStore } from '@/stores/appState'
import { useEnvironmentStore } from '@/stores/environment'
import type { KeyValue } from '@/types'

export function useRequest() {
  const tabsStore = useTabsStore()
  const responseStore = useResponseStore()
  const appState = useAppStateStore()
  const envStore = useEnvironmentStore()

  const activeTab = computed(() => tabsStore.activeTab)
  
  const responseState = computed(() => {
    if (!activeTab.value) return { status: 'idle' as const }
    return responseStore.getResponse(activeTab.value.id)
  })

  const isLoading = computed(() => responseState.value.status === 'loading')

  function buildUrl(baseUrl: string, params: KeyValue[]): string {
    let url = envStore.replaceVariables(baseUrl)
    
    try {
      const urlObj = new URL(url)
      for (const param of params) {
        if (param.enabled && param.key) {
          urlObj.searchParams.set(
            envStore.replaceVariables(param.key),
            envStore.replaceVariables(param.value)
          )
        }
      }
      return urlObj.toString()
    } catch {
      // If URL is invalid, just append params
      const queryParams = params
        .filter(p => p.enabled && p.key)
        .map(p => `${encodeURIComponent(envStore.replaceVariables(p.key))}=${encodeURIComponent(envStore.replaceVariables(p.value))}`)
        .join('&')
      
      if (queryParams) {
        return url.includes('?') ? `${url}&${queryParams}` : `${url}?${queryParams}`
      }
      return url
    }
  }

  function buildHeaders(headers: KeyValue[]): KeyValue[] {
    return headers
      .filter(h => h.enabled)
      .map(h => ({
        ...h,
        key: envStore.replaceVariables(h.key),
        value: envStore.replaceVariables(h.value),
      }))
  }

  async function sendRequest() {
    if (!activeTab.value?.url) return

    const tab = activeTab.value
    responseStore.setLoading(tab.id)

    try {
      const url = buildUrl(tab.url, tab.params)
      const headers = buildHeaders(tab.headers)
      const body = tab.bodyType !== 'none' ? envStore.replaceVariables(tab.body) : ''

      // TODO: Call Wails backend
      // const response = await RequestHandler.Execute({
      //   tabId: tab.id,
      //   method: tab.method,
      //   url,
      //   headers,
      //   body,
      //   bodyType: tab.bodyType,
      //   timeout: appState.requestTimeout,
      // })
      
      // Simulate for now
      await new Promise(resolve => setTimeout(resolve, 500))
      
      responseStore.setSuccess(tab.id, {
        statusCode: 200,
        status: '200 OK',
        headers: { 'Content-Type': 'application/json' },
        body: '{"message": "Response placeholder"}',
        size: 35,
        duration: 150,
      })
    } catch (error: any) {
      if (error.message?.includes('context canceled')) {
        responseStore.setCancelled(tab.id)
      } else if (error.message?.includes('timeout')) {
        responseStore.setTimeout(tab.id, appState.requestTimeout)
      } else {
        responseStore.setError(tab.id, error.message || 'Request failed')
      }
    }
  }

  function cancelRequest() {
    if (!activeTab.value) return
    
    // TODO: Call Wails backend
    // RequestHandler.CancelRequest(activeTab.value.id)
    
    responseStore.setCancelled(activeTab.value.id)
  }

  return {
    activeTab,
    responseState,
    isLoading,
    sendRequest,
    cancelRequest,
  }
}
