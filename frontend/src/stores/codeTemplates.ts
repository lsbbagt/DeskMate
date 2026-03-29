import { defineStore } from 'pinia'
import { ref } from 'vue'
import { 
  GetCodeFolders,
  SaveCodeTemplates,
  LoadCodeTemplates,
  ReadFileContent,
  OpenWithDefaultApp,
  SelectCodeFile
} from '../../wailsjs/go/main/App'

export interface CodeFile {
  id: string
  name: string
  path: string
  type: string
}

export interface CodeFolder {
  id: string
  name: string
  files: CodeFile[]
}

export const useCodeTemplatesStore = defineStore('codeTemplates', () => {
  const folders = ref<CodeFolder[]>([])
  const currentFileContent = ref('')
  const currentFilePath = ref('')

  // 加载代码模板
  const loadCodeTemplates = async () => {
    try {
      const data = await LoadCodeTemplates()
      folders.value = JSON.parse(data)
    } catch (error) {
      console.error('加载代码模板失败:', error)
      folders.value = []
    }
  }

  // 保存代码模板
  const saveCodeTemplates = async () => {
    try {
      await SaveCodeTemplates(JSON.stringify(folders.value))
    } catch (error) {
      console.error('保存代码模板失败:', error)
    }
  }

  // 添加文件夹
  const addFolder = async (name: string) => {
    const folder: CodeFolder = {
      id: Date.now().toString(),
      name,
      files: []
    }
    folders.value.push(folder)
    await saveCodeTemplates()
  }

  // 删除文件夹
  const deleteFolder = async (folderId: string) => {
    folders.value = folders.value.filter(f => f.id !== folderId)
    await saveCodeTemplates()
  }

  // 选择并添加文件到文件夹
  const selectAndAddFileToFolder = async (folderId: string) => {
    try {
      const filePath = await SelectCodeFile()
      
      if (filePath) {
        const fileName = filePath.split(/[/\\]/).pop() || ''
        const ext = fileName.split('.').pop()?.toLowerCase() || ''
        
        const file: CodeFile = {
          id: Date.now().toString(),
          name: fileName,
          path: filePath,
          type: ext
        }
        
        const folder = folders.value.find(f => f.id === folderId)
        if (folder) {
          folder.files.push(file)
          await saveCodeTemplates()
        }
      }
    } catch (error) {
      console.error('选择文件失败:', error)
    }
  }

  // 删除文件
  const deleteFile = async (folderId: string, fileId: string) => {
    const folder = folders.value.find(f => f.id === folderId)
    if (folder) {
      folder.files = folder.files.filter(f => f.id !== fileId)
      await saveCodeTemplates()
    }
  }

  // 读取文件内容
  const readFile = async (filePath: string) => {
    try {
      const content = await ReadFileContent(filePath)
      currentFileContent.value = content
      currentFilePath.value = filePath
    } catch (error) {
      console.error('读取文件失败:', error)
      currentFileContent.value = ''
    }
  }

  // 使用默认应用打开
  const openWithDefault = async (filePath: string) => {
    try {
      await OpenWithDefaultApp(filePath)
    } catch (error) {
      console.error('打开文件失败:', error)
    }
  }

  return {
    folders,
    currentFileContent,
    currentFilePath,
    loadCodeTemplates,
    saveCodeTemplates,
    addFolder,
    deleteFolder,
    selectAndAddFileToFolder,
    deleteFile,
    readFile,
    openWithDefault
  }
})

