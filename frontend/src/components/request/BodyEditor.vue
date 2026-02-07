<template>
  <div class="flex flex-col h-full">
    <!-- Body type selector -->
    <div class="flex gap-2 mb-4">
      <button
        v-for="type in bodyTypes"
        :key="type.value"
        @click="$emit('update:bodyType', type.value)"
        class="px-3 py-1.5 rounded-md text-sm font-medium transition-colors"
        :class="[
          bodyType === type.value
            ? 'bg-accent text-white'
            : (effectiveTheme === 'dark'
              ? 'bg-dark-surface text-gray-300 hover:bg-dark-hover'
              : 'bg-light-surface text-gray-600 hover:bg-light-hover')
        ]"
      >
        {{ type.label }}
      </button>
    </div>
    
    <!-- Form data editor (multipart) -->
    <div v-if="bodyType === 'form-data'" class="flex-1 min-h-0 overflow-auto">
      <p 
        class="text-xs mb-2"
        :class="effectiveTheme === 'dark' ? 'text-gray-500' : 'text-gray-400'"
      >
        Multipart form data - supports file uploads
      </p>
      <KeyValueEditor
        :items="formDataItems"
        @update:items="updateFormData"
        keyPlaceholder="Field name"
        valuePlaceholder="Value"
        :showFileUpload="true"
      />
    </div>
    
    <!-- URL Encoded form editor -->
    <div v-else-if="bodyType === 'x-www-form-urlencoded'" class="flex-1 min-h-0 overflow-auto">
      <p 
        class="text-xs mb-2"
        :class="effectiveTheme === 'dark' ? 'text-gray-500' : 'text-gray-400'"
      >
        URL encoded form data - text values only
      </p>
      <KeyValueEditor
        :items="formDataItems"
        @update:items="updateFormData"
        keyPlaceholder="Field name"
        valuePlaceholder="Value"
      />
    </div>
    
    <!-- Binary file picker -->
    <div v-else-if="bodyType === 'binary'" class="flex-1 flex flex-col items-center justify-center gap-4">
      <div 
        class="p-8 rounded-lg border-2 border-dashed text-center cursor-pointer transition-colors"
        :class="[
          effectiveTheme === 'dark' 
            ? 'border-dark-border hover:border-accent/50' 
            : 'border-light-border hover:border-accent/50'
        ]"
        @click="selectBinaryFile"
      >
        <DocumentArrowUpIcon class="w-12 h-12 mx-auto mb-3 text-gray-400" />
        <p 
          class="text-sm"
          :class="effectiveTheme === 'dark' ? 'text-gray-300' : 'text-gray-700'"
        >
          {{ body ? 'Click to change file' : 'Click to select a file' }}
        </p>
      </div>
      
      <div v-if="body" class="flex items-center gap-2">
        <span 
          class="text-sm truncate max-w-md"
          :class="effectiveTheme === 'dark' ? 'text-gray-300' : 'text-gray-700'"
        >
          {{ body }}
        </span>
        <button
          @click="$emit('update:body', '')"
          class="p-1 rounded hover:bg-red-500/20 text-red-500"
          title="Remove file"
        >
          <XMarkIcon class="w-4 h-4" />
        </button>
      </div>
    </div>
    
    <!-- Code editor for JSON/XML/Text -->
    <div v-else-if="bodyType !== 'none'" class="flex-1 min-h-0">
      <div 
        ref="editorContainer" 
        class="h-full min-h-[200px] rounded-md border overflow-hidden"
        :class="effectiveTheme === 'dark' ? 'border-dark-border bg-[#282c34]' : 'border-light-border bg-white'"
      />
    </div>
    
    <div v-else class="flex-1 flex items-center justify-center text-gray-500">
      This request does not have a body
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useAppStateStore } from '@/stores/appState'
import { EditorView, lineNumbers, highlightActiveLineGutter, highlightSpecialChars, drawSelection, dropCursor, rectangularSelection, crosshairCursor, highlightActiveLine, keymap } from '@codemirror/view'
import { EditorState, Prec } from '@codemirror/state'
import { foldGutter, indentOnInput, bracketMatching, foldKeymap } from '@codemirror/language'
import { history, defaultKeymap, historyKeymap } from '@codemirror/commands'
import { highlightSelectionMatches, search, searchKeymap, openSearchPanel } from '@codemirror/search'
import { closeBrackets, autocompletion, closeBracketsKeymap, completionKeymap } from '@codemirror/autocomplete'
import { lintKeymap } from '@codemirror/lint'
import { json } from '@codemirror/lang-json'
import { xml } from '@codemirror/lang-xml'
import { oneDark } from '@codemirror/theme-one-dark'
import { syntaxHighlighting, HighlightStyle, defaultHighlightStyle } from '@codemirror/language'
import { tags } from '@lezer/highlight'
import { emitKeyboardAction } from '@/composables/useKeyboardActions'
import KeyValueEditor from '@/components/common/KeyValueEditor.vue'
import { DocumentArrowUpIcon, XMarkIcon } from '@heroicons/vue/24/outline'
import type { KeyValue } from '@/types'

const props = defineProps<{
  body: string
  bodyType: string
}>()

const emit = defineEmits<{
  'update:body': [value: string]
  'update:bodyType': [value: string]
}>()

const appState = useAppStateStore()
const effectiveTheme = computed(() => appState.effectiveTheme)

const bodyTypes = [
  { value: 'none', label: 'None' },
  { value: 'form-data', label: 'Form Data' },
  { value: 'x-www-form-urlencoded', label: 'URL Encoded' },
  { value: 'binary', label: 'Binary' },
  { value: 'json', label: 'JSON' },
  { value: 'xml', label: 'XML' },
  { value: 'text', label: 'Text' },
]

// Light mode highlight style for JSON/XML (brighter colors for better visibility)
const lightHighlightStyle = HighlightStyle.define([
  { tag: tags.string, color: '#f97316' },      // Orange for strings
  { tag: tags.number, color: '#0d6efd' },      // Blue for numbers
  { tag: tags.bool, color: '#eab308' },        // Yellow for booleans
  { tag: tags.null, color: '#eab308' },        // Yellow for null
  { tag: tags.propertyName, color: '#22c55e' }, // Green for property names
  { tag: tags.separator, color: '#64748b' },
  { tag: tags.squareBracket, color: '#64748b' },
  { tag: tags.brace, color: '#64748b' },
  // XML tags
  { tag: tags.tagName, color: '#22c55e' },     // Green for tag names
  { tag: tags.attributeName, color: '#a855f7' }, // Light purple for attributes
  { tag: tags.attributeValue, color: '#f97316' },
  { tag: tags.angleBracket, color: '#64748b' },
  { tag: tags.content, color: '#475569' },
  { tag: tags.comment, color: '#94a3b8', fontStyle: 'italic' },
])

const editorContainer = ref<HTMLElement | null>(null)
let editor: EditorView | null = null

// Form data handling - stored as JSON array in body (for both form-data and x-www-form-urlencoded)
const formDataItems = computed<KeyValue[]>(() => {
  if ((props.bodyType !== 'form-data' && props.bodyType !== 'x-www-form-urlencoded') || !props.body) {
    return []
  }
  try {
    const parsed = JSON.parse(props.body)
    if (Array.isArray(parsed)) {
      return parsed.map(item => ({
        key: item.key || '',
        value: item.value || '',
        enabled: item.enabled !== false,
      }))
    }
  } catch {
    // Invalid JSON, return empty
  }
  return []
})

function updateFormData(items: KeyValue[]) {
  emit('update:body', JSON.stringify(items))
}

async function selectBinaryFile() {
  try {
    const { OpenFileDialog } = await import('../../../wailsjs/go/handlers/DialogHandler')
    const filePath = await OpenFileDialog('Select File')
    if (filePath) {
      emit('update:body', filePath)
    }
  } catch (error) {
    console.error('Failed to open file dialog:', error)
  }
}

function formatBody() {
  if (!editor) return
  
  const content = editor.state.doc.toString()
  let formatted = content
  
  try {
    if (props.bodyType === 'json') {
      formatted = JSON.stringify(JSON.parse(content), null, 2)
    }
    // XML formatting could be added here if needed
  } catch {
    // If parsing fails, leave content as-is
    return
  }
  
  if (formatted !== content) {
    editor.dispatch({
      changes: { from: 0, to: editor.state.doc.length, insert: formatted }
    })
  }
}

function createEditor() {
  if (!editorContainer.value || props.bodyType === 'none' || props.bodyType === 'form-data' || props.bodyType === 'x-www-form-urlencoded' || props.bodyType === 'binary') return
  
  // Build extensions manually (no basicSetup to avoid defaultHighlightStyle)
  const extensions = [
    // Custom keymap first
    keymap.of([
      {
        key: 'Ctrl-Enter',
        preventDefault: true,
        run: () => {
          emitKeyboardAction('send')
          return true
        },
      },
      {
        key: 'Mod-Enter',
        preventDefault: true,
        run: () => {
          emitKeyboardAction('send')
          return true
        },
      },
      {
        key: 'Ctrl-h',
        run: (view: EditorView) => {
          openSearchPanel(view)
          return true
        },
      },
      {
        key: 'Ctrl-Shift-f',
        run: () => {
          formatBody()
          return true
        },
      },
      {
        key: 'Mod-Shift-f',
        run: () => {
          formatBody()
          return true
        },
      },
    ]),
    // Core functionality (from basicSetup but without defaultHighlightStyle)
    lineNumbers(),
    highlightActiveLineGutter(),
    highlightSpecialChars(),
    history(),
    foldGutter(),
    drawSelection(),
    dropCursor(),
    EditorState.allowMultipleSelections.of(true),
    indentOnInput(),
    bracketMatching(),
    closeBrackets(),
    autocompletion(),
    rectangularSelection(),
    crosshairCursor(),
    highlightActiveLine(),
    highlightSelectionMatches(),
    search({ top: true }),
    keymap.of([
      ...closeBracketsKeymap,
      ...defaultKeymap,
      ...searchKeymap,
      ...historyKeymap,
      ...foldKeymap,
      ...completionKeymap,
      ...lintKeymap,
    ]),
    EditorView.updateListener.of((update) => {
      if (update.docChanged) {
        emit('update:body', update.state.doc.toString())
      }
    }),
  ]
  
  // Add language support FIRST (before highlighting)
  if (props.bodyType === 'json') {
    extensions.push(json())
  } else if (props.bodyType === 'xml') {
    extensions.push(xml())
  }
  
  // Add theme based on current mode
  if (effectiveTheme.value === 'dark') {
    extensions.push(oneDark)
  } else {
    // Light mode: use custom GitHub-style highlight with fallback
    extensions.push(syntaxHighlighting(lightHighlightStyle, { fallback: true }))
    // Add light mode editor theme
    extensions.push(EditorView.theme({
      '&': { backgroundColor: '#ffffff' },
      '.cm-content': { caretColor: '#24292e' },
      '.cm-cursor': { borderLeftColor: '#24292e' },
      '.cm-gutters': { backgroundColor: '#f5f5f5', color: '#6e7781', borderRight: '1px solid #e5e5e5' },
      '.cm-activeLineGutter': { backgroundColor: '#e8e8e8' },
      '.cm-activeLine': { backgroundColor: '#f0f0f0' },
      '&.cm-focused > .cm-scroller > .cm-selectionLayer .cm-selectionBackground, .cm-selectionBackground, .cm-content ::selection': { backgroundColor: '#b4d7ff' },
    }, { dark: false }))
  }

  // Set editor height to fill container with line wrapping
  extensions.push(EditorView.lineWrapping)
  extensions.push(EditorView.theme({
    '&': { height: '100%' },
    '.cm-scroller': { overflow: 'auto' }
  }))

  editor = new EditorView({
    state: EditorState.create({
      doc: props.body,
      extensions,
    }),
    parent: editorContainer.value,
  })
}

function destroyEditor() {
  if (editor) {
    editor.destroy()
    editor = null
  }
}

watch([() => props.bodyType, effectiveTheme], () => {
  destroyEditor()
  if (props.bodyType !== 'none' && props.bodyType !== 'form-data' && props.bodyType !== 'x-www-form-urlencoded' && props.bodyType !== 'binary') {
    setTimeout(createEditor, 0)
  }
})

watch(() => props.body, (newBody) => {
  if (editor && editor.state.doc.toString() !== newBody) {
    editor.dispatch({
      changes: { from: 0, to: editor.state.doc.length, insert: newBody }
    })
  }
})

onMounted(() => {
  if (props.bodyType !== 'none' && props.bodyType !== 'form-data' && props.bodyType !== 'x-www-form-urlencoded' && props.bodyType !== 'binary') {
    createEditor()
  }
})

onUnmounted(() => {
  destroyEditor()
})
</script>
