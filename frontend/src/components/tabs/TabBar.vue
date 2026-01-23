<template>
  <div 
    class="h-9 flex items-center border-b overflow-x-auto"
    :class="[
      effectiveTheme === 'dark'
        ? 'bg-dark-base border-dark-border'
        : 'bg-light-base border-light-border'
    ]"
  >
    <!-- Tabs -->
    <div 
      class="flex items-center min-w-0 flex-1"
      @dragover="onDragOverContainer"
      @dragleave="onDragLeaveContainer"
      @drop="onDropContainer"
    >
      <TabItem
        v-for="(tab, index) in tabs"
        :key="tab.id"
        :tab="tab"
        :isActive="tab.id === activeTabId"
        :class="dropIndicatorClass(index)"
        draggable="true"
        @dragstart="onDragStart($event, index)"
        @dragover.stop="onDragOver($event, index)"
        @dragleave="onDragLeave"
        @drop.stop="onDrop($event, index)"
        @dragend="onDragEnd"
        @click="setActiveTab(tab.id)"
        @close="closeTab(tab.id)"
        @dblclick="pinTab(tab.id)"
      />
      <!-- Drop zone after last tab -->
      <div 
        v-if="dragFromIndex !== null && dropTargetIndex === tabs.length"
        class="w-1 h-6 bg-accent rounded"
      />
    </div>
    
    <!-- Add tab button -->
    <button
      @click="addTab"
      class="flex-shrink-0 p-2 transition-colors"
      :class="[
        effectiveTheme === 'dark'
          ? 'hover:bg-dark-hover text-gray-400 hover:text-white'
          : 'hover:bg-light-hover text-gray-500 hover:text-gray-900'
      ]"
      title="New Tab (Ctrl+T)"
    >
      <PlusIcon class="w-4 h-4" />
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { PlusIcon } from '@heroicons/vue/24/outline'
import { useAppStateStore } from '@/stores/appState'
import { useTabsStore } from '@/stores/tabs'
import TabItem from './TabItem.vue'

const appState = useAppStateStore()
const tabsStore = useTabsStore()

const effectiveTheme = computed(() => appState.effectiveTheme)
const tabs = computed(() => tabsStore.tabs)
const activeTabId = computed(() => tabsStore.activeTabId)

// Drag state
const dragFromIndex = ref<number | null>(null)
const dropTargetIndex = ref<number | null>(null)
const dropPosition = ref<'before' | 'after'>('before')

function setActiveTab(id: string) {
  tabsStore.setActiveTab(id)
}

function closeTab(id: string) {
  tabsStore.closeTab(id)
}

function addTab() {
  tabsStore.addTab()
}

function pinTab(id: string) {
  tabsStore.pinTab(id)
}

function onDragStart(e: DragEvent, index: number) {
  dragFromIndex.value = index
  e.dataTransfer!.effectAllowed = 'move'
  e.dataTransfer!.setData('text/plain', String(index))
}

function onDragOver(e: DragEvent, index: number) {
  e.preventDefault()
  e.dataTransfer!.dropEffect = 'move'
  
  const rect = (e.currentTarget as HTMLElement).getBoundingClientRect()
  const x = e.clientX - rect.left
  dropPosition.value = x < rect.width / 2 ? 'before' : 'after'
  dropTargetIndex.value = index
}

function onDragLeave() {
  dropTargetIndex.value = null
}

function onDrop(e: DragEvent, toIndex: number) {
  e.preventDefault()
  if (dragFromIndex.value === null) return
  
  let adjustedToIndex = toIndex
  if (dropPosition.value === 'after') {
    adjustedToIndex = toIndex + 1
  }
  
  // Adjust for the fact that we're removing from the array first
  if (dragFromIndex.value < adjustedToIndex) {
    adjustedToIndex -= 1
  }
  
  tabsStore.reorderTabs(dragFromIndex.value, adjustedToIndex)
  onDragEnd()
}

function onDragEnd() {
  dragFromIndex.value = null
  dropTargetIndex.value = null
}

function onDragOverContainer(e: DragEvent) {
  if (dragFromIndex.value === null) return
  e.preventDefault()
  e.dataTransfer!.dropEffect = 'move'
  
  // Check if we're past the last tab
  const container = e.currentTarget as HTMLElement
  const tabs = container.querySelectorAll('[draggable="true"]')
  if (tabs.length === 0) return
  
  const lastTab = tabs[tabs.length - 1] as HTMLElement
  const lastTabRect = lastTab.getBoundingClientRect()
  
  // If cursor is past the right edge of the last tab, set drop target to after last
  if (e.clientX > lastTabRect.right) {
    dropTargetIndex.value = tabsStore.tabs.length
    dropPosition.value = 'after'
  }
}

function onDragLeaveContainer(e: DragEvent) {
  const container = e.currentTarget as HTMLElement
  const relatedTarget = e.relatedTarget as HTMLElement
  if (!container.contains(relatedTarget)) {
    dropTargetIndex.value = null
  }
}

function onDropContainer(e: DragEvent) {
  if (dragFromIndex.value === null) return
  if (dropTargetIndex.value !== tabsStore.tabs.length) return
  
  e.preventDefault()
  
  // Move to end
  const toIndex = tabsStore.tabs.length - 1
  if (dragFromIndex.value !== toIndex) {
    tabsStore.reorderTabs(dragFromIndex.value, toIndex)
  }
  onDragEnd()
}

function dropIndicatorClass(index: number) {
  if (dropTargetIndex.value !== index) return ''
  return dropPosition.value === 'before' ? 'border-l-2 border-accent' : 'border-r-2 border-accent'
}
</script>
