<script lang="ts" setup>
import { ref } from 'vue'
import { useBookmarksStore } from '../stores/bookmarks'
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime'
import type { Bookmark } from '../stores/bookmarks'

const props = defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
}>()

const bookmarksStore = useBookmarksStore()

const showAddDialog = ref(false)
const showEditDialog = ref(false)
const editingBookmark = ref<Bookmark | null>(null)
const newBookmarkName = ref('')
const newBookmarkUrl = ref('')
const newBookmarkOpenMode = ref<'internal' | 'external'>('internal')

const openAddDialog = () => {
  newBookmarkName.value = ''
  newBookmarkUrl.value = ''
  newBookmarkOpenMode.value = 'internal'
  showAddDialog.value = true
}

const addBookmark = async () => {
  if (newBookmarkName.value && newBookmarkUrl.value) {
    let url = newBookmarkUrl.value
    if (!url.startsWith('http://') && !url.startsWith('https://')) {
      url = 'https://' + url
    }
    
    await bookmarksStore.addBookmark({
      name: newBookmarkName.value,
      url: url,
      openMode: newBookmarkOpenMode.value
    })
    
    showAddDialog.value = false
    newBookmarkName.value = ''
    newBookmarkUrl.value = ''
  }
}

const openEditDialog = (bookmark: Bookmark) => {
  editingBookmark.value = { ...bookmark }
  showEditDialog.value = true
}

const saveEdit = async () => {
  if (editingBookmark.value) {
    await bookmarksStore.updateBookmark(editingBookmark.value.id, {
      name: editingBookmark.value.name,
      url: editingBookmark.value.url,
      openMode: editingBookmark.value.openMode
    })
    showEditDialog.value = false
    editingBookmark.value = null
  }
}

const deleteBookmark = async (id: string) => {
  await bookmarksStore.deleteBookmark(id)
}

const openBookmark = (bookmark: Bookmark) => {
  if (bookmark.openMode === 'external') {
    BrowserOpenURL(bookmark.url)
  } else {
    bookmarksStore.setCurrentUrl(bookmark.url)
  }
}

const openInExternal = (url: string) => {
  BrowserOpenURL(url)
}
</script>

<template>
  <v-navigation-drawer
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    width="320"
    location="right"
    class="right-panel"
    :permanent="true"
    color="surface"
  >
    <!-- 标题和操作按钮 -->
    <v-card flat class="pa-3">
      <div class="d-flex align-center justify-space-between mb-2">
        <div class="text-subtitle-1 font-weight-medium">
          收藏列表
        </div>
        <v-btn
          icon
          size="small"
          variant="text"
          color="primary"
          @click="openAddDialog"
        >
          <v-icon>mdi-plus</v-icon>
        </v-btn>
      </div>
    </v-card>

    <v-divider />

    <!-- 收藏列表 -->
    <v-list density="compact" class="py-0">
      <v-list-item
        v-for="bookmark in bookmarksStore.bookmarks"
        :key="bookmark.id"
        @click="openBookmark(bookmark)"
        class="bookmark-item"
        rounded="lg"
      >
        <template #prepend>
          <v-icon 
            :color="bookmark.openMode === 'external' ? 'warning' : 'primary'" 
            size="20"
          >
            {{ bookmark.openMode === 'external' ? 'mdi-open-in-new' : 'mdi-link' }}
          </v-icon>
        </template>

        <v-list-item-title class="d-flex align-center">
          {{ bookmark.name }}
          <v-chip
            :color="bookmark.openMode === 'external' ? 'warning' : 'primary'"
            size="x-small"
            label
            class="ml-2"
          >
            {{ bookmark.openMode === 'external' ? '外部' : '内部' }}
          </v-chip>
        </v-list-item-title>
        <v-list-item-subtitle class="text-caption">{{ bookmark.url }}</v-list-item-subtitle>

        <template #append>
          <v-btn
            icon
            size="x-small"
            variant="text"
            color="primary"
            @click.stop="openEditDialog(bookmark)"
            class="mr-1"
          >
            <v-icon size="14">mdi-pencil</v-icon>
          </v-btn>
          <v-btn
            icon
            size="x-small"
            variant="text"
            color="error"
            @click.stop="deleteBookmark(bookmark.id)"
          >
            <v-icon size="16">mdi-close</v-icon>
          </v-btn>
        </template>
      </v-list-item>
    </v-list>

    <!-- 空状态 -->
    <v-card v-if="bookmarksStore.bookmarks.length === 0" flat class="pa-6 text-center">
      <v-icon size="48" color="secondary">mdi-folder-open-outline</v-icon>
      <div class="text-body-2 mt-2 text-secondary">
        暂无收藏
      </div>
    </v-card>

    <!-- 添加对话框 -->
    <v-dialog v-model="showAddDialog" max-width="450">
      <v-card>
        <v-card-title class="text-h6">
          添加收藏
        </v-card-title>
        <v-card-text>
          <v-text-field
            v-model="newBookmarkName"
            label="名称"
            variant="outlined"
            density="compact"
            class="mb-3"
            hide-details
          />
          <v-text-field
            v-model="newBookmarkUrl"
            label="网址"
            variant="outlined"
            density="compact"
            hide-details
            placeholder="https://example.com"
            class="mb-3"
          />
          <div class="text-caption text-secondary mb-2">打开方式</div>
          <v-chip-group v-model="newBookmarkOpenMode" mandatory>
            <v-chip value="internal" color="primary" variant="outlined">
              <v-icon start>mdi-web</v-icon>
              内部打开
            </v-chip>
            <v-chip value="external" color="warning" variant="outlined">
              <v-icon start>mdi-open-in-new</v-icon>
              外部浏览器
            </v-chip>
          </v-chip-group>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn
            variant="text"
            @click="showAddDialog = false"
          >
            取消
          </v-btn>
          <v-btn
            color="primary"
            variant="flat"
            @click="addBookmark"
            :disabled="!newBookmarkName || !newBookmarkUrl"
          >
            添加
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- 编辑对话框 -->
    <v-dialog v-model="showEditDialog" max-width="450">
      <v-card v-if="editingBookmark">
        <v-card-title class="text-h6">
          编辑收藏
        </v-card-title>
        <v-card-text>
          <v-text-field
            v-model="editingBookmark.name"
            label="名称"
            variant="outlined"
            density="compact"
            class="mb-3"
            hide-details
          />
          <v-text-field
            v-model="editingBookmark.url"
            label="网址"
            variant="outlined"
            density="compact"
            hide-details
            class="mb-3"
          />
          <div class="text-caption text-secondary mb-2">打开方式</div>
          <v-chip-group v-model="editingBookmark.openMode" mandatory>
            <v-chip value="internal" color="primary" variant="outlined">
              <v-icon start>mdi-web</v-icon>
              内部打开
            </v-chip>
            <v-chip value="external" color="warning" variant="outlined">
              <v-icon start>mdi-open-in-new</v-icon>
              外部浏览器
            </v-chip>
          </v-chip-group>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn
            variant="text"
            @click="showEditDialog = false"
          >
            取消
          </v-btn>
          <v-btn
            color="primary"
            variant="flat"
            @click="saveEdit"
          >
            保存
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-navigation-drawer>
</template>

<style scoped>
.right-panel {
  border-left: 1px solid rgb(var(--v-theme-border));
}

.bookmark-item {
  margin: 4px 8px;
}
</style>
