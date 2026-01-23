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
    
    <!-- Body editor -->
    <div v-if="bodyType !== 'none'" class="flex-1 min-h-0">
      <div 
        ref="editorContainer" 
        class="h-full rounded-md border overflow-hidden"
        :class="effectiveTheme === 'dark' ? 'border-dark-border' : 'border-light-border'"
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
import { EditorView, basicSetup } from 'codemirror'
import { EditorState } from '@codemirror/state'
import { keymap } from '@codemirror/view'
import { json } from '@codemirror/lang-json'
import { xml } from '@codemirror/lang-xml'
import { oneDark } from '@codemirror/theme-one-dark'
import { search, openSearchPanel } from '@codemirror/search'
import { emitKeyboardAction } from '@/composables/useKeyboardActions'

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
  { value: 'json', label: 'JSON' },
  { value: 'xml', label: 'XML' },
  { value: 'text', label: 'Text' },
]

const editorContainer = ref<HTMLElement | null>(null)
let editor: EditorView | null = null

function createEditor() {
  if (!editorContainer.value || props.bodyType === 'none') return
  
  const extensions = [
    basicSetup,
    search({ top: true }),
    // Custom keymap to handle Ctrl+Enter and other shortcuts
    keymap.of([
      {
        key: 'Ctrl-Enter',
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
        run: (view: EditorView) => {
          openSearchPanel(view)
          return true
        },
      },
    ]),
    EditorView.updateListener.of((update) => {
      if (update.docChanged) {
        emit('update:body', update.state.doc.toString())
      }
    }),
  ]
  
  // Add language support
  if (props.bodyType === 'json') {
    extensions.push(json())
  } else if (props.bodyType === 'xml') {
    extensions.push(xml())
  }
  
  // Add dark theme if needed
  if (effectiveTheme.value === 'dark') {
    extensions.push(oneDark)
  }
  
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
  setTimeout(createEditor, 0)
})

watch(() => props.body, (newBody) => {
  if (editor && editor.state.doc.toString() !== newBody) {
    editor.dispatch({
      changes: { from: 0, to: editor.state.doc.length, insert: newBody }
    })
  }
})

onMounted(() => {
  if (props.bodyType !== 'none') {
    createEditor()
  }
})

onUnmounted(() => {
  destroyEditor()
})
</script>
