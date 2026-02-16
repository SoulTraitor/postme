# Real-time System Theme Following

## Problem

When theme is set to "system", the app only detects the system color scheme at startup. Changing the system theme (e.g., Windows dark/light mode) while the app is running has no effect until restart.

**Root cause:** `effectiveTheme` is a Vue `computed` that calls `window.matchMedia(...).matches` internally. Vue's reactivity system cannot track browser API changes. The existing workaround (`theme.value = 'system'` on change event) fails because Vue detects no actual value change and skips re-evaluation.

## Solution: Reactive System Preference Ref (Approach A)

### Changes (single file: `frontend/src/stores/appState.ts`)

1. **Add `systemIsDark` ref** initialized from `matchMedia('(prefers-color-scheme: dark)').matches`

2. **Register `matchMedia` change listener** that updates `systemIsDark.value = e.matches` on system theme change

3. **Modify `effectiveTheme` computed** to read `systemIsDark.value` instead of calling `matchMedia` directly

4. **Update `watch(effectiveTheme)` callback** to also call `window.runtime.WindowSetBackgroundColour()`:
   - dark: `(26, 26, 26, 255)` (#1a1a1a)
   - light: `(255, 255, 255, 255)` (#ffffff)

5. **Remove** the old ineffective "force reactivity" hack (lines 66-74)

### Data Flow

```
System theme change
  -> matchMedia 'change' event
  -> systemIsDark.value updated
  -> effectiveTheme recomputed (Vue tracks ref dependency)
  -> watch triggers:
     1. Toggle .dark class on <html>
     2. WindowSetBackgroundColour() updates Wails window background
  -> All dark:/light: Tailwind classes react via CSS
```

### Scope

- Only `frontend/src/stores/appState.ts` is modified
- No backend changes needed
- No new dependencies
