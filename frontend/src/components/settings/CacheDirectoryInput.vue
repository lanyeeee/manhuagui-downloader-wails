<script setup lang="ts">
import {useDownloaderStore} from "../../stores/downloader";
import {useNotification} from "naive-ui";
import {ChooseDirectory} from "../../../wailsjs/go/api/SettingsApi";

const store = useDownloaderStore()
const notification = useNotification()

async function onChooseCacheDirectory() {
  const response = await ChooseDirectory(store.cacheDirectory)
  if (response.code != 0) {
    notification.create({type: "error", title: "选择缓存目录失败", meta: response.msg,})
  } else if (response.data !== "") {
    store.cacheDirectory = response.data as string
  }
}

</script>

<template>
  <n-popover trigger="hover">
    <template #trigger>
      <n-input class="text-align-left"
               v-model:value="store.cacheDirectory"
               readonly
               placeholder=""
               @click="onChooseCacheDirectory">
        <template #prefix>缓存目录：</template>
      </n-input>
    </template>
    <span>{{ store.cacheDirectory }}</span>
  </n-popover>
</template>
