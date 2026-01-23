<template>
  <div class="space-y-2">
    <div 
      v-for="(item, index) in localItems" 
      :key="index"
      class="flex items-center gap-2"
    >
      <!-- Enabled checkbox -->
      <input
        type="checkbox"
        :checked="item.enabled"
        @change="updateItem(index, { enabled: ($event.target as HTMLInputElement).checked })"
        class="w-4 h-4 rounded border-gray-300 text-accent focus:ring-accent"
      />
      
      <!-- Key input -->
      <input
        type="text"
        :value="item.key"
        @input="updateItem(index, { key: ($event.target as HTMLInputElement).value })"
        :placeholder="keyPlaceholder"
        class="flex-1 px-3 py-1.5 rounded-md border outline-none text-sm"
        :class="[
          effectiveTheme === 'dark'
            ? 'bg-dark-surface border-dark-border text-white placeholder-gray-500 focus:border-accent'
            : 'bg-light-surface border-light-border text-gray-900 placeholder-gray-400 focus:border-accent'
        ]"
      />
      
      <!-- Value input -->
      <input
        type="text"
        :value="item.value"
        @input="updateItem(index, { value: ($event.target as HTMLInputElement).value })"
        :placeholder="valuePlaceholder"
        class="flex-1 px-3 py-1.5 rounded-md border outline-none text-sm"
        :class="[
          effectiveTheme === 'dark'
            ? 'bg-dark-surface border-dark-border text-white placeholder-gray-500 focus:border-accent'
            : 'bg-light-surface border-light-border text-gray-900 placeholder-gray-400 focus:border-accent'
        ]"
      />
      
      <!-- Delete button -->
      <button
        @click="removeItem(index)"
        class="p-1.5 rounded-md transition-colors"
        :class="effectiveTheme === 'dark' ? 'hover:bg-dark-hover text-gray-400 hover:text-red-400' : 'hover:bg-light-hover text-gray-500 hover:text-red-500'"
      >
        <TrashIcon class="w-4 h-4" />
      </button>
    </div>
    
    <!-- Add new item row -->
    <div class="flex items-center gap-2">
      <div class="w-4" /> <!-- Spacer for checkbox alignment -->
      <input
        type="text"
        v-model="newKey"
        :placeholder="keyPlaceholder"
        @keydown.enter="addItem"
        class="flex-1 px-3 py-1.5 rounded-md border outline-none text-sm"
        :class="[
          effectiveTheme === 'dark'
            ? 'bg-dark-surface border-dark-border text-white placeholder-gray-500 focus:border-accent'
            : 'bg-light-surface border-light-border text-gray-900 placeholder-gray-400 focus:border-accent'
        ]"
      />
      <input
        type="text"
        v-model="newValue"
        :placeholder="valuePlaceholder"
        @keydown.enter="addItem"
        class="flex-1 px-3 py-1.5 rounded-md border outline-none text-sm"
        :class="[
          effectiveTheme === 'dark'
            ? 'bg-dark-surface border-dark-border text-white placeholder-gray-500 focus:border-accent'
            : 'bg-light-surface border-light-border text-gray-900 placeholder-gray-400 focus:border-accent'
        ]"
      />
      <button
        @click="addItem"
        class="p-1.5 rounded-md transition-colors"
        :class="effectiveTheme === 'dark' ? 'hover:bg-dark-hover text-gray-400 hover:text-accent' : 'hover:bg-light-hover text-gray-500 hover:text-accent'"
      >
        <PlusIcon class="w-4 h-4" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { TrashIcon, PlusIcon } from '@heroicons/vue/24/outline'
import { useAppStateStore } from '@/stores/appState'
import type { KeyValue } from '@/types'

const props = defineProps<{
  items: KeyValue[]
  keyPlaceholder?: string
  valuePlaceholder?: string
}>()

const emit = defineEmits<{
  'update:items': [value: KeyValue[]]
}>()

const appState = useAppStateStore()
const effectiveTheme = computed(() => appState.effectiveTheme)

const localItems = ref<KeyValue[]>([...props.items])
const newKey = ref('')
const newValue = ref('')

watch(() => props.items, (newItems) => {
  localItems.value = [...newItems]
}, { deep: true })

function emitUpdate() {
  emit('update:items', [...localItems.value])
}

function updateItem(index: number, updates: Partial<KeyValue>) {
  localItems.value[index] = { ...localItems.value[index], ...updates }
  emitUpdate()
}

function removeItem(index: number) {
  localItems.value.splice(index, 1)
  emitUpdate()
}

function addItem() {
  if (newKey.value || newValue.value) {
    localItems.value.push({
      key: newKey.value,
      value: newValue.value,
      enabled: true,
    })
    newKey.value = ''
    newValue.value = ''
    emitUpdate()
  }
}
</script>
