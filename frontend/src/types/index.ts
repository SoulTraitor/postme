// TypeScript type definitions

// Key-Value pair
export interface KeyValue {
  key: string
  value: string
  enabled: boolean
  type?: 'text' | 'file'  // for form-data, file type means value is a file path
}

// Variable for environments
export interface Variable {
  key: string
  value: string
  secret: boolean
}

// HTTP Request
export interface Request {
  id: number
  collectionId: number
  folderId: number | null
  name: string
  method: string
  url: string
  headers: KeyValue[]
  params: KeyValue[]
  body: string
  bodyType: string
  sortOrder: number
  createdAt: string
  updatedAt: string
}

// HTTP Response
export interface Response {
  statusCode: number
  status: string
  headers: Record<string, string>
  body: string
  size: number
  duration: number
}

// Collection
export interface Collection {
  id: number
  name: string
  description: string
  sortOrder: number
  createdAt: string
  updatedAt: string
}

// Folder
export interface Folder {
  id: number
  collectionId: number
  name: string
  sortOrder: number
  createdAt: string
  updatedAt: string
}

// Environment
export interface Environment {
  id: number
  name: string
  variables: Variable[]
  createdAt: string
  updatedAt: string
}

// Global Variables
export interface GlobalVariables {
  id: number
  variables: Variable[]
  updatedAt: string
}

// History entry
export interface History {
  id: number
  requestId: number | null
  method: string
  url: string
  requestHeaders: string
  requestBody: string
  statusCode: number | null
  responseHeaders: string
  responseBody: string
  durationMs: number | null
  createdAt: string
}

// App State
export interface AppState {
  id: number
  windowWidth: number
  windowHeight: number
  windowX: number | null
  windowY: number | null
  windowMaximized: boolean
  sidebarOpen: boolean
  sidebarWidth: number
  layoutDirection: 'horizontal' | 'vertical'
  splitRatio: number
  theme: 'light' | 'dark' | 'system'
  activeEnvId: number | null
  requestTimeout: number
  autoLocateSidebar: boolean
  useSystemProxy: boolean
  requestPanelTab: 'params' | 'headers' | 'body'
  updatedAt: string
}

// Sidebar State
export interface SidebarState {
  id: number
  itemType: string
  itemId: number
  expanded: boolean
}

// Tab Session
export interface TabSession {
  id: number
  tabId: string
  requestId: number | null
  title: string
  sortOrder: number
  isActive: boolean
  isDirty: boolean
  method: string
  url: string
  headers: KeyValue[]
  params: KeyValue[]
  body: string
  bodyType: string
  createdAt: string
  updatedAt: string
}

// Collection Tree structures
export interface FolderTree {
  folder: Folder
  requests: Request[]
}

export interface CollectionTree {
  collection: Collection
  folders: FolderTree[]
  requests: Request[]
}

// Tab state for UI
export interface Tab {
  id: string
  requestId: number | null
  title: string
  method: string
  url: string
  headers: KeyValue[]
  params: KeyValue[]
  body: string
  bodyType: string
  isDirty: boolean
  isPreview: boolean
  // Original state for comparing dirty (null if new/unsaved tab)
  originalState?: {
    method: string
    url: string
    headers: KeyValue[]
    params: KeyValue[]
    body: string
    bodyType: string
  } | null
}

// Execute request params
export interface ExecuteRequestParams {
  tabId: string
  method: string
  url: string
  headers: KeyValue[]
  body: string
  bodyType: string
  timeout: number
}

// Response state
export type ResponseState = 
  | { status: 'idle' }
  | { status: 'loading' }
  | { status: 'cancelled' }
  | { status: 'timeout'; seconds: number }
  | { status: 'error'; message: string }
  | { status: 'success'; response: Response }

// Toast types
export type ToastType = 'success' | 'error' | 'warning' | 'info'

export interface Toast {
  id: string
  type: ToastType
  message: string
}

// Modal types
export type ModalType = 'confirm' | 'danger' | 'input' | 'select' | 'info'

// HTTP Methods
export const HTTP_METHODS = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'OPTIONS', 'HEAD'] as const
export type HttpMethod = typeof HTTP_METHODS[number]

// Body types
export const BODY_TYPES = ['none', 'json', 'xml', 'text', 'form-data', 'x-www-form-urlencoded'] as const
export type BodyType = typeof BODY_TYPES[number]
