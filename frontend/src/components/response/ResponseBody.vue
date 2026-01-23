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
    <div class="flex-1 overflow-auto">
      <div ref="editorContainer" class="h-full" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useAppStateStore } from '@/stores/appState'
import { EditorView, basicSetup } from 'codemirror'
import { EditorState } from '@codemirror/state'
import { json } from '@codemirror/lang-json'
import { xml } from '@codemirror/lang-xml'
import { html } from '@codemirror/lang-html'
import { oneDark } from '@codemirror/theme-one-dark'

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

function createEditor() {
  if (!editorContainer.value) return
  
  const extensions = [
    basicSetup,
    EditorState.readOnly.of(true),
  ]
  
  // Add language support based on content type
  if (isJSON.value) {
    extensions.push(json())
  } else if (isXML.value) {
    extensions.push(xml())
  } else if (isHTML.value) {
    extensions.push(html())
  }
  
  // Add dark theme if needed
  if (effectiveTheme.value === 'dark') {
    extensions.push(oneDark)
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
