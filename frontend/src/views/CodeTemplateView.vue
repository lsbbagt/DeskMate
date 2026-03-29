<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { useCodeTemplatesStore } from '../stores/codeTemplates'
import type { CodeFile } from '../stores/codeTemplates'

const codeStore = useCodeTemplatesStore()

const showAddFolderDialog = ref(false)
const showDeleteConfirm = ref(false)
const deleteType = ref<'folder' | 'file'>('folder')
const deleteTargetId = ref('')
const deleteFolderId = ref('')

onMounted(() => {
  codeStore.loadCodeTemplates()
})

// 添加文件夹
const newFolderName = ref('')
const openAddFolderDialog = () => {
  newFolderName.value = ''
  showAddFolderDialog.value = true
}

const addFolder = async () => {
  if (newFolderName.value) {
    await codeStore.addFolder(newFolderName.value)
    showAddFolderDialog.value = false
  }
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
    await codeStore.deleteFolder(deleteTargetId.value)
  } else {
    await codeStore.deleteFile(deleteFolderId.value, deleteTargetId.value)
  }
  showDeleteConfirm.value = false
}

// 导入文件
const importFile = async (folderId: string) => {
  await codeStore.selectAndAddFileToFolder(folderId)
}

// 查看文件
const viewFile = async (filePath: string) => {
  await codeStore.readFile(filePath)
}

// 使用默认应用打开
const openWithDefault = async (filePath: string) => {
  await codeStore.openWithDefault(filePath)
}

// 复制代码
const copyToClipboard = async () => {
  try {
    await navigator.clipboard.writeText(codeStore.currentFileContent)
  } catch (error) {
    console.error('复制失败:', error)
  }
}

// 获取文件图标
const getFileIcon = (type: string) => {
  const iconMap: Record<string, string> = {
    py: 'mdi-language-python',
    r: 'mdi-language-r',
    R: 'mdi-language-r',
    js: 'mdi-language-javascript',
    ts: 'mdi-language-typescript',
    java: 'mdi-language-java',
    cpp: 'mdi-language-cpp',
    c: 'mdi-language-c',
    go: 'mdi-language-go',
    exe: 'mdi-application'
  }
  return iconMap[type] || 'mdi-file-document'
}

// 获取语言类型（用于代码高亮显示）
const getLanguage = (type: string) => {
  const langMap: Record<string, string> = {
    py: 'Python',
    r: 'R',
    R: 'R',
    js: 'JavaScript',
    ts: 'TypeScript',
    java: 'Java',
    cpp: 'C++',
    c: 'C',
    go: 'Go'
  }
  return langMap[type] || 'Text'
}
</script>

<template>
  <div class="code-template-view">
    <v-row no-gutters class="fill-height">
      <!-- 中间：代码展示区域 -->
      <v-col class="d-flex flex-column code-display-area">
        <v-card flat class="fill-height">
          <!-- 顶部工具栏 -->
          <v-card-title class="d-flex align-center justify-space-between py-2 px-4">
            <div class="text-subtitle-1">
              {{ codeStore.currentFilePath ? codeStore.currentFilePath.split(/[/\\]/).pop() : '选择文件查看' }}
            </div>
            <v-btn
              v-if="codeStore.currentFileContent"
              color="primary"
              prepend-icon="mdi-content-copy"
              size="small"
              variant="flat"
              @click="copyToClipboard"
            >
              复制
            </v-btn>
          </v-card-title>

          <v-divider />

          <!-- 代码内容 -->
          <v-card-text class="code-content pa-0">
            <pre v-if="codeStore.currentFileContent" class="code-block">{{ codeStore.currentFileContent }}</pre>
            <div v-else class="empty-state text-center pa-8">
              <v-icon size="64" color="secondary">mdi-code-tags</v-icon>
              <div class="text-h6 mt-4 mb-2">暂无代码</div>
              <div class="text-body-2 text-secondary">
                从右侧选择文件查看代码内容
              </div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- 右侧：文件夹和文件列表 -->
      <v-col cols="auto">
        <v-navigation-drawer
          width="320"
          location="right"
          class="code-sidebar"
          :permanent="true"
          color="surface"
        >
          <!-- 标题和操作 -->
          <v-card flat class="pa-3">
            <div class="d-flex align-center justify-space-between">
              <div class="text-subtitle-1 font-weight-medium">代码模板</div>
              <v-btn
                icon
                size="small"
                variant="text"
                color="primary"
                @click="openAddFolderDialog"
              >
                <v-icon>mdi-folder-plus</v-icon>
              </v-btn>
            </div>
          </v-card>

          <v-divider />

          <!-- 文件夹列表 -->
          <v-list density="compact" class="py-0">
            <template v-for="folder in codeStore.folders" :key="folder.id">
              <v-list-group>
                <template #activator="{ props }">
                  <v-list-item
                    v-bind="props"
                    :prepend-icon="'mdi-folder'"
                    class="folder-item"
                  >
                    <v-list-item-title>{{ folder.name }}</v-list-item-title>
                    <template #append>
                      <v-btn
                        icon
                        size="x-small"
                        variant="text"
                        color="error"
                        @click.stop="confirmDeleteFolder(folder.id)"
                        class="mr-1"
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
                  :prepend-icon="getFileIcon(file.type)"
                  @click="viewFile(file.path)"
                  class="file-item"
                  rounded="lg"
                >
                  <v-list-item-title>{{ file.name }}</v-list-item-title>
                  <template #append>
                    <v-btn
                      icon
                      size="x-small"
                      variant="text"
                      color="primary"
                      @click.stop="openWithDefault(file.path)"
                      class="mr-1"
                    >
                      <v-icon size="14">mdi-launch</v-icon>
                    </v-btn>
                    <v-btn
                      icon
                      size="x-small"
                      variant="text"
                      color="error"
                      @click.stop="confirmDeleteFile(folder.id, file.id)"
                    >
                      <v-icon size="16">mdi-close</v-icon>
                    </v-btn>
                  </template>
                </v-list-item>

                <!-- 导入文件按钮 -->
                <v-list-item @click="importFile(folder.id)" class="import-item">
                  <v-list-item-title class="text-center text-primary">
                    <v-icon start>mdi-plus</v-icon>
                    导入文件
                  </v-list-item-title>
                </v-list-item>
              </v-list-group>
            </template>
          </v-list>

          <!-- 空状态 -->
          <v-card v-if="codeStore.folders.length === 0" flat class="pa-6 text-center">
            <v-icon size="48" color="secondary">mdi-folder-open-outline</v-icon>
            <div class="text-body-2 mt-2 text-secondary">
              暂无文件夹<br>点击右上角创建
            </div>
          </v-card>
        </v-navigation-drawer>
      </v-col>
    </v-row>

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
.code-template-view {
  height: 100%;
  overflow: hidden;
}

.code-display-area {
  height: 100%;
}

.code-sidebar {
  border-left: 1px solid rgb(var(--v-theme-border));
}

.code-content {
  height: calc(100% - 60px);
  overflow-y: auto;
  background-color: rgb(var(--v-theme-surface-variant));
}

.code-block {
  margin: 0;
  padding: 16px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.5;
  background-color: #1e1e1e;
  color: #d4d4d4;
  white-space: pre;
  overflow-x: auto;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.folder-item {
  margin: 4px 8px;
}

.file-item {
  margin: 2px 8px 2px 24px;
  padding-left: 8px;
}

.import-item {
  margin: 4px 8px 4px 24px;
  padding-left: 8px;
  border: 1px dashed rgb(var(--v-theme-border));
}
</style>
