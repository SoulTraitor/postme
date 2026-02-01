<template>
  <TransitionRoot :show="isOpen" as="template">
    <Dialog as="div" class="relative z-50" @close="close">
      <TransitionChild
        enter="ease-out duration-300"
        enter-from="opacity-0"
        enter-to="opacity-100"
        leave="ease-in duration-200"
        leave-from="opacity-100"
        leave-to="opacity-0"
        as="template"
      >
        <div
          class="fixed inset-0 modal-backdrop backdrop-blur-sm"
          style="transition: opacity 300ms ease-out, backdrop-filter 300ms ease-out"
          aria-hidden="true"
        />
      </TransitionChild>

      <div class="fixed inset-0 overflow-y-auto">
        <div class="flex min-h-full items-center justify-center p-4">
          <TransitionChild
            enter="ease-out duration-300"
            enter-from="opacity-0 scale-90"
            enter-to="opacity-100 scale-100"
            leave="ease-in duration-200"
            leave-from="opacity-100 scale-100"
            leave-to="opacity-0 scale-95"
          >
            <DialogPanel 
              class="w-full max-w-3xl rounded-lg p-6 shadow-xl"
              :class="effectiveTheme === 'dark' ? 'bg-dark-elevated' : 'bg-white'"
            >
              <DialogTitle 
                class="text-lg font-medium mb-6 flex items-center justify-between"
                :class="effectiveTheme === 'dark' ? 'text-white' : 'text-gray-900'"
              >
                <span>Manage Environments</span>
                <button
                  @click="close"
                  class="p-1 rounded hover:bg-gray-500/20"
                >
                  <XMarkIcon class="w-5 h-5" />
                </button>
              </DialogTitle>
              
              <div class="flex gap-4 h-96">
                <!-- Sidebar: Environment list -->
                <div 
                  class="w-48 flex-shrink-0 border-r pr-4 overflow-y-auto overflow-x-hidden"
                  :class="effectiveTheme === 'dark' ? 'border-dark-border' : 'border-light-border'"
                >
                  <div class="mb-2">
                    <button
                      @click="selectGlobals"
                      class="w-full text-left px-3 py-2 rounded-md text-sm font-medium"
                      :class="[
                        selectedEnvId === null 
                          ? 'bg-accent text-white' 
                          : (effectiveTheme === 'dark' ? 'hover:bg-dark-hover text-gray-300' : 'hover:bg-light-hover text-gray-700')
                      ]"
                    >
                      Global Variables
                    </button>
                  </div>
                  
                  <div class="border-t mb-2 pt-2" :class="effectiveTheme === 'dark' ? 'border-dark-border' : 'border-light-border'">
                    <div class="flex items-center justify-between mb-2">
                      <span class="text-xs font-medium text-gray-500 uppercase">Environments</span>
                      <button
                        @click="addEnvironment"
                        class="p-1 rounded hover:bg-gray-500/20"
                        title="Add Environment"
                      >
                        <PlusIcon class="w-4 h-4" />
                      </button>
                    </div>
                    
                    <div 
                      v-for="env in environments" 
                      :key="env.id"
                      class="flex items-center group"
                    >
                      <button
                        @click="selectEnv(env.id)"
                        class="flex-1 text-left px-3 py-2 rounded-md text-sm truncate"
                        :class="[
                          selectedEnvId === env.id 
                            ? 'bg-accent text-white' 
                            : (effectiveTheme === 'dark' ? 'hover:bg-dark-hover text-gray-300' : 'hover:bg-light-hover text-gray-700')
                        ]"
                      >
                        {{ env.name }}
                      </button>
                      <button
                        @click="deleteEnv(env.id)"
                        class="p-1 rounded opacity-0 group-hover:opacity-100 hover:bg-red-500/20 text-red-500"
                        title="Delete"
                      >
                        <TrashIcon class="w-4 h-4" />
                      </button>
                    </div>
                  </div>
                </div>
                
                  <!-- Main: Variable editor -->
                <div class="flex-1 min-w-[500px] overflow-hidden flex flex-col">
                  <!-- Environment name (if editing an environment) -->
                  <div v-if="selectedEnvId !== null" class="mb-4">
                    <label 
                      class="block text-sm font-medium mb-1"
                      :class="effectiveTheme === 'dark' ? 'text-gray-300' : 'text-gray-700'"
                    >
                      Environment Name
                    </label>
                    <input
                      v-model="currentName"
                      type="text"
                      class="w-full px-3 py-2 rounded-md border outline-none text-sm"
                      :class="[
                        effectiveTheme === 'dark'
                          ? 'bg-dark-surface border-dark-border text-white focus:border-accent'
                          : 'bg-white border-light-border text-gray-900 focus:border-accent'
                      ]"
                      @change="saveCurrentEnv"
                    />
                  </div>
                  <div v-else class="mb-4">
                    <h3 
                      class="text-sm font-medium"
                      :class="effectiveTheme === 'dark' ? 'text-gray-300' : 'text-gray-700'"
                    >
                      Global Variables
                    </h3>
                    <p class="text-xs text-gray-500">These variables are available in all environments</p>
                  </div>
                  
                  <!-- Variable list -->
                  <div class="flex-1 overflow-y-scroll">
                    <table class="w-full text-sm">
                      <thead>
                        <tr class="text-left text-xs font-medium text-gray-500 uppercase">
                          <th class="pb-2 w-1/3">Variable</th>
                          <th class="pb-2">Value</th>
                          <th class="pb-2 w-16">Secret</th>
                          <th class="pb-2 w-10"></th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-for="(variable, index) in currentVariables" :key="index" class="border-t" :class="effectiveTheme === 'dark' ? 'border-dark-border' : 'border-light-border'">
                          <td class="py-2 pr-2">
                            <input
                              v-model="variable.key"
                              type="text"
                              placeholder="Variable name"
                              class="w-full px-2 py-1 rounded border outline-none text-sm"
                              :class="[
                                isDuplicate(variable.key)
                                  ? 'border-red-500 focus:border-red-500'
                                  : (effectiveTheme === 'dark'
                                    ? 'bg-dark-surface border-dark-border text-white focus:border-accent'
                                    : 'bg-white border-light-border text-gray-900 focus:border-accent'),
                                effectiveTheme === 'dark' ? 'bg-dark-surface text-white' : 'bg-white text-gray-900'
                              ]"
                              @change="saveCurrentEnv"
                            />
                          </td>
                          <td class="py-2 pr-2">
                            <input
                              v-model="variable.value"
                              :type="variable.secret ? 'password' : 'text'"
                              placeholder="Value"
                              class="w-full px-2 py-1 rounded border outline-none text-sm"
                              :class="[
                                effectiveTheme === 'dark'
                                  ? 'bg-dark-surface border-dark-border text-white focus:border-accent'
                                  : 'bg-white border-light-border text-gray-900 focus:border-accent'
                              ]"
                              @change="saveCurrentEnv"
                            />
                          </td>
                          <td class="py-2 text-center">
                            <input
                              v-model="variable.secret"
                              type="checkbox"
                              class="rounded border-gray-300 text-accent focus:ring-accent"
                              @change="saveCurrentEnv"
                            />
                          </td>
                          <td class="py-2">
                            <button
                              @click.stop="removeVariable(index, $event)"
                              class="p-1 rounded hover:bg-red-500/20 text-red-500"
                            >
                              <TrashIcon class="w-4 h-4" />
                            </button>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    
                    <button
                      @click="addVariable"
                      class="mt-2 text-sm text-accent hover:text-accent-hover flex items-center gap-1"
                    >
                      <PlusIcon class="w-4 h-4" />
                      Add Variable
                    </button>
                  </div>
                </div>
              </div>
              
              <div class="flex justify-end gap-3 mt-6 pt-4 border-t" :class="effectiveTheme === 'dark' ? 'border-dark-border' : 'border-light-border'">
                <button
                  @click="close"
                  class="px-4 py-2 rounded-md font-medium text-white bg-accent hover:bg-accent-hover transition-colors"
                >
                  Done
                </button>
              </div>
            </DialogPanel>
          </TransitionChild>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { Dialog, DialogPanel, DialogTitle, TransitionRoot, TransitionChild } from '@headlessui/vue'
import { PlusIcon, TrashIcon, XMarkIcon } from '@heroicons/vue/24/outline'
import { useAppStateStore } from '@/stores/appState'
import { useEnvironmentStore } from '@/stores/environment'
import { api } from '@/services/api'
import type { Variable } from '@/types'

const props = defineProps<{
  isOpen: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const appState = useAppStateStore()
const envStore = useEnvironmentStore()

const effectiveTheme = computed(() => appState.effectiveTheme)

// Track if a nested modal (confirmation dialog) is open
const isNestedModalOpen = ref(false)
const environments = computed(() => envStore.environments)

const selectedEnvId = ref<number | null>(null)
const currentName = ref('')
const currentVariables = ref<Variable[]>([])

// Load variables when selection changes
watch(selectedEnvId, async (id) => {
  if (id === null) {
    currentName.value = ''
    currentVariables.value = [...envStore.globalVariables.map(v => ({ ...v }))]
  } else {
    const env = environments.value.find(e => e.id === id)
    if (env) {
      currentName.value = env.name
      currentVariables.value = env.variables.map(v => ({ ...v }))
    }
  }
})

// Initialize when modal opens
watch(() => props.isOpen, (isOpen, wasOpen) => {
  if (isOpen === wasOpen) return
  if (isOpen) {
    appState.addModalOpen()
    selectedEnvId.value = null
    currentVariables.value = [...envStore.globalVariables.map(v => ({ ...v }))]
  } else {
    appState.removeModalOpen()
  }
})

function selectGlobals() {
  selectedEnvId.value = null
}

function selectEnv(id: number) {
  selectedEnvId.value = id
}

async function addEnvironment() {
  try {
    const env = await api.createEnvironment('New Environment', [])
    envStore.addEnvironment(env)
    selectedEnvId.value = env.id
  } catch (error) {
    console.error('Failed to create environment:', error)
  }
}

async function deleteEnv(id: number) {
  const modal = (window as any).$modal
  if (!modal) return
  
  const env = environments.value.find(e => e.id === id)
  const envName = env?.name || 'this environment'
  
  isNestedModalOpen.value = true
  const confirmed = await modal.confirm({
    title: 'Delete Environment',
    message: `Are you sure you want to delete "${envName}"? All variables in this environment will be lost.`,
    confirmText: 'Delete',
    danger: true,
  })
  isNestedModalOpen.value = false
  
  if (!confirmed) return
  
  try {
    await api.deleteEnvironment(id)
    envStore.deleteEnvironment(id)
    if (selectedEnvId.value === id) {
      selectedEnvId.value = null
    }
  } catch (error) {
    console.error('Failed to delete environment:', error)
  }
}

function addVariable() {
  currentVariables.value.push({ key: '', value: '', secret: false })
}

async function removeVariable(index: number, event?: Event) {
  // Stop event propagation to prevent closing the parent modal
  event?.stopPropagation()
  
  const variable = currentVariables.value[index]
  
  // Only show confirmation if the variable has a key (non-empty)
  if (variable.key.trim()) {
    const modal = (window as any).$modal
    if (modal) {
      isNestedModalOpen.value = true
      const confirmed = await modal.confirm({
        title: 'Delete Variable',
        message: `Are you sure you want to delete the variable "${variable.key}"?`,
        confirmText: 'Delete',
        danger: true,
      })
      isNestedModalOpen.value = false
      
      if (!confirmed) return
    }
  }
  
  currentVariables.value.splice(index, 1)
  saveCurrentEnv()
}

function getDuplicateKeys(): Set<string> {
  const keys = currentVariables.value.map(v => v.key.trim().toLowerCase()).filter(k => k)
  const seen = new Set<string>()
  const duplicates = new Set<string>()
  for (const key of keys) {
    if (seen.has(key)) {
      duplicates.add(key)
    }
    seen.add(key)
  }
  return duplicates
}

function isDuplicate(key: string): boolean {
  const trimmedKey = key.trim().toLowerCase()
  if (!trimmedKey) return false
  const duplicates = getDuplicateKeys()
  return duplicates.has(trimmedKey)
}

async function saveCurrentEnv() {
  try {
    // Check for duplicate variable names
    const duplicates = getDuplicateKeys()
    if (duplicates.size > 0) {
      const toast = (window as any).$toast
      if (toast) {
        toast.error(`Duplicate variable names: ${[...duplicates].join(', ')}`)
      }
      return
    }
    
    if (selectedEnvId.value === null) {
      // Save global variables
      const vars = currentVariables.value.filter(v => v.key.trim())
      await api.updateGlobalVariables(vars)
      envStore.setGlobalVariables(vars)
    } else {
      // Save environment
      const env = environments.value.find(e => e.id === selectedEnvId.value)
      if (env) {
        const updated = {
          ...env,
          name: currentName.value,
          variables: currentVariables.value.filter(v => v.key.trim()),
        }
        await api.updateEnvironment(updated)
        envStore.updateEnvironment(updated)
      }
    }
  } catch (error) {
    console.error('Failed to save:', error)
  }
}

function close() {
  // Don't close if a nested modal (confirmation dialog) is open
  console.log('[EnvironmentModal] close() called, isNestedModalOpen:', isNestedModalOpen.value)
  if (isNestedModalOpen.value) {
    console.log('[EnvironmentModal] Blocked: nested modal is open')
    return
  }
  console.log('[EnvironmentModal] Closing...')
  emit('close')
}

// Handle ESC key: let it pass to confirm dialog, but prevent this modal from closing
// The close() function will block closing when isNestedModalOpen is true
onMounted(() => {
  // No keyboard listener needed - just rely on close() blocking logic
})

onUnmounted(() => {
  // Cleanup if needed
})
</script>

