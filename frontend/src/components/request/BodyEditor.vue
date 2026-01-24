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
  
  const extensions = [
    // Custom keymap FIRST to intercept before basicSetup
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
        run: (view: EditorView) => {
          formatBody()
          return true
        },
      },
      {
        key: 'Mod-Shift-f',
        run: (view: EditorView) => {
          formatBody()
          return true
        },
      },
    ]),
    basicSetup,
    search({ top: true }),
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
