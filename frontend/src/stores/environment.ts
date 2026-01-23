import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Environment, Variable } from '@/types'

export const useEnvironmentStore = defineStore('environment', () => {
  const environments = ref<Environment[]>([])
  const globalVariables = ref<Variable[]>([])
  const activeEnvId = ref<number | null>(null)
  const loading = ref(false)

  // Get active environment
  const activeEnvironment = computed(() => {
    if (!activeEnvId.value) return null
    return environments.value.find(e => e.id === activeEnvId.value) || null
  })

  // Get all variables (global + active env)
  const allVariables = computed(() => {
    const vars = new Map<string, string>()
    
    // Global variables first
    for (const v of globalVariables.value) {
      vars.set(v.key, v.value)
    }
    
    // Environment variables override
    if (activeEnvironment.value) {
      for (const v of activeEnvironment.value.variables) {
        vars.set(v.key, v.value)
      }
    }
    
    return vars
  })

  // Resolve a variable
  function resolveVariable(name: string): string | undefined {
    return allVariables.value.get(name)
  }

  // Replace variables in text
  function replaceVariables(text: string): string {
    return text.replace(/\{\{(\w+)\}\}/g, (match, name) => {
      return resolveVariable(name) ?? match
    })
  }

  // Set environments
  function setEnvironments(envs: Environment[]) {
    environments.value = envs
  }

  // Set global variables
  function setGlobalVariables(vars: Variable[]) {
    globalVariables.value = vars
  }

  // Set active environment
  function setActiveEnv(id: number | null) {
    activeEnvId.value = id
  }

  // Add environment
  function addEnvironment(env: Environment) {
    environments.value.push(env)
  }

  // Update environment
  function updateEnvironment(env: Environment) {
    const index = environments.value.findIndex(e => e.id === env.id)
    if (index !== -1) {
      environments.value[index] = env
    }
  }

  // Delete environment
  function deleteEnvironment(id: number) {
    const index = environments.value.findIndex(e => e.id === id)
    if (index !== -1) {
      environments.value.splice(index, 1)
      if (activeEnvId.value === id) {
        activeEnvId.value = null
      }
    }
  }

  return {
    environments,
    globalVariables,
    activeEnvId,
    loading,
    activeEnvironment,
    allVariables,
    resolveVariable,
    replaceVariables,
    setEnvironments,
    setGlobalVariables,
    setActiveEnv,
    addEnvironment,
    updateEnvironment,
    deleteEnvironment,
  }
})
