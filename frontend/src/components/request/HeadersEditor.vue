<template>
  <KeyValueEditor
    :items="headers"
    @update:items="$emit('update:headers', $event)"
    keyPlaceholder="Header"
    valuePlaceholder="Value"
  />
</template>

<script setup lang="ts">
import KeyValueEditor from '@/components/common/KeyValueEditor.vue'
import type { KeyValue } from '@/types'

defineProps<{
  headers: KeyValue[]
}>()

defineEmits<{
  'update:headers': [value: KeyValue[]]
}>()
</script>
