<script setup lang="ts">
import {NButton, TreeInst, TreeOption, useNotification} from "naive-ui";
import {defineModel, defineProps, nextTick, ref} from "vue";
import {DownloadOutline as DownloadIcon, PauseOutline as CancelIcon} from "@vicons/ionicons5";
import DownloadProgress from "./DownloadProgress.vue";
import {DownloadStatus} from "../../constants/download-constant";

enum DownloadButtonStatus {
  WAITING,
  DOWNLOADING,
  CANCELING,
}

const notification = useNotification();

const props = defineProps<{
  downloadTreeInst: TreeInst | null,
  downloadTreeOptions: TreeOption[],
  optionsToDownload: (TreeOption | null)[],
  downloadProgressRefs: (InstanceType<typeof DownloadProgress>)[],
}>();

const searchDisabled = defineModel<boolean>("searchDisabled", {required: true});
const downloadButtonStatus = ref<DownloadButtonStatus>(DownloadButtonStatus.WAITING);


async function onDownload() {
  if (props.optionsToDownload.length === 0) {
    notification.create({type: "error", title: "下载失败", content: "请选择要下载的章节", duration: 2000,});
    return;
  }

  try {
    searchDisabled.value = true;
    downloadButtonStatus.value = DownloadButtonStatus.DOWNLOADING;
    await downloadOptions();
  } finally {
    searchDisabled.value = false;
    downloadButtonStatus.value = DownloadButtonStatus.WAITING;
  }
}

async function downloadOptions() {
  // 先把所有的option都disable掉，将他们加入到下载队列中
  props.optionsToDownload.forEach(option => {
    if (option !== null) {
      option.disabled = true;
      option.suffix = () => DownloadStatus.WAITING;
    }
  });

  await nextTick();
  let first = true;
  while (props.downloadProgressRefs.length > 0 && downloadButtonStatus.value === DownloadButtonStatus.DOWNLOADING) {
    const progress = props.downloadProgressRefs[0];

    if (!first) {
      await progress.wait();
    } else {
      first = false;
    }

    await progress.download();
  }

  // 把没有下载完成的option重新enable
  props.optionsToDownload.forEach(option => {
    if (option !== null) {
      option.disabled = false;
      option.suffix = () => undefined;
    }
  });

}

</script>

<template>
  <n-button v-if="downloadButtonStatus == DownloadButtonStatus.WAITING"
            @click="onDownload"
            type="primary">
    开始下载
    <template #icon>
      <n-icon>
        <download-icon/>
      </n-icon>
    </template>
  </n-button>
  <n-button v-else-if="downloadButtonStatus == DownloadButtonStatus.DOWNLOADING"
            @click="downloadButtonStatus = DownloadButtonStatus.CANCELING"
            type="error">
    取消下载
    <template #icon>
      <n-icon>
        <cancel-icon/>
      </n-icon>
    </template>
  </n-button>
  <n-button v-else-if="downloadButtonStatus == DownloadButtonStatus.CANCELING"
            type="error"
            loading>
    将在本章节下载完成后取消
  </n-button>


</template>

