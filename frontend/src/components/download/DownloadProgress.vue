<script setup lang="ts">
import {defineProps, h, ref} from "vue"
import {NSpin, TreeOption, useNotification} from "naive-ui";
import {DownloadStatus} from "../../constants/download-constant";
import {EventsOff, EventsOn} from "../../../wailsjs/runtime";
import {DownloadChapter} from "../../../wailsjs/go/api/DownloadApi";
import {useDownloaderStore} from "../../stores/downloader";

const props = defineProps<{
  treeOption: TreeOption | null,
}>()

const notification = useNotification()
const store = useDownloaderStore()

const chapterProgressIndicator = ref<string>("")
const chapterProgressPercentage = ref<number>(0)

async function download() {
  if (props.treeOption === null) {
    return
  }

  chapterProgressIndicator.value = `即将下载 ${props.treeOption.label}`

  props.treeOption.prefix = () => h(NSpin, {size: 15})
  props.treeOption.suffix = () => DownloadStatus.DOWNLOADING

  const key = JSON.parse(props.treeOption.key as string) as { href: string, saveDir: string }
  const chapterUrl = "https://www.manhuagui.com" + key.href
  const saveDirectory = key.saveDir
  const concurrentCount = store.downloadConcurrentCount
  const proxyUrl = store.proxyUrl
  EventsOn("download", (msg: string, percentage: number) => {
    chapterProgressIndicator.value = msg
    chapterProgressPercentage.value = percentage
  })
  try {
    const response = await DownloadChapter(chapterUrl, saveDirectory, concurrentCount, proxyUrl)
    // 如果下载失败
    if (response.code != 0) {
      console.error(response.msg)
      notification.create({type: "error", title: "下载失败", meta: response.msg,})
      props.treeOption.disabled = false
      props.treeOption.prefix = undefined
      props.treeOption.suffix = () => DownloadStatus.FAILED
      return
    }
    // 下载成功
    props.treeOption.prefix = undefined
    props.treeOption.suffix = () => DownloadStatus.COMPLETED
  } finally {
    EventsOff("download")
  }
}

async function wait() {
  for (let i = 0; i < store.downloadInterval; i++) {
    await new Promise(resolve => setTimeout(resolve, 1000))
    chapterProgressIndicator.value = `将在${store.downloadInterval - i}秒后继续下载(防止被ban IP)`
  }

}

defineExpose({
  download,
  wait,
})

</script>

<template>
  <div class="flex">
    <n-text>{{ treeOption?.label }}</n-text>
    <n-progress class="flex-1"
                type="line"
                :percentage="chapterProgressPercentage"
                :height="22"
                indicator-placement="inside"
                indicator-text-color="black"
    >{{ chapterProgressIndicator }}
    </n-progress>
  </div>
</template>
