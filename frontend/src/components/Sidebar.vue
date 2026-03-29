<script lang="ts" setup>
import { ref } from 'vue'

const props = defineProps<{
  modelValue: boolean
  activeModule: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'update:activeModule', value: string): void
}>()

const collapsed = ref(false)

const modules = [
  { id: 'todos', name: '待办事项', icon: 'mdi-checkbox-marked-outline' },
  { id: 'codeTemplates', name: '代码模板', icon: 'mdi-code-tags' }
]

const selectModule = (id: string) => {
  emit('update:activeModule', id)
}

const toggleCollapse = () => {
  collapsed.value = !collapsed.value
}
</script>

<template>
  <v-navigation-drawer
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    :width="collapsed ? 60 : 200"
    class="sidebar"
    :permanent="true"
    color="background-secondary"
  >
    <v-list nav density="compact" class="py-2">
      <v-list-item
        v-for="module in modules"
        :key="module.id"
        :active="activeModule === module.id"
        :prepend-icon="module.icon"
        :title="collapsed ? '' : module.name"
        @click="selectModule(module.id)"
        class="module-item my-1"
        color="primary"
        rounded="lg"
        link
      />
    </v-list>
    
    <template #append>
      <v-divider />
      <v-btn
        block
        variant="text"
        :prepend-icon="collapsed ? 'mdi-chevron-right' : 'mdi-chevron-left'"
        @click="toggleCollapse"
        class="collapse-btn"
        color="secondary"
        size="small"
      >
        {{ collapsed ? '' : '收起' }}
      </v-btn>
    </template>
  </v-navigation-drawer>
</template>

<style scoped>
.sidebar {
  border-right: 1px solid rgb(var(--v-theme-border));
}

.module-item {
  margin: 4px 8px;
}
</style>
