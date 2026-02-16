# Real-time System Theme Following — Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Make the app's color mode follow system theme changes in real-time, not just at startup.

**Architecture:** Replace the non-reactive `matchMedia` call inside `effectiveTheme` computed with a reactive `systemIsDark` ref. Add Wails `WindowSetBackgroundColour` call to sync window background with theme.

**Tech Stack:** Vue 3 (Pinia store, `ref`/`computed`/`watch`), Wails v2 runtime API

---

## Problem

When theme is set to "system", the app only detects the system color scheme at startup. Changing the system theme while the app is running has no effect until restart.

**Root cause:** `effectiveTheme` is a Vue `computed` that calls `window.matchMedia(...).matches` internally. Vue's reactivity system cannot track browser API changes. The existing workaround (`theme.value = 'system'` on change event) fails because Vue detects no actual value change and skips re-evaluation.

---

### Task 1: Add reactive system preference tracking

**Files:**
- Modify: `frontend/src/stores/appState.ts:45-53`

**Step 1: Add `systemIsDark` ref after the `loading` ref (line 45)**

Add this line after `const loading = ref(true)`:

```typescript
  // Track system dark mode preference reactively
  const systemIsDark = ref(
    typeof window !== 'undefined'
      ? window.matchMedia('(prefers-color-scheme: dark)').matches
      : false
  )
```

**Step 2: Modify `effectiveTheme` computed to use `systemIsDark` ref**

Replace lines 48-53:

```typescript
  // Current effective theme
  const effectiveTheme = computed(() => {
    if (theme.value === 'system') {
      return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
    }
    return theme.value
  })
```

With:

```typescript
  // Current effective theme
  const effectiveTheme = computed(() => {
    if (theme.value === 'system') {
      return systemIsDark.value ? 'dark' : 'light'
    }
    return theme.value
  })
```

---

### Task 2: Replace broken listener with reactive one

**Files:**
- Modify: `frontend/src/stores/appState.ts:66-74`

**Step 1: Replace the old broken listener**

Replace lines 66-74:

```typescript
  // Listen to system theme changes
  if (typeof window !== 'undefined') {
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
      // Force reactivity update
      if (theme.value === 'system') {
        theme.value = 'system'
      }
    })
  }
```

With:

```typescript
  // Listen to system theme changes — update reactive ref
  if (typeof window !== 'undefined') {
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
      systemIsDark.value = e.matches
    })
  }
```

---

### Task 3: Sync Wails window background color with theme

**Files:**
- Modify: `frontend/src/stores/appState.ts:57-64`

**Step 1: Import Wails runtime function**

Add to top of file (after the existing imports on line 3):

```typescript
import { WindowSetBackgroundColour } from '../../wailsjs/runtime/runtime'
```

**Step 2: Update the `watch(effectiveTheme)` callback to also set window background**

Replace lines 57-64:

```typescript
  // Apply theme to document
  watch(effectiveTheme, (newTheme) => {
    if (newTheme === 'dark') {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }, { immediate: true })
```

With:

```typescript
  // Apply theme to document and window background
  watch(effectiveTheme, (newTheme) => {
    if (newTheme === 'dark') {
      document.documentElement.classList.add('dark')
      WindowSetBackgroundColour(26, 26, 26, 255)
    } else {
      document.documentElement.classList.remove('dark')
      WindowSetBackgroundColour(255, 255, 255, 255)
    }
  }, { immediate: true })
```

---

### Task 4: Build verification and commit

**Step 1: Run frontend build to verify no TypeScript errors**

Run: `cd frontend && npm run build`
Expected: Build succeeds with no errors

**Step 2: Run full Wails build to verify integration**

Run: `wails build`
Expected: Build succeeds

**Step 3: Manual test**

1. Launch the app
2. Set theme to "System" in settings
3. Change Windows system theme (Settings > Personalization > Colors > Choose your mode)
4. Verify: app theme switches in real-time without restart
5. Verify: window background color matches the theme (no dark flash in light mode)

**Step 4: Commit**

```bash
git add frontend/src/stores/appState.ts
git commit -m "fix: real-time system theme following with window background sync"
```
