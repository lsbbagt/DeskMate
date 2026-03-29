<script lang="ts" setup>
import { ref } from 'vue'
import { useBookmarksStore } from '../stores/bookmarks'

const props = defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
}>()

const bookmarksStore = useBookmarksStore()

const showAddDialog = ref(false)
const newBookmarkName = ref('')
const newBookmarkUrl = ref('')

const openAddDialog = () => {
  newBookmarkName.value = ''
  newBookmarkUrl.value = ''
  showAddDialog.value = true
}

const addBookmark = () => {
  if (newBookmarkName.value && newBookmarkUrl.value) {
    let url = newBookmarkUrl.value
    if (!url.startsWith('http://') && !url.startsWith('https://')) {
      url = 'https://' + url
    }
    
    bookmarksStore.addBookmark({
      name: newBookmarkName.value,
      url: url
    })
    
    showAddDialog.value = false
    newBookmarkName.value = ''
    newBookmarkUrl.value = ''
  }
}

const deleteBookmark = (id: string) => {
  bookmarksStore.deleteBookmark(id)
}

const openBookmark = (url: string) => {
  bookmarksStore.setCurrentUrl(url)
}
</script>

<template>
  <v-navigation-drawer
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    width="300"
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
        @click="openBookmark(bookmark.url)"
        class="bookmark-item"
        rounded="lg"
      >
        <template #prepend>
          <v-icon color="primary" size="20">mdi-link</v-icon>
        </template>

        <v-list-item-title>{{ bookmark.name }}</v-list-item-title>
        <v-list-item-subtitle class="text-caption">{{ bookmark.url }}</v-list-item-subtitle>

        <template #append>
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
    <v-dialog v-model="showAddDialog" max-width="400">
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
          />
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
