import { defineStore } from 'pinia'
import { ref } from 'vue'
import { LoadBookmarks, SaveBookmarks } from '../../wailsjs/go/main/App'

export interface Bookmark {
  id: string
  name: string
  url: string
  openMode: 'internal' | 'external' // 内部打开还是外部浏览器
  createdAt: number
}

export const useBookmarksStore = defineStore('bookmarks', () => {
  const bookmarks = ref<Bookmark[]>([])

  const currentUrl = ref<string>('about:blank')

  const addBookmark = async (bookmark: Omit<Bookmark, 'id' | 'createdAt'>) => {
    const newBookmark: Bookmark = {
      ...bookmark,
      id: Date.now().toString(),
      createdAt: Date.now()
    }
    bookmarks.value.push(newBookmark)
    await saveToStorage()
  }

  const updateBookmark = async (id: string, updates: Partial<Bookmark>) => {
    const index = bookmarks.value.findIndex(b => b.id === id)
    if (index !== -1) {
      bookmarks.value[index] = { ...bookmarks.value[index], ...updates }
      await saveToStorage()
    }
  }

  const deleteBookmark = async (id: string) => {
    const index = bookmarks.value.findIndex(b => b.id === id)
    if (index !== -1) {
      bookmarks.value.splice(index, 1)
      await saveToStorage()
    }
  }

  const setCurrentUrl = (url: string) => {
    currentUrl.value = url
  }

  const saveToStorage = async () => {
    try {
      const jsonData = JSON.stringify(bookmarks.value, null, 2)
      await SaveBookmarks(jsonData)
    } catch (error) {
      console.error('保存收藏失败:', error)
    }
  }

  const loadFromStorage = async () => {
    try {
      const jsonData = await LoadBookmarks()
      if (jsonData) {
        bookmarks.value = JSON.parse(jsonData)
      }
    } catch (error) {
      console.error('加载收藏失败:', error)
      bookmarks.value = []
    }
  }

  // 初始化时加载数据
  loadFromStorage()

  return {
    bookmarks,
    currentUrl,
    addBookmark,
    updateBookmark,
    deleteBookmark,
    setCurrentUrl,
    saveToStorage,
    loadFromStorage
  }
})
