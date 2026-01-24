<template>
  <Menu as="div" class="relative">
    <MenuButton
      class="flex items-center gap-1 px-2 py-1 rounded-md text-sm transition-colors"
      :class="[
        effectiveTheme === 'dark'
          ? 'hover:bg-dark-hover text-gray-300'
          : 'hover:bg-light-hover text-gray-600'
      ]"
    >
      <span v-if="activeEnvironment" class="truncate max-w-[120px]">{{ activeEnvironment.name }}</span>
      <span v-else class="text-gray-500">No Environment</span>
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
        class="absolute right-0 mt-1 w-48 rounded-md shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none z-50"
        :class="[
          effectiveTheme === 'dark'
            ? 'bg-dark-elevated'
            : 'bg-white'
        ]"
      >
        <div class="py-1">
          <!-- No Environment option -->
          <MenuItem v-slot="{ active }">
            <button
              @click="selectEnvironment(null)"
              class="w-full text-left px-4 py-2 text-sm flex items-center gap-2"
              :class="[
                active 
                  ? (effectiveTheme === 'dark' ? 'bg-dark-hover' : 'bg-light-hover')
                  : '',
                !activeEnvironment ? 'text-accent' : ''
              ]"
            >
              <CheckIcon v-if="!activeEnvironment" class="w-4 h-4" />
              <span :class="!activeEnvironment ? '' : 'pl-6'">No Environment</span>
            </button>
          </MenuItem>
          
          <div v-if="environments.length > 0" class="border-t my-1" :class="effectiveTheme === 'dark' ? 'border-dark-border' : 'border-light-border'" />
          
          <!-- Environment list -->
          <MenuItem v-for="env in environments" :key="env.id" v-slot="{ active }">
            <button
              @click="selectEnvironment(env.id)"
              class="w-full text-left px-4 py-2 text-sm flex items-center gap-2"
              :class="[
                active 
                  ? (effectiveTheme === 'dark' ? 'bg-dark-hover' : 'bg-light-hover')
                  : '',
                activeEnvironment?.id === env.id ? 'text-accent' : ''
              ]"
            >
              <CheckIcon v-if="activeEnvironment?.id === env.id" class="w-4 h-4" />
              <span class="truncate" :class="activeEnvironment?.id === env.id ? '' : 'pl-6'">{{ env.name }}</span>
            </button>
          </MenuItem>
          
          <div class="border-t my-1" :class="effectiveTheme === 'dark' ? 'border-dark-border' : 'border-light-border'" />
          
          <!-- Manage environments -->
          <MenuItem v-slot="{ active }">
            <button
              @click="manageEnvironments"
              class="w-full text-left px-4 py-2 text-sm"
              :class="active ? (effectiveTheme === 'dark' ? 'bg-dark-hover' : 'bg-light-hover') : ''"
            >
              Manage Environments
            </button>
          </MenuItem>
        </div>
      </MenuItems>
    </Transition>
  </Menu>
  
  <!-- Environment Management Modal -->
  <EnvironmentModal :isOpen="modalOpen" @close="modalOpen = false" />
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Menu, MenuButton, MenuItems, MenuItem } from '@headlessui/vue'
import { ChevronDownIcon, CheckIcon } from '@heroicons/vue/24/outline'
import { useAppStateStore } from '@/stores/appState'
import { useEnvironmentStore } from '@/stores/environment'
import EnvironmentModal from '@/components/modals/EnvironmentModal.vue'

const appState = useAppStateStore()
const envStore = useEnvironmentStore()

const effectiveTheme = computed(() => appState.effectiveTheme)
const environments = computed(() => envStore.environments)
const activeEnvironment = computed(() => envStore.activeEnvironment)

const modalOpen = ref(false)

function selectEnvironment(id: number | null) {
  envStore.setActiveEnv(id)
  appState.activeEnvId = id
}

function manageEnvironments() {
  modalOpen.value = true
}
</script>
