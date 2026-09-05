import { nextTick, onMounted, onUnmounted, type Ref } from 'vue'

export function useMacTapClick(button: Ref<HTMLButtonElement | null>, activate: () => void | Promise<void>) {
  const isMac = /macintosh|mac os x/i.test(navigator.userAgent)
  let pressedPointerId: number | null = null
  let release: { button: HTMLButtonElement; pointerId: number; x: number; y: number; time: number } | null = null
  let suppressPointerClick = false
  let gestureVersion = 0

  function isPrimaryMouse(event: PointerEvent) {
    return event.pointerType === 'mouse' && event.isPrimary && event.button === 0 && !event.ctrlKey
  }

  function handlePointerUp(event: PointerEvent) {
    if (!isPrimaryMouse(event)) return

    const hadPress = pressedPointerId === event.pointerId
    pressedPointerId = null
    release = null
    const target = button.value
    if (hadPress || !target || target.disabled || !(event.target instanceof Node) || !target.contains(event.target)) return

    release = {
      button: target,
      pointerId: event.pointerId,
      x: event.clientX,
      y: event.clientY,
      time: performance.now(),
    }
  }

  function handlePointerDown(event: PointerEvent) {
    const pendingRelease = release
    const version = ++gestureVersion
    release = null
    suppressPointerClick = false
    pressedPointerId = isPrimaryMouse(event) ? event.pointerId : null

    const target = button.value
    if (!isPrimaryMouse(event) || !pendingRelease || !target || target.disabled) return
    if (pendingRelease.button !== target || pendingRelease.pointerId !== event.pointerId) return
    if (!(event.target instanceof Node) || !target.contains(event.target)) return
    if (performance.now() - pendingRelease.time > 100) return
    if (Math.hypot(event.clientX - pendingRelease.x, event.clientY - pendingRelease.y) > 4) return

    // Some macOS IMEs deliver tap-to-click pointerup before pointerdown, with no click.
    // https://bugs.webkit.org/show_bug.cgi?id=219670
    pressedPointerId = null
    suppressPointerClick = true
    event.preventDefault()
    target.focus({ preventScroll: true })

    // Commit blur-driven parameter/form edits before activating the completed tap.
    void nextTick(() => {
      if (gestureVersion === version && button.value === target && target.isConnected && !target.disabled) {
        target.click()
      }
    })
  }

  function resetGesture() {
    gestureVersion++
    pressedPointerId = null
    release = null
    suppressPointerClick = false
  }

  onMounted(() => {
    if (!isMac) return
    // Observe presses outside the button too, so dragging onto it is not treated as a tap.
    document.addEventListener('pointerdown', handlePointerDown, true)
    document.addEventListener('pointerup', handlePointerUp, true)
    document.addEventListener('pointercancel', resetGesture, true)
    window.addEventListener('blur', resetGesture)
  })

  onUnmounted(() => {
    resetGesture()
    document.removeEventListener('pointerdown', handlePointerDown, true)
    document.removeEventListener('pointerup', handlePointerUp, true)
    document.removeEventListener('pointercancel', resetGesture, true)
    window.removeEventListener('blur', resetGesture)
  })

  return (event: MouseEvent) => {
    if (event.button !== 0 || (isMac && event.ctrlKey)) return
    if (event.detail > 0 && suppressPointerClick) return
    release = null
    return activate()
  }
}
