<script setup lang="ts">
import {defineProps, ref} from "vue"
import {TreeInst, TreeOption} from "naive-ui";
import DownloadProgress from "./DownloadProgress.vue";
import DownloadButton from "./DownloadButton.vue";
import {BrowserOpenURL} from "../../../wailsjs/runtime";
import {useDownloaderStore} from "../../stores/downloader";

const store = useDownloaderStore()

const props = defineProps<{
  downloadTreeInst: TreeInst | null,
  downloadTreeOptions: TreeOption[],
  optionsToDownload: (TreeOption | null)[],
  optionsDownloading: (TreeOption | null)[],
}>()

const downloadProgresses = ref<(InstanceType<typeof DownloadProgress>)[]>([])

const searchDisabled = defineModel<boolean>("searchDisabled", {required: true})

async function openCacheDirectory() {
  BrowserOpenURL(store.cacheDirectory)
}

</script>

<template>
  <div class="flex flex-col h-full gap-y-2">
    <n-scrollbar>
      <div class="flex flex-col gap-y-2">
        <n-h3>下载队列</n-h3>
        <download-progress v-for="option in props.optionsDownloading"
                           :key="option?.key"
                           ref="downloadProgresses"
                           :tree-option="option"
        />
      </div>
    </n-scrollbar>

    <download-button
        :download-tree-inst="downloadTreeInst"
        :download-tree-options="downloadTreeOptions"
        :options-to-download="optionsToDownload"
        :download-progresses="downloadProgresses"
        v-model:search-disabled="searchDisabled"
    />
    <n-button @click="openCacheDirectory">打开下载目录</n-button>
  </div>

</template>
