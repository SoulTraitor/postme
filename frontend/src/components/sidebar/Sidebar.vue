<template>
  <div
    class="flex flex-col h-full border-r transition-sidebar flex-shrink-0"
    :class="[
      effectiveTheme === 'dark'
        ? 'bg-dark-surface border-dark-border'
        : 'bg-light-surface border-light-border'
    ]"
    :style="{ width: `${appState.sidebarWidth}px` }"
  >
    <!-- Tab switcher -->
    <div class="flex border-b" :class="effectiveTheme === 'dark' ? 'border-dark-border' : 'border-light-border'">
      <button
        @click="appState.sidebarTab = 'collections'"
        class="flex-1 py-2 text-sm font-medium transition-colors"
        :class="[
          appState.sidebarTab === 'collections'
            ? 'text-accent border-b-2 border-accent'
            : (effectiveTheme === 'dark' ? 'text-gray-400 hover:text-white' : 'text-gray-500 hover:text-gray-900')
        ]"
      >
        Collections
      </button>
      <button
        @click="appState.sidebarTab = 'history'"
        class="flex-1 py-2 text-sm font-medium transition-colors"
        :class="[
          appState.sidebarTab === 'history'
            ? 'text-accent border-b-2 border-accent'
            : (effectiveTheme === 'dark' ? 'text-gray-400 hover:text-white' : 'text-gray-500 hover:text-gray-900')
        ]"
      >
        History
      </button>
      
      <!-- Collapse button -->
      <button
        @click="appState.toggleSidebar"
        class="px-2 transition-colors"
        :class="effectiveTheme === 'dark' ? 'hover:bg-dark-hover text-gray-400' : 'hover:bg-light-hover text-gray-500'"
        title="Collapse sidebar"
      >
        <ChevronLeftIcon class="w-4 h-4" />
      </button>
    </div>
    
    <!-- Content -->
    <div class="flex-1 overflow-hidden">
      <CollectionTree v-if="appState.sidebarTab === 'collections'" />
      <HistoryList v-else />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { ChevronLeftIcon } from '@heroicons/vue/24/outline'
import { useAppStateStore } from '@/stores/appState'
import CollectionTree from './CollectionTree.vue'
import HistoryList from './HistoryList.vue'

const appState = useAppStateStore()
const effectiveTheme = computed(() => appState.effectiveTheme)
</script>
