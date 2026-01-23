import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { CollectionTree, Collection, Folder, Request } from '@/types'

// Import api lazily to avoid circular dependency
let apiModule: typeof import('@/services/api') | null = null
async function getApi() {
  if (!apiModule) {
    apiModule = await import('@/services/api')
  }
  return apiModule.api
}

export const useCollectionStore = defineStore('collection', () => {
  const tree = ref<CollectionTree[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Flat list of all collections
  const collections = computed(() => tree.value.map(t => t.collection))

  // Get collection by ID
  function getCollection(id: number): Collection | undefined {
    return tree.value.find(t => t.collection.id === id)?.collection
  }

  // Get folder by ID
  function getFolder(id: number): Folder | undefined {
    for (const col of tree.value) {
      const folder = col.folders.find(f => f.folder.id === id)?.folder
      if (folder) return folder
    }
    return undefined
  }

  // Get request by ID
  function getRequest(id: number): Request | undefined {
    for (const col of tree.value) {
      // Check direct requests
      const directReq = col.requests.find(r => r.id === id)
      if (directReq) return directReq

      // Check folder requests
      for (const folder of col.folders) {
        const folderReq = folder.requests.find(r => r.id === id)
        if (folderReq) return folderReq
      }
    }
    return undefined
  }

  // Find the path to a request (collection and optional folder)
  function findRequestPath(requestId: number): { collectionId: number; folderId?: number } | null {
    for (const col of tree.value) {
      // Check direct requests
      if (col.requests.some(r => r.id === requestId)) {
        return { collectionId: col.collection.id }
      }

      // Check folder requests
      for (const folder of col.folders) {
        if (folder.requests.some(r => r.id === requestId)) {
          return { collectionId: col.collection.id, folderId: folder.folder.id }
        }
      }
    }
    return null
  }

  // Set the tree data
  function setTree(data: CollectionTree[]) {
    tree.value = data
  }

  // Load tree from API
  async function loadTree() {
    try {
      loading.value = true
      const api = await getApi()
      const data = await api.getCollectionTree()
      tree.value = data
    } catch (err) {
      console.error('Failed to load collection tree:', err)
    } finally {
      loading.value = false
    }
  }

  // Add a collection
  function addCollection(collection: Collection) {
    tree.value.push({
      collection,
      folders: [],
      requests: [],
    })
  }

  // Update a collection
  function updateCollection(collection: Collection) {
    const item = tree.value.find(t => t.collection.id === collection.id)
    if (item) {
      item.collection = collection
    }
  }

  // Delete a collection
  function deleteCollection(id: number) {
    const index = tree.value.findIndex(t => t.collection.id === id)
    if (index !== -1) {
      tree.value.splice(index, 1)
    }
  }

  // Add a folder to a collection
  function addFolder(collectionId: number, folder: Folder) {
    const col = tree.value.find(t => t.collection.id === collectionId)
    if (col) {
      col.folders.push({
        folder,
        requests: [],
      })
    }
  }

  // Update a folder
  function updateFolder(folder: Folder) {
    for (const col of tree.value) {
      const folderTree = col.folders.find(f => f.folder.id === folder.id)
      if (folderTree) {
        folderTree.folder = folder
        return
      }
    }
  }

  // Delete a folder
  function deleteFolder(id: number) {
    for (const col of tree.value) {
      const index = col.folders.findIndex(f => f.folder.id === id)
      if (index !== -1) {
        col.folders.splice(index, 1)
        return
      }
    }
  }

  // Add a request
  function addRequest(request: Request) {
    const col = tree.value.find(t => t.collection.id === request.collectionId)
    if (!col) return

    if (request.folderId) {
      const folder = col.folders.find(f => f.folder.id === request.folderId)
      if (folder) {
        folder.requests.push(request)
      }
    } else {
      col.requests.push(request)
    }
  }

  // Update a request
  function updateRequest(request: Request) {
    for (const col of tree.value) {
      // Check direct requests
      const directIndex = col.requests.findIndex(r => r.id === request.id)
      if (directIndex !== -1) {
        col.requests[directIndex] = request
        return
      }

      // Check folder requests
      for (const folder of col.folders) {
        const folderIndex = folder.requests.findIndex(r => r.id === request.id)
        if (folderIndex !== -1) {
          folder.requests[folderIndex] = request
          return
        }
      }
    }
  }

  // Delete a request
  function deleteRequest(id: number) {
    for (const col of tree.value) {
      // Check direct requests
      const directIndex = col.requests.findIndex(r => r.id === id)
      if (directIndex !== -1) {
        col.requests.splice(directIndex, 1)
        return
      }

      // Check folder requests
      for (const folder of col.folders) {
        const folderIndex = folder.requests.findIndex(r => r.id === id)
        if (folderIndex !== -1) {
          folder.requests.splice(folderIndex, 1)
          return
        }
      }
    }
  }

  return {
    tree,
    loading,
    error,
    collections,
    getCollection,
    getFolder,
    getRequest,
    findRequestPath,
    setTree,
    loadTree,
    addCollection,
    updateCollection,
    deleteCollection,
    addFolder,
    updateFolder,
    deleteFolder,
    addRequest,
    updateRequest,
    deleteRequest,
  }
})
