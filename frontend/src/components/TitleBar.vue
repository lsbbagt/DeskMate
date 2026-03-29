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
  <v-app-bar 
    flat 
    density="compact"
    height="40"
    class="title-bar"
    color="primary"
  >
    <!-- 左侧控制按钮 -->
    <v-btn
      icon
      size="small"
      variant="text"
      @click="$emit('toggle-sidebar')"
      class="control-btn mr-1"
    >
      <v-icon>mdi-menu</v-icon>
    </v-btn>

    <!-- 标题 -->
    <v-app-bar-title class="text-body-2 font-weight-medium title-text">
      StatBox
    </v-app-bar-title>

    <v-spacer />

    <!-- 右侧控制按钮 -->
    <v-btn
      icon
      size="small"
      variant="text"
      @click="$emit('toggle-right-panel')"
      class="control-btn mr-1"
    >
      <v-icon>mdi-folder-outline</v-icon>
    </v-btn>

    <!-- 窗口控制按钮 -->
    <v-btn
      icon
      size="small"
      variant="text"
      @click="minimize"
      class="window-btn"
    >
      <v-icon size="18">mdi-minus</v-icon>
    </v-btn>

    <v-btn
      icon
      size="small"
      variant="text"
      @click="maximize"
      class="window-btn"
    >
      <v-icon size="18">mdi-checkbox-blank-outline</v-icon>
    </v-btn>

    <v-btn
      icon
      size="small"
      variant="text"
      @click="close"
      class="window-btn close-btn"
    >
      <v-icon size="18">mdi-close</v-icon>
    </v-btn>
  </v-app-bar>
</template>

<style scoped>
.title-bar {
  -webkit-app-region: drag;
  user-select: none;
  padding: 0 !important;
}

.title-bar :deep(.v-toolbar__content) {
  padding: 0 8px;
}

.title-text {
  -webkit-app-region: drag;
  flex: none;
}

.control-btn,
.window-btn {
  -webkit-app-region: no-drag;
}

.window-btn {
  min-width: 46px !important;
  border-radius: 0 !important;
}

.close-btn:hover {
  background-color: rgba(239, 83, 80, 0.9) !important;
}
</style>
