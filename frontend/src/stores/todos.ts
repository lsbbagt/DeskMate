import { defineStore } from 'pinia'
import { ref } from 'vue'
import { AddTodo, GetTodos, DeleteTodo } from '../../wailsjs/go/main/App'

export interface Todo {
  id: string
  name: string
  description: string
  deadline: string
  completed: boolean
  createdAt: string
}

export const useTodosStore = defineStore('todos', () => {
  const todos = ref<Todo[]>([])

  const loadTodos = async () => {
    try {
      const result = await GetTodos()
      todos.value = result || []
    } catch (error) {
      console.error('加载待办失败:', error)
      todos.value = []
    }
  }

  const addTodo = async (todo: Omit<Todo, 'id' | 'completed' | 'createdAt'>) => {
    const newTodo: Todo = {
      id: Date.now().toString(),
      ...todo,
      completed: false,
      createdAt: new Date().toISOString()
    }
    
    todos.value.push(newTodo)
    
    try {
      await AddTodo(newTodo)
    } catch (error) {
      console.error('保存待办失败:', error)
    }
  }

  const deleteTodo = async (id: string) => {
    const index = todos.value.findIndex(t => t.id === id)
    if (index > -1) {
      todos.value.splice(index, 1)
      
      try {
        await DeleteTodo(id)
      } catch (error) {
        console.error('删除待办失败:', error)
      }
    }
  }

  const completeTodo = async (id: string) => {
    const todo = todos.value.find(t => t.id === id)
    if (todo) {
      todo.completed = true
      // 保存到后端
    }
  }

  return {
    todos,
    loadTodos,
    addTodo,
    deleteTodo,
    completeTodo
  }
})
