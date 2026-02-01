<template>
  <div class="flex flex-col h-full">
    <!-- Search and actions -->
    <div class="p-2 flex gap-2">
      <div class="flex-1 relative">
        <MagnifyingGlassIcon class="absolute left-2 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search..."
          class="w-full pl-8 pr-2 py-1.5 text-sm rounded-md border outline-none"
          :class="[
            effectiveTheme === 'dark'
              ? 'bg-dark-base border-dark-border text-white placeholder-gray-500 focus:border-accent'
              : 'bg-white border-light-border text-gray-900 placeholder-gray-400 focus:border-accent'
          ]"
        />
      </div>
      <button
        @click="createCollection"
        class="p-1.5 rounded-md transition-colors"
        :class="effectiveTheme === 'dark' ? 'hover:bg-dark-hover text-gray-400' : 'hover:bg-light-hover text-gray-500'"
        title="New Collection"
      >
        <PlusIcon class="w-5 h-5" />
      </button>
      <button
        @click="createFolder"
        class="p-1.5 rounded-md transition-colors"
        :class="effectiveTheme === 'dark' ? 'hover:bg-dark-hover text-gray-400' : 'hover:bg-light-hover text-gray-500'"
        title="New Folder"
      >
        <FolderPlusIcon class="w-5 h-5" />
      </button>
    </div>
    
    <!-- Tree -->
    <div class="flex-1 overflow-auto px-2 pb-2">
      <div v-if="filteredTree.length === 0" class="text-center py-8 text-gray-500 text-sm">
        No collections yet
      </div>
      
      <div v-for="col in filteredTree" :key="col.collection.id" class="mb-1">
        <!-- Collection -->
        <div
          class="flex items-center gap-1 px-2 py-1 rounded-md cursor-pointer group"
          :class="[
            effectiveTheme === 'dark' ? 'hover:bg-dark-hover' : 'hover:bg-light-hover',
            getDropIndicatorClass('collection', col.collection.id, dropTarget?.type === 'collection' && dropTarget?.id === col.collection.id ? dropTarget.position : 'inside')
          ]"
          draggable="true"
          @dragstart="(e) => onDragStartCollection(e, col.collection)"
          @dragover="(e) => onDragOverReorderable(e, 'collection', col.collection.id)"
          @dragleave="onDragLeave"
          @drop="(e) => onDropOnCollection(e, col.collection, dropTarget?.position || 'inside')"
          @dragend="onDragEnd"
          @click="toggleCollection(col.collection.id)"
          @contextmenu="(e) => showCollectionMenu(e, col.collection)"
        >
          <ChevronRightIcon 
            class="w-4 h-4 transition-transform flex-shrink-0"
            :class="{ 'rotate-90': isExpanded('collection', col.collection.id) }"
          />
          <FolderIcon class="w-4 h-4 text-accent flex-shrink-0" />
          <span class="flex-1 truncate text-sm" :class="effectiveTheme === 'dark' ? 'text-gray-200' : 'text-gray-800'">
            {{ col.collection.name }}
          </span>
        </div>
        
        <!-- Collection contents -->
        <div v-if="isExpanded('collection', col.collection.id)" class="ml-4">
          <!-- Folders -->
          <div v-for="folder in col.folders" :key="folder.folder.id" class="mb-1">
            <div
              class="flex items-center gap-1 px-2 py-1 rounded-md cursor-pointer group"
              :class="[
                effectiveTheme === 'dark' ? 'hover:bg-dark-hover' : 'hover:bg-light-hover',
                getDropIndicatorClass('folder', folder.folder.id, dropTarget?.type === 'folder' && dropTarget?.id === folder.folder.id ? dropTarget.position : 'inside')
              ]"
              draggable="true"
              @dragstart="(e) => onDragStartFolder(e, folder.folder, col.collection.id)"
              @dragover="(e) => onDragOverReorderable(e, 'folder', folder.folder.id)"
              @dragleave="onDragLeave"
              @drop="(e) => onDropOnFolder(e, folder.folder, col.collection.id, dropTarget?.position || 'inside')"
              @dragend="onDragEnd"
              @click="toggleFolder(folder.folder.id)"
              @contextmenu="(e) => showFolderMenu(e, folder.folder, col.collection.id)"
            >
              <ChevronRightIcon 
                class="w-4 h-4 transition-transform flex-shrink-0"
                :class="{ 'rotate-90': isExpanded('folder', folder.folder.id) }"
              />
              <FolderOpenIcon v-if="isExpanded('folder', folder.folder.id)" class="w-4 h-4 text-yellow-500 flex-shrink-0" />
              <FolderIcon v-else class="w-4 h-4 text-yellow-500 flex-shrink-0" />
              <span class="flex-1 truncate text-sm" :class="effectiveTheme === 'dark' ? 'text-gray-200' : 'text-gray-800'">
                {{ folder.folder.name }}
              </span>
            </div>
            
            <!-- Folder requests -->
            <div v-if="isExpanded('folder', folder.folder.id)" class="ml-4">
              <RequestItem
                v-for="req in folder.requests"
                :key="req.id"
                :request="req"
                :class="getRequestDropIndicatorClass(req.id)"
                draggable="true"
                @dragstart="(e: DragEvent) => onDragStartRequest(e, req)"
                @dragover="(e: DragEvent) => onDragOver(e, 'request', req.id, 'before')"
                @dragleave="onDragLeave"
                @drop="(e: DragEvent) => onDropOnRequest(e, req, 'before')"
                @dragend="onDragEnd"
                @click="openRequest(req)"
                @pin-request="openRequestInNewTab(req)"
                @contextmenu="(e: MouseEvent) => showRequestMenu(e, req)"
              />
            </div>
          </div>
          
          <!-- Direct requests -->
          <RequestItem
            v-for="req in col.requests"
            :key="req.id"
            :request="req"
            :class="getRequestDropIndicatorClass(req.id)"
            draggable="true"
            @dragstart="(e: DragEvent) => onDragStartRequest(e, req)"
            @dragover="(e: DragEvent) => onDragOver(e, 'request', req.id, 'before')"
            @dragleave="onDragLeave"
            @drop="(e: DragEvent) => onDropOnRequest(e, req, 'before')"
            @dragend="onDragEnd"
            @click="openRequest(req)"
            @pin-request="openRequestInNewTab(req)"
            @contextmenu="(e: MouseEvent) => showRequestMenu(e, req)"
          />
        </div>
      </div>
    </div>
    
    <!-- Context Menus -->
    <ContextMenu ref="collectionMenuRef" :items="collectionMenuItems" />
    <ContextMenu ref="folderMenuRef" :items="folderMenuItems" />
    <ContextMenu ref="requestMenuRef" :items="requestMenuItems" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { 
  MagnifyingGlassIcon, 
  PlusIcon, 
  FolderPlusIcon,
  ChevronRightIcon,
  FolderIcon,
  FolderOpenIcon,
  PencilIcon,
  TrashIcon,
  DocumentDuplicateIcon,
} from '@heroicons/vue/24/outline'
import { useAppStateStore } from '@/stores/appState'
import { useCollectionStore } from '@/stores/collection'
import { useTabsStore } from '@/stores/tabs'
import { api } from '@/services/api'
import type { Request, Collection, Folder, CollectionTree as CollectionTreeType, FolderTree } from '@/types'
import RequestItem from './RequestItem.vue'
import ContextMenu, { type ContextMenuItem } from '@/components/common/ContextMenu.vue'

// Drag and drop types
interface DragData {
  type: 'collection' | 'folder' | 'request'
  id: number
  collectionId?: number
  folderId?: number | null
}

// Single-click delay to allow double-click detection in RequestItem
let clickTimer: ReturnType<typeof setTimeout> | null = null
let pendingClickRequest: Request | null = null
const CLICK_DELAY = 300

const appState = useAppStateStore()
const collectionStore = useCollectionStore()
const tabsStore = useTabsStore()

const effectiveTheme = computed(() => appState.effectiveTheme)
const searchQuery = ref('')

// Drag and drop state
const dragData = ref<DragData | null>(null)
const dropTarget = ref<{ type: string; id: number; position: 'before' | 'after' | 'inside' } | null>(null)

// Context menu refs
const collectionMenuRef = ref<InstanceType<typeof ContextMenu> | null>(null)
const folderMenuRef = ref<InstanceType<typeof ContextMenu> | null>(null)
const requestMenuRef = ref<InstanceType<typeof ContextMenu> | null>(null)

// Current context menu targets
const selectedCollection = ref<Collection | null>(null)
const selectedFolder = ref<{ folder: Folder; collectionId: number } | null>(null)
const selectedRequest = ref<Request | null>(null)

// Context menu items
const collectionMenuItems = computed<ContextMenuItem[]>(() => [
  {
    id: 'add-folder',
    label: 'Add Folder',
    icon: FolderPlusIcon,
    action: () => addFolderToCollection(),
  },
  {
    id: 'add-request',
    label: 'Add Request',
    icon: PlusIcon,
    action: () => addRequestToCollection(),
  },
  {
    id: 'rename',
    label: 'Rename',
    icon: PencilIcon,
    action: () => renameCollection(),
  },
  {
    id: 'delete',
    label: 'Delete',
    icon: TrashIcon,
    danger: true,
    action: () => deleteCollection(),
  },
])

const folderMenuItems = computed<ContextMenuItem[]>(() => [
  {
    id: 'add-request',
    label: 'Add Request',
    icon: PlusIcon,
    action: () => addRequestToFolder(),
  },
  {
    id: 'rename',
    label: 'Rename',
    icon: PencilIcon,
    action: () => renameFolder(),
  },
  {
    id: 'delete',
    label: 'Delete',
    icon: TrashIcon,
    danger: true,
    action: () => deleteFolder(),
  },
])

const requestMenuItems = computed<ContextMenuItem[]>(() => [
  {
    id: 'open-new-tab',
    label: 'Open in New Tab',
    icon: DocumentDuplicateIcon,
    action: () => {
      if (selectedRequest.value) {
        openRequestInNewTab(selectedRequest.value)
      }
    },
  },
  {
    id: 'duplicate',
    label: 'Duplicate',
    icon: DocumentDuplicateIcon,
    action: () => duplicateRequest(),
  },
  {
    id: 'rename',
    label: 'Rename',
    icon: PencilIcon,
    action: () => renameRequest(),
  },
  {
    id: 'delete',
    label: 'Delete',
    icon: TrashIcon,
    danger: true,
    action: () => deleteRequest(),
  },
])

const filteredTree = computed(() => {
  if (!searchQuery.value) return collectionStore.tree
  
  const query = searchQuery.value.toLowerCase()
  return collectionStore.tree.filter(col => {
    // Check collection name
    if (col.collection.name.toLowerCase().includes(query)) return true
    
    // Check folder names
    if (col.folders.some(f => f.folder.name.toLowerCase().includes(query))) return true
    
    // Check request names/URLs
    if (col.requests.some(r => 
      r.name.toLowerCase().includes(query) || 
      r.url.toLowerCase().includes(query)
    )) return true
    
    // Check folder requests
    if (col.folders.some(f => 
      f.requests.some(r => 
        r.name.toLowerCase().includes(query) || 
        r.url.toLowerCase().includes(query)
      )
    )) return true
    
    return false
  })
})

function isExpanded(type: string, id: number) {
  return appState.isSidebarItemExpanded(type, id)
}

function toggleCollection(id: number) {
  appState.toggleSidebarItem('collection', id)
}

function toggleFolder(id: number) {
  appState.toggleSidebarItem('folder', id)
}

function openRequest(req: Request) {
  // Delay single-click to allow double-click detection in RequestItem
  if (clickTimer) {
    clearTimeout(clickTimer)
    clickTimer = null
  }
  pendingClickRequest = req
  clickTimer = setTimeout(() => {
    if (pendingClickRequest) {
      tabsStore.previewRequest(
        pendingClickRequest.id,
        pendingClickRequest.name,
        pendingClickRequest.method,
        pendingClickRequest.url,
        pendingClickRequest.headers,
        pendingClickRequest.params,
        pendingClickRequest.body,
        pendingClickRequest.bodyType
      )
      appState.highlightedRequestId = pendingClickRequest.id
      pendingClickRequest = null
    }
    clickTimer = null
  }, CLICK_DELAY)
}

function openRequestInNewTab(req: Request) {
  // Cancel pending single-click since we're pinning
  if (clickTimer) {
    clearTimeout(clickTimer)
    clickTimer = null
    pendingClickRequest = null
  }
  
  // Check if already open
  const existing = tabsStore.tabs.find(t => t.requestId === req.id)
  if (existing) {
    // Switch to the tab and pin it if it's a preview
    tabsStore.setActiveTab(existing.id)
    if (existing.isPreview) {
      tabsStore.pinTab(existing.id)
    }
    appState.highlightedRequestId = req.id
    return
  }
  
  // Not opened, open as pinned (non-preview)
  tabsStore.openRequest(
    req.id,
    req.name,
    req.method,
    req.url,
    req.headers,
    req.params,
    req.body,
    req.bodyType
  )
  appState.highlightedRequestId = req.id
}

async function createCollection() {
  const modal = (window as any).$modal
  if (!modal) return
  
  const name = await modal.input({
    title: 'New Collection',
    placeholder: 'Collection name',
    confirmText: 'Create',
  })
  
  if (name && name.trim()) {
    try {
      const collection = await api.createCollection(name.trim())
      collectionStore.addCollection(collection)
    } catch (error) {
      console.error('Failed to create collection:', error)
    }
  }
}

async function createFolder() {
  const modal = (window as any).$modal
  if (!modal) return
  
  // Need to select a collection first
  if (collectionStore.tree.length === 0) {
    await modal.confirm({
      title: 'No Collections',
      message: 'Please create a collection first before adding folders.',
      confirmText: 'OK',
    })
    return
  }
  
  let targetCollectionId: number
  
  if (collectionStore.tree.length === 1) {
    // Only one collection, use it directly
    targetCollectionId = collectionStore.tree[0].collection.id
  } else {
    // Multiple collections - let user choose
    const options = collectionStore.tree.map(c => ({ 
      value: c.collection.id, 
      label: c.collection.name 
    }))
    const selected = await modal.select({ 
      title: 'Select Collection', 
      options 
    })
    if (selected === null) return
    targetCollectionId = selected
  }
  
  const name = await modal.input({
    title: 'New Folder',
    placeholder: 'Folder name',
    confirmText: 'Create',
  })
  
  if (name && name.trim()) {
    try {
      const folder = await api.createFolder(targetCollectionId, name.trim())
      collectionStore.addFolder(targetCollectionId, folder)
      // Expand the collection to show the new folder
      appState.expandSidebarItem('collection', targetCollectionId)
    } catch (error) {
      console.error('Failed to create folder:', error)
    }
  }
}

// Context menu handlers
function showCollectionMenu(e: MouseEvent, collection: Collection) {
  e.preventDefault()
  selectedCollection.value = collection
  collectionMenuRef.value?.open(e)
}

function showFolderMenu(e: MouseEvent, folder: Folder, collectionId: number) {
  e.preventDefault()
  selectedFolder.value = { folder, collectionId }
  folderMenuRef.value?.open(e)
}

function showRequestMenu(e: MouseEvent, request: Request) {
  e.preventDefault()
  selectedRequest.value = request
  requestMenuRef.value?.open(e)
}

async function addFolderToCollection() {
  if (!selectedCollection.value) return
  const modal = (window as any).$modal
  if (!modal) return
  
  const name = await modal.input({
    title: 'New Folder',
    placeholder: 'Folder name',
    confirmText: 'Create',
  })
  
  if (name && name.trim()) {
    try {
      const folder = await api.createFolder(selectedCollection.value.id, name.trim())
      collectionStore.addFolder(selectedCollection.value.id, folder)
    } catch (error) {
      console.error('Failed to create folder:', error)
    }
  }
}

async function addRequestToCollection() {
  if (!selectedCollection.value) return
  try {
    const request = await api.createRequest({
      collectionId: selectedCollection.value.id,
      name: 'New Request',
      method: 'GET',
      url: '',
    })
    collectionStore.addRequest(request)
    tabsStore.openRequest(
      request.id,
      request.name,
      request.method,
      request.url,
      request.headers,
      request.params,
      request.body,
      request.bodyType
    )
  } catch (error) {
    console.error('Failed to create request:', error)
  }
}

async function addRequestToFolder() {
  if (!selectedFolder.value) return
  try {
    const request = await api.createRequest({
      collectionId: selectedFolder.value.collectionId,
      folderId: selectedFolder.value.folder.id,
      name: 'New Request',
      method: 'GET',
      url: '',
    })
    collectionStore.addRequest(request)
    tabsStore.openRequest(
      request.id,
      request.name,
      request.method,
      request.url,
      request.headers,
      request.params,
      request.body,
      request.bodyType
    )
  } catch (error) {
    console.error('Failed to create request:', error)
  }
}

async function renameCollection() {
  if (!selectedCollection.value) return
  const modal = (window as any).$modal
  if (!modal) return
  
  const name = await modal.input({
    title: 'Rename Collection',
    value: selectedCollection.value.name,
    confirmText: 'Save',
  })
  
  if (name && name.trim() && name.trim() !== selectedCollection.value.name) {
    try {
      const updated = { ...selectedCollection.value, name: name.trim() }
      await api.updateCollection(updated)
      collectionStore.updateCollection(updated)
    } catch (error) {
      console.error('Failed to rename collection:', error)
    }
  }
}

async function deleteCollection() {
  if (!selectedCollection.value) return
  const modal = (window as any).$modal
  if (!modal) return
  
  const confirmed = await modal.confirm({
    title: 'Delete Collection',
    message: `Are you sure you want to delete "${selectedCollection.value.name}"? This will delete all folders and requests inside it.`,
    confirmText: 'Delete',
    danger: true,
  })
  
  if (confirmed) {
    try {
      const collectionName = selectedCollection.value.name
      await api.deleteCollection(selectedCollection.value.id)
      collectionStore.deleteCollection(selectedCollection.value.id)

      const toast = (window as any).$toast
      if (toast) {
        toast.success(`Collection "${collectionName}" deleted`)
      }
    } catch (error) {
      console.error('Failed to delete collection:', error)
      const toast = (window as any).$toast
      if (toast) {
        toast.error('Failed to delete collection')
      }
    }
  }
}

async function renameFolder() {
  if (!selectedFolder.value) return
  const modal = (window as any).$modal
  if (!modal) return
  
  const name = await modal.input({
    title: 'Rename Folder',
    value: selectedFolder.value.folder.name,
    confirmText: 'Save',
  })
  
  if (name && name.trim() && name.trim() !== selectedFolder.value.folder.name) {
    try {
      const updated = { ...selectedFolder.value.folder, name: name.trim() }
      await api.updateFolder(updated)
      collectionStore.updateFolder(updated)
    } catch (error) {
      console.error('Failed to rename folder:', error)
    }
  }
}

async function deleteFolder() {
  if (!selectedFolder.value) return
  const modal = (window as any).$modal
  if (!modal) return
  
  const confirmed = await modal.confirm({
    title: 'Delete Folder',
    message: `Are you sure you want to delete "${selectedFolder.value.folder.name}"? This will delete all requests inside it.`,
    confirmText: 'Delete',
    danger: true,
  })
  
  if (confirmed) {
    try {
      const folderName = selectedFolder.value.folder.name
      await api.deleteFolder(selectedFolder.value.folder.id)
      collectionStore.deleteFolder(selectedFolder.value.folder.id)

      const toast = (window as any).$toast
      if (toast) {
        toast.success(`Folder "${folderName}" deleted`)
      }
    } catch (error) {
      console.error('Failed to delete folder:', error)
      const toast = (window as any).$toast
      if (toast) {
        toast.error('Failed to delete folder')
      }
    }
  }
}

async function renameRequest() {
  if (!selectedRequest.value) return
  const modal = (window as any).$modal
  if (!modal) return
  
  const name = await modal.input({
    title: 'Rename Request',
    value: selectedRequest.value.name,
    confirmText: 'Save',
  })
  
  if (name && name.trim() && name.trim() !== selectedRequest.value.name) {
    try {
      const updated = { ...selectedRequest.value, name: name.trim() }
      await api.updateRequest(updated)
      collectionStore.updateRequest(updated)
      // Sync tab title
      tabsStore.updateTabTitleByRequestId(updated.id, updated.name)
    } catch (error) {
      console.error('Failed to rename request:', error)
    }
  }
}

async function deleteRequest() {
  if (!selectedRequest.value) return
  const modal = (window as any).$modal
  if (!modal) return
  
  const confirmed = await modal.confirm({
    title: 'Delete Request',
    message: `Are you sure you want to delete "${selectedRequest.value.name}"?`,
    confirmText: 'Delete',
    danger: true,
  })
  
  if (confirmed) {
    try {
      const requestId = selectedRequest.value.id
      const requestName = selectedRequest.value.name
      await api.deleteRequest(requestId)
      collectionStore.deleteRequest(requestId)
      // Close the tab if it's open
      tabsStore.closeTabByRequestId(requestId)

      const toast = (window as any).$toast
      if (toast) {
        toast.success(`Request "${requestName}" deleted`)
      }
    } catch (error) {
      console.error('Failed to delete request:', error)
      const toast = (window as any).$toast
      if (toast) {
        toast.error('Failed to delete request')
      }
    }
  }
}

async function duplicateRequest() {
  if (!selectedRequest.value) return

  try {
    const duplicated = await api.duplicateRequest(selectedRequest.value.id)
    collectionStore.addRequest(duplicated)

    // Open the duplicated request in a new tab
    tabsStore.openRequest(
      duplicated.id,
      duplicated.name,
      duplicated.method,
      duplicated.url,
      duplicated.headers,
      duplicated.params,
      duplicated.body,
      duplicated.bodyType
    )

    const toast = (window as any).$toast
    if (toast) {
      toast.success(`Request duplicated as "${duplicated.name}"`)
    }
  } catch (error) {
    console.error('Failed to duplicate request:', error)
    const toast = (window as any).$toast
    if (toast) {
      toast.error('Failed to duplicate request')
    }
  }
}

// Drag and drop handlers
function onDragStartCollection(e: DragEvent, collection: Collection) {
  dragData.value = { type: 'collection', id: collection.id }
  e.dataTransfer!.effectAllowed = 'move'
  e.dataTransfer!.setData('text/plain', JSON.stringify(dragData.value))
}

function onDragStartFolder(e: DragEvent, folder: Folder, collectionId: number) {
  dragData.value = { type: 'folder', id: folder.id, collectionId }
  e.dataTransfer!.effectAllowed = 'move'
  e.dataTransfer!.setData('text/plain', JSON.stringify(dragData.value))
}

function onDragStartRequest(e: DragEvent, request: Request) {
  dragData.value = { 
    type: 'request', 
    id: request.id, 
    collectionId: request.collectionId, 
    folderId: request.folderId 
  }
  e.dataTransfer!.effectAllowed = 'move'
  e.dataTransfer!.setData('text/plain', JSON.stringify(dragData.value))
}

function onDragOver(e: DragEvent, targetType: string, targetId: number, defaultPosition: 'before' | 'after' | 'inside') {
  e.preventDefault()
  e.dataTransfer!.dropEffect = 'move'
  
  // For requests, calculate position based on cursor position
  let position = defaultPosition
  if (targetType === 'request' && defaultPosition !== 'inside') {
    const rect = (e.currentTarget as HTMLElement).getBoundingClientRect()
    const y = e.clientY - rect.top
    position = y < rect.height / 2 ? 'before' : 'after'
  }
  
  dropTarget.value = { type: targetType, id: targetId, position }
}

function onDragOverReorderable(e: DragEvent, targetType: string, targetId: number) {
  e.preventDefault()
  e.dataTransfer!.dropEffect = 'move'
  
  // Allow reordering for collections/folders when dragging same type
  if (dragData.value?.type === targetType) {
    const rect = (e.currentTarget as HTMLElement).getBoundingClientRect()
    const y = e.clientY - rect.top
    const position = y < rect.height / 2 ? 'before' : 'after'
    dropTarget.value = { type: targetType, id: targetId, position }
  } else {
    // For moving requests into collections/folders
    dropTarget.value = { type: targetType, id: targetId, position: 'inside' }
  }
}

function onDragLeave() {
  dropTarget.value = null
}

function onDragEnd() {
  dragData.value = null
  dropTarget.value = null
}

async function onDropOnCollection(e: DragEvent, targetCollection: Collection, position: 'before' | 'after' | 'inside') {
  e.preventDefault()
  if (!dragData.value) return
  
  try {
    if (dragData.value.type === 'collection' && position !== 'inside') {
      // Reorder collections
      const collections = collectionStore.tree.map(c => c.collection)
      const fromIndex = collections.findIndex(c => c.id === dragData.value!.id)
      const toIndex = collections.findIndex(c => c.id === targetCollection.id)
      if (fromIndex === -1 || toIndex === -1 || fromIndex === toIndex) return
      
      const newOrder = [...collections]
      const [moved] = newOrder.splice(fromIndex, 1)
      const insertIndex = position === 'before' ? toIndex : toIndex + 1
      newOrder.splice(insertIndex > fromIndex ? insertIndex - 1 : insertIndex, 0, moved)
      
      await api.reorderCollections(newOrder.map(c => c.id))
      await collectionStore.loadTree()
    } else if (dragData.value.type === 'folder' && position === 'inside') {
      // Move folder to another collection
      if (dragData.value.collectionId !== targetCollection.id) {
        await api.moveFolder(dragData.value.id, targetCollection.id)
        await collectionStore.loadTree()
      }
    } else if (dragData.value.type === 'request' && position === 'inside') {
      // Move request to collection root
      await api.moveRequest(dragData.value.id, targetCollection.id, null)
      await collectionStore.loadTree()
    }
  } catch (error) {
    console.error('Failed to drop:', error)
  }
  
  onDragEnd()
}

async function onDropOnFolder(e: DragEvent, folder: Folder, collectionId: number, position: 'before' | 'after' | 'inside') {
  e.preventDefault()
  if (!dragData.value) return
  
  try {
    if (dragData.value.type === 'folder' && position !== 'inside') {
      // Reorder folders within collection
      const col = collectionStore.tree.find(c => c.collection.id === collectionId)
      if (!col) return
      
      const folders = col.folders.map(f => f.folder)
      const fromIndex = folders.findIndex(f => f.id === dragData.value!.id)
      const toIndex = folders.findIndex(f => f.id === folder.id)
      if (fromIndex === -1 || toIndex === -1 || fromIndex === toIndex) return
      
      const newOrder = [...folders]
      const [moved] = newOrder.splice(fromIndex, 1)
      const insertIndex = position === 'before' ? toIndex : toIndex + 1
      newOrder.splice(insertIndex > fromIndex ? insertIndex - 1 : insertIndex, 0, moved)
      
      await api.reorderFolders(collectionId, newOrder.map(f => f.id))
      await collectionStore.loadTree()
    } else if (dragData.value.type === 'request' && position === 'inside') {
      // Move request to folder
      await api.moveRequest(dragData.value.id, collectionId, folder.id)
      await collectionStore.loadTree()
    }
  } catch (error) {
    console.error('Failed to drop:', error)
  }
  
  onDragEnd()
}

async function onDropOnRequest(e: DragEvent, targetRequest: Request, _defaultPosition: 'before' | 'after') {
  e.preventDefault()
  if (!dragData.value || dragData.value.type !== 'request') return
  
  // Use the position from dropTarget which was calculated in onDragOver
  const position = dropTarget.value?.position === 'after' ? 'after' : 'before'
  
  try {
    // Get the list of requests in the same container
    const col = collectionStore.tree.find(c => c.collection.id === targetRequest.collectionId)
    if (!col) return
    
    let requests: Request[]
    if (targetRequest.folderId) {
      const folder = col.folders.find(f => f.folder.id === targetRequest.folderId)
      if (!folder) return
      requests = folder.requests
    } else {
      requests = col.requests
    }
    
    const fromIndex = requests.findIndex(r => r.id === dragData.value!.id)
    const toIndex = requests.findIndex(r => r.id === targetRequest.id)
    
    // If moving from different container, first move the request
    if (fromIndex === -1) {
      await api.moveRequest(dragData.value.id, targetRequest.collectionId, targetRequest.folderId)
      await collectionStore.loadTree()
      // Get updated requests list
      const updatedCol = collectionStore.tree.find(c => c.collection.id === targetRequest.collectionId)
      if (!updatedCol) return
      
      if (targetRequest.folderId) {
        const folder = updatedCol.folders.find(f => f.folder.id === targetRequest.folderId)
        if (!folder) return
        requests = folder.requests
      } else {
        requests = updatedCol.requests
      }
    }
    
    // Now reorder
    const newFromIndex = requests.findIndex(r => r.id === dragData.value!.id)
    const newToIndex = requests.findIndex(r => r.id === targetRequest.id)
    if (newFromIndex === -1 || newToIndex === -1 || newFromIndex === newToIndex) return
    
    const newOrder = [...requests]
    const [moved] = newOrder.splice(newFromIndex, 1)
    const insertIndex = position === 'before' ? newToIndex : newToIndex + 1
    newOrder.splice(insertIndex > newFromIndex ? insertIndex - 1 : insertIndex, 0, moved)
    
    await api.reorderRequests(targetRequest.collectionId, targetRequest.folderId, newOrder.map(r => r.id))
    await collectionStore.loadTree()
  } catch (error) {
    console.error('Failed to drop:', error)
  }
  
  onDragEnd()
}

function getDropIndicatorClass(type: string, id: number, position: 'before' | 'after' | 'inside') {
  if (!dropTarget.value) return ''
  if (dropTarget.value.type === type && dropTarget.value.id === id && dropTarget.value.position === position) {
    if (position === 'inside') {
      return 'ring-2 ring-accent ring-inset'
    }
    return position === 'before' ? 'border-t-2 border-accent' : 'border-b-2 border-accent'
  }
  return ''
}

function getRequestDropIndicatorClass(id: number) {
  if (!dropTarget.value) return ''
  if (dropTarget.value.type === 'request' && dropTarget.value.id === id) {
    return dropTarget.value.position === 'before' ? 'border-t-2 border-accent' : 'border-b-2 border-accent'
  }
  return ''
}
</script>
