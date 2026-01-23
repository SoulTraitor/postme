import { ref, watch } from 'vue'
import type { KeyValue } from '@/types'

export function useKeyValue(initialItems: KeyValue[] = []) {
  const items = ref<KeyValue[]>([...initialItems])
  const newKey = ref('')
  const newValue = ref('')

  function addItem() {
    if (newKey.value || newValue.value) {
      items.value.push({
        key: newKey.value,
        value: newValue.value,
        enabled: true,
      })
      newKey.value = ''
      newValue.value = ''
    }
  }

  function removeItem(index: number) {
    items.value.splice(index, 1)
  }

  function updateItem(index: number, updates: Partial<KeyValue>) {
    items.value[index] = { ...items.value[index], ...updates }
  }

  function toggleItem(index: number) {
    items.value[index].enabled = !items.value[index].enabled
  }

  function setItems(newItems: KeyValue[]) {
    items.value = [...newItems]
  }

  function getEnabledItems(): KeyValue[] {
    return items.value.filter(item => item.enabled && item.key)
  }

  return {
    items,
    newKey,
    newValue,
    addItem,
    removeItem,
    updateItem,
    toggleItem,
    setItems,
    getEnabledItems,
  }
}
