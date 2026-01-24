import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { History } from '@/types'

export const useHistoryStore = defineStore('history', () => {
  const history = ref<History[]>([])
  const loading = ref(false)

  // Group history by date
  const groupedHistory = computed(() => {
    const groups: { date: string; label: string; items: History[] }[] = []
    const now = new Date()
    const todayKey = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`
    const yesterdayDate = new Date(now)
    yesterdayDate.setDate(yesterdayDate.getDate() - 1)
    const yesterdayKey = `${yesterdayDate.getFullYear()}-${String(yesterdayDate.getMonth() + 1).padStart(2, '0')}-${String(yesterdayDate.getDate()).padStart(2, '0')}`

    const dateGroups = new Map<string, History[]>()

    for (const item of history.value) {
      const itemDate = new Date(item.createdAt)
      const dateKey = `${itemDate.getFullYear()}-${String(itemDate.getMonth() + 1).padStart(2, '0')}-${String(itemDate.getDate()).padStart(2, '0')}`

      if (!dateGroups.has(dateKey)) {
        dateGroups.set(dateKey, [])
      }
      dateGroups.get(dateKey)!.push(item)
    }

    // Sort dates descending
    const sortedDates = Array.from(dateGroups.keys()).sort((a, b) => b.localeCompare(a))

    for (const dateKey of sortedDates) {
      let label: string

      if (dateKey === todayKey) {
        label = 'Today'
      } else if (dateKey === yesterdayKey) {
        label = 'Yesterday'
      } else {
        const [year, month, day] = dateKey.split('-').map(Number)
        const itemDate = new Date(year, month - 1, day)
        label = itemDate.toLocaleDateString('en-US', { 
          month: 'short', 
          day: 'numeric',
          year: itemDate.getFullYear() !== now.getFullYear() ? 'numeric' : undefined
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
