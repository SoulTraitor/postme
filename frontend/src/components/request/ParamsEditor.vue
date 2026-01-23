<template>
  <KeyValueEditor
    :items="params"
    @update:items="$emit('update:params', $event)"
    keyPlaceholder="Parameter"
    valuePlaceholder="Value"
  />
</template>

<script setup lang="ts">
import KeyValueEditor from '@/components/common/KeyValueEditor.vue'
import type { KeyValue } from '@/types'

defineProps<{
  params: KeyValue[]
}>()

defineEmits<{
  'update:params': [value: KeyValue[]]
}>()
</script>
