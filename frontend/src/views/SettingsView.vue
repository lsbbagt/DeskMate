<script lang="ts" setup>
import { ref } from 'vue'
import { OpenDataFolder } from '../../wailsjs/go/main/App'

const dataFolderPath = ref('')

const openDataFolder = async () => {
  try {
    const path = await OpenDataFolder()
    dataFolderPath.value = path
  } catch (error) {
    console.error('打开文件夹失败:', error)
  }
}
</script>

<template>
  <v-container fluid class="fill-height pa-6">
    <v-card flat max-width="600" class="mx-auto">
      <v-card-title class="text-h5 mb-4">
        <v-icon start color="primary">mdi-cog</v-icon>
        设置
      </v-card-title>

      <v-divider class="mb-4" />

      <v-card-text>
        <v-card variant="outlined" class="mb-4">
          <v-card-title class="text-subtitle-1">
            <v-icon start>mdi-folder</v-icon>
            数据存储位置
          </v-card-title>
          <v-card-text class="text-body-2 text-secondary">
            所有数据（收藏夹、待办事项、代码模板、Markdown文件）都存储在这个文件夹中。
            迁移到新设备时，只需复制整个文件夹即可。
          </v-card-text>
          <v-card-actions>
            <v-btn
              color="primary"
              prepend-icon="mdi-folder-open"
              variant="flat"
              @click="openDataFolder"
            >
              打开文件夹
            </v-btn>
          </v-card-actions>
        </v-card>

        <v-alert
          type="info"
          variant="tonal"
          density="compact"
          class="mt-4"
        >
          <template #prepend>
            <v-icon>mdi-information</v-icon>
          </template>
          数据文件夹位于用户主目录下的 <code>.deskmate</code> 文件夹
        </v-alert>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<style scoped>
code {
  background-color: rgb(var(--v-theme-surface-variant));
  padding: 2px 6px;
  border-radius: 4px;
}
</style>
