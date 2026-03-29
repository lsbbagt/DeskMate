import { defineStore } from 'pinia'
import { ref } from 'vue'
import { 
  SaveMarkdownFiles,
  LoadMarkdownFiles,
  ReadFileContent,
  SaveMarkdownContent,
  SelectMarkdownFile,
  ImportMarkdownFile,
  DeleteMarkdownFile,
  CreateMarkdownFile,
  DeleteMarkdownFolder,
  SyncMarkdownFiles
} from '../../wailsjs/go/main/App'

export interface MarkdownFile {
  id: string
  name: string
  path: string
  content?: string
}

export interface MarkdownFolder {
  id: string
  name: string
  files: MarkdownFile[]
}

export const useMarkdownFilesStore = defineStore('markdownFiles', () => {
  const folders = ref<MarkdownFolder[]>([])
  const currentFileContent = ref('')
  const currentFilePath = ref('')
  const currentFileName = ref('')

  // 加载 Markdown 文件
  const loadMarkdownFiles = async () => {
    try {
      // 先同步文件系统
      await SyncMarkdownFiles()
      // 再加载
      const data = await LoadMarkdownFiles()
      folders.value = JSON.parse(data)
    } catch (error) {
      console.error('加载 Markdown 文件失败:', error)
      folders.value = []
    }
  }

  // 保存 Markdown 文件列表
  const saveMarkdownFiles = async () => {
    try {
      await SaveMarkdownFiles(JSON.stringify(folders.value))
    } catch (error) {
      console.error('保存 Markdown 文件失败:', error)
    }
  }

  // 添加文件夹
  const addFolder = async (name: string) => {
    const folder: MarkdownFolder = {
      id: Date.now().toString(),
      name,
      files: []
    }
    folders.value.push(folder)
    await saveMarkdownFiles()
  }

  // 删除文件夹
  const deleteFolder = async (folderId: string) => {
    const folder = folders.value.find(f => f.id === folderId)
    if (folder) {
      try {
        // 使用文件夹名称删除实际文件夹
        await DeleteMarkdownFolder(folder.name)
      } catch (error) {
        console.error('删除文件夹失败:', error)
      }
    }
    folders.value = folders.value.filter(f => f.id !== folderId)
    await saveMarkdownFiles()
  }

  // 创建新的 Markdown 文件
  const createMarkdownFile = async (folderId: string, fileName: string) => {
    try {
      const folder = folders.value.find(f => f.id === folderId)
      if (!folder) return
      
      const finalName = fileName.endsWith('.md') ? fileName : `${fileName}.md`
      
      // 使用文件夹名称创建文件
      const targetPath = await CreateMarkdownFile(folder.name, finalName)
      
      const file: MarkdownFile = {
        id: Date.now().toString(),
        name: finalName,
        path: targetPath,
        content: ''
      }
      folder.files.push(file)
      await saveMarkdownFiles()
      
      // 自动选中新创建的文件
      currentFileContent.value = ''
      currentFilePath.value = targetPath
      currentFileName.value = finalName
    } catch (error) {
      console.error('创建 Markdown 文件失败:', error)
    }
  }

  // 选择并导入文件
  const selectAndImportFile = async (folderId: string) => {
    try {
      const sourcePath = await SelectMarkdownFile()
      
      if (sourcePath) {
        const folder = folders.value.find(f => f.id === folderId)
        if (!folder) return
        
        // 使用文件夹名称作为目录名
        const targetPath = await ImportMarkdownFile(sourcePath, folder.name)
        
        const fileName = sourcePath.split(/[/\\]/).pop() || ''
        
        const file: MarkdownFile = {
          id: Date.now().toString(),
          name: fileName,
          path: targetPath
        }
        
        folder.files.push(file)
        await saveMarkdownFiles()
      }
    } catch (error) {
      console.error('导入文件失败:', error)
    }
  }

  // 删除文件
  const deleteFile = async (folderId: string, fileId: string) => {
    const folder = folders.value.find(f => f.id === folderId)
    if (folder) {
      const file = folder.files.find(f => f.id === fileId)
      if (file && file.path) {
        // 删除实际文件
        try {
          await DeleteMarkdownFile(file.path)
        } catch (error) {
          console.error('删除文件失败:', error)
        }
      }
      
      folder.files = folder.files.filter(f => f.id !== fileId)
      await saveMarkdownFiles()
      
      // 如果删除的是当前文件，清空编辑器
      if (currentFilePath.value && folder.files.find(f => f.id === fileId)) {
        currentFileContent.value = ''
        currentFilePath.value = ''
        currentFileName.value = ''
      }
    }
  }

  // 打开文件
  const openFile = async (file: MarkdownFile) => {
    try {
      if (file.path) {
        // 从文件读取
        const content = await ReadFileContent(file.path)
        currentFileContent.value = content
        currentFilePath.value = file.path
      } else {
        // 如果没有路径（不应该发生），加载保存的内容
        currentFileContent.value = file.content || ''
        currentFilePath.value = ''
      }
      currentFileName.value = file.name
    } catch (error) {
      console.error('打开文件失败:', error)
      currentFileContent.value = ''
    }
  }

  // 保存文件内容
  const saveFile = async (content: string) => {
    try {
      if (currentFilePath.value) {
        // 保存到文件
        await SaveMarkdownContent(currentFilePath.value, content)
      } else {
        // 如果没有路径（不应该发生），保存到内存
        for (const folder of folders.value) {
          const file = folder.files.find(f => f.name === currentFileName.value)
          if (file) {
            file.content = content
            break
          }
        }
        await saveMarkdownFiles()
      }
    } catch (error) {
      console.error('保存文件失败:', error)
    }
  }

  return {
    folders,
    currentFileContent,
    currentFilePath,
    currentFileName,
    loadMarkdownFiles,
    saveMarkdownFiles,
    addFolder,
    deleteFolder,
    createMarkdownFile,
    selectAndImportFile,
    deleteFile,
    openFile,
    saveFile
  }
})
