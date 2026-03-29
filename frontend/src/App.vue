<script lang="ts" setup>
import { ref } from 'vue'
import TitleBar from './components/TitleBar.vue'
import Sidebar from './components/Sidebar.vue'
import RightPanel from './components/RightPanel.vue'
import BookmarksView from './views/BookmarksView.vue'

const sidebarVisible = ref(true)
const rightPanelVisible = ref(true)
const activeModule = ref('bookmarks')

const toggleSidebar = () => {
  sidebarVisible.value = !sidebarVisible.value
}

const toggleRightPanel = () => {
  rightPanelVisible.value = !rightPanelVisible.value
}
</script>

<template>
  <!-- 自定义标题栏 - 放在 v-app 外部 -->
  <TitleBar 
    @toggle-sidebar="toggleSidebar"
    @toggle-right-panel="toggleRightPanel"
  />
  
  <v-app>
    <!-- 左侧导航栏 -->
    <Sidebar 
      v-model="sidebarVisible"
      v-model:active-module="activeModule"
    />
    
    <!-- 主内容区 -->
    <v-main class="main-content">
      <v-container fluid class="fill-height pa-0">
        <v-row no-gutters class="fill-height">
          <!-- 中间主面板 -->
          <v-col class="d-flex flex-column">
            <BookmarksView />
          </v-col>
          
          <!-- 右侧面板 -->
          <RightPanel 
            v-model="rightPanelVisible"
          />
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<style>
html, body {
  overflow: hidden !important;
  margin: 0;
  padding: 0;
  height: 100%;
}

#app {
  height: 100%;
}

/* 侧边栏从标题栏下方开始 */
.v-navigation-drawer {
  top: 40px !important;
  height: calc(100% - 40px) !important;
}

.v-application {
  padding-top: 40px;
}

.main-content {
  padding-top: 0 !important;
  height: calc(100vh - 40px);
  max-height: calc(100vh - 40px);
  overflow: hidden !important;
}

.main-content > .v-container {
  overflow: hidden !important;
  padding: 0 !important;
}

.main-content .v-row {
  overflow: hidden !important;
  margin: 0 !important;
}

/* 中间主面板允许滚动 */
.main-content .v-col:first-child {
  overflow-y: auto !important;
  overflow-x: hidden !important;
}

/* 右侧面板不滚动 */
.main-content .v-col:last-child {
  overflow: hidden !important;
}
</style>
