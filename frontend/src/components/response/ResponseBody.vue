<template>
  <div class="h-full flex flex-col">
    <!-- Actions bar -->
    <div class="flex items-center gap-2 px-4 py-2 border-b" :class="effectiveTheme === 'dark' ? 'border-dark-border' : 'border-light-border'">
      <!-- Format button -->
      <button
        v-if="isJSON"
        @click="formatBody"
        class="px-2 py-1 text-xs rounded transition-colors"
        :class="effectiveTheme === 'dark' ? 'hover:bg-dark-hover text-gray-400' : 'hover:bg-light-hover text-gray-500'"
      >
        Format
      </button>
      
      <!-- Copy button -->
      <button
        @click="copyBody"
        class="px-2 py-1 text-xs rounded transition-colors"
        :class="effectiveTheme === 'dark' ? 'hover:bg-dark-hover text-gray-400' : 'hover:bg-light-hover text-gray-500'"
      >
        {{ copied ? 'Copied!' : 'Copy' }}
      </button>
      
      <span class="flex-1" />
      
      <!-- Content type -->
      <span class="text-xs text-gray-500">
        {{ contentType || 'Unknown' }}
      </span>
    </div>
    
    <!-- Body content -->
    <div 
      class="flex-1 overflow-auto"
      :class="effectiveTheme === 'dark' ? 'bg-[#282c34]' : 'bg-white'"
    >
      <div ref="editorContainer" class="h-full" />
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
import { highlightSelectionMatches, searchKeymap } from '@codemirror/search'
import { closeBrackets, autocompletion, closeBracketsKeymap, completionKeymap } from '@codemirror/autocomplete'
import { json } from '@codemirror/lang-json'
import { xml } from '@codemirror/lang-xml'
import { html } from '@codemirror/lang-html'
import { oneDark } from '@codemirror/theme-one-dark'
import { syntaxHighlighting, HighlightStyle, defaultHighlightStyle } from '@codemirror/language'
import { tags } from '@lezer/highlight'

const props = defineProps<{
  body: string
  contentType: string
}>()

const appState = useAppStateStore()
const effectiveTheme = computed(() => appState.effectiveTheme)

const editorContainer = ref<HTMLElement | null>(null)
const copied = ref(false)
let editor: EditorView | null = null

const isJSON = computed(() => 
  props.contentType.includes('json') || 
  (props.body.trim().startsWith('{') || props.body.trim().startsWith('['))
)

const isXML = computed(() => 
  props.contentType.includes('xml')
)

const isHTML = computed(() => 
  props.contentType.includes('html')
)

// Light mode highlight style for JSON/XML/HTML (brighter colors for better visibility)
const lightHighlightStyle = HighlightStyle.define([
  { tag: tags.string, color: '#f97316' },      // Orange for strings
  { tag: tags.number, color: '#0d6efd' },      // Blue for numbers
  { tag: tags.bool, color: '#eab308' },        // Yellow for booleans
  { tag: tags.null, color: '#eab308' },        // Yellow for null
  { tag: tags.propertyName, color: '#22c55e' }, // Green for property names
  { tag: tags.separator, color: '#64748b' },
  { tag: tags.squareBracket, color: '#64748b' },
  { tag: tags.brace, color: '#64748b' },
  // XML/HTML tags
  { tag: tags.tagName, color: '#22c55e' },     // Green for tag names
  { tag: tags.attributeName, color: '#a855f7' }, // Light purple for attributes
  { tag: tags.attributeValue, color: '#f97316' },
  { tag: tags.angleBracket, color: '#64748b' },
  { tag: tags.content, color: '#475569' },
  { tag: tags.comment, color: '#94a3b8', fontStyle: 'italic' },
])

const lightEditorTheme = EditorView.theme({
  '&': { backgroundColor: '#ffffff' },
  '.cm-content': { caretColor: '#24292e' },
  '.cm-cursor': { borderLeftColor: '#24292e' },
  '.cm-gutters': { backgroundColor: '#f5f5f5', color: '#6e7781', borderRight: '1px solid #e5e5e5' },
  '.cm-activeLineGutter': { backgroundColor: '#e8e8e8' },
  '.cm-activeLine': { backgroundColor: '#f0f0f0' },
  '&.cm-focused > .cm-scroller > .cm-selectionLayer .cm-selectionBackground, .cm-selectionBackground, .cm-content ::selection': { backgroundColor: '#b4d7ff' },
}, { dark: false })

function createEditor() {
  if (!editorContainer.value) return
  
  // Build extensions manually (no basicSetup to avoid defaultHighlightStyle conflict)
  const extensions = [
    // Core functionality
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
    EditorState.readOnly.of(true),
    keymap.of([
      ...closeBracketsKeymap,
      ...defaultKeymap,
      ...searchKeymap,
      ...historyKeymap,
      ...foldKeymap,
      ...completionKeymap,
    ]),
  ]
  
  // Add language support based on content type
  if (isJSON.value) {
    extensions.push(json())
  } else if (isXML.value) {
    extensions.push(xml())
  } else if (isHTML.value) {
    extensions.push(html())
  }
  
  // Add theme based on current mode
  if (effectiveTheme.value === 'dark') {
    extensions.push(oneDark)
  } else {
    // Light mode: use custom GitHub-style highlight with fallback
    extensions.push(syntaxHighlighting(lightHighlightStyle, { fallback: true }))
    extensions.push(lightEditorTheme)
  }
  
  let content = props.body
  
  // Try to format JSON
  if (isJSON.value) {
    try {
      content = JSON.stringify(JSON.parse(props.body), null, 2)
    } catch {}
  }
  
  editor = new EditorView({
    state: EditorState.create({
      doc: content,
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

function formatBody() {
  if (!editor || !isJSON.value) return
  
  try {
    const formatted = JSON.stringify(JSON.parse(props.body), null, 2)
    editor.dispatch({
      changes: { from: 0, to: editor.state.doc.length, insert: formatted }
    })
  } catch {}
}

async function copyBody() {
  try {
    await navigator.clipboard.writeText(props.body)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch {}
}

watch([() => props.body, effectiveTheme], () => {
  destroyEditor()
  setTimeout(createEditor, 0)
})

onMounted(createEditor)
onUnmounted(destroyEditor)
</script>
