// API service layer for Wails backend communication
import * as RequestHandler from '../../wailsjs/go/handlers/RequestHandler'
import * as CollectionHandler from '../../wailsjs/go/handlers/CollectionHandler'
import * as EnvironmentHandler from '../../wailsjs/go/handlers/EnvironmentHandler'
import * as HistoryHandler from '../../wailsjs/go/handlers/HistoryHandler'
import * as AppStateHandler from '../../wailsjs/go/handlers/AppStateHandler'
import { models, handlers, services } from '../../wailsjs/go/models'
import type { 
  CollectionTree, 
  FolderTree, 
  Request, 
  Collection, 
  Folder, 
  Environment, 
  History,
  AppState,
  SidebarState,
  TabSession,
  KeyValue,
  Variable,
  Response as ResponseType
} from '@/types'

// Type converters - convert Wails generated types to our frontend types

function convertKeyValue(kv: models.KeyValue): KeyValue {
  return {
    key: kv.key,
    value: kv.value,
    enabled: kv.enabled,
  }
}

function convertVariable(v: models.Variable): Variable {
  return {
    key: v.key,
    value: v.value,
    secret: v.secret,
  }
}

function convertRequest(req: models.Request): Request {
  return {
    id: req.id,
    collectionId: req.collectionId,
    folderId: req.folderId ?? null,
    name: req.name,
    method: req.method,
    url: req.url,
    headers: (req.headers || []).map(convertKeyValue),
    params: (req.params || []).map(convertKeyValue),
    body: req.body,
    bodyType: req.bodyType,
    sortOrder: req.sortOrder,
    createdAt: String(req.createdAt),
    updatedAt: String(req.updatedAt),
  }
}

function convertCollection(col: models.Collection): Collection {
  return {
    id: col.id,
    name: col.name,
    description: col.description,
    sortOrder: col.sortOrder,
    createdAt: String(col.createdAt),
    updatedAt: String(col.updatedAt),
  }
}

function convertFolder(folder: models.Folder): Folder {
  return {
    id: folder.id,
    collectionId: folder.collectionId,
    name: folder.name,
    sortOrder: folder.sortOrder,
    createdAt: String(folder.createdAt),
    updatedAt: String(folder.updatedAt),
  }
}

function convertEnvironment(env: models.Environment): Environment {
  return {
    id: env.id,
    name: env.name,
    variables: (env.variables || []).map(convertVariable),
    createdAt: String(env.createdAt),
    updatedAt: String(env.updatedAt),
  }
}

function convertHistory(h: models.History): History {
  return {
    id: h.id,
    requestId: h.requestId ?? null,
    method: h.method,
    url: h.url,
    requestHeaders: h.requestHeaders,
    requestBody: h.requestBody,
    statusCode: h.statusCode ?? null,
    responseHeaders: h.responseHeaders,
    responseBody: h.responseBody,
    durationMs: h.durationMs ?? null,
    createdAt: String(h.createdAt),
  }
}

function convertFolderTree(ft: services.FolderTree): FolderTree {
  return {
    folder: convertFolder(ft.folder),
    requests: (ft.requests || []).map(convertRequest),
  }
}

function convertCollectionTree(ct: services.CollectionTree): CollectionTree {
  return {
    collection: convertCollection(ct.collection),
    folders: (ct.folders || []).map(convertFolderTree),
    requests: (ct.requests || []).map(convertRequest),
  }
}

function convertAppState(state: models.AppState): AppState {
  return {
    id: state.id,
    windowWidth: state.windowWidth,
    windowHeight: state.windowHeight,
    windowX: state.windowX ?? null,
    windowY: state.windowY ?? null,
    windowMaximized: state.windowMaximized,
    sidebarOpen: state.sidebarOpen,
    sidebarWidth: state.sidebarWidth,
    layoutDirection: state.layoutDirection as 'horizontal' | 'vertical',
    splitRatio: state.splitRatio,
    theme: state.theme as 'light' | 'dark' | 'system',
    activeEnvId: state.activeEnvId ?? null,
    requestTimeout: state.requestTimeout,
    autoLocateSidebar: state.autoLocateSidebar,
    useSystemProxy: state.useSystemProxy,
    requestPanelTab: (state.requestPanelTab || 'params') as 'params' | 'headers' | 'body',
    updatedAt: String(state.updatedAt),
  }
}

function convertSidebarState(ss: models.SidebarState): SidebarState {
  return {
    id: ss.id,
    itemType: ss.itemType,
    itemId: ss.itemId,
    expanded: ss.expanded,
  }
}

function convertTabSession(ts: models.TabSession): TabSession {
  return {
    id: ts.id,
    tabId: ts.tabId,
    requestId: ts.requestId ?? null,
    title: ts.title,
    sortOrder: ts.sortOrder,
    isActive: ts.isActive,
    isDirty: ts.isDirty,
    method: ts.method,
    url: ts.url,
    headers: (ts.headers || []).map(convertKeyValue),
    params: (ts.params || []).map(convertKeyValue),
    body: ts.body,
    bodyType: ts.bodyType,
    createdAt: String(ts.createdAt),
    updatedAt: String(ts.updatedAt),
  }
}

function convertResponse(res: models.Response): ResponseType {
  return {
    statusCode: res.statusCode,
    status: res.status,
    headers: res.headers || {},
    body: res.body,
    size: res.size,
    duration: res.duration,
  }
}

// API functions

// Collections
export const api = {
  // Collection operations
  async getCollectionTree(): Promise<CollectionTree[]> {
    const tree = await CollectionHandler.GetTree()
    return (tree || []).map(convertCollectionTree)
  },

  async createCollection(name: string, description: string = ''): Promise<Collection> {
    const col = models.Collection.createFrom({
      name,
      description,
      sortOrder: 0,
    })
    const result = await CollectionHandler.Create(col)
    return convertCollection(result)
  },

  async updateCollection(collection: Collection): Promise<void> {
    const col = models.Collection.createFrom(collection)
    await CollectionHandler.Update(col)
  },

  async deleteCollection(id: number): Promise<void> {
    await CollectionHandler.Delete(id)
  },

  // Folder operations
  async createFolder(collectionId: number, name: string): Promise<Folder> {
    const folder = models.Folder.createFrom({
      collectionId,
      name,
      sortOrder: 0,
    })
    const result = await CollectionHandler.CreateFolder(folder)
    return convertFolder(result)
  },

  async updateFolder(folder: Folder): Promise<void> {
    const f = models.Folder.createFrom(folder)
    await CollectionHandler.UpdateFolder(f)
  },

  async deleteFolder(id: number): Promise<void> {
    await CollectionHandler.DeleteFolder(id)
  },

  async moveRequest(requestId: number, collectionId: number, folderId: number | null): Promise<void> {
    await CollectionHandler.MoveRequest(requestId, collectionId, folderId)
  },

  async moveFolder(folderId: number, collectionId: number): Promise<void> {
    await CollectionHandler.MoveFolder(folderId, collectionId)
  },

  async reorderCollections(ids: number[]): Promise<void> {
    await CollectionHandler.ReorderCollections(ids)
  },

  async reorderFolders(collectionId: number, ids: number[]): Promise<void> {
    await CollectionHandler.ReorderFolders(collectionId, ids)
  },

  async reorderRequests(collectionId: number, folderId: number | null, ids: number[]): Promise<void> {
    await CollectionHandler.ReorderRequests(collectionId, folderId, ids)
  },

  async exportCollection(id: number): Promise<void> {
    await CollectionHandler.ExportCollection(id)
  },

  async importCollection(): Promise<Collection | null> {
    const result = await CollectionHandler.ImportCollection()
    if (!result) return null
    return convertCollection(result)
  },

  // Request operations
  async createRequest(request: Partial<Request>): Promise<Request> {
    const req = models.Request.createFrom({
      collectionId: request.collectionId || 0,
      folderId: request.folderId,
      name: request.name || 'Untitled',
      method: request.method || 'GET',
      url: request.url || '',
      headers: (request.headers || []).map(h => models.KeyValue.createFrom(h)),
      params: (request.params || []).map(p => models.KeyValue.createFrom(p)),
      body: request.body || '',
      bodyType: request.bodyType || 'none',
      sortOrder: request.sortOrder || 0,
    })
    const result = await RequestHandler.Create(req)
    return convertRequest(result)
  },

  async updateRequest(request: Request): Promise<void> {
    const req = models.Request.createFrom({
      id: request.id,
      collectionId: request.collectionId,
      folderId: request.folderId,
      name: request.name,
      method: request.method,
      url: request.url,
      headers: request.headers.map(h => models.KeyValue.createFrom(h)),
      params: request.params.map(p => models.KeyValue.createFrom(p)),
      body: request.body,
      bodyType: request.bodyType,
      sortOrder: request.sortOrder,
    })
    await RequestHandler.Update(req)
  },

  async deleteRequest(id: number): Promise<void> {
    await RequestHandler.Delete(id)
  },

  async duplicateRequest(id: number): Promise<Request> {
    const req = await RequestHandler.Duplicate(id)
    return convertRequest(req)
  },

  async getRequest(id: number): Promise<Request> {
    const req = await RequestHandler.GetByID(id)
    return convertRequest(req)
  },

  // Execute request
  async executeRequest(params: {
    tabId: string
    method: string
    url: string
    headers: KeyValue[]
    body: string
    bodyType: string
    timeout: number
  }): Promise<ResponseType> {
    const execParams = handlers.ExecuteRequestParams.createFrom({
      tabId: params.tabId,
      method: params.method,
      url: params.url,
      headers: params.headers.map(h => models.KeyValue.createFrom(h)),
      body: params.body,
      bodyType: params.bodyType,
      timeout: params.timeout,
    })
    const result = await RequestHandler.Execute(execParams)
    return convertResponse(result)
  },

  async cancelRequest(tabId: string): Promise<void> {
    await RequestHandler.CancelRequest(tabId)
  },

  async setUseSystemProxy(useProxy: boolean): Promise<void> {
    await RequestHandler.SetUseSystemProxy(useProxy)
  },

  // Environment operations
  async getEnvironments(): Promise<Environment[]> {
    const envs = await EnvironmentHandler.GetAll()
    return (envs || []).map(convertEnvironment)
  },

  async createEnvironment(name: string, variables: Variable[] = []): Promise<Environment> {
    const env = models.Environment.createFrom({
      name,
      variables: variables.map(v => models.Variable.createFrom(v)),
    })
    const result = await EnvironmentHandler.Create(env)
    return convertEnvironment(result)
  },

  async updateEnvironment(env: Environment): Promise<void> {
    const e = models.Environment.createFrom({
      id: env.id,
      name: env.name,
      variables: env.variables.map(v => models.Variable.createFrom(v)),
    })
    await EnvironmentHandler.Update(e)
  },

  async deleteEnvironment(id: number): Promise<void> {
    await EnvironmentHandler.Delete(id)
  },

  async getGlobalVariables(): Promise<Variable[]> {
    const gv = await EnvironmentHandler.GetGlobalVariables()
    return (gv.variables || []).map(convertVariable)
  },

  async updateGlobalVariables(variables: Variable[]): Promise<void> {
    const vars = variables.map(v => models.Variable.createFrom(v))
    await EnvironmentHandler.UpdateGlobalVariables(vars)
  },

  // History operations
  async getHistory(): Promise<History[]> {
    const history = await HistoryHandler.GetAll()
    return (history || []).map(convertHistory)
  },

  async deleteHistoryItem(id: number): Promise<void> {
    await HistoryHandler.Delete(id)
  },

  async clearHistory(): Promise<void> {
    await HistoryHandler.Clear()
  },

  async addHistory(history: Partial<History>): Promise<History> {
    const h = models.History.createFrom({
      requestId: history.requestId,
      method: history.method || 'GET',
      url: history.url || '',
      requestHeaders: history.requestHeaders || '',
      requestBody: history.requestBody || '',
      statusCode: history.statusCode,
      responseHeaders: history.responseHeaders || '',
      responseBody: history.responseBody || '',
      durationMs: history.durationMs,
    })
    const result = await HistoryHandler.Create(h)
    return convertHistory(result)
  },

  // App state operations
  async getAppState(): Promise<AppState> {
    const state = await AppStateHandler.Get()
    return convertAppState(state)
  },

  async updateAppState(state: Partial<AppState>): Promise<void> {
    const currentState = await AppStateHandler.Get()
    const updated = models.AppState.createFrom({
      ...currentState,
      ...state,
    })
    await AppStateHandler.Update(updated)
  },

  async getSidebarState(): Promise<SidebarState[]> {
    const states = await AppStateHandler.GetSidebarState()
    return (states || []).map(convertSidebarState)
  },

  async setSidebarItemExpanded(itemType: string, itemId: number, expanded: boolean): Promise<void> {
    await AppStateHandler.SetSidebarItemExpanded(itemType, itemId, expanded)
  },

  async getTabSessions(): Promise<TabSession[]> {
    const sessions = await AppStateHandler.GetTabSessions()
    return (sessions || []).map(convertTabSession)
  },

  async saveTabSession(session: Partial<TabSession>): Promise<void> {
    const ts = models.TabSession.createFrom({
      tabId: session.tabId || '',
      requestId: session.requestId,
      title: session.title || 'Untitled',
      sortOrder: session.sortOrder || 0,
      isActive: session.isActive || false,
      isDirty: session.isDirty || false,
      method: session.method || 'GET',
      url: session.url || '',
      headers: (session.headers || []).map(h => models.KeyValue.createFrom(h)),
      params: (session.params || []).map(p => models.KeyValue.createFrom(p)),
      body: session.body || '',
      bodyType: session.bodyType || 'none',
    })
    await AppStateHandler.SaveTabSession(ts)
  },

  async deleteTabSession(tabId: string): Promise<void> {
    await AppStateHandler.DeleteTabSession(tabId)
  },

  async setActiveTab(tabId: string): Promise<void> {
    await AppStateHandler.SetActiveTab(tabId)
  },

  async clearTabSessions(): Promise<void> {
    await AppStateHandler.ClearTabSessions()
  },
}

export default api
