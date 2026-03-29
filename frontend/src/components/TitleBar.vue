<script lang="ts" setup>
import { WindowMinimise, WindowToggleMaximise, Quit } from '../../wailsjs/runtime/runtime'

const emit = defineEmits<{
  (e: 'toggle-sidebar'): void
  (e: 'toggle-right-panel'): void
}>()

const minimize = () => {
  WindowMinimise()
}

const maximize = () => {
  WindowToggleMaximise()
}

const close = () => {
  Quit()
}
</script>

<template>
  <div class="title-bar" style="--wails-draggable:drag">
    <div class="title-bar-content">
      <!-- 左侧控制按钮 -->
      <button class="title-btn" @click="$emit('toggle-sidebar')">
        <v-icon size="20">mdi-menu</v-icon>
      </button>

      <!-- 标题 -->
      <div class="title-text">StatBox</div>

      <div class="spacer"></div>

      <!-- 右侧控制按钮 -->
      <button class="title-btn" @click="$emit('toggle-right-panel')">
        <v-icon size="20">mdi-folder-outline</v-icon>
      </button>

      <!-- 窗口控制按钮 -->
      <button class="title-btn window-btn" @click.stop.prevent="minimize">
        <v-icon size="18">mdi-minus</v-icon>
      </button>

      <button class="title-btn window-btn" @click.stop.prevent="maximize">
        <v-icon size="18">mdi-checkbox-blank-outline</v-icon>
      </button>

      <button class="title-btn window-btn close-btn" @click.stop.prevent="close">
        <v-icon size="18">mdi-close</v-icon>
      </button>
    </div>
  </div>
</template>

<style scoped>
.title-bar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  width: 100%;
  height: 40px;
  background-color: rgb(var(--v-theme-primary));
  color: white;
  user-select: none;
  display: flex;
  align-items: center;
  z-index: 9999;
}

.title-bar-content {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  padding: 0 4px;
}

.title-text {
  font-size: 14px;
  font-weight: 500;
  margin-left: 8px;
  white-space: nowrap;
}

.spacer {
  flex: 1;
}

.title-btn {
  -webkit-app-region: no-drag;
  width: 46px;
  height: 40px;
  border: none;
  background: transparent;
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s;
  position: relative;
  z-index: 10;
}

.title-btn:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.title-btn:active {
  background-color: rgba(255, 255, 255, 0.2);
}

.window-btn {
  border-radius: 0;
}

.close-btn:hover {
  background-color: rgba(239, 83, 80, 0.9) !important;
}
</style>
