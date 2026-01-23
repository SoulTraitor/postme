import { ref } from 'vue'

// Event emitter for keyboard actions
type ActionCallback = () => void

const listeners: Map<string, Set<ActionCallback>> = new Map()

// Subscribe to an action
export function onKeyboardAction(action: string, callback: ActionCallback) {
  if (!listeners.has(action)) {
    listeners.set(action, new Set())
  }
  listeners.get(action)!.add(callback)
  
  // Return unsubscribe function
  return () => {
    listeners.get(action)?.delete(callback)
  }
}

// Emit an action
export function emitKeyboardAction(action: string) {
  listeners.get(action)?.forEach(cb => cb())
}

// Composable for registering actions
export function useKeyboardActions() {
  return {
    emit: emitKeyboardAction,
    on: onKeyboardAction,
  }
}
