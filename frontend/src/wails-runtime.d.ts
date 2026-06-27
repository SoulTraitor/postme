import type * as WailsRuntime from '../wailsjs/runtime/runtime'

declare global {
  interface Window {
    runtime?: typeof WailsRuntime
  }
}

export {}
