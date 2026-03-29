<script lang="ts" setup>
import { computed } from 'vue'
import { useBookmarksStore } from '../stores/bookmarks'
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime'

const bookmarksStore = useBookmarksStore()

const currentUrl = computed(() => bookmarksStore.currentUrl)

const goToHome = () => {
  bookmarksStore.setCurrentUrl('about:blank')
}

const openInExternal = () => {
  if (currentUrl.value && currentUrl.value !== 'about:blank') {
    BrowserOpenURL(currentUrl.value)
  }
}

const copyUrl = async () => {
  if (currentUrl.value && currentUrl.value !== 'about:blank') {
    try {
      await navigator.clipboard.writeText(currentUrl.value)
      // 可以添加提示
    } catch (error) {
      console.error('复制失败:', error)
    }
  }
}
</script>

<template>
  <div class="bookmarks-view d-flex flex-column fill-height">
    <!-- 浏览器工具栏 -->
    <v-toolbar flat density="compact" color="surface">
      <v-btn
        icon
        size="small"
        variant="text"
        @click="goToHome"
        title="回到首页"
      >
        <v-icon>mdi-home</v-icon>
      </v-btn>

      <v-text-field
        :model-value="currentUrl"
        readonly
        density="compact"
        variant="outlined"
        hide-details
        class="url-input mx-2"
        placeholder="输入网址..."
      />

      <v-btn
        icon
        size="small"
        variant="text"
        @click="openInExternal"
        title="在外部浏览器打开"
        :disabled="!currentUrl || currentUrl === 'about:blank'"
      >
        <v-icon>mdi-open-in-new</v-icon>
      </v-btn>

      <v-btn
        icon
        size="small"
        variant="text"
        @click="copyUrl"
        title="复制链接"
        :disabled="!currentUrl || currentUrl === 'about:blank'"
      >
        <v-icon>mdi-content-copy</v-icon>
      </v-btn>
    </v-toolbar>

    <v-divider />

    <!-- 浏览器内容 -->
    <div class="browser-container flex-grow-1">
      <iframe
        v-if="currentUrl && currentUrl !== 'about:blank'"
        :src="currentUrl"
        class="browser-frame"
        sandbox="allow-same-origin allow-scripts allow-popups allow-forms allow-top-navigation"
      />
      
      <!-- 空白页 -->
      <div v-else class="blank-page d-flex align-center justify-center">
        <v-card flat max-width="400" class="text-center pa-8">
          <v-icon size="64" color="primary">mdi-web</v-icon>
          <div class="text-h6 mt-4 mb-2">欢迎使用StatBox</div>
          <div class="text-body-2 text-secondary mb-4">
            从右侧选择收藏，或点击左上角菜单添加新收藏
          </div>
          <v-alert type="info" density="compact" variant="tonal" class="text-left">
            <div class="text-caption">
              <strong>提示：</strong>
              <ul class="mt-1 pl-3">
                <li>内部打开：在应用内浏览器中显示</li>
                <li>外部打开：使用系统浏览器打开</li>
                <li>部分网站可能无法在内部显示，建议选择外部打开</li>
              </ul>
            </div>
          </v-alert>
        </v-card>
      </div>
    </div>
  </div>
</template>

<style scoped>
.bookmarks-view {
  background-color: rgb(var(--v-theme-background));
}

.url-input {
  flex: 1;
  max-width: 800px;
}

.browser-container {
  position: relative;
  overflow: hidden;
}

.browser-frame {
  width: 100%;
  height: 100%;
  border: none;
  background: white;
}

.blank-page {
  height: 100%;
  background-color: rgb(var(--v-theme-background));
}
</style>
