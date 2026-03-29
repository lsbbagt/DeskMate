<script lang="ts" setup>
import { ref, onMounted, onUnmounted, watch, computed } from 'vue'
import { useMarkdownFilesStore } from '../stores/markdownFiles'
import type { MarkdownFile } from '../stores/markdownFiles'
import { renderMarkdown } from '../utils/markdownRenderer'
import 'katex/dist/katex.min.css'
import 'highlight.js/styles/github.css'

const props = defineProps<{
  rightPanelVisible: boolean
}>()

const markdownStore = useMarkdownFilesStore()

const showAddFolderDialog = ref(false)
const showAddFileDialog = ref(false)
const showDeleteConfirm = ref(false)
const deleteType = ref<'folder' | 'file'>('folder')
const deleteTargetId = ref('')
const deleteFolderId = ref('')
const selectedFolderId = ref('')
const newFileName = ref('')

const editorContent = ref('')
const previewVisible = ref(true)
const previewContent = computed(() => {
  return renderMarkdown(editorContent.value)
})

// 保存反馈状态
const saveFeedback = ref(false)

onMounted(() => {
  markdownStore.loadMarkdownFiles()
  // 添加 Ctrl+S 快捷键监听
  window.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  // 移除快捷键监听
  window.removeEventListener('keydown', handleKeyDown)
})

// 键盘事件处理
const handleKeyDown = (e: KeyboardEvent) => {
  // Ctrl+S 或 Cmd+S (Mac)
  if ((e.ctrlKey || e.metaKey) && e.key === 's') {
    e.preventDefault()
    if (markdownStore.currentFileName) {
      manualSave()
    }
  }
}

// 监听当前文件变化，更新编辑器内容
watch(() => markdownStore.currentFileContent, (newContent) => {
  editorContent.value = newContent
})

// 自动保存（防抖）
let saveTimer: ReturnType<typeof setTimeout> | null = null
watch(editorContent, (newContent) => {
  if (saveTimer) clearTimeout(saveTimer)
  saveTimer = setTimeout(() => {
    if (markdownStore.currentFileName) {
      markdownStore.saveFile(newContent)
    }
  }, 1000)
})

// 手动保存
const manualSave = async () => {
  if (!markdownStore.currentFileName) return
  
  await markdownStore.saveFile(editorContent.value)
  
  // 显示保存反馈
  saveFeedback.value = true
  setTimeout(() => {
    saveFeedback.value = false
  }, 2000)
}

// 添加文件夹
const newFolderName = ref('')
const openAddFolderDialog = () => {
  newFolderName.value = ''
  showAddFolderDialog.value = true
}

const addFolder = async () => {
  if (newFolderName.value) {
    await markdownStore.addFolder(newFolderName.value)
    showAddFolderDialog.value = false
  }
}

// 添加文件
const openAddFileDialog = (folderId: string) => {
  selectedFolderId.value = folderId
  newFileName.value = ''
  showAddFileDialog.value = true
}

const addFile = async () => {
  if (newFileName.value && selectedFolderId.value) {
    await markdownStore.createMarkdownFile(selectedFolderId.value, newFileName.value)
    showAddFileDialog.value = false
  }
}

// 导入文件
const importFile = async (folderId: string) => {
  await markdownStore.selectAndImportFile(folderId)
}

// 删除文件夹
const confirmDeleteFolder = (folderId: string) => {
  deleteType.value = 'folder'
  deleteTargetId.value = folderId
  showDeleteConfirm.value = true
}

// 删除文件
const confirmDeleteFile = (folderId: string, fileId: string) => {
  deleteType.value = 'file'
  deleteFolderId.value = folderId
  deleteTargetId.value = fileId
  showDeleteConfirm.value = true
}

const executeDelete = async () => {
  if (deleteType.value === 'folder') {
    await markdownStore.deleteFolder(deleteTargetId.value)
  } else {
    await markdownStore.deleteFile(deleteFolderId.value, deleteTargetId.value)
  }
  showDeleteConfirm.value = false
}

// 打开文件
const openFile = async (file: MarkdownFile) => {
  await markdownStore.openFile(file)
}
</script>

<template>
  <div class="markdown-editor-view">
    <div class="main-container">
      <!-- 左侧：编辑器 -->
      <div class="editor-area">
        <v-card flat class="fill-height">
          <v-card-title class="d-flex align-center justify-space-between py-2 px-4">
            <div class="text-subtitle-1">
              {{ markdownStore.currentFileName || '新建文档' }}
            </div>
            <div class="d-flex align-center">
              <!-- 保存反馈 -->
              <v-fade-transition>
                <div v-if="saveFeedback" class="save-feedback mr-2">
                  <v-icon color="success" size="16">mdi-check-circle</v-icon>
                  <span class="text-success text-body-2 ml-1">已保存</span>
                </div>
              </v-fade-transition>
              
              <v-btn
                v-if="markdownStore.currentFileName"
                color="primary"
                prepend-icon="mdi-content-save"
                size="small"
                variant="flat"
                @click="manualSave"
              >
                保存
              </v-btn>
            </div>
          </v-card-title>

          <v-divider />

          <textarea
            v-model="editorContent"
            class="editor-textarea"
            placeholder="在这里输入 Markdown 内容..."
          />
        </v-card>
      </div>

      <!-- 中间：预览区 -->
      <div v-if="previewVisible" class="preview-area">
        <v-card flat class="fill-height">
          <v-card-title class="d-flex align-center justify-space-between py-2 px-4">
            <div class="text-subtitle-1">预览</div>
            <v-btn
              icon
              size="small"
              variant="text"
              @click="previewVisible = false"
            >
              <v-icon>mdi-chevron-right</v-icon>
            </v-btn>
          </v-card-title>

          <v-divider />

          <div class="preview-content" v-html="previewContent" />
        </v-card>
      </div>

      <!-- 预览区折叠时的展开按钮 -->
      <div v-else class="preview-toggle-btn">
        <v-btn
          icon
          size="small"
          variant="text"
          @click="previewVisible = true"
        >
          <v-icon>mdi-chevron-left</v-icon>
        </v-btn>
      </div>

      <!-- 右侧：文件夹和文件列表 -->
      <div v-if="rightPanelVisible" class="markdown-sidebar">
        <div class="sidebar-content">
          <!-- 标题 -->
          <v-card flat class="pa-3">
            <div class="text-subtitle-1 font-weight-medium">Markdown 文件</div>
          </v-card>

          <v-divider />

          <!-- 文件夹列表 -->
          <v-list density="compact" class="py-0 file-list">
            <template v-for="folder in markdownStore.folders" :key="folder.id">
              <v-list-group>
                <template #activator="{ props }">
                  <v-list-item
                    v-bind="props"
                    :prepend-icon="'mdi-folder'"
                    class="folder-item"
                  >
                    <v-list-item-title class="text-truncate">{{ folder.name }}</v-list-item-title>
                    <template #append>
                      <v-btn
                        icon
                        size="x-small"
                        variant="text"
                        color="error"
                        @click.stop="confirmDeleteFolder(folder.id)"
                      >
                        <v-icon size="16">mdi-delete</v-icon>
                      </v-btn>
                    </template>
                  </v-list-item>
                </template>

                <!-- 文件列表 -->
                <v-list-item
                  v-for="file in folder.files"
                  :key="file.id"
                  :prepend-icon="'mdi-language-markdown'"
                  @click="openFile(file)"
                  class="file-item"
                  rounded="lg"
                >
                  <v-list-item-title class="text-truncate">{{ file.name }}</v-list-item-title>
                  <template #append>
                    <v-btn
                      icon
                      size="x-small"
                      variant="text"
                      color="error"
                      @click.stop="confirmDeleteFile(folder.id, file.id)"
                    >
                      <v-icon size="14">mdi-close</v-icon>
                    </v-btn>
                  </template>
                </v-list-item>

                <!-- 操作按钮 -->
                <v-list-item @click="openAddFileDialog(folder.id)" class="action-item">
                  <v-list-item-title class="text-center text-primary">
                    <v-icon start size="16">mdi-plus</v-icon>
                    新建文件
                  </v-list-item-title>
                </v-list-item>
                <v-list-item @click="importFile(folder.id)" class="action-item">
                  <v-list-item-title class="text-center text-primary">
                    <v-icon start size="16">mdi-file-import</v-icon>
                    导入文件
                  </v-list-item-title>
                </v-list-item>
              </v-list-group>
            </template>
          </v-list>

          <!-- 空状态 -->
          <v-card v-if="markdownStore.folders.length === 0" flat class="pa-6 text-center">
            <v-icon size="48" color="secondary">mdi-folder-open-outline</v-icon>
            <div class="text-body-2 mt-2 text-secondary">
              暂无文件夹<br>点击下方创建
            </div>
          </v-card>

          <!-- 添加文件夹按钮 -->
          <div class="pa-3">
            <v-btn
              block
              color="primary"
              variant="outlined"
              prepend-icon="mdi-folder-plus"
              @click="openAddFolderDialog"
            >
              创建文件夹
            </v-btn>
          </div>
        </div>
      </div>
    </div>

    <!-- 添加文件夹对话框 -->
    <v-dialog v-model="showAddFolderDialog" max-width="400">
      <v-card>
        <v-card-title class="text-h6">创建文件夹</v-card-title>
        <v-card-text>
          <v-text-field
            v-model="newFolderName"
            label="文件夹名称"
            variant="outlined"
            density="compact"
            hide-details
            @keyup.enter="addFolder"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showAddFolderDialog = false">取消</v-btn>
          <v-btn
            color="primary"
            variant="flat"
            @click="addFolder"
            :disabled="!newFolderName"
          >
            创建
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- 添加文件对话框 -->
    <v-dialog v-model="showAddFileDialog" max-width="400">
      <v-card>
        <v-card-title class="text-h6">新建 Markdown 文件</v-card-title>
        <v-card-text>
          <v-text-field
            v-model="newFileName"
            label="文件名称"
            variant="outlined"
            density="compact"
            hide-details
            hint="不需要添加 .md 后缀"
            @keyup.enter="addFile"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showAddFileDialog = false">取消</v-btn>
          <v-btn
            color="primary"
            variant="flat"
            @click="addFile"
            :disabled="!newFileName"
          >
            创建
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- 删除确认对话框 -->
    <v-dialog v-model="showDeleteConfirm" max-width="400">
      <v-card>
        <v-card-title class="text-h6">确认删除</v-card-title>
        <v-card-text>
          确定要删除这个{{ deleteType === 'folder' ? '文件夹及其所有文件' : '文件' }}吗？
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showDeleteConfirm = false">取消</v-btn>
          <v-btn
            color="error"
            variant="flat"
            @click="executeDelete"
          >
            删除
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<style scoped>
.markdown-editor-view {
  height: 100%;
  width: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.main-container {
  display: flex;
  height: 100%;
  width: 100%;
}

.editor-area {
  flex: 1;
  height: 100%;
  min-width: 0;
  display: flex;
}

.editor-area > .v-card {
  flex: 1;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.editor-textarea {
  flex: 1;
  width: 100%;
  padding: 16px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.6;
  border: none;
  outline: none;
  resize: none;
  background-color: #fafafa;
}

.preview-area {
  flex: 1;
  height: 100%;
  min-width: 0;
  border-left: 1px solid rgb(var(--v-theme-border));
  display: flex;
}

.preview-area > .v-card {
  flex: 1;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.preview-content {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
  background-color: #ffffff;
  line-height: 1.8;
}

/* Markdown 样式 */
.preview-content :deep(h1),
.preview-content :deep(h2),
.preview-content :deep(h3),
.preview-content :deep(h4),
.preview-content :deep(h5),
.preview-content :deep(h6) {
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
  line-height: 1.25;
}

.preview-content :deep(h1) {
  font-size: 2em;
  border-bottom: 1px solid #eaecef;
  padding-bottom: 0.3em;
}

.preview-content :deep(h2) {
  font-size: 1.5em;
  border-bottom: 1px solid #eaecef;
  padding-bottom: 0.3em;
}

.preview-content :deep(h3) { font-size: 1.25em; }
.preview-content :deep(h4) { font-size: 1em; }
.preview-content :deep(h5) { font-size: 0.875em; }
.preview-content :deep(h6) { font-size: 0.85em; color: #6a737d; }

.preview-content :deep(p) {
  margin-bottom: 16px;
}

.preview-content :deep(ul),
.preview-content :deep(ol) {
  padding-left: 2em;
  margin-bottom: 16px;
}

.preview-content :deep(li) {
  margin-bottom: 0.25em;
}

.preview-content :deep(blockquote) {
  padding: 0 1em;
  color: #6a737d;
  border-left: 0.25em solid #dfe2e5;
  margin-bottom: 16px;
}

.preview-content :deep(code) {
  padding: 0.2em 0.4em;
  margin: 0;
  font-size: 85%;
  background-color: rgba(27, 31, 35, 0.05);
  border-radius: 3px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
}

.preview-content :deep(pre) {
  padding: 16px;
  overflow: auto;
  font-size: 85%;
  line-height: 1.45;
  background-color: #f6f8fa;
  border-radius: 3px;
  margin-bottom: 16px;
}

.preview-content :deep(pre code) {
  background-color: transparent;
  padding: 0;
}

.preview-content :deep(table) {
  border-spacing: 0;
  border-collapse: collapse;
  margin-bottom: 16px;
  width: 100%;
}

.preview-content :deep(table th),
.preview-content :deep(table td) {
  padding: 6px 13px;
  border: 1px solid #dfe2e5;
}

.preview-content :deep(table th) {
  font-weight: 600;
  background-color: #f6f8fa;
}

.preview-content :deep(table tr:nth-child(2n)) {
  background-color: #f6f8fa;
}

.preview-content :deep(img) {
  max-width: 100%;
  box-sizing: content-box;
  background-color: #fff;
}

.preview-content :deep(a) {
  color: #0366d6;
  text-decoration: none;
}

.preview-content :deep(a:hover) {
  text-decoration: underline;
}

.preview-content :deep(hr) {
  height: 0.25em;
  padding: 0;
  margin: 24px 0;
  background-color: #e1e4e8;
  border: 0;
}

/* KaTeX 样式优化 */
.preview-content :deep(.katex-display) {
  margin: 1em 0;
  overflow-x: auto;
  overflow-y: hidden;
}

.preview-content :deep(.katex) {
  font-size: 1.1em;
}

/* 预览区折叠按钮 */
.preview-toggle-btn {
  display: flex;
  align-items: center;
  border-left: 1px solid rgb(var(--v-theme-border));
  background-color: rgb(var(--v-theme-surface));
  padding: 0 4px;
}

.preview-toggle-btn .v-btn {
  height: 100%;
}

.markdown-sidebar {
  width: 320px;
  min-width: 320px;
  border-left: 1px solid rgb(var(--v-theme-border));
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: rgb(var(--v-theme-surface));
}

.sidebar-content {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.file-list {
  padding-right: 8px;
}

.file-list :deep(.v-list-group__items) {
  padding-left: 0 !important;
  margin-left: 0 !important;
}

.file-list :deep(.v-list-group__items .v-list-item) {
  padding-left: 8px !important;
}

.folder-item,
.file-item,
.action-item {
  margin: 2px 0;
  padding: 0 8px;
}

.text-truncate {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.action-item {
  border: 1px dashed rgb(var(--v-theme-border));
  margin: 2px 8px;
}

/* 保存反馈动画 */
.save-feedback {
  display: flex;
  align-items: center;
}
</style>
