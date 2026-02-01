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
        class="kv-input w-1/3 min-w-0 px-3 py-1.5 rounded-md border outline-none text-sm"
        :class="[
          effectiveTheme === 'dark'
            ? 'bg-dark-surface border-dark-border text-white placeholder-gray-500 focus:border-accent'
            : 'bg-light-surface border-light-border text-gray-900 placeholder-gray-400 focus:border-accent'
        ]"
      />
      
      <!-- Type selector for form-data -->
      <select
        v-if="showFileUpload"
        :value="item.type || 'text'"
        @change="handleTypeChange(index, ($event.target as HTMLSelectElement).value as 'text' | 'file')"
        class="px-2 py-1.5 rounded-md border outline-none text-sm w-20"
        :class="[
          effectiveTheme === 'dark'
            ? 'bg-dark-surface border-dark-border text-white focus:border-accent'
            : 'bg-light-surface border-light-border text-gray-900 focus:border-accent'
        ]"
      >
        <option value="text">Text</option>
        <option value="file">File</option>
      </select>
      
      <!-- Value input (text) -->
      <input
        v-if="!showFileUpload || (item.type || 'text') === 'text'"
        type="text"
        :value="item.value"
        @input="updateItem(index, { value: ($event.target as HTMLInputElement).value })"
        :placeholder="valuePlaceholder"
        class="kv-input flex-1 px-3 py-1.5 rounded-md border outline-none text-sm"
        :class="[
          effectiveTheme === 'dark'
            ? 'bg-dark-surface border-dark-border text-white placeholder-gray-500 focus:border-accent'
            : 'bg-light-surface border-light-border text-gray-900 placeholder-gray-400 focus:border-accent'
        ]"
      />
      
      <!-- File picker (file type) -->
      <div v-else class="flex-1 flex items-center gap-2">
        <button
          @click="selectFile(index)"
          class="px-3 py-1.5 rounded-md border text-sm transition-colors"
          :class="[
            effectiveTheme === 'dark'
              ? 'bg-dark-surface border-dark-border text-gray-300 hover:bg-dark-hover'
              : 'bg-light-surface border-light-border text-gray-700 hover:bg-light-hover'
          ]"
        >
          Choose File
        </button>
        <span 
          class="text-sm truncate flex-1"
          :class="[
            item.value 
              ? (effectiveTheme === 'dark' ? 'text-gray-300' : 'text-gray-700')
              : (effectiveTheme === 'dark' ? 'text-gray-500' : 'text-gray-400')
          ]"
        >
          {{ item.value ? getFileName(item.value) : 'No file selected' }}
        </span>
      </div>
      
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
        @blur="onNewInputBlur"
        class="kv-input w-1/3 min-w-0 px-3 py-1.5 rounded-md border outline-none text-sm"
        :class="[
          effectiveTheme === 'dark'
            ? 'bg-dark-surface border-dark-border text-white placeholder-gray-500 focus:border-accent'
            : 'bg-light-surface border-light-border text-gray-900 placeholder-gray-400 focus:border-accent'
        ]"
      />
      <select
        v-if="showFileUpload"
        v-model="newType"
        class="px-2 py-1.5 rounded-md border outline-none text-sm w-20"
        :class="[
          effectiveTheme === 'dark'
            ? 'bg-dark-surface border-dark-border text-white focus:border-accent'
            : 'bg-light-surface border-light-border text-gray-900 focus:border-accent'
        ]"
      >
        <option value="text">Text</option>
        <option value="file">File</option>
      </select>
      <!-- Value input for new row (text type) -->
      <input
        v-if="!showFileUpload || newType === 'text'"
        type="text"
        v-model="newValue"
        :placeholder="valuePlaceholder"
        @keydown.enter="addItem"
        @blur="onNewInputBlur"
        class="kv-input flex-1 px-3 py-1.5 rounded-md border outline-none text-sm"
        :class="[
          effectiveTheme === 'dark'
            ? 'bg-dark-surface border-dark-border text-white placeholder-gray-500 focus:border-accent'
            : 'bg-light-surface border-light-border text-gray-900 placeholder-gray-400 focus:border-accent'
        ]"
      />
      <!-- File picker for new row (file type) -->
      <div v-else class="flex-1 flex items-center gap-2">
        <button
          @click="selectNewFile"
          class="px-3 py-1.5 rounded-md border text-sm transition-colors"
          :class="[
            effectiveTheme === 'dark'
              ? 'bg-dark-surface border-dark-border text-gray-300 hover:bg-dark-hover'
              : 'bg-light-surface border-light-border text-gray-700 hover:bg-light-hover'
          ]"
        >
          Choose File
        </button>
        <span 
          class="text-sm truncate flex-1"
          :class="[
            newValue 
              ? (effectiveTheme === 'dark' ? 'text-gray-300' : 'text-gray-700')
              : (effectiveTheme === 'dark' ? 'text-gray-500' : 'text-gray-400')
          ]"
        >
          {{ newValue ? getFileName(newValue) : 'No file selected' }}
        </span>
      </div>
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
  showFileUpload?: boolean
}>()

const emit = defineEmits<{
  'update:items': [value: KeyValue[]]
}>()

const appState = useAppStateStore()
const effectiveTheme = computed(() => appState.effectiveTheme)

const localItems = ref<KeyValue[]>([...props.items])
const newKey = ref('')
const newValue = ref('')
const newType = ref<'text' | 'file'>('text')

// Track if we just emitted to avoid overwriting with stale props
let justEmitted = false

watch(() => props.items, (newItems) => {
  if (justEmitted) {
    justEmitted = false
    return
  }
  localItems.value = [...newItems]
}, { deep: true })

function emitUpdate() {
  justEmitted = true
  emit('update:items', [...localItems.value])
}

function updateItem(index: number, updates: Partial<KeyValue>) {
  localItems.value[index] = { ...localItems.value[index], ...updates }
  emitUpdate()
}

function handleTypeChange(index: number, newType: 'text' | 'file') {
  const currentType = localItems.value[index].type || 'text'
  if (currentType !== newType) {
    // Clear value when switching types since file paths and text values are incompatible
    updateItem(index, { type: newType, value: '' })
  }
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
      type: props.showFileUpload ? newType.value : undefined,
    })
    newKey.value = ''
    newValue.value = ''
    newType.value = 'text'
    emitUpdate()
  }
}

function onNewInputBlur(event: FocusEvent) {
  // Check if focus is moving to another input in the same add row
  const relatedTarget = event.relatedTarget as HTMLElement
  const currentTarget = event.currentTarget as HTMLElement
  const parentRow = currentTarget.parentElement
  
  // If focus is moving to another element in the same row, don't add yet
  if (parentRow && relatedTarget && parentRow.contains(relatedTarget)) {
    return
  }
  
  // Add the item if there's any content
  addItem()
}

async function selectFile(index: number) {
  try {
    const { OpenFileDialog } = await import('../../../wailsjs/go/handlers/DialogHandler')
    const filePath = await OpenFileDialog('Select File')
    if (filePath) {
      updateItem(index, { value: filePath })
    }
  } catch (error) {
    console.error('Failed to open file dialog:', error)
  }
}

async function selectNewFile() {
  try {
    const { OpenFileDialog } = await import('../../../wailsjs/go/handlers/DialogHandler')
    const filePath = await OpenFileDialog('Select File')
    if (filePath) {
      newValue.value = filePath
    }
  } catch (error) {
    console.error('Failed to open file dialog:', error)
  }
}

function getFileName(path: string): string {
  // Extract filename from path
  const parts = path.replace(/\\/g, '/').split('/')
  return parts[parts.length - 1] || path
}
</script>
