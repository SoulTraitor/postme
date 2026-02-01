<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition ease-out duration-100"
      enter-from-class="opacity-0 scale-95"
      enter-to-class="opacity-100 scale-100"
      leave-active-class="transition ease-in duration-75"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-95"
    >
      <div
        v-if="isOpen"
        ref="menuRef"
        class="fixed z-50 min-w-48 rounded-md py-1 ring-1 ring-black ring-opacity-5 focus:outline-none"
        :class="[
          effectiveTheme === 'dark' ? 'bg-dark-elevated shadow-2xl' : 'bg-white shadow-xl',
        ]"
        :style="{
          top: `${position.y}px`,
          left: `${position.x}px`,
          boxShadow: effectiveTheme === 'dark'
            ? '0 20px 25px -5px rgba(0, 0, 0, 0.5), 0 10px 10px -5px rgba(0, 0, 0, 0.3)'
            : '0 20px 25px -5px rgba(0, 0, 0, 0.15), 0 10px 10px -5px rgba(0, 0, 0, 0.1)'
        }"
        @click.stop
      >
        <button
          v-for="item in items"
          :key="item.id"
          @click="handleClick(item)"
          class="w-full text-left px-4 py-2 text-sm flex items-center gap-2 transition-colors"
          :class="[
            item.danger
              ? 'text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20'
              : (effectiveTheme === 'dark' 
                  ? 'text-gray-300 hover:bg-dark-hover' 
                  : 'text-gray-700 hover:bg-gray-100')
          ]"
        >
          <component :is="item.icon" v-if="item.icon" class="w-4 h-4" />
          {{ item.label }}
        </button>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, type Component } from 'vue'
import { useAppStateStore } from '@/stores/appState'

export interface ContextMenuItem {
  id: string
  label: string
  icon?: Component
  danger?: boolean
  action: () => void
}

const props = defineProps<{
  items: ContextMenuItem[]
}>()

const appState = useAppStateStore()
const effectiveTheme = computed(() => appState.effectiveTheme)

const isOpen = ref(false)
const position = ref({ x: 0, y: 0 })
const menuRef = ref<HTMLElement | null>(null)

// Unique ID for this menu instance
const menuId = `context-menu-${Math.random().toString(36).substr(2, 9)}`

// Listen for global close event
function handleGlobalClose(event: CustomEvent) {
  if (event.detail !== menuId) {
    close()
  }
}

onMounted(() => {
  window.addEventListener('context-menu-open', handleGlobalClose as EventListener)
})

onUnmounted(() => {
  window.removeEventListener('context-menu-open', handleGlobalClose as EventListener)
  document.removeEventListener('click', handleClickOutside)
  document.removeEventListener('keydown', handleEscape)
})

function open(event: MouseEvent) {
  event.preventDefault()
  event.stopPropagation()
  
  // Dispatch global event to close other menus
  window.dispatchEvent(new CustomEvent('context-menu-open', { detail: menuId }))
  
  // Calculate position, keeping menu within viewport
  let x = event.clientX
  let y = event.clientY
  
  // Adjust if too close to right edge
  if (x + 200 > window.innerWidth) {
    x = window.innerWidth - 200
  }
  
  // Adjust if too close to bottom edge
  if (y + 200 > window.innerHeight) {
    y = window.innerHeight - 200
  }
  
  position.value = { x, y }
  isOpen.value = true
}

function close() {
  isOpen.value = false
}

function handleClick(item: ContextMenuItem) {
  item.action()
  close()
}

function handleClickOutside(event: MouseEvent) {
  if (menuRef.value && !menuRef.value.contains(event.target as Node)) {
    close()
  }
}

function handleEscape(event: KeyboardEvent) {
  if (event.key === 'Escape') {
    close()
  }
}

watch(isOpen, (open) => {
  if (open) {
    document.addEventListener('click', handleClickOutside)
    document.addEventListener('keydown', handleEscape)
    document.addEventListener('contextmenu', handleClickOutside)
  } else {
    document.removeEventListener('click', handleClickOutside)
    document.removeEventListener('keydown', handleEscape)
    document.removeEventListener('contextmenu', handleClickOutside)
  }
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  document.removeEventListener('keydown', handleEscape)
})

defineExpose({ open, close })
</script>
