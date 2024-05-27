<script setup lang="ts">
import {NButton, TreeInst, TreeOption, useNotification} from "naive-ui"
import {defineModel, defineProps, nextTick, ref} from "vue"
import {DownloadOutline as DownloadIcon} from "@vicons/ionicons5"
import {DownloadStatus} from "../../constants/download-constant"
import DownloadProgress from "./DownloadProgress.vue";

const notification = useNotification()

const props = defineProps<{
  downloadTreeInst: TreeInst | null,
  downloadTreeOptions: TreeOption[],
  optionsToDownload: (TreeOption | null)[],
  downloadProgressRefs: (InstanceType<typeof DownloadProgress>)[],
}>()

const searchDisabled = defineModel<boolean>("searchDisabled", {required: true})

const downloadButtonLoading = ref<boolean>(false)

// TODO: 增加取消下载的功能
async function onDownload() {
  if (props.optionsToDownload.length === 0) {
    notification.create({type: "error", title: "下载失败", content: "请选择要下载的章节", duration: 2000,})
    return
  }

  try {
    searchDisabled.value = true
    downloadButtonLoading.value = true
    await downloadOptions()
  } finally {
    searchDisabled.value = false
    downloadButtonLoading.value = false
  }
}

async function downloadOptions() {
  // 先把所有的option都disable掉，将他们加入到下载队列中
  props.optionsToDownload.forEach(option => {
    if (option !== null) {
      option.disabled = true
      option.suffix = () => DownloadStatus.WAITING
    }
  })

  await nextTick()
  let first = true
  while (props.downloadProgressRefs.length > 0) {
    const progress = props.downloadProgressRefs[0]

    if (!first) {
      await progress.wait()
    } else {
      first = false
    }

    await progress.download()
  }

}

</script>

<template>
  <n-button @click="onDownload" type="primary"
            :loading="downloadButtonLoading">
    开始下载
    <template #icon>
      <n-icon>
        <download-icon/>
      </n-icon>
    </template>
  </n-button>
</template>

