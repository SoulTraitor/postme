<template>
  <Menu as="div" class="relative">
    <MenuButton
      class="flex items-center gap-1 px-3 py-2 rounded-md font-medium text-sm min-w-[100px] justify-between"
      :class="[
        effectiveTheme === 'dark'
          ? 'bg-dark-surface border border-dark-border hover:bg-dark-hover'
          : 'bg-light-surface border border-light-border hover:bg-light-hover',
        methodColor
      ]"
    >
      {{ modelValue }}
      <ChevronDownIcon class="w-4 h-4" />
    </MenuButton>
    
    <Transition
      enter-active-class="transition duration-100 ease-out"
      enter-from-class="transform scale-95 opacity-0"
      enter-to-class="transform scale-100 opacity-100"
      leave-active-class="transition duration-75 ease-in"
      leave-from-class="transform scale-100 opacity-100"
      leave-to-class="transform scale-95 opacity-0"
    >
      <MenuItems
        class="absolute left-0 mt-1 w-32 rounded-md shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none z-50"
        :class="effectiveTheme === 'dark' ? 'bg-dark-elevated' : 'bg-white'"
      >
        <div class="py-1">
          <MenuItem v-for="method in methods" :key="method" v-slot="{ active }">
            <button
              @click="$emit('update:modelValue', method)"
              class="w-full text-left px-4 py-2 text-sm font-medium"
              :class="[
                active ? (effectiveTheme === 'dark' ? 'bg-dark-hover' : 'bg-light-hover') : '',
                getMethodColor(method)
              ]"
            >
              {{ method }}
            </button>
          </MenuItem>
        </div>
      </MenuItems>
    </Transition>
  </Menu>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Menu, MenuButton, MenuItems, MenuItem } from '@headlessui/vue'
import { ChevronDownIcon } from '@heroicons/vue/24/outline'
import { useAppStateStore } from '@/stores/appState'
import { HTTP_METHODS } from '@/types'

const props = defineProps<{
  modelValue: string
}>()

defineEmits<{
  'update:modelValue': [value: string]
}>()

const appState = useAppStateStore()
const effectiveTheme = computed(() => appState.effectiveTheme)
const methods = HTTP_METHODS

const methodColor = computed(() => getMethodColor(props.modelValue))

function getMethodColor(method: string) {
  switch (method.toUpperCase()) {
    case 'GET': return 'text-method-get'
    case 'POST': return 'text-method-post'
    case 'PUT': return 'text-method-put'
    case 'PATCH': return 'text-method-patch'
    case 'DELETE': return 'text-method-delete'
    default: return 'text-method-options'
  }
}
</script>
