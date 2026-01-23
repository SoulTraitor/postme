import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { History } from '@/types'

export const useHistoryStore = defineStore('history', () => {
  const history = ref<History[]>([])
  const loading = ref(false)

  // Group history by date
  const groupedHistory = computed(() => {
    const groups: { date: string; label: string; items: History[] }[] = []
    const today = new Date()
    today.setHours(0, 0, 0, 0)
    const yesterday = new Date(today)
    yesterday.setDate(yesterday.getDate() - 1)

    const dateGroups = new Map<string, History[]>()

    for (const item of history.value) {
      const itemDate = new Date(item.createdAt)
      itemDate.setHours(0, 0, 0, 0)
      const dateKey = itemDate.toISOString().split('T')[0]

      if (!dateGroups.has(dateKey)) {
        dateGroups.set(dateKey, [])
      }
      dateGroups.get(dateKey)!.push(item)
    }

    // Sort dates descending
    const sortedDates = Array.from(dateGroups.keys()).sort((a, b) => b.localeCompare(a))

    for (const dateKey of sortedDates) {
      const itemDate = new Date(dateKey)
      let label: string

      if (itemDate.getTime() === today.getTime()) {
        label = 'Today'
      } else if (itemDate.getTime() === yesterday.getTime()) {
        label = 'Yesterday'
      } else {
        label = itemDate.toLocaleDateString('en-US', { 
          month: 'short', 
          day: 'numeric',
          year: itemDate.getFullYear() !== today.getFullYear() ? 'numeric' : undefined
        })
      }

      groups.push({
        date: dateKey,
        label,
        items: dateGroups.get(dateKey)!,
      })
    }

    return groups
  })

  // Set history
  function setHistory(items: History[]) {
    history.value = items
  }

  // Add history item
  function addHistory(item: History) {
    history.value.unshift(item)
  }

  // Delete history item
  function deleteHistory(id: number) {
    const index = history.value.findIndex(h => h.id === id)
    if (index !== -1) {
      history.value.splice(index, 1)
    }
  }

  // Clear all history
  function clearHistory() {
    history.value = []
  }

  return {
    history,
    loading,
    groupedHistory,
    setHistory,
    addHistory,
    deleteHistory,
    clearHistory,
  }
})
