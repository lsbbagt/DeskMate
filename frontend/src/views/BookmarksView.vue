<script lang="ts" setup>
import { computed } from 'vue'
import { useBookmarksStore } from '../stores/bookmarks'
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime'

const bookmarksStore = useBookmarksStore()

const currentUrl = computed(() => bookmarksStore.currentUrl)

const goToHome = () => {
  bookmarksStore.setCurrentUrl('about:blank')
}

// 检查URL是否可以在iframe中打开
const canOpenInIframe = (url: string) => {
  // 这些网站通常不允许在iframe中打开
  const blockedDomains = [
    'github.com',
    'stackoverflow.com',
    'google.com',
    'youtube.com',
    'facebook.com',
    'twitter.com',
    'linkedin.com',
    'kaggle.com'
  ]
  
  try {
    const hostname = new URL(url).hostname
    return !blockedDomains.some(domain => hostname.includes(domain))
  } catch {
    return false
  }
}

// 打开URL
const openUrl = (url: string) => {
  if (canOpenInIframe(url)) {
    bookmarksStore.setCurrentUrl(url)
  } else {
    // 使用外部浏览器打开
    BrowserOpenURL(url)
  }
}

// 监听store中的URL变化
bookmarksStore.$onAction(({ name, after }) => {
  if (name === 'setCurrentUrl') {
    after((url) => {
      if (url && url !== 'about:blank' && !canOpenInIframe(url)) {
        BrowserOpenURL(url)
        bookmarksStore.setCurrentUrl('about:blank')
      }
    })
  }
})
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
          <div class="text-body-2 text-secondary">
            从右侧选择收藏，或点击左上角菜单添加新收藏
          </div>
          <div class="text-caption mt-4 text-secondary">
            注：部分网站（如GitHub）将在外部浏览器中打开
          </div>
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
