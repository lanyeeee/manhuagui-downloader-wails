<script lang="ts" setup>
import {useDownloaderStore} from "./stores/downloader"
import {onMounted} from "vue"
import DownloadPane from "./components/download/DownloadPane.vue";
import ExportPane from "./components/export/ExportPane.vue";
import SettingsPane from "./components/settings/SettingsPane.vue";

const store = useDownloaderStore()
onMounted(async () => {
  await store.init()
})

</script>

<template>
  <n-notification-provider placement="bottom-right">
    <n-tabs default-value="download" type="line" animated class="h-screen p-3">
      <n-tab-pane class="h-full p-0!" name="download" tab="下载">
        <download-pane></download-pane>
      </n-tab-pane>
      <n-tab-pane class="h-full p-0!" name="export" tab="导出">
        <export-pane></export-pane>
      </n-tab-pane>
      <n-tab-pane class="h-full p-0!" name="settings" tab="设置">
        <settings-pane></settings-pane>
      </n-tab-pane>
    </n-tabs>

  </n-notification-provider>
</template>

