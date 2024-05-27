<script setup lang="ts">
import {useDownloaderStore} from "../../stores/downloader";
import {ChooseDirectory} from "../../../wailsjs/go/api/SettingsApi";
import {useNotification} from "naive-ui";

const store = useDownloaderStore()
const notification = useNotification()

async function onChooseExportDirectory() {
  const response = await ChooseDirectory(store.exportDirectory)
  if (response.code != 0) {
    notification.create({type: "error", title: "选择导出目录失败", meta: response.msg,})
  } else if (response.data !== "") {
    store.exportDirectory = response.data as string
  }
}
</script>

<template>
  <n-popover trigger="hover">
    <template #trigger>
      <n-input class="text-align-left"
               v-model:value="store.exportDirectory"
               placeholder=""
               readonly
               @click="onChooseExportDirectory">
        <template #prefix>导出目录：</template>
      </n-input>
    </template>
    <span>{{ store.exportDirectory }}</span>
  </n-popover>
</template>
