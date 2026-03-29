<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { useTodosStore } from '../stores/todos'

const todosStore = useTodosStore()

const showAddDialog = ref(false)
const showEditDialog = ref(false)
const newTodoName = ref('')
const newTodoDescription = ref('')
const newTodoDeadline = ref('')
const editingTodo = ref<any>(null)

onMounted(() => {
  todosStore.loadTodos()
})

const openAddDialog = () => {
  newTodoName.value = ''
  newTodoDescription.value = ''
  newTodoDeadline.value = ''
  showAddDialog.value = true
}

const addTodo = async () => {
  if (newTodoName.value) {
    await todosStore.addTodo({
      name: newTodoName.value,
      description: newTodoDescription.value,
      deadline: newTodoDeadline.value
    })
    showAddDialog.value = false
  }
}

const openEditDialog = (todo: any) => {
  editingTodo.value = { ...todo }
  showEditDialog.value = true
}

const saveEdit = async () => {
  if (editingTodo.value) {
    await todosStore.updateTodo(editingTodo.value.id, {
      name: editingTodo.value.name,
      description: editingTodo.value.description,
      deadline: editingTodo.value.deadline
    })
    showEditDialog.value = false
    editingTodo.value = null
  }
}

const completeTodo = async (id: string) => {
  await todosStore.deleteTodo(id)
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN')
}
</script>

<template>
  <div class="todo-view pa-4">
    <!-- 标题和添加按钮 -->
    <div class="d-flex align-center justify-space-between mb-4">
      <div class="text-h5">待办事项</div>
      <v-btn
        color="primary"
        prepend-icon="mdi-plus"
        @click="openAddDialog"
      >
        添加待办
      </v-btn>
    </div>

    <!-- 待办卡片列表 -->
    <v-row v-if="todosStore.todos.length > 0">
      <v-col
        v-for="todo in todosStore.todos"
        :key="todo.id"
        cols="12"
        sm="6"
        md="4"
        lg="3"
      >
        <v-card class="todo-card" :class="{ 'completed': todo.completed }">
          <v-card-title class="text-h6 pb-2">
            {{ todo.name }}
          </v-card-title>
          
          <v-card-text class="pb-2">
            <div v-if="todo.description" class="text-body-2 mb-2">
              {{ todo.description }}
            </div>
            <div v-if="todo.deadline" class="text-caption text-secondary">
              <v-icon size="14" class="mr-1">mdi-clock-outline</v-icon>
              截止：{{ formatDate(todo.deadline) }}
            </div>
          </v-card-text>

          <v-card-actions>
            <v-spacer />
            <v-btn
              color="primary"
              variant="text"
              size="small"
              @click="openEditDialog(todo)"
            >
              <v-icon start>mdi-pencil</v-icon>
              编辑
            </v-btn>
            <v-btn
              color="success"
              variant="flat"
              size="small"
              @click="completeTodo(todo.id)"
            >
              <v-icon start>mdi-check</v-icon>
              完成
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>

    <!-- 空状态 -->
    <v-card v-else flat class="text-center pa-8">
      <v-icon size="64" color="secondary">mdi-clipboard-check-outline</v-icon>
      <div class="text-h6 mt-4 mb-2">暂无待办</div>
      <div class="text-body-2 text-secondary mb-4">
        点击上方"添加待办"创建新任务
      </div>
    </v-card>

    <!-- 添加对话框 -->
    <v-dialog v-model="showAddDialog" max-width="500">
      <v-card>
        <v-card-title class="text-h6">
          添加待办
        </v-card-title>
        <v-card-text>
          <v-text-field
            v-model="newTodoName"
            label="待办名称"
            variant="outlined"
            density="compact"
            class="mb-3"
            hide-details
            required
          />
          <v-textarea
            v-model="newTodoDescription"
            label="描述（可选）"
            variant="outlined"
            density="compact"
            rows="3"
            class="mb-3"
            hide-details
          />
          <v-text-field
            v-model="newTodoDeadline"
            label="截止日期"
            variant="outlined"
            density="compact"
            type="date"
            hide-details
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showAddDialog = false">
            取消
          </v-btn>
          <v-btn
            color="primary"
            variant="flat"
            @click="addTodo"
            :disabled="!newTodoName"
          >
            添加
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- 编辑对话框 -->
    <v-dialog v-model="showEditDialog" max-width="500">
      <v-card>
        <v-card-title class="text-h6">
          编辑待办
        </v-card-title>
        <v-card-text v-if="editingTodo">
          <v-text-field
            v-model="editingTodo.name"
            label="待办名称"
            variant="outlined"
            density="compact"
            class="mb-3"
            hide-details
            required
          />
          <v-textarea
            v-model="editingTodo.description"
            label="描述（可选）"
            variant="outlined"
            density="compact"
            rows="3"
            class="mb-3"
            hide-details
          />
          <v-text-field
            v-model="editingTodo.deadline"
            label="截止日期"
            variant="outlined"
            density="compact"
            type="date"
            hide-details
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showEditDialog = false">
            取消
          </v-btn>
          <v-btn
            color="primary"
            variant="flat"
            @click="saveEdit"
            :disabled="!editingTodo?.name"
          >
            保存
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<style scoped>
.todo-view {
  height: 100%;
  overflow-y: auto;
}

.todo-card {
  height: 100%;
  transition: all 0.3s;
}

.todo-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.todo-card.completed {
  opacity: 0.6;
}
</style>
