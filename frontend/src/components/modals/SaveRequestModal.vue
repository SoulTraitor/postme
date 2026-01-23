<template>
  <TransitionRoot :show="isOpen" as="template">
    <Dialog as="div" class="relative z-50" @close="close">
      <TransitionChild
        enter="ease-out duration-200"
        enter-from="opacity-0"
        enter-to="opacity-100"
        leave="ease-in duration-150"
        leave-from="opacity-100"
        leave-to="opacity-0"
      >
        <div class="fixed inset-0 bg-black/50" />
      </TransitionChild>

      <div class="fixed inset-0 overflow-y-auto">
        <div class="flex min-h-full items-center justify-center p-4">
          <TransitionChild
            enter="ease-out duration-200"
            enter-from="opacity-0 scale-95"
            enter-to="opacity-100 scale-100"
            leave="ease-in duration-150"
            leave-from="opacity-100 scale-100"
            leave-to="opacity-0 scale-95"
          >
            <DialogPanel 
              class="w-full max-w-md rounded-lg p-6 shadow-xl"
              :class="effectiveTheme === 'dark' ? 'bg-dark-elevated' : 'bg-white'"
            >
              <DialogTitle 
                class="text-lg font-medium mb-6"
                :class="effectiveTheme === 'dark' ? 'text-white' : 'text-gray-900'"
              >
                Save Request
              </DialogTitle>
              
              <div class="space-y-4">
                <!-- Request Name -->
                <div>
                  <label 
                    class="block text-sm font-medium mb-2"
                    :class="effectiveTheme === 'dark' ? 'text-gray-300' : 'text-gray-700'"
                  >
                    Name
                  </label>
                  <input
                    v-model="requestName"
                    type="text"
                    placeholder="Request name"
                    class="w-full px-3 py-2 rounded-md border outline-none text-sm"
                    :class="[
                      effectiveTheme === 'dark'
                        ? 'bg-dark-surface border-dark-border text-white placeholder-gray-500 focus:border-accent'
                        : 'bg-white border-light-border text-gray-900 placeholder-gray-400 focus:border-accent'
                    ]"
                  />
                </div>
                
                <!-- Collection Selector -->
                <div>
                  <label 
                    class="block text-sm font-medium mb-2"
                    :class="effectiveTheme === 'dark' ? 'text-gray-300' : 'text-gray-700'"
                  >
                    Collection
                  </label>
                  <Listbox v-model="selectedCollection">
                    <div class="relative">
                      <ListboxButton
                        class="w-full px-3 py-2 rounded-md border outline-none text-sm text-left flex items-center justify-between"
                        :class="[
                          effectiveTheme === 'dark'
                            ? 'bg-dark-surface border-dark-border text-white'
                            : 'bg-white border-light-border text-gray-900'
                        ]"
                      >
                        <span>{{ selectedCollection?.name || 'Select collection' }}</span>
                        <ChevronUpDownIcon class="w-4 h-4 text-gray-400" />
                      </ListboxButton>
                      
                      <TransitionRoot
                        leave="transition ease-in duration-100"
                        leave-from="opacity-100"
                        leave-to="opacity-0"
                      >
                        <ListboxOptions
                          class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md py-1 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none text-sm"
                          :class="effectiveTheme === 'dark' ? 'bg-dark-elevated' : 'bg-white'"
                        >
                          <ListboxOption
                            v-for="col in collectionStore.collections"
                            :key="col.id"
                            :value="col"
                            v-slot="{ active, selected }"
                          >
                            <li
                              class="cursor-pointer select-none px-3 py-2"
                              :class="[
                                active ? 'bg-accent text-white' : (effectiveTheme === 'dark' ? 'text-gray-300' : 'text-gray-900'),
                                selected ? 'font-medium' : ''
                              ]"
                            >
                              {{ col.name }}
                            </li>
                          </ListboxOption>
                        </ListboxOptions>
                      </TransitionRoot>
                    </div>
                  </Listbox>
                </div>
                
                <!-- Folder Selector (Optional) -->
                <div v-if="availableFolders.length > 0">
                  <label 
                    class="block text-sm font-medium mb-2"
                    :class="effectiveTheme === 'dark' ? 'text-gray-300' : 'text-gray-700'"
                  >
                    Folder (optional)
                  </label>
                  <Listbox v-model="selectedFolder">
                    <div class="relative">
                      <ListboxButton
                        class="w-full px-3 py-2 rounded-md border outline-none text-sm text-left flex items-center justify-between"
                        :class="[
                          effectiveTheme === 'dark'
                            ? 'bg-dark-surface border-dark-border text-white'
                            : 'bg-white border-light-border text-gray-900'
                        ]"
                      >
                        <span>{{ selectedFolder?.name || 'None' }}</span>
                        <ChevronUpDownIcon class="w-4 h-4 text-gray-400" />
                      </ListboxButton>
                      
                      <TransitionRoot
                        leave="transition ease-in duration-100"
                        leave-from="opacity-100"
                        leave-to="opacity-0"
                      >
                        <ListboxOptions
                          class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md py-1 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none text-sm"
                          :class="effectiveTheme === 'dark' ? 'bg-dark-elevated' : 'bg-white'"
                        >
                          <ListboxOption
                            :value="null"
                            v-slot="{ active }"
                          >
                            <li
                              class="cursor-pointer select-none px-3 py-2"
                              :class="active ? 'bg-accent text-white' : (effectiveTheme === 'dark' ? 'text-gray-300' : 'text-gray-900')"
                            >
                              None
                            </li>
                          </ListboxOption>
                          <ListboxOption
                            v-for="folder in availableFolders"
                            :key="folder.id"
                            :value="folder"
                            v-slot="{ active, selected }"
                          >
                            <li
                              class="cursor-pointer select-none px-3 py-2"
                              :class="[
                                active ? 'bg-accent text-white' : (effectiveTheme === 'dark' ? 'text-gray-300' : 'text-gray-900'),
                                selected ? 'font-medium' : ''
                              ]"
                            >
                              {{ folder.name }}
                            </li>
                          </ListboxOption>
                        </ListboxOptions>
                      </TransitionRoot>
                    </div>
                  </Listbox>
                </div>
              </div>
              
              <div class="flex justify-end gap-3 mt-6">
                <button
                  @click="close"
                  class="px-4 py-2 rounded-md font-medium transition-colors"
                  :class="effectiveTheme === 'dark' ? 'bg-dark-hover text-gray-300 hover:bg-dark-border' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'"
                >
                  Cancel
                </button>
                <button
                  @click="save"
                  :disabled="!canSave"
                  class="px-4 py-2 rounded-md font-medium text-white transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                  :class="canSave ? 'bg-accent hover:bg-accent-hover' : 'bg-accent'"
                >
                  Save
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
import { ref, computed, watch } from 'vue'
import { 
  Dialog, 
  DialogPanel, 
  DialogTitle, 
  TransitionRoot, 
  TransitionChild,
  Listbox,
  ListboxButton,
  ListboxOptions,
  ListboxOption,
} from '@headlessui/vue'
import { ChevronUpDownIcon } from '@heroicons/vue/24/outline'
import { useAppStateStore } from '@/stores/appState'
import { useCollectionStore } from '@/stores/collection'
import { useTabsStore } from '@/stores/tabs'
import { api } from '@/services/api'
import type { Collection, Folder } from '@/types'

const props = defineProps<{
  isOpen: boolean
  tabId?: string
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'saved', requestId: number): void
}>()

const appState = useAppStateStore()
const collectionStore = useCollectionStore()
const tabsStore = useTabsStore()

const effectiveTheme = computed(() => appState.effectiveTheme)

const requestName = ref('')
const selectedCollection = ref<Collection | null>(null)
const selectedFolder = ref<Folder | null>(null)

const availableFolders = computed(() => {
  if (!selectedCollection.value) return []
  const tree = collectionStore.tree.find(t => t.collection.id === selectedCollection.value?.id)
  return tree?.folders.map(f => f.folder) || []
})

const canSave = computed(() => {
  return requestName.value.trim() && selectedCollection.value
})

// Reset when modal opens
watch(() => props.isOpen, (isOpen) => {
  if (isOpen) {
    const tab = props.tabId ? tabsStore.getTab(props.tabId) : tabsStore.activeTab
    requestName.value = tab?.title !== 'Untitled' ? tab?.title || '' : ''
    
    // Pre-select first collection if available
    if (collectionStore.collections.length > 0) {
      selectedCollection.value = collectionStore.collections[0]
    } else {
      selectedCollection.value = null
    }
    selectedFolder.value = null
  }
})

// Clear folder when collection changes
watch(selectedCollection, () => {
  selectedFolder.value = null
})

function close() {
  emit('close')
}

async function save() {
  if (!canSave.value || !selectedCollection.value) return
  
  const tab = props.tabId ? tabsStore.getTab(props.tabId) : tabsStore.activeTab
  if (!tab) return
  
  try {
    const request = await api.createRequest({
      collectionId: selectedCollection.value.id,
      folderId: selectedFolder.value?.id || null,
      name: requestName.value.trim(),
      method: tab.method,
      url: tab.url,
      headers: tab.headers,
      params: tab.params,
      body: tab.body,
      bodyType: tab.bodyType,
    })
    
    // Add to store
    collectionStore.addRequest(request)
    
    // Update tab
    tabsStore.markSaved(tab.id, request.id, request.name)
    
    emit('saved', request.id)
    close()
  } catch (error) {
    console.error('Failed to save request:', error)
  }
}
</script>
